package layout

import (
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeSankeyLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	if len(g.SankeyLinks) == 0 {
		return &Layout{
			Kind:    g.Kind,
			Nodes:   map[string]*NodeLayout{},
			Width:   cfg.Sankey.ChartWidth + cfg.Sankey.PaddingX*2,
			Height:  cfg.Sankey.ChartHeight + cfg.Sankey.PaddingY*2,
			Diagram: SankeyData{},
		}
	}

	// Step 1: Collect unique node names preserving order of first appearance.
	nodeNames, nodeIndex := sankeyCollectNodes(g.SankeyLinks)

	// Step 2: Assign columns via longest path from sources.
	columns := sankeyAssignColumns(nodeNames, nodeIndex, g.SankeyLinks)
	maxCol := 0
	for _, c := range columns {
		if c > maxCol {
			maxCol = c
		}
	}

	// Step 3: Compute total flow per node (max of inflow vs outflow).
	totalFlow := sankeyComputeFlow(nodeNames, nodeIndex, g.SankeyLinks)

	// Step 4: Position nodes in columns.
	chartW := cfg.Sankey.ChartWidth
	chartH := cfg.Sankey.ChartHeight
	padX := cfg.Sankey.PaddingX
	padY := cfg.Sankey.PaddingY
	nodeWidth := cfg.Sankey.NodeWidth
	nodePad := cfg.Sankey.NodePadding

	// Group nodes by column.
	colNodes := make([][]int, maxCol+1)
	for i, c := range columns {
		colNodes[c] = append(colNodes[c], i)
	}

	// Compute horizontal spacing between columns.
	colWidth := float32(0)
	if maxCol > 0 {
		colWidth = (chartW - nodeWidth) / float32(maxCol)
	}

	// Position nodes within each column.
	nodes := make([]SankeyNodeLayout, len(nodeNames))
	for c, col := range colNodes {
		// Sort nodes in column by their original appearance order.
		sort.Slice(col, func(i, j int) bool { return col[i] < col[j] })

		x := padX + float32(c)*colWidth

		// Compute total flow for this column.
		colTotal := float64(0)
		for _, ni := range col {
			colTotal += totalFlow[ni]
		}

		// Scale heights so nodes fit in chartH minus inter-node padding.
		availH := chartH - nodePad*float32(len(col)-1)
		scale := float64(availH) / colTotal
		if colTotal == 0 {
			scale = 0
		}

		y := padY
		for _, ni := range col {
			h := float32(totalFlow[ni] * scale)
			if h < 2 {
				h = 2
			}
			nodes[ni] = SankeyNodeLayout{
				Label:      nodeNames[ni],
				X:          x,
				Y:          y,
				Width:      nodeWidth,
				Height:     h,
				ColorIndex: ni,
			}
			y += h + nodePad
		}
	}

	// Step 5: Compute link positions with running Y offsets for stacking.
	sourceOffsets := make([]float32, len(nodeNames))
	targetOffsets := make([]float32, len(nodeNames))

	links := make([]SankeyLinkLayout, len(g.SankeyLinks))
	for i, link := range g.SankeyLinks {
		si := nodeIndex[link.Source]
		ti := nodeIndex[link.Target]

		// Link width proportional to value relative to source node total flow.
		linkH := float32(0)
		if totalFlow[si] > 0 {
			linkH = nodes[si].Height * float32(link.Value/totalFlow[si])
		}
		if linkH < 1 {
			linkH = 1
		}

		links[i] = SankeyLinkLayout{
			SourceIdx: si,
			TargetIdx: ti,
			Value:     link.Value,
			SourceY:   nodes[si].Y + sourceOffsets[si],
			TargetY:   nodes[ti].Y + targetOffsets[ti],
			Width:     linkH,
		}
		sourceOffsets[si] += linkH
		targetOffsets[ti] += linkH
	}

	totalW := chartW + padX*2
	totalH := chartH + padY*2

	return &Layout{
		Kind:    g.Kind,
		Nodes:   map[string]*NodeLayout{},
		Width:   totalW,
		Height:  totalH,
		Diagram: SankeyData{Nodes: nodes, Links: links},
	}
}

// sankeyCollectNodes extracts unique node names from links, preserving the
// order of first appearance. It returns the name list and name-to-index map.
func sankeyCollectNodes(links []*ir.SankeyLink) ([]string, map[string]int) {
	index := make(map[string]int)
	var names []string
	for _, link := range links {
		if _, ok := index[link.Source]; !ok {
			index[link.Source] = len(names)
			names = append(names, link.Source)
		}
		if _, ok := index[link.Target]; !ok {
			index[link.Target] = len(names)
			names = append(names, link.Target)
		}
	}
	return names, index
}

// sankeyAssignColumns assigns each node to a column using the longest-path
// algorithm: source nodes (no incoming links) get column 0, others get
// max(source column) + 1.
func sankeyAssignColumns(names []string, index map[string]int, links []*ir.SankeyLink) []int {
	n := len(names)
	columns := make([]int, n)

	// Build incoming edge lists.
	incoming := make([][]int, n)
	for _, link := range links {
		si := index[link.Source]
		ti := index[link.Target]
		incoming[ti] = append(incoming[ti], si)
	}

	// Iterative longest-path relaxation.
	changed := true
	for changed {
		changed = false
		for i := 0; i < n; i++ {
			for _, src := range incoming[i] {
				if columns[src]+1 > columns[i] {
					columns[i] = columns[src] + 1
					changed = true
				}
			}
		}
	}

	return columns
}

// sankeyComputeFlow computes the total flow through each node as the max of
// its inflow and outflow sums.
func sankeyComputeFlow(names []string, index map[string]int, links []*ir.SankeyLink) []float64 {
	n := len(names)
	inflow := make([]float64, n)
	outflow := make([]float64, n)
	for _, link := range links {
		si := index[link.Source]
		ti := index[link.Target]
		outflow[si] += link.Value
		inflow[ti] += link.Value
	}
	totalFlow := make([]float64, n)
	for i := range totalFlow {
		if outflow[i] > inflow[i] {
			totalFlow[i] = outflow[i]
		} else {
			totalFlow[i] = inflow[i]
		}
	}
	return totalFlow
}
