package datatypes

import (
	"errors"
	"polymers/base"
)

type Connection struct {
	monomers [2]*Monomer
	ConnType ConnectionType
}

func NewConnection(monomers [2]*Monomer, connType ConnectionType) *Connection {
	newConn := new(Connection)
	copy(newConn.monomers[:], monomers[:])
	newConn.ConnType = connType
	return newConn
}

func (conn *Connection) GetOtherSide(currMonomer *Monomer) (*Monomer, error) {
	if MonomersAreEqual(conn.monomers[0], currMonomer) {
		return conn.monomers[1], nil
	} else if MonomersAreEqual(conn.monomers[1], currMonomer) {
		return conn.monomers[0], nil
	} else {
		return nil, errors.New("Connection is corrupted")
	}
}

func MakeConnection(mon1, mon2 *Monomer, connectionType ConnectionType) error {
	if mon1 == nil || mon2 == nil {
		return errors.New("one of monomers is nil")
	}
	side := GetSideByMonomers(mon1, mon2)
	if side == SideUndefined {
		return errors.New("cannot connect non-contiguous monomers")
	}
	conn := mon1.sides[side]
	if conn == nil {
		newConn := NewConnection([2]*Monomer{mon1, mon2}, connectionType)
		mon1.sides[side] = newConn
		mon2.sides[GetReversedSide(side)] = newConn
	} else {
		_, err := conn.GetOtherSide(mon1)
		if err != nil {
			panic("Two monomers occupy the same location")
		}
		conn.ConnType = connectionType
	}
	return nil
}

func MakeConnectionUnsafe(mon1, mon2 *Monomer, side Side, connectionType ConnectionType) error {
	if mon1 == nil || mon2 == nil {
		return errors.New("either of monomers is nil")
	}
	if side == SideUndefined {
		return errors.New("Side is undefined")
	}
	newConn := NewConnection([2]*Monomer{mon1, mon2}, connectionType)
	mon1.sides[side] = newConn
	mon2.sides[GetReversedSide(side)] = newConn
	return nil
}

func BreakConnection(mon1, mon2 *Monomer, side Side) error {

	if mon1 == nil || mon2 == nil {
		return errors.New("one of monomers is nil")
	}

	reversedSide := GetReversedSide(side)
	currConn, ok := mon1.sides[side]
	if !ok {
		return errors.New("could not receive the connection")
	}
	otherConn, ok := mon2.sides[reversedSide]
	if !ok {
		return errors.New("could not receive the connection on the other side")
	}
	if currConn != otherConn {
		return errors.New("two connection occupy the same place")
	}
	if MonomersAreEqual(mon1.NextMonomer, mon2) ||
		MonomersAreEqual(mon1.PrevMonomer, mon2) {
		currConn.ConnType = ConnectionTypeOne
	} else {
		currConn.ConnType = ConnectionTypeUndefined
	}

	return nil
}

func BreakConnection1(mon1, mon2 *Monomer) error {
	if mon1 == nil || mon2 == nil {
		return errors.New("one of monomers is nil")
	}
	side := mon1.GetSideOfSibling(mon2)
	if side == SideUndefined {
		return errors.New("the monomers are not sibiings")
	}
	mon1.sides[side].ConnType = ConnectionTypeUndefined
	return nil
}

func TierConnection(mon1, mon2 *Monomer, side Side) error {
	if mon1 == nil || mon2 == nil {
		return errors.New("one of monomers is nil")
	}

	reversedSide := GetReversedSide(side)
	currConn, ok := mon1.sides[side]
	if !ok {
		return errors.New("could not receive the connection")
	}
	otherConn, ok := mon2.sides[reversedSide]
	if !ok {
		return errors.New("could not receive the connection on the other side")
	}
	if currConn != otherConn {
		return errors.New("two connection occupy the same place")
	}
	currConn.ConnType = ConnectionTypeUndefined

	return nil
}

func GetConnectionType(sideOne, sideTwo *Monomer) ConnectionType {
	side := GetSideByMonomers(sideOne, sideTwo)
	if base.Contains(GetMovementSides(), side) {
		if (sideOne.NextMonomer != nil && sideOne.NextMonomer == sideTwo) ||
			(sideOne.PrevMonomer != nil && sideOne.PrevMonomer == sideTwo) {
			return ConnectionTypeOne
		} else {
			return ConnectionTypeCrossLinear
		}
	} else if base.Contains(GetSurfaceDiagonalSides(), side) {
		return ConnectionTypeCrossSurface
	} else if base.Contains(GetCubeDiagonalSides(), side) {
		return ConnectionTypeCrossSpacial
	} else {
		return ConnectionTypeUndefined
	}
}

type ConnectionJSON struct {
	monomers [2]base.Vector3DF
	ConnType ConnectionType
}

func (conn *Connection) ToJSON() ConnectionJSON {
	var jsonObj ConnectionJSON
	jsonObj.ConnType = conn.ConnType
	jsonObj.monomers[0] = conn.monomers[0].Coords()
	jsonObj.monomers[1] = conn.monomers[1].Coords()
	return jsonObj
}
