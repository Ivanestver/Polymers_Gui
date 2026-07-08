package trajectories

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"polymers/base"
	"polymers/cristallinity"
	"polymers/datatypes"
	"polymers/savers"
	"polymers/views"
	"strconv"
	"strings"
)

type _TrajectoriesApplier struct {
	globula               *views.GlobulaView
	trajectoriesFileNames []string
	fileReader            *bufio.Scanner
	monomers              map[int64]*datatypes.Monomer
	currentTimestep       int
	outputFilename        string
	percent               float64
}

func newTrajectoriesApplier(globula *views.GlobulaView, trajectoriesFilenames []string, outputFilename string, percent int) (*_TrajectoriesApplier, error) {
	applier := &_TrajectoriesApplier{
		globula:               globula,
		trajectoriesFileNames: trajectoriesFilenames,
	}
	applier.monomers = make(map[int64]*datatypes.Monomer)
	for polNumber := 0; polNumber < globula.Len(); polNumber++ {
		polymer := globula.GetPolymerByIdx(polNumber)
		for monNumber := 0; monNumber < polymer.Len(); monNumber++ {
			monomer := polymer.GetMonomerByIdx(monNumber)
			applier.monomers[monomer.Number] = monomer
		}
	}
	applier.currentTimestep = 0
	applier.outputFilename = outputFilename
	applier.percent = float64(percent) / 100.0
	return applier, nil
}

func (applier *_TrajectoriesApplier) apply() error {
	fileForS, err := os.Create(applier.outputFilename + "/s.csv")
	if err != nil {
		return err
	}
	defer fileForS.Close()
	fileForS.WriteString("Timestep;Cos (stick);Tensor (stick);Cos (C-to-C);Tensor (C-to-C);CristallinityRate\n")
	for _, filename := range applier.trajectoriesFileNames {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		applier.fileReader = bufio.NewScanner(file)
		numberOfAtoms, err := applier.getNumberOfAtoms()
		if err != nil {
			return err
		}
		for {
			timestep := applier.getNextTimestep()
			fmt.Printf("Timestep: %d\n", timestep)
			// Apply trajectories
			if err = applier.makeStep(numberOfAtoms); err != nil {
				break
			}
			if err := applier.calculateS(timestep, fileForS); err != nil {
				return err
			}
			// Save results
			if err := applier.saveTimestep(timestep); err != nil {
				return err
			}
		}
	}
	return nil
}

func (applier *_TrajectoriesApplier) getNextTimestep() int {
	applier.currentTimestep++
	return applier.currentTimestep
}

func (applier *_TrajectoriesApplier) getNumberOfAtoms() (int, error) {
	for applier.fileReader.Scan() {
		if strings.Contains(applier.fileReader.Text(), "TotAtms") {
			line := applier.fileReader.Text()
			line = strings.Trim(line, " ")
			line = base.ReduceSpacesToSingle(line)
			parts := strings.Split(line, " ")
			numberOfAtoms, err := strconv.Atoi(parts[4])
			if err == nil {
				return numberOfAtoms, nil
			} else {
				return -1, err
			}
		}
	}
	return -1, errors.New("больше нет секций NUMBER OF ATOMS")
}

func (applier *_TrajectoriesApplier) getBoxBounds() (box [3][2]float64, err error) {
	err = nil
	for applier.fileReader.Scan() && !strings.Contains(applier.fileReader.Text(), "BOX BOUNDS") {
	}
	for i := 0; i < 3; i++ {
		if applier.fileReader.Scan() {
			line := applier.fileReader.Text()
			prevLen := len(line)
			for {
				line = strings.ReplaceAll(line, "  ", " ")
				if prevLen == len(line) {
					break
				}
				prevLen = len(line)
			}
			line = strings.Trim(line, " ")
			parts := strings.Split(line, " ")
			if len(parts) != 2 {
				err = errors.New("в координатах должно быть 2 значения")
				return
			}
			for j := 0; j < len(parts); j++ {
				value, e := strconv.ParseFloat(strings.Trim(parts[j], " "), 64)
				if e != nil {
					err = e
				}
				box[i][j] = value
			}
		}
	}
	return
}

