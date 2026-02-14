package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeClassLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()

	// Size class nodes with UML compartments.
	nodes, compartments := sizeClassNodes(g, measurer, th, cfg)

	// Reuse Sugiyama pipeline.
	r := runSugiyama(g, nodes, cfg)

	return &Layout{
		Kind:   g.Kind,
		Nodes:  nodes,
		Edges:  r.Edges,
		Width:  r.Width,
		Height: r.Height,
		Diagram: ClassData{
			Compartments: compartments,
			Members:      g.Members,
			Annotations:  g.Annotations,
		},
	}
}

func sizeClassNodes(g *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) (map[string]*NodeLayout, map[string]ClassCompartment) {
	nodes := make(map[string]*NodeLayout, len(g.Nodes))
	compartments := make(map[string]ClassCompartment)

	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical
	lineH := th.FontSize * cfg.LabelLineHeight
	memberFontSize := cfg.Class.MemberFontSize
	compartmentPadY := cfg.Class.CompartmentPadY

	memberLineH := memberFontSize * cfg.LabelLineHeight

	for id, node := range g.Nodes {
		members := g.Members[id]

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
		if ann, ok := g.Annotations[id]; ok {
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
