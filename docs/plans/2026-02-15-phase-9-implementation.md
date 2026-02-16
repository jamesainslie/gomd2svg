# Phase 9: Requirement, Block & C4 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Requirement, Block, and C4 diagram support — three diagram types that reuse Sugiyama layout (Block uses hybrid grid/Sugiyama).

**Architecture:** Each diagram type follows the 4-layer pipeline (IR types -> Parser -> Layout -> Render). Requirement and C4 reuse `runSugiyama()` directly. Block uses grid layout when `columns` is specified, falling back to Sugiyama for connected graphs. All three add dispatch cases to `parser/parser.go`, `layout/layout.go`, and `render/svg.go`.

**Tech Stack:** Go, `textmetrics` for font measurement, `regexp` for parsing, SVG output via `svgBuilder`

---

### Task 1: Requirement IR Types

**Files:**
- Create: `ir/requirement.go`
- Create: `ir/requirement_test.go`
- Modify: `ir/graph.go` (add fields after Treemap section, ~line 135)

**Step 1: Create `ir/requirement.go`**

```go
package ir

// RequirementType classifies a requirement definition.
type RequirementType int

const (
	ReqTypeRequirement RequirementType = iota
	ReqTypeFunctional
	ReqTypeInterface
	ReqTypePerformance
	ReqTypePhysical
	ReqTypeDesignConstraint
)

func (t RequirementType) String() string {
	switch t {
	case ReqTypeFunctional:
		return "functionalRequirement"
	case ReqTypeInterface:
		return "interfaceRequirement"
	case ReqTypePerformance:
		return "performanceRequirement"
	case ReqTypePhysical:
		return "physicalRequirement"
	case ReqTypeDesignConstraint:
		return "designConstraint"
	default:
		return "requirement"
	}
}

// Stereotype returns the UML-style stereotype label for display.
func (t RequirementType) Stereotype() string {
	switch t {
	case ReqTypeFunctional:
		return "Functional Requirement"
	case ReqTypeInterface:
		return "Interface Requirement"
	case ReqTypePerformance:
		return "Performance Requirement"
	case ReqTypePhysical:
		return "Physical Requirement"
	case ReqTypeDesignConstraint:
		return "Design Constraint"
	default:
		return "Requirement"
	}
}

// RiskLevel represents a requirement's risk assessment.
type RiskLevel int

const (
	RiskNone RiskLevel = iota
	RiskLow
	RiskMedium
	RiskHigh
)

func (r RiskLevel) String() string {
	switch r {
	case RiskLow:
		return "Low"
	case RiskMedium:
		return "Medium"
	case RiskHigh:
		return "High"
	default:
		return ""
	}
}

// VerifyMethod represents how a requirement is verified.
type VerifyMethod int

const (
	VerifyNone VerifyMethod = iota
	VerifyAnalysis
	VerifyInspection
	VerifyTest
	VerifyDemonstration
)

func (v VerifyMethod) String() string {
	switch v {
	case VerifyAnalysis:
		return "Analysis"
	case VerifyInspection:
		return "Inspection"
	case VerifyTest:
		return "Test"
	case VerifyDemonstration:
		return "Demonstration"
	default:
		return ""
	}
}

// RequirementRelType classifies a relationship between requirements/elements.
type RequirementRelType int

const (
	ReqRelContains RequirementRelType = iota
	ReqRelCopies
	ReqRelDerives
	ReqRelSatisfies
	ReqRelVerifies
	ReqRelRefines
	ReqRelTraces
)

func (r RequirementRelType) String() string {
	switch r {
	case ReqRelContains:
		return "contains"
	case ReqRelCopies:
		return "copies"
	case ReqRelDerives:
		return "derives"
	case ReqRelSatisfies:
		return "satisfies"
	case ReqRelVerifies:
		return "verifies"
	case ReqRelRefines:
		return "refines"
	case ReqRelTraces:
		return "traces"
	default:
		return ""
	}
}

// RequirementDef represents a requirement block.
type RequirementDef struct {
	Name         string
	ID           string
	Text         string
	Type         RequirementType
	Risk         RiskLevel
	VerifyMethod VerifyMethod
}

// ElementDef represents an element block in a requirement diagram.
type ElementDef struct {
	Name   string
	Type   string
	DocRef string
}

// RequirementRel represents a relationship between two nodes.
type RequirementRel struct {
	Source  string
	Target string
	Type   RequirementRelType
}
```

**Step 2: Create `ir/requirement_test.go`**

```go
package ir

import "testing"

func TestRequirementType(t *testing.T) {
	tests := []struct {
		typ        RequirementType
		str        string
		stereotype string
	}{
		{ReqTypeRequirement, "requirement", "Requirement"},
		{ReqTypeFunctional, "functionalRequirement", "Functional Requirement"},
		{ReqTypeInterface, "interfaceRequirement", "Interface Requirement"},
		{ReqTypePerformance, "performanceRequirement", "Performance Requirement"},
		{ReqTypePhysical, "physicalRequirement", "Physical Requirement"},
		{ReqTypeDesignConstraint, "designConstraint", "Design Constraint"},
	}
	for _, tt := range tests {
		if tt.typ.String() != tt.str {
			t.Errorf("RequirementType(%d).String() = %q, want %q", tt.typ, tt.typ.String(), tt.str)
		}
		if tt.typ.Stereotype() != tt.stereotype {
			t.Errorf("RequirementType(%d).Stereotype() = %q, want %q", tt.typ, tt.typ.Stereotype(), tt.stereotype)
		}
	}
}

func TestRiskLevel(t *testing.T) {
	if RiskLow.String() != "Low" {
		t.Errorf("RiskLow = %q", RiskLow.String())
	}
	if RiskMedium.String() != "Medium" {
		t.Errorf("RiskMedium = %q", RiskMedium.String())
	}
	if RiskHigh.String() != "High" {
		t.Errorf("RiskHigh = %q", RiskHigh.String())
	}
	if RiskNone.String() != "" {
		t.Errorf("RiskNone = %q", RiskNone.String())
	}
}

func TestVerifyMethod(t *testing.T) {
	if VerifyTest.String() != "Test" {
		t.Errorf("VerifyTest = %q", VerifyTest.String())
	}
	if VerifyAnalysis.String() != "Analysis" {
		t.Errorf("VerifyAnalysis = %q", VerifyAnalysis.String())
	}
}

func TestRequirementRelType(t *testing.T) {
	if ReqRelSatisfies.String() != "satisfies" {
		t.Errorf("ReqRelSatisfies = %q", ReqRelSatisfies.String())
	}
	if ReqRelContains.String() != "contains" {
		t.Errorf("ReqRelContains = %q", ReqRelContains.String())
	}
}

func TestRequirementGraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Requirement
	g.Requirements = append(g.Requirements, &RequirementDef{
		Name: "test_req",
		ID:   "REQ-001",
		Text: "Must do something",
		Type: ReqTypeFunctional,
		Risk: RiskHigh,
	})
	g.ReqElements = append(g.ReqElements, &ElementDef{
		Name:   "test_element",
		Type:   "Simulation",
		DocRef: "DOC-001",
	})
	g.ReqRelationships = append(g.ReqRelationships, &RequirementRel{
		Source: "test_element",
		Target: "test_req",
		Type:   ReqRelSatisfies,
	})
	if len(g.Requirements) != 1 {
		t.Errorf("Requirements = %d", len(g.Requirements))
	}
	if len(g.ReqElements) != 1 {
		t.Errorf("ReqElements = %d", len(g.ReqElements))
	}
	if len(g.ReqRelationships) != 1 {
		t.Errorf("ReqRelationships = %d", len(g.ReqRelationships))
	}
}
```

**Step 3: Add fields to `ir/graph.go`**

After the Treemap fields block (~line 135), add:

```go
	// Requirement diagram fields
	Requirements     []*RequirementDef
	ReqElements      []*ElementDef
	ReqRelationships []*RequirementRel
```

**Step 4: Run tests**

Run: `go test ./ir/ -v -run TestRequirement`
Expected: All 5 tests PASS

**Step 5: Commit**

```bash
git add ir/requirement.go ir/requirement_test.go ir/graph.go
git commit -m "feat(ir): add Requirement diagram types"
```

---

### Task 2: Block IR Types

**Files:**
- Create: `ir/block.go`
- Create: `ir/block_test.go`
- Modify: `ir/graph.go` (add fields after Requirement section)

**Step 1: Create `ir/block.go`**

```go
package ir

// BlockDef represents a block in a block diagram.
type BlockDef struct {
	ID       string
	Label    string
	Shape    NodeShape
	Width    int        // column span (1 = default)
	Children []*BlockDef // nested blocks
}
```

**Step 2: Create `ir/block_test.go`**

