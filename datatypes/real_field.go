package datatypes

import (
	"fmt"
	"math"
	"polymers/base"
	"polymers/globaldata"
	"slices"
)

type _CmpResult int

const (
	_Less _CmpResult = iota
	_Exact
	_Greater
)

type _Node struct {
	Coordinate float64
	Subnodes   []_Node
	Monomer    *Monomer
}

type _RealFieldStorage struct {
	Nodes []_Node
}

func (realField *_RealFieldStorage) GetMonomer(coords base.Vector3DF) *Monomer {
	if leftX, _, isExact := realField.find(coords[base.AxisX], realField.Nodes); isExact {
		subnodes := realField.Nodes[leftX].Subnodes
		if leftY, _, isExact := realField.find(coords[base.AxisY], subnodes); isExact {
			subnodes := subnodes[leftY].Subnodes
			if leftZ, _, isExact := realField.find(coords[base.AxisZ], subnodes); isExact {
				return subnodes[leftZ].Monomer
			} else {
				newMonomer := NewMonomer(coords, base.C)
				realField.Nodes[leftX].Subnodes[leftY].Subnodes = append(realField.Nodes[leftX].Subnodes[leftY].Subnodes, _Node{
					Coordinate: coords[base.AxisZ],
					Monomer:    newMonomer,
				})
				slices.SortFunc(realField.Nodes[leftX].Subnodes[leftY].Subnodes, func(left, right _Node) int {
					if left.Coordinate < right.Coordinate {
						return -1
					} else if left.Coordinate > right.Coordinate {
						return 1
					} else {
						return 0
					}
				})
				return newMonomer
			}
		} else {
			newMonomer := NewMonomer(coords, base.C)
			realField.Nodes[leftY].Subnodes = append(realField.Nodes[leftY].Subnodes, _Node{
				Coordinate: coords[base.AxisY],
				Subnodes: []_Node{
					{
						Coordinate: coords[base.AxisZ],
						Monomer:    newMonomer,
					},
				},
			})
			slices.SortFunc(realField.Nodes[leftY].Subnodes, func(left, right _Node) int {
				if left.Coordinate < right.Coordinate {
					return -1
				} else if left.Coordinate > right.Coordinate {
					return 1
				} else {
					return 0
				}
			})
			return newMonomer
		}
	} else {
		newMonomer := NewMonomer(coords, base.C)
		realField.Nodes = append(realField.Nodes, _Node{
			Coordinate: coords[base.AxisX],
			Subnodes: []_Node{
				{
					Coordinate: coords[base.AxisY],
					Subnodes: []_Node{
						{
							Coordinate: coords[base.AxisZ],
							Monomer:    newMonomer,
						},
					},
				},
			},
		})
		slices.SortFunc(realField.Nodes, func(left, right _Node) int {
			if left.Coordinate < right.Coordinate {
				return -1
			} else if left.Coordinate > right.Coordinate {
				return 1
			} else {
				return 0
			}
		})
		return newMonomer
	}
}

func (realField *_RealFieldStorage) find(value float64, nodes []_Node) (left, right int, isExact bool) {
	left = 0
	right = len(nodes) - 1
	isExact = false
	if len(nodes) == 1 {
		isExact = base.CompareFloatWithE(value, nodes[left].Coordinate, 0.001)
		return
	}
	for left < right {
		m := (left + right) / 2
		middle := nodes[m].Coordinate
		if base.CompareFloatWithE(value, middle, 0.001) {
			left = m
			right = m
			isExact = true
			return
		} else if value < middle {
			right = m
		} else {
			left = m
		}
	}
	return
}

type RealField struct {
	monomers  map[base.Vector3DF]*Monomer
	monomers2 *_RealFieldStorage
}

func NewRealField(restrictions [3][2]float64) *RealField {
	realField := &RealField{}
	globaldata.SetSpaceDimention(globaldata.SpaceDimention{
		{Lower: restrictions[0][0], Higher: restrictions[0][1]},
		{Lower: restrictions[1][0], Higher: restrictions[1][1]},
		{Lower: restrictions[2][0], Higher: restrictions[2][1]},
	})
	realField.monomers = make(map[base.Vector3DF]*Monomer)
	realField.monomers2 = &_RealFieldStorage{
		Nodes: make([]_Node, 0),
	}
	return realField
}

func (realField *RealField) GetType() FieldType {
	return FieldTypeReal
}

func (realField *RealField) MakeFilled(monomer *Monomer) {
	if globaldata.GetGlobalData().SpaceDimention.PointInSpace(&monomer.coords) {
		roundedCoords := realField.roundCoords(monomer.coords)
		if _, ok := realField.monomers[roundedCoords]; !ok {
			realField.monomers[roundedCoords] = monomer
		}
	}
}

func (realField *RealField) MakeFree(monomer *Monomer) {
	if !globaldata.GetGlobalData().SpaceDimention.PointInSpace(&monomer.coords) {
		return
	}
	realField.monomers[realField.roundCoords(monomer.coords)] = nil
}

func (realField *RealField) IsFree(coords base.Vector3DF) bool {
	if !globaldata.GetGlobalData().SpaceDimention.PointInSpace(&coords) {
		return false
	}
	roundedCoords := realField.roundCoords(coords)
	if m, ok := realField.monomers[roundedCoords]; !ok {
		realField.monomers[roundedCoords] = nil
		return true
	} else {
		return m == nil
	}
}

func (realField *RealField) GetMonomerByCoords(coords base.Vector3DF) *Monomer {
	if !globaldata.GetGlobalData().SpaceDimention.PointInSpace(&coords) {
		return nil
	}
	// roundedCoords := realField.roundCoords(coords)
	// if m, ok := realField.monomers[roundedCoords]; ok {
	// 	return m
	// } else {
	// 	m1 := NewMonomer(roundedCoords, base.C)
	// 	realField.monomers[roundedCoords] = m1
	// 	return m1
	// }
	return realField.monomers2.GetMonomer(coords)
}

func (realField *RealField) roundCoords(coords base.Vector3DF) base.Vector3DF {
	return base.Vector3DF{
		math.Round(coords[base.AxisX]*10000) / 10000.0,
		math.Round(coords[base.AxisY]*10000) / 10000.0,
		math.Round(coords[base.AxisZ]*10000) / 10000.0,
	}
}

func (realField *RealField) MoveMonomer(monomer *Monomer, to base.Vector3DF) error {
	roundedCoords := realField.roundCoords(monomer.coords)
	if _, ok := realField.monomers[roundedCoords]; !ok {
		return fmt.Errorf("нет указанного мономера: %v", roundedCoords)
	}
	delete(realField.monomers, roundedCoords)
	monomer.coords = to
	realField.monomers[realField.roundCoords(to)] = monomer
	return nil
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

func (realField *RealField) DeepCopy() IField {
	newRealField := realField.CopyFieldOnly()
	for _, mon := range realField.monomers {
		newMon := newRealField.GetMonomerByCoords(mon.coords)
		newMon.DeepCopyFrom(mon, newRealField)
	}
	return newRealField
}

func (field *RealField) CopyFieldOnly() IField {
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	newRealField := NewRealField([3][2]float64{
		{spaceDimention[base.AxisX].Lower, spaceDimention[base.AxisX].Higher},
		{spaceDimention[base.AxisY].Lower, spaceDimention[base.AxisY].Higher},
		{spaceDimention[base.AxisZ].Lower, spaceDimention[base.AxisZ].Higher},
	})
	return newRealField
}
