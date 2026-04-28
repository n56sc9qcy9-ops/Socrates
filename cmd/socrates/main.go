package main

import (
	"flag"
	"fmt"
	"os"

	"socrates/internal/decipher"
)

func main() {
	// Define command flags
	decipherCmd := flag.NewFlagSet("decipher", flag.ExitOnError)
	_ = decipherCmd.Bool("descifer", false, "alias for decipher")

	flag.Usage = func() {
		fmt.Println("Socrates Language-Resonance Engine")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  socrates decipher <word|phrase>")
		fmt.Println("  socrates descifer <word|phrase>  (alias)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  socrates decipher inspired")
		fmt.Println("  socrates decipher energy")
		fmt.Println("  socrates decipher \"in the beginning was the word\"")
		fmt.Println("  socrates decipher רוח")
		fmt.Println("  socrates decipher प्राण")
		fmt.Println("  socrates decipher 道")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "decipher", "descifer":
		decipherCmd.Parse(os.Args[2:])
		args := decipherCmd.Args()

		if len(args) < 1 {
			fmt.Println("Error: please provide an input word or phrase")
			fmt.Println()
			flag.Usage()
			os.Exit(1)
		}

		input := args[0]

		// Run the engine
		engine, err := decipher.NewEngine()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize engine: %v\n", err)
			os.Exit(1)
		}
		reading := engine.Analyze(input)

		// Render the output
		output := decipher.RenderReading(reading)
		fmt.Print(output)

	case "help", "-h", "--help":
		flag.Usage()

	case "version", "-v", "--version":
		fmt.Println("Socrates Language-Resonance Engine v0.1.0")
		fmt.Println("A pattern engine for exploring word resonance across languages.")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}
}
