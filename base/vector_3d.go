package base

import (
	"math"
)

type Vector3D [AxisCount]int64

func InvalidVector() Vector3D {
	return Vector3D{
		-1,
		-1,
		-1,
	}
}

func (vector *Vector3D) IsInvalid() bool {
	return vector[AxisX] == -1 && vector[AxisY] == -1 && vector[AxisZ] == -1
}

func VectorsAreEqual(left, right *Vector3D) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	return left[AxisX] == right[AxisX] &&
		left[AxisY] == right[AxisY] &&
		left[AxisZ] == right[AxisZ]
}

func (vector *Vector3D) Add(other *Vector3D) {
	vector[AxisX] += other[AxisX]
	vector[AxisY] += other[AxisY]
	vector[AxisZ] += other[AxisZ]
}

func AddVec(left *Vector3D, right *Vector3D) *Vector3D {
	return &Vector3D{
		left[AxisX] + right[AxisX],
		left[AxisY] + right[AxisY],
		left[AxisZ] + right[AxisZ],
	}
}

type Point3DF struct {
	X, Y, Z float64
}

type Vector3DF [AxisCount]float64

func (vector *Vector3DF) Len() float64 {
	return math.Sqrt(
		vector[AxisX]*vector[AxisX] +
			vector[AxisY]*vector[AxisY] +
			vector[AxisZ]*vector[AxisZ])
}

func (vector *Vector3DF) IsInvalid() bool {
	return math.IsNaN(vector[AxisX]) &&
		math.IsNaN(vector[AxisY]) &&
		math.IsNaN(vector[AxisZ])
}

func (vector *Vector3DF) AddF(other *Vector3DF) {
	vector[AxisX] += other[AxisX]
	vector[AxisY] += other[AxisY]
	vector[AxisZ] += other[AxisZ]
}

func (vector *Vector3DF) MultiplyByConstantF(constant float64) {
	vector[AxisX] *= constant
	vector[AxisY] *= constant
	vector[AxisZ] *= constant
}

func (vector *Vector3DF) Normalized() Vector3DF {
	return MultiplyByConstantF(vector, 1.0/vector.Len())
}

func Vector3DToVector3DF(vector3D *Vector3D) Vector3DF {
	return Vector3DF{
		float64(vector3D[AxisX]),
		float64(vector3D[AxisY]),
		float64(vector3D[AxisZ]),
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
	return math.Abs(left[AxisX]-right[AxisX]) < EPSILON &&
		math.Abs(left[AxisY]-right[AxisY]) < EPSILON &&
		math.Abs(left[AxisZ]-right[AxisZ]) < EPSILON
}

func InvalidVectorF() Vector3DF {
	return Vector3DF{
		math.NaN(),
		math.NaN(),
		math.NaN(),
	}
}

func IdentityVectorF() Vector3DF {
	return Vector3DF{
		0.0,
		0.0,
		0.0,
	}
}

func AddVecF(vectors ...Vector3DF) Vector3DF {
	res := Vector3DF{0.0, 0.0, 0.0}
	for _, v := range vectors {
		res.AddF(&v)
	}
	return res
}

func SubtractVecF(from Vector3DF, what Vector3DF) Vector3DF {
	return Vector3DF{
		from[AxisX] - what[AxisX],
		from[AxisY] - what[AxisY],
		from[AxisZ] - what[AxisZ],
	}
}

func MultiplyByConstantF(v *Vector3DF, c float64) Vector3DF {
	return Vector3DF{
		v[AxisX] * c,
		v[AxisY] * c,
		v[AxisZ] * c,
	}
}

func RevertVecF(vec *Vector3DF) *Vector3DF {
	return &Vector3DF{
		-1 * vec[AxisX],
		-1 * vec[AxisY],
		-1 * vec[AxisZ],
	}
}

func MakeVectorF(from, to *Point3DF) Vector3DF {
	return Vector3DF{
		to.X - from.X,
		to.Y - from.Y,
		to.Z - from.Z,
	}
}

func RadIntoGrad(angleInRad float64) float64 {
	return angleInRad * 57.275
}

func GradIntoRad(angleInGrad float64) float64 {
	return angleInGrad / 57.275
}

func GetCos(v1, v2 Vector3DF) float64 {
	return DotProduct(v1, v2) / (v1.Len() * v2.Len())
}

func GetAngleInRad(v1, v2 Vector3DF) float64 {
	return math.Acos(GetCos(v1, v2))
}

func GetAngleInGrad(v1, v2 Vector3DF) float64 {
	return RadIntoGrad(GetAngleInRad(v1, v2))
}

func DotProduct(v1, v2 Vector3DF) float64 {
	prod := 0.0
	for i := range v1 {
		prod += v1[i] * v2[i]
	}
	return prod
}

func VectorProduct(v1, v2 Vector3DF) Vector3DF {
	return Vector3DF{
		v1[AxisY]*v2[AxisZ] - v1[AxisZ]*v2[AxisY],
		-(v1[AxisX]*v2[AxisZ] - v1[AxisZ]*v2[AxisX]),
		v1[AxisX]*v2[AxisY] - v1[AxisY]*v2[AxisX],
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
