package decipher

import (
	"fmt"
	"strings"
)

// RenderReading returns a formatted string representation of a reading using default mode.
func RenderReading(r Reading) string {
	return RenderReadingWithOptions(r, DefaultRenderOptions())
}

// RenderReadingWithOptions returns a formatted string representation of a reading with options.
func RenderReadingWithOptions(r Reading, opts RenderOptions) string {
	if opts.Mode == RenderModeDebug {
		return renderDebug(r)
	}
	return renderDefault(r)
}

// renderDefault produces concise output for default mode.
// Shows: input, top evidence paths, score summary, warnings.
// Does NOT show: full candidates, full fuzzy matches, generated forms internals.
func renderDefault(r Reading) string {
	var sb strings.Builder

	// Input
	sb.WriteString("Input:\n")
	sb.WriteString("  " + r.Input + "\n\n")

	// Top Evidence Paths (convergence-based, not candidate dumps)
	sb.WriteString("Evidence:\n")

	// Show top concepts from convergence with their evidence paths
	if len(r.Convergence.TopConcepts) > 0 {
		for _, tc := range r.Convergence.TopConcepts {
			// Show concept with strength and source evidence
			if len(tc.Sources) > 0 {
				sb.WriteString(fmt.Sprintf("  - %s [%.0f%%, via: %s]\n",
					tc.Concept, tc.Strength*100, strings.Join(tc.Sources, ", ")))
			} else {
				sb.WriteString(fmt.Sprintf("  - %s [%.0f%%]\n",
					tc.Concept, tc.Strength*100))
			}
		}
	}

	// Show top passage fields if available (summary, evidence-first)
	if len(r.PassageFields) > 0 {
		topFields := r.PassageFields.TopFields(3)
		if len(topFields) > 0 {
			sb.WriteString("  Top fields:\n")
			for _, field := range topFields {
				// Show field with strength, tokens, and a hint of evidence
				if len(field.TokenSources) > 0 {
					sb.WriteString(fmt.Sprintf("    - %s [%.0f%%, via: %s]\n",
						field.Concept, field.Strength*100, strings.Join(field.TokenSources, ", ")))
				} else {
					sb.WriteString(fmt.Sprintf("    - %s [%.0f%%]\n",
						field.Concept, field.Strength*100))
				}
			}
		}
	}

	// Show relation paths if available (evidence of graph connections)
	if len(r.Convergence.RelationPaths) > 0 {
		sb.WriteString("  Paths:\n")
		for _, path := range r.Convergence.RelationPaths {
			sb.WriteString(fmt.Sprintf("    - %s\n", path))
		}
	}

	// Show key channel signals (limited, not full dumps)
	if len(r.Channels) > 0 {
		hasSignals := false
		for _, ch := range r.Channels {
			if len(ch.Signals) > 0 {
				hasSignals = true
				break
			}
		}
		if hasSignals {
			sb.WriteString("  Channels:\n")
			for _, ch := range r.Channels {
				if len(ch.Signals) > 0 {
					// Show top 3 signals per channel
					maxShow := 3
					if len(ch.Signals) < maxShow {
						maxShow = len(ch.Signals)
					}
					for i := 0; i < maxShow; i++ {
						sig := ch.Signals[i]
						sb.WriteString(fmt.Sprintf("    - [%s] %s -> %s\n",
							ch.Name, sig.Text, sig.Target))
					}
					if len(ch.Signals) > maxShow {
						sb.WriteString(fmt.Sprintf("    ... +%d more\n", len(ch.Signals)-maxShow))
					}
				}
			}
		}
	}

	// Show converging patterns (evidence of multi-channel support)
	if len(r.ConvergingPatterns) > 0 {
		sb.WriteString("  Converging:\n")
		for _, p := range r.ConvergingPatterns {
			sb.WriteString(fmt.Sprintf("    - %s [%.0f%%]\n", p.Name, p.Strength*100))
		}
	}

	// Show weak signals / warnings
	if len(r.WeakSignals) > 0 {
		sb.WriteString("  Weak signals:\n")
		for _, p := range r.WeakSignals {
			sb.WriteString(fmt.Sprintf("    - %s [%.0f%%]\n", p.Name, p.Strength*100))
		}
	}

	sb.WriteString("\n")

	// Score summary
	sb.WriteString("Score:\n")
	sb.WriteString(fmt.Sprintf("  - Overall: %.2f\n", r.Score.Overall))
	sb.WriteString("  - Components:\n")
	sb.WriteString(fmt.Sprintf("    - exact: %.2f\n", r.Score.Components.ExactMatchScore))
	sb.WriteString(fmt.Sprintf("    - fuzzy: %.2f\n", r.Score.Components.FuzzyMatchScore))
	sb.WriteString(fmt.Sprintf("    - graph: %.2f\n", r.Score.Components.GraphExpansionScore))
	sb.WriteString(fmt.Sprintf("    - convergence: %.2f\n", r.Score.Components.PassageConvergenceScore))
	if r.Score.Components.MultiMethodBonus > 0 || r.Score.Components.ChannelDiversityBonus > 0 {
		sb.WriteString(fmt.Sprintf("    - multi_method: +%.2f\n", r.Score.Components.MultiMethodBonus))
		sb.WriteString(fmt.Sprintf("    - channel_diversity: +%.2f\n", r.Score.Components.ChannelDiversityBonus))
	}

	// Bounded work summary (if any discarded)
	if r.DiscardedCandidates > 0 || r.DiscardedComparisons > 0 {
		sb.WriteString("  - Bounded work:\n")
		if r.DiscardedCandidates > 0 {
			sb.WriteString(fmt.Sprintf("    - candidates: %d discarded\n", r.DiscardedCandidates))
		}
		if r.DiscardedComparisons > 0 {
			sb.WriteString(fmt.Sprintf("    - comparisons: %d discarded\n", r.DiscardedComparisons))
		}
	}

	sb.WriteString("\n")

	// Reading conclusion
	sb.WriteString("Reading:\n")
	sb.WriteString("  " + r.ConciseReading + "\n\n")

	// Warnings
	if len(r.Warnings) > 0 {
		sb.WriteString("Warnings:\n")
		for _, w := range r.Warnings {
			sb.WriteString("  - " + w + "\n")
		}
	}

	return sb.String()
}

