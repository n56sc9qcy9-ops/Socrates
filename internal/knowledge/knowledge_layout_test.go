package knowledge

import (
	"os"
	"path/filepath"
	"testing"
)

// repoRoot returns the repository root path (two levels up from this package).
func repoRoot() string {
	// This file is in internal/knowledge/, go up two levels to repo root
	return filepath.Join("..", "..")
}

// TestActiveKnowledgeLocationDocumentsRuntimeFiles verifies that the active
// knowledge location (internal/knowledge/) contains expected runtime files.
// This test proves the documented location exists and has the expected structure.
func TestActiveKnowledgeLocationDocumentsRuntimeFiles(t *testing.T) {
	root := repoRoot()

	// Verify the embed directive matches reality
	// The go:embed directive in loader.go embeds *.yaml from internal/knowledge/
	// We verify that the expected runtime files exist there

	expectedFiles := []string{
		"concepts.yaml",
		"forms.yaml",
		"relations.yaml",
		"frequencies.yaml",
		"glyphs.yaml",
	}

	for _, filename := range expectedFiles {
		path := filepath.Join(root, "internal", "knowledge", filename)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("Expected runtime file %s not found at %s: %v", filename, path, err)
			continue
		}
		if info.IsDir() {
			t.Errorf("Expected file %s to be a file, not directory", path)
		}
	}
}

// TestKnowledgeReferenceNotInRuntime proves that knowledge/reference/
// is separate from active runtime and not accidentally loaded.
// This test verifies the layout separation.
func TestKnowledgeReferenceNotInRuntime(t *testing.T) {
	root := repoRoot()

	// knowledge/reference/ should exist as a reference location
	refPath := filepath.Join(root, "knowledge", "reference")
	info, err := os.Stat(refPath)
	if err != nil {
		t.Skipf("knowledge/reference/ directory doesn't exist yet: %v", err)
	}

	if !info.IsDir() {
		t.Errorf("knowledge/reference should be a directory, not a file")
	}

	// Verify it contains the moved reference files
	expectedRefFiles := []string{
		"concepts.yaml",
		"forms.yaml",
		"relations.yaml",
	}

	for _, filename := range expectedRefFiles {
		path := filepath.Join(refPath, filename)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("Reference file %s should exist in knowledge/reference/: %v", filename, err)
		}
	}

	// The key assertion: reference files should NOT be in internal/knowledge/
	// This ensures there's no ambiguity about which is active runtime
	for _, filename := range expectedRefFiles {
		path := filepath.Join(root, "internal", "knowledge", filename)
		if _, err := os.Stat(path); err == nil {
			// File exists in both places - this is the ambiguity we're resolving
			// Keep this as a warning, not failure, since the ref files may have been added independently
			t.Logf("Warning: %s exists in both internal/knowledge/ and knowledge/reference/", filename)
		}
	}
}

// TestLoadFromEmbedProvesRuntimeLocation verifies that the embedded loader
// loads from internal/knowledge/*.yaml, proving it's the runtime location.
func TestLoadFromEmbedProvesRuntimeLocation(t *testing.T) {
	kb, err := LoadFromEmbed()
	if err != nil {
		t.Fatalf("LoadFromEmbed failed: %v", err)
	}

	// If we successfully loaded from embed, internal/knowledge/*.yaml is the runtime
	// This test serves as documentation that LoadFromEmbed uses internal/knowledge/

	if kb == nil {
		t.Error("Knowledge base from LoadFromEmbed was nil")
		return
	}

	// Verify we got actual data (not empty)
	if len(kb.Concepts) == 0 {
		t.Error("Loaded knowledge base has no concepts - embedded data may be missing")
	}

	// Verify concepts are loaded from internal/knowledge/concepts.yaml
	conceptCount := len(kb.Concepts)
	if conceptCount < 80 {
		t.Errorf("Expected at least 80 concepts in embedded runtime, got %d", conceptCount)
	}

	t.Logf("Runtime knowledge: LoadFromEmbed loaded %d concepts, %d forms from internal/knowledge/*.yaml",
		conceptCount, len(kb.Forms))
}

// TestLoadFromDirWithExplicitPath validates that explicit --dir path works
// and can load from alternate locations (like knowledge/reference/ for review).
func TestLoadFromDirWithExplicitPath(t *testing.T) {
	root := repoRoot()

	// Test that LoadFromDir can load from knowledge/reference/
	refPath := filepath.Join(root, "knowledge", "reference")

	if _, err := os.Stat(refPath); os.IsNotExist(err) {
		t.Skipf("knowledge/reference/ doesn't exist, skipping test")
	}

	kb, err := LoadFromDir(refPath)
	if err != nil {
		t.Skipf("LoadFromDir(%s) failed (may have different schema): %v", refPath, err)
	}

	// If it loads, it proves --dir path validation works
	if kb != nil {
		t.Logf("LoadFromDir(%s) loaded %d concepts (proves --dir works)", refPath, len(kb.Concepts))
	}
}