package render

import (
	"fmt"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderSankey(b *svgBuilder, l *layout.Layout, th *theme.Theme, _ *config.Layout) {
	sd, ok := l.Diagram.(layout.SankeyData)
	if !ok {
		return
	}

	nodeColors := th.SankeyNodeColors
	if len(nodeColors) == 0 {
		nodeColors = []string{"#4C78A8", "#72B7B2", "#EECA3B", "#F58518"}
	}
	linkColor := th.SankeyLinkColor
	if linkColor == "" {
		linkColor = "#888"
	}
	linkOpacity := th.SankeyLinkOpacity
	if linkOpacity == 0 {
		linkOpacity = 0.4
	}

	// Draw links first (behind nodes).
	for _, link := range sd.Links {
		if link.SourceIdx >= len(sd.Nodes) || link.TargetIdx >= len(sd.Nodes) {
			continue
		}
		src := sd.Nodes[link.SourceIdx]
		tgt := sd.Nodes[link.TargetIdx]

		// Cubic bezier from right side of source to left side of target.
		sx := src.X + src.Width
		sy := link.SourceY + link.Width/2
		tx := tgt.X
		ty := link.TargetY + link.Width/2
		midX := (sx + tx) / 2

		d := fmt.Sprintf("M %s,%s C %s,%s %s,%s %s,%s",
			fmtFloat(sx), fmtFloat(sy),
			fmtFloat(midX), fmtFloat(sy),
			fmtFloat(midX), fmtFloat(ty),
			fmtFloat(tx), fmtFloat(ty),
		)
		b.path(d,
			"fill", "none",
			"stroke", linkColor,
			"stroke-width", fmtFloat(link.Width),
			"stroke-opacity", fmt.Sprintf("%.2f", linkOpacity),
		)
	}

	// Draw node rectangles.
	for _, node := range sd.Nodes {
		color := nodeColors[node.ColorIndex%len(nodeColors)]
		b.rect(node.X, node.Y, node.Width, node.Height, 0,
			"fill", color,
		)

		// Label to the right of the node.
		labelX := node.X + node.Width + 4
		labelY := node.Y + node.Height/2 + th.FontSize/3
		b.text(labelX, labelY, node.Label,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize),
			"fill", th.TextColor,
		)
	}
}
