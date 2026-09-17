package buildglobula

import (
	"polymers/base"
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

type Plate []base.Vector3DF

func (alg Crystall2CalcAlg) Calc() []*datatypes.Polymer {
	plates := alg.createPlates()
	polymer := alg.turnPlatesIntoPolymer(plates)
	return []*datatypes.Polymer{polymer}
}

func (alg *Crystall2CalcAlg) createPlates() []Plate {
	panic("Not implemented")
}

func (alg *Crystall2CalcAlg) turnPlatesIntoPolymer(plates []Plate) *datatypes.Polymer {
	panic("Not implemented")
}
