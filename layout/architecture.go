package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeArchitectureLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	acfg := cfg.Architecture
	nodes := sizeArchNodes(g, measurer, th, cfg)

	// Build grid positions from edge directional hints.
	gridPos := make(map[string][2]int) // id -> [col, row]
	placed := make(map[string]bool)

	// Place first service/junction at origin.
	var firstID string
	if len(g.ArchServices) > 0 {
		firstID = g.ArchServices[0].ID
	} else if len(g.ArchJunctions) > 0 {
		firstID = g.ArchJunctions[0].ID
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

			for _, e := range g.ArchEdges {
				var neighbor string
				var dc, dr int
				if e.FromID == cur && !placed[e.ToID] {
					neighbor = e.ToID
					dc, dr = archSideOffset(e.FromSide)
				} else if e.ToID == cur && !placed[e.FromID] {
					neighbor = e.FromID
					dc, dr = archSideOffsetReverse(e.ToSide)
				}
				if neighbor != "" && !placed[neighbor] {
					gridPos[neighbor] = [2]int{curPos[0] + dc, curPos[1] + dr}
					placed[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}
	}

	// Place any unplaced nodes in a row below.
	nextRow := 0
	for _, pos := range gridPos {
		if pos[1] > nextRow {
			nextRow = pos[1]
		}
	}
	nextRow++
	nextCol := 0
	for _, svc := range g.ArchServices {
		if !placed[svc.ID] {
			gridPos[svc.ID] = [2]int{nextCol, nextRow}
			placed[svc.ID] = true
			nextCol++
		}
	}
	for _, junc := range g.ArchJunctions {
		if !placed[junc.ID] {
			gridPos[junc.ID] = [2]int{nextCol, nextRow}
			placed[junc.ID] = true
			nextCol++
		}
	}

	// Normalize grid to non-negative.
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

	// Convert grid to pixel coordinates.
	maxCol, maxRow := 0, 0
	for _, pos := range gridPos {
		if pos[0] > maxCol {
			maxCol = pos[0]
		}
		if pos[1] > maxRow {
			maxRow = pos[1]
		}
	}

	for id, pos := range gridPos {
		n := nodes[id]
		if n == nil {
			continue
		}
		n.X = acfg.PaddingX + float32(pos[0])*(acfg.ServiceWidth+acfg.ColumnGap) + acfg.ServiceWidth/2
		n.Y = acfg.PaddingY + float32(pos[1])*(acfg.ServiceHeight+acfg.RowGap) + acfg.ServiceHeight/2
	}

	// Compute junction layouts.
	var junctions []ArchJunctionLayout
	for _, junc := range g.ArchJunctions {
		n := nodes[junc.ID]
		if n == nil {
			continue
		}
		junctions = append(junctions, ArchJunctionLayout{
			ID:   junc.ID,
			X:    n.X,
			Y:    n.Y,
			Size: acfg.JunctionSize,
		})
	}

	// Compute group bounding rectangles.
	var groups []ArchGroupLayout
	for _, grp := range g.ArchGroups {
		gl := computeArchGroupBounds(grp, nodes, acfg)
		groups = append(groups, gl)
	}

	// Build edges with side-based anchor points.
	var edges []*EdgeLayout
	for _, e := range g.ArchEdges {
		src := nodes[e.FromID]
		dst := nodes[e.ToID]
		if src == nil || dst == nil {
			continue
		}
		sx, sy := archAnchorPoint(src, e.FromSide)
		dx, dy := archAnchorPoint(dst, e.ToSide)
		edges = append(edges, &EdgeLayout{
			From:       e.FromID,
			To:         e.ToID,
			Points:     [][2]float32{{sx, sy}, {dx, dy}},
			ArrowStart: e.ArrowLeft,
			ArrowEnd:   e.ArrowRight,
		})
	}

	totalW := acfg.PaddingX*2 + float32(maxCol+1)*acfg.ServiceWidth + float32(maxCol)*acfg.ColumnGap
	totalH := acfg.PaddingY*2 + float32(maxRow+1)*acfg.ServiceHeight + float32(maxRow)*acfg.RowGap

	// Build service info map for rendering (icon data).
	svcInfo := make(map[string]ArchServiceInfo, len(g.ArchServices))
	for _, svc := range g.ArchServices {
		svcInfo[svc.ID] = ArchServiceInfo{Icon: svc.Icon}
	}

	return &Layout{
		Kind:   g.Kind,
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

// archSideOffset returns the grid displacement when moving FROM the given side.
// If A's Right side connects, the neighbor goes to the right (+1 col).
func archSideOffset(side ir.ArchSide) (dc, dr int) {
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
func archSideOffsetReverse(side ir.ArchSide) (dc, dr int) {
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
func archAnchorPoint(n *NodeLayout, side ir.ArchSide) (float32, float32) {
	switch side {
	case ir.ArchLeft:
		return n.X - n.Width/2, n.Y
	case ir.ArchRight:
		return n.X + n.Width/2, n.Y
	case ir.ArchTop:
		return n.X, n.Y - n.Height/2
	case ir.ArchBottom:
		return n.X, n.Y + n.Height/2
	default:
		return n.X, n.Y
	}
}

func computeArchGroupBounds(grp *ir.ArchGroup, nodes map[string]*NodeLayout, acfg config.ArchitectureConfig) ArchGroupLayout {
	var minX, minY float32 = 1e9, 1e9
	var maxX, maxY float32 = -1e9, -1e9
	found := false

	for _, childID := range grp.Children {
		n := nodes[childID]
		if n == nil {
			continue
		}
		found = true
		left := n.X - n.Width/2
		right := n.X + n.Width/2
		top := n.Y - n.Height/2
		bottom := n.Y + n.Height/2
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

func sizeArchNodes(g *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(g.Nodes))
	fontSize := th.FontSize
	fontFamily := th.FontFamily
	for id, node := range g.Nodes {
		w := cfg.Architecture.ServiceWidth
		h := cfg.Architecture.ServiceHeight
		labelW := measurer.Width(node.Label, fontSize, fontFamily)
		if labelW+20 > w {
			w = labelW + 20
		}
		nodes[id] = &NodeLayout{
			ID:     id,
			Label:  TextBlock{Lines: []string{node.Label}, Width: labelW, Height: fontSize, FontSize: fontSize},
			Shape:  node.Shape,
			Width:  w,
			Height: h,
		}
	}
	return nodes
}
