package decipher

import "socrates/internal/knowledge"

// Confidence levels for signals.
const (
	ConfidenceVerified    = "verified"
	ConfidencePlausible   = "plausible"
	ConfidenceSpeculative = "speculative"
)

// Primitive represents a core symbolic/conceptual anchor.
type Primitive struct {
	ID        string
	Name      string
	Aliases   []string
	Neighbors []string
}

// FragmentSeed represents a reusable fragment with possible lenses.
type FragmentSeed struct {
	Fragment   string
	Lenses     []FragmentLens
	Confidence string
}

// FragmentLens represents a possible interpretation for a fragment.
type FragmentLens struct {
	Target     string
	Lens       string
	Confidence string
	BaseWeight float64
}

// ScriptType represents detected script categories.
type ScriptType string

const (
	ScriptLatin      ScriptType = "latin"
	ScriptHebrew     ScriptType = "hebrew"
	ScriptDevanagari ScriptType = "devanagari"
	ScriptHan        ScriptType = "han"
	ScriptArabic     ScriptType = "arabic"
	ScriptGreek      ScriptType = "greek"
	ScriptUnknown    ScriptType = "unknown"
)

// DetectScript determines the script type from Unicode ranges.
func DetectScript(s string) ScriptType {
	for _, r := range s {
		// Hebrew: 0x0590-0x05FF
		if r >= 0x0590 && r <= 0x05FF {
			return ScriptHebrew
		}
		// Devanagari: 0x0900-0x097F
		if r >= 0x0900 && r <= 0x097F {
			return ScriptDevanagari
		}
		// Han: 0x4E00-0x9FFF, Extension A: 0x3400-0x4DBF
		if (r >= 0x4E00 && r <= 0x9FFF) || (r >= 0x3400 && r <= 0x4DBF) {
			return ScriptHan
		}
		// Arabic: 0x0600-0x06FF
		if r >= 0x0600 && r <= 0x06FF {
			return ScriptArabic
		}
		// Greek: 0x0370-0x03FF
		if r >= 0x0370 && r <= 0x03FF {
			return ScriptGreek
		}
	}
	// Default to Latin for plain ASCII text
	return ScriptLatin
}

// Forms holds generated forms from input processing.
type Forms struct {
	Normalized   string
	Script       ScriptType
	Tokens       []string
	Runes        []rune
	PhoneticKeys []string
	Fragments    []FragmentPath
}

// FragmentPath represents a possible split of the input into fragments.
type FragmentPath struct {
	Parts      []string
	Method     string
	Confidence float64
}

// Signal represents a resonance signal found by a channel.
type Signal struct {
	Text       string
	Target     string
	Channel    string
	Lens       string
	Confidence string
	Weight     float64
}

// ChannelResult holds results from a single resonance channel.
type ChannelResult struct {
	Name    string
	Signals []Signal
	Score   float64
}

// Pattern represents a converging or weak resonance pattern.
type Pattern struct {
	Name     string
	Signals  []Signal
	Strength float64
}

// Score holds overall, per-channel resonance scores, and components.
type Score struct {
	Overall    float64
	ByChannel  map[string]float64
	Components ScoreComponents
}

// ScoreComponents holds individual score components.
type ScoreComponents struct {
	ExactMatchScore         float64
	FuzzyMatchScore         float64
	GraphExpansionScore     float64
	PassageConvergenceScore float64
	MultiMethodBonus        float64
	ChannelDiversityBonus   float64
}

// Reading is the complete output of the decipher engine.
type Reading struct {
	Input              string
	Forms              Forms
	Candidates         []CandidateForm                                // Phase A: Candidate form generation
	FuzzyMatches       []MatchEvidence                                // Phase B: Fuzzy match evidence
	ConceptExpansions  map[string][]knowledge.DecipherConceptRelation // Phase D: Concept graph expansions
	PassageSignals     []PassageSignal                                // Phase E: Passage-level signals
	Convergence        ConvergenceResult                              // Phase E: Convergence via generic activation
	Channels           []ChannelResult
	ConvergingPatterns []Pattern
	WeakSignals        []Pattern
	Score              Score
	ConciseReading     string
	Warnings           []string
}
