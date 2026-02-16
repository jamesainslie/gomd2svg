package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeQuadrantLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	chartW := cfg.Quadrant.ChartWidth
	chartH := cfg.Quadrant.ChartHeight
	padX := cfg.Quadrant.PaddingX
	padY := cfg.Quadrant.PaddingY

	// Title height.
	var titleHeight float32
	if g.QuadrantTitle != "" {
		titleHeight = th.FontSize + padY
	}

	// Y-axis label width (left side).
	var yAxisLabelWidth float32
	if g.YAxisBottom != "" || g.YAxisTop != "" {
		yAxisLabelWidth = padX
	}

	// Chart origin.
	chartX := padX + yAxisLabelWidth
	chartY := titleHeight + padY

	// X-axis label height (below chart).
	var xAxisLabelHeight float32
	if g.XAxisLeft != "" || g.XAxisRight != "" {
		xAxisLabelHeight = cfg.Quadrant.AxisLabelFontSize + padY/2
	}

	// Map normalized points to pixel positions.
	points := make([]QuadrantPointLayout, len(g.QuadrantPoints))
	for i, p := range g.QuadrantPoints {
		// X maps directly: 0 = left, 1 = right.
		px := chartX + float32(p.X)*chartW
		// Y is inverted: 0 = bottom (high Y), 1 = top (low Y).
		py := chartY + (1-float32(p.Y))*chartH
		points[i] = QuadrantPointLayout{
			Label: p.Label,
			X:     px,
			Y:     py,
		}
	}

	totalW := padX + yAxisLabelWidth + chartW + padX
	totalH := titleHeight + padY + chartH + xAxisLabelHeight + padY

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: QuadrantData{
			Points:      points,
			ChartX:      chartX,
			ChartY:      chartY,
			ChartWidth:  chartW,
			ChartHeight: chartH,
			Title:       g.QuadrantTitle,
			Labels:      g.QuadrantLabels,
			XAxisLeft:   g.XAxisLeft,
			XAxisRight:  g.XAxisRight,
			YAxisBottom: g.YAxisBottom,
			YAxisTop:    g.YAxisTop,
		},
	}
}
