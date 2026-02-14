package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestComputeStateLayoutSimple(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown

	g.EnsureNode("__start__", nil, nil)
	g.EnsureNode("First", nil, nil)
	g.EnsureNode("__end__", nil, nil)
	g.Edges = append(g.Edges,
		&ir.Edge{From: "__start__", To: "First", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "First", To: "__end__", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.State {
		t.Errorf("Kind = %v, want State", l.Kind)
	}
	if len(l.Nodes) != 3 {
		t.Errorf("nodes = %d, want 3", len(l.Nodes))
	}
	if _, ok := l.Diagram.(StateData); !ok {
		t.Errorf("Diagram data type = %T, want StateData", l.Diagram)
	}
}

func TestComputeStateLayoutComposite(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown

	g.EnsureNode("Outer", nil, nil)
	inner := ir.NewGraph()
	inner.Kind = ir.State
	inner.EnsureNode("__start__", nil, nil)
	inner.EnsureNode("inner1", nil, nil)
	inner.Edges = append(inner.Edges, &ir.Edge{From: "__start__", To: "inner1", Directed: true, ArrowEnd: true})
	g.CompositeStates["Outer"] = &ir.CompositeState{
		ID:    "Outer",
		Label: "Outer",
		Inner: inner,
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	outer := l.Nodes["Outer"]
	if outer == nil {
		t.Fatal("Outer node not in layout")
	}
	if outer.Width < 100 {
		t.Errorf("Outer width = %f, expected > 100 for composite", outer.Width)
	}
	sd, ok := l.Diagram.(StateData)
	if !ok {
		t.Fatal("expected StateData")
	}
	if sd.InnerLayouts["Outer"] == nil {
		t.Error("expected inner layout for Outer")
	}
}
