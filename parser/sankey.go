package parser

import (
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

func parseSankey(input string) (*ParseOutput, error) { //nolint:unparam // error return is part of the parser interface contract used by Parse().
	graph := ir.NewGraph()
	graph.Kind = ir.Sankey

	lines := preprocessInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: graph}, nil
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
		value, parseErr := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
		if parseErr != nil {
			continue
		}
		graph.SankeyLinks = append(graph.SankeyLinks, &ir.SankeyLink{
			Source: fields[0],
			Target: fields[1],
			Value:  value,
		})
	}

	return &ParseOutput{Graph: graph}, nil
}

// parseSankeyCSVLine parses an RFC 4180 CSV line with quoted field support.
func parseSankeyCSVLine(line string) []string {
	var fields []string
	pos := 0
	for pos < len(line) {
		if line[pos] == '"' {
			// Quoted field.
			pos++ // skip opening quote
			var field strings.Builder
			for pos < len(line) {
				if line[pos] == '"' {
					if pos+1 < len(line) && line[pos+1] == '"' {
						field.WriteByte('"')
						pos += 2
						continue
					}
					pos++ // skip closing quote
					break
				}
				field.WriteByte(line[pos])
				pos++
			}
			fields = append(fields, field.String())
			if pos < len(line) && line[pos] == ',' {
				pos++
			}
		} else {
			end := strings.IndexByte(line[pos:], ',')
			if end < 0 {
				fields = append(fields, strings.TrimSpace(line[pos:]))
				break
			}
			fields = append(fields, strings.TrimSpace(line[pos:pos+end]))
			pos += end + 1
		}
	}
	return fields
}
