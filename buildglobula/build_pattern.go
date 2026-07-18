package buildglobula

import (
	"errors"
	"polymers/datatypes"
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

func (alg *PatternCalcAlg) Calc() []*datatypes.Polymer {
	return nil
}
