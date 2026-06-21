package datatypes

import (
	"polymers/base"
)

type MoveDirection = int8
type Side = base.Vector3DF
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

func GetAxisColor(axis base.Axis) base.MendeleevTableElement {
	return map[base.Axis]base.MendeleevTableElement{
		base.AxisX: base.N,
		base.AxisY: base.O,
		base.AxisZ: base.F,
	}[axis]
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
