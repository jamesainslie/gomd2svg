package layout

import (
	"math"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Radar layout constants.
const (
	radarTitleHeight  float32 = 20
	radarDefaultValue float64 = 100
)

func computeRadarLayout(graph *ir.Graph, _ *theme.Theme, cfg *config.Layout) *Layout {
	radius := cfg.Radar.Radius
	padX := cfg.Radar.PaddingX
	padY := cfg.Radar.PaddingY
	labelOffset := cfg.Radar.LabelOffset

	numAxes := len(graph.RadarAxes)
	if numAxes == 0 {
		numAxes = 1
	}

	// Title height.
	var titleHeight float32
	if graph.RadarTitle != "" {
		titleHeight = radarTitleHeight + padY
	}

	centerX := padX + radius + labelOffset
	centerY := titleHeight + padY + radius + labelOffset

	// Determine max value.
	maxVal := graph.RadarMax
	if maxVal <= 0 {
		maxVal = radarAutoMax(graph)
	}
	minVal := graph.RadarMin

	valRange := maxVal - minVal
	if valRange <= 0 {
		valRange = 1
	}

	// Angle per axis (evenly distributed, starting from top = -pi/2).
	angleStep := 2 * math.Pi / float64(numAxes)

	// Build axis layouts.
	axes := make([]RadarAxisLayout, len(graph.RadarAxes))
	for i, ax := range graph.RadarAxes {
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
	curves := make([]RadarCurveLayout, len(graph.RadarCurves))
	for ci, curve := range graph.RadarCurves {
		var points [][2]float32
		for i, v := range curve.Values {
			if i >= len(graph.RadarAxes) {
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
	ticks := graph.RadarTicks
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
		Kind:   graph.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: RadarData{
			Axes:           axes,
			Curves:         curves,
			GraticuleRadii: graticuleRadii,
			GraticuleType:  graph.RadarGraticuleType,
			CenterX:        centerX,
			CenterY:        centerY,
			Radius:         radius,
			Title:          graph.RadarTitle,
			ShowLegend:     graph.RadarShowLegend,
			MaxValue:       maxVal,
			MinValue:       minVal,
		},
	}
}

func radarAutoMax(graph *ir.Graph) float64 {
	maxValue := 0.0
	for _, c := range graph.RadarCurves {
		for _, v := range c.Values {
			if v > maxValue {
				maxValue = v
			}
		}
	}
	if maxValue <= 0 {
		return radarDefaultValue
	}
	return niceMax(maxValue)
}
