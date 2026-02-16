# Phase 5: Pie & Quadrant Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Pie chart and Quadrant chart diagram support to gomd2svg with full structural parity to mermaid.js.

**Architecture:** Both diagrams use positioned geometry — no Sugiyama, no grid. Pie computes arc angles from percentage data; Quadrant maps normalized [0,1] coordinates to pixel space. Each gets IR types, parser, config, layout, and renderer following the established per-diagram-type pattern.

**Tech Stack:** Go stdlib (`math`, `fmt`, `regexp`, `strconv`, `strings`), existing `textmetrics`, `theme`, `config` packages.

---

### Task 1: IR types — Pie

**Files:**
- Create: `ir/pie.go`
- Create: `ir/pie_test.go`
- Modify: `ir/graph.go:72-76` (add Pie fields after Kanban/Packet)

**Step 1: Write the test**

```go
// ir/pie_test.go
package ir

import "testing"

func TestPieSliceDefaults(t *testing.T) {
	s := &PieSlice{Label: "Dogs", Value: 386}
	if s.Label != "Dogs" {
		t.Errorf("Label = %q, want %q", s.Label, "Dogs")
	}
	if s.Value != 386 {
		t.Errorf("Value = %f, want 386", s.Value)
	}
}

func TestGraphPieFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Pie
	g.PieTitle = "Pets"
	g.PieShowData = true
	g.PieSlices = append(g.PieSlices, &PieSlice{Label: "Dogs", Value: 386})

	if g.PieTitle != "Pets" {
		t.Errorf("PieTitle = %q, want %q", g.PieTitle, "Pets")
	}
	if !g.PieShowData {
		t.Error("PieShowData = false, want true")
	}
	if len(g.PieSlices) != 1 {
		t.Fatalf("PieSlices = %d, want 1", len(g.PieSlices))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestPieSlice -v && go test ./ir/ -run TestGraphPieFields -v`
Expected: FAIL — `PieSlice` undefined

**Step 3: Write minimal implementation**

```go
// ir/pie.go
package ir

// PieSlice represents a single slice of a pie chart.
type PieSlice struct {
	Label string
	Value float64
}
```

Add to `ir/graph.go` after the Packet diagram fields block:

```go
	// Pie diagram fields
	PieSlices   []*PieSlice
	PieTitle    string
	PieShowData bool
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestPieSlice -v && go test ./ir/ -run TestGraphPieFields -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/pie.go ir/pie_test.go ir/graph.go
git commit -m "feat(ir): add Pie diagram types"
```

---

### Task 2: IR types — Quadrant

**Files:**
- Create: `ir/quadrant.go`
- Create: `ir/quadrant_test.go`
- Modify: `ir/graph.go` (add Quadrant fields after Pie)

**Step 1: Write the test**

```go
// ir/quadrant_test.go
package ir

import "testing"

func TestQuadrantPointDefaults(t *testing.T) {
	p := &QuadrantPoint{Label: "Campaign A", X: 0.3, Y: 0.6}
	if p.Label != "Campaign A" {
		t.Errorf("Label = %q, want %q", p.Label, "Campaign A")
	}
	if p.X != 0.3 || p.Y != 0.6 {
		t.Errorf("X,Y = %f,%f, want 0.3,0.6", p.X, p.Y)
	}
}

func TestGraphQuadrantFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Quadrant
	g.QuadrantTitle = "Campaigns"
	g.QuadrantLabels = [4]string{"Expand", "Promote", "Re-evaluate", "Improve"}
	g.XAxisLeft = "Low Reach"
	g.XAxisRight = "High Reach"
	g.YAxisBottom = "Low Engagement"
	g.YAxisTop = "High Engagement"
	g.QuadrantPoints = append(g.QuadrantPoints, &QuadrantPoint{Label: "A", X: 0.3, Y: 0.6})

	if g.QuadrantTitle != "Campaigns" {
		t.Errorf("QuadrantTitle = %q, want %q", g.QuadrantTitle, "Campaigns")
	}
	if len(g.QuadrantPoints) != 1 {
		t.Fatalf("QuadrantPoints = %d, want 1", len(g.QuadrantPoints))
	}
	if g.QuadrantLabels[0] != "Expand" {
		t.Errorf("QuadrantLabels[0] = %q, want %q", g.QuadrantLabels[0], "Expand")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestQuadrant -v`
Expected: FAIL — `QuadrantPoint` undefined

**Step 3: Write minimal implementation**

```go
// ir/quadrant.go
package ir

// QuadrantPoint represents a data point in a quadrant chart.
// X and Y are normalized values in the range [0, 1].
type QuadrantPoint struct {
	Label string
	X     float64
	Y     float64
}
```

Add to `ir/graph.go` after Pie fields:

```go
	// Quadrant diagram fields
	QuadrantPoints []*QuadrantPoint
	QuadrantTitle  string
	QuadrantLabels [4]string // q1=top-right, q2=top-left, q3=bottom-left, q4=bottom-right
	XAxisLeft      string
	XAxisRight     string
	YAxisBottom    string
	YAxisTop       string
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestQuadrant -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/quadrant.go ir/quadrant_test.go ir/graph.go
git commit -m "feat(ir): add Quadrant diagram types"
```

---

### Task 3: Config — Pie and Quadrant

**Files:**
- Modify: `config/config.go:4-17` (add PieConfig and QuadrantConfig structs and Layout fields)
- Modify: `config/config.go:84-139` (add defaults in DefaultLayout)

**Step 1: Write the test**

Add to a new file or existing config test. Since `config/` has no test file yet, create one:

