package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderBlock(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Block
	g.BlockColumns = 2

	// Add 3 block nodes.
	for _, id := range []string{"a", "b", "c"} {
		label := strings.ToUpper(id)
		g.EnsureNode(id, &label, nil)
		g.Blocks = append(g.Blocks, &ir.BlockDef{ID: id, Label: label, Width: 1})
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
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
	g := ir.NewGraph()
	g.Kind = ir.Block

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
