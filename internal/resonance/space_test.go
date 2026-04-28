package resonance

import "testing"

func TestNewFrequencySpace(t *testing.T) {
	fs := NewFrequencySpace()
	if _, ok := fs.meanings["meaning:truth"]; !ok {
		t.Error("missing meaning:truth")
	}
}

func TestConvertToFrequency(t *testing.T) {
	fs := NewFrequencySpace()
	f := fs.ConvertToFrequency("truth")
	if f.Harmonic != 2.0 {
		t.Errorf("expected Harmonic 2.0, got %f", f.Harmonic)
	}
	f2 := fs.ConvertToFrequency("狗") // CN dog
	if f2.Harmonic != 2.4 {
		t.Errorf("expected 2.4, got %f", f2.Harmonic)
	}
}
