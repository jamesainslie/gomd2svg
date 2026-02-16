package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	radarAxisRe  = regexp.MustCompile(`(\w+)\["([^"]+)"\]`)
	radarCurveRe = regexp.MustCompile(`^curve\s+(\w+)(?:\["([^"]+)"\])?\s*\{([^}]+)\}`)
	radarKVRe    = regexp.MustCompile(`(\w+)\s*:\s*(-?[\d.]+)`)
)

func parseRadar(input string) (*ParseOutput, error) { //nolint:unparam // error return is part of the parser interface contract used by Parse().
	graph := ir.NewGraph()
	graph.Kind = ir.Radar

	lines := preprocessInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: graph}, nil
	}

	// Build axis ID -> index map for key-value curve resolution.
	var axisIDs []string

	for _, line := range lines[1:] {
		lower := strings.ToLower(strings.TrimSpace(line))
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(lower, "title"):
			graph.RadarTitle = extractQuotedText(trimmed[5:])

		case strings.HasPrefix(lower, "showlegend"):
			graph.RadarShowLegend = true

		case strings.HasPrefix(lower, "graticule"):
			rest := strings.TrimSpace(lower[len("graticule"):])
			if rest == "polygon" {
				graph.RadarGraticuleType = ir.RadarGraticulePolygon
			} else {
				graph.RadarGraticuleType = ir.RadarGraticuleCircle
			}

		case strings.HasPrefix(lower, "ticks"):
			if val, parseErr := strconv.Atoi(strings.TrimSpace(lower[5:])); parseErr == nil {
				graph.RadarTicks = val
			}

		case strings.HasPrefix(lower, "max"):
			if val, parseErr := strconv.ParseFloat(strings.TrimSpace(lower[3:]), 64); parseErr == nil {
				graph.RadarMax = val
			}

		case strings.HasPrefix(lower, "min"):
			if val, parseErr := strconv.ParseFloat(strings.TrimSpace(lower[3:]), 64); parseErr == nil {
				graph.RadarMin = val
			}

		case strings.HasPrefix(lower, "axis"):
			matches := radarAxisRe.FindAllStringSubmatch(trimmed, -1)
			for _, match := range matches {
				graph.RadarAxes = append(graph.RadarAxes, &ir.RadarAxis{
					ID:    match[1],
					Label: match[2],
				})
				axisIDs = append(axisIDs, match[1])
			}

		case strings.HasPrefix(lower, "curve"):
			if match := radarCurveRe.FindStringSubmatch(trimmed); match != nil {
				curve := &ir.RadarCurve{ID: match[1], Label: match[2]}
				valStr := match[3]

				// Check for key-value syntax.
				if kvMatches := radarKVRe.FindAllStringSubmatch(valStr, -1); len(kvMatches) > 0 {
					kvMap := make(map[string]float64)
					for _, kvEntry := range kvMatches {
						val, _ := strconv.ParseFloat(kvEntry[2], 64) //nolint:errcheck // regex guarantees digits.
						kvMap[kvEntry[1]] = val
					}
					// Map to axis order.
					curve.Values = make([]float64, len(axisIDs))
					for idx, axisID := range axisIDs {
						curve.Values[idx] = kvMap[axisID]
					}
				} else {
					// Positional values.
					parts := splitAndTrimCommas(valStr)
					for _, part := range parts {
						val, parseErr := strconv.ParseFloat(part, 64)
						if parseErr == nil {
							curve.Values = append(curve.Values, val)
						}
					}
				}
				graph.RadarCurves = append(graph.RadarCurves, curve)
			}
		}
	}

	return &ParseOutput{Graph: graph}, nil
}