```go
package ir

import "testing"

func TestBlockDef(t *testing.T) {
	b := &BlockDef{
		ID:    "a",
		Label: "Block A",
		Shape: Rectangle,
		Width: 2,
	}
	if b.Width != 2 {
		t.Errorf("Width = %d, want 2", b.Width)
	}
	if b.Shape != Rectangle {
		t.Errorf("Shape = %v, want Rectangle", b.Shape)
	}
}

func TestBlockNesting(t *testing.T) {
	parent := &BlockDef{
		ID:    "parent",
		Label: "Parent",
		Shape: Rectangle,
		Width: 1,
		Children: []*BlockDef{
			{ID: "child1", Label: "Child 1", Shape: Rectangle, Width: 1},
			{ID: "child2", Label: "Child 2", Shape: RoundedRect, Width: 1},
		},
	}
	if len(parent.Children) != 2 {
		t.Fatalf("Children = %d, want 2", len(parent.Children))
	}
	if parent.Children[1].Shape != RoundedRect {
		t.Errorf("child2 shape = %v, want RoundedRect", parent.Children[1].Shape)
	}
}

func TestBlockGraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Block
	g.BlockColumns = 3
	g.Blocks = append(g.Blocks, &BlockDef{ID: "a", Label: "A", Width: 1})
	g.Blocks = append(g.Blocks, &BlockDef{ID: "b", Label: "B", Width: 2})
	if g.BlockColumns != 3 {
		t.Errorf("BlockColumns = %d, want 3", g.BlockColumns)
	}
	if len(g.Blocks) != 2 {
		t.Errorf("Blocks = %d, want 2", len(g.Blocks))
	}
}
```

**Step 3: Add fields to `ir/graph.go`**

After the Requirement fields, add:

```go
	// Block diagram fields
	Blocks       []*BlockDef
	BlockColumns int
```

**Step 4: Run tests**

Run: `go test ./ir/ -v -run TestBlock`
Expected: All 3 tests PASS

**Step 5: Commit**

```bash
git add ir/block.go ir/block_test.go ir/graph.go
git commit -m "feat(ir): add Block diagram types"
```

---

### Task 3: C4 IR Types

**Files:**
- Create: `ir/c4.go`
- Create: `ir/c4_test.go`
- Modify: `ir/graph.go` (add fields after Block section)

**Step 1: Create `ir/c4.go`**

```go
package ir

// C4Kind identifies the specific C4 diagram subtype.
type C4Kind int

const (
	C4Context C4Kind = iota
	C4Container
	C4Component
	C4Dynamic
	C4Deployment
)

func (k C4Kind) String() string {
	switch k {
	case C4Container:
		return "C4Container"
	case C4Component:
		return "C4Component"
	case C4Dynamic:
		return "C4Dynamic"
	case C4Deployment:
		return "C4Deployment"
	default:
		return "C4Context"
	}
}

// C4ElementType identifies the kind of C4 element.
type C4ElementType int

const (
	C4Person C4ElementType = iota
	C4System
	C4SystemDb
	C4SystemQueue
	C4Container_
	C4ContainerDb
	C4ContainerQueue
	C4Component_
	C4ExternalPerson
	C4ExternalSystem
	C4ExternalSystemDb
	C4ExternalSystemQueue
	C4ExternalContainer
	C4ExternalContainerDb
	C4ExternalContainerQueue
	C4ExternalComponent
)

func (e C4ElementType) String() string {
	switch e {
	case C4Person:
		return "Person"
	case C4System:
		return "System"
	case C4SystemDb:
		return "SystemDb"
	case C4SystemQueue:
		return "SystemQueue"
	case C4Container_:
		return "Container"
	case C4ContainerDb:
		return "ContainerDb"
	case C4ContainerQueue:
		return "ContainerQueue"
	case C4Component_:
		return "Component"
	case C4ExternalPerson:
		return "Person_Ext"
	case C4ExternalSystem:
		return "System_Ext"
	case C4ExternalSystemDb:
		return "SystemDb_Ext"
	case C4ExternalSystemQueue:
		return "SystemQueue_Ext"
	case C4ExternalContainer:
		return "Container_Ext"
	case C4ExternalContainerDb:
		return "ContainerDb_Ext"
	case C4ExternalContainerQueue:
		return "ContainerQueue_Ext"
	case C4ExternalComponent:
		return "Component_Ext"
	default:
		return "System"
	}
}

// IsExternal returns true if the element type represents an external entity.
func (e C4ElementType) IsExternal() bool {
	return e >= C4ExternalPerson
}

// IsPerson returns true if the element type is a person (internal or external).
func (e C4ElementType) IsPerson() bool {
	return e == C4Person || e == C4ExternalPerson
}

// IsDatabase returns true if the element type is a database variant.
func (e C4ElementType) IsDatabase() bool {
	switch e {
	case C4SystemDb, C4ContainerDb, C4ExternalSystemDb, C4ExternalContainerDb:
		return true
	default:
		return false
	}
}

// IsQueue returns true if the element type is a queue variant.
func (e C4ElementType) IsQueue() bool {
	switch e {
	case C4SystemQueue, C4ContainerQueue, C4ExternalSystemQueue, C4ExternalContainerQueue:
		return true
	default:
		return false
	}
}

// C4Element represents a single element in a C4 diagram.
type C4Element struct {
	ID          string
	Label       string
	Technology  string
	Description string
	Type        C4ElementType
	BoundaryID  string // empty if top-level
}

// C4Boundary represents a boundary/grouping in a C4 diagram.
type C4Boundary struct {
	ID       string
	Label    string
	Type     string   // e.g. "Enterprise", "Software System"
	Children []string // element IDs within this boundary
}

// C4Rel represents a relationship between C4 elements.
type C4Rel struct {
	From        string
	To          string
	Label       string
	Technology  string
	Description string
}
```

**Step 2: Create `ir/c4_test.go`**

```go
package ir

import "testing"

func TestC4Kind(t *testing.T) {
	tests := []struct {
		kind C4Kind
		str  string
	}{
		{C4Context, "C4Context"},
		{C4Container, "C4Container"},
		{C4Component, "C4Component"},
		{C4Dynamic, "C4Dynamic"},
		{C4Deployment, "C4Deployment"},
	}
	for _, tt := range tests {
		if tt.kind.String() != tt.str {
			t.Errorf("C4Kind(%d).String() = %q, want %q", tt.kind, tt.kind.String(), tt.str)
		}
	}
}

func TestC4ElementType(t *testing.T) {
	if C4Person.String() != "Person" {
		t.Errorf("C4Person = %q", C4Person.String())
	}
	if C4Container_.String() != "Container" {
		t.Errorf("C4Container_ = %q", C4Container_.String())
	}
	if C4ExternalSystem.String() != "System_Ext" {
		t.Errorf("C4ExternalSystem = %q", C4ExternalSystem.String())
	}
}

func TestC4ElementTypePredicates(t *testing.T) {
	if !C4ExternalSystem.IsExternal() {
		t.Error("C4ExternalSystem should be external")
	}
	if C4System.IsExternal() {
		t.Error("C4System should not be external")
	}
	if !C4Person.IsPerson() {
		t.Error("C4Person should be person")
	}
	if !C4SystemDb.IsDatabase() {
		t.Error("C4SystemDb should be database")
	}
	if !C4ContainerQueue.IsQueue() {
		t.Error("C4ContainerQueue should be queue")
	}
}

func TestC4GraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = C4
	g.C4SubKind = C4Container
	g.C4Elements = append(g.C4Elements, &C4Element{
		ID:    "user",
		Label: "User",
		Type:  C4Person,
	})
	g.C4Boundaries = append(g.C4Boundaries, &C4Boundary{
		ID:       "system",
		Label:    "My System",
		Type:     "Software System",
		Children: []string{"webapp"},
	})
	g.C4Rels = append(g.C4Rels, &C4Rel{
		From:  "user",
		To:    "webapp",
		Label: "Uses",
	})
	if g.C4SubKind != C4Container {
		t.Errorf("C4SubKind = %v", g.C4SubKind)
	}
	if len(g.C4Elements) != 1 {
		t.Errorf("C4Elements = %d", len(g.C4Elements))
	}
	if len(g.C4Boundaries) != 1 {
		t.Errorf("C4Boundaries = %d", len(g.C4Boundaries))
	}
	if len(g.C4Rels) != 1 {
		t.Errorf("C4Rels = %d", len(g.C4Rels))
	}
}
```

**Step 3: Add fields to `ir/graph.go`**

After the Block fields, add:

```go
	// C4 diagram fields
	C4SubKind    C4Kind
	C4Elements   []*C4Element
	C4Boundaries []*C4Boundary
	C4Rels       []*C4Rel
```

**Step 4: Run tests**

Run: `go test ./ir/ -v -run TestC4`
Expected: All 4 tests PASS

**Step 5: Commit**

```bash
git add ir/c4.go ir/c4_test.go ir/graph.go
git commit -m "feat(ir): add C4 diagram types"
```

---

### Task 4: Config and Theme

**Files:**
- Modify: `config/config.go` (add 3 config structs + defaults)
- Modify: `config/config_test.go` (add 3 tests)
- Modify: `theme/theme.go` (add color fields to Theme struct + both presets)
- Modify: `theme/theme_test.go` (add 3 tests)

**Step 1: Add config structs to `config/config.go`**

Add to the `Layout` struct (~line 27):
```go
	Requirement RequirementConfig
	Block       BlockConfig
	C4          C4Config
```

