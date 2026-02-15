package ir

import "testing"

func TestJourneyTask(t *testing.T) {
	task := &JourneyTask{
		Name:    "Make tea",
		Score:   5,
		Actors:  []string{"Me"},
		Section: "Morning",
	}
	if task.Name != "Make tea" {
		t.Errorf("Name = %q, want %q", task.Name, "Make tea")
	}
	if task.Score != 5 {
		t.Errorf("Score = %d, want 5", task.Score)
	}
	if len(task.Actors) != 1 || task.Actors[0] != "Me" {
		t.Errorf("Actors = %v, want [Me]", task.Actors)
	}
	if task.Section != "Morning" {
		t.Errorf("Section = %q, want %q", task.Section, "Morning")
	}
}

func TestJourneySection(t *testing.T) {
	section := &JourneySection{
		Name:  "Morning",
		Tasks: []int{0, 1, 2},
	}
	if section.Name != "Morning" {
		t.Errorf("Name = %q, want %q", section.Name, "Morning")
	}
	if len(section.Tasks) != 3 {
		t.Errorf("Tasks len = %d, want 3", len(section.Tasks))
	}
	if section.Tasks[0] != 0 || section.Tasks[1] != 1 || section.Tasks[2] != 2 {
		t.Errorf("Tasks = %v, want [0 1 2]", section.Tasks)
	}
}
