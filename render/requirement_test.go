package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderRequirement(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Requirement

	// Add a requirement.
	g.Requirements = append(g.Requirements, &ir.RequirementDef{
		Name:         "test_req",
		ID:           "REQ-001",
		Text:         "Must do something",
		Type:         ir.ReqTypeFunctional,
		Risk:         ir.RiskHigh,
		VerifyMethod: ir.VerifyTest,
	})
	reqLabel := "test_req"
	g.EnsureNode("test_req", &reqLabel, nil)

	// Add an element.
	g.ReqElements = append(g.ReqElements, &ir.ElementDef{
		Name:   "test_element",
		Type:   "Simulation",
		DocRef: "DOC-001",
	})
	elemLabel := "test_element"
	g.EnsureNode("test_element", &elemLabel, nil)

	// Add a relationship edge.
	relLabel := "satisfies"
	g.ReqRelationships = append(g.ReqRelationships, &ir.RequirementRel{
		Source: "test_element",
		Target: "test_req",
		Type:   ir.ReqRelSatisfies,
	})
	g.Edges = append(g.Edges, &ir.Edge{
		From:     "test_element",
		To:       "test_req",
		Label:    &relLabel,
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
	if !strings.Contains(svg, "Functional Requirement") {
		t.Error("missing requirement stereotype")
	}
	if !strings.Contains(svg, "test_req") {
		t.Error("missing requirement name")
	}
	if !strings.Contains(svg, "REQ-001") {
		t.Error("missing requirement ID")
	}
	if !strings.Contains(svg, "Risk: High") {
		t.Error("missing risk metadata")
	}
	if !strings.Contains(svg, "Verify: Test") {
		t.Error("missing verify method metadata")
	}
	if !strings.Contains(svg, "element") {
		t.Error("missing element stereotype")
	}
	if !strings.Contains(svg, "Type: Simulation") {
		t.Error("missing element type")
	}
	if !strings.Contains(svg, "Doc: DOC-001") {
		t.Error("missing element docref")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing rectangles")
	}
	if !strings.Contains(svg, "satisfies") {
		t.Error("missing edge label")
	}
}

func TestRenderRequirementEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Requirement

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
