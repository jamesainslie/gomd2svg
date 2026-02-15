# Phase 4: Simple Grid — Kanban & Packet Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Kanban board and Packet header diagram support with full mermaid.js syntax parity.

**Architecture:** Each diagram type gets its own IR types, parser, layout function, and renderer. Neither uses the Sugiyama pipeline — Kanban uses column-major vertical stacking, Packet uses row-major bit-cell wrapping. Both are simple enough to not need a shared grid abstraction. The `detectDiagramKind` already recognizes `kanban` and `packet` keywords.

**Tech Stack:** Go 1.24+, existing `textmetrics`, `theme`, `config` packages.

---

### Task 1: IR types for Kanban and Packet

**Files:**
- Create: `ir/kanban.go`
- Create: `ir/kanban_test.go`
- Create: `ir/packet.go`
- Create: `ir/packet_test.go`
- Modify: `ir/graph.go` (add Columns and Fields)

**Step 1: Write failing tests**

```go
// ir/kanban_test.go
package ir

import "testing"

func TestKanbanPriorityString(t *testing.T) {
	tests := []struct {
		p    KanbanPriority
		want string
	}{
		{PriorityNone, ""},
		{PriorityVeryLow, "Very Low"},
		{PriorityLow, "Low"},
		{PriorityHigh, "High"},
		{PriorityVeryHigh, "Very High"},
	}
	for _, tt := range tests {
		if got := tt.p.String(); got != tt.want {
			t.Errorf("KanbanPriority(%d).String() = %q, want %q", tt.p, got, tt.want)
		}
	}
}

func TestKanbanColumnCards(t *testing.T) {
	col := &KanbanColumn{
		ID:    "todo",
		Label: "Todo",
		Cards: []*KanbanCard{
			{ID: "t1", Label: "Task 1", Priority: PriorityHigh},
			{ID: "t2", Label: "Task 2", Assigned: "alice"},
		},
	}
	if len(col.Cards) != 2 {
		t.Fatalf("len(Cards) = %d, want 2", len(col.Cards))
	}
	if col.Cards[0].Priority != PriorityHigh {
		t.Errorf("Cards[0].Priority = %v, want PriorityHigh", col.Cards[0].Priority)
	}
	if col.Cards[1].Assigned != "alice" {
		t.Errorf("Cards[1].Assigned = %q, want \"alice\"", col.Cards[1].Assigned)
	}
}
```

```go
// ir/packet_test.go
package ir

import "testing"

func TestPacketFieldRange(t *testing.T) {
	f := &PacketField{Start: 0, End: 15, Description: "Source Port"}
	if f.BitWidth() != 16 {
		t.Errorf("BitWidth() = %d, want 16", f.BitWidth())
	}
}

func TestGraphPacketFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Packet
	g.Fields = []*PacketField{
		{Start: 0, End: 15, Description: "Source Port"},
		{Start: 16, End: 31, Description: "Dest Port"},
	}
	if len(g.Fields) != 2 {
		t.Fatalf("len(Fields) = %d, want 2", len(g.Fields))
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./ir/ -run 'TestKanban|TestPacket' -v`
Expected: FAIL — types not defined.

**Step 3: Write implementation**

```go
// ir/kanban.go
package ir

// KanbanPriority represents the priority level of a Kanban card.
type KanbanPriority int

const (
	PriorityNone     KanbanPriority = iota
	PriorityVeryLow
	PriorityLow
	PriorityHigh
	PriorityVeryHigh
)

func (p KanbanPriority) String() string {
	switch p {
	case PriorityVeryLow:
		return "Very Low"
	case PriorityLow:
		return "Low"
	case PriorityHigh:
		return "High"
	case PriorityVeryHigh:
		return "Very High"
	default:
		return ""
	}
}

// KanbanCard represents a single card/task on a Kanban board.
type KanbanCard struct {
	ID          string
	Label       string
	Assigned    string
	Ticket      string
	Priority    KanbanPriority
	Icon        string
	Description string
}

// KanbanColumn represents a column on a Kanban board.
type KanbanColumn struct {
	ID    string
	Label string
	Cards []*KanbanCard
}
```

```go
// ir/packet.go
package ir

// PacketField represents a single field in a network packet header diagram.
type PacketField struct {
	Start       int
	End         int
	Description string
}

// BitWidth returns the number of bits this field spans.
func (f *PacketField) BitWidth() int {
	return f.End - f.Start + 1
}
```

Modify `ir/graph.go` — add to the `Graph` struct after the Sequence diagram fields:

```go
	// Kanban diagram fields
	Columns []*KanbanColumn

	// Packet diagram fields
	Fields []*PacketField
```

No `NewGraph()` changes needed (nil slices are correct zero values).

**Step 4: Run tests to verify they pass**

Run: `go test ./ir/ -run 'TestKanban|TestPacket' -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/kanban.go ir/kanban_test.go ir/packet.go ir/packet_test.go ir/graph.go
git commit -m "feat(ir): add Kanban and Packet diagram types"
```

---

### Task 2: Config for Kanban and Packet

**Files:**
- Modify: `config/config.go`

**Step 1: Write failing test**

```go
// In a temporary test or just verify compilation
// The config types are validated by the parser/layout tests in later tasks.
// For now, verify the defaults compile and have expected values.
```

