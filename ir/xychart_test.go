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
	graph := NewGraph()
	graph.Kind = XYChart
	graph.XYTitle = "Test"
	graph.XYSeries = append(graph.XYSeries, &XYSeries{
		Type:   XYSeriesBar,
		Values: []float64{1, 2, 3},
	})
	graph.XYXAxis = &XYAxis{
		Mode:       XYAxisBand,
		Title:      "Month",
		Categories: []string{"Jan", "Feb", "Mar"},
	}
	graph.XYYAxis = &XYAxis{
		Mode: XYAxisNumeric,
	}

	if graph.XYTitle != "Test" {
		t.Errorf("XYTitle = %q, want %q", graph.XYTitle, "Test")
	}
	if len(graph.XYSeries) != 1 {
		t.Fatalf("XYSeries len = %d, want 1", len(graph.XYSeries))
	}
	if graph.XYSeries[0].Type != XYSeriesBar {
		t.Errorf("series type = %v, want XYSeriesBar", graph.XYSeries[0].Type)
	}
	if graph.XYXAxis.Mode != XYAxisBand {
		t.Errorf("x-axis mode = %v, want XYAxisBand", graph.XYXAxis.Mode)
	}
}
