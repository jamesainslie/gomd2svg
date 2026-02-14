package theme

import (
	"math"
	"testing"
)

func TestParseHex3(t *testing.T) {
	h, s, l, ok := ParseColorToHSL("#fff")
	if !ok {
		t.Fatal("expected ok")
	}
	if math.Abs(float64(l)-100.0) > 0.1 {
		t.Errorf("l = %f, want ~100", l)
	}
	_ = h
	_ = s
}

func TestParseHex6(t *testing.T) {
	h, s, l, ok := ParseColorToHSL("#ECECFF")
	if !ok {
		t.Fatal("expected ok")
	}
	if h < 200 || h > 280 {
		t.Errorf("h = %f, expected ~240", h)
	}
	_ = s
	_ = l
}

func TestParseHSL(t *testing.T) {
	h, s, l, ok := ParseColorToHSL("hsl(240, 100%, 46.27%)")
	if !ok {
		t.Fatal("expected ok")
	}
	if math.Abs(float64(h)-240) > 0.1 {
		t.Errorf("h = %f, want 240", h)
	}
	if math.Abs(float64(s)-100) > 0.1 {
		t.Errorf("s = %f, want 100", s)
	}
	if math.Abs(float64(l)-46.27) > 0.1 {
		t.Errorf("l = %f, want 46.27", l)
	}
}

func TestAdjustColor(t *testing.T) {
	result := AdjustColor("#ECECFF", 0, 0, -10)
	if len(result) == 0 {
		t.Error("empty result")
	}
	if result[0:3] != "hsl" {
		t.Errorf("expected hsl(...), got %q", result)
	}
}

func TestAdjustColorInvalid(t *testing.T) {
	result := AdjustColor("not-a-color", 0, 0, -10)
	if result != "not-a-color" {
		t.Errorf("expected passthrough for invalid, got %q", result)
	}
}
