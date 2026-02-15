package render

import (
	"sort"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// renderBlock renders a block diagram: edges behind nodes, each node colored
// by cycling through the theme's BlockColors palette.
func renderBlock(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	_, ok := l.Diagram.(layout.BlockData)
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
	renderEdges(b, l, th)

	// Sort node IDs for deterministic output.
	ids := make([]string, 0, len(l.Nodes))
	for id := range l.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	// Render each node using renderNodeShape with color cycling.
	for i, id := range ids {
		n := l.Nodes[id]

		fill := colors[i%len(colors)]
		if n.Style.Fill != nil {
			fill = *n.Style.Fill
		}

		stroke := borderColor
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
