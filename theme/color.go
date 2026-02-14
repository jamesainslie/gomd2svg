package theme

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ParseColorToHSL parses a color string and returns its HSL components.
// Supported formats: "#RGB", "#RRGGBB", "hsl(h, s%, l%)".
// Returns ok=false if the color cannot be parsed.
func ParseColorToHSL(color string) (h, s, l float32, ok bool) {
	color = strings.TrimSpace(color)

	// Try HSL format: hsl(h, s%, l%)
	if strings.HasPrefix(color, "hsl(") && strings.HasSuffix(color, ")") {
		inner := color[4 : len(color)-1]
		parts := strings.Split(inner, ",")
		if len(parts) != 3 {
			return 0, 0, 0, false
		}

		hVal, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 32)
		if err != nil {
			return 0, 0, 0, false
		}

		sPart := strings.TrimSpace(parts[1])
		sPart = strings.TrimSuffix(sPart, "%")
		sVal, err := strconv.ParseFloat(strings.TrimSpace(sPart), 32)
		if err != nil {
			return 0, 0, 0, false
		}

		lPart := strings.TrimSpace(parts[2])
		lPart = strings.TrimSuffix(lPart, "%")
		lVal, err := strconv.ParseFloat(strings.TrimSpace(lPart), 32)
		if err != nil {
			return 0, 0, 0, false
		}

		return float32(hVal), float32(sVal), float32(lVal), true
	}

	// Try hex format
	if strings.HasPrefix(color, "#") {
		r, g, b, hexOK := parseHex(color[1:])
		if !hexOK {
			return 0, 0, 0, false
		}
		h, s, l = rgbToHSL(r, g, b)
		return h, s, l, true
	}

	return 0, 0, 0, false
}

// parseHex parses a 3-digit or 6-digit hex color string (without the leading #).
func parseHex(s string) (r, g, b int, ok bool) {
	switch len(s) {
	case 3:
		// Expand 3-digit hex: #RGB -> #RRGGBB
		rVal, err := strconv.ParseUint(string(s[0])+string(s[0]), 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		gVal, err := strconv.ParseUint(string(s[1])+string(s[1]), 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		bVal, err := strconv.ParseUint(string(s[2])+string(s[2]), 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		return int(rVal), int(gVal), int(bVal), true

	case 6:
		rVal, err := strconv.ParseUint(s[0:2], 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		gVal, err := strconv.ParseUint(s[2:4], 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		bVal, err := strconv.ParseUint(s[4:6], 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		return int(rVal), int(gVal), int(bVal), true

	default:
		return 0, 0, 0, false
	}
}

// rgbToHSL converts RGB values (0-255) to HSL (h: 0-360, s: 0-100, l: 0-100).
func rgbToHSL(r, g, b int) (h, s, l float32) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	maxC := math.Max(rf, math.Max(gf, bf))
	minC := math.Min(rf, math.Min(gf, bf))
	delta := maxC - minC

	// Lightness
	lf := (maxC + minC) / 2.0

	if delta == 0 {
		// Achromatic
		return 0, 0, float32(lf * 100.0)
	}

	// Saturation
	var sf float64
	if lf <= 0.5 {
		sf = delta / (maxC + minC)
	} else {
		sf = delta / (2.0 - maxC - minC)
	}

	// Hue
	var hf float64
	switch {
	case rf == maxC:
		hf = (gf - bf) / delta
		if gf < bf {
			hf += 6.0
		}
	case gf == maxC:
		hf = 2.0 + (bf-rf)/delta
	default:
		hf = 4.0 + (rf-gf)/delta
	}
	hf *= 60.0

	return float32(hf), float32(sf * 100.0), float32(lf * 100.0)
}

// AdjustColor parses the given color, applies HSL adjustments, and returns
// the result as an "hsl(h, s%, l%)" string. If the color cannot be parsed,
// the original string is returned unchanged.
func AdjustColor(color string, hueShift, satShift, lightShift float32) string {
	h, s, l, ok := ParseColorToHSL(color)
	if !ok {
		return color
	}

	h += hueShift
	s += satShift
	l += lightShift

	// Normalize hue to [0, 360)
	for h < 0 {
		h += 360
	}
	for h >= 360 {
		h -= 360
	}

	// Clamp saturation and lightness to [0, 100]
	if s < 0 {
		s = 0
	}
	if s > 100 {
		s = 100
	}
	if l < 0 {
		l = 0
	}
	if l > 100 {
		l = 100
	}

	return fmt.Sprintf("hsl(%.2f, %.2f%%, %.2f%%)", h, s, l)
}
