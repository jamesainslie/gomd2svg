package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

// c4SmallFontRatio is the ratio of small text size to the base font size.
const c4SmallFontRatio = 0.85

func computeC4Layout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	nodes := sizeC4Nodes(graph, measurer, th, cfg)

	result := runSugiyama(graph, nodes, cfg)

	elemMap := make(map[string]*ir.C4Element)
	for _, elem := range graph.C4Elements {
		elemMap[elem.ID] = elem
	}

	var boundaryLayouts []*C4BoundaryLayout
	for _, boundary := range graph.C4Boundaries {
		bl := computeC4BoundaryRect(boundary, nodes, cfg)
		if bl != nil {
			boundaryLayouts = append(boundaryLayouts, bl)
		}
	}

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  nodes,
		Edges:  result.Edges,
		Width:  result.Width,
		Height: result.Height,
		Diagram: C4Data{
			Elements:   elemMap,
			Boundaries: boundaryLayouts,
			SubKind:    graph.C4SubKind,
		},
	}
}

func sizeC4Nodes(graph *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(graph.Nodes))
	lineH := th.FontSize * cfg.LabelLineHeight
	smallFontSize := th.FontSize * c4SmallFontRatio
	smallLineH := smallFontSize * cfg.LabelLineHeight
	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical

	elemMap := make(map[string]*ir.C4Element)
	for _, elem := range graph.C4Elements {
		elemMap[elem.ID] = elem
	}

	for id, node := range graph.Nodes {
		elem := elemMap[id]
		var maxW, totalH float32

		labelW := measurer.Width(node.Label, th.FontSize, th.FontFamily)
		if labelW > maxW {
			maxW = labelW
		}
		totalH += lineH

		if elem != nil {
			if elem.Technology != "" {
				techText := "[" + elem.Technology + "]"
				techW := measurer.Width(techText, smallFontSize, th.FontFamily)
				if techW > maxW {
					maxW = techW
				}
				totalH += smallLineH
			}
			if elem.Description != "" {
				descW := measurer.Width(elem.Description, smallFontSize, th.FontFamily)
				if descW > maxW {
					maxW = descW
				}
				totalH += smallLineH
			}

			if elem.Type.IsPerson() {
				if maxW+2*padH < cfg.C4.PersonWidth {
					maxW = cfg.C4.PersonWidth - 2*padH
				}
				if totalH+2*padV < cfg.C4.PersonHeight {
					totalH = cfg.C4.PersonHeight - 2*padV
				}
			} else {
				if maxW+2*padH < cfg.C4.SystemWidth {
					maxW = cfg.C4.SystemWidth - 2*padH
				}
				if totalH+2*padV < cfg.C4.SystemHeight {
					totalH = cfg.C4.SystemHeight - 2*padV
				}
			}
		}

		nodes[id] = &NodeLayout{
			ID:     id,
			Label:  TextBlock{Lines: []string{node.Label}, Width: maxW, Height: lineH, FontSize: th.FontSize},
			Shape:  ir.Rectangle,
			Width:  maxW + 2*padH,
			Height: totalH + 2*padV,
		}
	}

	return nodes
}

func computeC4BoundaryRect(boundary *ir.C4Boundary, nodes map[string]*NodeLayout, cfg *config.Layout) *C4BoundaryLayout {
	const boundaryLabelHeight float32 = 20

	if len(boundary.Children) == 0 {
		return nil
	}
	pad := cfg.C4.BoundaryPadding

	var minX, minY, maxX, maxY float32
	first := true
	for _, childID := range boundary.Children {
		childNode, ok := nodes[childID]
		if !ok {
			continue
		}
		left := childNode.X - childNode.Width/2
		top := childNode.Y - childNode.Height/2
		right := childNode.X + childNode.Width/2
		bottom := childNode.Y + childNode.Height/2
		if first {
			minX, minY, maxX, maxY = left, top, right, bottom
			first = false
		} else {
			if left < minX {
				minX = left
			}
			if top < minY {
				minY = top
			}
			if right > maxX {
				maxX = right
			}
			if bottom > maxY {
				maxY = bottom
			}
		}
	}
	if first {
		return nil
	}

	return &C4BoundaryLayout{
		ID:     boundary.ID,
		Label:  boundary.Label,
		Type:   boundary.Type,
		X:      minX - pad,
		Y:      minY - pad - boundaryLabelHeight,
		Width:  (maxX - minX) + 2*pad,
		Height: (maxY - minY) + 2*pad + boundaryLabelHeight,
	}
}
