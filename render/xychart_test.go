package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderXYChart(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart
	g.XYTitle = "Revenue"
	g.XYXAxis = &ir.XYAxis{
		Mode:       ir.XYAxisBand,
		Categories: []string{"Q1", "Q2", "Q3"},
	}
	g.XYYAxis = &ir.XYAxis{Mode: ir.XYAxisNumeric, Min: 0, Max: 100}
	g.XYSeries = []*ir.XYSeries{
		{Type: ir.XYSeriesBar, Values: []float64{30, 60, 90}},
		{Type: ir.XYSeriesLine, Values: []float64{25, 55, 85}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Revenue") {
		t.Error("missing title text")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing bar rects")
	}
	if !strings.Contains(svg, "<polyline") || !strings.Contains(svg, "<circle") {
		t.Error("missing line series elements")
	}
}

func TestRenderXYChartValidSVG(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart
	g.XYSeries = []*ir.XYSeries{
		{Type: ir.XYSeriesLine, Values: []float64{1, 2, 3}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG doesn't start with <svg")
	}
	if !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
		t.Error("SVG doesn't end with </svg>")
	}
}
