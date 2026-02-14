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
