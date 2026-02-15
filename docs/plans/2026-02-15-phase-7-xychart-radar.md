# Phase 7: XYChart & Radar Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add XYChart (bar + line on cartesian axes) and Radar (spider/polar chart) diagram support to the mermaid-go renderer.

**Architecture:** Follow the established per-diagram pipeline: IR types → config/theme → parser → layout → renderer → integration tests. XYChart uses cartesian coordinate math (axis scaling, tick generation, bar/line positioning). Radar uses polar coordinate math (angles evenly distributed, values mapped to radii). Both share a color palette concept for multiple data series.

**Tech Stack:** Go stdlib only (math, regexp, strings, strconv, fmt). No external dependencies.

---

## Reference: Mermaid Syntax

### XYChart

```
xychart-beta
    title "Sales Revenue"
    x-axis [jan, feb, mar, apr, may]
    x-axis "Revenue" 100 --> 500
    y-axis "Amount" 0 --> 1000
    bar [100, 200, 300, 400, 500]
    line [150, 250, 350, 450, 550]
```

- Keyword: `xychart-beta` or `xychart`
- Orientation: `xychart-beta horizontal` for horizontal layout (default vertical)
- `title "text"` — optional title
- `x-axis` — categorical `[a, b, c]` or numeric `min --> max` with optional title
- `y-axis` — numeric `min --> max` with optional title (auto-range if omitted)
- `bar [values]` — bar series
- `line [values]` — line series
- Multiple bar/line series allowed

### Radar

```
radar-beta
    title "Skills"
    axis e["English"], f["French"], g["German"]
    curve a["User1"]{80, 60, 70}
    curve b["User2"]{60, 90, 50}
```

- Keyword: `radar-beta` or `radar`
- `title "text"` — optional title
- `axis id["Label"], id["Label"], ...` — defines radial axes
- `curve id["Label"]{v1, v2, ...}` — positional values
- `curve id{axis1: v, axis2: v}` — key-value pairs
- `showLegend` — toggle legend (default true)
- `graticule circle|polygon` — grid shape (default circle)
- `ticks N` — number of concentric grid rings
- `max N` — maximum scale value (auto if omitted)
- `min N` — minimum scale value (default 0)

---

## Existing Codebase Context

**Files you'll modify:**
- `ir/graph.go` — add XYChart and Radar fields to `Graph` struct (after GitGraph fields, ~line 108)
- `config/config.go` — add `XYChart` and `Radar` fields to `Layout` struct and `DefaultLayout()` (after GitGraph, ~line 21 and ~line 228)
- `theme/theme.go` — add XYChart and Radar theme fields to `Theme` struct and both presets (~line 102)
- `parser/parser.go` — add two cases to `Parse()` switch (~line 42)
- `layout/layout.go` — add two cases to `ComputeLayout()` switch
- `layout/types.go` — add `XYChartData` and `RadarData` types with `diagramData()` marker
- `render/svg.go` — add two cases to `RenderSVG()` switch (~line 64)

**Files you'll create:**
- `ir/xychart.go`, `ir/xychart_test.go`
- `ir/radar.go`, `ir/radar_test.go`
- `parser/xychart.go`, `parser/xychart_test.go`
- `parser/radar.go`, `parser/radar_test.go`
- `layout/xychart.go`, `layout/xychart_test.go`
- `layout/radar.go`, `layout/radar_test.go`
- `render/xychart.go`, `render/xychart_test.go`
- `render/radar.go`, `render/radar_test.go`
- `testdata/fixtures/xychart-basic.mmd`, `testdata/fixtures/xychart-horizontal.mmd`
- `testdata/fixtures/radar-basic.mmd`, `testdata/fixtures/radar-polygon.mmd`

**Patterns to follow:**
- iota enums with `String()` methods
- Table-driven tests, no test framework
- `preprocessInput()` for line cleanup
- Package-level `regexp.MustCompile` vars
- `DiagramData` sealed interface with `diagramData()` marker method
- Config-driven sizing via `cfg.XYChart.*` and `cfg.Radar.*`
- Safe comma-ok type assertions in renderers
- Empty color slice guards with fallback colors

---

### Task 1: IR Types — XYChart

**Files:**
- Create: `ir/xychart.go`
- Create: `ir/xychart_test.go`
- Modify: `ir/graph.go:107-108` (add XYChart fields after GitGraph block)

**Step 1: Write the failing test**

Create `ir/xychart_test.go`:

```go
package ir

import "testing"

func TestXYChartSeriesType(t *testing.T) {
	tests := []struct {
		st   XYSeriesType
		want string
	}{
		{XYSeriesBar, "bar"},
		{XYSeriesLine, "line"},
	}
	for _, tc := range tests {
		if got := tc.st.String(); got != tc.want {
			t.Errorf("XYSeriesType(%d).String() = %q, want %q", tc.st, got, tc.want)
		}
	}
}

func TestXYChartAxisMode(t *testing.T) {
	tests := []struct {
		mode XYAxisMode
		want string
	}{
		{XYAxisBand, "band"},
		{XYAxisNumeric, "numeric"},
	}
	for _, tc := range tests {
		if got := tc.mode.String(); got != tc.want {
			t.Errorf("XYAxisMode(%d).String() = %q, want %q", tc.mode, got, tc.want)
		}
	}
}

func TestXYChartGraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = XYChart
	g.XYTitle = "Test"
	g.XYSeries = append(g.XYSeries, &XYSeries{
		Type:   XYSeriesBar,
		Values: []float64{1, 2, 3},
	})
	g.XYXAxis = &XYAxis{
		Mode:       XYAxisBand,
		Title:      "Month",
		Categories: []string{"Jan", "Feb", "Mar"},
	}
	g.XYYAxis = &XYAxis{
		Mode: XYAxisNumeric,
	}

	if g.XYTitle != "Test" {
		t.Errorf("XYTitle = %q, want %q", g.XYTitle, "Test")
	}
	if len(g.XYSeries) != 1 {
		t.Fatalf("XYSeries len = %d, want 1", len(g.XYSeries))
	}
	if g.XYSeries[0].Type != XYSeriesBar {
		t.Errorf("series type = %v, want XYSeriesBar", g.XYSeries[0].Type)
	}
	if g.XYXAxis.Mode != XYAxisBand {
		t.Errorf("x-axis mode = %v, want XYAxisBand", g.XYXAxis.Mode)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestXYChart -v`
Expected: FAIL — types not defined

**Step 3: Write minimal implementation**

Create `ir/xychart.go`:

```go
package ir

// XYSeriesType distinguishes bar and line series.
type XYSeriesType int

const (
	XYSeriesBar XYSeriesType = iota
	XYSeriesLine
)

func (t XYSeriesType) String() string {
	switch t {
	case XYSeriesBar:
		return "bar"
	case XYSeriesLine:
		return "line"
	default:
		return "unknown"
	}
}

// XYAxisMode distinguishes categorical and numeric axes.
type XYAxisMode int

const (
	XYAxisBand XYAxisMode = iota
	XYAxisNumeric
)

func (m XYAxisMode) String() string {
	switch m {
	case XYAxisBand:
		return "band"
	case XYAxisNumeric:
		return "numeric"
	default:
		return "unknown"
	}
}

// XYAxis holds configuration for one axis of an XY chart.
type XYAxis struct {
	Mode       XYAxisMode
	Title      string
	Categories []string  // for band axis
	Min        float64   // for numeric axis
	Max        float64   // for numeric axis
}

// XYSeries holds one data series (bar or line).
type XYSeries struct {
	Type   XYSeriesType
	Values []float64
}
```

Add to `ir/graph.go` after the GitGraph fields block (~line 108):

