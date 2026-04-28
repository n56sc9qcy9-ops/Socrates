package resonance

import "testing"

func TestResonanceEngine_FullQuery(t *testing.T) {
	space := NewFrequencySpace()
	engine := NewResonanceEngine(space)

	resp, valid := engine.FullQuery("truth")
	if !valid {
		t.Error("expected valid for truth")
	}
	if resp == "" {
		t.Error("expected non-empty response")
	}

	_, validLie := engine.FullQuery("lie")
	if validLie {
		t.Error("expected invalid for lie")
	}
}
