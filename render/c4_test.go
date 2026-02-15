package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderC4(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.C4
	g.C4SubKind = ir.C4Context

	g.C4Elements = append(g.C4Elements, &ir.C4Element{
		ID:    "user",
		Label: "User",
		Type:  ir.C4Person,
	})
	g.C4Elements = append(g.C4Elements, &ir.C4Element{
		ID:         "webapp",
		Label:      "Web App",
		Type:       ir.C4System,
		Technology: "Go",
	})

	// Add nodes for each element.
	for _, elem := range g.C4Elements {
		g.EnsureNode(elem.ID, &elem.Label, nil)
	}

	g.Edges = append(g.Edges, &ir.Edge{
		From:     "user",
		To:       "webapp",
		Directed: true,
		ArrowEnd: true,
	})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
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
	g := ir.NewGraph()
	g.Kind = ir.C4

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
