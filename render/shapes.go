package render

import (
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
)

// Shape rendering constants.
const (
	roundRectRadius     float32 = 10
	defaultRectRadius   float32 = 3
	fallbackRectRadius  float32 = 6
	doubleCircleInset   float32 = 4
	cylinderMinRY       float32 = 6
	cylinderMaxRY       float32 = 14
	cylinderEllipseRate float32 = 0.12
	subroutineInset     float32 = 6
	asymmetricSlant     float32 = 0.22
	parallelogramSkew   float32 = 0.18
	lineHeightScale     float32 = 1.2
	textBaselineOffset  float32 = 0.75
	hexagonLeft         float32 = 0.25
	hexagonRight        float32 = 0.75
)

// renderNodeShape renders the SVG shape for a node and its centered text label.
// The node's X, Y are the center of the node.
func renderNodeShape(builder *svgBuilder, node *layout.NodeLayout, fill, stroke, textColor string) {
	strokeWidth := "1"
	if node.Style.StrokeWidth != nil {
		strokeWidth = fmtFloat(*node.Style.StrokeWidth)
	}

	dash := ""
	if node.Style.StrokeDasharray != nil {
		dash = *node.Style.StrokeDasharray
	}

	// Compute top-left corner from center coordinates.
	posX := node.X - node.Width/2
	posY := node.Y - node.Height/2
	width := node.Width
	height := node.Height

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

	switch node.Shape {
	case ir.Rectangle, ir.ForkJoin, ir.ActorBox:
		renderRectangle(builder, posX, posY, width, height, defaultRectRadius, baseAttrs)

	case ir.RoundRect:
		renderRectangle(builder, posX, posY, width, height, roundRectRadius, baseAttrs)

	case ir.Stadium:
		renderRectangle(builder, posX, posY, width, height, height/2, baseAttrs)

	case ir.Diamond:
		renderDiamond(builder, posX, posY, width, height, baseAttrs)

	case ir.Hexagon:
		renderHexagon(builder, posX, posY, width, height, baseAttrs)

	case ir.Circle, ir.DoubleCircle:
		renderCircleFromParams(renderCircleParams{
			builder: builder, node: node,
			posX: posX, posY: posY, width: width, height: height,
			attrs: baseAttrs, fill: fill, stroke: stroke,
		})

	case ir.Cylinder:
		renderCylinderFromParams(renderCylinderParams{
			builder: builder, posX: posX, posY: posY, width: width, height: height,
			fill: fill, stroke: stroke, strokeWidth: strokeWidth, dash: dash,
		})

	case ir.Subroutine:
		renderSubroutineFromParams(renderSubroutineParams{
			builder: builder, posX: posX, posY: posY, width: width, height: height,
			stroke: stroke, strokeWidth: strokeWidth, attrs: baseAttrs,
		})

	case ir.Asymmetric:
		renderAsymmetric(builder, posX, posY, width, height, baseAttrs)

	case ir.Parallelogram:
		renderParallelogram(builder, posX, posY, width, height, false, baseAttrs)

	case ir.ParallelogramAlt:
		renderParallelogram(builder, posX, posY, width, height, true, baseAttrs)

	case ir.Trapezoid:
		renderTrapezoid(builder, posX, posY, width, height, false, baseAttrs)

	case ir.TrapezoidAlt:
		renderTrapezoid(builder, posX, posY, width, height, true, baseAttrs)

	default:
		// Fallback to rectangle.
		renderRectangle(builder, posX, posY, width, height, fallbackRectRadius, baseAttrs)
	}

	// Render text label centered in the node.
	renderNodeLabel(builder, node, textColor)
}

// renderRectangle renders a <rect> with given corner radius.
func renderRectangle(builder *svgBuilder, posX, posY, width, height, rx float32, attrs []string) {
	all := make([]string, 0, 12+len(attrs)) //nolint:mnd // base capacity for rect attrs
	all = append(all,
		"x", fmtFloat(posX),
		"y", fmtFloat(posY),
		"width", fmtFloat(width),
		"height", fmtFloat(height),
		"rx", fmtFloat(rx),
		"ry", fmtFloat(rx),
	)
	all = append(all, attrs...)
	builder.selfClose("rect", all...)
}

// renderDiamond renders a <polygon> rotated square for Diamond shape.
func renderDiamond(builder *svgBuilder, posX, posY, width, height float32, attrs []string) {
	cx := posX + width/2
	cy := posY + height/2
	pts := [][2]float32{
		{cx, posY},
		{posX + width, cy},
		{cx, posY + height},
		{posX, cy},
	}
	all := make([]string, 0, 2+len(attrs))
	all = append(all, "points", formatPoints(pts))
	all = append(all, attrs...)
	builder.selfClose("polygon", all...)
}

