package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		fmt.Println("  socrates <word|phrase>   Analyze text for resonance (default)")
		fmt.Println("  socrates decipher <text>  Explicit decipher subcommand")
		fmt.Println("  socrates --help           Show full help with all commands and flags")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  decipher     Analyze a word or phrase for resonance")
		fmt.Println("  knowledge    Knowledge management (validate)")
		fmt.Println("  train        Training and evaluation")
		fmt.Println("  help         Show this help message")
		fmt.Println()
		fmt.Println("Quick start (root natural input):")
		fmt.Println("  socrates what is the purpose of my life")
		fmt.Println("  socrates love energy spirit")
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
		fmt.Println("Knowledge:")
		fmt.Println("  socrates knowledge validate")
		fmt.Println("  socrates knowledge validate --dir ./my-knowledge")
		fmt.Println()
		fmt.Println("Training:")
		fmt.Println("  socrates train")
		fmt.Println("  socrates train --debug")
		fmt.Println("  socrates train --examples ./my-examples.yaml")
		fmt.Println("  socrates train --heldout ./heldout.yaml")
		fmt.Println("  socrates train --weights ./weights.yaml")
		fmt.Println("  socrates train suggest-weights")
		fmt.Println("  socrates train apply-weights")
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

		// Check for suggest-weights subcommand
		args := trainCmd.Args()
		if len(args) > 0 && args[0] == "suggest-weights" {
			runSuggestWeights()
		} else if len(args) > 0 && args[0] == "apply-weights" {
			runApplyWeights()
		} else {
			runTrainEvaluate(*trainExamplesPath, *trainHeldOutPath, *trainWeightsPath, *trainDebugCmd)
		}

	case "suggest-weights":
		// Direct suggest-weights command
		runSuggestWeights()

	case "apply-weights":
		// Direct apply-weights command
		runApplyWeights()

	case "help", "-h", "--help":
		flag.Usage()

	case "version", "-v", "--version":
		fmt.Println("Socrates Language-Resonance Engine v0.1.0")
		fmt.Println("A pattern engine for exploring word resonance across languages.")

	default:
		// Root default: treat as decipher input
		// Join all remaining args into one passage
		input := command // first arg is part of input
		if len(os.Args) > 2 {
			// Join remaining arguments into one passage
			for i := 2; i < len(os.Args); i++ {
				if os.Args[i] != "" {
					if input != "" {
						input += " " + os.Args[i]
					} else {
						input = os.Args[i]
					}
				}
			}
		}

		// If input looks like a subcommand name (starts with "-" or matches known patterns),
		// show help instead
		if strings.HasPrefix(command, "-") {
			flag.Usage()
			os.Exit(1)
		}

		// Run decipher on the natural-language input
		if input == "" {
			fmt.Println("Error: please provide an input word or phrase")
			fmt.Println()
			fmt.Println("Usage: socrates <word|phrase>  - analyze text for resonance")
			fmt.Println("       socrates decipher <text>  - explicit decipher subcommand")
			fmt.Println()
			fmt.Println("For help: socrates --help")
			os.Exit(1)
		}

		// Run the engine
		engine, err := decipher.NewEngine()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize engine: %v\n", err)
			os.Exit(1)
		}
		reading := engine.Analyze(input)

		// Render the output based on mode (default to concise)
		output := decipher.RenderReadingWithOptions(reading, decipher.DefaultRenderOptions())
		fmt.Print(output)
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

