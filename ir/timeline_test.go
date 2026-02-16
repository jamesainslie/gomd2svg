package ir

import "testing"

func TestTimelineEventDefaults(t *testing.T) {
	e := &TimelineEvent{Text: "Launch"}
	if e.Text != "Launch" {
		t.Errorf("Text = %q, want %q", e.Text, "Launch")
	}
}

func TestTimelinePeriodDefaults(t *testing.T) {
	period := &TimelinePeriod{
		Title:  "2024 Q1",
		Events: []*TimelineEvent{{Text: "Start"}, {Text: "Hire"}},
	}
	if period.Title != "2024 Q1" {
		t.Errorf("Title = %q, want %q", period.Title, "2024 Q1")
	}
	if len(period.Events) != 2 {
		t.Errorf("Events = %d, want 2", len(period.Events))
	}
}

func TestGraphTimelineFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Timeline
	graph.TimelineTitle = "Project"
	graph.TimelineSections = append(graph.TimelineSections, &TimelineSection{
		Title: "Phase 1",
		Periods: []*TimelinePeriod{
			{Title: "Jan", Events: []*TimelineEvent{{Text: "Kickoff"}}},
		},
	})
	if graph.TimelineTitle != "Project" {
		t.Errorf("TimelineTitle = %q, want %q", graph.TimelineTitle, "Project")
	}
	if len(graph.TimelineSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(graph.TimelineSections))
	}
	if len(graph.TimelineSections[0].Periods) != 1 {
		t.Fatalf("Periods = %d, want 1", len(graph.TimelineSections[0].Periods))
	}
}
