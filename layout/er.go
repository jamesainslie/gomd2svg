package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeERLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()

	nodes, entityDims := sizeERNodes(g, measurer, th, cfg)

	r := runSugiyama(g, nodes, cfg)

	return &Layout{
		Kind:    g.Kind,
		Nodes:   nodes,
		Edges:   r.Edges,
		Width:   r.Width,
		Height:  r.Height,
		Diagram: ERData{EntityDims: entityDims, Entities: g.Entities},
	}
}

func sizeERNodes(g *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) (map[string]*NodeLayout, map[string]EntityDimensions) {
	nodes := make(map[string]*NodeLayout, len(g.Nodes))
	dims := make(map[string]EntityDimensions)

	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical
	lineH := th.FontSize * cfg.LabelLineHeight
	colPad := cfg.ER.ColumnPadding
	rowH := cfg.ER.AttributeRowHeight

	for id, node := range g.Nodes {
		entity := g.Entities[id]
		if entity == nil {
			nl := sizeNode(node, measurer, th, cfg)
			nodes[id] = nl
			continue
		}

		// Display name for header.
		displayName := entity.DisplayName()
		headerW := measurer.Width(displayName, th.FontSize, th.FontFamily)
		headerH := lineH + padV

		// Measure attribute columns.
		var maxTypeW, maxNameW, maxKeyW float32
		for _, attr := range entity.Attributes {
			tw := measurer.Width(attr.Type, th.FontSize, th.FontFamily)
			if tw > maxTypeW {
				maxTypeW = tw
			}
			nw := measurer.Width(attr.Name, th.FontSize, th.FontFamily)
			if nw > maxNameW {
				maxNameW = nw
			}
			var keyStr string
			for i, k := range attr.Keys {
				if i > 0 {
					keyStr += ","
				}
				keyStr += k.String()
			}
			kw := measurer.Width(keyStr, th.FontSize, th.FontFamily)
			if kw > maxKeyW {
				maxKeyW = kw
			}
		}

		rowCount := len(entity.Attributes)
		bodyH := rowH * float32(rowCount)
		if rowCount > 0 {
			bodyH += padV
		}

		rowW := maxTypeW + maxNameW + maxKeyW + 3*colPad // 3 columns with padding
		totalW := rowW + 2*padH
		if headerW+2*padH > totalW {
			totalW = headerW + 2*padH
		}
		totalH := headerH + bodyH + padV

		dims[id] = EntityDimensions{
			TypeColWidth: maxTypeW,
			NameColWidth: maxNameW,
			KeyColWidth:  maxKeyW,
			HeaderHeight: headerH,
			RowCount:     rowCount,
		}

		nodes[id] = &NodeLayout{
			ID:     id,
			Label:  TextBlock{Lines: []string{displayName}, Width: headerW, Height: headerH, FontSize: th.FontSize},
			Shape:  ir.Rectangle,
			Width:  totalW,
			Height: totalH,
		}
	}

	return nodes, dims
}
