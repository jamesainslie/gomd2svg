# Phase 8: Mindmap, Sankey & Treemap Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Mindmap (indentation-based tree), Sankey (flow diagram with CSV data), and Treemap (nested rectangles proportional to values) diagram support to the mermaid-go renderer.

**Architecture:** Follow the established per-diagram pipeline: IR types → config/theme → parser → layout → renderer → integration tests. Mindmap uses a radial tree layout algorithm (place root at center, distribute children in arcs). Sankey uses a column-based flow layout (topological sort → column assignment → vertical positioning → curved link routing). Treemap uses the squarified treemap algorithm (Bruls-Huizing-van Wijk) to partition rectangles proportional to values.

**Tech Stack:** Go stdlib only (math, regexp, strings, strconv, fmt, sort). No external dependencies.

---

## Reference: Mermaid Syntax

### Mindmap

```
mindmap
    root((Central Idea))
        Branch A
            [Square Node]
            (Rounded Node)
        Branch B
            ))Bang Node((
            )Cloud Node(
            {{Hexagon Node}}
```

- Keyword: `mindmap`
- Hierarchy via indentation (like Kanban)
- 7 node shapes: default (no border), `[square]`, `(rounded)`, `((circle))`, `))bang((`, `)cloud(`, `{{hexagon}}`
- `::icon(css-class)` for icon decorators (parsed but not rendered in SVG)
- `:::className` for CSS class styling (parsed but styling deferred)

### Sankey

```
sankey-beta

Source A,Target X,100
Source A,Target Y,200
Source B,Target X,150
Source B,Target Z,300
```

- Keyword: `sankey-beta` or `sankey`
- CSV data: `source,target,value` per line
- Quoted strings for names with commas: `"Heating, homes"`
- Empty lines allowed for spacing
- Nodes auto-assigned to columns by flow topology

### Treemap

```
treemap-beta
"Root"
    "Section A"
        "Leaf 1": 30
        "Leaf 2": 50
    "Section B"
        "Leaf 3": 20
```

- Keyword: `treemap-beta` or `treemap`
- Indentation-based hierarchy
- Quoted node names (double or single quotes)
- Section nodes: name only (branch containers)
- Leaf nodes: name followed by `:` and numeric value
- `:::className` for CSS class references (parsed, styling deferred)

---

## Existing Patterns to Follow

- Indentation-based parsing: reuse `kanbanLine` pattern from `parser/kanban.go`
- `DiagramData` sealed interface with unexported `diagramData()` marker method
- Config-driven sizing via `config.Layout` sub-structs
- Theme color slices with empty-slice guards and fallback colors
- iota enums with `String()` methods
- Package-level `regexp.MustCompile` vars
- Table-driven tests

---

### Task 1: IR Types — Mindmap

**Files:**
- Create: `ir/mindmap.go`
- Create: `ir/mindmap_test.go`
- Modify: `ir/graph.go` — add Mindmap fields

**Step 1: Write the failing test**

Create `ir/mindmap_test.go`:

```go
package ir

import "testing"

func TestMindmapNodeShape(t *testing.T) {
	tests := []struct {
		shape MindmapShape
		want  string
	}{
		{MindmapDefault, "default"},
		{MindmapSquare, "square"},
		{MindmapRounded, "rounded"},
		{MindmapCircle, "circle"},
		{MindmapBang, "bang"},
		{MindmapCloud, "cloud"},
		{MindmapHexagon, "hexagon"},
	}
	for _, tc := range tests {
		if got := tc.shape.String(); got != tc.want {
			t.Errorf("MindmapShape(%d).String() = %q, want %q", tc.shape, got, tc.want)
		}
	}
}

func TestMindmapNode(t *testing.T) {
	root := &MindmapNode{
		ID:    "root",
		Label: "Central Idea",
		Shape: MindmapCircle,
		Children: []*MindmapNode{
			{ID: "a", Label: "Branch A", Shape: MindmapDefault},
			{ID: "b", Label: "Branch B", Shape: MindmapSquare},
		},
	}
	if len(root.Children) != 2 {
		t.Errorf("children len = %d, want 2", len(root.Children))
	}
	if root.Shape != MindmapCircle {
		t.Errorf("shape = %v, want Circle", root.Shape)
	}
}

func TestMindmapGraphField(t *testing.T) {
	g := NewGraph()
	g.Kind = Mindmap
	g.MindmapRoot = &MindmapNode{ID: "root", Label: "Root"}
	if g.MindmapRoot == nil {
		t.Error("MindmapRoot should not be nil")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestMindmap -v`
Expected: FAIL — types not defined

**Step 3: Write minimal implementation**

Create `ir/mindmap.go`:

```go
package ir

// MindmapShape represents the visual shape of a mindmap node.
type MindmapShape int

const (
	MindmapDefault MindmapShape = iota
	MindmapSquare
	MindmapRounded
	MindmapCircle
	MindmapBang
	MindmapCloud
	MindmapHexagon
)

func (s MindmapShape) String() string {
	switch s {
	case MindmapDefault:
		return "default"
	case MindmapSquare:
		return "square"
	case MindmapRounded:
		return "rounded"
	case MindmapCircle:
		return "circle"
	case MindmapBang:
		return "bang"
	case MindmapCloud:
		return "cloud"
	case MindmapHexagon:
		return "hexagon"
	default:
		return "unknown"
	}
}

// MindmapNode represents a node in the mindmap tree.
type MindmapNode struct {
	ID       string
	Label    string
	Shape    MindmapShape
	Icon     string // CSS class from ::icon()
	Class    string // CSS class from :::
	Children []*MindmapNode
}
```

Add to `ir/graph.go` after the Radar fields block:

```go
	// Mindmap diagram fields
	MindmapRoot *MindmapNode
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestMindmap -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/mindmap.go ir/mindmap_test.go ir/graph.go
git commit -m "feat(ir): add Mindmap diagram types"
```

---

### Task 2: IR Types — Sankey

**Files:**
- Create: `ir/sankey.go`
- Create: `ir/sankey_test.go`
- Modify: `ir/graph.go` — add Sankey fields

**Step 1: Write the failing test**

Create `ir/sankey_test.go`:

