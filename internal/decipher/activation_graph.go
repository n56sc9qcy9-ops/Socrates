package decipher

import (
	"socrates/internal/knowledge"
)

// =============================================================================
// Phase 5: First-Class Activation Graph
// =============================================================================

// EvidencePath represents a single piece of evidence supporting a node or edge.
type EvidencePath struct {
	// Source information
	SourceToken string
	SourceForm  string
	SourceType  string // "direct_channel", "fuzzy_match", "passage_signal", "graph_expansion"
	// IsDirect is true if this evidence comes from direct form/glyph/script evidence,
	// false if it comes from symbolic neighbor expansion or other indirect sources.
	IsDirect    bool

	// Match information
	MatchForm  string
	MatchScore float64

	// Relation information (for edge evidence)
	RelationType string
	RelationFrom string
	RelationTo   string

	// Quality metrics
	Confidence string
	Weight     float64
}

// EvidenceID returns a deterministic identity for deduplication.
// Paths that are semantically the same should not inflate activation.
func (e EvidencePath) EvidenceID() string {
	switch e.SourceType {
	case "fuzzy_match":
		return "fuzzy|" + e.SourceToken + "|" + e.MatchForm
	case "passage_signal":
		return "passage|" + e.SourceToken + "|" + e.SourceForm + "|" + e.MatchForm
	case "graph_expansion":
		return "graph|" + e.RelationType + "|" + e.RelationFrom + "|" + e.RelationTo
	default:
		return e.SourceType + "|" + e.SourceToken + "|" + e.SourceForm
	}
}

// ActivationNode represents a concept node in the activation graph.
type ActivationNode struct {
	// Core identity
	Concept string

	// Activation metrics
	Strength     float64 // Accumulated activation strength
	BaseStrength float64 // Initial activation from direct evidence
	Confidence   string  // Highest confidence among evidence

	// Evidence tracking (deduplicated by EvidenceID)
	Evidence []EvidencePath

	// Graph position
	Depth   int  // Distance from direct evidence (0 = direct)
	Visited bool // For cycle detection during propagation
}

// EvidencePaths returns the evidence paths for this node.
func (n *ActivationNode) EvidencePaths() []EvidencePath {
	return n.Evidence
}

// AddEvidence adds an evidence path if not already present (deduplication).
// Returns true if the evidence was added, false if it was a duplicate.
func (n *ActivationNode) AddEvidence(evidence EvidencePath) bool {
	id := evidence.EvidenceID()
	for _, existing := range n.Evidence {
		if existing.EvidenceID() == id {
			return false // Duplicate
		}
	}
	n.Evidence = append(n.Evidence, evidence)
	return true
}

// ActivationEdge represents a relation edge between concepts.
type ActivationEdge struct {
	From         string
	To           string
	RelationType string
	Weight       float64 // From YAML knowledge base
	Strength     float64 // Propagated strength (with decay)
	Confidence   string
	Evidence     []EvidencePath // Evidence for this relation
}

// AddEvidence adds an evidence path to the edge.
func (e *ActivationEdge) AddEvidence(evidence EvidencePath) {
	e.Evidence = append(e.Evidence, evidence)
}

// ActivationGraph is the first-class graph structure for activation propagation.
type ActivationGraph struct {
	// Core graph data
	Nodes map[string]*ActivationNode // keyed by concept
	Edges []*ActivationEdge

	// Configuration
	MaxDepth    int     // Maximum propagation depth (default: 2)
	DecayFactor float64 // Decay per depth level (default: 0.5)

	// Metrics
	TotalStrength float64
	DirectCount   int // Number of directly activated concepts
}

// NewActivationGraph creates a new activation graph with defaults.
func NewActivationGraph() *ActivationGraph {
	return &ActivationGraph{
		Nodes:       make(map[string]*ActivationNode),
		Edges:       make([]*ActivationEdge, 0),
		MaxDepth:    2,
		DecayFactor: 0.5,
	}
}

// AddNode adds a node to the graph, merging evidence if it exists.
func (g *ActivationGraph) AddNode(concept string, strength float64, confidence string) *ActivationNode {
	if node, exists := g.Nodes[concept]; exists {
		// Merge: keep strongest activation
		if strength > node.BaseStrength {
			node.BaseStrength = strength
		}
		if strength > node.Strength {
			node.Strength = strength
		}
		// Upgrade confidence if higher
		if confidence == ConfidenceVerified && node.Confidence != ConfidenceVerified {
			node.Confidence = confidence
		}
		return node
	}

	node := &ActivationNode{
		Concept:      concept,
		Strength:     strength,
		BaseStrength: strength,
		Confidence:   confidence,
		Depth:        0,
		Evidence:     make([]EvidencePath, 0),
	}
	g.Nodes[concept] = node
	return node
}

