package gomd2svg

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/parser"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Options configures Mermaid rendering. Zero value uses Modern theme and default layout.
type Options struct {
	// ThemeName selects a built-in theme by name ("modern", "default", "dark",
	// "forest", "neutral"). Takes precedence over Theme if non-empty.
	ThemeName string
	// Theme provides a custom theme. Ignored if ThemeName is set.
	Theme  *theme.Theme
	Layout *config.Layout
}

func (o Options) resolveTheme(dir parser.Directive) *theme.Theme {
	switch {
	// CLI ThemeName takes highest precedence.
	case o.ThemeName != "":
		return theme.ByName(o.ThemeName)
	// Directive theme is second.
	case dir.Theme != "" || dir.ThemeVariables != (parser.ThemeVariables{}):
		var th *theme.Theme
		switch {
		case dir.Theme != "":
			th = theme.ByName(dir.Theme)
		case o.Theme != nil:
			th = o.Theme
		default:
			th = theme.Modern()
		}
		if dir.ThemeVariables != (parser.ThemeVariables{}) {
			ov := theme.Overrides{}
			if dir.ThemeVariables.FontFamily != "" {
				ov.FontFamily = &dir.ThemeVariables.FontFamily
			}
			if dir.ThemeVariables.Background != "" {
				ov.Background = &dir.ThemeVariables.Background
			}
			if dir.ThemeVariables.PrimaryColor != "" {
				ov.PrimaryColor = &dir.ThemeVariables.PrimaryColor
			}
			if dir.ThemeVariables.LineColor != "" {
				ov.LineColor = &dir.ThemeVariables.LineColor
			}
			if dir.ThemeVariables.TextColor != "" {
				ov.TextColor = &dir.ThemeVariables.TextColor
			}
			th = theme.WithOverrides(th, ov)
		}
		return th
	case o.Theme != nil:
		return o.Theme
	default:
		return theme.Modern()
	}
}

func (o Options) layoutOrDefault() *config.Layout {
	if o.Layout != nil {
		return o.Layout
	}
	return config.DefaultLayout()
}

// Result holds the rendered SVG and per-stage timing information.
type Result struct {
	SVG      string
	ParseUs  int64
	LayoutUs int64
	RenderUs int64
}

// TotalUs returns the total rendering time in microseconds.
func (r *Result) TotalUs() int64 {
	return r.ParseUs + r.LayoutUs + r.RenderUs
}

// TotalMs returns the total rendering time in milliseconds.
func (r *Result) TotalMs() float64 {
	const usPerMs = 1000.0
	return float64(r.TotalUs()) / usPerMs
}
