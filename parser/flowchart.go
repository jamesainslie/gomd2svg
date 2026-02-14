package parser

import (
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

// parseFlowchart parses a Mermaid flowchart/graph diagram.
func parseFlowchart(input string) (*ParseOutput, error) {
	graph := ir.NewGraph()
	graph.Kind = ir.Flowchart
	var subgraphStack []int

	lines := preprocessInput(input)

	for _, rawLine := range lines {
		for _, line := range splitStatements(rawLine) {
			if line == "" {
				continue
			}

			// Try header line (flowchart LR, graph TD, etc.)
			if caps := headerRe.FindStringSubmatch(line); caps != nil {
				if dir, ok := ir.DirectionFromToken(caps[2]); ok {
					graph.Direction = dir
				}
				continue
			}

			// Try "end" to close a subgraph.
			if line == "end" {
				if len(subgraphStack) > 0 {
					subgraphStack = subgraphStack[:len(subgraphStack)-1]
				}
				continue
			}

			// Try subgraph declaration.
			if caps := subgraphRe.FindStringSubmatch(line); caps != nil {
				rest := caps[1]
				id, label, _ := parseSubgraphHeader(rest)
				sg := &ir.Subgraph{
					ID:    id,
					Label: label,
				}
				graph.Subgraphs = append(graph.Subgraphs, sg)
				subgraphStack = append(subgraphStack, len(graph.Subgraphs)-1)
				continue
			}

			// Try direction line inside subgraph.
			if dir, ok := parseDirectionLine(line); ok {
				if len(subgraphStack) > 0 {
					idx := subgraphStack[len(subgraphStack)-1]
					if idx < len(graph.Subgraphs) {
						graph.Subgraphs[idx].Direction = &dir
					}
				} else {
					graph.Direction = dir
				}
				continue
			}

			// Skip classDef, class, style, linkStyle, click, accTitle, accDescr, title
			lowerLine := strings.ToLower(line)
			if strings.HasPrefix(lowerLine, "classdef") ||
				strings.HasPrefix(lowerLine, "class ") ||
				strings.HasPrefix(lowerLine, "style ") ||
				strings.HasPrefix(lowerLine, "linkstyle") ||
				strings.HasPrefix(lowerLine, "click ") ||
				strings.HasPrefix(lowerLine, "acctitle") ||
				strings.HasPrefix(lowerLine, "accdescr") ||
				strings.HasPrefix(lowerLine, "title ") {
				continue
			}

			// Try edge chain (A-->B-->C).
			if chainLines := splitEdgeChain(line); chainLines != nil {
				added := false
				for _, edgeLine := range chainLines {
					if addFlowchartEdge(edgeLine, graph, subgraphStack) {
						added = true
					}
				}
				if added {
					continue
				}
			}

			// Try single edge.
			if addFlowchartEdge(line, graph, subgraphStack) {
				continue
			}

			// Fallback: standalone node.
			if nodeID, nodeLabel, nodeShape, _, ok := parseNodeOnly(line); ok {
				graph.EnsureNode(nodeID, nodeLabel, nodeShape)
				addNodeToSubgraphs(graph, subgraphStack, nodeID)
			}
		}
	}

	return &ParseOutput{Graph: graph}, nil
}

// addFlowchartEdge parses a single edge line and adds the edge(s) to the graph.
func addFlowchartEdge(line string, graph *ir.Graph, subgraphStack []int) bool {
	left, label, right, meta, ok := parseEdgeLine(line)
	if !ok {
		return false
	}

	// Split on & for multi-source/target.
	sources := splitAndTrim(left)
	targets := splitAndTrim(right)

	var sourceIDs []string
	for _, source := range sources {
		id, lbl, shape, _ := parseNodeToken(source)
		graph.EnsureNode(id, lbl, shape)
		addNodeToSubgraphs(graph, subgraphStack, id)
		sourceIDs = append(sourceIDs, id)
	}

	var targetIDs []string
	for _, target := range targets {
		id, lbl, shape, _ := parseNodeToken(target)
		graph.EnsureNode(id, lbl, shape)
		addNodeToSubgraphs(graph, subgraphStack, id)
		targetIDs = append(targetIDs, id)
	}

	for _, srcID := range sourceIDs {
		for _, tgtID := range targetIDs {
			edge := &ir.Edge{
				From:            srcID,
				To:              tgtID,
				Label:           label,
				Directed:        meta.directed,
				ArrowStart:      meta.arrowStart,
				ArrowEnd:        meta.arrowEnd,
				ArrowStartKind:  meta.arrowStartKind,
				ArrowEndKind:    meta.arrowEndKind,
				StartDecoration: meta.startDecoration,
				EndDecoration:   meta.endDecoration,
				Style:           meta.style,
			}
			graph.Edges = append(graph.Edges, edge)
		}
	}

	return true
}

// splitAndTrim splits on & and trims whitespace, filtering empty parts.
func splitAndTrim(s string) []string {
	parts := strings.Split(s, "&")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// addNodeToSubgraphs adds a node to all subgraphs in the current stack.
func addNodeToSubgraphs(graph *ir.Graph, subgraphStack []int, nodeID string) {
	for _, idx := range subgraphStack {
		addNodeToSubgraph(graph, idx, nodeID)
	}
}

// addNodeToSubgraph adds a node to a specific subgraph, avoiding duplicates.
func addNodeToSubgraph(graph *ir.Graph, idx int, nodeID string) {
	if idx >= len(graph.Subgraphs) {
		return
	}
	sg := graph.Subgraphs[idx]
	for _, existing := range sg.Nodes {
		if existing == nodeID {
			return
		}
	}
	sg.Nodes = append(sg.Nodes, nodeID)
}
