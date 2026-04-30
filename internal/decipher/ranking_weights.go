package decipher

// =============================================================================
// Ranking Weights
// =============================================================================

// RankingWeights defines explicit weights for evidence-path scoring.
// These weights control how strongly each evidence type contributes to rankings.
// Default weights preserve current scoring behavior as closely as practical.
type RankingWeights struct {
	// Exact match weights
	ExactMatchBase float64 // Base weight for exact form matches

	// Fuzzy match weights
	FuzzyMatchBase     float64 // Base weight for fuzzy matches
	FuzzyDistancePenalty float64 // Penalty per edit distance

	// Confidence weights
	ConfidenceVerified   float64 // Weight multiplier for verified confidence
	ConfidencePlausible  float64 // Weight multiplier for plausible confidence
	ConfidenceSpeculative float64 // Weight multiplier for speculative confidence

	// Graph relation weights
	GraphExpansionWeight    float64 // Weight for concept graph expansion
	GraphDepthPenalty       float64 // Penalty factor per propagation depth
	GraphRelationBaseWeight float64 // Base weight for relation support

	// Passage and convergence weights
	PassageCoActivationWeight float64 // Weight for co-activation in passage
	PassageFieldBoost         float64 // Boost for passage field concepts

	// Harmonic profile weights
	HarmonicProfileWeight    float64 // Weight for harmonic profile support
	HarmonicRatioCompatibility float64 // Weight for ratio compatibility
	HarmonicArchetypeMatch    float64 // Weight for archetype match

	// Multi-method and channel weights
	MultiMethodBonus          float64 // Bonus for multiple methods agreeing
	MultiMethodThreshold      int     // Number of methods needed for bonus
	ChannelDiversityBonus     float64 // Bonus for channel diversity
	ChannelDiversityThreshold int     // Number of active channels for bonus

	// Penalty weights
	DuplicateNoisePenalty float64 // Penalty for duplicate evidence
	DissonancePenalty     float64 // Penalty for unresolved dissonance

	// Bounding weights
	MaxSpeculativeRatio   float64 // Maximum ratio of speculative to total candidates

	// Source/lens weights (additional boost for curated sources)
	SourceCurated      float64 // Multiplier for curated source
	SourceTraditional  float64 // Multiplier for traditional source
	SourceHumanReview  float64 // Multiplier for human_review source
	LensOrthographic   float64 // Weight for orthographic lens
	LensPhonetic       float64 // Weight for phonetic lens
	LensSemantic       float64 // Weight for semantic lens
}

// DefaultRankingWeights returns weights that preserve current scoring behavior.
// Current behavior is:
func DefaultRankingWeights() RankingWeights {
	return RankingWeights{
		// Exact match: 0.25 weight in final score (components range 0-1)
		ExactMatchBase: 1.0,

		// Fuzzy match: 0.20 weight with 0.8 multiplier
		FuzzyMatchBase: 0.8,
		FuzzyDistancePenalty: 0.1, // ~10% penalty per distance

		// Confidence multipliers
		ConfidenceVerified: 1.0,
		ConfidencePlausible: 0.8,
		ConfidenceSpeculative: 0.5,

		// Graph expansion: 0.15 weight, 0.7 multiplier, normalized by count
		GraphExpansionWeight: 0.7,
		GraphDepthPenalty: 0.1, // 10% penalty per depth
		GraphRelationBaseWeight: 0.5,

		// Passage: 0.25 weight
		PassageCoActivationWeight: 0.6,
		PassageFieldBoost: 0.4,

		// Harmonic profile support
		HarmonicProfileWeight: 1.0,
		HarmonicRatioCompatibility: 0.5,
		HarmonicArchetypeMatch: 0.3,

		// Multi-method: 0.15 bonus at 3+ methods, 0.1 at 2+
		MultiMethodBonus: 0.15,
		MultiMethodThreshold: 3,
		// Channel diversity: 0.2 bonus at 4+, 0.1 at 3+
		ChannelDiversityBonus: 0.2,
		ChannelDiversityThreshold: 4,

		// Penalties
		DuplicateNoisePenalty: 0.1,
		DissonancePenalty: 0.2,

		// Bounding: speculative capped at 60% of total
		MaxSpeculativeRatio: 0.6,

		// Source/lens
		SourceCurated: 1.0,
		SourceTraditional: 0.9,
		SourceHumanReview: 0.85,
		LensOrthographic: 0.7,
		LensPhonetic: 0.9,
		LensSemantic: 1.0,
	}
}