// runSuggestWeights generates weight change suggestions and writes them to review files.
func runSuggestWeights() {
	// Load knowledge base
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading knowledge: %v\n", err)
		os.Exit(1)
	}

	// Load ranking weights
	weights := decipher.DefaultRankingWeights()

	// Create engine and suggester
	engine := decipher.NewEngineWithKnowledge(kb)
	suggester := training.NewWeightSuggester(engine, weights)
	suggester.SetReviewPath("review/weight_suggestions.yaml")

	// Load training examples
	loader := training.NewLoader()
	trainExamples, err := loader.LoadFromFile("training/examples.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading training examples: %v\n", err)
		os.Exit(1)
	}

	// Load held-out examples (optional)
	heldOutExamples, err := loader.LoadHeldOutFromFile("training/heldout.yaml")
	if err != nil {
		// No held-out examples is OK - we'll just use train
		heldOutExamples = training.Examples{}
		fmt.Println("Note: No held-out examples found, using train only")
	}

	fmt.Printf("Loaded %d train examples, %d held-out examples\n", len(trainExamples), len(heldOutExamples))

	// Get baseline metrics by evaluating with current weights
	eval := training.NewEvaluatorWithWeights(engine, weights, "baseline")
	report := eval.EvaluateWithReport(trainExamples)

	baseTrainMetrics := training.MetricsDelta{
		ConceptPrecision: report.Train.ConceptPrec,
		ConceptRecall:    report.Train.ConceptRec,
		FieldPrecision:   report.Train.FieldPrec,
		FieldRecall:      report.Train.FieldRec,
	}

	var baseHeldOutMetrics training.MetricsDelta
	if len(heldOutExamples) > 0 {
		heldOutReport := &training.EvaluationReport{}
		eval.EvaluateHeldOut(heldOutExamples, heldOutReport)
		heldOutReport.Finalize()
		baseHeldOutMetrics = training.MetricsDelta{
			ConceptPrecision: heldOutReport.HeldOut.ConceptPrec,
			ConceptRecall:    heldOutReport.HeldOut.ConceptRec,
			FieldPrecision:   heldOutReport.HeldOut.FieldPrec,
			FieldRecall:      heldOutReport.HeldOut.FieldRec,
		}
	}

	fmt.Println("\nBaseline metrics:")
	fmt.Printf("  Train: precision=%.4f, recall=%.4f\n", baseTrainMetrics.FieldPrecision, baseTrainMetrics.FieldRecall)
	fmt.Printf("  HeldOut: precision=%.4f, recall=%.4f\n", baseHeldOutMetrics.FieldPrecision, baseHeldOutMetrics.FieldRecall)

	// Generate suggestions
	fmt.Println("\nGenerating weight change suggestions...")
	candidates := suggester.GenerateSuggestions(trainExamples, heldOutExamples, baseTrainMetrics, baseHeldOutMetrics)

	fmt.Printf("Generated %d candidate changes\n", len(candidates))

	if len(candidates) == 0 {
		fmt.Println("No suggestions found. Current weights are acceptable.")
	} else {
		fmt.Println("\nSuggestions:")
		for i, c := range candidates {
			fmt.Printf("  [%d] %s: %v -> %v\n", i+1, c.Path, c.Current, c.Suggested)
			fmt.Printf("       rationale: %s\n", c.Rationale)
			fmt.Printf("       heldout delta: fp=%.4f, rec=%.4f\n", c.HeldOutDelta.FieldPrecision, c.HeldOutDelta.FieldRecall)
		}
	}

	// Write to review file
	if err := suggester.WriteSuggestions(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing suggestions: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Suggestions written to %s\n", suggester.ReviewPath)
	fmt.Println("Review and manually apply accepted changes.")
	fmt.Println("Do NOT auto-apply - review each suggestion first.")

	os.Exit(0)
}

