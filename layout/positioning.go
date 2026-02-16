package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
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
	for rankIdx, layer := range layers {
		// Find the widest node in this rank.
		var maxWidth float32
		for _, id := range layer {
			if node, ok := nodes[id]; ok && node.Width > maxWidth {
				maxWidth = node.Width
			}
		}
		rankX[rankIdx] = cumX + maxWidth/2
		cumX += maxWidth + rankSpacing
	}

	// Second pass: assign coordinates.
	for rankIdx, layer := range layers {
		// Compute total height of this rank's column.
		var totalHeight float32
		for idx, id := range layer {
			if node, ok := nodes[id]; ok {
				totalHeight += node.Height
				if idx > 0 {
					totalHeight += nodeSpacing
				}
			}
		}

		// Center column vertically around 0, then shift by boundary padding.
		posY := -totalHeight/2 + layoutBoundaryPad

		for _, id := range layer {
			node, ok := nodes[id]
			if !ok {
				continue
			}
			posX := rankX[rankIdx] + layoutBoundaryPad
			if reverse {
				posX = cumX - rankX[rankIdx] + layoutBoundaryPad
			}
			node.X = posX
			node.Y = posY + node.Height/2
			posY += node.Height + nodeSpacing
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
	for rankIdx, layer := range layers {
		var maxHeight float32
		for _, id := range layer {
			if node, ok := nodes[id]; ok && node.Height > maxHeight {
				maxHeight = node.Height
			}
		}
		rankY[rankIdx] = cumY + maxHeight/2
		cumY += maxHeight + rankSpacing
	}

	// Second pass: assign coordinates.
	for rankIdx, layer := range layers {
		var totalWidth float32
		for idx, id := range layer {
			if node, ok := nodes[id]; ok {
				totalWidth += node.Width
				if idx > 0 {
					totalWidth += nodeSpacing
				}
			}
		}

		posX := -totalWidth/2 + layoutBoundaryPad

		for _, id := range layer {
			node, ok := nodes[id]
			if !ok {
				continue
			}
			posY := rankY[rankIdx] + layoutBoundaryPad
			if reverse {
				posY = cumY - rankY[rankIdx] + layoutBoundaryPad
			}
			node.X = posX + node.Width/2
			node.Y = posY
			posX += node.Width + nodeSpacing
		}
	}
}
