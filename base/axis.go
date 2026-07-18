package base

type Axis uint8

const (
	AxisX     Axis = 0
	AxisY     Axis = 1
	AxisZ     Axis = 2
	AxisCount Axis = 3
)

var AxisXVec = Vector3DF{1.0, 0.0, 0.0}
var AxisYVec = Vector3DF{0.0, 1.0, 0.0}
var AxisZVec = Vector3DF{0.0, 0.0, 1.0}
var AxisXVecReversed = Vector3DF{-1.0, 0.0, 0.0}
var AxisYVecReversed = Vector3DF{0.0, -1.0, 0.0}
var AxisZVecReversed = Vector3DF{0.0, 0.0, -1.0}
var AxisToVector map[Axis]Vector3DF = map[Axis]Vector3DF{
	AxisX: AxisXVec,
	AxisY: AxisYVec,
	AxisZ: AxisZVec,
}

func (axis Axis) ToString() string {
	switch axis {
	case AxisX:
		return "X"
	case AxisY:
		return "Y"
	case AxisZ:
		return "Z"
	default:
		return "Undefined"
	}
}
