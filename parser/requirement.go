package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	reqBlockStartRe  = regexp.MustCompile(`^(requirement|functionalRequirement|interfaceRequirement|performanceRequirement|physicalRequirement|designConstraint)\s+(\w+)\s*\{?\s*$`)
	elemBlockStartRe = regexp.MustCompile(`^element\s+(\w+)\s*\{?\s*$`)
	reqFieldRe       = regexp.MustCompile(`^\s*(\w+)\s*:\s*(.+?)\s*$`)
	reqRelRe         = regexp.MustCompile(`^(\w+)\s+-\s+(contains|copies|derives|satisfies|verifies|refines|traces)\s+->\s+(\w+)\s*$`)
)

func parseRequirement(input string) (*ParseOutput, error) {
	lines := preprocessInput(input)
	g := ir.NewGraph()
	g.Kind = ir.Requirement

	if len(lines) > 0 {
		lower := strings.ToLower(lines[0])
		if strings.HasPrefix(lower, "requirementdiagram") {
			lines = lines[1:]
		}
	}

	i := 0
	for i < len(lines) {
		line := lines[i]

		if m := reqBlockStartRe.FindStringSubmatch(line); m != nil {
			reqType := parseReqType(m[1])
			name := m[2]
			i++
			req := &ir.RequirementDef{Name: name, Type: reqType}
			for i < len(lines) && lines[i] != "}" {
				if fm := reqFieldRe.FindStringSubmatch(lines[i]); fm != nil {
					switch strings.ToLower(fm[1]) {
					case "id":
						req.ID = fm[2]
					case "text":
						req.Text = fm[2]
					case "risk":
						req.Risk = parseRiskLevel(fm[2])
					case "verifymethod":
						req.VerifyMethod = parseVerifyMethod(fm[2])
					}
				}
				i++
			}
			g.Requirements = append(g.Requirements, req)
			label := name
			g.EnsureNode(name, &label, nil)
			i++
			continue
		}

		if m := elemBlockStartRe.FindStringSubmatch(line); m != nil {
			name := m[1]
			i++
			elem := &ir.ElementDef{Name: name}
			for i < len(lines) && lines[i] != "}" {
				if fm := reqFieldRe.FindStringSubmatch(lines[i]); fm != nil {
					switch strings.ToLower(fm[1]) {
					case "type":
						elem.Type = fm[2]
					case "docref":
						elem.DocRef = fm[2]
					}
				}
				i++
			}
			g.ReqElements = append(g.ReqElements, elem)
			label := name
			g.EnsureNode(name, &label, nil)
			i++
			continue
		}

		if m := reqRelRe.FindStringSubmatch(line); m != nil {
			rel := &ir.RequirementRel{
				Source: m[1],
				Target: m[3],
				Type:   parseRelType(m[2]),
			}
			g.ReqRelationships = append(g.ReqRelationships, rel)
			relLabel := rel.Type.String()
			g.Edges = append(g.Edges, &ir.Edge{
				From:     rel.Source,
				To:       rel.Target,
				Label:    &relLabel,
				Directed: true,
				ArrowEnd: true,
			})
			i++
			continue
		}

		i++
	}

	return &ParseOutput{Graph: g}, nil
}

func parseReqType(s string) ir.RequirementType {
	switch strings.ToLower(s) {
	case "functionalrequirement":
		return ir.ReqTypeFunctional
	case "interfacerequirement":
		return ir.ReqTypeInterface
	case "performancerequirement":
		return ir.ReqTypePerformance
	case "physicalrequirement":
		return ir.ReqTypePhysical
	case "designconstraint":
		return ir.ReqTypeDesignConstraint
	default:
		return ir.ReqTypeRequirement
	}
}

func parseRiskLevel(s string) ir.RiskLevel {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "low":
		return ir.RiskLow
	case "medium":
		return ir.RiskMedium
	case "high":
		return ir.RiskHigh
	default:
		return ir.RiskNone
	}
}

func parseVerifyMethod(s string) ir.VerifyMethod {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "analysis":
		return ir.VerifyAnalysis
	case "inspection":
		return ir.VerifyInspection
	case "test":
		return ir.VerifyTest
	case "demonstration":
		return ir.VerifyDemonstration
	default:
		return ir.VerifyNone
	}
}

func parseRelType(s string) ir.RequirementRelType {
	switch strings.ToLower(s) {
	case "contains":
		return ir.ReqRelContains
	case "copies":
		return ir.ReqRelCopies
	case "derives":
		return ir.ReqRelDerives
	case "satisfies":
		return ir.ReqRelSatisfies
	case "verifies":
		return ir.ReqRelVerifies
	case "refines":
		return ir.ReqRelRefines
	case "traces":
		return ir.ReqRelTraces
	default:
		return ir.ReqRelContains
	}
}
