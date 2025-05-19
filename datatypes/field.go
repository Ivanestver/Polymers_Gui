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

func makeConnection(mon1, mon2 *Monomer, canMake bool) {
	if canMake {
		MakeConnection(mon1, mon2, CONNECTION_TYPE_UNDEFINED)
	}
}

func NewField(sphereRadius uint64) *Field {
	newField := &Field{}
	var globalData *global_data.GlobalData = global_data.GetGlobalData()
	shape := [...]int64{globalData.SpaceDimention, globalData.SpaceDimention, globalData.SpaceDimention}
	var i int64
	var j int64
	var k int64
	for i = 0; i < shape[0]; i++ {
		newField.field = append(newField.field, [][]*Monomer{})
		for j = 0; j < shape[1]; j++ {
			newField.field[i] = append(newField.field[i], []*Monomer{})
			for k = 0; k < shape[2]; k++ {
				newField.field[i][j] = append(newField.field[i][j], NewMonomer(base.Vector3D{X: i, Y: j, Z: k}, MONOMER_TYPE_UNDEFINED))
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
	monomer.MonomerType = MONOMER_TYPE_USUAL
}

func (field *Field) MakeFree(monomer *Monomer) {
	monomer.MonomerType = MONOMER_TYPE_UNDEFINED
}

func (field *Field) IsFree(coords base.Vector3D) bool {
	var monomer *Monomer = field.field[coords.X][coords.Y][coords.Z]
	return monomer.IsTypeOf(MONOMER_TYPE_UNDEFINED)
}

func (field *Field) GetSellWithinBorders(coords base.Vector3D) base.Vector3D {
	var globalData *global_data.GlobalData = global_data.GetGlobalData()
	return base.Vector3D{
		X: coords.X % globalData.SpaceDimention,
		Y: coords.Y % globalData.SpaceDimention,
		Z: coords.Z % globalData.SpaceDimention,
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

func (field *Field) GetMonomerByCoords(coords base.Vector3D) *Monomer {
	return field.field[coords.X][coords.Y][coords.Z]
}

func (field *Field) GetAvailableCells(currPos base.Vector3D) []*Monomer {
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

func (this *Field) DefineStartMonomer() *Monomer {
	globalData := global_data.GetGlobalData()
	startPosition := base.Vector3D{
		X: globalData.SpaceDimention / 2,
		Y: globalData.SpaceDimention / 2,
		Z: globalData.SpaceDimention / 2,
	}

	for !this.IsFree(startPosition) {
		available_cells := this.GetAvailableCells(startPosition)
		available_cells_cout := len(available_cells)
		if available_cells_cout != 0 {
			return available_cells[rand.Intn(available_cells_cout)]
		}

		randomAxis := Axis(rand.Intn(int(AXIS_COUNT)))
		rint := rand.Intn(int(DIRECTION_COUNT))
		randomDirection := DIRECTION_BACKWARD
		if rint%2 != 0 {
			randomDirection = DIRECTION_FORWARD
		}
		startMonomer, err := this.GetMonomerByCoords(startPosition).GetSibling(GetSide(randomAxis, randomDirection))
		if err != nil {
			continue
		}
		startPosition = startMonomer.Coords()
	}

	return this.GetMonomerByCoords(startPosition)
}

func (this *Field) MarshalJSON() ([]byte, error) {
	field := make([][][]MonomerJSON, len(this.field))
	for i, square := range this.field {
		field[i] = make([][]MonomerJSON, len(this.field[i]))
		for j, row := range square {
			field[i][j] = make([]MonomerJSON, len(this.field[i][j]))
			for k, item := range row {
				field[i][j][k] = item.ToJson()
			}
		}
	}
	return json.Marshal(&struct {
		SphereRadius uint64
		Field        [][][]MonomerJSON
	}{
		SphereRadius: this.sphereRadius,
		Field:        field,
	})
}
