package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeQuadrantLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	chartW := cfg.Quadrant.ChartWidth
	chartH := cfg.Quadrant.ChartHeight
	padX := cfg.Quadrant.PaddingX
	padY := cfg.Quadrant.PaddingY

	// Title height.
	var titleHeight float32
	if graph.QuadrantTitle != "" {
		titleHeight = th.FontSize + padY
	}

	// Y-axis label width (left side).
	var yAxisLabelWidth float32
	if graph.YAxisBottom != "" || graph.YAxisTop != "" {
		yAxisLabelWidth = padX
	}

	// Chart origin.
	chartX := padX + yAxisLabelWidth
	chartY := titleHeight + padY

	// X-axis label height (below chart).
	var xAxisLabelHeight float32
	if graph.XAxisLeft != "" || graph.XAxisRight != "" {
		xAxisLabelHeight = cfg.Quadrant.AxisLabelFontSize + padY/2
	}

	// Map normalized points to pixel positions.
	points := make([]QuadrantPointLayout, len(graph.QuadrantPoints))
	for idx, pt := range graph.QuadrantPoints {
		// X maps directly: 0 = left, 1 = right.
		px := chartX + float32(pt.X)*chartW
		// Y is inverted: 0 = bottom (high Y), 1 = top (low Y).
		py := chartY + (1-float32(pt.Y))*chartH
		points[idx] = QuadrantPointLayout{
			Label: pt.Label,
			X:     px,
			Y:     py,
		}
	}

	totalW := padX + yAxisLabelWidth + chartW + padX
	totalH := titleHeight + padY + chartH + xAxisLabelHeight + padY

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: QuadrantData{
			Points:      points,
			ChartX:      chartX,
			ChartY:      chartY,
			ChartWidth:  chartW,
			ChartHeight: chartH,
			Title:       graph.QuadrantTitle,
			Labels:      graph.QuadrantLabels,
			XAxisLeft:   graph.XAxisLeft,
			XAxisRight:  graph.XAxisRight,
			YAxisBottom: graph.YAxisBottom,
			YAxisTop:    graph.YAxisTop,
		},
	}
}
