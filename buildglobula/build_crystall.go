package buildglobula

import (
	"polymers/datatypes"
	"polymers/views"
)

type CrystallDataBuilder struct {
}

func (builder CrystallDataBuilder) CreateInputData(algType AlgType, defaultParams []string, particleName string) (ICalcAlgInputData, error) {
	return CrystallBuildInputData{}, nil
}

type CrystallBuildInputData struct {
}

func (inputData CrystallBuildInputData) GetGlobulaType() views.GlobulaProperty {
	return views.GlobulaCrystalType
}

type BuildCrystallAlg AbstractAlg[CrystallBuildInputData]

func (alg *BuildCrystallAlg) Calc() []*datatypes.Polymer {
	return nil
}
