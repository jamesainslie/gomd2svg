package theme

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color conversion constants.
const (
	hexDigits3       = 3   // length of a 3-digit hex color string
	hexDigits6       = 6   // length of a 6-digit hex color string
	rgbMax           = 255 // maximum value for an 8-bit color channel
	rgbMaxF          = 255.0
	hslHalfLightness = 0.5   // midpoint for lightness-based saturation formula
	percentScale     = 100.0 // conversion factor between fraction and percentage
	hueSegmentGreen  = 2.0   // hue offset for the green segment in HSL
	hueSegmentBlue   = 4.0   // hue offset for the blue segment in HSL
	hueSextant       = 6.0   // number of sextants in the hue circle
	hueDegrees       = 60.0  // degrees per sextant
	hueFullCircle    = 360   // full hue circle in degrees
	hueFullCircleF   = 360.0
	satMax           = 100 // maximum saturation percentage
	lumMax           = 100 // maximum lightness percentage
)

// ParseColorToHSL parses a color string and returns its HSL components.
// Supported formats: "#RGB", "#RRGGBB", "hsl(h, s%, l%)".
// Returns ok=false if the color cannot be parsed.
func ParseColorToHSL(color string) (float32, float32, float32, bool) {
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
		h, s, l := rgbToHSL(r, g, b)
		return h, s, l, true
	}

	return 0, 0, 0, false
}

// parseHex parses a 3-digit or 6-digit hex color string (without the leading #).
func parseHex(hex string) (int, int, int, bool) {
	switch len(hex) {
	case hexDigits3:
		// Expand 3-digit hex: #RGB -> #RRGGBB
		rVal, err := strconv.ParseUint(string(hex[0])+string(hex[0]), 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		gVal, err := strconv.ParseUint(string(hex[1])+string(hex[1]), 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		bVal, err := strconv.ParseUint(string(hex[2])+string(hex[2]), 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		return int(rVal), int(gVal), int(bVal), true

	case hexDigits6:
		rVal, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		gVal, err := strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		bVal, err := strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return 0, 0, 0, false
		}
		return int(rVal), int(gVal), int(bVal), true

	default:
		return 0, 0, 0, false
	}
}

// rgbToHSL converts RGB values (0-255) to HSL (h: 0-360, s: 0-100, l: 0-100).
func rgbToHSL(red, grn, blu int) (float32, float32, float32) {
	rf := float64(red) / rgbMaxF
	gf := float64(grn) / rgbMaxF
	bf := float64(blu) / rgbMaxF

	maxC := math.Max(rf, math.Max(gf, bf))
	minC := math.Min(rf, math.Min(gf, bf))
	delta := maxC - minC

	// Lightness
	lf := (maxC + minC) / hueSegmentGreen

	if delta == 0 {
		// Achromatic
		return 0, 0, float32(lf * percentScale)
	}

	// Saturation
	var sf float64
	if lf <= hslHalfLightness {
		sf = delta / (maxC + minC)
	} else {
		sf = delta / (hueSegmentGreen - maxC - minC)
	}

	// Hue
	var hf float64
	switch {
	case rf == maxC:
		hf = (gf - bf) / delta
		if gf < bf {
			hf += hueSextant
		}
	case gf == maxC:
		hf = hueSegmentGreen + (bf-rf)/delta
	default:
		hf = hueSegmentBlue + (rf-gf)/delta
	}
	hf *= hueDegrees

	return float32(hf), float32(sf * percentScale), float32(lf * percentScale)
}

// HSLToHex converts HSL values (h: 0-360, s: 0-100, l: 0-100) to a "#RRGGBB" hex string.
func HSLToHex(hue, sat, lum float32) string {
	red, grn, blu := hslToRGB(hue, sat, lum)
	return fmt.Sprintf("#%02X%02X%02X", red, grn, blu)
}

// hslToRGB converts HSL (h: 0-360, s: 0-100, l: 0-100) to RGB (0-255).
func hslToRGB(hue, sat, lum float32) (int, int, int) {
	sf := float64(sat) / percentScale
	lf := float64(lum) / percentScale
	hf := float64(hue)

	if sf == 0 {
		val := int(math.Round(lf * rgbMaxF))
		return val, val, val
	}

	var chroma2 float64
	if lf < hslHalfLightness {
		chroma2 = lf * (1.0 + sf)
	} else {
		chroma2 = lf + sf - lf*sf
	}
	chroma1 := hueSegmentGreen*lf - chroma2

	hNorm := hf / hueFullCircleF

	red := int(math.Round(hueToRGB(chroma1, chroma2, hNorm+1.0/3.0) * rgbMaxF))
	grn := int(math.Round(hueToRGB(chroma1, chroma2, hNorm) * rgbMaxF))
	blu := int(math.Round(hueToRGB(chroma1, chroma2, hNorm-1.0/3.0) * rgbMaxF))
	return red, grn, blu
}

func hueToRGB(chroma1, chroma2, tVal float64) float64 {
	if tVal < 0 {
		tVal++
	}
	if tVal > 1 {
		tVal--
	}
	switch {
	case tVal < 1.0/6.0:
		return chroma1 + (chroma2-chroma1)*6.0*tVal
	case tVal < 1.0/2.0:
		return chroma2
	case tVal < 2.0/3.0:
		return chroma1 + (chroma2-chroma1)*(2.0/3.0-tVal)*6.0
	default:
		return chroma1
	}
}

// AdjustColor parses the given color, applies HSL adjustments, and returns
// the result as an "hsl(h, s%, l%)" string. If the color cannot be parsed,
// the original string is returned unchanged.
func AdjustColor(color string, hueShift, satShift, lightShift float32) string {
	hue, sat, lum, ok := ParseColorToHSL(color)
	if !ok {
		return color
	}

	hue += hueShift
	sat += satShift
	lum += lightShift

	// Normalize hue to [0, 360)
	for hue < 0 {
		hue += hueFullCircle
	}
	for hue >= hueFullCircle {
		hue -= hueFullCircle
	}

	// Clamp saturation and lightness to [0, 100]
	if sat < 0 {
		sat = 0
	}
	if sat > satMax {
		sat = satMax
	}
	if lum < 0 {
		lum = 0
	}
	if lum > lumMax {
		lum = lumMax
	}

	return fmt.Sprintf("hsl(%.2f, %.2f%%, %.2f%%)", hue, sat, lum)
}
