package datatypes

import (
	"errors"
	"math/rand"
	"polymers/base"
	"polymers/globaldata"
)

type Field struct {
	sphereRadius uint64
	monomers     map[base.Vector3DF]*Monomer
}

func NewField(sphereRadius uint64) *Field {
	newField := &Field{}
	newField.monomers = make(map[base.Vector3DF]*Monomer)
	globalData := globaldata.GetGlobalData()
	lower := [...]int64{int64(globalData.SpaceDimention[base.AxisX].Lower), int64(globalData.SpaceDimention[base.AxisY].Lower), int64(globalData.SpaceDimention[base.AxisZ].Lower)}
	higher := [...]int64{int64(globalData.SpaceDimention[base.AxisX].Higher), int64(globalData.SpaceDimention[base.AxisY].Higher), int64(globalData.SpaceDimention[base.AxisZ].Higher)}
	var i int64
	var j int64
	var k int64
	for i = lower[0]; i <= higher[0]; i++ {
		for j = lower[1]; j <= higher[1]; j++ {
			for k = lower[2]; k <= higher[2]; k++ {
				coords := base.Vector3DF{float64(i), float64(j), float64(k)}
				newField.monomers[coords] = NewMonomer(coords, base.MendeleevTableElementUndefined)
			}
		}
	}
	for i = lower[0]; i <= higher[0]; i++ {
		for j = lower[1]; j <= higher[1]; j++ {
			for k = lower[2]; k <= higher[2]; k++ {
				coords := base.Vector3DF{float64(i), float64(j), float64(k)}
				monomer := newField.monomers[coords]
				if i < higher[0] {
					next := base.Vector3DF{float64(i + 1), float64(j), float64(k)}
					MakeConnection(monomer, newField.monomers[next], ConnectionTypeUndefined)
				}
				if i > lower[0] {
					next := base.Vector3DF{float64(i - 1), float64(j), float64(k)}
					MakeConnection(monomer, newField.monomers[next], ConnectionTypeUndefined)
				}
				if j < higher[1] {
					next := base.Vector3DF{float64(i), float64(j + 1), float64(k)}
					MakeConnection(monomer, newField.monomers[next], ConnectionTypeUndefined)
				}
				if j > lower[1] {
					next := base.Vector3DF{float64(i), float64(j - 1), float64(k)}
					MakeConnection(monomer, newField.monomers[next], ConnectionTypeUndefined)
				}
				if k < higher[2] {
					next := base.Vector3DF{float64(i), float64(j), float64(k + 1)}
					MakeConnection(monomer, newField.monomers[next], ConnectionTypeUndefined)
				}
				if k > lower[2] {
					next := base.Vector3DF{float64(i), float64(j), float64(k - 1)}
					MakeConnection(monomer, newField.monomers[next], ConnectionTypeUndefined)
				}
			}
		}
	}
	return newField
}

func (field *Field) MakeFilled(monomer *Monomer) {
	if monomer.IsTypeOf(base.MendeleevTableElementUndefined) {
		monomer.MonomerType = base.C
	}
}

func (field *Field) MakeFree(monomer *Monomer) {
	monomer.MonomerType = base.C
}

func (field *Field) IsFree(coords base.Vector3DF) bool {
	var monomer = field.monomers[coords]
	return monomer.IsTypeOf(base.MendeleevTableElementUndefined)
}

func (field *Field) GetSellWithinBorders(coords base.Vector3D) base.Vector3D {
	var globalData = globaldata.GetGlobalData()
	return base.Vector3D{
		coords[base.AxisX] % int64(globalData.SpaceDimention[base.AxisX].Higher),
		coords[base.AxisY] % int64(globalData.SpaceDimention[base.AxisY].Higher),
		coords[base.AxisZ] % int64(globalData.SpaceDimention[base.AxisZ].Higher),
	}
}

