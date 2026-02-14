// Package render produces SVG output from a computed layout.
package render

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// RenderSVG produces an SVG string from a computed layout, theme, and config.
func RenderSVG(l *layout.Layout, th *theme.Theme, cfg *config.Layout) string {
	var b svgBuilder

	width := l.Width
	if width < 1 {
		width = 1
	}
	height := l.Height
	if height < 1 {
		height = 1
	}

	// Open <svg> tag.
	b.openTag("svg",
		"xmlns", "http://www.w3.org/2000/svg",
		"width", fmtFloat(width),
		"height", fmtFloat(height),
		"viewBox", "0 0 "+fmtFloat(width)+" "+fmtFloat(height),
	)

	// Arrow marker definitions.
	renderDefs(&b, th)

	// Background.
	b.rect(0, 0, width, height, 0,
		"fill", th.Background,
	)

	// Dispatch based on diagram data type.
	switch l.Diagram.(type) {
	case layout.GraphData:
		renderGraph(&b, l, th, cfg)
	default:
		// For other diagram types, still render graph as a fallback.
		renderGraph(&b, l, th, cfg)
	}

	b.closeTag("svg")
	return b.String()
}

// renderDefs writes the <defs> block with reusable marker definitions.
func renderDefs(b *svgBuilder, th *theme.Theme) {
	b.openTag("defs")

	// Forward arrowhead marker.
	b.raw(`<marker id="arrowhead" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="userSpaceOnUse" markerWidth="8" markerHeight="8" orient="auto">`)
	b.selfClose("path",
		"d", "M 0 0 L 10 5 L 0 10 z",
		"fill", th.LineColor,
		"stroke", th.LineColor,
		"stroke-width", "1",
	)
	b.closeTag("marker")

	// Reverse arrowhead marker.
	b.raw(`<marker id="arrowhead-start" viewBox="0 0 10 10" refX="1" refY="5" markerUnits="userSpaceOnUse" markerWidth="8" markerHeight="8" orient="auto">`)
	b.selfClose("path",
		"d", "M 10 0 L 0 5 L 10 10 z",
		"fill", th.LineColor,
		"stroke", th.LineColor,
		"stroke-width", "1",
	)
	b.closeTag("marker")

	b.closeTag("defs")
}
