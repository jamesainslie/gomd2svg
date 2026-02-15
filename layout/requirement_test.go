package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRequirementLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Requirement
	g.Direction = ir.TopDown

	reqLabel := "test_req"
	elemLabel := "test_elem"
	g.EnsureNode("test_req", &reqLabel, nil)
	g.EnsureNode("test_elem", &elemLabel, nil)

	g.Requirements = append(g.Requirements, &ir.RequirementDef{
		Name: "test_req", ID: "REQ-001", Text: "Must work", Risk: ir.RiskHigh, VerifyMethod: ir.VerifyTest,
	})
	g.ReqElements = append(g.ReqElements, &ir.ElementDef{
		Name: "test_elem", Type: "Simulation",
	})

	relLabel := "satisfies"
	g.Edges = append(g.Edges, &ir.Edge{From: "test_elem", To: "test_req", Label: &relLabel, Directed: true, ArrowEnd: true})
	g.ReqRelationships = append(g.ReqRelationships, &ir.RequirementRel{Source: "test_elem", Target: "test_req", Type: ir.ReqRelSatisfies})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Requirement {
		t.Fatalf("Kind = %v", l.Kind)
	}
	rd, ok := l.Diagram.(RequirementData)
	if !ok {
		t.Fatal("Diagram is not RequirementData")
	}
	if len(rd.Requirements) != 1 {
		t.Errorf("Requirements = %d", len(rd.Requirements))
	}
	if len(rd.Elements) != 1 {
		t.Errorf("Elements = %d", len(rd.Elements))
	}
	if len(l.Nodes) != 2 {
		t.Errorf("Nodes = %d", len(l.Nodes))
	}
	if len(l.Edges) != 1 {
		t.Errorf("Edges = %d", len(l.Edges))
	}
}

func TestRequirementLayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Requirement
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)
	if len(l.Nodes) != 0 {
		t.Errorf("Nodes = %d", len(l.Nodes))
	}
}
