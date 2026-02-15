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
	case ir.Pie:
		return computePieLayout(g, th, cfg)
	case ir.Quadrant:
		return computeQuadrantLayout(g, th, cfg)
	case ir.Timeline:
		return computeTimelineLayout(g, th, cfg)
	case ir.Gantt:
		return computeGanttLayout(g, th, cfg)
	case ir.GitGraph:
		return computeGitGraphLayout(g, th, cfg)
	case ir.XYChart:
		return computeXYChartLayout(g, th, cfg)
	case ir.Radar:
		return computeRadarLayout(g, th, cfg)
	case ir.Mindmap:
		return computeMindmapLayout(g, th, cfg)
	case ir.Sankey:
		return computeSankeyLayout(g, th, cfg)
	case ir.Treemap:
		return computeTreemapLayout(g, th, cfg)
	case ir.Requirement:
		return computeRequirementLayout(g, th, cfg)
	case ir.Block:
		return computeBlockLayout(g, th, cfg)
	case ir.C4:
		return computeC4Layout(g, th, cfg)
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
	width, height := normalizeCoordinates(nodes, edges)
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

// normalizeCoordinates translates all node and edge positions so that the
// minimum coordinates are at layoutBoundaryPad, ensuring no content is clipped
// by the SVG viewBox "0 0 W H". Returns the final canvas width and height.
func normalizeCoordinates(nodes map[string]*NodeLayout, edges []*EdgeLayout) (float32, float32) {
	if len(nodes) == 0 {
		return 0, 0
	}

	// Find the bounding box of all content.
	var minX, minY float32
	var maxX, maxY float32
	first := true

	expandBounds := func(left, top, right, bottom float32) {
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

	for _, n := range nodes {
		expandBounds(n.X-n.Width/2, n.Y-n.Height/2, n.X+n.Width/2, n.Y+n.Height/2)
	}

	for _, e := range edges {
		for _, pt := range e.Points {
			expandBounds(pt[0], pt[1], pt[0], pt[1])
		}
	}

	// Compute the shift needed so that minX/minY become layoutBoundaryPad.
	dx := layoutBoundaryPad - minX
	dy := layoutBoundaryPad - minY

	// Translate all nodes.
	for _, n := range nodes {
		n.X += dx
		n.Y += dy
	}

	// Translate all edge points and label anchors.
	for _, e := range edges {
		for i := range e.Points {
			e.Points[i][0] += dx
			e.Points[i][1] += dy
		}
		e.LabelAnchor[0] += dx
		e.LabelAnchor[1] += dy
	}

	width := (maxX - minX) + 2*layoutBoundaryPad
	height := (maxY - minY) + 2*layoutBoundaryPad

	return width, height
}
