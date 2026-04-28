package resonance

import (
	"testing"
)

func TestAreInResonance(t *testing.T) {
	// Existing Harmonic-only (tones 0 implicit, cosine=0 but update expects true for ratios)
	f1 := Frequency{Harmonic: 2.0}
	f2 := Frequency{Harmonic: 2.0} // ratio 1
	if !AreInResonance(f1, f2) {
		t.Error("should resonate")
	}
	f3 := Frequency{Harmonic: 1.0} // ratio 2
	if !AreInResonance(f1, f3) {
		t.Error("should resonate octave")
	}
	f4 := Frequency{Harmonic: 0.1}
	if AreInResonance(f1, f4) {
		t.Error("should not resonate")
	}

	// New tonal tests
	truth1 := Truth
	truth1.Tone2 = 0
	truth1.Tone3 = 0
	truth1.Harmonic = 1.0

	truth2 := Truth
	truth2.Tone2 = 0
	truth2.Tone3 = 0

	if !AreInResonance(truth1, truth2) {
		t.Error("pure Truth should resonate (cosine=1, ratio=2)")
	}

	love := Love
	love.Tone1 = 0
	love.Tone3 = 0

	if AreInResonance(truth1, love) {
		t.Error("orthogonal Truth/Love should not (cosine=0)")
	}

	// Blend: wisdom-like high Tone1/Tone3
	wise := Wisdom

	if !AreInResonance(truth1, wise) {
		t.Error("Truth/wise blend cosine~0.785 >0.7")
	}

	// Dissonant low
	lie := Lie

	if AreInResonance(truth1, lie) {
		t.Error("Truth/lie cosine~0.1 <0.7")
	}

	// Same blend different harmonic ok
	power1 := Power

	power2 := Power
	power2.Harmonic = 1.0

	if !AreInResonance(power1, power2) {
		t.Error("same power tones, ratio~2.2 ok")
	}
}

func TestIsValidResponse(t *testing.T) {
	resonant := []Frequency{
		{Harmonic: 2.0},
		{Harmonic: 1.0},
	}
	if !IsValidResponse(resonant) {
		t.Error("should be valid")
	}
	dissonant := []Frequency{
		{Harmonic: 2.0},
		{Harmonic: 1.0},
		{Harmonic: 0.1},
	}
	if IsValidResponse(dissonant) {
		t.Error("should be invalid")
	}
}
