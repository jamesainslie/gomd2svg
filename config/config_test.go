package config

import "testing"

func TestDefaultLayout(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.NodeSpacing != 50 {
		t.Errorf("NodeSpacing = %f, want 50", cfg.NodeSpacing)
	}
	if cfg.RankSpacing != 70 {
		t.Errorf("RankSpacing = %f, want 70", cfg.RankSpacing)
	}
	if cfg.LabelLineHeight <= 0 {
		t.Error("LabelLineHeight should be > 0")
	}
}

func TestDefaultLayoutHasClassConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Class.CompartmentPadX <= 0 {
		t.Error("Class.CompartmentPadX should be > 0")
	}
	if cfg.State.CompositePadding <= 0 {
		t.Error("State.CompositePadding should be > 0")
	}
	if cfg.ER.AttributeRowHeight <= 0 {
		t.Error("ER.AttributeRowHeight should be > 0")
	}
}

func TestDefaultLayoutHasSequenceConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Sequence.ParticipantSpacing <= 0 {
		t.Error("Sequence.ParticipantSpacing should be > 0")
	}
	if cfg.Sequence.MessageSpacing <= 0 {
		t.Error("Sequence.MessageSpacing should be > 0")
	}
	if cfg.Sequence.ActivationWidth <= 0 {
		t.Error("Sequence.ActivationWidth should be > 0")
	}
	if cfg.Sequence.HeaderHeight <= 0 {
		t.Error("Sequence.HeaderHeight should be > 0")
	}
}

func TestDefaultLayoutPieConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Pie.Radius != 150 {
		t.Errorf("Pie.Radius = %f, want 150", cfg.Pie.Radius)
	}
	if cfg.Pie.TextPosition != 0.75 {
		t.Errorf("Pie.TextPosition = %f, want 0.75", cfg.Pie.TextPosition)
	}
	if cfg.Pie.PaddingX != 20 {
		t.Errorf("Pie.PaddingX = %f, want 20", cfg.Pie.PaddingX)
	}
}

func TestDefaultLayoutQuadrantConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Quadrant.ChartWidth != 400 {
		t.Errorf("Quadrant.ChartWidth = %f, want 400", cfg.Quadrant.ChartWidth)
	}
	if cfg.Quadrant.ChartHeight != 400 {
		t.Errorf("Quadrant.ChartHeight = %f, want 400", cfg.Quadrant.ChartHeight)
	}
	if cfg.Quadrant.PointRadius != 5 {
		t.Errorf("Quadrant.PointRadius = %f, want 5", cfg.Quadrant.PointRadius)
	}
}

func TestDefaultLayoutTimelineConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Timeline.PeriodWidth != 150 {
		t.Errorf("Timeline.PeriodWidth = %f, want 150", cfg.Timeline.PeriodWidth)
	}
	if cfg.Timeline.EventHeight != 30 {
		t.Errorf("Timeline.EventHeight = %f, want 30", cfg.Timeline.EventHeight)
	}
}

func TestDefaultLayoutGanttConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Gantt.BarHeight != 20 {
		t.Errorf("Gantt.BarHeight = %f, want 20", cfg.Gantt.BarHeight)
	}
	if cfg.Gantt.SidePadding != 75 {
		t.Errorf("Gantt.SidePadding = %f, want 75", cfg.Gantt.SidePadding)
	}
}

func TestDefaultLayoutGitGraphConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.GitGraph.CommitRadius != 8 {
		t.Errorf("GitGraph.CommitRadius = %f, want 8", cfg.GitGraph.CommitRadius)
	}
	if cfg.GitGraph.CommitSpacing != 60 {
		t.Errorf("GitGraph.CommitSpacing = %f, want 60", cfg.GitGraph.CommitSpacing)
	}
	if cfg.GitGraph.BranchSpacing != 40 {
		t.Errorf("GitGraph.BranchSpacing = %f, want 40", cfg.GitGraph.BranchSpacing)
	}
}

func TestXYChartConfigDefaults(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.XYChart.ChartWidth != 700 {
		t.Errorf("XYChart.ChartWidth = %v, want 700", cfg.XYChart.ChartWidth)
	}
	if cfg.XYChart.ChartHeight != 500 {
		t.Errorf("XYChart.ChartHeight = %v, want 500", cfg.XYChart.ChartHeight)
	}
	if cfg.XYChart.BarWidth != 0.6 {
		t.Errorf("XYChart.BarWidth = %v, want 0.6", cfg.XYChart.BarWidth)
	}
}

func TestRadarConfigDefaults(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Radar.Radius != 200 {
		t.Errorf("Radar.Radius = %v, want 200", cfg.Radar.Radius)
	}
	if cfg.Radar.PaddingX != 40 {
		t.Errorf("Radar.PaddingX = %v, want 40", cfg.Radar.PaddingX)
	}
	if cfg.Radar.DefaultTicks != 5 {
		t.Errorf("Radar.DefaultTicks = %v, want 5", cfg.Radar.DefaultTicks)
	}
}