```go
// config/config_test.go
package config

import "testing"

func TestDefaultLayoutPieConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Pie.Radius != 150 {
		t.Errorf("Pie.Radius = %f, want 150", cfg.Pie.Radius)
	}
	if cfg.Pie.TextPosition != 0.75 {
		t.Errorf("Pie.TextPosition = %f, want 0.75", cfg.Pie.TextPosition)
	}
	if cfg.Pie.PaddingX != 20 {
		t.Errorf("Pie.PaddingX = %f, want 20", cfg.Pie.PaddingX)
	}
}

func TestDefaultLayoutQuadrantConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Quadrant.ChartWidth != 400 {
		t.Errorf("Quadrant.ChartWidth = %f, want 400", cfg.Quadrant.ChartWidth)
	}
	if cfg.Quadrant.ChartHeight != 400 {
		t.Errorf("Quadrant.ChartHeight = %f, want 400", cfg.Quadrant.ChartHeight)
	}
	if cfg.Quadrant.PointRadius != 5 {
		t.Errorf("Quadrant.PointRadius = %f, want 5", cfg.Quadrant.PointRadius)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./config/ -v`
Expected: FAIL — `cfg.Pie` undefined

**Step 3: Write minimal implementation**

Add structs to `config/config.go`:

```go
// PieConfig holds pie chart layout options.
type PieConfig struct {
	Radius       float32
	InnerRadius  float32
	TextPosition float32
	PaddingX     float32
	PaddingY     float32
}

// QuadrantConfig holds quadrant chart layout options.
type QuadrantConfig struct {
	ChartWidth            float32
	ChartHeight           float32
	PointRadius           float32
	PaddingX              float32
	PaddingY              float32
	QuadrantLabelFontSize float32
	AxisLabelFontSize     float32
}
```

Add fields to `Layout` struct:

```go
	Pie      PieConfig
	Quadrant QuadrantConfig
```

Add defaults in `DefaultLayout()`:

```go
		Pie: PieConfig{
			Radius:       150,
			InnerRadius:  0,
			TextPosition: 0.75,
			PaddingX:     20,
			PaddingY:     20,
		},
		Quadrant: QuadrantConfig{
			ChartWidth:            400,
			ChartHeight:           400,
			PointRadius:           5,
			PaddingX:              40,
			PaddingY:              40,
			QuadrantLabelFontSize: 14,
			AxisLabelFontSize:     12,
		},
```

**Step 4: Run test to verify it passes**

Run: `go test ./config/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add config/config.go config/config_test.go
git commit -m "feat(config): add Pie and Quadrant config"
```

---

### Task 4: Theme — Pie color palette and Quadrant fills

**Files:**
- Modify: `theme/theme.go:5-69` (add PieColors and Quadrant fields to Theme struct)
- Modify: `theme/theme.go:72-139` (add values in Modern())
- Modify: `theme/theme.go:142-209` (add values in MermaidDefault())

**Step 1: Write the test**

```go
// theme/theme_test.go — add to existing or create
package theme

import "testing"

func TestModernPieColors(t *testing.T) {
	th := Modern()
	if len(th.PieColors) < 8 {
		t.Errorf("PieColors = %d, want >= 8", len(th.PieColors))
	}
}

func TestModernQuadrantFills(t *testing.T) {
	th := Modern()
	if th.QuadrantFill1 == "" {
		t.Error("QuadrantFill1 is empty")
	}
	if th.QuadrantFill2 == "" {
		t.Error("QuadrantFill2 is empty")
	}
	if th.QuadrantFill3 == "" {
		t.Error("QuadrantFill3 is empty")
	}
	if th.QuadrantFill4 == "" {
		t.Error("QuadrantFill4 is empty")
	}
	if th.QuadrantPointFill == "" {
		t.Error("QuadrantPointFill is empty")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./theme/ -v`
Expected: FAIL — `PieColors` undefined

**Step 3: Write minimal implementation**

Add to `Theme` struct in `theme/theme.go`:

```go
	PieColors []string

	// Quadrant chart colors
	QuadrantFill1     string
	QuadrantFill2     string
	QuadrantFill3     string
	QuadrantFill4     string
	QuadrantPointFill string
```

Add to `Modern()`:

```go
		PieColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},

		QuadrantFill1:     "#E8EFF5",
		QuadrantFill2:     "#F0F4F8",
		QuadrantFill3:     "#F5F5F5",
		QuadrantFill4:     "#FFF8E1",
		QuadrantPointFill: "#4C78A8",
```

Add to `MermaidDefault()`:

```go
		PieColors: []string{
			"#4C78A8", "#48A9A6", "#E4E36A", "#F4A261",
			"#E76F51", "#7FB069", "#D08AC0", "#F7B7A3",
		},

		QuadrantFill1:     "#f0ece8",
		QuadrantFill2:     "#e8f0e8",
		QuadrantFill3:     "#f0f0f0",
		QuadrantFill4:     "#ece8f0",
		QuadrantPointFill: "#4C78A8",
```

**Step 4: Run test to verify it passes**

Run: `go test ./theme/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add theme/theme.go theme/theme_test.go
git commit -m "feat(theme): add Pie color palette and Quadrant fills"
```

---

### Task 5: Parser — Pie

**Files:**
- Create: `parser/pie.go`
- Create: `parser/pie_test.go`
- Modify: `parser/parser.go:18-36` (add `case ir.Pie` to Parse switch)

**Step 1: Write the test**

