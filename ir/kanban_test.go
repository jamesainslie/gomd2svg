package ir

import "testing"

func TestKanbanPriorityString(t *testing.T) {
	tests := []struct {
		p    KanbanPriority
		want string
	}{
		{PriorityNone, ""},
		{PriorityVeryLow, "Very Low"},
		{PriorityLow, "Low"},
		{PriorityHigh, "High"},
		{PriorityVeryHigh, "Very High"},
	}
	for _, tt := range tests {
		if got := tt.p.String(); got != tt.want {
			t.Errorf("KanbanPriority(%d).String() = %q, want %q", tt.p, got, tt.want)
		}
	}
}

func TestKanbanColumnCards(t *testing.T) {
	col := &KanbanColumn{
		ID:    "todo",
		Label: "Todo",
		Cards: []*KanbanCard{
			{ID: "t1", Label: "Task 1", Priority: PriorityHigh},
			{ID: "t2", Label: "Task 2", Assigned: "alice"},
		},
	}
	if len(col.Cards) != 2 {
		t.Fatalf("len(Cards) = %d, want 2", len(col.Cards))
	}
	if col.Cards[0].Priority != PriorityHigh {
		t.Errorf("Cards[0].Priority = %v, want PriorityHigh", col.Cards[0].Priority)
	}
	if col.Cards[1].Assigned != "alice" {
		t.Errorf("Cards[1].Assigned = %q, want \"alice\"", col.Cards[1].Assigned)
	}
}
