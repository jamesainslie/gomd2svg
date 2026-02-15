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
