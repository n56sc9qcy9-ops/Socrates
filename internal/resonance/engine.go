package resonance

// ResonanceEngine core component.
type ResonanceEngine struct {
	space *FrequencySpace
}

// NewResonanceEngine creates engine with space.
func NewResonanceEngine(space *FrequencySpace) *ResonanceEngine {
	return &ResonanceEngine{space: space}
}

// Query input idea, returns resonating frequencies via goroutines.
func (re *ResonanceEngine) Query(inputIdea string) []Frequency {
	inputFreq := re.space.ConvertToFrequency(inputIdea)
	allFreqs := re.space.AllFrequencies()

	// Sympathetic resonance with concurrency
	resonating := SympatheticResonance(inputFreq, allFreqs)

	return resonating
}

// ValidateResponse ensures internal consistency and meaningful resonance.
func (re *ResonanceEngine) ValidateResponse(res []Frequency) bool {
	// Empty response means no resonance found - invalid
	if len(res) == 0 {
		return false
	}
	// Check internal consistency
	return IsValidResponse(res)
}

// FullQuery performs query + validate + translate.
func (re *ResonanceEngine) FullQuery(inputIdea string) (string, bool) {
	res := re.Query(inputIdea)
	valid := re.ValidateResponse(res)
	english := re.space.ToEnglish(res)
	return english, valid
}
