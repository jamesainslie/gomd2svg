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