func (applier *_TrajectoriesApplier) calculateNewMonomerCoords(currCoords base.Vector3DF, trjItem *_TrjItem, box _BoxBounds) base.Vector3DF {
	offsets := make([]float64, len(box))
	for i := range offsets {
		offsets[i] = box[i][1] - box[i][0]
	}
	newCoords := currCoords
	for i := range newCoords {
		newCoords[i] = offsets[i]*trjItem.S[i] + box[i][0]
		newCoords[i] = newCoords[i] + offsets[i]*float64(trjItem.K[i])
	}
	return newCoords
}

func (applier *_TrajectoriesApplier) makeStep(numberOfAtoms int) error {
	for applier.fileReader.Scan() {
		if strings.Contains(applier.fileReader.Text(), "Coordinates") {
			for range numberOfAtoms {
				line := base.ReduceSpacesToSingle(applier.fileReader.Text())
				line = strings.Trim(line, " ")
				if err := applier.setNewCoords(line); err != nil {
					return err
				}
				applier.fileReader.Scan()
			}
			return nil
		}
	}
	return errors.New("не найдено координат")
}

func (applier *_TrajectoriesApplier) setNewCoords(line string) error {
	parts := strings.Split(line, " ")
	numberStr := parts[len(parts)-4]
	number, err := strconv.ParseInt(numberStr, 10, 64)
	if err != nil {
		return err
	}
	xStr := parts[len(parts)-3]
	x, err := strconv.ParseFloat(xStr, 64)
	if err != nil {
		return err
	}
	yStr := parts[len(parts)-2]
	y, err := strconv.ParseFloat(yStr, 64)
	if err != nil {
		return err
	}
	zStr := parts[len(parts)-1]
	z, err := strconv.ParseFloat(zStr, 64)
	if err != nil {
		return err
	}

	if mon, ok := applier.monomers[number]; ok {
		field := applier.globula.GetPolymerByIdx(0).GetUnderlinedField()
		return field.MoveMonomer(mon, base.Vector3DF{x, y, z})
	} else {
		return fmt.Errorf("нет мономера с номером %d", number)
	}
}

func (applier *_TrajectoriesApplier) calculateS(timestep int, file *os.File) error {
	debugFilename := applier.outputFilename + "/timesteps/debug_timestamp_" + strconv.Itoa(timestep) + ".dump"
	SMeanCos, STensor, SMeanCosCarbonToCarbon, STensorCarbonToCarbon, cristallinityRate, err := cristallinity.AnalyzeToOutside(applier.globula, 3, cristallinity.Atomistic, base.C, applier.percent, "s_process.log", &debugFilename)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(file, "%d;%s;%s;%s;%s;%s\n",
		timestep,
		strconv.FormatFloat(SMeanCos, 'f', 3, 64),
		strconv.FormatFloat(STensor, 'f', 3, 64),
		strconv.FormatFloat(SMeanCosCarbonToCarbon, 'f', 3, 64),
		strconv.FormatFloat(STensorCarbonToCarbon, 'f', 3, 64),
		strconv.FormatFloat(cristallinityRate, 'f', 3, 64),
	); err != nil {
		return err
	}
	return nil
}

func (applier *_TrajectoriesApplier) saveTimestep(timestep int) error {
	if s, err := savers.SaveToLammps(applier.globula); err != nil {
		return err
	} else {
		if err := os.Mkdir(applier.outputFilename+"/timesteps", os.ModeDir); err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
		filename := applier.outputFilename + "/timesteps/timestamp_" + strconv.Itoa(timestep) + ".dump"
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(s); err != nil {
			return err
		}
	}
	return nil
}

func ApplyTrajectories(globula *views.GlobulaView, trajectoriesFilenames []string, outputFilename string, percent int) error {
	applier, err := newTrajectoriesApplier(globula, trajectoriesFilenames, outputFilename, percent)
	if err != nil {
		return err
	}
	err = applier.apply()
	if err != nil {
		return err
	}
	return nil
}
