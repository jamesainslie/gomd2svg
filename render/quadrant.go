package render

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderQuadrant(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	qd, ok := l.Diagram.(layout.QuadrantData)
	if !ok {
		return
	}

	cx := qd.ChartX
	cy := qd.ChartY
	w := qd.ChartWidth
	h := qd.ChartHeight
	halfW := w / 2
	halfH := h / 2

	qLabelSize := cfg.Quadrant.QuadrantLabelFontSize
	axisLabelSize := cfg.Quadrant.AxisLabelFontSize

	fills := [4]string{th.QuadrantFill1, th.QuadrantFill2, th.QuadrantFill3, th.QuadrantFill4}

	// Title.
	if qd.Title != "" {
		b.text(l.Width/2, th.FontSize+cfg.Quadrant.PaddingY/2, qd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize),
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	// Quadrant background rects: Q1=top-right, Q2=top-left, Q3=bottom-left, Q4=bottom-right.
	b.rect(cx+halfW, cy, halfW, halfH, 0, "fill", fills[0])       // Q1 top-right
	b.rect(cx, cy, halfW, halfH, 0, "fill", fills[1])             // Q2 top-left
	b.rect(cx, cy+halfH, halfW, halfH, 0, "fill", fills[2])       // Q3 bottom-left
	b.rect(cx+halfW, cy+halfH, halfW, halfH, 0, "fill", fills[3]) // Q4 bottom-right

	// Quadrant border.
	b.rect(cx, cy, w, h, 0,
		"fill", "none",
		"stroke", th.LineColor,
		"stroke-width", "1",
	)

	// Center cross lines.
	b.line(cx+halfW, cy, cx+halfW, cy+h,
		"stroke", th.LineColor, "stroke-width", "0.5", "stroke-dasharray", "4,4")
	b.line(cx, cy+halfH, cx+w, cy+halfH,
		"stroke", th.LineColor, "stroke-width", "0.5", "stroke-dasharray", "4,4")

	// Quadrant labels centered in each quadrant.
	if qd.Labels[0] != "" {
		b.text(cx+halfW+halfW/2, cy+halfH/2, qd.Labels[0],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}
	if qd.Labels[1] != "" {
		b.text(cx+halfW/2, cy+halfH/2, qd.Labels[1],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}
	if qd.Labels[2] != "" {
		b.text(cx+halfW/2, cy+halfH+halfH/2, qd.Labels[2],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}
	if qd.Labels[3] != "" {
		b.text(cx+halfW+halfW/2, cy+halfH+halfH/2, qd.Labels[3],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}

	// Axis labels.
	if qd.XAxisLeft != "" {
		b.text(cx, cy+h+axisLabelSize+4, qd.XAxisLeft,
			"text-anchor", "start", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor)
	}
	if qd.XAxisRight != "" {
		b.text(cx+w, cy+h+axisLabelSize+4, qd.XAxisRight,
			"text-anchor", "end", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor)
	}
	if qd.YAxisBottom != "" {
		b.text(cx-4, cy+h, qd.YAxisBottom,
			"text-anchor", "end", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor,
			"transform", "rotate(-90,"+fmtFloat(cx-4)+","+fmtFloat(cy+h)+")")
	}
	if qd.YAxisTop != "" {
		b.text(cx-4, cy, qd.YAxisTop,
			"text-anchor", "start", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor,
			"transform", "rotate(-90,"+fmtFloat(cx-4)+","+fmtFloat(cy)+")")
	}

	// Data points.
	pointR := cfg.Quadrant.PointRadius
	for _, p := range qd.Points {
		b.circle(p.X, p.Y, pointR,
			"fill", th.QuadrantPointFill,
			"stroke", th.LineColor,
			"stroke-width", "1",
		)
		b.text(p.X+pointR+3, p.Y+4, p.Label,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize),
			"fill", th.TextColor,
		)
	}
}
