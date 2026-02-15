package parser

import (
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

func parseSankey(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Sankey

	lines := preprocessInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	for _, line := range lines[1:] { // skip "sankey-beta" keyword
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		fields := parseSankeyCSVLine(trimmed)
		if len(fields) < 3 {
			continue
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
		if err != nil {
			continue
		}
		g.SankeyLinks = append(g.SankeyLinks, &ir.SankeyLink{
			Source: fields[0],
			Target: fields[1],
			Value:  value,
		})
	}

	return &ParseOutput{Graph: g}, nil
}

// parseSankeyCSVLine parses an RFC 4180 CSV line with quoted field support.
func parseSankeyCSVLine(line string) []string {
	var fields []string
	i := 0
	for i < len(line) {
		if line[i] == '"' {
			// Quoted field.
			i++ // skip opening quote
			var field strings.Builder
			for i < len(line) {
				if line[i] == '"' {
					if i+1 < len(line) && line[i+1] == '"' {
						field.WriteByte('"')
						i += 2
					} else {
						i++ // skip closing quote
						break
					}
				} else {
					field.WriteByte(line[i])
					i++
				}
			}
			fields = append(fields, field.String())
			if i < len(line) && line[i] == ',' {
				i++
			}
		} else {
			end := strings.IndexByte(line[i:], ',')
			if end < 0 {
				fields = append(fields, strings.TrimSpace(line[i:]))
				break
			}
			fields = append(fields, strings.TrimSpace(line[i:i+end]))
			i += end + 1
		}
	}
	return fields
}
