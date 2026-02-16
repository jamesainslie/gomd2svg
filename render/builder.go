package render

import (
	"strconv"
	"strings"
)

// svgBuilder wraps a strings.Builder to produce SVG markup.
type svgBuilder struct {
	buf strings.Builder
}

// openTag writes an opening XML tag with optional attributes.
// Attributes are passed as alternating key, value pairs.
func (b *svgBuilder) openTag(name string, attrs ...string) {
	b.buf.WriteByte('<')
	b.buf.WriteString(name)
	writeAttrs(&b.buf, attrs)
	b.buf.WriteByte('>')
}

// closeTag writes a closing XML tag.
func (b *svgBuilder) closeTag(name string) {
	b.buf.WriteString("</")
	b.buf.WriteString(name)
	b.buf.WriteByte('>')
}

// selfClose writes a self-closing XML tag with optional attributes.
func (b *svgBuilder) selfClose(name string, attrs ...string) {
	b.buf.WriteByte('<')
	b.buf.WriteString(name)
	writeAttrs(&b.buf, attrs)
	b.buf.WriteString("/>")
}

// content writes escaped text content.
func (b *svgBuilder) content(text string) {
	b.buf.WriteString(escapeXML(text))
}

// raw writes a raw string without escaping.
func (b *svgBuilder) raw(str string) {
	b.buf.WriteString(str)
}

// String returns the accumulated SVG markup.
func (b *svgBuilder) String() string {
	return b.buf.String()
}

// rect renders an SVG <rect> element.
func (b *svgBuilder) rect(posX, posY, width, height, rx float32, attrs ...string) {
	all := make([]string, 0, 6+len(attrs)) //nolint:mnd // base capacity for x,y,width,height pairs
	all = append(all,
		"x", fmtFloat(posX),
		"y", fmtFloat(posY),
		"width", fmtFloat(width),
		"height", fmtFloat(height),
	)
	if rx > 0 {
		all = append(all, "rx", fmtFloat(rx), "ry", fmtFloat(rx))
	}
	all = append(all, attrs...)
	b.selfClose("rect", all...)
}

// circle renders an SVG <circle> element.
func (b *svgBuilder) circle(cx, cy, radius float32, attrs ...string) {
	all := make([]string, 0, 6+len(attrs)) //nolint:mnd // base capacity for cx,cy,r pairs
	all = append(all,
		"cx", fmtFloat(cx),
		"cy", fmtFloat(cy),
		"r", fmtFloat(radius),
	)
	all = append(all, attrs...)
	b.selfClose("circle", all...)
}

// ellipse renders an SVG <ellipse> element.
func (b *svgBuilder) ellipse(cx, cy, rx, ry float32, attrs ...string) {
	all := make([]string, 0, 8+len(attrs)) //nolint:mnd // base capacity for cx,cy,rx,ry pairs
	all = append(all,
		"cx", fmtFloat(cx),
		"cy", fmtFloat(cy),
		"rx", fmtFloat(rx),
		"ry", fmtFloat(ry),
	)
	all = append(all, attrs...)
	b.selfClose("ellipse", all...)
}

// path renders an SVG <path> element.
func (b *svgBuilder) path(pathData string, attrs ...string) {
	all := make([]string, 0, 2+len(attrs))
	all = append(all, "d", pathData)
	all = append(all, attrs...)
	b.selfClose("path", all...)
}

// text renders an SVG <text> element with content.
func (b *svgBuilder) text(posX, posY float32, content string, attrs ...string) {
	all := make([]string, 0, 4+len(attrs)) //nolint:mnd // base capacity for x,y pairs
	all = append(all,
		"x", fmtFloat(posX),
		"y", fmtFloat(posY),
	)
	all = append(all, attrs...)
	b.openTag("text", all...)
	b.content(content)
	b.closeTag("text")
}

// line renders an SVG <line> element.
func (b *svgBuilder) line(x1, y1, x2, y2 float32, attrs ...string) {
	all := make([]string, 0, 8+len(attrs)) //nolint:mnd // base capacity for x1,y1,x2,y2 pairs
	all = append(all,
		"x1", fmtFloat(x1),
		"y1", fmtFloat(y1),
		"x2", fmtFloat(x2),
		"y2", fmtFloat(y2),
	)
	all = append(all, attrs...)
	b.selfClose("line", all...)
}

// polygon renders an SVG <polygon> element.
func (b *svgBuilder) polygon(points [][2]float32, attrs ...string) {
	all := make([]string, 0, 2+len(attrs))
	all = append(all, "points", formatPoints(points))
	all = append(all, attrs...)
	b.selfClose("polygon", all...)
}

// writeAttrs writes key="value" attribute pairs to a builder.
// Values are XML-escaped to prevent injection via user-controlled content.
func writeAttrs(buf *strings.Builder, attrs []string) {
	for idx := 0; idx+1 < len(attrs); idx += 2 {
		buf.WriteByte(' ')
		buf.WriteString(attrs[idx])
		buf.WriteString("=\"")
		buf.WriteString(escapeXMLAttr(attrs[idx+1]))
		buf.WriteByte('"')
	}
}

// escapeXMLAttr escapes special characters in XML attribute values.
func escapeXMLAttr(str string) string {
	str = strings.ReplaceAll(str, "&", "&amp;")
	str = strings.ReplaceAll(str, "<", "&lt;")
	str = strings.ReplaceAll(str, ">", "&gt;")
	str = strings.ReplaceAll(str, "\"", "&quot;")
	return str
}

// formatPoints formats a slice of [2]float32 as an SVG points string.
func formatPoints(pts [][2]float32) string {
	parts := make([]string, len(pts))
	for idx, pt := range pts {
		parts[idx] = fmtFloat(pt[0]) + "," + fmtFloat(pt[1])
	}
	return strings.Join(parts, " ")
}

// pointsToPath builds an SVG path "d" attribute from a slice of points.
func pointsToPath(pts [][2]float32) string {
	if len(pts) == 0 {
		return ""
	}
	var buf strings.Builder
	buf.WriteString("M ")
	buf.WriteString(fmtFloat(pts[0][0]))
	buf.WriteByte(',')
	buf.WriteString(fmtFloat(pts[0][1]))
	for _, pt := range pts[1:] {
		buf.WriteString(" L ")
		buf.WriteString(fmtFloat(pt[0]))
		buf.WriteByte(',')
		buf.WriteString(fmtFloat(pt[1]))
	}
	return buf.String()
}

// escapeXML replaces XML special characters in text content.
func escapeXML(str string) string {
	str = strings.ReplaceAll(str, "&", "&amp;")
	str = strings.ReplaceAll(str, "<", "&lt;")
	str = strings.ReplaceAll(str, ">", "&gt;")
	str = strings.ReplaceAll(str, "\"", "&quot;")
	str = strings.ReplaceAll(str, "'", "&#39;")
	return str
}

// fmtFloat formats a float32 with no trailing zeros for compact SVG output.
func fmtFloat(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

// isTransparentColor reports whether a CSS color string already includes
// transparency (rgba, hsla, or the "transparent" keyword).
func isTransparentColor(colorStr string) bool {
	lower := strings.ToLower(colorStr)
	return strings.HasPrefix(lower, "rgba") ||
		strings.HasPrefix(lower, "hsla") ||
		lower == "transparent"
}
