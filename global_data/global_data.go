package global_data

import "polymers/base"

type SpaceDimention base.Vector3DF

func (space *SpaceDimention) PointInSpace(point *base.Vector3DF) bool {
	return 0 <= point.X && point.X < space.X &&
		0 <= point.Y && point.Y < space.Y &&
		0 <= point.Z && point.Z < space.Z
}

type GlobalData struct {
	UpVector       base.Vector3D
	DownVector     base.Vector3D
	LeftVector     base.Vector3D
	RightVector    base.Vector3D
	ForwardVector  base.Vector3D
	BackwardVector base.Vector3D

	SpaceDimention SpaceDimention
}

var globalData GlobalData

func ConfigureGlobalData(spaceDimention SpaceDimention) {
	globalData.UpVector = base.Vector3D{X: 0, Y: 1, Z: 0}
	globalData.DownVector = base.Vector3D{X: 0, Y: -1, Z: 0}
	globalData.LeftVector = base.Vector3D{X: 1, Y: 0, Z: 0}
	globalData.RightVector = base.Vector3D{X: -1, Y: 0, Z: 0}
	globalData.ForwardVector = base.Vector3D{X: 0, Y: 0, Z: 1}
	globalData.BackwardVector = base.Vector3D{X: 0, Y: 0, Z: -1}

	globalData.SpaceDimention = spaceDimention
}

func GetGlobalData() *GlobalData {
	return &globalData
}
