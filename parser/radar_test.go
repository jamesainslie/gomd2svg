package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
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
	graph := out.Graph
	if graph.Kind != ir.Radar {
		t.Fatalf("Kind = %v, want Radar", graph.Kind)
	}
	if graph.RadarTitle != "Language Skills" {
		t.Errorf("Title = %q, want %q", graph.RadarTitle, "Language Skills")
	}
	if len(graph.RadarAxes) != 3 {
		t.Fatalf("RadarAxes len = %d, want 3", len(graph.RadarAxes))
	}
	if graph.RadarAxes[0].Label != "English" {
		t.Errorf("axis[0] label = %q, want %q", graph.RadarAxes[0].Label, "English")
	}
	if len(graph.RadarCurves) != 2 {
		t.Fatalf("RadarCurves len = %d, want 2", len(graph.RadarCurves))
	}
	if graph.RadarCurves[0].Label != "User1" {
		t.Errorf("curve[0] label = %q, want %q", graph.RadarCurves[0].Label, "User1")
	}
	if graph.RadarCurves[0].Values[0] != 80 {
		t.Errorf("curve[0] value[0] = %v, want 80", graph.RadarCurves[0].Values[0])
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
	graph := out.Graph
	if !graph.RadarShowLegend {
		t.Error("RadarShowLegend = false, want true")
	}
	if graph.RadarGraticuleType != ir.RadarGraticulePolygon {
		t.Errorf("Graticule = %v, want Polygon", graph.RadarGraticuleType)
	}
	if graph.RadarTicks != 4 {
		t.Errorf("Ticks = %d, want 4", graph.RadarTicks)
	}
	if graph.RadarMax != 100 {
		t.Errorf("Max = %v, want 100", graph.RadarMax)
	}
	if graph.RadarMin != 10 {
		t.Errorf("Min = %v, want 10", graph.RadarMin)
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
	graph := out.Graph
	if len(graph.RadarCurves) != 1 {
		t.Fatalf("curves len = %d, want 1", len(graph.RadarCurves))
	}
	// Key-value maps to axis order: x=20, y=30, z=10
	vals := graph.RadarCurves[0].Values
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
