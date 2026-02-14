package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
)

// positionNodes assigns X, Y coordinates to each node based on its rank,
// position within the rank, and the diagram direction. For LeftRight layouts
// the rank axis is horizontal (X) and the cross axis is vertical (Y).
// For TopDown layouts the rank axis is vertical (Y) and the cross axis is
// horizontal (X).
func positionNodes(
	layers [][]string,
	nodes map[string]*NodeLayout,
	direction ir.Direction,
	cfg *config.Layout,
) {
	if len(layers) == 0 {
		return
	}

	rankSpacing := cfg.RankSpacing
	nodeSpacing := cfg.NodeSpacing

	switch direction {
	case ir.LeftRight, ir.RightLeft:
		positionLR(layers, nodes, rankSpacing, nodeSpacing, direction == ir.RightLeft)
	default: // TopDown, BottomTop
		positionTD(layers, nodes, rankSpacing, nodeSpacing, direction == ir.BottomTop)
	}
}

// positionLR positions nodes in a left-to-right (or right-to-left) layout.
// Ranks map to X positions; cross-axis (within rank) maps to Y.
func positionLR(
	layers [][]string,
	nodes map[string]*NodeLayout,
	rankSpacing, nodeSpacing float32,
	reverse bool,
) {
	// First pass: compute X positions per rank (cumulative width + spacing).
	rankX := make([]float32, len(layers))
	var cumX float32
	for r, layer := range layers {
		// Find the widest node in this rank.
		var maxWidth float32
		for _, id := range layer {
			if n, ok := nodes[id]; ok && n.Width > maxWidth {
				maxWidth = n.Width
			}
		}
		rankX[r] = cumX + maxWidth/2
		cumX += maxWidth + rankSpacing
	}

	// Second pass: assign coordinates.
	for r, layer := range layers {
		// Compute total height of this rank's column.
		var totalHeight float32
		for i, id := range layer {
			if n, ok := nodes[id]; ok {
				totalHeight += n.Height
				if i > 0 {
					totalHeight += nodeSpacing
				}
			}
		}

		// Center column vertically around 0, then shift by boundary padding.
		y := -totalHeight/2 + layoutBoundaryPad

		for _, id := range layer {
			n, ok := nodes[id]
			if !ok {
				continue
			}
			x := rankX[r] + layoutBoundaryPad
			if reverse {
				x = cumX - rankX[r] + layoutBoundaryPad
			}
			n.X = x
			n.Y = y + n.Height/2
			y += n.Height + nodeSpacing
		}
	}
}

// positionTD positions nodes in a top-to-bottom (or bottom-to-top) layout.
// Ranks map to Y positions; cross-axis (within rank) maps to X.
func positionTD(
	layers [][]string,
	nodes map[string]*NodeLayout,
	rankSpacing, nodeSpacing float32,
	reverse bool,
) {
	// First pass: compute Y positions per rank.
	rankY := make([]float32, len(layers))
	var cumY float32
	for r, layer := range layers {
		var maxHeight float32
		for _, id := range layer {
			if n, ok := nodes[id]; ok && n.Height > maxHeight {
				maxHeight = n.Height
			}
		}
		rankY[r] = cumY + maxHeight/2
		cumY += maxHeight + rankSpacing
	}

	// Second pass: assign coordinates.
	for r, layer := range layers {
		var totalWidth float32
		for i, id := range layer {
			if n, ok := nodes[id]; ok {
				totalWidth += n.Width
				if i > 0 {
					totalWidth += nodeSpacing
				}
			}
		}

		x := -totalWidth/2 + layoutBoundaryPad

		for _, id := range layer {
			n, ok := nodes[id]
			if !ok {
				continue
			}
			y := rankY[r] + layoutBoundaryPad
			if reverse {
				y = cumY - rankY[r] + layoutBoundaryPad
			}
			n.X = x + n.Width/2
			n.Y = y
			x += n.Width + nodeSpacing
		}
	}
}
