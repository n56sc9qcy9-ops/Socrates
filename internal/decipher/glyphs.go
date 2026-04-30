package decipher

import "socrates/internal/knowledge"

// knowledgeBasedGlyphLookup returns glyph patterns from knowledge base.
// Returns empty slice if kb is nil or no patterns found.
func knowledgeBasedGlyphLookup(script ScriptType, kb *knowledge.Knowledge) []GlyphPatternSpec {
	if kb == nil {
		return nil
	}
	patterns := kb.GetGlyphPatternsByScript(string(script))
	result := make([]GlyphPatternSpec, 0, len(patterns))
	for _, p := range patterns {
		result = append(result, GlyphPatternSpec{
			Pattern:    p.Pattern,
			Rune:       p.Rune,
			Concept:    p.Concept,
			Lens:       p.Lens,
			Confidence: p.Confidence,
			Weight:     p.Weight,
		})
	}
	return result
}

// GlyphPatternSpec holds glyph pattern specification.
type GlyphPatternSpec struct {
	Pattern    string
	Rune       uint32
	Concept    string
	Lens       string
	Confidence string
	Weight     float64
}
