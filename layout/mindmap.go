package layout

import (
	"math"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

// mindmapLayoutNode is an internal type that wraps MindmapNodeLayout with
// tree-traversal fields used during the radial layout computation.
type mindmapLayoutNode struct {
	MindmapNodeLayout
	children    []*mindmapLayoutNode
	subtreeSpan float32
}

// computeMindmapLayout builds a radial tree layout for a mindmap diagram.
// The root node is placed at the center and children are distributed radially
// in concentric rings at increasing distances.
func computeMindmapLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	if g.MindmapRoot == nil {
		return &Layout{
			Kind:    g.Kind,
			Nodes:   map[string]*NodeLayout{},
			Width:   100,
			Height:  100,
			Diagram: MindmapData{},
		}
	}

	measurer := textmetrics.New()
	padX := cfg.Mindmap.PaddingX
	padY := cfg.Mindmap.PaddingY
	nodePad := cfg.Mindmap.NodePadding
	levelSpacing := cfg.Mindmap.LevelSpacing
	branchSpacing := cfg.Mindmap.BranchSpacing

	// Phase 1: Build layout tree with measured node sizes.
	root := mindmapBuildLayoutTree(g.MindmapRoot, measurer, th, cfg, nodePad, 0, 0)

	// Phase 2: Compute subtree angular spans bottom-up.
	mindmapComputeSubtreeSize(root, branchSpacing)

	// Phase 3: Position nodes radially from center.
	root.X = 0
	root.Y = 0
	if len(root.children) > 0 {
		var totalSpan float64
		for _, c := range root.children {
			totalSpan += float64(c.subtreeSpan)
		}
		startAngle := -math.Pi / 2
		for _, c := range root.children {
			fraction := float64(c.subtreeSpan) / totalSpan
			midAngle := startAngle + fraction*math.Pi
			dist := levelSpacing + root.Width/2 + c.Width/2
			c.X = float32(math.Cos(midAngle)) * dist
			c.Y = float32(math.Sin(midAngle)) * dist
			mindmapPositionChildren(c, midAngle, levelSpacing, 2)
			startAngle += fraction * 2 * math.Pi
		}
	}

	// Phase 4: Normalize to positive coordinates.
	minX, minY, maxX, maxY := mindmapBounds(root)
	shiftX := padX - minX
	shiftY := padY - minY
	mindmapShift(root, shiftX, shiftY)

	totalW := (maxX - minX) + padX*2
	totalH := (maxY - minY) + padY*2

	return &Layout{
		Kind:    g.Kind,
		Nodes:   map[string]*NodeLayout{},
		Width:   totalW,
		Height:  totalH,
		Diagram: MindmapData{Root: &root.MindmapNodeLayout},
	}
}

// mindmapBuildLayoutTree recursively constructs the internal layout tree from
// the IR mindmap tree. Each node is measured and sized, and children are linked
// both in the internal tree (for traversal) and in the exported
// MindmapNodeLayout.Children slice (for the renderer).
func mindmapBuildLayoutTree(
	node *ir.MindmapNode,
	measurer *textmetrics.Measurer,
	th *theme.Theme,
	cfg *config.Layout,
	nodePad float32,
	depth int,
	branchIdx int,
) *mindmapLayoutNode {
	textW := measurer.Width(node.Label, th.FontSize, th.FontFamily)
	textH := th.FontSize * cfg.LabelLineHeight

	ln := &mindmapLayoutNode{}
	ln.Label = node.Label
	ln.Shape = node.Shape
	ln.Icon = node.Icon
	ln.Width = textW + nodePad*2
	ln.Height = textH + nodePad*2
	ln.ColorIndex = branchIdx

	for i, child := range node.Children {
		childBranch := branchIdx
		if depth == 0 {
			childBranch = i
		}
		childNode := mindmapBuildLayoutTree(child, measurer, th, cfg, nodePad, depth+1, childBranch)
		ln.children = append(ln.children, childNode)
		ln.Children = append(ln.Children, &childNode.MindmapNodeLayout)
	}

	return ln
}

// mindmapComputeSubtreeSize computes the angular span each subtree requires,
// measured in pixels of arc at a reference distance. Leaf nodes get a span
// equal to their height plus branchSpacing; parent nodes sum their children.
func mindmapComputeSubtreeSize(node *mindmapLayoutNode, branchSpacing float32) {
	if len(node.children) == 0 {
		node.subtreeSpan = node.Height + branchSpacing
		return
	}
	var total float32
	for _, c := range node.children {
		mindmapComputeSubtreeSize(c, branchSpacing)
		total += c.subtreeSpan
	}
	node.subtreeSpan = total
}

// mindmapPositionChildren recursively positions a node's children in a cone
// centered on parentAngle. The cone narrows at deeper levels to avoid overlap.
func mindmapPositionChildren(parent *mindmapLayoutNode, parentAngle float64, levelSpacing float32, depth int) {
	if len(parent.children) == 0 {
		return
	}

	// Narrower spread at deeper levels.
	spreadAngle := math.Pi / float64(depth+1)

	var totalSpan float64
	for _, c := range parent.children {
		totalSpan += float64(c.subtreeSpan)
	}

	startAngle := parentAngle - spreadAngle/2
	for _, c := range parent.children {
		fraction := float64(c.subtreeSpan) / totalSpan
		midAngle := startAngle + fraction*spreadAngle/2
		dist := levelSpacing + parent.Width/2 + c.Width/2
		c.X = parent.X + float32(math.Cos(midAngle))*dist
		c.Y = parent.Y + float32(math.Sin(midAngle))*dist
		mindmapPositionChildren(c, midAngle, levelSpacing, depth+1)
		startAngle += fraction * spreadAngle
	}
}

// mindmapBounds recursively computes the bounding box of the entire tree,
// accounting for node half-widths and half-heights.
func mindmapBounds(node *mindmapLayoutNode) (minX, minY, maxX, maxY float32) {
	halfW := node.Width / 2
	halfH := node.Height / 2
	minX = node.X - halfW
	minY = node.Y - halfH
	maxX = node.X + halfW
	maxY = node.Y + halfH

	for _, c := range node.children {
		cMinX, cMinY, cMaxX, cMaxY := mindmapBounds(c)
		if cMinX < minX {
			minX = cMinX
		}
		if cMinY < minY {
			minY = cMinY
		}
		if cMaxX > maxX {
			maxX = cMaxX
		}
		if cMaxY > maxY {
			maxY = cMaxY
		}
	}
	return
}

// mindmapShift recursively offsets all node coordinates by (dx, dy), updating
// both the exported MindmapNodeLayout fields (via the embedded struct) for the
// renderer.
func mindmapShift(node *mindmapLayoutNode, dx, dy float32) {
	node.X += dx
	node.Y += dy
	for _, c := range node.children {
		mindmapShift(c, dx, dy)
	}
}
