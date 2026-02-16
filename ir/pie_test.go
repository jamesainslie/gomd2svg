package ir

import "testing"

func TestPieSliceDefaults(t *testing.T) {
	s := &PieSlice{Label: "Dogs", Value: 386}
	if s.Label != "Dogs" {
		t.Errorf("Label = %q, want %q", s.Label, "Dogs")
	}
	if s.Value != 386 {
		t.Errorf("Value = %f, want 386", s.Value)
	}
}

func TestGraphPieFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Pie
	graph.PieTitle = "Pets"
	graph.PieShowData = true
	graph.PieSlices = append(graph.PieSlices, &PieSlice{Label: "Dogs", Value: 386})

	if graph.PieTitle != "Pets" {
		t.Errorf("PieTitle = %q, want %q", graph.PieTitle, "Pets")
	}
	if !graph.PieShowData {
		t.Error("PieShowData = false, want true")
	}
	if len(graph.PieSlices) != 1 {
		t.Fatalf("PieSlices = %d, want 1", len(graph.PieSlices))
	}
}