```go
	// XYChart diagram fields
	XYSeries     []*XYSeries
	XYTitle      string
	XYXAxis      *XYAxis
	XYYAxis      *XYAxis
	XYHorizontal bool
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestXYChart -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/xychart.go ir/xychart_test.go ir/graph.go
git commit -m "feat(ir): add XYChart diagram types"
```

---

### Task 2: IR Types — Radar

**Files:**
- Create: `ir/radar.go`
- Create: `ir/radar_test.go`
- Modify: `ir/graph.go` (add Radar fields after XYChart block)

**Step 1: Write the failing test**

Create `ir/radar_test.go`:

```go
package ir

import "testing"

func TestRadarGraticuleType(t *testing.T) {
	tests := []struct {
		gt   RadarGraticule
		want string
	}{
		{RadarGraticuleCircle, "circle"},
		{RadarGraticulePolygon, "polygon"},
	}
	for _, tc := range tests {
		if got := tc.gt.String(); got != tc.want {
			t.Errorf("RadarGraticule(%d).String() = %q, want %q", tc.gt, got, tc.want)
		}
	}
}

func TestRadarGraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Radar
	g.RadarTitle = "Skills"
	g.RadarAxes = []*RadarAxis{
		{ID: "e", Label: "English"},
		{ID: "f", Label: "French"},
	}
	g.RadarCurves = []*RadarCurve{
		{ID: "a", Label: "User1", Values: []float64{80, 60}},
	}
	g.RadarGraticuleType = RadarGraticuleCircle

	if len(g.RadarAxes) != 2 {
		t.Fatalf("RadarAxes len = %d, want 2", len(g.RadarAxes))
	}
	if g.RadarAxes[0].Label != "English" {
		t.Errorf("axis label = %q, want %q", g.RadarAxes[0].Label, "English")
	}
	if len(g.RadarCurves) != 1 {
		t.Fatalf("RadarCurves len = %d, want 1", len(g.RadarCurves))
	}
	if g.RadarCurves[0].Values[0] != 80 {
		t.Errorf("curve value = %v, want 80", g.RadarCurves[0].Values[0])
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestRadar -v`
Expected: FAIL — types not defined

**Step 3: Write minimal implementation**

Create `ir/radar.go`:

```go
package ir

// RadarGraticule distinguishes graticule shapes.
type RadarGraticule int

const (
	RadarGraticuleCircle RadarGraticule = iota
	RadarGraticulePolygon
)

func (g RadarGraticule) String() string {
	switch g {
	case RadarGraticuleCircle:
		return "circle"
	case RadarGraticulePolygon:
		return "polygon"
	default:
		return "unknown"
	}
}

// RadarAxis defines one radial axis.
type RadarAxis struct {
	ID    string
	Label string
}

// RadarCurve defines one data series on the radar chart.
type RadarCurve struct {
	ID     string
	Label  string
	Values []float64
}
```

Add to `ir/graph.go` after XYChart fields:

```go
	// Radar diagram fields
	RadarAxes          []*RadarAxis
	RadarCurves        []*RadarCurve
	RadarTitle         string
	RadarGraticuleType RadarGraticule
	RadarShowLegend    bool
	RadarTicks         int
	RadarMax           float64
	RadarMin           float64
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestRadar -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/radar.go ir/radar_test.go ir/graph.go
git commit -m "feat(ir): add Radar diagram types"
```

---

### Task 3: Config and Theme

**Files:**
- Modify: `config/config.go` — add `XYChartConfig` and `RadarConfig`
- Modify: `config/config_test.go` — add tests
- Modify: `theme/theme.go` — add XYChart and Radar fields to Theme + both presets
- Modify: `theme/theme_test.go` — add tests

**Step 1: Write the failing tests**

Add to `config/config_test.go`:

```go
func TestXYChartConfigDefaults(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.XYChart.ChartWidth != 700 {
		t.Errorf("XYChart.ChartWidth = %v, want 700", cfg.XYChart.ChartWidth)
	}
	if cfg.XYChart.ChartHeight != 500 {
		t.Errorf("XYChart.ChartHeight = %v, want 500", cfg.XYChart.ChartHeight)
	}
	if cfg.XYChart.BarWidth != 0.6 {
		t.Errorf("XYChart.BarWidth = %v, want 0.6", cfg.XYChart.BarWidth)
	}
}

func TestRadarConfigDefaults(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Radar.Radius != 200 {
		t.Errorf("Radar.Radius = %v, want 200", cfg.Radar.Radius)
	}
	if cfg.Radar.PaddingX != 40 {
		t.Errorf("Radar.PaddingX = %v, want 40", cfg.Radar.PaddingX)
	}
	if cfg.Radar.DefaultTicks != 5 {
		t.Errorf("Radar.DefaultTicks = %v, want 5", cfg.Radar.DefaultTicks)
	}
}
```

Add to `theme/theme_test.go`:

```go
func TestModernXYChartColors(t *testing.T) {
	th := Modern()
	if len(th.XYChartColors) == 0 {
		t.Error("Modern theme XYChartColors is empty")
	}
	if th.XYChartAxisColor == "" {
		t.Error("Modern theme XYChartAxisColor is empty")
	}
	if th.XYChartGridColor == "" {
		t.Error("Modern theme XYChartGridColor is empty")
	}
}

func TestModernRadarColors(t *testing.T) {
	th := Modern()
	if len(th.RadarCurveColors) == 0 {
		t.Error("Modern theme RadarCurveColors is empty")
	}
	if th.RadarAxisColor == "" {
		t.Error("Modern theme RadarAxisColor is empty")
	}
	if th.RadarGraticuleColor == "" {
		t.Error("Modern theme RadarGraticuleColor is empty")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./config/ ./theme/ -run "TestXYChart|TestRadar" -v`
Expected: FAIL — fields not defined

**Step 3: Write minimal implementation**

Add to `config/config.go` — structs after `GitGraphConfig`:

```go
// XYChartConfig holds XY chart layout options.
type XYChartConfig struct {
	ChartWidth   float32
	ChartHeight  float32
	PaddingX     float32
	PaddingY     float32
	BarWidth     float32 // fraction of band width (0-1)
	TickLength   float32
	AxisFontSize float32
	TitleFontSize float32
}

// RadarConfig holds radar chart layout options.
type RadarConfig struct {
	Radius       float32
	PaddingX     float32
	PaddingY     float32
	DefaultTicks int
	LabelOffset  float32 // extra distance for axis labels beyond radius
	CurveOpacity float32
}
```

Add fields to `Layout` struct:

```go
	XYChart  XYChartConfig
	Radar    RadarConfig
```

Add defaults to `DefaultLayout()`:

```go
		XYChart: XYChartConfig{
			ChartWidth:    700,
			ChartHeight:   500,
			PaddingX:      60,
			PaddingY:      40,
			BarWidth:      0.6,
			TickLength:    5,
			AxisFontSize:  12,
			TitleFontSize: 16,
		},
		Radar: RadarConfig{
			Radius:       200,
			PaddingX:     40,
			PaddingY:     40,
			DefaultTicks: 5,
			LabelOffset:  20,
			CurveOpacity: 0.3,
		},
```

Add to `theme/theme.go` Theme struct (after GitGraph fields):

```go
	// XYChart colors
	XYChartColors    []string
	XYChartAxisColor string
	XYChartGridColor string

	// Radar colors
	RadarCurveColors    []string
	RadarAxisColor      string
	RadarGraticuleColor string
	RadarCurveOpacity   float32
```

Add to `Modern()`:

