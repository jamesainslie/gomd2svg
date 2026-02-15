package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseRequirementBasic(t *testing.T) {
	input := `requirementDiagram

requirement test_req {
id: 1
text: the test text.
risk: high
verifymethod: test
}

element test_entity {
type: simulation
docref: DOC-001
}

test_entity - satisfies -> test_req`

	out, err := parseRequirement(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Requirement {
		t.Fatalf("Kind = %v, want Requirement", g.Kind)
	}
	if len(g.Requirements) != 1 {
		t.Fatalf("Requirements = %d, want 1", len(g.Requirements))
	}
	req := g.Requirements[0]
	if req.Name != "test_req" {
		t.Errorf("req name = %q", req.Name)
	}
	if req.ID != "1" {
		t.Errorf("req id = %q", req.ID)
	}
	if req.Risk != ir.RiskHigh {
		t.Errorf("req risk = %v", req.Risk)
	}
	if req.VerifyMethod != ir.VerifyTest {
		t.Errorf("req verify = %v", req.VerifyMethod)
	}
	if len(g.ReqElements) != 1 {
		t.Fatalf("Elements = %d, want 1", len(g.ReqElements))
	}
	elem := g.ReqElements[0]
	if elem.Type != "simulation" {
		t.Errorf("elem type = %q", elem.Type)
	}
	if len(g.ReqRelationships) != 1 {
		t.Fatalf("Rels = %d, want 1", len(g.ReqRelationships))
	}
	if g.ReqRelationships[0].Type != ir.ReqRelSatisfies {
		t.Errorf("rel type = %v", g.ReqRelationships[0].Type)
	}
}

func TestParseRequirementMultiple(t *testing.T) {
	input := `requirementDiagram

functionalRequirement req1 {
id: FR-001
text: Must authenticate users
risk: medium
verifymethod: demonstration
}

requirement req2 {
id: REQ-002
text: Must log events
risk: low
verifymethod: analysis
}

req1 - derives -> req2`

	out, err := parseRequirement(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Requirements) != 2 {
		t.Fatalf("Requirements = %d, want 2", len(out.Graph.Requirements))
	}
	if out.Graph.Requirements[0].Type != ir.ReqTypeFunctional {
		t.Errorf("req1 type = %v, want Functional", out.Graph.Requirements[0].Type)
	}
}

func TestParseRequirementEmpty(t *testing.T) {
	input := `requirementDiagram`
	out, err := parseRequirement(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Requirements) != 0 {
		t.Errorf("Requirements = %d, want 0", len(out.Graph.Requirements))
	}
}
