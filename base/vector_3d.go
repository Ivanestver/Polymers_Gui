package base

import (
	"math"
)

type Vector3D struct {
	X int64
	Y int64
	Z int64
}

func InvalidVector() Vector3D {
	return Vector3D{
		X: -1,
		Y: -1,
		Z: -1,
	}
}

func (vector *Vector3D) IsInvalid() bool {
	return vector.X == -1 && vector.Y == -1 && vector.Z == -1
}

func VectorsAreEqual(left, right *Vector3D) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return left.X == right.X &&
		left.Y == right.Y &&
		left.Z == right.Z
}

func (vector *Vector3D) Add(other *Vector3D) {
	vector.X += other.X
	vector.Y += other.Y
	vector.Z += other.Z
}

func AddVec(left *Vector3D, right *Vector3D) *Vector3D {
	return &Vector3D{
		X: left.X + right.X,
		Y: left.Y + right.Y,
		Z: left.Z + right.Z,
	}
}

type Vector3DF struct {
	X, Y, Z float64
}

func (vector *Vector3DF) IsInvalid() bool {
	return math.IsNaN(vector.X) &&
		math.IsNaN(vector.Y) &&
		math.IsNaN(vector.Z)
}

func (vector *Vector3DF) AddF(other *Vector3DF) {
	vector.X += other.X
	vector.Y += other.Y
	vector.Z += other.Z
}

func (vector *Vector3DF) MultiplyByConstantF(constant float64) {
	vector.X *= constant
	vector.Y *= constant
	vector.Z *= constant
}

func Vector3D_To_Vector3DF(vector3D *Vector3D) Vector3DF {
	return Vector3DF{
		X: float64(vector3D.X),
		Y: float64(vector3D.Y),
		Z: float64(vector3D.Z),
	}
}

func VectorsAreEqualF(left, right *Vector3DF) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return left.X == right.X &&
		left.Y == right.Y &&
		left.Z == right.Z
}

func InvalidVectorF() Vector3DF {
	return Vector3DF{
		X: math.NaN(),
		Y: math.NaN(),
		Z: math.NaN(),
	}
}

func AddVecF(left *Vector3DF, right *Vector3DF) *Vector3DF {
	return &Vector3DF{
		X: left.X + right.X,
		Y: left.Y + right.Y,
		Z: left.Z + right.Z,
	}
}

func SubtractVecF(left *Vector3DF, right *Vector3DF) *Vector3DF {
	return &Vector3DF{
		X: left.X - right.X,
		Y: left.Y - right.Y,
		Z: left.Z - right.Z,
	}
}

func RevertVecF(vec *Vector3DF) *Vector3DF {
	return &Vector3DF{
		X: -1 * vec.X,
		Y: -1 * vec.Y,
		Z: -1 * vec.Z,
	}
}
