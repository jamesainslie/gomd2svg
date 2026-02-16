package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestXYChartLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart
	g.XYTitle = "Sales"
	g.XYXAxis = &ir.XYAxis{
		Mode:       ir.XYAxisBand,
		Categories: []string{"Jan", "Feb", "Mar"},
	}
	g.XYYAxis = &ir.XYAxis{
		Mode: ir.XYAxisNumeric,
		Min:  0,
		Max:  100,
	}
	g.XYSeries = []*ir.XYSeries{
		{Type: ir.XYSeriesBar, Values: []float64{30, 60, 90}},
		{Type: ir.XYSeriesLine, Values: []float64{20, 50, 80}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", l.Width, l.Height)
	}

	xyd, ok := l.Diagram.(XYChartData)
	if !ok {
		t.Fatal("Diagram is not XYChartData")
	}
	if xyd.Title != "Sales" {
		t.Errorf("Title = %q, want %q", xyd.Title, "Sales")
	}
	if len(xyd.Series) != 2 {
		t.Fatalf("Series len = %d, want 2", len(xyd.Series))
	}
	if xyd.Series[0].Type != ir.XYSeriesBar {
		t.Errorf("series[0] type = %v, want Bar", xyd.Series[0].Type)
	}
	if len(xyd.Series[0].Points) != 3 {
		t.Errorf("series[0] points len = %d, want 3", len(xyd.Series[0].Points))
	}
	if len(xyd.XLabels) != 3 {
		t.Errorf("XLabels len = %d, want 3", len(xyd.XLabels))
	}
	if len(xyd.YTicks) == 0 {
		t.Error("YTicks is empty")
	}
}

func TestXYChartAutoRange(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart
	g.XYSeries = []*ir.XYSeries{
		{Type: ir.XYSeriesBar, Values: []float64{10, 50, 30}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	xyd, ok := l.Diagram.(XYChartData)
	if !ok {
		t.Fatal("Diagram is not XYChartData")
	}
	if xyd.YMax <= 0 {
		t.Error("YMax should be auto-computed > 0")
	}
}
