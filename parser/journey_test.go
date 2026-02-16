package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseJourney(t *testing.T) {
	input := `journey
  title My Working Day
  section Go to work
    Make tea: 5: Me
    Go upstairs: 3: Me, Cat
    Do work: 1: Me, Cat
  section Go home
    Go downstairs: 5: Me
    Sit down: 5: Me`

	out, err := parseJourney(input)
	if err != nil {
		t.Fatalf("parseJourney() error: %v", err)
	}
	g := out.Graph

	if g.Kind != ir.Journey {
		t.Errorf("Kind = %v, want Journey", g.Kind)
	}
	if g.JourneyTitle != "My Working Day" {
		t.Errorf("JourneyTitle = %q, want %q", g.JourneyTitle, "My Working Day")
	}

	// Sections
	if len(g.JourneySections) != 2 {
		t.Fatalf("len(JourneySections) = %d, want 2", len(g.JourneySections))
	}
	if g.JourneySections[0].Name != "Go to work" {
		t.Errorf("Section[0].Name = %q, want %q", g.JourneySections[0].Name, "Go to work")
	}
	if len(g.JourneySections[0].Tasks) != 3 {
		t.Errorf("Section[0].Tasks len = %d, want 3", len(g.JourneySections[0].Tasks))
	}
	if g.JourneySections[1].Name != "Go home" {
		t.Errorf("Section[1].Name = %q, want %q", g.JourneySections[1].Name, "Go home")
	}
	if len(g.JourneySections[1].Tasks) != 2 {
		t.Errorf("Section[1].Tasks len = %d, want 2", len(g.JourneySections[1].Tasks))
	}

	// Tasks
	if len(g.JourneyTasks) != 5 {
		t.Fatalf("len(JourneyTasks) = %d, want 5", len(g.JourneyTasks))
	}

	// First task: Make tea
	task0 := g.JourneyTasks[0]
	if task0.Name != "Make tea" {
		t.Errorf("Task[0].Name = %q, want %q", task0.Name, "Make tea")
	}
	if task0.Score != 5 {
		t.Errorf("Task[0].Score = %d, want 5", task0.Score)
	}
	if len(task0.Actors) != 1 || task0.Actors[0] != "Me" {
		t.Errorf("Task[0].Actors = %v, want [Me]", task0.Actors)
	}
	if task0.Section != "Go to work" {
		t.Errorf("Task[0].Section = %q, want %q", task0.Section, "Go to work")
	}

	// Second task: Go upstairs
	task1 := g.JourneyTasks[1]
	if task1.Name != "Go upstairs" {
		t.Errorf("Task[1].Name = %q, want %q", task1.Name, "Go upstairs")
	}
	if task1.Score != 3 {
		t.Errorf("Task[1].Score = %d, want 3", task1.Score)
	}
	if len(task1.Actors) != 2 || task1.Actors[0] != "Me" || task1.Actors[1] != "Cat" {
		t.Errorf("Task[1].Actors = %v, want [Me Cat]", task1.Actors)
	}
}

func TestParseJourneyMinimal(t *testing.T) {
	input := `journey
  Make tea: 5: Me`

	out, err := parseJourney(input)
	if err != nil {
		t.Fatalf("parseJourney() error: %v", err)
	}
	g := out.Graph

	if len(g.JourneyTasks) != 1 {
		t.Fatalf("len(JourneyTasks) = %d, want 1", len(g.JourneyTasks))
	}
	if len(g.JourneySections) != 0 {
		t.Errorf("len(JourneySections) = %d, want 0", len(g.JourneySections))
	}

	task := g.JourneyTasks[0]
	if task.Name != "Make tea" {
		t.Errorf("Task.Name = %q, want %q", task.Name, "Make tea")
	}
	if task.Score != 5 {
		t.Errorf("Task.Score = %d, want 5", task.Score)
	}
	if len(task.Actors) != 1 || task.Actors[0] != "Me" {
		t.Errorf("Task.Actors = %v, want [Me]", task.Actors)
	}
	if task.Section != "" {
		t.Errorf("Task.Section = %q, want empty", task.Section)
	}
}

func TestParseJourneyNoActors(t *testing.T) {
	input := `journey
  section Test
    Task A: 3
    Task B: 4:`

	out, err := parseJourney(input)
	if err != nil {
		t.Fatalf("parseJourney() error: %v", err)
	}
	g := out.Graph

	if len(g.JourneyTasks) != 2 {
		t.Fatalf("len(JourneyTasks) = %d, want 2", len(g.JourneyTasks))
	}

	taskA := g.JourneyTasks[0]
	if taskA.Name != "Task A" {
		t.Errorf("Task[0].Name = %q, want %q", taskA.Name, "Task A")
	}
	if taskA.Score != 3 {
		t.Errorf("Task[0].Score = %d, want 3", taskA.Score)
	}
	if len(taskA.Actors) != 0 {
		t.Errorf("Task[0].Actors = %v, want empty", taskA.Actors)
	}

	taskB := g.JourneyTasks[1]
	if taskB.Name != "Task B" {
		t.Errorf("Task[1].Name = %q, want %q", taskB.Name, "Task B")
	}
	if taskB.Score != 4 {
		t.Errorf("Task[1].Score = %d, want 4", taskB.Score)
	}
	if len(taskB.Actors) != 0 {
		t.Errorf("Task[1].Actors = %v, want empty", taskB.Actors)
	}
}
