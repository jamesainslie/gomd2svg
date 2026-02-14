package layout

import "github.com/jamesainslie/mermaid-go/ir"

// routeEdges computes simple L-shaped polyline routes for all edges.
// For LR direction: start at right side of source, end at left side of target.
// For TD direction: start at bottom of source, end at top of target.
func routeEdges(
	edges []*ir.Edge,
	nodes map[string]*NodeLayout,
	direction ir.Direction,
) []*EdgeLayout {
	result := make([]*EdgeLayout, 0, len(edges))

	for _, e := range edges {
		src, srcOK := nodes[e.From]
		dst, dstOK := nodes[e.To]
		if !srcOK || !dstOK {
			continue
		}

		var points [][2]float32
		var labelAnchor [2]float32

		switch direction {
		case ir.LeftRight, ir.RightLeft:
			points, labelAnchor = routeLR(src, dst)
		default:
			points, labelAnchor = routeTD(src, dst)
		}

		var label *TextBlock
		if e.Label != nil {
			label = &TextBlock{
				Lines:    []string{*e.Label},
				FontSize: src.Label.FontSize,
			}
		}

		result = append(result, &EdgeLayout{
			From:           e.From,
			To:             e.To,
			Label:          label,
			Points:         points,
			LabelAnchor:    labelAnchor,
			Style:          e.Style,
			ArrowStart:     e.ArrowStart,
			ArrowEnd:       e.ArrowEnd,
			ArrowStartKind: e.ArrowStartKind,
			ArrowEndKind:   e.ArrowEndKind,
		})
	}

	return result
}

// routeLR creates an L-shaped route for a left-right edge.
// Start from the right side of source, end at the left side of target.
// Route: start -> horizontal to midX -> vertical to target Y -> end.
func routeLR(src, dst *NodeLayout) ([][2]float32, [2]float32) {
	startX := src.X + src.Width/2
	startY := src.Y
	endX := dst.X - dst.Width/2
	endY := dst.Y

	midX := (startX + endX) / 2

	points := [][2]float32{
		{startX, startY},
		{midX, startY},
		{midX, endY},
		{endX, endY},
	}

	labelAnchor := [2]float32{midX, (startY + endY) / 2}

	return points, labelAnchor
}

// routeTD creates an L-shaped route for a top-down edge.
// Start from the bottom of source, end at the top of target.
// Route: start -> vertical to midY -> horizontal to target X -> end.
func routeTD(src, dst *NodeLayout) ([][2]float32, [2]float32) {
	startX := src.X
	startY := src.Y + src.Height/2
	endX := dst.X
	endY := dst.Y - dst.Height/2

	midY := (startY + endY) / 2

	points := [][2]float32{
		{startX, startY},
		{startX, midY},
		{endX, midY},
		{endX, endY},
	}

	labelAnchor := [2]float32{(startX + endX) / 2, midY}

	return points, labelAnchor
}
