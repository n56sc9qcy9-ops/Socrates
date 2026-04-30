# Training Model

Socrates should be trainable and verifiable through a supervised ranking model over evidence paths.

Training must not turn Socrates into a black-box authority. The engine should continue to expose the evidence that caused each concept, meaning-frequency field, harmonic tone, or future counsellor correction field to activate.

## Goal

The training layer learns how strongly different evidence paths should contribute to a reading.

It does not learn truth directly. It learns ranking weights that help the engine choose the best-supported fields from curated knowledge.

```text
input language
  -> generated candidates
  -> matched forms
  -> activated concepts
  -> relation paths
  -> meaning-frequency fields
  -> harmonic field
  -> ranked reading
```

The training objective is:

```text
given a curated example,
rank the expected concepts and meaning-frequency fields above weaker or false activations,
while preserving every evidence path used to make that ranking.
```

## Training Data

Use human-curated examples first. Start small and high quality.

```yaml
examples:
  - id: ex.self_forgiveness.1
    input: "I cannot forgive myself"
    expected_concepts:
      - guilt
      - self_judgment
      - forgiveness
    expected_fields:
      - mf.guilt
      - mf.self_judgment
      - mf.forgiveness
    expected_quality: unresolved_tension
    c: curated
    src: human_review
```

Training examples should be stored separately from active knowledge. They verify the engine; they are not automatically added as truth.

## Evidence Features

The first training model should tune weights for existing evidence channels:

- exact form match
- fuzzy form match
- phonetic match
- script or glyph match
- confidence label
- form weight
- relation path strength
- graph propagation depth
- passage co-activation
- harmonic profile weight
- harmonic ratio compatibility
- shared archetype
- duplicate/noise penalty
- contradiction or dissonance penalty

The model should produce feature contributions, not only a final score.

## Ranking Shape

Use a simple supervised ranker before considering neural models.

```text
score(field) =
  w_exact              * exact_match
+ w_fuzzy              * fuzzy_match
+ w_phonetic           * phonetic_match
+ w_confidence         * confidence_weight
+ w_relation           * relation_support
+ w_depth              * graph_depth_penalty
+ w_coactivation       * passage_coactivation
+ w_harmonic_profile   * harmonic_profile_support
+ w_ratio              * ratio_compatibility
+ w_archetype          * shared_archetype
- w_duplicate_noise    * duplicate_noise
- w_dissonance         * unresolved_dissonance
```

The first implementation can use deterministic evaluation with configurable weights. Later work can add weight fitting if the evaluation report proves the features are stable.

## Verification

Every training change must be evaluated against held-out examples.

Minimum metrics:

- concept precision
- concept recall
- meaning-frequency field precision
- meaning-frequency field recall
- harmonic quality agreement
- false activation count
- missed expected field count
- evidence-path validity

A result is not acceptable if it ranks the right answer for the wrong reason. The report must show which evidence paths contributed to each expected and unexpected activation.

## Commands

Target CLI:

```sh
go run ./cmd/socrates train
go run ./cmd/socrates train --examples training/examples.yaml
go run ./cmd/socrates train --examples training/examples.yaml --debug
```

Do not add automatic training mutation at first. Evaluation comes before weight fitting.

## Boundaries

- Do not train directly on prose interpretations.
- Do not let AI silently write active knowledge.
- Do not let training examples override curated knowledge validation.
- Do not hide evidence paths behind a model score.
- Do not use a neural model as the authority for truth.

AI or embeddings may later suggest candidate forms, relations, or example labels. Those suggestions must remain speculative until reviewed.
