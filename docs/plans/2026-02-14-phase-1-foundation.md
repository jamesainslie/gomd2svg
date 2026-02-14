# Phase 1: Foundation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Render `flowchart LR; A-->B-->C` to valid SVG end-to-end, cross-validated against the Rust reference implementation.

**Architecture:** 4-stage pipeline: Parse -> IR -> Layout -> Render SVG. Each stage is a separate Go package. The parser detects diagram kind and dispatches to flowchart-specific parsing. The layout engine sizes nodes, assigns ranks, orders within ranks, positions in coordinate space, and routes edges with simple L-shaped paths (A* deferred to Phase 5). The renderer generates SVG via strings.Builder.

**Tech Stack:** Go 1.23+, `golang.org/x/image/font/sfnt` for text metrics, stdlib `regexp` for parsing, no other external dependencies.

**Reference:** Rust source at `/Volumes/Development/mermaid-rs-renderer/src/`

---

### Task 1: IR types — enums and core structs

**Files:**
- Create: `ir/diagram.go`
- Create: `ir/shapes.go`
- Create: `ir/graph.go`
- Create: `ir/style.go`
- Test: `ir/graph_test.go`

**Step 1: Write the failing test**

```go
// ir/graph_test.go
package ir

import "testing"

func TestNewGraph(t *testing.T) {
    g := NewGraph()
    if g.Kind != Flowchart {
        t.Errorf("Kind = %v, want Flowchart", g.Kind)
    }
    if g.Direction != TopDown {
        t.Errorf("Direction = %v, want TopDown", g.Direction)
    }
    if g.Nodes == nil {
        t.Error("Nodes is nil")
    }
    if g.Edges != nil {
        t.Error("Edges should be nil (zero-value slice)")
    }
}

func TestEnsureNode(t *testing.T) {
    g := NewGraph()
    g.EnsureNode("A", nil, nil)
    if len(g.Nodes) != 1 {
        t.Fatalf("Nodes = %d, want 1", len(g.Nodes))
    }
    n := g.Nodes["A"]
    if n.ID != "A" {
        t.Errorf("ID = %q, want %q", n.ID, "A")
    }
    if n.Label != "A" {
        t.Errorf("Label = %q, want %q", n.Label, "A")
    }
    if n.Shape != Rectangle {
        t.Errorf("Shape = %v, want Rectangle", n.Shape)
    }

    // Update with label and shape
    label := "Start"
    shape := Stadium
    g.EnsureNode("A", &label, &shape)
    n = g.Nodes["A"]
    if n.Label != "Start" {
        t.Errorf("Label = %q, want %q", n.Label, "Start")
    }
    if n.Shape != Stadium {
        t.Errorf("Shape = %v, want Stadium", n.Shape)
    }
    if len(g.Nodes) != 1 {
        t.Errorf("Nodes = %d, want 1 (should not duplicate)", len(g.Nodes))
    }
}

func TestEnsureNodeOrder(t *testing.T) {
    g := NewGraph()
    g.EnsureNode("C", nil, nil)
    g.EnsureNode("A", nil, nil)
    g.EnsureNode("B", nil, nil)
    if g.NodeOrder["C"] != 0 {
        t.Errorf("C order = %d, want 0", g.NodeOrder["C"])
    }
    if g.NodeOrder["A"] != 1 {
        t.Errorf("A order = %d, want 1", g.NodeOrder["A"])
    }
    if g.NodeOrder["B"] != 2 {
        t.Errorf("B order = %d, want 2", g.NodeOrder["B"])
    }
    // Re-ensure does not change order
    g.EnsureNode("C", nil, nil)
    if g.NodeOrder["C"] != 0 {
        t.Errorf("C order = %d after re-ensure, want 0", g.NodeOrder["C"])
    }
}
```

**Step 2: Run test to verify it fails**

Run: `cd /Volumes/Development/mermaid-go && go test ./ir/ -v`
Expected: FAIL — packages/types don't exist yet

**Step 3: Write minimal implementation**

`ir/diagram.go`:
```go
package ir

type DiagramKind int

const (
    Flowchart DiagramKind = iota
    Class
    State
    Sequence
    Er
    Pie
    Mindmap
    Journey
    Timeline
    Gantt
    Requirement
    GitGraph
    C4
    Sankey
    Quadrant
    ZenUML
    Block
    Packet
    Kanban
    Architecture
    Radar
    Treemap
    XYChart
)

type Direction int

const (
    TopDown Direction = iota
    LeftRight
    BottomTop
    RightLeft
)

func DirectionFromToken(token string) (Direction, bool) {
    switch token {
    case "TD", "TB":
        return TopDown, true
    case "LR":
        return LeftRight, true
    case "RL":
        return RightLeft, true
    case "BT":
        return BottomTop, true
    default:
        return TopDown, false
    }
}
```

