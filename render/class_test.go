package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderClassCompartments(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Class
	g.Direction = ir.TopDown
	g.EnsureNode("Animal", nil, nil)
	g.Members["Animal"] = &ir.ClassMembers{
		Attributes: []ir.ClassMember{
			{Name: "name", Type: "String", Visibility: ir.VisPublic},
		},
		Methods: []ir.ClassMember{
			{Name: "speak", IsMethod: true, Visibility: ir.VisPublic, Type: "void"},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "Animal") {
		t.Error("missing class name 'Animal'")
	}
	if !strings.Contains(svg, "+") {
		t.Error("missing visibility symbol '+'")
	}
	if strings.Count(svg, "<line") < 2 {
		t.Errorf("expected at least 2 divider lines, got %d", strings.Count(svg, "<line"))
	}
}

func TestRenderClassRelationshipMarkers(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Class
	g.Direction = ir.TopDown
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	closedTri := ir.ClosedTriangle
	g.Edges = append(g.Edges, &ir.Edge{
		From: "A", To: "B", Directed: true, ArrowEnd: true,
		ArrowEndKind: &closedTri,
	})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "url(#marker-closed-triangle)") {
		t.Error("missing closed triangle marker reference on edge")
	}
}
