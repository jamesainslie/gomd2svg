package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestJourneyLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Journey
	g.JourneyTitle = "Test Journey"
	g.JourneySections = []*ir.JourneySection{
		{Name: "Section 1", Tasks: []int{0, 1}},
		{Name: "Section 2", Tasks: []int{2}},
	}
	g.JourneyTasks = []*ir.JourneyTask{
		{Name: "Task A", Score: 5, Actors: []string{"Alice"}, Section: "Section 1"},
		{Name: "Task B", Score: 1, Actors: []string{"Alice", "Bob"}, Section: "Section 1"},
		{Name: "Task C", Score: 3, Actors: []string{"Bob"}, Section: "Section 2"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeJourneyLayout(g, th, cfg)

	if l.Kind != ir.Journey {
		t.Fatalf("Kind = %v, want Journey", l.Kind)
	}
	data, ok := l.Diagram.(JourneyData)
	if !ok {
		t.Fatal("Diagram is not JourneyData")
	}
	if data.Title != "Test Journey" {
		t.Errorf("Title = %q, want %q", data.Title, "Test Journey")
	}
	if len(data.Sections) != 2 {
		t.Fatalf("len(Sections) = %d, want 2", len(data.Sections))
	}
	if len(data.Sections[0].Tasks) != 2 {
		t.Errorf("Section 0 tasks = %d, want 2", len(data.Sections[0].Tasks))
	}
	if len(data.Sections[1].Tasks) != 1 {
		t.Errorf("Section 1 tasks = %d, want 1", len(data.Sections[1].Tasks))
	}
	// Task with score 5 should be higher (lower Y) than task with score 1
	if data.Sections[0].Tasks[0].Y >= data.Sections[0].Tasks[1].Y {
		t.Errorf("Score 5 task Y (%f) should be less than score 1 task Y (%f)",
			data.Sections[0].Tasks[0].Y, data.Sections[0].Tasks[1].Y)
	}
	if len(data.Actors) < 2 {
		t.Errorf("len(Actors) = %d, want >= 2", len(data.Actors))
	}
	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %f x %f", l.Width, l.Height)
	}
}

func TestJourneyLayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Journey

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeJourneyLayout(g, th, cfg)

	if l.Kind != ir.Journey {
		t.Fatalf("Kind = %v, want Journey", l.Kind)
	}
	data, ok := l.Diagram.(JourneyData)
	if !ok {
		t.Fatal("Diagram is not JourneyData")
	}
	if len(data.Sections) != 0 {
		t.Errorf("len(Sections) = %d, want 0", len(data.Sections))
	}
}
