package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	xyValuesRe   = regexp.MustCompile(`\[([^\]]+)\]`)
	xyNumAxisRe  = regexp.MustCompile(`^(?:"([^"]*)"?\s+)?(-?[\d.]+)\s*-->\s*(-?[\d.]+)$`)
	xyBandAxisRe = regexp.MustCompile(`^(?:"([^"]*)"?\s+)?\[([^\]]+)\]$`)
)

func parseXYChart(input string) (*ParseOutput, error) { //nolint:unparam // error return is part of the parser interface contract used by Parse().
	graph := ir.NewGraph()
	graph.Kind = ir.XYChart

	lines := preprocessInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: graph}, nil
	}

	// Check for horizontal orientation on the first line.
	first := strings.ToLower(lines[0])
	if strings.Contains(first, "horizontal") {
		graph.XYHorizontal = true
	}

	for _, line := range lines[1:] {
		lower := strings.ToLower(strings.TrimSpace(line))
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(lower, "title"):
			graph.XYTitle = extractQuotedText(trimmed[5:])

		case strings.HasPrefix(lower, "x-axis"):
			graph.XYXAxis = parseXYAxis(strings.TrimSpace(trimmed[6:]))

		case strings.HasPrefix(lower, "y-axis"):
			graph.XYYAxis = parseXYAxis(strings.TrimSpace(trimmed[6:]))

		case strings.HasPrefix(lower, "bar"):
			if vals := parseXYValues(trimmed); vals != nil {
				graph.XYSeries = append(graph.XYSeries, &ir.XYSeries{
					Type:   ir.XYSeriesBar,
					Values: vals,
				})
			}

		case strings.HasPrefix(lower, "line"):
			if vals := parseXYValues(trimmed); vals != nil {
				graph.XYSeries = append(graph.XYSeries, &ir.XYSeries{
					Type:   ir.XYSeriesLine,
					Values: vals,
				})
			}
		}
	}

	return &ParseOutput{Graph: graph}, nil
}

func parseXYAxis(str string) *ir.XYAxis {
	// Try numeric range: "Title" min --> max  or  min --> max.
	if match := xyNumAxisRe.FindStringSubmatch(str); match != nil {
		axis := &ir.XYAxis{Mode: ir.XYAxisNumeric, Title: match[1]}
		axis.Min, _ = strconv.ParseFloat(match[2], 64) //nolint:errcheck // regex guarantees digits.
		axis.Max, _ = strconv.ParseFloat(match[3], 64) //nolint:errcheck // regex guarantees digits.
		return axis
	}
	// Try band/categorical: "Title" [a, b, c]  or  [a, b, c].
	if match := xyBandAxisRe.FindStringSubmatch(str); match != nil {
		cats := splitAndTrimCommas(match[2])
		return &ir.XYAxis{Mode: ir.XYAxisBand, Title: match[1], Categories: cats}
	}
	// Title only (auto-range).
	title := extractQuotedText(str)
	if title == "" {
		title = strings.TrimSpace(str)
	}
	return &ir.XYAxis{Mode: ir.XYAxisNumeric, Title: title}
}

func parseXYValues(line string) []float64 {
	match := xyValuesRe.FindStringSubmatch(line)
	if match == nil {
		return nil
	}
	parts := splitAndTrimCommas(match[1])
	vals := make([]float64, 0, len(parts))
	for _, part := range parts {
		val, parseErr := strconv.ParseFloat(part, 64)
		if parseErr == nil {
			vals = append(vals, val)
		}
	}
	return vals
}
