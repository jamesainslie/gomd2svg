package theme

import "testing"

func TestModern(t *testing.T) {
	th := Modern()
	if th.FontSize != 14 {
		t.Errorf("FontSize = %f, want 14", th.FontSize)
	}
	if th.PrimaryColor == "" {
		t.Error("PrimaryColor is empty")
	}
	if th.FontFamily == "" {
		t.Error("FontFamily is empty")
	}
}

func TestMermaidDefault(t *testing.T) {
	th := MermaidDefault()
	if th.FontSize != 16 {
		t.Errorf("FontSize = %f, want 16", th.FontSize)
	}
	if th.PrimaryColor != "#ECECFF" {
		t.Errorf("PrimaryColor = %q, want #ECECFF", th.PrimaryColor)
	}
}

func TestModernThemeHasClassColors(t *testing.T) {
	th := Modern()
	if th.ClassHeaderBg == "" {
		t.Error("ClassHeaderBg empty")
	}
	if th.StateFill == "" {
		t.Error("StateFill empty")
	}
	if th.EntityHeaderBg == "" {
		t.Error("EntityHeaderBg empty")
	}
}

func TestModernPieColors(t *testing.T) {
	th := Modern()
	if len(th.PieColors) < 8 {
		t.Errorf("PieColors = %d, want >= 8", len(th.PieColors))
	}
}

func TestModernQuadrantFills(t *testing.T) {
	th := Modern()
	if th.QuadrantFill1 == "" {
		t.Error("QuadrantFill1 is empty")
	}
	if th.QuadrantFill2 == "" {
		t.Error("QuadrantFill2 is empty")
	}
	if th.QuadrantFill3 == "" {
		t.Error("QuadrantFill3 is empty")
	}
	if th.QuadrantFill4 == "" {
		t.Error("QuadrantFill4 is empty")
	}
	if th.QuadrantPointFill == "" {
		t.Error("QuadrantPointFill is empty")
	}
}

func TestModernTimelineColors(t *testing.T) {
	th := Modern()
	if len(th.TimelineSectionColors) < 4 {
		t.Errorf("TimelineSectionColors = %d, want >= 4", len(th.TimelineSectionColors))
	}
	if th.TimelineEventFill == "" {
		t.Error("TimelineEventFill is empty")
	}
}

func TestModernGanttColors(t *testing.T) {
	th := Modern()
	if th.GanttTaskFill == "" {
		t.Error("GanttTaskFill is empty")
	}
	if th.GanttCritFill == "" {
		t.Error("GanttCritFill is empty")
	}
	if len(th.GanttSectionColors) < 4 {
		t.Errorf("GanttSectionColors = %d, want >= 4", len(th.GanttSectionColors))
	}
}

func TestModernGitGraphColors(t *testing.T) {
	th := Modern()
	if len(th.GitBranchColors) < 8 {
		t.Errorf("GitBranchColors = %d, want >= 8", len(th.GitBranchColors))
	}
	if th.GitCommitFill == "" {
		t.Error("GitCommitFill is empty")
	}
}

func TestModernXYChartColors(t *testing.T) {
	th := Modern()
	if len(th.XYChartColors) == 0 {
		t.Error("Modern theme XYChartColors is empty")
	}
	if th.XYChartAxisColor == "" {
		t.Error("Modern theme XYChartAxisColor is empty")
	}
	if th.XYChartGridColor == "" {
		t.Error("Modern theme XYChartGridColor is empty")
	}
}

func TestModernMindmapColors(t *testing.T) {
	th := Modern()
	if len(th.MindmapBranchColors) == 0 {
		t.Error("MindmapBranchColors is empty")
	}
	if th.MindmapNodeBorder == "" {
		t.Error("MindmapNodeBorder is empty")
	}
}

func TestModernSankeyColors(t *testing.T) {
	th := Modern()
	if len(th.SankeyNodeColors) == 0 {
		t.Error("SankeyNodeColors is empty")
	}
	if th.SankeyLinkColor == "" {
		t.Error("SankeyLinkColor is empty")
	}
}

func TestModernTreemapColors(t *testing.T) {
	th := Modern()
	if len(th.TreemapColors) == 0 {
		t.Error("TreemapColors is empty")
	}
	if th.TreemapBorder == "" {
		t.Error("TreemapBorder is empty")
	}
}

func TestModernRadarColors(t *testing.T) {
	th := Modern()
	if len(th.RadarCurveColors) == 0 {
		t.Error("Modern theme RadarCurveColors is empty")
	}
	if th.RadarAxisColor == "" {
		t.Error("Modern theme RadarAxisColor is empty")
	}
	if th.RadarGraticuleColor == "" {
		t.Error("Modern theme RadarGraticuleColor is empty")
	}
}

func TestModernRequirementColors(t *testing.T) {
	th := Modern()
	if th.RequirementFill == "" {
		t.Error("RequirementFill empty")
	}
}

func TestModernBlockColors(t *testing.T) {
	th := Modern()
	if len(th.BlockColors) == 0 {
		t.Error("BlockColors empty")
	}
}

func TestModernC4Colors(t *testing.T) {
	th := Modern()
	if th.C4PersonColor == "" {
		t.Error("C4PersonColor empty")
	}
	if th.C4SystemColor == "" {
		t.Error("C4SystemColor empty")
	}
}
