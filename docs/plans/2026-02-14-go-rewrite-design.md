# Go Rewrite Design: mermaid-go

**Date:** 2026-02-14
**Status:** Approved
**Motivation:** Ecosystem integration -- embeddable Go library for Go services and tooling

## Constraints

- Full parity with all 23 diagram types from day one
- SVG + PNG output formats
- Library-first API with CLI binary included
- Pure Go font parsing (no CGo), fully cross-compilable
- Performance target: as close to Rust as possible

## 1. Project Structure

```
github.com/jamesainslie/mermaid-go/
├── cmd/mmdr/              # CLI binary
│   └── main.go
├── mermaid.go             # Public API: Render(), Parse(), Layout(), RenderSVG()
├── options.go             # RenderOptions, RenderResult types
├── ir/                    # Intermediate representation
│   ├── graph.go           # Graph, Node, Edge, Subgraph
│   ├── diagram.go         # DiagramKind, Direction enums
│   ├── sequence.go        # Sequence-specific IR types
│   ├── c4.go              # C4-specific IR types
│   ├── gitgraph.go        # Git graph IR types
│   ├── style.go           # NodeStyle, EdgeStyleOverride
│   └── shapes.go          # NodeShape, EdgeStyle, EdgeDecoration enums
├── parser/                # Mermaid syntax parsing
│   ├── parser.go          # Parse dispatch + detectDiagramKind + shared helpers
│   ├── flowchart.go       # Flowchart parser
│   ├── sequence.go        # Sequence diagram parser
│   ├── class.go           # Class diagram parser
│   ├── state.go           # State diagram parser
│   ├── er.go              # ER diagram parser
│   └── ...                # One file per diagram type (23 total)
├── layout/                # Layout computation
│   ├── layout.go          # ComputeLayout, computeGraphLayout (main orchestrator)
│   ├── types.go           # Layout, NodeLayout, EdgeLayout, DiagramData interface
│   ├── sizing.go          # sizeNodes, text block measurement
│   ├── ranking.go         # assignRanks (longest-path DAG layering)
│   ├── ordering.go        # orderWithinRanks (crossing minimization)
│   ├── positioning.go     # positionNodes (coordinate assignment)
│   ├── routing.go         # grid, A* search, routeEdges
│   ├── label.go           # placeLabels, label candidate search
│   ├── overlap.go         # resolveOverlaps
│   └── ...                # One file per diagram-specific layout
├── render/                # SVG generation
│   ├── svg.go             # RenderSVG + main rendering dispatch
│   ├── builder.go         # svgBuilder primitive methods
│   ├── shapes.go          # renderNodeShape, diamond(), hexagon(), cylinder()
│   ├── graph.go           # renderGraph (flowchart/class/state/ER)
│   ├── png.go             # WritePNG (build-tagged)
│   ├── png_stub.go        # Stub when png tag absent
│   └── ...                # One file per diagram-type renderer
├── config/                # Configuration
│   └── config.go          # Config, LayoutConfig structs with defaults
├── theme/                 # Theme management
│   ├── theme.go           # Theme struct, Modern(), MermaidDefault()
│   └── color.go           # HSL parsing, color adjustment utilities
├── textmetrics/           # Font measurement
│   ├── measurer.go        # Measurer with sync.Mutex + cache
│   └── font.go            # Pure Go font face parsing + glyph advance lookup
└── internal/              # Internal utilities
    └── grid/              # Occupancy grid for A* routing
        └── grid.go
```

Rust-to-Go idiom translations:

- Rust `enum` with variants -> Go `iota` constants with a named type
- Rust `BTreeMap<String, Node>` -> Go `map[string]*Node` with separate `[]string` for ordered keys
- Rust `Option<T>` -> Go pointer `*T` or comma-ok pattern
- Rust `anyhow::Result` -> Go `(T, error)` return pattern
- Rust `Lazy<Regex>` -> Go package-level `regexp.MustCompile`
- Rust sum-type enums -> Go sealed interface with unexported method

## 2. Public API

