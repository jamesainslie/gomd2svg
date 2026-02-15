package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderJourney(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Journey
	g.JourneyTitle = "My Day"
	g.JourneySections = []*ir.JourneySection{
		{Name: "Morning", Tasks: []int{0, 1}},
	}
	g.JourneyTasks = []*ir.JourneyTask{
		{Name: "Wake up", Score: 3, Actors: []string{"Me"}, Section: "Morning"},
		{Name: "Coffee", Score: 5, Actors: []string{"Me"}, Section: "Morning"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
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
	g := ir.NewGraph()
	g.Kind = ir.Journey

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
}
