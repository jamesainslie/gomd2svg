package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseGanttBasic(t *testing.T) {
	input := `gantt
    title A Gantt Diagram
    dateFormat YYYY-MM-DD
    section Development
        Design :d1, 2024-01-01, 10d
        Coding :d2, after d1, 20d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if graph.Kind != ir.Gantt {
		t.Errorf("Kind = %v, want Gantt", graph.Kind)
	}
	if graph.GanttTitle != "A Gantt Diagram" {
		t.Errorf("Title = %q", graph.GanttTitle)
	}
	if graph.GanttDateFormat != "YYYY-MM-DD" {
		t.Errorf("DateFormat = %q", graph.GanttDateFormat)
	}
	if len(graph.GanttSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(graph.GanttSections))
	}
	sec := graph.GanttSections[0]
	if sec.Title != "Development" {
		t.Errorf("Section.Title = %q", sec.Title)
	}
	if len(sec.Tasks) != 2 {
		t.Fatalf("Tasks = %d, want 2", len(sec.Tasks))
	}
	t1 := sec.Tasks[0]
	if t1.ID != "d1" || t1.Label != "Design" {
		t.Errorf("Task[0] = %+v", t1)
	}
	if t1.StartStr != "2024-01-01" || t1.EndStr != "10d" {
		t.Errorf("Task[0] start=%q end=%q", t1.StartStr, t1.EndStr)
	}
	t2 := sec.Tasks[1]
	if t2.ID != "d2" || len(t2.AfterIDs) != 1 || t2.AfterIDs[0] != "d1" {
		t.Errorf("Task[1] = %+v", t2)
	}
}

func TestParseGanttTags(t *testing.T) {
	input := `gantt
    dateFormat YYYY-MM-DD
    section Tasks
        Done task :done, 2024-01-01, 5d
        Critical :crit, active, t3, 2024-01-06, 10d
        Milestone :milestone, m1, 2024-02-01, 0d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	tasks := out.Graph.GanttSections[0].Tasks

	if len(tasks[0].Tags) != 1 || tasks[0].Tags[0] != "done" {
		t.Errorf("Task[0].Tags = %v", tasks[0].Tags)
	}
	if len(tasks[1].Tags) != 2 {
		t.Errorf("Task[1].Tags = %v, want 2 tags", tasks[1].Tags)
	}
	if tasks[1].ID != "t3" {
		t.Errorf("Task[1].ID = %q, want t3", tasks[1].ID)
	}
	if len(tasks[2].Tags) != 1 || tasks[2].Tags[0] != "milestone" {
		t.Errorf("Task[2].Tags = %v", tasks[2].Tags)
	}
	if tasks[2].ID != "m1" {
		t.Errorf("Task[2].ID = %q, want m1", tasks[2].ID)
	}
}

func TestParseGanttDirectives(t *testing.T) {
	input := `gantt
    dateFormat YYYY-MM-DD
    axisFormat %m/%d
    excludes weekends
    tickInterval 1week
    todayMarker off
    weekend friday
    section A
        Task :2024-01-01, 5d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if graph.GanttAxisFormat != "%m/%d" {
		t.Errorf("AxisFormat = %q", graph.GanttAxisFormat)
	}
	if len(graph.GanttExcludes) != 1 || graph.GanttExcludes[0] != "weekends" {
		t.Errorf("Excludes = %v", graph.GanttExcludes)
	}
	if graph.GanttTickInterval != "1week" {
		t.Errorf("TickInterval = %q", graph.GanttTickInterval)
	}
	if graph.GanttTodayMarker != "off" {
		t.Errorf("TodayMarker = %q", graph.GanttTodayMarker)
	}
	if graph.GanttWeekday != "friday" {
		t.Errorf("Weekday = %q", graph.GanttWeekday)
	}
}

func TestParseGanttMultipleAfter(t *testing.T) {
	input := `gantt
    dateFormat YYYY-MM-DD
    section A
        Task A :a1, 2024-01-01, 5d
        Task B :b1, 2024-01-01, 3d
        Task C :c1, after a1 b1, 10d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	tasks := out.Graph.GanttSections[0].Tasks
	if len(tasks[2].AfterIDs) != 2 {
		t.Errorf("Task[2].AfterIDs = %v, want 2", tasks[2].AfterIDs)
	}
}

func TestParseGanttNoSection(t *testing.T) {
	input := `gantt
    dateFormat YYYY-MM-DD
    Task A :2024-01-01, 5d
    Task B :2024-01-06, 3d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.GanttSections) != 1 {
		t.Fatalf("Sections = %d, want 1 (implicit)", len(out.Graph.GanttSections))
	}
	if len(out.Graph.GanttSections[0].Tasks) != 2 {
		t.Errorf("Tasks = %d, want 2", len(out.Graph.GanttSections[0].Tasks))
	}
}