```go
package mermaid

// One-liner API
func Render(input string) (string, error)
func RenderWithOptions(input string, opts Options) (string, error)

// Pipeline API for stage-level control
func Parse(input string) (*ParseOutput, error)
func ComputeLayout(g *ir.Graph, theme *theme.Theme, cfg *config.Layout) *layout.Layout
func RenderSVG(l *layout.Layout, theme *theme.Theme, cfg *config.Layout) string

// Timing API for benchmarking
func RenderWithTiming(input string, opts Options) (*Result, error)
```

`Options` zero value uses Modern theme + production defaults. Functional options
(`WithNodeSpacing`, `WithRankSpacing`, `WithTheme`, `WithPreferredAspectRatio`) for
builder-style construction.

`Result` contains SVG string + per-stage microsecond timings (`ParseUs`, `LayoutUs`, `RenderUs`).

## 3. IR Types

### Enums as typed constants

```go
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
```

Same pattern for `Direction`, `NodeShape`, `EdgeStyle`, `EdgeDecoration`, `EdgeArrowhead`,
and all diagram-specific enums.

### Graph struct

The Rust `Graph` is already a flat struct with diagram-specific fields. This maps directly
to Go. Universal fields (`Nodes`, `Edges`, `Subgraphs`) plus diagram-specific fields
(`SequenceParticipants`, `PieSlices`, `GitGraph`, `C4`, etc.).

Optional fields use pointers (`*string`, `*float32`, `*int`). Diagram-specific data structs
(`C4Data`, `GitGraphData`, etc.) are embedded by value since they have useful zero values.
`*BlockDiagram` is a pointer because nil means "no block diagram".

### Layout DiagramData

Sealed interface with unexported method to prevent external implementation:

```go
type DiagramData interface {
    diagramData()
}
```

One struct per variant (`GraphData`, `SequenceLayoutData`, `PieLayoutData`, etc.).
Consumers type-switch to access diagram-specific layout data.

## 4. Parser Architecture

Dispatcher detects diagram kind from first non-comment line, routes to type-specific parsers.

Each diagram type gets its own file (`parser/flowchart.go`, `parser/sequence.go`, etc.).
Shared helpers in `parser/parser.go`: `stripTrailingComment`, `ensureNode`, `parseNodeToken`,
`parseArrowToken`, `parseInitDirective`.

Regexps compiled at package level with `regexp.MustCompile` (panics on invalid pattern --
acceptable for compile-time constants).

Go `regexp` is RE2-based. For the few complex arrow-matching patterns, benchmark early and
replace with hand-written string scanning if they become hot paths.

## 5. Layout Engine

### Universal graph layout pipeline

1. Size nodes based on label text metrics
2. Assign ranks (DAG layering via longest-path)
3. Order nodes within ranks (crossing minimization)
4. Position nodes in coordinate space
5. Compute subgraph bounding boxes
6. Route edges (orthogonal A* on occupancy grid)
7. Place edge labels (collision-aware search)
8. Resolve overlaps (iterative passes)
9. Compute final bounding box

### Edge routing

A* search on a flat `[]bool` occupancy grid (single allocation, row-major). Falls back to
direct L-shaped path if A* exceeds budget. `sync.Pool` for A* nodes to reduce GC pressure.
`container/heap`-backed priority queue.

### Label placement

Search grid for non-overlapping positions. Each candidate scored by distance from ideal
midpoint, proximity to edge path, and collision penalty.

### Text measurement

`textmetrics.Measurer` wraps `golang.org/x/image/font/sfnt` for glyph advance lookup.
Thread-safe via `sync.Mutex`. Caches measured widths. System font discovery via known
platform paths + `os.UserCacheDir()` for font cache.

Created at `ComputeLayout` entry and threaded through (not global).

## 6. SVG Renderer

`svgBuilder` wraps `strings.Builder` with typed methods for SVG primitives (`rect`, `circle`,
`path`, `text`, `line`, `polyline`). `attr` key-value pairs for SVG attributes. Float formatting
trims trailing zeros for smaller SVG output.

Main `RenderSVG` dispatches to diagram-specific renderers via type-switch on `DiagramData`.
Each diagram type renders in its own file.

### PNG output

