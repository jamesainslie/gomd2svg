package ir

import "testing"

func TestBlockDef(t *testing.T) {
	b := &BlockDef{
		ID:    "a",
		Label: "Block A",
		Shape: Rectangle,
		Width: 2,
	}
	if b.Width != 2 {
		t.Errorf("Width = %d, want 2", b.Width)
	}
	if b.Shape != Rectangle {
		t.Errorf("Shape = %v, want Rectangle", b.Shape)
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
	if len(parent.Children) != 2 {
		t.Fatalf("Children = %d, want 2", len(parent.Children))
	}
	if parent.Children[1].Shape != RoundRect {
		t.Errorf("child2 shape = %v, want RoundRect", parent.Children[1].Shape)
	}
}

func TestBlockGraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Block
	g.BlockColumns = 3
	g.Blocks = append(g.Blocks, &BlockDef{ID: "a", Label: "A", Width: 1})
	g.Blocks = append(g.Blocks, &BlockDef{ID: "b", Label: "B", Width: 2})
	if g.BlockColumns != 3 {
		t.Errorf("BlockColumns = %d, want 3", g.BlockColumns)
	}
	if len(g.Blocks) != 2 {
		t.Errorf("Blocks = %d, want 2", len(g.Blocks))
	}
}