`ir/shapes.go`:
```go
package ir

type NodeShape int

const (
    Rectangle NodeShape = iota
    ForkJoin
    RoundRect
    Stadium
    Subroutine
    Cylinder
    ActorBox
    Circle
    DoubleCircle
    Diamond
    Hexagon
    Parallelogram
    ParallelogramAlt
    Trapezoid
    TrapezoidAlt
    Asymmetric
    MindmapDefault
    Text
)

type EdgeStyle int

const (
    Solid EdgeStyle = iota
    Dotted
    Thick
)

type EdgeDecoration int

const (
    DecCircle EdgeDecoration = iota
    DecCross
    DecDiamond
    DecDiamondFilled
    DecCrowsFootOne
    DecCrowsFootZeroOne
    DecCrowsFootMany
    DecCrowsFootZeroMany
)

type EdgeArrowhead int

const (
    OpenTriangle EdgeArrowhead = iota
    ClassDependency
)
```

`ir/style.go`:
```go
package ir

type NodeStyle struct {
    Fill            *string
    Stroke          *string
    TextColor       *string
    StrokeWidth     *float32
    StrokeDasharray *string
    LineColor       *string
}

type EdgeStyleOverride struct {
    Stroke      *string
    StrokeWidth *float32
    Dasharray   *string
    LabelColor  *string
}

type NodeLink struct {
    URL    string
    Title  *string
    Target *string
}
```

`ir/graph.go`:
```go
package ir

type Node struct {
    ID    string
    Label string
    Shape NodeShape
    Value *float32
    Icon  *string
}

type Edge struct {
    From            string
    To              string
    Label           *string
    StartLabel      *string
    EndLabel        *string
    Directed        bool
    ArrowStart      bool
    ArrowEnd        bool
    ArrowStartKind  *EdgeArrowhead
    ArrowEndKind    *EdgeArrowhead
    StartDecoration *EdgeDecoration
    EndDecoration   *EdgeDecoration
    Style           EdgeStyle
}

type Subgraph struct {
    ID        *string
    Label     string
    Nodes     []string
    Direction *Direction
    Icon      *string
}

type Graph struct {
    Kind      DiagramKind
    Direction Direction
    Nodes     map[string]*Node
    NodeOrder map[string]int
    Edges     []*Edge
    Subgraphs []*Subgraph

    ClassDefs        map[string]*NodeStyle
    NodeClasses      map[string][]string
    NodeStyles       map[string]*NodeStyle
    SubgraphStyles   map[string]*NodeStyle
    SubgraphClasses  map[string][]string
    NodeLinks        map[string]*NodeLink
    EdgeStyles       map[int]*EdgeStyleOverride
    EdgeStyleDefault *EdgeStyleOverride
}

func NewGraph() *Graph {
    return &Graph{
        Kind:           Flowchart,
        Direction:      TopDown,
        Nodes:          make(map[string]*Node),
        NodeOrder:      make(map[string]int),
        ClassDefs:      make(map[string]*NodeStyle),
        NodeClasses:    make(map[string][]string),
        NodeStyles:     make(map[string]*NodeStyle),
        SubgraphStyles: make(map[string]*NodeStyle),
        SubgraphClasses: make(map[string][]string),
        NodeLinks:      make(map[string]*NodeLink),
        EdgeStyles:     make(map[int]*EdgeStyleOverride),
    }
}

func (g *Graph) EnsureNode(id string, label *string, shape *NodeShape) {
    n, exists := g.Nodes[id]
    if !exists {
        n = &Node{
            ID:    id,
            Label: id,
            Shape: Rectangle,
        }
        g.Nodes[id] = n
        g.NodeOrder[id] = len(g.NodeOrder)
    }
    if label != nil {
        n.Label = *label
    }
    if shape != nil {
        n.Shape = *shape
    }
}
```

**Step 4: Run test to verify it passes**

Run: `cd /Volumes/Development/mermaid-go && go test ./ir/ -v`
Expected: PASS

**Step 5: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add ir/ && git commit -m "feat(ir): add intermediate representation types"
```

---

### Task 2: Theme package — colors and presets

**Files:**
- Create: `theme/theme.go`
- Create: `theme/color.go`
- Test: `theme/color_test.go`
- Test: `theme/theme_test.go`

**Step 1: Write the failing test**

```go
// theme/color_test.go
package theme

import (
    "math"
    "testing"
)

func TestParseHex3(t *testing.T) {
    h, s, l, ok := ParseColorToHSL("#fff")
    if !ok {
        t.Fatal("expected ok")
    }
    if math.Abs(float64(l)-100.0) > 0.1 {
        t.Errorf("l = %f, want ~100", l)
    }
    _ = h
    _ = s
}

func TestParseHex6(t *testing.T) {
    h, s, l, ok := ParseColorToHSL("#ECECFF")
    if !ok {
        t.Fatal("expected ok")
    }
    if h < 200 || h > 280 {
        t.Errorf("h = %f, expected ~240", h)
    }
    _ = s
    _ = l
}

func TestParseHSL(t *testing.T) {
    h, s, l, ok := ParseColorToHSL("hsl(240, 100%, 46.27%)")
    if !ok {
        t.Fatal("expected ok")
    }
    if math.Abs(float64(h)-240) > 0.1 {
        t.Errorf("h = %f, want 240", h)
    }
    if math.Abs(float64(s)-100) > 0.1 {
        t.Errorf("s = %f, want 100", s)
    }
    if math.Abs(float64(l)-46.27) > 0.1 {
        t.Errorf("l = %f, want 46.27", l)
    }
}