Add after TreemapConfig:
```go
// RequirementConfig holds requirement diagram layout options.
type RequirementConfig struct {
	NodeMinWidth     float32
	NodePadding      float32
	MetadataFontSize float32
	PaddingX         float32
	PaddingY         float32
}

// BlockConfig holds block diagram layout options.
type BlockConfig struct {
	ColumnGap   float32
	RowGap      float32
	NodePadding float32
	PaddingX    float32
	PaddingY    float32
}

// C4Config holds C4 diagram layout options.
type C4Config struct {
	PersonWidth     float32
	PersonHeight    float32
	SystemWidth     float32
	SystemHeight    float32
	BoundaryPadding float32
	PaddingX        float32
	PaddingY        float32
}
```

Add defaults in `DefaultLayout()`:
```go
		Requirement: RequirementConfig{
			NodeMinWidth:     180,
			NodePadding:      12,
			MetadataFontSize: 11,
			PaddingX:         10,
			PaddingY:         10,
		},
		Block: BlockConfig{
			ColumnGap:   20,
			RowGap:      20,
			NodePadding: 12,
			PaddingX:    20,
			PaddingY:    20,
		},
		C4: C4Config{
			PersonWidth:     160,
			PersonHeight:    180,
			SystemWidth:     200,
			SystemHeight:    120,
			BoundaryPadding: 20,
			PaddingX:        20,
			PaddingY:        20,
		},
```

**Step 2: Add theme fields to `theme/theme.go`**

Add to Theme struct:
```go
	// Requirement diagram colors
	RequirementFill   string
	RequirementBorder string

	// Block diagram colors
	BlockColors     []string
	BlockNodeBorder string

	// C4 diagram colors
	C4PersonColor    string
	C4SystemColor    string
	C4ContainerColor string
	C4ComponentColor string
	C4ExternalColor  string
	C4BoundaryColor  string
	C4TextColor      string
```

Add to `Modern()`:
```go
		RequirementFill:   "#F0F4F8",
		RequirementBorder: "#3B6492",

		BlockColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		BlockNodeBorder: "#3B6492",

		C4PersonColor:    "#08427B",
		C4SystemColor:    "#1168BD",
		C4ContainerColor: "#438DD5",
		C4ComponentColor: "#85BBF0",
		C4ExternalColor:  "#999999",
		C4BoundaryColor:  "#444444",
		C4TextColor:      "#FFFFFF",
```

Add to `MermaidDefault()`:
```go
		RequirementFill:   "#ECECFF",
		RequirementBorder: "#9370DB",

		BlockColors: []string{
			"#9370DB", "#E76F51", "#7FB069", "#F4A261",
			"#48A9A6", "#D08AC0", "#E4E36A", "#F7B7A3",
		},
		BlockNodeBorder: "#9370DB",

		C4PersonColor:    "#08427B",
		C4SystemColor:    "#1168BD",
		C4ContainerColor: "#438DD5",
		C4ComponentColor: "#85BBF0",
		C4ExternalColor:  "#999999",
		C4BoundaryColor:  "#444444",
		C4TextColor:      "#FFFFFF",
```

**Step 3: Add tests to `config/config_test.go` and `theme/theme_test.go`**

Config tests:
```go
func TestDefaultLayoutRequirementConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Requirement.NodeMinWidth <= 0 {
		t.Error("NodeMinWidth should be positive")
	}
}

func TestDefaultLayoutBlockConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Block.ColumnGap <= 0 {
		t.Error("ColumnGap should be positive")
	}
}

func TestDefaultLayoutC4Config(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.C4.PersonWidth <= 0 {
		t.Error("PersonWidth should be positive")
	}
}
```

Theme tests:
```go
func TestModernRequirementColors(t *testing.T) {
	th := Modern()
	if th.RequirementFill == "" {
		t.Error("RequirementFill empty")
	}
}

func TestModernBlockColors(t *testing.T) {
	th := Modern()
	if len(th.BlockColors) == 0 {
		t.Error("BlockColors empty")
	}
}

func TestModernC4Colors(t *testing.T) {
	th := Modern()
	if th.C4PersonColor == "" {
		t.Error("C4PersonColor empty")
	}
	if th.C4SystemColor == "" {
		t.Error("C4SystemColor empty")
	}
}
```

**Step 4: Run tests**

Run: `go test ./config/ ./theme/ -v -run "Requirement|Block|C4"`
Expected: All 6 new tests PASS

**Step 5: Commit**

```bash
git add config/config.go config/config_test.go theme/theme.go theme/theme_test.go
git commit -m "feat(config,theme): add Requirement, Block, and C4 config and theme"
```

---

### Task 5: Requirement Parser

**Files:**
- Create: `parser/requirement.go`
- Create: `parser/requirement_test.go`
- Modify: `parser/parser.go` (add dispatch case)

**Step 1: Create `parser/requirement.go`**

