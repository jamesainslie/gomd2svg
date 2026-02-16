package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	mindmapIconRe  = regexp.MustCompile(`::icon\(([^)]+)\)`)
	mindmapClassRe = regexp.MustCompile(`:::(\S+)`)
)

func parseMindmap(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap

	lines := preprocessMindmapInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	type stackEntry struct {
		node   *ir.MindmapNode
		indent int
	}
	var stack []stackEntry
	nodeCount := 0

	for _, entry := range lines {
		text := entry.text
		indent := entry.indent

		if strings.EqualFold(strings.TrimSpace(text), "mindmap") {
			continue
		}

		node := parseMindmapNodeText(text, nodeCount)
		nodeCount++

		// Pop stack until we find a parent with smaller indentation.
		for len(stack) > 0 && stack[len(stack)-1].indent >= indent {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			g.MindmapRoot = node
		} else {
			parent := stack[len(stack)-1].node
			parent.Children = append(parent.Children, node)
		}

		stack = append(stack, stackEntry{node: node, indent: indent})
	}

	return &ParseOutput{Graph: g}, nil
}

func parseMindmapNodeText(text string, index int) *ir.MindmapNode {
	node := &ir.MindmapNode{}

	// Extract and strip icon decorator.
	if m := mindmapIconRe.FindStringSubmatch(text); m != nil {
		node.Icon = m[1]
		text = strings.TrimSpace(mindmapIconRe.ReplaceAllString(text, ""))
	}

	// Extract and strip class decorator.
	if m := mindmapClassRe.FindStringSubmatch(text); m != nil {
		node.Class = m[1]
		text = strings.TrimSpace(mindmapClassRe.ReplaceAllString(text, ""))
	}

	text = strings.TrimSpace(text)
	node.Shape, node.Label = parseMindmapShape(text)

	// Generate an ID.
	node.ID = fmt.Sprintf("mm_%d", index)

	return node
}

func parseMindmapShape(text string) (ir.MindmapShape, string) {
	// Handle shapes with possible ID prefix: e.g., "root((Central))", "A[Square]"
	// Check for double-char delimiters first, then single-char.

	// Bang: either standalone ))text(( or id))text((
	if idx := strings.Index(text, "))"); idx >= 0 {
		if strings.HasSuffix(text, "((") {
			label := text[idx+2 : len(text)-2]
			return ir.MindmapBang, label
		}
	}
	// Circle: either standalone ((text)) or id((text))
	if idx := strings.Index(text, "(("); idx >= 0 {
		if strings.HasSuffix(text, "))") {
			label := text[idx+2 : len(text)-2]
			return ir.MindmapCircle, label
		}
	}
	// Hexagon: either standalone {{text}} or id{{text}}
	if idx := strings.Index(text, "{{"); idx >= 0 {
		if strings.HasSuffix(text, "}}") {
			label := text[idx+2 : len(text)-2]
			return ir.MindmapHexagon, label
		}
	}
	// Cloud: either standalone )text( or id)text(
	if idx := strings.Index(text, ")"); idx >= 0 && !strings.HasPrefix(text, "))") {
		if strings.HasSuffix(text, "(") && !strings.HasSuffix(text, "((") {
			label := text[idx+1 : len(text)-1]
			return ir.MindmapCloud, label
		}
	}
	// Square: either standalone [text] or id[text]
	if idx := strings.Index(text, "["); idx >= 0 {
		if strings.HasSuffix(text, "]") {
			label := text[idx+1 : len(text)-1]
			return ir.MindmapSquare, label
		}
	}
	// Rounded: either standalone (text) or id(text)
	if idx := strings.Index(text, "("); idx >= 0 && !strings.HasPrefix(text[idx:], "((") {
		if strings.HasSuffix(text, ")") && !strings.HasSuffix(text, "))") {
			label := text[idx+1 : len(text)-1]
			return ir.MindmapRounded, label
		}
	}
	// Default: bare text
	return ir.MindmapShapeDefault, text
}

// preprocessMindmapInput preserves indentation for hierarchy detection.
// Uses shared indentedLine type.
func preprocessMindmapInput(input string) []indentedLine {
	var result []indentedLine
	for _, rawLine := range strings.Split(input, "\n") {
		indent := 0
		for _, ch := range rawLine {
			switch ch {
			case ' ':
				indent++
			case '\t':
				indent += 2
			default:
				goto done
			}
		}
	done:
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
		result = append(result, indentedLine{text: without, indent: indent})
	}
	return result
}
