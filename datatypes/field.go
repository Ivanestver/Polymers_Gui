package datatypes

import (
	"encoding/json"
	"math/rand"
	"polymers/base"
	"polymers/global_data"
)

type Field struct {
	sphereRadius uint64
	field        [][][]*Monomer
}

func NewField(sphereRadius uint64) *Field {
	newField := &Field{}
	var globalData *global_data.GlobalData = global_data.GetGlobalData()
	shape := [...]int64{int64(globalData.SpaceDimention[base.X_AXIS].Higher), int64(globalData.SpaceDimention[base.Y_AXIS].Higher), int64(globalData.SpaceDimention[base.Z_AXIS].Higher)}
	var i int64
	var j int64
	var k int64
	for i = 0; i < shape[0]; i++ {
		newField.field = append(newField.field, [][]*Monomer{})
		for j = 0; j < shape[1]; j++ {
			newField.field[i] = append(newField.field[i], []*Monomer{})
			for k = 0; k < shape[2]; k++ {
				newField.field[i][j] = append(newField.field[i][j], NewMonomer(base.Vector3DF{X: float64(i), Y: float64(j), Z: float64(k)}, MONOMER_TYPE_UNDEFINED))
			}
		}
	}
	for i = 0; i < shape[0]; i++ {
		for j = 0; j < shape[1]; j++ {
			for k = 0; k < shape[2]; k++ {
				monomer := newField.field[i][j][k]
				if i < shape[0]-1 {
					MakeConnection(monomer, newField.field[i+1][j][k], CONNECTION_TYPE_UNDEFINED)
				}
				if i > 0 {
					MakeConnection(monomer, newField.field[i-1][j][k], CONNECTION_TYPE_UNDEFINED)
				}
				if j < shape[1]-1 {
					MakeConnection(monomer, newField.field[i][j+1][k], CONNECTION_TYPE_UNDEFINED)
				}
				if j > 0 {
					MakeConnection(monomer, newField.field[i][j-1][k], CONNECTION_TYPE_UNDEFINED)
				}
				if k < shape[2]-1 {
					MakeConnection(monomer, newField.field[i][j][k+1], CONNECTION_TYPE_UNDEFINED)
				}
				if k > 0 {
					MakeConnection(monomer, newField.field[i][j][k-1], CONNECTION_TYPE_UNDEFINED)
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
	var monomer *Monomer = field.field[int(coords.X)][int(coords.Y)][int(coords.Z)]
	return monomer.IsTypeOf(MONOMER_TYPE_UNDEFINED)
}

func (field *Field) GetSellWithinBorders(coords base.Vector3D) base.Vector3D {
	var globalData *global_data.GlobalData = global_data.GetGlobalData()
	return base.Vector3D{
		X: coords.X % int64(globalData.SpaceDimention[base.X_AXIS].Higher),
		Y: coords.Y % int64(globalData.SpaceDimention[base.Y_AXIS].Higher),
		Z: coords.Z % int64(globalData.SpaceDimention[base.Z_AXIS].Higher),
	}
}

func (field *Field) IsBusy() bool {
	for _, v_x := range field.field {
		for _, v_y := range v_x {
			for _, val := range v_y {
				if val.IsTypeOf(MONOMER_TYPE_UNDEFINED) {
					return false
				}
			}
		}
	}
	return true
}

func (field *Field) GetMonomerByCoords(coords base.Vector3DF) *Monomer {
	return field.field[int64(coords.X)][int64(coords.Y)][int64(coords.Z)]
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
	newField := make([][][]MonomerJSON, len(field.field))
	for i, square := range field.field {
		newField[i] = make([][]MonomerJSON, len(field.field[i]))
		for j, row := range square {
			newField[i][j] = make([]MonomerJSON, len(field.field[i][j]))
			for k, item := range row {
				newField[i][j][k] = item.ToJson()
			}
		}
	}
	return json.Marshal(&struct {
		SphereRadius uint64
		Field        [][][]MonomerJSON
	}{
		SphereRadius: field.sphereRadius,
		Field:        newField,
	})
}

func (field *Field) DeepCopy() *Field {
	newField := NewField(field.sphereRadius)
	for x, monX := range field.field {
		for y, monY := range monX {
			for z := range monY {
				newField.field[x][y][z].DeepCopyFrom(field.field[x][y][z], newField)
			}
		}
	}
	return newField
}

func (field *Field) Waterize() {
	for x, monX := range field.field {
		for y, monY := range monX {
			for z := range monY {
				mon := field.field[x][y][z]
				if mon.MonomerType == MONOMER_TYPE_UNDEFINED {
					mon.MonomerType = MONOMER_TYPE_WATER
				}
			}
		}
	}
}
