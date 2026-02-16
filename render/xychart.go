package render

import (
	"fmt"
	"strings"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// XY chart rendering constants.
const (
	xyTickLabelPad float32 = 4
)

func renderXYChart(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	xyd, ok := lay.Diagram.(layout.XYChartData)
	if !ok {
		return
	}

	// Title.
	if xyd.Title != "" {
		builder.text(lay.Width/2, cfg.XYChart.TitleFontSize, xyd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(cfg.XYChart.TitleFontSize),
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	cx := xyd.ChartX
	cy := xyd.ChartY
	cw := xyd.ChartWidth
	ch := xyd.ChartHeight

	// Grid lines and Y-axis ticks.
	gridColor := th.XYChartGridColor
	if gridColor == "" {
		gridColor = "#E0E0E0"
	}
	axisColor := th.XYChartAxisColor
	if axisColor == "" {
		axisColor = "#333"
	}

	for _, tick := range xyd.YTicks {
		// Horizontal grid line.
		builder.line(cx, tick.Y, cx+cw, tick.Y,
			"stroke", gridColor, "stroke-width", "0.5")
		// Tick label.
		builder.text(cx-xyTickLabelPad, tick.Y+xyTickLabelPad, tick.Label,
			"text-anchor", "end",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(cfg.XYChart.AxisFontSize),
			"fill", th.TextColor,
		)
	}

	// Axis lines.
	builder.line(cx, cy, cx, cy+ch, "stroke", axisColor, "stroke-width", "1")       // Y-axis
	builder.line(cx, cy+ch, cx+cw, cy+ch, "stroke", axisColor, "stroke-width", "1") // X-axis

	// X-axis labels.
	for _, label := range xyd.XLabels {
		builder.text(label.X, cy+ch+cfg.XYChart.AxisFontSize+xyTickLabelPad, label.Text,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(cfg.XYChart.AxisFontSize),
			"fill", th.TextColor,
		)
	}

	// Render each series.
	for _, series := range xyd.Series {
		color := "#4C78A8" // fallback
		if len(th.XYChartColors) > 0 {
			color = th.XYChartColors[series.ColorIndex%len(th.XYChartColors)]
		}

		switch series.Type {
		case ir.XYSeriesBar:
			for _, pt := range series.Points {
				builder.rect(pt.X, pt.Y, pt.Width, pt.Height, 0,
					"fill", color,
					"stroke", "none",
				)
			}
		case ir.XYSeriesLine:
			// Polyline.
			var pointStrs []string
			for _, pt := range series.Points {
				pointStrs = append(pointStrs, fmt.Sprintf("%s,%s", fmtFloat(pt.X), fmtFloat(pt.Y)))
			}
			if len(pointStrs) > 0 {
				builder.selfClose("polyline",
					"points", strings.Join(pointStrs, " "),
					"fill", "none",
					"stroke", color,
					"stroke-width", "2",
				)
			}
			// Data point circles.
			for _, pt := range series.Points {
				builder.circle(pt.X, pt.Y, 3,
					"fill", color,
					"stroke", th.Background,
					"stroke-width", "1",
				)
			}
		}
	}
}
