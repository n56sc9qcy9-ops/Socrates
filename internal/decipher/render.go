package decipher

import (
	"fmt"
	"strings"
)

// RenderReading returns a formatted string representation of a reading.
func RenderReading(r Reading) string {
	var sb strings.Builder

	// Input
	sb.WriteString("Input:\n")
	sb.WriteString("  " + r.Input + "\n\n")

	// Generated Forms
	sb.WriteString("Generated Forms:\n")
	sb.WriteString("  - normalized: " + r.Forms.Normalized + "\n")
	sb.WriteString("  - script: " + string(r.Forms.Script) + "\n")

	if len(r.Forms.Tokens) > 0 {
		sb.WriteString("  - tokens: " + strings.Join(r.Forms.Tokens, " | ") + "\n")
	}

	if len(r.Forms.Runes) > 0 {
		runeStrs := make([]string, len(r.Forms.Runes))
		for i, r := range r.Forms.Runes {
			runeStrs[i] = fmt.Sprintf("%U", r)
		}
		sb.WriteString("  - runes: " + strings.Join(runeStrs, " | ") + "\n")
	}

	if len(r.Forms.PhoneticKeys) > 0 {
		sb.WriteString("  - phonetic keys: " + strings.Join(r.Forms.PhoneticKeys, ", ") + "\n")
	}

	if len(r.Forms.Fragments) > 0 {
		sb.WriteString("  - fragments: ")
		var fragStrs []string
		for _, f := range r.Forms.Fragments {
			if len(fragStrs) < 5 { // Limit display
				fragStrs = append(fragStrs, strings.Join(f.Parts, " | "))
			}
		}
		sb.WriteString(strings.Join(fragStrs, ", "))
		if len(r.Forms.Fragments) > 5 {
			sb.WriteString(" ... (and " + itoa(len(r.Forms.Fragments)-5) + " more)")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Phase A: Candidate Neighbors
	if len(r.Candidates) > 0 {
		sb.WriteString("Candidate Neighbors:\n")
		for i, c := range r.Candidates {
			if i < 15 { // Limit display
				sb.WriteString(fmt.Sprintf("  - %s [method: %s, distance: %.2f, %s]\n",
					c.Form, c.Method, c.Distance, c.Confidence))
			}
		}
		if len(r.Candidates) > 15 {
			sb.WriteString(fmt.Sprintf("  ... (and %d more)\n", len(r.Candidates)-15))
		}
		sb.WriteString("\n")
	}

	// Phase B: Fuzzy Matches (capped at 10 + weaker-match count)
	if len(r.FuzzyMatches) > 0 {
		sb.WriteString("Fuzzy Matches:\n")
		maxDisplay := 10
		for i, m := range r.FuzzyMatches {
			if i < maxDisplay {
				sb.WriteString(fmt.Sprintf("  - %s -> %s [method: %s, distance: %.2f, weight: %.2f]\n",
					m.InputForm, m.AnchorForm, m.Method, m.Distance, m.Weight))
			}
		}
		if len(r.FuzzyMatches) > maxDisplay {
			sb.WriteString(fmt.Sprintf("  ... and %d weaker matches\n", len(r.FuzzyMatches)-maxDisplay))
		}
		sb.WriteString("\n")
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

	// Resonance Channels
	sb.WriteString("Resonance Channels:\n")
	for _, ch := range r.Channels {
		sb.WriteString("  - " + ch.Name + ": [score: " + fmt.Sprintf("%.2f", ch.Score) + "]\n")
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

	// Scores
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
