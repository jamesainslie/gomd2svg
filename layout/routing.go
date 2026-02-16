package layout

import (
	"container/heap"
	"math"

	"github.com/jamesainslie/mermaid-go/ir"
)

const (
	// Default grid cell size for A* routing.
	defaultCellSize float32 = 8
	// Default padding around nodes in the obstacle grid.
	defaultNodePad float32 = 4
	// Margin around the bounding box for the grid.
	gridMargin float32 = 40
)

// routeEdges computes edge routes using A* pathfinding that avoids node overlap.
// Falls back to L-shaped routes when A* cannot find a path.
func routeEdges(
	edges []*ir.Edge,
	nodes map[string]*NodeLayout,
	direction ir.Direction,
) []*EdgeLayout {
	result := make([]*EdgeLayout, 0, len(edges))

	g := buildGrid(nodes, defaultCellSize, defaultNodePad)

	for _, e := range edges {
		src, srcOK := nodes[e.From]
		dst, dstOK := nodes[e.To]
		if !srcOK || !dstOK {
			continue
		}

		var points [][2]float32
		var labelAnchor [2]float32

		// Compute start/end points on node boundaries.
		startX, startY, endX, endY := edgeEndpoints(src, dst, direction)

		// Try A* routing.
		astarPath := g.findPath(startX, startY, endX, endY, e.From, e.To)
		if astarPath != nil {
			points = simplifyPath(astarPath)
			labelAnchor = pathMidpoint(points)
		} else {
			// Fallback to L-shaped routing.
			switch direction {
			case ir.LeftRight, ir.RightLeft:
				points, labelAnchor = routeLR(src, dst)
			default:
				points, labelAnchor = routeTD(src, dst)
			}
		}

		var label *TextBlock
		if e.Label != nil {
			label = &TextBlock{
				Lines:    []string{*e.Label},
				FontSize: src.Label.FontSize,
			}
		}

		result = append(result, &EdgeLayout{
			From:           e.From,
			To:             e.To,
			Label:          label,
			Points:         points,
			LabelAnchor:    labelAnchor,
			Style:          e.Style,
			ArrowStart:     e.ArrowStart,
			ArrowEnd:       e.ArrowEnd,
			ArrowStartKind: e.ArrowStartKind,
			ArrowEndKind:   e.ArrowEndKind,
		})
	}

	return result
}

// edgeEndpoints computes the start and end points on node boundaries
// based on the diagram direction.
func edgeEndpoints(src, dst *NodeLayout, direction ir.Direction) (startX, startY, endX, endY float32) {
	switch direction {
	case ir.LeftRight:
		startX = src.X + src.Width/2
		startY = src.Y
		endX = dst.X - dst.Width/2
		endY = dst.Y
	case ir.RightLeft:
		startX = src.X - src.Width/2
		startY = src.Y
		endX = dst.X + dst.Width/2
		endY = dst.Y
	case ir.BottomTop:
		startX = src.X
		startY = src.Y - src.Height/2
		endX = dst.X
		endY = dst.Y + dst.Height/2
	default: // TopDown
		startX = src.X
		startY = src.Y + src.Height/2
		endX = dst.X
		endY = dst.Y - dst.Height/2
	}
	return
}

// grid represents a 2D obstacle grid for A* pathfinding.
type grid struct {
	blocked  [][]bool
	nodeIDs  [][][]string // which node IDs block each cell (empty if free)
	originX  float32
	originY  float32
	cellSize float32
	cols     int
	rows     int
}