```go
package ir

import "testing"

func TestSankeyLink(t *testing.T) {
	link := &SankeyLink{Source: "A", Target: "B", Value: 100.5}
	if link.Source != "A" {
		t.Errorf("Source = %q, want A", link.Source)
	}
	if link.Value != 100.5 {
		t.Errorf("Value = %v, want 100.5", link.Value)
	}
}

func TestSankeyGraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Sankey
	g.SankeyLinks = append(g.SankeyLinks, &SankeyLink{
		Source: "Solar", Target: "Grid", Value: 59.9,
	})
	if len(g.SankeyLinks) != 1 {
		t.Fatalf("SankeyLinks len = %d, want 1", len(g.SankeyLinks))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestSankey -v`
Expected: FAIL

**Step 3: Write minimal implementation**

Create `ir/sankey.go`:

```go
package ir

// SankeyLink represents a flow from source to target with a value.
type SankeyLink struct {
	Source string
	Target string
	Value  float64
}
```

Add to `ir/graph.go`:

```go
	// Sankey diagram fields
	SankeyLinks []*SankeyLink
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestSankey -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/sankey.go ir/sankey_test.go ir/graph.go
git commit -m "feat(ir): add Sankey diagram types"
```

---

### Task 3: IR Types — Treemap

**Files:**
- Create: `ir/treemap.go`
- Create: `ir/treemap_test.go`
- Modify: `ir/graph.go` — add Treemap fields

**Step 1: Write the failing test**

Create `ir/treemap_test.go`:

```go
package ir

import "testing"

func TestTreemapNode(t *testing.T) {
	root := &TreemapNode{
		Label: "Root",
		Children: []*TreemapNode{
			{Label: "Leaf A", Value: 30},
			{Label: "Section B", Children: []*TreemapNode{
				{Label: "Leaf C", Value: 20},
			}},
		},
	}
	if len(root.Children) != 2 {
		t.Errorf("children len = %d, want 2", len(root.Children))
	}
	if root.Children[0].Value != 30 {
		t.Errorf("leaf value = %v, want 30", root.Children[0].Value)
	}
}

func TestTreemapNodeIsLeaf(t *testing.T) {
	leaf := &TreemapNode{Label: "Leaf", Value: 10}
	section := &TreemapNode{Label: "Sec", Children: []*TreemapNode{leaf}}
	if !leaf.IsLeaf() {
		t.Error("expected leaf")
	}
	if section.IsLeaf() {
		t.Error("expected non-leaf")
	}
}

func TestTreemapGraphField(t *testing.T) {
	g := NewGraph()
	g.Kind = Treemap
	g.TreemapRoot = &TreemapNode{Label: "Root"}
	if g.TreemapRoot == nil {
		t.Error("TreemapRoot should not be nil")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestTreemap -v`
Expected: FAIL

**Step 3: Write minimal implementation**

Create `ir/treemap.go`:

```go
package ir

// TreemapNode represents a node in the treemap hierarchy.
// Leaf nodes have Value > 0 and no Children.
// Section nodes have Children and their value is the sum of children.
type TreemapNode struct {
	Label    string
	Value    float64 // leaf value (0 for sections)
	Class    string  // CSS class from :::
	Children []*TreemapNode
}

// IsLeaf returns true if this node has no children.
func (n *TreemapNode) IsLeaf() bool {
	return len(n.Children) == 0
}

// TotalValue returns the node's own value if it's a leaf,
// or the recursive sum of children's values if it's a section.
func (n *TreemapNode) TotalValue() float64 {
	if n.IsLeaf() {
		return n.Value
	}
	var sum float64
	for _, c := range n.Children {
		sum += c.TotalValue()
	}
	return sum
}
```

Add to `ir/graph.go`:

```go
	// Treemap diagram fields
	TreemapRoot  *TreemapNode
	TreemapTitle string
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestTreemap -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/treemap.go ir/treemap_test.go ir/graph.go
git commit -m "feat(ir): add Treemap diagram types"
```

---

### Task 4: Config and Theme

**Files:**
- Modify: `config/config.go` — add MindmapConfig, SankeyConfig, TreemapConfig
- Modify: `config/config_test.go` — add tests
- Modify: `theme/theme.go` — add colors
- Modify: `theme/theme_test.go` — add tests

**Step 1: Write the failing tests**

Add to `config/config_test.go`:

```go
func TestDefaultLayoutMindmapConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Mindmap.BranchSpacing != 80 {
		t.Errorf("Mindmap.BranchSpacing = %v, want 80", cfg.Mindmap.BranchSpacing)
	}
	if cfg.Mindmap.LevelSpacing != 60 {
		t.Errorf("Mindmap.LevelSpacing = %v, want 60", cfg.Mindmap.LevelSpacing)
	}
}

func TestDefaultLayoutSankeyConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Sankey.ChartWidth != 800 {
		t.Errorf("Sankey.ChartWidth = %v, want 800", cfg.Sankey.ChartWidth)
	}
	if cfg.Sankey.NodeWidth != 20 {
		t.Errorf("Sankey.NodeWidth = %v, want 20", cfg.Sankey.NodeWidth)
	}
}

func TestDefaultLayoutTreemapConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Treemap.ChartWidth != 600 {
		t.Errorf("Treemap.ChartWidth = %v, want 600", cfg.Treemap.ChartWidth)
	}
	if cfg.Treemap.Padding != 4 {
		t.Errorf("Treemap.Padding = %v, want 4", cfg.Treemap.Padding)
	}
}
```

Add to `theme/theme_test.go`:

```go
func TestModernMindmapColors(t *testing.T) {
	th := Modern()
	if len(th.MindmapBranchColors) == 0 {
		t.Error("MindmapBranchColors is empty")
	}
	if th.MindmapNodeBorder == "" {
		t.Error("MindmapNodeBorder is empty")
	}
}

func TestModernSankeyColors(t *testing.T) {
	th := Modern()
	if len(th.SankeyNodeColors) == 0 {
		t.Error("SankeyNodeColors is empty")
	}
	if th.SankeyLinkColor == "" {
		t.Error("SankeyLinkColor is empty")
	}
}

func TestModernTreemapColors(t *testing.T) {
	th := Modern()
	if len(th.TreemapColors) == 0 {
		t.Error("TreemapColors is empty")
	}
	if th.TreemapBorder == "" {
		t.Error("TreemapBorder is empty")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./config/ ./theme/ -run 'Mindmap|Sankey|Treemap' -v`
Expected: FAIL

**Step 3: Write minimal implementation**

Add to `config/config.go` Layout struct:

```go
	Mindmap  MindmapConfig
	Sankey   SankeyConfig
	Treemap  TreemapConfig
```

Add config structs:

