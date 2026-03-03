package datatypes

import "polymers/base"

type RealField struct {
	monomers map[base.Vector3DF]*Monomer
}

func NewRealField() *RealField {
	realField := &RealField{}
	realField.monomers = make(map[base.Vector3DF]*Monomer)
	return realField
}

func (realField *RealField) MakeFilled(monomer *Monomer) {
	if _, ok := realField.monomers[monomer.coords]; !ok {
		realField.monomers[monomer.coords] = monomer
	}
}

func (realField *RealField) MakeFree(monomer *Monomer) {
	realField.monomers[monomer.coords] = nil
}

func (realField *RealField) IsFree(coords base.Vector3DF) bool {
	if m, ok := realField.monomers[coords]; !ok {
		realField.monomers[coords] = nil
		return true
	} else {
		return m == nil
	}
}

func (realField *RealField) GetMonomerByCoords(coords base.Vector3DF) *Monomer {
	if m, ok := realField.monomers[coords]; ok {
		return m
	} else {
		m1 := NewMonomer(coords, MONOMER_TYPE_USUAL)
		realField.monomers[coords] = m1
		return m1
	}
}
