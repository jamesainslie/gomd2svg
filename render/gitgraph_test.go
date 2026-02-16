package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderGitGraph(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.GitGraph
	graph.GitMainBranch = "main"
	graph.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1", Tag: "v1.0"},
		&ir.GitBranch{Name: "develop"},
		&ir.GitCheckout{Branch: "develop"},
		&ir.GitCommit{ID: "c2"},
		&ir.GitCheckout{Branch: "main"},
		&ir.GitMerge{Branch: "develop"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "v1.0") {
		t.Error("missing tag label")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing commit circles")
	}
	if !strings.Contains(svg, "main") {
		t.Error("missing branch label")
	}
}
