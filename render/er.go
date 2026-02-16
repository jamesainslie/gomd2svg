package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// ER diagram rendering constants.
const (
	erBorderRadius  float32 = 4
	erTextBaseline  float32 = 0.35
	erAttrFontScale float32 = 0.85
)

// renderER renders all ER diagram elements: edges with labels, then entity boxes.
func renderER(builder *svgBuilder, computed *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	erData, ok := computed.Diagram.(layout.ERData)
	if !ok {
		return
	}

	// Render edges first (behind entities).
	renderEREdges(builder, computed, th)

	// Render entity boxes (on top of edges).
	renderEREntities(builder, computed, erData, th, cfg)
}

// renderEREdges renders all ER diagram edges as SVG paths with optional labels.
// ER diagrams use plain lines without arrow markers; crow's foot notation
// decorations are not yet rendered.
func renderEREdges(builder *svgBuilder, computed *layout.Layout, th *theme.Theme) {
	for edgeIdx, edge := range computed.Edges {
		if len(edge.Points) < 2 {
			continue
		}

		pathData := pointsToPath(edge.Points)
		edgeID := fmt.Sprintf("edge-%d", edgeIdx)

		builder.selfClose("path",
			"id", edgeID,
			"class", "edgePath",
			"d", pathData,
			"fill", "none",
			"stroke", th.LineColor,
			"stroke-width", "1.5",
			"stroke-linecap", "round",
			"stroke-linejoin", "round",
		)

		// Render edge label if present.
		if edge.Label != nil && len(edge.Label.Lines) > 0 {
			renderEdgeLabel(builder, edge, th)
		}
	}
}

// renderEREntities renders entity boxes sorted by ID for deterministic output.
func renderEREntities(builder *svgBuilder, computed *layout.Layout, erData layout.ERData, th *theme.Theme, cfg *config.Layout) {
	ids := make([]string, 0, len(computed.Nodes))
	for id := range computed.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		node := computed.Nodes[id]
		entity := erData.Entities[id]
		dims := erData.EntityDims[id]

		if entity == nil {
			// Fallback: render as a plain rectangle node.
			renderNodeShape(builder, node, th.PrimaryColor, th.PrimaryBorderColor, th.PrimaryTextColor)
			continue
		}

		renderEntityBox(builder, node, entity, dims, th, cfg)
	}
}

// renderEntityBox renders a single ER entity as a table-like box with a
// colored header and attribute rows.
func renderEntityBox(builder *svgBuilder, node *layout.NodeLayout, entity *ir.Entity, dims layout.EntityDimensions, th *theme.Theme, cfg *config.Layout) {
	posX := node.X - node.Width/2
	posY := node.Y - node.Height/2
	width := node.Width
	height := node.Height

	padH := cfg.Padding.NodeHorizontal
	lineH := th.FontSize * cfg.LabelLineHeight
	rowH := cfg.ER.AttributeRowHeight
	if rowH <= 0 {
		rowH = lineH
	}

	// Outer border rect (full entity box).
	builder.rect(posX, posY, width, height, erBorderRadius,
		"fill", th.EntityBodyBg,
		"stroke", th.EntityBorder,
		"stroke-width", "1",
	)

	// Header rect (filled with primary/header color).
	headerH := dims.HeaderHeight
	builder.rect(posX, posY, width, headerH, 0,
		"fill", th.EntityHeaderBg,
		"stroke", th.EntityBorder,
		"stroke-width", "1",
	)

	// Header text: entity display name, centered.
	displayName := entity.DisplayName()
	headerTextY := posY + headerH/2 + th.FontSize*erTextBaseline
	builder.text(posX+width/2, headerTextY, displayName,
		"text-anchor", "middle",
		"fill", th.PrimaryTextColor,
		"font-size", fmtFloat(th.FontSize),
		"font-weight", "bold",
	)

	// Separator line below header.
	builder.line(posX, posY+headerH, posX+width, posY+headerH,
		"stroke", th.EntityBorder,
		"stroke-width", "1",
	)

	// Attribute rows.
	bodyY := posY + headerH
	colPad := cfg.ER.ColumnPadding
	if colPad <= 0 {
		colPad = padH
	}

	for attrIdx, attr := range entity.Attributes {
		rowY := bodyY + float32(attrIdx)*rowH

		// Alternating row background for readability.
		if attrIdx%2 == 1 {
			builder.rect(posX+1, rowY, width-2, rowH, 0,
				"fill", th.EntityBodyBg,
				"stroke", "none",
				"opacity", "0.5",
			)
		}

		// Horizontal separator between rows (skip first row, it has the header line).
		if attrIdx > 0 {
			builder.line(posX, rowY, posX+width, rowY,
				"stroke", th.EntityBorder,
				"stroke-width", "0.5",
				"opacity", "0.4",
			)
		}

		textY := rowY + rowH/2 + th.FontSize*erTextBaseline

		// Column 1: Key constraints (PK, FK, UK).
		keyX := posX + colPad
		keyParts := make([]string, 0, len(attr.Keys))
		for _, keyVal := range attr.Keys {
			keyParts = append(keyParts, keyVal.String())
		}
		keyStr := strings.Join(keyParts, ",")
		if keyStr != "" {
			builder.text(keyX, textY, keyStr,
				"fill", th.TextColor,
				"font-size", fmtFloat(th.FontSize*erAttrFontScale),
				"font-weight", "bold",
			)
		}

		// Column 2: Type.
		typeX := keyX + dims.KeyColWidth + colPad
		builder.text(typeX, textY, attr.Type,
			"fill", th.TextColor,
			"font-size", fmtFloat(th.FontSize*erAttrFontScale),
			"font-style", "italic",
		)

		// Column 3: Name.
		nameX := typeX + dims.TypeColWidth + colPad
		builder.text(nameX, textY, attr.Name,
			"fill", th.TextColor,
			"font-size", fmtFloat(th.FontSize*erAttrFontScale),
		)
	}
}
