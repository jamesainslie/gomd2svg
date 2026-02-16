package layout

import (
	"container/heap"
	"math"

	"github.com/jamesainslie/gomd2svg/ir"
)

const (
	// defaultCellSize is the grid cell size for A* routing.
	defaultCellSize float32 = 8
	// defaultNodePad is the padding around nodes in the obstacle grid.
	defaultNodePad float32 = 4
	// gridMargin is the margin around the bounding box for the grid.
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

	obstacleGrid := buildGrid(nodes, defaultCellSize, defaultNodePad)

	for _, edge := range edges {
		src, srcOK := nodes[edge.From]
		dst, dstOK := nodes[edge.To]
		if !srcOK || !dstOK {
			continue
		}

		var points [][2]float32
		var labelAnchor [2]float32

		// Compute start/end points on node boundaries.
		startX, startY, endX, endY := edgeEndpoints(src, dst, direction)

		// Try A* routing.
		astarPath := obstacleGrid.findPath(startX, startY, endX, endY, edge.From, edge.To)
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
		if edge.Label != nil {
			label = &TextBlock{
				Lines:    []string{*edge.Label},
				FontSize: src.Label.FontSize,
			}
		}

		result = append(result, &EdgeLayout{
			From:           edge.From,
			To:             edge.To,
			Label:          label,
			Points:         points,
			LabelAnchor:    labelAnchor,
			Style:          edge.Style,
			ArrowStart:     edge.ArrowStart,
			ArrowEnd:       edge.ArrowEnd,
			ArrowStartKind: edge.ArrowStartKind,
			ArrowEndKind:   edge.ArrowEndKind,
		})
	}

	return result
}

