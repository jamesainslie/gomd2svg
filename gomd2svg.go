// Package mermaid renders Mermaid diagram text to SVG.
//
// The public API wires the full pipeline: Parse -> Layout -> Render SVG.
package gomd2svg

import (
	"fmt"
	"strings"
	"time"

	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/parser"
	"github.com/jamesainslie/gomd2svg/render"
)

// Render parses a Mermaid diagram string and returns SVG output using default options.
func Render(input string) (string, error) {
	return RenderWithOptions(input, Options{})
}

// RenderWithOptions parses a Mermaid diagram string and returns SVG output using the given options.
func RenderWithOptions(input string, opts Options) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", fmt.Errorf("mermaid: empty input")
	}

	cfg := opts.layoutOrDefault()

	parsed, err := parser.Parse(input)
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}

	th := opts.resolveTheme(parsed.Directive)
	l := layout.ComputeLayout(parsed.Graph, th, cfg)
	svg := render.RenderSVG(l, th, cfg)
	return svg, nil
}

// RenderWithTiming parses and renders a Mermaid diagram with per-stage timing.
func RenderWithTiming(input string, opts Options) (*Result, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("mermaid: empty input")
	}

	cfg := opts.layoutOrDefault()

	t0 := time.Now()
	parsed, err := parser.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	parseUs := time.Since(t0).Microseconds()

	th := opts.resolveTheme(parsed.Directive)

	t1 := time.Now()
	l := layout.ComputeLayout(parsed.Graph, th, cfg)
	layoutUs := time.Since(t1).Microseconds()

	t2 := time.Now()
	svg := render.RenderSVG(l, th, cfg)
	renderUs := time.Since(t2).Microseconds()

	return &Result{
		SVG:      svg,
		ParseUs:  parseUs,
		LayoutUs: layoutUs,
		RenderUs: renderUs,
	}, nil
}
