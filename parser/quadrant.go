package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var quadrantPointRe = regexp.MustCompile(`^\s*(.+?):\s*\[([0-9.]+),\s*([0-9.]+)\]\s*$`)

func parseQuadrant(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Quadrant

	lines := preprocessInput(input)

	for _, line := range lines {
		lower := strings.ToLower(line)

		if strings.HasPrefix(lower, "quadrantchart") {
			continue
		}

		// Title.
		if strings.HasPrefix(lower, "title ") {
			g.QuadrantTitle = strings.TrimSpace(line[len("title "):])
			continue
		}

		// X-axis: "x-axis Left --> Right" or "x-axis Label"
		if strings.HasPrefix(lower, "x-axis ") {
			rest := strings.TrimSpace(line[len("x-axis "):])
			if parts := strings.SplitN(rest, "-->", 2); len(parts) == 2 {
				g.XAxisLeft = strings.TrimSpace(parts[0])
				g.XAxisRight = strings.TrimSpace(parts[1])
			} else {
				g.XAxisLeft = rest
			}
			continue
		}

		// Y-axis: "y-axis Bottom --> Top" or "y-axis Label"
		if strings.HasPrefix(lower, "y-axis ") {
			rest := strings.TrimSpace(line[len("y-axis "):])
			if parts := strings.SplitN(rest, "-->", 2); len(parts) == 2 {
				g.YAxisBottom = strings.TrimSpace(parts[0])
				g.YAxisTop = strings.TrimSpace(parts[1])
			} else {
				g.YAxisBottom = rest
			}
			continue
		}

		// Quadrant labels.
		if strings.HasPrefix(lower, "quadrant-1 ") {
			g.QuadrantLabels[0] = strings.TrimSpace(line[len("quadrant-1 "):])
			continue
		}
		if strings.HasPrefix(lower, "quadrant-2 ") {
			g.QuadrantLabels[1] = strings.TrimSpace(line[len("quadrant-2 "):])
			continue
		}
		if strings.HasPrefix(lower, "quadrant-3 ") {
			g.QuadrantLabels[2] = strings.TrimSpace(line[len("quadrant-3 "):])
			continue
		}
		if strings.HasPrefix(lower, "quadrant-4 ") {
			g.QuadrantLabels[3] = strings.TrimSpace(line[len("quadrant-4 "):])
			continue
		}

		// Data point: "Label: [x, y]"
		if m := quadrantPointRe.FindStringSubmatch(line); m != nil {
			x, _ := strconv.ParseFloat(m[2], 64) // regex guarantees digits
			y, _ := strconv.ParseFloat(m[3], 64) // regex guarantees digits
			// Clamp to valid [0, 1] range.
			if x < 0 {
				x = 0
			} else if x > 1 {
				x = 1
			}
			if y < 0 {
				y = 0
			} else if y > 1 {
				y = 1
			}
			g.QuadrantPoints = append(g.QuadrantPoints, &ir.QuadrantPoint{
				Label: strings.TrimSpace(m[1]),
				X:     x,
				Y:     y,
			})
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}
