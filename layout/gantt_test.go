package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestGanttLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Gantt
	graph.GanttDateFormat = "YYYY-MM-DD"
	graph.GanttTitle = "Project"
	graph.GanttSections = []*ir.GanttSection{
		{
			Title: "Dev",
			Tasks: []*ir.GanttTask{
				{ID: "t1", Label: "Design", StartStr: "2024-01-01", EndStr: "10d"},
				{ID: "t2", Label: "Code", StartStr: "2024-01-11", EndStr: "20d"},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.Gantt {
		t.Errorf("Kind = %v, want Gantt", lay.Kind)
	}
	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("dimensions = %f x %f", lay.Width, lay.Height)
	}

	gd, ok := lay.Diagram.(GanttData)
	if !ok {
		t.Fatalf("Diagram type = %T, want GanttData", lay.Diagram)
	}
	if len(gd.Sections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(gd.Sections))
	}
	if len(gd.Sections[0].Tasks) != 2 {
		t.Fatalf("Tasks = %d, want 2", len(gd.Sections[0].Tasks))
	}
	// Task 1 should start before Task 2.
	if gd.Sections[0].Tasks[0].X >= gd.Sections[0].Tasks[1].X {
		t.Errorf("Task1.X=%f >= Task2.X=%f", gd.Sections[0].Tasks[0].X, gd.Sections[0].Tasks[1].X)
	}
	// Task 2 should be wider (20d vs 10d).
	if gd.Sections[0].Tasks[1].Width <= gd.Sections[0].Tasks[0].Width {
		t.Errorf("Task2.Width=%f <= Task1.Width=%f", gd.Sections[0].Tasks[1].Width, gd.Sections[0].Tasks[0].Width)
	}
}

func TestGanttLayoutAfterDependency(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Gantt
	graph.GanttDateFormat = "YYYY-MM-DD"
	graph.GanttSections = []*ir.GanttSection{
		{
			Tasks: []*ir.GanttTask{
				{ID: "a", Label: "Task A", StartStr: "2024-01-01", EndStr: "5d"},
				{ID: "b", Label: "Task B", StartStr: "after a", EndStr: "3d", AfterIDs: []string{"a"}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	gd, ok := lay.Diagram.(GanttData)
	if !ok {
		t.Fatal("expected GanttData")
	}
	tasks := gd.Sections[0].Tasks
	// Task B should start where Task A ends.
	if tasks[1].X <= tasks[0].X+tasks[0].Width-1 {
		t.Errorf("TaskB.X=%f should start after TaskA ends at %f", tasks[1].X, tasks[0].X+tasks[0].Width)
	}
}

func TestGanttLayoutMilestone(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Gantt
	graph.GanttDateFormat = "YYYY-MM-DD"
	graph.GanttSections = []*ir.GanttSection{
		{
			Tasks: []*ir.GanttTask{
				{Label: "Release", StartStr: "2024-02-01", EndStr: "0d", Tags: []string{"milestone"}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	gd, ok := lay.Diagram.(GanttData)
	if !ok {
		t.Fatal("expected GanttData")
	}
	task := gd.Sections[0].Tasks[0]
	if !task.IsMilestone {
		t.Error("expected milestone flag")
	}
}

func TestGanttLayoutExcludesWeekends(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Gantt
	graph.GanttDateFormat = "YYYY-MM-DD"
	graph.GanttExcludes = []string{"weekends"}
	graph.GanttSections = []*ir.GanttSection{
		{
			Tasks: []*ir.GanttTask{
				{ID: "t1", Label: "Work", StartStr: "2024-01-01", EndStr: "5d"},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	// Should succeed without panic.
	if lay.Width <= 0 {
		t.Errorf("Width = %f", lay.Width)
	}
}
