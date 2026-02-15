package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

// kanbanLine holds a preprocessed line with its indentation level.
type kanbanLine struct {
	text   string
	indent int
}

// kanbanCardRe matches id[Label] with optional @{...} metadata.
var kanbanCardRe = regexp.MustCompile(`^(\w+)\[([^\]]+)\](?:\s*@\{(.+)\})?$`)

// kanbanColumnHeaderRe matches id[Label] for column headers.
var kanbanColumnHeaderRe = regexp.MustCompile(`^(\w+)\[([^\]]+)\]$`)

func parseKanban(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Kanban

	lines := preprocessKanbanInput(input)

	var currentCol *ir.KanbanColumn
	colIndent := -1

	for _, entry := range lines {
		line := entry.text
		indent := entry.indent

		// Skip the kanban header line.
		trimmed := strings.TrimSpace(line)
		if strings.EqualFold(trimmed, "kanban") {
			continue
		}
		if trimmed == "" {
			continue
		}

		// Determine if this is a column or card based on indentation.
		// The first non-header line sets the column indent level.
		if colIndent < 0 || indent <= colIndent {
			// New column
			colIndent = indent
			id, label := parseKanbanColumnHeader(trimmed)
			currentCol = &ir.KanbanColumn{ID: id, Label: label}
			g.Columns = append(g.Columns, currentCol)
			continue
		}

		// Card line (more indented than column)
		if currentCol != nil {
			card := parseKanbanCard(trimmed)
			if card != nil {
				currentCol.Cards = append(currentCol.Cards, card)
			}
		}
	}

	return &ParseOutput{Graph: g}, nil
}

// preprocessKanbanInput splits the input into lines, strips comments and blank
// lines, but preserves leading whitespace so indentation can be measured.
// Each returned kanbanLine has the cleaned text and the number of leading spaces.
func preprocessKanbanInput(input string) []kanbanLine {
	var result []kanbanLine
	for _, rawLine := range strings.Split(input, "\n") {
		// Count leading whitespace (tabs count as 1 indent unit each,
		// spaces count as 1 each — consistent with mermaid.js behavior).
		indent := 0
		for _, ch := range rawLine {
			switch ch {
			case ' ':
				indent++
			case '\t':
				indent += 2 // treat tab as 2 spaces
			default:
				goto done
			}
		}
	done:

		trimmed := strings.TrimSpace(rawLine)
		if trimmed == "" {
			continue
		}
		// Skip full-line comments.
		if strings.HasPrefix(trimmed, "%%") {
			continue
		}
		// Strip trailing comments.
		without := stripTrailingComment(trimmed)
		if without == "" {
			continue
		}

		result = append(result, kanbanLine{
			text:   without,
			indent: indent,
		})
	}
	return result
}

// parseKanbanColumnHeader parses a column header line.
// It handles both bare "ColumnName" and "id[Column Label]" syntax.
func parseKanbanColumnHeader(line string) (id, label string) {
	if m := kanbanColumnHeaderRe.FindStringSubmatch(line); m != nil {
		return m[1], m[2]
	}
	// Bare column name — ID and label are both the trimmed line.
	trimmed := strings.TrimSpace(line)
	return trimmed, trimmed
}

// parseKanbanCard parses a card line like "id[Label]" or "id[Label]@{ key: 'val' }".
func parseKanbanCard(line string) *ir.KanbanCard {
	m := kanbanCardRe.FindStringSubmatch(line)
	if m == nil {
		return nil
	}

	card := &ir.KanbanCard{
		ID:    m[1],
		Label: m[2],
	}

	// Parse optional metadata block.
	if m[3] != "" {
		assigned, ticket, icon, description, priority := parseKanbanMetadata(m[3])
		card.Assigned = assigned
		card.Ticket = ticket
		card.Icon = icon
		card.Description = description
		card.Priority = priority
	}

	return card
}

// parseKanbanMetadata parses the key-value pairs inside @{ ... }.
// Format: key: 'value', key2: 'value2'
// Values use single quotes. Keys are unquoted identifiers.
func parseKanbanMetadata(raw string) (assigned, ticket, icon, description string, priority ir.KanbanPriority) {
	// Scan through the raw string extracting key: 'value' pairs.
	s := strings.TrimSpace(raw)
	for len(s) > 0 {
		// Skip whitespace and commas.
		s = strings.TrimLeft(s, " \t,")
		if len(s) == 0 {
			break
		}

		// Extract key (everything up to ':').
		colonIdx := strings.Index(s, ":")
		if colonIdx < 0 {
			break
		}
		key := strings.TrimSpace(s[:colonIdx])
		s = s[colonIdx+1:]

		// Skip whitespace before the value.
		s = strings.TrimLeft(s, " \t")
		if len(s) == 0 {
			break
		}

		var value string
		if s[0] == '\'' {
			// Single-quoted value — find the closing quote.
			endIdx := strings.Index(s[1:], "'")
			if endIdx < 0 {
				// Unterminated quote, take rest.
				value = s[1:]
				s = ""
			} else {
				value = s[1 : endIdx+1]
				s = s[endIdx+2:]
			}
		} else {
			// Unquoted value — read until comma or end.
			commaIdx := strings.Index(s, ",")
			if commaIdx < 0 {
				value = strings.TrimSpace(s)
				s = ""
			} else {
				value = strings.TrimSpace(s[:commaIdx])
				s = s[commaIdx+1:]
			}
		}

		switch strings.ToLower(key) {
		case "assigned":
			assigned = value
		case "ticket":
			ticket = value
		case "icon":
			icon = value
		case "description":
			description = value
		case "priority":
			priority = parsePriorityValue(value)
		}
	}

	return
}

// parsePriorityValue maps a priority string to an ir.KanbanPriority.
// Matching is case-insensitive.
func parsePriorityValue(val string) ir.KanbanPriority {
	switch strings.ToLower(val) {
	case "very high":
		return ir.PriorityVeryHigh
	case "high":
		return ir.PriorityHigh
	case "low":
		return ir.PriorityLow
	case "very low":
		return ir.PriorityVeryLow
	default:
		return ir.PriorityNone
	}
}
