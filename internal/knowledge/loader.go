package knowledge

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader handles loading knowledge data from YAML files.
type Loader struct {
	embedFS  embed.FS
	useEmbed bool
	baseDir  string
}

// NewLoader creates a new knowledge loader.
// If useEmbed is true, loads from embedded data; otherwise loads from baseDir.
func NewLoader(embedFS embed.FS, baseDir string, useEmbed bool) *Loader {
	return &Loader{
		embedFS:  embedFS,
		useEmbed: useEmbed,
		baseDir:  baseDir,
	}
}

// LoadAll loads all knowledge data files and returns a Knowledge struct.
// Returns error if any file fails to load or parse.
func (l *Loader) LoadAll() (*Knowledge, error) {
	kb := NewKnowledgeBuilder()

	// Load concepts
	if err := l.loadConcepts(kb); err != nil {
		return nil, err
	}

	// Load forms
	if err := l.loadForms(kb); err != nil {
		return nil, err
	}

	// Load relations
	if err := l.loadRelations(kb); err != nil {
		return nil, err
	}

	// Load transmute (counsellor/transmutation field relations)
	if err := l.loadTransmutations(kb); err != nil {
		return nil, err
	}

	// Load glyphs
	if err := l.loadGlyphs(kb); err != nil {
		return nil, err
	}

	// Load frequency profiles
	if err := l.loadFrequencies(kb); err != nil {
		return nil, err
	}

	return kb.Build(), nil
}

// loadConcepts loads concepts.yaml and populates the knowledge builder.
func (l *Loader) loadConcepts(kb *KnowledgeBuilder) error {
	data, err := l.readFile("concepts.yaml")
	if err != nil {
		return err
	}

	var doc conceptsDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}

	for _, c := range doc.Concepts {
		kb.AddConcept(c.ToConcept())
	}

	return nil
}

// loadForms loads forms.yaml and populates the knowledge builder.
func (l *Loader) loadForms(kb *KnowledgeBuilder) error {
	data, err := l.readFile("forms.yaml")
	if err != nil {
		return err
	}

	var doc formsDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}

	for _, f := range doc.Fragments {
		kb.AddFragment(f.ToFragment())
	}

	for _, f := range doc.Forms {
		kb.AddFragment(f.ToForm())
	}

	for _, w := range doc.ScriptWords {
		kb.AddScriptWord(w.ToScriptWord())
	}

	return nil
}

// loadRelations loads relations.yaml and populates the knowledge builder.
func (l *Loader) loadRelations(kb *KnowledgeBuilder) error {
	data, err := l.readFile("relations.yaml")
	if err != nil {
		return err
	}

	var doc relationsDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}

	for _, r := range doc.Relations {
		kb.AddRelation(r.ToRelation())
	}

	return nil
}

// loadTransmutations loads transmute.yaml for counsellor/transmutation fields.
func (l *Loader) loadTransmutations(kb *KnowledgeBuilder) error {
	data, err := l.readFile("transmute.yaml")
	if err != nil {
		// File may not exist in all builds — not an error
		return nil
	}

	var doc transmuteDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("failed to parse transmute.yaml: %w", err)
	}

	for _, e := range doc.TransmuteEntries {
		kb.AddTransmutation(e.ToTransmutation())
	}

	return nil
}

