package decipher

import (
	"sort"
	"strings"

	"socrates/internal/knowledge"
)

// ScriptWord represents a complete word in a specific script for decipher output.
type ScriptWord struct {
	Word       string
	Input      []rune
	Meanings   []string
	Confidence string
	Weight     float64
}

// knowledgeFragment represents a fragment from the knowledge base.
type knowledgeFragment struct {
	Form       string
	Concept    string
	Lens       string
	Confidence string
	Weight     float64
}

// knowledgePrimitive represents a primitive from the knowledge base.
type knowledgePrimitive struct {
	ID        string
	Name      string
	Aliases   []string
	Neighbors []string
}

// getKnowledgeFragments returns fragment mappings from the knowledge base.
// Returns nil if knowledge base is not available.
func getKnowledgeFragments(text string, kb *knowledge.Knowledge) []FragmentSeed {
	if kb == nil {
		return nil
	}

	forms := kb.GetFormsByText(strings.ToLower(text))
	if len(forms) == 0 {
		return nil
	}

	// Convert to FragmentSeed format
	seeds := make([]FragmentSeed, 0, len(forms))
	seenForms := make(map[string]bool)

	for _, f := range forms {
		// Group by form
		formKey := f.Form
		if !seenForms[formKey] {
			seenForms[formKey] = true
			seeds = append(seeds, FragmentSeed{
				Fragment:   f.Form,
				Lenses:     []FragmentLens{},
				Confidence: f.Confidence,
			})
		}

		// Find the seed and add the lens
		for i := range seeds {
			if seeds[i].Fragment == formKey {
				seeds[i].Lenses = append(seeds[i].Lenses, FragmentLens{
					Target:     f.Concept,
					Lens:       f.Lens,
					Confidence: f.Confidence,
					BaseWeight: f.Weight,
				})
				break
			}
		}
	}

	// Sort for deterministic output
	sort.Slice(seeds, func(i, j int) bool {
		return seeds[i].Fragment < seeds[j].Fragment
	})

	return seeds
}

// getKnowledgePrimitives returns primitive mappings from the knowledge base.
// Returns nil if knowledge base is not available or concept not found.
func getKnowledgePrimitives(text string, kb *knowledge.Knowledge) []Primitive {
	if kb == nil {
		return nil
	}

	textLower := strings.ToLower(text)

	// Check by alias
	concept, ok := kb.GetConceptByAlias(textLower)
	if ok {
		return []Primitive{
			{
				ID:        concept.ID,
				Name:      concept.Name,
				Aliases:   concept.Aliases,
				Neighbors: concept.Neighbors,
			},
		}
	}

	// Check by ID
	concept, ok = kb.GetConceptByID(textLower)
	if ok {
		return []Primitive{
			{
				ID:        concept.ID,
				Name:      concept.Name,
				Aliases:   concept.Aliases,
				Neighbors: concept.Neighbors,
			},
		}
	}

	return nil
}

// getKnowledgeScriptWords returns script words from the knowledge base.
// Returns nil if knowledge base is not available or script word not found.
func getKnowledgeScriptWords(script string, runes []rune, kb *knowledge.Knowledge) []knowledge.ScriptWord {
	if kb == nil {
		return nil
	}

	words := kb.GetScriptWordsByScript(script)
	if len(words) == 0 {
		return nil
	}

	// Filter by rune match
	result := make([]knowledge.ScriptWord, 0)
	for _, w := range words {
		// Check if runes match
		if len(w.Runes) != len(runes) {
			continue
		}
		match := true
		for i, r := range runes {
			if uint32(r) != w.Runes[i] {
				match = false
				break
			}
		}
		if match {
			result = append(result, w)
		}
	}

	return result
}

// knowledgeBasedFragmentLookup is the data-driven fragment lookup.
// Returns nil if knowledge base is not available or fragment not found.
func knowledgeBasedFragmentLookup(text string, kb *knowledge.Knowledge) []FragmentSeed {
	return getKnowledgeFragments(text, kb)
}

// knowledgeBasedPrimitiveLookup is the data-driven primitive lookup.
// Returns nil if knowledge base is not available or concept not found.
func knowledgeBasedPrimitiveLookup(text string, kb *knowledge.Knowledge) []Primitive {
	return getKnowledgePrimitives(text, kb)
}

// knowledgeBasedScriptWordLookup is the data-driven script word lookup.
// Returns nil if knowledge base is not available or script word not found.
func knowledgeBasedScriptWordLookup(script string, runes []rune, kb *knowledge.Knowledge) []knowledge.ScriptWord {
	return getKnowledgeScriptWords(script, runes, kb)
}