// AddEdge adds an edge to the graph.
// Deduplicates by (from, to, relationType) - first edge wins.
func (g *ActivationGraph) AddEdge(from, to, relationType string, weight float64, confidence string) *ActivationEdge {
	// Deduplicate: check if similar edge exists
	for _, existing := range g.Edges {
		if existing.From == from && existing.To == to && existing.RelationType == relationType {
			return existing // Return existing, don't add duplicate
		}
	}

	edge := &ActivationEdge{
		From:         from,
		To:           to,
		RelationType: relationType,
		Weight:       weight,
		Strength:     weight, // Initial strength is the YAML weight
		Confidence:   confidence,
		Evidence:     make([]EvidencePath, 0),
	}
	g.Edges = append(g.Edges, edge)
	return edge
}

// GetNode returns a node by concept, or nil if not found.
func (g *ActivationGraph) GetNode(concept string) *ActivationNode {
	return g.Nodes[concept]
}

// GetEdgesFrom returns all edges originating from a concept.
func (g *ActivationGraph) GetEdgesFrom(concept string) []*ActivationEdge {
	var result []*ActivationEdge
	for _, edge := range g.Edges {
		if edge.From == concept {
			result = append(result, edge)
		}
	}
	return result
}

// GetEdgesTo returns all edges pointing to a concept.
func (g *ActivationGraph) GetEdgesTo(concept string) []*ActivationEdge {
	var result []*ActivationEdge
	for _, edge := range g.Edges {
		if edge.To == concept {
			result = append(result, edge)
		}
	}
	return result
}

// NodeCount returns the number of nodes in the graph.
func (g *ActivationGraph) NodeCount() int {
	return len(g.Nodes)
}

// EdgeCount returns the number of edges in the graph.
func (g *ActivationGraph) EdgeCount() int {
	return len(g.Edges)
}

// =============================================================================
// Graph Building from Existing Evidence
// =============================================================================

// BuildGraphFromEvidence constructs an activation graph from existing engine data.
// This replaces ad-hoc activation tracking with a first-class graph structure.
func BuildGraphFromEvidence(
	channels []ChannelResult,
	passageSignals []PassageSignal,
	fuzzyMatches []MatchEvidence,
	conceptExpansions map[string][]knowledge.DecipherConceptRelation,
	kb *knowledge.Knowledge,
) *ActivationGraph {
	g := NewActivationGraph()

	// Phase 1: Add direct channel concepts as depth-0 nodes
	conceptSet := make(map[string]bool)
	for _, ch := range channels {
		for _, sig := range ch.Signals {
			conceptSet[sig.Target] = true
		}
	}

	for concept := range conceptSet {
		// Collect evidence for this concept from channels
		var totalWeight float64
		var maxConfidence string
		// Use map to deduplicate by channel (same target via different channels = different evidence)
		channelEvidence := make(map[string]EvidencePath)

		for _, ch := range channels {
			for _, sig := range ch.Signals {
				if sig.Target == concept {
					totalWeight += sig.Weight
					// Use channel + token as key to preserve different-channel evidence
					key := ch.Name + "|" + sig.Text
					if _, exists := channelEvidence[key]; !exists {
						channelEvidence[key] = EvidencePath{
							SourceToken: sig.Text,
							SourceForm:  ch.Name + "|" + sig.Text,
							SourceType:  "direct_channel",
							Confidence:  sig.Confidence,
							Weight:      sig.Weight,
							IsDirect:    sig.IsDirect,
						}
					}
					if sig.Confidence == ConfidenceVerified {
						maxConfidence = sig.Confidence
					} else if maxConfidence == "" {
						maxConfidence = sig.Confidence
					}
				}
			}
		}

		if maxConfidence == "" {
			maxConfidence = ConfidencePlausible
		}

		node := g.AddNode(concept, totalWeight, maxConfidence)
		for _, ev := range channelEvidence {
			node.AddEvidence(ev)
		}
	}

	// Phase 2: Add passage signals as depth-0 nodes
	for _, sig := range passageSignals {
		node := g.AddNode(sig.Concept, sig.Weight, sig.Confidence)
		evidence := EvidencePath{
			SourceToken: sig.Token,
			SourceForm:  sig.Token,
			SourceType:  "passage_signal",
			MatchForm:   sig.MatchForm,
			MatchScore:  sig.MatchScore,
			Confidence:  sig.Confidence,
			Weight:      sig.Weight,
			IsDirect:    true, // Passage signals are direct evidence
		}
		node.AddEvidence(evidence)
	}

	// Phase 3: Add fuzzy match evidence to existing nodes
	for _, match := range fuzzyMatches {
		// Find the concept from the matched anchor
		anchors := GetAllAnchors(kb)
		var matchedConcept string
		for _, anchor := range anchors {
			if anchor.Form == match.AnchorForm {
				matchedConcept = anchor.Concept
				break
			}
		}

		if matchedConcept != "" {
			node := g.AddNode(matchedConcept, match.Weight, anchorConfidence(match.Weight))
			evidence := EvidencePath{
				SourceToken: match.InputForm,
				SourceForm:  match.InputForm,
				SourceType:  "fuzzy_match",
				MatchForm:   match.AnchorForm,
				MatchScore:  match.Weight,
				Confidence:  anchorConfidence(match.Weight),
				Weight:      match.Weight,
				IsDirect:    false, // Fuzzy matches are indirect
			}
			node.AddEvidence(evidence)
		}
	}

	// Phase 4: Add graph-expanded concepts from YAML relations
	// These nodes have depth 1 (not direct) so they don't inflate direct counts
	// Only new nodes get depth 1 - existing direct nodes stay direct
	for fromConcept, relations := range conceptExpansions {
		for _, rel := range relations {
			// Add the related concept if not already present
			if _, exists := g.Nodes[rel.To]; !exists {
				node := g.AddNode(rel.To, rel.Weight*0.5, ConfidencePlausible)
				node.Depth = 1 // Graph-expanded, not direct
			}
			// Note: if node already exists (was direct), don't change its depth

			// Add relation edge with deduplication
			g.AddEdge(fromConcept, rel.To, rel.Type, rel.Weight, ConfidencePlausible)
		}
	}

	// Phase 5: Reset visited flags for propagation
	g.resetVisited()

	return g
}

