package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseFlowchartSimpleChain(t *testing.T) {
	out, err := Parse("flowchart LR; A-->B-->C")
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Flowchart {
		t.Errorf("Kind = %v, want Flowchart", g.Kind)
	}
	if g.Direction != ir.LeftRight {
		t.Errorf("Direction = %v, want LeftRight", g.Direction)
	}
	if len(g.Nodes) != 3 {
		t.Errorf("Nodes = %d, want 3", len(g.Nodes))
	}
	if len(g.Edges) != 2 {
		t.Errorf("Edges = %d, want 2", len(g.Edges))
	}
}

func TestParseFlowchartWithLabels(t *testing.T) {
	out, err := Parse("flowchart TD\n  A[Start] --> B{Decision}\n  B -->|Yes| C[OK]")
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	g := out.Graph
	if g.Nodes["A"].Label != "Start" {
		t.Errorf("A label = %q, want Start", g.Nodes["A"].Label)
	}
	if g.Nodes["A"].Shape != ir.Rectangle {
		t.Errorf("A shape = %v, want Rectangle", g.Nodes["A"].Shape)
	}
	if g.Nodes["B"].Shape != ir.Diamond {
		t.Errorf("B shape = %v, want Diamond", g.Nodes["B"].Shape)
	}
	if len(g.Edges) != 2 {
		t.Errorf("Edges = %d, want 2", len(g.Edges))
	}
	var labelEdge *ir.Edge
	for _, e := range g.Edges {
		if e.Label != nil && *e.Label == "Yes" {
			labelEdge = e
		}
	}
	if labelEdge == nil {
		t.Error("expected edge with label 'Yes'")
	}
}

func TestParseFlowchartDottedEdge(t *testing.T) {
	out, err := Parse("flowchart LR; A-.->B")
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(out.Graph.Edges) != 1 {
		t.Fatalf("Edges = %d, want 1", len(out.Graph.Edges))
	}
	if out.Graph.Edges[0].Style != ir.Dotted {
		t.Errorf("Style = %v, want Dotted", out.Graph.Edges[0].Style)
	}
}

func TestParseFlowchartThickEdge(t *testing.T) {
	out, err := Parse("flowchart LR; A==>B")
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.Edges[0].Style != ir.Thick {
		t.Errorf("Style = %v, want Thick", out.Graph.Edges[0].Style)
	}
}

func TestParseFlowchartBidirectional(t *testing.T) {
	out, err := Parse("flowchart LR; A<-->B")
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	e := out.Graph.Edges[0]
	if !e.ArrowStart || !e.ArrowEnd {
		t.Errorf("ArrowStart=%v ArrowEnd=%v, want both true", e.ArrowStart, e.ArrowEnd)
	}
}

func TestParseFlowchartSubgraph(t *testing.T) {
	input := "flowchart TD\n  subgraph sg1[Group]\n    A-->B\n  end\n  C-->A"
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(out.Graph.Subgraphs) != 1 {
		t.Fatalf("Subgraphs = %d, want 1", len(out.Graph.Subgraphs))
	}
	sg := out.Graph.Subgraphs[0]
	if sg.Label != "Group" {
		t.Errorf("Subgraph label = %q, want Group", sg.Label)
	}
}

func TestParseFlowchartShapes(t *testing.T) {
	tests := []struct {
		input string
		shape ir.NodeShape
	}{
		{"flowchart LR; A[rect]", ir.Rectangle},
		{"flowchart LR; A(round)", ir.RoundRect},
		{"flowchart LR; A([stadium])", ir.Stadium},
		{"flowchart LR; A{diamond}", ir.Diamond},
		{"flowchart LR; A{{hexagon}}", ir.Hexagon},
		{"flowchart LR; A[[subroutine]]", ir.Subroutine},
		{"flowchart LR; A[(cylinder)]", ir.Cylinder},
		{"flowchart LR; A((circle))", ir.DoubleCircle},
		{"flowchart LR; A>asym]", ir.Asymmetric},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			out, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error: %v", err)
			}
			if out.Graph.Nodes["A"].Shape != tt.shape {
				t.Errorf("shape = %v, want %v", out.Graph.Nodes["A"].Shape, tt.shape)
			}
		})
	}
}
