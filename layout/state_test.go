package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestComputeStateLayoutSimple(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.State
	graph.Direction = ir.TopDown

	graph.EnsureNode("__start__", nil, nil)
	graph.EnsureNode("First", nil, nil)
	graph.EnsureNode("__end__", nil, nil)
	graph.Edges = append(graph.Edges,
		&ir.Edge{From: "__start__", To: "First", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "First", To: "__end__", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.State {
		t.Errorf("Kind = %v, want State", lay.Kind)
	}
	if len(lay.Nodes) != 3 {
		t.Errorf("nodes = %d, want 3", len(lay.Nodes))
	}
	if _, ok := lay.Diagram.(StateData); !ok {
		t.Errorf("Diagram data type = %T, want StateData", lay.Diagram)
	}
}

func TestComputeStateLayoutComposite(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.State
	graph.Direction = ir.TopDown

	graph.EnsureNode("Outer", nil, nil)
	inner := ir.NewGraph()
	inner.Kind = ir.State
	inner.EnsureNode("__start__", nil, nil)
	inner.EnsureNode("inner1", nil, nil)
	inner.Edges = append(inner.Edges, &ir.Edge{From: "__start__", To: "inner1", Directed: true, ArrowEnd: true})
	graph.CompositeStates["Outer"] = &ir.CompositeState{
		ID:    "Outer",
		Label: "Outer",
		Inner: inner,
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	outer := lay.Nodes["Outer"]
	if outer == nil {
		t.Fatal("Outer node not in layout")
	}
	if outer.Width < 100 {
		t.Errorf("Outer width = %f, expected > 100 for composite", outer.Width)
	}
	sd, ok := lay.Diagram.(StateData)
	if !ok {
		t.Fatal("expected StateData")
	}
	if sd.InnerLayouts["Outer"] == nil {
		t.Error("expected inner layout for Outer")
	}
}
