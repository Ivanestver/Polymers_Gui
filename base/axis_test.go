package base

import "testing"

func TestAxisToString(t *testing.T) {
	expectedAxisXStr := "X"
	realAxisXStr := AxisX.ToString()
	if realAxisXStr != expectedAxisXStr {
		t.Fatalf("Expected: %v, got: %v", expectedAxisXStr, realAxisXStr)
	}
	expectedAxisYStr := "Y"
	realAxisYStr := AxisY.ToString()
	if realAxisYStr != expectedAxisYStr {
		t.Fatalf("Expected: %v, got: %v", expectedAxisYStr, realAxisYStr)
	}
	expectedAxisZStr := "Z"
	realAxisZStr := AxisZ.ToString()
	if realAxisZStr != expectedAxisZStr {
		t.Fatalf("Expected: %v, got: %v", expectedAxisZStr, realAxisZStr)
	}
}
