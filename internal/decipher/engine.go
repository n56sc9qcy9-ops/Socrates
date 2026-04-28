package decipher

import (
	"fmt"
	"sort"
	"strings"

	"socrates/internal/knowledge"
)

// Engine is the core resonance pattern engine.
type Engine struct {
	// Knowledge is the loaded knowledge base for data-driven lookups
	Knowledge *knowledge.Knowledge
}

// NewEngine creates a new decipher engine with embedded knowledge.
// Returns nil if knowledge cannot be loaded.
func NewEngine() (*Engine, error) {
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		return nil, err
	}
	return &Engine{Knowledge: kb}, nil
}

// NewEngineWithKnowledge creates a new engine with explicit knowledge.
func NewEngineWithKnowledge(kb *knowledge.Knowledge) *Engine {
	return &Engine{Knowledge: kb}
}

// Analyze runs all channels and produces a complete reading.
func (e *Engine) Analyze(input string) Reading {
	// Generate forms
	forms := GenerateForms(input)

	// Generate candidate forms for neighbor discovery
	candidates := GenerateCandidateForms(input)

	// Run all channels
	channels := RunAllChannels(forms, e.Knowledge)

	// Run fuzzy anchor matching
	anchors := GetAllAnchors(e.Knowledge)
	fuzzyMatches := FuzzyMatchEvidence(candidates, anchors)

	// Expand concepts through graph
	directConcepts := extractDirectConcepts(channels)
	conceptExpansions := ExpandConcepts(directConcepts, 0.4, e.Knowledge)

	// Analyze passage-level convergence using generic concept activation
	passageSignals := AnalyzePassageTokens(forms.Tokens, e.Knowledge)
	convergence := DetectConvergence(passageSignals, directConcepts, e.Knowledge)

	// Collect all signals
	allSignals := collectAllSignals(channels)

	// Build signal graph by target
	signalGraph := buildSignalGraph(allSignals)

	// Find convergence patterns
	converging, weakSignals := findPatterns(signalGraph)

	// Calculate scores
	scoreComponents := CalculateScoreComponents(candidates, fuzzyMatches, conceptExpansions, convergence, len(channels))
	baseScore := calculateOverallScore(channels, converging, allSignals)
	discoveryScore := CalculateFinalScore(scoreComponents)
	finalScore := (baseScore.Overall + discoveryScore) / 2.0

	// Generate concise reading
	reading := generateConciseReading(input, converging, weakSignals, convergence)

	// Generate warnings
	warnings := generateWarnings(input, converging, weakSignals)

	return Reading{
		Input:              input,
		Forms:              forms,
		Candidates:         candidates,
		FuzzyMatches:       fuzzyMatches,
		ConceptExpansions:  conceptExpansions,
		PassageSignals:     passageSignals,
		Convergence:        convergence,
		Channels:           channels,
		ConvergingPatterns: converging,
		WeakSignals:        weakSignals,
		Score: Score{
			Overall:    finalScore,
			ByChannel:  calculateChannelScores(channels),
			Components: scoreComponents,
		},
		ConciseReading: reading,
		Warnings:       warnings,
	}
}

// extractDirectConcepts extracts direct concept targets from signals.
func extractDirectConcepts(channels []ChannelResult) []string {
	conceptSet := make(map[string]bool)
	for _, ch := range channels {
		for _, sig := range ch.Signals {
			conceptSet[sig.Target] = true
		}
	}

	concepts := make([]string, 0, len(conceptSet))
	for c := range conceptSet {
		concepts = append(concepts, c)
	}
	return concepts
}

// calculateChannelScores computes per-channel scores.
func calculateChannelScores(channels []ChannelResult) map[string]float64 {
	scores := make(map[string]float64)
	for _, ch := range channels {
		scores[ch.Name] = ch.Score
	}
	return scores
}

// collectAllSignals gathers signals from all channels.
func collectAllSignals(channels []ChannelResult) []Signal {
	signals := make([]Signal, 0)
	for _, ch := range channels {
		signals = append(signals, ch.Signals...)
	}
	return signals
}

// SignalNode represents a node in the signal graph.
type SignalNode struct {
	Target       string
	Signals      []Signal
	Strength     float64
	ChannelCount int
	SignalCount  int
}

// buildSignalGraph groups signals by target concept.
func buildSignalGraph(signals []Signal) map[string]SignalNode {
	graph := make(map[string]SignalNode)

	for _, sig := range signals {
		// Normalize target for grouping
		target := normalizeTarget(sig.Target)

		if node, exists := graph[target]; exists {
			node.Signals = append(node.Signals, sig)
			node.Strength += sig.Weight
			node.ChannelCount = countUniqueChannels(node.Signals)
			node.SignalCount = len(node.Signals)
			graph[target] = node
		} else {
			graph[target] = SignalNode{
				Target:       target,
				Signals:      []Signal{sig},
				Strength:     sig.Weight,
				ChannelCount: 1,
				SignalCount:  1,
			}
		}
	}

	return graph
}

