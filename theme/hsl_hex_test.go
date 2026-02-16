package theme

import "testing"

func TestHSLToHex(t *testing.T) {
	tests := []struct {
		name    string
		h, s, l float32
		wantHex string
	}{
		{name: "white", h: 0, s: 0, l: 100, wantHex: "#FFFFFF"},
		{name: "black", h: 0, s: 0, l: 0, wantHex: "#000000"},
		{name: "red", h: 0, s: 100, l: 50, wantHex: "#FF0000"},
		{name: "green", h: 120, s: 100, l: 50, wantHex: "#00FF00"},
		{name: "blue", h: 240, s: 100, l: 50, wantHex: "#0000FF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HSLToHex(tt.h, tt.s, tt.l)
			if got != tt.wantHex {
				t.Errorf("HSLToHex(%v, %v, %v) = %q, want %q",
					tt.h, tt.s, tt.l, got, tt.wantHex)
			}
		})
	}
}

func TestColorRoundTrip(t *testing.T) {
	// Parse hex -> HSL -> back to hex should be stable.
	colors := []string{"#4C78A8", "#FF0000", "#00FF00", "#000000", "#FFFFFF"}
	for _, color := range colors {
		hue, sat, lum, ok := ParseColorToHSL(color)
		if !ok {
			t.Fatalf("ParseColorToHSL(%q) failed", color)
		}
		got := HSLToHex(hue, sat, lum)
		if got != color {
			t.Errorf("round-trip %q -> HSL(%v,%v,%v) -> %q", color, hue, sat, lum, got)
		}
	}
}
