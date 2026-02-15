package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRadarLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Radar
	g.RadarTitle = "Skills"
	g.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "Speed"},
		{ID: "b", Label: "Power"},
		{ID: "c", Label: "Magic"},
	}
	g.RadarCurves = []*ir.RadarCurve{
		{ID: "p1", Label: "Player1", Values: []float64{80, 60, 40}},
	}
	g.RadarMax = 100

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", l.Width, l.Height)
	}

	rd, ok := l.Diagram.(RadarData)
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
	g := ir.NewGraph()
	g.Kind = ir.Radar
	g.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
	}
	g.RadarCurves = []*ir.RadarCurve{
		{ID: "c", Label: "C", Values: []float64{50, 80}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	rd, ok := l.Diagram.(RadarData)
	if !ok {
		t.Fatal("Diagram is not RadarData")
	}
	if rd.MaxValue <= 0 {
		t.Error("MaxValue should be auto-computed > 0")
	}
}