func TestAdjustColor(t *testing.T) {
    result := AdjustColor("#ECECFF", 0, 0, -10)
    // Should return an hsl() string
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
```

```go
// theme/theme_test.go
package theme

import "testing"

func TestModern(t *testing.T) {
    th := Modern()
    if th.FontSize != 14 {
        t.Errorf("FontSize = %f, want 14", th.FontSize)
    }
    if th.PrimaryColor == "" {
        t.Error("PrimaryColor is empty")
    }
    if th.FontFamily == "" {
        t.Error("FontFamily is empty")
    }
}

func TestMermaidDefault(t *testing.T) {
    th := MermaidDefault()
    if th.FontSize != 16 {
        t.Errorf("FontSize = %f, want 16", th.FontSize)
    }
    if th.PrimaryColor != "#ECECFF" {
        t.Errorf("PrimaryColor = %q, want #ECECFF", th.PrimaryColor)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `cd /Volumes/Development/mermaid-go && go test ./theme/ -v`
Expected: FAIL

**Step 3: Write minimal implementation**

Port `theme/color.go` from Rust `src/theme.rs` lines 198-287 (parseHex, parseHSL, rgbToHSL, AdjustColor, ParseColorToHSL).

Port `theme/theme.go` from Rust `src/theme.rs` lines 36-178 (Theme struct with all fields, Modern() and MermaidDefault() constructors).

**Step 4: Run test to verify it passes**

Run: `cd /Volumes/Development/mermaid-go && go test ./theme/ -v`
Expected: PASS

**Step 5: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add theme/ && git commit -m "feat(theme): add theme types and color utilities"
```

---

### Task 3: Config package — layout configuration

**Files:**
- Create: `config/config.go`
- Test: `config/config_test.go`

**Step 1: Write the failing test**

```go
// config/config_test.go
package config

import "testing"

func TestDefaultLayout(t *testing.T) {
    cfg := DefaultLayout()
    if cfg.NodeSpacing != 50 {
        t.Errorf("NodeSpacing = %f, want 50", cfg.NodeSpacing)
    }
    if cfg.RankSpacing != 70 {
        t.Errorf("RankSpacing = %f, want 70", cfg.RankSpacing)
    }
    if cfg.LabelLineHeight <= 0 {
        t.Error("LabelLineHeight should be > 0")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `cd /Volumes/Development/mermaid-go && go test ./config/ -v`

**Step 3: Write minimal implementation**

`config/config.go` — Port the Layout/LayoutConfig struct from Rust `src/config.rs`. Only include fields needed for flowchart rendering in Phase 1: `NodeSpacing`, `RankSpacing`, `LabelLineHeight`, `PreferredAspectRatio`. Stub diagram-specific configs.

**Step 4: Run test to verify it passes**

Run: `cd /Volumes/Development/mermaid-go && go test ./config/ -v`

**Step 5: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add config/ && git commit -m "feat(config): add layout configuration with defaults"
```

---

### Task 4: Text metrics — font measurement

**Files:**
- Create: `textmetrics/measurer.go`
- Test: `textmetrics/measurer_test.go`

**Step 1: Write the failing test**

```go
// textmetrics/measurer_test.go
package textmetrics

import "testing"

func TestMeasureEmpty(t *testing.T) {
    m := New()
    w := m.Width("", 14, "sans-serif")
    if w != 0 {
        t.Errorf("Width of empty = %f, want 0", w)
    }
}

func TestMeasureNonEmpty(t *testing.T) {
    m := New()
    w := m.Width("Hello", 14, "sans-serif")
    if w <= 0 {
        t.Errorf("Width of 'Hello' = %f, want > 0", w)
    }
}

func TestMeasureLongerIsWider(t *testing.T) {
    m := New()
    short := m.Width("Hi", 14, "sans-serif")
    long := m.Width("Hello World", 14, "sans-serif")
    if long <= short {
        t.Errorf("long (%f) should be > short (%f)", long, short)
    }
}

func TestMeasureLargerFontIsWider(t *testing.T) {
    m := New()
    small := m.Width("Hello", 10, "sans-serif")
    big := m.Width("Hello", 20, "sans-serif")
    if big <= small {
        t.Errorf("big (%f) should be > small (%f)", big, small)
    }
}

func TestAverageCharWidth(t *testing.T) {
    m := New()
    w := m.AverageCharWidth("sans-serif", 14)
    if w <= 0 {
        t.Errorf("AverageCharWidth = %f, want > 0", w)
    }
    if w > 14 {
        t.Errorf("AverageCharWidth = %f, unexpectedly large", w)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `cd /Volumes/Development/mermaid-go && go test ./textmetrics/ -v`

**Step 3: Write minimal implementation**

`textmetrics/measurer.go`:
- `Measurer` struct with `sync.Mutex`, font data cache, and loaded flag
- `New()` constructor
- `Width(text, fontSize, fontFamily)` — loads system font via `golang.org/x/image/font/sfnt`, caches parsed font, iterates glyphs for advance widths. Falls back to `fontSize * 0.6 * len(text)` if font not found.
- `AverageCharWidth(fontFamily, fontSize)` — measures sample alphabet, returns mean

Need to add dependency: `go get golang.org/x/image`

**Step 4: Run test to verify it passes**

Run: `cd /Volumes/Development/mermaid-go && go test ./textmetrics/ -v`

**Step 5: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add textmetrics/ go.mod go.sum && git commit -m "feat(textmetrics): add pure Go font measurement"
```

---

### Task 5: Parser — diagram detection and flowchart parsing

**Files:**
- Create: `parser/parser.go`
- Create: `parser/flowchart.go`
- Create: `parser/helpers.go`
- Test: `parser/parser_test.go`
- Test: `parser/flowchart_test.go`

**Step 1: Write the failing tests**

```go
// parser/parser_test.go
package parser

import (
    "testing"

    "github.com/jamesainslie/mermaid-go/ir"
)

func TestDetectDiagramKind(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  ir.DiagramKind
    }{
        {"flowchart LR", "flowchart LR\n  A-->B", ir.Flowchart},
        {"graph TD", "graph TD\n  A-->B", ir.Flowchart},
        {"sequenceDiagram", "sequenceDiagram\n  Alice->>Bob: Hi", ir.Sequence},
        {"classDiagram", "classDiagram\n  A <|-- B", ir.Class},
        {"stateDiagram", "stateDiagram-v2\n  [*] --> A", ir.State},
        {"pie", "pie\n  \"A\" : 10", ir.Pie},
        {"skip comments", "%%{init}%%\nflowchart LR\n  A-->B", ir.Flowchart},
        {"skip empty lines", "\n\n  flowchart TD\n  A-->B", ir.Flowchart},
        {"default flowchart", "A-->B", ir.Flowchart},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := detectDiagramKind(tt.input)
            if got != tt.want {
                t.Errorf("detectDiagramKind() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

```go
// parser/flowchart_test.go
package parser

import (
    "testing"

    "github.com/jamesainslie/mermaid-go/ir"
)

func TestParseFlowchartSimpleChain(t *testing.T) {
    out, err := Parse("flowchart LR; A-->B-->C")
    if err != nil {
        t.Fatalf("Parse() error: %v", err)
    }
    g := out.Graph
    if g.Kind != ir.Flowchart {
        t.Errorf("Kind = %v, want Flowchart", g.Kind)
    }
    if g.Direction != ir.LeftRight {
        t.Errorf("Direction = %v, want LeftRight", g.Direction)
    }
    if len(g.Nodes) != 3 {
        t.Errorf("Nodes = %d, want 3", len(g.Nodes))
    }
    if len(g.Edges) != 2 {
        t.Errorf("Edges = %d, want 2", len(g.Edges))
    }
}

func TestParseFlowchartWithLabels(t *testing.T) {
    out, err := Parse("flowchart TD\n  A[Start] --> B{Decision}\n  B -->|Yes| C[OK]")
    if err != nil {
        t.Fatalf("Parse() error: %v", err)
    }
    g := out.Graph
    if g.Nodes["A"].Label != "Start" {
        t.Errorf("A label = %q, want Start", g.Nodes["A"].Label)
    }
    if g.Nodes["A"].Shape != ir.Rectangle {
        t.Errorf("A shape = %v, want Rectangle", g.Nodes["A"].Shape)
    }
    if g.Nodes["B"].Shape != ir.Diamond {
        t.Errorf("B shape = %v, want Diamond", g.Nodes["B"].Shape)
    }
    if len(g.Edges) != 2 {
        t.Errorf("Edges = %d, want 2", len(g.Edges))
    }
    // Second edge should have label "Yes"
    var labelEdge *ir.Edge
    for _, e := range g.Edges {
        if e.Label != nil && *e.Label == "Yes" {
            labelEdge = e
        }
    }
    if labelEdge == nil {
        t.Error("expected edge with label 'Yes'")
    }
}

func TestParseFlowchartDottedEdge(t *testing.T) {
    out, err := Parse("flowchart LR; A-.->B")
    if err != nil {
        t.Fatalf("Parse() error: %v", err)
    }
    if len(out.Graph.Edges) != 1 {
        t.Fatalf("Edges = %d, want 1", len(out.Graph.Edges))
    }
    if out.Graph.Edges[0].Style != ir.Dotted {
        t.Errorf("Style = %v, want Dotted", out.Graph.Edges[0].Style)
    }
}

func TestParseFlowchartThickEdge(t *testing.T) {
    out, err := Parse("flowchart LR; A==>B")
    if err != nil {
        t.Fatalf("Parse() error: %v", err)
    }
    if out.Graph.Edges[0].Style != ir.Thick {
        t.Errorf("Style = %v, want Thick", out.Graph.Edges[0].Style)
    }
}

func TestParseFlowchartBidirectional(t *testing.T) {
    out, err := Parse("flowchart LR; A<-->B")
    if err != nil {
        t.Fatalf("Parse() error: %v", err)
    }
    e := out.Graph.Edges[0]
    if !e.ArrowStart || !e.ArrowEnd {
        t.Errorf("ArrowStart=%v ArrowEnd=%v, want both true", e.ArrowStart, e.ArrowEnd)
    }
}

func TestParseFlowchartSubgraph(t *testing.T) {
    input := "flowchart TD\n  subgraph sg1[Group]\n    A-->B\n  end\n  C-->A"
    out, err := Parse(input)
    if err != nil {
        t.Fatalf("Parse() error: %v", err)
    }
    if len(out.Graph.Subgraphs) != 1 {
        t.Fatalf("Subgraphs = %d, want 1", len(out.Graph.Subgraphs))
    }
    sg := out.Graph.Subgraphs[0]
    if sg.Label != "Group" {
        t.Errorf("Subgraph label = %q, want Group", sg.Label)
    }
}

func TestParseFlowchartShapes(t *testing.T) {
    tests := []struct {
        input string
        shape ir.NodeShape
    }{
        {"flowchart LR; A[rect]", ir.Rectangle},
        {"flowchart LR; A(round)", ir.RoundRect},
        {"flowchart LR; A([stadium])", ir.Stadium},
        {"flowchart LR; A{diamond}", ir.Diamond},
        {"flowchart LR; A{{hexagon}}", ir.Hexagon},
        {"flowchart LR; A[[subroutine]]", ir.Subroutine},
        {"flowchart LR; A[(cylinder)]", ir.Cylinder},
        {"flowchart LR; A((circle))", ir.DoubleCircle},
        {"flowchart LR; A>asym]", ir.Asymmetric},
    }
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            out, err := Parse(tt.input)
            if err != nil {
                t.Fatalf("Parse() error: %v", err)
            }
            if out.Graph.Nodes["A"].Shape != tt.shape {
                t.Errorf("shape = %v, want %v", out.Graph.Nodes["A"].Shape, tt.shape)
            }
        })
    }
}
```

**Step 2: Run tests to verify they fail**

Run: `cd /Volumes/Development/mermaid-go && go test ./parser/ -v`

**Step 3: Write implementation**

Port from Rust `src/parser.rs`:

- `parser/parser.go`: `Parse()` dispatcher, `detectDiagramKind()`, `ParseOutput` struct
- `parser/helpers.go`: `stripTrailingComment()`, `splitStatements()`, `maskBracketContent()`, `splitEdgeChain()`, `parseEdgeLine()`, `parseEdgeMeta()`, `parseNodeToken()`, `splitIDLabel()` — helper functions
- `parser/flowchart.go`: `parseFlowchart()`, `addFlowchartEdge()` — flowchart-specific parsing

Regex patterns compiled at package level with `regexp.MustCompile`.

**Step 4: Run tests to verify they pass**

Run: `cd /Volumes/Development/mermaid-go && go test ./parser/ -v`

**Step 5: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add parser/ && git commit -m "feat(parser): add flowchart parser with edge chain support"
```

---

### Task 6: Layout types and core layout orchestrator

**Files:**
- Create: `layout/types.go`
- Create: `layout/layout.go`
- Create: `layout/sizing.go`
- Create: `layout/ranking.go`
- Create: `layout/ordering.go`
- Create: `layout/positioning.go`
- Create: `layout/routing.go`
- Test: `layout/ranking_test.go`
- Test: `layout/layout_test.go`

**Step 1: Write the failing tests**

```go
// layout/ranking_test.go
package layout

import (
    "testing"

    "github.com/jamesainslie/mermaid-go/ir"
)

func edge(from, to string) *ir.Edge {
    return &ir.Edge{
        From: from, To: to, Directed: true, ArrowEnd: true, Style: ir.Solid,
    }
}

func TestComputeRanksLinearChain(t *testing.T) {
    nodes := []string{"A", "B", "C"}
    edges := []*ir.Edge{edge("A", "B"), edge("B", "C")}
    ranks := computeRanks(nodes, edges, map[string]int{})
    if ranks["A"] != 0 {
        t.Errorf("A = %d, want 0", ranks["A"])
    }
    if ranks["B"] != 1 {
        t.Errorf("B = %d, want 1", ranks["B"])
    }
    if ranks["C"] != 2 {
        t.Errorf("C = %d, want 2", ranks["C"])
    }
}

func TestComputeRanksDiamond(t *testing.T) {
    nodes := []string{"A", "B", "C", "D"}
    edges := []*ir.Edge{edge("A", "B"), edge("A", "C"), edge("B", "D"), edge("C", "D")}
    ranks := computeRanks(nodes, edges, map[string]int{})
    if ranks["A"] != 0 {
        t.Errorf("A = %d, want 0", ranks["A"])
    }
    if ranks["B"] != 1 {
        t.Errorf("B = %d, want 1", ranks["B"])
    }
    if ranks["C"] != 1 {
        t.Errorf("C = %d, want 1", ranks["C"])
    }
    if ranks["D"] != 2 {
        t.Errorf("D = %d, want 2", ranks["D"])
    }
}

func TestComputeRanksCycle(t *testing.T) {
    nodes := []string{"A", "B", "C"}
    edges := []*ir.Edge{edge("A", "B"), edge("B", "C"), edge("C", "A")}
    ranks := computeRanks(nodes, edges, map[string]int{})
    if len(ranks) != 3 {
        t.Errorf("len(ranks) = %d, want 3", len(ranks))
    }
}
```

```go
// layout/layout_test.go
package layout

import (
    "testing"

    "github.com/jamesainslie/mermaid-go/config"
    "github.com/jamesainslie/mermaid-go/ir"
    "github.com/jamesainslie/mermaid-go/theme"
)

func TestComputeLayoutSimple(t *testing.T) {
    g := ir.NewGraph()
    g.Kind = ir.Flowchart
    g.Direction = ir.LeftRight
    g.EnsureNode("A", nil, nil)
    g.EnsureNode("B", nil, nil)
    g.EnsureNode("C", nil, nil)
    g.Edges = []*ir.Edge{edge("A", "B"), edge("B", "C")}

    th := theme.Modern()
    cfg := config.DefaultLayout()
    l := ComputeLayout(g, th, cfg)

    if l.Kind != ir.Flowchart {
        t.Errorf("Kind = %v, want Flowchart", l.Kind)
    }
    if len(l.Nodes) != 3 {
        t.Errorf("Nodes = %d, want 3", len(l.Nodes))
    }
    if len(l.Edges) != 2 {
        t.Errorf("Edges = %d, want 2", len(l.Edges))
    }
    if l.Width <= 0 {
        t.Errorf("Width = %f, want > 0", l.Width)
    }
    if l.Height <= 0 {
        t.Errorf("Height = %f, want > 0", l.Height)
    }

    // In LR direction, nodes should be positioned left to right
    ax := l.Nodes["A"].X
    bx := l.Nodes["B"].X
    cx := l.Nodes["C"].X
    if ax >= bx || bx >= cx {
        t.Errorf("expected A.x < B.x < C.x, got %f, %f, %f", ax, bx, cx)
    }
}
```

**Step 2: Run tests to verify they fail**

Run: `cd /Volumes/Development/mermaid-go && go test ./layout/ -v`

**Step 3: Write implementation**

Port from Rust `src/layout/`:

- `layout/types.go`: `Layout`, `NodeLayout`, `EdgeLayout`, `SubgraphLayout`, `TextBlock`, `DiagramData` interface + `GraphData` struct
- `layout/layout.go`: `ComputeLayout()` dispatcher, `computeGraphLayout()` pipeline
- `layout/sizing.go`: `sizeNodes()` — compute node width/height from text metrics
- `layout/ranking.go`: `computeRanks()` — topological sort with cycle handling (from Rust `ranking.rs`)
- `layout/ordering.go`: `orderRankNodes()` — median-based crossing minimization (from Rust `ranking.rs`)
- `layout/positioning.go`: `positionNodes()` — assign x,y coordinates based on rank and order
- `layout/routing.go`: Simple L-shaped edge routing (NOT full A* — that's Phase 5). Connect edge endpoints with 2-3 point polylines.

**Step 4: Run tests to verify they pass**

Run: `cd /Volumes/Development/mermaid-go && go test ./layout/ -v`

**Step 5: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add layout/ && git commit -m "feat(layout): add graph layout pipeline with ranking and positioning"
```

---

### Task 7: SVG renderer — flowchart output

**Files:**
- Create: `render/svg.go`
- Create: `render/builder.go`
- Create: `render/shapes.go`
- Create: `render/graph.go`
- Test: `render/svg_test.go`

**Step 1: Write the failing tests**

```go
// render/svg_test.go
package render

import (
    "strings"
    "testing"

    "github.com/jamesainslie/mermaid-go/config"
    "github.com/jamesainslie/mermaid-go/ir"
    "github.com/jamesainslie/mermaid-go/layout"
    "github.com/jamesainslie/mermaid-go/theme"
)

func simpleLayout() *layout.Layout {
    g := ir.NewGraph()
    g.Kind = ir.Flowchart
    g.Direction = ir.LeftRight
    g.EnsureNode("A", nil, nil)
    g.EnsureNode("B", nil, nil)
    g.Edges = []*ir.Edge{{
        From: "A", To: "B", Directed: true, ArrowEnd: true, Style: ir.Solid,
    }}
    th := theme.Modern()
    cfg := config.DefaultLayout()
    return layout.ComputeLayout(g, th, cfg)
}

func TestRenderSVGContainsSVGTags(t *testing.T) {
    l := simpleLayout()
    th := theme.Modern()
    cfg := config.DefaultLayout()
    svg := RenderSVG(l, th, cfg)
    if !strings.Contains(svg, "<svg") {
        t.Error("missing <svg tag")
    }
    if !strings.Contains(svg, "</svg>") {
        t.Error("missing </svg> tag")
    }
}

func TestRenderSVGContainsNodes(t *testing.T) {
    l := simpleLayout()
    th := theme.Modern()
    cfg := config.DefaultLayout()
    svg := RenderSVG(l, th, cfg)
    if !strings.Contains(svg, "<rect") {
        t.Error("missing <rect for node shapes")
    }
    if !strings.Contains(svg, "A") {
        t.Error("missing node label A")
    }
    if !strings.Contains(svg, "B") {
        t.Error("missing node label B")
    }
}

func TestRenderSVGContainsEdge(t *testing.T) {
    l := simpleLayout()
    th := theme.Modern()
    cfg := config.DefaultLayout()
    svg := RenderSVG(l, th, cfg)
    if !strings.Contains(svg, "<path") || !strings.Contains(svg, "edgePath") {
        t.Error("missing edge path")
    }
}

func TestRenderSVGHasViewBox(t *testing.T) {
    l := simpleLayout()
    th := theme.Modern()
    cfg := config.DefaultLayout()
    svg := RenderSVG(l, th, cfg)
    if !strings.Contains(svg, "viewBox") {
        t.Error("missing viewBox attribute")
    }
}
```

**Step 2: Run tests to verify they fail**

Run: `cd /Volumes/Development/mermaid-go && go test ./render/ -v`

**Step 3: Write implementation**

Port from Rust `src/render.rs`:

- `render/builder.go`: `svgBuilder` struct wrapping `strings.Builder` with helper methods (`openTag`, `closeTag`, `selfClose`, `writeAttr`, `rect`, `circle`, `path`, `text`, `polyline`, `escapeXML`)
- `render/shapes.go`: `renderNodeShape()` — dispatches to shape-specific SVG based on `NodeShape`
- `render/graph.go`: `renderGraph()` — renders subgraphs, nodes, edges, labels for flowchart/graph diagrams
- `render/svg.go`: `RenderSVG()` — opens SVG document, sets viewBox/dimensions, renders background, dispatches to diagram-specific renderer, closes SVG

**Step 4: Run tests to verify they pass**

Run: `cd /Volumes/Development/mermaid-go && go test ./render/ -v`

**Step 5: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add render/ && git commit -m "feat(render): add SVG renderer for flowchart diagrams"
```

---

### Task 8: Public API — wire up the full pipeline

**Files:**
- Create: `mermaid.go`
- Create: `options.go`
- Test: `mermaid_test.go`

**Step 1: Write the failing tests**

```go
// mermaid_test.go
package mermaid

import (
    "strings"
    "testing"
)

func TestRender(t *testing.T) {
    svg, err := Render("flowchart LR; A-->B-->C")
    if err != nil {
        t.Fatalf("Render() error: %v", err)
    }
    if !strings.Contains(svg, "<svg") {
        t.Error("missing <svg")
    }
    if !strings.Contains(svg, "</svg>") {
        t.Error("missing </svg>")
    }
}

func TestRenderWithOptions(t *testing.T) {
    opts := Options{}
    svg, err := RenderWithOptions("flowchart TD; X-->Y", opts)
    if err != nil {
        t.Fatalf("RenderWithOptions() error: %v", err)
    }
    if !strings.Contains(svg, "<svg") {
        t.Error("missing <svg")
    }
}

func TestRenderWithTiming(t *testing.T) {
    result, err := RenderWithTiming("flowchart LR; A-->B", Options{})
    if err != nil {
        t.Fatalf("RenderWithTiming() error: %v", err)
    }
    if !strings.Contains(result.SVG, "<svg") {
        t.Error("missing <svg")
    }
    if result.TotalUs() <= 0 {
        t.Error("TotalUs should be > 0")
    }
}

func TestRenderInvalidInput(t *testing.T) {
    _, err := Render("")
    if err == nil {
        t.Error("expected error for empty input")
    }
}

func TestRenderContainsNodeLabels(t *testing.T) {
    svg, err := Render("flowchart LR\n  A[Start] --> B[End]")
    if err != nil {
        t.Fatalf("Render() error: %v", err)
    }
    if !strings.Contains(svg, "Start") {
        t.Error("missing label 'Start'")
    }
    if !strings.Contains(svg, "End") {
        t.Error("missing label 'End'")
    }
}
```

**Step 2: Run tests to verify they fail**

Run: `cd /Volumes/Development/mermaid-go && go test -v`

**Step 3: Write implementation**

`mermaid.go`:
```go
package mermaid

import (
    "fmt"
    "time"

    "github.com/jamesainslie/mermaid-go/config"
    "github.com/jamesainslie/mermaid-go/layout"
    "github.com/jamesainslie/mermaid-go/parser"
    "github.com/jamesainslie/mermaid-go/render"
    "github.com/jamesainslie/mermaid-go/theme"
)

func Render(input string) (string, error) {
    return RenderWithOptions(input, Options{})
}

func RenderWithOptions(input string, opts Options) (string, error) {
    th := opts.theme()
    cfg := opts.layout()
    parsed, err := parser.Parse(input)
    if err != nil {
        return "", fmt.Errorf("parse: %w", err)
    }
    l := layout.ComputeLayout(parsed.Graph, th, cfg)
    svg := render.RenderSVG(l, th, cfg)
    return svg, nil
}

func RenderWithTiming(input string, opts Options) (*Result, error) {
    th := opts.theme()
    cfg := opts.layout()

    t0 := time.Now()
    parsed, err := parser.Parse(input)
    if err != nil {
        return nil, fmt.Errorf("parse: %w", err)
    }
    parseUs := time.Since(t0).Microseconds()

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
```

`options.go`:
```go
package mermaid

import (
    "github.com/jamesainslie/mermaid-go/config"
    "github.com/jamesainslie/mermaid-go/theme"
)

type Options struct {
    Theme  *theme.Theme
    Layout *config.Layout
}

func (o Options) theme() *theme.Theme {
    if o.Theme != nil {
        return o.Theme
    }
    return theme.Modern()
}

func (o Options) layout() *config.Layout {
    if o.Layout != nil {
        return o.Layout
    }
    return config.DefaultLayout()
}

type Result struct {
    SVG      string
    ParseUs  int64
    LayoutUs int64
    RenderUs int64
}

func (r *Result) TotalUs() int64 {
    return r.ParseUs + r.LayoutUs + r.RenderUs
}

func (r *Result) TotalMs() float64 {
    return float64(r.TotalUs()) / 1000.0
}
```

**Step 4: Run tests to verify they pass**

Run: `cd /Volumes/Development/mermaid-go && go test -v`

**Step 5: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add mermaid.go options.go mermaid_test.go && git commit -m "feat: wire up full Render pipeline"
```

---

### Task 9: Golden test — cross-validate against Rust reference

**Files:**
- Create: `testdata/fixtures/flowchart-simple.mmd`
- Modify: `mermaid_test.go` (add golden test)

**Step 1: Create fixture files**

```
testdata/fixtures/flowchart-simple.mmd:
flowchart LR; A-->B-->C

testdata/fixtures/flowchart-labels.mmd:
flowchart TD
    A[Start] --> B{Decision}
    B -->|Yes| C[OK]
    B -->|No| D[Cancel]

testdata/fixtures/flowchart-shapes.mmd:
flowchart LR
    A[Rectangle] --> B(Rounded)
    B --> C([Stadium])
    C --> D{Diamond}
    D --> E((Circle))
```

**Step 2: Add golden test**

```go
// In mermaid_test.go
func TestGoldenFlowchartSimple(t *testing.T) {
    input := "flowchart LR; A-->B-->C"
    svg, err := Render(input)
    if err != nil {
        t.Fatalf("Render() error: %v", err)
    }
    // Structural checks — the SVG should contain the expected elements
    if !strings.Contains(svg, "viewBox") {
        t.Error("missing viewBox")
    }
    if strings.Count(svg, "<rect") < 3 {
        t.Errorf("expected at least 3 rects (nodes), got %d", strings.Count(svg, "<rect"))
    }
    if strings.Count(svg, "edgePath") < 2 {
        t.Errorf("expected at least 2 edge paths, got %d", strings.Count(svg, "edgePath"))
    }
    // Node labels present
    for _, label := range []string{"A", "B", "C"} {
        if !strings.Contains(svg, ">"+label+"<") {
            t.Errorf("missing node label %q in SVG", label)
        }
    }
}
```

**Step 3: Run test**

Run: `cd /Volumes/Development/mermaid-go && go test -v -run TestGolden`

**Step 4: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add testdata/ mermaid_test.go && git commit -m "test: add golden tests for flowchart SVG output"
```

---

### Task 10: Benchmark — establish baseline performance

**Files:**
- Create: `mermaid_bench_test.go`

**Step 1: Write benchmarks**

```go
// mermaid_bench_test.go
package mermaid

import "testing"

func BenchmarkRenderSimple(b *testing.B) {
    input := "flowchart LR; A-->B-->C"
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _, err := Render(input)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkRenderMedium(b *testing.B) {
    input := `flowchart TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Process]
    B -->|No| D[Cancel]
    C --> E[End]
    D --> E
    E --> F[Cleanup]
    F --> G[Done]`
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        _, err := Render(input)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

**Step 2: Run benchmarks**

Run: `cd /Volumes/Development/mermaid-go && go test -bench=. -benchmem -count=3`

Record baseline numbers.

**Step 3: Commit**

```bash
cd /Volumes/Development/mermaid-go && git add mermaid_bench_test.go && git commit -m "bench: add baseline render benchmarks"
```

---

## Summary

| Task | Package | What |
|------|---------|------|
| 1 | ir | Enums, Node, Edge, Graph, EnsureNode |
| 2 | theme | Theme struct, Modern/MermaidDefault, color utils |
| 3 | config | LayoutConfig with defaults |
| 4 | textmetrics | Pure Go font measurement |
| 5 | parser | Diagram detection + flowchart parser |
| 6 | layout | Ranking, ordering, positioning, L-shaped routing |
| 7 | render | SVG builder + flowchart renderer |
| 8 | root | Public API: Render(), RenderWithOptions(), RenderWithTiming() |
| 9 | testdata | Golden tests against Rust reference |
| 10 | root | Baseline benchmarks |

**Exit criteria:** `mermaid.Render("flowchart LR; A-->B-->C")` produces valid SVG with correct structure. All tests pass. Benchmark baseline recorded.
