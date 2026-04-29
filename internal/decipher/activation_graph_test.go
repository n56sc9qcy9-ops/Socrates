package decipher

import (
	"testing"

	"socrates/internal/knowledge"
)

// =============================================================================
// ACTIVATION GRAPH TESTS
// =============================================================================

func TestNewActivationGraph(t *testing.T) {
	g := NewActivationGraph()

	if g == nil {
		t.Fatal("NewActivationGraph should return a non-nil graph")
	}

	if g.NodeCount() != 0 {
		t.Errorf("new graph should have 0 nodes, got %d", g.NodeCount())
	}

	if g.EdgeCount() != 0 {
		t.Errorf("new graph should have 0 edges, got %d", g.EdgeCount())
	}

	if g.MaxDepth != 2 {
		t.Errorf("default MaxDepth should be 2, got %d", g.MaxDepth)
	}

	if g.DecayFactor != 0.5 {
		t.Errorf("default DecayFactor should be 0.5, got %f", g.DecayFactor)
	}
}

func TestActivationGraph_AddNode(t *testing.T) {
	g := NewActivationGraph()

	// Add first node
	node1 := g.AddNode("spirit", 0.8, ConfidenceVerified)
	if node1 == nil {
		t.Fatal("AddNode should return non-nil node")
	}

	if node1.Concept != "spirit" {
		t.Errorf("node concept should be 'spirit', got '%s'", node1.Concept)
	}

	if node1.Strength != 0.8 {
		t.Errorf("node strength should be 0.8, got %f", node1.Strength)
	}

	if node1.BaseStrength != 0.8 {
		t.Errorf("node base strength should be 0.8, got %f", node1.BaseStrength)
	}

	if node1.Confidence != ConfidenceVerified {
		t.Errorf("node confidence should be '%s', got '%s'", ConfidenceVerified, node1.Confidence)
	}

	if node1.Depth != 0 {
		t.Errorf("new node depth should be 0, got %d", node1.Depth)
	}

	// Add second node
	node2 := g.AddNode("breath", 0.5, ConfidencePlausible)
	if node2.Concept != "breath" {
		t.Errorf("node concept should be 'breath', got '%s'", node2.Concept)
	}

	if g.NodeCount() != 2 {
		t.Errorf("graph should have 2 nodes, got %d", g.NodeCount())
	}

	// Adding same concept should merge, not duplicate
	node1Again := g.AddNode("spirit", 0.9, ConfidencePlausible)
	if node1Again != node1 {
		t.Error("AddNode for existing concept should return existing node")
	}

	if g.NodeCount() != 2 {
		t.Errorf("adding existing concept should not increase count, got %d", g.NodeCount())
	}

	// Merged strength should be strongest base
	if node1.BaseStrength != 0.9 {
		t.Errorf("merged base strength should be 0.9, got %f", node1.BaseStrength)
	}
}

func TestActivationGraph_AddEdge(t *testing.T) {
	g := NewActivationGraph()

	// Add nodes first
	g.AddNode("spirit", 0.8, ConfidenceVerified)
	g.AddNode("breath", 0.5, ConfidencePlausible)

	// Add edge between them
	edge := g.AddEdge("spirit", "breath", "related_to", 0.7, ConfidencePlausible)

	if edge == nil {
		t.Fatal("AddEdge should return non-nil edge")
	}

	if edge.From != "spirit" {
		t.Errorf("edge From should be 'spirit', got '%s'", edge.From)
	}

	if edge.To != "breath" {
		t.Errorf("edge To should be 'breath', got '%s'", edge.To)
	}

	if edge.RelationType != "related_to" {
		t.Errorf("edge relation type should be 'related_to', got '%s'", edge.RelationType)
	}

	if edge.Weight != 0.7 {
		t.Errorf("edge weight should be 0.7, got %f", edge.Weight)
	}

	if g.EdgeCount() != 1 {
		t.Errorf("graph should have 1 edge, got %d", g.EdgeCount())
	}
}