// edgeEndpoints computes the start and end points on node boundaries
// based on the diagram direction.
func edgeEndpoints(src, dst *NodeLayout, direction ir.Direction) (float32, float32, float32, float32) {
	switch direction {
	case ir.LeftRight:
		return src.X + src.Width/2, src.Y, dst.X - dst.Width/2, dst.Y
	case ir.RightLeft:
		return src.X - src.Width/2, src.Y, dst.X + dst.Width/2, dst.Y
	case ir.BottomTop:
		return src.X, src.Y - src.Height/2, dst.X, dst.Y + dst.Height/2
	default: // TopDown
		return src.X, src.Y + src.Height/2, dst.X, dst.Y - dst.Height/2
	}
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
func buildGrid(nodes map[string]*NodeLayout, cellSize, nodePad float32) *grid { //nolint:unparam // nodePad kept as parameter for testability
	if len(nodes) == 0 {
		return &grid{cellSize: cellSize}
	}

	// Find bounding box of all nodes.
	var minX, minY, maxX, maxY float32
	first := true
	for _, node := range nodes {
		left := node.X - node.Width/2 - nodePad
		right := node.X + node.Width/2 + nodePad
		top := node.Y - node.Height/2 - nodePad
		bottom := node.Y + node.Height/2 + nodePad
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

	numCols := int(math.Ceil(float64((maxX - minX) / cellSize)))
	numRows := int(math.Ceil(float64((maxY - minY) / cellSize)))
	if numCols < 1 {
		numCols = 1
	}
	if numRows < 1 {
		numRows = 1
	}

	blocked := make([][]bool, numRows)
	nodeIDGrid := make([][][]string, numRows)
	for row := range blocked {
		blocked[row] = make([]bool, numCols)
		nodeIDGrid[row] = make([][]string, numCols)
	}

	obstacleGrid := &grid{
		blocked:  blocked,
		nodeIDs:  nodeIDGrid,
		originX:  minX,
		originY:  minY,
		cellSize: cellSize,
		cols:     numCols,
		rows:     numRows,
	}

	// Mark cells overlapping with node bounds (+ padding) as blocked.
	for _, node := range nodes {
		left := node.X - node.Width/2 - nodePad
		right := node.X + node.Width/2 + nodePad
		top := node.Y - node.Height/2 - nodePad
		bottom := node.Y + node.Height/2 + nodePad

		rMin, cMin := obstacleGrid.worldToCell(left, top)
		rMax, cMax := obstacleGrid.worldToCell(right, bottom)

		for row := rMin; row <= rMax; row++ {
			for col := cMin; col <= cMax; col++ {
				if row >= 0 && row < numRows && col >= 0 && col < numCols {
					obstacleGrid.blocked[row][col] = true
					obstacleGrid.nodeIDs[row][col] = append(obstacleGrid.nodeIDs[row][col], node.ID)
				}
			}
		}
	}

	return obstacleGrid
}

// worldToCell converts world coordinates to grid cell indices.
func (g *grid) worldToCell(worldX, worldY float32) (int, int) {
	gridCol := int((worldX - g.originX) / g.cellSize)
	gridRow := int((worldY - g.originY) / g.cellSize)
	return gridRow, gridCol
}

// cellToWorld converts grid cell indices to world coordinates (cell center).
func (g *grid) cellToWorld(row, col int) (float32, float32) {
	wx := g.originX + float32(col)*g.cellSize + g.cellSize/2
	wy := g.originY + float32(row)*g.cellSize + g.cellSize/2
	return wx, wy
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

	startRow, startCol := g.worldToCell(startX, startY)
	endRow, endCol := g.worldToCell(endX, endY)

	// Clamp to grid bounds.
	startRow = clampInt(startRow, 0, g.rows-1)
	startCol = clampInt(startCol, 0, g.cols-1)
	endRow = clampInt(endRow, 0, g.rows-1)
	endCol = clampInt(endCol, 0, g.cols-1)

	if startRow == endRow && startCol == endCol {
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

	start := cell{startRow, startCol}
	end := cell{endRow, endCol}

	gScore[start] = 0

	heuristic := func(target cell) float32 {
		deltaRow := target.row - end.row
		deltaCol := target.col - end.col
		if deltaRow < 0 {
			deltaRow = -deltaRow
		}
		if deltaCol < 0 {
			deltaCol = -deltaCol
		}
		return float32(deltaRow + deltaCol)
	}

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{row: startRow, col: startCol, f: heuristic(start)})
	dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	found := false
	for pq.Len() > 0 {
		cur, _ := heap.Pop(pq).(*pqItem) //nolint:errcheck // type assertion is safe for our priority queue
		currentCell := cell{cur.row, cur.col}

		if visited[currentCell] {
			continue
		}
		visited[currentCell] = true

		if currentCell == end {
			found = true
			break
		}

		curG := gScore[currentCell]

		for _, dir := range dirs {
			neighborRow, neighborCol := cur.row+dir[0], cur.col+dir[1]
			if neighborRow < 0 || neighborRow >= g.rows || neighborCol < 0 || neighborCol >= g.cols {
				continue
			}
			if g.isBlockedExcluding(neighborRow, neighborCol, fromNode, toNode) {
				continue
			}

			next := cell{neighborRow, neighborCol}
			newG := curG + 1

			if prev, ok := gScore[next]; ok && newG >= prev {
				continue
			}

			gScore[next] = newG
			parent[next] = currentCell
			fScore := newG + heuristic(next)
			heap.Push(pq, &pqItem{row: neighborRow, col: neighborCol, f: fScore})
		}
	}

	if !found {
		return nil
	}

	// Reconstruct path.
	var path []cell
	currentCell := end
	for currentCell != start {
		path = append(path, currentCell)
		currentCell = parent[currentCell]
	}
	path = append(path, start)

	// Reverse path.
	for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
		path[left], path[right] = path[right], path[left]
	}

	// Convert to world coordinates. Use exact start/end for first/last points.
	points := make([][2]float32, len(path))
	for idx, pathCell := range path {
		worldX, worldY := g.cellToWorld(pathCell.row, pathCell.col)
		points[idx] = [2]float32{worldX, worldY}
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
	for idx := 1; idx < len(pts)-1; idx++ {
		prev := result[len(result)-1]
		next := pts[idx+1]
		cur := pts[idx]
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
	for idx := 1; idx < len(pts); idx++ {
		deltaX := pts[idx][0] - pts[idx-1][0]
		deltaY := pts[idx][1] - pts[idx-1][1]
		totalLen += float32(math.Sqrt(float64(deltaX*deltaX + deltaY*deltaY)))
	}

	// Walk to the halfway point.
	halfLen := totalLen / 2
	var walked float32
	for idx := 1; idx < len(pts); idx++ {
		deltaX := pts[idx][0] - pts[idx-1][0]
		deltaY := pts[idx][1] - pts[idx-1][1]
		segLen := float32(math.Sqrt(float64(deltaX*deltaX + deltaY*deltaY)))
		if walked+segLen >= halfLen && segLen > 0 {
			frac := (halfLen - walked) / segLen
			return [2]float32{
				pts[idx-1][0] + deltaX*frac,
				pts[idx-1][1] + deltaY*frac,
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

func clampInt(val, low, high int) int { //nolint:unparam // low kept as parameter for generality
	if val < low {
		return low
	}
	if val > high {
		return high
	}
	return val
}

// Priority queue for A*.
type pqItem struct {
	row, col int
	f        float32
	index    int
}

type priorityQueue []*pqItem

func (pq priorityQueue) Len() int                 { return len(pq) }
func (pq priorityQueue) Less(idxA, idxB int) bool { return pq[idxA].f < pq[idxB].f }
func (pq priorityQueue) Swap(idxA, idxB int) {
	pq[idxA], pq[idxB] = pq[idxB], pq[idxA]
	pq[idxA].index = idxA
	pq[idxB].index = idxB
}

func (pq *priorityQueue) Push(x any) {
	item, _ := x.(*pqItem) //nolint:errcheck // type assertion is safe for our priority queue
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	length := len(old)
	item := old[length-1]
	old[length-1] = nil
	item.index = -1
	*pq = old[:length-1]
	return item
}
