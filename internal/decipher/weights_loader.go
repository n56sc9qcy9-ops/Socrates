package decipher

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RankingWeightsEntry is the YAML structure for ranking weights.
type RankingWeightsEntry struct {
	Version     string `yaml:"version"`
	Description string `yaml:"description,omitempty"`

	ExactMatch struct {
		Base float64 `yaml:"base"`
	} `yaml:"exact_match"`

	FuzzyMatch struct {
		Base          float64 `yaml:"base"`
		DistancePenalty float64 `yaml:"distance_penalty"`
	} `yaml:"fuzzy_match"`

	Confidence struct {
		Verified    float64 `yaml:"verified"`
		Plausible   float64 `yaml:"plausible"`
		Speculative float64 `yaml:"speculative"`
	} `yaml:"confidence"`

	Graph struct {
		ExpansionWeight    float64 `yaml:"expansion_weight"`
		DepthPenalty       float64 `yaml:"depth_penalty"`
		RelationBaseWeight float64 `yaml:"relation_base_weight"`
	} `yaml:"graph"`

	Passage struct {
		CoActivationWeight float64 `yaml:"co_activation_weight"`
		FieldBoost          float64 `yaml:"field_boost"`
	} `yaml:"passage"`

	Harmonic struct {
		ProfileWeight       float64 `yaml:"profile_weight"`
		RatioCompatibility  float64 `yaml:"ratio_compatibility"`
		ArchetypeMatch      float64 `yaml:"archetype_match"`
	} `yaml:"harmonic"`

	MultiMethod struct {
		Bonus     float64 `yaml:"bonus"`
		Threshold int     `yaml:"threshold"`
	} `yaml:"multi_method"`

	ChannelDiversity struct {
		Bonus     float64 `yaml:"bonus"`
		Threshold int     `yaml:"threshold"`
	} `yaml:"channel_diversity"`

	Penalties struct {
		DuplicateNoise float64 `yaml:"duplicate_noise"`
		Dissonance     float64 `yaml:"dissonance"`
	} `yaml:"penalties"`

	Bounding struct {
		MaxSpeculativeRatio float64 `yaml:"max_speculative_ratio"`
	} `yaml:"bounding"`

	Source struct {
		Curated      float64 `yaml:"curated"`
		Traditional  float64 `yaml:"traditional"`
		HumanReview  float64 `yaml:"human_review"`
	} `yaml:"source"`

	Lens struct {
		Orthographic float64 `yaml:"orthographic"`
		Phonetic     float64 `yaml:"phonetic"`
		Semantic     float64 `yaml:"semantic"`
	} `yaml:"lens"`
}

