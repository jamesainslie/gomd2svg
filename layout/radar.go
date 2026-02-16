package layout

import (
	"math"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeRadarLayout(g *ir.Graph, _ *theme.Theme, cfg *config.Layout) *Layout {
	radius := cfg.Radar.Radius
	padX := cfg.Radar.PaddingX
	padY := cfg.Radar.PaddingY
	labelOffset := cfg.Radar.LabelOffset

	numAxes := len(g.RadarAxes)
	if numAxes == 0 {
		numAxes = 1
	}

	// Title height.
	var titleHeight float32
	if g.RadarTitle != "" {
		titleHeight = 20 + padY
	}

	centerX := padX + radius + labelOffset
	centerY := titleHeight + padY + radius + labelOffset

	// Determine max value.
	maxVal := g.RadarMax
	if maxVal <= 0 {
		maxVal = radarAutoMax(g)
	}
	minVal := g.RadarMin

	valRange := maxVal - minVal
	if valRange <= 0 {
		valRange = 1
	}

	// Angle per axis (evenly distributed, starting from top = -pi/2).
	angleStep := 2 * math.Pi / float64(numAxes)

	// Build axis layouts.
	axes := make([]RadarAxisLayout, len(g.RadarAxes))
	for i, ax := range g.RadarAxes {
		angle := -math.Pi/2 + float64(i)*angleStep
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))
		axes[i] = RadarAxisLayout{
			Label:  ax.Label,
			EndX:   centerX + radius*cos,
			EndY:   centerY + radius*sin,
			LabelX: centerX + (radius+labelOffset)*cos,
			LabelY: centerY + (radius+labelOffset)*sin,
		}
	}

	// Build curve layouts.
	curves := make([]RadarCurveLayout, len(g.RadarCurves))
	for ci, curve := range g.RadarCurves {
		var points [][2]float32
		for i, v := range curve.Values {
			if i >= len(g.RadarAxes) {
				break
			}
			angle := -math.Pi/2 + float64(i)*angleStep
			frac := (v - minVal) / valRange
			if frac < 0 {
				frac = 0
			}
			if frac > 1 {
				frac = 1
			}
			r := float32(frac) * radius
			px := centerX + r*float32(math.Cos(angle))
			py := centerY + r*float32(math.Sin(angle))
			points = append(points, [2]float32{px, py})
		}
		curves[ci] = RadarCurveLayout{
			Label:      curve.Label,
			Points:     points,
			ColorIndex: ci,
		}
	}

	// Graticule radii (concentric rings).
	ticks := g.RadarTicks
	if ticks <= 0 {
		ticks = cfg.Radar.DefaultTicks
	}
	graticuleRadii := make([]float32, ticks)
	for i := range ticks {
		graticuleRadii[i] = radius * float32(i+1) / float32(ticks)
	}

	totalW := (padX + labelOffset + radius) * 2
	totalH := titleHeight + (padY+labelOffset+radius)*2

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: RadarData{
			Axes:           axes,
			Curves:         curves,
			GraticuleRadii: graticuleRadii,
			GraticuleType:  g.RadarGraticuleType,
			CenterX:        centerX,
			CenterY:        centerY,
			Radius:         radius,
			Title:          g.RadarTitle,
			ShowLegend:     g.RadarShowLegend,
			MaxValue:       maxVal,
			MinValue:       minVal,
		},
	}
}

func radarAutoMax(g *ir.Graph) float64 {
	max := 0.0
	for _, c := range g.RadarCurves {
		for _, v := range c.Values {
			if v > max {
				max = v
			}
		}
	}
	if max <= 0 {
		return 100
	}
	return niceMax(max)
}
