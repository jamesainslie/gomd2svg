# Phase 10: Journey & Architecture Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Journey (user journey map) and Architecture (architecture-beta) diagram support to mermaid-go.

**Architecture:** Both diagrams use custom positioning (not Sugiyama). Journey uses a horizontal score track; Architecture uses constraint-based grid placement from directional edge hints. Follow the established per-diagram pattern: IR types → config/theme → parser → layout → renderer → integration tests.

**Tech Stack:** Go, stdlib regexp, existing textmetrics/theme/config packages.

---

### Task 1: Journey IR Types

**Files:**
- Create: `ir/journey.go`
- Create: `ir/journey_test.go`
- Modify: `ir/graph.go` — add Journey fields

**Step 1: Create ir/journey.go**

```go
package ir

// JourneyTask represents a single task in a journey diagram.
type JourneyTask struct {
	Name    string
	Score   int      // 1-5 satisfaction score
	Actors  []string // participating actors
	Section string   // section name this task belongs to
}

// JourneySection represents a named section grouping tasks.
type JourneySection struct {
	Name  string
	Tasks []int // indices into Graph.JourneyTasks
}
```

**Step 2: Create ir/journey_test.go**

Table-driven test verifying JourneyTask and JourneySection struct creation and field access.

**Step 3: Add Graph fields in ir/graph.go**

Add after the C4 fields block:
```go
// Journey diagram fields
JourneyTitle    string
JourneyTasks    []*JourneyTask
JourneySections []*JourneySection
```

**Step 4: Verify**

Run: `go build ./ir/... && go test ./ir/... -v`

---

### Task 2: Architecture IR Types

**Files:**
- Create: `ir/architecture.go`
- Create: `ir/architecture_test.go`
- Modify: `ir/graph.go` — add Architecture fields

**Step 1: Create ir/architecture.go**

```go
package ir

// ArchSide represents a connection side on a service or junction.
type ArchSide int

const (
	ArchLeft ArchSide = iota
	ArchRight
	ArchTop
	ArchBottom
)

func (s ArchSide) String() string { ... }

// ArchService represents a service node in an architecture diagram.
type ArchService struct {
	ID      string
	Label   string
	Icon    string
	GroupID string // empty if top-level
}

// ArchGroup represents a grouping container.
type ArchGroup struct {
	ID       string
	Label    string
	Icon     string
	ParentID string   // for nested groups
	Children []string // service/junction IDs
}

// ArchJunction is a connection point between edges.
type ArchJunction struct {
	ID      string
	GroupID string
}

// ArchEdge represents a connection between services/junctions.
type ArchEdge struct {
	FromID     string
	FromSide   ArchSide
	ToID       string
	ToSide     ArchSide
	ArrowLeft  bool
	ArrowRight bool
}
```

**Step 2: Create ir/architecture_test.go**

Table-driven test for ArchSide.String() and struct creation.

**Step 3: Add Graph fields in ir/graph.go**

```go
// Architecture diagram fields
ArchServices  []*ArchService
ArchGroups    []*ArchGroup
ArchJunctions []*ArchJunction
ArchEdges     []*ArchEdge
```

**Step 4: Verify**

Run: `go build ./ir/... && go test ./ir/... -v`

---

### Task 3: Config and Theme

**Files:**
- Modify: `config/config.go` — add JourneyConfig, ArchitectureConfig
- Modify: `theme/theme.go` — add Journey/Architecture colors

**Step 1: Add config types and defaults**

Add to Layout struct:
```go
Journey      JourneyConfig
Architecture ArchitectureConfig
```

Config structs:
```go
type JourneyConfig struct {
	TaskWidth   float32
	TaskHeight  float32
	TaskSpacing float32
	TrackHeight float32
	SectionGap  float32
	PaddingX    float32
	PaddingY    float32
}

type ArchitectureConfig struct {
	ServiceWidth  float32
	ServiceHeight float32
	GroupPadding  float32
	JunctionSize  float32
	ColumnGap     float32
	RowGap        float32
	PaddingX      float32
	PaddingY      float32
}
```

Defaults: TaskWidth=120, TaskHeight=50, TaskSpacing=20, TrackHeight=200, SectionGap=10, PaddingX=30, PaddingY=40. ServiceWidth=120, ServiceHeight=80, GroupPadding=30, JunctionSize=10, ColumnGap=60, RowGap=60, PaddingX=30, PaddingY=30.

