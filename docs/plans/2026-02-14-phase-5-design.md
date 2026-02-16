# Phase 5: Pie & Quadrant — Design

**Date:** 2026-02-14
**Status:** Approved

## Goal

Add Pie chart and Quadrant chart diagram support to gomd2svg. Both use positioned geometry (no edges, no grid, no Sugiyama).

## IR Types

### Pie (`ir/pie.go`)

- `PieSlice` struct: `Label string`, `Value float64`
- Graph fields: `PieSlices []*PieSlice`, `PieTitle string`, `PieShowData bool`

### Quadrant (`ir/quadrant.go`)

- `QuadrantPoint` struct: `Label string`, `X float64`, `Y float64` (0.0-1.0 normalized)
- Graph fields: `QuadrantPoints []*QuadrantPoint`, `QuadrantTitle string`, `QuadrantLabels [4]string` (q1=top-right, q2=top-left, q3=bottom-left, q4=bottom-right), `XAxisLeft string`, `XAxisRight string`, `YAxisBottom string`, `YAxisTop string`

## Parsers

### Pie

Line-by-line parsing:
- First content line: `pie` keyword, optional `showData` suffix
- `title <text>` line sets PieTitle
- Data lines: `"Label" : value` — regex `^\s*"([^"]+)"\s*:\s*(\d+\.?\d*)\s*$`
- Values must be positive (> 0)

### Quadrant

Line-by-line parsing:
- First content line: `quadrantChart`
- `title <text>` — sets QuadrantTitle
- `x-axis <left> --> <right>` or `x-axis <left>` — sets axis labels
- `y-axis <bottom> --> <top>` or `y-axis <bottom>` — sets axis labels
- `quadrant-1` through `quadrant-4` — sets quadrant labels
- Data points: `<label>: [x, y]` — regex `^\s*(.+?):\s*\[([0-9.]+),\s*([0-9.]+)\]\s*$`
- x, y values are 0.0-1.0

## Config

### PieConfig

| Field | Type | Default | Purpose |
|-------|------|---------|---------|
| Radius | float32 | 150 | Pie circle radius |
| InnerRadius | float32 | 0 | Donut hole (0 = solid) |
| TextPosition | float32 | 0.75 | Label position along radius (0=center, 1=edge) |
| PaddingX | float32 | 20 | Horizontal canvas padding |
| PaddingY | float32 | 20 | Vertical canvas padding |

### QuadrantConfig

| Field | Type | Default | Purpose |
|-------|------|---------|---------|
| ChartWidth | float32 | 400 | Quadrant area width |
| ChartHeight | float32 | 400 | Quadrant area height |
| PointRadius | float32 | 5 | Default data point radius |
| PaddingX | float32 | 40 | Horizontal canvas padding |
| PaddingY | float32 | 40 | Vertical canvas padding |
| QuadrantLabelFontSize | float32 | 14 | Font size for quadrant names |
| AxisLabelFontSize | float32 | 12 | Font size for axis labels |

## Layout

### Pie (`layout/pie.go`)

Compute cumulative angles from percentage of total value. Each slice gets:
- StartAngle, EndAngle (radians, clockwise from top)
- LabelX, LabelY — positioned at TextPosition along the radius at the midpoint angle
- No Sugiyama — pure trigonometric geometry

Layout types: `PieData` (implements DiagramData), `PieSliceLayout` (angles + label position).

### Quadrant (`layout/quadrant.go`)

Map [0,1] normalized coordinates to pixel positions within the chart area. Layout computes:
- Chart origin (top-left of quadrant area, after padding for title and axis labels)
- Four quadrant background rects
- Axis label positions
- Data point pixel positions
- Title position

Layout types: `QuadrantData` (implements DiagramData), `QuadrantPointLayout` (pixel x, y + label).

## Rendering

### Pie (`render/pie.go`)

- SVG `<path>` with arc commands (`M`, `L`, `A`) for each slice
- Slice colors from theme palette (cycle through)
- Labels at computed positions (slice name, optionally value if showData)
- Title centered above pie

### Quadrant (`render/quadrant.go`)

- Four background rects with distinct quadrant fills
- Quadrant labels centered in each quadrant
- Axis lines (horizontal and vertical center lines)
- Axis labels at edges
- Data points as `<circle>` elements
- Point labels as `<text>` next to each point
- Title centered above chart

## Deferred

- Per-point styling (radius, color, stroke) — Phase 12 Theme System
- classDef for quadrant points — Phase 12 Theme System
- Pie donut mode (InnerRadius > 0) — available via config but not parsed from syntax
