package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

// reqTypeField is the field key for requirement type.
const reqTypeField = "type"

var (
	reqBlockStartRe  = regexp.MustCompile(`^(requirement|functionalRequirement|interfaceRequirement|performanceRequirement|physicalRequirement|designConstraint)\s+(\w+)\s*\{?\s*$`)
	elemBlockStartRe = regexp.MustCompile(`^element\s+(\w+)\s*\{?\s*$`)
	reqFieldRe       = regexp.MustCompile(`^\s*(\w+)\s*:\s*(.+?)\s*$`)
	reqRelRe         = regexp.MustCompile(`^(\w+)\s+-\s+(contains|copies|derives|satisfies|verifies|refines|traces)\s+->\s+(\w+)\s*$`)
)

func parseRequirement(input string) (*ParseOutput, error) { //nolint:unparam // error return is part of the parser interface contract used by Parse().
	lines := preprocessInput(input)
	graph := ir.NewGraph()
	graph.Kind = ir.Requirement

	if len(lines) > 0 {
		lower := strings.ToLower(lines[0])
		if strings.HasPrefix(lower, "requirementdiagram") {
			lines = lines[1:]
		}
	}

	idx := 0
	for idx < len(lines) {
		line := lines[idx]

		if match := reqBlockStartRe.FindStringSubmatch(line); match != nil {
			reqType := parseReqType(match[1])
			name := match[2]
			idx++
			req := &ir.RequirementDef{Name: name, Type: reqType}
			for idx < len(lines) && lines[idx] != "}" {
				if fieldMatch := reqFieldRe.FindStringSubmatch(lines[idx]); fieldMatch != nil {
					switch strings.ToLower(fieldMatch[1]) {
					case "id":
						req.ID = fieldMatch[2]
					case "text":
						req.Text = fieldMatch[2]
					case "risk":
						req.Risk = parseRiskLevel(fieldMatch[2])
					case "verifymethod":
						req.VerifyMethod = parseVerifyMethod(fieldMatch[2])
					}
				}
				idx++
			}
			graph.Requirements = append(graph.Requirements, req)
			label := name
			graph.EnsureNode(name, &label, nil)
			idx++
			continue
		}

		if match := elemBlockStartRe.FindStringSubmatch(line); match != nil {
			name := match[1]
			idx++
			elem := &ir.ElementDef{Name: name}
			for idx < len(lines) && lines[idx] != "}" {
				if fieldMatch := reqFieldRe.FindStringSubmatch(lines[idx]); fieldMatch != nil {
					switch strings.ToLower(fieldMatch[1]) {
					case reqTypeField:
						elem.Type = fieldMatch[2]
					case "docref":
						elem.DocRef = fieldMatch[2]
					}
				}
				idx++
			}
			graph.ReqElements = append(graph.ReqElements, elem)
			label := name
			graph.EnsureNode(name, &label, nil)
			idx++
			continue
		}

		if match := reqRelRe.FindStringSubmatch(line); match != nil {
			rel := &ir.RequirementRel{
				Source: match[1],
				Target: match[3],
				Type:   parseRelType(match[2]),
			}
			graph.ReqRelationships = append(graph.ReqRelationships, rel)
			relLabel := rel.Type.String()
			graph.Edges = append(graph.Edges, &ir.Edge{
				From:     rel.Source,
				To:       rel.Target,
				Label:    &relLabel,
				Directed: true,
				ArrowEnd: true,
			})
			idx++
			continue
		}

		idx++
	}

	return &ParseOutput{Graph: graph}, nil
}

func parseReqType(str string) ir.RequirementType {
	switch strings.ToLower(str) {
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

func parseRiskLevel(str string) ir.RiskLevel {
	switch strings.ToLower(strings.TrimSpace(str)) {
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

func parseVerifyMethod(str string) ir.VerifyMethod {
	switch strings.ToLower(strings.TrimSpace(str)) {
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

func parseRelType(str string) ir.RequirementRelType {
	switch strings.ToLower(str) {
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
