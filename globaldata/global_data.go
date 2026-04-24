package globaldata

import (
	"polymers/base"
)

type RealFieldRestriction struct{ Lower, Higher float64 }

type SpaceDimention [base.AxisCount]RealFieldRestriction

// func (spaceDimention *SpaceDimention) PointInSpace(point *base.Vector3DF) bool {
// 	return 0 <= point.X && point.X < spaceDimention.X &&
// 		0 <= point.Y && point.Y < spaceDimention.Y &&
// 		0 <= point.Z && point.Z < spaceDimention.Z
// }

func (spaceDimention *SpaceDimention) PointInSpace(coords *base.Vector3DF) bool {
	return base.PointInSpace(coords,
		&base.Vector3DF{
			spaceDimention[base.AxisX].Lower,
			spaceDimention[base.AxisY].Lower,
			spaceDimention[base.AxisZ].Lower,
		},
		&base.Vector3DF{
			spaceDimention[base.AxisX].Higher,
			spaceDimention[base.AxisY].Higher,
			spaceDimention[base.AxisZ].Higher,
		})
	// return (spaceDimention[base.X_AXIS].Lower < coords.X || base.CompareFloat(spaceDimention[base.X_AXIS].Lower, coords.X)) &&
	// 	(coords.X < spaceDimention[base.X_AXIS].Higher || base.CompareFloat(spaceDimention[base.X_AXIS].Higher, coords.X)) &&
	// 	(spaceDimention[base.Y_AXIS].Lower < coords.Y || base.CompareFloat(spaceDimention[base.Y_AXIS].Lower, coords.Y)) &&
	// 	(coords.Y < spaceDimention[base.Y_AXIS].Higher || base.CompareFloat(spaceDimention[base.Y_AXIS].Lower, coords.Y)) &&
	// 	(spaceDimention[base.Z_AXIS].Lower < coords.Z || base.CompareFloat(spaceDimention[base.Z_AXIS].Lower, coords.Z)) &&
	// 	(coords.Z < spaceDimention[base.Z_AXIS].Higher || base.CompareFloat(spaceDimention[base.Z_AXIS].Lower, coords.Z))
}

func (spaceDimention *SpaceDimention) GetCenter() base.Vector3DF {
	return base.Vector3DF{
		(spaceDimention[base.AxisX].Lower + spaceDimention[base.AxisX].Higher) / 2,
		(spaceDimention[base.AxisY].Lower + spaceDimention[base.AxisY].Higher) / 2,
		(spaceDimention[base.AxisZ].Lower + spaceDimention[base.AxisZ].Higher) / 2,
	}
}

func (spaceDimention *SpaceDimention) GetV() float64 {
	return (spaceDimention[base.AxisX].Higher - spaceDimention[base.AxisX].Lower) *
		(spaceDimention[base.AxisY].Higher - spaceDimention[base.AxisY].Lower) *
		(spaceDimention[base.AxisZ].Higher - spaceDimention[base.AxisZ].Lower)
}

type GlobalData struct {
	SpaceDimention SpaceDimention
}

var globalData GlobalData

func SetSpaceDimention(spaceDimention SpaceDimention) {
	globalData.SpaceDimention = spaceDimention
}

func GetGlobalData() *GlobalData {
	return &globalData
}
