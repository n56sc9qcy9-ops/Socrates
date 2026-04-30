package decipher

import (
	"socrates/internal/knowledge"
)

// testKB returns a knowledge fixture for tests.
// Uses embedded knowledge if available, otherwise nil.
func testKB() *knowledge.Knowledge {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		return nil
	}
	return kb
}
