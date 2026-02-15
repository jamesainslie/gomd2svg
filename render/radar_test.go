package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderRadar(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Radar
	g.RadarTitle = "Skills"
	g.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "Speed"},
		{ID: "b", Label: "Power"},
		{ID: "c", Label: "Magic"},
	}
	g.RadarCurves = []*ir.RadarCurve{
		{ID: "p1", Label: "Player1", Values: []float64{80, 60, 40}},
	}
	g.RadarMax = 100

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
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
	g := ir.NewGraph()
	g.Kind = ir.Radar
	g.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
		{ID: "c", Label: "C"},
	}
	g.RadarCurves = []*ir.RadarCurve{
		{ID: "x", Label: "X", Values: []float64{50, 50, 50}},
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