// buildGrid constructs an obstacle grid from positioned nodes.
func buildGrid(nodes map[string]*NodeLayout, cellSize, nodePad float32) *grid {
	if len(nodes) == 0 {
		return &grid{cellSize: cellSize}
	}

	// Find bounding box of all nodes.
	var minX, minY, maxX, maxY float32
	first := true
	for _, n := range nodes {
		left := n.X - n.Width/2 - nodePad
		right := n.X + n.Width/2 + nodePad
		top := n.Y - n.Height/2 - nodePad
		bottom := n.Y + n.Height/2 + nodePad
		if first {
			minX, minY, maxX, maxY = left, top, right, bottom
			first = false
		} else {
			if left < minX {
				minX = left
			}
			if top < minY {
				minY = top
			}
			if right > maxX {
				maxX = right
			}
			if bottom > maxY {
				maxY = bottom
			}
		}
	}

	// Expand by margin.
	minX -= gridMargin
	minY -= gridMargin
	maxX += gridMargin
	maxY += gridMargin

	cols := int(math.Ceil(float64((maxX - minX) / cellSize)))
	rows := int(math.Ceil(float64((maxY - minY) / cellSize)))
	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}

	blocked := make([][]bool, rows)
	nodeIDGrid := make([][][]string, rows)
	for r := range blocked {
		blocked[r] = make([]bool, cols)
		nodeIDGrid[r] = make([][]string, cols)
	}

	g := &grid{
		blocked:  blocked,
		nodeIDs:  nodeIDGrid,
		originX:  minX,
		originY:  minY,
		cellSize: cellSize,
		cols:     cols,
		rows:     rows,
	}

	// Mark cells overlapping with node bounds (+ padding) as blocked.
	for _, n := range nodes {
		left := n.X - n.Width/2 - nodePad
		right := n.X + n.Width/2 + nodePad
		top := n.Y - n.Height/2 - nodePad
		bottom := n.Y + n.Height/2 + nodePad

		rMin, cMin := g.worldToCell(left, top)
		rMax, cMax := g.worldToCell(right, bottom)

		for r := rMin; r <= rMax; r++ {
			for c := cMin; c <= cMax; c++ {
				if r >= 0 && r < rows && c >= 0 && c < cols {
					g.blocked[r][c] = true
					g.nodeIDs[r][c] = append(g.nodeIDs[r][c], n.ID)
				}
			}
		}
	}

	return g
}

// worldToCell converts world coordinates to grid cell indices.
func (g *grid) worldToCell(wx, wy float32) (row, col int) {
	col = int((wx - g.originX) / g.cellSize)
	row = int((wy - g.originY) / g.cellSize)
	return
}

// cellToWorld converts grid cell indices to world coordinates (cell center).
func (g *grid) cellToWorld(row, col int) (wx, wy float32) {
	wx = g.originX + float32(col)*g.cellSize + g.cellSize/2
	wy = g.originY + float32(row)*g.cellSize + g.cellSize/2
	return
}

// isBlocked returns true if the cell at (row, col) is an obstacle.
func (g *grid) isBlocked(row, col int) bool {
	if row < 0 || row >= g.rows || col < 0 || col >= g.cols {
		return true // out of bounds = blocked
	}
	return g.blocked[row][col]
}

// isBlockedExcluding returns true if the cell is blocked by a node other than
// the specified excluded IDs (used to allow paths through source/target nodes).
// A cell is passable if any of its occupying node IDs match an excluded ID,
// because the edge is allowed to traverse its own source/target nodes' cells
// even when other nodes' padding overlaps.
func (g *grid) isBlockedExcluding(row, col int, excludeA, excludeB string) bool {
	if row < 0 || row >= g.rows || col < 0 || col >= g.cols {
		return true
	}
	if !g.blocked[row][col] {
		return false
	}
	for _, id := range g.nodeIDs[row][col] {
		if id == excludeA || id == excludeB {
			return false // passable: source or target node occupies this cell
		}
	}
	return true
}

// findPath runs A* from (startX, startY) to (endX, endY), treating cells
// belonging to fromNode or toNode as passable.
func (g *grid) findPath(startX, startY, endX, endY float32, fromNode, toNode string) [][2]float32 {
	if g.rows == 0 || g.cols == 0 {
		return nil
	}

	sr, sc := g.worldToCell(startX, startY)
	er, ec := g.worldToCell(endX, endY)

	// Clamp to grid bounds.
	sr = clampInt(sr, 0, g.rows-1)
	sc = clampInt(sc, 0, g.cols-1)
	er = clampInt(er, 0, g.rows-1)
	ec = clampInt(ec, 0, g.cols-1)

	if sr == er && sc == ec {
		return [][2]float32{{startX, startY}, {endX, endY}}
	}

	// A* with 4-directional movement.
	type cell struct {
		row, col int
	}

	// Cost and parent tracking.
	gScore := make(map[cell]float32)
	parent := make(map[cell]cell)
	visited := make(map[cell]bool)

	start := cell{sr, sc}
	end := cell{er, ec}

	gScore[start] = 0

	heuristic := func(c cell) float32 {
		dr := c.row - end.row
		dc := c.col - end.col
		if dr < 0 {
			dr = -dr
		}
		if dc < 0 {
			dc = -dc
		}
		return float32(dr + dc)
	}

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{row: sr, col: sc, f: heuristic(start)})

	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	found := false
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*pqItem)
		c := cell{cur.row, cur.col}

		if visited[c] {
			continue
		}
		visited[c] = true

		if c == end {
			found = true
			break
		}

		curG := gScore[c]

		for _, d := range dirs {
			nr, nc := cur.row+d[0], cur.col+d[1]
			if nr < 0 || nr >= g.rows || nc < 0 || nc >= g.cols {
				continue
			}
			if g.isBlockedExcluding(nr, nc, fromNode, toNode) {
				continue
			}

			next := cell{nr, nc}
			newG := curG + 1

			if prev, ok := gScore[next]; ok && newG >= prev {
				continue
			}

			gScore[next] = newG
			parent[next] = c
			f := newG + heuristic(next)
			heap.Push(pq, &pqItem{row: nr, col: nc, f: f})
		}
	}

	if !found {
		return nil
	}

	// Reconstruct path.
	var path []cell
	c := end
	for c != start {
		path = append(path, c)
		c = parent[c]
	}
	path = append(path, start)

	// Reverse path.
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	// Convert to world coordinates. Use exact start/end for first/last points.
	points := make([][2]float32, len(path))
	for i, c := range path {
		wx, wy := g.cellToWorld(c.row, c.col)
		points[i] = [2]float32{wx, wy}
	}
	points[0] = [2]float32{startX, startY}
	points[len(points)-1] = [2]float32{endX, endY}

	return points
}

