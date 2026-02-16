package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var pieDataRe = regexp.MustCompile(`^\s*"([^"]+)"\s*:\s*(\d+\.?\d*)\s*$`)

func parsePie(input string) (*ParseOutput, error) { //nolint:unparam // error return is part of the parser interface contract used by Parse().
	graph := ir.NewGraph()
	graph.Kind = ir.Pie

	lines := preprocessInput(input)

	for _, line := range lines {
		lower := strings.ToLower(line)

		// Skip the declaration line, extract showData flag and inline title.
		if strings.HasPrefix(lower, "pie") {
			if strings.Contains(lower, "showdata") {
				graph.PieShowData = true
			}
			// Extract inline title: "pie title My Title" or "pie showData title My Title".
			if idx := strings.Index(lower, "title "); idx >= 0 {
				graph.PieTitle = strings.TrimSpace(line[idx+len("title "):])
			}
			continue
		}

		// Title line.
		if strings.HasPrefix(lower, "title ") {
			graph.PieTitle = strings.TrimSpace(line[len("title "):])
			continue
		}

		// Data line: "Label" : value.
		if match := pieDataRe.FindStringSubmatch(line); match != nil {
			val, _ := strconv.ParseFloat(match[2], 64) //nolint:errcheck // regex guarantees digits.
			graph.PieSlices = append(graph.PieSlices, &ir.PieSlice{
				Label: match[1],
				Value: val,
			})
			continue
		}
	}

	return &ParseOutput{Graph: graph}, nil
}
