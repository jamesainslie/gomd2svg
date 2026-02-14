package render

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// renderGraph renders all flowchart/graph elements: subgraphs, edges, and nodes.
func renderGraph(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	// Render subgraphs first (they appear behind nodes and edges).
	renderSubgraphs(b, l, th)

	// Render edges.
	renderEdges(b, l, th)

	// Render nodes (on top of edges).
	renderNodes(b, l, th)
}

// renderSubgraphs renders subgraph containers as rectangles with labels.
func renderSubgraphs(b *svgBuilder, l *layout.Layout, th *theme.Theme) {
	for _, sg := range l.Subgraphs {
		b.rect(sg.X, sg.Y, sg.Width, sg.Height, 4,
			"fill", th.ClusterBackground,
			"stroke", th.ClusterBorder,
			"stroke-width", "1",
			"stroke-dasharray", "5,5",
		)

		// Subgraph label at top-left.
		if sg.Label != "" {
			labelX := sg.X + 8
			labelY := sg.Y + 16
			b.text(labelX, labelY, sg.Label,
				"fill", th.TextColor,
				"font-size", fmtFloat(th.FontSize*0.9),
				"font-weight", "bold",
			)
		}
	}
}

// renderEdges renders all edges as SVG paths with optional arrow markers.
func renderEdges(b *svgBuilder, l *layout.Layout, th *theme.Theme) {
	for i, edge := range l.Edges {
		if len(edge.Points) < 2 {
			continue
		}

		d := pointsToPath(edge.Points)
		edgeID := fmt.Sprintf("edge-%d", i)

		strokeColor := th.LineColor
		strokeWidth := "1.5"

		attrs := []string{
			"id", edgeID,
			"class", "edgePath",
			"d", d,
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

		b.selfClose("path", attrs...)

		// Render edge label if present.
		if edge.Label != nil && len(edge.Label.Lines) > 0 {
			renderEdgeLabel(b, edge, th)
		}
	}
}

// renderEdgeLabel renders the label for an edge at its anchor point.
func renderEdgeLabel(b *svgBuilder, edge *layout.EdgeLayout, th *theme.Theme) {
	label := edge.Label
	anchorX := edge.LabelAnchor[0]
	anchorY := edge.LabelAnchor[1]

	padX := float32(4)
	padY := float32(2)

	// Background rect.
	bgW := label.Width + padX*2
	bgH := label.Height + padY*2
	b.rect(anchorX-bgW/2, anchorY-bgH/2, bgW, bgH, 2,
		"fill", th.EdgeLabelBackground,
		"stroke", "none",
	)

	// Label text.
	fontSize := label.FontSize
	if fontSize <= 0 {
		fontSize = th.FontSize * 0.85
	}
	lineHeight := fontSize * 1.2
	totalH := lineHeight * float32(len(label.Lines))
	startY := anchorY - totalH/2 + lineHeight*0.75

	for i, line := range label.Lines {
		ly := startY + float32(i)*lineHeight
		b.text(anchorX, ly, line,
			"text-anchor", "middle",
			"fill", th.LabelTextColor,
			"font-size", fmtFloat(fontSize),
		)
	}
}

// renderNodes renders all nodes sorted by ID for deterministic output.
func renderNodes(b *svgBuilder, l *layout.Layout, th *theme.Theme) {
	// Sort node IDs for deterministic rendering order.
	ids := make([]string, 0, len(l.Nodes))
	for id := range l.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		n := l.Nodes[id]

		// Determine colors: use node style overrides if set, otherwise theme defaults.
		fill := th.PrimaryColor
		if n.Style.Fill != nil {
			fill = *n.Style.Fill
		}

		stroke := th.PrimaryBorderColor
		if n.Style.Stroke != nil {
			stroke = *n.Style.Stroke
		}

		textColor := th.PrimaryTextColor
		if n.Style.TextColor != nil {
			textColor = *n.Style.TextColor
		}

		renderNodeShape(b, n, fill, stroke, textColor)
	}
}
