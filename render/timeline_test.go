package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderTimeline(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Timeline
	graph.TimelineTitle = "History"
	graph.TimelineSections = []*ir.TimelineSection{
		{
			Title: "Early",
			Periods: []*ir.TimelinePeriod{
				{Title: "2002", Events: []*ir.TimelineEvent{{Text: "LinkedIn"}}},
				{Title: "2004", Events: []*ir.TimelineEvent{{Text: "Facebook"}}},
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
	if !strings.Contains(svg, "History") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "LinkedIn") {
		t.Error("missing event text")
	}
	if !strings.Contains(svg, "2002") {
		t.Error("missing period label")
	}
}