**Step 2: Add config structs and defaults**

Add to `config/config.go`:

```go
// KanbanConfig holds Kanban diagram layout options.
type KanbanConfig struct {
	Padding      float32
	SectionWidth float32
	CardSpacing  float32
	HeaderHeight float32
}

// PacketConfig holds Packet diagram layout options.
type PacketConfig struct {
	RowHeight  float32
	BitWidth   float32
	BitsPerRow int
	ShowBits   bool
	PaddingX   float32
	PaddingY   float32
}
```

Add fields to `Layout` struct:

```go
	Kanban  KanbanConfig
	Packet  PacketConfig
```

Add defaults in `DefaultLayout()`:

```go
		Kanban: KanbanConfig{
			Padding:      8,
			SectionWidth: 200,
			CardSpacing:  8,
			HeaderHeight: 36,
		},
		Packet: PacketConfig{
			RowHeight:  32,
			BitWidth:   32,
			BitsPerRow: 32,
			ShowBits:   true,
			PaddingX:   5,
			PaddingY:   5,
		},
```

**Step 3: Verify compilation**

Run: `go build ./...`
Expected: success

**Step 4: Commit**

```bash
git add config/config.go
git commit -m "feat(config): add Kanban and Packet layout config"
```

---

### Task 3: Kanban parser

**Files:**
- Create: `parser/kanban.go`
- Create: `parser/kanban_test.go`
- Modify: `parser/parser.go` (add case to `Parse` switch)

**Step 1: Write failing tests**

```go
// parser/kanban_test.go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseKanbanBasic(t *testing.T) {
	input := `kanban
  Todo
    task1[Create tests]
    task2[Write docs]
  Done
    task3[Ship feature]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Kanban {
		t.Fatalf("Kind = %v, want Kanban", g.Kind)
	}
	if len(g.Columns) != 2 {
		t.Fatalf("len(Columns) = %d, want 2", len(g.Columns))
	}
	if g.Columns[0].Label != "Todo" {
		t.Errorf("Columns[0].Label = %q, want \"Todo\"", g.Columns[0].Label)
	}
	if len(g.Columns[0].Cards) != 2 {
		t.Errorf("len(Columns[0].Cards) = %d, want 2", len(g.Columns[0].Cards))
	}
	if g.Columns[0].Cards[0].ID != "task1" {
		t.Errorf("Cards[0].ID = %q, want \"task1\"", g.Columns[0].Cards[0].ID)
	}
	if g.Columns[0].Cards[0].Label != "Create tests" {
		t.Errorf("Cards[0].Label = %q, want \"Create tests\"", g.Columns[0].Cards[0].Label)
	}
	if g.Columns[1].Label != "Done" {
		t.Errorf("Columns[1].Label = %q, want \"Done\"", g.Columns[1].Label)
	}
	if len(g.Columns[1].Cards) != 1 {
		t.Errorf("len(Columns[1].Cards) = %d, want 1", len(g.Columns[1].Cards))
	}
}

func TestParseKanbanMetadata(t *testing.T) {
	input := `kanban
  Backlog
    t1[Fix bug]@{ assigned: 'alice', ticket: 'BUG-42', priority: 'High' }`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	g := out.Graph
	if len(g.Columns) != 1 {
		t.Fatalf("len(Columns) = %d, want 1", len(g.Columns))
	}
	card := g.Columns[0].Cards[0]
	if card.Assigned != "alice" {
		t.Errorf("Assigned = %q, want \"alice\"", card.Assigned)
	}
	if card.Ticket != "BUG-42" {
		t.Errorf("Ticket = %q, want \"BUG-42\"", card.Ticket)
	}
	if card.Priority != ir.PriorityHigh {
		t.Errorf("Priority = %v, want PriorityHigh", card.Priority)
	}
}

