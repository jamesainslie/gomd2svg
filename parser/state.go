package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	stateHeaderRe       = regexp.MustCompile(`(?i)^stateDiagram`)
	stateTransitionRe   = regexp.MustCompile(`^(\S+)\s*-->\s*(\S+)(\s*:\s*(.+))?$`)
	stateDescRe         = regexp.MustCompile(`^(\w+)\s*:\s*(.+)$`)
	stateAsRe           = regexp.MustCompile(`^state\s+"([^"]+)"\s+as\s+(\w+)$`)
	stateCompositeRe    = regexp.MustCompile(`^state\s+(\w+)\s*\{$`)
	stateAnnotRe        = regexp.MustCompile(`^state\s+(\w+)\s+<<(\w+)>>$`)
	stateBracketAnnotRe = regexp.MustCompile(`^state\s+(\w+)\s+\[\[(\w+)\]\]$`)
	stateNoteInlineRe   = regexp.MustCompile(`^note\s+(right of|left of)\s+(\w+)\s*:\s*(.+)$`)
	stateNoteBlockRe    = regexp.MustCompile(`^note\s+(right of|left of)\s+(\w+)\s*$`)
)

// parseState parses a Mermaid state diagram.
func parseState(input string) (*ParseOutput, error) {
	graph := ir.NewGraph()
	graph.Kind = ir.State

	lines := preprocessInput(input)

	// Filter out the header line
	var bodyLines []string
	for _, line := range lines {
		if stateHeaderRe.MatchString(line) {
			continue
		}
		bodyLines = append(bodyLines, line)
	}

	parseStateBody(bodyLines, graph)

	// Validate brace balance: count open and close braces across all body lines.
	depth := 0
	for _, line := range bodyLines {
		for _, ch := range line {
			switch ch {
			case '{':
				depth++
			case '}':
				depth--
			}
		}
	}
	if depth > 0 {
		return nil, &ParseError{
			Diagram: "state",
			Message: "unclosed composite state (missing \"}\")",
		}
	}
	if depth < 0 {
		return nil, &ParseError{
			Diagram: "state",
			Message: "unexpected \"}\" without matching \"{\"",
		}
	}

	return &ParseOutput{Graph: graph}, nil
}

