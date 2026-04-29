package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"socrates/internal/decipher"
	"socrates/internal/knowledge"
)

func main() {
	// Define command flags
	decipherCmd := flag.NewFlagSet("decipher", flag.ExitOnError)
	_ = decipherCmd.Bool("descifer", false, "alias for decipher")
	debugMode := decipherCmd.Bool("debug", false, "show debug output including candidates, fuzzy matches, and internal details")

	// Knowledge subcommand
	knowledgeCmd := flag.NewFlagSet("knowledge", flag.ExitOnError)
	validateCmd := knowledgeCmd.Bool("validate", false, "validate knowledge data")
	validateDir := knowledgeCmd.String("dir", "", "directory containing knowledge YAML files to validate")

	flag.Usage = func() {
		fmt.Println("Socrates Language-Resonance Engine")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  socrates decipher <word|phrase> [flags]")
		fmt.Println("  socrates descifer <word|phrase>  (alias)")
		fmt.Println("  socrates knowledge validate [--dir <path>]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  decipher, descifer  Analyze a word or phrase for resonance")
		fmt.Println("  knowledge          Knowledge management commands")
		fmt.Println("  help               Show this help message")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  socrates decipher inspired")
		fmt.Println("  socrates decipher energy")
		fmt.Println("  socrates decipher \"in the beginning was the word\"")
		fmt.Println("  socrates decipher רוח")
		fmt.Println("  socrates decipher प्राण")
		fmt.Println("  socrates decipher 道")
		fmt.Println("  socrates decipher skal --debug")
		fmt.Println("  socrates knowledge validate")
		fmt.Println("  socrates knowledge validate --dir ./my-knowledge")
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

	case "knowledge":
		knowledgeCmd.Parse(os.Args[2:])
		
		if *validateCmd {
			runKnowledgeValidate(*validateDir)
		} else {
			fmt.Println("Knowledge management commands:")
			fmt.Println("  validate  Validate knowledge data")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  socrates knowledge validate [--dir <path>]")
			fmt.Println()
			fmt.Println("Examples:")
			fmt.Println("  socrates knowledge validate          # validate embedded knowledge")
			fmt.Println("  socrates knowledge validate --dir .  # validate local YAML files")
			os.Exit(1)
		}

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

// runKnowledgeValidate runs the knowledge validation command.
func runKnowledgeValidate(dir string) {
	var kb *knowledge.Knowledge
	var err error

	if dir != "" {
		// Load from directory
		absDir, err := filepath.Abs(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid directory path: %v\n", err)
			os.Exit(1)
		}
		
		kb, err = knowledge.LoadFromDir(absDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading knowledge from %s: %v\n", absDir, err)
			os.Exit(1)
		}
		fmt.Printf("Validating knowledge from directory: %s\n", absDir)
	} else {
		// Load embedded knowledge
		kb, err = knowledge.LoadFromEmbed()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading embedded knowledge: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Validating embedded knowledge:")
	}

	// Run validation
	result := knowledge.ValidateKnowledge(kb)

	// Print results
	fmt.Println()
	
	if len(result.Errors) == 0 && len(result.Warnings) == 0 {
		fmt.Println("✓ Knowledge is valid (no errors or warnings)")
		os.Exit(0)
	}

	if len(result.Errors) > 0 {
		fmt.Printf("✗ Found %d error(s):\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("  %s: %s\n", e.Field, e.Message)
		}
		fmt.Println()
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("⚠ Found %d warning(s):\n", len(result.Warnings))
		for _, w := range result.Warnings {
			fmt.Printf("  %s: %s\n", w.Field, w.Message)
		}
		fmt.Println()
	}

	if len(result.Errors) > 0 {
		fmt.Println("Validation FAILED")
		os.Exit(1)
	} else {
		fmt.Println("Validation passed with warnings")
		os.Exit(0)
	}
}
