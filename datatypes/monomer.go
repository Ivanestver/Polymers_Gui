package datatypes

import (
	"errors"
	"polymers/base"
)

type Monomer struct {
	coords      base.Vector3D
	MonomerType MonomerType
	PrevMonomer *Monomer
	NextMonomer *Monomer
	sides       map[Side]*Connection
	Number      int64
}

func NewMonomer(coords base.Vector3D, monomerType MonomerType) *Monomer {
	newMonomer := new(Monomer)
	newMonomer.coords = coords
	newMonomer.MonomerType = monomerType
	newMonomer.PrevMonomer = nil
	newMonomer.NextMonomer = nil
	newMonomer.Number = -1
	newMonomer.sides = make(map[Side]*Connection)
	return newMonomer
}

func (mon *Monomer) IsTypeOf(monType MonomerType) bool {
	return mon.MonomerType == monType
}

func (mon *Monomer) IsNotTypeOf(monType MonomerType) bool {
	return mon.MonomerType != monType
}

func (mon *Monomer) GetSibling(side Side) (*Monomer, error) {
	conn := mon.sides[side]
	if conn != nil {
		return conn.GetOtherSide(mon)
	} else {
		return nil, errors.New("No sibling")
	}
}

func (mon *Monomer) GetSideOfSibling(otherMon *Monomer) Side {
	for side, conn := range mon.sides {
		if conn == nil {
			continue
		}

		sibling, _ := conn.GetOtherSide(mon)
		if otherMon == sibling {
			return side
		}
	}

	return SIDE_Undefined
}

func (this *Monomer) GetTypeOfConnectionWithSide(side Side) ConnectionType {
	conn, ok := this.sides[side]
	if conn != nil && ok {
		return conn.ConnType
	} else {
		return CONNECTION_TYPE_UNDEFINED
	}
}

func (mon *Monomer) Coords() base.Vector3D {
	return mon.coords
}

func (this *Monomer) Copy() *Monomer {
	newMon := NewMonomer(this.coords, this.MonomerType)
	newMon.NextMonomer = this.NextMonomer
	newMon.PrevMonomer = this.PrevMonomer
	newMon.Number = this.Number
	newMon.sides = this.sides
	return newMon
}

func MonomersAreEqual(left, right *Monomer) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	return base.VectorsAreEqual(&left.coords, &right.coords) &&
		left.Number == right.Number
}

func GetSideByMonomers(from, to *Monomer) Side {
	if from == nil || to == nil {
		return SIDE_Undefined
	}

	result := base.Vector3D{
		X: to.coords.X - from.coords.X,
		Y: to.coords.Y - from.coords.Y,
		Z: to.coords.Z - from.coords.Z,
	}
	if result.X == 0 {
		if result.Y == 0 {
			if result.Z == 0 {
				return SIDE_Undefined
			} else if result.Z > 0 {
				return SIDE_Up
			} else {
				return SIDE_Down
			}
		} else if result.Y > 0 {
			if result.Z == 0 {
				return SIDE_Left
			} else if result.Z > 0 {
				return SIDE_LeftUp
			} else {
				return SIDE_LeftDown
			}
		} else if result.Y < 0 {
			if result.Z == 0 {
				return SIDE_Right
			} else if result.Z > 0 {
				return SIDE_RightUp
			} else {
				return SIDE_RightDown
			}
		}
	} else if result.X > 0 {
		if result.Y == 0 {
			if result.Z == 0 {
				return SIDE_Forward
			} else if result.Z > 0 {
				return SIDE_UpForward
			} else {
				return SIDE_DownForward
			}
		} else if result.Y > 0 {
			if result.Z == 0 {
				return SIDE_LeftForward
			} else if result.Z > 0 {
				return SIDE_UpLeftForward
			} else {
				return SIDE_DownLeftForward
			}
		} else if result.Y < 0 {
			if result.Z == 0 {
				return SIDE_RightForward
			} else if result.Z > 0 {
				return SIDE_UpRightForward
			} else {
				return SIDE_DownRightForward
			}
		}
	} else { // result.X < 0
		if result.Y == 0 {
			if result.Z == 0 {
				return SIDE_Backward
			} else if result.Z > 0 {
				return SIDE_UpBackward
			} else {
				return SIDE_DownBackward
			}
		} else if result.Y > 0 {
			if result.Z == 0 {
				return SIDE_LeftBackward
			} else if result.Z > 0 {
				return SIDE_UpLeftBackward
			} else {
				return SIDE_DownLeftBackward
			}
		} else if result.Y < 0 {
			if result.Z == 0 {
				return SIDE_RightBackward
			} else if result.Z > 0 {
				return SIDE_UpRightBackward
			} else {
				return SIDE_DownRightBackward
			}
		}
	}

	return SIDE_Undefined
}

type MonomerJSON struct {
	Coords      base.Vector3D
	MonomerType MonomerType
	PrevMonomer base.Vector3D
	NextMonomer base.Vector3D
	Sides       map[Side]ConnectionJSON
	Number      int64
}

func (this *Monomer) ToJson() MonomerJSON {
	var obj MonomerJSON
	obj.Coords = this.coords
	obj.MonomerType = this.MonomerType
	if this.PrevMonomer != nil {
		obj.PrevMonomer = this.PrevMonomer.coords
	}
	if this.NextMonomer != nil {
		obj.NextMonomer = this.NextMonomer.coords
	}
	obj.Sides = make(map[Side]ConnectionJSON)
	for key, value := range this.sides {
		obj.Sides[key] = value.ToJson()
	}
	obj.Number = this.Number
	return obj
}
