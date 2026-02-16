package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	c4ElementRe  = regexp.MustCompile(`^(Person|Person_Ext|System|System_Ext|SystemDb|SystemDb_Ext|SystemQueue|SystemQueue_Ext|Container|Container_Ext|ContainerDb|ContainerDb_Ext|ContainerQueue|ContainerQueue_Ext|Component|Component_Ext)\s*\((.+)\)\s*$`)
	c4BoundaryRe = regexp.MustCompile(`^(Enterprise_Boundary|System_Boundary|Container_Boundary|Boundary)\s*\(([^,]+),\s*"([^"]*)"(?:\s*,\s*"([^"]*)")?\)\s*\{?\s*$`)
	c4RelRe      = regexp.MustCompile(`^(Rel|Rel_Back|Rel_Neighbor|Rel_Back_Neighbor|BiRel|BiRel_Neighbor)\s*\((.+)\)\s*$`)
)

func parseC4(input string) (*ParseOutput, error) {
	lines := preprocessInput(input)
	g := ir.NewGraph()
	g.Kind = ir.C4

	if len(lines) > 0 {
		g.C4SubKind = parseC4Kind(lines[0])
		lines = lines[1:]
	}

	var boundaryStack []*ir.C4Boundary

	for _, line := range lines {
		if strings.TrimSpace(line) == "}" {
			if len(boundaryStack) > 0 {
				boundaryStack = boundaryStack[:len(boundaryStack)-1]
			}
			continue
		}

		if m := c4BoundaryRe.FindStringSubmatch(line); m != nil {
			boundaryType := m[1]
			id := strings.TrimSpace(m[2])
			label := m[3]
			bType := ""
			switch boundaryType {
			case "Enterprise_Boundary":
				bType = "Enterprise"
			case "System_Boundary":
				bType = "Software System"
			case "Container_Boundary":
				bType = "Container"
			default:
				bType = m[4]
			}

			boundary := &ir.C4Boundary{
				ID:    id,
				Label: label,
				Type:  bType,
			}
			g.C4Boundaries = append(g.C4Boundaries, boundary)
			boundaryStack = append(boundaryStack, boundary)
			continue
		}

		if m := c4ElementRe.FindStringSubmatch(line); m != nil {
			elemType := parseC4ElementType(m[1])
			args := parseC4Args(m[2])
			if len(args) < 2 {
				continue
			}
			elem := &ir.C4Element{
				ID:    args[0],
				Label: args[1],
				Type:  elemType,
			}
			if len(args) > 2 {
				if elemType.IsPerson() {
					elem.Description = args[2]
				} else {
					elem.Technology = args[2]
					if len(args) > 3 {
						elem.Description = args[3]
					}
				}
			}

			if len(boundaryStack) > 0 {
				parent := boundaryStack[len(boundaryStack)-1]
				elem.BoundaryID = parent.ID
				parent.Children = append(parent.Children, elem.ID)
			}

			g.C4Elements = append(g.C4Elements, elem)
			label := elem.Label
			g.EnsureNode(elem.ID, &label, nil)
			continue
		}

		if m := c4RelRe.FindStringSubmatch(line); m != nil {
			args := parseC4Args(m[2])
			if len(args) < 3 {
				continue
			}
			rel := &ir.C4Rel{
				From:  args[0],
				To:    args[1],
				Label: args[2],
			}
			if len(args) > 3 {
				rel.Technology = args[3]
			}
			g.C4Rels = append(g.C4Rels, rel)
			relLabel := rel.Label
			g.Edges = append(g.Edges, &ir.Edge{
				From:     rel.From,
				To:       rel.To,
				Label:    &relLabel,
				Directed: true,
				ArrowEnd: true,
			})
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}

func parseC4Kind(line string) ir.C4Kind {
	lower := strings.ToLower(strings.TrimSpace(line))
	switch {
	case strings.HasPrefix(lower, "c4container"):
		return ir.C4Container
	case strings.HasPrefix(lower, "c4component"):
		return ir.C4Component
	case strings.HasPrefix(lower, "c4dynamic"):
		return ir.C4Dynamic
	case strings.HasPrefix(lower, "c4deployment"):
		return ir.C4Deployment
	default:
		return ir.C4Context
	}
}

func parseC4ElementType(s string) ir.C4ElementType {
	switch s {
	case "Person":
		return ir.C4Person
	case "Person_Ext":
		return ir.C4ExternalPerson
	case "System":
		return ir.C4System
	case "System_Ext":
		return ir.C4ExternalSystem
	case "SystemDb":
		return ir.C4SystemDb
	case "SystemDb_Ext":
		return ir.C4ExternalSystemDb
	case "SystemQueue":
		return ir.C4SystemQueue
	case "SystemQueue_Ext":
		return ir.C4ExternalSystemQueue
	case "Container":
		return ir.C4ContainerPlain
	case "Container_Ext":
		return ir.C4ExternalContainer
	case "ContainerDb":
		return ir.C4ContainerDb
	case "ContainerDb_Ext":
		return ir.C4ExternalContainerDb
	case "ContainerQueue":
		return ir.C4ContainerQueue
	case "ContainerQueue_Ext":
		return ir.C4ExternalContainerQueue
	case "Component":
		return ir.C4ComponentPlain
	case "Component_Ext":
		return ir.C4ExternalComponent
	default:
		return ir.C4System
	}
}

func parseC4Args(s string) []string {
	var args []string
	var current strings.Builder
	inQuote := false

	for _, ch := range s {
		switch {
		case ch == '"':
			inQuote = !inQuote
		case ch == ',' && !inQuote:
			args = append(args, strings.TrimSpace(current.String()))
			current.Reset()
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, strings.TrimSpace(current.String()))
	}
	return args
}
