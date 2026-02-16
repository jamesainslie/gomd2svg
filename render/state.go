package render

import (
	"sort"
	"strings"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// State diagram rendering constants.
const (
	stateEndInnerScale    float32 = 0.6
	stateCompositeRadius  float32 = 8
	stateLabelOffsetX     float32 = 10
	stateLabelPadY        float32 = 4
	stateRegularRadius    float32 = 10
	stateLineHeightScale  float32 = 1.2
	stateNamePadY         float32 = 4
	stateDividerPadY      float32 = 6
	stateDividerInsetX    float32 = 4
	stateDescFontScale    float32 = 0.9
	stateNodeBorderRadius float32 = 6
)

// renderState renders all state diagram elements: edges, then nodes.
func renderState(builder *svgBuilder, computed *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	sd, ok := computed.Diagram.(layout.StateData)
	if !ok {
		return
	}

	renderStateEdges(builder, computed, th)
	renderStateNodes(builder, computed, th, cfg, &sd)
}

// renderStateEdges renders state transitions. Reuses the same edge rendering
// logic as the flowchart renderer.
func renderStateEdges(builder *svgBuilder, computed *layout.Layout, th *theme.Theme) {
	renderEdges(builder, computed, th)
}

// renderStateNodes renders state nodes sorted by ID for deterministic output.
func renderStateNodes(builder *svgBuilder, computed *layout.Layout, th *theme.Theme, cfg *config.Layout, sd *layout.StateData) {
	ids := make([]string, 0, len(computed.Nodes))
	for id := range computed.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		node := computed.Nodes[id]

		// Start pseudo-state: filled black circle.
		if strings.HasPrefix(id, "__start__") {
			renderStartState(builder, node, th)
			continue
		}

		// End pseudo-state: bullseye (outer ring + inner filled circle).
		if strings.HasPrefix(id, "__end__") {
			renderEndState(builder, node, th)
			continue
		}

		// Fork/join annotation: horizontal bar.
		if ann, ok := sd.Annotations[id]; ok {
			switch ann {
			case ir.StateFork, ir.StateJoin:
				renderForkJoinState(builder, node, th)
				continue
			case ir.StateChoice:
				renderChoiceState(builder, node, th)
				continue
			}
		}

		// Composite state: outer container with inner layout.
		if cs, ok := sd.CompositeStates[id]; ok {
			if innerLayout, hasInner := sd.InnerLayouts[id]; hasInner {
				renderCompositeState(builder, node, innerLayout, cs.Label, th, cfg)
				continue
			}
		}

		// Regular state.
		desc := sd.Descriptions[id]
		renderRegularState(builder, node, desc, th)
	}
}

// renderStartState renders a filled black circle for the start pseudo-state.
func renderStartState(builder *svgBuilder, node *layout.NodeLayout, th *theme.Theme) {
	cx := node.X
	cy := node.Y
	radius := node.Width / 2
	builder.circle(cx, cy, radius,
		"fill", th.StateStartEnd,
		"stroke", th.StateStartEnd,
		"stroke-width", "1",
	)
}

// renderEndState renders a bullseye for the end pseudo-state:
// an outer circle with stroke only and an inner filled circle.
func renderEndState(builder *svgBuilder, node *layout.NodeLayout, th *theme.Theme) {
	cx := node.X
	cy := node.Y
	outerR := node.Width / 2
	innerR := outerR * stateEndInnerScale

	// Outer circle (stroke only).
	builder.circle(cx, cy, outerR,
		"fill", "none",
		"stroke", th.StateStartEnd,
		"stroke-width", "1.5",
	)

	// Inner filled circle.
	builder.circle(cx, cy, innerR,
		"fill", th.StateStartEnd,
		"stroke", th.StateStartEnd,
		"stroke-width", "1",
	)
}

// renderForkJoinState renders a filled horizontal bar for fork/join pseudo-states.
func renderForkJoinState(builder *svgBuilder, node *layout.NodeLayout, th *theme.Theme) {
	posX := node.X - node.Width/2
	posY := node.Y - node.Height/2
	builder.rect(posX, posY, node.Width, node.Height, 2,
		"fill", th.StateStartEnd,
		"stroke", th.StateStartEnd,
		"stroke-width", "1",
	)
}

