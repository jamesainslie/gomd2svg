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
	case layout.ClassData:
		renderClass(&b, l, th, cfg)
	case layout.ERData:
		renderER(&b, l, th, cfg)
	case layout.StateData:
		renderState(&b, l, th, cfg)
	case layout.SequenceData:
		renderSequence(&b, l, th, cfg)
	case layout.KanbanData:
		renderKanban(&b, l, th, cfg)
	case layout.PacketData:
		renderPacket(&b, l, th, cfg)
	case layout.PieData:
		renderPie(&b, l, th, cfg)
	case layout.QuadrantData:
		renderQuadrant(&b, l, th, cfg)
	case layout.TimelineData:
		renderTimeline(&b, l, th, cfg)
	case layout.GanttData:
		renderGantt(&b, l, th, cfg)
	case layout.GitGraphData:
		renderGitGraph(&b, l, th, cfg)
	case layout.XYChartData:
		renderXYChart(&b, l, th, cfg)
	case layout.RadarData:
		renderRadar(&b, l, th, cfg)
	case layout.MindmapData:
		renderMindmap(&b, l, th, cfg)
	case layout.SankeyData:
		renderSankey(&b, l, th, cfg)
	case layout.TreemapData:
		renderTreemap(&b, l, th, cfg)
	case layout.RequirementData:
		renderRequirement(&b, l, th, cfg)
	case layout.BlockData:
		renderBlock(&b, l, th, cfg)
	case layout.C4Data:
		renderC4(&b, l, th, cfg)
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

	// Closed triangle (inheritance/realization) — forward
	b.raw(`<marker id="marker-closed-triangle" viewBox="0 0 20 20" refX="18" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	b.selfClose("path", "d", "M 0 0 L 20 10 L 0 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
	b.closeTag("marker")

	// Closed triangle — reverse
	b.raw(`<marker id="marker-closed-triangle-start" viewBox="0 0 20 20" refX="2" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	b.selfClose("path", "d", "M 20 0 L 0 10 L 20 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
	b.closeTag("marker")

	// Filled diamond (composition) — forward
	b.raw(`<marker id="marker-filled-diamond" viewBox="0 0 20 20" refX="18" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	b.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.LineColor, "stroke", th.LineColor, "stroke-width", "1")
	b.closeTag("marker")

	// Filled diamond — reverse
	b.raw(`<marker id="marker-filled-diamond-start" viewBox="0 0 20 20" refX="2" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	b.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.LineColor, "stroke", th.LineColor, "stroke-width", "1")
	b.closeTag("marker")

	// Open diamond (aggregation) — forward
	b.raw(`<marker id="marker-open-diamond" viewBox="0 0 20 20" refX="18" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	b.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
	b.closeTag("marker")

	// Open diamond — reverse
	b.raw(`<marker id="marker-open-diamond-start" viewBox="0 0 20 20" refX="2" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	b.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
	b.closeTag("marker")

	// Open arrowhead (async messages) — forward
	b.raw(`<marker id="marker-open-arrow" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="userSpaceOnUse" markerWidth="8" markerHeight="8" orient="auto">`)
	b.selfClose("path", "d", "M 0 0 L 10 5 L 0 10", "fill", "none", "stroke", th.LineColor, "stroke-width", "1.5")
	b.closeTag("marker")

	// Cross end (termination messages) — forward
	b.raw(`<marker id="marker-cross" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" orient="auto">`)
	b.selfClose("path", "d", "M 2 2 L 8 8 M 8 2 L 2 8", "fill", "none", "stroke", th.LineColor, "stroke-width", "1.5")
	b.closeTag("marker")

	b.closeTag("defs")
}
