package render

import (
	"fmt"
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
func (b *svgBuilder) raw(s string) {
	b.buf.WriteString(s)
}

// String returns the accumulated SVG markup.
func (b *svgBuilder) String() string {
	return b.buf.String()
}

// rect renders an SVG <rect> element.
func (b *svgBuilder) rect(x, y, w, h, rx float32, attrs ...string) {
	all := []string{
		"x", fmtFloat(x),
		"y", fmtFloat(y),
		"width", fmtFloat(w),
		"height", fmtFloat(h),
	}
	if rx > 0 {
		all = append(all, "rx", fmtFloat(rx), "ry", fmtFloat(rx))
	}
	all = append(all, attrs...)
	b.selfClose("rect", all...)
}

// circle renders an SVG <circle> element.
func (b *svgBuilder) circle(cx, cy, r float32, attrs ...string) {
	all := []string{
		"cx", fmtFloat(cx),
		"cy", fmtFloat(cy),
		"r", fmtFloat(r),
	}
	all = append(all, attrs...)
	b.selfClose("circle", all...)
}

// ellipse renders an SVG <ellipse> element.
func (b *svgBuilder) ellipse(cx, cy, rx, ry float32, attrs ...string) {
	all := []string{
		"cx", fmtFloat(cx),
		"cy", fmtFloat(cy),
		"rx", fmtFloat(rx),
		"ry", fmtFloat(ry),
	}
	all = append(all, attrs...)
	b.selfClose("ellipse", all...)
}

// path renders an SVG <path> element.
func (b *svgBuilder) path(d string, attrs ...string) {
	all := []string{"d", d}
	all = append(all, attrs...)
	b.selfClose("path", all...)
}

// text renders an SVG <text> element with content.
func (b *svgBuilder) text(x, y float32, content string, attrs ...string) {
	all := []string{
		"x", fmtFloat(x),
		"y", fmtFloat(y),
	}
	all = append(all, attrs...)
	b.openTag("text", all...)
	b.content(content)
	b.closeTag("text")
}

// line renders an SVG <line> element.
func (b *svgBuilder) line(x1, y1, x2, y2 float32, attrs ...string) {
	all := []string{
		"x1", fmtFloat(x1),
		"y1", fmtFloat(y1),
		"x2", fmtFloat(x2),
		"y2", fmtFloat(y2),
	}
	all = append(all, attrs...)
	b.selfClose("line", all...)
}

// polygon renders an SVG <polygon> element.
func (b *svgBuilder) polygon(points [][2]float32, attrs ...string) {
	all := []string{"points", formatPoints(points)}
	all = append(all, attrs...)
	b.selfClose("polygon", all...)
}

// writeAttrs writes key="value" attribute pairs to a builder.
func writeAttrs(buf *strings.Builder, attrs []string) {
	for i := 0; i+1 < len(attrs); i += 2 {
		buf.WriteByte(' ')
		buf.WriteString(attrs[i])
		buf.WriteString("=\"")
		buf.WriteString(attrs[i+1])
		buf.WriteByte('"')
	}
}

// formatPoints formats a slice of [2]float32 as an SVG points string.
func formatPoints(pts [][2]float32) string {
	parts := make([]string, len(pts))
	for i, p := range pts {
		parts[i] = fmtFloat(p[0]) + "," + fmtFloat(p[1])
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
	for _, p := range pts[1:] {
		buf.WriteString(" L ")
		buf.WriteString(fmtFloat(p[0]))
		buf.WriteByte(',')
		buf.WriteString(fmtFloat(p[1]))
	}
	return buf.String()
}

// escapeXML replaces XML special characters in text content.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// fmtFloat formats a float32 with no trailing zeros for compact SVG output.
func fmtFloat(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

// fmtFloat2 formats a float32 to 2 decimal places.
func fmtFloat2(f float32) string {
	return fmt.Sprintf("%.2f", f)
}