```go
// parser/pie_test.go
package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParsePieBasic(t *testing.T) {
	input := `pie
    title Pets adopted by volunteers
    "Dogs" : 386
    "Cats" : 85
    "Rats" : 15`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if out.Graph.Kind != ir.Pie {
		t.Errorf("Kind = %v, want Pie", out.Graph.Kind)
	}
	if out.Graph.PieTitle != "Pets adopted by volunteers" {
		t.Errorf("PieTitle = %q", out.Graph.PieTitle)
	}
	if out.Graph.PieShowData {
		t.Error("PieShowData = true, want false")
	}
	if len(out.Graph.PieSlices) != 3 {
		t.Fatalf("PieSlices = %d, want 3", len(out.Graph.PieSlices))
	}
	if out.Graph.PieSlices[0].Label != "Dogs" || out.Graph.PieSlices[0].Value != 386 {
		t.Errorf("slice[0] = %+v", out.Graph.PieSlices[0])
	}
}

func TestParsePieShowData(t *testing.T) {
	input := `pie showData
    title Budget
    "Engineering" : 45.50
    "Sales" : 25.25`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if !out.Graph.PieShowData {
		t.Error("PieShowData = false, want true")
	}
	if len(out.Graph.PieSlices) != 2 {
		t.Fatalf("PieSlices = %d, want 2", len(out.Graph.PieSlices))
	}
	if out.Graph.PieSlices[0].Value != 45.50 {
		t.Errorf("slice[0].Value = %f, want 45.5", out.Graph.PieSlices[0].Value)
	}
}

func TestParsePieComments(t *testing.T) {
	input := `pie
    %% This is a comment
    "A" : 10
    "B" : 20 %% trailing comment`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.PieSlices) != 2 {
		t.Fatalf("PieSlices = %d, want 2", len(out.Graph.PieSlices))
	}
}

func TestParsePieNoTitle(t *testing.T) {
	input := `pie
    "X" : 50
    "Y" : 50`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if out.Graph.PieTitle != "" {
		t.Errorf("PieTitle = %q, want empty", out.Graph.PieTitle)
	}
	if len(out.Graph.PieSlices) != 2 {
		t.Fatalf("PieSlices = %d, want 2", len(out.Graph.PieSlices))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParsePie -v`
Expected: FAIL — falls through to `parseFlowchart`

**Step 3: Write minimal implementation**

```go
// parser/pie.go
package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var pieDataRe = regexp.MustCompile(`^\s*"([^"]+)"\s*:\s*(\d+\.?\d*)\s*$`)

