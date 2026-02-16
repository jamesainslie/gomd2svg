package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestKanbanLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Kanban
	graph.Columns = []*ir.KanbanColumn{
		{ID: "todo", Label: "Todo", Cards: []*ir.KanbanCard{
			{ID: "t1", Label: "Task 1"},
			{ID: "t2", Label: "Task 2"},
		}},
		{ID: "done", Label: "Done", Cards: []*ir.KanbanCard{
			{ID: "t3", Label: "Task 3"},
		}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	kd, ok := lay.Diagram.(KanbanData)
	if !ok {
		t.Fatalf("Diagram type = %T, want KanbanData", lay.Diagram)
	}
	if len(kd.Columns) != 2 {
		t.Fatalf("len(Columns) = %d, want 2", len(kd.Columns))
	}

	// Columns should be side by side
	if kd.Columns[1].X <= kd.Columns[0].X {
		t.Errorf("Column[1].X (%v) should be > Column[0].X (%v)", kd.Columns[1].X, kd.Columns[0].X)
	}

	// First column should have 2 cards
	if len(kd.Columns[0].Cards) != 2 {
		t.Errorf("len(Columns[0].Cards) = %d, want 2", len(kd.Columns[0].Cards))
	}

	// Cards should be stacked vertically
	if kd.Columns[0].Cards[1].Y <= kd.Columns[0].Cards[0].Y {
		t.Errorf("Card[1].Y (%v) should be > Card[0].Y (%v)",
			kd.Columns[0].Cards[1].Y, kd.Columns[0].Cards[0].Y)
	}

	// Diagram should have positive dimensions
	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("dimensions = %v x %v, want positive", lay.Width, lay.Height)
	}
}

func TestKanbanLayoutEmptyColumn(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Kanban
	graph.Columns = []*ir.KanbanColumn{
		{ID: "empty", Label: "Empty"},
		{ID: "has", Label: "Has Cards", Cards: []*ir.KanbanCard{
			{ID: "t1", Label: "Task"},
		}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	kd, ok := lay.Diagram.(KanbanData)
	if !ok {
		t.Fatal("expected KanbanData")
	}
	if len(kd.Columns[0].Cards) != 0 {
		t.Errorf("empty column should have 0 cards, got %d", len(kd.Columns[0].Cards))
	}
}
