package datatypes

import (
	"errors"
	"polymers/base"
)

type MoveDirection = int8
type Side int16
type MonomerType int8
type ConnectionType int8
type GlobulaViewType int

const (
	DirectionBackward MoveDirection = -1
	DirectionForward  MoveDirection = 1
	DirectionCount    MoveDirection = 2
)

const (
	SideUndefined         = -10
	SideForward           = Side(int8(base.AxisX+1) * DirectionForward)
	SideBackward          = Side(int8(base.AxisX+1) * DirectionBackward)
	SideLeft              = Side(int8(base.AxisY+1) * DirectionForward)
	SideRight             = Side(int8(base.AxisY+1) * DirectionBackward)
	SideUp                = Side(int8(base.AxisZ+1) * DirectionForward)
	SideDown              = Side(int8(base.AxisZ+1) * DirectionBackward)
	SideUpLeftForward     = 1000
	SideUpForward         = SideUpLeftForward + 1
	SideUpRightForward    = SideUpForward + 1
	SideLeftForward       = SideUpRightForward + 1
	SideRightForward      = SideLeftForward + 1
	SideDownLeftForward   = SideRightForward + 1
	SideDownForward       = SideDownLeftForward + 1
	SideDownRightForward  = SideDownForward + 1
	SideLeftUp            = SideDownRightForward + 1
	SideLeftDown          = SideLeftUp + 1
	SideRightUp           = SideLeftDown + 1
	SideRightDown         = SideRightUp + 1
	SideUpLeftBackward    = SideRightDown + 1
	SideUpBackward        = SideUpLeftBackward + 1
	SideUpRightBackward   = SideUpBackward + 1
	SideLeftBackward      = SideUpRightBackward + 1
	SideRightBackward     = SideLeftBackward + 1
	SideDownLeftBackward  = SideRightBackward + 1
	SideDownBackward      = SideDownLeftBackward + 1
	SideDownRightBackward = SideDownBackward + 1
)

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
	return Side(uint8(axis+1) * uint8(moveDirection))
}

func GetReversedSide(side Side) Side {
	switch side {
	case SideForward:
		return SideBackward
	case SideBackward:
		return SideForward
	case SideUp:
		return SideDown
	case SideDown:
		return SideUp
	case SideLeft:
		return SideRight
	case SideRight:
		return SideLeft
	case SideUpLeftForward:
		return SideDownRightBackward
	case SideUpForward:
		return SideDownBackward
	case SideUpRightForward:
		return SideDownLeftBackward
	case SideLeftForward:
		return SideRightBackward
	case SideRightForward:
		return SideLeftBackward
	case SideDownLeftForward:
		return SideUpRightBackward
	case SideDownForward:
		return SideUpBackward
	case SideDownRightForward:
		return SideUpLeftBackward
	case SideLeftUp:
		return SideRightDown
	case SideLeftDown:
		return SideRightUp
	case SideRightUp:
		return SideLeftDown
	case SideRightDown:
		return SideLeftUp
	case SideUpLeftBackward:
		return SideDownRightForward
	case SideUpBackward:
		return SideDownForward
	case SideUpRightBackward:
		return SideDownLeftForward
	case SideLeftBackward:
		return SideRightForward
	case SideRightBackward:
		return SideLeftForward
	case SideDownLeftBackward:
		return SideUpRightForward
	case SideDownBackward:
		return SideUpForward
	case SideDownRightBackward:
		return SideUpLeftForward
	default:
		panic("Unappropriate side")
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
