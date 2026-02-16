package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderMindmap(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Mindmap
	graph.MindmapRoot = &ir.MindmapNode{
		ID: "root", Label: "Central", Shape: ir.MindmapCircle,
		Children: []*ir.MindmapNode{
			{ID: "a", Label: "Square", Shape: ir.MindmapSquare},
			{ID: "b", Label: "Rounded", Shape: ir.MindmapRounded},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Central") {
		t.Error("missing root label")
	}
	if !strings.Contains(svg, "Square") {
		t.Error("missing child label")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing circle for root node")
	}
}

func TestRenderMindmapEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Mindmap

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
