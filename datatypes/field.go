package datatypes

import (
	"encoding/json"
	"math/rand"
	"polymers/base"
	"polymers/global_data"
)

type Field struct {
	sphereRadius uint64
	field        map[base.Vector3DF]*Monomer
}

func NewField(sphereRadius uint64) *Field {
	newField := &Field{}
	newField.field = make(map[base.Vector3DF]*Monomer)
	globalData := global_data.GetGlobalData()
	lower := [...]int64{int64(globalData.SpaceDimention[base.X_AXIS].Lower), int64(globalData.SpaceDimention[base.Y_AXIS].Lower), int64(globalData.SpaceDimention[base.Z_AXIS].Lower)}
	higher := [...]int64{int64(globalData.SpaceDimention[base.X_AXIS].Higher), int64(globalData.SpaceDimention[base.Y_AXIS].Higher), int64(globalData.SpaceDimention[base.Z_AXIS].Higher)}
	var i int64
	var j int64
	var k int64
	for i = lower[0]; i <= higher[0]; i++ {
		for j = lower[1]; j <= higher[1]; j++ {
			for k = lower[2]; k <= higher[2]; k++ {
				coords := base.Vector3DF{X: float64(i), Y: float64(j), Z: float64(k)}
				newField.field[coords] = NewMonomer(coords, MONOMER_TYPE_UNDEFINED)
			}
		}
	}
	for i = lower[0]; i <= higher[0]; i++ {
		for j = lower[1]; j <= higher[1]; j++ {
			for k = lower[2]; k <= higher[2]; k++ {
				coords := base.Vector3DF{X: float64(i), Y: float64(j), Z: float64(k)}
				monomer := newField.field[coords]
				if i < higher[0] {
					next := base.Vector3DF{X: float64(i + 1), Y: float64(j), Z: float64(k)}
					MakeConnection(monomer, newField.field[next], CONNECTION_TYPE_UNDEFINED)
				}
				if i > lower[0] {
					next := base.Vector3DF{X: float64(i - 1), Y: float64(j), Z: float64(k)}
					MakeConnection(monomer, newField.field[next], CONNECTION_TYPE_UNDEFINED)
				}
				if j < higher[1] {
					next := base.Vector3DF{X: float64(i), Y: float64(j + 1), Z: float64(k)}
					MakeConnection(monomer, newField.field[next], CONNECTION_TYPE_UNDEFINED)
				}
				if j > lower[1] {
					next := base.Vector3DF{X: float64(i), Y: float64(j - 1), Z: float64(k)}
					MakeConnection(monomer, newField.field[next], CONNECTION_TYPE_UNDEFINED)
				}
				if k < higher[2] {
					next := base.Vector3DF{X: float64(i), Y: float64(j), Z: float64(k + 1)}
					MakeConnection(monomer, newField.field[next], CONNECTION_TYPE_UNDEFINED)
				}
				if k > lower[2] {
					next := base.Vector3DF{X: float64(i), Y: float64(j), Z: float64(k - 1)}
					MakeConnection(monomer, newField.field[next], CONNECTION_TYPE_UNDEFINED)
				}
			}
		}
	}
	return newField
}

func (field *Field) MakeFilled(monomer *Monomer) {
	if monomer.MonomerType == MONOMER_TYPE_UNDEFINED {
		monomer.MonomerType = MONOMER_TYPE_USUAL
	}
}

func (field *Field) MakeFree(monomer *Monomer) {
	monomer.MonomerType = MONOMER_TYPE_UNDEFINED
}

func (field *Field) IsFree(coords base.Vector3DF) bool {
	var monomer = field.field[coords]
	return monomer.IsTypeOf(MONOMER_TYPE_UNDEFINED)
}

func (field *Field) GetSellWithinBorders(coords base.Vector3D) base.Vector3D {
	var globalData = global_data.GetGlobalData()
	return base.Vector3D{
		X: coords.X % int64(globalData.SpaceDimention[base.X_AXIS].Higher),
		Y: coords.Y % int64(globalData.SpaceDimention[base.Y_AXIS].Higher),
		Z: coords.Z % int64(globalData.SpaceDimention[base.Z_AXIS].Higher),
	}
}

func (field *Field) IsBusy() bool {
	for _, mon := range field.field {
		if mon.IsTypeOf(MONOMER_TYPE_UNDEFINED) {
			return false
		}
	}
	return true
}

func (field *Field) GetMonomerByCoords(coords base.Vector3DF) *Monomer {
	return field.field[coords]
}

func (field *Field) GetMonomersWithin(lower, higher base.Vector3DF) []*Monomer {
	isIn := func(point base.Vector3DF) bool {
		return base.PointInSpace(&point, &lower, &higher)
	}
	monomers := make([]*Monomer, 0)
	for point, mon := range field.field {
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
		if err == nil && sibling.IsTypeOf(MONOMER_TYPE_UNDEFINED) {
			availableCells = append(availableCells, sibling)
		}
	}

	return availableCells
}

func (field *Field) DefineStartMonomer() *Monomer {
	globalData := global_data.GetGlobalData()
	startPosition := globalData.SpaceDimention.GetCenter()

	for !field.IsFree(startPosition) {
		available_cells := field.GetAvailableCells(startPosition)
		available_cells_cout := len(available_cells)
		if available_cells_cout != 0 {
			return available_cells[rand.Intn(available_cells_cout)]
		}

		randomAxis := base.Axis(rand.Intn(int(base.AXIS_COUNT)))
		rint := rand.Intn(int(DIRECTION_COUNT))
		randomDirection := DIRECTION_BACKWARD
		if rint%2 != 0 {
			randomDirection = DIRECTION_FORWARD
		}
		startMonomer, err := field.GetMonomerByCoords(startPosition).GetSibling(GetSide(randomAxis, randomDirection))
		if err != nil {
			continue
		}
		startPosition = startMonomer.Coords()
	}

	return field.GetMonomerByCoords(startPosition)
}

func (field *Field) MarshalJSON() ([]byte, error) {
	newField := make([]MonomerJSON, 0)
	for _, mon := range field.field {
		newField = append(newField, mon.ToJson())
	}
	return json.Marshal(&struct {
		SphereRadius uint64
		Field        []MonomerJSON
	}{
		SphereRadius: field.sphereRadius,
		Field:        newField,
	})
}

func (field *Field) DeepCopy() *Field {
	newField := NewField(field.sphereRadius)
	for coords, mon := range field.field {
		newField.field[coords].DeepCopyFrom(mon, newField)
	}
	return newField
}

func (field *Field) Waterize() {
	for _, mon := range field.field {
		if mon.MonomerType == MONOMER_TYPE_UNDEFINED {
			mon.MonomerType = MONOMER_TYPE_WATER
		}
	}
}
