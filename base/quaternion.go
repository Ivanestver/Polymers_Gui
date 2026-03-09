package base

import "math"

type Quaternion struct {
	W float64
	V Vector3DF
}

func MakeQuaternionFromFloat(w, x, y, z float64) Quaternion {
	return Quaternion{
		W: x,
		V: Vector3DF{
			X: x,
			Y: y,
			Z: z,
		},
	}
}

func MakeQuaternionFromVector(w float64, v Vector3DF) Quaternion {
	return Quaternion{
		W: w,
		V: v,
	}
}

func (q *Quaternion) Conjugate() Quaternion {
	return Quaternion{
		W: q.W,
		V: Vector3DF{
			X: -q.V.X,
			Y: -q.V.Y,
			Z: -q.V.Z,
		},
	}
}

func (q *Quaternion) Len() float64 {
	return math.Sqrt(
		q.W*q.W +
			q.V.X*q.V.X +
			q.V.Y*q.V.Y +
			q.V.Z*q.V.Z)
}

func MultiplyQuaternions(q1, q2 Quaternion) Quaternion {
	return Quaternion{
		W: (q1.W*q2.W - DotProduct(q1.V, q2.V)),
		V: AddVecF(AddVecF(MultiplyByConstantF(&q2.V, q1.W), MultiplyByConstantF(&q1.V, q2.W)), VectorProduct(q1.V, q2.V)),
	}
}
