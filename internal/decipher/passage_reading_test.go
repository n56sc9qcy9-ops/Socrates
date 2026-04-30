package decipher

import (
	"testing"
)

// PASSAGE FIELDS IN READING TESTS
// =============================================================================

func TestReading_ContainsPassageFields(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Multi-token passage
	reading := engine.Analyze("love truth light")

	// Reading MUST contain PassageFields
	if reading.PassageFields == nil {
		t.Fatal("Reading.PassageFields should not be nil for multi-token input")
	}

	if len(reading.PassageFields) == 0 {
		t.Error("Reading.PassageFields should not be empty for multi-token input")
	}
}

func TestReading_PassageFieldsNotNilForSingleToken(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Single token
	reading := engine.Analyze("love")

	// Reading should contain PassageFields (even for single token)
	if reading.PassageFields == nil {
		t.Fatal("Reading.PassageFields should not be nil")
	}
}

// =============================================================================
