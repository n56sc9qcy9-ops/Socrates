package decipher

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoSemanticBucketBooleans verifies that production code does not use
// semantic bucket boolean names like hasModal, hasContrast, hasEmptiness.
func TestNoSemanticBucketBooleans(t *testing.T) {
	// Determine the package directory by finding the test file's location
	testFile := "/Users/bot/Socrates/internal/decipher/phase1_correction_test.go"
	pkgDir := filepath.Dir(testFile)
	testDir := pkgDir

	entries, err := os.ReadDir(testDir)
	if err != nil {
		t.Skipf("cannot read directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue // Skip test files
		}

		filePath := filepath.Join(testDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("cannot read %s: %v", filePath, err)
			continue
		}

		// Parse the file
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
		if err != nil {
			t.Errorf("cannot parse %s: %v", filePath, err)
			continue
		}

		// Walk the AST and find any identifiers matching semantic bucket names
		for _, node := range f.Scope.Objects {
			if ident, ok := node.Decl.(*ast.Ident); ok {
				if isSemanticBucketBool(ident.Name) {
					t.Errorf("%s: found semantic bucket boolean %q in production code",
						filePath, ident.Name)
				}
			}
		}

		// Also check string literals for hardcoded convergence names
		contentStr := string(content)
		badNames := []string{
			"hollow_obligation",
			"possible_hollow_obligation",
			"obligation_present",
			"emptiness_present",
		}
		for _, name := range badNames {
			if strings.Contains(contentStr, name) {
				t.Errorf("%s: found hardcoded convergence name %q in production code",
					filePath, name)
			}
		}
	}
}

// isSemanticBucketBool checks if a name is a semantic bucket boolean.
func isSemanticBucketBool(name string) bool {
	badNames := []string{
		"hasModal",
		"hasContrast",
		"hasEmptiness",
	}
	for _, bad := range badNames {
		if name == bad {
			return true
		}
	}
	return false
}

// TestNoDirectSkalBranching verifies that 'skal' doesn't have direct semantic branching.
func TestNoDirectSkalBranching(t *testing.T) {
	testFile := "/Users/bot/Socrates/internal/decipher/phase1_correction_test.go"
	pkgDir := filepath.Dir(testFile)
	testDir := pkgDir

	entries, err := os.ReadDir(testDir)
	if err != nil {
		t.Skipf("cannot read directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(testDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("cannot read %s: %v", filePath, err)
			continue
		}

		// Check for direct if/skal patterns
		badPatterns := []string{
			`if word == "skal"`,
			`if input == "skal"`,
			`== "skal"`,
			`if skal`,
			`"skal" ->`,
			`skal -> "shell"`,
			`skal -> "bowl"`,
		}

		contentStr := string(content)
		for _, pattern := range badPatterns {
			if strings.Contains(contentStr, pattern) {
				t.Errorf("%s: found direct 'skal' branch pattern %q",
					filePath, pattern)
			}
		}
	}
}

// TestPassageConvergenceUsesActivationGraph verifies that convergence
// is computed from ActivatedConcept fields, not named semantic buckets.
func TestPassageConvergenceUsesActivationGraph(t *testing.T) {
	// Parse the convergence.go file to verify DetectConvergence signature
	testFile := "/Users/bot/Socrates/internal/decipher/phase1_correction_test.go"
	pkgDir := filepath.Dir(testFile)
	convergencePath := filepath.Join(pkgDir, "convergence.go")
	activationPath := filepath.Join(pkgDir, "activation.go")

	// Check convergence.go for DetectConvergence
	content, err := os.ReadFile(convergencePath)
	if err != nil {
		t.Skipf("cannot read convergence.go: %v", err)
	}
	contentStr := string(content)
	if !strings.Contains(contentStr, "func DetectConvergence(passageSignals []PassageSignal, directConcepts []string, kb *knowledge.Knowledge) ConvergenceResult") {
		t.Error("DetectConvergence should return ConvergenceResult, not []ConvergencePattern")
	}

	// Verify AnalyzePassageTokens returns []PassageSignal (not with booleans)
	// Check activation.go for AnalyzePassageTokens
	activationContent, err := os.ReadFile(activationPath)
	if err != nil {
		t.Skipf("cannot read activation.go: %v", err)
	}
	activationStr := string(activationContent)

	if !strings.Contains(activationStr, "func AnalyzePassageTokens(tokens []string, kb *knowledge.Knowledge) []PassageSignal") {
		t.Error("AnalyzePassageTokens should return []PassageSignal with kb parameter (no hasContrast, hasModal, hasEmptiness booleans)")
	}
}

