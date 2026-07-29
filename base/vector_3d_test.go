package base

import (
	"math"
	"testing"
)

func TestVectorAdd(t *testing.T) {
	origin := Vector3DF{1.0, 0.0, 1.0}
	whatToAdd := Vector3DF{1.0, 0.0, -1.0}
	origin.AddF(whatToAdd)
	expected := Vector3DF{2.0, 0.0, 0.0}
	if !VectorsAreEqualF(origin, expected) {
		t.Fatalf("Expected: %v, got: %v", expected, origin)
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
	if !VectorsAreEqualF(expected, real) {
		t.Fatalf("Incorrect vector product. Expected: %v, got: %v", expected, real)
	}
}

func TestRotateVector(t *testing.T) {
	v := Vector3DF{1, 0, 0}
	angle := math.Pi / 6
	rotationVector := Vector3DF{0, 0, 1}
	expected := Vector3DF{0.866, 0.5, 0}
	real := RotateVector(v, angle, rotationVector)
	if !VectorsAreEqualF(expected, real) {
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

func TestCosAngle(t *testing.T) {
	v1 := Vector3DF{1, 0, 0}
	v2 := Vector3DF{5, 0, 0}
	expected := 1.0
	real := GetCos(v1, v2)
	if !CompareFloat(expected, real) {
		t.Fatalf("Неверный расчёт косинуса. Ожидаемый: %f, фактический: %f", expected, real)
	}
}
