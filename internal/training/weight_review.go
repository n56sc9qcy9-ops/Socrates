package training

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// WeightSuggestion represents a parsed suggestion from a review file.
type WeightSuggestion struct {
	ID               string      `yaml:"id"`
	Path             string      `yaml:"path"`
	Current          interface{} `yaml:"current"`
	Suggested        interface{} `yaml:"suggested"`
	Delta            string      `yaml:"delta"`
	Rationale        string      `yaml:"rationale"`
	Status           string      `yaml:"status"` // "pending", "accepted", "rejected"
	TrainDelta       MetricsDelta `yaml:"train_delta"`
	HeldOutDelta     MetricsDelta `yaml:"heldout_delta"`
}

// ReviewFile represents a parsed review file with metadata.
type ReviewFile struct {
	Version    string             `yaml:"version"`
	Generated  string             `yaml:"generated"`
	Suggestions []WeightSuggestion `yaml:"suggestions"`
	RawYAML    string             `yaml:"-"`
}

// ReviewValidationError captures errors during review file parsing.
type ReviewValidationError struct {
	Field   string
	Message string
}

func (e ReviewValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Field, e.Message)
}

// AllowedWeightFields lists the weight fields that can be modified.
// These are the flat key versions that map to nested YAML paths.
// This restricts suggestions to ranking weight parameters only.
var AllowedWeightFields = map[string]bool{
	"exact_match.base":              true,
	"fuzzy_match.base":              true,
	"fuzzy_match.distance_penalty": true,
	"confidence.verified":           true,
	"confidence.plausible":          true,
	"confidence.speculative":        true,
	"graph.expansion_weight":        true,
	"graph.depth_penalty":           true,
	"graph.relation_base_weight":    true,
	"passage.co_activation_weight":  true,
	"passage.field_boost":           true,
	"harmonic.profile_weight":       true,
	"harmonic.ratio_compatibility":  true,
	"harmonic.archetype_match":       true,
	"multi_method.bonus":            true,
	"multi_method.threshold":        true,
	"channel_diversity.bonus":       true,
	"channel_diversity.threshold":   true,
	"penalties.duplicate_noise":     true,
	"penalties.dissonance":          true,
	"bounding.max_speculative_ratio": true,
	"source.curated":               true,
	"source.traditional":            true,
	"source.human_review":          true,
	"lens.orthographic":             true,
	"lens.phonetic":                true,
	"lens.semantic":                true,
}

// FlatToNestedKey converts a flat key to nested YAML structure.
// E.g., "fuzzy_match.base" -> "fuzzy_match.base"
// The key format is already nested with dots.
func FlatToNestedKey(flat string) string {
	return flat
}

// ParseReviewFile reads and parses a weight suggestion review file.
func ParseReviewFile(path string) (*ReviewFile, []ReviewValidationError) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, []ReviewValidationError{{Field: "file", Message: err.Error()}}
	}

	var raw struct {
		Version    string                   `yaml:"version"`
		Generated  string                   `yaml:"generated"`
		Suggestions []map[string]interface{} `yaml:"suggestions"`
	}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, []ReviewValidationError{{Field: "yaml", Message: err.Error()}}
	}

	rf := &ReviewFile{
		Version:   raw.Version,
		Generated: raw.Generated,
		RawYAML:   string(data),
	}

	var errors []ReviewValidationError

	for i, rawSuggestion := range raw.Suggestions {
		suggestion, errs := parseSuggestion(i, rawSuggestion)
		if len(errs) > 0 {
			errors = append(errors, errs...)
		} else {
			rf.Suggestions = append(rf.Suggestions, *suggestion)
		}
	}

	return rf, errors
}

