# Phase 9 Design: Requirement, Block, C4 Diagrams

## Goal

Add Requirement, Block, and C4 diagram support to gomd2svg, completing the Sugiyama-reuse phase of the roadmap.

## Architecture

All three diagram types follow the established 4-layer pipeline (IR -> Parser -> Layout -> Render). Requirement and C4 reuse Sugiyama directly (like Class/ER). Block uses a hybrid: grid layout when `columns` is specified, Sugiyama when blocks have connections without a grid.

## IR Types

### Requirement (`ir/requirement.go`)

- `RequirementType` enum (iota): Requirement, FunctionalRequirement, InterfaceRequirement, PerformanceRequirement, PhysicalRequirement, DesignConstraint
- `RiskLevel` enum (iota): Low, Medium, High
- `VerifyMethod` enum (iota): Analysis, Inspection, Test, Demonstration
- `RequirementDef` struct: ID, Label, Text, Type, Risk, VerifyMethod
- `ElementDef` struct: ID, Label, Type, DocRef
- `RequirementRelType` enum: Contains, Copies, Derives, Satisfies, Verifies, Refines, Traces
- Graph fields: `Requirements []*RequirementDef`, `Elements []*ElementDef`

### Block (`ir/block.go`)

- `BlockDef` struct: ID, Label, Shape (reuse ir.Shape), Width (column span int), Children []*BlockDef
- Graph fields: `Blocks []*BlockDef`, `BlockColumns int`
- Reuses existing `ir.Shape` enum and edge types from flowchart

### C4 (`ir/c4.go`)

- `C4Kind` enum: C4Context, C4Container, C4Component, C4Dynamic, C4Deployment
- `C4ElementType` enum: Person, System, SystemDb, SystemQueue, Container, ContainerDb, ContainerQueue, Component, ExternalPerson, ExternalSystem, ExternalSystemDb, ExternalSystemQueue, ExternalContainer, ExternalContainerDb, ExternalContainerQueue, ExternalComponent
- `C4Element` struct: ID, Label, Technology, Description, ElementType, BoundaryID
- `C4Boundary` struct: ID, Label, Type, Children []string
- Graph fields: `C4Elements []*C4Element`, `C4Boundaries []*C4Boundary`, `C4Kind C4Kind`

## Parsers

### Requirement (`parser/requirement.go`)

Brace-delimited blocks (same pattern as Class parser). Parse `requirement name { ... }` and `element name { ... }` blocks, extract metadata fields (id, text, risk, verifymethod, type, docref). Parse relationship lines: `source - relType -> target`.

### Block (`parser/block.go`)

Parse `columns N` directive. Parse block definitions with optional shape syntax (reuse flowchart shape patterns), width spans (`:N`), nesting (indentation or braces), and edge syntax (`-->`, `---`, `-->|label|`).

### C4 (`parser/c4.go`)

Function-call syntax. Parse `Person(id, "name", "desc")`, `System(id, "name", "desc")`, `Container_Boundary(id, "label") { ... }`, `Rel(from, to, "label")`. Support named params with `$key=value` for styling overrides.

## Layouts

### Requirement (`layout/requirement.go`)

Nodes sized with text metrics (label + metadata rows). `runSugiyama()` for positioning. `RequirementData` carries requirement/element metadata for rendering.

### Block (`layout/block.go`) — Hybrid

- If `columns > 0`: grid layout. Place blocks left-to-right, wrap at column count, respect span widths. Draw edges over grid.
- If `columns == 0` and edges exist: `runSugiyama()`.
- If `columns == 0` and no edges: single-column vertical stack.
- `BlockData` carries column count and nesting info.

### C4 (`layout/c4.go`)

Nodes sized with text metrics (label + technology + description multi-line). `runSugiyama()` for positioning within boundaries. Boundaries rendered as grouping rectangles (similar to subgraphs/namespaces). `C4Data` carries element types and boundaries.

## Renderers

### Requirement (`render/requirement.go`)

Rounded rectangles with type stereotype header (`<<requirement>>`), metadata rows (id, text, risk, verifymethod). Relationship labels on edges. Styling similar to Class diagram UML boxes.

### Block (`render/block.go`)

Shape-specific rendering (reuse flowchart shape drawing patterns). Nested blocks rendered as containers with child blocks inside.

### C4 (`render/c4.go`)

C4-specific shapes: Person as rounded box with person icon at top, Systems as large boxes, Containers as boxes with technology subtitle, dashed boundary rectangles. Blue/gray C4 color palette from theme.

## Config & Theme

Config structs added to `config/config.go`:

- `RequirementConfig`: NodeWidth, NodePadding, MetadataFontSize, PaddingX, PaddingY
- `BlockConfig`: ColumnGap, RowGap, DefaultColumns, PaddingX, PaddingY, NodePadding
- `C4Config`: PersonWidth, PersonHeight, SystemWidth, SystemHeight, BoundaryPadding, PaddingX, PaddingY

Theme fields added to `theme/theme.go`:

- Requirement: RequirementNodeFill, RequirementNodeBorder, RequirementTextColor
- Block: BlockColors []string (8-color palette), BlockNodeFill, BlockNodeBorder
- C4: C4PersonColor, C4SystemColor, C4ContainerColor, C4ComponentColor, C4BoundaryColor, C4ExternalColor, C4TextColor

## Testing

Same pattern as all prior phases:

- Unit tests per layer (IR, parser, layout, render)
- Integration fixture tests in `mermaid_test.go`
- 2 fixtures per diagram type (6 total)
- Table-driven tests where applicable
