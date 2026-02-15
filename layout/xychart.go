package layout

import (
	"fmt"
	"math"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeXYChartLayout(g *ir.Graph, _ *theme.Theme, cfg *config.Layout) *Layout {
	padX := cfg.XYChart.PaddingX
	padY := cfg.XYChart.PaddingY
	chartW := cfg.XYChart.ChartWidth
	chartH := cfg.XYChart.ChartHeight

	// Title height.
	var titleHeight float32
	if g.XYTitle != "" {
		titleHeight = cfg.XYChart.TitleFontSize + padY
	}

	// Determine Y range.
	yMin, yMax := xyDataRange(g)
	if g.XYYAxis != nil && (g.XYYAxis.Min != 0 || g.XYYAxis.Max != 0) {
		yMin = g.XYYAxis.Min
		yMax = g.XYYAxis.Max
	}
	if yMax <= yMin {
		yMax = yMin + 1
	}
	yRange := yMax - yMin

	chartX := padX
	chartY := titleHeight + padY

	// valueToY maps a data value to a pixel Y position.
	valueToY := func(v float64) float32 {
		frac := (v - yMin) / yRange
		return chartY + chartH - float32(frac)*chartH
	}

	// Generate Y-axis ticks.
	yTicks := generateYTicks(yMin, yMax, 5, chartY, chartH)

	// Generate X-axis labels and series points.
	var xLabels []XYAxisLabel
	numPoints := xyMaxPoints(g)
	if numPoints == 0 {
		numPoints = 1
	}

	bandW := chartW / float32(numPoints)

	// X-axis labels from categories or numeric range.
	if g.XYXAxis != nil && g.XYXAxis.Mode == ir.XYAxisBand {
		for i, cat := range g.XYXAxis.Categories {
			xLabels = append(xLabels, XYAxisLabel{
				Text: cat,
				X:    chartX + float32(i)*bandW + bandW/2,
			})
		}
	} else {
		for i := range numPoints {
			xLabels = append(xLabels, XYAxisLabel{
				Text: fmt.Sprintf("%d", i+1),
				X:    chartX + float32(i)*bandW + bandW/2,
			})
		}
	}

	// Count bar series for grouping.
	barCount := 0
	for _, s := range g.XYSeries {
		if s.Type == ir.XYSeriesBar {
			barCount++
		}
	}

	barGroupWidth := bandW * cfg.XYChart.BarWidth
	var singleBarW float32
	if barCount > 0 {
		singleBarW = barGroupWidth / float32(barCount)
	}

	// Build series layouts.
	barIdx := 0
	var series []XYSeriesLayout
	for si, s := range g.XYSeries {
		var points []XYPointLayout
		for i, v := range s.Values {
			cx := chartX + float32(i)*bandW + bandW/2
			py := valueToY(v)

			switch s.Type {
			case ir.XYSeriesBar:
				barX := cx - barGroupWidth/2 + float32(barIdx)*singleBarW
				baseY := valueToY(math.Max(yMin, 0))
				h := baseY - py
				if h < 0 {
					h = -h
					py = baseY
				}
				points = append(points, XYPointLayout{
					X: barX, Y: py, Width: singleBarW, Height: h, Value: v,
				})
			case ir.XYSeriesLine:
				points = append(points, XYPointLayout{
					X: cx, Y: py, Value: v,
				})
			}
		}
		if s.Type == ir.XYSeriesBar {
			barIdx++
		}
		series = append(series, XYSeriesLayout{
			Type:       s.Type,
			Points:     points,
			ColorIndex: si,
		})
	}

	totalW := padX*2 + chartW
	totalH := titleHeight + padY*2 + chartH + cfg.XYChart.AxisFontSize + padY

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: XYChartData{
			Series:      series,
			XLabels:     xLabels,
			YTicks:      yTicks,
			Title:       g.XYTitle,
			ChartX:      chartX,
			ChartY:      chartY,
			ChartWidth:  chartW,
			ChartHeight: chartH,
			YMin:        yMin,
			YMax:        yMax,
			Horizontal:  g.XYHorizontal,
		},
	}
}

func xyDataRange(g *ir.Graph) (float64, float64) {
	min, max := math.Inf(1), math.Inf(-1)
	for _, s := range g.XYSeries {
		for _, v := range s.Values {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}
	if math.IsInf(min, 1) {
		min = 0
	}
	if math.IsInf(max, -1) {
		max = 100
	}
	if min > 0 {
		min = 0
	}
	// Round max up to a nice number.
	max = niceMax(max)
	return min, max
}

func niceMax(v float64) float64 {
	if v <= 0 {
		return 1
	}
	magnitude := math.Pow(10, math.Floor(math.Log10(v)))
	normalized := v / magnitude
	if normalized <= 1 {
		return magnitude
	} else if normalized <= 2 {
		return 2 * magnitude
	} else if normalized <= 5 {
		return 5 * magnitude
	}
	return 10 * magnitude
}

func generateYTicks(yMin, yMax float64, count int, chartY, chartH float32) []XYAxisTick {
	yRange := yMax - yMin
	step := yRange / float64(count)
	var ticks []XYAxisTick
	for i := range count + 1 {
		v := yMin + float64(i)*step
		frac := (v - yMin) / yRange
		y := chartY + chartH - float32(frac)*chartH
		ticks = append(ticks, XYAxisTick{
			Label: fmt.Sprintf("%.4g", v),
			Y:     y,
		})
	}
	return ticks
}

func xyMaxPoints(g *ir.Graph) int {
	max := 0
	for _, s := range g.XYSeries {
		if len(s.Values) > max {
			max = len(s.Values)
		}
	}
	return max
}