```go
package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	reqBlockStartRe = regexp.MustCompile(`^(requirement|functionalRequirement|interfaceRequirement|performanceRequirement|physicalRequirement|designConstraint)\s+(\w+)\s*\{?\s*$`)
	elemBlockStartRe = regexp.MustCompile(`^element\s+(\w+)\s*\{?\s*$`)
	reqFieldRe       = regexp.MustCompile(`^\s*(\w+)\s*:\s*(.+?)\s*$`)
	reqRelRe         = regexp.MustCompile(`^(\w+)\s+-\s+(contains|copies|derives|satisfies|verifies|refines|traces)\s+->\s+(\w+)\s*$`)
)

func parseRequirement(input string) (*ParseOutput, error) {
	lines := preprocessInput(input)
	g := ir.NewGraph()
	g.Kind = ir.Requirement

	if len(lines) > 0 {
		lower := strings.ToLower(lines[0])
		if strings.HasPrefix(lower, "requirementdiagram") {
			lines = lines[1:]
		}
	}

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Try requirement/element block
		if m := reqBlockStartRe.FindStringSubmatch(line); m != nil {
			reqType := parseReqType(m[1])
			name := m[2]
			i++
			req := &ir.RequirementDef{Name: name, Type: reqType}
			for i < len(lines) && lines[i] != "}" {
				if fm := reqFieldRe.FindStringSubmatch(lines[i]); fm != nil {
					switch strings.ToLower(fm[1]) {
					case "id":
						req.ID = fm[2]
					case "text":
						req.Text = fm[2]
					case "risk":
						req.Risk = parseRiskLevel(fm[2])
					case "verifymethod":
						req.VerifyMethod = parseVerifyMethod(fm[2])
					}
				}
				i++
			}
			g.Requirements = append(g.Requirements, req)
			// Create a graph node for layout
			label := name
			g.EnsureNode(name, &label, nil)
			i++ // skip closing brace
			continue
		}

		if m := elemBlockStartRe.FindStringSubmatch(line); m != nil {
			name := m[1]
			i++
			elem := &ir.ElementDef{Name: name}
			for i < len(lines) && lines[i] != "}" {
				if fm := reqFieldRe.FindStringSubmatch(lines[i]); fm != nil {
					switch strings.ToLower(fm[1]) {
					case "type":
						elem.Type = fm[2]
					case "docref":
						elem.DocRef = fm[2]
					}
				}
				i++
			}
			g.ReqElements = append(g.ReqElements, elem)
			label := name
			g.EnsureNode(name, &label, nil)
			i++ // skip closing brace
			continue
		}

		// Try relationship
		if m := reqRelRe.FindStringSubmatch(line); m != nil {
			rel := &ir.RequirementRel{
				Source: m[1],
				Target: m[3],
				Type:   parseRelType(m[2]),
			}
			g.ReqRelationships = append(g.ReqRelationships, rel)
			// Add edge for layout
			relLabel := rel.Type.String()
			g.Edges = append(g.Edges, &ir.Edge{
				From:     rel.Source,
				To:       rel.Target,
				Label:    &relLabel,
				Directed: true,
				ArrowEnd: true,
			})
			i++
			continue
		}

		i++
	}

	return &ParseOutput{Graph: g}, nil
}

func parseReqType(s string) ir.RequirementType {
	switch strings.ToLower(s) {
	case "functionalrequirement":
		return ir.ReqTypeFunctional
	case "interfacerequirement":
		return ir.ReqTypeInterface
	case "performancerequirement":
		return ir.ReqTypePerformance
	case "physicalrequirement":
		return ir.ReqTypePhysical
	case "designconstraint":
		return ir.ReqTypeDesignConstraint
	default:
		return ir.ReqTypeRequirement
	}
}

func parseRiskLevel(s string) ir.RiskLevel {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "low":
		return ir.RiskLow
	case "medium":
		return ir.RiskMedium
	case "high":
		return ir.RiskHigh
	default:
		return ir.RiskNone
	}
}

func parseVerifyMethod(s string) ir.VerifyMethod {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "analysis":
		return ir.VerifyAnalysis
	case "inspection":
		return ir.VerifyInspection
	case "test":
		return ir.VerifyTest
	case "demonstration":
		return ir.VerifyDemonstration
	default:
		return ir.VerifyNone
	}
}

func parseRelType(s string) ir.RequirementRelType {
	switch strings.ToLower(s) {
	case "contains":
		return ir.ReqRelContains
	case "copies":
		return ir.ReqRelCopies
	case "derives":
		return ir.ReqRelDerives
	case "satisfies":
		return ir.ReqRelSatisfies
	case "verifies":
		return ir.ReqRelVerifies
	case "refines":
		return ir.ReqRelRefines
	case "traces":
		return ir.ReqRelTraces
	default:
		return ir.ReqRelContains
	}
}
```

**Step 2: Create `parser/requirement_test.go`**

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseRequirementBasic(t *testing.T) {
	input := `requirementDiagram

requirement test_req {
id: 1
text: the test text.
risk: high
verifymethod: test
}

element test_entity {
type: simulation
docref: DOC-001
}

test_entity - satisfies -> test_req`

	out, err := parseRequirement(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Requirement {
		t.Fatalf("Kind = %v, want Requirement", g.Kind)
	}
	if len(g.Requirements) != 1 {
		t.Fatalf("Requirements = %d, want 1", len(g.Requirements))
	}
	req := g.Requirements[0]
	if req.Name != "test_req" {
		t.Errorf("req name = %q", req.Name)
	}
	if req.ID != "1" {
		t.Errorf("req id = %q", req.ID)
	}
	if req.Risk != ir.RiskHigh {
		t.Errorf("req risk = %v", req.Risk)
	}
	if req.VerifyMethod != ir.VerifyTest {
		t.Errorf("req verify = %v", req.VerifyMethod)
	}
	if len(g.ReqElements) != 1 {
		t.Fatalf("Elements = %d, want 1", len(g.ReqElements))
	}
	elem := g.ReqElements[0]
	if elem.Type != "simulation" {
		t.Errorf("elem type = %q", elem.Type)
	}
	if len(g.ReqRelationships) != 1 {
		t.Fatalf("Rels = %d, want 1", len(g.ReqRelationships))
	}
	if g.ReqRelationships[0].Type != ir.ReqRelSatisfies {
		t.Errorf("rel type = %v", g.ReqRelationships[0].Type)
	}
}

func TestParseRequirementMultiple(t *testing.T) {
	input := `requirementDiagram

functionalRequirement req1 {
id: FR-001
text: Must authenticate users
risk: medium
verifymethod: demonstration
}

requirement req2 {
id: REQ-002
text: Must log events
risk: low
verifymethod: analysis
}

req1 - derives -> req2`

	out, err := parseRequirement(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Requirements) != 2 {
		t.Fatalf("Requirements = %d, want 2", len(out.Graph.Requirements))
	}
	if out.Graph.Requirements[0].Type != ir.ReqTypeFunctional {
		t.Errorf("req1 type = %v, want Functional", out.Graph.Requirements[0].Type)
	}
}

func TestParseRequirementEmpty(t *testing.T) {
	input := `requirementDiagram`
	out, err := parseRequirement(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Requirements) != 0 {
		t.Errorf("Requirements = %d, want 0", len(out.Graph.Requirements))
	}
}
```

**Step 3: Add dispatch case to `parser/parser.go`**

Add after the Treemap case:
```go
	case ir.Requirement:
		return parseRequirement(input)
```

**Step 4: Run tests**

Run: `go test ./parser/ -v -run TestParseRequirement`
Expected: All 3 tests PASS

**Step 5: Commit**

```bash
git add parser/requirement.go parser/requirement_test.go parser/parser.go
git commit -m "feat(parser): add Requirement diagram parser"
```

---

### Task 6: Block Parser

**Files:**
- Create: `parser/block.go`
- Create: `parser/block_test.go`
- Modify: `parser/parser.go` (add dispatch case)

**Step 1: Create `parser/block.go`**

Parse `block-beta` header, `columns N` directive, block definitions with optional shapes (reuse flowchart shape syntax), width spans (`:N`), edges (`-->`, `---`), and nesting via indentation.

```go
package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	blockColumnsRe = regexp.MustCompile(`^columns\s+(\d+)$`)
	blockEdgeRe    = regexp.MustCompile(`^(\w+)\s*(-->|---)\s*(?:\|"?([^"|]*)"?\|\s*)?(\w+)$`)
	blockDefRe     = regexp.MustCompile(`^(\w+)(?:\["([^"]*)"\]|\("([^"]*)"\)|\(\("([^"]*)"\)\)|\{"([^"]*)"\}|>\["([^"]*)"\])?(?::(\d+))?\s*$`)
)

func parseBlock(input string) (*ParseOutput, error) {
	lines := preprocessInput(input)
	g := ir.NewGraph()
	g.Kind = ir.Block

	if len(lines) > 0 {
		lower := strings.ToLower(lines[0])
		if strings.HasPrefix(lower, "block") {
			lines = lines[1:]
		}
	}

	for _, line := range lines {
		// columns directive
		if m := blockColumnsRe.FindStringSubmatch(line); m != nil {
			cols, _ := strconv.Atoi(m[1])
			g.BlockColumns = cols
			continue
		}

		// edge
		if m := blockEdgeRe.FindStringSubmatch(line); m != nil {
			from, arrow, label, to := m[1], m[2], m[3], m[4]
			edge := &ir.Edge{
				From:     from,
				To:       to,
				Directed: arrow == "-->",
				ArrowEnd: arrow == "-->",
			}
			if label != "" {
				edge.Label = &label
			}
			g.Edges = append(g.Edges, edge)
			// Ensure nodes exist
			g.EnsureNode(from, nil, nil)
			g.EnsureNode(to, nil, nil)
			continue
		}

		// block definition(s) — may have multiple on one line
		parseBlockDefs(line, g)
	}

	return &ParseOutput{Graph: g}, nil
}

func parseBlockDefs(line string, g *ir.Graph) {
	// Split line by whitespace to handle multiple blocks on one line
	tokens := strings.Fields(line)
	for _, token := range tokens {
		m := blockDefRe.FindStringSubmatch(token)
		if m == nil {
			continue
		}
		id := m[1]
		// Determine label and shape from capture groups
		label := id
		shape := ir.Rectangle
		switch {
		case m[2] != "": // ["label"] = square
			label = m[2]
			shape = ir.Rectangle
		case m[3] != "": // ("label") = rounded
			label = m[3]
			shape = ir.RoundedRect
		case m[4] != "": // (("label")) = circle
			label = m[4]
			shape = ir.Circle
		case m[5] != "": // {"label"} = diamond
			label = m[5]
			shape = ir.Diamond
		case m[6] != "": // >["label"] = asymmetric
			label = m[6]
			shape = ir.Asymmetric
		}

		width := 1
		if m[7] != "" {
			width, _ = strconv.Atoi(m[7])
		}

		block := &ir.BlockDef{
			ID:    id,
			Label: label,
			Shape: shape,
			Width: width,
		}
		g.Blocks = append(g.Blocks, block)
		g.EnsureNode(id, &label, &shape)
	}
}
```

**Step 2: Create `parser/block_test.go`**

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseBlockBasic(t *testing.T) {
	input := `block-beta
columns 3
a["A"] b["B"] c["C"]
d["D"]:2 e["E"]`

	out, err := parseBlock(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Block {
		t.Fatalf("Kind = %v, want Block", g.Kind)
	}
	if g.BlockColumns != 3 {
		t.Fatalf("BlockColumns = %d, want 3", g.BlockColumns)
	}
	if len(g.Blocks) != 5 {
		t.Fatalf("Blocks = %d, want 5", len(g.Blocks))
	}
	// d spans 2 columns
	if g.Blocks[3].Width != 2 {
		t.Errorf("d width = %d, want 2", g.Blocks[3].Width)
	}
}

func TestParseBlockEdges(t *testing.T) {
	input := `block-beta
a["Source"] b["Target"]
a --> b`

	out, err := parseBlock(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Edges) != 1 {
		t.Fatalf("Edges = %d, want 1", len(out.Graph.Edges))
	}
	e := out.Graph.Edges[0]
	if e.From != "a" || e.To != "b" {
		t.Errorf("Edge = %s -> %s", e.From, e.To)
	}
}

func TestParseBlockShapes(t *testing.T) {
	input := `block-beta
a["Square"] b("Rounded") c(("Circle")) d{"Diamond"}`

	out, err := parseBlock(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Blocks) != 4 {
		t.Fatalf("Blocks = %d, want 4", len(out.Graph.Blocks))
	}
	shapes := []ir.NodeShape{ir.Rectangle, ir.RoundedRect, ir.Circle, ir.Diamond}
	for i, want := range shapes {
		if out.Graph.Blocks[i].Shape != want {
			t.Errorf("block[%d] shape = %v, want %v", i, out.Graph.Blocks[i].Shape, want)
		}
	}
}

func TestParseBlockEmpty(t *testing.T) {
	input := `block-beta`
	out, err := parseBlock(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.Blocks) != 0 {
		t.Errorf("Blocks = %d, want 0", len(out.Graph.Blocks))
	}
}
```

**Step 3: Add dispatch case to `parser/parser.go`**

```go
	case ir.Block:
		return parseBlock(input)
```

**Step 4: Run tests**

Run: `go test ./parser/ -v -run TestParseBlock`
Expected: All 4 tests PASS

**Step 5: Commit**

```bash
git add parser/block.go parser/block_test.go parser/parser.go
git commit -m "feat(parser): add Block diagram parser"
```

---

### Task 7: C4 Parser

**Files:**
- Create: `parser/c4.go`
- Create: `parser/c4_test.go`
- Modify: `parser/parser.go` (add dispatch case)

**Step 1: Create `parser/c4.go`**

Parse function-call syntax: `Person(id, "label", "description")`, `Container_Boundary(id, "label") {`, `Rel(from, to, "label", "tech")`, boundary nesting with braces.

```go
package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	c4ElementRe  = regexp.MustCompile(`^(Person|Person_Ext|System|System_Ext|SystemDb|SystemDb_Ext|SystemQueue|SystemQueue_Ext|Container|Container_Ext|ContainerDb|ContainerDb_Ext|ContainerQueue|ContainerQueue_Ext|Component|Component_Ext)\s*\((.+)\)\s*$`)
	c4BoundaryRe = regexp.MustCompile(`^(Enterprise_Boundary|System_Boundary|Container_Boundary|Boundary)\s*\(([^,]+),\s*"([^"]*)"(?:\s*,\s*"([^"]*)")?\)\s*\{?\s*$`)
	c4RelRe      = regexp.MustCompile(`^(Rel|Rel_Back|Rel_Neighbor|Rel_Back_Neighbor|BiRel|BiRel_Neighbor)\s*\((.+)\)\s*$`)
)

func parseC4(input string) (*ParseOutput, error) {
	lines := preprocessInput(input)
	g := ir.NewGraph()
	g.Kind = ir.C4

	if len(lines) > 0 {
		g.C4SubKind = parseC4Kind(lines[0])
		lines = lines[1:]
	}

	var boundaryStack []*ir.C4Boundary

	for _, line := range lines {
		// Closing brace — pop boundary
		if strings.TrimSpace(line) == "}" {
			if len(boundaryStack) > 0 {
				boundaryStack = boundaryStack[:len(boundaryStack)-1]
			}
			continue
		}

		// Boundary
		if m := c4BoundaryRe.FindStringSubmatch(line); m != nil {
			boundaryType := m[1]
			id := strings.TrimSpace(m[2])
			label := m[3]
			bType := ""
			switch boundaryType {
			case "Enterprise_Boundary":
				bType = "Enterprise"
			case "System_Boundary":
				bType = "Software System"
			case "Container_Boundary":
				bType = "Container"
			default:
				bType = m[4] // Boundary(id, "label", "type")
			}

			boundary := &ir.C4Boundary{
				ID:    id,
				Label: label,
				Type:  bType,
			}
			g.C4Boundaries = append(g.C4Boundaries, boundary)
			boundaryStack = append(boundaryStack, boundary)
			continue
		}

		// Element
		if m := c4ElementRe.FindStringSubmatch(line); m != nil {
			elemType := parseC4ElementType(m[1])
			args := parseC4Args(m[2])
			if len(args) < 2 {
				continue
			}
			elem := &ir.C4Element{
				ID:    args[0],
				Label: args[1],
				Type:  elemType,
			}
			if len(args) > 2 {
				// For Person: args[2] is description
				// For Container/Component: args[2] is technology, args[3] is description
				if elemType.IsPerson() {
					elem.Description = args[2]
				} else {
					elem.Technology = args[2]
					if len(args) > 3 {
						elem.Description = args[3]
					}
				}
			}

			// Assign to current boundary
			if len(boundaryStack) > 0 {
				parent := boundaryStack[len(boundaryStack)-1]
				elem.BoundaryID = parent.ID
				parent.Children = append(parent.Children, elem.ID)
			}

			g.C4Elements = append(g.C4Elements, elem)
			label := elem.Label
			g.EnsureNode(elem.ID, &label, nil)
			continue
		}

		// Relationship
		if m := c4RelRe.FindStringSubmatch(line); m != nil {
			args := parseC4Args(m[2])
			if len(args) < 3 {
				continue
			}
			rel := &ir.C4Rel{
				From:  args[0],
				To:    args[1],
				Label: args[2],
			}
			if len(args) > 3 {
				rel.Technology = args[3]
			}
			g.C4Rels = append(g.C4Rels, rel)
			relLabel := rel.Label
			g.Edges = append(g.Edges, &ir.Edge{
				From:     rel.From,
				To:       rel.To,
				Label:    &relLabel,
				Directed: true,
				ArrowEnd: true,
			})
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}

func parseC4Kind(line string) ir.C4Kind {
	lower := strings.ToLower(strings.TrimSpace(line))
	switch {
	case strings.HasPrefix(lower, "c4container"):
		return ir.C4Container
	case strings.HasPrefix(lower, "c4component"):
		return ir.C4Component
	case strings.HasPrefix(lower, "c4dynamic"):
		return ir.C4Dynamic
	case strings.HasPrefix(lower, "c4deployment"):
		return ir.C4Deployment
	default:
		return ir.C4Context
	}
}

func parseC4ElementType(s string) ir.C4ElementType {
	switch s {
	case "Person":
		return ir.C4Person
	case "Person_Ext":
		return ir.C4ExternalPerson
	case "System":
		return ir.C4System
	case "System_Ext":
		return ir.C4ExternalSystem
	case "SystemDb":
		return ir.C4SystemDb
	case "SystemDb_Ext":
		return ir.C4ExternalSystemDb
	case "SystemQueue":
		return ir.C4SystemQueue
	case "SystemQueue_Ext":
		return ir.C4ExternalSystemQueue
	case "Container":
		return ir.C4Container_
	case "Container_Ext":
		return ir.C4ExternalContainer
	case "ContainerDb":
		return ir.C4ContainerDb
	case "ContainerDb_Ext":
		return ir.C4ExternalContainerDb
	case "ContainerQueue":
		return ir.C4ContainerQueue
	case "ContainerQueue_Ext":
		return ir.C4ExternalContainerQueue
	case "Component":
		return ir.C4Component_
	case "Component_Ext":
		return ir.C4ExternalComponent
	default:
		return ir.C4System
	}
}

// parseC4Args splits a C4 argument list, respecting quoted strings.
func parseC4Args(s string) []string {
	var args []string
	var current strings.Builder
	inQuote := false

	for _, ch := range s {
		switch {
		case ch == '"':
			inQuote = !inQuote
		case ch == ',' && !inQuote:
			args = append(args, strings.TrimSpace(current.String()))
			current.Reset()
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, strings.TrimSpace(current.String()))
	}
	return args
}
```

**Step 2: Create `parser/c4_test.go`**

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseC4Context(t *testing.T) {
	input := `C4Context
Person(user, "User", "A user of the system")
System(webapp, "Web App", "The main web application")
Rel(user, webapp, "Uses", "HTTPS")`

	out, err := parseC4(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.C4 {
		t.Fatalf("Kind = %v, want C4", g.Kind)
	}
	if g.C4SubKind != ir.C4Context {
		t.Fatalf("C4SubKind = %v, want C4Context", g.C4SubKind)
	}
	if len(g.C4Elements) != 2 {
		t.Fatalf("Elements = %d, want 2", len(g.C4Elements))
	}
	if g.C4Elements[0].Type != ir.C4Person {
		t.Errorf("elem[0] type = %v, want Person", g.C4Elements[0].Type)
	}
	if g.C4Elements[0].Description != "A user of the system" {
		t.Errorf("elem[0] desc = %q", g.C4Elements[0].Description)
	}
	if len(g.C4Rels) != 1 {
		t.Fatalf("Rels = %d, want 1", len(g.C4Rels))
	}
	if g.C4Rels[0].Technology != "HTTPS" {
		t.Errorf("rel tech = %q", g.C4Rels[0].Technology)
	}
}

func TestParseC4Container(t *testing.T) {
	input := `C4Container
Person(user, "User", "End user")
Container_Boundary(system, "My System") {
Container(api, "API", "Go", "REST API")
ContainerDb(db, "Database", "PostgreSQL", "Stores data")
}
Rel(user, api, "Calls")
Rel(api, db, "Reads/Writes")`

	out, err := parseC4(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.C4SubKind != ir.C4Container {
		t.Fatalf("C4SubKind = %v, want C4Container", g.C4SubKind)
	}
	if len(g.C4Boundaries) != 1 {
		t.Fatalf("Boundaries = %d, want 1", len(g.C4Boundaries))
	}
	b := g.C4Boundaries[0]
	if b.Label != "My System" {
		t.Errorf("boundary label = %q", b.Label)
	}
	if len(b.Children) != 2 {
		t.Fatalf("boundary children = %d, want 2", len(b.Children))
	}
	// Check container has technology
	var api *ir.C4Element
	for _, e := range g.C4Elements {
		if e.ID == "api" {
			api = e
		}
	}
	if api == nil {
		t.Fatal("api element not found")
	}
	if api.Technology != "Go" {
		t.Errorf("api tech = %q, want Go", api.Technology)
	}
	if api.BoundaryID != "system" {
		t.Errorf("api boundary = %q, want system", api.BoundaryID)
	}
}

func TestParseC4Empty(t *testing.T) {
	input := `C4Context`
	out, err := parseC4(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.C4Elements) != 0 {
		t.Errorf("Elements = %d, want 0", len(out.Graph.C4Elements))
	}
}

func TestParseC4Args(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{`user, "User", "Description"`, 3},
		{`api, "API", "Go", "REST API"`, 4},
		{`a, b`, 2},
	}
	for _, tt := range tests {
		args := parseC4Args(tt.input)
		if len(args) != tt.want {
			t.Errorf("parseC4Args(%q) = %d args, want %d", tt.input, len(args), tt.want)
		}
	}
}
```

**Step 3: Add dispatch case to `parser/parser.go`**

```go
	case ir.C4:
		return parseC4(input)
```

**Step 4: Run tests**

Run: `go test ./parser/ -v -run TestParseC4`
Expected: All 4 tests PASS

**Step 5: Commit**

```bash
git add parser/c4.go parser/c4_test.go parser/parser.go
git commit -m "feat(parser): add C4 diagram parser"
```

---

### Task 8: Requirement Layout

**Files:**
- Create: `layout/requirement.go`
- Create: `layout/requirement_test.go`
- Modify: `layout/types.go` (add RequirementData)
- Modify: `layout/layout.go` (add dispatch case)

**Step 1: Add types to `layout/types.go`**

```go
// RequirementData holds requirement-diagram-specific layout data.
type RequirementData struct {
	Requirements map[string]*ir.RequirementDef
	Elements     map[string]*ir.ElementDef
	NodeKinds    map[string]string // node ID -> "requirement" or "element"
}

func (RequirementData) diagramData() {}
```

**Step 2: Create `layout/requirement.go`**

```go
package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeRequirementLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	nodes := sizeRequirementNodes(g, measurer, th, cfg)

	r := runSugiyama(g, nodes, cfg)

	reqMap := make(map[string]*ir.RequirementDef)
	for _, req := range g.Requirements {
		reqMap[req.Name] = req
	}
	elemMap := make(map[string]*ir.ElementDef)
	for _, elem := range g.ReqElements {
		elemMap[elem.Name] = elem
	}
	nodeKinds := make(map[string]string)
	for _, req := range g.Requirements {
		nodeKinds[req.Name] = "requirement"
	}
	for _, elem := range g.ReqElements {
		nodeKinds[elem.Name] = "element"
	}

	return &Layout{
		Kind:   g.Kind,
		Nodes:  nodes,
		Edges:  r.Edges,
		Width:  r.Width,
		Height: r.Height,
		Diagram: RequirementData{
			Requirements: reqMap,
			Elements:     elemMap,
			NodeKinds:    nodeKinds,
		},
	}
}

func sizeRequirementNodes(g *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(g.Nodes))
	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical
	lineH := th.FontSize * cfg.LabelLineHeight
	metaFontSize := cfg.Requirement.MetadataFontSize
	metaLineH := metaFontSize * cfg.LabelLineHeight
	minW := cfg.Requirement.NodeMinWidth

	reqMap := make(map[string]*ir.RequirementDef)
	for _, req := range g.Requirements {
		reqMap[req.Name] = req
	}
	elemMap := make(map[string]*ir.ElementDef)
	for _, elem := range g.ReqElements {
		elemMap[elem.Name] = elem
	}

	for id, node := range g.Nodes {
		var maxW float32
		var totalH float32

		// Stereotype line
		stereotypeText := ""
		if req, ok := reqMap[id]; ok {
			stereotypeText = "\u00AB" + req.Type.Stereotype() + "\u00BB"
		} else {
			stereotypeText = "\u00ABelement\u00BB"
		}
		stW := measurer.Width(stereotypeText, metaFontSize, th.FontFamily)
		if stW > maxW {
			maxW = stW
		}
		totalH += lineH // stereotype

		// Name line
		nameW := measurer.Width(node.Label, th.FontSize, th.FontFamily)
		if nameW > maxW {
			maxW = nameW
		}
		totalH += lineH // name

		// Metadata lines
		if req, ok := reqMap[id]; ok {
			lines := 0
			if req.ID != "" {
				w := measurer.Width("Id: "+req.ID, metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			if req.Text != "" {
				w := measurer.Width("Text: "+req.Text, metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			if req.Risk != ir.RiskNone {
				w := measurer.Width("Risk: "+req.Risk.String(), metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			if req.VerifyMethod != ir.VerifyNone {
				w := measurer.Width("Verify: "+req.VerifyMethod.String(), metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			totalH += metaLineH * float32(lines)
		} else if elem, ok := elemMap[id]; ok {
			lines := 0
			if elem.Type != "" {
				w := measurer.Width("Type: "+elem.Type, metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			if elem.DocRef != "" {
				w := measurer.Width("Doc: "+elem.DocRef, metaFontSize, th.FontFamily)
				if w > maxW {
					maxW = w
				}
				lines++
			}
			totalH += metaLineH * float32(lines)
		}

		w := maxW + 2*padH
		if w < minW {
			w = minW
		}
		h := totalH + 2*padV

		nodes[id] = &NodeLayout{
			ID:     id,
			Label:  TextBlock{Lines: []string{node.Label}, Width: maxW, Height: lineH, FontSize: th.FontSize},
			Shape:  ir.Rectangle,
			Width:  w,
			Height: h,
		}
	}

	return nodes
}
```

**Step 3: Create `layout/requirement_test.go`**

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRequirementLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Requirement
	g.Direction = ir.TopDown

	reqLabel := "test_req"
	elemLabel := "test_elem"
	g.EnsureNode("test_req", &reqLabel, nil)
	g.EnsureNode("test_elem", &elemLabel, nil)

	g.Requirements = append(g.Requirements, &ir.RequirementDef{
		Name: "test_req", ID: "REQ-001", Text: "Must work", Risk: ir.RiskHigh, VerifyMethod: ir.VerifyTest,
	})
	g.ReqElements = append(g.ReqElements, &ir.ElementDef{
		Name: "test_elem", Type: "Simulation",
	})

	relLabel := "satisfies"
	g.Edges = append(g.Edges, &ir.Edge{From: "test_elem", To: "test_req", Label: &relLabel, Directed: true, ArrowEnd: true})
	g.ReqRelationships = append(g.ReqRelationships, &ir.RequirementRel{Source: "test_elem", Target: "test_req", Type: ir.ReqRelSatisfies})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Requirement {
		t.Fatalf("Kind = %v", l.Kind)
	}
	rd, ok := l.Diagram.(RequirementData)
	if !ok {
		t.Fatal("Diagram is not RequirementData")
	}
	if len(rd.Requirements) != 1 {
		t.Errorf("Requirements = %d", len(rd.Requirements))
	}
	if len(rd.Elements) != 1 {
		t.Errorf("Elements = %d", len(rd.Elements))
	}
	if len(l.Nodes) != 2 {
		t.Errorf("Nodes = %d", len(l.Nodes))
	}
	if len(l.Edges) != 1 {
		t.Errorf("Edges = %d", len(l.Edges))
	}
}

func TestRequirementLayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Requirement
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)
	if len(l.Nodes) != 0 {
		t.Errorf("Nodes = %d", len(l.Nodes))
	}
}
```

**Step 4: Add dispatch case to `layout/layout.go`**

```go
	case ir.Requirement:
		return computeRequirementLayout(g, th, cfg)
```

**Step 5: Run tests**

Run: `go test ./layout/ -v -run TestRequirement`
Expected: Both tests PASS

**Step 6: Commit**

```bash
git add layout/requirement.go layout/requirement_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Requirement diagram layout"
```

---

### Task 9: Block Layout (Hybrid Grid/Sugiyama)

**Files:**
- Create: `layout/block.go`
- Create: `layout/block_test.go`
- Modify: `layout/types.go` (add BlockData)
- Modify: `layout/layout.go` (add dispatch case)

**Step 1: Add types to `layout/types.go`**

```go
// BlockData holds block-diagram-specific layout data.
type BlockData struct {
	Columns    int
	BlockInfos map[string]BlockInfo
}

func (BlockData) diagramData() {}

// BlockInfo stores per-block layout metadata.
type BlockInfo struct {
	Span     int
	HasChildren bool
}
```

**Step 2: Create `layout/block.go`**

Grid layout when columns > 0, Sugiyama fallback when connections exist without columns.

```go
package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeBlockLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	nodes := sizeBlockNodes(g, measurer, th, cfg)

	blockInfos := make(map[string]BlockInfo)
	for _, b := range g.Blocks {
		blockInfos[b.ID] = BlockInfo{Span: b.Width, HasChildren: len(b.Children) > 0}
	}

	// Decide layout strategy
	if g.BlockColumns > 0 {
		return blockGridLayout(g, nodes, blockInfos, th, cfg)
	}
	if len(g.Edges) > 0 {
		// Sugiyama fallback for connected graphs
		r := runSugiyama(g, nodes, cfg)
		return &Layout{
			Kind:    g.Kind,
			Nodes:   nodes,
			Edges:   r.Edges,
			Width:   r.Width,
			Height:  r.Height,
			Diagram: BlockData{Columns: 0, BlockInfos: blockInfos},
		}
	}
	// No columns, no edges — single column stack
	return blockGridLayout(g, nodes, blockInfos, th, cfg)
}

func blockGridLayout(g *ir.Graph, nodes map[string]*NodeLayout, blockInfos map[string]BlockInfo, th *theme.Theme, cfg *config.Layout) *Layout {
	cols := g.BlockColumns
	if cols <= 0 {
		cols = 1
	}
	colGap := cfg.Block.ColumnGap
	rowGap := cfg.Block.RowGap
	padX := cfg.Block.PaddingX
	padY := cfg.Block.PaddingY

	// Measure max cell dimensions
	var maxCellW, maxCellH float32
	for _, n := range nodes {
		if n.Width > maxCellW {
			maxCellW = n.Width
		}
		if n.Height > maxCellH {
			maxCellH = n.Height
		}
	}

	// Place blocks in grid order (following g.Blocks order)
	col := 0
	row := 0
	for _, blk := range g.Blocks {
		n, ok := nodes[blk.ID]
		if !ok {
			continue
		}

		span := blk.Width
		if span <= 0 {
			span = 1
		}
		// Wrap if this block doesn't fit in current row
		if col+span > cols {
			col = 0
			row++
		}

		cellW := maxCellW*float32(span) + colGap*float32(span-1)
		n.Width = cellW
		n.X = padX + float32(col)*(maxCellW+colGap) + cellW/2
		n.Y = padY + float32(row)*(maxCellH+rowGap) + maxCellH/2

		col += span
		if col >= cols {
			col = 0
			row++
		}
	}

	// Route any edges over the grid
	var edges []*EdgeLayout
	for _, e := range g.Edges {
		src := nodes[e.From]
		dst := nodes[e.To]
		if src == nil || dst == nil {
			continue
		}
		edges = append(edges, &EdgeLayout{
			From:     e.From,
			To:       e.To,
			Points:   [][2]float32{{src.X, src.Y}, {dst.X, dst.Y}},
			ArrowEnd: e.ArrowEnd,
		})
	}

	totalW := padX*2 + float32(cols)*maxCellW + float32(cols-1)*colGap
	totalRows := row + 1
	if col == 0 && row > 0 {
		totalRows = row
	}
	totalH := padY*2 + float32(totalRows)*maxCellH + float32(totalRows-1)*rowGap

	return &Layout{
		Kind:    g.Kind,
		Nodes:   nodes,
		Edges:   edges,
		Width:   totalW,
		Height:  totalH,
		Diagram: BlockData{Columns: cols, BlockInfos: blockInfos},
	}
}

func sizeBlockNodes(g *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(g.Nodes))
	for id, node := range g.Nodes {
		nl := sizeNode(node, measurer, th, cfg)
		nodes[id] = nl
	}
	return nodes
}
```

**Step 3: Create `layout/block_test.go`**

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestBlockGridLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Block
	g.BlockColumns = 3

	for _, id := range []string{"a", "b", "c", "d", "e"} {
		label := id
		g.EnsureNode(id, &label, nil)
		g.Blocks = append(g.Blocks, &ir.BlockDef{ID: id, Label: id, Shape: ir.Rectangle, Width: 1})
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Block {
		t.Fatalf("Kind = %v", l.Kind)
	}
	bd, ok := l.Diagram.(BlockData)
	if !ok {
		t.Fatal("Diagram is not BlockData")
	}
	if bd.Columns != 3 {
		t.Errorf("Columns = %d, want 3", bd.Columns)
	}
	if len(l.Nodes) != 5 {
		t.Errorf("Nodes = %d, want 5", len(l.Nodes))
	}
	// First row: a, b, c should have same Y
	ay := l.Nodes["a"].Y
	by := l.Nodes["b"].Y
	cy := l.Nodes["c"].Y
	if ay != by || by != cy {
		t.Errorf("Row 1 Y mismatch: a=%v b=%v c=%v", ay, by, cy)
	}
	// Second row: d, e should have same Y, different from first row
	dy := l.Nodes["d"].Y
	if dy == ay {
		t.Error("Row 2 should have different Y than row 1")
	}
}

func TestBlockSpanLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Block
	g.BlockColumns = 3

	aLabel := "A"
	bLabel := "B"
	g.EnsureNode("a", &aLabel, nil)
	g.EnsureNode("b", &bLabel, nil)
	g.Blocks = append(g.Blocks,
		&ir.BlockDef{ID: "a", Label: "A", Shape: ir.Rectangle, Width: 2},
		&ir.BlockDef{ID: "b", Label: "B", Shape: ir.Rectangle, Width: 1},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	// a spans 2 columns, so should be wider than b
	if l.Nodes["a"].Width <= l.Nodes["b"].Width {
		t.Errorf("a width (%v) should be > b width (%v)", l.Nodes["a"].Width, l.Nodes["b"].Width)
	}
}

func TestBlockSugiyamaFallback(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Block
	// No columns directive, but has edges

	aLabel := "A"
	bLabel := "B"
	g.EnsureNode("a", &aLabel, nil)
	g.EnsureNode("b", &bLabel, nil)
	g.Blocks = append(g.Blocks,
		&ir.BlockDef{ID: "a", Label: "A", Width: 1},
		&ir.BlockDef{ID: "b", Label: "B", Width: 1},
	)
	g.Edges = append(g.Edges, &ir.Edge{From: "a", To: "b", Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if len(l.Edges) != 1 {
		t.Errorf("Edges = %d, want 1", len(l.Edges))
	}
}

func TestBlockLayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Block
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)
	if len(l.Nodes) != 0 {
		t.Errorf("Nodes = %d", len(l.Nodes))
	}
}
```

**Step 4: Add dispatch case to `layout/layout.go`**

```go
	case ir.Block:
		return computeBlockLayout(g, th, cfg)
```

**Step 5: Run tests**

Run: `go test ./layout/ -v -run TestBlock`
Expected: All 4 tests PASS

**Step 6: Commit**

```bash
git add layout/block.go layout/block_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Block diagram layout with hybrid grid/Sugiyama"
```

---

### Task 10: C4 Layout

**Files:**
- Create: `layout/c4.go`
- Create: `layout/c4_test.go`
- Modify: `layout/types.go` (add C4Data)
- Modify: `layout/layout.go` (add dispatch case)

**Step 1: Add types to `layout/types.go`**

```go
// C4Data holds C4-diagram-specific layout data.
type C4Data struct {
	Elements     map[string]*ir.C4Element
	Boundaries   []*C4BoundaryLayout
	SubKind      ir.C4Kind
}

func (C4Data) diagramData() {}

// C4BoundaryLayout stores positioned boundary rectangles.
type C4BoundaryLayout struct {
	ID     string
	Label  string
	Type   string
	X, Y   float32
	Width  float32
	Height float32
}
```

**Step 2: Create `layout/c4.go`**

```go
package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/textmetrics"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeC4Layout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	nodes := sizeC4Nodes(g, measurer, th, cfg)

	r := runSugiyama(g, nodes, cfg)

	// Build element map
	elemMap := make(map[string]*ir.C4Element)
	for _, elem := range g.C4Elements {
		elemMap[elem.ID] = elem
	}

	// Compute boundary rectangles
	var boundaryLayouts []*C4BoundaryLayout
	for _, b := range g.C4Boundaries {
		bl := computeC4BoundaryRect(b, nodes, cfg)
		if bl != nil {
			boundaryLayouts = append(boundaryLayouts, bl)
		}
	}

	return &Layout{
		Kind:   g.Kind,
		Nodes:  nodes,
		Edges:  r.Edges,
		Width:  r.Width,
		Height: r.Height,
		Diagram: C4Data{
			Elements:   elemMap,
			Boundaries: boundaryLayouts,
			SubKind:    g.C4SubKind,
		},
	}
}

func sizeC4Nodes(g *ir.Graph, measurer *textmetrics.Measurer, th *theme.Theme, cfg *config.Layout) map[string]*NodeLayout {
	nodes := make(map[string]*NodeLayout, len(g.Nodes))
	lineH := th.FontSize * cfg.LabelLineHeight
	smallFontSize := th.FontSize * 0.85
	smallLineH := smallFontSize * cfg.LabelLineHeight
	padH := cfg.Padding.NodeHorizontal
	padV := cfg.Padding.NodeVertical

	elemMap := make(map[string]*ir.C4Element)
	for _, elem := range g.C4Elements {
		elemMap[elem.ID] = elem
	}

	for id, node := range g.Nodes {
		elem := elemMap[id]
		var maxW, totalH float32

		// Label
		labelW := measurer.Width(node.Label, th.FontSize, th.FontFamily)
		if labelW > maxW {
			maxW = labelW
		}
		totalH += lineH

		if elem != nil {
			// Technology line
			if elem.Technology != "" {
				techText := "[" + elem.Technology + "]"
				techW := measurer.Width(techText, smallFontSize, th.FontFamily)
				if techW > maxW {
					maxW = techW
				}
				totalH += smallLineH
			}
			// Description
			if elem.Description != "" {
				descW := measurer.Width(elem.Description, smallFontSize, th.FontFamily)
				if descW > maxW {
					maxW = descW
				}
				totalH += smallLineH
			}

			// Apply min dimensions based on element type
			if elem.Type.IsPerson() {
				if maxW+2*padH < cfg.C4.PersonWidth {
					maxW = cfg.C4.PersonWidth - 2*padH
				}
				if totalH+2*padV < cfg.C4.PersonHeight {
					totalH = cfg.C4.PersonHeight - 2*padV
				}
			} else {
				if maxW+2*padH < cfg.C4.SystemWidth {
					maxW = cfg.C4.SystemWidth - 2*padH
				}
				if totalH+2*padV < cfg.C4.SystemHeight {
					totalH = cfg.C4.SystemHeight - 2*padV
				}
			}
		}

		nodes[id] = &NodeLayout{
			ID:     id,
			Label:  TextBlock{Lines: []string{node.Label}, Width: maxW, Height: lineH, FontSize: th.FontSize},
			Shape:  ir.Rectangle,
			Width:  maxW + 2*padH,
			Height: totalH + 2*padV,
		}
	}

	return nodes
}

func computeC4BoundaryRect(b *ir.C4Boundary, nodes map[string]*NodeLayout, cfg *config.Layout) *C4BoundaryLayout {
	if len(b.Children) == 0 {
		return nil
	}
	pad := cfg.C4.BoundaryPadding

	var minX, minY, maxX, maxY float32
	first := true
	for _, childID := range b.Children {
		n, ok := nodes[childID]
		if !ok {
			continue
		}
		left := n.X - n.Width/2
		top := n.Y - n.Height/2
		right := n.X + n.Width/2
		bottom := n.Y + n.Height/2
		if first {
			minX, minY, maxX, maxY = left, top, right, bottom
			first = false
		} else {
			if left < minX {
				minX = left
			}
			if top < minY {
				minY = top
			}
			if right > maxX {
				maxX = right
			}
			if bottom > maxY {
				maxY = bottom
			}
		}
	}
	if first {
		return nil
	}

	return &C4BoundaryLayout{
		ID:     b.ID,
		Label:  b.Label,
		Type:   b.Type,
		X:      minX - pad,
		Y:      minY - pad - 20, // extra space for label
		Width:  (maxX - minX) + 2*pad,
		Height: (maxY - minY) + 2*pad + 20,
	}
}
```

**Step 3: Create `layout/c4_test.go`**

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestC4Layout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.C4
	g.C4SubKind = ir.C4Context

	userLabel := "User"
	sysLabel := "System"
	g.EnsureNode("user", &userLabel, nil)
	g.EnsureNode("sys", &sysLabel, nil)

	g.C4Elements = append(g.C4Elements,
		&ir.C4Element{ID: "user", Label: "User", Type: ir.C4Person},
		&ir.C4Element{ID: "sys", Label: "System", Type: ir.C4System},
	)
	relLabel := "Uses"
	g.Edges = append(g.Edges, &ir.Edge{From: "user", To: "sys", Label: &relLabel, Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.C4 {
		t.Fatalf("Kind = %v", l.Kind)
	}
	cd, ok := l.Diagram.(C4Data)
	if !ok {
		t.Fatal("Diagram is not C4Data")
	}
	if len(cd.Elements) != 2 {
		t.Errorf("Elements = %d", len(cd.Elements))
	}
	if cd.SubKind != ir.C4Context {
		t.Errorf("SubKind = %v", cd.SubKind)
	}
}

func TestC4LayoutWithBoundary(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.C4
	g.C4SubKind = ir.C4Container

	apiLabel := "API"
	dbLabel := "DB"
	g.EnsureNode("api", &apiLabel, nil)
	g.EnsureNode("db", &dbLabel, nil)

	g.C4Elements = append(g.C4Elements,
		&ir.C4Element{ID: "api", Label: "API", Type: ir.C4Container_, BoundaryID: "sys"},
		&ir.C4Element{ID: "db", Label: "DB", Type: ir.C4ContainerDb, BoundaryID: "sys"},
	)
	g.C4Boundaries = append(g.C4Boundaries, &ir.C4Boundary{
		ID: "sys", Label: "My System", Type: "Software System", Children: []string{"api", "db"},
	})
	relLabel := "reads"
	g.Edges = append(g.Edges, &ir.Edge{From: "api", To: "db", Label: &relLabel, Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	cd, ok := l.Diagram.(C4Data)
	if !ok {
		t.Fatal("Diagram is not C4Data")
	}
	if len(cd.Boundaries) != 1 {
		t.Fatalf("Boundaries = %d, want 1", len(cd.Boundaries))
	}
	b := cd.Boundaries[0]
	if b.Label != "My System" {
		t.Errorf("boundary label = %q", b.Label)
	}
	if b.Width <= 0 || b.Height <= 0 {
		t.Errorf("boundary size = %vx%v", b.Width, b.Height)
	}
}

func TestC4LayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.C4
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)
	if len(l.Nodes) != 0 {
		t.Errorf("Nodes = %d", len(l.Nodes))
	}
}
```

**Step 4: Add dispatch case to `layout/layout.go`**

```go
	case ir.C4:
		return computeC4Layout(g, th, cfg)
```

**Step 5: Run tests**

Run: `go test ./layout/ -v -run TestC4`
Expected: All 3 tests PASS

**Step 6: Commit**

```bash
git add layout/c4.go layout/c4_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add C4 diagram layout"
```

---

### Task 11: Requirement Renderer

**Files:**
- Create: `render/requirement.go`
- Create: `render/requirement_test.go`
- Modify: `render/svg.go` (add dispatch case)

**Step 1: Create `render/requirement.go`**

Renders requirement/element nodes as rounded rectangles with stereotype header, name, and metadata rows. Edges with relationship labels.

**Step 2: Create `render/requirement_test.go`** — 2 tests (basic render, empty)

**Step 3: Add `case layout.RequirementData` to `render/svg.go`**

**Step 4: Run tests, commit**

```bash
git commit -m "feat(render): add Requirement diagram SVG renderer"
```

---

### Task 12: Block Renderer

**Files:**
- Create: `render/block.go`
- Create: `render/block_test.go`
- Modify: `render/svg.go` (add dispatch case)

**Step 1: Create `render/block.go`**

Renders block nodes using shape-specific rendering (reuse `renderNodeShape` from graph renderer). Edges drawn as simple arrows over grid.

**Step 2: Create `render/block_test.go`** — 2 tests (grid render, empty)

**Step 3: Add `case layout.BlockData` to `render/svg.go`**

**Step 4: Run tests, commit**

```bash
git commit -m "feat(render): add Block diagram SVG renderer"
```

---

### Task 13: C4 Renderer

**Files:**
- Create: `render/c4.go`
- Create: `render/c4_test.go`
- Modify: `render/svg.go` (add dispatch case)

**Step 1: Create `render/c4.go`**

Renders C4 elements with type-specific styling: Person elements as rounded boxes with person icon, Systems as blue boxes, Containers with technology subtitle, boundaries as dashed rectangles.

**Step 2: Create `render/c4_test.go`** — 2 tests (context render, empty)

**Step 3: Add `case layout.C4Data` to `render/svg.go`**

**Step 4: Run tests, commit**

```bash
git commit -m "feat(render): add C4 diagram SVG renderer"
```

---

### Task 14: Integration Tests and Fixtures

**Files:**
- Create: `testdata/fixtures/requirement-basic.mmd`
- Create: `testdata/fixtures/requirement-multiple.mmd`
- Create: `testdata/fixtures/block-grid.mmd`
- Create: `testdata/fixtures/block-edges.mmd`
- Create: `testdata/fixtures/c4-context.mmd`
- Create: `testdata/fixtures/c4-container.mmd`
- Modify: `mermaid_test.go` (add 6 integration tests)

**Step 1: Create fixture files**

`requirement-basic.mmd`:
```
requirementDiagram

requirement test_req {
id: 1
text: the test text.
risk: high
verifymethod: test
}

element test_entity {
type: simulation
}

test_entity - satisfies -> test_req
```

`block-grid.mmd`:
```
block-beta
columns 3
a["Block A"] b["Block B"] c["Block C"]
d["Block D"]:2 e["Block E"]
```

`c4-context.mmd`:
```
C4Context
Person(user, "User", "A user of the system")
System(webapp, "Web Application", "Main web app")
System_Ext(email, "Email System", "Sends emails")
Rel(user, webapp, "Uses", "HTTPS")
Rel(webapp, email, "Sends notifications")
```

**Step 2: Add integration tests following existing pattern**

**Step 3: Run `go test ./...`**
Expected: All 8 packages PASS

**Step 4: Commit**

```bash
git commit -m "test: add integration tests and fixtures for Requirement, Block, and C4"
```

---

### Task 15: Final Validation

**Step 1:** Run `go vet ./...`
**Step 2:** Run `go build ./...`
**Step 3:** Run `gofmt -l .` — fix any formatting issues
**Step 4:** Run `go test -race ./...`
**Step 5:** Run go-code-reviewer agent
**Step 6:** Commit any fixes, close beads
