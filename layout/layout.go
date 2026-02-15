package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

// ComputeLayout dispatches to the appropriate layout algorithm based on
// the diagram kind. Currently only flowchart/graph layout is implemented.
func ComputeLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	switch g.Kind {
	case ir.Flowchart:
		return computeGraphLayout(g, th, cfg)
	case ir.Class:
		return computeClassLayout(g, th, cfg)
	case ir.Er:
		return computeERLayout(g, th, cfg)
	case ir.State:
		return computeStateLayout(g, th, cfg)
	case ir.Sequence:
		return computeSequenceLayout(g, th, cfg)
	case ir.Kanban:
		return computeKanbanLayout(g, th, cfg)
	case ir.Packet:
		return computePacketLayout(g, th, cfg)
	default:
		// For unsupported diagram kinds, return a minimal layout.
		return computeGraphLayout(g, th, cfg)
	}
}

// sugiyamaResult holds the outputs of the shared Sugiyama pipeline.
type sugiyamaResult struct {
	Edges  []*EdgeLayout
	Width  float32
	Height float32
}

// runSugiyama runs the shared ranking, ordering, positioning, routing, and
// bounding box pipeline steps.
func runSugiyama(g *ir.Graph, nodes map[string]*NodeLayout, cfg *config.Layout) sugiyamaResult {
	nodeIDs := sortedNodeIDs(g.Nodes, g.NodeOrder)
	ranks := computeRanks(nodeIDs, g.Edges, g.NodeOrder)
	layers := orderRankNodes(ranks, g.Edges, cfg.Flowchart.OrderPasses)
	positionNodes(layers, nodes, g.Direction, cfg)
	edges := routeEdges(g.Edges, nodes, g.Direction)
	width, height := computeBoundingBox(nodes)
	return sugiyamaResult{Edges: edges, Width: width, Height: height}
}

// computeGraphLayout runs the full Sugiyama-style layout pipeline:
// 1. Size nodes based on text metrics
// 2. Run Sugiyama ranking, ordering, positioning, routing, and bounding box
func computeGraphLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()

	// Step 1: Size all nodes.
	nodes := sizeNodes(g.Nodes, measurer, th, cfg)

	// Step 2: Run Sugiyama pipeline.
	r := runSugiyama(g, nodes, cfg)

	return &Layout{
		Kind:    g.Kind,
		Nodes:   nodes,
		Edges:   r.Edges,
		Width:   r.Width,
		Height:  r.Height,
		Diagram: GraphData{},
	}
}

// computeBoundingBox finds the smallest rectangle containing all nodes,
// including the layout boundary padding.
func computeBoundingBox(nodes map[string]*NodeLayout) (float32, float32) {
	if len(nodes) == 0 {
		return 0, 0
	}

	var minX, minY float32
	var maxX, maxY float32
	first := true

	for _, n := range nodes {
		left := n.X - n.Width/2
		right := n.X + n.Width/2
		top := n.Y - n.Height/2
		bottom := n.Y + n.Height/2

		if first {
			minX = left
			maxX = right
			minY = top
			maxY = bottom
			first = false
		} else {
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
	}

	return maxX - minX + 2*layoutBoundaryPad, maxY - minY + 2*layoutBoundaryPad
}