**Step 2: Add theme fields**

Journey:
- `JourneySectionColors []string` — cycling fill colors
- `JourneyTaskFill`, `JourneyTaskBorder`, `JourneyTaskText` string
- `JourneyScoreColors [5]string` — red to green for scores 1-5

Architecture:
- `ArchServiceFill`, `ArchServiceBorder`, `ArchServiceText` string
- `ArchGroupFill`, `ArchGroupBorder`, `ArchGroupText` string
- `ArchEdgeColor`, `ArchJunctionFill` string

Initialize in both Modern() and MermaidDefault().

**Step 3: Add config_test.go assertion**

Verify DefaultLayout().Journey and DefaultLayout().Architecture have non-zero defaults.

**Step 4: Add theme_test.go assertion**

Verify Modern() and MermaidDefault() have non-empty Journey/Architecture color fields.

**Step 5: Verify**

Run: `go test ./config/... ./theme/... -v`

---

### Task 4: Journey Parser

**Files:**
- Create: `parser/journey.go`
- Create: `parser/journey_test.go`
- Modify: `parser/parser.go` — add Journey dispatch case

**Step 1: Create parser/journey.go**

```go
func parseJourney(input string) (*ParseOutput, error)
```

Parser logic:
1. `preprocessInput()` to strip comments/blanks
2. Skip the `journey` keyword line
3. Parse `title <text>` → sets JourneyTitle
4. Parse `section <name>` → creates new JourneySection
5. Parse `<name>: <score>: <actor1>, <actor2>` → creates JourneyTask
6. Task regex: `^\s*(.+?):\s*(\d+)\s*(?::\s*(.*))?$`
7. Score parsed as int, actors split by comma and trimmed

**Step 2: Create parser/journey_test.go**

Tests:
- TestParseJourney: full example with title, 2 sections, tasks with scores and actors
- TestParseJourneyMinimal: single task, no section
- TestParseJourneyNoActors: tasks with empty actor lists

**Step 3: Add dispatch in parser/parser.go**

```go
case ir.Journey:
    return parseJourney(input)
```

**Step 4: Verify**

Run: `go test ./parser/... -v -run Journey`

---

### Task 5: Architecture Parser

**Files:**
- Create: `parser/architecture.go`
- Create: `parser/architecture_test.go`
- Modify: `parser/parser.go` — add Architecture dispatch case

**Step 1: Create parser/architecture.go**

```go
func parseArchitecture(input string) (*ParseOutput, error)
```

Parser logic:
1. `preprocessInput()` to strip comments/blanks
2. Skip `architecture-beta` keyword line
3. Parse `group <id>(<icon>)[<label>]( in <parent>)?` → ArchGroup
4. Parse `service <id>(<icon>)[<label>]( in <parent>)?` → ArchService, also create Node
5. Parse `junction <id>( in <parent>)?` → ArchJunction, also create Node
6. Parse `<id>:<side> <arrow> <side>:<id>` → ArchEdge
7. Side regex for L/R/T/B, arrow patterns: `--`, `-->`, `<--`, `<-->`
8. Group regex: `^group\s+(\w+)(?:\((\w+)\))?\[([^\]]+)\](?:\s+in\s+(\w+))?$`
9. Service regex: `^service\s+(\w+)(?:\(([^)]+)\))?\[([^\]]+)\](?:\s+in\s+(\w+))?$`
10. Edge regex: `^(\w+)(?:\{group\})?:(L|R|T|B)\s*(<)?--(>)?\s*(L|R|T|B):(\w+)(?:\{group\})?$`

**Step 2: Create parser/architecture_test.go**

Tests:
- TestParseArchitecture: groups, services, junctions, edges
- TestParseArchitectureMinimal: two services with one edge
- TestParseArchitectureNested: nested groups with `in` keyword

**Step 3: Add dispatch in parser/parser.go**

```go
case ir.Architecture:
    return parseArchitecture(input)
```

**Step 4: Verify**

Run: `go test ./parser/... -v -run Architecture`

---

### Task 6: Journey Layout

