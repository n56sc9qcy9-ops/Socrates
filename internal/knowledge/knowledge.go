package knowledge

// ============================================================
// Frequency Profile Types
// ============================================================

// FrequencyProfile represents a curated meaning-frequency identity.
// All values are integer-based for harmonic precision.
// Multiple language forms can converge on the same profile.
type FrequencyProfile struct {
	// MeaningFrequencyID is a stable identity (e.g., "breath-vibration").
	// Forms from different languages that share this ID converge on the same frequency.
	MeaningFrequencyID string

	// Concepts is the list of concept IDs that share this meaning-frequency.
	Concepts []string

	// Vector is a 3D tone vector [Tone1, Tone2, Tone3] as integers.
	Vector []int

	// Ratio is the foundational frequency ratio as [numerator, denominator].
	Ratio []int

	// Archetype is the reference archetype ID (e.g., "pythagorean_triple_1").
	Archetype string

	// Labels provide optional semantic labels as integers.
	Labels IntFrequencyLabels

	// Confidence indicates the profile's reliability.
	Confidence string

	// Source describes the provenance of this profile.
	Source string

	// Lens is the interpretive framework used.
	Lens string

	// Weight is the importance (0-100, integer).
	Weight int
}

// IntFrequencyLabels provides optional semantic labels as integers.
type IntFrequencyLabels struct {
	Note  int `yaml:"note"`
	Color int `yaml:"color"`
	Field int `yaml:"field"`
}

// ============================================================
// Knowledge holds all loaded knowledge data and indexes.
// ============================================================

// Knowledge holds all loaded knowledge data and indexes.
// This is the main data structure returned by the loader.
type Knowledge struct {
	// Raw data
	Concepts          []Concept
	Forms             []Form
	ScriptWords       []ScriptWord
	Relations         []Relation
	Transmutations    []TransmutationRelation
	GlyphPatterns     []GlyphPattern
	FrequencyProfiles []FrequencyProfile

	// Indexes for fast lookup
	formByText                 map[string][]Form
	formsByConcept             map[string][]Form
	scriptWordsByScript        map[string][]ScriptWord
	conceptsByID               map[string]Concept
	conceptByName              map[string]Concept
	aliasesToConcept           map[string]Concept
	relationsFrom              map[string][]Relation
	relationsTo                map[string][]Relation
	transmutationsFrom         map[string][]TransmutationRelation
	transmutationsTo           map[string][]TransmutationRelation
	glyphPatternsByScript      map[string][]GlyphPattern
	glyphPatternsByRune        map[uint32][]GlyphPattern
	frequencyProfilesByConcept map[string][]FrequencyProfile
	frequencyProfilesByID     map[string]FrequencyProfile
}

