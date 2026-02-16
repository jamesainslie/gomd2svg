package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeBlockLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	nodes := sizeBlockNodes(g, measurer, th, cfg)

	blockInfos := make(map[string]BlockInfo)
	for _, b := range g.Blocks {
		blockInfos[b.ID] = BlockInfo{Span: b.Width, HasChildren: len(b.Children) > 0}
	}

	// Decide layout strategy
	if g.BlockColumns > 0 {
		return blockGridLayout(g, nodes, blockInfos, cfg)
	}
	if len(g.Edges) > 0 {
		r := runSugiyama(g, nodes, cfg)
		return &Layout{
			Kind:    g.Kind,
			Nodes:   nodes,
			Edges:   r.Edges,
			Width:   r.Width,
			Height:  r.Height,
			Diagram: BlockData{Columns: 0, BlockInfos: blockInfos},
		}
	}
	return blockGridLayout(g, nodes, blockInfos, cfg)
}

func blockGridLayout(g *ir.Graph, nodes map[string]*NodeLayout, blockInfos map[string]BlockInfo, cfg *config.Layout) *Layout {
	cols := g.BlockColumns
	if cols <= 0 {
		cols = 1
	}
	colGap := cfg.Block.ColumnGap
	rowGap := cfg.Block.RowGap
	padX := cfg.Block.PaddingX
	padY := cfg.Block.PaddingY

	var maxCellW, maxCellH float32
	for _, n := range nodes {
		if n.Width > maxCellW {
			maxCellW = n.Width
		}
		if n.Height > maxCellH {
			maxCellH = n.Height
		}
	}

	col := 0
	row := 0
	for _, blk := range g.Blocks {
		n, ok := nodes[blk.ID]
		if !ok {
			continue
		}

		span := blk.Width
		if span <= 0 {
			span = 1
		}
		if col+span > cols {
			col = 0
			row++
		}

		cellW := maxCellW*float32(span) + colGap*float32(span-1)
		n.Width = cellW
		n.X = padX + float32(col)*(maxCellW+colGap) + cellW/2
		n.Y = padY + float32(row)*(maxCellH+rowGap) + maxCellH/2

		col += span
		if col >= cols {
			col = 0
			row++
		}
	}

	var edges []*EdgeLayout
	for _, e := range g.Edges {
		src := nodes[e.From]
		dst := nodes[e.To]
		if src == nil || dst == nil {
			continue
		}
		edges = append(edges, &EdgeLayout{
			From:     e.From,
			To:       e.To,
			Points:   [][2]float32{{src.X, src.Y}, {dst.X, dst.Y}},
			ArrowEnd: e.ArrowEnd,
		})
	}

	totalW := padX*2 + float32(cols)*maxCellW + float32(cols-1)*colGap
	totalRows := row + 1
	if col == 0 && row > 0 {
		totalRows = row
	}
	totalH := padY*2 + float32(totalRows)*maxCellH + float32(totalRows-1)*rowGap

	return &Layout{
		Kind:    g.Kind,
		Nodes:   nodes,
		Edges:   edges,
		Width:   totalW,
		Height:  totalH,
		Diagram: BlockData{Columns: cols, BlockInfos: blockInfos},
	}
}

func sizeBlockNodes(g *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(g.Nodes))
	for id, node := range g.Nodes {
		nl := sizeNode(node, measurer, th, cfg)
		nodes[id] = nl
	}
	return nodes
}
