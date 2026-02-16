package layout

import (
	"math"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestPieLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Pie
	graph.PieTitle = "Pets"
	graph.PieSlices = []*ir.PieSlice{
		{Label: "Dogs", Value: 50},
		{Label: "Cats", Value: 50},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.Pie {
		t.Errorf("Kind = %v, want Pie", lay.Kind)
	}
	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("dimensions = %f x %f, want > 0", lay.Width, lay.Height)
	}

	pd, ok := lay.Diagram.(PieData)
	if !ok {
		t.Fatalf("Diagram type = %T, want PieData", lay.Diagram)
	}
	if len(pd.Slices) != 2 {
		t.Fatalf("Slices = %d, want 2", len(pd.Slices))
	}

	// Two equal slices: each should span pi radians.
	s0 := pd.Slices[0]
	s1 := pd.Slices[1]
	span0 := s0.EndAngle - s0.StartAngle
	span1 := s1.EndAngle - s1.StartAngle
	if math.Abs(float64(span0-span1)) > 0.01 {
		t.Errorf("spans differ: %f vs %f", span0, span1)
	}
	if math.Abs(float64(span0)-math.Pi) > 0.01 {
		t.Errorf("span0 = %f, want ~pi", span0)
	}
}

func TestPieLayoutSingleSlice(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Pie
	graph.PieSlices = []*ir.PieSlice{
		{Label: "All", Value: 100},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	pd, ok := lay.Diagram.(PieData)
	if !ok {
		t.Fatal("expected PieData")
	}
	if len(pd.Slices) != 1 {
		t.Fatalf("Slices = %d, want 1", len(pd.Slices))
	}
	span := pd.Slices[0].EndAngle - pd.Slices[0].StartAngle
	if math.Abs(float64(span)-2*math.Pi) > 0.01 {
		t.Errorf("span = %f, want ~2*pi", span)
	}
}
