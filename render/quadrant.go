package render

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Quadrant diagram rendering constants.
const (
	quadrantAxisLabelPad float32 = 4
)

func renderQuadrant(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	qd, ok := lay.Diagram.(layout.QuadrantData)
	if !ok {
		return
	}

	cx := qd.ChartX
	cy := qd.ChartY
	width := qd.ChartWidth
	height := qd.ChartHeight
	halfW := width / 2
	halfH := height / 2

	qLabelSize := cfg.Quadrant.QuadrantLabelFontSize
	axisLabelSize := cfg.Quadrant.AxisLabelFontSize

	fills := [4]string{th.QuadrantFill1, th.QuadrantFill2, th.QuadrantFill3, th.QuadrantFill4}

	// Title.
	if qd.Title != "" {
		builder.text(lay.Width/2, th.FontSize+cfg.Quadrant.PaddingY/2, qd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize),
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	// Quadrant background rects: Q1=top-right, Q2=top-left, Q3=bottom-left, Q4=bottom-right.
	builder.rect(cx+halfW, cy, halfW, halfH, 0, "fill", fills[0])       // Q1 top-right
	builder.rect(cx, cy, halfW, halfH, 0, "fill", fills[1])             // Q2 top-left
	builder.rect(cx, cy+halfH, halfW, halfH, 0, "fill", fills[2])       // Q3 bottom-left
	builder.rect(cx+halfW, cy+halfH, halfW, halfH, 0, "fill", fills[3]) // Q4 bottom-right

	// Quadrant border.
	builder.rect(cx, cy, width, height, 0,
		"fill", "none",
		"stroke", th.LineColor,
		"stroke-width", "1",
	)

	// Center cross lines.
	builder.line(cx+halfW, cy, cx+halfW, cy+height,
		"stroke", th.LineColor, "stroke-width", "0.5", "stroke-dasharray", "4,4")
	builder.line(cx, cy+halfH, cx+width, cy+halfH,
		"stroke", th.LineColor, "stroke-width", "0.5", "stroke-dasharray", "4,4")

	// Quadrant labels centered in each quadrant.
	if qd.Labels[0] != "" {
		builder.text(cx+halfW+halfW/2, cy+halfH/2, qd.Labels[0],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}
	if qd.Labels[1] != "" {
		builder.text(cx+halfW/2, cy+halfH/2, qd.Labels[1],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}
	if qd.Labels[2] != "" {
		builder.text(cx+halfW/2, cy+halfH+halfH/2, qd.Labels[2],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}
	if qd.Labels[3] != "" {
		builder.text(cx+halfW+halfW/2, cy+halfH+halfH/2, qd.Labels[3],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}

	// Axis labels.
	if qd.XAxisLeft != "" {
		builder.text(cx, cy+height+axisLabelSize+quadrantAxisLabelPad, qd.XAxisLeft,
			"text-anchor", "start", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor)
	}
	if qd.XAxisRight != "" {
		builder.text(cx+width, cy+height+axisLabelSize+quadrantAxisLabelPad, qd.XAxisRight,
			"text-anchor", "end", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor)
	}
	if qd.YAxisBottom != "" {
		builder.text(cx-quadrantAxisLabelPad, cy+height, qd.YAxisBottom,
			"text-anchor", "end", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor,
			"transform", "rotate(-90,"+fmtFloat(cx-quadrantAxisLabelPad)+","+fmtFloat(cy+height)+")")
	}
	if qd.YAxisTop != "" {
		builder.text(cx-quadrantAxisLabelPad, cy, qd.YAxisTop,
			"text-anchor", "start", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor,
			"transform", "rotate(-90,"+fmtFloat(cx-quadrantAxisLabelPad)+","+fmtFloat(cy)+")")
	}

	// Data points.
	pointR := cfg.Quadrant.PointRadius
	for _, point := range qd.Points {
		builder.circle(point.X, point.Y, pointR,
			"fill", th.QuadrantPointFill,
			"stroke", th.LineColor,
			"stroke-width", "1",
		)
		builder.text(point.X+pointR+3, point.Y+quadrantAxisLabelPad, point.Label,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize),
			"fill", th.TextColor,
		)
	}
}
