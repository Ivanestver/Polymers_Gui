package datatypes

import "errors"

type Axis = uint8
type MoveDirection = int8
type Side int16
type MonomerType int8
type ConnectionType int8

const (
	X_AXIS     Axis = 0
	Y_AXIS     Axis = 1
	Z_AXIS     Axis = 2
	AXIS_COUNT Axis = 3
)

const (
	DIRECTION_BACKWARD MoveDirection = -1
	DIRECTION_FORWARD  MoveDirection = 1
	DIRECTION_COUNT    MoveDirection = 2
)

const (
	SIDE_Undefined         = -10
	SIDE_Forward           = Side(int8(X_AXIS+1) * DIRECTION_FORWARD)
	SIDE_Backward          = Side(int8(X_AXIS+1) * DIRECTION_BACKWARD)
	SIDE_Left              = Side(int8(Y_AXIS+1) * DIRECTION_FORWARD)
	SIDE_Right             = Side(int8(Y_AXIS+1) * DIRECTION_BACKWARD)
	SIDE_Up                = Side(int8(Z_AXIS+1) * DIRECTION_FORWARD)
	SIDE_Down              = Side(int8(Z_AXIS+1) * DIRECTION_BACKWARD)
	SIDE_UpLeftForward     = 1000
	SIDE_UpForward         = SIDE_UpLeftForward + 1
	SIDE_UpRightForward    = SIDE_UpForward + 1
	SIDE_LeftForward       = SIDE_UpRightForward + 1
	SIDE_RightForward      = SIDE_LeftForward + 1
	SIDE_DownLeftForward   = SIDE_RightForward + 1
	SIDE_DownForward       = SIDE_DownLeftForward + 1
	SIDE_DownRightForward  = SIDE_DownForward + 1
	SIDE_LeftUp            = SIDE_DownRightForward + 1
	SIDE_LeftDown          = SIDE_LeftUp + 1
	SIDE_RightUp           = SIDE_LeftDown + 1
	SIDE_RightDown         = SIDE_RightUp + 1
	SIDE_UpLeftBackward    = SIDE_RightDown + 1
	SIDE_UpBackward        = SIDE_UpLeftBackward + 1
	SIDE_UpRightBackward   = SIDE_UpBackward + 1
	SIDE_LeftBackward      = SIDE_UpRightBackward + 1
	SIDE_RightBackward     = SIDE_LeftBackward + 1
	SIDE_DownLeftBackward  = SIDE_RightBackward + 1
	SIDE_DownBackward      = SIDE_DownLeftBackward + 1
	SIDE_DownRightBackward = SIDE_DownBackward + 1
)

const (
	MONOMER_TYPE_UNDEFINED MonomerType = -1
	MONOMER_TYPE_USUAL     MonomerType = 0
	MONOMER_TYPE_OWISE     MonomerType = 1
	MONOMER_TYPE_NWISE     MonomerType = 2
	MONOMER_TYPE_FWISE     MonomerType = 3
	MONOMER_TYPE_CLWISE    MonomerType = 4
	MONOMER_TYPE_HWISE     MonomerType = 5
)

const (
	CONNECTION_TYPE_UNDEFINED = 0
	CONNECTION_TYPE_ONE       = 1
	CONNECTION_TYPE_TWO       = 2
	CONNECTION_TYPE_THREE     = 3
)

func GetMovementSides() []Side {
	return []Side{SIDE_Left, SIDE_Right, SIDE_Forward, SIDE_Backward, SIDE_Up, SIDE_Down}
}

func GetSurfaceDiagonalSides() []Side {
	return []Side{SIDE_UpForward, SIDE_UpBackward, SIDE_DownForward, SIDE_DownBackward, SIDE_LeftUp, SIDE_LeftDown, SIDE_RightUp, SIDE_RightDown, SIDE_LeftForward, SIDE_LeftBackward, SIDE_RightForward, SIDE_RightBackward}
}

func GetCubeDiagonalSides() []Side {
	return []Side{SIDE_UpLeftForward, SIDE_UpLeftBackward, SIDE_UpRightForward, SIDE_UpRightBackward, SIDE_DownLeftForward, SIDE_DownLeftBackward, SIDE_DownRightForward, SIDE_DownRightBackward}
}

func GetAdditionalSides() []Side {
	return append(GetSurfaceDiagonalSides(), GetCubeDiagonalSides()...)
}

func GetAllSides() []Side {
	return append(GetMovementSides(), GetAdditionalSides()...)
}

func GetSide(axis Axis, moveDirection MoveDirection) Side {
	return Side((axis + 1) * uint8(moveDirection))
}

func GetReversedSide(side Side) Side {
	switch side {
	case SIDE_Forward:
		return SIDE_Backward
	case SIDE_Backward:
		return SIDE_Forward
	case SIDE_Up:
		return SIDE_Down
	case SIDE_Down:
		return SIDE_Up
	case SIDE_Left:
		return SIDE_Right
	case SIDE_Right:
		return SIDE_Left
	case SIDE_UpLeftForward:
		return SIDE_DownRightBackward
	case SIDE_UpForward:
		return SIDE_DownBackward
	case SIDE_UpRightForward:
		return SIDE_DownLeftBackward
	case SIDE_LeftForward:
		return SIDE_RightBackward
	case SIDE_RightForward:
		return SIDE_LeftBackward
	case SIDE_DownLeftForward:
		return SIDE_UpRightBackward
	case SIDE_DownForward:
		return SIDE_UpBackward
	case SIDE_DownRightForward:
		return SIDE_UpLeftBackward
	case SIDE_LeftUp:
		return SIDE_RightDown
	case SIDE_LeftDown:
		return SIDE_RightUp
	case SIDE_RightUp:
		return SIDE_LeftDown
	case SIDE_RightDown:
		return SIDE_LeftUp
	case SIDE_UpLeftBackward:
		return SIDE_DownRightForward
	case SIDE_UpBackward:
		return SIDE_DownForward
	case SIDE_UpRightBackward:
		return SIDE_DownLeftForward
	case SIDE_LeftBackward:
		return SIDE_RightForward
	case SIDE_RightBackward:
		return SIDE_LeftForward
	case SIDE_DownLeftBackward:
		return SIDE_UpRightForward
	case SIDE_DownBackward:
		return SIDE_UpForward
	case SIDE_DownRightBackward:
		return SIDE_UpLeftForward
	default:
		panic("Unappropriate side")
	}
}

func GetAxisColor(axis Axis) MonomerType {
	return map[Axis]MonomerType{
		X_AXIS: MONOMER_TYPE_OWISE,
		Y_AXIS: MONOMER_TYPE_NWISE,
		Z_AXIS: MONOMER_TYPE_FWISE,
	}[axis]
}

func (monType *MonomerType) ToLiteral() (string, error) {
	switch *monType {
	case MONOMER_TYPE_UNDEFINED:
		return "", errors.New("There is no letter for Undefined monomer")
	case MONOMER_TYPE_USUAL:
		return "C", nil
	case MONOMER_TYPE_OWISE:
		return "O", nil
	case MONOMER_TYPE_NWISE:
		return "N", nil
	case MONOMER_TYPE_FWISE:
		return "F", nil
	case MONOMER_TYPE_CLWISE:
		return "Cl", nil
	case MONOMER_TYPE_HWISE:
		return "H", nil
	default:
		return "", errors.New("There is no letter for this type")
	}
}