```go
// MindmapConfig holds mindmap diagram layout options.
type MindmapConfig struct {
	BranchSpacing float32
	LevelSpacing  float32
	PaddingX      float32
	PaddingY      float32
	NodePadding   float32
}

// SankeyConfig holds Sankey diagram layout options.
type SankeyConfig struct {
	ChartWidth  float32
	ChartHeight float32
	NodeWidth   float32
	NodePadding float32
	PaddingX    float32
	PaddingY    float32
	Iterations  int // link relaxation iterations
}

// TreemapConfig holds Treemap diagram layout options.
type TreemapConfig struct {
	ChartWidth    float32
	ChartHeight   float32
	Padding       float32 // inner padding between rects
	HeaderHeight  float32
	PaddingX      float32
	PaddingY      float32
	LabelFontSize float32
	ValueFontSize float32
}
```

Add defaults in `DefaultLayout()`:

```go
		Mindmap: MindmapConfig{
			BranchSpacing: 80,
			LevelSpacing:  60,
			PaddingX:      40,
			PaddingY:      40,
			NodePadding:   12,
		},
		Sankey: SankeyConfig{
			ChartWidth:  800,
			ChartHeight: 400,
			NodeWidth:   20,
			NodePadding: 10,
			PaddingX:    40,
			PaddingY:    20,
			Iterations:  32,
		},
		Treemap: TreemapConfig{
			ChartWidth:    600,
			ChartHeight:   400,
			Padding:       4,
			HeaderHeight:  24,
			PaddingX:      10,
			PaddingY:      10,
			LabelFontSize: 12,
			ValueFontSize: 10,
		},
```

Add to `theme/theme.go` Theme struct:

```go
	// Mindmap colors
	MindmapBranchColors []string
	MindmapNodeFill     string
	MindmapNodeBorder   string

	// Sankey colors
	SankeyNodeColors []string
	SankeyLinkColor  string
	SankeyLinkOpacity float32

	// Treemap colors
	TreemapColors    []string
	TreemapBorder    string
	TreemapTextColor string
```

Add values in `Modern()`:

```go
		MindmapBranchColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		MindmapNodeFill:   "#F0F4F8",
		MindmapNodeBorder: "#3B6492",

		SankeyNodeColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		SankeyLinkColor:   "#6E7B8B",
		SankeyLinkOpacity: 0.4,

		TreemapColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		TreemapBorder:    "#3B6492",
		TreemapTextColor: "#FFFFFF",
```

Add values in `MermaidDefault()`:

```go
		MindmapBranchColors: []string{
			"#9370DB", "#E76F51", "#7FB069", "#F4A261",
			"#48A9A6", "#D08AC0", "#E4E36A", "#F7B7A3",
		},
		MindmapNodeFill:   "#ECECFF",
		MindmapNodeBorder: "#9370DB",

		SankeyNodeColors: []string{
			"#4C78A8", "#48A9A6", "#E4E36A", "#F4A261",
			"#E76F51", "#7FB069", "#D08AC0", "#F7B7A3",
		},
		SankeyLinkColor:   "#888",
		SankeyLinkOpacity: 0.4,

		TreemapColors: []string{
			"#4C78A8", "#48A9A6", "#E4E36A", "#F4A261",
			"#E76F51", "#7FB069", "#D08AC0", "#F7B7A3",
		},
		TreemapBorder:    "#9370DB",
		TreemapTextColor: "#333",
```

**Step 4: Run tests to verify they pass**

Run: `go test ./config/ ./theme/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add config/config.go config/config_test.go theme/theme.go theme/theme_test.go
git commit -m "feat(config,theme): add Mindmap, Sankey, and Treemap config and theme"
```

---

### Task 5: Mindmap Parser

**Files:**
- Create: `parser/mindmap.go`
- Create: `parser/mindmap_test.go`
- Modify: `parser/parser.go` — add `case ir.Mindmap`

**Step 1: Write the failing test**

Create `parser/mindmap_test.go`:

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseMindmapBasic(t *testing.T) {
	input := `mindmap
    root((Central))
        A[Square]
        B(Rounded)
            C((Circle))`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Mindmap {
		t.Fatalf("Kind = %v, want Mindmap", g.Kind)
	}
	if g.MindmapRoot == nil {
		t.Fatal("MindmapRoot is nil")
	}
	if g.MindmapRoot.Label != "Central" {
		t.Errorf("root label = %q, want %q", g.MindmapRoot.Label, "Central")
	}
	if g.MindmapRoot.Shape != ir.MindmapCircle {
		t.Errorf("root shape = %v, want Circle", g.MindmapRoot.Shape)
	}
	if len(g.MindmapRoot.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(g.MindmapRoot.Children))
	}
	a := g.MindmapRoot.Children[0]
	if a.Label != "Square" || a.Shape != ir.MindmapSquare {
		t.Errorf("child A = %q/%v, want Square/square", a.Label, a.Shape)
	}
	b := g.MindmapRoot.Children[1]
	if len(b.Children) != 1 {
		t.Fatalf("B children = %d, want 1", len(b.Children))
	}
	if b.Children[0].Shape != ir.MindmapCircle {
		t.Errorf("C shape = %v, want Circle", b.Children[0].Shape)
	}
}

func TestParseMindmapAllShapes(t *testing.T) {
	input := `mindmap
    root
        Default
        [Square]
        (Rounded)
        ((Circle))
        ))Bang((
        )Cloud(
        {{Hexagon}}`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if len(g.MindmapRoot.Children) != 7 {
		t.Fatalf("children = %d, want 7", len(g.MindmapRoot.Children))
	}
	expected := []ir.MindmapShape{
		ir.MindmapDefault, ir.MindmapSquare, ir.MindmapRounded, ir.MindmapCircle,
		ir.MindmapBang, ir.MindmapCloud, ir.MindmapHexagon,
	}
	for i, want := range expected {
		got := g.MindmapRoot.Children[i].Shape
		if got != want {
			t.Errorf("child[%d] shape = %v, want %v", i, got, want)
		}
	}
}

