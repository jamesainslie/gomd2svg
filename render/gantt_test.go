package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderGantt(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Gantt
	graph.GanttTitle = "Project"
	graph.GanttDateFormat = "YYYY-MM-DD"
	graph.GanttSections = []*ir.GanttSection{
		{
			Title: "Dev",
			Tasks: []*ir.GanttTask{
				{ID: "t1", Label: "Design", StartStr: "2024-01-01", EndStr: "10d"},
				{ID: "t2", Label: "Code", StartStr: "2024-01-11", EndStr: "20d", Tags: []string{"crit"}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "Project") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Design") {
		t.Error("missing task label")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing task bars")
	}
}
