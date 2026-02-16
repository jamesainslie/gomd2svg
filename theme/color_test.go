package theme

import (
	"math"
	"testing"
)

func TestParseHex3(t *testing.T) {
	hue, sat, lum, ok := ParseColorToHSL("#fff")
	if !ok {
		t.Fatal("expected ok")
	}
	if math.Abs(float64(lum)-100.0) > 0.1 {
		t.Errorf("l = %f, want ~100", lum)
	}
	_ = hue
	_ = sat
}

func TestParseHex6(t *testing.T) {
	hue, sat, lum, ok := ParseColorToHSL("#ECECFF")
	if !ok {
		t.Fatal("expected ok")
	}
	if hue < 200 || hue > 280 {
		t.Errorf("h = %f, expected ~240", hue)
	}
	_ = sat
	_ = lum
}

func TestParseHSL(t *testing.T) {
	hue, sat, lum, ok := ParseColorToHSL("hsl(240, 100%, 46.27%)")
	if !ok {
		t.Fatal("expected ok")
	}
	if math.Abs(float64(hue)-240) > 0.1 {
		t.Errorf("h = %f, want 240", hue)
	}
	if math.Abs(float64(sat)-100) > 0.1 {
		t.Errorf("s = %f, want 100", sat)
	}
	if math.Abs(float64(lum)-46.27) > 0.1 {
		t.Errorf("l = %f, want 46.27", lum)
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
