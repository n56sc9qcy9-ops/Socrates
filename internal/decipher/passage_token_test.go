package decipher

import (
	"testing"
)

// TOKENIZATION TESTS
// =============================================================================

func TestTokenizePassage_Basic(t *testing.T) {
	tokens := tokenizePassage("love truth")
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0] != "love" {
		t.Errorf("expected first token 'love', got '%s'", tokens[0])
	}
	if tokens[1] != "truth" {
		t.Errorf("expected second token 'truth', got '%s'", tokens[1])
	}
}

func TestTokenizePassage_Whitespace(t *testing.T) {
	tokens := tokenizePassage("  love   truth  ")
	if len(tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(tokens))
	}
}

func TestTokenizePassage_Newlines(t *testing.T) {
	tokens := tokenizePassage("love\ntruth\rlight")
	if len(tokens) != 3 {
		t.Errorf("expected 3 tokens, got %d", len(tokens))
	}
}

func TestTokenizePassage_SingleToken(t *testing.T) {
	tokens := tokenizePassage("love")
	if len(tokens) != 1 {
		t.Errorf("expected 1 token, got %d", len(tokens))
	}
}

func TestTokenizePassage_Empty(t *testing.T) {
	tokens := tokenizePassage("")
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

// =============================================================================