```go
		XYChartColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		XYChartAxisColor: "#6E7B8B",
		XYChartGridColor: "#E0E0E0",

		RadarCurveColors: []string{
			"#4C78A8", "#E45756", "#54A24B", "#F58518",
			"#72B7B2", "#B279A2", "#EECA3B", "#FF9DA6",
		},
		RadarAxisColor:      "#6E7B8B",
		RadarGraticuleColor: "#E0E0E0",
		RadarCurveOpacity:   0.3,
```

Add to `MermaidDefault()`:

```go
		XYChartColors: []string{
			"#4C78A8", "#48A9A6", "#E4E36A", "#F4A261",
			"#E76F51", "#7FB069", "#D08AC0", "#F7B7A3",
		},
		XYChartAxisColor: "#333",
		XYChartGridColor: "#ddd",

		RadarCurveColors: []string{
			"#9370DB", "#E76F51", "#7FB069", "#F4A261",
			"#48A9A6", "#D08AC0", "#E4E36A", "#F7B7A3",
		},
		RadarAxisColor:      "#888",
		RadarGraticuleColor: "#ddd",
		RadarCurveOpacity:   0.3,
```

**Step 4: Run tests to verify they pass**

Run: `go test ./config/ ./theme/ -run "TestXYChart|TestRadar" -v`
Expected: PASS

**Step 5: Commit**

```bash
git add config/config.go config/config_test.go theme/theme.go theme/theme_test.go
git commit -m "feat(config,theme): add XYChart and Radar config and theme"
```

---

### Task 4: XYChart Parser

**Files:**
- Create: `parser/xychart.go`
- Create: `parser/xychart_test.go`
- Modify: `parser/parser.go` — add `case ir.XYChart`

**Step 1: Write the failing test**

Create `parser/xychart_test.go`:

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseXYChartBasic(t *testing.T) {
	input := `xychart-beta
    title "Sales Revenue"
    x-axis [jan, feb, mar, apr, may]
    y-axis "Revenue" 0 --> 1000
    bar [100, 200, 300, 400, 500]
    line [150, 250, 350, 450, 550]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.XYChart {
		t.Fatalf("Kind = %v, want XYChart", g.Kind)
	}
	if g.XYTitle != "Sales Revenue" {
		t.Errorf("Title = %q, want %q", g.XYTitle, "Sales Revenue")
	}
	if g.XYXAxis == nil {
		t.Fatal("XYXAxis is nil")
	}
	if g.XYXAxis.Mode != ir.XYAxisBand {
		t.Errorf("x-axis mode = %v, want XYAxisBand", g.XYXAxis.Mode)
	}
	if len(g.XYXAxis.Categories) != 5 {
		t.Errorf("x-axis categories len = %d, want 5", len(g.XYXAxis.Categories))
	}
	if g.XYYAxis == nil {
		t.Fatal("XYYAxis is nil")
	}
	if g.XYYAxis.Max != 1000 {
		t.Errorf("y-axis max = %v, want 1000", g.XYYAxis.Max)
	}
	if len(g.XYSeries) != 2 {
		t.Fatalf("XYSeries len = %d, want 2", len(g.XYSeries))
	}
	if g.XYSeries[0].Type != ir.XYSeriesBar {
		t.Errorf("series[0] type = %v, want Bar", g.XYSeries[0].Type)
	}
	if g.XYSeries[1].Type != ir.XYSeriesLine {
		t.Errorf("series[1] type = %v, want Line", g.XYSeries[1].Type)
	}
}

func TestParseXYChartHorizontal(t *testing.T) {
	input := `xychart-beta horizontal
    bar [10, 20, 30]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if !out.Graph.XYHorizontal {
		t.Error("XYHorizontal = false, want true")
	}
}

func TestParseXYChartNumericXAxis(t *testing.T) {
	input := `xychart-beta
    x-axis "Time" 0 --> 100
    bar [10, 20, 30]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if out.Graph.XYXAxis.Mode != ir.XYAxisNumeric {
		t.Errorf("x-axis mode = %v, want XYAxisNumeric", out.Graph.XYXAxis.Mode)
	}
	if out.Graph.XYXAxis.Min != 0 {
		t.Errorf("x-axis min = %v, want 0", out.Graph.XYXAxis.Min)
	}
	if out.Graph.XYXAxis.Max != 100 {
		t.Errorf("x-axis max = %v, want 100", out.Graph.XYXAxis.Max)
	}
}

func TestParseXYChartMinimal(t *testing.T) {
	input := `xychart-beta
    line [1.5, 2.3, 0.8]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.XYSeries) != 1 {
		t.Fatalf("series len = %d, want 1", len(out.Graph.XYSeries))
	}
	if out.Graph.XYSeries[0].Values[0] != 1.5 {
		t.Errorf("value[0] = %v, want 1.5", out.Graph.XYSeries[0].Values[0])
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParseXYChart -v`
Expected: FAIL — `parseXYChart` not defined

**Step 3: Write minimal implementation**

Create `parser/xychart.go`:

```go
package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	xyValuesRe   = regexp.MustCompile(`\[([^\]]+)\]`)
	xyNumAxisRe  = regexp.MustCompile(`^(?:"([^"]*)"?\s+)?(-?[\d.]+)\s*-->\s*(-?[\d.]+)$`)
	xyBandAxisRe = regexp.MustCompile(`^(?:"([^"]*)"?\s+)?\[([^\]]+)\]$`)
)

func parseXYChart(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart

	lines := preprocessInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	// Check for horizontal orientation on the first line.
	first := strings.ToLower(lines[0])
	if strings.Contains(first, "horizontal") {
		g.XYHorizontal = true
	}

	for _, line := range lines[1:] {
		lower := strings.ToLower(strings.TrimSpace(line))
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(lower, "title"):
			g.XYTitle = extractQuotedText(trimmed[5:])

		case strings.HasPrefix(lower, "x-axis"):
			g.XYXAxis = parseXYAxis(strings.TrimSpace(trimmed[6:]))

		case strings.HasPrefix(lower, "y-axis"):
			g.XYYAxis = parseXYAxis(strings.TrimSpace(trimmed[6:]))

		case strings.HasPrefix(lower, "bar"):
			if vals := parseXYValues(trimmed); vals != nil {
				g.XYSeries = append(g.XYSeries, &ir.XYSeries{
					Type:   ir.XYSeriesBar,
					Values: vals,
				})
			}

		case strings.HasPrefix(lower, "line"):
			if vals := parseXYValues(trimmed); vals != nil {
				g.XYSeries = append(g.XYSeries, &ir.XYSeries{
					Type:   ir.XYSeriesLine,
					Values: vals,
				})
			}
		}
	}

	return &ParseOutput{Graph: g}, nil
}

func parseXYAxis(s string) *ir.XYAxis {
	// Try numeric range: "Title" min --> max  or  min --> max
	if m := xyNumAxisRe.FindStringSubmatch(s); m != nil {
		axis := &ir.XYAxis{Mode: ir.XYAxisNumeric, Title: m[1]}
		axis.Min, _ = strconv.ParseFloat(m[2], 64) // regex guarantees digits
		axis.Max, _ = strconv.ParseFloat(m[3], 64) // regex guarantees digits
		return axis
	}
	// Try band/categorical: "Title" [a, b, c]  or  [a, b, c]
	if m := xyBandAxisRe.FindStringSubmatch(s); m != nil {
		cats := splitAndTrim(m[2])
		return &ir.XYAxis{Mode: ir.XYAxisBand, Title: m[1], Categories: cats}
	}
	// Title only (auto-range).
	title := extractQuotedText(s)
	if title == "" {
		title = strings.TrimSpace(s)
	}
	return &ir.XYAxis{Mode: ir.XYAxisNumeric, Title: title}
}

