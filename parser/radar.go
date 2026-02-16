package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	radarAxisRe  = regexp.MustCompile(`(\w+)\["([^"]+)"\]`)
	radarCurveRe = regexp.MustCompile(`^curve\s+(\w+)(?:\["([^"]+)"\])?\s*\{([^}]+)\}`)
	radarKVRe    = regexp.MustCompile(`(\w+)\s*:\s*(-?[\d.]+)`)
)

func parseRadar(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Radar

	lines := preprocessInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	// Build axis ID -> index map for key-value curve resolution.
	var axisIDs []string

	for _, line := range lines[1:] {
		lower := strings.ToLower(strings.TrimSpace(line))
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(lower, "title"):
			g.RadarTitle = extractQuotedText(trimmed[5:])

		case strings.HasPrefix(lower, "showlegend"):
			g.RadarShowLegend = true

		case strings.HasPrefix(lower, "graticule"):
			rest := strings.TrimSpace(lower[len("graticule"):])
			if rest == "polygon" {
				g.RadarGraticuleType = ir.RadarGraticulePolygon
			} else {
				g.RadarGraticuleType = ir.RadarGraticuleCircle
			}

		case strings.HasPrefix(lower, "ticks"):
			if v, err := strconv.Atoi(strings.TrimSpace(lower[5:])); err == nil {
				g.RadarTicks = v
			}

		case strings.HasPrefix(lower, "max"):
			if v, err := strconv.ParseFloat(strings.TrimSpace(lower[3:]), 64); err == nil {
				g.RadarMax = v
			}

		case strings.HasPrefix(lower, "min"):
			if v, err := strconv.ParseFloat(strings.TrimSpace(lower[3:]), 64); err == nil {
				g.RadarMin = v
			}

		case strings.HasPrefix(lower, "axis"):
			matches := radarAxisRe.FindAllStringSubmatch(trimmed, -1)
			for _, m := range matches {
				g.RadarAxes = append(g.RadarAxes, &ir.RadarAxis{
					ID:    m[1],
					Label: m[2],
				})
				axisIDs = append(axisIDs, m[1])
			}

		case strings.HasPrefix(lower, "curve"):
			if m := radarCurveRe.FindStringSubmatch(trimmed); m != nil {
				curve := &ir.RadarCurve{ID: m[1], Label: m[2]}
				valStr := m[3]

				// Check for key-value syntax.
				if kvMatches := radarKVRe.FindAllStringSubmatch(valStr, -1); len(kvMatches) > 0 {
					kvMap := make(map[string]float64)
					for _, kv := range kvMatches {
						v, _ := strconv.ParseFloat(kv[2], 64) // regex guarantees digits
						kvMap[kv[1]] = v
					}
					// Map to axis order.
					curve.Values = make([]float64, len(axisIDs))
					for i, id := range axisIDs {
						curve.Values[i] = kvMap[id]
					}
				} else {
					// Positional values.
					parts := splitAndTrimCommas(valStr)
					for _, p := range parts {
						v, err := strconv.ParseFloat(p, 64)
						if err == nil {
							curve.Values = append(curve.Values, v)
						}
					}
				}
				g.RadarCurves = append(g.RadarCurves, curve)
			}
		}
	}

	return &ParseOutput{Graph: g}, nil
}
