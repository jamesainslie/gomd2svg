package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRadarLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Radar
	graph.RadarTitle = "Skills"
	graph.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "Speed"},
		{ID: "b", Label: "Power"},
		{ID: "c", Label: "Magic"},
	}
	graph.RadarCurves = []*ir.RadarCurve{
		{ID: "p1", Label: "Player1", Values: []float64{80, 60, 40}},
	}
	graph.RadarMax = 100

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", lay.Width, lay.Height)
	}

	rd, ok := lay.Diagram.(RadarData)
	if !ok {
		t.Fatal("Diagram is not RadarData")
	}
	if rd.Title != "Skills" {
		t.Errorf("Title = %q, want %q", rd.Title, "Skills")
	}
	if len(rd.Axes) != 3 {
		t.Fatalf("Axes len = %d, want 3", len(rd.Axes))
	}
	if len(rd.Curves) != 1 {
		t.Fatalf("Curves len = %d, want 1", len(rd.Curves))
	}
	if len(rd.Curves[0].Points) != 3 {
		t.Errorf("curve[0] points len = %d, want 3", len(rd.Curves[0].Points))
	}
	if len(rd.GraticuleRadii) == 0 {
		t.Error("GraticuleRadii is empty")
	}
}

func TestRadarAutoMax(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Radar
	graph.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
	}
	graph.RadarCurves = []*ir.RadarCurve{
		{ID: "c", Label: "C", Values: []float64{50, 80}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	rd, ok := lay.Diagram.(RadarData)
	if !ok {
		t.Fatal("Diagram is not RadarData")
	}
	if rd.MaxValue <= 0 {
		t.Error("MaxValue should be auto-computed > 0")
	}
}