// loadGlyphs loads glyphs.yaml and populates the knowledge builder.
func (l *Loader) loadGlyphs(kb *KnowledgeBuilder) error {
	data, err := l.readFile("glyphs.yaml")
	if err != nil {
		return err
	}

	var doc glyphsDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}

	for _, e := range doc.Latin.Bigrams {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "latin",
			Pattern:    e.Pattern,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     float64(e.Weight),
		})
	}
	for _, e := range doc.Latin.Prefixes {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "latin",
			Pattern:    e.Pattern,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     float64(e.Weight),
		})
	}
	for _, e := range doc.Latin.Suffixes {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "latin",
			Pattern:    e.Pattern,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     float64(e.Weight),
		})
	}
	kb.AddGlyphPattern(GlyphPattern{
		Script:     "latin",
		Pattern:    "heavy",
		Concept:    doc.Latin.VowelStructure.Heavy.Concept,
		Confidence: doc.Latin.VowelStructure.Heavy.Confidence,
		Weight:     float64(doc.Latin.VowelStructure.Heavy.Weight),
		Lens:       "vowel-structure",
	})
	kb.AddGlyphPattern(GlyphPattern{
		Script:     "latin",
		Pattern:    "light",
		Concept:    doc.Latin.VowelStructure.Light.Concept,
		Confidence: doc.Latin.VowelStructure.Light.Confidence,
		Weight:     float64(doc.Latin.VowelStructure.Light.Weight),
		Lens:       "vowel-structure",
	})
	for _, e := range doc.Latin.RepeatedLetters {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "latin",
			Pattern:    e.Pattern,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     float64(e.Weight),
			Lens:       "repetition",
		})
	}
	for _, e := range doc.Hebrew.Letters {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "hebrew",
			Rune:       e.Rune,
			Pattern:    e.Name,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     float64(e.Weight),
		})
	}
	for _, e := range doc.Devanagari.SpecialChars {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "devanagari",
			Rune:       e.Rune,
			Pattern:    e.Name,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     float64(e.Weight),
		})
	}
	for _, e := range doc.Han.Characters {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "han",
			Rune:       e.Rune,
			Pattern:    e.Pattern,
			Readings:   e.Readings,
			Confidence: e.Confidence,
			Weight:     float64(e.Weight),
		})
	}

	return nil
}

// loadFrequencies loads frequencies.yaml and populates the knowledge builder.
// CRITICAL: This function enforces strict parsing of frequency profiles.
// Float harmonic meaning fields (frequency_hz, pitch, color_rgb, em_band_id, etc.)
// are NOT modeled and must be rejected to prevent silent data model corruption.
// Use scanForUnknownFreqFields to detect these fields before unmarshaling.
func (l *Loader) loadFrequencies(kb *KnowledgeBuilder) error {
	data, err := l.readFile("frequencies.yaml")
	if err != nil {
		// frequencies.yaml is optional - if it doesn't exist, just return
		return nil
	}

	// Strict check: reject float harmonic meaning fields that aren't modeled.
	// Fields like frequency_hz: 528.0, pitch: 432.0, color_rgb: [...], or
	// em_band_id: ... must not silently pass through YAML unmarshaling.
	if unknown, err := hasUnknownFreqFields(data); err == nil && len(unknown) > 0 {
		return fmt.Errorf("frequencies.yaml contains unknown harmonic identity fields: %v; model these explicitly or remove", unknown)
	}

	var doc frequenciesDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}

	for _, e := range doc.FrequencyProfiles {
		kb.AddFrequencyProfile(e.ToFrequencyProfile())
	}

	return nil
}

// readFile reads a file from embedded FS or filesystem.
func (l *Loader) readFile(filename string) ([]byte, error) {
	if l.useEmbed {
		return l.embedFS.ReadFile(filename)
	}
	path := filename
	if l.baseDir != "" {
		path = filepath.Join(l.baseDir, filename)
	}
	return os.ReadFile(path)
}

//go:embed *.yaml
var embedFS embed.FS

// LoadFromEmbed is a convenience function that loads from embedded data.
func LoadFromEmbed() (*Knowledge, error) {
	loader := NewLoader(embedFS, "", true)
	return loader.LoadAll()
}

// LoadFromDir loads knowledge from a directory on the filesystem.
func LoadFromDir(baseDir string) (*Knowledge, error) {
	loader := NewLoader(embed.FS{}, baseDir, false)
	return loader.LoadAll()
}

// LoadOrPanic loads knowledge or panics on error.
func LoadOrPanic() *Knowledge {
	kb, err := LoadFromEmbed()
	if err != nil {
		panic("failed to load knowledge: " + err.Error())
	}
	return kb
}

// ============================================================
// YAML Document Structures
// ============================================================

// conceptsDoc represents the YAML structure for concepts.
type conceptsDoc struct {
	Concepts []ConceptEntry `yaml:"concepts"`
}

// ConceptEntry represents a single concept in YAML.
type ConceptEntry struct {
	ID        string   `yaml:"id"`
	Name      string   `yaml:"name"`
	Aliases   []string `yaml:"aliases"`
	Neighbors []string `yaml:"neighbors"`
}