// NewKnowledgeBuilder builds a Knowledge struct incrementally.
type KnowledgeBuilder struct {
	concepts          []Concept
	forms             []Form
	scriptWords       []ScriptWord
	relations         []Relation
	transmutations    []TransmutationRelation
	glyphPatterns     []GlyphPattern
	frequencyProfiles []FrequencyProfile
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

// AddTransmutation adds a transmutation relation to the builder.
func (b *KnowledgeBuilder) AddTransmutation(t TransmutationRelation) {
	b.transmutations = append(b.transmutations, t)
}

// AddGlyphPattern adds a glyph pattern to the builder.
func (b *KnowledgeBuilder) AddGlyphPattern(g GlyphPattern) {
	b.glyphPatterns = append(b.glyphPatterns, g)
}

// AddFrequencyProfile adds a frequency profile to the builder.
func (b *KnowledgeBuilder) AddFrequencyProfile(fp FrequencyProfile) {
	b.frequencyProfiles = append(b.frequencyProfiles, fp)
}

// NewKnowledgeBuilder creates a new builder.
func NewKnowledgeBuilder() *KnowledgeBuilder {
	return &KnowledgeBuilder{
		concepts:          make([]Concept, 0),
		forms:             make([]Form, 0),
		scriptWords:       make([]ScriptWord, 0),
		relations:         make([]Relation, 0),
		transmutations:    make([]TransmutationRelation, 0),
		glyphPatterns:     make([]GlyphPattern, 0),
		frequencyProfiles: make([]FrequencyProfile, 0),
	}
}

// Build constructs the final Knowledge struct with indexes.
func (b *KnowledgeBuilder) Build() *Knowledge {
	kb := &Knowledge{
		Concepts:          b.concepts,
		Forms:             b.forms,
		ScriptWords:       b.scriptWords,
		Relations:         b.relations,
		Transmutations:    b.transmutations,
		GlyphPatterns:     b.glyphPatterns,
		FrequencyProfiles: b.frequencyProfiles,

		// Initialize maps
		formByText:                 make(map[string][]Form),
		formsByConcept:             make(map[string][]Form),
		scriptWordsByScript:        make(map[string][]ScriptWord),
		conceptsByID:               make(map[string]Concept),
		conceptByName:              make(map[string]Concept),
		aliasesToConcept:           make(map[string]Concept),
		relationsFrom:              make(map[string][]Relation),
		relationsTo:                make(map[string][]Relation),
		transmutationsFrom:         make(map[string][]TransmutationRelation),
		transmutationsTo:           make(map[string][]TransmutationRelation),
		glyphPatternsByScript:      make(map[string][]GlyphPattern),
		glyphPatternsByRune:        make(map[uint32][]GlyphPattern),
		frequencyProfilesByConcept: make(map[string][]FrequencyProfile),
		frequencyProfilesByID:      make(map[string]FrequencyProfile),
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

	// Index script words as forms by text so AnalyzePassageTokens can find them.
	// Enables non-Latin script words to flow through passage/harmonic field activation.
	for _, w := range kb.ScriptWords {
		for _, meaning := range w.Meanings {
			f := Form{
				Form:       w.Word,
				Concept:    meaning,
				Source:     w.Source,
				Confidence: w.Confidence,
				Weight:     w.Weight,
			}
			kb.formByText[w.Word] = append(kb.formByText[w.Word], f)
			kb.formsByConcept[meaning] = append(kb.formsByConcept[meaning], f)
		}
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

	// Index frequency profiles
	for _, fp := range kb.FrequencyProfiles {
		kb.frequencyProfilesByID[fp.MeaningFrequencyID] = fp
		for _, concept := range fp.Concepts {
			kb.frequencyProfilesByConcept[concept] = append(kb.frequencyProfilesByConcept[concept], fp)
		}
	}

	// Index transmutations
	kb.transmutationsFrom = make(map[string][]TransmutationRelation)
	kb.transmutationsTo = make(map[string][]TransmutationRelation)
	for _, t := range kb.Transmutations {
		kb.transmutationsFrom[t.From] = append(kb.transmutationsFrom[t.From], t)
		kb.transmutationsTo[t.To] = append(kb.transmutationsTo[t.To], t)
	}

	return kb
}

// BuildIndexes rebuilds the lookup indexes from the current data.
// Call this after directly constructing a Knowledge struct (not via KnowledgeBuilder).
func (kb *Knowledge) BuildIndexes() {
	// Initialize maps
	kb.formByText = make(map[string][]Form)
	kb.formsByConcept = make(map[string][]Form)
	kb.scriptWordsByScript = make(map[string][]ScriptWord)
	kb.conceptsByID = make(map[string]Concept)
	kb.conceptByName = make(map[string]Concept)
	kb.aliasesToConcept = make(map[string]Concept)
	kb.relationsFrom = make(map[string][]Relation)
	kb.relationsTo = make(map[string][]Relation)
	kb.transmutationsFrom = make(map[string][]TransmutationRelation)
	kb.transmutationsTo = make(map[string][]TransmutationRelation)
	kb.glyphPatternsByScript = make(map[string][]GlyphPattern)
	kb.glyphPatternsByRune = make(map[uint32][]GlyphPattern)
	kb.frequencyProfilesByConcept = make(map[string][]FrequencyProfile)
	kb.frequencyProfilesByID = make(map[string]FrequencyProfile)

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

	// Index frequency profiles
	for _, fp := range kb.FrequencyProfiles {
		kb.frequencyProfilesByID[fp.MeaningFrequencyID] = fp
		for _, concept := range fp.Concepts {
			kb.frequencyProfilesByConcept[concept] = append(kb.frequencyProfilesByConcept[concept], fp)
		}
	}
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
// Source describes provenance (e.g., "curated", "traditional", "human_review").
// Confidence describes truth/support level: "verified", "plausible", or "speculative".
type Form struct {
	Form       string
	Concept    string
	Lens       string
	Source     string
	Confidence string
	Weight     float64
}

// ScriptWord represents a complete word in a specific script.
// Source describes provenance (e.g., "curated", "traditional").
// Confidence describes truth/support level: "verified", "plausible", or "speculative".
type ScriptWord struct {
	Script     string
	Word       string
	Runes      []uint32
	Meanings   []string
	Source     string
	Confidence string
	Weight     float64
}

// Relation represents a directed relation between concepts.
// Source describes provenance.
type Relation struct {
	From   string
	To     string
	Type   string
	Source string
	Weight float64
}

// GlyphPattern represents a glyph/script pattern mapping.
// Source describes provenance (e.g., "curated", "traditional").
// Confidence describes truth/support level: "verified", "plausible", or "speculative".
type GlyphPattern struct {
	Script     string
	Rune       uint32
	Pattern    string
	Concept    string
	Readings   []string
	Source     string
	Lens       string
	Confidence string
	Weight     float64
}

// TransmutationRelation represents a data-backed counsellor/transmutation relation.
// Maps a tension/contracted field (From) to a correction field (To).
// Not fortune telling — curated relational guidance from knowledge.
type TransmutationRelation struct {
	From       string
	To         string
	Kind       string
	Confidence string
	Source     string
	Lens       string
	Weight     float64
	Notes      string
}

// ============================================================
// Lookup Methods (Index-based, fast)
// ============================================================

// HasConcept returns true if the concept exists in the knowledge base.
func (k *Knowledge) HasConcept(conceptID string) bool {
	_, ok := k.conceptsByID[conceptID]
	return ok
}

// HasFrequencyProfile returns true if the frequency profile exists.
func (k *Knowledge) HasFrequencyProfile(id string) bool {
	_, ok := k.frequencyProfilesByID[id]
	return ok
}

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

// GetTransmutationsFrom returns all transmutation relations outgoing from a concept.
func (k *Knowledge) GetTransmutationsFrom(concept string) []TransmutationRelation {
	return k.transmutationsFrom[concept]
}

// GetTransmutationsTo returns all transmutation relations incoming to a concept.
func (k *Knowledge) GetTransmutationsTo(concept string) []TransmutationRelation {
	return k.transmutationsTo[concept]
}

// AllTransmutations returns all loaded transmutation relations.
func (k *Knowledge) AllTransmutations() []TransmutationRelation {
	return k.Transmutations
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
				Form:       f.Form,
				Concept:    f.Concept,
				Confidence: f.Confidence,
				Weight:     f.Weight,
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
	Form       string
	Concept    string
	Confidence string
	Weight     float64
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

// ============================================================
// Frequency Profile Lookup Methods
// ============================================================

// GetFrequencyProfilesByConcept returns all frequency profiles for a concept.
func (k *Knowledge) GetFrequencyProfilesByConcept(concept string) []FrequencyProfile {
	return k.frequencyProfilesByConcept[concept]
}

// GetFrequencyProfileByID returns a frequency profile by its meaning-frequency ID.
func (k *Knowledge) GetFrequencyProfileByID(id string) (FrequencyProfile, bool) {
	fp, ok := k.frequencyProfilesByID[id]
	return fp, ok
}

// GetMeaningFrequencyIDForConcept returns the primary meaning-frequency ID for a concept.
// Returns empty string if no profile exists.
func (k *Knowledge) GetMeaningFrequencyIDForConcept(concept string) string {
	profiles := k.GetFrequencyProfilesByConcept(concept)
	if len(profiles) == 0 {
		return ""
	}
	// Return the profile with highest weight
	var best FrequencyProfile
	bestWeight := -1
	for _, fp := range profiles {
		if fp.Weight > bestWeight {
			bestWeight = fp.Weight
			best = fp
		}
	}
	return best.MeaningFrequencyID
}

// AllFrequencyProfiles returns all loaded frequency profiles.
func (k *Knowledge) AllFrequencyProfiles() []FrequencyProfile {
	return k.FrequencyProfiles
}

// FrequencyProfilesByConcepts returns frequency profiles for multiple concepts.
func (k *Knowledge) FrequencyProfilesByConcepts(concepts []string) []FrequencyProfile {
	seen := make(map[string]bool)
	var profiles []FrequencyProfile
	for _, concept := range concepts {
		for _, fp := range k.GetFrequencyProfilesByConcept(concept) {
			if !seen[fp.MeaningFrequencyID] {
				seen[fp.MeaningFrequencyID] = true
				profiles = append(profiles, fp)
			}
		}
	}
	return profiles
}
