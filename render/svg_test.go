package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func simpleLayout() *layout.Layout {
	graph := ir.NewGraph()
	graph.Kind = ir.Flowchart
	graph.Direction = ir.LeftRight
	graph.EnsureNode("A", nil, nil)
	graph.EnsureNode("B", nil, nil)
	graph.Edges = []*ir.Edge{{
		From: "A", To: "B", Directed: true, ArrowEnd: true, Style: ir.Solid,
	}}
	th := theme.Modern()
	cfg := config.DefaultLayout()
	return layout.ComputeLayout(graph, th, cfg)
}

func TestRenderSVGContainsSVGTags(t *testing.T) {
	l := simpleLayout()
	th := theme.Modern()
	cfg := config.DefaultLayout()
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("missing </svg> tag")
	}
}

func TestRenderSVGContainsNodes(t *testing.T) {
	l := simpleLayout()
	th := theme.Modern()
	cfg := config.DefaultLayout()
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, "<rect") {
		t.Error("missing <rect for node shapes")
	}
	if !strings.Contains(svg, "A") {
		t.Error("missing node label A")
	}
	if !strings.Contains(svg, "B") {
		t.Error("missing node label B")
	}
}

func TestRenderSVGContainsEdge(t *testing.T) {
	l := simpleLayout()
	th := theme.Modern()
	cfg := config.DefaultLayout()
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, "<path") || !strings.Contains(svg, "edgePath") {
		t.Error("missing edge path")
	}
}

func TestRenderDefsHasAllMarkers(t *testing.T) {
	th := theme.Modern()
	var b svgBuilder
	renderDefs(&b, th)
	svg := b.String()

	markers := []string{
		"arrowhead", "arrowhead-start",
		"marker-closed-triangle", "marker-closed-triangle-start",
		"marker-filled-diamond", "marker-filled-diamond-start",
		"marker-open-diamond", "marker-open-diamond-start",
		"marker-open-arrow", "marker-cross",
	}
	for _, id := range markers {
		if !strings.Contains(svg, `id="`+id+`"`) {
			t.Errorf("missing marker definition: %s", id)
		}
	}
}

func TestRenderSVGHasViewBox(t *testing.T) {
	l := simpleLayout()
	th := theme.Modern()
	cfg := config.DefaultLayout()
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, "viewBox") {
		t.Error("missing viewBox attribute")
	}
}