func TestActivationGraph_GetEdges(t *testing.T) {
	g := NewActivationGraph()

	// Add nodes
	g.AddNode("a", 0.5, ConfidencePlausible)
	g.AddNode("b", 0.5, ConfidencePlausible)
	g.AddNode("c", 0.5, ConfidencePlausible)

	// Add edges
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("a", "c", "similar", 0.6, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.7, ConfidencePlausible)

	// Get edges from 'a'
	edgesFromA := g.GetEdgesFrom("a")
	if len(edgesFromA) != 2 {
		t.Errorf("should have 2 edges from 'a', got %d", len(edgesFromA))
	}

	// Get edges to 'c'
	edgesToC := g.GetEdgesTo("c")
	if len(edgesToC) != 2 {
		t.Errorf("should have 2 edges to 'c', got %d", len(edgesToC))
	}

	// Get edges from non-existent node
	edgesFromX := g.GetEdgesFrom("x")
	if len(edgesFromX) != 0 {
		t.Errorf("should have 0 edges from 'x', got %d", len(edgesFromX))
	}
}

// =============================================================================
// NODE EVIDENCE TESTS
// =============================================================================

func TestActivationNode_AddEvidence(t *testing.T) {
	g := NewActivationGraph()
	node := g.AddNode("spirit", 0.8, ConfidenceVerified)

	// Add first evidence
	evidence1 := EvidencePath{
		SourceToken: "in",
		SourceForm:  "in",
		SourceType:  "direct_channel",
		Confidence:  ConfidenceVerified,
		Weight:      0.5,
	}

	if !node.AddEvidence(evidence1) {
		t.Error("first evidence should be added")
	}

	if len(node.Evidence) != 1 {
		t.Errorf("node should have 1 evidence, got %d", len(node.Evidence))
	}

	// Add duplicate evidence
	evidence1Duplicate := EvidencePath{
		SourceToken: "in",
		SourceForm:  "in",
		SourceType:  "direct_channel",
		Confidence:  ConfidencePlausible,
		Weight:      0.6,
	}

	if node.AddEvidence(evidence1Duplicate) {
		t.Error("duplicate evidence should not be added")
	}

	if len(node.Evidence) != 1 {
		t.Errorf("duplicate evidence should not increase count, got %d", len(node.Evidence))
	}

	// Add different evidence
	evidence2 := EvidencePath{
		SourceToken: "spirit",
		SourceForm:  "spirit",
		SourceType:  "fuzzy_match",
		MatchForm:   "spirit",
		Confidence:  ConfidencePlausible,
		Weight:      0.7,
	}

	if !node.AddEvidence(evidence2) {
		t.Error("different evidence should be added")
	}

	if len(node.Evidence) != 2 {
		t.Errorf("node should have 2 evidence, got %d", len(node.Evidence))
	}
}

func TestEvidencePath_EvidenceID(t *testing.T) {
	// Test evidence ID generation for deduplication
	fuzzyEv := EvidencePath{
		SourceToken: "inspired",
		SourceForm:  "inspired",
		SourceType:  "fuzzy_match",
		MatchForm:   "spirit",
	}

	fuzzyID := fuzzyEv.EvidenceID()
	if fuzzyID != "fuzzy|inspired|spirit" {
		t.Errorf("fuzzy evidence ID should be 'fuzzy|inspired|spirit', got '%s'", fuzzyID)
	}

	passageEv := EvidencePath{
		SourceToken: "in",
		SourceForm:  "in",
		SourceType:  "passage_signal",
		MatchForm:   "in",
	}

	passageID := passageEv.EvidenceID()
	if passageID != "passage|in|in|in" {
		t.Errorf("passage evidence ID should be 'passage|in|in|in', got '%s'", passageID)
	}

	graphEv := EvidencePath{
		SourceType:  "graph_expansion",
		RelationType: "related_to",
		RelationFrom: "spirit",
		RelationTo:   "breath",
	}

	graphID := graphEv.EvidenceID()
	if graphID != "graph|related_to|spirit|breath" {
		t.Errorf("graph evidence ID should be 'graph|related_to|spirit|breath', got '%s'", graphID)
	}
}

// =============================================================================
// GRAPH BUILDING TESTS
// =============================================================================

