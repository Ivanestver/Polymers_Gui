package buildglobula

import (
	"polymers/datatypes"
	"polymers/views"
)

type CrystallInputDataBuilder struct {
}

func (builder CrystallInputDataBuilder) CreateInputData(defaultParams []string, particleName string) (ICalcAlgInputData, error) {
	return CrystallBuildInputData{}, nil
}

type CrystallBuildInputData struct {
}

func (inputData CrystallBuildInputData) GetGlobulaType() views.GlobulaProperty {
	return views.GlobulaCrystalType
}

type CrystallCalcAlg AbstractAlg[CrystallBuildInputData]

func (alg *CrystallCalcAlg) Calc() []*datatypes.Polymer {
	return nil
}
