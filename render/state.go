package render

import (
	"sort"
	"strings"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// renderState renders all state diagram elements: edges, then nodes.
func renderState(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	sd, ok := l.Diagram.(layout.StateData)
	if !ok {
		return
	}

	renderStateEdges(b, l, th)
	renderStateNodes(b, l, th, cfg, &sd)
}

// renderStateEdges renders state transitions. Reuses the same edge rendering
// logic as the flowchart renderer.
func renderStateEdges(b *svgBuilder, l *layout.Layout, th *theme.Theme) {
	renderEdges(b, l, th)
}

// renderStateNodes renders state nodes sorted by ID for deterministic output.
func renderStateNodes(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout, sd *layout.StateData) {
	ids := make([]string, 0, len(l.Nodes))
	for id := range l.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		n := l.Nodes[id]

		// Start pseudo-state: filled black circle.
		if strings.HasPrefix(id, "__start__") {
			renderStartState(b, n, th)
			continue
		}

		// End pseudo-state: bullseye (outer ring + inner filled circle).
		if strings.HasPrefix(id, "__end__") {
			renderEndState(b, n, th)
			continue
		}

		// Fork/join annotation: horizontal bar.
		if ann, ok := sd.Annotations[id]; ok {
			switch ann {
			case ir.StateFork, ir.StateJoin:
				renderForkJoinState(b, n, th)
				continue
			case ir.StateChoice:
				renderChoiceState(b, n, th)
				continue
			}
		}

		// Composite state: outer container with inner layout.
		if cs, ok := sd.CompositeStates[id]; ok {
			if innerLayout, hasInner := sd.InnerLayouts[id]; hasInner {
				renderCompositeState(b, n, innerLayout, cs.Label, th, cfg)
				continue
			}
		}

		// Regular state.
		desc := sd.Descriptions[id]
		renderRegularState(b, n, desc, th)
	}
}

// renderStartState renders a filled black circle for the start pseudo-state.
func renderStartState(b *svgBuilder, n *layout.NodeLayout, th *theme.Theme) {
	cx := n.X
	cy := n.Y
	r := n.Width / 2
	b.circle(cx, cy, r,
		"fill", th.StateStartEnd,
		"stroke", th.StateStartEnd,
		"stroke-width", "1",
	)
}

// renderEndState renders a bullseye for the end pseudo-state:
// an outer circle with stroke only and an inner filled circle.
func renderEndState(b *svgBuilder, n *layout.NodeLayout, th *theme.Theme) {
	cx := n.X
	cy := n.Y
	outerR := n.Width / 2
	innerR := outerR * 0.6

	// Outer circle (stroke only).
	b.circle(cx, cy, outerR,
		"fill", "none",
		"stroke", th.StateStartEnd,
		"stroke-width", "1.5",
	)

	// Inner filled circle.
	b.circle(cx, cy, innerR,
		"fill", th.StateStartEnd,
		"stroke", th.StateStartEnd,
		"stroke-width", "1",
	)
}

// renderForkJoinState renders a filled horizontal bar for fork/join pseudo-states.
func renderForkJoinState(b *svgBuilder, n *layout.NodeLayout, th *theme.Theme) {
	x := n.X - n.Width/2
	y := n.Y - n.Height/2
	b.rect(x, y, n.Width, n.Height, 2,
		"fill", th.StateStartEnd,
		"stroke", th.StateStartEnd,
		"stroke-width", "1",
	)
}

// renderChoiceState renders a diamond polygon for choice pseudo-states.
func renderChoiceState(b *svgBuilder, n *layout.NodeLayout, th *theme.Theme) {
	cx := n.X
	cy := n.Y
	hw := n.Width / 2
	hh := n.Height / 2
	pts := [][2]float32{
		{cx, cy - hh},
		{cx + hw, cy},
		{cx, cy + hh},
		{cx - hw, cy},
	}
	b.polygon(pts,
		"fill", th.StateFill,
		"stroke", th.StateBorder,
		"stroke-width", "1.5",
	)
}

// renderCompositeState renders a composite state as a rounded rect container
// with a label at top, then recursively renders the inner layout offset to
// fit inside the container.
func renderCompositeState(b *svgBuilder, n *layout.NodeLayout, innerLayout *layout.Layout, label string, th *theme.Theme, cfg *config.Layout) {
	x := n.X - n.Width/2
	y := n.Y - n.Height/2

	// Outer rounded rect (dashed border like subgraph).
	b.rect(x, y, n.Width, n.Height, 8,
		"fill", th.ClusterBackground,
		"stroke", th.StateBorder,
		"stroke-width", "1.5",
		"stroke-dasharray", "5,5",
	)

	// Label at top-left inside the container.
	labelX := x + 10
	labelY := y + th.FontSize + 4
	b.text(labelX, labelY, label,
		"fill", th.TextColor,
		"font-size", fmtFloat(th.FontSize),
		"font-weight", "bold",
	)

	// Render inner layout offset to fit inside the container.
	// The inner content starts below the label area.
	labelAreaH := th.FontSize*cfg.LabelLineHeight + cfg.Padding.NodeVertical
	offsetX := x + cfg.Padding.NodeHorizontal
	offsetY := y + labelAreaH

	renderInnerLayout(b, innerLayout, th, cfg, offsetX, offsetY)
}

// renderInnerLayout renders a nested layout at a given offset by translating
// all node and edge coordinates.
func renderInnerLayout(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout, offsetX, offsetY float32) {
	b.openTag("g",
		"transform", "translate("+fmtFloat(offsetX)+","+fmtFloat(offsetY)+")",
	)

	// Render edges.
	renderEdges(b, l, th)

	// Render nodes — delegate to appropriate renderer based on diagram type.
	switch d := l.Diagram.(type) {
	case layout.StateData:
		renderStateNodes(b, l, th, cfg, &d)
	default:
		renderNodes(b, l, th)
	}

	b.closeTag("g")
}

// renderRegularState renders a regular state node as a rounded rect with
// the state name centered. If a description is present, a divider line
// separates the name from the description text below.
func renderRegularState(b *svgBuilder, n *layout.NodeLayout, description string, th *theme.Theme) {
	x := n.X - n.Width/2
	y := n.Y - n.Height/2

	// Rounded rect background.
	b.rect(x, y, n.Width, n.Height, 10,
		"fill", th.StateFill,
		"stroke", th.StateBorder,
		"stroke-width", "1.5",
	)

	fontSize := n.Label.FontSize
	if fontSize <= 0 {
		fontSize = th.FontSize
	}
	lineHeight := fontSize * 1.2

	if description == "" {
		// No description: center the label vertically.
		renderNodeLabel(b, n, th.TextColor)
	} else {
		// With description: name at top, divider, description below.
		nameY := y + lineHeight + 4
		b.text(n.X, nameY, n.Label.Lines[0],
			"text-anchor", "middle",
			"dominant-baseline", "auto",
			"fill", th.TextColor,
			"font-size", fmtFloat(fontSize),
			"font-weight", "bold",
		)

		// Divider line.
		dividerY := nameY + 6
		b.line(x+4, dividerY, x+n.Width-4, dividerY,
			"stroke", th.StateBorder,
			"stroke-width", "0.5",
		)

		// Description text below divider.
		descY := dividerY + lineHeight
		b.text(n.X, descY, description,
			"text-anchor", "middle",
			"dominant-baseline", "auto",
			"fill", th.TextColor,
			"font-size", fmtFloat(fontSize*0.9),
		)
	}
}
