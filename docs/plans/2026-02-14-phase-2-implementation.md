# Phase 2: Core Graph Variants Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add full mermaid-js parity for class diagrams, state diagrams (v2), and ER diagrams.

**Architecture:** Extend the existing 4-stage pipeline (Parse -> IR -> Layout -> Render). Class and ER reuse the Sugiyama graph layout with custom node sizing. State gets a recursive layout for composite/nested states. Each diagram type gets its own parser, layout adapter, and renderer file.

**Tech Stack:** Go 1.23+, stdlib only (regexp, strings, math), `golang.org/x/image/font/sfnt` for text metrics, `github.com/google/go-cmp` for test diffs.

---

### Task 1: IR Types — Class Diagram

**Files:**
- Create: `ir/class.go`
- Create: `ir/class_test.go`
- Modify: `ir/graph.go:35-51` (add class fields to Graph struct)

**Step 1: Write the failing test**

Create `ir/class_test.go`:

```go
package ir

import "testing"

func TestClassMemberIsMethod(t *testing.T) {
	tests := []struct {
		name     string
		member   ClassMember
		isMethod bool
	}{
		{"method with parens", ClassMember{Name: "getAge", IsMethod: true}, true},
		{"attribute", ClassMember{Name: "age", IsMethod: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.member.IsMethod != tt.isMethod {
				t.Errorf("IsMethod = %v, want %v", tt.member.IsMethod, tt.isMethod)
			}
		})
	}
}

func TestVisibilitySymbol(t *testing.T) {
	tests := []struct {
		vis  Visibility
		want string
	}{
		{VisPublic, "+"},
		{VisPrivate, "-"},
		{VisProtected, "#"},
		{VisPackage, "~"},
		{VisNone, ""},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.vis.Symbol(); got != tt.want {
				t.Errorf("Symbol() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestClassMember -v`
Expected: FAIL — `ClassMember`, `Visibility` undefined

**Step 3: Write minimal implementation**

Create `ir/class.go`:

```go
package ir

// Visibility represents the access level of a class member.
type Visibility int

const (
	VisNone      Visibility = iota
	VisPublic               // +
	VisPrivate              // -
	VisProtected            // #
	VisPackage              // ~
)

// Symbol returns the single-character symbol for this visibility.
func (v Visibility) Symbol() string {
	switch v {
	case VisPublic:
		return "+"
	case VisPrivate:
		return "-"
	case VisProtected:
		return "#"
	case VisPackage:
		return "~"
	default:
		return ""
	}
}

// MemberClassifier marks a member as abstract or static.
type MemberClassifier int

const (
	ClassifierNone     MemberClassifier = iota
	ClassifierAbstract                  // *
	ClassifierStatic                    // $
)

// ClassMember represents a single attribute or method of a class.
type ClassMember struct {
	Name       string
	Type       string           // attribute type or return type
	Params     string           // raw parameter string (empty for attributes)
	IsMethod   bool             // true if member has parentheses
	Visibility Visibility
	Classifier MemberClassifier
	Generic    string           // generic type parameter e.g. "T"
}

// ClassMembers groups a class's attributes and methods.
type ClassMembers struct {
	Attributes []ClassMember
	Methods    []ClassMember
}

// Namespace groups classes under a shared name.
type Namespace struct {
	Name    string
	Classes []string // node IDs
}

// DiagramNote is a text note attached to a node or floating.
type DiagramNote struct {
	Text     string
	Position string // "right of", "left of", or "" for floating
	Target   string // node ID, empty for floating notes
}
```

Add fields to `Graph` struct in `ir/graph.go` (after `EdgeStyleDefault`):

```go
	// Class diagram fields
	Members     map[string]*ClassMembers
	Annotations map[string]string   // node ID -> stereotype e.g. "interface"
	Namespaces  []*Namespace
	Notes       []*DiagramNote
```

Update `NewGraph()` to initialize: `Members: make(map[string]*ClassMembers)`

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestClassMember -v && go test ./ir/ -run TestVisibility -v`
Expected: PASS

**Step 5: Commit**

```
git add ir/class.go ir/class_test.go ir/graph.go
git commit -m "feat(ir): add class diagram types"
```

---

### Task 2: IR Types — State Diagram

**Files:**
- Create: `ir/state.go`
- Create: `ir/state_test.go`
- Modify: `ir/graph.go` (add state fields to Graph struct)

**Step 1: Write the failing test**

Create `ir/state_test.go`:

```go
package ir

import "testing"

