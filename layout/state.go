package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

// State layout constants.
const (
	stateStartEndPadding = 4
	stateChoicePadding   = 24
)

func computeStateLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	innerLayouts := make(map[string]*Layout)

	// Size state nodes. Special handling for __start__/__end__/fork/choice.
	nodes := sizeStateNodes(graph, measurer, th, cfg, innerLayouts)

	result := runSugiyama(graph, nodes, cfg)

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  nodes,
		Edges:  result.Edges,
		Width:  result.Width,
		Height: result.Height,
		Diagram: StateData{
			InnerLayouts:    innerLayouts,
			Descriptions:    graph.StateDescriptions,
			Annotations:     graph.StateAnnotations,
			CompositeStates: graph.CompositeStates,
		},
	}
}

func sizeStateNodes(graph *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout, innerLayouts map[string]*Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(graph.Nodes))

	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical
	lineH := th.FontSize * cfg.LabelLineHeight

	for id, node := range graph.Nodes {
		// Special nodes: __start__ and __end__ are small circles.
		if id == "__start__" || id == "__end__" {
			size := cfg.State.StartEndRadius*2 + stateStartEndPadding
			shape := ir.Circle
			if id == "__end__" {
				shape = ir.DoubleCircle
			}
			nodes[id] = &NodeLayout{
				ID:     id,
				Label:  TextBlock{FontSize: th.FontSize},
				Shape:  shape,
				Width:  size,
				Height: size,
			}
			continue
		}

		// Fork/join annotations: narrow bar shape.
		if ann, ok := graph.StateAnnotations[id]; ok {
			switch ann {
			case ir.StateFork, ir.StateJoin:
				nodes[id] = &NodeLayout{
					ID:     id,
					Label:  TextBlock{FontSize: th.FontSize},
					Shape:  ir.ForkJoin,
					Width:  cfg.State.ForkBarWidth,
					Height: cfg.State.ForkBarHeight,
				}
				continue
			case ir.StateChoice:
				choiceSize := cfg.State.StartEndRadius*2 + stateChoicePadding
				nodes[id] = &NodeLayout{
					ID:     id,
					Label:  TextBlock{FontSize: th.FontSize},
					Shape:  ir.Diamond,
					Width:  choiceSize,
					Height: choiceSize,
				}
				continue
			}
		}

		// Composite states: recursively layout inner graph.
		if cs, ok := graph.CompositeStates[id]; ok && cs.Inner != nil {
			innerLayout := computeStateLayout(cs.Inner, th, cfg)
			innerLayouts[id] = innerLayout

			// Size the composite node to contain its inner layout.
			labelW := measurer.Width(node.Label, th.FontSize, th.FontFamily)
			labelH := lineH + padV

			innerW := innerLayout.Width
			innerH := innerLayout.Height

			totalW := innerW + 2*cfg.State.CompositePadding
			if labelW+2*cfg.State.CompositePadding > totalW {
				totalW = labelW + 2*cfg.State.CompositePadding
			}
			totalH := labelH + innerH + cfg.State.CompositePadding

			nodes[id] = &NodeLayout{
				ID:     id,
				Label:  TextBlock{Lines: []string{node.Label}, Width: labelW, Height: labelH, FontSize: th.FontSize},
				Shape:  ir.Rectangle,
				Width:  totalW,
				Height: totalH,
			}
			continue
		}

		// Regular state node with optional description.
		nodeLayout := sizeNode(node, measurer, th, cfg)

		// Add description height if present.
		if desc, ok := graph.StateDescriptions[id]; ok {
			descW := measurer.Width(desc, th.FontSize, th.FontFamily)
			nodeLayout.Height += lineH + padV
			if descW+2*padH > nodeLayout.Width {
				nodeLayout.Width = descW + 2*padH
			}
		}

		// Apply rounded corners style for state nodes.
		nodeLayout.Shape = ir.RoundRect
		nodes[id] = nodeLayout
	}

	return nodes
}