// renderHexagon renders a <polygon> with 6 vertices.
func renderHexagon(builder *svgBuilder, posX, posY, width, height float32, attrs []string) {
	leftX := posX + width*hexagonLeft
	rightX := posX + width*hexagonRight
	yMid := posY + height/2
	pts := [][2]float32{
		{leftX, posY},
		{rightX, posY},
		{posX + width, yMid},
		{rightX, posY + height},
		{leftX, posY + height},
		{posX, yMid},
	}
	all := make([]string, 0, 2+len(attrs))
	all = append(all, "points", formatPoints(pts))
	all = append(all, attrs...)
	builder.selfClose("polygon", all...)
}

// renderCircleParams holds the parameters for rendering a circle shape to reduce argument count.
type renderCircleParams struct {
	builder *svgBuilder
	node    *layout.NodeLayout
	posX    float32
	posY    float32
	width   float32
	height  float32
	attrs   []string
	fill    string
	stroke  string
}

// renderCircleFromParams renders a <circle> and optionally an inner circle for DoubleCircle.
func renderCircleFromParams(params renderCircleParams) {
	cx := params.posX + params.width/2
	cy := params.posY + params.height/2
	radius := min(params.width, params.height) / 2

	all := make([]string, 0, 6+len(params.attrs)) //nolint:mnd // base capacity for cx,cy,r pairs
	all = append(all,
		"cx", fmtFloat(cx),
		"cy", fmtFloat(cy),
		"r", fmtFloat(radius),
	)
	all = append(all, params.attrs...)
	params.builder.selfClose("circle", all...)

	if params.node.Shape == ir.DoubleCircle {
		innerRadius := radius - doubleCircleInset
		if innerRadius > 0 {
			params.builder.selfClose("circle",
				"cx", fmtFloat(cx),
				"cy", fmtFloat(cy),
				"r", fmtFloat(innerRadius),
				"fill", "none",
				"stroke", params.stroke,
				"stroke-width", "1",
				"stroke-linejoin", "round",
				"stroke-linecap", "round",
			)
		}
	}
}

// renderCylinderParams holds parameters for cylinder rendering to reduce argument count.
type renderCylinderParams struct {
	builder     *svgBuilder
	posX        float32
	posY        float32
	width       float32
	height      float32
	fill        string
	stroke      string
	strokeWidth string
	dash        string
}

// renderCylinderFromParams renders a cylinder shape using ellipses and a rect.
func renderCylinderFromParams(params renderCylinderParams) {
	cx := params.posX + params.width/2
	ry := clamp(params.height*cylinderEllipseRate, cylinderMinRY, cylinderMaxRY)
	rx := params.width / 2

	joinAttrs := []string{
		"fill", params.fill,
		"stroke", params.stroke,
		"stroke-width", params.strokeWidth,
		"stroke-linejoin", "round",
		"stroke-linecap", "round",
	}
	if params.dash != "" {
		joinAttrs = append(joinAttrs, "stroke-dasharray", params.dash)
	}

	// Top ellipse (filled).
	topAttrs := make([]string, 0, 8+len(joinAttrs)) //nolint:mnd // base capacity for ellipse attrs
	topAttrs = append(topAttrs,
		"cx", fmtFloat(cx),
		"cy", fmtFloat(params.posY+ry),
		"rx", fmtFloat(rx),
		"ry", fmtFloat(ry),
	)
	topAttrs = append(topAttrs, joinAttrs...)
	params.builder.selfClose("ellipse", topAttrs...)

	// Body rect.
	bodyH := params.height - 2*ry
	if bodyH < 0 {
		bodyH = 0
	}
	bodyAttrs := make([]string, 0, 8+len(joinAttrs)) //nolint:mnd // base capacity for rect attrs
	bodyAttrs = append(bodyAttrs,
		"x", fmtFloat(params.posX),
		"y", fmtFloat(params.posY+ry),
		"width", fmtFloat(params.width),
		"height", fmtFloat(bodyH),
	)
	bodyAttrs = append(bodyAttrs, joinAttrs...)
	params.builder.selfClose("rect", bodyAttrs...)

	// Bottom ellipse (stroke only).
	bottomAttrs := []string{
		"cx", fmtFloat(cx),
		"cy", fmtFloat(params.posY + params.height - ry),
		"rx", fmtFloat(rx),
		"ry", fmtFloat(ry),
		"fill", "none",
		"stroke", params.stroke,
		"stroke-width", params.strokeWidth,
		"stroke-linejoin", "round",
		"stroke-linecap", "round",
	}
	if params.dash != "" {
		bottomAttrs = append(bottomAttrs, "stroke-dasharray", params.dash)
	}
	params.builder.selfClose("ellipse", bottomAttrs...)
}

