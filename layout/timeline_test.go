package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestTimelineLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Timeline
	graph.TimelineTitle = "History"
	graph.TimelineSections = []*ir.TimelineSection{
		{
			Title: "Early",
			Periods: []*ir.TimelinePeriod{
				{Title: "2002", Events: []*ir.TimelineEvent{{Text: "LinkedIn"}}},
				{Title: "2004", Events: []*ir.TimelineEvent{{Text: "Facebook"}, {Text: "Google"}}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.Timeline {
		t.Errorf("Kind = %v, want Timeline", lay.Kind)
	}
	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("dimensions = %f x %f", lay.Width, lay.Height)
	}

	td, ok := lay.Diagram.(TimelineData)
	if !ok {
		t.Fatalf("Diagram type = %T, want TimelineData", lay.Diagram)
	}
	if len(td.Sections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(td.Sections))
	}
	if len(td.Sections[0].Periods) != 2 {
		t.Fatalf("Periods = %d, want 2", len(td.Sections[0].Periods))
	}
	// Second period has 2 events, so should be taller.
	p0 := td.Sections[0].Periods[0]
	p1 := td.Sections[0].Periods[1]
	if p0.X >= p1.X {
		t.Errorf("Period[0].X=%f >= Period[1].X=%f, want left-to-right", p0.X, p1.X)
	}
}

func TestTimelineLayoutMultipleSections(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Timeline
	graph.TimelineSections = []*ir.TimelineSection{
		{Title: "S1", Periods: []*ir.TimelinePeriod{
			{Title: "P1", Events: []*ir.TimelineEvent{{Text: "E1"}}},
		}},
		{Title: "S2", Periods: []*ir.TimelinePeriod{
			{Title: "P2", Events: []*ir.TimelineEvent{{Text: "E2"}}},
		}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	td, ok := lay.Diagram.(TimelineData)
	if !ok {
		t.Fatal("expected TimelineData")
	}
	if len(td.Sections) != 2 {
		t.Fatalf("Sections = %d, want 2", len(td.Sections))
	}
	// S2 should be below S1.
	if td.Sections[0].Y >= td.Sections[1].Y {
		t.Errorf("S1.Y=%f >= S2.Y=%f, want S1 above S2", td.Sections[0].Y, td.Sections[1].Y)
	}
}
