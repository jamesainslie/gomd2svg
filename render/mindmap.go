package render

import (
	"fmt"
	"math"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Mindmap rendering constants.
const (
	mindmapRoundedRadius float32 = 8
	mindmapDefaultRadius float32 = 4
	mindmapHexagonSides  int     = 6
	mindmapHexAngleDelta float64 = math.Pi / 6
	mindmapCloudScale    float32 = 1.1
)

func renderMindmap(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, _ *config.Layout) {
	md, ok := lay.Diagram.(layout.MindmapData)
	if !ok || md.Root == nil {
		return
	}

	branchColors := th.MindmapBranchColors
	if len(branchColors) == 0 {
		branchColors = []string{"#4C78A8", "#72B7B2", "#EECA3B", "#F58518"}
	}

	// Draw connections first (behind nodes).
	renderMindmapConnections(builder, md.Root, branchColors)

	// Draw nodes.
	renderMindmapNode(builder, md.Root, branchColors, th)
}

func renderMindmapConnections(builder *svgBuilder, node *layout.MindmapNodeLayout, colors []string) {
	for _, child := range node.Children {
		color := colors[child.ColorIndex%len(colors)]
		// Draw a curved connection from parent center to child center
		// using a cubic bezier with horizontal control points.
		midX := (node.X + child.X) / 2
		pathData := fmt.Sprintf("M %s %s C %s %s, %s %s, %s %s",
			fmtFloat(node.X), fmtFloat(node.Y),
			fmtFloat(midX), fmtFloat(node.Y),
			fmtFloat(midX), fmtFloat(child.Y),
			fmtFloat(child.X), fmtFloat(child.Y),
		)
		builder.path(pathData,
			"fill", "none",
			"stroke", color,
			"stroke-width", "2",
			"opacity", "0.6",
		)
		renderMindmapConnections(builder, child, colors)
	}
}

func renderMindmapNode(builder *svgBuilder, node *layout.MindmapNodeLayout, colors []string, th *theme.Theme) {
	color := colors[node.ColorIndex%len(colors)]
	cx := node.X
	cy := node.Y
	hw := node.Width / 2
	hh := node.Height / 2

	switch node.Shape {
	case ir.MindmapSquare:
		builder.rect(cx-hw, cy-hh, node.Width, node.Height, 0,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
		)
	case ir.MindmapRounded:
		builder.rect(cx-hw, cy-hh, node.Width, node.Height, mindmapRoundedRadius,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
		)
	case ir.MindmapCircle:
		radius := hw
		if hh > radius {
			radius = hh
		}
		builder.circle(cx, cy, radius,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
		)
	case ir.MindmapHexagon:
		// Draw hexagon using polygon.
		pts := make([][2]float32, mindmapHexagonSides)
		for idx := range mindmapHexagonSides {
			angle := float64(idx)*math.Pi/3 - mindmapHexAngleDelta
			pts[idx] = [2]float32{
				cx + hw*float32(math.Cos(angle)),
				cy + hh*float32(math.Sin(angle)),
			}
		}
		builder.polygon(pts,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
		)
	case ir.MindmapBang:
		// Star-burst shape using a larger circle with a thicker border.
		radius := hw
		if hh > radius {
			radius = hh
		}
		builder.circle(cx, cy, radius,
			"fill", color,
			"stroke", color,
			"stroke-width", "4",
		)
	case ir.MindmapCloud:
		// Cloud-like shape: ellipse with dashed stroke.
		builder.ellipse(cx, cy, hw*mindmapCloudScale, hh*mindmapCloudScale,
			"fill", th.MindmapNodeFill,
			"stroke", color,
			"stroke-width", "2",
			"stroke-dasharray", "4 2",
		)
	default: // MindmapShapeDefault - no border, just text background.
		builder.rect(cx-hw, cy-hh, node.Width, node.Height, mindmapDefaultRadius,
			"fill", th.MindmapNodeFill,
			"stroke", "none",
		)
	}

	// Draw label.
	textColor := th.TextColor
	if node.Shape == ir.MindmapBang {
		textColor = "#FFFFFF"
	}
	builder.text(cx, cy+th.FontSize/3, node.Label,
		"text-anchor", "middle",
		"font-family", th.FontFamily,
		"font-size", fmtFloat(th.FontSize),
		"fill", textColor,
	)

	// Recursively render children.
	for _, child := range node.Children {
		renderMindmapNode(builder, child, colors, th)
	}
}
