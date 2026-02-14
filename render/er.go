package render

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// renderER renders all ER diagram elements: edges with labels, then entity boxes.
func renderER(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	erData, ok := l.Diagram.(layout.ERData)
	if !ok {
		return
	}

	// Render edges first (behind entities).
	renderEREdges(b, l, th)

	// Render entity boxes (on top of edges).
	renderEREntities(b, l, erData, th, cfg)
}

// renderEREdges renders all ER diagram edges as SVG paths with optional labels.
// ER diagrams use plain lines without arrow markers; crow's foot notation
// decorations are not yet rendered.
func renderEREdges(b *svgBuilder, l *layout.Layout, th *theme.Theme) {
	for i, edge := range l.Edges {
		if len(edge.Points) < 2 {
			continue
		}

		d := pointsToPath(edge.Points)
		edgeID := fmt.Sprintf("edge-%d", i)

		b.selfClose("path",
			"id", edgeID,
			"class", "edgePath",
			"d", d,
			"fill", "none",
			"stroke", th.LineColor,
			"stroke-width", "1.5",
			"stroke-linecap", "round",
			"stroke-linejoin", "round",
		)

		// Render edge label if present.
		if edge.Label != nil && len(edge.Label.Lines) > 0 {
			renderEdgeLabel(b, edge, th)
		}
	}
}

// renderEREntities renders entity boxes sorted by ID for deterministic output.
func renderEREntities(b *svgBuilder, l *layout.Layout, erData layout.ERData, th *theme.Theme, cfg *config.Layout) {
	ids := make([]string, 0, len(l.Nodes))
	for id := range l.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		n := l.Nodes[id]
		entity := erData.Entities[id]
		dims := erData.EntityDims[id]

		if entity == nil {
			// Fallback: render as a plain rectangle node.
			renderNodeShape(b, n, th.PrimaryColor, th.PrimaryBorderColor, th.PrimaryTextColor)
			continue
		}

		renderEntityBox(b, n, entity, dims, th, cfg)
	}
}

// renderEntityBox renders a single ER entity as a table-like box with a
// colored header and attribute rows.
func renderEntityBox(b *svgBuilder, n *layout.NodeLayout, entity *ir.Entity, dims layout.EntityDimensions, th *theme.Theme, cfg *config.Layout) {
	x := n.X - n.Width/2
	y := n.Y - n.Height/2
	w := n.Width
	h := n.Height

	padH := cfg.Padding.NodeHorizontal
	lineH := th.FontSize * cfg.LabelLineHeight
	rowH := cfg.ER.AttributeRowHeight
	if rowH <= 0 {
		rowH = lineH
	}

	// Outer border rect (full entity box).
	b.rect(x, y, w, h, 4,
		"fill", th.EntityBodyBg,
		"stroke", th.EntityBorder,
		"stroke-width", "1",
	)

	// Header rect (filled with primary/header color).
	headerH := dims.HeaderHeight
	b.rect(x, y, w, headerH, 0,
		"fill", th.EntityHeaderBg,
		"stroke", th.EntityBorder,
		"stroke-width", "1",
	)

	// Header text: entity display name, centered.
	displayName := entity.DisplayName()
	headerTextY := y + headerH/2 + th.FontSize*0.35
	b.text(x+w/2, headerTextY, displayName,
		"text-anchor", "middle",
		"fill", th.PrimaryTextColor,
		"font-size", fmtFloat(th.FontSize),
		"font-weight", "bold",
	)

	// Separator line below header.
	b.line(x, y+headerH, x+w, y+headerH,
		"stroke", th.EntityBorder,
		"stroke-width", "1",
	)

	// Attribute rows.
	bodyY := y + headerH
	colPad := cfg.ER.ColumnPadding
	if colPad <= 0 {
		colPad = padH
	}

	for i, attr := range entity.Attributes {
		rowY := bodyY + float32(i)*rowH

		// Alternating row background for readability.
		if i%2 == 1 {
			b.rect(x+1, rowY, w-2, rowH, 0,
				"fill", th.EntityBodyBg,
				"stroke", "none",
				"opacity", "0.5",
			)
		}

		// Horizontal separator between rows (skip first row, it has the header line).
		if i > 0 {
			b.line(x, rowY, x+w, rowY,
				"stroke", th.EntityBorder,
				"stroke-width", "0.5",
				"opacity", "0.4",
			)
		}

		textY := rowY + rowH/2 + th.FontSize*0.35

		// Column 1: Key constraints (PK, FK, UK).
		keyX := x + colPad
		var keyStr string
		for j, k := range attr.Keys {
			if j > 0 {
				keyStr += ","
			}
			keyStr += k.String()
		}
		if keyStr != "" {
			b.text(keyX, textY, keyStr,
				"fill", th.TextColor,
				"font-size", fmtFloat(th.FontSize*0.85),
				"font-weight", "bold",
			)
		}

		// Column 2: Type.
		typeX := keyX + dims.KeyColWidth + colPad
		b.text(typeX, textY, attr.Type,
			"fill", th.TextColor,
			"font-size", fmtFloat(th.FontSize*0.85),
			"font-style", "italic",
		)

		// Column 3: Name.
		nameX := typeX + dims.TypeColWidth + colPad
		b.text(nameX, textY, attr.Name,
			"fill", th.TextColor,
			"font-size", fmtFloat(th.FontSize*0.85),
		)
	}
}
