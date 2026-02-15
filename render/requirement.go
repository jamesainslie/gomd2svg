package render

import (
	"sort"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// renderRequirement renders all requirement diagram elements: edges and
// requirement/element nodes with stereotype headers and metadata lines.
func renderRequirement(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	rd, ok := l.Diagram.(layout.RequirementData)
	if !ok {
		return
	}

	// Render edges (reuse shared edge rendering with arrow markers and labels).
	renderEdges(b, l, th)

	// Render nodes sorted by ID for deterministic output.
	ids := make([]string, 0, len(l.Nodes))
	for id := range l.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	metaFontSize := cfg.Requirement.MetadataFontSize
	if metaFontSize <= 0 {
		metaFontSize = 11
	}
	lineH := th.FontSize * cfg.LabelLineHeight
	metaLineH := metaFontSize * cfg.LabelLineHeight
	padX := cfg.Requirement.NodePadding
	if padX <= 0 {
		padX = 12
	}

	for _, id := range ids {
		n := l.Nodes[id]

		// Top-left from center coordinates.
		x := n.X - n.Width/2
		y := n.Y - n.Height/2

		// Outer rounded rectangle.
		b.rect(x, y, n.Width, n.Height, 4,
			"fill", th.RequirementFill,
			"stroke", th.RequirementBorder,
			"stroke-width", "1",
		)

		kind := rd.NodeKinds[id]
		curY := y + cfg.Padding.NodeVertical

		// Stereotype header line.
		stereotype := "\u00ABelement\u00BB"
		if kind == "requirement" {
			if req, ok := rd.Requirements[id]; ok {
				stereotype = "\u00AB" + req.Type.Stereotype() + "\u00BB"
			}
		}
		curY += lineH * 0.75
		b.text(x+n.Width/2, curY, stereotype,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(metaFontSize),
			"font-style", "italic",
			"fill", th.TextColor,
		)
		curY += lineH * 0.25

		// Divider line after stereotype.
		b.line(x, curY, x+n.Width, curY,
			"stroke", th.RequirementBorder,
			"stroke-width", "0.5",
		)

		// Name line (bold, centered).
		curY += lineH * 0.75
		name := ""
		if len(n.Label.Lines) > 0 {
			name = n.Label.Lines[0]
		}
		b.text(x+n.Width/2, curY, name,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize),
			"font-weight", "bold",
			"fill", th.TextColor,
		)
		curY += lineH * 0.25

		// Divider line after name.
		b.line(x, curY, x+n.Width, curY,
			"stroke", th.RequirementBorder,
			"stroke-width", "0.5",
		)

		// Metadata lines.
		if kind == "requirement" {
			if req, ok := rd.Requirements[id]; ok {
				renderRequirementMeta(b, x+padX, &curY, metaFontSize, metaLineH, th, req)
			}
		} else {
			if elem, ok := rd.Elements[id]; ok {
				renderElementMeta(b, x+padX, &curY, metaFontSize, metaLineH, th, elem)
			}
		}
	}
}

// renderRequirementMeta renders metadata lines for a requirement node.
func renderRequirementMeta(b *svgBuilder, x float32, curY *float32, fontSize, lineH float32, th *theme.Theme, req *ir.RequirementDef) {
	if req.ID != "" {
		*curY += lineH
		b.text(x, *curY, "Id: "+req.ID,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(fontSize),
			"fill", th.TextColor,
		)
	}
	if req.Text != "" {
		*curY += lineH
		b.text(x, *curY, "Text: "+req.Text,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(fontSize),
			"fill", th.TextColor,
		)
	}
	if req.Risk != ir.RiskNone {
		*curY += lineH
		b.text(x, *curY, "Risk: "+req.Risk.String(),
			"font-family", th.FontFamily,
			"font-size", fmtFloat(fontSize),
			"fill", th.TextColor,
		)
	}
	if req.VerifyMethod != ir.VerifyNone {
		*curY += lineH
		b.text(x, *curY, "Verify: "+req.VerifyMethod.String(),
			"font-family", th.FontFamily,
			"font-size", fmtFloat(fontSize),
			"fill", th.TextColor,
		)
	}
}

// renderElementMeta renders metadata lines for an element node.
func renderElementMeta(b *svgBuilder, x float32, curY *float32, fontSize, lineH float32, th *theme.Theme, elem *ir.ElementDef) {
	if elem.Type != "" {
		*curY += lineH
		b.text(x, *curY, "Type: "+elem.Type,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(fontSize),
			"fill", th.TextColor,
		)
	}
	if elem.DocRef != "" {
		*curY += lineH
		b.text(x, *curY, "Doc: "+elem.DocRef,
			"font-family", th.FontFamily,
			"font-size", fmtFloat(fontSize),
			"fill", th.TextColor,
		)
	}
}
