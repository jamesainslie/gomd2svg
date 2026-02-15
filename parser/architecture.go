package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	archGroupRe    = regexp.MustCompile(`^group\s+(\w+)(?:\(([^)]*)\))?\[([^\]]+)\](?:\s+in\s+(\w+))?$`)
	archServiceRe  = regexp.MustCompile(`^service\s+(\w+)(?:\(([^)]*)\))?\[([^\]]+)\](?:\s+in\s+(\w+))?$`)
	archJunctionRe = regexp.MustCompile(`^junction\s+(\w+)(?:\s+in\s+(\w+))?$`)
	archEdgeRe     = regexp.MustCompile(`^(\w+)(?:\{group\})?:(L|R|T|B)\s*(<)?--(>)?\s*(L|R|T|B):(\w+)(?:\{group\})?$`)
)

func parseArchSide(s string) ir.ArchSide {
	switch s {
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

func parseArchitecture(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Architecture

	lines := preprocessInput(input)
	groupChildren := make(map[string][]string)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if lower == "architecture-beta" || lower == "architecture" {
			continue
		}

		// Try group
		if m := archGroupRe.FindStringSubmatch(trimmed); m != nil {
			grp := &ir.ArchGroup{
				ID:       m[1],
				Icon:     m[2],
				Label:    m[3],
				ParentID: m[4],
			}
			g.ArchGroups = append(g.ArchGroups, grp)
			if grp.ParentID != "" {
				groupChildren[grp.ParentID] = append(groupChildren[grp.ParentID], grp.ID)
			}
			continue
		}

		// Try service
		if m := archServiceRe.FindStringSubmatch(trimmed); m != nil {
			svc := &ir.ArchService{
				ID:      m[1],
				Icon:    m[2],
				Label:   m[3],
				GroupID: m[4],
			}
			g.ArchServices = append(g.ArchServices, svc)
			label := svc.Label
			g.EnsureNode(svc.ID, &label, nil)
			if svc.GroupID != "" {
				groupChildren[svc.GroupID] = append(groupChildren[svc.GroupID], svc.ID)
			}
			continue
		}

		// Try junction
		if m := archJunctionRe.FindStringSubmatch(trimmed); m != nil {
			junc := &ir.ArchJunction{
				ID:      m[1],
				GroupID: m[2],
			}
			g.ArchJunctions = append(g.ArchJunctions, junc)
			g.EnsureNode(junc.ID, nil, nil)
			if junc.GroupID != "" {
				groupChildren[junc.GroupID] = append(groupChildren[junc.GroupID], junc.ID)
			}
			continue
		}

		// Try edge
		if m := archEdgeRe.FindStringSubmatch(trimmed); m != nil {
			edge := &ir.ArchEdge{
				FromID:     m[1],
				FromSide:   parseArchSide(m[2]),
				ArrowLeft:  m[3] == "<",
				ArrowRight: m[4] == ">",
				ToSide:     parseArchSide(m[5]),
				ToID:       m[6],
			}
			g.ArchEdges = append(g.ArchEdges, edge)
			g.Edges = append(g.Edges, &ir.Edge{
				From:       edge.FromID,
				To:         edge.ToID,
				Directed:   edge.ArrowRight,
				ArrowEnd:   edge.ArrowRight,
				ArrowStart: edge.ArrowLeft,
			})
			continue
		}
	}

	// Populate group Children from accumulated map
	for _, grp := range g.ArchGroups {
		grp.Children = groupChildren[grp.ID]
	}

	return &ParseOutput{Graph: g}, nil
}
