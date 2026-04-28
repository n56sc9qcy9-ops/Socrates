package resonance

import (
	"math"
	"sync"
)

// Frequency represents a meaning as tonal frequencies.
type Frequency struct {
	Tone1    float64 // Truth (Root)
	Tone2    float64 // Love (Harmonic Binding)
	Tone3    float64 // Being (Resonant Field)
	Harmonic float64 // Combined resonance
}

var (
	Truth  = Frequency{Tone1: 1.0, Tone2: 0.5, Tone3: 0.5, Harmonic: 2.0}
	Love   = Frequency{Tone1: 0.5, Tone2: 1.0, Tone3: 0.5, Harmonic: 2.0}
	Being  = Frequency{Tone1: 0.5, Tone2: 0.5, Tone3: 1.0, Harmonic: 2.0}
	Power  = Frequency{Tone1: 0.6, Tone2: 0.4, Tone3: 1.0, Harmonic: 2.2}
	Wisdom = Frequency{Tone1: 0.8, Tone2: 0.0, Tone3: 0.8, Harmonic: 2.5}

	Lie = Frequency{Tone1: 0.1, Tone2: 0.1, Tone3: 0.1, Harmonic: 0.3}
)

// AreInResonance checks tonal cosine similarity >= 0.7 AND harmonic ratio.
func AreInResonance(f1, f2 Frequency) bool {
	if f1.Harmonic == 0 || f2.Harmonic == 0 {
		return true
	}
	// Legacy zero-tone support: if both zero-norm, use harmonic only
	n1 := vectorNorm(f1.Tone1, f1.Tone2, f1.Tone3)
	n2 := vectorNorm(f2.Tone1, f2.Tone2, f2.Tone3)
	if n1 == 0 || n2 == 0 {
		ratio := f1.Harmonic / f2.Harmonic
		if ratio < 1 {
			ratio = 1 / ratio
		}
		return isHarmonicRatio(ratio) && f1.Harmonic >= 0.5 && f2.Harmonic >= 0.5
	}
	cosine := cosineSimilarity(f1.Tone1, f1.Tone2, f1.Tone3, f2.Tone1, f2.Tone2, f2.Tone3)
	ratio := f1.Harmonic / f2.Harmonic
	if ratio < 1 {
		ratio = 1 / ratio // Normalize to >=1
	}
	return cosine >= 0.7 && isHarmonicRatio(ratio)
}

// dotProduct computes dot product of two 3D tone vectors.
func dotProduct(t11, t12, t13, t21, t22, t23 float64) float64 {
	return t11*t21 + t12*t22 + t13*t23
}

// vectorNorm computes Euclidean norm of 3D tone vector.
func vectorNorm(t1, t2, t3 float64) float64 {
	return math.Sqrt(t1*t1 + t2*t2 + t3*t3)
}

// isHarmonicRatio is unexported helper (used by AreInResonance).
// Checks common musical ratios (unison, octave, fifth, thirds, golden).
func isHarmonicRatio(r float64) bool {
	if r < 1 {
		r = 1 / r
	}
	return r >= 0.99 && r <= 1.01 || // Unison ~1.0
		r >= 1.8 && r <= 2.5 || // Octave ~2.0 (extended)

		r >= 1.48 && r <= 1.52 || // Perfect Fifth 1.5
		r >= 1.25 && r <= 1.35 || // Major/Minor Thirds
		math.Abs(r-1.618) <= 0.02 // Golden Ratio
}

// cosineSimilarity computes cosine sim between two tone vectors (requires math).
func cosineSimilarity(t11, t12, t13, t21, t22, t23 float64) float64 {
	n1 := vectorNorm(t11, t12, t13)
	n2 := vectorNorm(t21, t22, t23)
	if n1 == 0 || n2 == 0 {
		return 0
	}
	return dotProduct(t11, t12, t13, t21, t22, t23) / (n1 * n2)
}

// SympatheticResonance finds all frequencies resonating with input using goroutines.
func SympatheticResonance(inputFreq Frequency, allFrequencies []Frequency) []Frequency {
	var resonating []Frequency
	var mu sync.Mutex
	var wg sync.WaitGroup

	results := make(chan Frequency, len(allFrequencies))

	for _, freq := range allFrequencies {
		wg.Add(1)
		go func(f Frequency) {
			defer wg.Done()
			if AreInResonance(inputFreq, f) {
				results <- f
			}
		}(freq)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for freq := range results {
		mu.Lock()
		resonating = append(resonating, freq)
		mu.Unlock()
	}

	return resonating
}

// IsValidResponse checks if all frequencies in set are mutually resonant (no dissonance).
func IsValidResponse(frequencies []Frequency) bool {
	for i := 0; i < len(frequencies); i++ {
		for j := i + 1; j < len(frequencies); j++ {
			if !AreInResonance(frequencies[i], frequencies[j]) {
				return false
			}
		}
	}
	return true
}
