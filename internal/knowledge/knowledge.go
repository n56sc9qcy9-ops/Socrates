package knowledge

// Knowledge holds all loaded knowledge data and indexes.
// This is the main data structure returned by the loader.
type Knowledge struct {
	// Raw data
	Concepts   []Concept
	Forms      []Form
	ScriptWords []ScriptWord
	Relations  []Relation
	GlyphPatterns []GlyphPattern

	// Indexes for fast lookup
	formByText   map[string][]Form
	formsByConcept map[string][]Form
	scriptWordsByScript map[string][]ScriptWord
	conceptsByID map[string]Concept
	conceptByName map[string]Concept
	aliasesToConcept map[string]Concept
	relationsFrom map[string][]Relation
	relationsTo map[string][]Relation
	glyphPatternsByScript map[string][]GlyphPattern
	glyphPatternsByRune map[uint32][]GlyphPattern
}

// NewKnowledgeBuilder builds a Knowledge struct incrementally.
type KnowledgeBuilder struct {
	concepts      []Concept
	forms         []Form
	scriptWords   []ScriptWord
	relations     []Relation
	glyphPatterns []GlyphPattern
}

// AddConcept adds a concept to the builder.
func (b *KnowledgeBuilder) AddConcept(c Concept) {
	b.concepts = append(b.concepts, c)
}

// AddFragment adds a form/fragment to the builder.
func (b *KnowledgeBuilder) AddFragment(f Form) {
	b.forms = append(b.forms, f)
}

// AddScriptWord adds a script word to the builder.
func (b *KnowledgeBuilder) AddScriptWord(w ScriptWord) {
	b.scriptWords = append(b.scriptWords, w)
}

// AddRelation adds a relation to the builder.
func (b *KnowledgeBuilder) AddRelation(r Relation) {
	b.relations = append(b.relations, r)
}

// AddGlyphPattern adds a glyph pattern to the builder.
func (b *KnowledgeBuilder) AddGlyphPattern(g GlyphPattern) {
	b.glyphPatterns = append(b.glyphPatterns, g)
}

// NewKnowledgeBuilder creates a new builder.
func NewKnowledgeBuilder() *KnowledgeBuilder {
	return &KnowledgeBuilder{
		concepts:   make([]Concept, 0),
		forms:      make([]Form, 0),
		scriptWords: make([]ScriptWord, 0),
		relations:  make([]Relation, 0),
		glyphPatterns: make([]GlyphPattern, 0),
	}
}

// Build constructs the final Knowledge struct with indexes.
func (b *KnowledgeBuilder) Build() *Knowledge {
	kb := &Knowledge{
		Concepts:   b.concepts,
		Forms:      b.forms,
		ScriptWords: b.scriptWords,
		Relations:  b.relations,
		GlyphPatterns: b.glyphPatterns,

		// Initialize maps
		formByText:    make(map[string][]Form),
		formsByConcept: make(map[string][]Form),
		scriptWordsByScript: make(map[string][]ScriptWord),
		conceptsByID:  make(map[string]Concept),
		conceptByName: make(map[string]Concept),
		aliasesToConcept: make(map[string]Concept),
		relationsFrom: make(map[string][]Relation),
		relationsTo:   make(map[string][]Relation),
		glyphPatternsByScript: make(map[string][]GlyphPattern),
		glyphPatternsByRune: make(map[uint32][]GlyphPattern),
	}

	// Index concepts
	for _, c := range kb.Concepts {
		kb.conceptsByID[c.ID] = c
		kb.conceptByName[c.Name] = c
		for _, alias := range c.Aliases {
			kb.aliasesToConcept[alias] = c
		}
	}

	// Index forms by text and concept
	for _, f := range kb.Forms {
		kb.formByText[f.Form] = append(kb.formByText[f.Form], f)
		kb.formsByConcept[f.Concept] = append(kb.formsByConcept[f.Concept], f)
	}

	// Index script words by script
	for _, w := range kb.ScriptWords {
		kb.scriptWordsByScript[w.Script] = append(kb.scriptWordsByScript[w.Script], w)
	}

	// Index relations
	for _, r := range kb.Relations {
		kb.relationsFrom[r.From] = append(kb.relationsFrom[r.From], r)
		kb.relationsTo[r.To] = append(kb.relationsTo[r.To], r)
	}

	// Index glyph patterns
	for _, g := range kb.GlyphPatterns {
		kb.glyphPatternsByScript[g.Script] = append(kb.glyphPatternsByScript[g.Script], g)
		if g.Rune != 0 {
			kb.glyphPatternsByRune[g.Rune] = append(kb.glyphPatternsByRune[g.Rune], g)
		}
	}

	return kb
}

// ============================================================
// Core Types
// ============================================================

// Concept represents a core symbolic/conceptual anchor.
type Concept struct {
	ID        string
	Name      string
	Aliases   []string
	Neighbors []string
}

// Form represents a form-to-concept mapping (fragment).
type Form struct {
	Form      string
	Concept   string
	Lens      string
	Confidence string
	Weight    float64
}

// ScriptWord represents a complete word in a specific script.
type ScriptWord struct {
	Script    string
	Word      string
	Runes     []uint32
	Meanings  []string
	Confidence string
	Weight    float64
}

// Relation represents a directed relation between concepts.
type Relation struct {
	From   string
	To     string
	Type   string
	Weight float64
}

// GlyphPattern represents a glyph/script pattern mapping.
type GlyphPattern struct {
	Script     string
	Rune       uint32
	Pattern    string
	Concept    string
	Readings   []string
	Confidence string
	Weight     float64
}

// ============================================================
// Lookup Methods (Index-based, fast)
// ============================================================

