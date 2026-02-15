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
