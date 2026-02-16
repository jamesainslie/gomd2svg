package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderSankey(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Sankey
	graph.SankeyLinks = []*ir.SankeyLink{
		{Source: "Solar", Target: "Grid", Value: 60},
		{Source: "Wind", Target: "Grid", Value: 290},
		{Source: "Grid", Target: "Industry", Value: 350},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Solar") {
		t.Error("missing node label Solar")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing node rects")
	}
	if !strings.Contains(svg, "<path") {
		t.Error("missing link paths")
	}
}

func TestRenderSankeyEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Sankey

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