func parsePie(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Pie

	lines := preprocessInput(input)

	for _, line := range lines {
		lower := strings.ToLower(line)

		// Skip the declaration line, extract showData flag.
		if strings.HasPrefix(lower, "pie") {
			if strings.Contains(lower, "showdata") {
				g.PieShowData = true
			}
			continue
		}

		// Title line.
		if strings.HasPrefix(lower, "title ") {
			g.PieTitle = strings.TrimSpace(line[len("title "):])
			continue
		}

		// Data line: "Label" : value
		if m := pieDataRe.FindStringSubmatch(line); m != nil {
			val, _ := strconv.ParseFloat(m[2], 64) // regex guarantees digits
			g.PieSlices = append(g.PieSlices, &ir.PieSlice{
				Label: m[1],
				Value: val,
			})
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}
```

Add to `parser/parser.go` Parse switch:

```go
	case ir.Pie:
		return parsePie(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParsePie -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/pie.go parser/pie_test.go parser/parser.go
git commit -m "feat(parser): add Pie chart parser"
```

---

### Task 6: Parser — Quadrant

**Files:**
- Create: `parser/quadrant.go`
- Create: `parser/quadrant_test.go`
- Modify: `parser/parser.go:18-36` (add `case ir.Quadrant` to Parse switch)

**Step 1: Write the test**

```go
// parser/quadrant_test.go
package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseQuadrantBasic(t *testing.T) {
	input := `quadrantChart
    title Reach and engagement of campaigns
    x-axis Low Reach --> High Reach
    y-axis Low Engagement --> High Engagement
    quadrant-1 We should expand
    quadrant-2 Need to promote
    quadrant-3 Re-evaluate
    quadrant-4 May be improved
    Campaign A: [0.3, 0.6]
    Campaign B: [0.45, 0.23]
    Campaign C: [0.57, 0.69]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Quadrant {
		t.Errorf("Kind = %v, want Quadrant", g.Kind)
	}
	if g.QuadrantTitle != "Reach and engagement of campaigns" {
		t.Errorf("Title = %q", g.QuadrantTitle)
	}
	if g.XAxisLeft != "Low Reach" || g.XAxisRight != "High Reach" {
		t.Errorf("XAxis = %q / %q", g.XAxisLeft, g.XAxisRight)
	}
	if g.YAxisBottom != "Low Engagement" || g.YAxisTop != "High Engagement" {
		t.Errorf("YAxis = %q / %q", g.YAxisBottom, g.YAxisTop)
	}
	if g.QuadrantLabels[0] != "We should expand" {
		t.Errorf("Q1 = %q", g.QuadrantLabels[0])
	}
	if g.QuadrantLabels[2] != "Re-evaluate" {
		t.Errorf("Q3 = %q", g.QuadrantLabels[2])
	}
	if len(g.QuadrantPoints) != 3 {
		t.Fatalf("Points = %d, want 3", len(g.QuadrantPoints))
	}
	p := g.QuadrantPoints[0]
	if p.Label != "Campaign A" || p.X != 0.3 || p.Y != 0.6 {
		t.Errorf("point[0] = %+v", p)
	}
}

func TestParseQuadrantSingleAxisLabel(t *testing.T) {
	input := `quadrantChart
    x-axis Effort
    y-axis Impact
    Task A: [0.5, 0.5]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.XAxisLeft != "Effort" || g.XAxisRight != "" {
		t.Errorf("XAxis = %q / %q", g.XAxisLeft, g.XAxisRight)
	}
	if g.YAxisBottom != "Impact" || g.YAxisTop != "" {
		t.Errorf("YAxis = %q / %q", g.YAxisBottom, g.YAxisTop)
	}
}

func TestParseQuadrantMinimal(t *testing.T) {
	input := `quadrantChart
    Point: [0.1, 0.9]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.QuadrantPoints) != 1 {
		t.Fatalf("Points = %d, want 1", len(out.Graph.QuadrantPoints))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParseQuadrant -v`
Expected: FAIL — falls through to `parseFlowchart`

**Step 3: Write minimal implementation**

```go
// parser/quadrant.go
package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var quadrantPointRe = regexp.MustCompile(`^\s*(.+?):\s*\[([0-9.]+),\s*([0-9.]+)\]\s*$`)

func parseQuadrant(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Quadrant

	lines := preprocessInput(input)

	for _, line := range lines {
		lower := strings.ToLower(line)

		if strings.HasPrefix(lower, "quadrantchart") {
			continue
		}

		// Title.
		if strings.HasPrefix(lower, "title ") {
			g.QuadrantTitle = strings.TrimSpace(line[len("title "):])
			continue
		}

		// X-axis: "x-axis Left --> Right" or "x-axis Label"
		if strings.HasPrefix(lower, "x-axis ") {
			rest := strings.TrimSpace(line[len("x-axis "):])
			if parts := strings.SplitN(rest, "-->", 2); len(parts) == 2 {
				g.XAxisLeft = strings.TrimSpace(parts[0])
				g.XAxisRight = strings.TrimSpace(parts[1])
			} else {
				g.XAxisLeft = rest
			}
			continue
		}

		// Y-axis: "y-axis Bottom --> Top" or "y-axis Label"
		if strings.HasPrefix(lower, "y-axis ") {
			rest := strings.TrimSpace(line[len("y-axis "):])
			if parts := strings.SplitN(rest, "-->", 2); len(parts) == 2 {
				g.YAxisBottom = strings.TrimSpace(parts[0])
				g.YAxisTop = strings.TrimSpace(parts[1])
			} else {
				g.YAxisBottom = rest
			}
			continue
		}

		// Quadrant labels.
		if strings.HasPrefix(lower, "quadrant-1 ") {
			g.QuadrantLabels[0] = strings.TrimSpace(line[len("quadrant-1 "):])
			continue
		}
		if strings.HasPrefix(lower, "quadrant-2 ") {
			g.QuadrantLabels[1] = strings.TrimSpace(line[len("quadrant-2 "):])
			continue
		}
		if strings.HasPrefix(lower, "quadrant-3 ") {
			g.QuadrantLabels[2] = strings.TrimSpace(line[len("quadrant-3 "):])
			continue
		}
		if strings.HasPrefix(lower, "quadrant-4 ") {
			g.QuadrantLabels[3] = strings.TrimSpace(line[len("quadrant-4 "):])
			continue
		}

		// Data point: "Label: [x, y]"
		if m := quadrantPointRe.FindStringSubmatch(line); m != nil {
			x, _ := strconv.ParseFloat(m[2], 64) // regex guarantees digits
			y, _ := strconv.ParseFloat(m[3], 64) // regex guarantees digits
			g.QuadrantPoints = append(g.QuadrantPoints, &ir.QuadrantPoint{
				Label: strings.TrimSpace(m[1]),
				X:     x,
				Y:     y,
			})
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}
```

Add to `parser/parser.go` Parse switch:

```go
	case ir.Quadrant:
		return parseQuadrant(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParseQuadrant -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/quadrant.go parser/quadrant_test.go parser/parser.go
git commit -m "feat(parser): add Quadrant chart parser"
```

---

### Task 7: Layout — Pie

**Files:**
- Create: `layout/pie.go`
- Create: `layout/pie_test.go`
- Modify: `layout/types.go` (add PieData and PieSliceLayout types)
- Modify: `layout/layout.go:13-31` (add `case ir.Pie` to ComputeLayout switch)

**Step 1: Write the test**

```go
// layout/pie_test.go
package layout

import (
	"math"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestPieLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Pie
	g.PieTitle = "Pets"
	g.PieSlices = []*ir.PieSlice{
		{Label: "Dogs", Value: 50},
		{Label: "Cats", Value: 50},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Pie {
		t.Errorf("Kind = %v, want Pie", l.Kind)
	}
	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("dimensions = %f x %f, want > 0", l.Width, l.Height)
	}

	pd, ok := l.Diagram.(PieData)
	if !ok {
		t.Fatalf("Diagram type = %T, want PieData", l.Diagram)
	}
	if len(pd.Slices) != 2 {
		t.Fatalf("Slices = %d, want 2", len(pd.Slices))
	}

	// Two equal slices: each should span pi radians.
	s0 := pd.Slices[0]
	s1 := pd.Slices[1]
	span0 := s0.EndAngle - s0.StartAngle
	span1 := s1.EndAngle - s1.StartAngle
	if math.Abs(float64(span0-span1)) > 0.01 {
		t.Errorf("spans differ: %f vs %f", span0, span1)
	}
	if math.Abs(float64(span0)-math.Pi) > 0.01 {
		t.Errorf("span0 = %f, want ~pi", span0)
	}
}

func TestPieLayoutSingleSlice(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Pie
	g.PieSlices = []*ir.PieSlice{
		{Label: "All", Value: 100},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	pd := l.Diagram.(PieData)
	if len(pd.Slices) != 1 {
		t.Fatalf("Slices = %d, want 1", len(pd.Slices))
	}
	span := pd.Slices[0].EndAngle - pd.Slices[0].StartAngle
	if math.Abs(float64(span)-2*math.Pi) > 0.01 {
		t.Errorf("span = %f, want ~2*pi", span)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestPieLayout -v`
Expected: FAIL — `PieData` undefined

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// PieData holds pie-chart-specific layout data.
type PieData struct {
	Slices   []PieSliceLayout
	CenterX  float32
	CenterY  float32
	Radius   float32
	Title    string
	ShowData bool
}

func (PieData) diagramData() {}

// PieSliceLayout holds computed angles and label position for one slice.
type PieSliceLayout struct {
	Label      string
	Value      float64
	Percentage float32
	StartAngle float32
	EndAngle   float32
	LabelX     float32
	LabelY     float32
	ColorIndex int
}
```

Create `layout/pie.go`:

```go
package layout

import (
	"math"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computePieLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	radius := cfg.Pie.Radius
	padX := cfg.Pie.PaddingX
	padY := cfg.Pie.PaddingY
	textPos := cfg.Pie.TextPosition

	// Compute total value.
	var total float64
	for _, s := range g.PieSlices {
		total += s.Value
	}
	if total <= 0 {
		total = 1
	}

	// Title height.
	var titleHeight float32
	if g.PieTitle != "" {
		titleHeight = th.PieTitleTextSize + padY
	}

	centerX := padX + radius
	centerY := titleHeight + padY + radius

	// Compute slice angles (clockwise from top = -pi/2).
	slices := make([]PieSliceLayout, len(g.PieSlices))
	var angle float32 = -math.Pi / 2 // start at top

	for i, s := range g.PieSlices {
		frac := float32(s.Value / total)
		span := frac * 2 * math.Pi

		midAngle := angle + span/2
		labelR := radius * textPos
		labelX := centerX + labelR*float32(math.Cos(float64(midAngle)))
		labelY := centerY + labelR*float32(math.Sin(float64(midAngle)))

		slices[i] = PieSliceLayout{
			Label:      s.Label,
			Value:      s.Value,
			Percentage: frac * 100,
			StartAngle: angle,
			EndAngle:   angle + span,
			LabelX:     labelX,
			LabelY:     labelY,
			ColorIndex: i,
		}

		angle += span
	}

	width := 2*padX + 2*radius
	height := titleHeight + 2*padY + 2*radius

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  width,
		Height: height,
		Diagram: PieData{
			Slices:   slices,
			CenterX:  centerX,
			CenterY:  centerY,
			Radius:   radius,
			Title:    g.PieTitle,
			ShowData: g.PieShowData,
		},
	}
}
```

Add to `layout/layout.go` ComputeLayout switch:

```go
	case ir.Pie:
		return computePieLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./layout/ -run TestPieLayout -v`
Expected: PASS

**Step 5: Commit**

```bash
git add layout/pie.go layout/pie_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Pie chart layout"
```

---

### Task 8: Layout — Quadrant

**Files:**
- Create: `layout/quadrant.go`
- Create: `layout/quadrant_test.go`
- Modify: `layout/types.go` (add QuadrantData and QuadrantPointLayout types)
- Modify: `layout/layout.go:13-31` (add `case ir.Quadrant` to ComputeLayout switch)

**Step 1: Write the test**

```go
// layout/quadrant_test.go
package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestQuadrantLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Quadrant
	g.QuadrantTitle = "Campaigns"
	g.XAxisLeft = "Low"
	g.XAxisRight = "High"
	g.YAxisBottom = "Low"
	g.YAxisTop = "High"
	g.QuadrantLabels = [4]string{"Q1", "Q2", "Q3", "Q4"}
	g.QuadrantPoints = []*ir.QuadrantPoint{
		{Label: "A", X: 0.0, Y: 0.0},
		{Label: "B", X: 1.0, Y: 1.0},
		{Label: "C", X: 0.5, Y: 0.5},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Quadrant {
		t.Errorf("Kind = %v, want Quadrant", l.Kind)
	}
	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("dimensions = %f x %f", l.Width, l.Height)
	}

	qd, ok := l.Diagram.(QuadrantData)
	if !ok {
		t.Fatalf("Diagram type = %T, want QuadrantData", l.Diagram)
	}
	if len(qd.Points) != 3 {
		t.Fatalf("Points = %d, want 3", len(qd.Points))
	}

	// Point A (0,0) should be at bottom-left; Point B (1,1) at top-right.
	// In pixel space: A.X < B.X and A.Y > B.Y (SVG Y is inverted).
	if qd.Points[0].X >= qd.Points[1].X {
		t.Errorf("A.X=%f >= B.X=%f, want A left of B", qd.Points[0].X, qd.Points[1].X)
	}
	if qd.Points[0].Y <= qd.Points[1].Y {
		t.Errorf("A.Y=%f <= B.Y=%f, want A below B (higher Y)", qd.Points[0].Y, qd.Points[1].Y)
	}

	// Point C (0.5,0.5) should be at center.
	midX := (qd.Points[0].X + qd.Points[1].X) / 2
	midY := (qd.Points[0].Y + qd.Points[1].Y) / 2
	if abs32(qd.Points[2].X-midX) > 1 {
		t.Errorf("C.X=%f not near midX=%f", qd.Points[2].X, midX)
	}
	if abs32(qd.Points[2].Y-midY) > 1 {
		t.Errorf("C.Y=%f not near midY=%f", qd.Points[2].Y, midY)
	}
}

func TestQuadrantLayoutNoPoints(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Quadrant
	g.QuadrantLabels = [4]string{"Q1", "Q2", "Q3", "Q4"}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	qd := l.Diagram.(QuadrantData)
	if len(qd.Points) != 0 {
		t.Errorf("Points = %d, want 0", len(qd.Points))
	}
	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("dimensions = %f x %f", l.Width, l.Height)
	}
}

func abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestQuadrantLayout -v`
Expected: FAIL — `QuadrantData` undefined

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// QuadrantData holds quadrant-chart-specific layout data.
type QuadrantData struct {
	Points        []QuadrantPointLayout
	ChartX        float32 // top-left X of the quadrant area
	ChartY        float32 // top-left Y of the quadrant area
	ChartWidth    float32
	ChartHeight   float32
	Title         string
	Labels        [4]string
	XAxisLeft     string
	XAxisRight    string
	YAxisBottom   string
	YAxisTop      string
}

func (QuadrantData) diagramData() {}

// QuadrantPointLayout holds the pixel position of a data point.
type QuadrantPointLayout struct {
	Label string
	X     float32
	Y     float32
}
```

Create `layout/quadrant.go`:

```go
package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeQuadrantLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	chartW := cfg.Quadrant.ChartWidth
	chartH := cfg.Quadrant.ChartHeight
	padX := cfg.Quadrant.PaddingX
	padY := cfg.Quadrant.PaddingY

	// Title height.
	var titleHeight float32
	if g.QuadrantTitle != "" {
		titleHeight = th.FontSize + padY
	}

	// Y-axis label width (left side).
	var yAxisLabelWidth float32
	if g.YAxisBottom != "" || g.YAxisTop != "" {
		yAxisLabelWidth = padX
	}

	// Chart origin.
	chartX := padX + yAxisLabelWidth
	chartY := titleHeight + padY

	// X-axis label height (below chart).
	var xAxisLabelHeight float32
	if g.XAxisLeft != "" || g.XAxisRight != "" {
		xAxisLabelHeight = cfg.Quadrant.AxisLabelFontSize + padY/2
	}

	// Map normalized points to pixel positions.
	points := make([]QuadrantPointLayout, len(g.QuadrantPoints))
	for i, p := range g.QuadrantPoints {
		// X maps directly: 0 = left, 1 = right.
		px := chartX + float32(p.X)*chartW
		// Y is inverted: 0 = bottom (high Y), 1 = top (low Y).
		py := chartY + (1-float32(p.Y))*chartH
		points[i] = QuadrantPointLayout{
			Label: p.Label,
			X:     px,
			Y:     py,
		}
	}

	totalW := padX + yAxisLabelWidth + chartW + padX
	totalH := titleHeight + padY + chartH + xAxisLabelHeight + padY

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: QuadrantData{
			Points:      points,
			ChartX:      chartX,
			ChartY:      chartY,
			ChartWidth:  chartW,
			ChartHeight: chartH,
			Title:       g.QuadrantTitle,
			Labels:      g.QuadrantLabels,
			XAxisLeft:   g.XAxisLeft,
			XAxisRight:  g.XAxisRight,
			YAxisBottom: g.YAxisBottom,
			YAxisTop:    g.YAxisTop,
		},
	}
}
```

Add to `layout/layout.go` ComputeLayout switch:

```go
	case ir.Quadrant:
		return computeQuadrantLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./layout/ -run TestQuadrantLayout -v`
Expected: PASS

**Step 5: Commit**

```bash
git add layout/quadrant.go layout/quadrant_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Quadrant chart layout"
```

---

### Task 9: Renderer — Pie

**Files:**
- Create: `render/pie.go`
- Create: `render/pie_test.go`
- Modify: `render/svg.go:40-58` (add `case layout.PieData` to RenderSVG switch)

**Step 1: Write the test**

```go
// render/pie_test.go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderPie(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Pie
	g.PieTitle = "Pets"
	g.PieSlices = []*ir.PieSlice{
		{Label: "Dogs", Value: 386},
		{Label: "Cats", Value: 85},
		{Label: "Rats", Value: 15},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "Pets") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Dogs") {
		t.Error("missing slice label Dogs")
	}
	if !strings.Contains(svg, "Cats") {
		t.Error("missing slice label Cats")
	}
	// Should contain arc paths.
	if !strings.Contains(svg, "<path") {
		t.Error("missing <path> for arcs")
	}
}

func TestRenderPieShowData(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Pie
	g.PieShowData = true
	g.PieSlices = []*ir.PieSlice{
		{Label: "A", Value: 60},
		{Label: "B", Value: 40},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// When showData is true, values should appear in the output.
	if !strings.Contains(svg, "60") {
		t.Error("missing value 60 with showData=true")
	}
}

func TestRenderPieSingleSlice(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Pie
	g.PieSlices = []*ir.PieSlice{
		{Label: "All", Value: 100},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Single slice = full circle.
	if !strings.Contains(svg, "<circle") || !strings.Contains(svg, "<path") {
		// Either a circle element or a path with two arcs is acceptable.
		if !strings.Contains(svg, "All") {
			t.Error("missing label for single slice")
		}
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./render/ -run TestRenderPie -v`
Expected: FAIL — falls through to default case

**Step 3: Write minimal implementation**

Create `render/pie.go`:

```go
package render

import (
	"fmt"
	"math"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderPie(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	pd := l.Diagram.(layout.PieData)

	// Title.
	if pd.Title != "" {
		b.text(l.Width/2, th.PieTitleTextSize, pd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.PieTitleTextSize),
			"font-weight", "bold",
			"fill", th.PieTitleTextColor,
		)
	}

	cx := pd.CenterX
	cy := pd.CenterY
	r := pd.Radius

	for _, s := range pd.Slices {
		color := th.PieColors[s.ColorIndex%len(th.PieColors)]
		opacity := fmt.Sprintf("%.2f", th.PieOpacity)

		// Full circle special case.
		span := s.EndAngle - s.StartAngle
		if span >= 2*math.Pi-0.01 {
			b.circle(cx, cy, r,
				"fill", color,
				"fill-opacity", opacity,
				"stroke", th.PieStrokeColor,
				"stroke-width", fmtFloat(th.PieStrokeWidth),
			)
		} else {
			// Arc path.
			x1 := cx + r*float32(math.Cos(float64(s.StartAngle)))
			y1 := cy + r*float32(math.Sin(float64(s.StartAngle)))
			x2 := cx + r*float32(math.Cos(float64(s.EndAngle)))
			y2 := cy + r*float32(math.Sin(float64(s.EndAngle)))

			largeArc := "0"
			if span > math.Pi {
				largeArc = "1"
			}

			d := fmt.Sprintf("M %s,%s L %s,%s A %s,%s 0 %s,1 %s,%s Z",
				fmtFloat(cx), fmtFloat(cy),
				fmtFloat(x1), fmtFloat(y1),
				fmtFloat(r), fmtFloat(r),
				largeArc,
				fmtFloat(x2), fmtFloat(y2),
			)

			b.path(d,
				"fill", color,
				"fill-opacity", opacity,
				"stroke", th.PieStrokeColor,
				"stroke-width", fmtFloat(th.PieStrokeWidth),
			)
		}

		// Slice label.
		labelText := s.Label
		if pd.ShowData {
			labelText = fmt.Sprintf("%s (%.0f)", s.Label, s.Value)
		}
		b.text(s.LabelX, s.LabelY, labelText,
			"text-anchor", "middle",
			"dominant-baseline", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.PieSectionTextSize),
			"fill", th.PieSectionTextColor,
		)
	}

	// Outer stroke.
	if th.PieOuterStrokeWidth > 0 {
		b.circle(cx, cy, r,
			"fill", "none",
			"stroke", th.PieOuterStrokeColor,
			"stroke-width", fmtFloat(th.PieOuterStrokeWidth),
		)
	}
}
```

Add to `render/svg.go` RenderSVG switch:

```go
	case layout.PieData:
		renderPie(&b, l, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./render/ -run TestRenderPie -v`
Expected: PASS

**Step 5: Commit**

```bash
git add render/pie.go render/pie_test.go render/svg.go
git commit -m "feat(render): add Pie chart SVG renderer"
```

---

### Task 10: Renderer — Quadrant

**Files:**
- Create: `render/quadrant.go`
- Create: `render/quadrant_test.go`
- Modify: `render/svg.go:40-58` (add `case layout.QuadrantData` to RenderSVG switch)

**Step 1: Write the test**

```go
// render/quadrant_test.go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderQuadrant(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Quadrant
	g.QuadrantTitle = "Campaigns"
	g.XAxisLeft = "Low Reach"
	g.XAxisRight = "High Reach"
	g.YAxisBottom = "Low Engagement"
	g.YAxisTop = "High Engagement"
	g.QuadrantLabels = [4]string{"Expand", "Promote", "Re-evaluate", "Improve"}
	g.QuadrantPoints = []*ir.QuadrantPoint{
		{Label: "Campaign A", X: 0.3, Y: 0.6},
		{Label: "Campaign B", X: 0.7, Y: 0.4},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "Campaigns") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Expand") {
		t.Error("missing quadrant label Q1")
	}
	if !strings.Contains(svg, "Low Reach") {
		t.Error("missing x-axis left label")
	}
	if !strings.Contains(svg, "High Engagement") {
		t.Error("missing y-axis top label")
	}
	if !strings.Contains(svg, "Campaign A") {
		t.Error("missing point label")
	}
	// Should contain circles for data points.
	if !strings.Contains(svg, "<circle") {
		t.Error("missing <circle> for data points")
	}
	// Should contain rects for quadrant backgrounds.
	count := strings.Count(svg, "<rect")
	// 1 background + 4 quadrant rects = at least 5
	if count < 5 {
		t.Errorf("rect count = %d, want >= 5", count)
	}
}

func TestRenderQuadrantMinimal(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Quadrant
	g.QuadrantPoints = []*ir.QuadrantPoint{
		{Label: "P", X: 0.5, Y: 0.5},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing data point circle")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./render/ -run TestRenderQuadrant -v`
Expected: FAIL — falls through to default case

**Step 3: Write minimal implementation**

Create `render/quadrant.go`:

```go
package render

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderQuadrant(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	qd := l.Diagram.(layout.QuadrantData)

	cx := qd.ChartX
	cy := qd.ChartY
	w := qd.ChartWidth
	h := qd.ChartHeight
	halfW := w / 2
	halfH := h / 2

	qLabelSize := cfg.Quadrant.QuadrantLabelFontSize
	axisLabelSize := cfg.Quadrant.AxisLabelFontSize

	fills := [4]string{th.QuadrantFill1, th.QuadrantFill2, th.QuadrantFill3, th.QuadrantFill4}

	// Title.
	if qd.Title != "" {
		b.text(l.Width/2, th.FontSize+cfg.Quadrant.PaddingY/2, qd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize),
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	// Quadrant background rects: Q2(top-left), Q1(top-right), Q3(bottom-left), Q4(bottom-right).
	// q1=top-right, q2=top-left, q3=bottom-left, q4=bottom-right
	b.rect(cx+halfW, cy, halfW, halfH, 0, "fill", fills[0]) // Q1 top-right
	b.rect(cx, cy, halfW, halfH, 0, "fill", fills[1])        // Q2 top-left
	b.rect(cx, cy+halfH, halfW, halfH, 0, "fill", fills[2])  // Q3 bottom-left
	b.rect(cx+halfW, cy+halfH, halfW, halfH, 0, "fill", fills[3]) // Q4 bottom-right

	// Quadrant border.
	b.rect(cx, cy, w, h, 0,
		"fill", "none",
		"stroke", th.LineColor,
		"stroke-width", "1",
	)

	// Center cross lines.
	b.line(cx+halfW, cy, cx+halfW, cy+h,
		"stroke", th.LineColor, "stroke-width", "0.5", "stroke-dasharray", "4,4")
	b.line(cx, cy+halfH, cx+w, cy+halfH,
		"stroke", th.LineColor, "stroke-width", "0.5", "stroke-dasharray", "4,4")

	// Quadrant labels centered in each quadrant.
	if qd.Labels[0] != "" {
		b.text(cx+halfW+halfW/2, cy+halfH/2, qd.Labels[0],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}
	if qd.Labels[1] != "" {
		b.text(cx+halfW/2, cy+halfH/2, qd.Labels[1],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}
	if qd.Labels[2] != "" {
		b.text(cx+halfW/2, cy+halfH+halfH/2, qd.Labels[2],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}
	if qd.Labels[3] != "" {
		b.text(cx+halfW+halfW/2, cy+halfH+halfH/2, qd.Labels[3],
			"text-anchor", "middle", "dominant-baseline", "middle",
			"font-family", th.FontFamily, "font-size", fmtFloat(qLabelSize),
			"fill", th.TextColor, "opacity", "0.6")
	}

	// Axis labels.
	if qd.XAxisLeft != "" {
		b.text(cx, cy+h+axisLabelSize+4, qd.XAxisLeft,
			"text-anchor", "start", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor)
	}
	if qd.XAxisRight != "" {
		b.text(cx+w, cy+h+axisLabelSize+4, qd.XAxisRight,
			"text-anchor", "end", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor)
	}
	if qd.YAxisBottom != "" {
		b.text(cx-4, cy+h, qd.YAxisBottom,
			"text-anchor", "end", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor,
			"transform", "rotate(-90,"+fmtFloat(cx-4)+","+fmtFloat(cy+h)+")")
	}
	if qd.YAxisTop != "" {
		b.text(cx-4, cy, qd.YAxisTop,
			"text-anchor", "start", "font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize), "fill", th.TextColor,
			"transform", "rotate(-90,"+fmtFloat(cx-4)+","+fmtFloat(cy)+")")
	}

	// Data points.
	pointR := cfg.Quadrant.PointRadius
	for _, p := range qd.Points {
		b.circle(p.X, p.Y, pointR,
			"fill", th.QuadrantPointFill,
			"stroke", th.LineColor,
			"stroke-width", "1",
		)
		b.text(p.X+pointR+3, p.Y+4, p.Label,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(axisLabelSize),
			"fill", th.TextColor,
		)
	}
}
```

Add to `render/svg.go` RenderSVG switch:

```go
	case layout.QuadrantData:
		renderQuadrant(&b, l, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./render/ -run TestRenderQuadrant -v`
Expected: PASS

**Step 5: Commit**

```bash
git add render/quadrant.go render/quadrant_test.go render/svg.go
git commit -m "feat(render): add Quadrant chart SVG renderer"
```

---

### Task 11: Integration tests and fixtures

**Files:**
- Create: `testdata/fixtures/pie-basic.mmd`
- Create: `testdata/fixtures/pie-showdata.mmd`
- Create: `testdata/fixtures/quadrant-campaigns.mmd`
- Create: `testdata/fixtures/quadrant-minimal.mmd`

**Step 1: Create fixture files**

`testdata/fixtures/pie-basic.mmd`:
```
pie title Pets adopted by volunteers
    "Dogs" : 386
    "Cats" : 85
    "Rats" : 15
```

`testdata/fixtures/pie-showdata.mmd`:
```
pie showData
    title Key elements in Product X
    "Calcium" : 42.96
    "Potassium" : 50.05
    "Magnesium" : 10.01
    "Iron" : 5
```

`testdata/fixtures/quadrant-campaigns.mmd`:
```
quadrantChart
    title Reach and engagement of campaigns
    x-axis Low Reach --> High Reach
    y-axis Low Engagement --> High Engagement
    quadrant-1 We should expand
    quadrant-2 Need to promote
    quadrant-3 Re-evaluate
    quadrant-4 May be improved
    Campaign A: [0.3, 0.6]
    Campaign B: [0.45, 0.23]
    Campaign C: [0.57, 0.69]
    Campaign D: [0.78, 0.34]
    Campaign E: [0.40, 0.34]
    Campaign F: [0.35, 0.78]
```

`testdata/fixtures/quadrant-minimal.mmd`:
```
quadrantChart
    Point A: [0.1, 0.9]
    Point B: [0.9, 0.1]
```

**Step 2: Write integration test**

Add to the top-level integration test file (create if needed):

```go
// mermaid_test.go (top-level package)
// Add test cases for Pie and Quadrant to the existing integration test pattern.

func TestRenderPieFixture(t *testing.T) {
	input := readFixture(t, "testdata/fixtures/pie-basic.mmd")
	svg, err := mermaid.Render(input, nil)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if !strings.Contains(svg, "Dogs") {
		t.Error("missing Dogs label")
	}
	if !strings.Contains(svg, "<path") {
		t.Error("missing arc path")
	}
}

func TestRenderQuadrantFixture(t *testing.T) {
	input := readFixture(t, "testdata/fixtures/quadrant-campaigns.mmd")
	svg, err := mermaid.Render(input, nil)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if !strings.Contains(svg, "Campaign A") {
		t.Error("missing Campaign A")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing data point circles")
	}
}
```

**Step 3: Run tests**

Run: `go test ./... -v`
Expected: ALL PASS

**Step 4: Commit**

```bash
git add testdata/fixtures/pie-basic.mmd testdata/fixtures/pie-showdata.mmd \
       testdata/fixtures/quadrant-campaigns.mmd testdata/fixtures/quadrant-minimal.mmd \
       mermaid_test.go
git commit -m "test: add integration tests and fixtures for Pie and Quadrant"
```

---

### Task 12: Final validation

**Step 1: Run full test suite**

Run: `go test ./... -v`
Expected: ALL PASS

**Step 2: Run go vet and gofmt**

Run: `go vet ./... && gofmt -l .`
Expected: Clean

**Step 3: Verify build**

Run: `go build ./...`
Expected: Clean
