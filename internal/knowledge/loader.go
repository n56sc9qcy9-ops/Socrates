package knowledge

import (
	"embed"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader handles loading knowledge data from YAML files.
type Loader struct {
	embedFS embed.FS
	useEmbed bool
	baseDir string
}

// NewLoader creates a new knowledge loader.
// If useEmbed is true, loads from embedded data; otherwise loads from baseDir.
func NewLoader(embedFS embed.FS, baseDir string, useEmbed bool) *Loader {
	return &Loader{
		embedFS: embedFS,
		useEmbed: useEmbed,
		baseDir: baseDir,
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

	// Load glyphs
	if err := l.loadGlyphs(kb); err != nil {
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
			Weight:     e.Weight,
		})
	}
	for _, e := range doc.Latin.Prefixes {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "latin",
			Pattern:    e.Pattern,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     e.Weight,
		})
	}
	for _, e := range doc.Latin.Suffixes {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "latin",
			Pattern:    e.Pattern,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     e.Weight,
		})
	}
	kb.AddGlyphPattern(GlyphPattern{
		Script:     "latin",
		Pattern:    "heavy",
		Concept:    doc.Latin.VowelStructure.Heavy.Concept,
		Confidence: doc.Latin.VowelStructure.Heavy.Confidence,
		Weight:     doc.Latin.VowelStructure.Heavy.Weight,
	})
	kb.AddGlyphPattern(GlyphPattern{
		Script:     "latin",
		Pattern:    "light",
		Concept:    doc.Latin.VowelStructure.Light.Concept,
		Confidence: doc.Latin.VowelStructure.Light.Confidence,
		Weight:     doc.Latin.VowelStructure.Light.Weight,
	})
	for _, e := range doc.Hebrew.Letters {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "hebrew",
			Rune:       e.Rune,
			Pattern:    e.Name,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     e.Weight,
		})
	}
	for _, e := range doc.Devanagari.SpecialChars {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "devanagari",
			Rune:       e.Rune,
			Pattern:    e.Name,
			Concept:    e.Concept,
			Confidence: e.Confidence,
			Weight:     e.Weight,
		})
	}
	for _, e := range doc.Han.Characters {
		kb.AddGlyphPattern(GlyphPattern{
			Script:     "han",
			Rune:       e.Rune,
			Pattern:    e.Pattern,
			Readings:   e.Readings,
			Confidence: e.Confidence,
			Weight:     e.Weight,
		})
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
	ScriptWords []ScriptWordEntry `yaml:"script_words"`
}

// FragmentEntry represents a fragment entry in YAML.
type FragmentEntry struct {
	Form       string `yaml:"form"`
	Concept    string `yaml:"concept"`
	Lens       string `yaml:"lens"`
	Confidence string `yaml:"confidence"`
	Weight     float64 `yaml:"weight"`
}

// ToFragment converts to internal Form type.
func (f FragmentEntry) ToFragment() Form {
	return Form{
		Form:      f.Form,
		Concept:   f.Concept,
		Lens:      f.Lens,
		Confidence: f.Confidence,
		Weight:    f.Weight,
	}
}

// ScriptWordEntry represents a script word entry in YAML.
type ScriptWordEntry struct {
	Script    string   `yaml:"script"`
	Word      string   `yaml:"word"`
	Runes     []uint32 `yaml:"runes"`
	Meanings  []string `yaml:"meanings"`
	Confidence string  `yaml:"confidence"`
	Weight    float64  `yaml:"weight"`
}

// ToScriptWord converts to internal ScriptWord type.
func (s ScriptWordEntry) ToScriptWord() ScriptWord {
	return ScriptWord{
		Script:    s.Script,
		Word:      s.Word,
		Runes:     s.Runes,
		Meanings:  s.Meanings,
		Confidence: s.Confidence,
		Weight:    s.Weight,
	}
}

// relationsDoc represents the YAML structure for relations.
type relationsDoc struct {
	Relations []RelationEntry `yaml:"relations"`
}

// RelationEntry represents a relation entry in YAML.
type RelationEntry struct {
	From   string  `yaml:"from"`
	To     string  `yaml:"to"`
	Type   string  `yaml:"type"`
	Weight float64 `yaml:"weight"`
}

// ToRelation converts to internal Relation type.
func (r RelationEntry) ToRelation() Relation {
	return Relation{
		From:   r.From,
		To:     r.To,
		Type:   r.Type,
		Weight: r.Weight,
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
	Bigrams        []LatinBigramEntry   `yaml:"bigrams"`
	Prefixes       []LatinPrefixEntry   `yaml:"prefixes"`
	Suffixes       []LatinSuffixEntry   `yaml:"suffixes"`
	VowelStructure VowelStructureEntry `yaml:"vowel_structure"`
}

type LatinBigramEntry struct {
	Pattern    string  `yaml:"pattern"`
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     float64 `yaml:"weight"`
}

type LatinPrefixEntry struct {
	Pattern    string  `yaml:"pattern"`
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     float64 `yaml:"weight"`
}

type LatinSuffixEntry struct {
	Pattern    string  `yaml:"pattern"`
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     float64 `yaml:"weight"`
}

type VowelStructureEntry struct {
	Heavy VowelStructureItem `yaml:"heavy"`
	Light VowelStructureItem `yaml:"light"`
}

type VowelStructureItem struct {
	Concept    string  `yaml:"concept"`
	Confidence string  `yaml:"confidence"`
	Weight     float64 `yaml:"weight"`
}

type HebrewGlyphs struct {
	Letters []HebrewLetterEntry `yaml:"letters"`
}

type HebrewLetterEntry struct {
	Rune       uint32 `yaml:"rune"`
	Name       string `yaml:"name"`
	Concept    string `yaml:"concept"`
	Confidence string `yaml:"confidence"`
	Weight     float64 `yaml:"weight"`
}

type DevanagariGlyphs struct {
	SpecialChars []DevanagariCharEntry `yaml:"special_chars"`
}

type DevanagariCharEntry struct {
	Rune       uint32 `yaml:"rune"`
	Name       string `yaml:"name"`
	Concept    string `yaml:"concept"`
	Confidence string `yaml:"confidence"`
	Weight     float64 `yaml:"weight"`
}

type HanGlyphs struct {
	Characters []HanCharEntry `yaml:"characters"`
}

type HanCharEntry struct {
	Rune       uint32   `yaml:"rune"`
	Pattern    string   `yaml:"pattern"`
	Readings   []string `yaml:"readings"`
	Confidence string   `yaml:"confidence"`
	Weight     float64  `yaml:"weight"`
}