// simplifyPath removes collinear intermediate points from an axis-aligned polyline.
// It assumes all segments are horizontal or vertical (as produced by 4-directional A*).
func simplifyPath(pts [][2]float32) [][2]float32 {
	if len(pts) <= 2 {
		return pts
	}

	result := [][2]float32{pts[0]}
	for i := 1; i < len(pts)-1; i++ {
		prev := result[len(result)-1]
		next := pts[i+1]
		cur := pts[i]
		// Keep point if direction changes (not collinear).
		sameX := prev[0] == cur[0] && cur[0] == next[0]
		sameY := prev[1] == cur[1] && cur[1] == next[1]
		if !sameX && !sameY {
			result = append(result, cur)
		}
	}
	result = append(result, pts[len(pts)-1])
	return result
}

// pathMidpoint returns the point at the middle of the path's total length.
func pathMidpoint(pts [][2]float32) [2]float32 {
	if len(pts) == 0 {
		return [2]float32{}
	}
	if len(pts) == 1 {
		return pts[0]
	}

	// Compute total path length.
	var totalLen float32
	for i := 1; i < len(pts); i++ {
		dx := pts[i][0] - pts[i-1][0]
		dy := pts[i][1] - pts[i-1][1]
		totalLen += float32(math.Sqrt(float64(dx*dx + dy*dy)))
	}

	// Walk to the halfway point.
	halfLen := totalLen / 2
	var walked float32
	for i := 1; i < len(pts); i++ {
		dx := pts[i][0] - pts[i-1][0]
		dy := pts[i][1] - pts[i-1][1]
		segLen := float32(math.Sqrt(float64(dx*dx + dy*dy)))
		if walked+segLen >= halfLen && segLen > 0 {
			t := (halfLen - walked) / segLen
			return [2]float32{
				pts[i-1][0] + dx*t,
				pts[i-1][1] + dy*t,
			}
		}
		walked += segLen
	}

	return pts[len(pts)-1]
}

// routeLR creates an L-shaped fallback route for left-right edges.
func routeLR(src, dst *NodeLayout) ([][2]float32, [2]float32) {
	startX := src.X + src.Width/2
	startY := src.Y
	endX := dst.X - dst.Width/2
	endY := dst.Y

	midX := (startX + endX) / 2

	points := [][2]float32{
		{startX, startY},
		{midX, startY},
		{midX, endY},
		{endX, endY},
	}

	labelAnchor := [2]float32{midX, (startY + endY) / 2}

	return points, labelAnchor
}

// routeTD creates an L-shaped fallback route for top-down edges.
func routeTD(src, dst *NodeLayout) ([][2]float32, [2]float32) {
	startX := src.X
	startY := src.Y + src.Height/2
	endX := dst.X
	endY := dst.Y - dst.Height/2

	midY := (startY + endY) / 2

	points := [][2]float32{
		{startX, startY},
		{startX, midY},
		{endX, midY},
		{endX, endY},
	}

	labelAnchor := [2]float32{(startX + endX) / 2, midY}

	return points, labelAnchor
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Priority queue for A*.
type pqItem struct {
	row, col int
	f        float32
	index    int
}

type priorityQueue []*pqItem

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].f < pq[j].f }
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	item := x.(*pqItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[:n-1]
	return item
}