// ToConcept converts to internal Concept type.
func (c ConceptEntry) ToConcept() Concept {
	return Concept{
		ID:        c.ID,
		Name:      c.Name,
		Aliases:   c.Aliases,
		Neighbors: c.Neighbors,
	}
}

// formsDoc represents the YAML structure for forms.
type formsDoc struct {
	Fragments   []FragmentEntry   `yaml:"fragments"`
	Forms       []FormEntry       `yaml:"forms"`
	ScriptWords []ScriptWordEntry `yaml:"script_words"`
}

// FormEntry represents a form entry in YAML (whole-word matches).
type FormEntry struct {
	Form       string  `yaml:"form"`
	Script     string  `yaml:"script"`
	Concept    string  `yaml:"concept"`
	Source     string  `yaml:"source"`
	Confidence string  `yaml:"confidence"`
	Weight     int `yaml:"weight"`
}

// ToForm converts to internal Form type.
func (f FormEntry) ToForm() Form {
	return Form{
		Form:       f.Form,
		Concept:    f.Concept,
		Lens:       f.Script, // script field used as lens
		Source:     f.Source,
		Confidence: f.Confidence,
		Weight:     float64(f.Weight),
	}
}

// FragmentEntry represents a fragment entry in YAML.
type FragmentEntry struct {
	Form       string  `yaml:"form"`
	Concept    string  `yaml:"concept"`
	Lens       string  `yaml:"lens"`
	Source     string  `yaml:"source"`
	Confidence string  `yaml:"confidence"`
	Weight     int `yaml:"weight"`
}

// ToFragment converts to internal Form type.
func (f FragmentEntry) ToFragment() Form {
	return Form{
		Form:       f.Form,
		Concept:    f.Concept,
		Lens:       f.Lens,
		Source:     f.Source,
		Confidence: f.Confidence,
		Weight:     float64(f.Weight),
	}
}

// ScriptWordEntry represents a script word entry in YAML.
type ScriptWordEntry struct {
	Script     string   `yaml:"script"`
	Word       string   `yaml:"word"`
	Runes      []uint32 `yaml:"runes"`
	Meanings   []string `yaml:"meanings"`
	Source     string   `yaml:"source"`
	Confidence string   `yaml:"confidence"`
	Weight     int  `yaml:"weight"`
}

// ToScriptWord converts to internal ScriptWord type.
func (s ScriptWordEntry) ToScriptWord() ScriptWord {
	return ScriptWord{
		Script:     s.Script,
		Word:       s.Word,
		Runes:      s.Runes,
		Meanings:   s.Meanings,
		Source:     s.Source,
		Confidence: s.Confidence,
		Weight:     float64(s.Weight),
	}
}

// relationsDoc represents the YAML structure for relations.
type relationsDoc struct {
	Relations []RelationEntry `yaml:"relations"`
}

// transmuteDoc represents the YAML structure for counsellor/transmutation fields.
type transmuteDoc struct {
	TransmuteEntries []TransmuteEntry `yaml:"transmute_entries"`
}

// TransmuteEntry represents a transmutation entry in YAML.
type TransmuteEntry struct {
	From       string  `yaml:"from"`
	To         string  `yaml:"to"`
	Kind       string  `yaml:"kind"`
	Confidence string  `yaml:"confidence"`
	Source     string  `yaml:"source"`
	Lens       string  `yaml:"lens"`
	Weight     int `yaml:"weight"`
	Notes      string  `yaml:"notes"`
}

// ToTransmutation converts to decipher.TransmutationRelation type.
func (e TransmuteEntry) ToTransmutation() TransmutationRelation {
	return TransmutationRelation{
		From:       e.From,
		To:         e.To,
		Kind:       e.Kind,
		Confidence: e.Confidence,
		Source:     e.Source,
		Lens:       e.Lens,
		Weight:     float64(e.Weight),
		Notes:      e.Notes,
	}
}

// RelationEntry represents a relation entry in YAML.
type RelationEntry struct {
	From   string  `yaml:"from"`
	To     string  `yaml:"to"`
	Type   string  `yaml:"type"`
	Source string  `yaml:"source"`
	Weight int `yaml:"weight"`
}

// ToRelation converts to internal Relation type.
func (r RelationEntry) ToRelation() Relation {
	return Relation{
		From:   r.From,
		To:     r.To,
		Type:   r.Type,
		Source: r.Source,
		Weight: float64(r.Weight),
	}
}

