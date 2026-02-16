package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderC4(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.C4
	graph.C4SubKind = ir.C4Context

	graph.C4Elements = append(graph.C4Elements, &ir.C4Element{
		ID:    "user",
		Label: "User",
		Type:  ir.C4Person,
	})
	graph.C4Elements = append(graph.C4Elements, &ir.C4Element{
		ID:         "webapp",
		Label:      "Web App",
		Type:       ir.C4System,
		Technology: "Go",
	})

	// Add nodes for each element.
	for _, elem := range graph.C4Elements {
		graph.EnsureNode(elem.ID, &elem.Label, nil)
	}

	graph.Edges = append(graph.Edges, &ir.Edge{
		From:     "user",
		To:       "webapp",
		Directed: true,
		ArrowEnd: true,
	})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "User") {
		t.Error("missing User label")
	}
	if !strings.Contains(svg, "Web App") {
		t.Error("missing Web App label")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing rectangles")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing circle (person icon head)")
	}
}

func TestRenderC4Empty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.C4

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