// renderSubroutineParams holds parameters for subroutine rendering to reduce argument count.
type renderSubroutineParams struct {
	builder     *svgBuilder
	posX        float32
	posY        float32
	width       float32
	height      float32
	stroke      string
	strokeWidth string
	attrs       []string
}

// renderSubroutineFromParams renders a rect with double vertical lines at the sides.
func renderSubroutineFromParams(params renderSubroutineParams) {
	// Main rect.
	all := make([]string, 0, 12+len(params.attrs)) //nolint:mnd // base capacity for rect attrs
	all = append(all,
		"x", fmtFloat(params.posX),
		"y", fmtFloat(params.posY),
		"width", fmtFloat(params.width),
		"height", fmtFloat(params.height),
		"rx", "6",
		"ry", "6",
	)
	all = append(all, params.attrs...)
	params.builder.selfClose("rect", all...)

	// Inner vertical lines.
	lineTopY := params.posY + 2
	lineBottomY := params.posY + params.height - 2
	lineLeftX := params.posX + subroutineInset
	lineRightX := params.posX + params.width - subroutineInset

	lineAttrs := []string{
		"stroke", params.stroke,
		"stroke-width", params.strokeWidth,
		"stroke-linejoin", "round",
		"stroke-linecap", "round",
	}

	params.builder.line(lineLeftX, lineTopY, lineLeftX, lineBottomY, lineAttrs...)
	params.builder.line(lineRightX, lineTopY, lineRightX, lineBottomY, lineAttrs...)
}

// renderAsymmetric renders a flag-shaped polygon.
func renderAsymmetric(builder *svgBuilder, posX, posY, width, height float32, attrs []string) {
	slant := width * asymmetricSlant
	pts := [][2]float32{
		{posX, posY},
		{posX + width - slant, posY},
		{posX + width, posY + height/2},
		{posX + width - slant, posY + height},
		{posX, posY + height},
	}
	all := make([]string, 0, 2+len(attrs))
	all = append(all, "points", formatPoints(pts))
	all = append(all, attrs...)
	builder.selfClose("polygon", all...)
}

// renderParallelogram renders a parallelogram polygon.
func renderParallelogram(builder *svgBuilder, posX, posY, width, height float32, alt bool, attrs []string) {
	offset := width * parallelogramSkew
	var pts [][2]float32
	if !alt {
		pts = [][2]float32{
			{posX + offset, posY},
			{posX + width, posY},
			{posX + width - offset, posY + height},
			{posX, posY + height},
		}
	} else {
		pts = [][2]float32{
			{posX, posY},
			{posX + width - offset, posY},
			{posX + width, posY + height},
			{posX + offset, posY + height},
		}
	}
	all := make([]string, 0, 2+len(attrs))
	all = append(all, "points", formatPoints(pts))
	all = append(all, attrs...)
	builder.selfClose("polygon", all...)
}

// renderTrapezoid renders a trapezoid polygon.
func renderTrapezoid(builder *svgBuilder, posX, posY, width, height float32, alt bool, attrs []string) {
	offset := width * parallelogramSkew
	var pts [][2]float32
	if !alt {
		pts = [][2]float32{
			{posX + offset, posY},
			{posX + width - offset, posY},
			{posX + width, posY + height},
			{posX, posY + height},
		}
	} else {
		pts = [][2]float32{
			{posX, posY},
			{posX + width, posY},
			{posX + width - offset, posY + height},
			{posX + offset, posY + height},
		}
	}
	all := make([]string, 0, 2+len(attrs))
	all = append(all, "points", formatPoints(pts))
	all = append(all, attrs...)
	builder.selfClose("polygon", all...)
}

// renderNodeLabel renders text lines centered within a node.
func renderNodeLabel(builder *svgBuilder, node *layout.NodeLayout, textColor string) {
	if len(node.Label.Lines) == 0 {
		return
	}

	fontSize := node.Label.FontSize
	if fontSize <= 0 {
		fontSize = 14
	}

	lineHeight := fontSize * lineHeightScale
	totalTextHeight := lineHeight * float32(len(node.Label.Lines))
	// Start Y so that text block is vertically centered in node.
	startY := node.Y - totalTextHeight/2 + lineHeight*textBaselineOffset

	for idx, line := range node.Label.Lines {
		ly := startY + float32(idx)*lineHeight
		builder.text(node.X, ly, line,
			"text-anchor", "middle",
			"dominant-baseline", "auto",
			"fill", textColor,
			"font-size", fmtFloat(fontSize),
		)
	}
}

// clamp restricts val to the range [lo, hi].
func clamp(val, lo, hi float32) float32 {
	if val < lo {
		return lo
	}
	if val > hi {
		return hi
	}
	return val
}
