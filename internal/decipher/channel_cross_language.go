package decipher

import (
	"strings"
)

// runCrossLanguageChannel analyzes cross-language echoes.
func runCrossLanguageChannel(forms Forms) ChannelResult {
	signals := make([]Signal, 0)

	input := forms.Normalized

	// Cross-language phonetic echoes
	crossEchos := map[string][]string{
		"chi":    {"chinese", "greek"},
		"qi":     {"chinese", "life-force"},
		"ruach":  {"hebrew", "breath"},
		"prana":  {"sanskrit", "breath"},
		"logos":  {"greek", "word"},
		"dao":    {"chinese", "way"},
		"om":     {"sanskrit", "sacred-sound"},
		"mantra": {"sanskrit", "sacred-utterance"},
	}

	for frag, langs := range crossEchos {
		if strings.Contains(input, frag) {
			for _, lang := range langs {
				signals = append(signals, Signal{
					Text:       frag + " cross-language echo",
					Target:     frag,
					Channel:    "Cross-Language",
					Lens:       lang,
					Confidence: ConfidencePlausible,
					Weight:     0.5,
				})
			}
		}
	}

	// Phonetic similarity to known roots
	phoneticRoots := map[string]string{
		"spir":  "spirit",
		"spire": "spirit",
		"生气":    "life-breath",
	}

	for inputPhonetic, target := range phoneticRoots {
		if forms.Normalized == inputPhonetic {
			signals = append(signals, Signal{
				Text:       "phonetic root: " + target,
				Target:     target,
				Channel:    "Cross-Language",
				Lens:       "cross-language",
				Confidence: ConfidencePlausible,
				Weight:     0.5,
			})
		}
	}

	score := calculateChannelScore(signals)

	return ChannelResult{
		Name:    "Cross-Language",
		Signals: signals,
		Score:   score,
	}
}
