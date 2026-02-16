package gomd2svg

import (
	"strings"
	"testing"
)

func TestDirectiveThemeOverride(t *testing.T) {
	input := "%%{init: {\"theme\": \"dark\"}}%%\nflowchart LR\n  A-->B"
	svg, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	// Dark theme has #1A1A2E background
	if !strings.Contains(svg, "#1A1A2E") {
		t.Error("expected dark background color in SVG")
	}
}

func TestDirectiveOverriddenByCLI(t *testing.T) {
	input := "%%{init: {\"theme\": \"dark\"}}%%\nflowchart LR\n  A-->B"
	svg, err := RenderWithOptions(input, Options{ThemeName: "forest"})
	if err != nil {
		t.Fatal(err)
	}
	// CLI flag (forest) should win over directive (dark).
	// Forest uses #2D6A4F as primary color — unique to this theme.
	if !strings.Contains(svg, "#2D6A4F") {
		t.Error("expected forest primary color #2D6A4F in SVG")
	}
	// Forest background is #FFFFFF; dark is #1A1A2E.
	// The first <rect> is the background — verify it uses the forest background.
	if !strings.Contains(svg, `fill="#FFFFFF"`) {
		t.Error("expected forest background #FFFFFF in SVG")
	}
}

func TestDirectiveFontOverride(t *testing.T) {
	input := "%%{init: {\"themeVariables\": {\"fontFamily\": \"Fira Code\"}}}%%\nflowchart LR\n  A-->B"
	svg, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "Fira Code") {
		t.Error("expected Fira Code font in SVG")
	}
}

func TestDirectiveSingleQuotes(t *testing.T) {
	input := "%%{init: {'theme': 'dark'}}%%\nflowchart LR\n  A-->B"
	svg, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "#1A1A2E") {
		t.Error("expected dark background from single-quote directive")
	}
}