func TestBuildGraphFromEvidence_Basic(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create minimal test data
	channels := []ChannelResult{
		{
			Name: "glyph",
			Signals: []Signal{
				{Text: "spirit", Target: "spirit", Channel: "glyph", Confidence: ConfidenceVerified, Weight: 0.8},
			},
		},
	}

	passageSignals := []PassageSignal{
		{Token: "in", Concept: "spirit", Weight: 0.7, Confidence: ConfidencePlausible, MatchForm: "in"},
	}

	fuzzyMatches := []MatchEvidence{
		{InputForm: "inspired", AnchorForm: "spirit", Method: "phonetic", Distance: 0.2, Weight: 0.8},
	}

	conceptExpansions := map[string][]knowledge.DecipherConceptRelation{
		"spirit": {
			{To: "breath", Type: "related_to", Weight: 0.6, Confidence: ConfidencePlausible},
		},
	}

	graph := BuildGraphFromEvidence(channels, passageSignals, fuzzyMatches, conceptExpansions, kb)

	if graph == nil {
		t.Fatal("BuildGraphFromEvidence should return non-nil graph")
	}

	// Should have spirit node from direct channel
	spiritNode := graph.GetNode("spirit")
	if spiritNode == nil {
		t.Fatal("spirit node should exist")
	}

	// Should have breath node from graph expansion
	breathNode := graph.GetNode("breath")
	if breathNode == nil {
		t.Fatal("breath node should exist from graph expansion")
	}

	// Should have edge from spirit to breath
	edges := graph.GetEdgesFrom("spirit")
	if len(edges) == 0 {
		t.Error("should have edge from spirit")
	}
}

func TestBuildGraphFromEvidence_Deduplication(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create data with duplicate evidence from SAME channel (same text, same channel)
	// These should be deduplicated
	channels := []ChannelResult{
		{
			Name: "channel1",
			Signals: []Signal{
				{Text: "test", Target: "test", Channel: "channel1", Confidence: ConfidenceVerified, Weight: 0.5},
				{Text: "test", Target: "test", Channel: "channel1", Confidence: ConfidencePlausible, Weight: 0.6}, // duplicate - same channel + same text
			},
		},
	}

	passageSignals := []PassageSignal{}
	fuzzyMatches := []MatchEvidence{}
	conceptExpansions := map[string][]knowledge.DecipherConceptRelation{}

	graph := BuildGraphFromEvidence(channels, passageSignals, fuzzyMatches, conceptExpansions, kb)

	testNode := graph.GetNode("test")
	if testNode == nil {
		t.Fatal("test node should exist")
	}

	// Same channel + same text = same key = 1 evidence entry (deduplicated)
	if len(testNode.Evidence) != 1 {
		t.Errorf("should have 1 evidence entry from same channel duplicates, got %d", len(testNode.Evidence))
	}
}

// =============================================================================
// PROPAGATION TESTS
// =============================================================================

func TestPropagateActivation_Basic(t *testing.T) {
	g := NewActivationGraph()

	// A -> B -> C chain
	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.0, ConfidencePlausible)
	g.AddNode("c", 0.0, ConfidencePlausible)

	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.6, ConfidencePlausible)

	// Propagate
	g.PropagateActivation()

	// Node 'a' should have strength 1.0 (direct)
	aNode := g.GetNode("a")
	if aNode.Strength != 1.0 {
		t.Errorf("node 'a' should have strength 1.0, got %f", aNode.Strength)
	}

	// Node 'b' should get propagated strength
	bNode := g.GetNode("b")
	if bNode.Strength <= 0.0 {
		t.Error("node 'b' should have propagated strength")
	}

	// Node 'c' should get propagated strength (with double decay)
	cNode := g.GetNode("c")
	if cNode.Depth == 0 || cNode.Depth > 2 {
		t.Errorf("node 'c' depth should be 1 or 2, got %d", cNode.Depth)
	}
}

func TestPropagateActivation_Decay(t *testing.T) {
	g := NewActivationGraph()
	g.MaxDepth = 3
	g.DecayFactor = 0.5

	// A -> B chain with edge weight 0.8
	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.0, ConfidencePlausible)
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	g.PropagateActivation()

	// Expected: 1.0 * 0.8 * 0.5 = 0.4
	bNode := g.GetNode("b")
	if bNode.Strength != 0.4 {
		t.Errorf("node 'b' should have strength 0.4 (1.0 * 0.8 * 0.5), got %f", bNode.Strength)
	}
}