// glyphsDoc represents the YAML structure for glyph patterns.
type glyphsDoc struct {
	Latin      LatinGlyphs      `yaml:"latin"`
	Hebrew     HebrewGlyphs     `yaml:"hebrew"`
	Devanagari DevanagariGlyphs `yaml:"devanagari"`
	Han        HanGlyphs        `yaml:"han"`
}

type LatinGlyphs struct {
	Bigrams         []LatinBigramEntry   `yaml:"bigrams"`
	Prefixes        []LatinPrefixEntry   `yaml:"prefixes"`
	Suffixes        []LatinSuffixEntry   `yaml:"suffixes"`
	VowelStructure  VowelStructureEntry  `yaml:"vowel_structure"`
	RepeatedLetters []LatinRepeatedEntry `yaml:"repeated_letters"`
}

type LatinBigramEntry struct {
	Pattern    string  `yaml:"pattern"`
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     int `yaml:"weight"`
}

type LatinPrefixEntry struct {
	Pattern    string  `yaml:"pattern"`
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     int `yaml:"weight"`
}

type LatinSuffixEntry struct {
	Pattern    string  `yaml:"pattern"`
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     int `yaml:"weight"`
}

type LatinRepeatedEntry struct {
	Pattern    string  `yaml:"pattern"`
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     int `yaml:"weight"`
	Lens       string  `yaml:"lens"`
}

type VowelStructureEntry struct {
	Heavy VowelStructureItem `yaml:"heavy"`
	Light VowelStructureItem `yaml:"light"`
}

type VowelStructureItem struct {
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     int `yaml:"weight"`
	Lens       string  `yaml:"lens"`
}

type HebrewGlyphs struct {
	Letters []HebrewLetterEntry `yaml:"letters"`
}

type HebrewLetterEntry struct {
	Rune       uint32  `yaml:"rune"`
	Name       string  `yaml:"name"`
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     int `yaml:"weight"`
}

type DevanagariGlyphs struct {
	SpecialChars []DevanagariCharEntry `yaml:"special_chars"`
}

type DevanagariCharEntry struct {
	Rune       uint32  `yaml:"rune"`
	Name       string  `yaml:"name"`
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     int `yaml:"weight"`
}

type HanGlyphs struct {
	Characters []HanCharEntry `yaml:"characters"`
}

type HanCharEntry struct {
	Rune       uint32   `yaml:"rune"`
	Pattern    string   `yaml:"pattern"`
	Readings   []string `yaml:"readings"`
	Confidence string   `yaml:"confidence"`
	Weight     int  `yaml:"weight"`
}

// ============================================================
// Frequency Profile YAML Structures
// ============================================================

// frequenciesDoc represents the YAML structure for frequency profiles.
type frequenciesDoc struct {
	FrequencyProfiles []FrequencyProfileEntry `yaml:"frequency_profiles"`
}

// FrequencyProfileEntry represents a frequency profile entry in YAML.
// All values are integer-based for harmonic precision.
// Source describes provenance (e.g., "curated", "human_review", "physics").
// Confidence describes truth/support level: "verified", "plausible", or "speculative".
type FrequencyProfileEntry struct {
	MeaningFrequencyID string        `yaml:"meaning_frequency_id"`
	Concepts           []string      `yaml:"concepts"`
	Vector             []int         `yaml:"vector"`
	Ratio              []int         `yaml:"ratio"`
	Archetype          string        `yaml:"archetype"`
	Labels             YAMLIntLabels `yaml:"labels"`
	Confidence         string        `yaml:"confidence"`
	Source             string        `yaml:"source"`
	Lens               string        `yaml:"lens"`
	Weight             int           `yaml:"weight"`
}

// YAMLIntLabels is used for YAML unmarshaling of integer labels.
type YAMLIntLabels struct {
	Note  int `yaml:"note"`
	Color int `yaml:"color"`
	Field int `yaml:"field"`
}

