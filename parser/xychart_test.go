package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
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
	g := out.Graph
	if g.Kind != ir.XYChart {
		t.Fatalf("Kind = %v, want XYChart", g.Kind)
	}
	if g.XYTitle != "Sales Revenue" {
		t.Errorf("Title = %q, want %q", g.XYTitle, "Sales Revenue")
	}
	if g.XYXAxis == nil {
		t.Fatal("XYXAxis is nil")
	}
	if g.XYXAxis.Mode != ir.XYAxisBand {
		t.Errorf("x-axis mode = %v, want XYAxisBand", g.XYXAxis.Mode)
	}
	if len(g.XYXAxis.Categories) != 5 {
		t.Errorf("x-axis categories len = %d, want 5", len(g.XYXAxis.Categories))
	}
	if g.XYYAxis == nil {
		t.Fatal("XYYAxis is nil")
	}
	if g.XYYAxis.Max != 1000 {
		t.Errorf("y-axis max = %v, want 1000", g.XYYAxis.Max)
	}
	if len(g.XYSeries) != 2 {
		t.Fatalf("XYSeries len = %d, want 2", len(g.XYSeries))
	}
	if g.XYSeries[0].Type != ir.XYSeriesBar {
		t.Errorf("series[0] type = %v, want Bar", g.XYSeries[0].Type)
	}
	if g.XYSeries[1].Type != ir.XYSeriesLine {
		t.Errorf("series[1] type = %v, want Line", g.XYSeries[1].Type)
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
