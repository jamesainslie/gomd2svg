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
// mindmapEmptySize is the default width/height for an empty mindmap layout.
const mindmapEmptySize float32 = 100

func computeMindmapLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	if graph.MindmapRoot == nil {
		return &Layout{
			Kind:    graph.Kind,
			Nodes:   map[string]*NodeLayout{},
			Width:   mindmapEmptySize,
			Height:  mindmapEmptySize,
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
	root := mindmapBuildLayoutTree(graph.MindmapRoot, measurer, th, cfg, nodePad, 0, 0)

	// Phase 2: Compute subtree angular spans bottom-up.
	mindmapComputeSubtreeSize(root, branchSpacing)

	// Phase 3: Position nodes radially from center.
	root.X = 0
	root.Y = 0
	if len(root.children) > 0 {
		var totalSpan float64
		for _, child := range root.children {
			totalSpan += float64(child.subtreeSpan)
		}
		startAngle := -math.Pi / 2
		for _, child := range root.children {
			fraction := float64(child.subtreeSpan) / totalSpan
			midAngle := startAngle + fraction*math.Pi
			dist := levelSpacing + root.Width/2 + child.Width/2
			child.X = float32(math.Cos(midAngle)) * dist
			child.Y = float32(math.Sin(midAngle)) * dist
			mindmapPositionChildren(child, midAngle, levelSpacing, 2)
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
		Kind:    graph.Kind,
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
	for _, child := range node.children {
		mindmapComputeSubtreeSize(child, branchSpacing)
		total += child.subtreeSpan
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
	for _, child := range parent.children {
		totalSpan += float64(child.subtreeSpan)
	}

	startAngle := parentAngle - spreadAngle/2
	for _, child := range parent.children {
		fraction := float64(child.subtreeSpan) / totalSpan
		midAngle := startAngle + fraction*spreadAngle/2
		dist := levelSpacing + parent.Width/2 + child.Width/2
		child.X = parent.X + float32(math.Cos(midAngle))*dist
		child.Y = parent.Y + float32(math.Sin(midAngle))*dist
		mindmapPositionChildren(child, midAngle, levelSpacing, depth+1)
		startAngle += fraction * spreadAngle
	}
}

// mindmapBounds recursively computes the bounding box of the entire tree,
// accounting for node half-widths and half-heights.
func mindmapBounds(node *mindmapLayoutNode) (float32, float32, float32, float32) {
	halfW := node.Width / 2
	halfH := node.Height / 2
	boundsMinX := node.X - halfW
	boundsMinY := node.Y - halfH
	boundsMaxX := node.X + halfW
	boundsMaxY := node.Y + halfH

	for _, child := range node.children {
		cMinX, cMinY, cMaxX, cMaxY := mindmapBounds(child)
		if cMinX < boundsMinX {
			boundsMinX = cMinX
		}
		if cMinY < boundsMinY {
			boundsMinY = cMinY
		}
		if cMaxX > boundsMaxX {
			boundsMaxX = cMaxX
		}
		if cMaxY > boundsMaxY {
			boundsMaxY = cMaxY
		}
	}
	return boundsMinX, boundsMinY, boundsMaxX, boundsMaxY
}

// mindmapShift recursively offsets all node coordinates by (dx, dy), updating
// both the exported MindmapNodeLayout fields (via the embedded struct) for the
// renderer.
func mindmapShift(node *mindmapLayoutNode, dx, dy float32) {
	node.X += dx
	node.Y += dy
	for _, child := range node.children {
		mindmapShift(child, dx, dy)
	}
}