func (field *Field) IsBusy() bool {
	for _, mon := range field.monomers {
		if mon.IsTypeOf(base.MendeleevTableElementUndefined) {
			return false
		}
	}
	return true
}

func (field *Field) GetMonomerByCoords(coords base.Vector3DF) *Monomer {
	return field.monomers[coords]
}

func (field *Field) MoveMonomer(monomer *Monomer, to base.Vector3DF) error {
	return errors.New("Невозможно двигать мономеры на решётке")
}

func (field *Field) GetMonomersWithin(lower, higher base.Vector3DF) []*Monomer {
	isIn := func(point base.Vector3DF) bool {
		return base.PointInSpace(&point, &lower, &higher)
	}
	monomers := make([]*Monomer, 0)
	for point, mon := range field.monomers {
		if isIn(point) {
			monomers = append(monomers, mon)
		}
	}
	return monomers
}

func (field *Field) GetAvailableCells(currPos base.Vector3DF) []*Monomer {
	availableCells := make([]*Monomer, 0)
	monomer := field.GetMonomerByCoords(currPos)
	for _, side := range GetMovementSides() {
		sibling, err := monomer.GetSibling(side)
		if err == nil && sibling.IsTypeOf(base.MendeleevTableElementUndefined) {
			availableCells = append(availableCells, sibling)
		}
	}

	return availableCells
}

func (field *Field) DefineStartMonomer() *Monomer {
	globalData := globaldata.GetGlobalData()
	startPosition := globalData.SpaceDimention.GetCenter()

	for !field.IsFree(startPosition) {
		availableCells := field.GetAvailableCells(startPosition)
		availableCellsCout := len(availableCells)
		if availableCellsCout != 0 {
			return availableCells[rand.Intn(availableCellsCout)]
		}

		randomAxis := base.Axis(rand.Intn(int(base.AxisCount)))
		rint := rand.Intn(int(DirectionCount))
		randomDirection := DirectionBackward
		if rint%2 != 0 {
			randomDirection = DirectionForward
		}
		startMonomer, err := field.GetMonomerByCoords(startPosition).GetSibling(GetSide(randomAxis, randomDirection))
		if err != nil {
			continue
		}
		startPosition = startMonomer.Coords()
	}

	return field.GetMonomerByCoords(startPosition)
}

// func (field *Field) MarshalJSON() ([]byte, error) {
// 	newField := make([]MonomerJSON, 0)
// 	for _, mon := range field.monomers {
// 		newField = append(newField, mon.ToJSON())
// 	}
// 	return json.Marshal(&struct {
// 		SphereRadius uint64
// 		Field        []MonomerJSON
// 	}{
// 		SphereRadius: field.sphereRadius,
// 		Field:        newField,
// 	})
// }

func (field *Field) DeepCopy() IField {
	newField := NewField(field.sphereRadius)
	for coords, mon := range field.monomers {
		newField.monomers[coords].DeepCopyFrom(mon, newField)
	}
	return newField
}

func (field *Field) Waterize() {
	for _, mon := range field.monomers {
		if mon.IsTypeOf(base.MendeleevTableElementUndefined) {
			mon.MonomerType = base.N
		}
	}
}

func (field *Field) GetMinMonomersByAxis(axis base.Axis) []*Monomer {
	minMonomers := make([]*Monomer, 0)
	for point, mon := range field.monomers {
		if len(minMonomers) == 0 || point[axis] < minMonomers[0].Coords()[axis] {
			minMonomers = []*Monomer{mon}
		} else {
			minMonomers = append(minMonomers, mon)
		}
	}
	return minMonomers
}

func (field *Field) GetMaxMonomersByAxis(axis base.Axis) []*Monomer {
	minMonomers := make([]*Monomer, 0)
	for point, mon := range field.monomers {
		if len(minMonomers) == 0 || point[axis] > minMonomers[0].Coords()[axis] {
			minMonomers = []*Monomer{mon}
		} else {
			minMonomers = append(minMonomers, mon)
		}
	}
	return minMonomers
}