// parseSuggestion converts a raw map to a WeightSuggestion with validation.
func parseSuggestion(index int, raw map[string]interface{}) (*WeightSuggestion, []ReviewValidationError) {
	s := WeightSuggestion{}
	var errors []ReviewValidationError
	prefix := fmt.Sprintf("suggestions[%d]", index)

	// Parse ID
	if v, ok := raw["id"].(string); ok {
		s.ID = v
	} else {
		s.ID = fmt.Sprintf("suggestion_%d", index+1)
	}

	// Parse message (informational-only suggestions)
	if msg, ok := raw["message"].(string); ok {
		// Check if this is a "no suggestions" placeholder
		if msg == "No suggestions generated. Current weights are acceptable." || msg == "No suggestions found. Current weights are acceptable." {
			s.ID = fmt.Sprintf("placeholder_%d", index+1)
			s.Rationale = msg
			s.Status = "none" // Special status for placeholder messages
			return &s, nil
		}
		s.Rationale = msg
	}

	// Parse path
	if v, ok := raw["path"].(string); ok {
		s.Path = v
	} else if s.Status != "none" {
		// Only error if not a placeholder
		errors = append(errors, ReviewValidationError{
			Field:   prefix + ".path",
			Message: "missing or invalid path",
		})
	}

	// Parse current
	s.Current = raw["current"]

	// Parse suggested
	s.Suggested = raw["suggested"]

	// Parse delta
	if v, ok := raw["delta"].(string); ok {
		s.Delta = v
	}

	// Parse rationale
	if v, ok := raw["rationale"].(string); ok {
		s.Rationale = v
	}

	// Parse status
	if v, ok := raw["status"].(string); ok {
		s.Status = v
	} else {
		errors = append(errors, ReviewValidationError{
			Field:   prefix + ".status",
			Message: "missing or invalid status",
		})
	}

	// Parse train_delta
	if td, ok := raw["train_delta"].(map[string]interface{}); ok {
		s.TrainDelta = parseMetricsDelta(td, prefix+".train_delta")
	}

	// Parse heldout_delta
	if hd, ok := raw["heldout_delta"].(map[string]interface{}); ok {
		s.HeldOutDelta = parseMetricsDelta(hd, prefix+".heldout_delta")
	}

	return &s, errors
}

// parseMetricsDelta converts a raw map to a MetricsDelta.
func parseMetricsDelta(raw map[string]interface{}, prefix string) MetricsDelta {
	md := MetricsDelta{}

	if v, ok := raw["concept_precision"].(float64); ok {
		md.ConceptPrecision = v
	} else if v, ok := raw["concept_precision"].(int); ok {
		md.ConceptPrecision = float64(v)
	}

	if v, ok := raw["concept_recall"].(float64); ok {
		md.ConceptRecall = v
	} else if v, ok := raw["concept_recall"].(int); ok {
		md.ConceptRecall = float64(v)
	}

	if v, ok := raw["field_precision"].(float64); ok {
		md.FieldPrecision = v
	} else if v, ok := raw["field_precision"].(int); ok {
		md.FieldPrecision = float64(v)
	}

	if v, ok := raw["field_recall"].(float64); ok {
		md.FieldRecall = v
	} else if v, ok := raw["field_recall"].(int); ok {
		md.FieldRecall = float64(v)
	}

	if v, ok := raw["pass_rate_delta"].(float64); ok {
		md.PassRateDelta = v
	} else if v, ok := raw["pass_rate_delta"].(int); ok {
		md.PassRateDelta = float64(v)
	}

	return md
}

// ValidateReviewFile checks the semantic validity of a review file.
// It validates:
// - Status values are valid ("pending", "accepted", "rejected", "none")
// - Target paths are restricted to known ranking weight fields
// - Required fields are present for each suggestion
func ValidateReviewFile(rf *ReviewFile) []ReviewValidationError {
	var errors []ReviewValidationError

	for i, s := range rf.Suggestions {
		prefix := fmt.Sprintf("suggestions[%d]", i)

		// Skip placeholders with "none" status
		if s.Status == "none" {
			continue
		}

		// Validate status
		if s.Status != "pending" && s.Status != "accepted" && s.Status != "rejected" {
			errors = append(errors, ReviewValidationError{
				Field:   prefix + ".status",
				Message: fmt.Sprintf("invalid status '%s', expected 'pending', 'accepted', or 'rejected'", s.Status),
			})
		}

		// Validate path is a known weight field
		if !AllowedWeightFields[s.Path] {
			errors = append(errors, ReviewValidationError{
				Field:   prefix + ".path",
				Message: fmt.Sprintf("unknown weight field '%s'", s.Path),
			})
		}

		// Validate current/suggested present
		if s.Current == nil {
			errors = append(errors, ReviewValidationError{
				Field:   prefix + ".current",
				Message: "missing current value",
			})
		}
		if s.Suggested == nil {
			errors = append(errors, ReviewValidationError{
				Field:   prefix + ".suggested",
				Message: "missing suggested value",
			})
		}
	}

	return errors
}

