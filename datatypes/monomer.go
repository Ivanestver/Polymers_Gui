package datatypes

import (
	"errors"
	"polymers/base"
)

type Monomer struct {
	coords      base.Vector3DF
	MonomerType MonomerType
	PrevMonomer *Monomer
	NextMonomer *Monomer
	sides       map[Side]*Connection
	Number      int64
}

func NewMonomer(coords base.Vector3DF, monomerType MonomerType) *Monomer {
	newMonomer := new(Monomer)
	newMonomer.coords = coords
	newMonomer.MonomerType = monomerType
	newMonomer.PrevMonomer = nil
	newMonomer.NextMonomer = nil
	newMonomer.Number = -1
	newMonomer.sides = make(map[Side]*Connection)
	return newMonomer
}

func (monomer *Monomer) IsTypeOf(monType MonomerType) bool {
	return monomer.MonomerType == monType
}

func (monomer *Monomer) IsNotTypeOf(monType MonomerType) bool {
	return monomer.MonomerType != monType
}

func (monomer *Monomer) GetSibling(side Side) (*Monomer, error) {
	conn := monomer.sides[side]
	if conn != nil {
		return conn.GetOtherSide(monomer)
	} else {
		return nil, errors.New("no sibling")
	}
}

func (monomer *Monomer) GetSideOfSibling(otherMon *Monomer) Side {
	for side, conn := range monomer.sides {
		if conn == nil {
			continue
		}

		sibling, _ := conn.GetOtherSide(monomer)
		if MonomersAreEqual(otherMon, sibling) {
			return side
		}
	}

	return SideUndefined
}

func (monomer *Monomer) GetTypeOfConnectionWithSide(side Side) ConnectionType {
	conn, ok := monomer.sides[side]
	if conn != nil && ok {
		return conn.ConnType
	} else {
		return ConnectionTypeUndefined
	}
}

func (monomer *Monomer) Coords() base.Vector3DF {
	return monomer.coords
}

func (monomer *Monomer) Copy() *Monomer {
	newMon := NewMonomer(monomer.coords, monomer.MonomerType)
	newMon.NextMonomer = monomer.NextMonomer
	newMon.PrevMonomer = monomer.PrevMonomer
	newMon.Number = monomer.Number
	newMon.sides = monomer.sides
	return newMon
}

func (monomer *Monomer) DeepCopy(field *Field) *Monomer {
	newMon := NewMonomer(monomer.coords, monomer.MonomerType)

	if monomer.NextMonomer != nil {
		newMon.NextMonomer = field.GetMonomerByCoords(monomer.NextMonomer.coords)
	} else {
		newMon.NextMonomer = nil
	}

	if monomer.PrevMonomer != nil {
		newMon.PrevMonomer = field.GetMonomerByCoords(monomer.PrevMonomer.coords)
	} else {
		newMon.PrevMonomer = nil
	}

	newMon.Number = monomer.Number
	newMon.sides = monomer.sides
	return newMon
}

func (monomer *Monomer) DeepCopyFrom(other *Monomer, field *Field) {
	if monomer.coords != other.coords {
		return
	}
	if other.NextMonomer != nil {
		monomer.NextMonomer = field.GetMonomerByCoords(other.NextMonomer.coords)
	} else {
		monomer.NextMonomer = nil
	}

	if other.PrevMonomer != nil {
		monomer.PrevMonomer = field.GetMonomerByCoords(other.PrevMonomer.coords)
	} else {
		monomer.PrevMonomer = nil
	}

	monomer.Number = other.Number
	for _, conn := range other.sides {
		otherSideOfOther, err := conn.GetOtherSide(other)
		if err != nil {
			continue
		}
		MakeConnection(monomer, field.GetMonomerByCoords(otherSideOfOther.coords), conn.ConnType)
	}
	monomer.MonomerType = other.MonomerType
}

func MonomersAreEqual(left, right *Monomer) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	return base.VectorsAreEqualF(&left.coords, &right.coords) &&
		left.Number == right.Number
}

func GetSideByMonomers(from, to *Monomer) Side {
	if from == nil || to == nil {
		return SideUndefined
	}
	return base.SubtractVecF(to.coords, from.coords)
}

type MonomerJSON struct {
	Coords      base.Vector3DF
	MonomerType MonomerType
	PrevMonomer base.Vector3DF
	NextMonomer base.Vector3DF
	Sides       map[Side]ConnectionJSON
	Number      int64
}

func (monomer *Monomer) ToJSON() MonomerJSON {
	var obj MonomerJSON
	obj.Coords = monomer.coords
	obj.MonomerType = monomer.MonomerType
	if monomer.PrevMonomer != nil {
		obj.PrevMonomer = monomer.PrevMonomer.coords
	}
	if monomer.NextMonomer != nil {
		obj.NextMonomer = monomer.NextMonomer.coords
	}
	obj.Sides = make(map[Side]ConnectionJSON)
	for key, value := range monomer.sides {
		obj.Sides[key] = value.ToJSON()
	}
	obj.Number = monomer.Number
	return obj
}