// runApplyWeights applies accepted weight suggestions from a review file.
func runApplyWeights() {
	// Default paths
	reviewPath := "review/weight_suggestions.yaml"
	targetPath := "training/ranking_weights.yaml"
	auditPath := "review/weight_audit.yaml"

	// Parse any --review and --target flags if provided
	// For simplicity, use fixed paths (can be extended with flags later)

	fmt.Println("=== Apply Weight Suggestions ===")
	fmt.Printf("Review file: %s\n", reviewPath)
	fmt.Printf("Target config: %s\n", targetPath)

	// Check if review file exists
	if _, err := os.Stat(reviewPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: review file not found at %s\n", reviewPath)
		fmt.Println("Run 'socrates suggest-weights' first to generate suggestions.")
		os.Exit(1)
	}

	// Parse review file
	rf, parseErrors := training.ParseReviewFile(reviewPath)
	if len(parseErrors) > 0 {
		fmt.Fprintf(os.Stderr, "Error parsing review file:\n")
		for _, e := range parseErrors {
			fmt.Fprintf(os.Stderr, "  %s\n", e.Error())
		}
		os.Exit(1)
	}

	// Validate review file
	validationErrors := training.ValidateReviewFile(rf)
	if len(validationErrors) > 0 {
		fmt.Fprintf(os.Stderr, "Error validating review file:\n")
		for _, e := range validationErrors {
			fmt.Fprintf(os.Stderr, "  %s\n", e.Error())
		}
		os.Exit(1)
	}

	// Get pending, accepted, and rejected counts
	pending := training.GetPendingSuggestions(rf)
	accepted := training.GetAcceptedSuggestions(rf)
	rejected := training.GetRejectedSuggestions(rf)

	fmt.Printf("\nReview file status:\n")
	fmt.Printf("  Pending: %d\n", len(pending))
	fmt.Printf("  Accepted: %d\n", len(accepted))
	fmt.Printf("  Rejected: %d\n", len(rejected))

	// Check for accepted suggestions
	if len(accepted) == 0 {
		fmt.Println("\nNo accepted suggestions to apply.")
		fmt.Println("Edit the review file and set status to 'accepted' for suggestions you want to apply.")
		os.Exit(0)
	}

	// Display accepted suggestions
	fmt.Println("\nAccepted suggestions to apply:")
	for i, s := range accepted {
		fmt.Printf("  [%d] %s: %.2v -> %.2v\n", i+1, s.Path, s.Current, s.Suggested)
		if s.Rationale != "" {
			fmt.Printf("      (%s)\n", s.Rationale)
		}
	}

	// Confirm before applying
	fmt.Println("\nThis will update the target config file.")
	fmt.Print("Proceed with apply? (y/N): ")

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		fmt.Println("Aborted. No changes made.")
		os.Exit(0)
	}

	// Read current weights from target
	var currentWeights map[string]interface{}
	if data, err := os.ReadFile(targetPath); err == nil {
		// Parse existing weights
		if parsed, err := training.ParseWeightsYAML(string(data)); err == nil {
			currentWeights = parsed
		}
	}

	if currentWeights == nil {
		// Start with default weights if target doesn't exist
		currentWeights = make(map[string]interface{})
		fmt.Printf("Note: No existing config at %s, using default weights as baseline.\n", targetPath)
		// Initialize with default ranking weights
		defaultWeights := decipher.DefaultRankingWeights()
		currentWeights["exact_match_base"] = defaultWeights.ExactMatchBase
		currentWeights["fuzzy_match_base"] = defaultWeights.FuzzyMatchBase
		currentWeights["fuzzy_distance_penalty"] = defaultWeights.FuzzyDistancePenalty
		currentWeights["confidence_verified"] = defaultWeights.ConfidenceVerified
		currentWeights["confidence_plausible"] = defaultWeights.ConfidencePlausible
		currentWeights["confidence_speculative"] = defaultWeights.ConfidenceSpeculative
		currentWeights["graph_expansion_weight"] = defaultWeights.GraphExpansionWeight
		currentWeights["graph_depth_penalty"] = defaultWeights.GraphDepthPenalty
		currentWeights["graph_relation_base_weight"] = defaultWeights.GraphRelationBaseWeight
		currentWeights["passage_co_activation_weight"] = defaultWeights.PassageCoActivationWeight
		currentWeights["passage_field_boost"] = defaultWeights.PassageFieldBoost
		currentWeights["harmonic_profile_weight"] = defaultWeights.HarmonicProfileWeight
		currentWeights["harmonic_ratio_compatibility"] = defaultWeights.HarmonicRatioCompatibility
		currentWeights["harmonic_archetype_match"] = defaultWeights.HarmonicArchetypeMatch
		currentWeights["multi_method_bonus"] = defaultWeights.MultiMethodBonus
		currentWeights["multi_method_threshold"] = defaultWeights.MultiMethodThreshold
		currentWeights["channel_diversity_bonus"] = defaultWeights.ChannelDiversityBonus
		currentWeights["channel_diversity_threshold"] = defaultWeights.ChannelDiversityThreshold
		currentWeights["duplicate_noise_penalty"] = defaultWeights.DuplicateNoisePenalty
		currentWeights["dissonance_penalty"] = defaultWeights.DissonancePenalty
		currentWeights["max_speculative_ratio"] = defaultWeights.MaxSpeculativeRatio
		currentWeights["source_curated"] = defaultWeights.SourceCurated
		currentWeights["source_traditional"] = defaultWeights.SourceTraditional
		currentWeights["source_human_review"] = defaultWeights.SourceHumanReview
		currentWeights["lens_orthographic"] = defaultWeights.LensOrthographic
		currentWeights["lens_phonetic"] = defaultWeights.LensPhonetic
		currentWeights["lens_semantic"] = defaultWeights.LensSemantic
	}

	// Read existing audit trail
	existingAudit, _ := training.ReadAuditTrail(auditPath)

	// Apply accepted suggestions
	newAudit, applyErrors := training.ApplyWeightSuggestions(reviewPath, targetPath, currentWeights, existingAudit)
	if len(applyErrors) > 0 {
		fmt.Fprintf(os.Stderr, "Error applying suggestions:\n")
		for _, e := range applyErrors {
			fmt.Fprintf(os.Stderr, "  %s\n", e.Error())
		}
		fmt.Println("\nRefusing to apply due to drift or errors.")
		os.Exit(1)
	}

	// Write audit trail
	if err := training.WriteAuditTrail(auditPath, newAudit); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write audit trail: %v\n", err)
	}

	fmt.Printf("\n✓ Applied %d weight change(s)\n", len(accepted))
	fmt.Printf("  Target: %s\n", targetPath)
	fmt.Printf("  Audit: %s\n", auditPath)

	// Run validation on the new config
	fmt.Println("\nValidating updated weights...")
	if err := validateWeightsFile(targetPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: weight validation failed: %v\n", err)
	} else {
		fmt.Println("✓ Weight validation passed")
	}

	// Run training evaluation with new weights
	fmt.Println("\nRunning training evaluation with new weights...")

	// Load knowledge base
	kb, err := knowledge.LoadFromEmbed()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading knowledge: %v\n", err)
		os.Exit(1)
	}

	engine := decipher.NewEngineWithKnowledge(kb)

	// Load new weights
	newWeights, err := decipher.LoadRankingWeights(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading new weights: %v\n", err)
		os.Exit(1)
	}

	// Run evaluation
	eval := training.NewEvaluatorWithWeights(engine, *newWeights, targetPath)

	loader := training.NewLoader()
	trainExamples, err := loader.LoadFromFile("training/examples.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading training examples: %v\n", err)
		os.Exit(1)
	}

	report := eval.EvaluateWithReport(trainExamples)

	fmt.Println("\nTraining results with new weights:")
	fmt.Printf("  Concept: prec=%.4f, rec=%.4f\n", report.Train.ConceptPrec, report.Train.ConceptRec)
	fmt.Printf("  Field: prec=%.4f, rec=%.4f\n", report.Train.FieldPrec, report.Train.FieldRec)
	fmt.Printf("  Passed: %d/%d\n", report.Train.Passed, report.Train.Examples)

	os.Exit(0)
}

// validateWeightsFile checks that a weights YAML file is valid.
func validateWeightsFile(path string) error {
	// Parse the file
	if data, err := os.ReadFile(path); err != nil {
		// File doesn't exist yet
		if os.IsNotExist(err) {
			return nil
		}
		return err
	} else {
		_, err := training.ParseWeightsYAML(string(data))
		return err
	}
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