**Files:**
- Create: `layout/journey.go`
- Create: `layout/journey_test.go`
- Modify: `layout/types.go` — add JourneyData and related types
- Modify: `layout/layout.go` — add Journey dispatch case

**Step 1: Add layout types to layout/types.go**

```go
type JourneyData struct {
	Sections []JourneySectionLayout
	Title    string
	Actors   []JourneyActorLayout
	TrackY   float32
	TrackH   float32
}
func (JourneyData) diagramData() {}

type JourneySectionLayout struct {
	Label  string
	X, Y   float32
	Width  float32
	Height float32
	Color  string
	Tasks  []JourneyTaskLayout
}

type JourneyTaskLayout struct {
	Label  string
	Score  int
	X, Y   float32
	Width  float32
	Height float32
}

type JourneyActorLayout struct {
	Name       string
	ColorIndex int
}
```

**Step 2: Create layout/journey.go**

```go
func computeJourneyLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout
```

Layout algorithm:
1. Measure task labels with textmetrics
2. Position sections left-to-right, each section width = sum of task widths + spacing
3. Within each section, position tasks left-to-right
4. Task Y = trackY + trackH * (1 - (score-1)/4) — score 5 at top, score 1 at bottom
5. Collect unique actors across all tasks
6. Compute total width = sum of section widths + gaps + padding
7. Compute total height = title + track + actor legend + padding
8. Return Layout with JourneyData and no Nodes/Edges (all data in Diagram field)

**Step 3: Create layout/journey_test.go**

Tests:
- TestJourneyLayout: 2 sections, verify section/task positions
- TestJourneyLayoutEmpty: empty graph

**Step 4: Add dispatch in layout/layout.go**

```go
case ir.Journey:
    return computeJourneyLayout(g, th, cfg)
```

**Step 5: Verify**

Run: `go test ./layout/... -v -run Journey`

---

### Task 7: Architecture Layout

**Files:**
- Create: `layout/architecture.go`
- Create: `layout/architecture_test.go`
- Modify: `layout/types.go` — add ArchitectureData and related types
- Modify: `layout/layout.go` — add Architecture dispatch case

**Step 1: Add layout types to layout/types.go**

```go
type ArchitectureData struct {
	Groups    []ArchGroupLayout
	Junctions []ArchJunctionLayout
}
func (ArchitectureData) diagramData() {}

type ArchGroupLayout struct {
	ID     string
	Label  string
	Icon   string
	X, Y   float32
	Width  float32
	Height float32
}

type ArchJunctionLayout struct {
	ID   string
	X, Y float32
	Size float32
}
```

**Step 2: Create layout/architecture.go**

```go
func computeArchitectureLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout
```

Layout algorithm (constraint-based grid):
1. Size service nodes using textmetrics
2. Build adjacency map from ArchEdges with directional hints
3. BFS from first service: place at grid (0,0)
4. For each edge, derive neighbor grid position from sides:
   - A:R -- L:B → B is one column right of A
   - A:B -- T:B → B is one row below A
   - A:L -- R:B → B is one column left of A
   - A:T -- B:B → B is one row above A
5. Normalize grid to non-negative coordinates
6. Convert grid to pixel positions: x = col * (serviceWidth + colGap), y = row * (serviceHeight + rowGap)
7. Place junction nodes at grid positions (small circle)
8. Compute group bounding rectangles from children + padding
9. Build EdgeLayout with anchor points on correct sides of nodes
10. Return Layout with ArchitectureData

**Step 3: Create layout/architecture_test.go**

Tests:
- TestArchitectureLayout: 3 services with edges, verify positions
- TestArchitectureLayoutGroups: services in groups, verify group bounds
- TestArchitectureLayoutEmpty: empty graph

**Step 4: Add dispatch in layout/layout.go**

```go
case ir.Architecture:
    return computeArchitectureLayout(g, th, cfg)
```

**Step 5: Verify**

Run: `go test ./layout/... -v -run Architecture`

---

### Task 8: Journey Renderer

**Files:**
- Create: `render/journey.go`
- Create: `render/journey_test.go`
- Modify: `render/svg.go` — add JourneyData dispatch case

**Step 1: Create render/journey.go**

```go
func renderJourney(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout)
```

