package ir

import "testing"

func TestRadarGraticuleType(t *testing.T) {
	tests := []struct {
		gt   RadarGraticule
		want string
	}{
		{RadarGraticuleCircle, "circle"},
		{RadarGraticulePolygon, "polygon"},
	}
	for _, tc := range tests {
		if got := tc.gt.String(); got != tc.want {
			t.Errorf("RadarGraticule(%d).String() = %q, want %q", tc.gt, got, tc.want)
		}
	}
}

func TestRadarGraphFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Radar
	graph.RadarTitle = "Skills"
	graph.RadarAxes = []*RadarAxis{
		{ID: "e", Label: "English"},
		{ID: "f", Label: "French"},
	}
	graph.RadarCurves = []*RadarCurve{
		{ID: "a", Label: "User1", Values: []float64{80, 60}},
	}
	graph.RadarGraticuleType = RadarGraticuleCircle

	if len(graph.RadarAxes) != 2 {
		t.Fatalf("RadarAxes len = %d, want 2", len(graph.RadarAxes))
	}
	if graph.RadarAxes[0].Label != "English" {
		t.Errorf("axis label = %q, want %q", graph.RadarAxes[0].Label, "English")
	}
	if len(graph.RadarCurves) != 1 {
		t.Fatalf("RadarCurves len = %d, want 1", len(graph.RadarCurves))
	}
	if graph.RadarCurves[0].Values[0] != 80 {
		t.Errorf("curve value = %v, want 80", graph.RadarCurves[0].Values[0])
	}
}
