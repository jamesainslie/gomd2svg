package render

import (
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
)

// renderNodeShape renders the SVG shape for a node and its centered text label.
// The node's X, Y are the center of the node.
func renderNodeShape(b *svgBuilder, n *layout.NodeLayout, fill, stroke, textColor string) {
	strokeWidth := "1"
	if n.Style.StrokeWidth != nil {
		strokeWidth = fmtFloat(*n.Style.StrokeWidth)
	}

	dash := ""
	if n.Style.StrokeDasharray != nil {
		dash = *n.Style.StrokeDasharray
	}

	// Compute top-left corner from center coordinates.
	x := n.X - n.Width/2
	y := n.Y - n.Height/2
	w := n.Width
	h := n.Height

	join := "round"

	baseAttrs := []string{
		"fill", fill,
		"stroke", stroke,
		"stroke-width", strokeWidth,
		"stroke-linejoin", join,
		"stroke-linecap", join,
	}
	if dash != "" {
		baseAttrs = append(baseAttrs, "stroke-dasharray", dash)
	}

	switch n.Shape {
	case ir.Rectangle, ir.ForkJoin, ir.ActorBox:
		renderRectangle(b, x, y, w, h, 3, baseAttrs)

	case ir.RoundRect:
		renderRectangle(b, x, y, w, h, 10, baseAttrs)

	case ir.Stadium:
		renderRectangle(b, x, y, w, h, h/2, baseAttrs)

	case ir.Diamond:
		renderDiamond(b, x, y, w, h, baseAttrs)

	case ir.Hexagon:
		renderHexagon(b, x, y, w, h, baseAttrs)

	case ir.Circle, ir.DoubleCircle:
		renderCircle(b, n, x, y, w, h, baseAttrs, fill, stroke)

	case ir.Cylinder:
		renderCylinder(b, x, y, w, h, fill, stroke, strokeWidth, dash)

	case ir.Subroutine:
		renderSubroutine(b, x, y, w, h, stroke, strokeWidth, baseAttrs)

	case ir.Asymmetric:
		renderAsymmetric(b, x, y, w, h, baseAttrs)

	case ir.Parallelogram:
		renderParallelogram(b, x, y, w, h, false, baseAttrs)

	case ir.ParallelogramAlt:
		renderParallelogram(b, x, y, w, h, true, baseAttrs)

	case ir.Trapezoid:
		renderTrapezoid(b, x, y, w, h, false, baseAttrs)

	case ir.TrapezoidAlt:
		renderTrapezoid(b, x, y, w, h, true, baseAttrs)

	default:
		// Fallback to rectangle.
		renderRectangle(b, x, y, w, h, 6, baseAttrs)
	}

	// Render text label centered in the node.
	renderNodeLabel(b, n, textColor)
}

// renderRectangle renders a <rect> with given corner radius.
func renderRectangle(b *svgBuilder, x, y, w, h, rx float32, attrs []string) {
	all := []string{
		"x", fmtFloat(x),
		"y", fmtFloat(y),
		"width", fmtFloat(w),
		"height", fmtFloat(h),
		"rx", fmtFloat(rx),
		"ry", fmtFloat(rx),
	}
	all = append(all, attrs...)
	b.selfClose("rect", all...)
}

// renderDiamond renders a <polygon> rotated square for Diamond shape.
func renderDiamond(b *svgBuilder, x, y, w, h float32, attrs []string) {
	cx := x + w/2
	cy := y + h/2
	pts := [][2]float32{
		{cx, y},
		{x + w, cy},
		{cx, y + h},
		{x, cy},
	}
	all := []string{"points", formatPoints(pts)}
	all = append(all, attrs...)
	b.selfClose("polygon", all...)
}

// renderHexagon renders a <polygon> with 6 vertices.
func renderHexagon(b *svgBuilder, x, y, w, h float32, attrs []string) {
	x1 := x + w*0.25
	x2 := x + w*0.75
	yMid := y + h/2
	pts := [][2]float32{
		{x1, y},
		{x2, y},
		{x + w, yMid},
		{x2, y + h},
		{x1, y + h},
		{x, yMid},
	}
	all := []string{"points", formatPoints(pts)}
	all = append(all, attrs...)
	b.selfClose("polygon", all...)
}

// renderCircle renders a <circle> and optionally an inner circle for DoubleCircle.
func renderCircle(b *svgBuilder, n *layout.NodeLayout, x, y, w, h float32, attrs []string, fill, stroke string) {
	cx := x + w/2
	cy := y + h/2
	r := min(w, h) / 2

	all := []string{
		"cx", fmtFloat(cx),
		"cy", fmtFloat(cy),
		"r", fmtFloat(r),
	}
	all = append(all, attrs...)
	b.selfClose("circle", all...)

	if n.Shape == ir.DoubleCircle {
		r2 := r - 4
		if r2 > 0 {
			b.selfClose("circle",
				"cx", fmtFloat(cx),
				"cy", fmtFloat(cy),
				"r", fmtFloat(r2),
				"fill", "none",
				"stroke", stroke,
				"stroke-width", "1",
				"stroke-linejoin", "round",
				"stroke-linecap", "round",
			)
		}
	}
}

