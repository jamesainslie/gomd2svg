package render

import (
	"fmt"
	"math"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderMindmap(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	md, ok := l.Diagram.(layout.MindmapData)
	if !ok || md.Root == nil {
		return
	}

	branchColors := th.MindmapBranchColors
	if len(branchColors) == 0 {
		branchColors = []string{"#4C78A8", "#72B7B2", "#EECA3B", "#F58518"}
	}

	// Draw connections first (behind nodes).
	renderMindmapConnections(b, md.Root, branchColors, th)

	// Draw nodes.
	renderMindmapNode(b, md.Root, branchColors, th)
}

func renderMindmapConnections(b *svgBuilder, node *layout.MindmapNodeLayout, colors []string, th *theme.Theme) {
	for _, child := range node.Children {
		color := colors[child.ColorIndex%len(colors)]
		// Draw a curved connection from parent center to child center
		// using a cubic bezier with horizontal control points.
		midX := (node.X + child.X) / 2
		d := fmt.Sprintf("M %s %s C %s %s, %s %s, %s %s",
			fmtFloat(node.X), fmtFloat(node.Y),
			fmtFloat(midX), fmtFloat(node.Y),
			fmtFloat(midX), fmtFloat(child.Y),
			fmtFloat(child.X), fmtFloat(child.Y),
		)
		b.path(d,
			"fill", "none",
			"stroke", color,
			"stroke-width", "2",
			"opacity", "0.6",
		)
		renderMindmapConnections(b, child, colors, th)
	}
}

func renderMindmapNode(b *svgBuilder, node *layout.MindmapNodeLayout, colors []string, th *theme.Theme) {
	color := colors[node.ColorIndex%len(colors)]
	cx := node.X
	cy := node.Y
	hw := node.Width / 2
	hh := node.Height / 2

	switch node.Shape {
	case ir.MindmapSquare:
		b.rect(cx-hw, cy-hh, node.Width, node.Height, 0,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
		)
	case ir.MindmapRounded:
		b.rect(cx-hw, cy-hh, node.Width, node.Height, 8,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
		)
	case ir.MindmapCircle:
		r := hw
		if hh > r {
			r = hh
		}
		b.circle(cx, cy, r,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
		)
	case ir.MindmapHexagon:
		// Draw hexagon using polygon.
		pts := make([][2]float32, 6)
		for i := 0; i < 6; i++ {
			angle := float64(i)*math.Pi/3 - math.Pi/6
			pts[i] = [2]float32{
				cx + hw*float32(math.Cos(angle)),
				cy + hh*float32(math.Sin(angle)),
			}
		}
		b.polygon(pts,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
		)
	case ir.MindmapBang:
		// Star-burst shape using a larger circle with a thicker border.
		r := hw
		if hh > r {
			r = hh
		}
		b.circle(cx, cy, r,
			"fill", color,
			"stroke", color,
			"stroke-width", "4",
		)
	case ir.MindmapCloud:
		// Cloud-like shape: ellipse with dashed stroke.
		b.ellipse(cx, cy, hw*1.1, hh*1.1,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
			"stroke-dasharray", "4 2",
		)
	default: // MindmapShapeDefault - no border, just text background.
		b.rect(cx-hw, cy-hh, node.Width, node.Height, 4,
			"fill", th.MindmapNodeFill,
			"stroke", "none",
		)
	}

	// Draw label.
	textColor := th.TextColor
	if node.Shape == ir.MindmapBang {
		textColor = "#FFFFFF"
	}
	b.text(cx, cy+th.FontSize/3, node.Label,
		"text-anchor", "middle",
		"font-family", th.FontFamily,
		"font-size", fmtFloat(th.FontSize),
		"fill", textColor,
	)

	// Recursively render children.
	for _, child := range node.Children {
		renderMindmapNode(b, child, colors, th)
	}
}