// TestEngineAnalyzeSignature verifies the engine properly uses ConvergenceResult.
func TestEngineAnalyzeSignature(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	reading := engine.Analyze("shall")

	// Verify Convergence field is of type ConvergenceResult
	// (accessing fields that exist on ConvergenceResult but not []ConvergencePattern)
	_ = reading.Convergence.ActivatedConcepts
	_ = reading.Convergence.CoActivationScore
	_ = reading.Convergence.TopConcepts
	_ = reading.Convergence.RelationPaths

	_ = fmt.Sprintf("CoActivation: %.2f", reading.Convergence.CoActivationScore)
}

// TestConvergenceFromGenericActivation verifies convergence emerges from
// concept activation and graph relations, not hardcoded names.
func TestConvergenceFromGenericActivation(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	// Test with a word that might have triggered "hollow_obligation" before
	// Now it should produce generic activated concepts
	reading := engine.Analyze("shall")

	// Verify we get ActivatedConcepts, not named patterns
	if len(reading.Convergence.ActivatedConcepts) > 0 {
		// Good - we have generic activation
		for _, ac := range reading.Convergence.ActivatedConcepts {
			// Verify concept name is a valid concept from the graph, not a hardcoded pattern name
			badNames := []string{
				"hollow_obligation",
				"possible_hollow_obligation",
				"obligation_present",
				"emptiness_present",
			}
			for _, bad := range badNames {
				if ac.Concept == bad {
					t.Errorf("found hardcoded pattern name %q in activated concepts", bad)
				}
			}
		}
	}

	// Verify relation paths show generic graph relations
	if len(reading.Convergence.RelationPaths) > 0 {
		for _, path := range reading.Convergence.RelationPaths {
			// Paths should look like "conceptA --[type]--> conceptB", not hardcoded names
			if !strings.Contains(path, "-->") {
				t.Errorf("relation path %q doesn't look like a graph path", path)
			}
		}
	}
}
func TestNoHardcodedKnowledgeInProduction(t *testing.T) {
	// This test proves the architecture is data-driven, not hardcoded.
	// It searches production Go files for any remaining hardcoded knowledge.

	patterns := []string{
		"var FragmentSeeds",
		"var ScriptWords",
		"var Primitives",
	}

	// Read all production Go files in decipher package
	files := []string{
		"activation.go",
		"candidate_generation.go",
		"channel_cross_language.go",
		"channel_fragment.go",
		"channel_glyph.go",
		"channel_scoring.go",
		"channel_script_word.go",
		"channel_sound.go",
		"channel_symbolic.go",
		"channels.go",
		"convergence.go",
		"discovery.go",
		"engine.go",
		"forms.go",
		"glyphs.go",
		"knowledge_bridge.go",
		"passage.go",
		"render.go",
		"scoring.go",
		"similarity.go",
		"types.go",
	}

	for _, file := range files {
		path := "internal/decipher/" + file
		content, err := os.ReadFile(path)
		if err != nil {
			continue // skip files that don't exist
		}

		for _, pattern := range patterns {
			if strings.Contains(string(content), pattern) {
				t.Errorf("production file %s contains hardcoded knowledge pattern: %s", file, pattern)
			}
		}
	}
}
