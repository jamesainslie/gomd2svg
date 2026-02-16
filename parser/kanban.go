package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

// indentedLine holds a preprocessed line with its indentation level.
type indentedLine struct {
	text   string
	indent int
}

// kanbanCardRe matches id[Label] with optional @{...} metadata.
var kanbanCardRe = regexp.MustCompile(`^(\w+)\[([^\]]+)\](?:\s*@\{(.+)\})?$`)

// kanbanColumnHeaderRe matches id[Label] for column headers.
var kanbanColumnHeaderRe = regexp.MustCompile(`^(\w+)\[([^\]]+)\]$`)

func parseKanban(input string) (*ParseOutput, error) { //nolint:unparam // error return is part of the parser interface contract used by Parse().
	graph := ir.NewGraph()
	graph.Kind = ir.Kanban

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
			// New column.
			colIndent = indent
			colID, colLabel := parseKanbanColumnHeader(trimmed)
			currentCol = &ir.KanbanColumn{ID: colID, Label: colLabel}
			graph.Columns = append(graph.Columns, currentCol)
			continue
		}

		// Card line (more indented than column).
		if currentCol != nil {
			card := parseKanbanCard(trimmed)
			if card != nil {
				currentCol.Cards = append(currentCol.Cards, card)
			}
		}
	}

	return &ParseOutput{Graph: graph}, nil
}

// preprocessKanbanInput splits the input into lines, strips comments and blank
// lines, but preserves leading whitespace so indentation can be measured.
// Each returned indentedLine has the cleaned text and the number of leading spaces.
func preprocessKanbanInput(input string) []indentedLine {
	var result []indentedLine
	for _, rawLine := range strings.Split(input, "\n") {
		// Count leading whitespace (tabs count as 1 indent unit each,
		// spaces count as 1 each -- consistent with mermaid.js behavior).
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

		result = append(result, indentedLine{
			text:   without,
			indent: indent,
		})
	}
	return result
}

// parseKanbanColumnHeader parses a column header line.
// It handles both bare "ColumnName" and "id[Column Label]" syntax.
func parseKanbanColumnHeader(line string) (colID, colLabel string) { //nolint:nonamedreturns // named returns clarify the multi-value return.
	if match := kanbanColumnHeaderRe.FindStringSubmatch(line); match != nil {
		return match[1], match[2]
	}
	// Bare column name -- ID and label are both the trimmed line.
	trimmed := strings.TrimSpace(line)
	return trimmed, trimmed
}

// parseKanbanCard parses a card line like "id[Label]" or "id[Label]@{ key: 'val' }".
func parseKanbanCard(line string) *ir.KanbanCard {
	match := kanbanCardRe.FindStringSubmatch(line)
	if match == nil {
		return nil
	}

	card := &ir.KanbanCard{
		ID:    match[1],
		Label: match[2],
	}

	// Parse optional metadata block.
	if match[3] != "" {
		assigned, ticket, icon, description, priority := parseKanbanMetadata(match[3])
		card.Assigned = assigned
		card.Ticket = ticket
		card.Icon = icon
		card.Description = description
		card.Priority = priority
	}

	return card
}

// kanbanMetadataResult holds the parsed metadata for a kanban card.
type kanbanMetadataResult struct {
	Assigned    string
	Ticket      string
	Icon        string
	Description string
	Priority    ir.KanbanPriority
}

// parseKanbanMetadata parses the key-value pairs inside @{ ... }.
// Format: key: 'value', key2: 'value2'
// Values use single quotes. Keys are unquoted identifiers.
func parseKanbanMetadata(raw string) (string, string, string, string, ir.KanbanPriority) { //nolint:revive // five return values needed for backward compatibility.
	result := kanbanMetadataResult{}
	// Scan through the raw string extracting key: 'value' pairs.
	str := strings.TrimSpace(raw)
	for len(str) > 0 {
		// Skip whitespace and commas.
		str = strings.TrimLeft(str, " \t,")
		if len(str) == 0 {
			break
		}

		// Extract key (everything up to ':').
		colonIdx := strings.Index(str, ":")
		if colonIdx < 0 {
			break
		}
		key := strings.TrimSpace(str[:colonIdx])
		str = str[colonIdx+1:]

		// Skip whitespace before the value.
		str = strings.TrimLeft(str, " \t")
		if len(str) == 0 {
			break
		}

		var value string
		if str[0] == '\'' {
			// Single-quoted value -- find the closing quote.
			endIdx := strings.Index(str[1:], "'")
			if endIdx < 0 {
				// Unterminated quote, take rest.
				value = str[1:]
				str = ""
			} else {
				value = str[1 : endIdx+1]
				str = str[endIdx+2:]
			}
		} else {
			// Unquoted value -- read until comma or end.
			commaIdx := strings.Index(str, ",")
			if commaIdx < 0 {
				value = strings.TrimSpace(str)
				str = ""
			} else {
				value = strings.TrimSpace(str[:commaIdx])
				str = str[commaIdx+1:]
			}
		}

		switch strings.ToLower(key) {
		case "assigned":
			result.Assigned = value
		case "ticket":
			result.Ticket = value
		case "icon":
			result.Icon = value
		case "description":
			result.Description = value
		case "priority":
			result.Priority = parsePriorityValue(value)
		}
	}

	return result.Assigned, result.Ticket, result.Icon, result.Description, result.Priority
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