func parseXYValues(line string) []float64 {
	m := xyValuesRe.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	parts := splitAndTrim(m[1])
	vals := make([]float64, 0, len(parts))
	for _, p := range parts {
		v, err := strconv.ParseFloat(p, 64)
		if err == nil {
			vals = append(vals, v)
		}
	}
	return vals
}

func extractQuotedText(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' {
		if end := strings.Index(s[1:], "\""); end >= 0 {
			return s[1 : end+1]
		}
	}
	return ""
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, "\"")
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
```

Add to `parser/parser.go` Parse() switch (after `case ir.GitGraph`):

```go
	case ir.XYChart:
		return parseXYChart(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParseXYChart -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/xychart.go parser/xychart_test.go parser/parser.go
git commit -m "feat(parser): add XYChart parser"
```

---

### Task 5: Radar Parser

**Files:**
- Create: `parser/radar.go`
- Create: `parser/radar_test.go`
- Modify: `parser/parser.go` — add `case ir.Radar`

**Step 1: Write the failing test**

Create `parser/radar_test.go`:

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseRadarBasic(t *testing.T) {
	input := `radar-beta
    title "Language Skills"
    axis e["English"], f["French"], g["German"]
    curve a["User1"]{80, 60, 70}
    curve b["User2"]{60, 90, 50}`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Radar {
		t.Fatalf("Kind = %v, want Radar", g.Kind)
	}
	if g.RadarTitle != "Language Skills" {
		t.Errorf("Title = %q, want %q", g.RadarTitle, "Language Skills")
	}
	if len(g.RadarAxes) != 3 {
		t.Fatalf("RadarAxes len = %d, want 3", len(g.RadarAxes))
	}
	if g.RadarAxes[0].Label != "English" {
		t.Errorf("axis[0] label = %q, want %q", g.RadarAxes[0].Label, "English")
	}
	if len(g.RadarCurves) != 2 {
		t.Fatalf("RadarCurves len = %d, want 2", len(g.RadarCurves))
	}
	if g.RadarCurves[0].Label != "User1" {
		t.Errorf("curve[0] label = %q, want %q", g.RadarCurves[0].Label, "User1")
	}
	if g.RadarCurves[0].Values[0] != 80 {
		t.Errorf("curve[0] value[0] = %v, want 80", g.RadarCurves[0].Values[0])
	}
}

func TestParseRadarConfig(t *testing.T) {
	input := `radar-beta
    showLegend
    graticule polygon
    ticks 4
    max 100
    min 10
    axis a["A"], b["B"]
    curve c["C"]{50, 60}`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if !g.RadarShowLegend {
		t.Error("RadarShowLegend = false, want true")
	}
	if g.RadarGraticuleType != ir.RadarGraticulePolygon {
		t.Errorf("Graticule = %v, want Polygon", g.RadarGraticuleType)
	}
	if g.RadarTicks != 4 {
		t.Errorf("Ticks = %d, want 4", g.RadarTicks)
	}
	if g.RadarMax != 100 {
		t.Errorf("Max = %v, want 100", g.RadarMax)
	}
	if g.RadarMin != 10 {
		t.Errorf("Min = %v, want 10", g.RadarMin)
	}
}

func TestParseRadarKeyValueCurve(t *testing.T) {
	input := `radar-beta
    axis x["X"], y["Y"], z["Z"]
    curve d["D"]{y: 30, x: 20, z: 10}`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if len(g.RadarCurves) != 1 {
		t.Fatalf("curves len = %d, want 1", len(g.RadarCurves))
	}
	// Key-value maps to axis order: x=20, y=30, z=10
	vals := g.RadarCurves[0].Values
	if len(vals) != 3 {
		t.Fatalf("values len = %d, want 3", len(vals))
	}
	if vals[0] != 20 {
		t.Errorf("vals[0] = %v, want 20 (x)", vals[0])
	}
	if vals[1] != 30 {
		t.Errorf("vals[1] = %v, want 30 (y)", vals[1])
	}
	if vals[2] != 10 {
		t.Errorf("vals[2] = %v, want 10 (z)", vals[2])
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParseRadar -v`
Expected: FAIL — `parseRadar` not defined

**Step 3: Write minimal implementation**

Create `parser/radar.go`:

```go
package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	radarAxisRe  = regexp.MustCompile(`(\w+)\["([^"]+)"\]`)
	radarCurveRe = regexp.MustCompile(`^curve\s+(\w+)(?:\["([^"]+)"\])?\s*\{([^}]+)\}`)
	radarKVRe    = regexp.MustCompile(`(\w+)\s*:\s*(-?[\d.]+)`)
)

