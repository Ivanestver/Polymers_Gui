package datatypes

import (
	"polymers/base"
	"polymers/globaldata"
)

type RealField struct {
	monomers map[base.Vector3DF]*Monomer
}

func NewRealField(restrictions [3][2]float64) *RealField {
	realField := &RealField{}
	globaldata.ConfigureGlobalData(globaldata.SpaceDimention{
		{Lower: restrictions[0][0], Higher: restrictions[0][1]},
		{Lower: restrictions[1][0], Higher: restrictions[1][1]},
		{Lower: restrictions[2][0], Higher: restrictions[2][1]},
	})
	realField.monomers = make(map[base.Vector3DF]*Monomer)
	return realField
}

func (realField *RealField) MakeFilled(monomer *Monomer) {
	if globaldata.GetGlobalData().SpaceDimention.PointInSpace(&monomer.coords) {
		if _, ok := realField.monomers[monomer.coords]; !ok {
			realField.monomers[monomer.coords] = monomer
		}
	}
}

func (realField *RealField) MakeFree(monomer *Monomer) {
	if !globaldata.GetGlobalData().SpaceDimention.PointInSpace(&monomer.coords) {
		return
	}
	realField.monomers[monomer.coords] = nil
}

func (realField *RealField) IsFree(coords base.Vector3DF) bool {
	if !globaldata.GetGlobalData().SpaceDimention.PointInSpace(&coords) {
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
	if !globaldata.GetGlobalData().SpaceDimention.PointInSpace(&coords) {
		return nil
	}
	if m, ok := realField.monomers[coords]; ok {
		return m
	} else {
		m1 := NewMonomer(coords, MonomerTypeUsual)
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

func (realField *RealField) GetMinMonomersByAxis(axis base.Axis) []*Monomer {
	minMonomers := make([]*Monomer, 0)
	for point, mon := range realField.monomers {
		if len(minMonomers) == 0 || point[axis] < minMonomers[0].Coords()[axis] {
			minMonomers = []*Monomer{mon}
		} else if base.CompareFloat(point[axis], minMonomers[0].Coords()[axis]) {
			minMonomers = append(minMonomers, mon)
		} else {
			continue
		}
	}
	return minMonomers
}

func (realField *RealField) GetMaxMonomersByAxis(axis base.Axis) []*Monomer {
	minMonomers := make([]*Monomer, 0)
	for point, mon := range realField.monomers {
		if len(minMonomers) == 0 || point[axis] > minMonomers[0].Coords()[axis] {
			minMonomers = []*Monomer{mon}
		} else if base.CompareFloat(point[axis], minMonomers[0].Coords()[axis]) {
			minMonomers = append(minMonomers, mon)
		} else {
			continue
		}
	}
	return minMonomers
}