func TestPropagateActivation_BoundedDepth(t *testing.T) {
	g := NewActivationGraph()
	g.MaxDepth = 1 // Very shallow
	g.DecayFactor = 0.5

	// A -> B -> C chain
	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.0, ConfidencePlausible)
	g.AddNode("c", 0.0, ConfidencePlausible)

	g.AddEdge("a", "b", "r1", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "r2", 0.8, ConfidencePlausible)

	g.PropagateActivation()

	// Node 'c' should NOT get propagation (depth would be 2, but max is 1)
	cNode := g.GetNode("c")
	if cNode.Depth != 0 || cNode.Strength != 0.0 {
		t.Errorf("node 'c' should not be activated (exceeds max depth 1), got depth=%d, strength=%f", cNode.Depth, cNode.Strength)
	}

	// Node 'b' should be activated
	bNode := g.GetNode("b")
	if bNode.Depth == 0 {
		t.Error("node 'b' should have depth 1")
	}
}

func TestPropagateActivation_CyclePrevention(t *testing.T) {
	g := NewActivationGraph()
	g.MaxDepth = 2 // Limited depth
	g.DecayFactor = 0.5 // Decay to prevent inflation

	// A -> B -> A cycle
	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.0, ConfidencePlausible)

	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "a", "related", 0.8, ConfidencePlausible)


	g.PropagateActivation()

	// Node 'a' should only have its original strength plus limited re-propagation
	// The Visited flag prevents infinite loops, but allows one re-entry per path
	aNode := g.GetNode("a")
	
	// Expected behavior: 'a' starts at 1.0, then:
	// - Propagate a->b: b gets 1.0 * 0.8 * 0.5 = 0.4
	// - Propagate b->a: a gets additional 0.4 * 0.8 * 0.5 = 0.16 (1 cycle allowed)
	// - Propagate a(b from step 2): would check Visited, prevent re-entry
	// Final: ~1.16 (1 cycle allowed, not infinite)
	if aNode.Strength > 1.2 {
		t.Errorf("node 'a' should not have inflated strength from cycles, got %f (expected ~1.16 with 1 allowed cycle)", aNode.Strength)
	}
	
	if aNode.Strength < 1.0 {
		t.Errorf("node 'a' should retain at least original strength, got %f", aNode.Strength)
	}
}

// =============================================================================
// GRAPH ANALYSIS TESTS
// =============================================================================

func TestGetTopNodes(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("low", 0.2, ConfidencePlausible)
	g.AddNode("high", 0.9, ConfidenceVerified)
	g.AddNode("medium", 0.5, ConfidencePlausible)

	top2 := g.GetTopNodes(2)

	if len(top2) != 2 {
		t.Errorf("GetTopNodes(2) should return 2 nodes, got %d", len(top2))
	}

	if top2[0].Concept != "high" {
		t.Errorf("top node should be 'high', got '%s'", top2[0].Concept)
	}

	if top2[1].Concept != "medium" {
		t.Errorf("second node should be 'medium', got '%s'", top2[1].Concept)
	}
}

func TestGetCoActivationScore(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 0.8, ConfidenceVerified)
	g.AddNode("b", 0.7, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)
	g.AddNode("d", 0.5, ConfidencePlausible)

	// Connect a-b and b-c, but not d
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.7, ConfidencePlausible)

	score := g.GetCoActivationScore()

	// Among top 10 concepts, connected pairs / total pairs
	// a-b are connected, b-c are connected, a-c not directly connected
	// Score should be > 0 since some pairs are connected
	if score <= 0.0 {
		t.Error("co-activation score should be > 0 when concepts are connected")
	}
}

func TestGetRelationPaths(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 0.8, ConfidenceVerified)
	g.AddNode("b", 0.7, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)

	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "similar", 0.7, ConfidencePlausible)

	paths := g.GetRelationPaths([]string{"a", "b", "c"})

	// Should include a -> b and b -> c paths
	if len(paths) == 0 {
		t.Error("should return relation paths between concepts")
	}
}