func parseRadar(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Radar

	lines := preprocessInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	// Build axis ID -> index map for key-value curve resolution.
	var axisIDs []string

	for _, line := range lines[1:] {
		lower := strings.ToLower(strings.TrimSpace(line))
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(lower, "title"):
			g.RadarTitle = extractQuotedText(trimmed[5:])

		case strings.HasPrefix(lower, "showlegend"):
			g.RadarShowLegend = true

		case strings.HasPrefix(lower, "graticule"):
			rest := strings.TrimSpace(lower[len("graticule"):])
			if rest == "polygon" {
				g.RadarGraticuleType = ir.RadarGraticulePolygon
			} else {
				g.RadarGraticuleType = ir.RadarGraticuleCircle
			}

		case strings.HasPrefix(lower, "ticks"):
			if v, err := strconv.Atoi(strings.TrimSpace(lower[5:])); err == nil {
				g.RadarTicks = v
			}

		case strings.HasPrefix(lower, "max"):
			if v, err := strconv.ParseFloat(strings.TrimSpace(lower[3:]), 64); err == nil {
				g.RadarMax = v
			}

		case strings.HasPrefix(lower, "min"):
			if v, err := strconv.ParseFloat(strings.TrimSpace(lower[3:]), 64); err == nil {
				g.RadarMin = v
			}

		case strings.HasPrefix(lower, "axis"):
			matches := radarAxisRe.FindAllStringSubmatch(trimmed, -1)
			for _, m := range matches {
				g.RadarAxes = append(g.RadarAxes, &ir.RadarAxis{
					ID:    m[1],
					Label: m[2],
				})
				axisIDs = append(axisIDs, m[1])
			}

		case strings.HasPrefix(lower, "curve"):
			if m := radarCurveRe.FindStringSubmatch(trimmed); m != nil {
				curve := &ir.RadarCurve{ID: m[1], Label: m[2]}
				valStr := m[3]

				// Check for key-value syntax.
				if kvMatches := radarKVRe.FindAllStringSubmatch(valStr, -1); len(kvMatches) > 0 {
					kvMap := make(map[string]float64)
					for _, kv := range kvMatches {
						v, _ := strconv.ParseFloat(kv[2], 64) // regex guarantees digits
						kvMap[kv[1]] = v
					}
					// Map to axis order.
					curve.Values = make([]float64, len(axisIDs))
					for i, id := range axisIDs {
						curve.Values[i] = kvMap[id]
					}
				} else {
					// Positional values.
					parts := splitAndTrim(valStr)
					for _, p := range parts {
						v, err := strconv.ParseFloat(p, 64)
						if err == nil {
							curve.Values = append(curve.Values, v)
						}
					}
				}
				g.RadarCurves = append(g.RadarCurves, curve)
			}
		}
	}

	return &ParseOutput{Graph: g}, nil
}
```

Add to `parser/parser.go` Parse() switch (after XYChart case):

```go
	case ir.Radar:
		return parseRadar(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParseRadar -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/radar.go parser/radar_test.go parser/parser.go
git commit -m "feat(parser): add Radar chart parser"
```

---

### Task 6: XYChart Layout

**Files:**
- Create: `layout/xychart.go`
- Create: `layout/xychart_test.go`
- Modify: `layout/types.go` — add `XYChartData` types
- Modify: `layout/layout.go` — add `case ir.XYChart`

**Step 1: Write the failing test**

Create `layout/xychart_test.go`:

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestXYChartLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart
	g.XYTitle = "Sales"
	g.XYXAxis = &ir.XYAxis{
		Mode:       ir.XYAxisBand,
		Categories: []string{"Jan", "Feb", "Mar"},
	}
	g.XYYAxis = &ir.XYAxis{
		Mode: ir.XYAxisNumeric,
		Min:  0,
		Max:  100,
	}
	g.XYSeries = []*ir.XYSeries{
		{Type: ir.XYSeriesBar, Values: []float64{30, 60, 90}},
		{Type: ir.XYSeriesLine, Values: []float64{20, 50, 80}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", l.Width, l.Height)
	}

	xyd, ok := l.Diagram.(XYChartData)
	if !ok {
		t.Fatal("Diagram is not XYChartData")
	}
	if xyd.Title != "Sales" {
		t.Errorf("Title = %q, want %q", xyd.Title, "Sales")
	}
	if len(xyd.Series) != 2 {
		t.Fatalf("Series len = %d, want 2", len(xyd.Series))
	}
	if xyd.Series[0].Type != ir.XYSeriesBar {
		t.Errorf("series[0] type = %v, want Bar", xyd.Series[0].Type)
	}
	if len(xyd.Series[0].Points) != 3 {
		t.Errorf("series[0] points len = %d, want 3", len(xyd.Series[0].Points))
	}
	if len(xyd.XLabels) != 3 {
		t.Errorf("XLabels len = %d, want 3", len(xyd.XLabels))
	}
	if len(xyd.YTicks) == 0 {
		t.Error("YTicks is empty")
	}
}

func TestXYChartAutoRange(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart
	g.XYSeries = []*ir.XYSeries{
		{Type: ir.XYSeriesBar, Values: []float64{10, 50, 30}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	xyd, ok := l.Diagram.(XYChartData)
	if !ok {
		t.Fatal("Diagram is not XYChartData")
	}
	if xyd.YMax <= 0 {
		t.Error("YMax should be auto-computed > 0")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestXYChart -v`
Expected: FAIL — types not defined

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// XYChartData holds XY chart layout data.
type XYChartData struct {
	Series       []XYSeriesLayout
	XLabels      []XYAxisLabel
	YTicks       []XYAxisTick
	Title        string
	ChartX       float32
	ChartY       float32
	ChartWidth   float32
	ChartHeight  float32
	YMin         float64
	YMax         float64
	Horizontal   bool
}

func (XYChartData) diagramData() {}

// XYSeriesLayout holds one positioned data series.
type XYSeriesLayout struct {
	Type       ir.XYSeriesType
	Points     []XYPointLayout
	ColorIndex int
}

// XYPointLayout holds the pixel position and value of one data point.
type XYPointLayout struct {
	X      float32
	Y      float32
	Width  float32 // bar width (0 for line points)
	Height float32 // bar height (0 for line points)
	Value  float64
}

// XYAxisLabel holds a label on the x-axis.
type XYAxisLabel struct {
	Text string
	X    float32
}

// XYAxisTick holds a tick mark on the y-axis.
type XYAxisTick struct {
	Label string
	Y     float32
}
```

Create `layout/xychart.go`:

```go
package layout

import (
	"fmt"
	"math"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeXYChartLayout(g *ir.Graph, _ *theme.Theme, cfg *config.Layout) *Layout {
	padX := cfg.XYChart.PaddingX
	padY := cfg.XYChart.PaddingY
	chartW := cfg.XYChart.ChartWidth
	chartH := cfg.XYChart.ChartHeight

	// Title height.
	var titleHeight float32
	if g.XYTitle != "" {
		titleHeight = cfg.XYChart.TitleFontSize + padY
	}

	// Determine Y range.
	yMin, yMax := xyDataRange(g)
	if g.XYYAxis != nil && (g.XYYAxis.Min != 0 || g.XYYAxis.Max != 0) {
		yMin = g.XYYAxis.Min
		yMax = g.XYYAxis.Max
	}
	if yMax <= yMin {
		yMax = yMin + 1
	}
	yRange := yMax - yMin

	chartX := padX
	chartY := titleHeight + padY

	// valueToY maps a data value to a pixel Y position.
	valueToY := func(v float64) float32 {
		frac := (v - yMin) / yRange
		return chartY + chartH - float32(frac)*chartH
	}

	// Generate Y-axis ticks.
	yTicks := generateYTicks(yMin, yMax, 5, chartY, chartH)

	// Generate X-axis labels and series points.
	var xLabels []XYAxisLabel
	numPoints := xyMaxPoints(g)
	if numPoints == 0 {
		numPoints = 1
	}

	bandW := chartW / float32(numPoints)

	// X-axis labels from categories or numeric range.
	if g.XYXAxis != nil && g.XYXAxis.Mode == ir.XYAxisBand {
		for i, cat := range g.XYXAxis.Categories {
			xLabels = append(xLabels, XYAxisLabel{
				Text: cat,
				X:    chartX + float32(i)*bandW + bandW/2,
			})
		}
	} else {
		for i := range numPoints {
			xLabels = append(xLabels, XYAxisLabel{
				Text: fmt.Sprintf("%d", i+1),
				X:    chartX + float32(i)*bandW + bandW/2,
			})
		}
	}

	// Count bar series for grouping.
	barCount := 0
	for _, s := range g.XYSeries {
		if s.Type == ir.XYSeriesBar {
			barCount++
		}
	}

	barGroupWidth := bandW * cfg.XYChart.BarWidth
	var singleBarW float32
	if barCount > 0 {
		singleBarW = barGroupWidth / float32(barCount)
	}

	// Build series layouts.
	barIdx := 0
	var series []XYSeriesLayout
	for si, s := range g.XYSeries {
		var points []XYPointLayout
		for i, v := range s.Values {
			cx := chartX + float32(i)*bandW + bandW/2
			py := valueToY(v)

			switch s.Type {
			case ir.XYSeriesBar:
				barX := cx - barGroupWidth/2 + float32(barIdx)*singleBarW
				baseY := valueToY(math.Max(yMin, 0))
				h := baseY - py
				if h < 0 {
					h = -h
					py = baseY
				}
				points = append(points, XYPointLayout{
					X: barX, Y: py, Width: singleBarW, Height: h, Value: v,
				})
			case ir.XYSeriesLine:
				points = append(points, XYPointLayout{
					X: cx, Y: py, Value: v,
				})
			}
		}
		if s.Type == ir.XYSeriesBar {
			barIdx++
		}
		series = append(series, XYSeriesLayout{
			Type:       s.Type,
			Points:     points,
			ColorIndex: si,
		})
	}

	totalW := padX*2 + chartW
	totalH := titleHeight + padY*2 + chartH + cfg.XYChart.AxisFontSize + padY

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: XYChartData{
			Series:      series,
			XLabels:     xLabels,
			YTicks:      yTicks,
			Title:       g.XYTitle,
			ChartX:      chartX,
			ChartY:      chartY,
			ChartWidth:  chartW,
			ChartHeight: chartH,
			YMin:        yMin,
			YMax:        yMax,
			Horizontal:  g.XYHorizontal,
		},
	}
}

func xyDataRange(g *ir.Graph) (float64, float64) {
	min, max := math.Inf(1), math.Inf(-1)
	for _, s := range g.XYSeries {
		for _, v := range s.Values {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}
	if math.IsInf(min, 1) {
		min = 0
	}
	if math.IsInf(max, -1) {
		max = 100
	}
	if min > 0 {
		min = 0
	}
	// Round max up to a nice number.
	max = niceMax(max)
	return min, max
}

func niceMax(v float64) float64 {
	if v <= 0 {
		return 1
	}
	magnitude := math.Pow(10, math.Floor(math.Log10(v)))
	normalized := v / magnitude
	if normalized <= 1 {
		return magnitude
	} else if normalized <= 2 {
		return 2 * magnitude
	} else if normalized <= 5 {
		return 5 * magnitude
	}
	return 10 * magnitude
}

func generateYTicks(yMin, yMax float64, count int, chartY, chartH float32) []XYAxisTick {
	yRange := yMax - yMin
	step := yRange / float64(count)
	var ticks []XYAxisTick
	for i := range count + 1 {
		v := yMin + float64(i)*step
		frac := (v - yMin) / yRange
		y := chartY + chartH - float32(frac)*chartH
		ticks = append(ticks, XYAxisTick{
			Label: fmt.Sprintf("%.4g", v),
			Y:     y,
		})
	}
	return ticks
}

func xyMaxPoints(g *ir.Graph) int {
	max := 0
	for _, s := range g.XYSeries {
		if len(s.Values) > max {
			max = len(s.Values)
		}
	}
	return max
}
```

Add to `layout/layout.go` ComputeLayout switch:

```go
	case ir.XYChart:
		return computeXYChartLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./layout/ -run TestXYChart -v`
Expected: PASS

**Step 5: Commit**

```bash
git add layout/xychart.go layout/xychart_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add XYChart layout"
```

---

### Task 7: Radar Layout

**Files:**
- Create: `layout/radar.go`
- Create: `layout/radar_test.go`
- Modify: `layout/types.go` — add `RadarData` types
- Modify: `layout/layout.go` — add `case ir.Radar`

**Step 1: Write the failing test**

Create `layout/radar_test.go`:

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRadarLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Radar
	g.RadarTitle = "Skills"
	g.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "Speed"},
		{ID: "b", Label: "Power"},
		{ID: "c", Label: "Magic"},
	}
	g.RadarCurves = []*ir.RadarCurve{
		{ID: "p1", Label: "Player1", Values: []float64{80, 60, 40}},
	}
	g.RadarMax = 100

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", l.Width, l.Height)
	}

	rd, ok := l.Diagram.(RadarData)
	if !ok {
		t.Fatal("Diagram is not RadarData")
	}
	if rd.Title != "Skills" {
		t.Errorf("Title = %q, want %q", rd.Title, "Skills")
	}
	if len(rd.Axes) != 3 {
		t.Fatalf("Axes len = %d, want 3", len(rd.Axes))
	}
	if len(rd.Curves) != 1 {
		t.Fatalf("Curves len = %d, want 1", len(rd.Curves))
	}
	if len(rd.Curves[0].Points) != 3 {
		t.Errorf("curve[0] points len = %d, want 3", len(rd.Curves[0].Points))
	}
	if len(rd.GraticuleRadii) == 0 {
		t.Error("GraticuleRadii is empty")
	}
}

func TestRadarAutoMax(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Radar
	g.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
	}
	g.RadarCurves = []*ir.RadarCurve{
		{ID: "c", Label: "C", Values: []float64{50, 80}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	rd, ok := l.Diagram.(RadarData)
	if !ok {
		t.Fatal("Diagram is not RadarData")
	}
	if rd.MaxValue <= 0 {
		t.Error("MaxValue should be auto-computed > 0")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestRadar -v`
Expected: FAIL — types not defined

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// RadarData holds radar chart layout data.
type RadarData struct {
	Axes            []RadarAxisLayout
	Curves          []RadarCurveLayout
	GraticuleRadii  []float32
	GraticuleType   ir.RadarGraticule
	CenterX         float32
	CenterY         float32
	Radius          float32
	Title           string
	ShowLegend      bool
	MaxValue        float64
	MinValue        float64
}

func (RadarData) diagramData() {}

// RadarAxisLayout holds the endpoint and label position of one axis.
type RadarAxisLayout struct {
	Label  string
	EndX   float32
	EndY   float32
	LabelX float32
	LabelY float32
}

// RadarCurveLayout holds one data series polygon.
type RadarCurveLayout struct {
	Label      string
	Points     [][2]float32
	ColorIndex int
}
```

Create `layout/radar.go`:

```go
package layout

import (
	"math"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeRadarLayout(g *ir.Graph, _ *theme.Theme, cfg *config.Layout) *Layout {
	radius := cfg.Radar.Radius
	padX := cfg.Radar.PaddingX
	padY := cfg.Radar.PaddingY
	labelOffset := cfg.Radar.LabelOffset

	numAxes := len(g.RadarAxes)
	if numAxes == 0 {
		numAxes = 1
	}

	// Title height.
	var titleHeight float32
	if g.RadarTitle != "" {
		titleHeight = 20 + padY
	}

	centerX := padX + radius + labelOffset
	centerY := titleHeight + padY + radius + labelOffset

	// Determine max value.
	maxVal := g.RadarMax
	if maxVal <= 0 {
		maxVal = radarAutoMax(g)
	}
	minVal := g.RadarMin

	valRange := maxVal - minVal
	if valRange <= 0 {
		valRange = 1
	}

	// Angle per axis (evenly distributed, starting from top = -pi/2).
	angleStep := 2 * math.Pi / float64(numAxes)

	// Build axis layouts.
	axes := make([]RadarAxisLayout, len(g.RadarAxes))
	for i, ax := range g.RadarAxes {
		angle := -math.Pi/2 + float64(i)*angleStep
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))
		axes[i] = RadarAxisLayout{
			Label:  ax.Label,
			EndX:   centerX + radius*cos,
			EndY:   centerY + radius*sin,
			LabelX: centerX + (radius+labelOffset)*cos,
			LabelY: centerY + (radius+labelOffset)*sin,
		}
	}

	// Build curve layouts.
	curves := make([]RadarCurveLayout, len(g.RadarCurves))
	for ci, curve := range g.RadarCurves {
		var points [][2]float32
		for i, v := range curve.Values {
			if i >= len(g.RadarAxes) {
				break
			}
			angle := -math.Pi/2 + float64(i)*angleStep
			frac := (v - minVal) / valRange
			if frac < 0 {
				frac = 0
			}
			if frac > 1 {
				frac = 1
			}
			r := float32(frac) * radius
			px := centerX + r*float32(math.Cos(angle))
			py := centerY + r*float32(math.Sin(angle))
			points = append(points, [2]float32{px, py})
		}
		curves[ci] = RadarCurveLayout{
			Label:      curve.Label,
			Points:     points,
			ColorIndex: ci,
		}
	}

	// Graticule radii (concentric rings).
	ticks := g.RadarTicks
	if ticks <= 0 {
		ticks = cfg.Radar.DefaultTicks
	}
	graticuleRadii := make([]float32, ticks)
	for i := range ticks {
		graticuleRadii[i] = radius * float32(i+1) / float32(ticks)
	}

	totalW := (padX + labelOffset + radius) * 2
	totalH := titleHeight + (padY+labelOffset+radius)*2

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: RadarData{
			Axes:           axes,
			Curves:         curves,
			GraticuleRadii: graticuleRadii,
			GraticuleType:  g.RadarGraticuleType,
			CenterX:        centerX,
			CenterY:        centerY,
			Radius:         radius,
			Title:          g.RadarTitle,
			ShowLegend:     g.RadarShowLegend,
			MaxValue:       maxVal,
			MinValue:       minVal,
		},
	}
}

