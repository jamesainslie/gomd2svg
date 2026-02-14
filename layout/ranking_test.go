package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func edge(from, to string) *ir.Edge {
	return &ir.Edge{From: from, To: to, Directed: true, ArrowEnd: true, Style: ir.Solid}
}

func TestComputeRanksLinearChain(t *testing.T) {
	nodes := []string{"A", "B", "C"}
	edges := []*ir.Edge{edge("A", "B"), edge("B", "C")}
	ranks := computeRanks(nodes, edges, map[string]int{})
	if ranks["A"] != 0 {
		t.Errorf("A = %d, want 0", ranks["A"])
	}
	if ranks["B"] != 1 {
		t.Errorf("B = %d, want 1", ranks["B"])
	}
	if ranks["C"] != 2 {
		t.Errorf("C = %d, want 2", ranks["C"])
	}
}

func TestComputeRanksDiamond(t *testing.T) {
	nodes := []string{"A", "B", "C", "D"}
	edges := []*ir.Edge{edge("A", "B"), edge("A", "C"), edge("B", "D"), edge("C", "D")}
	ranks := computeRanks(nodes, edges, map[string]int{})
	if ranks["A"] != 0 {
		t.Errorf("A = %d, want 0", ranks["A"])
	}
	if ranks["B"] != 1 {
		t.Errorf("B = %d, want 1", ranks["B"])
	}
	if ranks["C"] != 1 {
		t.Errorf("C = %d, want 1", ranks["C"])
	}
	if ranks["D"] != 2 {
		t.Errorf("D = %d, want 2", ranks["D"])
	}
}

func TestComputeRanksCycle(t *testing.T) {
	nodes := []string{"A", "B", "C"}
	edges := []*ir.Edge{edge("A", "B"), edge("B", "C"), edge("C", "A")}
	ranks := computeRanks(nodes, edges, map[string]int{})
	if len(ranks) != 3 {
		t.Errorf("len(ranks) = %d, want 3", len(ranks))
	}
}

func TestComputeRanksDisconnected(t *testing.T) {
	nodes := []string{"A", "B", "C"}
	edges := []*ir.Edge{edge("A", "B")}
	ranks := computeRanks(nodes, edges, map[string]int{})
	if ranks["A"] != 0 {
		t.Errorf("A = %d, want 0", ranks["A"])
	}
	if ranks["B"] != 1 {
		t.Errorf("B = %d, want 1", ranks["B"])
	}
	if ranks["C"] != 0 {
		t.Errorf("C = %d, want 0 (disconnected node)", ranks["C"])
	}
}

func TestComputeRanksWithNodeOrder(t *testing.T) {
	// When cycle-breaking, nodeOrder should determine which node is picked first
	nodes := []string{"X", "Y", "Z"}
	edges := []*ir.Edge{edge("X", "Y"), edge("Y", "Z"), edge("Z", "X")}
	nodeOrder := map[string]int{"X": 0, "Y": 1, "Z": 2}
	ranks := computeRanks(nodes, edges, nodeOrder)
	// All nodes should get a rank even with cycles
	for _, n := range nodes {
		if _, ok := ranks[n]; !ok {
			t.Errorf("node %s has no rank", n)
		}
	}
}

func TestComputeRanksSingleNode(t *testing.T) {
	nodes := []string{"A"}
	var edges []*ir.Edge
	ranks := computeRanks(nodes, edges, map[string]int{})
	if ranks["A"] != 0 {
		t.Errorf("A = %d, want 0", ranks["A"])
	}
}
