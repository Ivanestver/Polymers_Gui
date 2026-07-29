package base

import (
	"math"
)

type Quaternion struct {
	W float64
	V Vector3DF
}

func MakeQuaternionFromFloat(w, x, y, z float64) Quaternion {
	return Quaternion{
		W: x,
		V: Vector3DF{
			x,
			y,
			z,
		},
	}
}

func MakeQuaternionFromVector(w float64, v Vector3DF) Quaternion {
	return Quaternion{
		W: w,
		V: v,
	}
}

func (q Quaternion) Conjugate() Quaternion {
	return Quaternion{
		W: q.W,
		V: Vector3DF{
			-q.V[AxisX],
			-q.V[AxisY],
			-q.V[AxisZ],
		},
	}
}

func (q Quaternion) Len() float64 {
	prod := 0.0
	for i := range q.V {
		prod += q.V[i] * q.V[i]
	}
	return math.Sqrt(q.W*q.W + prod)
}

func MultiplyQuaternions(q1, q2 Quaternion) Quaternion {
	return Quaternion{
		W: (q1.W*q2.W - DotProduct(q1.V, q2.V)),
		V: AddVecF(AddVecF(MultiplyByConstantF(q2.V, q1.W), MultiplyByConstantF(q1.V, q2.W)), VectorProduct(q1.V, q2.V)),
	}
}
