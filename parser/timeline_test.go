package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseTimelineBasic(t *testing.T) {
	input := `timeline
    title History of Social Media
    2002 : LinkedIn
    2004 : Facebook : Google
    2005 : YouTube`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if graph.Kind != ir.Timeline {
		t.Errorf("Kind = %v, want Timeline", graph.Kind)
	}
	if graph.TimelineTitle != "History of Social Media" {
		t.Errorf("Title = %q", graph.TimelineTitle)
	}
	// No explicit sections, so 1 implicit section.
	if len(graph.TimelineSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(graph.TimelineSections))
	}
	sec := graph.TimelineSections[0]
	if len(sec.Periods) != 3 {
		t.Fatalf("Periods = %d, want 3", len(sec.Periods))
	}
	if sec.Periods[0].Title != "2002" {
		t.Errorf("Period[0].Title = %q", sec.Periods[0].Title)
	}
	if len(sec.Periods[0].Events) != 1 {
		t.Fatalf("Period[0].Events = %d, want 1", len(sec.Periods[0].Events))
	}
	if sec.Periods[0].Events[0].Text != "LinkedIn" {
		t.Errorf("Period[0].Events[0] = %q", sec.Periods[0].Events[0].Text)
	}
	// 2004 has 2 events.
	if len(sec.Periods[1].Events) != 2 {
		t.Fatalf("Period[1].Events = %d, want 2", len(sec.Periods[1].Events))
	}
}

func TestParseTimelineSections(t *testing.T) {
	input := `timeline
    title Product Timeline
    section Phase 1
        January : Research
        February : Design
    section Phase 2
        March : Development`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if len(graph.TimelineSections) != 2 {
		t.Fatalf("Sections = %d, want 2", len(graph.TimelineSections))
	}
	if graph.TimelineSections[0].Title != "Phase 1" {
		t.Errorf("Section[0].Title = %q", graph.TimelineSections[0].Title)
	}
	if len(graph.TimelineSections[0].Periods) != 2 {
		t.Errorf("Section[0].Periods = %d, want 2", len(graph.TimelineSections[0].Periods))
	}
	if graph.TimelineSections[1].Title != "Phase 2" {
		t.Errorf("Section[1].Title = %q", graph.TimelineSections[1].Title)
	}
}

func TestParseTimelineContinuationEvents(t *testing.T) {
	input := `timeline
    2023 : Event A
         : Event B
         : Event C`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	sec := graph.TimelineSections[0]
	if len(sec.Periods) != 1 {
		t.Fatalf("Periods = %d, want 1", len(sec.Periods))
	}
	if len(sec.Periods[0].Events) != 3 {
		t.Fatalf("Events = %d, want 3", len(sec.Periods[0].Events))
	}
	if sec.Periods[0].Events[2].Text != "Event C" {
		t.Errorf("Events[2] = %q", sec.Periods[0].Events[2].Text)
	}
}

func TestParseTimelineMinimal(t *testing.T) {
	input := `timeline
    2024 : Launch`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.TimelineSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(out.Graph.TimelineSections))
	}
	if len(out.Graph.TimelineSections[0].Periods) != 1 {
		t.Fatalf("Periods = %d, want 1", len(out.Graph.TimelineSections[0].Periods))
	}
}
