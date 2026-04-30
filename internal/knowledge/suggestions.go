package knowledge

// Suggestion represents a suggested addition to knowledge that has not been vetted.
// Suggestions are NOT automatically added to trusted knowledge - they require human review.
type Suggestion struct {
	// Type of suggestion: "concept", "form", "relation", "script_word", "glyph_pattern"
	Type string

	// Proposed data
	ProposedID   string // For concepts
	ProposedName string // For concepts
	ProposedForm string // For forms
	ProposedWord string // For script words
	ProposedFrom string // For relations
	ProposedTo   string // For relations
	ProposedType string // For relations

	// Source and evidence
	Source     string  // Where this suggestion came from (e.g., "user-input:skall", "pattern:phonetic")
	Evidence   string  // Evidence for this suggestion
	Confidence string  // "verified", "plausible", or "speculative"
	Weight     float64 // Suggested weight

	// Review metadata
	Rationale  string // Why this might be valid
	Notes      string // Additional notes
	Reviewed   bool   // Whether this has been reviewed
	Accepted   bool   // Whether this was accepted (set after review)
	ReviewedBy string // Who reviewed this
}

// SuggestionReviewer provides functionality for reviewing suggestions.
type SuggestionReviewer struct {
	suggestions []Suggestion
}

// NewSuggestionReviewer creates a new suggestion reviewer.
func NewSuggestionReviewer() *SuggestionReviewer {
	return &SuggestionReviewer{
		suggestions: make([]Suggestion, 0),
	}
}

// AddSuggestion adds a new suggestion to the review queue.
func (r *SuggestionReviewer) AddSuggestion(s Suggestion) {
	r.suggestions = append(r.suggestions, s)
}

// GetSuggestions returns all suggestions.
func (r *SuggestionReviewer) GetSuggestions() []Suggestion {
	return r.suggestions
}

// GetUnreviewedSuggestions returns suggestions that haven't been reviewed yet.
func (r *SuggestionReviewer) GetUnreviewedSuggestions() []Suggestion {
	var unreviewed []Suggestion
	for _, s := range r.suggestions {
		if !s.Reviewed {
			unreviewed = append(unreviewed, s)
		}
	}
	return unreviewed
}

// MarkReviewed marks a suggestion as reviewed.
func (r *SuggestionReviewer) MarkReviewed(index int, accepted bool, reviewer string) {
	if index >= 0 && index < len(r.suggestions) {
		r.suggestions[index].Reviewed = true
		r.suggestions[index].Accepted = accepted
		r.suggestions[index].ReviewedBy = reviewer
	}
}

// SuggestionForUnknownInput analyzes unknown input and suggests possible knowledge additions.
// This does NOT auto-add to trusted knowledge - it creates a suggestion for review.
func SuggestionForUnknownInput(input string, kb *Knowledge) *Suggestion {
	s := &Suggestion{
		Source:   "user-input",
		Evidence: "unknown input: " + input,
	}

	// Check if this could be a new form for an existing concept
	existingForms := kb.GetFormsByText(input)
	if len(existingForms) > 0 {
		// Not unknown - get concept
		return nil
	}

	// Check for potential phonetic matches
	// This is just a suggestion, not a determination
	s.Type = "form"
	s.ProposedForm = input
	s.Confidence = "speculative"
	s.Rationale = "user input that didn't match existing forms"

	return s
}

// SuggestionForWeaklyMatchedInput creates suggestions for weakly matched inputs.
// These can help a curator decide whether to add new knowledge.
func SuggestionForWeaklyMatchedInput(input string, matches []Form, kb *Knowledge) []Suggestion {
	var suggestions []Suggestion

	// If input only weakly matches, suggest it might need a new form
	if len(matches) == 0 || (len(matches) == 1 && matches[0].Weight < 0.5) {
		suggestions = append(suggestions, Suggestion{
			Type:         "form",
			ProposedForm: input,
			Source:       "weak-match",
			Evidence:     "input matched no or very weak forms",
			Confidence:   "speculative",
			Weight:       0.3,
			Rationale:    "suggest adding explicit form mapping for this input",
		})
	}

	return suggestions
}

// FormatSuggestionsForReview formats suggestions for human review.
// Output is suitable for a review file or console output.
func FormatSuggestionsForReview(suggestions []Suggestion) string {
	if len(suggestions) == 0 {
		return "No suggestions to review.\n"
	}

	result := "# Knowledge Addition Suggestions\n"
	result += "# Review Status: pending review | accepted | rejected\n"
	result += "# Format: Type | Proposed Data | Source | Evidence | Confidence | Weight | Rationale\n"
	result += "#\n"
	result += "# To accept a suggestion, add it to the appropriate YAML file with the specified data.\n"
	result += "# To reject, add \"# REJECTED\" comment and your reasoning below.\n"
	result += "#\n\n"

	for i, s := range suggestions {
		result += formatSuggestion(i, s)
		result += "\n"
	}

	return result
}

func formatSuggestion(index int, s Suggestion) string {
	result := "## Suggestion " + itoa(index+1) + "\n"
	result += "Type: " + s.Type + "\n"

	switch s.Type {
	case "concept":
		result += "Proposed ID: " + s.ProposedID + "\n"
		result += "Proposed Name: " + s.ProposedName + "\n"
	case "form":
		result += "Proposed Form: " + s.ProposedForm + "\n"
	case "relation":
		result += "Proposed From: " + s.ProposedFrom + "\n"
		result += "Proposed To: " + s.ProposedTo + "\n"
		result += "Proposed Type: " + s.ProposedType + "\n"
	}

	result += "Source: " + s.Source + "\n"
	result += "Evidence: " + s.Evidence + "\n"
	result += "Confidence: " + s.Confidence + "\n"
	if s.Weight > 0 {
		result += "Weight: " + ftos(s.Weight) + "\n"
	}
	result += "Rationale: " + s.Rationale + "\n"

	if s.Notes != "" {
		result += "Notes: " + s.Notes + "\n"
	}

	if s.Reviewed {
		if s.Accepted {
			result += "Status: ACCEPTED by " + s.ReviewedBy + "\n"
		} else {
			result += "Status: REJECTED by " + s.ReviewedBy + "\n"
		}
	} else {
		result += "Status: PENDING REVIEW\n"
	}

	return result
}

// MinimumReviewRecord returns the minimum fields required for a proper review record.
// This ensures consistency when manually adding knowledge after review.
type MinimumReviewRecord struct {
	// What type of addition
	Type string

	// For concept additions
	ConceptID   string
	ConceptName string
	Aliases     []string

	// For form additions
	Form       string
	Concept    string
	Lens       string
	Weight     float64
	Confidence string

	// For relation additions
	From      string
	To        string
	RelType   string
	RelWeight float64

	// Review metadata
	EvidenceSource  string
	ConfidenceLevel string
	WeightValue     float64 // Using WeightValue to avoid confusion with field weights
	Rationale       string
	Notes           string
	Accepted        bool
	ReviewedBy      string
}
