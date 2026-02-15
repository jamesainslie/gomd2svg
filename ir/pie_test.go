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
	g := NewGraph()
	g.Kind = Pie
	g.PieTitle = "Pets"
	g.PieShowData = true
	g.PieSlices = append(g.PieSlices, &PieSlice{Label: "Dogs", Value: 386})

	if g.PieTitle != "Pets" {
		t.Errorf("PieTitle = %q, want %q", g.PieTitle, "Pets")
	}
	if !g.PieShowData {
		t.Error("PieShowData = false, want true")
	}
	if len(g.PieSlices) != 1 {
		t.Fatalf("PieSlices = %d, want 1", len(g.PieSlices))
	}
}
