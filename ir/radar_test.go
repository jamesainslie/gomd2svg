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
	g := NewGraph()
	g.Kind = Radar
	g.RadarTitle = "Skills"
	g.RadarAxes = []*RadarAxis{
		{ID: "e", Label: "English"},
		{ID: "f", Label: "French"},
	}
	g.RadarCurves = []*RadarCurve{
		{ID: "a", Label: "User1", Values: []float64{80, 60}},
	}
	g.RadarGraticuleType = RadarGraticuleCircle

	if len(g.RadarAxes) != 2 {
		t.Fatalf("RadarAxes len = %d, want 2", len(g.RadarAxes))
	}
	if g.RadarAxes[0].Label != "English" {
		t.Errorf("axis label = %q, want %q", g.RadarAxes[0].Label, "English")
	}
	if len(g.RadarCurves) != 1 {
		t.Fatalf("RadarCurves len = %d, want 1", len(g.RadarCurves))
	}
	if g.RadarCurves[0].Values[0] != 80 {
		t.Errorf("curve value = %v, want 80", g.RadarCurves[0].Values[0])
	}
}
