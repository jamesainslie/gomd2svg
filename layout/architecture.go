package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

// archLabelPadding is the horizontal padding added around a service label.
const archLabelPadding float32 = 20

func computeArchitectureLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	acfg := cfg.Architecture
	nodes := sizeArchNodes(graph, measurer, th, cfg)

	// Build grid positions from edge directional hints.
	gridPos := archBuildGridPositions(graph)

	// Convert grid to pixel coordinates.
	maxCol, maxRow := archGridMaxExtents(gridPos)

	for id, pos := range gridPos {
		node := nodes[id]
		if node == nil {
			continue
		}
		node.X = acfg.PaddingX + float32(pos[0])*(acfg.ServiceWidth+acfg.ColumnGap) + acfg.ServiceWidth/2
		node.Y = acfg.PaddingY + float32(pos[1])*(acfg.ServiceHeight+acfg.RowGap) + acfg.ServiceHeight/2
	}

	// Compute junction layouts.
	var junctions []ArchJunctionLayout
	for _, junc := range graph.ArchJunctions {
		node := nodes[junc.ID]
		if node == nil {
			continue
		}
		junctions = append(junctions, ArchJunctionLayout{
			ID:   junc.ID,
			X:    node.X,
			Y:    node.Y,
			Size: acfg.JunctionSize,
		})
	}

	// Compute group bounding rectangles.
	groups := make([]ArchGroupLayout, 0, len(graph.ArchGroups))
	for _, grp := range graph.ArchGroups {
		groupLayout := computeArchGroupBounds(grp, nodes, acfg)
		groups = append(groups, groupLayout)
	}

	// Build edges with side-based anchor points.
	var edges []*EdgeLayout
	for _, archEdge := range graph.ArchEdges {
		src := nodes[archEdge.FromID]
		dst := nodes[archEdge.ToID]
		if src == nil || dst == nil {
			continue
		}
		sx, sy := archAnchorPoint(src, archEdge.FromSide)
		dx, dy := archAnchorPoint(dst, archEdge.ToSide)
		edges = append(edges, &EdgeLayout{
			From:       archEdge.FromID,
			To:         archEdge.ToID,
			Points:     [][2]float32{{sx, sy}, {dx, dy}},
			ArrowStart: archEdge.ArrowLeft,
			ArrowEnd:   archEdge.ArrowRight,
		})
	}

	totalW := acfg.PaddingX*2 + float32(maxCol+1)*acfg.ServiceWidth + float32(maxCol)*acfg.ColumnGap
	totalH := acfg.PaddingY*2 + float32(maxRow+1)*acfg.ServiceHeight + float32(maxRow)*acfg.RowGap

	// Build service info map for rendering (icon data).
	svcInfo := make(map[string]ArchServiceInfo, len(graph.ArchServices))
	for _, svc := range graph.ArchServices {
		svcInfo[svc.ID] = ArchServiceInfo{Icon: svc.Icon}
	}

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  nodes,
		Edges:  edges,
		Width:  totalW,
		Height: totalH,
		Diagram: ArchitectureData{
			Groups:    groups,
			Junctions: junctions,
			Services:  svcInfo,
		},
	}
}

