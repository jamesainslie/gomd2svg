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
	default:
		// For unsupported diagram kinds, return a minimal layout.
		return computeGraphLayout(g, th, cfg)
	}
}

// computeGraphLayout runs the full Sugiyama-style layout pipeline:
// 1. Size nodes based on text metrics
// 2. Compute rank assignments via topological sort
// 3. Order nodes within ranks to minimize crossings
// 4. Assign X, Y coordinates
// 5. Route edges with simple L-shaped polylines
// 6. Compute the bounding box
func computeGraphLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()

	// Step 1: Size all nodes.
	nodes := sizeNodes(g.Nodes, measurer, th, cfg)

	// Step 2: Compute ranks using topological ordering.
	nodeIDs := sortedNodeIDs(g.Nodes, g.NodeOrder)
	ranks := computeRanks(nodeIDs, g.Edges, g.NodeOrder)

	// Step 3: Order nodes within each rank (crossing minimization).
	layers := orderRankNodes(ranks, g.Edges, cfg.Flowchart.OrderPasses)

	// Step 4: Assign X, Y positions.
	positionNodes(layers, nodes, g.Direction, cfg)

	// Step 5: Route edges.
	edges := routeEdges(g.Edges, nodes, g.Direction)

	// Step 6: Compute bounding box.
	width, height := computeBoundingBox(nodes)

	return &Layout{
		Kind:    g.Kind,
		Nodes:   nodes,
		Edges:   edges,
		Width:   width,
		Height:  height,
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
