package layout

import (
	"strings"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

// sizeNodes computes the width and height for each node based on its label
// text, shape, and padding configuration. It returns a map of NodeLayout
// keyed by node ID.
func sizeNodes(
	nodes map[string]*ir.Node,
	measurer *textmetrics.Measurer,
	th *theme.Theme,
	cfg *config.Layout,
) map[string]*NodeLayout {
	result := make(map[string]*NodeLayout, len(nodes))
	for id, node := range nodes {
		nl := sizeNode(node, measurer, th, cfg)
		result[id] = nl
	}
	return result
}

// sizeNode computes layout dimensions for a single node.
func sizeNode(
	node *ir.Node,
	measurer *textmetrics.Measurer,
	th *theme.Theme,
	cfg *config.Layout,
) *NodeLayout {
	fontSize := th.FontSize
	fontFamily := th.FontFamily

	// Split label into lines and measure each.
	lines := strings.Split(node.Label, "\n")
	lineHeight := fontSize * cfg.LabelLineHeight

	var maxLineWidth float32
	for _, line := range lines {
		w := measurer.Width(line, fontSize, fontFamily)
		if w > maxLineWidth {
			maxLineWidth = w
		}
	}

	textWidth := maxLineWidth
	textHeight := lineHeight * float32(len(lines))

	tb := TextBlock{
		Lines:    lines,
		Width:    textWidth,
		Height:   textHeight,
		FontSize: fontSize,
	}

	// Apply padding.
	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical
	width := textWidth + 2*padH
	height := textHeight + 2*padV

	// Shape-specific adjustments.
	switch node.Shape {
	case ir.Diamond:
		// Diamond inscribes the text box; needs to be a square with
		// side = sqrt(2) * max(textW, textH) to contain the content.
		side := width
		if height > side {
			side = height
		}
		side *= 1.42 // approx sqrt(2)
		width = side
		height = side

	case ir.Hexagon:
		// Hexagons need extra horizontal space for the angled sides.
		width += flowchartPadCross

	case ir.Parallelogram, ir.ParallelogramAlt:
		// Parallelograms need extra horizontal space for the skew.
		width += flowchartPadCross * 0.5

	case ir.Circle, ir.DoubleCircle:
		// Circle must contain the text; diameter = diagonal of text box.
		diag := width
		if height > diag {
			diag = height
		}
		diag *= 1.15 // slight padding beyond inscribed
		width = diag
		height = diag
	}

	return &NodeLayout{
		ID:     node.ID,
		Label:  tb,
		Shape:  node.Shape,
		Width:  width,
		Height: height,
	}
}

// Layout constants matching the Rust reference implementation.
const (
	flowchartPadMain  = 40.0 // main-axis padding
	flowchartPadCross = 30.0 // cross-axis padding
	layoutBoundaryPad = 8.0  // final canvas padding
)
