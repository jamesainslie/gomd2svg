package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderStateSimple(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.State
	graph.Direction = ir.TopDown
	graph.EnsureNode("__start__", nil, nil)
	graph.EnsureNode("First", nil, nil)
	graph.EnsureNode("__end__", nil, nil)
	graph.StateDescriptions["First"] = "First state"
	graph.Edges = append(graph.Edges,
		&ir.Edge{From: "__start__", To: "First", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "First", To: "__end__", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
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
	graph := ir.NewGraph()
	graph.Kind = ir.State
	graph.Direction = ir.TopDown
	graph.EnsureNode("Outer", nil, nil)

	inner := ir.NewGraph()
	inner.Kind = ir.State
	inner.EnsureNode("__start__", nil, nil)
	inner.EnsureNode("inner1", nil, nil)
	inner.Edges = append(inner.Edges, &ir.Edge{From: "__start__", To: "inner1", Directed: true, ArrowEnd: true})
	graph.CompositeStates["Outer"] = &ir.CompositeState{ID: "Outer", Label: "Outer", Inner: inner}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "Outer") {
		t.Error("missing composite label 'Outer'")
	}
	if !strings.Contains(svg, "inner1") {
		t.Error("missing inner state 'inner1'")
	}
}

func TestRenderStateForkJoin(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.State
	graph.Direction = ir.TopDown
	graph.EnsureNode("fork1", nil, nil)
	graph.StateAnnotations["fork1"] = ir.StateFork
	graph.EnsureNode("A", nil, nil)
	graph.EnsureNode("B", nil, nil)
	graph.Edges = append(graph.Edges,
		&ir.Edge{From: "fork1", To: "A", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "fork1", To: "B", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
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
	graph := ir.NewGraph()
	graph.Kind = ir.State
	graph.Direction = ir.TopDown
	graph.EnsureNode("choice1", nil, nil)
	graph.StateAnnotations["choice1"] = ir.StateChoice
	graph.EnsureNode("A", nil, nil)
	graph.EnsureNode("B", nil, nil)
	graph.Edges = append(graph.Edges,
		&ir.Edge{From: "choice1", To: "A", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "choice1", To: "B", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Choice should render as a diamond polygon
	if !strings.Contains(svg, "<polygon") {
		t.Error("missing polygon element for choice diamond")
	}
}

func TestRenderStateEndBullseye(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.State
	graph.Direction = ir.TopDown
	graph.EnsureNode("__start__", nil, nil)
	graph.EnsureNode("__end__", nil, nil)
	graph.Edges = append(graph.Edges,
		&ir.Edge{From: "__start__", To: "__end__", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// End state should have multiple circle elements (bullseye)
	count := strings.Count(svg, "<circle")
	if count < 2 {
		t.Errorf("expected at least 2 circle elements for bullseye end state, got %d", count)
	}
}
