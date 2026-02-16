// Package render produces SVG output from a computed layout.
package render

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// RenderSVG produces an SVG string from a computed layout, theme, and config.
func RenderSVG(computed *layout.Layout, th *theme.Theme, cfg *config.Layout) string {
	var builder svgBuilder

	width := computed.Width
	if width < 1 {
		width = 1
	}
	height := computed.Height
	if height < 1 {
		height = 1
	}

	// Compute accessibility attributes.
	ariaLabel := computed.Title
	if ariaLabel == "" {
		ariaLabel = computed.Kind.String() + " diagram"
	}

	// Open <svg> tag. Set font-family so all <text> elements inherit it.
	builder.openTag("svg",
		"xmlns", "http://www.w3.org/2000/svg",
		"width", fmtFloat(width),
		"height", fmtFloat(height),
		"viewBox", "0 0 "+fmtFloat(width)+" "+fmtFloat(height),
		"font-family", th.FontFamily,
		"role", "img",
		"aria-label", ariaLabel,
	)

	// Arrow marker definitions.
	renderDefs(&builder, th)

	// Background.
	builder.rect(0, 0, width, height, 0,
		"fill", th.Background,
	)

	// Title element for accessibility (only when a diagram title is set).
	if computed.Title != "" {
		builder.openTag("title")
		builder.content(computed.Title)
		builder.closeTag("title")
	}

	// Dispatch based on diagram data type.
	switch computed.Diagram.(type) {
	case layout.GraphData:
		renderGraph(&builder, computed, th, cfg)
	case layout.ClassData:
		renderClass(&builder, computed, th, cfg)
	case layout.ERData:
		renderER(&builder, computed, th, cfg)
	case layout.StateData:
		renderState(&builder, computed, th, cfg)
	case layout.SequenceData:
		renderSequence(&builder, computed, th, cfg)
	case layout.KanbanData:
		renderKanban(&builder, computed, th, cfg)
	case layout.PacketData:
		renderPacket(&builder, computed, th, cfg)
	case layout.PieData:
		renderPie(&builder, computed, th, cfg)
	case layout.QuadrantData:
		renderQuadrant(&builder, computed, th, cfg)
	case layout.TimelineData:
		renderTimeline(&builder, computed, th, cfg)
	case layout.GanttData:
		renderGantt(&builder, computed, th, cfg)
	case layout.GitGraphData:
		renderGitGraph(&builder, computed, th, cfg)
	case layout.XYChartData:
		renderXYChart(&builder, computed, th, cfg)
	case layout.RadarData:
		renderRadar(&builder, computed, th, cfg)
	case layout.MindmapData:
		renderMindmap(&builder, computed, th, cfg)
	case layout.SankeyData:
		renderSankey(&builder, computed, th, cfg)
	case layout.TreemapData:
		renderTreemap(&builder, computed, th, cfg)
	case layout.RequirementData:
		renderRequirement(&builder, computed, th, cfg)
	case layout.BlockData:
		renderBlock(&builder, computed, th, cfg)
	case layout.C4Data:
		renderC4(&builder, computed, th, cfg)
	case layout.JourneyData:
		renderJourney(&builder, computed, th, cfg)
	case layout.ArchitectureData:
		renderArchitecture(&builder, computed, th, cfg)
	default:
		// For other diagram types, still render graph as a fallback.
		renderGraph(&builder, computed, th, cfg)
	}

	builder.closeTag("svg")
	return builder.String()
}

// renderDefs writes the <defs> block with reusable marker definitions.
func renderDefs(builder *svgBuilder, th *theme.Theme) {
	builder.openTag("defs")

	// Forward arrowhead marker.
	builder.raw(`<marker id="arrowhead" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="userSpaceOnUse" markerWidth="8" markerHeight="8" orient="auto">`)
	builder.selfClose("path",
		"d", "M 0 0 L 10 5 L 0 10 z",
		"fill", th.LineColor,
		"stroke", th.LineColor,
		"stroke-width", "1",
	)
	builder.closeTag("marker")

	// Reverse arrowhead marker.
	builder.raw(`<marker id="arrowhead-start" viewBox="0 0 10 10" refX="1" refY="5" markerUnits="userSpaceOnUse" markerWidth="8" markerHeight="8" orient="auto">`)
	builder.selfClose("path",
		"d", "M 10 0 L 0 5 L 10 10 z",
		"fill", th.LineColor,
		"stroke", th.LineColor,
		"stroke-width", "1",
	)
	builder.closeTag("marker")

	// Closed triangle (inheritance/realization) — forward
	builder.raw(`<marker id="marker-closed-triangle" viewBox="0 0 20 20" refX="18" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	builder.selfClose("path", "d", "M 0 0 L 20 10 L 0 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
	builder.closeTag("marker")

	// Closed triangle — reverse
	builder.raw(`<marker id="marker-closed-triangle-start" viewBox="0 0 20 20" refX="2" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	builder.selfClose("path", "d", "M 20 0 L 0 10 L 20 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
	builder.closeTag("marker")

	// Filled diamond (composition) — forward
	builder.raw(`<marker id="marker-filled-diamond" viewBox="0 0 20 20" refX="18" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	builder.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.LineColor, "stroke", th.LineColor, "stroke-width", "1")
	builder.closeTag("marker")

	// Filled diamond — reverse
	builder.raw(`<marker id="marker-filled-diamond-start" viewBox="0 0 20 20" refX="2" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	builder.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.LineColor, "stroke", th.LineColor, "stroke-width", "1")
	builder.closeTag("marker")

	// Open diamond (aggregation) — forward
	builder.raw(`<marker id="marker-open-diamond" viewBox="0 0 20 20" refX="18" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	builder.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
	builder.closeTag("marker")

	// Open diamond — reverse
	builder.raw(`<marker id="marker-open-diamond-start" viewBox="0 0 20 20" refX="2" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
	builder.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
	builder.closeTag("marker")

	// Open arrowhead (async messages) — forward
	builder.raw(`<marker id="marker-open-arrow" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="userSpaceOnUse" markerWidth="8" markerHeight="8" orient="auto">`)
	builder.selfClose("path", "d", "M 0 0 L 10 5 L 0 10", "fill", "none", "stroke", th.LineColor, "stroke-width", "1.5")
	builder.closeTag("marker")

	// Cross end (termination messages) — forward
	builder.raw(`<marker id="marker-cross" viewBox="0 0 10 10" refX="9" refY="5" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" orient="auto">`)
	builder.selfClose("path", "d", "M 2 2 L 8 8 M 8 2 L 2 8", "fill", "none", "stroke", th.LineColor, "stroke-width", "1.5")
	builder.closeTag("marker")

	builder.closeTag("defs")
}