Build-tagged (`//go:build png`). Uses `oksvg` + `rasterx` for pure Go SVG-to-PNG. Stub
file provides error message when tag is absent.

Our SVG subset (rects, paths, circles, text, basic transforms) is well within oksvg's
capabilities.

## 7. Configuration, Theming & CLI

### Config

`config.Layout` with explicit `Default*()` constructors. Zero-value `Options{}` calls
defaults internally.

### Theme

Two built-in presets (`Modern`, `MermaidDefault`). JSON tags on all fields for
deserialization from config files.

Color utilities: `AdjustColor`, `ParseColorToHSL`, `parseHex`, `rgbToHSL`.

### CLI

`cmd/mmdr/main.go` using stdlib `flag` package. No subcommands needed.

Flags: `-i` (input), `-o` (output), `-e` (format: svg/png), `-c` (config file),
`--timing`, `--theme`, `--dump-layout`.

## 8. Testing Strategy

### Layer 1: Unit tests (table-driven)

Every package gets table-driven tests with `t.Run` subtests and `t.Parallel()`.

### Layer 2: Golden file / snapshot tests

Use existing Rust test fixtures (~40 diagrams) as inputs. `-update` flag regenerates
golden SVG files. Compare with `go-cmp`.

### Layer 3: Cross-implementation SVG diff

Compare Go output against Rust output structurally (node count, edge count, bounding box
within tolerance). Runs in CI nightly.

### Benchmarks

Per-stage benchmarks (`BenchmarkParse`, `BenchmarkLayout`, `BenchmarkRenderSVG`,
`BenchmarkRouteEdges`, `BenchmarkPlaceLabels`). Track regressions with `benchstat`.

## 9. Dependencies

### Production (SVG-only: zero external deps)

| Dependency | Purpose |
|------------|---------|
| `golang.org/x/image/font/sfnt` | TTF/OTF parsing for glyph advances |
| stdlib (`regexp`, `encoding/json`, `strings`, `sync`, `math`, `image`, `image/png`) | Core operations |

### Production (PNG: 2 external deps, build-tagged)

| Dependency | Purpose |
|------------|---------|
| `github.com/srwiley/oksvg` | SVG parsing for PNG conversion |
| `github.com/srwiley/rasterx` | SVG rasterization |

### Dev/test

| Dependency | Purpose |
|------------|---------|
| `github.com/google/go-cmp` | Structural diffs in tests |

### Build

- Single `go.mod`, not multi-module
- Go 1.23+
- Build tag `png` for PNG support
- Cross-compiles cleanly (no CGo)

## 10. Implementation Order

### Phase 1: Foundation

IR types, theme, config, textmetrics, flowchart parser, graph layout pipeline, SVG renderer.

**Exit criteria:** `mermaid.Render("flowchart LR; A-->B-->C")` produces valid SVG.
Cross-validate against Rust output.

### Phase 2: Core graph variants

Class, state, ER parsers + rendering. These share the universal graph layout pipeline.

### Phase 3: Sequence diagrams

Parser, layout (lifelines, activations, frames, notes), renderer.

### Phase 4: Remaining diagram types

18 remaining types. Each is self-contained (parser + layout + render file). Parallelizable.

### Phase 5: Edge routing & label placement optimization

A* with occupancy grid, collision-aware label search. Phases 1-4 can use simpler fallback
routing (direct L-shaped paths).

### Phase 6: CLI & PNG

`cmd/mmdr` binary, `oksvg`/`rasterx` PNG integration.

### Phase 7: Performance tuning

Profile with `pprof`. Sync.Pool for A* nodes. Pre-allocated slices. Bit-packed grid.
Benchmark against Rust, target within 2-3x.

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Go `regexp` too slow for complex arrow patterns | Hand-written scanner for hot regexps. Benchmark in Phase 1 |
| `oksvg` can't handle our SVG subset | Test early. Fallback: shell out to `resvg` binary |
| Font metrics differ between Go and Rust | Cross-validate text widths. Tolerance-based comparison |
| A* routing slower than Rust | Profile in Phase 5. sync.Pool, cache-friendly grid |
| `x/image/font/sfnt` can't find system fonts | Platform-specific font discovery with known paths |