// GetConceptByID returns a concept by its ID.
func (k *Knowledge) GetConceptByID(id string) (Concept, bool) {
	c, ok := k.conceptsByID[id]
	return c, ok
}

// GetConceptByAlias returns a concept by any of its aliases.
func (k *Knowledge) GetConceptByAlias(alias string) (Concept, bool) {
	c, ok := k.aliasesToConcept[alias]
	return c, ok
}

// GetFormsByText returns all forms matching the given text.
func (k *Knowledge) GetFormsByText(text string) []Form {
	return k.formByText[text]
}

// GetFormsByConcept returns all forms that map to a concept.
func (k *Knowledge) GetFormsByConcept(concept string) []Form {
	return k.formsByConcept[concept]
}

// GetScriptWordsByScript returns all script words for a given script.
func (k *Knowledge) GetScriptWordsByScript(script string) []ScriptWord {
	return k.scriptWordsByScript[script]
}

// GetRelationsFrom returns all relations outgoing from a concept.
func (k *Knowledge) GetRelationsFrom(concept string) []Relation {
	return k.relationsFrom[concept]
}

// GetRelationsTo returns all relations incoming to a concept.
func (k *Knowledge) GetRelationsTo(concept string) []Relation {
	return k.relationsTo[concept]
}

// GetNeighborConcepts returns all concepts that are neighbors of the given concept.
func (k *Knowledge) GetNeighborConcepts(conceptID string) []Concept {
	concept, ok := k.GetConceptByID(conceptID)
	if !ok {
		return nil
	}

	var neighbors []Concept
	seen := make(map[string]bool)
	for _, neighborID := range concept.Neighbors {
		if !seen[neighborID] {
			if neighbor, ok := k.GetConceptByID(neighborID); ok {
				neighbors = append(neighbors, neighbor)
				seen[neighborID] = true
			}
		}
	}
	return neighbors
}

// GetConceptRelations returns all relations (incoming + outgoing) for a concept.
func (k *Knowledge) GetConceptRelations(conceptID string) []Relation {
	var result []Relation
	result = append(result, k.relationsFrom[conceptID]...)
	result = append(result, k.relationsTo[conceptID]...)
	return result
}

// AllConcepts returns all loaded concepts.
func (k *Knowledge) AllConcepts() []Concept {
	return k.Concepts
}

// AllForms returns all loaded forms.
func (k *Knowledge) AllForms() []Form {
	return k.Forms
}

// AllRelations returns all loaded relations.
func (k *Knowledge) AllRelations() []Relation {
	return k.Relations
}

// GetAllFormsAsAnchors returns all forms as anchor concepts for fuzzy matching.
func (k *Knowledge) GetAllFormsAsAnchors() []AnchorConcept {
	anchors := make([]AnchorConcept, 0, len(k.Forms))
	seen := make(map[string]bool)

	for _, f := range k.Forms {
		key := f.Form + ":" + f.Concept
		if !seen[key] {
			seen[key] = true
			anchors = append(anchors, AnchorConcept{
				Form:      f.Form,
				Concept:   f.Concept,
				Confidence: f.Confidence,
				Weight:    f.Weight,
			})
		}
	}

	return anchors
}

// GetGlyphPatternsByScript returns all glyph patterns for a given script.
func (k *Knowledge) GetGlyphPatternsByScript(script string) []GlyphPattern {
	return k.glyphPatternsByScript[script]
}

// GetGlyphPatternsByRune returns all glyph patterns for a given rune.
func (k *Knowledge) GetGlyphPatternsByRune(rune uint32) []GlyphPattern {
	return k.glyphPatternsByRune[rune]
}

// ExpandConceptFromRelation expands a concept using the relation graph.
func (k *Knowledge) ExpandConceptFromRelation(conceptID string) []Relation {
	return k.relationsFrom[conceptID]
}

// AnchorConcept represents a form-to-concept mapping for anchor lookups.
type AnchorConcept struct {
	Form      string
	Concept   string
	Confidence string
	Weight    float64
}

// DecipherConceptRelation is the relation type exposed to the decipher package.
type DecipherConceptRelation struct {
	From       string
	To         string
	Type       string
	Weight     float64
	Confidence string
}

// GetConceptRelationsAsDecipher converts relations to DecipherConceptRelation format.
func (k *Knowledge) GetConceptRelationsAsDecipher(conceptID string) []DecipherConceptRelation {
	relations := k.GetConceptRelations(conceptID)
	result := make([]DecipherConceptRelation, 0, len(relations))

	for _, r := range relations {
		conf := "verified"
		if r.Weight < 0.7 {
			conf = "plausible"
		}
		if r.Weight < 0.5 {
			conf = "speculative"
		}

		result = append(result, DecipherConceptRelation{
			From:       r.From,
			To:         r.To,
			Type:       r.Type,
			Weight:     r.Weight,
			Confidence: conf,
		})
	}

	return result
}

// ExpandConcepts expands multiple concepts using the knowledge relation graph.
func (k *Knowledge) ExpandConcepts(conceptIDs []string, minWeight float64) map[string][]DecipherConceptRelation {
	result := make(map[string][]DecipherConceptRelation)

	for _, conceptID := range conceptIDs {
		allRels := k.GetConceptRelations(conceptID)
		var expansions []DecipherConceptRelation
		for _, rel := range allRels {
			if rel.Weight >= minWeight {
				conf := "verified"
				if rel.Weight < 0.7 {
					conf = "plausible"
				}
				if rel.Weight < 0.5 {
					conf = "speculative"
				}

				expansions = append(expansions, DecipherConceptRelation{
					From:       rel.From,
					To:         rel.To,
					Type:       rel.Type,
					Weight:     rel.Weight,
					Confidence: conf,
				})
			}
		}
		if len(expansions) > 0 {
			result[conceptID] = expansions
		}
	}

	return result
}
