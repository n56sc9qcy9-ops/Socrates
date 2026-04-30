# Knowledge Curation Guide

This document describes how to add new knowledge to Socrates's curated knowledge base.

## Authoritative Runtime Location

**Runtime knowledge** is embedded in `internal/knowledge/*.yaml` and loaded at build time.

**Source/reference material** in `knowledge/reference/*.yaml` is non-runtime reference material for curation.

## Knowledge Types

### 1. Concepts
Core semantic units. Each concept has:
- `id`: Unique identifier (e.g., "breath")
- `name`: Display name
- `aliases`: Alternative names that map to this concept
- `neighbors`: Related concept IDs

### 2. Forms
Mappings from text patterns to concepts:
- `form`: The text pattern (e.g., "skall")
- `concept`: Target concept ID
- `lens`: Interpretation lens
- `weight`: Reliability (0, 1]
- `confidence`: "verified", "plausible", or "speculative"

### 3. Relations
Directed links between concepts:
- `from`, `to`: Concept IDs
- `type`: "synonym", "related", "opposite", "part-of", "kind-of", or extended types
- `weight`: Strength of relation (0, 1]

### 4. Script Words
Complete words in specific scripts:
- `script`: "hebrew", "devanagari", etc.
- `word`: The written word
- `runes`: Unicode code points
- `meanings`: Concept IDs

### 5. Glyph Patterns
Script-level pattern mappings (Latin bigrams, Hebrew letters, etc.)

## Adding New Knowledge

### Step 1: Analyze the Input
Use `socrates decipher <word>` to understand what matches currently exist.

### Step 2: Validate Existing Data
Run validation:
```bash
socrates knowledge validate
```
Fix any errors or warnings before proceeding.

### Step 3: Create Suggestion for Review
For unknown inputs, the engine outputs suggestions that you can capture:
```bash
socrates decipher unknown_word > suggestions.txt
```

Review the suggestion and gather evidence from etymology, linguistics, etc.

### Step 4: Document the Addition
Create a minimum review record:
```
Type: form
Proposed Form: newword
Target Concept: existing-concept
Evidence Source: etymology reference
Confidence: plausible
Weight: 0.6
Rationale: explained in X source
Reviewed By: curator-name
```

### Step 5: Add to YAML Files
Add the curated knowledge to the appropriate file in `internal/knowledge/`:
- `concepts.yaml` for new concepts
- `forms.yaml` for new form mappings
- `relations.yaml` for new relations
- `script_words.yaml` for script words
- `glyphs.yaml` for glyph patterns

### Step 6: Re-validate
```bash
socrates knowledge validate
```
Ensure no errors.

## Validation Rules

### Hard Errors (Must Fix)
- Duplicate concept IDs
- Empty concept IDs or names
- Form/script word target concept doesn't exist (and isn't marked external/speculative)
- Relation endpoints don't exist
- Invalid weights (not in (0, 1])
- Invalid confidence levels

### Warnings (Consider Fixing)
- Unknown concept references (mark as "external" if valid)
- Extended relation types (consider using standard types)
- Non-standard patterns

## Confidence Levels

- **verified**: Confirmed by authoritative sources
- **plausible**: Supported by evidence, but not definitive
- **speculative**: Possible but requires further research

## Review Checklist

Before adding new knowledge, verify:
- [ ] Source is cited or evidence is documented
- [ ] Concept ID is unique and descriptive
- [ ] Aliases don't create ambiguity
- [ ] Weight is justified (based on evidence strength)
- [ ] Confidence level is appropriate
- [ ] No conflicts with existing knowledge
- [ ] Tests pass after addition

## Testing

Run tests after any knowledge change:
```bash
go test ./...
```

Specific validation tests:
```bash
go test ./internal/knowledge/... -v -run TestValidate
```