// renderChoiceState renders a diamond polygon for choice pseudo-states.
func renderChoiceState(builder *svgBuilder, node *layout.NodeLayout, th *theme.Theme) {
	cx := node.X
	cy := node.Y
	hw := node.Width / 2
	hh := node.Height / 2
	pts := [][2]float32{
		{cx, cy - hh},
		{cx + hw, cy},
		{cx, cy + hh},
		{cx - hw, cy},
	}
	builder.polygon(pts,
		"fill", th.StateFill,
		"stroke", th.StateBorder,
		"stroke-width", "1.5",
	)
}

// renderCompositeState renders a composite state as a rounded rect container
// with a label at top, then recursively renders the inner layout offset to
// fit inside the container.
func renderCompositeState(builder *svgBuilder, node *layout.NodeLayout, innerLayout *layout.Layout, label string, th *theme.Theme, cfg *config.Layout) {
	posX := node.X - node.Width/2
	posY := node.Y - node.Height/2

	// Outer rounded rect (dashed border like subgraph).
	builder.rect(posX, posY, node.Width, node.Height, stateCompositeRadius,
		"fill", th.ClusterBackground,
		"stroke", th.StateBorder,
		"stroke-width", "1.5",
		"stroke-dasharray", "5,5",
	)

	// Label at top-left inside the container.
	labelX := posX + stateLabelOffsetX
	labelY := posY + th.FontSize + stateLabelPadY
	builder.text(labelX, labelY, label,
		"fill", th.TextColor,
		"font-size", fmtFloat(th.FontSize),
		"font-weight", "bold",
	)

	// Render inner layout offset to fit inside the container.
	// The inner content starts below the label area.
	labelAreaH := th.FontSize*cfg.LabelLineHeight + cfg.Padding.NodeVertical
	offsetX := posX + cfg.Padding.NodeHorizontal
	offsetY := posY + labelAreaH

	renderInnerLayout(builder, innerLayout, th, cfg, offsetX, offsetY)
}

// renderInnerLayout renders a nested layout at a given offset by translating
// all node and edge coordinates.
func renderInnerLayout(builder *svgBuilder, computed *layout.Layout, th *theme.Theme, cfg *config.Layout, offsetX, offsetY float32) {
	builder.openTag("g",
		"transform", "translate("+fmtFloat(offsetX)+","+fmtFloat(offsetY)+")",
	)

	// Render edges.
	renderEdges(builder, computed, th)

	// Render nodes — delegate to appropriate renderer based on diagram type.
	switch diag := computed.Diagram.(type) {
	case layout.StateData:
		renderStateNodes(builder, computed, th, cfg, &diag)
	default:
		renderNodes(builder, computed, th)
	}

	builder.closeTag("g")
}

// renderRegularState renders a regular state node as a rounded rect with
// the state name centered. If a description is present, a divider line
// separates the name from the description text below.
func renderRegularState(builder *svgBuilder, node *layout.NodeLayout, description string, th *theme.Theme) {
	posX := node.X - node.Width/2
	posY := node.Y - node.Height/2

	// Rounded rect background.
	builder.rect(posX, posY, node.Width, node.Height, stateRegularRadius,
		"fill", th.StateFill,
		"stroke", th.StateBorder,
		"stroke-width", "1.5",
	)

	fontSize := node.Label.FontSize
	if fontSize <= 0 {
		fontSize = th.FontSize
	}
	lineHeight := fontSize * stateLineHeightScale

	if description == "" {
		// No description: center the label vertically.
		renderNodeLabel(builder, node, th.TextColor)
	} else {
		// With description: name at top, divider, description below.
		nameY := posY + lineHeight + stateNamePadY
		builder.text(node.X, nameY, node.Label.Lines[0],
			"text-anchor", "middle",
			"dominant-baseline", "auto",
			"fill", th.TextColor,
			"font-size", fmtFloat(fontSize),
			"font-weight", "bold",
		)

		// Divider line.
		dividerY := nameY + stateDividerPadY
		builder.line(posX+stateDividerInsetX, dividerY, posX+node.Width-stateDividerInsetX, dividerY,
			"stroke", th.StateBorder,
			"stroke-width", "0.5",
		)

		// Description text below divider.
		descY := dividerY + lineHeight
		builder.text(node.X, descY, description,
			"text-anchor", "middle",
			"dominant-baseline", "auto",
			"fill", th.TextColor,
			"font-size", fmtFloat(fontSize*stateDescFontScale),
		)
	}
}