// normalizeTarget groups similar targets together.
func normalizeTarget(target string) string {
	// Map common targets to canonical forms
	aliases := map[string]string{
		"breath":        "breath",
		"wind":          "breath",
		"air":           "breath",
		"prana":         "breath",
		"ruach":         "breath",
		"qi":            "breath",
		"spirit":        "spirit",
		"soul":          "spirit",
		"mind":          "mind",
		"consciousness": "mind",
		"word":          "word",
		"speech":        "word",
		"logos":         "word",
		"utterance":     "word",
		"truth":         "truth",
		"real":          "truth",
		"light":         "light",
		"life":          "life",
		"vitality":      "life",
		"living":        "life",
		"one":           "one",
		"unity":         "one",
		"give":          "give",
		"gift":          "give",
		"path":          "path",
		"way":           "path",
		"source":        "source",
		"origin":        "source",
		"being":         "being",
		"existence":     "being",
	}

	if canon, ok := aliases[target]; ok {
		return canon
	}
	return target
}

// countUniqueChannels counts unique Channel values in signals.
// This is the correct way to count channels - by the Channel field, not Lens.
func countUniqueChannels(signals []Signal) int {
	channels := make(map[string]bool)
	for _, s := range signals {
		channels[s.Channel] = true
	}
	return len(channels)
}

// countUniqueLenses counts unique Lens values (for informational purposes only).
func countUniqueLenses(signals []Signal) int {
	lenses := make(map[string]bool)
	for _, s := range signals {
		lenses[s.Lens] = true
	}
	return len(lenses)
}

// findPatterns identifies converging patterns and weak signals.
func findPatterns(graph map[string]SignalNode) ([]Pattern, []Pattern) {
	converging := make([]Pattern, 0)
	weak := make([]Pattern, 0)

	// Sort targets by strength
	type targetScore struct {
		target string
		node   SignalNode
		score  float64
	}
	var ranked []targetScore

	for _, node := range graph {
		score := calculateTargetScore(node)
		ranked = append(ranked, targetScore{target: node.Target, node: node, score: score})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].score > ranked[j].score
	})

	// Top patterns are converging if they have multiple channels OR strong evidence
	for _, ts := range ranked {
		if ts.node.ChannelCount >= 2 || ts.score > 0.8 {
			converging = append(converging, Pattern{
				Name:     ts.target,
				Signals:  ts.node.Signals,
				Strength: ts.score,
			})
		} else if ts.node.ChannelCount == 1 && ts.score < 0.5 {
			weak = append(weak, Pattern{
				Name:     ts.target,
				Signals:  ts.node.Signals,
				Strength: ts.score,
			})
		} else if ts.score > 0.3 {
			// Middle ground - include but mark as weak if speculative only
			hasVerified := false
			for _, s := range ts.node.Signals {
				if s.Confidence == ConfidenceVerified {
					hasVerified = true
					break
				}
			}
			if !hasVerified {
				weak = append(weak, Pattern{
					Name:     ts.target,
					Signals:  ts.node.Signals,
					Strength: ts.score,
				})
			} else {
				converging = append(converging, Pattern{
					Name:     ts.target,
					Signals:  ts.node.Signals,
					Strength: ts.score,
				})
			}
		}
	}

	return converging, weak
}

