package base

type Axis uint8

const (
	X_AXIS     Axis = 0
	Y_AXIS     Axis = 1
	Z_AXIS     Axis = 2
	AXIS_COUNT Axis = 3
)

func (axis Axis) ToString() string {
	switch axis {
	case X_AXIS:
		return "X"
	case Y_AXIS:
		return "Y"
	case Z_AXIS:
		return "Z"
	default:
		return "Undefined"
	}
}
