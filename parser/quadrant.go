package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var quadrantPointRe = regexp.MustCompile(`^\s*(.+?):\s*\[([0-9.]+),\s*([0-9.]+)\]\s*$`)

func parseQuadrant(input string) (*ParseOutput, error) { //nolint:unparam // error return is part of the parser interface contract used by Parse().
	graph := ir.NewGraph()
	graph.Kind = ir.Quadrant

	lines := preprocessInput(input)

	for _, line := range lines {
		lower := strings.ToLower(line)

		if strings.HasPrefix(lower, "quadrantchart") {
			continue
		}

		// Title.
		if strings.HasPrefix(lower, "title ") {
			graph.QuadrantTitle = strings.TrimSpace(line[len("title "):])
			continue
		}

		// X-axis: "x-axis Left --> Right" or "x-axis Label".
		if strings.HasPrefix(lower, "x-axis ") {
			rest := strings.TrimSpace(line[len("x-axis "):])
			if parts := strings.SplitN(rest, "-->", 2); len(parts) == 2 {
				graph.XAxisLeft = strings.TrimSpace(parts[0])
				graph.XAxisRight = strings.TrimSpace(parts[1])
			} else {
				graph.XAxisLeft = rest
			}
			continue
		}

		// Y-axis: "y-axis Bottom --> Top" or "y-axis Label".
		if strings.HasPrefix(lower, "y-axis ") {
			rest := strings.TrimSpace(line[len("y-axis "):])
			if parts := strings.SplitN(rest, "-->", 2); len(parts) == 2 {
				graph.YAxisBottom = strings.TrimSpace(parts[0])
				graph.YAxisTop = strings.TrimSpace(parts[1])
			} else {
				graph.YAxisBottom = rest
			}
			continue
		}

		// Quadrant labels.
		if strings.HasPrefix(lower, "quadrant-1 ") {
			graph.QuadrantLabels[0] = strings.TrimSpace(line[len("quadrant-1 "):])
			continue
		}
		if strings.HasPrefix(lower, "quadrant-2 ") {
			graph.QuadrantLabels[1] = strings.TrimSpace(line[len("quadrant-2 "):])
			continue
		}
		if strings.HasPrefix(lower, "quadrant-3 ") {
			graph.QuadrantLabels[2] = strings.TrimSpace(line[len("quadrant-3 "):])
			continue
		}
		if strings.HasPrefix(lower, "quadrant-4 ") {
			graph.QuadrantLabels[3] = strings.TrimSpace(line[len("quadrant-4 "):])
			continue
		}

		// Data point: "Label: [x, y]".
		if match := quadrantPointRe.FindStringSubmatch(line); match != nil {
			xCoord, _ := strconv.ParseFloat(match[2], 64) //nolint:errcheck // regex guarantees digits.
			yCoord, _ := strconv.ParseFloat(match[3], 64) //nolint:errcheck // regex guarantees digits.
			// Clamp to valid [0, 1] range.
			if xCoord < 0 {
				xCoord = 0
			} else if xCoord > 1 {
				xCoord = 1
			}
			if yCoord < 0 {
				yCoord = 0
			} else if yCoord > 1 {
				yCoord = 1
			}
			graph.QuadrantPoints = append(graph.QuadrantPoints, &ir.QuadrantPoint{
				Label: strings.TrimSpace(match[1]),
				X:     xCoord,
				Y:     yCoord,
			})
			continue
		}
	}

	return &ParseOutput{Graph: graph}, nil
}
