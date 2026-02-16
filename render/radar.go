package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Radar diagram rendering constants.
const (
	radarTitleOffsetY      float32 = 20
	radarLegendWidth       float32 = 100
	radarLegendRowHeight   float32 = 20
	radarLegendSwatchW     float32 = 12
	radarLegendSwatchH     float32 = 12
	radarLegendTextOff     float32 = 16
	radarLegendBaselineOff float32 = 10
	radarFallbackColor             = "#4C78A8"
)

func renderRadar(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	rd, ok := lay.Diagram.(layout.RadarData)
	if !ok {
		return
	}

	cx := rd.CenterX
	cy := rd.CenterY
	numAxes := len(rd.Axes)

	// Title.
	if rd.Title != "" {
		builder.text(lay.Width/2, radarTitleOffsetY, rd.Title,
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
	for _, radius := range rd.GraticuleRadii {
		if rd.GraticuleType == ir.RadarGraticulePolygon && numAxes >= 3 {
			// Polygon graticule.
			var points []string
			angleStep := 2 * math.Pi / float64(numAxes)
			for idx := range numAxes {
				angle := -math.Pi/2 + float64(idx)*angleStep
				px := cx + radius*float32(math.Cos(angle))
				py := cy + radius*float32(math.Sin(angle))
				points = append(points, fmt.Sprintf("%s,%s", fmtFloat(px), fmtFloat(py)))
			}
			builder.selfClose("polygon",
				"points", strings.Join(points, " "),
				"fill", "none",
				"stroke", graticuleColor,
				"stroke-width", "0.5",
			)
		} else {
			// Circle graticule.
			builder.circle(cx, cy, radius,
				"fill", "none",
				"stroke", graticuleColor,
				"stroke-width", "0.5",
			)
		}
	}

	// Axis lines.
	for _, ax := range rd.Axes {
		builder.line(cx, cy, ax.EndX, ax.EndY,
			"stroke", axisColor, "stroke-width", "1")
		// Axis label.
		anchor := "middle"
		if ax.LabelX > cx+5 {
			anchor = "start"
		} else if ax.LabelX < cx-5 {
			anchor = "end"
		}
		builder.text(ax.LabelX, ax.LabelY, ax.Label,
			"text-anchor", anchor,
			"dominant-baseline", "middle",
			"font-family", th.FontFamily,
			"font-size", "12",
			"fill", th.TextColor,
		)
	}

	// Curve polygons.
	opacity := th.RadarCurveOpacity
	if opacity <= 0 {
		opacity = cfg.Radar.CurveOpacity
	}
	curveOpacity := fmt.Sprintf("%.2f", opacity)
	for _, curve := range rd.Curves {
		color := radarFallbackColor
		if len(th.RadarCurveColors) > 0 {
			color = th.RadarCurveColors[curve.ColorIndex%len(th.RadarCurveColors)]
		}

		var points []string
		for _, pt := range curve.Points {
			points = append(points, fmt.Sprintf("%s,%s", fmtFloat(pt[0]), fmtFloat(pt[1])))
		}
		if len(points) > 0 {
			builder.selfClose("polygon",
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
		legendX := lay.Width - cfg.Radar.PaddingX - radarLegendWidth
		legendY := cfg.Radar.PaddingY
		for idx, curve := range rd.Curves {
			color := radarFallbackColor
			if len(th.RadarCurveColors) > 0 {
				color = th.RadarCurveColors[curve.ColorIndex%len(th.RadarCurveColors)]
			}
			posY := legendY + float32(idx)*radarLegendRowHeight
			builder.rect(legendX, posY, radarLegendSwatchW, radarLegendSwatchH, 0, "fill", color)
			builder.text(legendX+radarLegendTextOff, posY+radarLegendBaselineOff, curve.Label,
				"font-family", th.FontFamily,
				"font-size", "12",
				"fill", th.TextColor,
			)
		}
	}
}
