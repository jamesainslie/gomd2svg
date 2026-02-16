package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestSVGHasRoleImg(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Flowchart
	label := "Start"
	g.EnsureNode("A", &label, nil)
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, `role="img"`) {
		t.Error("missing role=img")
	}
}

func TestSVGFallbackAriaLabel(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Flowchart
	label := "Start"
	g.EnsureNode("A", &label, nil)
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, `aria-label="Flowchart diagram"`) {
		t.Error("missing fallback aria-label")
	}
}

func TestSVGTitleElement(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Pie
	g.PieTitle = "Browser Market Share"
	g.PieSlices = append(g.PieSlices, &ir.PieSlice{Label: "Chrome", Value: 60})
	g.PieSlices = append(g.PieSlices, &ir.PieSlice{Label: "Firefox", Value: 40})
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)
	if !strings.Contains(svg, "<title>Browser Market Share</title>") {
		t.Error("missing <title> element")
	}
	if !strings.Contains(svg, `aria-label="Browser Market Share"`) {
		t.Error("missing aria-label with title")
	}
}

func TestSVGNoTitleNoElement(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Flowchart
	label := "Start"
	g.EnsureNode("A", &label, nil)
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)
	if strings.Contains(svg, "<title>") {
		t.Error("<title> should not be present when no diagram title is set")
	}
}
