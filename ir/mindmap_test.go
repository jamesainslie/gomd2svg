package ir

import "testing"

func TestMindmapNodeShape(t *testing.T) {
	tests := []struct {
		shape MindmapShape
		want  string
	}{
		{MindmapShapeDefault, "default"},
		{MindmapSquare, "square"},
		{MindmapRounded, "rounded"},
		{MindmapCircle, "circle"},
		{MindmapBang, "bang"},
		{MindmapCloud, "cloud"},
		{MindmapHexagon, "hexagon"},
	}
	for _, tc := range tests {
		if got := tc.shape.String(); got != tc.want {
			t.Errorf("MindmapShape(%d).String() = %q, want %q", tc.shape, got, tc.want)
		}
	}
}

func TestMindmapNode(t *testing.T) {
	root := &MindmapNode{
		ID:    "root",
		Label: "Central Idea",
		Shape: MindmapCircle,
		Children: []*MindmapNode{
			{ID: "a", Label: "Branch A", Shape: MindmapShapeDefault},
			{ID: "b", Label: "Branch B", Shape: MindmapSquare},
		},
	}
	if root.ID != "root" {
		t.Errorf("ID = %q, want %q", root.ID, "root")
	}
	if root.Label != "Central Idea" {
		t.Errorf("Label = %q, want %q", root.Label, "Central Idea")
	}
	if len(root.Children) != 2 {
		t.Errorf("children len = %d, want 2", len(root.Children))
	}
	if root.Shape != MindmapCircle {
		t.Errorf("shape = %v, want Circle", root.Shape)
	}
}

func TestMindmapGraphField(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Mindmap
	graph.MindmapRoot = &MindmapNode{ID: "root", Label: "Root"}
	if graph.MindmapRoot == nil {
		t.Error("MindmapRoot should not be nil")
	}
}
