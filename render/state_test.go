package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderStateSimple(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown
	g.EnsureNode("__start__", nil, nil)
	g.EnsureNode("First", nil, nil)
	g.EnsureNode("__end__", nil, nil)
	g.StateDescriptions["First"] = "First state"
	g.Edges = append(g.Edges,
		&ir.Edge{From: "__start__", To: "First", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "First", To: "__end__", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<circle") {
		t.Error("missing circle element for start/end state")
	}
	if !strings.Contains(svg, "First") {
		t.Error("missing state label 'First'")
	}
	// Description should appear below the state name
	if !strings.Contains(svg, "First state") {
		t.Error("missing state description 'First state'")
	}
	// State box should use rounded rect
	if !strings.Contains(svg, "<rect") {
		t.Error("missing rect element for state box")
	}
	// Edges should have arrowheads
	if !strings.Contains(svg, "marker-end") {
		t.Error("missing arrowhead marker on edges")
	}
}

func TestRenderStateComposite(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown
	g.EnsureNode("Outer", nil, nil)

	inner := ir.NewGraph()
	inner.Kind = ir.State
	inner.EnsureNode("__start__", nil, nil)
	inner.EnsureNode("inner1", nil, nil)
	inner.Edges = append(inner.Edges, &ir.Edge{From: "__start__", To: "inner1", Directed: true, ArrowEnd: true})
	g.CompositeStates["Outer"] = &ir.CompositeState{ID: "Outer", Label: "Outer", Inner: inner}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "Outer") {
		t.Error("missing composite label 'Outer'")
	}
	if !strings.Contains(svg, "inner1") {
		t.Error("missing inner state 'inner1'")
	}
}

func TestRenderStateForkJoin(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown
	g.EnsureNode("fork1", nil, nil)
	g.StateAnnotations["fork1"] = ir.StateFork
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	g.Edges = append(g.Edges,
		&ir.Edge{From: "fork1", To: "A", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "fork1", To: "B", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Fork should render as a filled rect bar
	if !strings.Contains(svg, "<rect") {
		t.Error("missing rect element for fork bar")
	}
	if !strings.Contains(svg, "A") {
		t.Error("missing state label 'A'")
	}
}

func TestRenderStateChoice(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown
	g.EnsureNode("choice1", nil, nil)
	g.StateAnnotations["choice1"] = ir.StateChoice
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	g.Edges = append(g.Edges,
		&ir.Edge{From: "choice1", To: "A", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "choice1", To: "B", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Choice should render as a diamond polygon
	if !strings.Contains(svg, "<polygon") {
		t.Error("missing polygon element for choice diamond")
	}
}

func TestRenderStateEndBullseye(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown
	g.EnsureNode("__start__", nil, nil)
	g.EnsureNode("__end__", nil, nil)
	g.Edges = append(g.Edges,
		&ir.Edge{From: "__start__", To: "__end__", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// End state should have multiple circle elements (bullseye)
	count := strings.Count(svg, "<circle")
	if count < 2 {
		t.Errorf("expected at least 2 circle elements for bullseye end state, got %d", count)
	}
}
