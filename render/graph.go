package render

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Graph rendering constants.
const (
	subgraphBorderRadius   float32 = 4
	subgraphLabelOffsetX   float32 = 8
	subgraphLabelOffsetY   float32 = 16
	subgraphFontScale      float32 = 0.9
	edgeLabelPadX          float32 = 4
	edgeLabelPadY          float32 = 2
	edgeLabelFontScale     float32 = 0.85
	edgeLabelLineHeight    float32 = 1.2
	edgeLabelBaselineShift float32 = 0.75
)

// renderGraph renders all flowchart/graph elements: subgraphs, edges, and nodes.
func renderGraph(builder *svgBuilder, computed *layout.Layout, th *theme.Theme, _ *config.Layout) {
	// Render subgraphs first (they appear behind nodes and edges).
	renderSubgraphs(builder, computed, th)

	// Render edges.
	renderEdges(builder, computed, th)

	// Render nodes (on top of edges).
	renderNodes(builder, computed, th)
}

// renderSubgraphs renders subgraph containers as rectangles with labels.
func renderSubgraphs(builder *svgBuilder, computed *layout.Layout, th *theme.Theme) {
	for _, sg := range computed.Subgraphs {
		builder.rect(sg.X, sg.Y, sg.Width, sg.Height, subgraphBorderRadius,
			"fill", th.ClusterBackground,
			"stroke", th.ClusterBorder,
			"stroke-width", "1",
			"stroke-dasharray", "5,5",
		)

		// Subgraph label at top-left.
		if sg.Label != "" {
			labelX := sg.X + subgraphLabelOffsetX
			labelY := sg.Y + subgraphLabelOffsetY
			builder.text(labelX, labelY, sg.Label,
				"fill", th.TextColor,
				"font-size", fmtFloat(th.FontSize*subgraphFontScale),
				"font-weight", "bold",
			)
		}
	}
}

// renderEdges renders all edges as SVG paths with optional arrow markers.
func renderEdges(builder *svgBuilder, computed *layout.Layout, th *theme.Theme) {
	for edgeIdx, edge := range computed.Edges {
		if len(edge.Points) < 2 {
			continue
		}

		pathData := pointsToPath(edge.Points)
		edgeID := fmt.Sprintf("edge-%d", edgeIdx)

		strokeColor := th.LineColor
		strokeWidth := "1.5"

		attrs := []string{
			"id", edgeID,
			"class", "edgePath",
			"d", pathData,
			"fill", "none",
			"stroke", strokeColor,
			"stroke-width", strokeWidth,
			"stroke-linecap", "round",
			"stroke-linejoin", "round",
		}

		// Edge style: dotted or thick.
		switch edge.Style {
		case ir.Dotted:
			attrs = append(attrs, "stroke-dasharray", "5,5")
		case ir.Thick:
			attrs = append(attrs, "stroke-width", "3")
		}

		// Arrow marker references.
		if edge.ArrowEnd {
			attrs = append(attrs, "marker-end", "url(#arrowhead)")
		}
		if edge.ArrowStart {
			attrs = append(attrs, "marker-start", "url(#arrowhead-start)")
		}

		builder.selfClose("path", attrs...)

		// Render edge label if present.
		if edge.Label != nil && len(edge.Label.Lines) > 0 {
			renderEdgeLabel(builder, edge, th)
		}
	}
}

// renderEdgeLabel renders the label for an edge at its anchor point.
func renderEdgeLabel(builder *svgBuilder, edge *layout.EdgeLayout, th *theme.Theme) {
	label := edge.Label
	anchorX := edge.LabelAnchor[0]
	anchorY := edge.LabelAnchor[1]

	// Background rect.
	bgW := label.Width + edgeLabelPadX*2
	bgH := label.Height + edgeLabelPadY*2
	builder.rect(anchorX-bgW/2, anchorY-bgH/2, bgW, bgH, 2,
		"fill", th.EdgeLabelBackground,
		"stroke", "none",
	)

	// Label text.
	fontSize := label.FontSize
	if fontSize <= 0 {
		fontSize = th.FontSize * edgeLabelFontScale
	}
	lineHeight := fontSize * edgeLabelLineHeight
	totalH := lineHeight * float32(len(label.Lines))
	startY := anchorY - totalH/2 + lineHeight*edgeLabelBaselineShift

	for idx, line := range label.Lines {
		ly := startY + float32(idx)*lineHeight
		builder.text(anchorX, ly, line,
			"text-anchor", "middle",
			"fill", th.LabelTextColor,
			"font-size", fmtFloat(fontSize),
		)
	}
}

// renderNodes renders all nodes sorted by ID for deterministic output.
func renderNodes(builder *svgBuilder, computed *layout.Layout, th *theme.Theme) {
	// Sort node IDs for deterministic rendering order.
	ids := make([]string, 0, len(computed.Nodes))
	for id := range computed.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		node := computed.Nodes[id]

		// Determine colors: use node style overrides if set, otherwise theme defaults.
		fill := th.PrimaryColor
		if node.Style.Fill != nil {
			fill = *node.Style.Fill
		}

		stroke := th.PrimaryBorderColor
		if node.Style.Stroke != nil {
			stroke = *node.Style.Stroke
		}

		textColor := th.PrimaryTextColor
		if node.Style.TextColor != nil {
			textColor = *node.Style.TextColor
		}

		renderNodeShape(builder, node, fill, stroke, textColor)
	}
}