func TestParseKanbanColumnNoCards(t *testing.T) {
	input := `kanban
  EmptyCol
  WithCards
    t1[Do something]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(out.Graph.Columns) != 2 {
		t.Fatalf("len(Columns) = %d, want 2", len(out.Graph.Columns))
	}
	if len(out.Graph.Columns[0].Cards) != 0 {
		t.Errorf("EmptyCol should have 0 cards, got %d", len(out.Graph.Columns[0].Cards))
	}
}

func TestParseKanbanColumnWithBracketLabel(t *testing.T) {
	input := `kanban
  col1[In Progress]
    t1[Working on it]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if out.Graph.Columns[0].Label != "In Progress" {
		t.Errorf("Label = %q, want \"In Progress\"", out.Graph.Columns[0].Label)
	}
	if out.Graph.Columns[0].ID != "col1" {
		t.Errorf("ID = %q, want \"col1\"", out.Graph.Columns[0].ID)
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./parser/ -run TestParseKanban -v`
Expected: FAIL

**Step 3: Write implementation**

Key design notes for the implementer:
- **Indentation matters** for Kanban. The existing `preprocessInput()` strips all whitespace, so you CANNOT use it. Instead, process raw lines, stripping only comments and blank lines while preserving leading whitespace.
- Column lines have less indentation than card lines. Use a simple heuristic: if a line has MORE leading whitespace than the previous column line, it's a card; otherwise it's a new column.
- Card syntax: `id[Label]` optionally followed by `@{ key: 'value', ... }`
- Column syntax: either bare `ColumnName` or `id[Column Label]`
- The `@{...}` metadata uses a simple key-value format with single-quoted string values. Parse with string scanning, not `encoding/json`.

Parser structure in `parser/kanban.go`:

```go
package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var kanbanCardRe = regexp.MustCompile(`^(\w+)\[([^\]]+)\](?:\s*@\{(.+)\})?$`)

func parseKanban(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Kanban

	lines := preprocessKanbanInput(input)

	var currentCol *ir.KanbanColumn
	colIndent := -1

	for _, entry := range lines {
		line := entry.text
		indent := entry.indent

		// Skip header
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "kanban") {
			continue
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Determine if this is a column or card based on indentation.
		// First non-header line sets the column indent level.
		if colIndent < 0 || indent <= colIndent {
			// New column
			colIndent = indent
			id, label := parseKanbanColumnHeader(trimmed)
			currentCol = &ir.KanbanColumn{ID: id, Label: label}
			g.Columns = append(g.Columns, currentCol)
			continue
		}

		// Card line (more indented than column)
		if currentCol != nil {
			card := parseKanbanCard(trimmed)
			if card != nil {
				currentCol.Cards = append(currentCol.Cards, card)
			}
		}
	}

	return &ParseOutput{Graph: g}, nil
}
```

Also implement:
- `preprocessKanbanInput(input string) []kanbanLine` — returns lines with indent levels, strips comments/blanks
- `parseKanbanColumnHeader(line string) (id, label string)` — handles bare `Name` and `id[Label]`
- `parseKanbanCard(line string) *ir.KanbanCard` — parses `id[Label]@{...}` with metadata
- `parseKanbanMetadata(raw string) (assigned, ticket, icon, description string, priority ir.KanbanPriority)` — simple key-value scanner

Add to `parser/parser.go` `Parse()` switch:

```go
	case ir.Kanban:
		return parseKanban(input)
```

**Step 4: Run tests to verify they pass**

Run: `go test ./parser/ -run TestParseKanban -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/kanban.go parser/kanban_test.go parser/parser.go
git commit -m "feat(parser): add Kanban diagram parser"
```

---

### Task 4: Packet parser

**Files:**
- Create: `parser/packet.go`
- Create: `parser/packet_test.go`
- Modify: `parser/parser.go` (add case to `Parse` switch)

**Step 1: Write failing tests**

```go
// parser/packet_test.go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParsePacketRangeNotation(t *testing.T) {
	input := `packet
0-15: "Source Port"
16-31: "Destination Port"
32-63: "Sequence Number"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Packet {
		t.Fatalf("Kind = %v, want Packet", g.Kind)
	}
	if len(g.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(g.Fields))
	}
	if g.Fields[0].Start != 0 || g.Fields[0].End != 15 {
		t.Errorf("Fields[0] = %d-%d, want 0-15", g.Fields[0].Start, g.Fields[0].End)
	}
	if g.Fields[0].Description != "Source Port" {
		t.Errorf("Fields[0].Description = %q, want \"Source Port\"", g.Fields[0].Description)
	}
	if g.Fields[2].Start != 32 || g.Fields[2].End != 63 {
		t.Errorf("Fields[2] = %d-%d, want 32-63", g.Fields[2].Start, g.Fields[2].End)
	}
}

func TestParsePacketBitCountNotation(t *testing.T) {
	input := `packet
+16: "Source Port"
+16: "Destination Port"
+32: "Sequence Number"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	g := out.Graph
	if len(g.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(g.Fields))
	}
	if g.Fields[0].Start != 0 || g.Fields[0].End != 15 {
		t.Errorf("Fields[0] = %d-%d, want 0-15", g.Fields[0].Start, g.Fields[0].End)
	}
	if g.Fields[1].Start != 16 || g.Fields[1].End != 31 {
		t.Errorf("Fields[1] = %d-%d, want 16-31", g.Fields[1].Start, g.Fields[1].End)
	}
	if g.Fields[2].Start != 32 || g.Fields[2].End != 63 {
		t.Errorf("Fields[2] = %d-%d, want 32-63", g.Fields[2].Start, g.Fields[2].End)
	}
}