// GetAcceptedSuggestions returns only suggestions with status "accepted".
func GetAcceptedSuggestions(rf *ReviewFile) []WeightSuggestion {
	var accepted []WeightSuggestion
	for _, s := range rf.Suggestions {
		if s.Status == "accepted" {
			accepted = append(accepted, s)
		}
	}
	return accepted
}

// GetPendingSuggestions returns only suggestions with status "pending".
func GetPendingSuggestions(rf *ReviewFile) []WeightSuggestion {
	var pending []WeightSuggestion
	for _, s := range rf.Suggestions {
		if s.Status == "pending" {
			pending = append(pending, s)
		}
	}
	return pending
}

// GetRejectedSuggestions returns only suggestions with status "rejected".
func GetRejectedSuggestions(rf *ReviewFile) []WeightSuggestion {
	var rejected []WeightSuggestion
	for _, s := range rf.Suggestions {
		if s.Status == "rejected" {
			rejected = append(rejected, s)
		}
	}
	return rejected
}

// AuditRecord represents a single applied weight change for audit trail.
type AuditRecord struct {
	Timestamp      string            `yaml:"timestamp"`
	SuggestionID   string            `yaml:"suggestion_id"`
	Path           string            `yaml:"path"`
	OldValue       interface{}       `yaml:"old_value"`
	NewValue       interface{}       `yaml:"new_value"`
	Rationale      string            `yaml:"rationale"`
	TrainDelta     MetricsDelta      `yaml:"train_delta"`
	HeldOutDelta   MetricsDelta      `yaml:"heldout_delta"`
	AppliedAt      string            `yaml:"applied_at"`
	ReviewFilePath string            `yaml:"review_file"`
}

// AuditTrail is a collection of audit records.
type AuditTrail struct {
	Records []AuditRecord `yaml:"audit"`
}

// ApplyWeightSuggestions applies accepted suggestions from a review file to a target config.
func ApplyWeightSuggestions(
	reviewPath string,
	targetPath string,
	currentWeights map[string]interface{},
	existingAudit []AuditRecord,
) ([]AuditRecord, []ReviewValidationError) {
	var errors []ReviewValidationError
	var newAudit []AuditRecord

	// Parse and validate review file
	rf, parseErrors := ParseReviewFile(reviewPath)
	if len(parseErrors) > 0 {
		return nil, parseErrors
	}

	validationErrors := ValidateReviewFile(rf)
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	// Get accepted suggestions
	accepted := GetAcceptedSuggestions(rf)
	if len(accepted) == 0 {
		errors = append(errors, ReviewValidationError{
			Field:   "status",
			Message: "no accepted suggestions found in review file",
		})
		return nil, errors
	}

	// Check for drift (current value in config doesn't match recorded current value)
	for i, s := range accepted {
		currentVal, exists := currentWeights[s.Path]
		if !exists {
			errors = append(errors, ReviewValidationError{
				Field:   fmt.Sprintf("suggestions[%d].path", i),
				Message: fmt.Sprintf("path '%s' not found in target config", s.Path),
			})
			continue
		}

		// Compare values (handle float/int casting)
		if !valuesMatch(currentVal, s.Current) {
			errors = append(errors, ReviewValidationError{
				Field:   fmt.Sprintf("suggestions[%d].current", i),
				Message: fmt.Sprintf("drift: recorded current %.2v but config has %.2v", s.Current, currentVal),
			})
		}
	}

	// If any errors, don't apply anything
	if len(errors) > 0 {
		return nil, errors
	}

	// Apply accepted suggestions
	updatedWeights := make(map[string]interface{})
	for k, v := range currentWeights {
		updatedWeights[k] = v
	}

	for _, s := range accepted {
		updatedWeights[s.Path] = s.Suggested

		// Create audit record
		audit := AuditRecord{
			Timestamp:      time.Now().Format(time.RFC3339),
			SuggestionID:   s.ID,
			Path:           s.Path,
			OldValue:       currentWeights[s.Path],
			NewValue:       s.Suggested,
			Rationale:      s.Rationale,
			TrainDelta:     s.TrainDelta,
			HeldOutDelta:   s.HeldOutDelta,
			AppliedAt:      time.Now().Format(time.RFC3339),
			ReviewFilePath: reviewPath,
		}
		newAudit = append(newAudit, audit)
	}

	// Write updated weights to target config
	if err := writeWeightsYAML(targetPath, updatedWeights); err != nil {
		return nil, []ReviewValidationError{{Field: "write", Message: err.Error()}}
	}

	// Combine with existing audit
	allAudit := append(existingAudit, newAudit...)

	return allAudit, nil
}