// archBuildGridPositions uses BFS to assign grid positions to all architecture
// services and junctions based on edge directional hints.
func archBuildGridPositions(graph *ir.Graph) map[string][2]int {
	gridPos := make(map[string][2]int)
	placed := make(map[string]bool)

	// Find the first node to place at origin.
	var firstID string
	if len(graph.ArchServices) > 0 {
		firstID = graph.ArchServices[0].ID
	} else if len(graph.ArchJunctions) > 0 {
		firstID = graph.ArchJunctions[0].ID
	}

	if firstID != "" {
		gridPos[firstID] = [2]int{0, 0}
		placed[firstID] = true

		// BFS to place connected nodes based on edge sides.
		queue := []string{firstID}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			curPos := gridPos[cur]

			for _, archEdge := range graph.ArchEdges {
				neighbor, dc, dr := archEdgeNeighborOffset(archEdge, cur, placed)
				if neighbor != "" && !placed[neighbor] {
					gridPos[neighbor] = [2]int{curPos[0] + dc, curPos[1] + dr}
					placed[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}
	}

	// Place any unplaced nodes in a row below.
	nextRow := archGridMaxRow(gridPos) + 1
	nextCol := 0
	for _, svc := range graph.ArchServices {
		if !placed[svc.ID] {
			gridPos[svc.ID] = [2]int{nextCol, nextRow}
			placed[svc.ID] = true
			nextCol++
		}
	}
	for _, junc := range graph.ArchJunctions {
		if !placed[junc.ID] {
			gridPos[junc.ID] = [2]int{nextCol, nextRow}
			placed[junc.ID] = true
			nextCol++
		}
	}

	// Normalize grid to non-negative.
	archNormalizeGrid(gridPos)

	return gridPos
}

// archEdgeNeighborOffset returns the unplaced neighbor and grid offset for an edge,
// given the current node ID.
func archEdgeNeighborOffset(archEdge *ir.ArchEdge, cur string, placed map[string]bool) (string, int, int) {
	if archEdge.FromID == cur && !placed[archEdge.ToID] {
		dc, dr := archSideOffset(archEdge.FromSide)
		return archEdge.ToID, dc, dr
	}
	if archEdge.ToID == cur && !placed[archEdge.FromID] {
		dc, dr := archSideOffsetReverse(archEdge.ToSide)
		return archEdge.FromID, dc, dr
	}
	return "", 0, 0
}

// archGridMaxRow returns the maximum row value in the grid.
func archGridMaxRow(gridPos map[string][2]int) int {
	maxRow := 0
	for _, pos := range gridPos {
		if pos[1] > maxRow {
			maxRow = pos[1]
		}
	}
	return maxRow
}

// archNormalizeGrid shifts all positions so that the minimum col and row are 0.
func archNormalizeGrid(gridPos map[string][2]int) {
	minCol, minRow := 0, 0
	for _, pos := range gridPos {
		if pos[0] < minCol {
			minCol = pos[0]
		}
		if pos[1] < minRow {
			minRow = pos[1]
		}
	}
	for id, pos := range gridPos {
		gridPos[id] = [2]int{pos[0] - minCol, pos[1] - minRow}
	}
}

// archGridMaxExtents returns the maximum column and row values.
func archGridMaxExtents(gridPos map[string][2]int) (int, int) {
	maxCol, maxRow := 0, 0
	for _, pos := range gridPos {
		if pos[0] > maxCol {
			maxCol = pos[0]
		}
		if pos[1] > maxRow {
			maxRow = pos[1]
		}
	}
	return maxCol, maxRow
}

// archSideOffset returns the grid displacement when moving FROM the given side.
// If A's Right side connects, the neighbor goes to the right (+1 col).
func archSideOffset(side ir.ArchSide) (int, int) {
	switch side {
	case ir.ArchRight:
		return 1, 0
	case ir.ArchLeft:
		return -1, 0
	case ir.ArchBottom:
		return 0, 1
	case ir.ArchTop:
		return 0, -1
	default:
		return 1, 0
	}
}

// archSideOffsetReverse returns the grid displacement when moving TO the given side.
// If the neighbor's Left side receives, the neighbor goes to the left.
func archSideOffsetReverse(side ir.ArchSide) (int, int) {
	switch side {
	case ir.ArchLeft:
		return -1, 0
	case ir.ArchRight:
		return 1, 0
	case ir.ArchTop:
		return 0, -1
	case ir.ArchBottom:
		return 0, 1
	default:
		return -1, 0
	}
}

// archAnchorPoint returns the pixel coordinate on a node's side.
func archAnchorPoint(node *NodeLayout, side ir.ArchSide) (float32, float32) {
	switch side {
	case ir.ArchLeft:
		return node.X - node.Width/2, node.Y
	case ir.ArchRight:
		return node.X + node.Width/2, node.Y
	case ir.ArchTop:
		return node.X, node.Y - node.Height/2
	case ir.ArchBottom:
		return node.X, node.Y + node.Height/2
	default:
		return node.X, node.Y
	}
}

func computeArchGroupBounds(grp *ir.ArchGroup, nodes map[string]*NodeLayout, acfg config.ArchitectureConfig) ArchGroupLayout {
	var minX, minY float32 = 1e9, 1e9
	var maxX, maxY float32 = -1e9, -1e9
	found := false

	for _, childID := range grp.Children {
		node := nodes[childID]
		if node == nil {
			continue
		}
		found = true
		left := node.X - node.Width/2
		right := node.X + node.Width/2
		top := node.Y - node.Height/2
		bottom := node.Y + node.Height/2
		if left < minX {
			minX = left
		}
		if right > maxX {
			maxX = right
		}
		if top < minY {
			minY = top
		}
		if bottom > maxY {
			maxY = bottom
		}
	}

	if !found {
		return ArchGroupLayout{ID: grp.ID, Label: grp.Label, Icon: grp.Icon}
	}

	pad := acfg.GroupPadding
	return ArchGroupLayout{
		ID:     grp.ID,
		Label:  grp.Label,
		Icon:   grp.Icon,
		X:      minX - pad,
		Y:      minY - pad,
		Width:  (maxX - minX) + 2*pad,
		Height: (maxY - minY) + 2*pad,
	}
}

func sizeArchNodes(graph *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(graph.Nodes))
	fontSize := th.FontSize
	fontFamily := th.FontFamily
	for id, node := range graph.Nodes {
		width := cfg.Architecture.ServiceWidth
		height := cfg.Architecture.ServiceHeight
		labelW := measurer.Width(node.Label, fontSize, fontFamily)
		if labelW+archLabelPadding > width {
			width = labelW + archLabelPadding
		}
		nodes[id] = &NodeLayout{
			ID:     id,
			Label:  TextBlock{Lines: []string{node.Label}, Width: labelW, Height: fontSize, FontSize: fontSize},
			Shape:  node.Shape,
			Width:  width,
			Height: height,
		}
	}
	return nodes
}
