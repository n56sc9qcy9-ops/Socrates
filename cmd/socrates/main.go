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
	debugMode := decipherCmd.Bool("debug", false, "show debug output including candidates, fuzzy matches, and internal details")

	flag.Usage = func() {
		fmt.Println("Socrates Language-Resonance Engine")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  socrates decipher <word|phrase> [flags]")
		fmt.Println("  socrates descifer <word|phrase>  (alias)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  socrates decipher inspired")
		fmt.Println("  socrates decipher energy")
		fmt.Println("  socrates decipher \"in the beginning was the word\"")
		fmt.Println("  socrates decipher רוח")
		fmt.Println("  socrates decipher प्राण")
		fmt.Println("  socrates decipher 道")
		fmt.Println("  socrates decipher skal --debug")
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
		// Parse flags from command line (after the subcommand)
		// Note: in Go's flag package, non-flag arguments come AFTER flags
		// So for "decipher --debug word", args come as ["--debug", "word"]
		// For "decipher word --debug", the input comes first
		decipherCmd.Parse(os.Args[2:])
		
		// Get positional args (non-flag arguments)
		args := decipherCmd.Args()
		
		var input string
		if len(args) >= 1 {
			input = args[0]
		} else {
			// Try to find input as a non-flag arg in the raw arguments
			for _, arg := range os.Args[2:] {
				if len(arg) > 0 && arg[0] != '-' {
					input = arg
					break
				}
			}
		}

		if input == "" {
			fmt.Println("Error: please provide an input word or phrase")
			fmt.Println()
			flag.Usage()
			os.Exit(1)
		}

		// Run the engine
		engine, err := decipher.NewEngine()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize engine: %v\n", err)
			os.Exit(1)
		}
		reading := engine.Analyze(input)

		// Render the output based on mode
		var output string
		if *debugMode {
			output = decipher.RenderReadingWithOptions(reading, decipher.DebugRenderOptions())
		} else {
			output = decipher.RenderReadingWithOptions(reading, decipher.DefaultRenderOptions())
		}
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
