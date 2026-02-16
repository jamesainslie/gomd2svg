package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

// ComputeLayout dispatches to the appropriate layout algorithm based on
// the diagram kind and sets the diagram title on the result.
func ComputeLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	lay := computeLayout(graph, th, cfg)
	lay.Title = diagramTitle(graph)
	return lay
}

// diagramTitle extracts the title from diagram-kind-specific fields on the graph.
func diagramTitle(graph *ir.Graph) string {
	switch graph.Kind {
	case ir.Pie:
		return graph.PieTitle
	case ir.Quadrant:
		return graph.QuadrantTitle
	case ir.Timeline:
		return graph.TimelineTitle
	case ir.Gantt:
		return graph.GanttTitle
	case ir.XYChart:
		return graph.XYTitle
	case ir.Radar:
		return graph.RadarTitle
	case ir.Treemap:
		return graph.TreemapTitle
	case ir.Journey:
		return graph.JourneyTitle
	default:
		return ""
	}
}

// computeLayout dispatches to the appropriate layout algorithm based on
// the diagram kind.
func computeLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	switch graph.Kind {
	case ir.Flowchart:
		return computeGraphLayout(graph, th, cfg)
	case ir.Class:
		return computeClassLayout(graph, th, cfg)
	case ir.Er:
		return computeERLayout(graph, th, cfg)
	case ir.State:
		return computeStateLayout(graph, th, cfg)
	case ir.Sequence:
		return computeSequenceLayout(graph, th, cfg)
	case ir.Kanban:
		return computeKanbanLayout(graph, th, cfg)
	case ir.Packet:
		return computePacketLayout(graph, th, cfg)
	case ir.Pie:
		return computePieLayout(graph, th, cfg)
	case ir.Quadrant:
		return computeQuadrantLayout(graph, th, cfg)
	case ir.Timeline:
		return computeTimelineLayout(graph, th, cfg)
	case ir.Gantt:
		return computeGanttLayout(graph, th, cfg)
	case ir.GitGraph:
		return computeGitGraphLayout(graph, th, cfg)
	case ir.XYChart:
		return computeXYChartLayout(graph, th, cfg)
	case ir.Radar:
		return computeRadarLayout(graph, th, cfg)
	case ir.Mindmap:
		return computeMindmapLayout(graph, th, cfg)
	case ir.Sankey:
		return computeSankeyLayout(graph, th, cfg)
	case ir.Treemap:
		return computeTreemapLayout(graph, th, cfg)
	case ir.Requirement:
		return computeRequirementLayout(graph, th, cfg)
	case ir.Block:
		return computeBlockLayout(graph, th, cfg)
	case ir.C4:
		return computeC4Layout(graph, th, cfg)
	case ir.Journey:
		return computeJourneyLayout(graph, th, cfg)
	case ir.Architecture:
		return computeArchitectureLayout(graph, th, cfg)
	case ir.ZenUML:
		return computeSequenceLayout(graph, th, cfg)
	default:
		// For unsupported diagram kinds, return a minimal layout.
		return computeGraphLayout(graph, th, cfg)
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
func runSugiyama(graph *ir.Graph, nodes map[string]*NodeLayout, cfg *config.Layout) sugiyamaResult {
	nodeIDs := sortedNodeIDs(graph.Nodes, graph.NodeOrder)
	ranks := computeRanks(nodeIDs, graph.Edges, graph.NodeOrder)
	layers := orderRankNodes(ranks, graph.Edges, cfg.Flowchart.OrderPasses)
	positionNodes(layers, nodes, graph.Direction, cfg)
	edges := routeEdges(graph.Edges, nodes, graph.Direction)
	width, height := normalizeCoordinates(nodes, edges)
	return sugiyamaResult{Edges: edges, Width: width, Height: height}
}

// computeGraphLayout runs the full Sugiyama-style layout pipeline:
// 1. Size nodes based on text metrics.
// 2. Run Sugiyama ranking, ordering, positioning, routing, and bounding box.
func computeGraphLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()

	// Step 1: Size all nodes.
	nodes := sizeNodes(graph.Nodes, measurer, th, cfg)

	// Step 2: Run Sugiyama pipeline.
	result := runSugiyama(graph, nodes, cfg)

	return &Layout{
		Kind:    graph.Kind,
		Nodes:   nodes,
		Edges:   result.Edges,
		Width:   result.Width,
		Height:  result.Height,
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

	for _, node := range nodes {
		expandBounds(node.X-node.Width/2, node.Y-node.Height/2, node.X+node.Width/2, node.Y+node.Height/2)
	}

	for _, edge := range edges {
		for _, pt := range edge.Points {
			expandBounds(pt[0], pt[1], pt[0], pt[1])
		}
	}

	// Compute the shift needed so that minX/minY become layoutBoundaryPad.
	dx := layoutBoundaryPad - minX
	dy := layoutBoundaryPad - minY

	// Translate all nodes.
	for _, node := range nodes {
		node.X += dx
		node.Y += dy
	}

	// Translate all edge points and label anchors.
	for _, edge := range edges {
		for idx := range edge.Points {
			edge.Points[idx][0] += dx
			edge.Points[idx][1] += dy
		}
		edge.LabelAnchor[0] += dx
		edge.LabelAnchor[1] += dy
	}

	width := (maxX - minX) + 2*layoutBoundaryPad
	height := (maxY - minY) + 2*layoutBoundaryPad

	return width, height
}