Rendering:
1. Title text centered at top (if present)
2. Dashed horizontal lines for score levels (1-5)
3. Section background rectangles with fill colors cycling through JourneySectionColors
4. Task rounded rectangles with label centered, fill from JourneyTaskFill
5. Score indicator: small filled circle using JourneyScoreColors[score-1]
6. Actor legend at bottom: colored dots + names

**Step 2: Create render/journey_test.go**

Tests:
- TestRenderJourney: verify SVG contains section rects, task rects, title, score indicators
- TestRenderJourneyEmpty: empty diagram outputs valid SVG

**Step 3: Add dispatch in render/svg.go**

```go
case layout.JourneyData:
    renderJourney(&b, l, th, cfg)
```

**Step 4: Verify**

Run: `go test ./render/... -v -run Journey`

---

### Task 9: Architecture Renderer

**Files:**
- Create: `render/architecture.go`
- Create: `render/architecture_test.go`
- Modify: `render/svg.go` — add ArchitectureData dispatch case

**Step 1: Create render/architecture.go**

```go
func renderArchitecture(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout)
```

Rendering:
1. Group rectangles (rounded, dashed border) with icon + label header
2. Service rectangles (rounded, solid border) with icon + label
3. Junction circles (small filled dots)
4. Edges using node side anchor points:
   - L: (x - w/2, y), R: (x + w/2, y), T: (x, y - h/2), B: (x, y + h/2)
   - Arrow markers based on ArrowLeft/ArrowRight
5. Built-in icon rendering: simple SVG shapes for cloud/database/disk/internet/server

Helper:
```go
func renderArchIcon(b *svgBuilder, icon string, cx, cy, size float32)
```

Renders simple icons:
- cloud: ellipse
- database: cylinder (rect + arcs)
- disk: circle
- internet: circle with cross lines
- server: rectangle with horizontal lines

**Step 2: Create render/architecture_test.go**

Tests:
- TestRenderArchitecture: services + group + edges
- TestRenderArchitectureEmpty: empty diagram

**Step 3: Add dispatch in render/svg.go**

```go
case layout.ArchitectureData:
    renderArchitecture(&b, l, th, cfg)
```

**Step 4: Verify**

Run: `go test ./render/... -v -run Architecture`

---

### Task 10: Integration Tests and Fixtures

**Files:**
- Create: `testdata/fixtures/journey-basic.mmd`
- Create: `testdata/fixtures/journey-multiactor.mmd`
- Create: `testdata/fixtures/architecture-basic.mmd`
- Create: `testdata/fixtures/architecture-groups.mmd`
- Modify: `mermaid_test.go` — add 4 integration tests

**Step 1: Create fixture files**

journey-basic.mmd:
```
journey
  title My Working Day
  section Go to work
    Make tea: 5: Me
    Go upstairs: 3: Me
    Do work: 1: Me, Cat
  section Go home
    Go downstairs: 5: Me
    Sit down: 5: Me
```

journey-multiactor.mmd:
```
journey
  title Shopping Trip
  section Browse
    Enter store: 4: Alice, Bob
    Look at items: 3: Alice
  section Purchase
    Checkout: 2: Alice, Bob
    Pay: 1: Alice
```

architecture-basic.mmd:
```
architecture-beta
  service db(database)[Database]
  service server(server)[Server]
  service web(cloud)[WebApp]
  db:R -- L:server
  server:R --> L:web
```

architecture-groups.mmd:
```
architecture-beta
  group api(cloud)[API]
  service db(database)[Database] in api
  service server(server)[Server] in api
  db:R -- L:server
```

**Step 2: Add integration tests to mermaid_test.go**

4 tests following the existing pattern:
- TestRenderJourneyBasicFixture
- TestRenderJourneyMultiactorFixture
- TestRenderArchitectureBasicFixture
- TestRenderArchitectureGroupsFixture

Each reads fixture, calls Render(), verifies non-empty SVG with diagram-specific content.

**Step 3: Verify**

Run: `go test -v -run "Journey|Architecture"` — all pass.

---

### Task 11: Final Validation

**Step 1:** `go vet ./...`
**Step 2:** `go build ./...`
**Step 3:** `gofmt -l .` — should output nothing
**Step 4:** `go test -race ./...` — all packages pass
**Step 5:** Run go-code-reviewer agent
**Step 6:** Fix any P1 issues from review
**Step 7:** Commit all fixes
