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
	Original     string // Original input string (preserved before tokenization)
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
	// IsDirect indicates this signal comes from direct form/glyph/script evidence
	// (not from symbolic neighbor expansion). Used to distinguish evidence provenance.
	IsDirect bool
}

// EvidenceID returns a deterministic identity for this signal.
// Used for deduplication before scoring.
func (s Signal) EvidenceID() string {
	// Include channel, target, and text for meaningful identity
	// Same target via different channels is independent evidence
	return s.Channel + "|" + s.Target + "|" + s.Text
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

// CounsellorField contains data-backed transmutation suggestions.
// Built from curated transmute relations — not from hardcoded Go.
type CounsellorField struct {
	// SourceFields are the activated tension/contracted fields.
	SourceFields []TransmutationSource

	// Suggestions are the data-backed correction paths.
	Suggestions []TransmutationSuggestion

	// EvidencePaths trace from activated concepts to suggestions.
	EvidencePaths []TransmutationEvidence
}

// TransmutationSource is a tension/contracted field detected in the reading.
type TransmutationSource struct {
	// Concept is the concept ID of the source field.
	Concept string

	// Strength is the activation strength (0-1).
	Strength float64

	// Depth is the propagation depth (0 = direct, 1+ = graph-expanded/neighbor).
	Depth int

	// IsDirectEvidence is true if this source has direct form/glyph/script evidence.
	// False means evidence comes only from symbolic neighbor expansion.
	IsDirectEvidence bool

	// EvidenceSources are where this field was activated from.
	EvidenceSources []string
}

// TransmutationSuggestion is a data-backed correction path.
type TransmutationSuggestion struct {
	// SourceConcept is the source field concept ID.
	SourceConcept string

	// TargetConcept is the correction target concept ID.
	TargetConcept string

	// Kind is the transmutation kind (transmutes_to, softens_through, etc.).
	Kind string

	// Confidence is the confidence level.
	Confidence string

	// Source is the provenance (curated, traditional, human_review).
	Source string

	// Lens is the interpretive lens.
	Lens string

	// Weight is the relation weight (0-1).
	Weight float64

	// Strength is the combined strength (source field × relation weight).
	Strength float64
}

// TransmutationEvidence traces a suggestion back to its data source.
type TransmutationEvidence struct {
	SourceConcept  string
	TargetConcept  string
	Kind           string
	DataSourceFile string
	Notes          string
	IsDirectSource bool // true if from direct concept, false if from graph-propagated neighbor
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
	PassageFields      PassageFields                                  // Phase E: Passage fields via activation graph
	HarmonicField      *HarmonicField                                 // Phase F: Harmonic field from frequency profiles
	CounsellorField    *CounsellorField                               // Phase G: Data-backed transmutation/counsellor fields
	Channels           []ChannelResult
	ConvergingPatterns []Pattern
	WeakSignals        []Pattern
	Score              Score
	ConciseReading     string
	Warnings           []string
	// Discarded counts for bounded work
	DiscardedCandidates  int `json:"discardedCandidates,omitempty"`  // Candidates skipped due to bounds
	DiscardedComparisons int `json:"discardedComparisons,omitempty"` // Fuzzy comparisons skipped due to bounds
}

// HarmonicField represents the harmonic field derived from active concepts.
// Built from activated concept frequency profiles - not from hardcoded Go data.
type HarmonicField struct {
	// Tones are the individual tones in the harmonic field.
	Tones []HarmonicTone

	// Coherence is the overall harmonic coherence (0-1).
	Coherence float64

	// Consonance is the consonance score (compatible ratios).
	Consonance float64

	// Dissonance is the dissonance score (conflicting ratios).
	Dissonance float64

	// EvidencePaths trace from activated concepts to field tones.
	EvidencePaths []HarmonicEvidence
}

// HarmonicTone represents a single tone in the harmonic field.
type HarmonicTone struct {
	// MeaningFrequencyID is the stable identity of this tone.
	MeaningFrequencyID string

	// SourceConcepts are the activated concepts that contributed to this tone.
	SourceConcepts []string

	// Vector is the 3D tone vector [Tone1, Tone2, Tone3].
	Vector [3]int

	// Ratio is the frequency ratio [numerator, denominator].
	Ratio [2]int

	// Archetype is the archetype reference ID.
	Archetype string

	// Labels are the semantic labels.
	Labels knowledge.IntFrequencyLabels

	// Strength is the combined strength from all contributing concepts.
	Strength float64

	// Confidence is the weighted confidence from source profiles.
	Confidence string
}

// HarmonicEvidence traces from activated concept to harmonic tone.
type HarmonicEvidence struct {
	// SourceConcept is the activated concept.
	SourceConcept string

	// SourceConceptWeight is the passage field weight of the source concept.
	SourceConceptWeight float64

	// MeaningFrequencyID is the matched profile's ID.
	MeaningFrequencyID string

	// ProfileWeight is the frequency profile's weight (0-100).
	ProfileWeight int

	// Confidence is the profile's confidence.
	Confidence string
}

// CandidateBounds defines limits for candidate generation.
// High-confidence candidates (normalized, skeleton, phonetic) are always kept.
// Speculative candidates (edit variants, n-grams, prefixes, suffixes) are capped.
type CandidateBounds struct {
	// MaxCandidates is the maximum total candidates allowed.
	MaxCandidates int
	// MaxSpeculativeCandidates caps speculative candidates (edit variants, n-grams, etc.).
	// If 0, uses MaxCandidates as default.
	MaxSpeculativeCandidates int
}

// DefaultCandidateBounds returns sensible defaults for bounded candidate generation.
func DefaultCandidateBounds() CandidateBounds {
	return CandidateBounds{
		MaxCandidates:            50,
		MaxSpeculativeCandidates: 30,
	}
}

// FuzzyBounds defines limits for fuzzy matching comparisons.
type FuzzyBounds struct {
	// MaxComparisons is the maximum candidate-anchor comparisons allowed.
	MaxComparisons int
	// MaxMatches caps the output fuzzy matches.
	MaxMatches int
}

// DefaultFuzzyBounds returns sensible defaults for bounded fuzzy matching.
func DefaultFuzzyBounds() FuzzyBounds {
	return FuzzyBounds{
		MaxComparisons: 500,
		MaxMatches:     20,
	}
}

// RenderMode controls what level of detail is included in rendered output.
type RenderMode int

const (
	// RenderModeDefault shows concise output: input, top evidence paths, score summary, warnings.
	RenderModeDefault RenderMode = iota
	// RenderModeDebug shows full output including candidates, fuzzy matches, generated forms, and internal details.
	RenderModeDebug
)

// RenderOptions controls rendering behavior.
type RenderOptions struct {
	Mode RenderMode
}

// DefaultRenderOptions returns default render options (default mode).
func DefaultRenderOptions() RenderOptions {
	return RenderOptions{
		Mode: RenderModeDefault,
	}
}

// DebugRenderOptions returns render options for debug output.
func DebugRenderOptions() RenderOptions {
	return RenderOptions{
		Mode: RenderModeDebug,
	}
}