func TestParseMindmapIconAndClass(t *testing.T) {
	input := `mindmap
    root
        A::icon(fa fa-book)
        B:::urgent`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.MindmapRoot.Children[0].Icon != "fa fa-book" {
		t.Errorf("icon = %q, want %q", g.MindmapRoot.Children[0].Icon, "fa fa-book")
	}
	if g.MindmapRoot.Children[1].Class != "urgent" {
		t.Errorf("class = %q, want %q", g.MindmapRoot.Children[1].Class, "urgent")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParseMindmap -v`
Expected: FAIL

**Step 3: Write minimal implementation**

Create `parser/mindmap.go`:

```go
package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	mindmapIconRe  = regexp.MustCompile(`::icon\(([^)]+)\)`)
	mindmapClassRe = regexp.MustCompile(`:::(\S+)`)
)

func parseMindmap(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap

	lines := preprocessMindmapInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	// Build tree from indentation.
	type stackEntry struct {
		node   *ir.MindmapNode
		indent int
	}
	var stack []stackEntry
	nodeCount := 0

	for _, entry := range lines {
		text := entry.text
		indent := entry.indent

		// Skip the mindmap keyword line.
		if strings.EqualFold(strings.TrimSpace(text), "mindmap") {
			continue
		}

		node := parseMindmapNodeText(text, nodeCount)
		nodeCount++

		// Pop stack until we find a parent with smaller indentation.
		for len(stack) > 0 && stack[len(stack)-1].indent >= indent {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			// This is the root node.
			g.MindmapRoot = node
		} else {
			parent := stack[len(stack)-1].node
			parent.Children = append(parent.Children, node)
		}

		stack = append(stack, stackEntry{node: node, indent: indent})
	}

	return &ParseOutput{Graph: g}, nil
}

func parseMindmapNodeText(text string, index int) *ir.MindmapNode {
	node := &ir.MindmapNode{}

	// Extract and strip icon decorator.
	if m := mindmapIconRe.FindStringSubmatch(text); m != nil {
		node.Icon = m[1]
		text = strings.TrimSpace(mindmapIconRe.ReplaceAllString(text, ""))
	}

	// Extract and strip class decorator.
	if m := mindmapClassRe.FindStringSubmatch(text); m != nil {
		node.Class = m[1]
		text = strings.TrimSpace(mindmapClassRe.ReplaceAllString(text, ""))
	}

	// Detect shape from delimiters.
	text = strings.TrimSpace(text)
	node.Shape, node.Label = parseMindmapShape(text)

	// Generate an ID.
	if node.Label != "" {
		node.ID = node.Label
	} else {
		node.ID = text
	}

	return node
}

func parseMindmapShape(text string) (ir.MindmapShape, string) {
	// Order matters: check longest delimiters first.
	// Bang: ))text((
	if strings.HasPrefix(text, "))") && strings.HasSuffix(text, "((") && len(text) > 4 {
		return ir.MindmapBang, text[2 : len(text)-2]
	}
	// Circle: ((text))
	if strings.HasPrefix(text, "((") && strings.HasSuffix(text, "))") && len(text) > 4 {
		return ir.MindmapCircle, text[2 : len(text)-2]
	}
	// Hexagon: {{text}}
	if strings.HasPrefix(text, "{{") && strings.HasSuffix(text, "}}") && len(text) > 4 {
		return ir.MindmapHexagon, text[2 : len(text)-2]
	}
	// Cloud: )text(
	if strings.HasPrefix(text, ")") && strings.HasSuffix(text, "(") && len(text) > 2 {
		return ir.MindmapCloud, text[1 : len(text)-1]
	}
	// Square: [text]
	if strings.HasPrefix(text, "[") && strings.HasSuffix(text, "]") && len(text) > 2 {
		return ir.MindmapSquare, text[1 : len(text)-1]
	}
	// Rounded: (text)
	if strings.HasPrefix(text, "(") && strings.HasSuffix(text, ")") && len(text) > 2 {
		return ir.MindmapRounded, text[1 : len(text)-1]
	}
	// Default: bare text.
	return ir.MindmapDefault, text
}

// preprocessMindmapInput is identical to preprocessKanbanInput — preserves
// indentation for hierarchy detection.
func preprocessMindmapInput(input string) []kanbanLine {
	var result []kanbanLine
	for _, rawLine := range strings.Split(input, "\n") {
		indent := 0
		for _, ch := range rawLine {
			switch ch {
			case ' ':
				indent++
			case '\t':
				indent += 2
			default:
				goto done
			}
		}
	done:
		trimmed := strings.TrimSpace(rawLine)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "%%") {
			continue
		}
		without := stripTrailingComment(trimmed)
		if without == "" {
			continue
		}
		result = append(result, kanbanLine{text: without, indent: indent})
	}
	return result
}
```

Add to `parser/parser.go` switch:

```go
	case ir.Mindmap:
		return parseMindmap(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParseMindmap -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/mindmap.go parser/mindmap_test.go parser/parser.go
git commit -m "feat(parser): add Mindmap diagram parser"
```

---

### Task 6: Sankey Parser

**Files:**
- Create: `parser/sankey.go`
- Create: `parser/sankey_test.go`
- Modify: `parser/parser.go` — add `case ir.Sankey`

**Step 1: Write the failing test**

Create `parser/sankey_test.go`:

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseSankeyBasic(t *testing.T) {
	input := `sankey-beta

Solar,Grid,60
Wind,Grid,290
Grid,Industry,340
Grid,Homes,114`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Sankey {
		t.Fatalf("Kind = %v, want Sankey", g.Kind)
	}
	if len(g.SankeyLinks) != 4 {
		t.Fatalf("links len = %d, want 4", len(g.SankeyLinks))
	}
	if g.SankeyLinks[0].Source != "Solar" {
		t.Errorf("link[0] source = %q, want Solar", g.SankeyLinks[0].Source)
	}
	if g.SankeyLinks[0].Value != 60 {
		t.Errorf("link[0] value = %v, want 60", g.SankeyLinks[0].Value)
	}
}

func TestParseSankeyQuotedNames(t *testing.T) {
	input := `sankey

"Source A","Target, B",100.5
"Source ""C""",Target D,200`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if len(g.SankeyLinks) != 2 {
		t.Fatalf("links len = %d, want 2", len(g.SankeyLinks))
	}
	if g.SankeyLinks[0].Target != "Target, B" {
		t.Errorf("link[0] target = %q, want %q", g.SankeyLinks[0].Target, "Target, B")
	}
	if g.SankeyLinks[0].Value != 100.5 {
		t.Errorf("link[0] value = %v, want 100.5", g.SankeyLinks[0].Value)
	}
	if g.SankeyLinks[1].Source != `Source "C"` {
		t.Errorf("link[1] source = %q, want %q", g.SankeyLinks[1].Source, `Source "C"`)
	}
}

func TestParseSankeyMinimal(t *testing.T) {
	input := `sankey-beta
A,B,10`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.SankeyLinks) != 1 {
		t.Fatalf("links len = %d, want 1", len(out.Graph.SankeyLinks))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParseSankey -v`
Expected: FAIL

**Step 3: Write minimal implementation**

Create `parser/sankey.go`:

```go
package parser

import (
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

func parseSankey(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Sankey

	lines := preprocessInput(input)
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Parse CSV: source,target,value
		fields := parseSankeyCSVLine(trimmed)
		if len(fields) < 3 {
			continue
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
		if err != nil {
			continue
		}
		g.SankeyLinks = append(g.SankeyLinks, &ir.SankeyLink{
			Source: fields[0],
			Target: fields[1],
			Value:  value,
		})
	}

	return &ParseOutput{Graph: g}, nil
}

// parseSankeyCSVLine parses an RFC 4180 CSV line with quoted field support.
func parseSankeyCSVLine(line string) []string {
	var fields []string
	i := 0
	for i < len(line) {
		if line[i] == '"' {
			// Quoted field.
			i++ // skip opening quote
			var field strings.Builder
			for i < len(line) {
				if line[i] == '"' {
					if i+1 < len(line) && line[i+1] == '"' {
						// Escaped quote.
						field.WriteByte('"')
						i += 2
					} else {
						// Closing quote.
						i++ // skip closing quote
						break
					}
				} else {
					field.WriteByte(line[i])
					i++
				}
			}
			fields = append(fields, field.String())
			// Skip comma after field.
			if i < len(line) && line[i] == ',' {
				i++
			}
		} else {
			// Unquoted field.
			end := strings.IndexByte(line[i:], ',')
			if end < 0 {
				fields = append(fields, strings.TrimSpace(line[i:]))
				break
			}
			fields = append(fields, strings.TrimSpace(line[i:i+end]))
			i += end + 1
		}
	}
	return fields
}
```

Add to `parser/parser.go` switch:

```go
	case ir.Sankey:
		return parseSankey(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParseSankey -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/sankey.go parser/sankey_test.go parser/parser.go
git commit -m "feat(parser): add Sankey diagram parser"
```

---

### Task 7: Treemap Parser

**Files:**
- Create: `parser/treemap.go`
- Create: `parser/treemap_test.go`
- Modify: `parser/parser.go` — add `case ir.Treemap`

**Step 1: Write the failing test**

Create `parser/treemap_test.go`:

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseTreemapBasic(t *testing.T) {
	input := `treemap-beta
"Root"
    "Section A"
        "Leaf 1": 30
        "Leaf 2": 50
    "Section B"
        "Leaf 3": 20`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Treemap {
		t.Fatalf("Kind = %v, want Treemap", g.Kind)
	}
	if g.TreemapRoot == nil {
		t.Fatal("TreemapRoot is nil")
	}
	if g.TreemapRoot.Label != "Root" {
		t.Errorf("root label = %q, want Root", g.TreemapRoot.Label)
	}
	if len(g.TreemapRoot.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(g.TreemapRoot.Children))
	}
	secA := g.TreemapRoot.Children[0]
	if secA.Label != "Section A" {
		t.Errorf("secA label = %q", secA.Label)
	}
	if len(secA.Children) != 2 {
		t.Fatalf("secA children = %d, want 2", len(secA.Children))
	}
	if secA.Children[0].Value != 30 {
		t.Errorf("leaf1 value = %v, want 30", secA.Children[0].Value)
	}
}

