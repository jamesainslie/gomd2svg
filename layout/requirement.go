package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeRequirementLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	nodes := sizeRequirementNodes(g, measurer, th, cfg)

	r := runSugiyama(g, nodes, cfg)

	reqMap := make(map[string]*ir.RequirementDef)
	for _, req := range g.Requirements {
		reqMap[req.Name] = req
	}
	elemMap := make(map[string]*ir.ElementDef)
	for _, elem := range g.ReqElements {
		elemMap[elem.Name] = elem
	}
	nodeKinds := make(map[string]string)
	for _, req := range g.Requirements {
		nodeKinds[req.Name] = "requirement"
	}
	for _, elem := range g.ReqElements {
		nodeKinds[elem.Name] = "element"
	}

	return &Layout{
		Kind:   g.Kind,
		Nodes:  nodes,
		Edges:  r.Edges,
		Width:  r.Width,
		Height: r.Height,
		Diagram: RequirementData{
			Requirements: reqMap,
			Elements:     elemMap,
			NodeKinds:    nodeKinds,
		},
	}
}

func sizeRequirementNodes(g *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(g.Nodes))
	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical
	lineH := th.FontSize * cfg.LabelLineHeight
	metaFontSize := cfg.Requirement.MetadataFontSize
	metaLineH := metaFontSize * cfg.LabelLineHeight
	minW := cfg.Requirement.NodeMinWidth

	reqMap := make(map[string]*ir.RequirementDef)
	for _, req := range g.Requirements {
		reqMap[req.Name] = req
	}
	elemMap := make(map[string]*ir.ElementDef)
	for _, elem := range g.ReqElements {
		elemMap[elem.Name] = elem
	}

	for id, node := range g.Nodes {
		var maxW float32
		var totalH float32

		// Stereotype line
		stereotypeText := ""
		if req, ok := reqMap[id]; ok {
			stereotypeText = "\u00AB" + req.Type.Stereotype() + "\u00BB"
		} else {
			stereotypeText = "\u00ABelement\u00BB"
		}
		stW := measurer.Width(stereotypeText, metaFontSize, th.FontFamily)
		if stW > maxW {
			maxW = stW
		}
		totalH += lineH

		// Name line
		nameW := measurer.Width(node.Label, th.FontSize, th.FontFamily)
		if nameW > maxW {
			maxW = nameW
		}
		totalH += lineH

		// Metadata lines
		if req, ok := reqMap[id]; ok {
			lines := 0
			if req.ID != "" {
				w := measurer.Width("Id: "+req.ID, metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			if req.Text != "" {
				w := measurer.Width("Text: "+req.Text, metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			if req.Risk != ir.RiskNone {
				w := measurer.Width("Risk: "+req.Risk.String(), metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			if req.VerifyMethod != ir.VerifyNone {
				w := measurer.Width("Verify: "+req.VerifyMethod.String(), metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			totalH += metaLineH * float32(lines)
		} else if elem, ok := elemMap[id]; ok {
			lines := 0
			if elem.Type != "" {
				w := measurer.Width("Type: "+elem.Type, metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			if elem.DocRef != "" {
				w := measurer.Width("Doc: "+elem.DocRef, metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			totalH += metaLineH * float32(lines)
		}

		w := maxW + 2*padH
		if w < minW {
			w = minW
		}
		h := totalH + 2*padV

		nodes[id] = &NodeLayout{
			ID:     id,
			Label:  TextBlock{Lines: []string{node.Label}, Width: maxW, Height: lineH, FontSize: th.FontSize},
			Shape:  ir.Rectangle,
			Width:  w,
			Height: h,
		}
	}

	return nodes
}