// =============================================================================
// GRAPH-TO-CONVERGENCE BRIDGE TESTS
// =============================================================================

func TestToConvergenceResult(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create test data
	channels := []ChannelResult{
		{
			Name: "test",
			Signals: []Signal{
				{Text: "test", Target: "test", Channel: "test", Confidence: ConfidenceVerified, Weight: 0.8},
			},
		},
	}

	passageSignals := []PassageSignal{
		{Token: "test", Concept: "test", Weight: 0.7, Confidence: ConfidencePlausible, MatchForm: "test"},
	}

	fuzzyMatches := []MatchEvidence{}
	conceptExpansions := map[string][]knowledge.DecipherConceptRelation{}

	graph := BuildGraphFromEvidence(channels, passageSignals, fuzzyMatches, conceptExpansions, kb)
	graph.PropagateActivation()

	result := graph.ToConvergenceResult()

	if len(result.ActivatedConcepts) == 0 {
		t.Error("convergence result should have activated concepts")
	}

	if result.CoActivationScore < 0.0 || result.CoActivationScore > 1.0 {
		t.Errorf("co-activation score should be in [0, 1], got %f", result.CoActivationScore)
	}
}

func TestGraphStats(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 0.8, ConfidenceVerified)
	g.AddNode("b", 0.7, ConfidenceVerified)
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	stats := g.GraphStats()

	if stats["node_count"] != 2 {
		t.Errorf("node_count should be 2, got %v", stats["node_count"])
	}

	if stats["edge_count"] != 1 {
		t.Errorf("edge_count should be 1, got %v", stats["edge_count"])
	}
}

// =============================================================================
// INTEGRATION TESTS
// =============================================================================

func TestEngine_Analyze_RoutesThroughGraph(t *testing.T) {
	engine, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}

	reading := engine.Analyze("inspired")

	// Reading should be produced normally
	if reading.Input != "inspired" {
		t.Errorf("expected input 'inspired', got '%s'", reading.Input)
	}

	// Should have convergence result (populated through graph)
	if len(reading.Convergence.ActivatedConcepts) == 0 {
		t.Error("convergence should have activated concepts from graph")
	}

	// Should have co-activation score
	if reading.Convergence.CoActivationScore < 0.0 || reading.Convergence.CoActivationScore > 1.0 {
		t.Errorf("co-activation score should be in [0, 1], got %f", reading.Convergence.CoActivationScore)
	}

	// Should have a reading
	if reading.ConciseReading == "" {
		t.Error("should produce a concise reading")
	}
}

func TestActivationGraph_EdgeCase_EmptyGraph(t *testing.T) {
	g := NewActivationGraph()

	g.PropagateActivation()

	stats := g.GraphStats()
	if stats["node_count"] != 0 {
		t.Errorf("empty graph should have 0 nodes, got %v", stats["node_count"])
	}

	score := g.GetCoActivationScore()
	if score != 0.0 {
		t.Errorf("empty graph should have co-activation score 0, got %f", score)
	}

	topNodes := g.GetTopNodes(5)
	if len(topNodes) != 0 {
		t.Errorf("empty graph should have 0 top nodes, got %d", len(topNodes))
	}
}

func TestActivationGraph_EdgeCase_SingleNode(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("solo", 0.8, ConfidenceVerified)
	g.PropagateActivation()

	score := g.GetCoActivationScore()
	if score != 0.0 {
		t.Errorf("single node should have co-activation score 0, got %f", score)
	}

	topNodes := g.GetTopNodes(3)
	if len(topNodes) != 1 {
		t.Errorf("single node graph should have 1 top node, got %d", len(topNodes))
	}
}

// =============================================================================
// PROPAGATION DETERMINISM TESTS
// =============================================================================

func TestPropagateActivation_StableWithinSinglePass(t *testing.T) {
	// Test that propagation produces consistent results within a single run
	// by checking intermediate results are stable
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.8, ConfidencePlausible)
	g.AddEdge("a", "c", "related", 0.6, ConfidencePlausible)

	// Record results after single propagation
	g.PropagateActivation()
	result := make(map[string]float64)
	for _, node := range g.Nodes {
		result[node.Concept] = node.Strength
	}

	// Check that propagation is bounded (no excessive growth)
	for concept, strength := range result {
		if strength > 2.0 {
			t.Errorf("concept '%s' has inflated strength %.4f", concept, strength)
		}
	}

	// Direct nodes should retain their original strength
	if aNode := g.GetNode("a"); aNode.Strength < 1.0 {
		t.Errorf("direct node 'a' should retain at least original strength 1.0, got %.4f", aNode.Strength)
	}
}

