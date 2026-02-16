package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	archGroupRe    = regexp.MustCompile(`^group\s+(\w+)(?:\(([^)]*)\))?\[([^\]]+)\](?:\s+in\s+(\w+))?$`)
	archServiceRe  = regexp.MustCompile(`^service\s+(\w+)(?:\(([^)]*)\))?\[([^\]]+)\](?:\s+in\s+(\w+))?$`)
	archJunctionRe = regexp.MustCompile(`^junction\s+(\w+)(?:\s+in\s+(\w+))?$`)
	archEdgeRe     = regexp.MustCompile(`^(\w+)(?:\{group\})?:(L|R|T|B)\s*(<)?--(>)?\s*(L|R|T|B):(\w+)(?:\{group\})?$`)
)

func parseArchSide(side string) ir.ArchSide {
	switch side {
	case "R":
		return ir.ArchRight
	case "T":
		return ir.ArchTop
	case "B":
		return ir.ArchBottom
	default:
		return ir.ArchLeft
	}
}

//nolint:unparam // error return is part of the parser interface contract used by Parse().
func parseArchitecture(input string) (*ParseOutput, error) {
	graph := ir.NewGraph()
	graph.Kind = ir.Architecture

	lines := preprocessInput(input)
	groupChildren := make(map[string][]string)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if lower == "architecture-beta" || lower == "architecture" {
			continue
		}

		// Try group.
		if match := archGroupRe.FindStringSubmatch(trimmed); match != nil {
			grp := &ir.ArchGroup{
				ID:       match[1],
				Icon:     match[2],
				Label:    match[3],
				ParentID: match[4],
			}
			graph.ArchGroups = append(graph.ArchGroups, grp)
			if grp.ParentID != "" {
				groupChildren[grp.ParentID] = append(groupChildren[grp.ParentID], grp.ID)
			}
			continue
		}

		// Try service.
		if match := archServiceRe.FindStringSubmatch(trimmed); match != nil {
			svc := &ir.ArchService{
				ID:      match[1],
				Icon:    match[2],
				Label:   match[3],
				GroupID: match[4],
			}
			graph.ArchServices = append(graph.ArchServices, svc)
			label := svc.Label
			graph.EnsureNode(svc.ID, &label, nil)
			if svc.GroupID != "" {
				groupChildren[svc.GroupID] = append(groupChildren[svc.GroupID], svc.ID)
			}
			continue
		}

		// Try junction.
		if match := archJunctionRe.FindStringSubmatch(trimmed); match != nil {
			junc := &ir.ArchJunction{
				ID:      match[1],
				GroupID: match[2],
			}
			graph.ArchJunctions = append(graph.ArchJunctions, junc)
			graph.EnsureNode(junc.ID, nil, nil)
			if junc.GroupID != "" {
				groupChildren[junc.GroupID] = append(groupChildren[junc.GroupID], junc.ID)
			}
			continue
		}

		// Try edge.
		if match := archEdgeRe.FindStringSubmatch(trimmed); match != nil {
			edge := &ir.ArchEdge{
				FromID:     match[1],
				FromSide:   parseArchSide(match[2]),
				ArrowLeft:  match[3] == "<",
				ArrowRight: match[4] == ">",
				ToSide:     parseArchSide(match[5]),
				ToID:       match[6],
			}
			graph.ArchEdges = append(graph.ArchEdges, edge)
			graph.Edges = append(graph.Edges, &ir.Edge{
				From:       edge.FromID,
				To:         edge.ToID,
				Directed:   edge.ArrowRight,
				ArrowEnd:   edge.ArrowRight,
				ArrowStart: edge.ArrowLeft,
			})
			continue
		}
	}

	// Populate group Children from accumulated map.
	for _, grp := range graph.ArchGroups {
		grp.Children = groupChildren[grp.ID]
	}

	return &ParseOutput{Graph: graph}, nil
}
