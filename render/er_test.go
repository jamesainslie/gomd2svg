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
	g := ir.NewGraph()
	g.Kind = ir.Er
	g.Direction = ir.TopDown
	g.EnsureNode("CUSTOMER", nil, nil)
	g.Entities["CUSTOMER"] = &ir.Entity{
		ID: "CUSTOMER",
		Attributes: []ir.EntityAttribute{
			{Type: "string", Name: "name"},
			{Type: "int", Name: "id", Keys: []ir.AttributeKey{ir.KeyPrimary}},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
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
	g := ir.NewGraph()
	g.Kind = ir.Er
	g.Direction = ir.TopDown
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	g.Entities["A"] = &ir.Entity{ID: "A"}
	g.Entities["B"] = &ir.Entity{ID: "B"}

	startDec := ir.DecCrowsFootOne
	endDec := ir.DecCrowsFootZeroMany
	label := "has"
	g.Edges = append(g.Edges, &ir.Edge{
		From: "A", To: "B",
		StartDecoration: &startDec,
		EndDecoration:   &endDec,
		Label:           &label,
	})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "has") {
		t.Error("missing relationship label 'has'")
	}
	if !strings.Contains(svg, "edgePath") {
		t.Error("missing edge path")
	}
}
