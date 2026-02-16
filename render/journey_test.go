package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderJourney(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Journey
	graph.JourneyTitle = "My Day"
	graph.JourneySections = []*ir.JourneySection{
		{Name: "Morning", Tasks: []int{0, 1}},
	}
	graph.JourneyTasks = []*ir.JourneyTask{
		{Name: "Wake up", Score: 3, Actors: []string{"Me"}, Section: "Morning"},
		{Name: "Coffee", Score: 5, Actors: []string{"Me"}, Section: "Morning"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "My Day") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Wake up") {
		t.Error("missing task label")
	}
	if !strings.Contains(svg, "Morning") {
		t.Error("missing section label")
	}
}

func TestRenderJourneyEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Journey

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(graph, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
