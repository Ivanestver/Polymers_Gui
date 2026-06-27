package base

import "testing"

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
