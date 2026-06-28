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

type _TrajectoriesApplierOld struct {
	globula              *views.GlobulaView
	trajectoriesFileName string
	fileReader           *bufio.Scanner
	monomers             map[int64]*datatypes.Monomer
}

func newTrajectoriesApplierOld(globula *views.GlobulaView, trajectoriesFilename string) (*_TrajectoriesApplierOld, error) {
	applier := &_TrajectoriesApplierOld{
		globula:              globula,
		trajectoriesFileName: trajectoriesFilename,
	}
	file, err := os.Open(trajectoriesFilename)
	if err != nil {
		return nil, err
	}
	applier.fileReader = bufio.NewScanner(file)
	applier.monomers = make(map[int64]*datatypes.Monomer)
	for polNumber := 0; polNumber < globula.Len(); polNumber++ {
		polymer := globula.GetPolymerByIdx(polNumber)
		for monNumber := 0; monNumber < polymer.Len(); monNumber++ {
			monomer := polymer.GetMonomerByIdx(monNumber)
			applier.monomers[monomer.Number] = monomer
		}
	}
	return applier, nil
}

func (applier *_TrajectoriesApplierOld) apply() error {
	fileForS, err := os.Create(applier.trajectoriesFileName + "_s.log")
	if err != nil {
		return err
	}
	defer fileForS.Close()
	for {
		timestep, err := applier.getNextTimestep()
		if err != nil {
			return err
		} else if timestep == -1 {
			break
		}
		fmt.Printf("Timestep: %d\n", timestep)
		numberOfAtoms, err := applier.getNumberOfAtoms()
		if err != nil {
			return err
		}
		box, err := applier.getBoxBounds()
		if err != nil {
			return err
		}
		// Apply trajectories
		if err = applier.makeStep(numberOfAtoms, box); err != nil {
			return err
		}
		if err := applier.calculateS(fileForS); err != nil {
			return err
		}
		// Save results
		if err := applier.saveTimestep(timestep); err != nil {
			return err
		}
	}
	return nil
}

func (applier *_TrajectoriesApplierOld) getNextTimestep() (int, error) {
	for applier.fileReader.Scan() && !strings.Contains(applier.fileReader.Text(), "TIMESTEP") {
	}
	if applier.fileReader.Scan() {
		line := applier.fileReader.Text()
		line = strings.Trim(line, " ")
		timestep, err := strconv.Atoi(line)
		if err == nil {
			return timestep, nil
		} else {
			return -1, err
		}
	}
	return -1, nil
}

func (applier *_TrajectoriesApplierOld) getNumberOfAtoms() (int, error) {
	for applier.fileReader.Scan() && !strings.Contains(applier.fileReader.Text(), "NUMBER OF ATOMS") {
	}
	if applier.fileReader.Scan() {
		line := applier.fileReader.Text()
		line = strings.Trim(line, " ")
		numberOfAtoms, err := strconv.Atoi(line)
		if err == nil {
			return numberOfAtoms, nil
		} else {
			return -1, err
		}
	}
	return -1, errors.New("больше нет секций NUMBER OF ATOMS")
}

func (applier *_TrajectoriesApplierOld) getBoxBounds() (box [3][2]float64, err error) {
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

func (applier *_TrajectoriesApplierOld) calculateNewMonomerCoords(currCoords base.Vector3DF, trjItem *_TrjItem, box _BoxBounds) base.Vector3DF {
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

func (applier *_TrajectoriesApplierOld) makeStep(numberOfAtoms int, box _BoxBounds) error {
	field := applier.globula.GetPolymerByIdx(0).GetUnderlinedField()
	for applier.fileReader.Scan() && !strings.Contains(applier.fileReader.Text(), "id type") {
	}
	for range numberOfAtoms {
		if !applier.fileReader.Scan() {
			return errors.New("нет координат")
		}
		s := applier.fileReader.Text()
		trjItem := trjItemFromString(strings.Trim(s, " "))
		if mon, ok := applier.monomers[int64(trjItem.AtomID)]; ok {
			newMonomerCoords := applier.calculateNewMonomerCoords(mon.Coords(), trjItem, box)
			if err := field.MoveMonomer(mon, newMonomerCoords); err != nil {
				return err
			}
		}
	}
	return nil
}

func (applier *_TrajectoriesApplierOld) calculateS(file *os.File) error {
	_, S, _, _, err := cristallinity.AnalyzeToOutside(applier.globula, 3, cristallinity.Atomistic, base.C, 15.0, "")
	if err != nil {
		return err
	}
	if _, err = file.WriteString(strconv.FormatFloat(S, 'f', 3, 64) + "\n"); err != nil {
		return err
	}
	return nil
}

func (applier *_TrajectoriesApplierOld) saveTimestep(timestep int) error {
	if s, err := savers.SaveToLammps(applier.globula); err != nil {
		return err
	} else {
		filename := "timestep_" + applier.trajectoriesFileName + "_" + strconv.Itoa(timestep) + ".dump"
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

func ApplyTrajectoriesOld(globula *views.GlobulaView, trajectoryID int, trajectoriesFilename string) error {
	applier, err := newTrajectoriesApplierOld(globula, trajectoriesFilename)
	if err != nil {
		return err
	}
	err = applier.apply()
	if err != nil {
		return err
	}
	return nil
}
