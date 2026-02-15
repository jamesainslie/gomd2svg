package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseRadarBasic(t *testing.T) {
	input := `radar-beta
    title "Language Skills"
    axis e["English"], f["French"], g["German"]
    curve a["User1"]{80, 60, 70}
    curve b["User2"]{60, 90, 50}`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Radar {
		t.Fatalf("Kind = %v, want Radar", g.Kind)
	}
	if g.RadarTitle != "Language Skills" {
		t.Errorf("Title = %q, want %q", g.RadarTitle, "Language Skills")
	}
	if len(g.RadarAxes) != 3 {
		t.Fatalf("RadarAxes len = %d, want 3", len(g.RadarAxes))
	}
	if g.RadarAxes[0].Label != "English" {
		t.Errorf("axis[0] label = %q, want %q", g.RadarAxes[0].Label, "English")
	}
	if len(g.RadarCurves) != 2 {
		t.Fatalf("RadarCurves len = %d, want 2", len(g.RadarCurves))
	}
	if g.RadarCurves[0].Label != "User1" {
		t.Errorf("curve[0] label = %q, want %q", g.RadarCurves[0].Label, "User1")
	}
	if g.RadarCurves[0].Values[0] != 80 {
		t.Errorf("curve[0] value[0] = %v, want 80", g.RadarCurves[0].Values[0])
	}
}

func TestParseRadarConfig(t *testing.T) {
	input := `radar-beta
    showLegend
    graticule polygon
    ticks 4
    max 100
    min 10
    axis a["A"], b["B"]
    curve c["C"]{50, 60}`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if !g.RadarShowLegend {
		t.Error("RadarShowLegend = false, want true")
	}
	if g.RadarGraticuleType != ir.RadarGraticulePolygon {
		t.Errorf("Graticule = %v, want Polygon", g.RadarGraticuleType)
	}
	if g.RadarTicks != 4 {
		t.Errorf("Ticks = %d, want 4", g.RadarTicks)
	}
	if g.RadarMax != 100 {
		t.Errorf("Max = %v, want 100", g.RadarMax)
	}
	if g.RadarMin != 10 {
		t.Errorf("Min = %v, want 10", g.RadarMin)
	}
}

func TestParseRadarKeyValueCurve(t *testing.T) {
	input := `radar-beta
    axis x["X"], y["Y"], z["Z"]
    curve d["D"]{y: 30, x: 20, z: 10}`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if len(g.RadarCurves) != 1 {
		t.Fatalf("curves len = %d, want 1", len(g.RadarCurves))
	}
	// Key-value maps to axis order: x=20, y=30, z=10
	vals := g.RadarCurves[0].Values
	if len(vals) != 3 {
		t.Fatalf("values len = %d, want 3", len(vals))
	}
	if vals[0] != 20 {
		t.Errorf("vals[0] = %v, want 20 (x)", vals[0])
	}
	if vals[1] != 30 {
		t.Errorf("vals[1] = %v, want 30 (y)", vals[1])
	}
	if vals[2] != 10 {
		t.Errorf("vals[2] = %v, want 10 (z)", vals[2])
	}
}