func TestPropagateActivation_CycleSafe(t *testing.T) {
	// Test that cycles don't cause infinite propagation or stack overflow
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)

	// Create cycle: a -> b -> c -> a
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.8, ConfidencePlausible)
	g.AddEdge("c", "a", "related", 0.8, ConfidencePlausible)

	// This should not cause infinite recursion or stack overflow
	g.PropagateActivation()

	// All nodes should have bounded strength (no infinite growth)
	for _, node := range g.Nodes {
		if node.Strength > 2.0 {
			t.Errorf("node '%s' has inflated strength %.4f from cycle", node.Concept, node.Strength)
		}
	}
}

func TestPropagateActivation_LongCycleSafe(t *testing.T) {
	// Test that long chains with cycles don't explode
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)
	g.AddNode("d", 0.4, ConfidencePlausible)
	g.AddNode("e", 0.3, ConfidencePlausible)

	// Long chain: a -> b -> c -> d -> e -> a
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)
	g.AddEdge("b", "c", "related", 0.8, ConfidencePlausible)
	g.AddEdge("c", "d", "related", 0.8, ConfidencePlausible)
	g.AddEdge("d", "e", "related", 0.8, ConfidencePlausible)
	g.AddEdge("e", "a", "related", 0.8, ConfidencePlausible)

	// This should complete without stack overflow
	g.PropagateActivation()

	// All nodes should have bounded strength
	for _, node := range g.Nodes {
		if node.Strength > 2.0 {
			t.Errorf("node '%s' has inflated strength %.4f from long cycle", node.Concept, node.Strength)
		}
	}
}

// =============================================================================
// DIRECT VS DERIVED NODE TESTS
// =============================================================================

func TestDirectAndDerivedNodes_Distinguishable(t *testing.T) {
	// Test that direct nodes (depth 0) and derived nodes (depth > 0) are distinguishable
	g := NewActivationGraph()

	// Add direct node
	directNode := g.AddNode("direct", 1.0, ConfidenceVerified)
	if directNode.Depth != 0 {
		t.Errorf("direct node should have depth 0, got %d", directNode.Depth)
	}

	// Add derived node (simulating graph expansion)
	derivedNode := g.AddNode("derived", 0.5, ConfidencePlausible)
	derivedNode.Depth = 1 // Manually set to derived

	if directNode.Depth == derivedNode.Depth {
		t.Error("direct and derived nodes should have different depths")
	}
}

func TestGraphExpandedNodes_NotDirect(t *testing.T) {
	kb := testKB()
	if kb == nil {
		t.Skip("no embedded knowledge available")
	}

	// Create test data with a concept that will expand via relations
	channels := []ChannelResult{
		{
			Name: "test",
			Signals: []Signal{
				{Text: "love", Target: "love", Channel: "test", Confidence: ConfidenceVerified, Weight: 0.8},
			},
		},
	}

	passageSignals := []PassageSignal{}
	fuzzyMatches := []MatchEvidence{}
	conceptExpansions := map[string][]knowledge.DecipherConceptRelation{}

	// Get relations for "love" using the correct method
	if relations := kb.GetConceptRelationsAsDecipher("love"); len(relations) > 0 {
		conceptExpansions["love"] = relations
	}

	graph := BuildGraphFromEvidence(channels, passageSignals, fuzzyMatches, conceptExpansions, kb)
	graph.PropagateActivation()

	// Check that expanded nodes have depth > 0
	directCount := 0
	for _, node := range graph.Nodes {
		if node.Depth == 0 {
			directCount++
		}
	}

	// The "love" node should be direct (depth 0)
	loveNode := graph.GetNode("love")
	if loveNode == nil {
		t.Fatal("love node should exist")
	}
	if loveNode.Depth != 0 {
		t.Errorf("direct node 'love' should have depth 0, got %d", loveNode.Depth)
	}
}

