package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func renderRadar(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	rd, ok := l.Diagram.(layout.RadarData)
	if !ok {
		return
	}

	cx := rd.CenterX
	cy := rd.CenterY
	numAxes := len(rd.Axes)

	// Title.
	if rd.Title != "" {
		b.text(l.Width/2, 20, rd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", "16",
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	graticuleColor := th.RadarGraticuleColor
	if graticuleColor == "" {
		graticuleColor = "#E0E0E0"
	}
	axisColor := th.RadarAxisColor
	if axisColor == "" {
		axisColor = "#333"
	}

	// Graticule (concentric rings).
	for _, r := range rd.GraticuleRadii {
		if rd.GraticuleType == ir.RadarGraticulePolygon && numAxes >= 3 {
			// Polygon graticule.
			var points []string
			angleStep := 2 * math.Pi / float64(numAxes)
			for i := range numAxes {
				angle := -math.Pi/2 + float64(i)*angleStep
				px := cx + r*float32(math.Cos(angle))
				py := cy + r*float32(math.Sin(angle))
				points = append(points, fmt.Sprintf("%s,%s", fmtFloat(px), fmtFloat(py)))
			}
			b.selfClose("polygon",
				"points", strings.Join(points, " "),
				"fill", "none",
				"stroke", graticuleColor,
				"stroke-width", "0.5",
			)
		} else {
			// Circle graticule.
			b.circle(cx, cy, r,
				"fill", "none",
				"stroke", graticuleColor,
				"stroke-width", "0.5",
			)
		}
	}

	// Axis lines.
	for _, ax := range rd.Axes {
		b.line(cx, cy, ax.EndX, ax.EndY,
			"stroke", axisColor, "stroke-width", "1")
		// Axis label.
		anchor := "middle"
		if ax.LabelX > cx+5 {
			anchor = "start"
		} else if ax.LabelX < cx-5 {
			anchor = "end"
		}
		b.text(ax.LabelX, ax.LabelY, ax.Label,
			"text-anchor", anchor,
			"dominant-baseline", "middle",
			"font-family", th.FontFamily,
			"font-size", "12",
			"fill", th.TextColor,
		)
	}

	// Curve polygons.
	curveOpacity := fmt.Sprintf("%.2f", cfg.Radar.CurveOpacity)
	for _, curve := range rd.Curves {
		color := "#4C78A8" // fallback
		if len(th.RadarCurveColors) > 0 {
			color = th.RadarCurveColors[curve.ColorIndex%len(th.RadarCurveColors)]
		}

		var points []string
		for _, p := range curve.Points {
			points = append(points, fmt.Sprintf("%s,%s", fmtFloat(p[0]), fmtFloat(p[1])))
		}
		if len(points) > 0 {
			b.selfClose("polygon",
				"points", strings.Join(points, " "),
				"fill", color,
				"fill-opacity", curveOpacity,
				"stroke", color,
				"stroke-width", "2",
			)
		}
	}

	// Legend.
	if rd.ShowLegend && len(rd.Curves) > 0 {
		legendX := l.Width - cfg.Radar.PaddingX - 100
		legendY := cfg.Radar.PaddingY
		for i, curve := range rd.Curves {
			color := "#4C78A8"
			if len(th.RadarCurveColors) > 0 {
				color = th.RadarCurveColors[curve.ColorIndex%len(th.RadarCurveColors)]
			}
			y := legendY + float32(i)*20
			b.rect(legendX, y, 12, 12, 0, "fill", color)
			b.text(legendX+16, y+10, curve.Label,
				"font-family", th.FontFamily,
				"font-size", "12",
				"fill", th.TextColor,
			)
		}
	}
}