// anchorConfidence determines confidence from match weight.
func anchorConfidence(weight float64) string {
	if weight > 0.7 {
		return ConfidencePlausible
	}
	return ConfidenceSpeculative
}

// resetVisited resets the visited flag on all nodes.
func (g *ActivationGraph) resetVisited() {
	for _, node := range g.Nodes {
		node.Visited = false
	}
}

// =============================================================================
// Activation Propagation
// =============================================================================

// PropagateActivation propagates activation through graph edges with decay.
// Uses path-local visited tracking for deterministic cycle-safe propagation.
func (g *ActivationGraph) PropagateActivation() {
	// Propagate from each depth-0 (direct) node with a fresh path-local visited set
	for _, node := range g.Nodes {
		if node.Depth == 0 {
			visitedInPath := make(map[string]bool)
			visitedInPath[node.Concept] = true
			g.propagateFrom(node, node.Strength, 0, visitedInPath)
		}
	}

	// Recalculate total strength
	g.TotalStrength = 0
	g.DirectCount = 0
	for _, node := range g.Nodes {
		g.TotalStrength += node.Strength
		if node.Depth == 0 {
			g.DirectCount++
		}
	}
}

// propagateFrom propagates activation from a source node through its edges.
// Uses path-local visited tracking for deterministic cycle-safe propagation.
func (g *ActivationGraph) propagateFrom(source *ActivationNode, strength float64, depth int, visitedInPath map[string]bool) {
	if depth >= g.MaxDepth {
		return // Bounded depth
	}

	edges := g.GetEdgesFrom(source.Concept)
	for _, edge := range edges {
		target := g.GetNode(edge.To)
		if target == nil {
			continue
		}

		// Path-local cycle detection: prevent revisiting same node in current path
		if visitedInPath[edge.To] {
			continue
		}

		// Calculate propagated strength with decay
		decayedStrength := strength * edge.Weight * g.DecayFactor

		// Only propagate if strength is meaningful
		if decayedStrength < 0.05 {
			continue
		}

		// Add to path-local visited set
		visitedInPath[edge.To] = true

		// Update target strength
		target.Strength += decayedStrength

		// Set depth if not already set (first path wins)
		if target.Depth == 0 || depth+1 < target.Depth {
			target.Depth = depth + 1
		}

		// Recursively propagate with the path-local visited set
		g.propagateFrom(target, decayedStrength, depth+1, visitedInPath)
	}
}

// =============================================================================
// Graph Analysis
// =============================================================================

