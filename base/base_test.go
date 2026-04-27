package base

import (
	"math"
	"testing"
)

func TestQuaternionMultiplication(t *testing.T) {
	q1 := MakeQuaternionFromFloat(1, 1, 0, 0)
	q2 := MakeQuaternionFromFloat(0, 0, 1, 0)
	qExpected := MakeQuaternionFromFloat(0, 0, 1, 1)
	qReal := MultiplyQuaternions(q1, q2)
	if qExpected.W != qReal.W ||
		!VectorsAreEqualF(&qExpected.V, &qReal.V) {
		t.Fatalf("Incorrect multiplication. Expected: %v, got: %v", qExpected, qReal)
	}
}

func TestDotProduct(t *testing.T) {
	v1 := Vector3DF{1, -2, 4}
	v2 := Vector3DF{3, 0, 5}
	expected := 23.0
	real := DotProduct(v1, v2)
	if math.Abs(float64(expected)-real) > 0.0001 {
		t.Fatalf("Incorrect dot product. Expected: %f, got: %f", expected, real)
	}
}

func TestVectorProduct(t *testing.T) {
	v1 := Vector3DF{1, 2, 3}
	v2 := Vector3DF{4, 5, 6}
	expected := Vector3DF{-3, 6, -3}
	real := VectorProduct(v1, v2)
	if !VectorsAreEqualF(&expected, &real) {
		t.Fatalf("Incorrect vector product. Expected: %v, got: %v", expected, real)
	}
}

func TestRotateVector(t *testing.T) {
	v := Vector3DF{1, 0, 0}
	angle := math.Pi / 6
	rotationVector := Vector3DF{0, 0, 1}
	expected := Vector3DF{0.866, 0.5, 0}
	real := RotateVector(v, angle, rotationVector)
	if !VectorsAreEqualF(&expected, &real) {
		t.Fatalf("Incorrect rotated vector. Expected: %v, got: %v", expected, real)
	}
}

func TestRotateVectorCorrectness(t *testing.T) {
	v := Vector3DF{1, 0, 0}
	angle := math.Pi / 6
	rotationVector := Vector3DF{0, 0, 1}
	real := RotateVector(v, angle, rotationVector)
	realAngle := GetAngleInRad(v, real)
	if math.Abs(realAngle-angle) > 0.1 {
		t.Fatalf("Incorrect rotated vector correctness. Expected: %f, got: %f", angle, realAngle)
	}
}