// renderDebug produces verbose output for debug mode.
// Shows: everything including candidates, fuzzy matches, generated forms, internal details.
func renderDebug(r Reading) string {
	var sb strings.Builder

	// Input
	sb.WriteString("Input:\n")
	sb.WriteString("  " + r.Input + "\n\n")

	// Generated Forms (full detail)
	sb.WriteString("Generated Forms:\n")
	sb.WriteString("  - normalized: " + r.Forms.Normalized + "\n")
	sb.WriteString("  - script: " + string(r.Forms.Script) + "\n")

	if len(r.Forms.Tokens) > 0 {
		sb.WriteString("  - tokens: " + strings.Join(r.Forms.Tokens, " | ") + "\n")
	}

	if len(r.Forms.Runes) > 0 {
		runeStrs := make([]string, len(r.Forms.Runes))
		for i, ru := range r.Forms.Runes {
			runeStrs[i] = fmt.Sprintf("%U", ru)
		}
		sb.WriteString("  - runes: " + strings.Join(runeStrs, " | ") + "\n")
	}

	if len(r.Forms.PhoneticKeys) > 0 {
		sb.WriteString("  - phonetic keys: " + strings.Join(r.Forms.PhoneticKeys, ", ") + "\n")
	}

	if len(r.Forms.Fragments) > 0 {
		sb.WriteString("  - fragments:\n")
		for _, f := range r.Forms.Fragments {
			sb.WriteString(fmt.Sprintf("    - %s [method: %s, conf: %.2f]\n",
				strings.Join(f.Parts, " | "), f.Method, f.Confidence))
		}
	}
	sb.WriteString("\n")

	// Phase A: Candidate Neighbors (full list)
	if len(r.Candidates) > 0 {
		sb.WriteString(fmt.Sprintf("Candidate Neighbors (%d total):\n", len(r.Candidates)))
		for _, c := range r.Candidates {
			sb.WriteString(fmt.Sprintf("  - %s [method: %s, distance: %.2f, %s]\n",
				c.Form, c.Method, c.Distance, c.Confidence))
		}
		sb.WriteString("\n")
	}

	// Bounded work info
	if r.DiscardedCandidates > 0 {
		sb.WriteString(fmt.Sprintf("Discarded Candidates: %d\n\n", r.DiscardedCandidates))
	}

	// Phase B: Fuzzy Matches (full list)
	sb.WriteString(fmt.Sprintf("Fuzzy Matches (%d total):\n", len(r.FuzzyMatches)))
	for _, m := range r.FuzzyMatches {
		sb.WriteString(fmt.Sprintf("  - %s -> %s [method: %s, distance: %.2f, weight: %.2f]\n",
			m.InputForm, m.AnchorForm, m.Method, m.Distance, m.Weight))
	}
	if len(r.FuzzyMatches) == 0 {
		sb.WriteString("  (none)\n")
	}
	sb.WriteString("\n")

	// Bounded comparisons info
	if r.DiscardedComparisons > 0 {
		sb.WriteString(fmt.Sprintf("Discarded Comparisons: %d\n\n", r.DiscardedComparisons))
	}

	// Phase D: Concept Expansions
	if len(r.ConceptExpansions) > 0 {
		sb.WriteString("Graph Expansions:\n")
		for from, exps := range r.ConceptExpansions {
			for _, e := range exps {
				sb.WriteString(fmt.Sprintf("  - %s -> %s [type: %s, weight: %.2f, %s]\n",
					from, e.To, e.Type, e.Weight, e.Confidence))
			}
		}
		sb.WriteString("\n")
	}

	// Phase E: Passage Signals
	if len(r.PassageSignals) > 0 {
		sb.WriteString("Passage Signals (Activated Concepts via Fuzzy Matching):\n")
		for _, ps := range r.PassageSignals {
			sb.WriteString(fmt.Sprintf("  - '%s' -> %s [weight: %.2f, matched: %s (%.2f), %s]\n",
				ps.Token, ps.Concept, ps.Weight, ps.MatchForm, ps.MatchScore, ps.Confidence))
		}
		sb.WriteString("\n")
	}

	// Phase E: Convergence via Generic Activation
	if len(r.Convergence.ActivatedConcepts) > 0 {
		sb.WriteString("Convergence (Generic Activation):\n")
		sb.WriteString(fmt.Sprintf("  Co-Activation Score: %.0f%%\n", r.Convergence.CoActivationScore*100))
		if len(r.Convergence.TopConcepts) > 0 {
			sb.WriteString("  Top Concepts:\n")
			for _, tc := range r.Convergence.TopConcepts {
				sources := strings.Join(tc.Sources, ", ")
				sb.WriteString(fmt.Sprintf("    - %s [strength: %.2f, sources: %s, %s]\n",
					tc.Concept, tc.Strength, sources, tc.Confidence))
			}
		}
		if len(r.Convergence.RelationPaths) > 0 {
			sb.WriteString("  Relation Paths:\n")
			for _, path := range r.Convergence.RelationPaths {
				sb.WriteString(fmt.Sprintf("    - %s\n", path))
			}
		}
		sb.WriteString("\n")
	}

	// Phase E: Passage Fields (full detail)
	if len(r.PassageFields) > 0 {
		sb.WriteString(fmt.Sprintf("Passage Fields (%d total):\n", len(r.PassageFields)))
		for _, field := range r.PassageFields {
			depthStr := "direct"
			if field.Depth > 0 {
				depthStr = fmt.Sprintf("depth-%d", field.Depth)
			}
			if len(field.TokenSources) > 0 {
				sb.WriteString(fmt.Sprintf("  - %s [strength: %.2f, %s, via: %s]\n",
					field.Concept, field.Strength, depthStr, strings.Join(field.TokenSources, ", ")))
			} else {
				sb.WriteString(fmt.Sprintf("  - %s [strength: %.2f, %s]\n",
					field.Concept, field.Strength, depthStr))
			}
			// Show evidence paths
			for _, ev := range field.EvidencePaths {
				sb.WriteString(fmt.Sprintf("    - evidence: %s [%s, %.2f]\n",
						ev.SourceToken, ev.SourceType, ev.Weight))
			}
			// Show relation paths
			for _, path := range field.RelationPaths {
				sb.WriteString(fmt.Sprintf("    - path: %s\n", path))
			}
		}
		sb.WriteString("\n")
	}

	// Resonance Channels (full detail)
	sb.WriteString("Resonance Channels:\n")
	for _, ch := range r.Channels {
		sb.WriteString(fmt.Sprintf("  - %s: [score: %.2f]\n", ch.Name, ch.Score))
		if len(ch.Signals) > 0 {
			for _, sig := range ch.Signals {
				sb.WriteString(fmt.Sprintf("    - %s -> %s [channel: %s, lens: %s, %s]\n",
					sig.Text, sig.Target, sig.Channel, sig.Lens, sig.Confidence))
			}
		}
	}
	sb.WriteString("\n")

	// Converging Patterns
	sb.WriteString("Direct Converging Patterns:\n")
	if len(r.ConvergingPatterns) == 0 {
		sb.WriteString("  (none)\n")
	} else {
		for _, p := range r.ConvergingPatterns {
			sb.WriteString(fmt.Sprintf("  - %s [strength: %.2f]\n", p.Name, p.Strength))
			for _, sig := range p.Signals {
				sb.WriteString(fmt.Sprintf("    - %s [channel: %s, lens: %s, %s]\n",
					sig.Text, sig.Channel, sig.Lens, sig.Confidence))
			}
		}
	}
	sb.WriteString("\n")

	// Weak Signals
	sb.WriteString("Weak Signals / Conflicts:\n")
	if len(r.WeakSignals) == 0 {
		sb.WriteString("  (none)\n")
	} else {
		for _, p := range r.WeakSignals {
			sb.WriteString(fmt.Sprintf("  - %s [strength: %.2f]\n", p.Name, p.Strength))
			for _, sig := range p.Signals {
				sb.WriteString(fmt.Sprintf("    - %s [channel: %s, lens: %s, %s]\n",
					sig.Text, sig.Channel, sig.Lens, sig.Confidence))
			}
		}
	}
	sb.WriteString("\n")

	// Scores (full detail)
	sb.WriteString("Resonance Score:\n")
	sb.WriteString(fmt.Sprintf("  - Overall: %.2f\n", r.Score.Overall))
	sb.WriteString("  - Components:\n")
	sb.WriteString(fmt.Sprintf("    - exact_match: %.2f\n", r.Score.Components.ExactMatchScore))
	sb.WriteString(fmt.Sprintf("    - fuzzy_match: %.2f\n", r.Score.Components.FuzzyMatchScore))
	sb.WriteString(fmt.Sprintf("    - graph_expansion: %.2f\n", r.Score.Components.GraphExpansionScore))
	sb.WriteString(fmt.Sprintf("    - passage_convergence: %.2f\n", r.Score.Components.PassageConvergenceScore))
	sb.WriteString(fmt.Sprintf("    - multi_method_bonus: %.2f\n", r.Score.Components.MultiMethodBonus))
	sb.WriteString(fmt.Sprintf("    - channel_diversity_bonus: %.2f\n", r.Score.Components.ChannelDiversityBonus))
	sb.WriteString("  - By channel:\n")
	for name, score := range r.Score.ByChannel {
		sb.WriteString(fmt.Sprintf("    - %s: %.2f\n", name, score))
	}
	sb.WriteString("\n")

	// Reading
	sb.WriteString("Reading:\n")
	sb.WriteString("  " + r.ConciseReading + "\n\n")

	// Warnings
	if len(r.Warnings) > 0 {
		sb.WriteString("Warnings:\n")
		for _, w := range r.Warnings {
			sb.WriteString("  - " + w + "\n")
		}
	}

	return sb.String()
}
