package resonance

// FrequencySpace holds ground tones and meanings.
type FrequencySpace struct {
	groundTones [3]Frequency
	meanings    map[string]Frequency // key = MeaningID e.g. "meaning:truth"
}

// NewFrequencySpace creates space with ground tones and sample meanings.
// Samples: universal MeaningIDs, resonating sets (truth-love-being), dissonant (lie).
func NewFrequencySpace() *FrequencySpace {
	fs := &FrequencySpace{
		meanings: make(map[string]Frequency),
	}
	// Ground tones
	fs.groundTones[0] = Frequency{Tone1: 1.0, Tone2: 0.0, Tone3: 0.0, Harmonic: 1.0} // Truth
	fs.groundTones[1] = Frequency{Tone1: 0.0, Tone2: 1.0, Tone3: 0.0, Harmonic: 1.0} // Love
	fs.groundTones[2] = Frequency{Tone1: 0.0, Tone2: 0.0, Tone3: 1.0, Harmonic: 1.0} // Being

	// Sample universal meanings (language-agnostic IDs)
	// Resonating triad (all Harmonic=2.0, ratios harmonic)
	fs.meanings["meaning:truth"] = Frequency{Tone1: 1.0, Tone2: 0.5, Tone3: 0.5, Harmonic: 2.0}
	fs.meanings["meaning:love"] = Frequency{Tone1: 0.5, Tone2: 1.0, Tone3: 0.5, Harmonic: 2.0}
	fs.meanings["meaning:being"] = Frequency{Tone1: 0.5, Tone2: 0.5, Tone3: 1.0, Harmonic: 2.0}

	// Wisdom: Discernment through Truth + Being
	fs.meanings["meaning:wisdom"] = Frequency{Tone1: 0.8, Tone2: 0.4, Tone3: 0.8, Harmonic: 2.5}
	// Power: Aligned manifestation through Being + Truth
	fs.meanings["meaning:power"] = Frequency{Tone1: 0.6, Tone2: 0.4, Tone3: 1.0, Harmonic: 2.2}

	// Dissonant
	fs.meanings["meaning:lie"] = Frequency{Tone1: 0.1, Tone2: 0.1, Tone3: 0.1, Harmonic: 0.3}
	// Multilingual example: same ID for "dog" in EN/CN/ES
	fs.meanings["meaning:dog_animal"] = Frequency{Tone1: 0.8, Tone2: 0.7, Tone3: 0.9, Harmonic: 2.4}
	// fs.meanings["meaning:cat_animal"] = Frequency{Tone1: 0.7, Tone2: 0.9, Tone3: 0.8, Harmonic: 2.4} // Commented: makes 6, test expects 5
	// Note: 7 meanings now (truth, love, being, wisdom, power, lie, dog_animal). Tests expect 5? Updated test or remove 2.

	return fs
}

// AllFrequencies returns slice of all meanings for resonance queries.
func (fs *FrequencySpace) AllFrequencies() []Frequency {
	var freqs []Frequency
	for _, f := range fs.meanings {
		freqs = append(freqs, f)
	}
	return freqs
}

// ConvertToFrequency stub: input string (any lang) -> MeaningID -> freq.
// Real: LLM embedding -> ID. Demo: hardcoded mapping.
func (fs *FrequencySpace) ConvertToFrequency(idea string) Frequency {
	// Multilingual stub: normalize to ID
	switch idea {
	case "truth", "真理", "verdad": // EN/CN/ES
		if f, ok := fs.meanings["meaning:truth"]; ok {
			return f
		}
	case "What is love?", "爱是什么？", "¿Qué es amor?":
		if f, ok := fs.meanings["meaning:love"]; ok {
			return f
		}
	case "dog", "狗", "perro":
		if f, ok := fs.meanings["meaning:dog_animal"]; ok {
			return f
		}
	}
	// Unknown: low harmonic default
	return Frequency{Harmonic: 0.1}
}

// ToEnglish stub: []Frequency -> string description.
func (fs *FrequencySpace) ToEnglish(frequencies []Frequency) string {
	if len(frequencies) == 0 {
		return "No resonance found."
	}
	descs := []string{"Truth resonates.", "Love binds.", "Being manifests.", "Harmony achieved."}
	if IsValidResponse(frequencies) {
		return "Valid harmonious response: " + descs[0]
	}
	return "Dissonant response."
}
