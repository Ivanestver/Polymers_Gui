package buildglobula

import (
	"bufio"
	"errors"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/globaldata"
	"polymers/outputformat"
	"polymers/views"
)

type PatternInputDataBuilder struct {
}

func (builder *PatternInputDataBuilder) CreateInputData(algType AlgType, defaultParams []string, particleName string) (ICalcAlgInputData, error) {
	if len(defaultParams) == 0 {
		return nil, errors.New("не задан файл с паттерном")
	}
	return &PatternAlgInputData{
		patternFileName: defaultParams[0],
	}, nil
}

type PatternAlgInputData struct {
	patternFileName string
}

func (inputData *PatternAlgInputData) GetGlobulaType() views.GlobulaProperty {
	return views.GlobulaPatternType
}

type PatternCalcAlg struct {
	inputData *PatternAlgInputData
}

/*
Y
^
|
|*-* *-* *-*
|| | | | | |
|* * * * * * *
|| | | | | | |
|* *-* *-* *-*
 ----------------------------------> X
*/

func (alg *PatternCalcAlg) Calc() []*datatypes.Polymer {
	file, err := os.Open(alg.inputData.patternFileName)
	if err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return nil
	}
	polymers := make([]*datatypes.Polymer, 1)
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	field := datatypes.NewField(uint64(
		max(
			spaceDimention[base.AxisX].Higher,
			spaceDimention[base.AxisX].Lower,
			spaceDimention[base.AxisY].Higher,
			spaceDimention[base.AxisY].Lower,
			spaceDimention[base.AxisZ].Higher,
			spaceDimention[base.AxisZ].Lower)))
	polymers[0] = datatypes.NewPolymer(field, 0)
	direction := base.AxisYVec
	currPoint := base.AxisYVecReversed
	scanner := bufio.NewScanner(file)
	for scanner.Scan() && scanner.Err() == nil {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		for _, beadType := range line {
			currPoint = base.AddVecF(currPoint, direction)
			mon := field.GetMonomerByCoords(currPoint)
			mon.MonomerType = base.MendeleevTableElement(string(beadType))
			polymers[0].AddMonomer(mon)
		}
		currPoint = base.AddVecF(base.AddVecF(currPoint, direction), base.AxisXVec)
		direction.MultiplyByConstantF(-1.0)
	}
	return polymers
}
