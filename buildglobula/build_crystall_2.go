package buildglobula

import (
	"polymers/datatypes"
	"polymers/views"
)

type Crystall2InputDataBuilder struct {
}

func CreateInputData(defaultParams []string) (ICalcAlgInputData, error) {
	panic("Not implemented")
}

type Crystall2InputData struct {
}

func (inputData Crystall2InputData) GetGlobulaType() views.GlobulaProperty {
	panic("Not implemented")
}

type Crystall2CalcAlg AbstractAlg[Crystall2InputData]

func (alg Crystall2CalcAlg) Calc() []*datatypes.Polymer {
	panic("Not implemented")
}
