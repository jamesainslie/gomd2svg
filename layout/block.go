package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeBlockLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	nodes := sizeBlockNodes(graph, measurer, th, cfg)

	blockInfos := make(map[string]BlockInfo)
	for _, blk := range graph.Blocks {
		blockInfos[blk.ID] = BlockInfo{Span: blk.Width, HasChildren: len(blk.Children) > 0}
	}

	// Decide layout strategy
	if graph.BlockColumns > 0 {
		return blockGridLayout(graph, nodes, blockInfos, cfg)
	}
	if len(graph.Edges) > 0 {
		result := runSugiyama(graph, nodes, cfg)
		return &Layout{
			Kind:    graph.Kind,
			Nodes:   nodes,
			Edges:   result.Edges,
			Width:   result.Width,
			Height:  result.Height,
			Diagram: BlockData{Columns: 0, BlockInfos: blockInfos},
		}
	}
	return blockGridLayout(graph, nodes, blockInfos, cfg)
}

func blockGridLayout(graph *ir.Graph, nodes map[string]*NodeLayout, blockInfos map[string]BlockInfo, cfg *config.Layout) *Layout {
	cols := graph.BlockColumns
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
	for _, blk := range graph.Blocks {
		node, ok := nodes[blk.ID]
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
		node.Width = cellW
		node.X = padX + float32(col)*(maxCellW+colGap) + cellW/2
		node.Y = padY + float32(row)*(maxCellH+rowGap) + maxCellH/2

		col += span
		if col >= cols {
			col = 0
			row++
		}
	}

	var edges []*EdgeLayout
	for _, edge := range graph.Edges {
		src := nodes[edge.From]
		dst := nodes[edge.To]
		if src == nil || dst == nil {
			continue
		}
		edges = append(edges, &EdgeLayout{
			From:     edge.From,
			To:       edge.To,
			Points:   [][2]float32{{src.X, src.Y}, {dst.X, dst.Y}},
			ArrowEnd: edge.ArrowEnd,
		})
	}

	totalW := padX*2 + float32(cols)*maxCellW + float32(cols-1)*colGap
	totalRows := row + 1
	if col == 0 && row > 0 {
		totalRows = row
	}
	totalH := padY*2 + float32(totalRows)*maxCellH + float32(totalRows-1)*rowGap

	return &Layout{
		Kind:    graph.Kind,
		Nodes:   nodes,
		Edges:   edges,
		Width:   totalW,
		Height:  totalH,
		Diagram: BlockData{Columns: cols, BlockInfos: blockInfos},
	}
}

func sizeBlockNodes(graph *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(graph.Nodes))
	for id, node := range graph.Nodes {
		nl := sizeNode(node, measurer, th, cfg)
		nodes[id] = nl
	}
	return nodes
}
