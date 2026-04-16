package datatypes

import (
	"errors"
	"polymers/base"
)

type MoveDirection = int8
type Side = base.Vector3DF
type MonomerType int8
type ConnectionType int8
type GlobulaViewType int

const (
	DirectionBackward MoveDirection = -1
	DirectionForward  MoveDirection = 1
	DirectionCount    MoveDirection = 2
)

var SideUndefined = Side{0.0, 0.0, 0.0}
var SideForward = base.AxisXVec
var SideBackward = base.MultiplyByConstantF(&SideForward, float64(DirectionBackward))
var SideLeft = base.AxisYVec
var SideRight = base.MultiplyByConstantF(&SideLeft, float64(DirectionBackward))
var SideUp = base.AxisZVec
var SideDown = base.MultiplyByConstantF(&SideUp, float64(DirectionBackward))
var SideUpLeftForward = base.AddVecF(base.AddVecF(SideUp, SideLeft), SideForward)
var SideUpForward = base.AddVecF(SideUp, SideForward)
var SideUpRightForward = base.AddVecF(SideUp, SideRight, SideForward)
var SideLeftForward = base.AddVecF(SideLeft, SideForward)
var SideRightForward = base.AddVecF(SideRight, SideForward)
var SideDownLeftForward = base.AddVecF(SideDown, SideLeft, SideForward)
var SideDownForward = base.AddVecF(SideDown, SideForward)
var SideDownRightForward = base.AddVecF(SideDown, SideRight, SideForward)
var SideLeftUp = base.AddVecF(SideLeft, SideUp)
var SideLeftDown = base.AddVecF(SideLeft, SideDown)
var SideRightUp = base.AddVecF(SideRight, SideUp)
var SideRightDown = base.AddVecF(SideRight, SideDown)
var SideUpLeftBackward = base.AddVecF(SideUp, SideLeft, SideBackward)
var SideUpBackward = base.AddVecF(SideUp, SideBackward)
var SideUpRightBackward = base.AddVecF(SideUp, SideRight, SideBackward)
var SideLeftBackward = base.AddVecF(SideLeft, SideBackward)
var SideRightBackward = base.AddVecF(SideRight, SideBackward)
var SideDownLeftBackward = base.AddVecF(SideDown, SideLeft, SideBackward)
var SideDownBackward = base.AddVecF(SideDown, SideBackward)
var SideDownRightBackward = base.AddVecF(SideDown, SideRight, SideBackward)

const (
	MonomerTypeUndefined   MonomerType = -1
	MonomerTypeUsual       MonomerType = 0
	MonomerTypeOContaining MonomerType = 1
	MonomerTypeVynil       MonomerType = 2
	MonomerTypeCrosslinked MonomerType = 3
	MonomerTypeFwise       MonomerType = 4
	MonomerTypeClwise      MonomerType = 5
	MonomerTypeWater       MonomerType = 6
	MonomerTypeS           MonomerType = 7
)

const (
	ConnectionTypeUndefined ConnectionType = iota
	ConnectionTypeOne
	ConnectionTypeCrosslinks
	ConnectionTypeCrossLinear
	ConnectionTypeCrossSurface
	ConnectionTypeCrossSpacial
)

const (
	GlobulaViewTypeGlobula GlobulaViewType = iota
	GlobulaViewTypeThread
)

func GetMovementSides() []Side {
	return []Side{SideLeft, SideRight, SideForward, SideBackward, SideUp, SideDown}
}

func GetSurfaceDiagonalSides() []Side {
	return []Side{SideUpForward, SideUpBackward, SideDownForward, SideDownBackward, SideLeftUp, SideLeftDown, SideRightUp, SideRightDown, SideLeftForward, SideLeftBackward, SideRightForward, SideRightBackward}
}

func GetCubeDiagonalSides() []Side {
	return []Side{SideUpLeftForward, SideUpLeftBackward, SideUpRightForward, SideUpRightBackward, SideDownLeftForward, SideDownLeftBackward, SideDownRightForward, SideDownRightBackward}
}

func GetAdditionalSides() []Side {
	return append(GetSurfaceDiagonalSides(), GetCubeDiagonalSides()...)
}

func GetAllSides() []Side {
	return append(GetMovementSides(), GetAdditionalSides()...)
}

func GetSide(axis base.Axis, moveDirection MoveDirection) Side {
	v := base.AxisToVector[axis]
	v.MultiplyByConstantF(float64(moveDirection))
	return v
}

func GetReversedSide(side Side) Side {
	if side == SideUndefined {
		panic("Unappropriate side")
	} else {
		return base.MultiplyByConstantF(&side, -1)
	}
}

func GetAxisColor(axis base.Axis) MonomerType {
	return map[base.Axis]MonomerType{
		base.AxisX: MonomerTypeVynil,
		base.AxisY: MonomerTypeOContaining,
		base.AxisZ: MonomerTypeFwise,
	}[axis]
}

func (monType *MonomerType) ToLiteral() (string, error) {
	switch *monType {
	case MonomerTypeUndefined:
		return "", errors.New("there is no letter for Undefined monomer")
	case MonomerTypeUsual:
		return "C", nil
	case MonomerTypeVynil:
		return "O", nil
	case MonomerTypeOContaining:
		return "N", nil
	case MonomerTypeFwise:
		return "F", nil
	case MonomerTypeClwise:
		return "Cl", nil
	case MonomerTypeCrosslinked:
		return "H", nil
	case MonomerTypeS:
		return "S", nil
	default:
		return "", errors.New("there is no letter for this type")
	}
}

func GetNormalByAxis(axis base.Axis) []Side {
	switch axis {
	case base.AxisX:
		return []Side{
			SideForward,
			SideLeft,
			SideRight,
			SideUp,
			SideDown,
		}
	case base.AxisY:
		return []Side{
			SideForward,
			SideBackward,
			SideLeft,
			SideUp,
			SideDown,
		}
	case base.AxisZ:
		return []Side{
			SideForward,
			SideBackward,
			SideLeft,
			SideRight,
			SideUp,
		}
	default:
		return []Side{}
	}
}
