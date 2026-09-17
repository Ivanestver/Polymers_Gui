package buildglobula

import (
	"polymers/datatypes"
	"polymers/views"
)

type Crystall2InputDataBuilder struct {
}

func (builder Crystall2InputDataBuilder) CreateInputData(defaultParams []string) (ICalcAlgInputData, error) {
	return Crystall2InputData{}, nil
}

type Crystall2InputData struct {
}

func (inputData Crystall2InputData) GetGlobulaType() views.GlobulaProperty {
	return views.GlobulaCrystal2Type
}

type Crystall2CalcAlg AbstractAlg[Crystall2InputData]

func (alg Crystall2CalcAlg) Calc() []*datatypes.Polymer {
	panic("Not implemented")
}
