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
		return errors.New("One of monomers is nil")
	}
	side := GetSideByMonomers(mon1, mon2)
	conn, _ := mon1.sides[side]
	if conn == nil {
		newConn := NewConnection([2]*Monomer{mon1, mon2}, connectionType)
		mon1.sides[side] = newConn
		mon2.sides[GetReversedSide(side)] = newConn
	} else {
		_, err := conn.GetOtherSide(mon1)
		if err != nil {
			panic("Two monomers occupy the same location&")
		}
		conn.ConnType = connectionType
	}
	return nil
}

func BreakConnection(mon1, mon2 *Monomer, side Side) error {

	if mon1 == nil || mon2 == nil {
		return errors.New("One of monomers is nil")
	}

	reversed_side := GetReversedSide(side)
	curr_conn, ok := mon1.sides[side]
	if !ok {
		return errors.New("Could not receive the connection")
	}
	other_conn, ok := mon2.sides[reversed_side]
	if !ok {
		return errors.New("Could not receive the connection on the other side")
	}
	if curr_conn != other_conn {
		return errors.New("Two connection occupy the same place")
	}
	if MonomersAreEqual(mon1.NextMonomer, mon2) ||
		MonomersAreEqual(mon1.PrevMonomer, mon2) {
		curr_conn.ConnType = CONNECTION_TYPE_ONE
	} else {
		curr_conn.ConnType = CONNECTION_TYPE_UNDEFINED
	}

	return nil
}

func TierConnection(mon1, mon2 *Monomer, side Side) error {
	if mon1 == nil || mon2 == nil {
		return errors.New("One of monomers is nil")
	}

	reversed_side := GetReversedSide(side)
	curr_conn, ok := mon1.sides[side]
	if !ok {
		return errors.New("Could not receive the connection")
	}
	other_conn, ok := mon2.sides[reversed_side]
	if !ok {
		return errors.New("Could not receive the connection on the other side")
	}
	if curr_conn != other_conn {
		return errors.New("Two connection occupy the same place")
	}
	curr_conn.ConnType = CONNECTION_TYPE_UNDEFINED

	return nil
}

func GetConnectionType(sideOne, sideTwo *Monomer) ConnectionType {
	side := GetSideByMonomers(sideOne, sideTwo)
	if base.Contains[Side](GetMovementSides(), side) {
		return CONNECTION_TYPE_ONE
	} else if base.Contains[Side](GetSurfaceDiagonalSides(), side) {
		return CONNECTION_TYPE_TWO
	} else if base.Contains[Side](GetCubeDiagonalSides(), side) {
		return CONNECTION_TYPE_THREE
	} else {
		return CONNECTION_TYPE_UNDEFINED
	}
}

type ConnectionJSON struct {
	monomers [2]base.Vector3D
	ConnType ConnectionType
}

func (this *Connection) ToJson() ConnectionJSON {
	var jsonObj ConnectionJSON
	jsonObj.ConnType = this.ConnType
	jsonObj.monomers[0] = this.monomers[0].Coords()
	jsonObj.monomers[1] = this.monomers[1].Coords()
	return jsonObj
}
