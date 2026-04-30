package decipher

import (
	"testing"
)

// CONFIDENCE MERGE BUG FIX TESTS
// =============================================================================

// TestConfidenceMergeBugFix verifies the bug fix for confidence downgrading.
// The previous isHigherConfidence() had inverted semantics that could downgrade
// verified confidence when merging with plausible/speculative evidence.
func TestConfidenceMergeBugFix_VerifiedNotDowngraded(t *testing.T) {
	fields1 := PassageFields{
		{Concept: "truth", Strength: 0.5, Confidence: ConfidenceVerified},
	}
	fields2 := PassageFields{
		{Concept: "truth", Strength: 0.3, Confidence: ConfidencePlausible},
	}

	merged := fields1.Merge(fields2)

	if len(merged) != 1 {
		t.Fatalf("expected 1 field, got %d", len(merged))
	}

	// Verified should NOT be downgraded to plausible
	if merged[0].Confidence != ConfidenceVerified {
		t.Errorf("verified confidence should not be downgraded: got %s, want %s",
			merged[0].Confidence, ConfidenceVerified)
	}
}

func TestConfidenceMergeBugFix_VerifiedWinsOverSpeculative(t *testing.T) {
	fields1 := PassageFields{
		{Concept: "truth", Strength: 0.5, Confidence: ConfidenceVerified},
	}
	fields2 := PassageFields{
		{Concept: "truth", Strength: 0.3, Confidence: ConfidenceSpeculative},
	}

	merged := fields1.Merge(fields2)

	// Verified should win over speculative
	if merged[0].Confidence != ConfidenceVerified {
		t.Errorf("verified confidence should win: got %s, want %s",
			merged[0].Confidence, ConfidenceVerified)
	}
}

func TestConfidenceMergeBugFix_PlausibleWinsOverSpeculative(t *testing.T) {
	fields1 := PassageFields{
		{Concept: "truth", Strength: 0.5, Confidence: ConfidencePlausible},
	}
	fields2 := PassageFields{
		{Concept: "truth", Strength: 0.3, Confidence: ConfidenceSpeculative},
	}

	merged := fields1.Merge(fields2)

	// Plausible should win over speculative
	if merged[0].Confidence != ConfidencePlausible {
		t.Errorf("plausible confidence should win: got %s, want %s",
			merged[0].Confidence, ConfidencePlausible)
	}
}

func TestKeepHigherConfidence_Direct(t *testing.T) {
	tests := []struct {
		a, b, want string
	}{
		{ConfidenceVerified, ConfidenceVerified, ConfidenceVerified},
		{ConfidenceVerified, ConfidencePlausible, ConfidenceVerified},
		{ConfidenceVerified, ConfidenceSpeculative, ConfidenceVerified},
		{ConfidencePlausible, ConfidenceVerified, ConfidenceVerified},
		{ConfidencePlausible, ConfidencePlausible, ConfidencePlausible},
		{ConfidencePlausible, ConfidenceSpeculative, ConfidencePlausible},
		{ConfidenceSpeculative, ConfidenceVerified, ConfidenceVerified},
		{ConfidenceSpeculative, ConfidencePlausible, ConfidencePlausible},
		{ConfidenceSpeculative, ConfidenceSpeculative, ConfidenceSpeculative},
	}

	for _, tc := range tests {
		got := keepHigherConfidence(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("keepHigherConfidence(%s, %s): got %s, want %s",
				tc.a, tc.b, got, tc.want)
		}
	}
}

// =============================================================================
