package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"socrates/internal/decipher"
	"socrates/internal/knowledge"
	"socrates/internal/training"
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

	// Training subcommand
	trainCmd := flag.NewFlagSet("train", flag.ExitOnError)
	trainExamplesPath := trainCmd.String("examples", "", "path to training examples YAML file")
	trainHeldOutPath := trainCmd.String("heldout", "", "path to held-out examples YAML file (optional)")
	trainWeightsPath := trainCmd.String("weights", "", "path to ranking weights YAML file (optional, uses defaults)")
	trainDebugCmd := trainCmd.Bool("debug", false, "show detailed evaluation output")

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
		fmt.Println("  socrates train")
		fmt.Println("  socrates train --debug")
		fmt.Println("  socrates train --examples ./my-examples.yaml --debug")
		fmt.Println("  socrates train --heldout ./heldout.yaml")
		fmt.Println("  socrates train --weights ./weights.yaml")
		fmt.Println()
		fmt.Println("Training flags:")
		fmt.Println("  --examples <path>  Path to training examples YAML (default: training/examples.yaml)")
		fmt.Println("  --heldout <path>   Path to held-out examples YAML (optional)")
		fmt.Println("  --weights <path>   Path to ranking weights YAML (optional, uses default weights)")
		fmt.Println("  --debug            Show detailed evaluation output")
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

		// Check for positional subcommand
		args := knowledgeCmd.Args()
		if len(args) > 0 && args[0] == "validate" {
			runKnowledgeValidate(*validateDir)
		} else if *validateCmd {
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

	case "train":
		trainCmd.Parse(os.Args[2:])

		// Run evaluation (train always evaluates)
		runTrainEvaluate(*trainExamplesPath, *trainHeldOutPath, *trainWeightsPath, *trainDebugCmd)

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
		// Group warnings by category
		counts := result.WarningCountByCategory()
		categoryNames := map[knowledge.WarningCategory]string{
			knowledge.CategoryDuplicateForm:      "duplicate forms",
			knowledge.CategoryAliasResolution:    "alias resolutions",
			knowledge.CategoryUnknownConcept:     "unknown concepts",
			knowledge.CategoryUnusualProvenance:  "unusual provenance",
			knowledge.CategoryNonStandardRelation: "non-standard relations",
			knowledge.CategoryDataQuality:         "data quality issues",
		}

		fmt.Printf("⚠ Found %d warning(s) in %d categories:\n", len(result.Warnings), len(counts))
		for cat, count := range counts {
			name := categoryNames[cat]
			if name == "" {
				name = string(cat)
			}
			fmt.Printf("  %s: %d\n", name, count)
		}
		fmt.Println()

		// Group and print warnings by category
		warningsByCategory := make(map[knowledge.WarningCategory][]string)
		for _, w := range result.Warnings {
			warningsByCategory[w.Category] = append(warningsByCategory[w.Category], fmt.Sprintf("  %s: %s", w.Field, w.Message))
		}

		// Print warnings grouped by category
		for cat := range counts {
			name := categoryNames[cat]
			if name == "" {
				name = string(cat)
			}
			fmt.Printf("[%s]\n", name)
			for _, msg := range warningsByCategory[cat] {
				fmt.Println(msg)
			}
			fmt.Println()
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("Validation FAILED")
		os.Exit(1)
	}

	// Final report
	if len(result.Warnings) == 0 {
		fmt.Println("✓ Validation passed with no warnings")
	} else {
		catCount := len(result.WarningCountByCategory())
		fmt.Printf("✓ Validation passed with %d warning(s) in %d categories\n", len(result.Warnings), catCount)
	}
	os.Exit(0)
}

