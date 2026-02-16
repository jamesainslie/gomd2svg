package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeClassLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()

	// Size class nodes with UML compartments.
	nodes, compartments := sizeClassNodes(graph, measurer, th, cfg)

	// Reuse Sugiyama pipeline.
	result := runSugiyama(graph, nodes, cfg)

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  nodes,
		Edges:  result.Edges,
		Width:  result.Width,
		Height: result.Height,
		Diagram: ClassData{
			Compartments: compartments,
			Members:      graph.Members,
			Annotations:  graph.Annotations,
		},
	}
}

func sizeClassNodes(graph *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) (map[string]*NodeLayout, map[string]ClassCompartment) {
	nodes := make(map[string]*NodeLayout, len(graph.Nodes))
	compartments := make(map[string]ClassCompartment)

	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical
	lineH := th.FontSize * cfg.LabelLineHeight
	memberFontSize := cfg.Class.MemberFontSize
	compartmentPadY := cfg.Class.CompartmentPadY

	memberLineH := memberFontSize * cfg.LabelLineHeight

	for id, node := range graph.Nodes {
		members := graph.Members[id]

		if members == nil || (len(members.Attributes) == 0 && len(members.Methods) == 0) {
			// Simple node — no compartments, just measure label.
			nl := sizeNode(node, measurer, th, cfg)
			nodes[id] = nl
			continue
		}

		// Measure header (class name).
		headerW := measurer.Width(node.Label, th.FontSize, th.FontFamily)
		headerH := lineH + padV

		// Measure annotation if present.
		if ann, ok := graph.Annotations[id]; ok {
			annW := measurer.Width("<<"+ann+">>", th.FontSize, th.FontFamily)
			if annW > headerW {
				headerW = annW
			}
			headerH += lineH
		}

		// Measure attributes.
		var attrH float32
		maxW := headerW
		for _, attr := range members.Attributes {
			text := attr.Visibility.Symbol() + attr.Type + " " + attr.Name
			w := measurer.Width(text, memberFontSize, th.FontFamily)
			if w > maxW {
				maxW = w
			}
			attrH += memberLineH
		}
		if len(members.Attributes) > 0 {
			attrH += compartmentPadY // section padding
		}

		// Measure methods.
		var methH float32
		for _, meth := range members.Methods {
			text := meth.Visibility.Symbol() + meth.Name + "(" + meth.Params + ")"
			if meth.Type != "" {
				text += " : " + meth.Type
			}
			w := measurer.Width(text, memberFontSize, th.FontFamily)
			if w > maxW {
				maxW = w
			}
			methH += memberLineH
		}
		if len(members.Methods) > 0 {
			methH += compartmentPadY
		}

		totalW := maxW + 2*padH
		totalH := headerH + attrH + methH + padV

		compartments[id] = ClassCompartment{
			HeaderHeight:    headerH,
			AttributeHeight: attrH,
			MethodHeight:    methH,
		}

		nodes[id] = &NodeLayout{
			ID:     id,
			Label:  TextBlock{Lines: []string{node.Label}, Width: maxW, Height: headerH, FontSize: th.FontSize},
			Shape:  ir.Rectangle,
			Width:  totalW,
			Height: totalH,
		}
	}

	return nodes, compartments
}
