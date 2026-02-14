package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	erRelRe = regexp.MustCompile(
		`^(\S+)\s+(\|[|o]|\}[|o]|o[|{]|\|[{])(--|\.\.)(\|[|o]|\}[|o]|o[|{]|\|[{])\s+(\S+)\s*:\s*(.+)$`,
	)
	erEntityOpenRe = regexp.MustCompile(`^(\S+?)(?:\["([^"]+)"\])?\s*\{$`)
)

// mapCardinality converts an ER cardinality token to an EdgeDecoration.
func mapCardinality(token string) ir.EdgeDecoration {
	switch token {
	case "||":
		return ir.DecCrowsFootOne
	case "|o", "o|":
		return ir.DecCrowsFootZeroOne
	case "}|", "|{":
		return ir.DecCrowsFootMany
	case "}o", "o{":
		return ir.DecCrowsFootZeroMany
	default:
		return ir.DecCrowsFootOne
	}
}

// parseKeyToken converts a key string like "PK", "FK", "UK", or composite
// "PK,FK" into a slice of AttributeKey values.
func parseKeyToken(token string) []ir.AttributeKey {
	if token == "" {
		return nil
	}
	parts := strings.Split(token, ",")
	var keys []ir.AttributeKey
	for _, p := range parts {
		p = strings.TrimSpace(p)
		switch p {
		case "PK":
			keys = append(keys, ir.KeyPrimary)
		case "FK":
			keys = append(keys, ir.KeyForeign)
		case "UK":
			keys = append(keys, ir.KeyUnique)
		}
	}
	return keys
}

// parseERAttribute parses a single attribute line inside an entity block.
// Format: type name [keys] ["comment"]
func parseERAttribute(line string) ir.EntityAttribute {
	attr := ir.EntityAttribute{}

	// Try to extract a quoted comment at the end.
	var comment string
	if idx := strings.Index(line, `"`); idx >= 0 {
		rest := line[idx+1:]
		if endIdx := strings.Index(rest, `"`); endIdx >= 0 {
			comment = rest[:endIdx]
			line = strings.TrimSpace(line[:idx])
		}
	}

	fields := strings.Fields(line)
	if len(fields) >= 1 {
		attr.Type = fields[0]
	}
	if len(fields) >= 2 {
		attr.Name = fields[1]
	}
	if len(fields) >= 3 {
		attr.Keys = parseKeyToken(fields[2])
	}
	attr.Comment = comment
	return attr
}

// parseEntityID parses an entity identifier, handling alias syntax like:
//
//	p["Person"] or entityName[Alias]
//
// Returns the ID and label (empty if no alias).
func parseEntityID(token string) (id string, label string) {
	if bracketIdx := strings.Index(token, "["); bracketIdx >= 0 {
		id = token[:bracketIdx]
		rest := token[bracketIdx+1:]
		if strings.HasSuffix(rest, "]") {
			label = rest[:len(rest)-1]
			label = stripQuotes(label)
		}
		return id, label
	}
	return token, ""
}

// parseER parses an ER diagram from Mermaid syntax into a ParseOutput.
func parseER(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Er

	lines := preprocessInput(input)

	inEntity := false
	var currentEntityID string
	braceDepth := 0

	for _, line := range lines {
		lower := strings.ToLower(line)

		// Skip header line.
		if strings.HasPrefix(lower, "erdiagram") {
			continue
		}

		// Handle direction.
		if dir, ok := parseDirectionLine(line); ok {
			g.Direction = dir
			continue
		}

		// Inside entity block: collect attributes.
		if inEntity {
			trimmed := strings.TrimSpace(line)
			if trimmed == "}" {
				braceDepth--
				if braceDepth == 0 {
					inEntity = false
					currentEntityID = ""
				}
				continue
			}
			// Count nested braces (unlikely in ER but be safe).
			braceDepth += strings.Count(trimmed, "{") - strings.Count(trimmed, "}")
			if braceDepth <= 0 {
				inEntity = false
				currentEntityID = ""
				continue
			}
			// Parse attribute line.
			attr := parseERAttribute(trimmed)
			if ent, ok := g.Entities[currentEntityID]; ok {
				ent.Attributes = append(ent.Attributes, attr)
			}
			continue
		}

		// Try entity block opening: ENTITY_NAME { or ENTITY_NAME["Alias"] {
		if m := erEntityOpenRe.FindStringSubmatch(line); m != nil {
			rawID := m[1]
			alias := m[2]
			entityID, entityLabel := parseEntityID(rawID)
			// If regex captured the alias in group 2, use that.
			if alias != "" {
				entityLabel = alias
			}
			if _, exists := g.Entities[entityID]; !exists {
				g.Entities[entityID] = &ir.Entity{
					ID:    entityID,
					Label: entityLabel,
				}
			} else {
				if entityLabel != "" {
					g.Entities[entityID].Label = entityLabel
				}
			}
			g.EnsureNode(entityID, nil, nil)
			inEntity = true
			currentEntityID = entityID
			braceDepth = 1
			continue
		}

		// Try relationship line.
		if m := erRelRe.FindStringSubmatch(line); m != nil {
			leftEntity := m[1]
			leftCard := m[2]
			lineStyle := m[3]
			rightCard := m[4]
			rightEntity := m[5]
			label := strings.TrimSpace(m[6])

			// Ensure nodes exist.
			g.EnsureNode(leftEntity, nil, nil)
			g.EnsureNode(rightEntity, nil, nil)

			// Ensure entities exist.
			if _, ok := g.Entities[leftEntity]; !ok {
				g.Entities[leftEntity] = &ir.Entity{ID: leftEntity}
			}
			if _, ok := g.Entities[rightEntity]; !ok {
				g.Entities[rightEntity] = &ir.Entity{ID: rightEntity}
			}

			// Map cardinality.
			startDec := mapCardinality(leftCard)
			endDec := mapCardinality(rightCard)

			// Determine line style.
			var style ir.EdgeStyle
			if lineStyle == ".." {
				style = ir.Dotted
			} else {
				style = ir.Solid
			}

			edge := &ir.Edge{
				From:            leftEntity,
				To:              rightEntity,
				Label:           &label,
				StartDecoration: &startDec,
				EndDecoration:   &endDec,
				Style:           style,
				Directed:        false,
			}
			g.Edges = append(g.Edges, edge)
			continue
		}

		// Standalone entity name (no braces, no relationship).
		// Only if it's a single word or has alias syntax.
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && (!strings.Contains(trimmed, " ") || strings.Contains(trimmed, "[")) {
			entityID, entityLabel := parseEntityID(strings.Fields(trimmed)[0])
			if _, exists := g.Entities[entityID]; !exists {
				g.Entities[entityID] = &ir.Entity{
					ID:    entityID,
					Label: entityLabel,
				}
			}
			g.EnsureNode(entityID, nil, nil)
		}
	}

	return &ParseOutput{Graph: g}, nil
}
