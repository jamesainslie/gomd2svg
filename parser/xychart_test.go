package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseXYChartBasic(t *testing.T) {
	input := `xychart-beta
    title "Sales Revenue"
    x-axis [jan, feb, mar, apr, may]
    y-axis "Revenue" 0 --> 1000
    bar [100, 200, 300, 400, 500]
    line [150, 250, 350, 450, 550]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if graph.Kind != ir.XYChart {
		t.Fatalf("Kind = %v, want XYChart", graph.Kind)
	}
	if graph.XYTitle != "Sales Revenue" {
		t.Errorf("Title = %q, want %q", graph.XYTitle, "Sales Revenue")
	}
	if graph.XYXAxis == nil {
		t.Fatal("XYXAxis is nil")
	}
	if graph.XYXAxis.Mode != ir.XYAxisBand {
		t.Errorf("x-axis mode = %v, want XYAxisBand", graph.XYXAxis.Mode)
	}
	if len(graph.XYXAxis.Categories) != 5 {
		t.Errorf("x-axis categories len = %d, want 5", len(graph.XYXAxis.Categories))
	}
	if graph.XYYAxis == nil {
		t.Fatal("XYYAxis is nil")
	}
	if graph.XYYAxis.Max != 1000 {
		t.Errorf("y-axis max = %v, want 1000", graph.XYYAxis.Max)
	}
	if len(graph.XYSeries) != 2 {
		t.Fatalf("XYSeries len = %d, want 2", len(graph.XYSeries))
	}
	if graph.XYSeries[0].Type != ir.XYSeriesBar {
		t.Errorf("series[0] type = %v, want Bar", graph.XYSeries[0].Type)
	}
	if graph.XYSeries[1].Type != ir.XYSeriesLine {
		t.Errorf("series[1] type = %v, want Line", graph.XYSeries[1].Type)
	}
}

func TestParseXYChartHorizontal(t *testing.T) {
	input := `xychart-beta horizontal
    bar [10, 20, 30]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if !out.Graph.XYHorizontal {
		t.Error("XYHorizontal = false, want true")
	}
}

func TestParseXYChartNumericXAxis(t *testing.T) {
	input := `xychart-beta
    x-axis "Time" 0 --> 100
    bar [10, 20, 30]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if out.Graph.XYXAxis.Mode != ir.XYAxisNumeric {
		t.Errorf("x-axis mode = %v, want XYAxisNumeric", out.Graph.XYXAxis.Mode)
	}
	if out.Graph.XYXAxis.Min != 0 {
		t.Errorf("x-axis min = %v, want 0", out.Graph.XYXAxis.Min)
	}
	if out.Graph.XYXAxis.Max != 100 {
		t.Errorf("x-axis max = %v, want 100", out.Graph.XYXAxis.Max)
	}
}

func TestParseXYChartMinimal(t *testing.T) {
	input := `xychart-beta
    line [1.5, 2.3, 0.8]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.XYSeries) != 1 {
		t.Fatalf("series len = %d, want 1", len(out.Graph.XYSeries))
	}
	if out.Graph.XYSeries[0].Values[0] != 1.5 {
		t.Errorf("value[0] = %v, want 1.5", out.Graph.XYSeries[0].Values[0])
	}
}