// calculateTargetScore computes a score for a target based on multiple factors.
// Score reflects "strength of resonance pattern".
func calculateTargetScore(node SignalNode) float64 {
	// Factors:
	// - More signals = higher score
	// - More ACTUAL CHANNELS = higher score (not just lenses)
	// - Higher weights = higher score
	// - Verified signals = higher score
	// - Path simplicity (fewer parts for same coverage)

	signalCount := float64(node.SignalCount)

	// Channel factor - count actual channels, not lenses
	channelFactor := float64(node.ChannelCount) / 5.0 // Normalize to 0-1

	var totalWeight float64
	var verifiedWeight float64
	for _, s := range node.Signals {
		totalWeight += s.Weight
		if s.Confidence == ConfidenceVerified {
			verifiedWeight += s.Weight
		}
	}

	// Base score from strength
	score := node.Strength

	// Boost for multiple signals
	signalBoost := signalCount * 0.1
	if signalBoost > 0.3 {
		signalBoost = 0.3
	}
	score += signalBoost

	// Boost for MULTIPLE CHANNELS (not just different lenses)
	// This is the key fix - channels are actual data sources, lenses are just metadata
	score += channelFactor * 0.3

	// Boost for verified signals
	if totalWeight > 0 {
		verifiedRatio := verifiedWeight / totalWeight
		score += verifiedRatio * 0.2
	}

	// Normalize to 0-1
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateOverallScore computes the overall, per-channel, and per-pattern scores.
func calculateOverallScore(channels []ChannelResult, converging []Pattern, allSignals []Signal) Score {
	byChannel := make(map[string]float64)
	var total float64

	for _, ch := range channels {
		byChannel[ch.Name] = ch.Score
		total += ch.Score
	}

	// Weight by number of channels
	channelCount := len(channels)
	if channelCount == 0 {
		channelCount = 1
	}
	overall := total / float64(channelCount)

	// Boost for convergence across ACTUAL channels
	if len(converging) > 0 {
		var convStrength float64
		var totalConvChannels int
		for _, p := range converging {
			convStrength += p.Strength
			totalConvChannels += countUniqueChannels(p.Signals)
		}
		avgConv := convStrength / float64(len(converging))

		// Channel convergence bonus: multi-channel convergence scores higher
		// than single-channel signals, even verified ones
		channelBonus := float64(totalConvChannels) / float64(channelCount) * 0.2
		avgConv += channelBonus

		overall = (overall + avgConv) / 2.0
	}

	return Score{
		Overall:   overall,
		ByChannel: byChannel,
	}
}

// generateConciseReading creates a narrative reading from patterns.
// Uses ConvergenceResult for generic concept activation instead of semantic buckets.
func generateConciseReading(input string, converging []Pattern, weakSignals []Pattern, convergence ConvergenceResult) string {
	if len(converging) == 0 && len(weakSignals) == 0 && len(convergence.ActivatedConcepts) == 0 {
		return "Limited resonance patterns detected for this input."
	}

	var parts []string

	// Add passage-level activated concepts
	if len(convergence.TopConcepts) > 0 {
		var conceptNames []string
		for _, ac := range convergence.TopConcepts {
			if ac.Strength > 0.3 {
				conceptNames = append(conceptNames, ac.Concept)
			}
		}
		if len(conceptNames) > 0 {
			parts = append(parts, "Activated concepts: "+strings.Join(conceptNames, ", ")+".")
		}
	}

	// Add relation paths found
	if len(convergence.RelationPaths) > 0 {
		var pathSummaries []string
		for _, p := range convergence.RelationPaths {
			if len(pathSummaries) < 2 { // Limit to first 2 paths
				pathSummaries = append(pathSummaries, p)
			}
		}
		if len(pathSummaries) > 0 {
			parts = append(parts, "Relation paths: "+strings.Join(pathSummaries, "; ")+".")
		}
	}

	// Add co-activation score context
	if convergence.CoActivationScore > 0.3 {
		parts = append(parts, fmt.Sprintf("Co-activation strength: %.0f%%.", convergence.CoActivationScore*100))
	}

	// Identify main themes from converging patterns
	if len(converging) > 0 {
		var themes []string
		for _, p := range converging {
			if p.Strength > 0.5 {
				themes = append(themes, p.Name)
			}
		}
		if len(themes) > 0 {
			parts = append(parts, "This word shows resonance around "+strings.Join(themes, ", ")+".")
		}
	}

	// Add notes about weak signals
	if len(weakSignals) > 0 {
		var weakThemes []string
		for _, p := range weakSignals {
			weakThemes = append(weakThemes, p.Name)
		}
		if len(weakThemes) > 0 {
			parts = append(parts, "Secondary resonance with "+strings.Join(weakThemes[:min(3, len(weakThemes))], ", ")+" (speculative).")
		}
	}

	if len(parts) == 0 {
		return "Complex resonance patterns detected. Multiple interpretations possible."
	}

	return strings.Join(parts, " ")
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// generateWarnings generates appropriate warnings based on analysis.
func generateWarnings(input string, converging []Pattern, weakSignals []Pattern) []string {
	warnings := make([]string, 0)

	// Always include the disclaimer
	warnings = append(warnings, "This is a resonance reading, not verified etymology.")

	// Check for heavy reliance on speculative signals
	speculativeCount := 0
	totalSignals := 0
	for _, p := range converging {
		for _, s := range p.Signals {
			totalSignals++
			if s.Confidence == ConfidenceSpeculative {
				speculativeCount++
			}
		}
	}
	for _, p := range weakSignals {
		for _, s := range p.Signals {
			totalSignals++
			if s.Confidence == ConfidenceSpeculative {
				speculativeCount++
			}
		}
	}

	if totalSignals > 0 && float64(speculativeCount)/float64(totalSignals) > 0.5 {
		warnings = append(warnings, "Many signals are speculative. Interpretation may not reflect historical meaning.")
	}

	// Check for weak convergence
	if len(converging) == 0 && len(weakSignals) > 0 {
		warnings = append(warnings, "No strong convergence detected. This word may have limited resonance patterns in the system.")
	}

	// Check for unknown script
	script := DetectScript(input)
	if script == ScriptUnknown {
		warnings = append(warnings, "Unknown script detected. Analysis may be incomplete.")
	}

	return warnings
}
