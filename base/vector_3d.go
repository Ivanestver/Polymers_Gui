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

type Point3DF struct {
	X, Y, Z float64
}

type Vector3DF Point3DF

func (vector *Vector3DF) Len() float64 {
	return math.Sqrt(
		vector.X*vector.X +
			vector.Y*vector.Y +
			vector.Z*vector.Z)
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

func (vector *Vector3DF) Normalized() Vector3DF {
	return MultiplyByConstantF(vector, 1.0/vector.Len())
}

func Vector3DToVector3DF(vector3D *Vector3D) Vector3DF {
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
	const EPSILON float64 = 0.0001
	return math.Abs(left.X-right.X) < EPSILON &&
		math.Abs(left.Y-right.Y) < EPSILON &&
		math.Abs(left.Z-right.Z) < EPSILON
}

func InvalidVectorF() Vector3DF {
	return Vector3DF{
		X: math.NaN(),
		Y: math.NaN(),
		Z: math.NaN(),
	}
}

func IndentityVectorF() Vector3DF {
	return Vector3DF{
		X: 0.0,
		Y: 0.0,
		Z: 0.0,
	}
}

func AddVecF(left Vector3DF, right Vector3DF) Vector3DF {
	return Vector3DF{
		X: left.X + right.X,
		Y: left.Y + right.Y,
		Z: left.Z + right.Z,
	}
}

func SubtractVecF(left Vector3DF, right Vector3DF) Vector3DF {
	return Vector3DF{
		X: left.X - right.X,
		Y: left.Y - right.Y,
		Z: left.Z - right.Z,
	}
}

func MultiplyByConstantF(left *Vector3DF, c float64) Vector3DF {
	return Vector3DF{
		X: left.X * c,
		Y: left.Y * c,
		Z: left.Z * c,
	}
}

func RevertVecF(vec *Vector3DF) *Vector3DF {
	return &Vector3DF{
		X: -1 * vec.X,
		Y: -1 * vec.Y,
		Z: -1 * vec.Z,
	}
}

func MakeVectorF(from, to *Point3DF) Vector3DF {
	return Vector3DF{
		X: to.X - from.X,
		Y: to.Y - from.Y,
		Z: to.Z - from.Z,
	}
}

func GetAngle(v1, v2 Vector3DF) float64 {
	return math.Acos(DotProduct(v1, v2) / (v1.Len() * v2.Len()))
}

func DotProduct(v1, v2 Vector3DF) float64 {
	return v1.X*v2.X +
		v1.Y*v2.Y +
		v1.Z*v2.Z
}

func VectorProduct(v1, v2 Vector3DF) Vector3DF {
	return Vector3DF{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: -(v1.X*v2.Z - v1.Z*v2.X),
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

func RotateVector(originVector Vector3DF, angle float64, rotationVector Vector3DF) Vector3DF {
	cosAngle := math.Cos(angle / 2)
	sinAngle := math.Sin(angle / 2)
	r := rotationVector.Normalized()
	q := MakeQuaternionFromVector(cosAngle, MultiplyByConstantF(&r, sinAngle))
	qConjugate := q.Conjugate()
	v := MakeQuaternionFromVector(0, originVector)
	rotatedVector := MultiplyQuaternions(MultiplyQuaternions(q, v), qConjugate)
	return rotatedVector.V
}
