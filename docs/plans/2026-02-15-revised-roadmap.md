# Revised Roadmap

**Date:** 2026-02-15
**Status:** Approved

## Completed Phases

### Phase 1: Foundation (complete)
- Flowchart diagrams: TD/LR/RL/BT, subgraphs, shapes, edge styles
- Pipeline: Parse -> IR -> Layout (Sugiyama) -> Render SVG
- Text metrics via sfnt, stdlib regexp parsing

### Phase 2: Core Graph Variants (complete)
- Class diagrams: classes, interfaces, members, relationships, namespaces
- State diagrams: states, transitions, start/end/fork/join/choice, composites
- ER diagrams: entities, attributes, relationships with cardinality
- Reused Sugiyama layout; state got recursive layout for composites

### Phase 3: Sequence Diagrams (complete)
- All 10 arrow types, activations, notes, frames, boxes, autonumber, create/destroy
- Timeline-based layout (not Sugiyama)
- Participant types: box, actor, boundary, control, entity, database, collections, queue

## Upcoming Phases

### Phase 4: Simple Grid — Kanban & Packet

**Layout engine:** 2D grid/table — cells in rows and columns, no edges.

**Diagrams:**
- **Kanban**: Columns with cards. Column widths based on widest card, cards stacked vertically.
- **Packet**: Network packet header. Fixed-width rows, fields fill left-to-right, wrap at 32 bits.

**Shared code:** `layout/grid.go` — grid cell layout with text measurement for sizing.

**Estimated scope:** ~800-1000 lines.

---

### Phase 5: Pie & Quadrant

**Layout engine:** Positioned geometry — no edges, no grid. Items placed by computed coordinates.

**Diagrams:**
- **Pie**: Sectors from percentage data. SVG arc paths + legend.
- **Quadrant**: 2x2 quadrant with labeled axes and positioned points.

**Shared code:** Text measurement for labels, self-contained positioning.

**Estimated scope:** ~700-900 lines.

---

### Phase 6: Timeline, Gantt & GitGraph

**Layout engine:** Horizontal timeline — items positioned along a time axis.

**Diagrams:**
- **Timeline**: Sections with events along a horizontal axis. Simple time periods.
- **Gantt**: Tasks with start/end dates, dependencies, sections, milestones. Date parsing and dependency-aware bar positioning.
- **GitGraph**: Commits on branches along a horizontal axis. Branch/merge/cherry-pick as swim-lane commit graph.

**Shared code:** `layout/timeline.go` — horizontal axis math, vertical swim lanes.

**Estimated scope:** ~1500-2000 lines. Gantt is the most complex (date math, dependency arrows).

---

### Phase 7: Charts — XYChart & Radar

**Layout engine:** Cartesian/polar axes with data series rendering.

**Diagrams:**
- **XYChart**: Bar and line charts on x/y axes. Axis scales, tick marks, bars/lines.
- **Radar** (spider chart): Radial axes with polygon data overlay. Polar coordinates.

**Shared code:** Axis scale computation, tick generation, data series normalization.

**Estimated scope:** ~1000-1300 lines.

---

### Phase 8: Hierarchical — Mindmap, Sankey & Treemap

**Layout engine:** Tree/DAG layout with parent-child relationships.

**Diagrams:**
- **Mindmap**: Radial tree from a root node. Indentation-based hierarchy, curved branches.
- **Sankey**: Weighted flow between nodes across columns. Variable-width curve paths.
- **Treemap**: Nested rectangles proportional to values. Squarified treemap algorithm.

**Shared code:** Hierarchical data structures, parent-child traversal.

**Estimated scope:** ~1800-2200 lines. Most complex phase.

---

### Phase 9: Graph Variants — Requirement, Block & C4

**Layout engine:** Sugiyama (reuse existing `runSugiyama()`).

**Diagrams:**
- **Requirement**: UML requirement stereotypes + relationships (contains, derives, satisfies, verifies, refines, traces, copies).
- **Block**: Block definition diagrams with containers, ports, typed edges. Nested blocks similar to subgraphs.
- **C4**: C4 architecture model (System Context, Container, Component, Dynamic, Deployment). Stereotyped nodes per C4 level.

**Shared code:** Existing Sugiyama layout — new parsers + renderers only.

**Estimated scope:** ~1200-1500 lines.

---

### Phase 10: Journey & Architecture

**Layout engine:** Custom positioning.

**Diagrams:**
- **Journey**: User journey map with sections, tasks, satisfaction scores on a horizontal track.
- **Architecture**: Manual positioning with explicit x,y coordinates for services/groups.

**Shared code:** Minimal — both are unique custom layouts.

**Estimated scope:** ~800-1000 lines.

---

### Phase 11: ZenUML

**Layout engine:** Sequence variant (reuse Phase 3's timeline layout).

**Diagrams:**
- **ZenUML**: Code-like sequence diagram syntax. Different parser, same layout and renderer.

**Shared code:** Reuses `computeSequenceLayout()` and sequence renderer entirely. Parser only.

**Estimated scope:** ~400-500 lines.

---

## Cross-Cutting Concerns

These are deferred improvements that apply across all diagram types:

- **A* edge routing** (deferred from Phase 1): Replace L-shaped paths with smarter routing that avoids node overlap.
- **Parser error handling** (P2 bugs from Phase 3): `parseSequence` never returns errors for structural issues; JSON parse errors silently swallowed. Pattern should be fixed across all parsers.
- **Accessibility**: SVG `<title>` and `aria-label` attributes for screen readers.
- **CLI tool**: Command-line interface for rendering `.mmd` files to SVG.
