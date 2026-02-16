package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	xyValuesRe   = regexp.MustCompile(`\[([^\]]+)\]`)
	xyNumAxisRe  = regexp.MustCompile(`^(?:"([^"]*)"?\s+)?(-?[\d.]+)\s*-->\s*(-?[\d.]+)$`)
	xyBandAxisRe = regexp.MustCompile(`^(?:"([^"]*)"?\s+)?\[([^\]]+)\]$`)
)

func parseXYChart(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart

	lines := preprocessInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	// Check for horizontal orientation on the first line.
	first := strings.ToLower(lines[0])
	if strings.Contains(first, "horizontal") {
		g.XYHorizontal = true
	}

	for _, line := range lines[1:] {
		lower := strings.ToLower(strings.TrimSpace(line))
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(lower, "title"):
			g.XYTitle = extractQuotedText(trimmed[5:])

		case strings.HasPrefix(lower, "x-axis"):
			g.XYXAxis = parseXYAxis(strings.TrimSpace(trimmed[6:]))

		case strings.HasPrefix(lower, "y-axis"):
			g.XYYAxis = parseXYAxis(strings.TrimSpace(trimmed[6:]))

		case strings.HasPrefix(lower, "bar"):
			if vals := parseXYValues(trimmed); vals != nil {
				g.XYSeries = append(g.XYSeries, &ir.XYSeries{
					Type:   ir.XYSeriesBar,
					Values: vals,
				})
			}

		case strings.HasPrefix(lower, "line"):
			if vals := parseXYValues(trimmed); vals != nil {
				g.XYSeries = append(g.XYSeries, &ir.XYSeries{
					Type:   ir.XYSeriesLine,
					Values: vals,
				})
			}
		}
	}

	return &ParseOutput{Graph: g}, nil
}

func parseXYAxis(s string) *ir.XYAxis {
	// Try numeric range: "Title" min --> max  or  min --> max
	if m := xyNumAxisRe.FindStringSubmatch(s); m != nil {
		axis := &ir.XYAxis{Mode: ir.XYAxisNumeric, Title: m[1]}
		axis.Min, _ = strconv.ParseFloat(m[2], 64) // regex guarantees digits
		axis.Max, _ = strconv.ParseFloat(m[3], 64) // regex guarantees digits
		return axis
	}
	// Try band/categorical: "Title" [a, b, c]  or  [a, b, c]
	if m := xyBandAxisRe.FindStringSubmatch(s); m != nil {
		cats := splitAndTrimCommas(m[2])
		return &ir.XYAxis{Mode: ir.XYAxisBand, Title: m[1], Categories: cats}
	}
	// Title only (auto-range).
	title := extractQuotedText(s)
	if title == "" {
		title = strings.TrimSpace(s)
	}
	return &ir.XYAxis{Mode: ir.XYAxisNumeric, Title: title}
}

func parseXYValues(line string) []float64 {
	m := xyValuesRe.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	parts := splitAndTrimCommas(m[1])
	vals := make([]float64, 0, len(parts))
	for _, p := range parts {
		v, err := strconv.ParseFloat(p, 64)
		if err == nil {
			vals = append(vals, v)
		}
	}
	return vals
}
