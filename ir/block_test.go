package ir

import "testing"

func TestBlockDef(t *testing.T) {
	block := &BlockDef{
		ID:    "a",
		Label: "Block A",
		Shape: Rectangle,
		Width: 2,
	}
	if block.ID != "a" {
		t.Errorf("ID = %q, want %q", block.ID, "a")
	}
	if block.Label != "Block A" {
		t.Errorf("Label = %q, want %q", block.Label, "Block A")
	}
	if block.Width != 2 {
		t.Errorf("Width = %d, want 2", block.Width)
	}
	if block.Shape != Rectangle {
		t.Errorf("Shape = %v, want Rectangle", block.Shape)
	}
}

func TestBlockNesting(t *testing.T) {
	parent := &BlockDef{
		ID:    "parent",
		Label: "Parent",
		Shape: Rectangle,
		Width: 1,
		Children: []*BlockDef{
			{ID: "child1", Label: "Child 1", Shape: Rectangle, Width: 1},
			{ID: "child2", Label: "Child 2", Shape: RoundRect, Width: 1},
		},
	}
	if parent.ID != "parent" {
		t.Errorf("ID = %q, want %q", parent.ID, "parent")
	}
	if parent.Label != "Parent" {
		t.Errorf("Label = %q, want %q", parent.Label, "Parent")
	}
	if parent.Shape != Rectangle {
		t.Errorf("Shape = %v, want Rectangle", parent.Shape)
	}
	if parent.Width != 1 {
		t.Errorf("Width = %d, want 1", parent.Width)
	}
	if len(parent.Children) != 2 {
		t.Fatalf("Children = %d, want 2", len(parent.Children))
	}
	if parent.Children[1].Shape != RoundRect {
		t.Errorf("child2 shape = %v, want RoundRect", parent.Children[1].Shape)
	}
}

func TestBlockGraphFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Block
	graph.BlockColumns = 3
	graph.Blocks = append(graph.Blocks, &BlockDef{ID: "a", Label: "A", Width: 1})
	graph.Blocks = append(graph.Blocks, &BlockDef{ID: "b", Label: "B", Width: 2})
	if graph.BlockColumns != 3 {
		t.Errorf("BlockColumns = %d, want 3", graph.BlockColumns)
	}
	if len(graph.Blocks) != 2 {
		t.Errorf("Blocks = %d, want 2", len(graph.Blocks))
	}
}
