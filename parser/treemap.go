package parser

import (
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

func parseTreemap(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Treemap

	lines := preprocessMindmapInput(input) // reuse indentation-aware preprocessor
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	type stackEntry struct {
		node   *ir.TreemapNode
		indent int
	}
	var stack []stackEntry

	for _, entry := range lines {
		text := entry.text
		indent := entry.indent
		lower := strings.ToLower(strings.TrimSpace(text))

		// Skip keyword line.
		if strings.HasPrefix(lower, "treemap") {
			continue
		}

		// Handle title directive.
		if strings.HasPrefix(lower, "title ") || strings.HasPrefix(lower, "title\t") {
			g.TreemapTitle = strings.TrimSpace(text[6:])
			continue
		}

		// Parse node: "Label": value  or  "Label"
		label, value, hasValue, class := parseTreemapNodeLine(text)
		if label == "" {
			continue
		}

		node := &ir.TreemapNode{Label: label, Class: class}
		if hasValue {
			node.Value = value
		}

		// Pop stack to find parent.
		for len(stack) > 0 && stack[len(stack)-1].indent >= indent {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			g.TreemapRoot = node
		} else {
			parent := stack[len(stack)-1].node
			parent.Children = append(parent.Children, node)
		}

		stack = append(stack, stackEntry{node: node, indent: indent})
	}

	return &ParseOutput{Graph: g}, nil
}

// parseTreemapNodeLine parses a line like `"Label": 30` or `"Label"`.
func parseTreemapNodeLine(line string) (label string, value float64, hasValue bool, class string) {
	line = strings.TrimSpace(line)

	// Strip optional :::class decorator.
	if idx := strings.Index(line, ":::"); idx >= 0 {
		class = strings.TrimSpace(line[idx+3:])
		line = strings.TrimSpace(line[:idx])
	}

	if len(line) < 2 {
		return "", 0, false, ""
	}
	quote := line[0]
	if quote != '"' && quote != '\'' {
		return "", 0, false, ""
	}
	end := strings.IndexByte(line[1:], quote)
	if end < 0 {
		return "", 0, false, ""
	}
	label = line[1 : end+1]
	rest := strings.TrimSpace(line[end+2:])

	if len(rest) > 0 && (rest[0] == ':' || rest[0] == ',') {
		valStr := strings.TrimSpace(rest[1:])
		if v, err := strconv.ParseFloat(valStr, 64); err == nil {
			return label, v, true, class
		}
	}

	return label, 0, false, class
}