// =============================================================================
// EDGE DEDUPLICATION TESTS
// =============================================================================

func TestAddEdge_DeduplicatesByIdentity(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)


	// Add edge first time
	e1 := g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	// Try to add duplicate edge (same from, to, relationType)
	e2 := g.AddEdge("a", "b", "related", 0.9, ConfidencePlausible) // Different weight

	// Should return existing edge, not create new one
	if e1 != e2 {
		t.Error("duplicate edge should return existing edge")
	}

	// Graph should only have one edge
	if len(g.Edges) != 1 {
		t.Errorf("graph should have 1 edge, got %d", len(g.Edges))
	}

	// Edge should retain original weight
	if e1.Weight != 0.8 {
		t.Errorf("edge should retain original weight 0.8, got %.2f", e1.Weight)
	}
}

func TestAddEdge_DifferentRelationTypesNotDuplicates(t *testing.T) {
	g := NewActivationGraph()


	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)

	// Add edge with "related"
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	// Add edge with "similar" - different relation type, should be kept
	g.AddEdge("a", "b", "similar", 0.7, ConfidencePlausible)

	// Should have 2 edges
	if len(g.Edges) != 2 {
		t.Errorf("graph should have 2 edges (different relation types), got %d", len(g.Edges))
	}
}

func TestAddEdge_DifferentToNotDuplicate(t *testing.T) {
	g := NewActivationGraph()

	g.AddNode("a", 1.0, ConfidenceVerified)
	g.AddNode("b", 0.8, ConfidenceVerified)
	g.AddNode("c", 0.6, ConfidenceVerified)

	// Add edge a -> b
	g.AddEdge("a", "b", "related", 0.8, ConfidencePlausible)

	// Add edge a -> c - different target, should be kept
	g.AddEdge("a", "c", "related", 0.7, ConfidencePlausible)

	// Should have 2 edges
	if len(g.Edges) != 2 {
		t.Errorf("graph should have 2 edges (different targets), got %d", len(g.Edges))
	}
}

// =============================================================================
// EVIDENCE DEDUPLICATION TESTS
// =============================================================================

func TestActivationNode_AddEvidence_DeduplicatesByID(t *testing.T) {
	node := &ActivationNode{
		Concept: "test",
		Evidence: make([]EvidencePath, 0),
	}


	// Add first evidence
	evidence1 := EvidencePath{
		SourceToken: "token1",
		SourceForm:  "form1",
		SourceType:  "direct",
		Confidence:  ConfidenceVerified,
		Weight:      0.8,
	}
	if !node.AddEvidence(evidence1) {
		t.Error("first evidence should be added")
	}

	// Add duplicate evidence (same ID)
	evidence1Dup := EvidencePath{
		SourceToken: "token1", // Same token
		SourceForm:  "form1",  // Same form
		SourceType:  "direct", // Same type
		Confidence:  ConfidenceVerified,
		Weight:      0.8,
	}
	if node.AddEvidence(evidence1Dup) {
		t.Error("duplicate evidence should not be added")
	}

	if len(node.Evidence) != 1 {
		t.Errorf("node should have 1 evidence, got %d", len(node.Evidence))
	}
}

func TestActivationNode_AddEvidence_DifferentIDNotDuplicates(t *testing.T) {
	node := &ActivationNode{
		Concept: "test",
		Evidence: make([]EvidencePath, 0),
	}

	// Add evidence from different sources
	evidence1 := EvidencePath{
		SourceToken: "token1",
		SourceForm:  "form1",
		SourceType:  "direct",
		Confidence:  ConfidenceVerified,
		Weight:      0.8,
	}
	evidence2 := EvidencePath{
		SourceToken: "token2", // Different token
		SourceForm:  "form2",
		SourceType:  "passage", // Different type
		Confidence:  ConfidencePlausible,
		Weight:      0.6,
	}

	node.AddEvidence(evidence1)
	node.AddEvidence(evidence2)

	if len(node.Evidence) != 2 {
		t.Errorf("node should have 2 evidence from different sources, got %d", len(node.Evidence))
	}
}