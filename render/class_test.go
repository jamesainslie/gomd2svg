package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderClassCompartments(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Class
	graph.Direction = ir.TopDown
	graph.EnsureNode("Animal", nil, nil)
	graph.Members["Animal"] = &ir.ClassMembers{
		Attributes: []ir.ClassMember{
			{Name: "name", Type: "String", Visibility: ir.VisPublic},
		},
		Methods: []ir.ClassMember{
			{Name: "speak", IsMethod: true, Visibility: ir.VisPublic, Type: "void"},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
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
	graph := ir.NewGraph()
	graph.Kind = ir.Class
	graph.Direction = ir.TopDown
	graph.EnsureNode("A", nil, nil)
	graph.EnsureNode("B", nil, nil)
	closedTri := ir.ClosedTriangle
	graph.Edges = append(graph.Edges, &ir.Edge{
		From: "A", To: "B", Directed: true, ArrowEnd: true,
		ArrowEndKind: &closedTri,
	})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "url(#marker-closed-triangle)") {
		t.Error("missing closed triangle marker reference on edge")
	}
}