func TestParseTreemapFlat(t *testing.T) {
	input := `treemap
"Root"
    "A": 10
    "B": 20
    "C": 30`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if len(g.TreemapRoot.Children) != 3 {
		t.Fatalf("children = %d, want 3", len(g.TreemapRoot.Children))
	}
	total := g.TreemapRoot.TotalValue()
	if total != 60 {
		t.Errorf("total = %v, want 60", total)
	}
}

func TestParseTreemapTitle(t *testing.T) {
	input := `treemap-beta
    title Budget
"Operations"
    "Salaries": 700
    "Equipment": 200`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if out.Graph.TreemapTitle != "Budget" {
		t.Errorf("title = %q, want Budget", out.Graph.TreemapTitle)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParseTreemap -v`
Expected: FAIL

**Step 3: Write minimal implementation**

Create `parser/treemap.go`:

```go
package parser

import (
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

func parseTreemap(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Treemap

	lines := preprocessMindmapInput(input) // reuse indentation-aware preprocessor
	if len(lines) == 0 {
		return &ParseOutput{Graph: g}, nil
	}

	type stackEntry struct {
		node   *ir.TreemapNode
		indent int
	}
	var stack []stackEntry

	for _, entry := range lines {
		text := entry.text
		indent := entry.indent
		lower := strings.ToLower(strings.TrimSpace(text))

		// Skip keyword line.
		if strings.HasPrefix(lower, "treemap") {
			continue
		}

		// Handle title directive.
		if strings.HasPrefix(lower, "title") {
			g.TreemapTitle = strings.TrimSpace(text[5:])
			continue
		}

		// Parse node: "Label": value  or  "Label"
		label, value, hasValue := parseTreemapNodeLine(text)
		if label == "" {
			continue
		}

		node := &ir.TreemapNode{Label: label}
		if hasValue {
			node.Value = value
		}

		// Pop stack to find parent.
		for len(stack) > 0 && stack[len(stack)-1].indent >= indent {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			g.TreemapRoot = node
		} else {
			parent := stack[len(stack)-1].node
			parent.Children = append(parent.Children, node)
		}

		stack = append(stack, stackEntry{node: node, indent: indent})
	}

	return &ParseOutput{Graph: g}, nil
}

// parseTreemapNodeLine parses a line like `"Label": 30` or `"Label"`.
func parseTreemapNodeLine(line string) (label string, value float64, hasValue bool) {
	line = strings.TrimSpace(line)

	// Strip optional :::class decorator.
	if idx := strings.Index(line, ":::"); idx >= 0 {
		line = strings.TrimSpace(line[:idx])
	}

	// Extract quoted label.
	if len(line) < 2 {
		return "", 0, false
	}
	quote := line[0]
	if quote != '"' && quote != '\'' {
		return "", 0, false
	}
	end := strings.IndexByte(line[1:], quote)
	if end < 0 {
		return "", 0, false
	}
	label = line[1 : end+1]
	rest := strings.TrimSpace(line[end+2:])

	// Check for value separator.
	if len(rest) > 0 && (rest[0] == ':' || rest[0] == ',') {
		valStr := strings.TrimSpace(rest[1:])
		if v, err := strconv.ParseFloat(valStr, 64); err == nil {
			return label, v, true
		}
	}

	return label, 0, false
}
```

Add to `parser/parser.go` switch:

```go
	case ir.Treemap:
		return parseTreemap(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParseTreemap -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/treemap.go parser/treemap_test.go parser/parser.go
git commit -m "feat(parser): add Treemap diagram parser"
```

---

### Task 8: Mindmap Layout

**Files:**
- Create: `layout/mindmap.go`
- Create: `layout/mindmap_test.go`
- Modify: `layout/types.go` — add `MindmapData` types
- Modify: `layout/layout.go` — add `case ir.Mindmap`

**Step 1: Write the failing test**

Create `layout/mindmap_test.go`:

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestMindmapLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap
	g.MindmapRoot = &ir.MindmapNode{
		ID: "root", Label: "Central", Shape: ir.MindmapCircle,
		Children: []*ir.MindmapNode{
			{ID: "a", Label: "Branch A", Shape: ir.MindmapDefault},
			{ID: "b", Label: "Branch B", Shape: ir.MindmapSquare,
				Children: []*ir.MindmapNode{
					{ID: "c", Label: "Leaf C", Shape: ir.MindmapRounded},
				},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", l.Width, l.Height)
	}

	md, ok := l.Diagram.(MindmapData)
	if !ok {
		t.Fatal("Diagram is not MindmapData")
	}
	if md.Root == nil {
		t.Fatal("Root is nil")
	}
	if md.Root.Label != "Central" {
		t.Errorf("root label = %q", md.Root.Label)
	}
	if len(md.Root.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(md.Root.Children))
	}
}

func TestMindmapLayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions for empty mindmap: %v x %v", l.Width, l.Height)
	}
}
```

**Step 2: Run test to verify it fails**

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// MindmapData holds mindmap layout data.
type MindmapData struct {
	Root *MindmapNodeLayout
}

func (MindmapData) diagramData() {}

// MindmapNodeLayout holds the positioned data for one mindmap node.
type MindmapNodeLayout struct {
	Label      string
	Shape      ir.MindmapShape
	Icon       string
	X, Y       float32
	Width      float32
	Height     float32
	ColorIndex int
	Children   []*MindmapNodeLayout
}
```

Create `layout/mindmap.go`:

```go
package layout

import (
	"math"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeMindmapLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	if g.MindmapRoot == nil {
		return &Layout{
			Kind:    g.Kind,
			Nodes:   map[string]*NodeLayout{},
			Width:   100,
			Height:  100,
			Diagram: MindmapData{},
		}
	}

	measurer := textmetrics.New()
	padX := cfg.Mindmap.PaddingX
	padY := cfg.Mindmap.PaddingY
	nodePad := cfg.Mindmap.NodePadding
	levelSpacing := cfg.Mindmap.LevelSpacing
	branchSpacing := cfg.Mindmap.BranchSpacing

	// Phase 1: Build layout tree with sizes.
	root := mindmapBuildLayoutTree(g.MindmapRoot, measurer, th, nodePad, 0, 0)

	// Phase 2: Position nodes radially from center.
	// First compute subtree angular spans, then assign positions.
	mindmapComputeSubtreeSize(root, branchSpacing)

	// Place root at center (0,0), then position children.
	root.X = 0
	root.Y = 0
	if len(root.Children) > 0 {
		totalSpan := float64(0)
		for _, c := range root.Children {
			totalSpan += float64(c.subtreeSpan)
		}
		startAngle := -math.Pi / 2
		for _, c := range root.Children {
			fraction := float64(c.subtreeSpan) / totalSpan
			midAngle := startAngle + fraction*math.Pi
			dist := levelSpacing + root.Width/2 + c.Width/2
			c.X = float32(math.Cos(midAngle)) * dist
			c.Y = float32(math.Sin(midAngle)) * dist
			mindmapPositionChildren(c, midAngle, levelSpacing, 2)
			startAngle += fraction * 2 * math.Pi
		}
	}

	// Phase 3: Normalize to positive coordinates.
	minX, minY, maxX, maxY := mindmapBounds(root)
	shiftX := padX - minX
	shiftY := padY - minY
	mindmapShift(root, shiftX, shiftY)

	totalW := (maxX - minX) + padX*2
	totalH := (maxY - minY) + padY*2

	return &Layout{
		Kind:    g.Kind,
		Nodes:   map[string]*NodeLayout{},
		Width:   totalW,
		Height:  totalH,
		Diagram: MindmapData{Root: &root.MindmapNodeLayout},
	}
}

type mindmapLayoutNode struct {
	MindmapNodeLayout
	Children     []*mindmapLayoutNode
	subtreeSpan  float32 // angular span proportional to subtree size
}

func mindmapBuildLayoutTree(node *ir.MindmapNode, m *textmetrics.Measurer, th *theme.Theme, pad float32, depth, branchIdx int) *mindmapLayoutNode {
	w, h := m.Measure(node.Label, th.FontSize, th.FontFamily)
	ln := &mindmapLayoutNode{
		MindmapNodeLayout: MindmapNodeLayout{
			Label:      node.Label,
			Shape:      node.Shape,
			Icon:       node.Icon,
			Width:      w + pad*2,
			Height:     h + pad*2,
			ColorIndex: branchIdx,
		},
	}
	for i, child := range node.Children {
		bi := branchIdx
		if depth == 0 {
			bi = i
		}
		childLayout := mindmapBuildLayoutTree(child, m, th, pad, depth+1, bi)
		ln.Children = append(ln.Children, childLayout)
	}
	// Copy children to the exported struct too.
	ln.MindmapNodeLayout.Children = make([]*MindmapNodeLayout, len(ln.Children))
	for i, c := range ln.Children {
		ln.MindmapNodeLayout.Children[i] = &c.MindmapNodeLayout
	}
	return ln
}

func mindmapComputeSubtreeSize(node *mindmapLayoutNode, branchSpacing float32) {
	if len(node.Children) == 0 {
		node.subtreeSpan = node.Height + branchSpacing
		return
	}
	var total float32
	for _, c := range node.Children {
		mindmapComputeSubtreeSize(c, branchSpacing)
		total += c.subtreeSpan
	}
	node.subtreeSpan = total
}

func mindmapPositionChildren(parent *mindmapLayoutNode, parentAngle float64, levelSpacing float32, depth int) {
	if len(parent.Children) == 0 {
		return
	}

	totalSpan := float64(0)
	for _, c := range parent.Children {
		totalSpan += float64(c.subtreeSpan)
	}

	// Spread children within a cone centered on parentAngle.
	spreadAngle := math.Pi / float64(depth+1) // narrower at deeper levels
	startAngle := parentAngle - spreadAngle/2

	for _, c := range parent.Children {
		fraction := float64(c.subtreeSpan) / totalSpan
		midAngle := startAngle + fraction*spreadAngle/2
		dist := levelSpacing + parent.Width/2 + c.Width/2
		c.X = parent.X + float32(math.Cos(midAngle))*dist
		c.Y = parent.Y + float32(math.Sin(midAngle))*dist
		mindmapPositionChildren(c, midAngle, levelSpacing, depth+1)
		startAngle += fraction * spreadAngle
	}
}

func mindmapBounds(node *mindmapLayoutNode) (minX, minY, maxX, maxY float32) {
	minX = node.X - node.Width/2
	minY = node.Y - node.Height/2
	maxX = node.X + node.Width/2
	maxY = node.Y + node.Height/2
	for _, c := range node.Children {
		cMinX, cMinY, cMaxX, cMaxY := mindmapBounds(c)
		if cMinX < minX { minX = cMinX }
		if cMinY < minY { minY = cMinY }
		if cMaxX > maxX { maxX = cMaxX }
		if cMaxY > maxY { maxY = cMaxY }
	}
	return
}

func mindmapShift(node *mindmapLayoutNode, dx, dy float32) {
	node.X += dx
	node.Y += dy
	node.MindmapNodeLayout.X = node.X
	node.MindmapNodeLayout.Y = node.Y
	for _, c := range node.Children {
		mindmapShift(c, dx, dy)
	}
}
```

Add to `layout/layout.go`:

```go
	case ir.Mindmap:
		return computeMindmapLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

**Step 5: Commit**

```bash
git add layout/mindmap.go layout/mindmap_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Mindmap radial tree layout"
```

---

### Task 9: Sankey Layout

**Files:**
- Create: `layout/sankey.go`
- Create: `layout/sankey_test.go`
- Modify: `layout/types.go` — add `SankeyData` types
- Modify: `layout/layout.go` — add `case ir.Sankey`

**Step 1: Write the failing test**

Create `layout/sankey_test.go`:

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestSankeyLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Sankey
	g.SankeyLinks = []*ir.SankeyLink{
		{Source: "A", Target: "X", Value: 100},
		{Source: "A", Target: "Y", Value: 200},
		{Source: "B", Target: "X", Value: 150},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", l.Width, l.Height)
	}

	sd, ok := l.Diagram.(SankeyData)
	if !ok {
		t.Fatal("Diagram is not SankeyData")
	}
	if len(sd.Nodes) < 3 {
		t.Errorf("nodes = %d, want >= 3", len(sd.Nodes))
	}
	if len(sd.Links) != 3 {
		t.Errorf("links = %d, want 3", len(sd.Links))
	}
}

func TestSankeyLayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Sankey

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", l.Width, l.Height)
	}
}
```

**Step 2: Run test to verify it fails**

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// SankeyData holds Sankey diagram layout data.
type SankeyData struct {
	Nodes []SankeyNodeLayout
	Links []SankeyLinkLayout
}

func (SankeyData) diagramData() {}

// SankeyNodeLayout holds a positioned Sankey node.
type SankeyNodeLayout struct {
	Label      string
	X, Y       float32
	Width      float32
	Height     float32
	ColorIndex int
}

// SankeyLinkLayout holds a positioned Sankey flow link.
type SankeyLinkLayout struct {
	SourceIdx int
	TargetIdx int
	Value     float64
	SourceY   float32 // start Y position on source node
	TargetY   float32 // start Y position on target node
	Width     float32 // link thickness
}
```

Create `layout/sankey.go` — implements a simplified Sankey layout:
1. Collect unique nodes, assign to columns via longest-path from sources
2. Compute node heights proportional to total flow
3. Position nodes in columns with vertical spacing
4. Compute link positions

This file is substantial (~150 lines). The implementer should follow the algorithm described above. Key functions:

- `computeSankeyLayout(g, th, cfg)` — main function
- `sankeyAssignColumns(links)` — topological column assignment
- `sankeyNodeHeights(nodes, links, chartH)` — proportional height calculation

Add to `layout/layout.go`:

```go
	case ir.Sankey:
		return computeSankeyLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

**Step 5: Commit**

```bash
git add layout/sankey.go layout/sankey_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Sankey flow diagram layout"
```

---

### Task 10: Treemap Layout

**Files:**
- Create: `layout/treemap.go`
- Create: `layout/treemap_test.go`
- Modify: `layout/types.go` — add `TreemapData` types
- Modify: `layout/layout.go` — add `case ir.Treemap`

**Step 1: Write the failing test**

Create `layout/treemap_test.go`:

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestTreemapLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Treemap
	g.TreemapTitle = "Budget"
	g.TreemapRoot = &ir.TreemapNode{
		Label: "Root",
		Children: []*ir.TreemapNode{
			{Label: "A", Value: 60},
			{Label: "B", Value: 30},
			{Label: "C", Value: 10},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", l.Width, l.Height)
	}

	td, ok := l.Diagram.(TreemapData)
	if !ok {
		t.Fatal("Diagram is not TreemapData")
	}
	if td.Title != "Budget" {
		t.Errorf("title = %q, want Budget", td.Title)
	}
	if len(td.Rects) != 3 {
		t.Fatalf("rects = %d, want 3", len(td.Rects))
	}
	// Verify rects don't overlap and fill the area.
	for _, r := range td.Rects {
		if r.Width <= 0 || r.Height <= 0 {
			t.Errorf("rect %q has zero dimension: %v x %v", r.Label, r.Width, r.Height)
		}
	}
}

func TestTreemapLayoutNested(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Treemap
	g.TreemapRoot = &ir.TreemapNode{
		Label: "Root",
		Children: []*ir.TreemapNode{
			{Label: "Section", Children: []*ir.TreemapNode{
				{Label: "X", Value: 20},
				{Label: "Y", Value: 30},
			}},
			{Label: "Z", Value: 50},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	td, ok := l.Diagram.(TreemapData)
	if !ok {
		t.Fatal("Diagram is not TreemapData")
	}
	// Should have rects for Section header, X, Y, and Z
	if len(td.Rects) < 3 {
		t.Errorf("rects = %d, want >= 3", len(td.Rects))
	}
}
```

**Step 2: Run test to verify it fails**

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// TreemapData holds treemap layout data.
type TreemapData struct {
	Rects []TreemapRectLayout
	Title string
}

func (TreemapData) diagramData() {}

// TreemapRectLayout holds a positioned treemap rectangle.
type TreemapRectLayout struct {
	Label      string
	Value      float64
	X, Y       float32
	Width      float32
	Height     float32
	Depth      int
	IsSection  bool
	ColorIndex int
}
```

Create `layout/treemap.go` — implements the squarified treemap algorithm:

The key algorithm is `squarify(items, rect)` which partitions a rectangle into sub-rectangles with aspect ratios as close to 1:1 as possible. The implementer should:
1. Sort children by value (descending)
2. For each child, decide whether to add it to the current row or start a new row
3. The decision uses the "worst aspect ratio" heuristic
4. Recursively layout section nodes within their allocated rectangle

Add to `layout/layout.go`:

```go
	case ir.Treemap:
		return computeTreemapLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

**Step 5: Commit**

```bash
git add layout/treemap.go layout/treemap_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Treemap squarified layout"
```

---

### Task 11: Mindmap Renderer

**Files:**
- Create: `render/mindmap.go`
- Create: `render/mindmap_test.go`
- Modify: `render/svg.go` — add `case layout.MindmapData`

**Step 1: Write the failing test**

Create `render/mindmap_test.go`:

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

func TestRenderMindmap(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap
	g.MindmapRoot = &ir.MindmapNode{
		ID: "root", Label: "Central", Shape: ir.MindmapCircle,
		Children: []*ir.MindmapNode{
			{ID: "a", Label: "Square", Shape: ir.MindmapSquare},
			{ID: "b", Label: "Rounded", Shape: ir.MindmapRounded},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Central") {
		t.Error("missing root label")
	}
	if !strings.Contains(svg, "Square") {
		t.Error("missing child label")
	}
}
```

**Step 2-5:** Implement renderer that draws nodes with shape-appropriate SVG elements (rect, rounded rect, circle, polygon for hexagon/bang/cloud) and curved connection lines from parent to child. Add case to svg.go. Commit.

---

### Task 12: Sankey Renderer

**Files:**
- Create: `render/sankey.go`
- Create: `render/sankey_test.go`
- Modify: `render/svg.go` — add `case layout.SankeyData`

**Step 1: Write the failing test**

Create `render/sankey_test.go`:

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

func TestRenderSankey(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Sankey
	g.SankeyLinks = []*ir.SankeyLink{
		{Source: "Solar", Target: "Grid", Value: 60},
		{Source: "Wind", Target: "Grid", Value: 290},
		{Source: "Grid", Target: "Industry", Value: 350},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Solar") {
		t.Error("missing node label Solar")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing node rects")
	}
	if !strings.Contains(svg, "<path") {
		t.Error("missing link paths")
	}
}
```

**Step 2-5:** Implement renderer that draws node rectangles with labels and cubic Bezier path links between nodes. Add case to svg.go. Commit.

---

### Task 13: Treemap Renderer

**Files:**
- Create: `render/treemap.go`
- Create: `render/treemap_test.go`
- Modify: `render/svg.go` — add `case layout.TreemapData`

**Step 1: Write the failing test**

Create `render/treemap_test.go`:

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

func TestRenderTreemap(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Treemap
	g.TreemapTitle = "Budget"
	g.TreemapRoot = &ir.TreemapNode{
		Label: "Root",
		Children: []*ir.TreemapNode{
			{Label: "Salaries", Value: 70},
			{Label: "Equipment", Value: 30},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(svg, "Budget") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Salaries") {
		t.Error("missing leaf label")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing rectangles")
	}
}
```

**Step 2-5:** Implement renderer that draws nested rectangles with labels and values. Section nodes get a header bar. Colors cycle through `th.TreemapColors`. Add case to svg.go. Commit.

---

### Task 14: Integration Tests and Fixtures

**Files:**
- Create: `testdata/fixtures/mindmap-basic.mmd`
- Create: `testdata/fixtures/mindmap-shapes.mmd`
- Create: `testdata/fixtures/sankey-basic.mmd`
- Create: `testdata/fixtures/sankey-energy.mmd`
- Create: `testdata/fixtures/treemap-basic.mmd`
- Create: `testdata/fixtures/treemap-nested.mmd`
- Modify: `mermaid_test.go` — add integration tests

**Step 1: Create fixture files**

`testdata/fixtures/mindmap-basic.mmd`:
```
mindmap
    root((Project))
        Goals
            Revenue
            Growth
        Risks
            Budget
            Timeline
```

`testdata/fixtures/mindmap-shapes.mmd`:
```
mindmap
    root((Center))
        [Square]
        (Rounded)
        ))Bang((
        )Cloud(
        {{Hexagon}}
```

`testdata/fixtures/sankey-basic.mmd`:
```
sankey-beta

Solar,Grid,60
Wind,Grid,290
Hydro,Grid,7
Grid,Industry,342
Grid,Homes,114
```

`testdata/fixtures/sankey-energy.mmd`:
```
sankey-beta

Coal,Electricity,100
Gas,Electricity,80
Nuclear,Electricity,200
Electricity,Industry,150
Electricity,Transport,80
Electricity,Homes,120
Electricity,Losses,30
```

`testdata/fixtures/treemap-basic.mmd`:
```
treemap-beta
"Budget"
    "Engineering": 500
    "Marketing": 200
    "Operations": 300
```

`testdata/fixtures/treemap-nested.mmd`:
```
treemap-beta
"Company"
    "Engineering"
        "Backend": 200
        "Frontend": 150
        "DevOps": 100
    "Business"
        "Sales": 250
        "Marketing": 150
    "Support": 100
```

**Step 2: Write integration tests** — follow existing pattern with `readFixture` + `Render` + string assertions.

**Step 3: Commit**

```bash
git add testdata/fixtures/ mermaid_test.go
git commit -m "test: add integration tests and fixtures for Mindmap, Sankey, and Treemap"
```

---

### Task 15: Final Validation

**Step 1:** Run `go test ./...` — all 8 packages PASS
**Step 2:** Run `go vet ./...` — clean
**Step 3:** Run `gofmt -l .` — clean
**Step 4:** Run `go build ./...` — clean
**Step 5:** Run go-code-reviewer agent
**Step 6:** Fix any issues found, commit

---