// runTrainEvaluate runs the training evaluation command.
func runTrainEvaluate(examplesPath string, heldOutPath string, weightsPath string, debugMode bool) {
	// Load knowledge base
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading knowledge: %v\n", err)
		os.Exit(1)
	}

	// Load ranking weights if specified
	var weights decipher.RankingWeights
	weightsSource := "default"

	if weightsPath != "" {
		weightsSource = weightsPath
		var loadErr error
		loadedWeights, loadErr := decipher.LoadRankingWeights(weightsPath)
		if loadErr != nil {
			fmt.Fprintf(os.Stderr, "Error loading ranking weights from %s: %v\n", weightsPath, loadErr)
			os.Exit(1)
		}
		weights = *loadedWeights
		fmt.Printf("Loaded ranking weights from %s\n", weightsPath)
	} else {
		weights = decipher.DefaultRankingWeights()
	}

	// Load training examples
	loader := training.NewLoader()
	var examples training.Examples
	var heldOutExamples training.Examples

	if examplesPath != "" {
		// Load from specified file
		absPath, err := filepath.Abs(examplesPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid examples path: %v\n", err)
			os.Exit(1)
		}

		examples, err = loader.LoadFromFile(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading training examples from %s: %v\n", absPath, err)
			os.Exit(1)
		}
		fmt.Printf("Loaded %d examples from %s\n", len(examples), absPath)
	} else {
		// Load from default training/examples.yaml
		defaultPath := "training/examples.yaml"

		if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: no examples file specified and %s not found\n", defaultPath)
			fmt.Println("Use --examples <path> to specify a training examples file")
			os.Exit(1)
		}

		examples, err = loader.LoadFromFile(defaultPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading default training examples: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Loaded %d examples from default file\n", len(examples))
	}

	// Load held-out examples if specified
	if heldOutPath != "" {
		absPath, err := filepath.Abs(heldOutPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid held-out path: %v\n", err)
			os.Exit(1)
		}

		heldOutExamples, err = loader.LoadHeldOutFromFile(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading held-out examples from %s: %v\n", absPath, err)
			os.Exit(1)
		}
		fmt.Printf("Loaded %d held-out examples from %s\n", len(heldOutExamples), absPath)
	}

	// Validate examples against knowledge
	fmt.Println("\nValidating examples against knowledge base...")
	validationResult := training.ValidateExamples(examples, kb)

	if !validationResult.IsValid() {
		fmt.Printf("✗ Found %d validation error(s):\n", len(validationResult.Errors))
		for _, e := range validationResult.Errors {
			fmt.Printf("  [%s] %s: %s\n", e.ExampleID, e.Field, e.Message)
		}
		fmt.Println("\nFix validation errors before running evaluation")
		os.Exit(1)
	}
	fmt.Printf("✓ All %d examples validated\n", len(examples))

	// Validate held-out examples if loaded
	if len(heldOutExamples) > 0 {
		heldOutValidationResult := training.ValidateExamples(heldOutExamples, kb)
		if !heldOutValidationResult.IsValid() {
			fmt.Printf("✗ Found %d validation error(s) in held-out examples:\n", len(heldOutValidationResult.Errors))
			for _, e := range heldOutValidationResult.Errors {
				fmt.Printf("  [%s] %s: %s\n", e.ExampleID, e.Field, e.Message)
			}
			fmt.Println("\nFix validation errors before running evaluation")
			os.Exit(1)
		}
		fmt.Printf("✓ All %d held-out examples validated\n", len(heldOutExamples))
	}

	// Create engine and evaluator with weights
	engine, err := decipher.NewEngine()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating engine: %v\n", err)
		os.Exit(1)
	}

	evaluator := training.NewEvaluatorWithWeights(engine, weights, weightsSource)

	// Run full evaluation with train and held-out splits
	fmt.Println("\nRunning evaluation...")

	var report *training.EvaluationReport
	if len(heldOutExamples) > 0 {
		report = evaluator.RunFullEvaluation(examples, heldOutExamples)
	} else {
		report = evaluator.EvaluateWithReport(examples)
	}

	// Print the evaluation report
	fmt.Print(report.FormatReport())

	// Print detailed results in debug mode
	if debugMode {
		fmt.Println("\n--- DETAILED TRAIN RESULTS ---")
		for _, ex := range report.Train.Detailed {
			fmt.Printf("[%s] %s\n", ex.ExampleID, ex.Input)
			fmt.Printf("  Passed: %v\n", ex.Passed)
			fmt.Printf("  Concepts: %d hits, %d misses, %d false pos (prec=%.2f, rec=%.2f)\n",
				ex.ConceptHits, ex.ConceptMisses, ex.ConceptFalsePos,
				ex.ConceptPrecision, ex.ConceptRecall)
			fmt.Printf("  Fields: %d hits, %d misses, %d false pos (prec=%.2f, rec=%.2f)\n",
				ex.FieldHits, ex.FieldMisses, ex.FieldFalsePos,
				ex.FieldPrecision, ex.FieldRecall)
			if len(ex.MissedExpectedFields) > 0 {
				fmt.Printf("  Missed fields: %v\n", ex.MissedExpectedFields)
			}
			if len(ex.FalseActivatedFields) > 0 {
				fmt.Printf("  False fields: %v\n", ex.FalseActivatedFields)
			}
			fmt.Println()
		}

		if len(heldOutExamples) > 0 {
			fmt.Println("\n--- DETAILED HELD-OUT RESULTS ---")
			for _, ex := range report.HeldOut.Detailed {
				fmt.Printf("[%s] %s\n", ex.ExampleID, ex.Input)
				fmt.Printf("  Passed: %v\n", ex.Passed)
				fmt.Printf("  Concepts: %d hits, %d misses, %d false pos (prec=%.2f, rec=%.2f)\n",
					ex.ConceptHits, ex.ConceptMisses, ex.ConceptFalsePos,
					ex.ConceptPrecision, ex.ConceptRecall)
				fmt.Printf("  Fields: %d hits, %d misses, %d false pos (prec=%.2f, rec=%.2f)\n",
					ex.FieldHits, ex.FieldMisses, ex.FieldFalsePos,
					ex.FieldPrecision, ex.FieldRecall)
				if len(ex.MissedExpectedFields) > 0 {
					fmt.Printf("  Missed fields: %v\n", ex.MissedExpectedFields)
				}
				if len(ex.FalseActivatedFields) > 0 {
					fmt.Printf("  False fields: %v\n", ex.FalseActivatedFields)
				}
				fmt.Println()
			}
		}
	}

	// Exit with appropriate code
	if report.TotalFailed > 0 {
		fmt.Printf("\n%d example(s) failed evaluation\n", report.TotalFailed)
		os.Exit(1)
	} else {
		fmt.Println("\n✓ All examples passed evaluation")
		os.Exit(0)
	}
}
