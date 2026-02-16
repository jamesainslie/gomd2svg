package render

import (
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// renderBlock renders a block diagram: edges behind nodes, each node colored
// by cycling through the theme's BlockColors palette.
func renderBlock(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, _ *config.Layout) {
	_, ok := lay.Diagram.(layout.BlockData)
	if !ok {
		return
	}

	colors := th.BlockColors
	if len(colors) == 0 {
		colors = []string{"#D4E6F1", "#D5F5E3", "#FCF3CF", "#FADBD8"}
	}

	borderColor := th.BlockNodeBorder
	if borderColor == "" {
		borderColor = "#3B6492"
	}

	// Render edges first so they appear behind nodes.
	renderEdges(builder, lay, th)

	// Sort node IDs for deterministic output.
	ids := make([]string, 0, len(lay.Nodes))
	for id := range lay.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	// Render each node using renderNodeShape with color cycling.
	for i, id := range ids {
		node := lay.Nodes[id]

		fill := colors[i%len(colors)]
		if node.Style.Fill != nil {
			fill = *node.Style.Fill
		}

		stroke := borderColor
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
