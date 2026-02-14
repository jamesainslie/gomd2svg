package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeStateLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	innerLayouts := make(map[string]*Layout)

	// Size state nodes. Special handling for __start__/__end__/fork/choice.
	nodes := sizeStateNodes(g, measurer, th, cfg, innerLayouts)

	r := runSugiyama(g, nodes, cfg)

	return &Layout{
		Kind:   g.Kind,
		Nodes:  nodes,
		Edges:  r.Edges,
		Width:  r.Width,
		Height: r.Height,
		Diagram: StateData{
			InnerLayouts:    innerLayouts,
			Descriptions:    g.StateDescriptions,
			Annotations:     g.StateAnnotations,
			CompositeStates: g.CompositeStates,
		},
	}
}

func sizeStateNodes(g *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout, innerLayouts map[string]*Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(g.Nodes))

	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical
	lineH := th.FontSize * cfg.LabelLineHeight

	for id, node := range g.Nodes {
		// Special nodes: __start__ and __end__ are small circles.
		if id == "__start__" || id == "__end__" {
			size := cfg.State.StartEndRadius*2 + 4
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
		if ann, ok := g.StateAnnotations[id]; ok {
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
				choiceSize := cfg.State.StartEndRadius*2 + 24
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
		if cs, ok := g.CompositeStates[id]; ok && cs.Inner != nil {
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
		nl := sizeNode(node, measurer, th, cfg)

		// Add description height if present.
		if desc, ok := g.StateDescriptions[id]; ok {
			descW := measurer.Width(desc, th.FontSize, th.FontFamily)
			nl.Height += lineH + padV
			if descW+2*padH > nl.Width {
				nl.Width = descW + 2*padH
			}
		}

		// Apply rounded corners style for state nodes.
		nl.Shape = ir.RoundRect
		nodes[id] = nl
	}

	return nodes
}
