package parser

import (
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

// ParseOutput holds the result of parsing a Mermaid diagram.
type ParseOutput struct {
	Graph      *ir.Graph
	InitConfig map[string]interface{}
}

// Parse detects the diagram kind and dispatches to the appropriate parser.
func Parse(input string) (*ParseOutput, error) {
	kind := detectDiagramKind(input)
	switch kind {
	case ir.Flowchart:
		return parseFlowchart(input)
	case ir.Class:
		return parseClass(input)
	case ir.State:
		return parseState(input)
	case ir.Er:
		return parseER(input)
	case ir.Sequence:
		return parseSequence(input)
	default:
		return parseFlowchart(input)
	}
}

// detectDiagramKind scans lines, skipping comments and empty lines,
// and matches the first keyword case-insensitively.
func detectDiagramKind(input string) ir.DiagramKind {
	for _, rawLine := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(rawLine)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "%%") {
			continue
		}
		without := stripTrailingComment(trimmed)
		if without == "" {
			continue
		}
		lower := strings.ToLower(without)

		if strings.HasPrefix(lower, "sequencediagram") {
			return ir.Sequence
		}
		if strings.HasPrefix(lower, "classdiagram") {
			return ir.Class
		}
		if strings.HasPrefix(lower, "statediagram") {
			return ir.State
		}
		if strings.HasPrefix(lower, "erdiagram") {
			return ir.Er
		}
		if strings.HasPrefix(lower, "pie") {
			return ir.Pie
		}
		if strings.HasPrefix(lower, "mindmap") {
			return ir.Mindmap
		}
		if strings.HasPrefix(lower, "journey") {
			return ir.Journey
		}
		if strings.HasPrefix(lower, "timeline") {
			return ir.Timeline
		}
		if strings.HasPrefix(lower, "gantt") {
			return ir.Gantt
		}
		if strings.HasPrefix(lower, "requirementdiagram") {
			return ir.Requirement
		}
		if strings.HasPrefix(lower, "gitgraph") {
			return ir.GitGraph
		}
		if strings.HasPrefix(lower, "c4") {
			return ir.C4
		}
		if strings.HasPrefix(lower, "sankey") {
			return ir.Sankey
		}
		if strings.HasPrefix(lower, "quadrantchart") {
			return ir.Quadrant
		}
		if strings.HasPrefix(lower, "zenuml") {
			return ir.ZenUML
		}
		if strings.HasPrefix(lower, "block") {
			return ir.Block
		}
		if strings.HasPrefix(lower, "packet") {
			return ir.Packet
		}
		if strings.HasPrefix(lower, "kanban") {
			return ir.Kanban
		}
		if strings.HasPrefix(lower, "architecture") {
			return ir.Architecture
		}
		if strings.HasPrefix(lower, "radar") {
			return ir.Radar
		}
		if strings.HasPrefix(lower, "treemap") {
			return ir.Treemap
		}
		if strings.HasPrefix(lower, "xychart") {
			return ir.XYChart
		}
		if strings.HasPrefix(lower, "flowchart") || strings.HasPrefix(lower, "graph") {
			return ir.Flowchart
		}
		// First non-comment, non-empty line didn't match any keyword.
		// Default to Flowchart.
		break
	}
	return ir.Flowchart
}

// preprocessInput filters out comments and empty lines, strips trailing comments,
// and returns the cleaned lines.
func preprocessInput(input string) []string {
	var lines []string
	for _, rawLine := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(rawLine)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "%%") {
			continue
		}
		without := stripTrailingComment(trimmed)
		if without == "" {
			continue
		}
		lines = append(lines, without)
	}
	return lines
}
