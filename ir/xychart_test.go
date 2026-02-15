package ir

import "testing"

func TestXYChartSeriesType(t *testing.T) {
	tests := []struct {
		st   XYSeriesType
		want string
	}{
		{XYSeriesBar, "bar"},
		{XYSeriesLine, "line"},
	}
	for _, tc := range tests {
		if got := tc.st.String(); got != tc.want {
			t.Errorf("XYSeriesType(%d).String() = %q, want %q", tc.st, got, tc.want)
		}
	}
}

func TestXYChartAxisMode(t *testing.T) {
	tests := []struct {
		mode XYAxisMode
		want string
	}{
		{XYAxisBand, "band"},
		{XYAxisNumeric, "numeric"},
	}
	for _, tc := range tests {
		if got := tc.mode.String(); got != tc.want {
			t.Errorf("XYAxisMode(%d).String() = %q, want %q", tc.mode, got, tc.want)
		}
	}
}

func TestXYChartGraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = XYChart
	g.XYTitle = "Test"
	g.XYSeries = append(g.XYSeries, &XYSeries{
		Type:   XYSeriesBar,
		Values: []float64{1, 2, 3},
	})
	g.XYXAxis = &XYAxis{
		Mode:       XYAxisBand,
		Title:      "Month",
		Categories: []string{"Jan", "Feb", "Mar"},
	}
	g.XYYAxis = &XYAxis{
		Mode: XYAxisNumeric,
	}

	if g.XYTitle != "Test" {
		t.Errorf("XYTitle = %q, want %q", g.XYTitle, "Test")
	}
	if len(g.XYSeries) != 1 {
		t.Fatalf("XYSeries len = %d, want 1", len(g.XYSeries))
	}
	if g.XYSeries[0].Type != XYSeriesBar {
		t.Errorf("series type = %v, want XYSeriesBar", g.XYSeries[0].Type)
	}
	if g.XYXAxis.Mode != XYAxisBand {
		t.Errorf("x-axis mode = %v, want XYAxisBand", g.XYXAxis.Mode)
	}
}
