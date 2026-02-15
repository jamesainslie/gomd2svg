package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseKanbanBasic(t *testing.T) {
	input := `kanban
  Todo
    task1[Create tests]
    task2[Write docs]
  Done
    task3[Ship feature]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Kanban {
		t.Fatalf("Kind = %v, want Kanban", g.Kind)
	}
	if len(g.Columns) != 2 {
		t.Fatalf("len(Columns) = %d, want 2", len(g.Columns))
	}
	if g.Columns[0].Label != "Todo" {
		t.Errorf("Columns[0].Label = %q, want \"Todo\"", g.Columns[0].Label)
	}
	if len(g.Columns[0].Cards) != 2 {
		t.Errorf("len(Columns[0].Cards) = %d, want 2", len(g.Columns[0].Cards))
	}
	if g.Columns[0].Cards[0].ID != "task1" {
		t.Errorf("Cards[0].ID = %q, want \"task1\"", g.Columns[0].Cards[0].ID)
	}
	if g.Columns[0].Cards[0].Label != "Create tests" {
		t.Errorf("Cards[0].Label = %q, want \"Create tests\"", g.Columns[0].Cards[0].Label)
	}
	if g.Columns[1].Label != "Done" {
		t.Errorf("Columns[1].Label = %q, want \"Done\"", g.Columns[1].Label)
	}
	if len(g.Columns[1].Cards) != 1 {
		t.Errorf("len(Columns[1].Cards) = %d, want 1", len(g.Columns[1].Cards))
	}
}

func TestParseKanbanMetadata(t *testing.T) {
	input := `kanban
  Backlog
    t1[Fix bug]@{ assigned: 'alice', ticket: 'BUG-42', priority: 'High' }`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	g := out.Graph
	if len(g.Columns) != 1 {
		t.Fatalf("len(Columns) = %d, want 1", len(g.Columns))
	}
	card := g.Columns[0].Cards[0]
	if card.Assigned != "alice" {
		t.Errorf("Assigned = %q, want \"alice\"", card.Assigned)
	}
	if card.Ticket != "BUG-42" {
		t.Errorf("Ticket = %q, want \"BUG-42\"", card.Ticket)
	}
	if card.Priority != ir.PriorityHigh {
		t.Errorf("Priority = %v, want PriorityHigh", card.Priority)
	}
}

func TestParseKanbanColumnNoCards(t *testing.T) {
	input := `kanban
  EmptyCol
  WithCards
    t1[Do something]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(out.Graph.Columns) != 2 {
		t.Fatalf("len(Columns) = %d, want 2", len(out.Graph.Columns))
	}
	if len(out.Graph.Columns[0].Cards) != 0 {
		t.Errorf("EmptyCol should have 0 cards, got %d", len(out.Graph.Columns[0].Cards))
	}
}

func TestParseKanbanColumnWithBracketLabel(t *testing.T) {
	input := `kanban
  col1[In Progress]
    t1[Working on it]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if out.Graph.Columns[0].Label != "In Progress" {
		t.Errorf("Label = %q, want \"In Progress\"", out.Graph.Columns[0].Label)
	}
	if out.Graph.Columns[0].ID != "col1" {
		t.Errorf("ID = %q, want \"col1\"", out.Graph.Columns[0].ID)
	}
}