// ToFrequencyProfile converts a YAML entry to internal FrequencyProfile.
func (e FrequencyProfileEntry) ToFrequencyProfile() FrequencyProfile {
	labels := IntFrequencyLabels{
		Note:  e.Labels.Note,
		Color: e.Labels.Color,
		Field: e.Labels.Field,
	}
	return FrequencyProfile{
		MeaningFrequencyID: e.MeaningFrequencyID,
		Concepts:           e.Concepts,
		Vector:             e.Vector,
		Ratio:              e.Ratio,
		Archetype:          e.Archetype,
		Labels:             labels,
		Confidence:         e.Confidence,
		Source:             e.Source,
		Lens:               e.Lens,
		Weight:             e.Weight,
	}
}

// ============================================================
// Strict Frequency Field Detection
// ============================================================
//
// The harmonic meaning-frequency model is strictly integer-based.
// Float values (frequency_hz, pitch, color_rgb, em_band_id) are NOT modeled
// and must NOT silently pass through the YAML loader.
//
// This section implements a pre-unmarshaling scan that detects these fields
// and returns an error before any silent data corruption can occur.
//
// CRITICAL: YAML silently ignores unknown fields when unmarshaling into a
// typed struct. A field like `frequency_hz: 528.0` in the YAML would simply
// be dropped by yaml.Unmarshal, making it impossible to detect after the fact.
// Therefore we scan the raw YAML before unmarshaling.

// floatHarmonicKeys are unmodeled float-meaning fields that must be rejected.
// These represent physical float measurements that do NOT belong in the
// integer harmonic meaning-frequency substrate.
var floatHarmonicKeys = map[string]bool{
	"frequency_hz":    true,
	"frequency_hertz": true,
	"pitch":           true,
	"color_rgb":       true,
	"rgb":             true,
	"rgba":            true,
	"hex_color":       true,
	"colour":          true,
	"em_band_id":      true,
	"em_frequency":    true,
	"em_wavelength":   true,
	"wave_length":     true,
	"wavelength_nm":   true,
	"wavelength_hz":   true,
	"wavelength":      true,
}

// structuralDocKeys are document-level structural keys (not profile fields).
var structuralDocKeys = map[string]bool{
	"frequency_profiles": true,
	"harmonic_systems":   true,
}

// knownFreqProfileKeys are modeled fields within a frequency profile entry.
var knownFreqProfileKeys = map[string]bool{
	"meaning_frequency_id": true,
	"concepts":             true,
	"vector":               true,
	"ratio":                true,
	"archetype":            true,
	"labels":               true,
	"confidence":           true,
	"source":               true,
	"lens":                 true,
	"weight":               true,
	// labels sub-fields
	"note":  true,
	"color": true,
	"field": true,
}

// hasUnknownFreqFields scans YAML data for unmodeled float-harmonic fields.
// Returns field paths (e.g., "frequency_profiles[0].frequency_hz") for any
// float-harmonic fields found inside frequency profile entries.
// Returns nil if no unmodeled fields are found (valid data).
func hasUnknownFreqFields(data []byte) ([]string, error) {
	var raw interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return scanForUnknownFreqFields(raw, ""), nil
}

// scanForUnknownFreqFields recursively scans a parsed YAML node.
// path is the dot-separated path to the current node.
// Returns paths to any unmodeled float-harmonic fields found.
func scanForUnknownFreqFields(node interface{}, path string) []string {
	var unknown []string

	switch v := node.(type) {
	case map[string]interface{}:
		for key, val := range v {
			isStructural := structuralDocKeys[key]
			isKnownField := knownFreqProfileKeys[key]
			isFloatHarmonic := floatHarmonicKeys[key]

			var fullPath string
			if path == "" {
				fullPath = key
			} else {
				fullPath = path + "." + key
			}

			if isStructural && path == "" {
				// Top-level structural key (e.g., "frequency_profiles")
				unknown = append(unknown, scanForUnknownFreqFields(val, fullPath)...)
			} else if isKnownField {
				// Known field - recurse but don't flag
				unknown = append(unknown, scanForUnknownFreqFields(val, fullPath)...)
			} else if isFloatHarmonic {
				// Unmodeled float-harmonic field - this is an error
				unknown = append(unknown, fullPath)
			} else {
				// Unknown key - recurse to find nested float-harmonic fields
				unknown = append(unknown, scanForUnknownFreqFields(val, fullPath)...)
			}
		}
	case []interface{}:
		for i, item := range v {
			unknown = append(unknown, scanForUnknownFreqFields(item, path+"["+itoa(i)+"]")...)
		}
	}

	return unknown
}