// renderCylinder renders a cylinder shape using ellipses and a rect.
func renderCylinder(b *svgBuilder, x, y, w, h float32, fill, stroke, strokeWidth, dash string) {
	cx := x + w/2
	ry := clamp(h*0.12, 6, 14)
	rx := w / 2

	joinAttrs := []string{
		"fill", fill,
		"stroke", stroke,
		"stroke-width", strokeWidth,
		"stroke-linejoin", "round",
		"stroke-linecap", "round",
	}
	if dash != "" {
		joinAttrs = append(joinAttrs, "stroke-dasharray", dash)
	}

	// Top ellipse (filled).
	topAttrs := []string{
		"cx", fmtFloat(cx),
		"cy", fmtFloat(y + ry),
		"rx", fmtFloat(rx),
		"ry", fmtFloat(ry),
	}
	topAttrs = append(topAttrs, joinAttrs...)
	b.selfClose("ellipse", topAttrs...)

	// Body rect.
	bodyH := h - 2*ry
	if bodyH < 0 {
		bodyH = 0
	}
	bodyAttrs := []string{
		"x", fmtFloat(x),
		"y", fmtFloat(y + ry),
		"width", fmtFloat(w),
		"height", fmtFloat(bodyH),
	}
	bodyAttrs = append(bodyAttrs, joinAttrs...)
	b.selfClose("rect", bodyAttrs...)

	// Bottom ellipse (stroke only).
	bottomAttrs := []string{
		"cx", fmtFloat(cx),
		"cy", fmtFloat(y + h - ry),
		"rx", fmtFloat(rx),
		"ry", fmtFloat(ry),
		"fill", "none",
		"stroke", stroke,
		"stroke-width", strokeWidth,
		"stroke-linejoin", "round",
		"stroke-linecap", "round",
	}
	if dash != "" {
		bottomAttrs = append(bottomAttrs, "stroke-dasharray", dash)
	}
	b.selfClose("ellipse", bottomAttrs...)
}

// renderSubroutine renders a rect with double vertical lines at the sides.
func renderSubroutine(b *svgBuilder, x, y, w, h float32, stroke, strokeWidth string, attrs []string) {
	// Main rect.
	all := []string{
		"x", fmtFloat(x),
		"y", fmtFloat(y),
		"width", fmtFloat(w),
		"height", fmtFloat(h),
		"rx", "6",
		"ry", "6",
	}
	all = append(all, attrs...)
	b.selfClose("rect", all...)

	// Inner vertical lines.
	inset := float32(6)
	y1 := y + 2
	y2 := y + h - 2
	x1 := x + inset
	x2 := x + w - inset

	lineAttrs := []string{
		"stroke", stroke,
		"stroke-width", strokeWidth,
		"stroke-linejoin", "round",
		"stroke-linecap", "round",
	}

	b.line(x1, y1, x1, y2, lineAttrs...)
	b.line(x2, y1, x2, y2, lineAttrs...)
}

// renderAsymmetric renders a flag-shaped polygon.
func renderAsymmetric(b *svgBuilder, x, y, w, h float32, attrs []string) {
	slant := w * 0.22
	pts := [][2]float32{
		{x, y},
		{x + w - slant, y},
		{x + w, y + h/2},
		{x + w - slant, y + h},
		{x, y + h},
	}
	all := []string{"points", formatPoints(pts)}
	all = append(all, attrs...)
	b.selfClose("polygon", all...)
}

// renderParallelogram renders a parallelogram polygon.
func renderParallelogram(b *svgBuilder, x, y, w, h float32, alt bool, attrs []string) {
	offset := w * 0.18
	var pts [][2]float32
	if !alt {
		pts = [][2]float32{
			{x + offset, y},
			{x + w, y},
			{x + w - offset, y + h},
			{x, y + h},
		}
	} else {
		pts = [][2]float32{
			{x, y},
			{x + w - offset, y},
			{x + w, y + h},
			{x + offset, y + h},
		}
	}
	all := []string{"points", formatPoints(pts)}
	all = append(all, attrs...)
	b.selfClose("polygon", all...)
}

// renderTrapezoid renders a trapezoid polygon.
func renderTrapezoid(b *svgBuilder, x, y, w, h float32, alt bool, attrs []string) {
	offset := w * 0.18
	var pts [][2]float32
	if !alt {
		pts = [][2]float32{
			{x + offset, y},
			{x + w - offset, y},
			{x + w, y + h},
			{x, y + h},
		}
	} else {
		pts = [][2]float32{
			{x, y},
			{x + w, y},
			{x + w - offset, y + h},
			{x + offset, y + h},
		}
	}
	all := []string{"points", formatPoints(pts)}
	all = append(all, attrs...)
	b.selfClose("polygon", all...)
}

// renderNodeLabel renders text lines centered within a node.
func renderNodeLabel(b *svgBuilder, n *layout.NodeLayout, textColor string) {
	if len(n.Label.Lines) == 0 {
		return
	}

	fontSize := n.Label.FontSize
	if fontSize <= 0 {
		fontSize = 14
	}

	lineHeight := fontSize * 1.2
	totalTextHeight := lineHeight * float32(len(n.Label.Lines))
	// Start Y so that text block is vertically centered in node.
	startY := n.Y - totalTextHeight/2 + lineHeight*0.75

	for i, line := range n.Label.Lines {
		ly := startY + float32(i)*lineHeight
		b.text(n.X, ly, line,
			"text-anchor", "middle",
			"dominant-baseline", "auto",
			"fill", textColor,
			"font-size", fmtFloat(fontSize),
		)
	}
}

// clamp restricts v to the range [lo, hi].
func clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
