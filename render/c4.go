package render

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// renderC4 renders all C4 diagram elements: boundaries, edges, and element nodes
// with type-specific colors and person icons.
func renderC4(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	cd, ok := l.Diagram.(layout.C4Data)
	if !ok {
		return
	}

	smallFontSize := th.FontSize * 0.85
	lineH := th.FontSize * cfg.LabelLineHeight
	smallLineH := smallFontSize * cfg.LabelLineHeight

	// 1. Render boundaries first (behind everything).
	for _, boundary := range cd.Boundaries {
		// Dashed rectangle.
		b.rect(boundary.X, boundary.Y, boundary.Width, boundary.Height, 4,
			"fill", "none",
			"stroke", th.C4BoundaryColor,
			"stroke-width", "1",
			"stroke-dasharray", "5,5",
		)
		// Boundary label at top-left.
		b.text(boundary.X+8, boundary.Y+16, boundary.Label,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize*0.9),
			"fill", th.C4BoundaryColor,
			"font-weight", "bold",
		)
		// Type subtitle.
		if boundary.Type != "" {
			b.text(boundary.X+8, boundary.Y+30, "["+boundary.Type+"]",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize*0.8),
				"fill", th.C4BoundaryColor,
			)
		}
	}

	// 2. Render edges.
	renderEdges(b, l, th)

	// 3. Render nodes (elements) sorted by ID for deterministic order.
	ids := make([]string, 0, len(l.Nodes))
	for id := range l.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		n := l.Nodes[id]
		elem := cd.Elements[id]

		color := c4ElementColor(elem, th)

		// Top-left from center coordinates.
		x := n.X - n.Width/2
		y := n.Y - n.Height/2

		if elem != nil && elem.Type.IsPerson() {
			renderC4Person(b, x, y, n.Width, n.Height, color, th)
		} else {
			// Rounded rectangle for non-person elements.
			b.rect(x, y, n.Width, n.Height, 6,
				"fill", color,
				"stroke", "none",
			)
		}

		// Render text inside the element.
		cx := n.X
		curY := y + n.Height/2 - lineH*0.25

		// For person elements, shift text down to account for the icon.
		if elem != nil && elem.Type.IsPerson() {
			curY = y + 50
		}

		// Label (bold, white).
		b.text(cx, curY, n.Label.Lines[0],
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize),
			"font-weight", "bold",
			"fill", th.C4TextColor,
		)
		curY += lineH

		if elem != nil {
			// Technology line in brackets.
			if elem.Technology != "" {
				b.text(cx, curY, fmt.Sprintf("[%s]", elem.Technology),
					"text-anchor", "middle",
					"font-family", th.FontFamily,
					"font-size", fmtFloat(smallFontSize),
					"fill", th.C4TextColor,
				)
				curY += smallLineH
			}

			// Description line.
			if elem.Description != "" {
				b.text(cx, curY, elem.Description,
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
func renderC4Person(b *svgBuilder, x, y, w, h float32, color string, th *theme.Theme) {
	// Body rectangle (rounded).
	b.rect(x, y, w, h, 6,
		"fill", color,
		"stroke", "none",
	)

	// Person icon: head circle.
	cx := x + w/2
	headR := float32(12)
	headCY := y + 18
	b.circle(cx, headCY, headR,
		"fill", th.C4TextColor,
	)

	// Person icon: body arc (simple path).
	bodyTop := headCY + headR + 2
	bodyW := float32(20)
	d := fmt.Sprintf("M %s,%s Q %s,%s %s,%s",
		fmtFloat(cx-bodyW), fmtFloat(bodyTop),
		fmtFloat(cx), fmtFloat(bodyTop+24),
		fmtFloat(cx+bodyW), fmtFloat(bodyTop),
	)
	b.path(d,
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
	case ir.C4Container_, ir.C4ContainerDb, ir.C4ContainerQueue:
		return th.C4ContainerColor
	case ir.C4Component_:
		return th.C4ComponentColor
	default:
		return th.C4SystemColor
	}
}
