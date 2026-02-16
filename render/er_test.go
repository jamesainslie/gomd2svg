package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderEREntities(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Er
	graph.Direction = ir.TopDown
	graph.EnsureNode("CUSTOMER", nil, nil)
	graph.Entities["CUSTOMER"] = &ir.Entity{
		ID: "CUSTOMER",
		Attributes: []ir.EntityAttribute{
			{Type: "string", Name: "name"},
			{Type: "int", Name: "id", Keys: []ir.AttributeKey{ir.KeyPrimary}},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "CUSTOMER") {
		t.Error("missing entity name 'CUSTOMER'")
	}
	if !strings.Contains(svg, "PK") {
		t.Error("missing key constraint 'PK'")
	}
	if !strings.Contains(svg, "name") {
		t.Error("missing attribute 'name'")
	}
}

func TestRenderERRelationship(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Er
	graph.Direction = ir.TopDown
	graph.EnsureNode("A", nil, nil)
	graph.EnsureNode("B", nil, nil)
	graph.Entities["A"] = &ir.Entity{ID: "A"}
	graph.Entities["B"] = &ir.Entity{ID: "B"}

	startDec := ir.DecCrowsFootOne
	endDec := ir.DecCrowsFootZeroMany
	label := "has"
	graph.Edges = append(graph.Edges, &ir.Edge{
		From: "A", To: "B",
		StartDecoration: &startDec,
		EndDecoration:   &endDec,
		Label:           &label,
	})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "has") {
		t.Error("missing relationship label 'has'")
	}
	if !strings.Contains(svg, "edgePath") {
		t.Error("missing edge path")
	}
}
