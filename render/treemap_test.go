package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderTreemap(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Treemap
	g.TreemapTitle = "Budget"
	g.TreemapRoot = &ir.TreemapNode{
		Label: "Root",
		Children: []*ir.TreemapNode{
			{Label: "Salaries", Value: 70},
			{Label: "Equipment", Value: 30},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Budget") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Salaries") {
		t.Error("missing leaf label")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing rectangles")
	}
}

func TestRenderTreemapEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Treemap

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
