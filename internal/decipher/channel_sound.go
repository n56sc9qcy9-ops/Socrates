package decipher

import (
	"strings"
)

// runSoundChannel analyzes phonetic patterns.
func runSoundChannel(forms Forms) ChannelResult {
	signals := make([]Signal, 0)

	input := forms.Normalized
	script := forms.Script

	if script != ScriptLatin {
		// For non-Latin scripts, sound channel produces script identification
		// NOT automatic spiritual meaning assignment
		signals = append(signals, Signal{
			Text:       "Non-Latin phonetic analysis",
			Target:     "unknown",
			Channel:    "Sound",
			Lens:       string(script) + "-sound",
			Confidence: ConfidenceSpeculative,
			Weight:     0.2,
		})
		return ChannelResult{Name: "Sound", Signals: signals, Score: calculateChannelScore(signals)}
	}

	// Consonant skeleton analysis
	skeleton := consonantSkeleton(input)
	if len(skeleton) > 0 {
		signals = append(signals, Signal{
			Text:       "consonant skeleton: " + skeleton,
			Target:     "phonetic-structure",
			Channel:    "Sound",
			Lens:       "sound",
			Confidence: ConfidenceVerified,
			Weight:     0.5,
		})
	}

	// Phonetic variant analysis
	for _, key := range forms.PhoneticKeys {
		if strings.HasPrefix(key, "phonetic:") {
			phonetic := key[9:]
			signals = append(signals, Signal{
				Text:       "phonetic form: " + phonetic,
				Target:     "phonetic-observation",
				Channel:    "Sound",
				Lens:       "sound",
				Confidence: ConfidencePlausible,
				Weight:     0.4,
			})
		}
	}

	score := calculateChannelScore(signals)

	return ChannelResult{
		Name:    "Sound",
		Signals: signals,
		Score:   score,
	}
}
