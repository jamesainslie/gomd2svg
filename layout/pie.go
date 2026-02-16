package layout

import (
	"math"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computePieLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	radius := cfg.Pie.Radius
	padX := cfg.Pie.PaddingX
	padY := cfg.Pie.PaddingY
	textPos := cfg.Pie.TextPosition

	// Compute total value.
	var total float64
	for _, s := range g.PieSlices {
		total += s.Value
	}
	if total <= 0 {
		total = 1
	}

	// Title height.
	var titleHeight float32
	if g.PieTitle != "" {
		titleHeight = th.PieTitleTextSize + padY
	}

	centerX := padX + radius
	centerY := titleHeight + padY + radius

	// Compute slice angles (clockwise from top = -pi/2).
	slices := make([]PieSliceLayout, len(g.PieSlices))
	var angle float32 = -math.Pi / 2 // start at top

	for i, s := range g.PieSlices {
		frac := float32(s.Value / total)
		span := frac * 2 * math.Pi

		midAngle := angle + span/2
		labelR := radius * textPos
		labelX := centerX + labelR*float32(math.Cos(float64(midAngle)))
		labelY := centerY + labelR*float32(math.Sin(float64(midAngle)))

		slices[i] = PieSliceLayout{
			Label:      s.Label,
			Value:      s.Value,
			Percentage: frac * 100,
			StartAngle: angle,
			EndAngle:   angle + span,
			LabelX:     labelX,
			LabelY:     labelY,
			ColorIndex: i,
		}

		angle += span
	}

	width := 2*padX + 2*radius
	height := titleHeight + 2*padY + 2*radius

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  width,
		Height: height,
		Diagram: PieData{
			Slices:   slices,
			CenterX:  centerX,
			CenterY:  centerY,
			Radius:   radius,
			Title:    g.PieTitle,
			ShowData: g.PieShowData,
		},
	}
}