func TestParsePacketMixedNotation(t *testing.T) {
	input := `packet
0-15: "Source Port"
+16: "Destination Port"
32-63: "Sequence Number"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(out.Graph.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(out.Graph.Fields))
	}
	// +16 after 0-15 should be 16-31
	f := out.Graph.Fields[1]
	if f.Start != 16 || f.End != 31 {
		t.Errorf("Fields[1] = %d-%d, want 16-31", f.Start, f.End)
	}
}

func TestParsePacketSingleBit(t *testing.T) {
	input := `packet
0-3: "Version"
+1: "Flag"
+1: "Flag2"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(out.Graph.Fields) != 3 {
		t.Fatalf("len(Fields) = %d, want 3", len(out.Graph.Fields))
	}
	// +1 after 0-3 should be 4-4
	if out.Graph.Fields[1].Start != 4 || out.Graph.Fields[1].End != 4 {
		t.Errorf("Fields[1] = %d-%d, want 4-4", out.Graph.Fields[1].Start, out.Graph.Fields[1].End)
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./parser/ -run TestParsePacket -v`
Expected: FAIL

**Step 3: Write implementation**

```go
// parser/packet.go
package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	packetRangeRe    = regexp.MustCompile(`^(\d+)-(\d+)\s*:\s*"([^"]*)"$`)
	packetBitCountRe = regexp.MustCompile(`^\+(\d+)\s*:\s*"([^"]*)"$`)
)

func parsePacket(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Packet

	lines := preprocessInput(input)
	nextBit := 0

	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "packet") {
			continue
		}

		// Try range notation: 0-15: "Source Port"
		if m := packetRangeRe.FindStringSubmatch(line); m != nil {
			start, _ := strconv.Atoi(m[1])
			end, _ := strconv.Atoi(m[2])
			desc := m[3]
			g.Fields = append(g.Fields, &ir.PacketField{
				Start: start, End: end, Description: desc,
			})
			nextBit = end + 1
			continue
		}

		// Try bit count notation: +16: "Source Port"
		if m := packetBitCountRe.FindStringSubmatch(line); m != nil {
			count, _ := strconv.Atoi(m[1])
			desc := m[2]
			start := nextBit
			end := start + count - 1
			g.Fields = append(g.Fields, &ir.PacketField{
				Start: start, End: end, Description: desc,
			})
			nextBit = end + 1
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}
```

Add to `parser/parser.go` `Parse()` switch:

```go
	case ir.Packet:
		return parsePacket(input)
```

**Step 4: Run tests to verify they pass**

Run: `go test ./parser/ -run TestParsePacket -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/packet.go parser/packet_test.go parser/parser.go
git commit -m "feat(parser): add Packet diagram parser"
```

---

### Task 5: Layout types and Kanban layout

**Files:**
- Modify: `layout/types.go` (add KanbanData, PacketData)
- Create: `layout/kanban.go`
- Create: `layout/kanban_test.go`
- Modify: `layout/layout.go` (add Kanban case)

**Step 1: Write failing tests**

```go
// layout/kanban_test.go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestKanbanLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Kanban
	g.Columns = []*ir.KanbanColumn{
		{ID: "todo", Label: "Todo", Cards: []*ir.KanbanCard{
			{ID: "t1", Label: "Task 1"},
			{ID: "t2", Label: "Task 2"},
		}},
		{ID: "done", Label: "Done", Cards: []*ir.KanbanCard{
			{ID: "t3", Label: "Task 3"},
		}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	kd, ok := l.Diagram.(KanbanData)
	if !ok {
		t.Fatalf("Diagram type = %T, want KanbanData", l.Diagram)
	}
	if len(kd.Columns) != 2 {
		t.Fatalf("len(Columns) = %d, want 2", len(kd.Columns))
	}

	// Columns should be side by side
	if kd.Columns[1].X <= kd.Columns[0].X {
		t.Errorf("Column[1].X (%v) should be > Column[0].X (%v)", kd.Columns[1].X, kd.Columns[0].X)
	}

	// First column should have 2 cards
	if len(kd.Columns[0].Cards) != 2 {
		t.Errorf("len(Columns[0].Cards) = %d, want 2", len(kd.Columns[0].Cards))
	}

	// Cards should be stacked vertically
	if kd.Columns[0].Cards[1].Y <= kd.Columns[0].Cards[0].Y {
		t.Errorf("Card[1].Y (%v) should be > Card[0].Y (%v)",
			kd.Columns[0].Cards[1].Y, kd.Columns[0].Cards[0].Y)
	}

	// Diagram should have positive dimensions
	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("dimensions = %v x %v, want positive", l.Width, l.Height)
	}
}

func TestKanbanLayoutEmptyColumn(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Kanban
	g.Columns = []*ir.KanbanColumn{
		{ID: "empty", Label: "Empty"},
		{ID: "has", Label: "Has Cards", Cards: []*ir.KanbanCard{
			{ID: "t1", Label: "Task"},
		}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	kd := l.Diagram.(KanbanData)
	if len(kd.Columns[0].Cards) != 0 {
		t.Errorf("empty column should have 0 cards, got %d", len(kd.Columns[0].Cards))
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./layout/ -run TestKanban -v`
Expected: FAIL

**Step 3: Write implementation**

Add to `layout/types.go`:

```go
// KanbanData holds Kanban-diagram-specific layout data.
type KanbanData struct {
	Columns []KanbanColumnLayout
}

func (KanbanData) diagramData() {}

// KanbanColumnLayout holds the position of a Kanban column.
type KanbanColumnLayout struct {
	ID     string
	Label  TextBlock
	X, Y   float32
	Width  float32
	Height float32
	Cards  []KanbanCardLayout
}

// KanbanCardLayout holds the position of a single Kanban card.
type KanbanCardLayout struct {
	ID       string
	Label    TextBlock
	Priority ir.KanbanPriority
	X, Y     float32
	Width    float32
	Height   float32
	Metadata map[string]string
}

// PacketData holds Packet-diagram-specific layout data.
type PacketData struct {
	Rows       []PacketRowLayout
	BitsPerRow int
	ShowBits   bool
}

func (PacketData) diagramData() {}

// PacketRowLayout holds the position of a row of packet fields.
type PacketRowLayout struct {
	Y      float32
	Height float32
	Fields []PacketFieldLayout
}

// PacketFieldLayout holds the position of a single packet field cell.
type PacketFieldLayout struct {
	Label    TextBlock
	X, Y     float32
	Width    float32
	Height   float32
	StartBit int
	EndBit   int
}
```

Create `layout/kanban.go`:

```go
package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeKanbanLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	kc := cfg.Kanban
	pad := kc.Padding

	columns := make([]KanbanColumnLayout, len(g.Columns))
	cursorX := pad

	maxColHeight := float32(0)

	for i, col := range g.Columns {
		// Measure column header
		headerTW := measurer.Width(col.Label, th.FontSize, th.FontFamily)
		headerTB := TextBlock{
			Lines: []string{col.Label}, Width: headerTW,
			Height: th.FontSize * cfg.LabelLineHeight, FontSize: th.FontSize,
		}

		// Layout cards vertically within column
		cardY := kc.HeaderHeight + pad
		cards := make([]KanbanCardLayout, len(col.Cards))

		for j, card := range col.Cards {
			cardTW := measurer.Width(card.Label, th.FontSize, th.FontFamily)
			cardH := th.FontSize*cfg.LabelLineHeight + 2*pad
			cardTB := TextBlock{
				Lines: []string{card.Label}, Width: cardTW,
				Height: th.FontSize * cfg.LabelLineHeight, FontSize: th.FontSize,
			}

			// Build metadata map for renderer
			meta := make(map[string]string)
			if card.Assigned != "" {
				meta["assigned"] = card.Assigned
			}
			if card.Ticket != "" {
				meta["ticket"] = card.Ticket
			}
			if card.Icon != "" {
				meta["icon"] = card.Icon
			}

			cards[j] = KanbanCardLayout{
				ID:       card.ID,
				Label:    cardTB,
				Priority: card.Priority,
				X:        cursorX + pad,
				Y:        cardY,
				Width:    kc.SectionWidth - 2*pad,
				Height:   cardH,
				Metadata: meta,
			}

			cardY += cardH + kc.CardSpacing
		}

		colHeight := cardY + pad
		if colHeight < kc.HeaderHeight+2*pad {
			colHeight = kc.HeaderHeight + 2*pad
		}

		columns[i] = KanbanColumnLayout{
			ID:     col.ID,
			Label:  headerTB,
			X:      cursorX,
			Y:      0,
			Width:  kc.SectionWidth,
			Height: colHeight,
			Cards:  cards,
		}

		if colHeight > maxColHeight {
			maxColHeight = colHeight
		}

		cursorX += kc.SectionWidth + pad
	}

	// Normalize all columns to same height
	for i := range columns {
		columns[i].Height = maxColHeight
	}

	totalW := cursorX
	totalH := maxColHeight + pad

	return &Layout{
		Kind:   g.Kind,
		Width:  totalW,
		Height: totalH,
		Diagram: KanbanData{
			Columns: columns,
		},
	}
}
```

Add to `layout/layout.go` `ComputeLayout()` switch:

```go
	case ir.Kanban:
		return computeKanbanLayout(g, th, cfg)
```

**Step 4: Run tests to verify they pass**

Run: `go test ./layout/ -run TestKanban -v`
Expected: PASS

**Step 5: Commit**

```bash
git add layout/types.go layout/kanban.go layout/kanban_test.go layout/layout.go
git commit -m "feat(layout): add Kanban diagram layout"
```

---

### Task 6: Packet layout

**Files:**
- Create: `layout/packet.go`
- Create: `layout/packet_test.go`
- Modify: `layout/layout.go` (add Packet case)

**Step 1: Write failing tests**

```go
// layout/packet_test.go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestPacketLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Packet
	g.Fields = []*ir.PacketField{
		{Start: 0, End: 15, Description: "Source Port"},
		{Start: 16, End: 31, Description: "Dest Port"},
		{Start: 32, End: 63, Description: "Sequence Number"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	pd, ok := l.Diagram.(PacketData)
	if !ok {
		t.Fatalf("Diagram type = %T, want PacketData", l.Diagram)
	}

	// 32 bits per row: first row has 2 fields (0-15, 16-31), second row has 1 (32-63)
	if len(pd.Rows) != 2 {
		t.Fatalf("len(Rows) = %d, want 2", len(pd.Rows))
	}
	if len(pd.Rows[0].Fields) != 2 {
		t.Errorf("len(Rows[0].Fields) = %d, want 2", len(pd.Rows[0].Fields))
	}
	if len(pd.Rows[1].Fields) != 1 {
		t.Errorf("len(Rows[1].Fields) = %d, want 1", len(pd.Rows[1].Fields))
	}

	// Field widths should be proportional to bit count
	f0 := pd.Rows[0].Fields[0]
	f1 := pd.Rows[0].Fields[1]
	if f0.Width != f1.Width {
		t.Errorf("16-bit fields should have equal width: %v vs %v", f0.Width, f1.Width)
	}

	// 32-bit field should span full row width
	f2 := pd.Rows[1].Fields[0]
	expectedFullWidth := float32(cfg.Packet.BitsPerRow) * cfg.Packet.BitWidth
	if f2.Width != expectedFullWidth {
		t.Errorf("32-bit field width = %v, want %v", f2.Width, expectedFullWidth)
	}

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("dimensions = %v x %v, want positive", l.Width, l.Height)
	}
}

func TestPacketLayoutSingleBit(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Packet
	g.Fields = []*ir.PacketField{
		{Start: 0, End: 0, Description: "Flag"},
		{Start: 1, End: 31, Description: "Rest"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	pd := l.Diagram.(PacketData)
	if len(pd.Rows) != 1 {
		t.Fatalf("len(Rows) = %d, want 1", len(pd.Rows))
	}
	if pd.Rows[0].Fields[0].Width != cfg.Packet.BitWidth {
		t.Errorf("1-bit field width = %v, want %v", pd.Rows[0].Fields[0].Width, cfg.Packet.BitWidth)
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./layout/ -run TestPacket -v`
Expected: FAIL

**Step 3: Write implementation**

```go
// layout/packet.go
package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computePacketLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	pc := cfg.Packet
	bitsPerRow := pc.BitsPerRow
	bitW := pc.BitWidth
	rowH := pc.RowHeight
	padX := pc.PaddingX
	padY := pc.PaddingY

	totalRowW := float32(bitsPerRow) * bitW

	// Group fields into rows based on bit positions.
	type rowBuilder struct {
		fields []PacketFieldLayout
		y      float32
	}
	var rows []rowBuilder

	cursorY := padY
	// Optional: if ShowBits, reserve space for bit number labels at top
	if pc.ShowBits {
		cursorY += th.FontSize*cfg.LabelLineHeight + padY
	}

	for _, field := range g.Fields {
		startRow := field.Start / bitsPerRow
		endRow := field.End / bitsPerRow

		for row := startRow; row <= endRow; row++ {
			// Ensure we have enough rows
			for len(rows) <= row {
				rows = append(rows, rowBuilder{
					y: cursorY + float32(len(rows))*(rowH+padY),
				})
			}

			// Compute the bit range within this row
			rowStartBit := row * bitsPerRow
			fieldStartInRow := field.Start
			if fieldStartInRow < rowStartBit {
				fieldStartInRow = rowStartBit
			}
			fieldEndInRow := field.End
			if fieldEndInRow >= rowStartBit+bitsPerRow {
				fieldEndInRow = rowStartBit + bitsPerRow - 1
			}

			bitsInRow := fieldEndInRow - fieldStartInRow + 1
			offsetInRow := fieldStartInRow - rowStartBit

			x := padX + float32(offsetInRow)*bitW
			w := float32(bitsInRow) * bitW

			tw := measurer.Width(field.Description, th.FontSize, th.FontFamily)
			tb := TextBlock{
				Lines:    []string{field.Description},
				Width:    tw,
				Height:   th.FontSize * cfg.LabelLineHeight,
				FontSize: th.FontSize,
			}

			rows[row].fields = append(rows[row].fields, PacketFieldLayout{
				Label:    tb,
				X:        x,
				Y:        rows[row].y,
				Width:    w,
				Height:   rowH,
				StartBit: fieldStartInRow,
				EndBit:   fieldEndInRow,
			})
		}
	}

	// Build final row layouts
	resultRows := make([]PacketRowLayout, len(rows))
	for i, rb := range rows {
		resultRows[i] = PacketRowLayout{
			Y:      rb.y,
			Height: rowH,
			Fields: rb.fields,
		}
	}

	totalH := cursorY + float32(len(rows))*(rowH+padY) + padY
	totalW := totalRowW + 2*padX

	return &Layout{
		Kind:   g.Kind,
		Width:  totalW,
		Height: totalH,
		Diagram: PacketData{
			Rows:       resultRows,
			BitsPerRow: bitsPerRow,
			ShowBits:   pc.ShowBits,
		},
	}
}
```

Add to `layout/layout.go` `ComputeLayout()` switch:

```go
	case ir.Packet:
		return computePacketLayout(g, th, cfg)
```

**Step 4: Run tests to verify they pass**

Run: `go test ./layout/ -run TestPacket -v`
Expected: PASS

**Step 5: Commit**

```bash
git add layout/packet.go layout/packet_test.go layout/layout.go
git commit -m "feat(layout): add Packet diagram layout"
```

---

### Task 7: Kanban renderer

**Files:**
- Create: `render/kanban.go`
- Create: `render/kanban_test.go`
- Modify: `render/svg.go` (add KanbanData case)

**Step 1: Write failing tests**

```go
// render/kanban_test.go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderKanbanContainsColumns(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Kanban
	g.Columns = []*ir.KanbanColumn{
		{ID: "todo", Label: "Todo", Cards: []*ir.KanbanCard{
			{ID: "t1", Label: "Task One"},
		}},
		{ID: "done", Label: "Done"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "Todo") {
		t.Error("SVG should contain column label 'Todo'")
	}
	if !strings.Contains(svg, "Done") {
		t.Error("SVG should contain column label 'Done'")
	}
	if !strings.Contains(svg, "Task One") {
		t.Error("SVG should contain card label 'Task One'")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("SVG should contain rect elements for cards")
	}
}

func TestRenderKanbanValidSVG(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Kanban
	g.Columns = []*ir.KanbanColumn{
		{ID: "col", Label: "Column", Cards: []*ir.KanbanCard{
			{ID: "c1", Label: "Card"},
		}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG should start with <svg")
	}
	if !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
		t.Error("SVG should end with </svg>")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./render/ -run TestRenderKanban -v`
Expected: FAIL

**Step 3: Write implementation**

```go
// render/kanban.go
package render

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func renderKanban(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	kd, ok := l.Diagram.(layout.KanbanData)
	if !ok {
		return
	}

	for _, col := range kd.Columns {
		// Column background
		b.rect(col.X, col.Y, col.Width, col.Height, 4,
			"fill", th.NodeBackground,
			"stroke", th.NodeBorderColor,
			"stroke-width", "1",
		)

		// Column header text
		headerX := col.X + col.Width/2
		headerY := col.Y + cfg.Kanban.HeaderHeight/2 + col.Label.FontSize/3
		b.text(headerX, headerY, col.Label.Lines[0],
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(col.Label.FontSize),
			"font-weight", "bold",
			"fill", th.PrimaryTextColor,
		)

		// Header divider line
		divY := col.Y + cfg.Kanban.HeaderHeight
		b.line(col.X, divY, col.X+col.Width, divY,
			"stroke", th.NodeBorderColor,
			"stroke-width", "1",
		)

		// Cards
		for _, card := range col.Cards {
			b.rect(card.X, card.Y, card.Width, card.Height, 3,
				"fill", th.Background,
				"stroke", th.NodeBorderColor,
				"stroke-width", "1",
			)

			// Card label
			textX := card.X + cfg.Kanban.Padding
			textY := card.Y + card.Height/2 + card.Label.FontSize/3
			b.text(textX, textY, card.Label.Lines[0],
				"font-family", th.FontFamily,
				"font-size", fmtFloat(card.Label.FontSize),
				"fill", th.PrimaryTextColor,
			)
		}
	}
}
```

Add to `render/svg.go` `RenderSVG()` switch:

```go
	case layout.KanbanData:
		renderKanban(&b, l, th, cfg)
```

**Step 4: Run tests to verify they pass**

Run: `go test ./render/ -run TestRenderKanban -v`
Expected: PASS

**Step 5: Commit**

```bash
git add render/kanban.go render/kanban_test.go render/svg.go
git commit -m "feat(render): add Kanban diagram SVG renderer"
```

---

### Task 8: Packet renderer

**Files:**
- Create: `render/packet.go`
- Create: `render/packet_test.go`
- Modify: `render/svg.go` (add PacketData case)

**Step 1: Write failing tests**

```go
// render/packet_test.go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderPacketContainsFields(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Packet
	g.Fields = []*ir.PacketField{
		{Start: 0, End: 15, Description: "Source Port"},
		{Start: 16, End: 31, Description: "Dest Port"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "Source Port") {
		t.Error("SVG should contain 'Source Port'")
	}
	if !strings.Contains(svg, "Dest Port") {
		t.Error("SVG should contain 'Dest Port'")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("SVG should contain rect elements for fields")
	}
}

func TestRenderPacketBitNumbers(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Packet
	g.Fields = []*ir.PacketField{
		{Start: 0, End: 31, Description: "Data"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Should contain bit number "0" and "31"
	if !strings.Contains(svg, ">0<") {
		t.Error("SVG should contain bit number '0'")
	}
	if !strings.Contains(svg, ">31<") {
		t.Error("SVG should contain bit number '31'")
	}
}

func TestRenderPacketValidSVG(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Packet
	g.Fields = []*ir.PacketField{
		{Start: 0, End: 15, Description: "Field"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG should start with <svg")
	}
	if !strings.HasSuffix(strings.TrimSpace(svg), "</svg>") {
		t.Error("SVG should end with </svg>")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./render/ -run TestRenderPacket -v`
Expected: FAIL

**Step 3: Write implementation**

```go
// render/packet.go
package render

import (
	"fmt"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func renderPacket(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	pd, ok := l.Diagram.(layout.PacketData)
	if !ok {
		return
	}

	pc := cfg.Packet
	smallFontSize := th.FontSize * 0.7

	// Render bit numbers at top if enabled
	if pd.ShowBits {
		for bit := 0; bit < pd.BitsPerRow; bit++ {
			x := pc.PaddingX + float32(bit)*pc.BitWidth + pc.BitWidth/2
			y := th.FontSize * cfg.LabelLineHeight
			b.text(x, y, fmt.Sprintf("%d", bit),
				"text-anchor", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(smallFontSize),
				"fill", th.SecondaryTextColor,
			)
		}
	}

	// Render rows and fields
	for _, row := range pd.Rows {
		for _, field := range row.Fields {
			// Field rectangle
			b.rect(field.X, field.Y, field.Width, field.Height, 0,
				"fill", th.NodeBackground,
				"stroke", th.NodeBorderColor,
				"stroke-width", "1",
			)

			// Field label (centered)
			textX := field.X + field.Width/2
			textY := field.Y + field.Height/2 + field.Label.FontSize/3
			b.text(textX, textY, field.Label.Lines[0],
				"text-anchor", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(field.Label.FontSize),
				"fill", th.PrimaryTextColor,
			)

			// Bit range labels at bottom-left and bottom-right of field
			bitY := field.Y + field.Height - 2
			b.text(field.X+2, bitY, fmt.Sprintf("%d", field.StartBit),
				"font-family", th.FontFamily,
				"font-size", fmtFloat(smallFontSize),
				"fill", th.SecondaryTextColor,
			)
			if field.EndBit != field.StartBit {
				b.text(field.X+field.Width-2, bitY, fmt.Sprintf("%d", field.EndBit),
					"text-anchor", "end",
					"font-family", th.FontFamily,
					"font-size", fmtFloat(smallFontSize),
					"fill", th.SecondaryTextColor,
				)
			}
		}
	}
}
```

Add to `render/svg.go` `RenderSVG()` switch:

```go
	case layout.PacketData:
		renderPacket(&b, l, th, cfg)
```

**Step 4: Run tests to verify they pass**

Run: `go test ./render/ -run TestRenderPacket -v`
Expected: PASS

**Step 5: Commit**

```bash
git add render/packet.go render/packet_test.go render/svg.go
git commit -m "feat(render): add Packet diagram SVG renderer"
```

---

### Task 9: Integration tests, fixtures, and benchmarks

**Files:**
- Create: `testdata/fixtures/kanban-basic.mmd`
- Create: `testdata/fixtures/kanban-metadata.mmd`
- Create: `testdata/fixtures/packet-tcp.mmd`
- Create: `testdata/fixtures/packet-bitcount.mmd`
- Modify: `mermaid_test.go` (add integration tests)
- Modify: `mermaid_bench_test.go` (add benchmarks)

**Step 1: Create fixture files**

```
// testdata/fixtures/kanban-basic.mmd
kanban
  Todo
    task1[Write tests]
    task2[Write docs]
  InProgress
    task3[Build feature]
  Done
    task4[Ship v1]
```

```
// testdata/fixtures/kanban-metadata.mmd
kanban
  Backlog
    t1[Fix login bug]@{ assigned: 'alice', ticket: 'BUG-101', priority: 'High' }
    t2[Update deps]@{ priority: 'Low' }
  InProgress
    t3[Add auth]@{ assigned: 'bob', ticket: 'FEAT-42', priority: 'Very High' }
  Done
    t4[Deploy v2]@{ ticket: 'REL-10' }
```

```
// testdata/fixtures/packet-tcp.mmd
packet
0-15: "Source Port"
16-31: "Destination Port"
32-63: "Sequence Number"
64-95: "Acknowledgment Number"
96-99: "Data Offset"
100-105: "Reserved"
106-111: "Flags"
112-127: "Window"
128-143: "Checksum"
144-159: "Urgent Pointer"
```

```
// testdata/fixtures/packet-bitcount.mmd
packet
+16: "Source Port"
+16: "Destination Port"
+16: "Length"
+16: "Checksum"
```

**Step 2: Write integration tests**

Add to `mermaid_test.go`:

```go
func TestRenderKanbanDiagram(t *testing.T) {
	input := readFixture(t, "kanban-basic.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Todo") {
		t.Error("expected column label 'Todo'")
	}
	if !strings.Contains(svg, "Write tests") {
		t.Error("expected card label 'Write tests'")
	}
}

func TestRenderKanbanMetadata(t *testing.T) {
	input := readFixture(t, "kanban-metadata.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "Fix login bug") {
		t.Error("expected card label")
	}
}

func TestRenderPacketDiagram(t *testing.T) {
	input := readFixture(t, "packet-tcp.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Source Port") {
		t.Error("expected field label 'Source Port'")
	}
}

func TestRenderPacketBitCount(t *testing.T) {
	input := readFixture(t, "packet-bitcount.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "Checksum") {
		t.Error("expected field label 'Checksum'")
	}
}
```

**Step 3: Write benchmarks**

Add to `mermaid_bench_test.go`:

```go
func BenchmarkRenderKanban(b *testing.B) {
	input := readBenchFixture(b, "kanban-basic.mmd")
	b.ReportAllocs()
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderPacket(b *testing.B) {
	input := readBenchFixture(b, "packet-tcp.mmd")
	b.ReportAllocs()
	for b.Loop() {
		_, err := Render(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
```

**Step 4: Run all tests**

Run: `go test ./... -v`
Expected: ALL PASS

**Step 5: Run benchmarks**

Run: `go test -bench=BenchmarkRender -benchmem -count=1`

**Step 6: Commit**

```bash
git add testdata/fixtures/kanban-basic.mmd testdata/fixtures/kanban-metadata.mmd testdata/fixtures/packet-tcp.mmd testdata/fixtures/packet-bitcount.mmd mermaid_test.go mermaid_bench_test.go
git commit -m "test: add integration tests, fixtures, and benchmarks for Kanban and Packet"
```

---

### Task 10: Final validation

**Step 1: Run full test suite**

Run: `go test ./... -v -count=1`
Expected: ALL PASS

**Step 2: Check formatting**

Run: `gofmt -l .`
Expected: no files listed (all formatted)

**Step 3: Run go vet**

Run: `go vet ./...`
Expected: no issues

**Step 4: Run benchmarks**

Run: `go test -bench=BenchmarkRender -benchmem -count=1`
Expected: Kanban and Packet benchmarks run successfully
