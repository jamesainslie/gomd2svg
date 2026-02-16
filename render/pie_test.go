package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderPie(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Pie
	g.PieTitle = "Pets"
	g.PieSlices = []*ir.PieSlice{
		{Label: "Dogs", Value: 386},
		{Label: "Cats", Value: 85},
		{Label: "Rats", Value: 15},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "Pets") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Dogs") {
		t.Error("missing slice label Dogs")
	}
	if !strings.Contains(svg, "Cats") {
		t.Error("missing slice label Cats")
	}
	// Should contain arc paths.
	if !strings.Contains(svg, "<path") {
		t.Error("missing <path> for arcs")
	}
}

func TestRenderPieShowData(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Pie
	g.PieShowData = true
	g.PieSlices = []*ir.PieSlice{
		{Label: "A", Value: 60},
		{Label: "B", Value: 40},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// When showData is true, values should appear in the output.
	if !strings.Contains(svg, "60") {
		t.Error("missing value 60 with showData=true")
	}
}

func TestRenderPieSingleSlice(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Pie
	g.PieSlices = []*ir.PieSlice{
		{Label: "All", Value: 100},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Single slice = full circle, should use <circle> element.
	if !strings.Contains(svg, "<circle") {
		t.Error("missing <circle> for single slice")
	}
	if !strings.Contains(svg, "All") {
		t.Error("missing label for single slice")
	}
}
