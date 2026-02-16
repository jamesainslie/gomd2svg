package layout

import (
	"math"
	"strconv"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

// XY chart layout constants.
const (
	xyChartDefaultTickCount = 5
	xyChartTickBase         = 10
	xyChartNiceStepFive     = 5
)

func computeXYChartLayout(graph *ir.Graph, _ *theme.Theme, cfg *config.Layout) *Layout {
	padX := cfg.XYChart.PaddingX
	padY := cfg.XYChart.PaddingY
	chartW := cfg.XYChart.ChartWidth
	chartH := cfg.XYChart.ChartHeight

	// Title height.
	var titleHeight float32
	if graph.XYTitle != "" {
		titleHeight = cfg.XYChart.TitleFontSize + padY
	}

	// Determine Y range.
	yMin, yMax := xyDataRange(graph)
	if graph.XYYAxis != nil && (graph.XYYAxis.Min != 0 || graph.XYYAxis.Max != 0) {
		yMin = graph.XYYAxis.Min
		yMax = graph.XYYAxis.Max
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
	yTicks := generateYTicks(yMin, yMax, xyChartDefaultTickCount, chartY, chartH)

	// Generate X-axis labels and series points.
	var xLabels []XYAxisLabel
	numPoints := xyMaxPoints(graph)
	if numPoints == 0 {
		numPoints = 1
	}

	bandW := chartW / float32(numPoints)

	// X-axis labels from categories or numeric range.
	if graph.XYXAxis != nil && graph.XYXAxis.Mode == ir.XYAxisBand {
		for i, cat := range graph.XYXAxis.Categories {
			xLabels = append(xLabels, XYAxisLabel{
				Text: cat,
				X:    chartX + float32(i)*bandW + bandW/2,
			})
		}
	} else {
		for i := range numPoints {
			xLabels = append(xLabels, XYAxisLabel{
				Text: strconv.Itoa(i + 1),
				X:    chartX + float32(i)*bandW + bandW/2,
			})
		}
	}

	// Count bar series for grouping.
	barCount := 0
	for _, series := range graph.XYSeries {
		if series.Type == ir.XYSeriesBar {
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
	seriesLayouts := make([]XYSeriesLayout, 0, len(graph.XYSeries))
	for si, xySeries := range graph.XYSeries {
		var points []XYPointLayout
		for i, val := range xySeries.Values {
			cx := chartX + float32(i)*bandW + bandW/2
			py := valueToY(val)

			switch xySeries.Type {
			case ir.XYSeriesBar:
				barX := cx - barGroupWidth/2 + float32(barIdx)*singleBarW
				baseY := valueToY(math.Max(yMin, 0))
				barHeight := baseY - py
				if barHeight < 0 {
					barHeight = -barHeight
					py = baseY
				}
				points = append(points, XYPointLayout{
					X: barX, Y: py, Width: singleBarW, Height: barHeight, Value: val,
				})
			case ir.XYSeriesLine:
				points = append(points, XYPointLayout{
					X: cx, Y: py, Value: val,
				})
			}
		}
		if xySeries.Type == ir.XYSeriesBar {
			barIdx++
		}
		seriesLayouts = append(seriesLayouts, XYSeriesLayout{
			Type:       xySeries.Type,
			Points:     points,
			ColorIndex: si,
		})
	}

	totalW := padX*2 + chartW
	totalH := titleHeight + padY*2 + chartH + cfg.XYChart.AxisFontSize + padY

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: XYChartData{
			Series:      seriesLayouts,
			XLabels:     xLabels,
			YTicks:      yTicks,
			Title:       graph.XYTitle,
			ChartX:      chartX,
			ChartY:      chartY,
			ChartWidth:  chartW,
			ChartHeight: chartH,
			YMin:        yMin,
			YMax:        yMax,
			Horizontal:  graph.XYHorizontal,
		},
	}
}

func xyDataRange(graph *ir.Graph) (float64, float64) {
	minVal, maxVal := math.Inf(1), math.Inf(-1)
	for _, xySeries := range graph.XYSeries {
		for _, val := range xySeries.Values {
			if val < minVal {
				minVal = val
			}
			if val > maxVal {
				maxVal = val
			}
		}
	}
	if math.IsInf(minVal, 1) {
		minVal = 0
	}
	if math.IsInf(maxVal, -1) {
		maxVal = 100
	}
	if minVal > 0 {
		minVal = 0
	}
	// Round max up to a nice number.
	maxVal = niceMax(maxVal)
	return minVal, maxVal
}

func niceMax(v float64) float64 {
	if v <= 0 {
		return 1
	}
	magnitude := math.Pow(xyChartTickBase, math.Floor(math.Log10(v)))
	normalized := v / magnitude
	switch {
	case normalized <= 1:
		return magnitude
	case normalized <= 2:
		return 2 * magnitude
	case normalized <= xyChartNiceStepFive:
		return xyChartNiceStepFive * magnitude
	default:
		return xyChartTickBase * magnitude
	}
}

func generateYTicks(yMin, yMax float64, count int, chartY, chartH float32) []XYAxisTick {
	yRange := yMax - yMin
	step := yRange / float64(count)
	ticks := make([]XYAxisTick, 0, count+1)
	for i := range count + 1 {
		v := yMin + float64(i)*step
		frac := (v - yMin) / yRange
		y := chartY + chartH - float32(frac)*chartH
		ticks = append(ticks, XYAxisTick{
			Label: strconv.FormatFloat(v, 'g', 4, 64),
			Y:     y,
		})
	}
	return ticks
}

func xyMaxPoints(graph *ir.Graph) int {
	maxPoints := 0
	for _, xySeries := range graph.XYSeries {
		if len(xySeries.Values) > maxPoints {
			maxPoints = len(xySeries.Values)
		}
	}
	return maxPoints
}
