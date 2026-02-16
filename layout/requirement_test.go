package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRequirementLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Requirement
	graph.Direction = ir.TopDown

	reqLabel := "test_req"
	elemLabel := "test_elem"
	graph.EnsureNode("test_req", &reqLabel, nil)
	graph.EnsureNode("test_elem", &elemLabel, nil)

	graph.Requirements = append(graph.Requirements, &ir.RequirementDef{
		Name: "test_req", ID: "REQ-001", Text: "Must work", Risk: ir.RiskHigh, VerifyMethod: ir.VerifyTest,
	})
	graph.ReqElements = append(graph.ReqElements, &ir.ElementDef{
		Name: "test_elem", Type: "Simulation",
	})

	relLabel := "satisfies"
	graph.Edges = append(graph.Edges, &ir.Edge{From: "test_elem", To: "test_req", Label: &relLabel, Directed: true, ArrowEnd: true})
	graph.ReqRelationships = append(graph.ReqRelationships, &ir.RequirementRel{Source: "test_elem", Target: "test_req", Type: ir.ReqRelSatisfies})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.Requirement {
		t.Fatalf("Kind = %v", lay.Kind)
	}
	rd, ok := lay.Diagram.(RequirementData)
	if !ok {
		t.Fatal("Diagram is not RequirementData")
	}
	if len(rd.Requirements) != 1 {
		t.Errorf("Requirements = %d", len(rd.Requirements))
	}
	if len(rd.Elements) != 1 {
		t.Errorf("Elements = %d", len(rd.Elements))
	}
	if len(lay.Nodes) != 2 {
		t.Errorf("Nodes = %d", len(lay.Nodes))
	}
	if len(lay.Edges) != 1 {
		t.Errorf("Edges = %d", len(lay.Edges))
	}
}

func TestRequirementLayoutEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Requirement
	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)
	if len(lay.Nodes) != 0 {
		t.Errorf("Nodes = %d", len(lay.Nodes))
	}
}
