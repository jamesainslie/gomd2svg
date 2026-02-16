package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestComputeLayoutSimple(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Flowchart
	graph.Direction = ir.LeftRight
	graph.EnsureNode("A", nil, nil)
	graph.EnsureNode("B", nil, nil)
	graph.EnsureNode("C", nil, nil)
	graph.Edges = []*ir.Edge{edge("A", "B"), edge("B", "C")}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.Flowchart {
		t.Errorf("Kind = %v, want Flowchart", lay.Kind)
	}
	if len(lay.Nodes) != 3 {
		t.Errorf("Nodes = %d, want 3", len(lay.Nodes))
	}
	if len(lay.Edges) != 2 {
		t.Errorf("Edges = %d, want 2", len(lay.Edges))
	}
	if lay.Width <= 0 {
		t.Errorf("Width = %f, want > 0", lay.Width)
	}
	if lay.Height <= 0 {
		t.Errorf("Height = %f, want > 0", lay.Height)
	}

	// In LR direction, nodes should be positioned left to right
	ax := lay.Nodes["A"].X
	bx := lay.Nodes["B"].X
	cx := lay.Nodes["C"].X
	if ax >= bx || bx >= cx {
		t.Errorf("expected A.x < B.x < C.x, got %f, %f, %f", ax, bx, cx)
	}
}

func TestComputeLayoutTopDown(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Flowchart
	graph.Direction = ir.TopDown
	graph.EnsureNode("A", nil, nil)
	graph.EnsureNode("B", nil, nil)
	graph.Edges = []*ir.Edge{edge("A", "B")}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	// In TD direction, A should be above B (smaller Y)
	ay := lay.Nodes["A"].Y
	by := lay.Nodes["B"].Y
	if ay >= by {
		t.Errorf("expected A.y < B.y in TopDown, got A.y=%f B.y=%f", ay, by)
	}
}

func TestComputeLayoutEdgeRouting(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Flowchart
	graph.Direction = ir.LeftRight
	graph.EnsureNode("A", nil, nil)
	graph.EnsureNode("B", nil, nil)
	graph.Edges = []*ir.Edge{edge("A", "B")}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if len(lay.Edges) != 1 {
		t.Fatalf("Edges = %d, want 1", len(lay.Edges))
	}
	edge := lay.Edges[0]
	if len(edge.Points) < 2 {
		t.Errorf("Edge points = %d, want >= 2", len(edge.Points))
	}
	if edge.From != "A" {
		t.Errorf("Edge.From = %q, want %q", edge.From, "A")
	}
	if edge.To != "B" {
		t.Errorf("Edge.To = %q, want %q", edge.To, "B")
	}
}

func TestComputeLayoutDiamondShape(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Flowchart
	graph.Direction = ir.TopDown
	diamondShape := ir.Diamond
	graph.EnsureNode("D", nil, &diamondShape)
	graph.EnsureNode("A", nil, nil)
	graph.Edges = []*ir.Edge{edge("A", "D")}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	d := lay.Nodes["D"]
	// Diamond nodes should be sized as a square (width == height)
	if d.Width != d.Height {
		t.Errorf("Diamond node Width=%f Height=%f, want square", d.Width, d.Height)
	}
}

func TestComputeLayoutNonNegativeCoordinates(t *testing.T) {
	// All Sugiyama-based layouts must produce non-negative coordinates so that
	// the SVG viewBox "0 0 W H" doesn't clip content.
	tests := []struct {
		name      string
		direction ir.Direction
		nodeCount int
	}{
		{"LR-2", ir.LeftRight, 2},
		{"TD-2", ir.TopDown, 2},
		{"LR-5", ir.LeftRight, 5},
		{"TD-5", ir.TopDown, 5},
		{"BT-3", ir.BottomTop, 3},
		{"RL-3", ir.RightLeft, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := ir.NewGraph()
			graph.Kind = ir.Flowchart
			graph.Direction = tt.direction

			// Create a chain: A->B->C->...
			ids := make([]string, tt.nodeCount)
			for i := range ids {
				ids[i] = string(rune('A' + i))
				graph.EnsureNode(ids[i], nil, nil)
			}
			for i := range len(ids) - 1 {
				graph.Edges = append(graph.Edges, edge(ids[i], ids[i+1]))
			}

			th := theme.Modern()
			cfg := config.DefaultLayout()
			lay := ComputeLayout(graph, th, cfg)

			// Every node must have non-negative bounding box.
			for id, n := range lay.Nodes {
				left := n.X - n.Width/2
				top := n.Y - n.Height/2
				if left < 0 {
					t.Errorf("node %s left edge = %f, want >= 0", id, left)
				}
				if top < 0 {
					t.Errorf("node %s top edge = %f, want >= 0", id, top)
				}
			}

			// Every edge point and label anchor must be non-negative.
			for idx, edge := range lay.Edges {
				for j, pt := range edge.Points {
					if pt[0] < 0 {
						t.Errorf("edge %d point %d X = %f, want >= 0", idx, j, pt[0])
					}
					if pt[1] < 0 {
						t.Errorf("edge %d point %d Y = %f, want >= 0", idx, j, pt[1])
					}
				}
				if edge.LabelAnchor[0] < 0 {
					t.Errorf("edge %d LabelAnchor X = %f, want >= 0", idx, edge.LabelAnchor[0])
				}
				if edge.LabelAnchor[1] < 0 {
					t.Errorf("edge %d LabelAnchor Y = %f, want >= 0", idx, edge.LabelAnchor[1])
				}
			}
		})
	}
}

func TestComputeLayoutEmptyGraph(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Flowchart
	graph.Direction = ir.LeftRight

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if len(lay.Nodes) != 0 {
		t.Errorf("Nodes = %d, want 0", len(lay.Nodes))
	}
	if len(lay.Edges) != 0 {
		t.Errorf("Edges = %d, want 0", len(lay.Edges))
	}
}
