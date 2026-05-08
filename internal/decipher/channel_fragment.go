package decipher

import (
	"socrates/internal/knowledge"
)

// standaloneWeightMultiplier reduces the weight of signals from standalone function words.
// This prevents short prepositions/function words from creating overconfident structural
// fields through fragment/bigram matching (e.g., standalone "in" should not dominate
// with "inward"/"into" fields when used as a preposition).
const standaloneWeightMultiplier = 0.1

// applyStandaloneWeightReduction reduces weights for signals from standalone tokens.
// Returns a new slice with adjusted weights; original signals are unchanged.
func applyStandaloneWeightReduction(signals []Signal) []Signal {
	if len(signals) == 0 {
		return signals
	}

	// Check if any signal is from a standalone token
	hasStandalone := false
	for _, sig := range signals {
		if sig.IsStandaloneToken {
			hasStandalone = true
			break
		}
	}
	if !hasStandalone {
		return signals
	}

	// Apply reduction to signals from standalone tokens
	reduced := make([]Signal, len(signals))
	for i, sig := range signals {
		reduced[i] = sig
		if sig.IsStandaloneToken {
			reduced[i].Weight = sig.Weight * standaloneWeightMultiplier
			reduced[i].IsStandaloneToken = true // Keep flag for downstream
		}
	}
	return reduced
}

// runFragmentChannel analyzes fragment paths.
// For Latin script: includes exact whole-token matching.
// For non-Latin scripts: fragment matching is not applicable (see runScriptWordChannel).
// Uses knowledge base exclusively for data-driven lookups.
func runFragmentChannel(forms Forms, kb *knowledge.Knowledge) ChannelResult {
	signals := make([]Signal, 0)

	// Track primitives to avoid emitting duplicate primitive signals
	// (same primitive found via different fragment parts should only appear once)
	emittedPrimitives := make(map[string]bool)

	// Determine which tokens are standalone function words.
	// Signals from these tokens will be flagged as IsStandaloneToken.
	standaloneTokens := make(map[string]bool)
	for _, token := range forms.Tokens {
		if IsPrepositionOrFunctionWord(token) {
			standaloneTokens[token] = true
		}
	}

	// First: exact whole-token matching
	exactSignals := runWholeTokenMatching(forms, kb, standaloneTokens)
	signals = append(signals, exactSignals...)

	// For non-Latin scripts (but not Han, Hebrew, Devanagari which use ScriptWord)
	// check if the whole form matches a fragment seed (e.g., transliterated words)
	nonLatinSignals := runNonLatinWholeTokenMatching(forms, kb, standaloneTokens)
	signals = append(signals, nonLatinSignals...)

	// Second: fragment path analysis (knowledge base only).
	// Fragment matches from standalone tokens are flagged as IsStandaloneToken.
	// After collecting all signals, apply weight reduction for standalone token signals
	// to prevent function words from creating overconfident structural fields.
	for _, path := range forms.Fragments {
		for _, part := range path.Parts {
			// Determine if this part came from a standalone token.
			// A fragment part is standalone only if it equals a standalone token exactly
			// (e.g., 'in' in 'in my heart' → standalone; 'in' in 'inside' → not standalone).
			isStandalone := standaloneTokens[part]

			// Lookup fragment seeds from knowledge base
			seeds := knowledgeBasedFragmentLookup(part, kb)
			for _, seed := range seeds {
				for _, lens := range seed.Lenses {
					sig := Signal{
						Text:              "fragment '" + part + "' -> " + lens.Target,
						Target:            lens.Target,
						Channel:           "Fragment",
						Lens:              lens.Lens,
						Confidence:        lens.Confidence,
						Weight:            lens.Weight(path.Confidence),
						IsStandaloneToken: isStandalone,
					}
					signals = append(signals, sig)
				}
			}

			// Lookup primitives from knowledge base
			prims := knowledgeBasedPrimitiveLookup(part, kb)
			for _, prim := range prims {
				// Only emit one signal per primitive concept, regardless of fragment part
				if !emittedPrimitives[prim.ID] {
					emittedPrimitives[prim.ID] = true
					signals = append(signals, Signal{
						Text:              "primitive '" + prim.Name + "' matches " + prim.Name,
						Target:            prim.ID,
						Channel:           "Fragment",
						Lens:              "primitive",
						Confidence:        ConfidenceVerified,
						Weight:            0.6,
						IsStandaloneToken: isStandalone,
					})
				}
			}
		}
	}

	// Apply weight reduction for signals from standalone tokens
	// to prevent function words from creating overconfident structural fields.
	signals = applyStandaloneWeightReduction(signals)

	score := calculateChannelScore(signals)

	return ChannelResult{
		Name:    "Fragment",
		Signals: signals,
		Score:   score,
	}
}

// runWholeTokenMatching checks if the entire input matches a known word/seed.
// This provides exact whole-token matching for words like prana, ruach, logos, mantra.
// ONLY processes Latin script - non-Latin scripts use runScriptWordChannel.
func runWholeTokenMatching(forms Forms, kb *knowledge.Knowledge, standaloneTokens map[string]bool) []Signal {
	signals := make([]Signal, 0)

	// Only process Latin script - non-Latin scripts use runScriptWordChannel
	script := forms.Script
	if script != ScriptLatin {
		return signals
	}

	input := forms.Normalized
	isStandalone := standaloneTokens[input]

	// Check fragment seeds for exact whole-token match from knowledge base
	seeds := knowledgeBasedFragmentLookup(input, kb)
	for _, seed := range seeds {
		if seed.Fragment == input {
			for _, lens := range seed.Lenses {
				signals = append(signals, Signal{
					Text:              "exact whole-token match: " + input,
					Target:            lens.Target,
					Channel:           "Fragment",
					Lens:              lens.Lens,
					Confidence:        lens.Confidence,
					Weight:            lens.BaseWeight,
					IsDirect:          true, // Direct form evidence from knowledge base
					IsStandaloneToken: isStandalone,
				})
			}
		}
	}

	return signals
}

// runNonLatinWholeTokenMatching checks non-Latin script whole-token matches.
// This handles exact matches for non-Latin forms stored as fragments in the knowledge base.
// Unlike runWholeTokenMatching which only handles Latin script, this handles forms
// like Han 愛, Hebrew transliterated words, etc. that are stored in the fragment
// section of forms.yaml rather than as ScriptWords.
func runNonLatinWholeTokenMatching(forms Forms, kb *knowledge.Knowledge, standaloneTokens map[string]bool) []Signal {
	signals := make([]Signal, 0)

	// Only process non-Latin scripts
	script := forms.Script
	if script == ScriptLatin {
		return signals
	}

	input := forms.Normalized
	isStandalone := standaloneTokens[input]

	// Check fragment seeds for exact whole-token match from knowledge base
	seeds := knowledgeBasedFragmentLookup(input, kb)
	for _, seed := range seeds {
		if seed.Fragment == input {
			for _, lens := range seed.Lenses {
				signals = append(signals, Signal{
					Text:              "exact whole-token match: " + input,
					Target:            lens.Target,
					Channel:           "Fragment",
					Lens:              lens.Lens,
					Confidence:        lens.Confidence,
					Weight:            lens.BaseWeight,
					IsDirect:          true, // Direct form evidence from knowledge base
					IsStandaloneToken: isStandalone,
				})
			}
		}
	}

	return signals
}

// Weight adjusts weight by a multiplier (used by fragment confidence).
func (l FragmentLens) Weight(multiplier float64) float64 {
	weight := l.BaseWeight * multiplier
	if weight > 1.0 {
		return 1.0
	}
	return weight
}