// GetTopNodes returns the top N nodes by activation strength.
func (g *ActivationGraph) GetTopNodes(n int) []*ActivationNode {
	nodes := make([]*ActivationNode, 0, len(g.Nodes))
	for _, node := range g.Nodes {
		nodes = append(nodes, node)
	}

	// Sort by strength descending
	for i := 0; i < len(nodes)-1; i++ {
		for j := i + 1; j < len(nodes); j++ {
			if nodes[j].Strength > nodes[i].Strength {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
		}
	}

	if len(nodes) > n {
		return nodes[:n]
	}
	return nodes
}

// GetCoActivationScore computes the co-activation score based on graph structure.
// Higher score means more concepts are connected through the graph.
func (g *ActivationGraph) GetCoActivationScore() float64 {
	if g.NodeCount() < 2 {
		return 0.0
	}

	// Count concepts that are connected through graph edges
	conceptSet := make(map[string]bool)
	for _, node := range g.Nodes {
		conceptSet[node.Concept] = true
	}

	var connectedPairs float64
	var totalPairs float64

	// Get top concepts for comparison
	topConcepts := g.GetTopNodes(10)
	conceptList := make([]string, 0, len(topConcepts))
	for _, n := range topConcepts {
		conceptList = append(conceptList, n.Concept)
	}

	// Check each pair for graph connectivity
	for i := 0; i < len(conceptList); i++ {
		for j := i + 1; j < len(conceptList); j++ {
			totalPairs++

			// Check if there's a graph edge between these concepts
			conceptA := conceptList[i]
			conceptB := conceptList[j]

			if g.hasConnection(conceptA, conceptB) {
				connectedPairs++
			}
		}
	}

	if totalPairs == 0 {
		return 0.0
	}

	return connectedPairs / totalPairs
}

// hasConnection checks if there's a direct graph edge between two concepts.
func (g *ActivationGraph) hasConnection(from, to string) bool {
	edges := g.GetEdgesFrom(from)
	for _, edge := range edges {
		if edge.To == to {
			return true
		}
	}
	return false
}

// GetRelationPaths returns graph-backed relation paths between top concepts.
func (g *ActivationGraph) GetRelationPaths(concepts []string) []string {
	paths := make([]string, 0)
	seen := make(map[string]bool)

	for i := 0; i < len(concepts); i++ {
		for j := i + 1; j < len(concepts); j++ {
			conceptA := concepts[i]
			conceptB := concepts[j]

			edges := g.GetEdgesFrom(conceptA)
			for _, edge := range edges {
				if edge.To == conceptB {
					path := conceptA + " --[" + edge.RelationType + "]--> " + conceptB
					if !seen[path] {
						paths = append(paths, path)
						seen[path] = true
					}
				}
			}

			// Also check reverse
			edges = g.GetEdgesFrom(conceptB)
			for _, edge := range edges {
				if edge.To == conceptA {
					path := conceptB + " --[" + edge.RelationType + "]--> " + conceptA
					if !seen[path] {
						paths = append(paths, path)
						seen[path] = true
					}
				}
			}
		}
	}

	return paths
}

// =============================================================================
// Graph-to-Convergence Bridge
// =============================================================================

// ToConvergenceResult converts the graph to a ConvergenceResult for backward compatibility.
func (g *ActivationGraph) ToConvergenceResult() ConvergenceResult {
	// Get top concepts
	topNodes := g.GetTopNodes(3)
	topConcepts := make([]ActivatedConcept, len(topNodes))
	for i, node := range topNodes {
		topConcepts[i] = ActivatedConcept{
			Concept:    node.Concept,
			Strength:   node.Strength,
			Confidence: node.Confidence,
		}
	}

	// Collect all activated concepts
	activated := make([]ActivatedConcept, 0, len(g.Nodes))
	for _, node := range g.Nodes {
		// Collect sources from evidence
		sources := make([]string, 0)
		for _, ev := range node.EvidencePaths() {
			if ev.SourceToken != "" {
				sources = append(sources, ev.SourceToken)
			}
		}
		if len(sources) == 0 {
			sources = []string{"graph"}
		}

		activated = append(activated, ActivatedConcept{
			Concept:    node.Concept,
			Strength:   node.Strength,
			Sources:    sources,
			Confidence: node.Confidence,
		})
	}

	// Get relation paths from top concepts
	conceptIDs := make([]string, len(topConcepts))
	for i, tc := range topConcepts {
		conceptIDs[i] = tc.Concept
	}
	paths := g.GetRelationPaths(conceptIDs)

	return ConvergenceResult{
		ActivatedConcepts: activated,
		CoActivationScore: g.GetCoActivationScore(),
		TopConcepts:       topConcepts,
		RelationPaths:     paths,
	}
}

// =============================================================================
// Debug Output
// =============================================================================

// GraphStats returns summary statistics for the graph.
func (g *ActivationGraph) GraphStats() map[string]interface{} {
	return map[string]interface{}{
		"node_count":     g.NodeCount(),
		"edge_count":     g.EdgeCount(),
		"total_strength": g.TotalStrength,
		"direct_count":   g.DirectCount,
		"max_depth":      g.MaxDepth,
		"decay_factor":   g.DecayFactor,
	}
}