func radarAutoMax(g *ir.Graph) float64 {
	max := 0.0
	for _, c := range g.RadarCurves {
		for _, v := range c.Values {
			if v > max {
				max = v
			}
		}
	}
	if max <= 0 {
		return 100
	}
	return niceMax(max)
}
```

Add to `layout/layout.go` ComputeLayout switch:

```go
	case ir.Radar:
		return computeRadarLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./layout/ -run TestRadar -v`
Expected: PASS

**Step 5: Commit**

```bash
git add layout/radar.go layout/radar_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Radar chart layout"
```

---

### Task 8: XYChart Renderer

**Files:**
- Create: `render/xychart.go`
- Create: `render/xychart_test.go`
- Modify: `render/svg.go` — add `case layout.XYChartData`

**Step 1: Write the failing test**

Create `render/xychart_test.go`:

```go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderXYChart(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart
	g.XYTitle = "Revenue"
	g.XYXAxis = &ir.XYAxis{
		Mode:       ir.XYAxisBand,
		Categories: []string{"Q1", "Q2", "Q3"},
	}
	g.XYYAxis = &ir.XYAxis{Mode: ir.XYAxisNumeric, Min: 0, Max: 100}
	g.XYSeries = []*ir.XYSeries{
		{Type: ir.XYSeriesBar, Values: []float64{30, 60, 90}},
		{Type: ir.XYSeriesLine, Values: []float64{25, 55, 85}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Revenue") {
		t.Error("missing title text")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing bar rects")
	}
	if !strings.Contains(svg, "<polyline") || !strings.Contains(svg, "<circle") {
		t.Error("missing line series elements")
	}
}

func TestRenderXYChartValidSVG(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.XYChart
	g.XYSeries = []*ir.XYSeries{
		{Type: ir.XYSeriesLine, Values: []float64{1, 2, 3}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG doesn't start with <svg")
	}
	if !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
		t.Error("SVG doesn't end with </svg>")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./render/ -run TestRenderXYChart -v`
Expected: FAIL — falls through to default `renderGraph`

**Step 3: Write minimal implementation**

Create `render/xychart.go`:

```go
package render

import (
	"fmt"
	"strings"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func renderXYChart(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	xyd, ok := l.Diagram.(layout.XYChartData)
	if !ok {
		return
	}

	// Title.
	if xyd.Title != "" {
		b.text(l.Width/2, cfg.XYChart.TitleFontSize, xyd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(cfg.XYChart.TitleFontSize),
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	cx := xyd.ChartX
	cy := xyd.ChartY
	cw := xyd.ChartWidth
	ch := xyd.ChartHeight

	// Grid lines and Y-axis ticks.
	gridColor := th.XYChartGridColor
	if gridColor == "" {
		gridColor = "#E0E0E0"
	}
	axisColor := th.XYChartAxisColor
	if axisColor == "" {
		axisColor = "#333"
	}

	for _, tick := range xyd.YTicks {
		// Horizontal grid line.
		b.line(cx, tick.Y, cx+cw, tick.Y,
			"stroke", gridColor, "stroke-width", "0.5")
		// Tick label.
		b.text(cx-4, tick.Y+4, tick.Label,
			"text-anchor", "end",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(cfg.XYChart.AxisFontSize),
			"fill", th.TextColor,
		)
	}

	// Axis lines.
	b.line(cx, cy, cx, cy+ch, "stroke", axisColor, "stroke-width", "1") // Y-axis
	b.line(cx, cy+ch, cx+cw, cy+ch, "stroke", axisColor, "stroke-width", "1") // X-axis

	// X-axis labels.
	for _, label := range xyd.XLabels {
		b.text(label.X, cy+ch+cfg.XYChart.AxisFontSize+4, label.Text,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(cfg.XYChart.AxisFontSize),
			"fill", th.TextColor,
		)
	}

	// Render each series.
	for _, s := range xyd.Series {
		color := "#4C78A8" // fallback
		if len(th.XYChartColors) > 0 {
			color = th.XYChartColors[s.ColorIndex%len(th.XYChartColors)]
		}

		switch s.Type {
		case layout.XYSeriesBar:
			for _, p := range s.Points {
				b.rect(p.X, p.Y, p.Width, p.Height, 0,
					"fill", color,
					"stroke", "none",
				)
			}
		case layout.XYSeriesLine:
			// Polyline.
			var pointStrs []string
			for _, p := range s.Points {
				pointStrs = append(pointStrs, fmt.Sprintf("%s,%s", fmtFloat(p.X), fmtFloat(p.Y)))
			}
			if len(pointStrs) > 0 {
				b.selfClose("polyline",
					"points", strings.Join(pointStrs, " "),
					"fill", "none",
					"stroke", color,
					"stroke-width", "2",
				)
			}
			// Data point circles.
			for _, p := range s.Points {
				b.circle(p.X, p.Y, 3,
					"fill", color,
					"stroke", th.Background,
					"stroke-width", "1",
				)
			}
		}
	}
}
```

Note: We need to re-export the `ir.XYSeriesType` constants via the layout package. Instead, update the `XYSeriesLayout` to use `ir.XYSeriesType` and import `ir` in the renderer. The type assertion in `render/xychart.go` uses `layout.XYChartData`, and the series `Type` field is `ir.XYSeriesType` which is already importable. Replace `layout.XYSeriesBar` / `layout.XYSeriesLine` with `ir.XYSeriesBar` / `ir.XYSeriesLine` in the renderer.

Add to `render/svg.go` switch (after `case layout.GitGraphData`):

```go
	case layout.XYChartData:
		renderXYChart(&b, l, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./render/ -run TestRenderXYChart -v`
Expected: PASS

**Step 5: Commit**

```bash
git add render/xychart.go render/xychart_test.go render/svg.go
git commit -m "feat(render): add XYChart SVG renderer"
```

---

### Task 9: Radar Renderer

**Files:**
- Create: `render/radar.go`
- Create: `render/radar_test.go`
- Modify: `render/svg.go` — add `case layout.RadarData`

**Step 1: Write the failing test**

Create `render/radar_test.go`:

```go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderRadar(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Radar
	g.RadarTitle = "Skills"
	g.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "Speed"},
		{ID: "b", Label: "Power"},
		{ID: "c", Label: "Magic"},
	}
	g.RadarCurves = []*ir.RadarCurve{
		{ID: "p1", Label: "Player1", Values: []float64{80, 60, 40}},
	}
	g.RadarMax = 100

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Skills") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "<polygon") || !strings.Contains(svg, "<line") {
		t.Error("missing radar elements (polygon or axis lines)")
	}
}

func TestRenderRadarValidSVG(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Radar
	g.RadarAxes = []*ir.RadarAxis{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
		{ID: "c", Label: "C"},
	}
	g.RadarCurves = []*ir.RadarCurve{
		{ID: "x", Label: "X", Values: []float64{50, 50, 50}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG doesn't start with <svg")
	}
	if !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
		t.Error("SVG doesn't end with </svg>")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./render/ -run TestRenderRadar -v`
Expected: FAIL — falls through to default

**Step 3: Write minimal implementation**

Create `render/radar.go`:

```go
package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func renderRadar(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	rd, ok := l.Diagram.(layout.RadarData)
	if !ok {
		return
	}

	cx := rd.CenterX
	cy := rd.CenterY
	numAxes := len(rd.Axes)

	// Title.
	if rd.Title != "" {
		b.text(l.Width/2, 20, rd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", "16",
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	graticuleColor := th.RadarGraticuleColor
	if graticuleColor == "" {
		graticuleColor = "#E0E0E0"
	}
	axisColor := th.RadarAxisColor
	if axisColor == "" {
		axisColor = "#333"
	}

	// Graticule (concentric rings).
	for _, r := range rd.GraticuleRadii {
		if rd.GraticuleType == ir.RadarGraticulePolygon && numAxes >= 3 {
			// Polygon graticule.
			var points []string
			angleStep := 2 * math.Pi / float64(numAxes)
			for i := range numAxes {
				angle := -math.Pi/2 + float64(i)*angleStep
				px := cx + r*float32(math.Cos(angle))
				py := cy + r*float32(math.Sin(angle))
				points = append(points, fmt.Sprintf("%s,%s", fmtFloat(px), fmtFloat(py)))
			}
			b.selfClose("polygon",
				"points", strings.Join(points, " "),
				"fill", "none",
				"stroke", graticuleColor,
				"stroke-width", "0.5",
			)
		} else {
			// Circle graticule.
			b.circle(cx, cy, r,
				"fill", "none",
				"stroke", graticuleColor,
				"stroke-width", "0.5",
			)
		}
	}

	// Axis lines.
	for _, ax := range rd.Axes {
		b.line(cx, cy, ax.EndX, ax.EndY,
			"stroke", axisColor, "stroke-width", "1")
		// Axis label.
		anchor := "middle"
		if ax.LabelX > cx+5 {
			anchor = "start"
		} else if ax.LabelX < cx-5 {
			anchor = "end"
		}
		b.text(ax.LabelX, ax.LabelY, ax.Label,
			"text-anchor", anchor,
			"dominant-baseline", "middle",
			"font-family", th.FontFamily,
			"font-size", "12",
			"fill", th.TextColor,
		)
	}

	// Curve polygons.
	curveOpacity := fmt.Sprintf("%.2f", cfg.Radar.CurveOpacity)
	for _, curve := range rd.Curves {
		color := "#4C78A8" // fallback
		if len(th.RadarCurveColors) > 0 {
			color = th.RadarCurveColors[curve.ColorIndex%len(th.RadarCurveColors)]
		}

		var points []string
		for _, p := range curve.Points {
			points = append(points, fmt.Sprintf("%s,%s", fmtFloat(p[0]), fmtFloat(p[1])))
		}
		if len(points) > 0 {
			b.selfClose("polygon",
				"points", strings.Join(points, " "),
				"fill", color,
				"fill-opacity", curveOpacity,
				"stroke", color,
				"stroke-width", "2",
			)
		}
	}

	// Legend.
	if rd.ShowLegend && len(rd.Curves) > 0 {
		legendX := l.Width - cfg.Radar.PaddingX - 100
		legendY := cfg.Radar.PaddingY
		for i, curve := range rd.Curves {
			color := "#4C78A8"
			if len(th.RadarCurveColors) > 0 {
				color = th.RadarCurveColors[curve.ColorIndex%len(th.RadarCurveColors)]
			}
			y := legendY + float32(i)*20
			b.rect(legendX, y, 12, 12, 0, "fill", color)
			b.text(legendX+16, y+10, curve.Label,
				"font-family", th.FontFamily,
				"font-size", "12",
				"fill", th.TextColor,
			)
		}
	}
}
```

Add to `render/svg.go` switch (after XYChartData):

```go
	case layout.RadarData:
		renderRadar(&b, l, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./render/ -run TestRenderRadar -v`
Expected: PASS

**Step 5: Commit**

```bash
git add render/radar.go render/radar_test.go render/svg.go
git commit -m "feat(render): add Radar chart SVG renderer"
```

---

### Task 10: Integration Tests and Fixtures

**Files:**
- Create: `testdata/fixtures/xychart-basic.mmd`
- Create: `testdata/fixtures/xychart-horizontal.mmd`
- Create: `testdata/fixtures/radar-basic.mmd`
- Create: `testdata/fixtures/radar-polygon.mmd`
- Modify: `mermaid_test.go` — add integration tests

**Step 1: Create fixture files**

`testdata/fixtures/xychart-basic.mmd`:
```
xychart-beta
    title "Sales Revenue"
    x-axis [jan, feb, mar, apr, may, jun]
    y-axis "Revenue ($)" 0 --> 12000
    bar [5000, 6000, 7500, 8200, 9500, 10500]
    line [5000, 6000, 7500, 8200, 9500, 10500]
```

`testdata/fixtures/xychart-horizontal.mmd`:
```
xychart-beta horizontal
    title "Performance"
    x-axis [A, B, C]
    bar [30, 60, 90]
```

`testdata/fixtures/radar-basic.mmd`:
```
radar-beta
    title "Language Skills"
    axis e["English"], f["French"], g["German"], s["Spanish"], d["Dutch"]
    curve a["User1"]{80, 60, 70, 50, 40}
    curve b["User2"]{60, 90, 50, 80, 70}
```

`testdata/fixtures/radar-polygon.mmd`:
```
radar-beta
    title "Team Comparison"
    graticule polygon
    ticks 4
    max 100
    showLegend
    axis sp["Speed"], po["Power"], st["Stamina"], sk["Skill"]
    curve a["Team A"]{80, 60, 70, 90}
    curve b["Team B"]{60, 90, 80, 50}
```

**Step 2: Write integration tests**

Add to `mermaid_test.go`:

```go
func TestRenderXYChartFixture(t *testing.T) {
	svg := renderFixture(t, "testdata/fixtures/xychart-basic.mmd")
	assertContains(t, svg, "<svg")
	assertContains(t, svg, "Sales Revenue")
	assertContains(t, svg, "<rect")    // bars
	assertContains(t, svg, "<polyline") // line
}

func TestRenderXYChartHorizontalFixture(t *testing.T) {
	svg := renderFixture(t, "testdata/fixtures/xychart-horizontal.mmd")
	assertContains(t, svg, "<svg")
	assertContains(t, svg, "Performance")
}

func TestRenderRadarFixture(t *testing.T) {
	svg := renderFixture(t, "testdata/fixtures/radar-basic.mmd")
	assertContains(t, svg, "<svg")
	assertContains(t, svg, "Language Skills")
	assertContains(t, svg, "<polygon") // curves
	assertContains(t, svg, "<line")    // axes
}

func TestRenderRadarPolygonFixture(t *testing.T) {
	svg := renderFixture(t, "testdata/fixtures/radar-polygon.mmd")
	assertContains(t, svg, "<svg")
	assertContains(t, svg, "Team Comparison")
	assertContains(t, svg, "<polygon") // polygon graticule + curves
}
```

**Step 3: Run tests to verify they pass**

Run: `go test ./... -v`
Expected: ALL PASS

**Step 4: Commit**

```bash
git add testdata/fixtures/xychart-basic.mmd testdata/fixtures/xychart-horizontal.mmd \
    testdata/fixtures/radar-basic.mmd testdata/fixtures/radar-polygon.mmd mermaid_test.go
git commit -m "test: add integration tests and fixtures for XYChart and Radar"
```

---

### Task 11: Final Validation

**Step 1: Run full test suite**

Run: `go test ./...`
Expected: All 8 packages PASS

**Step 2: Run go vet**

Run: `go vet ./...`
Expected: Clean

**Step 3: Run gofmt check**

Run: `gofmt -l .`
Expected: No output (all files formatted)

**Step 4: Run go build**

Run: `go build ./...`
Expected: Clean

**Step 5: Run go-code-reviewer**

Dispatch the `go-code-reviewer` agent to review all new files for quality issues.

**Step 6: Fix any issues found**

Apply fixes, commit.

---
