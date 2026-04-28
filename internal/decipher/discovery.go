package decipher

// discovery.go holds the deprecated monolithic file that has been split into:
// - candidate_generation.go: CandidateForm, GenerateCandidateForms, and form generation helpers
// - similarity.go: MatchEvidence, FuzzyMatchEvidence, and fuzzy matching logic
// - activation.go: ExpandConcept, AnalyzePassageTokens, and activation logic
// - convergence.go: ConvergenceResult, DetectConvergence
// - scoring.go: CalculateScoreComponents, CalculateFinalScore
// - types.go: Type definitions shared across the package
// - forms.go: Script detection and skeleton functions
