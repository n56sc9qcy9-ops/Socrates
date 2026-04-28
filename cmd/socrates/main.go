package main

import (
	"fmt"

	"socrates/internal/resonance"
)

func main() {
	// Init
	space := resonance.NewFrequencySpace()
	engine := resonance.NewResonanceEngine(space)

	// Demos
	queries := []string{
		"truth", // EN
		"真理",    // CN truth
		"dog",   // EN
		"狗",     // CN dog
		"lie",   // Dissonant trigger
	}

	for _, q := range queries {
		resp, valid := engine.FullQuery(q)
		status := "VALID"
		if !valid {
			status = "DISSONANT"
		}
		fmt.Printf("Query: '%s'\nResponse: %s (%s)\n\n", q, resp, status)
	}

	fmt.Println("Socrates Resonance Demo Complete. Frequencies resonate mathematically!")
}
