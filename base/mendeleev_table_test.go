package base

import "testing"

func TestRecognizeElementValid(t *testing.T) {
	expected := Na
	got := RecognizeElement("Na")
	if expected != got {
		t.Fatalf("expected: %v, got %v", expected, got)
	}
}
func TestRecognizeElementInvalid(t *testing.T) {
	expected := Na
	got := RecognizeElement("Nas")
	if expected != got {
		t.Fatalf("expected: %v, got %v", expected, got)
	}
}