// valuesMatch compares two values for equality, handling float/int casting.
func valuesMatch(a, b interface{}) bool {
	// Handle float/int comparison
	af, aOk := toFloat64(a)
	bf, bOk := toFloat64(b)
	if aOk && bOk {
		// Allow small floating point difference
		return af == bf || (af > bf - 0.001 && af < bf + 0.001)
	}

	// Fall back to string comparison
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// toFloat64 converts a value to float64 if possible.
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

// writeWeightsYAML writes weights to a YAML file in the nested format.
func writeWeightsYAML(path string, weights map[string]interface{}) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Build YAML content in nested format matching training/ranking_weights.yaml
	content := "# Ranking Weights Configuration\n"
	content += "# Auto-generated by weight suggestion apply workflow\n"
	content += fmt.Sprintf("# Updated: %s\n\n", time.Now().Format(time.RFC3339))
	content += "# DO NOT EDIT MANUALLY - use review/apply workflow\n\n"
	content += "version: \"1.0\"\n\n"

	// Group weights by prefix
	sections := map[string][]string{
		"exact_match":    {},
		"fuzzy_match":    {},
		"confidence":     {},
		"graph":          {},
		"passage":        {},
		"harmonic":       {},
		"multi_method":   {},
		"channel_diversity": {},
		"penalties":      {},
		"bounding":       {},
		"source":         {},
		"lens":           {},
	}

	for k, v := range weights {
		// Parse the key prefix (everything before the last dot)
		parts := strings.Split(k, ".")
		if len(parts) < 2 {
			continue
		}
		prefix := parts[0]
		suffix := parts[1]

		var line string
		switch val := v.(type) {
		case float64:
			line = fmt.Sprintf("  %s: %.2f\n", suffix, val)
		case int:
			line = fmt.Sprintf("  %s: %d\n", suffix, val)
		default:
			line = fmt.Sprintf("  %s: %v\n", suffix, val)
		}

		if section, ok := sections[prefix]; ok {
			sections[prefix] = append(section, line)
		}
	}

	// Write sections in order
	for _, prefix := range []string{"exact_match", "fuzzy_match", "confidence", "graph",
		"passage", "harmonic", "multi_method", "channel_diversity", "penalties",
		"bounding", "source", "lens"} {
		lines := sections[prefix]
		if len(lines) == 0 {
			continue
		}
		content += fmt.Sprintf("%s:\n", prefix)
		for _, line := range lines {
			content += line
		}
		content += "\n"
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write weights file: %w", err)
	}

	return nil
}

// WriteAuditTrail writes an audit trail to a file.
func WriteAuditTrail(path string, trail []AuditRecord) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	at := AuditTrail{Records: trail}
	data, err := yaml.Marshal(at)
	if err != nil {
		return fmt.Errorf("failed to marshal audit trail: %w", err)
	}

	header := "# Weight Suggestion Audit Trail\n"
	header += "# This file records all applied weight changes.\n"
	header += "# Do not edit manually.\n\n"

	if err := os.WriteFile(path, []byte(header+string(data)), 0644); err != nil {
		return fmt.Errorf("failed to write audit trail: %w", err)
	}

	return nil
}

// ReadAuditTrail reads an existing audit trail from a file.
func ReadAuditTrail(path string) ([]AuditRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// File doesn't exist yet - return empty
		return []AuditRecord{}, nil
	}

	var at AuditTrail
	if err := yaml.Unmarshal(data, &at); err != nil {
		return nil, fmt.Errorf("failed to parse audit trail: %w", err)
	}

	return at.Records, nil
}

