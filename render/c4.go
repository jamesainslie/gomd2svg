package render

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// C4 diagram rendering constants.
const (
	c4SmallFontScale      float32 = 0.85
	c4BoundaryRadius      float32 = 4
	c4BoundaryLabelOffset float32 = 8
	c4BoundaryLabelY      float32 = 16
	c4BoundaryTypeY       float32 = 30
	c4BoundaryFontScale   float32 = 0.9
	c4BoundaryTypeFScale  float32 = 0.8
	c4NodeBorderRadius    float32 = 6
	c4TextBaselineShift   float32 = 0.25
	c4PersonTextOffsetY   float32 = 50
	c4PersonNodeRadius    float32 = 6
	c4PersonHeadRadius    float32 = 12
	c4PersonHeadCenterY   float32 = 18
	c4PersonBodyWidth     float32 = 20
	c4PersonBodyArcHeight float32 = 24
)

// renderC4 renders all C4 diagram elements: boundaries, edges, and element nodes
// with type-specific colors and person icons.
func renderC4(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	cd, ok := lay.Diagram.(layout.C4Data)
	if !ok {
		return
	}

	smallFontSize := th.FontSize * c4SmallFontScale
	lineH := th.FontSize * cfg.LabelLineHeight
	smallLineH := smallFontSize * cfg.LabelLineHeight

	// 1. Render boundaries first (behind everything).
	for _, boundary := range cd.Boundaries {
		// Dashed rectangle.
		builder.rect(boundary.X, boundary.Y, boundary.Width, boundary.Height, c4BoundaryRadius,
			"fill", "none",
			"stroke", th.C4BoundaryColor,
			"stroke-width", "1",
			"stroke-dasharray", "5,5",
		)
		// Boundary label at top-left.
		builder.text(boundary.X+c4BoundaryLabelOffset, boundary.Y+c4BoundaryLabelY, boundary.Label,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize*c4BoundaryFontScale),
			"fill", th.C4BoundaryColor,
			"font-weight", "bold",
		)
		// Type subtitle.
		if boundary.Type != "" {
			builder.text(boundary.X+c4BoundaryLabelOffset, boundary.Y+c4BoundaryTypeY, "["+boundary.Type+"]",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize*c4BoundaryTypeFScale),
				"fill", th.C4BoundaryColor,
			)
		}
	}

	// 2. Render edges.
	renderEdges(builder, lay, th)

	// 3. Render nodes (elements) sorted by ID for deterministic order.
	ids := make([]string, 0, len(lay.Nodes))
	for id := range lay.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		node := lay.Nodes[id]
		elem := cd.Elements[id]

		color := c4ElementColor(elem, th)

		// Top-left from center coordinates.
		posX := node.X - node.Width/2
		posY := node.Y - node.Height/2

		if elem != nil && elem.Type.IsPerson() {
			renderC4Person(builder, posX, posY, node.Width, node.Height, color, th)
		} else {
			// Rounded rectangle for non-person elements.
			builder.rect(posX, posY, node.Width, node.Height, c4NodeBorderRadius,
				"fill", color,
				"stroke", "none",
			)
		}

		// Render text inside the element.
		cx := node.X
		curY := posY + node.Height/2 - lineH*c4TextBaselineShift

		// For person elements, shift text down to account for the icon.
		if elem != nil && elem.Type.IsPerson() {
			curY = posY + c4PersonTextOffsetY
		}

		// Label (bold, white).
		if len(node.Label.Lines) > 0 {
			builder.text(cx, curY, node.Label.Lines[0],
				"text-anchor", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize),
				"font-weight", "bold",
				"fill", th.C4TextColor,
			)
			curY += lineH
		}

		if elem != nil {
			// Technology line in brackets.
			if elem.Technology != "" {
				builder.text(cx, curY, fmt.Sprintf("[%s]", elem.Technology),
					"text-anchor", "middle",
					"font-family", th.FontFamily,
					"font-size", fmtFloat(smallFontSize),
					"fill", th.C4TextColor,
				)
				curY += smallLineH
			}

			// Description line.
			if elem.Description != "" {
				builder.text(cx, curY, elem.Description,
					"text-anchor", "middle",
					"font-family", th.FontFamily,
					"font-size", fmtFloat(smallFontSize),
					"fill", th.C4TextColor,
				)
			}
		}
	}
}

// renderC4Person draws a C4 person shape: a filled rounded rectangle with a
// person icon (circle head + arc body) centered at the top.
func renderC4Person(builder *svgBuilder, posX, posY, width, height float32, color string, th *theme.Theme) {
	// Body rectangle (rounded).
	builder.rect(posX, posY, width, height, c4PersonNodeRadius,
		"fill", color,
		"stroke", "none",
	)

	// Person icon: head circle.
	cx := posX + width/2
	headCY := posY + c4PersonHeadCenterY
	builder.circle(cx, headCY, c4PersonHeadRadius,
		"fill", th.C4TextColor,
	)

	// Person icon: body arc (simple path).
	bodyTop := headCY + c4PersonHeadRadius + 2
	pathData := fmt.Sprintf("M %s,%s Q %s,%s %s,%s",
		fmtFloat(cx-c4PersonBodyWidth), fmtFloat(bodyTop),
		fmtFloat(cx), fmtFloat(bodyTop+c4PersonBodyArcHeight),
		fmtFloat(cx+c4PersonBodyWidth), fmtFloat(bodyTop),
	)
	builder.path(pathData,
		"fill", "none",
		"stroke", th.C4TextColor,
		"stroke-width", "2",
		"stroke-linecap", "round",
	)
}

// c4ElementColor returns the fill color for a C4 element based on its type.
func c4ElementColor(elem *ir.C4Element, th *theme.Theme) string {
	if elem == nil {
		return th.C4SystemColor
	}
	if elem.Type.IsExternal() {
		return th.C4ExternalColor
	}
	if elem.Type.IsPerson() {
		return th.C4PersonColor
	}
	switch elem.Type {
	case ir.C4ContainerPlain, ir.C4ContainerDb, ir.C4ContainerQueue:
		return th.C4ContainerColor
	case ir.C4ComponentPlain:
		return th.C4ComponentColor
	default:
		return th.C4SystemColor
	}
}
