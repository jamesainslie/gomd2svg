package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderBlock(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Block
	graph.BlockColumns = 2

	// Add 3 block nodes.
	for _, id := range []string{"a", "b", "c"} {
		label := strings.ToUpper(id)
		graph.EnsureNode(id, &label, nil)
		graph.Blocks = append(graph.Blocks, &ir.BlockDef{ID: id, Label: label, Width: 1})
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing <rect elements")
	}
	for _, label := range []string{"A", "B", "C"} {
		if !strings.Contains(svg, label) {
			t.Errorf("missing node label %q", label)
		}
	}
}

func TestRenderBlockEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Block

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
