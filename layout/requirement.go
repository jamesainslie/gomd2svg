package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeRequirementLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	nodes := sizeRequirementNodes(graph, measurer, th, cfg)

	result := runSugiyama(graph, nodes, cfg)

	reqMap := make(map[string]*ir.RequirementDef)
	for _, req := range graph.Requirements {
		reqMap[req.Name] = req
	}
	elemMap := make(map[string]*ir.ElementDef)
	for _, elem := range graph.ReqElements {
		elemMap[elem.Name] = elem
	}
	nodeKinds := make(map[string]string)
	for _, req := range graph.Requirements {
		nodeKinds[req.Name] = "requirement"
	}
	for _, elem := range graph.ReqElements {
		nodeKinds[elem.Name] = "element"
	}

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  nodes,
		Edges:  result.Edges,
		Width:  result.Width,
		Height: result.Height,
		Diagram: RequirementData{
			Requirements: reqMap,
			Elements:     elemMap,
			NodeKinds:    nodeKinds,
		},
	}
}

func sizeRequirementNodes(graph *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(graph.Nodes))
	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical
	lineH := th.FontSize * cfg.LabelLineHeight
	metaFontSize := cfg.Requirement.MetadataFontSize
	metaLineH := metaFontSize * cfg.LabelLineHeight
	minW := cfg.Requirement.NodeMinWidth

	reqMap := make(map[string]*ir.RequirementDef)
	for _, req := range graph.Requirements {
		reqMap[req.Name] = req
	}
	elemMap := make(map[string]*ir.ElementDef)
	for _, elem := range graph.ReqElements {
		elemMap[elem.Name] = elem
	}

	for id, node := range graph.Nodes {
		var maxW float32
		var totalH float32

		// Stereotype line
		var stereotypeText string
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
			metaW, metaLines := reqMetadataSize(req, measurer, metaFontSize, th.FontFamily)
			if metaW > maxW {
				maxW = metaW
			}
			totalH += metaLineH * float32(metaLines)
		} else if elem, ok := elemMap[id]; ok {
			metaW, metaLines := elemMetadataSize(elem, measurer, metaFontSize, th.FontFamily)
			if metaW > maxW {
				maxW = metaW
			}
			totalH += metaLineH * float32(metaLines)
		}

		nodeW := maxW + 2*padH
		if nodeW < minW {
			nodeW = minW
		}
		nodeH := totalH + 2*padV

		nodes[id] = &NodeLayout{
			ID:     id,
			Label:  TextBlock{Lines: []string{node.Label}, Width: maxW, Height: lineH, FontSize: th.FontSize},
			Shape:  ir.Rectangle,
			Width:  nodeW,
			Height: nodeH,
		}
	}

	return nodes
}

// reqMetadataSize measures the width and line count of requirement metadata fields.
func reqMetadataSize(req *ir.RequirementDef, measurer *textmetrics.Measurer, fontSize float32, fontFamily string) (float32, int) {
	var maxW float32
	lines := 0
	if req.ID != "" {
		lineW := measurer.Width("Id: "+req.ID, fontSize, fontFamily)
		if lineW > maxW {
			maxW = lineW
		}
		lines++
	}
	if req.Text != "" {
		lineW := measurer.Width("Text: "+req.Text, fontSize, fontFamily)
		if lineW > maxW {
			maxW = lineW
		}
		lines++
	}
	if req.Risk != ir.RiskNone {
		lineW := measurer.Width("Risk: "+req.Risk.String(), fontSize, fontFamily)
		if lineW > maxW {
			maxW = lineW
		}
		lines++
	}
	if req.VerifyMethod != ir.VerifyNone {
		lineW := measurer.Width("Verify: "+req.VerifyMethod.String(), fontSize, fontFamily)
		if lineW > maxW {
			maxW = lineW
		}
		lines++
	}
	return maxW, lines
}

// elemMetadataSize measures the width and line count of element metadata fields.
func elemMetadataSize(elem *ir.ElementDef, measurer *textmetrics.Measurer, fontSize float32, fontFamily string) (float32, int) {
	var maxW float32
	lines := 0
	if elem.Type != "" {
		lineW := measurer.Width("Type: "+elem.Type, fontSize, fontFamily)
		if lineW > maxW {
			maxW = lineW
		}
		lines++
	}
	if elem.DocRef != "" {
		lineW := measurer.Width("Doc: "+elem.DocRef, fontSize, fontFamily)
		if lineW > maxW {
			maxW = lineW
		}
		lines++
	}
	return maxW, lines
}
