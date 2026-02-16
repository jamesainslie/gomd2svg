package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderRadar(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Radar
	graph.RadarTitle = "Skills"
	graph.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "Speed"},
		{ID: "b", Label: "Power"},
		{ID: "c", Label: "Magic"},
	}
	graph.RadarCurves = []*ir.RadarCurve{
		{ID: "p1", Label: "Player1", Values: []float64{80, 60, 40}},
	}
	graph.RadarMax = 100

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Skills") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "<polygon") || !strings.Contains(svg, "<line") {
		t.Error("missing radar elements (polygon or axis lines)")
	}
}

func TestRenderRadarValidSVG(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Radar
	graph.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
		{ID: "c", Label: "C"},
	}
	graph.RadarCurves = []*ir.RadarCurve{
		{ID: "x", Label: "X", Values: []float64{50, 50, 50}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG doesn't start with <svg")
	}
	if !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
		t.Error("SVG doesn't end with </svg>")
	}
}
