package ir

// RequirementType classifies a requirement definition.
type RequirementType int

const (
	ReqTypeRequirement RequirementType = iota
	ReqTypeFunctional
	ReqTypeInterface
	ReqTypePerformance
	ReqTypePhysical
	ReqTypeDesignConstraint
)

func (t RequirementType) String() string {
	switch t {
	case ReqTypeFunctional:
		return "functionalRequirement"
	case ReqTypeInterface:
		return "interfaceRequirement"
	case ReqTypePerformance:
		return "performanceRequirement"
	case ReqTypePhysical:
		return "physicalRequirement"
	case ReqTypeDesignConstraint:
		return "designConstraint"
	default:
		return "requirement"
	}
}

// Stereotype returns the UML-style stereotype label for display.
func (t RequirementType) Stereotype() string {
	switch t {
	case ReqTypeFunctional:
		return "Functional Requirement"
	case ReqTypeInterface:
		return "Interface Requirement"
	case ReqTypePerformance:
		return "Performance Requirement"
	case ReqTypePhysical:
		return "Physical Requirement"
	case ReqTypeDesignConstraint:
		return "Design Constraint"
	default:
		return "Requirement"
	}
}

// RiskLevel represents a requirement's risk assessment.
type RiskLevel int

const (
	RiskNone RiskLevel = iota
	RiskLow
	RiskMedium
	RiskHigh
)

func (r RiskLevel) String() string {
	switch r {
	case RiskLow:
		return "Low"
	case RiskMedium:
		return "Medium"
	case RiskHigh:
		return "High"
	default:
		return ""
	}
}

// VerifyMethod represents how a requirement is verified.
type VerifyMethod int

const (
	VerifyNone VerifyMethod = iota
	VerifyAnalysis
	VerifyInspection
	VerifyTest
	VerifyDemonstration
)

func (v VerifyMethod) String() string {
	switch v {
	case VerifyAnalysis:
		return "Analysis"
	case VerifyInspection:
		return "Inspection"
	case VerifyTest:
		return "Test"
	case VerifyDemonstration:
		return "Demonstration"
	default:
		return ""
	}
}

// RequirementRelType classifies a relationship between requirements/elements.
type RequirementRelType int

const (
	ReqRelContains RequirementRelType = iota
	ReqRelCopies
	ReqRelDerives
	ReqRelSatisfies
	ReqRelVerifies
	ReqRelRefines
	ReqRelTraces
)

func (r RequirementRelType) String() string {
	switch r {
	case ReqRelContains:
		return "contains"
	case ReqRelCopies:
		return "copies"
	case ReqRelDerives:
		return "derives"
	case ReqRelSatisfies:
		return "satisfies"
	case ReqRelVerifies:
		return "verifies"
	case ReqRelRefines:
		return "refines"
	case ReqRelTraces:
		return "traces"
	default:
		return ""
	}
}

// RequirementDef represents a requirement block.
type RequirementDef struct {
	Name         string
	ID           string
	Text         string
	Type         RequirementType
	Risk         RiskLevel
	VerifyMethod VerifyMethod
}

// ElementDef represents an element block in a requirement diagram.
type ElementDef struct {
	Name   string
	Type   string
	DocRef string
}

// RequirementRel represents a relationship between two nodes.
type RequirementRel struct {
	Source string
	Target string
	Type   RequirementRelType
}