// parseStateBody parses the body lines of a state diagram into the given graph.
func parseStateBody(lines []string, graph *ir.Graph) {
	i := 0
	for i < len(lines) {
		line := lines[i]

		// Direction
		if dir, ok := parseDirectionLine(line); ok {
			graph.Direction = dir
			i++
			continue
		}

		// Note (inline): note right of State1 : text
		if caps := stateNoteInlineRe.FindStringSubmatch(line); caps != nil {
			graph.Notes = append(graph.Notes, &ir.DiagramNote{
				Position: caps[1],
				Target:   caps[2],
				Text:     strings.TrimSpace(caps[3]),
			})
			i++
			continue
		}

		// Note (block): note right of State1 ... end note
		if caps := stateNoteBlockRe.FindStringSubmatch(line); caps != nil {
			position := caps[1]
			target := caps[2]
			var noteLines []string
			i++
			for i < len(lines) {
				if strings.TrimSpace(lines[i]) == "end note" {
					i++
					break
				}
				noteLines = append(noteLines, lines[i])
				i++
			}
			graph.Notes = append(graph.Notes, &ir.DiagramNote{
				Position: position,
				Target:   target,
				Text:     strings.Join(noteLines, "\n"),
			})
			continue
		}

		// State annotation with angle brackets: state name <<choice>>
		if caps := stateAnnotRe.FindStringSubmatch(line); caps != nil {
			name := caps[1]
			annType := strings.ToLower(caps[2])
			if ann, ok := parseAnnotationType(annType); ok {
				graph.StateAnnotations[name] = ann
				graph.EnsureNode(name, nil, nil)
			}
			i++
			continue
		}

		// State annotation with brackets: state name [[fork]]
		if caps := stateBracketAnnotRe.FindStringSubmatch(line); caps != nil {
			name := caps[1]
			annType := strings.ToLower(caps[2])
			if ann, ok := parseAnnotationType(annType); ok {
				graph.StateAnnotations[name] = ann
				graph.EnsureNode(name, nil, nil)
			}
			i++
			continue
		}

		// Composite state: state Name {
		if caps := stateCompositeRe.FindStringSubmatch(line); caps != nil {
			stateName := caps[1]
			// Collect inner lines until matching }
			innerLines, endIdx := collectBraceBlock(lines, i+1)
			i = endIdx

			// Check for concurrent regions (lines that are just "--")
			regions := splitRegions(innerLines)
			if len(regions) > 1 {
				// Concurrent state with regions
				cs := &ir.CompositeState{
					ID:    stateName,
					Label: stateName,
				}
				for _, regionLines := range regions {
					regionGraph := ir.NewGraph()
					regionGraph.Kind = ir.State
					parseStateBody(regionLines, regionGraph)
					cs.Regions = append(cs.Regions, regionGraph)
				}
				graph.CompositeStates[stateName] = cs
				graph.EnsureNode(stateName, nil, nil)
			} else {
				// Simple composite state
				innerGraph := ir.NewGraph()
				innerGraph.Kind = ir.State
				parseStateBody(innerLines, innerGraph)
				cs := &ir.CompositeState{
					ID:    stateName,
					Label: stateName,
					Inner: innerGraph,
				}
				graph.CompositeStates[stateName] = cs
				graph.EnsureNode(stateName, nil, nil)
			}
			continue
		}

		// State "description" as alias
		if caps := stateAsRe.FindStringSubmatch(line); caps != nil {
			desc := caps[1]
			alias := caps[2]
			graph.EnsureNode(alias, nil, nil)
			graph.StateDescriptions[alias] = desc
			i++
			continue
		}

		// Transition: A --> B or A --> B : label
		if caps := stateTransitionRe.FindStringSubmatch(line); caps != nil {
			from := mapStarToken(caps[1], true)
			to := mapStarToken(caps[2], false)
			edge := &ir.Edge{
				From:     from,
				To:       to,
				Directed: true,
				ArrowEnd: true,
			}
			if caps[4] != "" {
				label := strings.TrimSpace(caps[4])
				edge.Label = &label
			}
			graph.EnsureNode(from, nil, nil)
			graph.EnsureNode(to, nil, nil)
			graph.Edges = append(graph.Edges, edge)
			i++
			continue
		}

		// State description: stateId : description
		if caps := stateDescRe.FindStringSubmatch(line); caps != nil {
			stateID := caps[1]
			desc := strings.TrimSpace(caps[2])
			// Don't match "state" keyword lines or "note" lines
			if stateID != "state" && stateID != "note" {
				graph.StateDescriptions[stateID] = desc
				graph.EnsureNode(stateID, nil, nil)
				i++
				continue
			}
		}

		// Standalone node reference
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.Contains(trimmed, " ") {
			graph.EnsureNode(trimmed, nil, nil)
		}

		i++
	}
}

// mapStarToken converts [*] to __start__ or __end__ based on context.
func mapStarToken(token string, isSource bool) string {
	if token == "[*]" {
		if isSource {
			return "__start__"
		}
		return "__end__"
	}
	return token
}

// parseAnnotationType converts an annotation string to a StateAnnotation.
func parseAnnotationType(annType string) (ir.StateAnnotation, bool) {
	switch annType {
	case "choice":
		return ir.StateChoice, true
	case "fork":
		return ir.StateFork, true
	case "join":
		return ir.StateJoin, true
	default:
		return ir.StateChoice, false
	}
}

// collectBraceBlock collects lines from startIdx until the matching closing brace.
// Returns the inner lines (excluding the closing brace line) and the index after the closing brace.
func collectBraceBlock(lines []string, startIdx int) ([]string, int) {
	depth := 1
	var inner []string
	i := startIdx
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Count braces
		for _, ch := range trimmed {
			switch ch {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					// If there's content before the }, include it
					beforeBrace := strings.TrimSpace(strings.TrimRight(trimmed, "}"))
					if beforeBrace != "" {
						inner = append(inner, beforeBrace)
					}
					return inner, i + 1
				}
			}
		}

		if depth > 0 {
			inner = append(inner, line)
		}
		i++
	}
	return inner, i
}

// splitRegions splits lines by the "--" separator into concurrent regions.
func splitRegions(lines []string) [][]string {
	var regions [][]string
	var current []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "--" {
			regions = append(regions, current)
			current = nil
			continue
		}
		current = append(current, line)
	}

	if len(current) > 0 || len(regions) > 0 {
		regions = append(regions, current)
	}

	return regions
}