// Validate checks if weights are within reasonable ranges.
func (w RankingWeights) Validate() []string {
	var errors []string

	if w.ExactMatchBase < 0 || w.ExactMatchBase > 10 {
		errors = append(errors, "ExactMatchBase must be in [0, 10]")
	}
	if w.FuzzyMatchBase < 0 || w.FuzzyMatchBase > 10 {
		errors = append(errors, "FuzzyMatchBase must be in [0, 10]")
	}
	if w.FuzzyDistancePenalty < 0 || w.FuzzyDistancePenalty > 1 {
		errors = append(errors, "FuzzyDistancePenalty must be in [0, 1]")
	}

	if w.ConfidenceVerified < 0 || w.ConfidenceVerified > 10 {
		errors = append(errors, "ConfidenceVerified must be in [0, 10]")
	}
	if w.ConfidencePlausible < 0 || w.ConfidencePlausible > 10 {
		errors = append(errors, "ConfidencePlausible must be in [0, 10]")
	}
	if w.ConfidenceSpeculative < 0 || w.ConfidenceSpeculative > 10 {
		errors = append(errors, "ConfidenceSpeculative must be in [0, 10]")
	}

	if w.GraphExpansionWeight < 0 || w.GraphExpansionWeight > 10 {
		errors = append(errors, "GraphExpansionWeight must be in [0, 10]")
	}
	if w.GraphDepthPenalty < 0 || w.GraphDepthPenalty > 1 {
		errors = append(errors, "GraphDepthPenalty must be in [0, 1]")
	}

	if w.PassageCoActivationWeight < 0 || w.PassageCoActivationWeight > 10 {
		errors = append(errors, "PassageCoActivationWeight must be in [0, 10]")
	}

	if w.HarmonicProfileWeight < 0 || w.HarmonicProfileWeight > 10 {
		errors = append(errors, "HarmonicProfileWeight must be in [0, 10]")
	}

	if w.MultiMethodBonus < 0 || w.MultiMethodBonus > 1 {
		errors = append(errors, "MultiMethodBonus must be in [0, 1]")
	}
	if w.ChannelDiversityBonus < 0 || w.ChannelDiversityBonus > 1 {
		errors = append(errors, "ChannelDiversityBonus must be in [0, 1]")
	}

	if w.DuplicateNoisePenalty < 0 || w.DuplicateNoisePenalty > 1 {
		errors = append(errors, "DuplicateNoisePenalty must be in [0, 1]")
	}
	if w.DissonancePenalty < 0 || w.DissonancePenalty > 1 {
		errors = append(errors, "DissonancePenalty must be in [0, 1]")
	}

	if w.MaxSpeculativeRatio < 0 || w.MaxSpeculativeRatio > 1 {
		errors = append(errors, "MaxSpeculativeRatio must be in [0, 1]")
	}

	return errors
}

// ConfidenceMultiplier returns the weight multiplier for a given confidence level.
func (w RankingWeights) ConfidenceMultiplier(confidence string) float64 {
	switch confidence {
	case "verified":
		return w.ConfidenceVerified
	case "plausible":
		return w.ConfidencePlausible
	case "speculative":
		return w.ConfidenceSpeculative
	default:
		return w.ConfidencePlausible
	}
}

// SourceMultiplier returns the weight multiplier for a given source.
func (w RankingWeights) SourceMultiplier(source string) float64 {
	switch source {
	case "curated":
		return w.SourceCurated
	case "traditional":
		return w.SourceTraditional
	case "human_review":
		return w.SourceHumanReview
	default:
		return 1.0
	}
}

// LensWeight returns the weight for a given lens.
func (w RankingWeights) LensWeight(lens string) float64 {
	switch lens {
	case "orthographic":
		return w.LensOrthographic
	case "phonetic":
		return w.LensPhonetic
	case "semantic":
		return w.LensSemantic
	default:
		return 0.5
	}
}