package ir

import "testing"

func TestRequirementType(t *testing.T) {
	tests := []struct {
		typ        RequirementType
		str        string
		stereotype string
	}{
		{ReqTypeRequirement, "requirement", "Requirement"},
		{ReqTypeFunctional, "functionalRequirement", "Functional Requirement"},
		{ReqTypeInterface, "interfaceRequirement", "Interface Requirement"},
		{ReqTypePerformance, "performanceRequirement", "Performance Requirement"},
		{ReqTypePhysical, "physicalRequirement", "Physical Requirement"},
		{ReqTypeDesignConstraint, "designConstraint", "Design Constraint"},
	}
	for _, tt := range tests {
		if tt.typ.String() != tt.str {
			t.Errorf("RequirementType(%d).String() = %q, want %q", tt.typ, tt.typ.String(), tt.str)
		}
		if tt.typ.Stereotype() != tt.stereotype {
			t.Errorf("RequirementType(%d).Stereotype() = %q, want %q", tt.typ, tt.typ.Stereotype(), tt.stereotype)
		}
	}
}

func TestRiskLevel(t *testing.T) {
	if RiskLow.String() != "Low" {
		t.Errorf("RiskLow = %q", RiskLow.String())
	}
	if RiskMedium.String() != "Medium" {
		t.Errorf("RiskMedium = %q", RiskMedium.String())
	}
	if RiskHigh.String() != "High" {
		t.Errorf("RiskHigh = %q", RiskHigh.String())
	}
	if RiskNone.String() != "" {
		t.Errorf("RiskNone = %q", RiskNone.String())
	}
}

func TestVerifyMethod(t *testing.T) {
	if VerifyTest.String() != "Test" {
		t.Errorf("VerifyTest = %q", VerifyTest.String())
	}
	if VerifyAnalysis.String() != "Analysis" {
		t.Errorf("VerifyAnalysis = %q", VerifyAnalysis.String())
	}
}

func TestRequirementRelType(t *testing.T) {
	if ReqRelSatisfies.String() != "satisfies" {
		t.Errorf("ReqRelSatisfies = %q", ReqRelSatisfies.String())
	}
	if ReqRelContains.String() != "contains" {
		t.Errorf("ReqRelContains = %q", ReqRelContains.String())
	}
}

func TestRequirementGraphFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Requirement
	graph.Requirements = append(graph.Requirements, &RequirementDef{
		Name: "test_req",
		ID:   "REQ-001",
		Text: "Must do something",
		Type: ReqTypeFunctional,
		Risk: RiskHigh,
	})
	graph.ReqElements = append(graph.ReqElements, &ElementDef{
		Name:   "test_element",
		Type:   "Simulation",
		DocRef: "DOC-001",
	})
	graph.ReqRelationships = append(graph.ReqRelationships, &RequirementRel{
		Source: "test_element",
		Target: "test_req",
		Type:   ReqRelSatisfies,
	})
	if len(graph.Requirements) != 1 {
		t.Errorf("Requirements = %d", len(graph.Requirements))
	}
	if len(graph.ReqElements) != 1 {
		t.Errorf("ReqElements = %d", len(graph.ReqElements))
	}
	if len(graph.ReqRelationships) != 1 {
		t.Errorf("ReqRelationships = %d", len(graph.ReqRelationships))
	}
}
