package render

import (
	"fmt"
	"strings"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func renderXYChart(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	xyd, ok := l.Diagram.(layout.XYChartData)
	if !ok {
		return
	}

	// Title.
	if xyd.Title != "" {
		b.text(l.Width/2, cfg.XYChart.TitleFontSize, xyd.Title,
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
		b.line(cx, tick.Y, cx+cw, tick.Y,
			"stroke", gridColor, "stroke-width", "0.5")
		// Tick label.
		b.text(cx-4, tick.Y+4, tick.Label,
			"text-anchor", "end",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(cfg.XYChart.AxisFontSize),
			"fill", th.TextColor,
		)
	}

	// Axis lines.
	b.line(cx, cy, cx, cy+ch, "stroke", axisColor, "stroke-width", "1")       // Y-axis
	b.line(cx, cy+ch, cx+cw, cy+ch, "stroke", axisColor, "stroke-width", "1") // X-axis

	// X-axis labels.
	for _, label := range xyd.XLabels {
		b.text(label.X, cy+ch+cfg.XYChart.AxisFontSize+4, label.Text,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(cfg.XYChart.AxisFontSize),
			"fill", th.TextColor,
		)
	}

	// Render each series.
	for _, s := range xyd.Series {
		color := "#4C78A8" // fallback
		if len(th.XYChartColors) > 0 {
			color = th.XYChartColors[s.ColorIndex%len(th.XYChartColors)]
		}

		switch s.Type {
		case ir.XYSeriesBar:
			for _, p := range s.Points {
				b.rect(p.X, p.Y, p.Width, p.Height, 0,
					"fill", color,
					"stroke", "none",
				)
			}
		case ir.XYSeriesLine:
			// Polyline.
			var pointStrs []string
			for _, p := range s.Points {
				pointStrs = append(pointStrs, fmt.Sprintf("%s,%s", fmtFloat(p.X), fmtFloat(p.Y)))
			}
			if len(pointStrs) > 0 {
				b.selfClose("polyline",
					"points", strings.Join(pointStrs, " "),
					"fill", "none",
					"stroke", color,
					"stroke-width", "2",
				)
			}
			// Data point circles.
			for _, p := range s.Points {
				b.circle(p.X, p.Y, 3,
					"fill", color,
					"stroke", th.Background,
					"stroke-width", "1",
				)
			}
		}
	}
}
