package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var pieDataRe = regexp.MustCompile(`^\s*"([^"]+)"\s*:\s*(\d+\.?\d*)\s*$`)

func parsePie(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Pie

	lines := preprocessInput(input)

	for _, line := range lines {
		lower := strings.ToLower(line)

		// Skip the declaration line, extract showData flag and inline title.
		if strings.HasPrefix(lower, "pie") {
			if strings.Contains(lower, "showdata") {
				g.PieShowData = true
			}
			// Extract inline title: "pie title My Title" or "pie showData title My Title"
			if idx := strings.Index(lower, "title "); idx >= 0 {
				g.PieTitle = strings.TrimSpace(line[idx+len("title "):])
			}
			continue
		}

		// Title line.
		if strings.HasPrefix(lower, "title ") {
			g.PieTitle = strings.TrimSpace(line[len("title "):])
			continue
		}

		// Data line: "Label" : value
		if m := pieDataRe.FindStringSubmatch(line); m != nil {
			val, _ := strconv.ParseFloat(m[2], 64) // regex guarantees digits
			g.PieSlices = append(g.PieSlices, &ir.PieSlice{
				Label: m[1],
				Value: val,
			})
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}
