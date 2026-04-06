package datatypes

import (
	"polymers/base"
	"polymers/global_data"
)

type RealField struct {
	monomers map[base.Vector3DF]*Monomer
}

func NewRealField(restrictions [3][2]float64) *RealField {
	realField := &RealField{}
	global_data.ConfigureGlobalData(global_data.SpaceDimention{
		{Lower: restrictions[0][0], Higher: restrictions[0][1]},
		{Lower: restrictions[1][0], Higher: restrictions[1][1]},
		{Lower: restrictions[2][0], Higher: restrictions[2][1]},
	})
	realField.monomers = make(map[base.Vector3DF]*Monomer)
	return realField
}

func (realField *RealField) MakeFilled(monomer *Monomer) {
	if global_data.GetGlobalData().SpaceDimention.PointInSpace(&monomer.coords) {
		if _, ok := realField.monomers[monomer.coords]; !ok {
			realField.monomers[monomer.coords] = monomer
		}
	}
}

func (realField *RealField) MakeFree(monomer *Monomer) {
	if !global_data.GetGlobalData().SpaceDimention.PointInSpace(&monomer.coords) {
		return
	}
	realField.monomers[monomer.coords] = nil
}

func (realField *RealField) IsFree(coords base.Vector3DF) bool {
	if !global_data.GetGlobalData().SpaceDimention.PointInSpace(&coords) {
		return false
	}
	if m, ok := realField.monomers[coords]; !ok {
		realField.monomers[coords] = nil
		return true
	} else {
		return m == nil
	}
}

func (realField *RealField) GetMonomerByCoords(coords base.Vector3DF) *Monomer {
	if !global_data.GetGlobalData().SpaceDimention.PointInSpace(&coords) {
		return nil
	}
	if m, ok := realField.monomers[coords]; ok {
		return m
	} else {
		m1 := NewMonomer(coords, MONOMER_TYPE_USUAL)
		realField.monomers[coords] = m1
		return m1
	}
}

func (realField *RealField) GetMonomersWithin(lower, higher base.Vector3DF) []*Monomer {
	isIn := func(point base.Vector3DF) bool {
		return base.PointInSpace(&point, &lower, &higher)
	}
	monomers := make([]*Monomer, 0)
	for point, mon := range realField.monomers {
		if isIn(point) {
			monomers = append(monomers, mon)
		}
	}
	return monomers
}