// ParseWeightsYAML parses a ranking weights YAML file.
func ParseWeightsYAML(content string) (map[string]interface{}, error) {
	type WeightsEntry struct {
		ExactMatch struct {
			Base float64 `yaml:"base"`
		} `yaml:"exact_match"`
		FuzzyMatch struct {
			Base           float64 `yaml:"base"`
			DistancePenalty float64 `yaml:"distance_penalty"`
		} `yaml:"fuzzy_match"`
		Confidence struct {
			Verified    float64 `yaml:"verified"`
			Plausible   float64 `yaml:"plausible"`
			Speculative float64 `yaml:"speculative"`
		} `yaml:"confidence"`
		Graph struct {
			ExpansionWeight   float64 `yaml:"expansion_weight"`
			DepthPenalty      float64 `yaml:"depth_penalty"`
			RelationBaseWeight float64 `yaml:"relation_base_weight"`
		} `yaml:"graph"`
		Passage struct {
			CoActivationWeight float64 `yaml:"co_activation_weight"`
			FieldBoost         float64 `yaml:"field_boost"`
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

	var entry WeightsEntry
	if err := yaml.Unmarshal([]byte(content), &entry); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	// Convert to flat map with nested keys
	result := make(map[string]interface{})

	if entry.ExactMatch.Base > 0 {
		result["exact_match.base"] = entry.ExactMatch.Base
	}
	if entry.FuzzyMatch.Base > 0 {
		result["fuzzy_match.base"] = entry.FuzzyMatch.Base
	}
	if entry.FuzzyMatch.DistancePenalty > 0 {
		result["fuzzy_match.distance_penalty"] = entry.FuzzyMatch.DistancePenalty
	}
	if entry.Confidence.Verified > 0 {
		result["confidence.verified"] = entry.Confidence.Verified
	}
	if entry.Confidence.Plausible > 0 {
		result["confidence.plausible"] = entry.Confidence.Plausible
	}
	if entry.Confidence.Speculative > 0 {
		result["confidence.speculative"] = entry.Confidence.Speculative
	}
	if entry.Graph.ExpansionWeight > 0 {
		result["graph.expansion_weight"] = entry.Graph.ExpansionWeight
	}
	if entry.Graph.DepthPenalty > 0 {
		result["graph.depth_penalty"] = entry.Graph.DepthPenalty
	}
	if entry.Graph.RelationBaseWeight > 0 {
		result["graph.relation_base_weight"] = entry.Graph.RelationBaseWeight
	}
	if entry.Passage.CoActivationWeight > 0 {
		result["passage.co_activation_weight"] = entry.Passage.CoActivationWeight
	}
	if entry.Passage.FieldBoost > 0 {
		result["passage.field_boost"] = entry.Passage.FieldBoost
	}
	if entry.Harmonic.ProfileWeight > 0 {
		result["harmonic.profile_weight"] = entry.Harmonic.ProfileWeight
	}
	if entry.Harmonic.RatioCompatibility > 0 {
		result["harmonic.ratio_compatibility"] = entry.Harmonic.RatioCompatibility
	}
	if entry.Harmonic.ArchetypeMatch > 0 {
		result["harmonic.archetype_match"] = entry.Harmonic.ArchetypeMatch
	}
	if entry.MultiMethod.Bonus > 0 {
		result["multi_method.bonus"] = entry.MultiMethod.Bonus
	}
	if entry.MultiMethod.Threshold > 0 {
		result["multi_method.threshold"] = entry.MultiMethod.Threshold
	}
	if entry.ChannelDiversity.Bonus > 0 {
		result["channel_diversity.bonus"] = entry.ChannelDiversity.Bonus
	}
	if entry.ChannelDiversity.Threshold > 0 {
		result["channel_diversity.threshold"] = entry.ChannelDiversity.Threshold
	}
	if entry.Penalties.DuplicateNoise > 0 {
		result["penalties.duplicate_noise"] = entry.Penalties.DuplicateNoise
	}
	if entry.Penalties.Dissonance > 0 {
		result["penalties.dissonance"] = entry.Penalties.Dissonance
	}
	if entry.Bounding.MaxSpeculativeRatio > 0 {
		result["bounding.max_speculative_ratio"] = entry.Bounding.MaxSpeculativeRatio
	}
	if entry.Source.Curated > 0 {
		result["source.curated"] = entry.Source.Curated
	}
	if entry.Source.Traditional > 0 {
		result["source.traditional"] = entry.Source.Traditional
	}
	if entry.Source.HumanReview > 0 {
		result["source.human_review"] = entry.Source.HumanReview
	}
	if entry.Lens.Orthographic > 0 {
		result["lens.orthographic"] = entry.Lens.Orthographic
	}
	if entry.Lens.Phonetic > 0 {
		result["lens.phonetic"] = entry.Lens.Phonetic
	}
	if entry.Lens.Semantic > 0 {
		result["lens.semantic"] = entry.Lens.Semantic
	}

	return result, nil
}