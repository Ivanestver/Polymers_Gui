package buildglobula

import (
	"errors"
	"math/rand/v2"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/outputformat"
	"polymers/views"
	"strconv"
)

type CrystallInputDataBuilder struct {
}

func (builder CrystallInputDataBuilder) CreateInputData(defaultParams []string) (ICalcAlgInputData, error) {
	if len(defaultParams) == 0 {
		return nil, errors.New("Не заданы параметры")
	}
	maxAmorphousPartSize, err := strconv.Atoi(defaultParams[0])
	if err != nil {
		return nil, err
	}
	return CrystallBuildInputData{
		maxAmorphousPartSize: maxAmorphousPartSize,
	}, nil
}

type CrystallBuildInputData struct {
	maxAmorphousPartSize int
}

func (inputData CrystallBuildInputData) GetGlobulaType() views.GlobulaProperty {
	return views.GlobulaCrystalType
}

type CrystallCalcAlg AbstractAlg[CrystallBuildInputData]

func (alg *CrystallCalcAlg) Calc() []*datatypes.Polymer {
	filename := "crystallalg.txt"
	polymers, err := alg.getPolymers(filename)
	if err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return nil
	}
	if err = alg.processPolymers(&polymers); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return nil
	}
	if err := deleteFile(filename); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return nil
	}
	return polymers
}

func (alg *CrystallCalcAlg) getPolymers(filename string) ([]*datatypes.Polymer, error) {
	algbuilder := CreateInputDataBuilder(PatternBuildAlg)
	inputData, err := algbuilder.CreateInputData([]string{filename, "false", "4"})
	if err != nil {
		return nil, err
	}
	patternAlg := CreateCalcAlg(inputData, PatternBuildAlg)
	if patternAlg == nil {
		return nil, err
	}
	if err := createPatternFile(filename); err != nil {
		return nil, err
	}
	return patternAlg.Calc(), nil
}

func (alg *CrystallCalcAlg) processPolymers(polymers *[]*datatypes.Polymer) error {
	if err := alg.breakToPieces(polymers); err != nil {
		return err
	}
	if err := alg.vaccinateRandomAmorphousParts(*polymers); err != nil {
		return err
	}
	return nil
}

func (alg *CrystallCalcAlg) breakToPieces(polymers *[]*datatypes.Polymer) error {
	newPolymers := make([]*datatypes.Polymer, 0)
	for _, polymer := range *polymers {
		if pieces, err := alg.breakConnectionAndGetPieces(polymer); err == nil {
			newPolymers = append(newPolymers, pieces...)
		} else {
			return err
		}
	}
	*polymers = newPolymers
	return nil
}

func (alg *CrystallCalcAlg) breakConnectionAndGetPieces(polymer *datatypes.Polymer) ([]*datatypes.Polymer, error) {
	if polymer.Len() <= 2 {
		return nil, errors.New("Слишком короткий полимер")
	}
	newPolymers := make([]*datatypes.Polymer, 0)
	direction := base.SubtractVecF(polymer.GetMonomerByIdx(1).Coords(), polymer.GetMonomerByIdx(0).Coords())
	field := polymer.GetField().(*datatypes.Field)
	newPolymer := datatypes.NewPolymer(field, int64(len(newPolymers)))
	newPolymer.AddMonomerNoConnection(polymer.GetMonomerByIdx(0))
	for i := 2; i < polymer.Len(); i++ {
		prevMonomer := polymer.GetMonomerByIdx(i - 1)
		currMonomer := polymer.GetMonomerByIdx(i)
		currDirection := base.SubtractVecF(currMonomer.Coords(), prevMonomer.Coords())
		newPolymer.AddMonomerNoConnection(prevMonomer)
		if !base.VectorsAreEqualF(direction, currDirection) {
			direction = base.RevertVecF(direction)
			if rand.IntN(100)%2 == 0 {
				newPolymers = append(newPolymers, newPolymer)
				newPolymer = datatypes.NewPolymer(field, int64(len(newPolymers)))
				if err := datatypes.BreakConnection(prevMonomer, currMonomer); err != nil {
					return nil, err
				}
			}
		}
	}
	if newPolymer.Len() > 0 {
		newPolymer.AddMonomerNoConnection(polymer.LastMonomer())
		newPolymers = append(newPolymers, newPolymer)
	}
	return newPolymers, nil
}

func (alg *CrystallCalcAlg) vaccinateRandomAmorphousParts(polymers []*datatypes.Polymer) error {
	for _, polymer := range polymers {
		if err := alg.vaccinateRandomAmorphousPartsToPolymer(polymer); err != nil {
			return err
		}
	}
	return nil
}

func (alg *CrystallCalcAlg) vaccinateRandomAmorphousPartsToPolymer(polymer *datatypes.Polymer) error {
	field := polymer.GetField().(*datatypes.Field)
	if err := alg.growAmorphousPartFromMonomer(
		field,
		func() *datatypes.Monomer {
			return polymer.LastMonomer()
		},
		func(newMonomer *datatypes.Monomer) {
			polymer.AddMonomer(newMonomer)
		}); err != nil {
		return err
	}
	if err := alg.growAmorphousPartFromMonomer(
		field,
		func() *datatypes.Monomer {
			return polymer.GetMonomerByIdx(0)
		},
		func(newMonomer *datatypes.Monomer) {
			polymer.AddMonomerAtStart(newMonomer)
		}); err != nil {
		return err
	}
	return nil
}

func (alg *CrystallCalcAlg) growAmorphousPartFromMonomer(field *datatypes.Field, getCurrMonomer func() *datatypes.Monomer, addMonomer func(*datatypes.Monomer)) error {
	for i := 0; i < alg.inputData.maxAmorphousPartSize; i++ {
		monomer := getCurrMonomer()
		availableCells := field.GetAvailableCells(monomer.Coords())
		if len(availableCells) == 0 {
			return nil
		}
		randomInt := rand.IntN(len(availableCells))
		newCell := availableCells[randomInt]
		newCell.MonomerType = base.O
		addMonomer(newCell)
	}
	return nil
}

func createPatternFile(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	buf := []byte("CCCCC\nCCCCC\nCCCCC\nCCCCC\nCCCCC\nCCCCC\nCCCCC")

	return os.WriteFile(filename, buf, 0666)
}

func deleteFile(filename string) error {
	err := os.Remove(filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