// LoadRankingWeights loads ranking weights from a YAML file.
func LoadRankingWeights(path string) (*RankingWeights, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read weights file: %w", err)
	}

	var entry RankingWeightsEntry
	if err := yaml.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("failed to parse weights file: %w", err)
	}

	weights := DefaultRankingWeights()

	// Apply loaded values where specified (0 means not set)
	if entry.ExactMatch.Base > 0 {
		weights.ExactMatchBase = entry.ExactMatch.Base
	}
	if entry.FuzzyMatch.Base > 0 {
		weights.FuzzyMatchBase = entry.FuzzyMatch.Base
	}
	if entry.FuzzyMatch.DistancePenalty > 0 {
		weights.FuzzyDistancePenalty = entry.FuzzyMatch.DistancePenalty
	}

	if entry.Confidence.Verified > 0 {
		weights.ConfidenceVerified = entry.Confidence.Verified
	}
	if entry.Confidence.Plausible > 0 {
		weights.ConfidencePlausible = entry.Confidence.Plausible
	}
	if entry.Confidence.Speculative > 0 {
		weights.ConfidenceSpeculative = entry.Confidence.Speculative
	}

	if entry.Graph.ExpansionWeight > 0 {
		weights.GraphExpansionWeight = entry.Graph.ExpansionWeight
	}
	if entry.Graph.DepthPenalty > 0 {
		weights.GraphDepthPenalty = entry.Graph.DepthPenalty
	}
	if entry.Graph.RelationBaseWeight > 0 {
		weights.GraphRelationBaseWeight = entry.Graph.RelationBaseWeight
	}

	if entry.Passage.CoActivationWeight > 0 {
		weights.PassageCoActivationWeight = entry.Passage.CoActivationWeight
	}
	if entry.Passage.FieldBoost > 0 {
		weights.PassageFieldBoost = entry.Passage.FieldBoost
	}

	if entry.Harmonic.ProfileWeight > 0 {
		weights.HarmonicProfileWeight = entry.Harmonic.ProfileWeight
	}
	if entry.Harmonic.RatioCompatibility > 0 {
		weights.HarmonicRatioCompatibility = entry.Harmonic.RatioCompatibility
	}
	if entry.Harmonic.ArchetypeMatch > 0 {
		weights.HarmonicArchetypeMatch = entry.Harmonic.ArchetypeMatch
	}

	if entry.MultiMethod.Bonus > 0 {
		weights.MultiMethodBonus = entry.MultiMethod.Bonus
	}
	if entry.MultiMethod.Threshold > 0 {
		weights.MultiMethodThreshold = entry.MultiMethod.Threshold
	}

	if entry.ChannelDiversity.Bonus > 0 {
		weights.ChannelDiversityBonus = entry.ChannelDiversity.Bonus
	}
	if entry.ChannelDiversity.Threshold > 0 {
		weights.ChannelDiversityThreshold = entry.ChannelDiversity.Threshold
	}

	if entry.Penalties.DuplicateNoise > 0 {
		weights.DuplicateNoisePenalty = entry.Penalties.DuplicateNoise
	}
	if entry.Penalties.Dissonance > 0 {
		weights.DissonancePenalty = entry.Penalties.Dissonance
	}

	if entry.Bounding.MaxSpeculativeRatio > 0 {
		weights.MaxSpeculativeRatio = entry.Bounding.MaxSpeculativeRatio
	}

	if entry.Source.Curated > 0 {
		weights.SourceCurated = entry.Source.Curated
	}
	if entry.Source.Traditional > 0 {
		weights.SourceTraditional = entry.Source.Traditional
	}
	if entry.Source.HumanReview > 0 {
		weights.SourceHumanReview = entry.Source.HumanReview
	}

	if entry.Lens.Orthographic > 0 {
		weights.LensOrthographic = entry.Lens.Orthographic
	}
	if entry.Lens.Phonetic > 0 {
		weights.LensPhonetic = entry.Lens.Phonetic
	}
	if entry.Lens.Semantic > 0 {
		weights.LensSemantic = entry.Lens.Semantic
	}

	// Validate
	if errors := weights.Validate(); len(errors) > 0 {
		return nil, fmt.Errorf("invalid weights: %v", errors)
	}

	return &weights, nil
}

// SaveRankingWeights saves ranking weights to a YAML file.
func SaveRankingWeights(path string, weights RankingWeights) error {
	entry := RankingWeightsEntry{
		Version: "1.0",
		Description: "Explicit ranking weights for Socrates evidence-path scoring.",
	}

	entry.ExactMatch.Base = weights.ExactMatchBase
	entry.FuzzyMatch.Base = weights.FuzzyMatchBase
	entry.FuzzyMatch.DistancePenalty = weights.FuzzyDistancePenalty

	entry.Confidence.Verified = weights.ConfidenceVerified
	entry.Confidence.Plausible = weights.ConfidencePlausible
	entry.Confidence.Speculative = weights.ConfidenceSpeculative

	entry.Graph.ExpansionWeight = weights.GraphExpansionWeight
	entry.Graph.DepthPenalty = weights.GraphDepthPenalty
	entry.Graph.RelationBaseWeight = weights.GraphRelationBaseWeight

	entry.Passage.CoActivationWeight = weights.PassageCoActivationWeight
	entry.Passage.FieldBoost = weights.PassageFieldBoost

	entry.Harmonic.ProfileWeight = weights.HarmonicProfileWeight
	entry.Harmonic.RatioCompatibility = weights.HarmonicRatioCompatibility
	entry.Harmonic.ArchetypeMatch = weights.HarmonicArchetypeMatch

	entry.MultiMethod.Bonus = weights.MultiMethodBonus
	entry.MultiMethod.Threshold = weights.MultiMethodThreshold

	entry.ChannelDiversity.Bonus = weights.ChannelDiversityBonus
	entry.ChannelDiversity.Threshold = weights.ChannelDiversityThreshold

	entry.Penalties.DuplicateNoise = weights.DuplicateNoisePenalty
	entry.Penalties.Dissonance = weights.DissonancePenalty

	entry.Bounding.MaxSpeculativeRatio = weights.MaxSpeculativeRatio

	entry.Source.Curated = weights.SourceCurated
	entry.Source.Traditional = weights.SourceTraditional
	entry.Source.HumanReview = weights.SourceHumanReview

	entry.Lens.Orthographic = weights.LensOrthographic
	entry.Lens.Phonetic = weights.LensPhonetic
	entry.Lens.Semantic = weights.LensSemantic

	data, err := yaml.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal weights: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write weights file: %w", err)
	}

	return nil
}