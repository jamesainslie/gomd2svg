package mermaid

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/theme"
)

// Options configures Mermaid rendering. Zero value uses Modern theme and default layout.
type Options struct {
	Theme  *theme.Theme
	Layout *config.Layout
}

func (o Options) themeOrDefault() *theme.Theme {
	if o.Theme != nil {
		return o.Theme
	}
	return theme.Modern()
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
	return float64(r.TotalUs()) / 1000.0
}
