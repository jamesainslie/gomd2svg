package ir

import "testing"

func TestGanttTaskDefaults(t *testing.T) {
	task := &GanttTask{
		ID:       "t1",
		Label:    "Design",
		StartStr: "2024-01-01",
		EndStr:   "10d",
		Tags:     []string{"crit"},
	}
	if task.ID != "t1" {
		t.Errorf("ID = %q, want %q", task.ID, "t1")
	}
	if task.Label != "Design" {
		t.Errorf("Label = %q, want %q", task.Label, "Design")
	}
	if task.StartStr != "2024-01-01" {
		t.Errorf("StartStr = %q, want %q", task.StartStr, "2024-01-01")
	}
	if task.EndStr != "10d" {
		t.Errorf("EndStr = %q, want %q", task.EndStr, "10d")
	}
	if len(task.Tags) != 1 || task.Tags[0] != "crit" {
		t.Errorf("Tags = %v, want [crit]", task.Tags)
	}
}

func TestGanttSectionDefaults(t *testing.T) {
	section := &GanttSection{
		Title: "Development",
		Tasks: []*GanttTask{
			{ID: "d1", Label: "Code"},
		},
	}
	if section.Title != "Development" {
		t.Errorf("Title = %q", section.Title)
	}
	if len(section.Tasks) != 1 {
		t.Errorf("Tasks = %d, want 1", len(section.Tasks))
	}
}

func TestGraphGanttFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Gantt
	graph.GanttTitle = "Project"
	graph.GanttDateFormat = "YYYY-MM-DD"
	graph.GanttAxisFormat = "%Y-%m-%d"
	graph.GanttExcludes = []string{"weekends"}
	graph.GanttSections = append(graph.GanttSections, &GanttSection{
		Title: "Dev",
		Tasks: []*GanttTask{{ID: "t1", Label: "Code"}},
	})

	if graph.GanttTitle != "Project" {
		t.Errorf("GanttTitle = %q", graph.GanttTitle)
	}
	if graph.GanttDateFormat != "YYYY-MM-DD" {
		t.Errorf("GanttDateFormat = %q", graph.GanttDateFormat)
	}
	if len(graph.GanttExcludes) != 1 {
		t.Errorf("GanttExcludes = %d, want 1", len(graph.GanttExcludes))
	}
	if len(graph.GanttSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(graph.GanttSections))
	}
}