func TestStateAnnotationString(t *testing.T) {
	tests := []struct {
		ann  StateAnnotation
		want string
	}{
		{StateChoice, "choice"},
		{StateFork, "fork"},
		{StateJoin, "join"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.ann.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestStateAnnotation -v`
Expected: FAIL — `StateAnnotation` undefined

**Step 3: Write minimal implementation**

Create `ir/state.go`:

```go
package ir

// StateAnnotation marks a state as a special type.
type StateAnnotation int

const (
	StateChoice StateAnnotation = iota
	StateFork
	StateJoin
)

// String returns the annotation keyword.
func (a StateAnnotation) String() string {
	switch a {
	case StateChoice:
		return "choice"
	case StateFork:
		return "fork"
	case StateJoin:
		return "join"
	default:
		return ""
	}
}

// CompositeState holds a nested state machine within a parent state.
type CompositeState struct {
	ID        string
	Label     string
	Inner     *Graph   // primary nested state machine
	Regions   []*Graph // concurrent regions (separated by --)
	Direction *Direction
}
```

Add fields to `Graph` struct in `ir/graph.go`:

```go
	// State diagram fields
	CompositeStates   map[string]*CompositeState
	StateDescriptions map[string]string          // state ID -> description text
	StateAnnotations  map[string]StateAnnotation // state ID -> choice/fork/join
```

Update `NewGraph()`:
```go
CompositeStates:  make(map[string]*CompositeState),
StateDescriptions: make(map[string]string),
StateAnnotations: make(map[string]StateAnnotation),
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestStateAnnotation -v`
Expected: PASS

**Step 5: Commit**

```
git add ir/state.go ir/state_test.go ir/graph.go
git commit -m "feat(ir): add state diagram types"
```

---

### Task 3: IR Types — ER Diagram

**Files:**
- Create: `ir/er.go`
- Create: `ir/er_test.go`
- Modify: `ir/graph.go` (add Entities field)

**Step 1: Write the failing test**

Create `ir/er_test.go`:

```go
package ir

import "testing"

func TestAttributeKeyString(t *testing.T) {
	tests := []struct {
		key  AttributeKey
		want string
	}{
		{KeyPrimary, "PK"},
		{KeyForeign, "FK"},
		{KeyUnique, "UK"},
		{KeyNone, ""},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.key.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestAttributeKey -v`
Expected: FAIL — `AttributeKey` undefined

**Step 3: Write minimal implementation**

Create `ir/er.go`:

```go
package ir

// AttributeKey represents a key constraint on an entity attribute.
type AttributeKey int

const (
	KeyNone    AttributeKey = iota
	KeyPrimary              // PK
	KeyForeign              // FK
	KeyUnique               // UK
)

// String returns the key constraint abbreviation.
func (k AttributeKey) String() string {
	switch k {
	case KeyPrimary:
		return "PK"
	case KeyForeign:
		return "FK"
	case KeyUnique:
		return "UK"
	default:
		return ""
	}
}

// EntityAttribute is a single attribute in an ER entity.
type EntityAttribute struct {
	Type    string
	Name    string
	Keys    []AttributeKey
	Comment string
}

// Entity represents an ER diagram entity with typed attributes.
type Entity struct {
	ID         string
	Label      string // display name (alias), empty means use ID
	Attributes []EntityAttribute
}

// DisplayName returns Label if set, otherwise ID.
func (e *Entity) DisplayName() string {
	if e.Label != "" {
		return e.Label
	}
	return e.ID
}
```

Add to `Graph` struct in `ir/graph.go`:

```go
	// ER diagram fields
	Entities map[string]*Entity
```

Update `NewGraph()`: `Entities: make(map[string]*Entity)`

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestAttributeKey -v`
Expected: PASS

**Step 5: Commit**

```
git add ir/er.go ir/er_test.go ir/graph.go
git commit -m "feat(ir): add ER diagram types"
```

---

### Task 4: IR Types — New Edge Arrowheads

**Files:**
- Modify: `ir/shapes.go:47-52` (add new EdgeArrowhead constants)

**Step 1: Write the failing test**

Add to `ir/graph_test.go`:

```go
func TestEdgeArrowheadValues(t *testing.T) {
	// Verify new arrowhead constants exist and have distinct values.
	heads := []ir.EdgeArrowhead{
		ir.OpenTriangle,
		ir.ClassDependency,
		ir.ClosedTriangle,
		ir.FilledDiamond,
		ir.OpenDiamond,
		ir.Lollipop,
	}
	seen := make(map[ir.EdgeArrowhead]bool)
	for _, h := range heads {
		if seen[h] {
			t.Errorf("duplicate arrowhead value: %d", h)
		}
		seen[h] = true
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestEdgeArrowhead -v`
Expected: FAIL — `ClosedTriangle`, `FilledDiamond`, `OpenDiamond`, `Lollipop` undefined

**Step 3: Write minimal implementation**

Modify `ir/shapes.go`, extend the `EdgeArrowhead` const block:

```go
const (
	OpenTriangle    EdgeArrowhead = iota
	ClassDependency
	ClosedTriangle  // inheritance, realization
	FilledDiamond   // composition
	OpenDiamond     // aggregation
	Lollipop        // provided interface
)
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestEdgeArrowhead -v`
Expected: PASS

**Step 5: Commit**

```
git add ir/shapes.go ir/graph_test.go
git commit -m "feat(ir): add class/ER edge arrowhead types"
```

---

### Task 5: Theme & Config Extensions

**Files:**
- Modify: `theme/theme.go:5-53` (add class/state/ER color fields)
- Modify: `theme/theme.go:56-110` (Modern preset)
- Modify: `theme/theme.go:113-167` (MermaidDefault preset)
- Modify: `config/config.go` (add ClassConfig, StateConfig, ERConfig)

**Step 1: Write the failing test**

Add to `theme/theme_test.go`:

```go
func TestModernThemeHasClassColors(t *testing.T) {
	th := Modern()
	if th.ClassHeaderBg == "" {
		t.Error("ClassHeaderBg should not be empty")
	}
	if th.StateFill == "" {
		t.Error("StateFill should not be empty")
	}
	if th.EntityHeaderBg == "" {
		t.Error("EntityHeaderBg should not be empty")
	}
}
```

Add to `config/config_test.go`:

```go
func TestDefaultLayoutHasClassConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Class.CompartmentPadX <= 0 {
		t.Error("Class.CompartmentPadX should be > 0")
	}
	if cfg.State.CompositePadding <= 0 {
		t.Error("State.CompositePadding should be > 0")
	}
	if cfg.ER.AttributeRowHeight <= 0 {
		t.Error("ER.AttributeRowHeight should be > 0")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./theme/ -run TestModernThemeHasClassColors -v && go test ./config/ -run TestDefaultLayoutHasClassConfig -v`
Expected: FAIL — fields not found

**Step 3: Write minimal implementation**

Add to `Theme` struct in `theme/theme.go`:

```go
	// Class diagram colors
	ClassHeaderBg string
	ClassBodyBg   string
	ClassBorder   string

	// State diagram colors
	StateFill         string
	StateBorder       string
	StateStartEnd     string
	CompositeHeaderBg string

	// ER diagram colors
	EntityHeaderBg string
	EntityBodyBg   string
	EntityBorder   string
```

Add values to `Modern()`:
```go
ClassHeaderBg:     "#4C78A8",
ClassBodyBg:       "#F0F4F8",
ClassBorder:       "#3B6492",
StateFill:         "#F0F4F8",
StateBorder:       "#3B6492",
StateStartEnd:     "#333344",
CompositeHeaderBg: "#E8EFF5",
EntityHeaderBg:    "#4C78A8",
EntityBodyBg:      "#F0F4F8",
EntityBorder:      "#3B6492",
```

Add values to `MermaidDefault()`:
```go
ClassHeaderBg:     "#ECECFF",
ClassBodyBg:       "#FFFFFF",
ClassBorder:       "#9370DB",
StateFill:         "#ECECFF",
StateBorder:       "#9370DB",
StateStartEnd:     "#333",
CompositeHeaderBg: "#f4f4f4",
EntityHeaderBg:    "#ECECFF",
EntityBodyBg:      "#FFFFFF",
EntityBorder:      "#9370DB",
```

Add to `config/config.go`:

```go
type ClassConfig struct {
	CompartmentPadX float32
	CompartmentPadY float32
	MemberFontSize  float32
}

type StateConfig struct {
	CompositePadding   float32
	RegionSeparatorPad float32
	StartEndRadius     float32
	ForkBarWidth       float32
	ForkBarHeight      float32
}

type ERConfig struct {
	AttributeRowHeight float32
	ColumnPadding      float32
	HeaderPadding      float32
}
```

Add fields to `Layout` struct:
```go
Class ClassConfig
State StateConfig
ER    ERConfig
```

Add defaults to `DefaultLayout()`:
```go
Class: ClassConfig{
	CompartmentPadX: 12,
	CompartmentPadY: 6,
	MemberFontSize:  12,
},
State: StateConfig{
	CompositePadding:   20,
	RegionSeparatorPad: 10,
	StartEndRadius:     8,
	ForkBarWidth:       80,
	ForkBarHeight:      6,
},
ER: ERConfig{
	AttributeRowHeight: 22,
	ColumnPadding:      10,
	HeaderPadding:      8,
},
```

**Step 4: Run tests to verify they pass**

Run: `go test ./theme/ -v && go test ./config/ -v`
Expected: PASS

**Step 5: Commit**

```
git add theme/theme.go theme/theme_test.go config/config.go config/config_test.go
git commit -m "feat(theme,config): add class/state/ER colors and config"
```

---

### Task 6: Class Diagram Parser

**Files:**
- Create: `parser/class.go`
- Create: `parser/class_test.go`
- Modify: `parser/parser.go:18-23` (add Class case to switch)

**Step 1: Write the failing tests**

Create `parser/class_test.go` with table-driven tests:

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseClassSimple(t *testing.T) {
	input := `classDiagram
    class Animal {
        +String name
        +int age
        +isMammal() bool
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.Kind != ir.Class {
		t.Errorf("Kind = %v, want Class", out.Graph.Kind)
	}
	members, ok := out.Graph.Members["Animal"]
	if !ok {
		t.Fatal("Animal not in Members map")
	}
	if len(members.Attributes) != 2 {
		t.Errorf("attributes count = %d, want 2", len(members.Attributes))
	}
	if len(members.Methods) != 1 {
		t.Errorf("methods count = %d, want 1", len(members.Methods))
	}
	if members.Methods[0].Type != "bool" {
		t.Errorf("method return type = %q, want %q", members.Methods[0].Type, "bool")
	}
}

func TestParseClassRelationships(t *testing.T) {
	tests := []struct {
		name  string
		input string
		edges int
	}{
		{"inheritance", "classDiagram\n    Animal <|-- Dog", 1},
		{"composition", "classDiagram\n    Car *-- Engine", 1},
		{"aggregation", "classDiagram\n    Library o-- Book", 1},
		{"association", "classDiagram\n    Student --> Course", 1},
		{"dependency", "classDiagram\n    Class1 ..> Class2", 1},
		{"realization", "classDiagram\n    Animal ..|> Walkable", 1},
		{"link solid", "classDiagram\n    A -- B", 1},
		{"link dashed", "classDiagram\n    A .. B", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error: %v", err)
			}
			if len(out.Graph.Edges) != tt.edges {
				t.Errorf("edges = %d, want %d", len(out.Graph.Edges), tt.edges)
			}
		})
	}
}

func TestParseClassAnnotation(t *testing.T) {
	input := `classDiagram
    class Shape
    <<interface>> Shape`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	ann, ok := out.Graph.Annotations["Shape"]
	if !ok {
		t.Fatal("Shape not in Annotations")
	}
	if ann != "interface" {
		t.Errorf("annotation = %q, want %q", ann, "interface")
	}
}

func TestParseClassDirection(t *testing.T) {
	input := `classDiagram
    direction LR
    A <|-- B`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.Direction != ir.LeftRight {
		t.Errorf("Direction = %v, want LeftRight", out.Graph.Direction)
	}
}

func TestParseClassVisibility(t *testing.T) {
	input := `classDiagram
    class MyClass {
        +publicField
        -privateField
        #protectedField
        ~packageField
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	members := out.Graph.Members["MyClass"]
	if members == nil {
		t.Fatal("MyClass not in Members")
	}
	vis := []ir.Visibility{ir.VisPublic, ir.VisPrivate, ir.VisProtected, ir.VisPackage}
	for i, want := range vis {
		if i >= len(members.Attributes) {
			t.Fatalf("only %d attributes, want at least %d", len(members.Attributes), i+1)
		}
		if members.Attributes[i].Visibility != want {
			t.Errorf("attr[%d] visibility = %d, want %d", i, members.Attributes[i].Visibility, want)
		}
	}
}

func TestParseClassCardinality(t *testing.T) {
	input := `classDiagram
    Customer "1" --> "*" Order : places`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(out.Graph.Edges) != 1 {
		t.Fatalf("edges = %d, want 1", len(out.Graph.Edges))
	}
	e := out.Graph.Edges[0]
	if e.StartLabel == nil || *e.StartLabel != "1" {
		t.Errorf("StartLabel = %v, want '1'", e.StartLabel)
	}
	if e.EndLabel == nil || *e.EndLabel != "*" {
		t.Errorf("EndLabel = %v, want '*'", e.EndLabel)
	}
	if e.Label == nil || *e.Label != "places" {
		t.Errorf("Label = %v, want 'places'", e.Label)
	}
}

func TestParseClassNamespace(t *testing.T) {
	input := `classDiagram
    namespace BaseShapes {
        class Triangle
        class Rectangle
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(out.Graph.Namespaces) != 1 {
		t.Fatalf("namespaces = %d, want 1", len(out.Graph.Namespaces))
	}
	ns := out.Graph.Namespaces[0]
	if ns.Name != "BaseShapes" {
		t.Errorf("namespace name = %q, want %q", ns.Name, "BaseShapes")
	}
	if len(ns.Classes) != 2 {
		t.Errorf("namespace classes = %d, want 2", len(ns.Classes))
	}
}

func TestParseClassGeneric(t *testing.T) {
	input := `classDiagram
    class List~T~ {
        +add(T item)
        +get(int index) T
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if _, ok := out.Graph.Nodes["List"]; !ok {
		t.Error("expected node List")
	}
}

func TestParseClassColon(t *testing.T) {
	input := `classDiagram
    Animal : +int age
    Animal : +String name`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	members := out.Graph.Members["Animal"]
	if members == nil {
		t.Fatal("Animal not in Members")
	}
	if len(members.Attributes) != 2 {
		t.Errorf("attributes = %d, want 2", len(members.Attributes))
	}
}

func TestParseClassNote(t *testing.T) {
	input := `classDiagram
    class Animal
    note for Animal "This is a note"`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(out.Graph.Notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(out.Graph.Notes))
	}
	if out.Graph.Notes[0].Target != "Animal" {
		t.Errorf("note target = %q, want %q", out.Graph.Notes[0].Target, "Animal")
	}
}

func TestParseClassClassifier(t *testing.T) {
	input := `classDiagram
    class MyClass {
        +abstractMethod()*
        +staticMethod()$
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	members := out.Graph.Members["MyClass"]
	if members == nil {
		t.Fatal("MyClass not in Members")
	}
	if len(members.Methods) != 2 {
		t.Fatalf("methods = %d, want 2", len(members.Methods))
	}
	if members.Methods[0].Classifier != ir.ClassifierAbstract {
		t.Errorf("method[0] classifier = %d, want Abstract", members.Methods[0].Classifier)
	}
	if members.Methods[1].Classifier != ir.ClassifierStatic {
		t.Errorf("method[1] classifier = %d, want Static", members.Methods[1].Classifier)
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./parser/ -run TestParseClass -v`
Expected: FAIL — `parseClass` not found

**Step 3: Write implementation**

Create `parser/class.go` implementing `parseClass(input string) (*ParseOutput, error)`. Key pieces:

- Header regex: `(?i)^classdiagram(-v2)?`
- Relationship regex covering all 8 types + bidirectional + cardinality
- Class body parsing (brace-delimited block)
- Member parsing: visibility prefix detection, method detection (parentheses), return type, classifiers (`*`, `$`)
- Namespace blocks (brace-delimited)
- Annotation parsing: `<<keyword>> ClassName` or inline `class Name <<keyword>>`
- Colon member syntax: `ClassName : member`
- Note parsing: `note "text"` and `note for ClassName "text"`
- Direction parsing: reuse existing `parseDirectionLine()`
- Directive skip: classDef, style, cssClass, click, callback, link

Wire into `parser/parser.go` switch:
```go
case ir.Class:
    return parseClass(input)
```

**Step 4: Run tests to verify they pass**

Run: `go test ./parser/ -run TestParseClass -v`
Expected: ALL PASS

**Step 5: Commit**

```
git add parser/class.go parser/class_test.go parser/parser.go
git commit -m "feat(parser): add class diagram parser with full syntax support"
```

---

### Task 7: State Diagram Parser

**Files:**
- Create: `parser/state.go`
- Create: `parser/state_test.go`
- Modify: `parser/parser.go:18-23` (add State case)

**Step 1: Write the failing tests**

Create `parser/state_test.go`:

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseStateSimple(t *testing.T) {
	input := `stateDiagram-v2
    [*] --> First
    First --> Second
    Second --> [*]`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.Kind != ir.State {
		t.Errorf("Kind = %v, want State", out.Graph.Kind)
	}
	if len(out.Graph.Edges) != 3 {
		t.Errorf("edges = %d, want 3", len(out.Graph.Edges))
	}
}

func TestParseStateDescription(t *testing.T) {
	input := `stateDiagram-v2
    s1 : This is state s1`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	desc, ok := out.Graph.StateDescriptions["s1"]
	if !ok {
		t.Fatal("s1 not in StateDescriptions")
	}
	if desc != "This is state s1" {
		t.Errorf("description = %q, want %q", desc, "This is state s1")
	}
}

func TestParseStateAsKeyword(t *testing.T) {
	input := `stateDiagram-v2
    state "Moving state" as s1
    [*] --> s1`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if _, ok := out.Graph.Nodes["s1"]; !ok {
		t.Fatal("s1 not in Nodes")
	}
	desc, ok := out.Graph.StateDescriptions["s1"]
	if !ok {
		t.Fatal("s1 not in StateDescriptions")
	}
	if desc != "Moving state" {
		t.Errorf("description = %q, want %q", desc, "Moving state")
	}
}

func TestParseStateComposite(t *testing.T) {
	input := `stateDiagram-v2
    state First {
        [*] --> fir
        fir --> [*]
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	cs, ok := out.Graph.CompositeStates["First"]
	if !ok {
		t.Fatal("First not in CompositeStates")
	}
	if cs.Inner == nil {
		t.Fatal("Inner graph should not be nil")
	}
	if len(cs.Inner.Edges) != 2 {
		t.Errorf("inner edges = %d, want 2", len(cs.Inner.Edges))
	}
}

func TestParseStateChoice(t *testing.T) {
	input := `stateDiagram-v2
    state if_state <<choice>>
    [*] --> if_state`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	ann, ok := out.Graph.StateAnnotations["if_state"]
	if !ok {
		t.Fatal("if_state not in StateAnnotations")
	}
	if ann != ir.StateChoice {
		t.Errorf("annotation = %v, want StateChoice", ann)
	}
}

func TestParseStateForkJoin(t *testing.T) {
	input := `stateDiagram-v2
    state fork_state <<fork>>
    state join_state <<join>>`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.StateAnnotations["fork_state"] != ir.StateFork {
		t.Error("expected fork_state to be StateFork")
	}
	if out.Graph.StateAnnotations["join_state"] != ir.StateJoin {
		t.Error("expected join_state to be StateJoin")
	}
}

func TestParseStateConcurrent(t *testing.T) {
	input := `stateDiagram-v2
    state Active {
        [*] --> NumLockOff
        --
        [*] --> CapsLockOff
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	cs, ok := out.Graph.CompositeStates["Active"]
	if !ok {
		t.Fatal("Active not in CompositeStates")
	}
	if len(cs.Regions) != 2 {
		t.Errorf("regions = %d, want 2", len(cs.Regions))
	}
}

func TestParseStateTransitionLabel(t *testing.T) {
	input := `stateDiagram-v2
    s1 --> s2 : A transition`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(out.Graph.Edges) != 1 {
		t.Fatalf("edges = %d, want 1", len(out.Graph.Edges))
	}
	if out.Graph.Edges[0].Label == nil || *out.Graph.Edges[0].Label != "A transition" {
		t.Errorf("edge label = %v, want 'A transition'", out.Graph.Edges[0].Label)
	}
}

func TestParseStateDirection(t *testing.T) {
	input := `stateDiagram-v2
    direction LR
    [*] --> s1`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.Direction != ir.LeftRight {
		t.Errorf("Direction = %v, want LeftRight", out.Graph.Direction)
	}
}

func TestParseStateNote(t *testing.T) {
	input := `stateDiagram-v2
    State1
    note right of State1 : Important info`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if len(out.Graph.Notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(out.Graph.Notes))
	}
	if out.Graph.Notes[0].Position != "right of" {
		t.Errorf("note position = %q, want %q", out.Graph.Notes[0].Position, "right of")
	}
}

func TestParseStateBracketAnnotation(t *testing.T) {
	input := `stateDiagram-v2
    state fork_state [[fork]]`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.StateAnnotations["fork_state"] != ir.StateFork {
		t.Error("expected fork_state to be StateFork")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./parser/ -run TestParseState -v`
Expected: FAIL — `parseState` not found

**Step 3: Write implementation**

Create `parser/state.go` implementing `parseState(input string) (*ParseOutput, error)`. Key pieces:

- Header regex: `(?i)^statediagram(-v2)?`
- Transition regex: `id --> id` with optional `: label`
- `[*]` start/end state handling (mapped to special node IDs like `__start__` and `__end__`)
- State description: colon syntax (`s1 : desc`) and `state "desc" as s1` syntax
- Composite state: `state Name { ... }` with recursive parsing — collect lines between braces, call `parseStateBody()` recursively
- Annotations: `<<choice>>`, `<<fork>>`, `<<join>>` and `[[choice]]`, `[[fork]]`, `[[join]]`
- Concurrent regions: split composite body on standalone `--` lines
- Notes: `note right of State1 : text` and multi-line `note right of State1 ... end note`
- Direction directives within composites
- Wire into parser dispatch

**Step 4: Run tests to verify they pass**

Run: `go test ./parser/ -run TestParseState -v`
Expected: ALL PASS

**Step 5: Commit**

```
git add parser/state.go parser/state_test.go parser/parser.go
git commit -m "feat(parser): add state diagram parser with composite states"
```

---

### Task 8: ER Diagram Parser

**Files:**
- Create: `parser/er.go`
- Create: `parser/er_test.go`
- Modify: `parser/parser.go:18-23` (add Er case)

**Step 1: Write the failing tests**

Create `parser/er_test.go`:

```go
package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseERSimple(t *testing.T) {
	input := `erDiagram
    CUSTOMER ||--o{ ORDER : places`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.Kind != ir.Er {
		t.Errorf("Kind = %v, want Er", out.Graph.Kind)
	}
	if len(out.Graph.Edges) != 1 {
		t.Errorf("edges = %d, want 1", len(out.Graph.Edges))
	}
	if _, ok := out.Graph.Nodes["CUSTOMER"]; !ok {
		t.Error("CUSTOMER not in Nodes")
	}
	if _, ok := out.Graph.Nodes["ORDER"]; !ok {
		t.Error("ORDER not in Nodes")
	}
}

func TestParseERAttributes(t *testing.T) {
	input := `erDiagram
    CUSTOMER {
        string name
        int age PK
        string email UK "User email"
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	ent, ok := out.Graph.Entities["CUSTOMER"]
	if !ok {
		t.Fatal("CUSTOMER not in Entities")
	}
	if len(ent.Attributes) != 3 {
		t.Fatalf("attributes = %d, want 3", len(ent.Attributes))
	}
	// First attr: no key
	if len(ent.Attributes[0].Keys) != 0 {
		t.Errorf("attr[0] keys = %d, want 0", len(ent.Attributes[0].Keys))
	}
	// Second attr: PK
	if len(ent.Attributes[1].Keys) != 1 || ent.Attributes[1].Keys[0] != ir.KeyPrimary {
		t.Errorf("attr[1] expected PK")
	}
	// Third attr: UK with comment
	if len(ent.Attributes[2].Keys) != 1 || ent.Attributes[2].Keys[0] != ir.KeyUnique {
		t.Errorf("attr[2] expected UK")
	}
	if ent.Attributes[2].Comment != "User email" {
		t.Errorf("attr[2] comment = %q, want %q", ent.Attributes[2].Comment, "User email")
	}
}

func TestParseERCardinality(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		startDec   ir.EdgeDecoration
		endDec     ir.EdgeDecoration
	}{
		{"one-to-many", "erDiagram\n    A ||--o{ B : has", ir.DecCrowsFootOne, ir.DecCrowsFootZeroMany},
		{"one-to-one", "erDiagram\n    A ||--|| B : is", ir.DecCrowsFootOne, ir.DecCrowsFootOne},
		{"zero-one-to-many", "erDiagram\n    A |o--o{ B : has", ir.DecCrowsFootZeroOne, ir.DecCrowsFootZeroMany},
		{"many-to-many", "erDiagram\n    A }|--|{ B : has", ir.DecCrowsFootMany, ir.DecCrowsFootMany},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error: %v", err)
			}
			if len(out.Graph.Edges) != 1 {
				t.Fatalf("edges = %d, want 1", len(out.Graph.Edges))
			}
			e := out.Graph.Edges[0]
			if e.StartDecoration == nil || *e.StartDecoration != tt.startDec {
				t.Errorf("start decoration = %v, want %v", e.StartDecoration, tt.startDec)
			}
			if e.EndDecoration == nil || *e.EndDecoration != tt.endDec {
				t.Errorf("end decoration = %v, want %v", e.EndDecoration, tt.endDec)
			}
		})
	}
}

func TestParseERLineStyle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		style ir.EdgeStyle
	}{
		{"solid", "erDiagram\n    A ||--|| B : is", ir.Solid},
		{"dashed", "erDiagram\n    A ||..|| B : refs", ir.Dotted},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error: %v", err)
			}
			if out.Graph.Edges[0].Style != tt.style {
				t.Errorf("style = %v, want %v", out.Graph.Edges[0].Style, tt.style)
			}
		})
	}
}

func TestParseERLabel(t *testing.T) {
	input := `erDiagram
    CUSTOMER ||--o{ ORDER : places`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.Edges[0].Label == nil || *out.Graph.Edges[0].Label != "places" {
		t.Errorf("label = %v, want 'places'", out.Graph.Edges[0].Label)
	}
}

func TestParseERAlias(t *testing.T) {
	input := `erDiagram
    p["Person"] {
        string firstName
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	ent, ok := out.Graph.Entities["p"]
	if !ok {
		t.Fatal("p not in Entities")
	}
	if ent.Label != "Person" {
		t.Errorf("label = %q, want %q", ent.Label, "Person")
	}
}

func TestParseERCompositeKey(t *testing.T) {
	input := `erDiagram
    T {
        int id PK,FK
    }`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	attr := out.Graph.Entities["T"].Attributes[0]
	if len(attr.Keys) != 2 {
		t.Fatalf("keys = %d, want 2", len(attr.Keys))
	}
	if attr.Keys[0] != ir.KeyPrimary || attr.Keys[1] != ir.KeyForeign {
		t.Errorf("keys = %v, want [PK, FK]", attr.Keys)
	}
}

func TestParseERDirection(t *testing.T) {
	input := `erDiagram
    direction LR
    A ||--|| B : is`
	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if out.Graph.Direction != ir.LeftRight {
		t.Errorf("Direction = %v, want LeftRight", out.Graph.Direction)
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./parser/ -run TestParseER -v`
Expected: FAIL — `parseER` not found

**Step 3: Write implementation**

Create `parser/er.go` implementing `parseER(input string) (*ParseOutput, error)`. Key pieces:

- Header regex: `(?i)^erdiagram`
- Relationship regex matching cardinality markers on both sides: `(|||\|o|o\||}\||}\||\|{|o{|}o)\s*(--|\.\.)\s*(|||\|o|o\||}\||\|{|o{|}o)`
- Map cardinality tokens to `ir.EdgeDecoration` values
- Entity body parsing (brace-delimited): `type name [keys] ["comment"]`
- Entity alias: `name["Display"]` or `name[Alias]`
- Composite keys: `PK,FK` parsed by splitting on comma
- Direction and style directives
- Wire into parser dispatch

**Step 4: Run tests to verify they pass**

Run: `go test ./parser/ -run TestParseER -v`
Expected: ALL PASS

**Step 5: Commit**

```
git add parser/er.go parser/er_test.go parser/parser.go
git commit -m "feat(parser): add ER diagram parser with cardinality support"
```

---

### Task 9: Layout — Class Diagram Sizing + DiagramData

**Files:**
- Create: `layout/class.go`
- Create: `layout/class_test.go`
- Modify: `layout/types.go:64-67` (add ClassData type)
- Modify: `layout/layout.go:13-19` (add Class case)

**Step 1: Write the failing tests**

Create `layout/class_test.go`:

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestComputeClassLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Class
	g.Direction = ir.TopDown

	g.EnsureNode("Animal", nil, nil)
	g.EnsureNode("Dog", nil, nil)
	g.Members["Animal"] = &ir.ClassMembers{
		Attributes: []ir.ClassMember{
			{Name: "name", Type: "String", Visibility: ir.VisPublic},
		},
		Methods: []ir.ClassMember{
			{Name: "speak", IsMethod: true, Visibility: ir.VisPublic, Type: "void"},
		},
	}
	g.Edges = append(g.Edges, &ir.Edge{From: "Dog", To: "Animal", Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Class {
		t.Errorf("Kind = %v, want Class", l.Kind)
	}
	if len(l.Nodes) != 2 {
		t.Errorf("nodes = %d, want 2", len(l.Nodes))
	}
	if len(l.Edges) != 1 {
		t.Errorf("edges = %d, want 1", len(l.Edges))
	}
	// Animal should be taller than Dog (it has members).
	animal := l.Nodes["Animal"]
	dog := l.Nodes["Dog"]
	if animal.Height <= dog.Height {
		t.Errorf("Animal height (%f) should be > Dog height (%f)", animal.Height, dog.Height)
	}
	// Verify ClassData is returned.
	if _, ok := l.Diagram.(ClassData); !ok {
		t.Errorf("Diagram data type = %T, want ClassData", l.Diagram)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestComputeClassLayout -v`
Expected: FAIL — `ClassData` undefined

**Step 3: Write implementation**

Add to `layout/types.go`:

```go
// ClassData holds class-diagram-specific layout data.
type ClassData struct {
	Compartments map[string]ClassCompartment // node ID -> compartment dims
}

func (ClassData) diagramData() {}

// ClassCompartment stores section dimensions for UML class boxes.
type ClassCompartment struct {
	HeaderHeight    float32
	AttributeHeight float32
	MethodHeight    float32
}
```

Create `layout/class.go`:

```go
package layout

// computeClassLayout lays out a class diagram using the Sugiyama pipeline
// with class-specific node sizing for UML compartment boxes.
func computeClassLayout(...) *Layout { ... }
```

This function:
1. Creates a `textmetrics.Measurer`
2. Calls `sizeClassNodes()` — measures UML compartment boxes (header + attributes + methods)
3. Reuses `computeRanks()`, `orderRankNodes()`, `positionNodes()`, `routeEdges()`
4. Returns `Layout` with `ClassData`

The `sizeClassNodes()` function measures:
- Header: class name (+ annotation if present), full width
- Attribute section: each attribute line (visibility symbol + type + name)
- Method section: each method line (visibility symbol + name + params + return type)
- Width: max of all sections + 2 * CompartmentPadX
- Height: sum of sections + divider heights + 2 * CompartmentPadY

Wire into `layout/layout.go` switch:
```go
case ir.Class:
    return computeClassLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./layout/ -run TestComputeClassLayout -v`
Expected: PASS

**Step 5: Commit**

```
git add layout/class.go layout/class_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add class diagram layout with UML compartment sizing"
```

---

### Task 10: Layout — ER Diagram Sizing + DiagramData

**Files:**
- Create: `layout/er.go`
- Create: `layout/er_test.go`
- Modify: `layout/types.go` (add ERData type)
- Modify: `layout/layout.go` (add Er case)

**Step 1: Write the failing test**

Create `layout/er_test.go`:

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestComputeERLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Er
	g.Direction = ir.TopDown

	g.EnsureNode("CUSTOMER", nil, nil)
	g.EnsureNode("ORDER", nil, nil)
	g.Entities["CUSTOMER"] = &ir.Entity{
		ID: "CUSTOMER",
		Attributes: []ir.EntityAttribute{
			{Type: "string", Name: "name"},
			{Type: "int", Name: "id", Keys: []ir.AttributeKey{ir.KeyPrimary}},
		},
	}
	g.Entities["ORDER"] = &ir.Entity{
		ID: "ORDER",
		Attributes: []ir.EntityAttribute{
			{Type: "int", Name: "id", Keys: []ir.AttributeKey{ir.KeyPrimary}},
		},
	}
	g.Edges = append(g.Edges, &ir.Edge{From: "CUSTOMER", To: "ORDER"})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Er {
		t.Errorf("Kind = %v, want Er", l.Kind)
	}
	if len(l.Nodes) != 2 {
		t.Errorf("nodes = %d, want 2", len(l.Nodes))
	}
	// CUSTOMER should be taller (more attributes).
	cust := l.Nodes["CUSTOMER"]
	order := l.Nodes["ORDER"]
	if cust.Height <= order.Height {
		t.Errorf("CUSTOMER height (%f) should be > ORDER height (%f)", cust.Height, order.Height)
	}
	if _, ok := l.Diagram.(ERData); !ok {
		t.Errorf("Diagram data type = %T, want ERData", l.Diagram)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestComputeERLayout -v`
Expected: FAIL — `ERData` undefined

**Step 3: Write implementation**

Add to `layout/types.go`:

```go
// ERData holds ER-diagram-specific layout data.
type ERData struct {
	EntityDims map[string]EntityDimensions
}

func (ERData) diagramData() {}

// EntityDimensions stores column widths for entity rendering.
type EntityDimensions struct {
	TypeColWidth float32
	NameColWidth float32
	KeyColWidth  float32
	HeaderHeight float32
	RowCount     int
}
```

Create `layout/er.go`:
- `computeERLayout()` — same pipeline as class, but with `sizeERNodes()`
- `sizeERNodes()` measures entity boxes: header row + attribute rows with column widths
- Wire into dispatch switch

**Step 4: Run test to verify it passes**

Run: `go test ./layout/ -run TestComputeERLayout -v`
Expected: PASS

**Step 5: Commit**

```
git add layout/er.go layout/er_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add ER diagram layout with entity box sizing"
```

---

### Task 11: Layout — State Diagram (Recursive)

**Files:**
- Create: `layout/state.go`
- Create: `layout/state_test.go`
- Modify: `layout/types.go` (add StateData type)
- Modify: `layout/layout.go` (add State case)

**Step 1: Write the failing tests**

Create `layout/state_test.go`:

```go
package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestComputeStateLayoutSimple(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown

	g.EnsureNode("__start__", nil, nil)
	g.EnsureNode("First", nil, nil)
	g.EnsureNode("__end__", nil, nil)
	g.Edges = append(g.Edges,
		&ir.Edge{From: "__start__", To: "First", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "First", To: "__end__", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.State {
		t.Errorf("Kind = %v, want State", l.Kind)
	}
	if len(l.Nodes) != 3 {
		t.Errorf("nodes = %d, want 3", len(l.Nodes))
	}
	if _, ok := l.Diagram.(StateData); !ok {
		t.Errorf("Diagram data type = %T, want StateData", l.Diagram)
	}
}

func TestComputeStateLayoutComposite(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown

	g.EnsureNode("Outer", nil, nil)
	inner := ir.NewGraph()
	inner.Kind = ir.State
	inner.EnsureNode("__start__", nil, nil)
	inner.EnsureNode("inner1", nil, nil)
	inner.Edges = append(inner.Edges, &ir.Edge{From: "__start__", To: "inner1", Directed: true, ArrowEnd: true})
	g.CompositeStates["Outer"] = &ir.CompositeState{
		ID:    "Outer",
		Label: "Outer",
		Inner: inner,
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	outer := l.Nodes["Outer"]
	if outer == nil {
		t.Fatal("Outer node not in layout")
	}
	// Composite should be bigger than a simple state.
	if outer.Width < 100 {
		t.Errorf("Outer width = %f, expected > 100 for composite", outer.Width)
	}
	sd, ok := l.Diagram.(StateData)
	if !ok {
		t.Fatal("expected StateData")
	}
	if sd.InnerLayouts["Outer"] == nil {
		t.Error("expected inner layout for Outer")
	}
}

func TestComputeStateLayoutAnnotations(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown

	forkShape := ir.ForkJoin
	diamondShape := ir.Diamond
	g.EnsureNode("fork1", nil, &forkShape)
	g.EnsureNode("choice1", nil, &diamondShape)
	g.StateAnnotations["fork1"] = ir.StateFork
	g.StateAnnotations["choice1"] = ir.StateChoice

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	fork := l.Nodes["fork1"]
	if fork == nil {
		t.Fatal("fork1 not in layout")
	}
	if fork.Shape != ir.ForkJoin {
		t.Errorf("fork shape = %v, want ForkJoin", fork.Shape)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestComputeStateLayout -v`
Expected: FAIL — `StateData` undefined

**Step 3: Write implementation**

Add to `layout/types.go`:

```go
// StateData holds state-diagram-specific layout data.
type StateData struct {
	InnerLayouts map[string]*Layout // composite state ID -> inner layout
}

func (StateData) diagramData() {}
```

Create `layout/state.go`:

`computeStateLayout()`:
1. Identify special node shapes: `__start__` -> small circle, `__end__` -> bullseye, fork/join annotations -> bar, choice -> diamond
2. For each composite state, recursively call `computeStateLayout()` on `Inner` graph
3. Size composite node to contain its inner layout + padding + label height
4. For concurrent regions, layout each region, stack vertically, compute total size
5. Run Sugiyama pipeline on top-level with sized nodes
6. Translate inner layout coordinates relative to composite position
7. Return `Layout` with `StateData` containing inner layouts

Wire into dispatch switch.

**Step 4: Run tests to verify they pass**

Run: `go test ./layout/ -run TestComputeStateLayout -v`
Expected: ALL PASS

**Step 5: Commit**

```
git add layout/state.go layout/state_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add state diagram recursive layout"
```

---

### Task 12: Renderer — SVG Markers for Class/ER

**Files:**
- Modify: `render/svg.go:52-77` (add new marker definitions to renderDefs)

**Step 1: Write the failing test**

Add to `render/svg_test.go`:

```go
func TestRenderDefsHasClassMarkers(t *testing.T) {
	th := theme.Modern()
	var b svgBuilder
	renderDefs(&b, th)
	svg := b.String()

	markers := []string{
		"marker-closed-triangle",
		"marker-filled-diamond",
		"marker-open-diamond",
	}
	for _, id := range markers {
		if !strings.Contains(svg, id) {
			t.Errorf("missing marker definition: %s", id)
		}
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./render/ -run TestRenderDefsHasClassMarkers -v`
Expected: FAIL — markers not found

**Step 3: Write implementation**

Add to `renderDefs()` in `render/svg.go`:

```go
// Closed triangle (inheritance/realization) — forward
b.raw(`<marker id="marker-closed-triangle" viewBox="0 0 20 20" refX="18" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
b.selfClose("path", "d", "M 0 0 L 20 10 L 0 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
b.closeTag("marker")

// Closed triangle — reverse
b.raw(`<marker id="marker-closed-triangle-start" viewBox="0 0 20 20" refX="2" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
b.selfClose("path", "d", "M 20 0 L 0 10 L 20 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
b.closeTag("marker")

// Filled diamond (composition) — forward
b.raw(`<marker id="marker-filled-diamond" viewBox="0 0 20 20" refX="18" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
b.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.LineColor, "stroke", th.LineColor, "stroke-width", "1")
b.closeTag("marker")

// Filled diamond — reverse
b.raw(`<marker id="marker-filled-diamond-start" viewBox="0 0 20 20" refX="2" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
b.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.LineColor, "stroke", th.LineColor, "stroke-width", "1")
b.closeTag("marker")

// Open diamond (aggregation) — forward
b.raw(`<marker id="marker-open-diamond" viewBox="0 0 20 20" refX="18" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
b.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
b.closeTag("marker")

// Open diamond — reverse
b.raw(`<marker id="marker-open-diamond-start" viewBox="0 0 20 20" refX="2" refY="10" markerUnits="userSpaceOnUse" markerWidth="12" markerHeight="12" orient="auto">`)
b.selfClose("path", "d", "M 0 10 L 10 0 L 20 10 L 10 20 z", "fill", th.Background, "stroke", th.LineColor, "stroke-width", "1")
b.closeTag("marker")
```

**Step 4: Run test to verify it passes**

Run: `go test ./render/ -run TestRenderDefsHasClassMarkers -v`
Expected: PASS

**Step 5: Commit**

```
git add render/svg.go render/svg_test.go
git commit -m "feat(render): add SVG marker definitions for class/ER relationships"
```

---

### Task 13: Renderer — Class Diagram

**Files:**
- Create: `render/class.go`
- Create: `render/class_test.go`
- Modify: `render/svg.go:40-46` (add ClassData case)

**Step 1: Write the failing tests**

Create `render/class_test.go`:

```go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderClassCompartments(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Class
	g.Direction = ir.TopDown
	g.EnsureNode("Animal", nil, nil)
	g.Members["Animal"] = &ir.ClassMembers{
		Attributes: []ir.ClassMember{
			{Name: "name", Type: "String", Visibility: ir.VisPublic},
		},
		Methods: []ir.ClassMember{
			{Name: "speak", IsMethod: true, Visibility: ir.VisPublic, Type: "void"},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Should contain the class name.
	if !strings.Contains(svg, "Animal") {
		t.Error("missing class name 'Animal'")
	}
	// Should contain visibility symbols.
	if !strings.Contains(svg, "+") {
		t.Error("missing visibility symbol '+'")
	}
	// Should have divider lines (rendered as <line> elements).
	if strings.Count(svg, "<line") < 2 {
		t.Errorf("expected at least 2 divider lines, got %d", strings.Count(svg, "<line"))
	}
}

func TestRenderClassRelationshipMarkers(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Class
	g.Direction = ir.TopDown
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	closedTri := ir.ClosedTriangle
	g.Edges = append(g.Edges, &ir.Edge{
		From: "A", To: "B", Directed: true, ArrowEnd: true,
		ArrowEndKind: &closedTri,
	})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "marker-closed-triangle") {
		t.Error("missing closed triangle marker reference")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./render/ -run TestRenderClass -v`
Expected: FAIL — ClassData case not handled

**Step 3: Write implementation**

Create `render/class.go`:

- `renderClass(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout)`:
  1. Render namespace backgrounds (like subgraphs)
  2. Render edges with appropriate markers based on `ArrowEndKind`/`ArrowStartKind`
  3. Render class nodes as compartment boxes:
     - Header section: class name centered, annotation above in `<<guillemets>>`
     - Horizontal divider line
     - Attributes section: left-aligned, visibility prefix + type + name
     - Horizontal divider line
     - Methods section: left-aligned, visibility prefix + name + (params) + return type

- `renderClassNode(b *svgBuilder, n *layout.NodeLayout, members *ir.ClassMembers, comp layout.ClassCompartment, th *theme.Theme, cfg *config.Layout)`:
  - Draw outer rounded rect
  - Draw header background rect
  - Draw header text (class name)
  - Draw divider line
  - Draw each attribute line
  - Draw divider line
  - Draw each method line

- `renderClassEdges()` — extends edge rendering to use relationship-specific markers

Wire into `render/svg.go` switch:
```go
case layout.ClassData:
    renderClass(b, l, th, cfg)
```

**Step 4: Run tests to verify they pass**

Run: `go test ./render/ -run TestRenderClass -v`
Expected: ALL PASS

**Step 5: Commit**

```
git add render/class.go render/class_test.go render/svg.go
git commit -m "feat(render): add class diagram renderer with UML compartments"
```

---

### Task 14: Renderer — State Diagram

**Files:**
- Create: `render/state.go`
- Create: `render/state_test.go`
- Modify: `render/svg.go` (add StateData case)

**Step 1: Write the failing tests**

Create `render/state_test.go`:

```go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderStateSimple(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown
	g.EnsureNode("__start__", nil, nil)
	g.EnsureNode("First", nil, nil)
	g.EnsureNode("__end__", nil, nil)
	g.StateDescriptions["First"] = "First state"
	g.Edges = append(g.Edges,
		&ir.Edge{From: "__start__", To: "First", Directed: true, ArrowEnd: true},
		&ir.Edge{From: "First", To: "__end__", Directed: true, ArrowEnd: true},
	)

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Start state should be a filled circle.
	if !strings.Contains(svg, "<circle") {
		t.Error("missing circle element for start/end state")
	}
	// Should contain state label.
	if !strings.Contains(svg, "First") {
		t.Error("missing state label 'First'")
	}
}

func TestRenderStateComposite(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.State
	g.Direction = ir.TopDown
	g.EnsureNode("Outer", nil, nil)

	inner := ir.NewGraph()
	inner.Kind = ir.State
	inner.EnsureNode("__start__", nil, nil)
	inner.EnsureNode("inner1", nil, nil)
	inner.Edges = append(inner.Edges, &ir.Edge{From: "__start__", To: "inner1", Directed: true, ArrowEnd: true})
	g.CompositeStates["Outer"] = &ir.CompositeState{ID: "Outer", Label: "Outer", Inner: inner}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	// Composite renders a larger container rect + inner content.
	if !strings.Contains(svg, "Outer") {
		t.Error("missing composite label 'Outer'")
	}
	if !strings.Contains(svg, "inner1") {
		t.Error("missing inner state 'inner1'")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./render/ -run TestRenderState -v`
Expected: FAIL — StateData case not handled

**Step 3: Write implementation**

Create `render/state.go`:

- `renderState(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout)`:
  1. Render edges (transitions) with arrows
  2. Render nodes based on type:
     - `__start__`: filled black circle (small, fixed radius)
     - `__end__`: filled black circle with outer ring (bullseye)
     - Fork/join annotated: wide horizontal bar
     - Choice annotated: diamond
     - Composite: rounded rect container with label, recursively render inner layout
     - Regular state: rounded rect with optional description divider
  3. Render notes

- `renderStateNode()`: dispatch based on node type
- `renderStartState()`: filled circle
- `renderEndState()`: bullseye (filled circle + outer ring)
- `renderCompositeState()`: container rect + recursive inner render
- `renderRegularState()`: rounded rect with optional description
- `renderConcurrentRegions()`: horizontal dashed divider lines

Wire into dispatch.

**Step 4: Run tests to verify they pass**

Run: `go test ./render/ -run TestRenderState -v`
Expected: ALL PASS

**Step 5: Commit**

```
git add render/state.go render/state_test.go render/svg.go
git commit -m "feat(render): add state diagram renderer with composite states"
```

---

### Task 15: Renderer — ER Diagram

**Files:**
- Create: `render/er.go`
- Create: `render/er_test.go`
- Modify: `render/svg.go` (add ERData case)

**Step 1: Write the failing tests**

Create `render/er_test.go`:

```go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderEREntities(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Er
	g.Direction = ir.TopDown
	g.EnsureNode("CUSTOMER", nil, nil)
	g.Entities["CUSTOMER"] = &ir.Entity{
		ID: "CUSTOMER",
		Attributes: []ir.EntityAttribute{
			{Type: "string", Name: "name"},
			{Type: "int", Name: "id", Keys: []ir.AttributeKey{ir.KeyPrimary}},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "CUSTOMER") {
		t.Error("missing entity name 'CUSTOMER'")
	}
	if !strings.Contains(svg, "PK") {
		t.Error("missing key constraint 'PK'")
	}
	// Should have attribute rows.
	if !strings.Contains(svg, "name") {
		t.Error("missing attribute 'name'")
	}
}

func TestRenderERRelationship(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Er
	g.Direction = ir.TopDown
	g.EnsureNode("A", nil, nil)
	g.EnsureNode("B", nil, nil)
	g.Entities["A"] = &ir.Entity{ID: "A"}
	g.Entities["B"] = &ir.Entity{ID: "B"}

	startDec := ir.DecCrowsFootOne
	endDec := ir.DecCrowsFootZeroMany
	label := "has"
	g.Edges = append(g.Edges, &ir.Edge{
		From: "A", To: "B",
		StartDecoration: &startDec,
		EndDecoration:   &endDec,
		Label:           &label,
	})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "has") {
		t.Error("missing relationship label 'has'")
	}
	// Should have edge paths.
	if !strings.Contains(svg, "edgePath") {
		t.Error("missing edge path")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./render/ -run TestRenderER -v`
Expected: FAIL — ERData case not handled

**Step 3: Write implementation**

Create `render/er.go`:

- `renderER(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout)`:
  1. Render edges with crow's foot decorations at endpoints
  2. Render entity boxes:
     - Header rect (filled) with entity name
     - Attribute rows: type | name | key columns, with alternating background
     - Row separator lines
  3. Render edge labels at midpoints

- `renderEntityBox()`: header + attribute rows
- `renderCrowsFoot()`: draw crow's foot decoration SVG paths at edge endpoints
  - `||` exactly one: two short perpendicular lines
  - `o|`/`|o` zero or one: circle + perpendicular line
  - `|{`/`}|` one or more: perpendicular line + 3-line fork
  - `o{`/`}o` zero or more: circle + 3-line fork

Wire into dispatch.

**Step 4: Run tests to verify they pass**

Run: `go test ./render/ -run TestRenderER -v`
Expected: ALL PASS

**Step 5: Commit**

```
git add render/er.go render/er_test.go render/svg.go
git commit -m "feat(render): add ER diagram renderer with crow's foot notation"
```

---

### Task 16: Integration Tests & Fixtures

**Files:**
- Create: `testdata/fixtures/class-simple.mmd`
- Create: `testdata/fixtures/class-relationships.mmd`
- Create: `testdata/fixtures/state-simple.mmd`
- Create: `testdata/fixtures/state-composite.mmd`
- Create: `testdata/fixtures/er-simple.mmd`
- Create: `testdata/fixtures/er-attributes.mmd`
- Modify: `mermaid_test.go` (add integration tests)

**Step 1: Create test fixtures**

`testdata/fixtures/class-simple.mmd`:
```
classDiagram
    class Animal {
        +String name
        +int age
        +isMammal() bool
        +mate()
    }
    class Dog {
        +String breed
        +bark() void
    }
    Animal <|-- Dog
```

`testdata/fixtures/class-relationships.mmd`:
```
classDiagram
    Animal <|-- Dog : extends
    Car *-- Engine
    Library o-- Book
    Student --> Course
    Class1 ..> Class2
    Interface1 ..|> Impl1
```

`testdata/fixtures/state-simple.mmd`:
```
stateDiagram-v2
    [*] --> Still
    Still --> Moving
    Moving --> Crash
    Crash --> [*]
```

`testdata/fixtures/state-composite.mmd`:
```
stateDiagram-v2
    [*] --> First
    state First {
        [*] --> Second
        Second --> [*]
    }
    First --> Third
    Third --> [*]
```

`testdata/fixtures/er-simple.mmd`:
```
erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE-ITEM : contains
```

`testdata/fixtures/er-attributes.mmd`:
```
erDiagram
    CUSTOMER {
        string name
        int custNumber PK
        string sector
    }
    ORDER {
        int orderNumber PK
        string deliveryAddress
    }
    CUSTOMER ||--o{ ORDER : places
```

**Step 2: Write failing integration tests**

Add to `mermaid_test.go`:

```go
func TestRenderClassDiagram(t *testing.T) {
	input := `classDiagram
    class Animal {
        +String name
        +isMammal() bool
    }
    class Dog {
        +String breed
    }
    Animal <|-- Dog`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
	for _, label := range []string{"Animal", "Dog", "name", "breed"} {
		if !strings.Contains(svg, label) {
			t.Errorf("missing label %q", label)
		}
	}
}

func TestRenderStateDiagram(t *testing.T) {
	input := `stateDiagram-v2
    [*] --> Still
    Still --> Moving
    Moving --> Crash
    Crash --> [*]`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
	for _, label := range []string{"Still", "Moving", "Crash"} {
		if !strings.Contains(svg, label) {
			t.Errorf("missing label %q", label)
		}
	}
}

func TestRenderERDiagram(t *testing.T) {
	input := `erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE-ITEM : contains`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
	for _, label := range []string{"CUSTOMER", "ORDER", "LINE-ITEM"} {
		if !strings.Contains(svg, label) {
			t.Errorf("missing label %q", label)
		}
	}
}

func TestGoldenClassRelationships(t *testing.T) {
	input := `classDiagram
    Animal <|-- Dog : extends
    Car *-- Engine
    Library o-- Book
    Student --> Course
    Class1 ..> Class2
    Interface1 ..|> Impl1`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Count(svg, "edgePath") < 6 {
		t.Errorf("expected at least 6 edge paths, got %d", strings.Count(svg, "edgePath"))
	}
}

func TestGoldenStateComposite(t *testing.T) {
	input := `stateDiagram-v2
    [*] --> First
    state First {
        [*] --> Second
        Second --> [*]
    }
    First --> Third
    Third --> [*]`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "First") {
		t.Error("missing composite state 'First'")
	}
	if !strings.Contains(svg, "Second") {
		t.Error("missing inner state 'Second'")
	}
}

func TestGoldenERAttributes(t *testing.T) {
	input := `erDiagram
    CUSTOMER {
        string name
        int custNumber PK
    }
    ORDER {
        int orderNumber PK
    }
    CUSTOMER ||--o{ ORDER : places`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "PK") {
		t.Error("missing PK annotation")
	}
	if !strings.Contains(svg, "places") {
		t.Error("missing relationship label 'places'")
	}
}
```

**Step 3: Run all integration tests**

Run: `go test ./... -v`
Expected: ALL PASS (all integration tests pass end-to-end)

**Step 4: Commit**

```
git add testdata/fixtures/ mermaid_test.go
git commit -m "test: add integration tests and fixtures for class/state/ER diagrams"
```

---

### Task 17: Benchmarks

**Files:**
- Modify: `mermaid_bench_test.go` (add benchmarks for new diagram types)

**Step 1: Add benchmarks**

Add to `mermaid_bench_test.go`:

```go
func BenchmarkRenderClassSimple(b *testing.B) {
	input := `classDiagram
    class Animal {
        +String name
        +int age
        +isMammal() bool
    }
    class Dog {
        +String breed
        +bark() void
    }
    Animal <|-- Dog`
	for b.Loop() {
		_, _ = Render(input)
	}
}

func BenchmarkRenderStateDiagram(b *testing.B) {
	input := `stateDiagram-v2
    [*] --> Still
    Still --> Moving
    Moving --> Crash
    Crash --> [*]`
	for b.Loop() {
		_, _ = Render(input)
	}
}

func BenchmarkRenderERDiagram(b *testing.B) {
	input := `erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE-ITEM : contains
    CUSTOMER {
        string name
        int id PK
    }`
	for b.Loop() {
		_, _ = Render(input)
	}
}
```

**Step 2: Run benchmarks**

Run: `go test -bench=. -benchmem ./...`
Expected: Benchmarks run successfully with timing output

**Step 3: Commit**

```
git add mermaid_bench_test.go
git commit -m "bench: add benchmarks for class, state, and ER diagrams"
```

---

### Task 18: Final Validation

**Step 1: Run all tests**

Run: `go test ./... -v -count=1`
Expected: ALL PASS across all packages

**Step 2: Run linter**

Run: `golangci-lint run ./...`
Expected: No errors

**Step 3: Run benchmarks**

Run: `go test -bench=. -benchmem ./...`
Expected: All benchmarks complete

**Step 4: Manual smoke test**

Test each diagram type via Go test:
```go
// Quick manual verification of SVG output structure
go test -run TestRenderClassDiagram -v
go test -run TestRenderStateDiagram -v
go test -run TestRenderERDiagram -v
```

**Step 5: Commit any fixes and push**

```
git push origin <branch>
```
