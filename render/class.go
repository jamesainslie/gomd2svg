package render

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// renderClass renders all class diagram elements: edges and UML class nodes.
func renderClass(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	cd, ok := l.Diagram.(layout.ClassData)
	if !ok {
		renderGraph(b, l, th, cfg)
		return
	}

	// Render edges with class-specific markers.
	renderClassEdges(b, l, th)

	// Render class nodes as UML compartment boxes.
	renderClassNodes(b, l, &cd, th, cfg)
}

// markerRefForArrowKind returns the SVG marker ID for a given arrowhead kind.
func markerRefForArrowKind(kind *ir.EdgeArrowhead, start bool) string {
	base := "arrowhead"
	if kind != nil {
		switch *kind {
		case ir.ClosedTriangle:
			base = "marker-closed-triangle"
		case ir.FilledDiamond:
			base = "marker-filled-diamond"
		case ir.OpenDiamond:
			base = "marker-open-diamond"
		case ir.OpenTriangle, ir.ClassDependency, ir.Lollipop:
			base = "arrowhead"
		}
	}
	if start {
		return base + "-start"
	}
	return base
}

// renderClassEdges renders edges using arrowhead kind to pick the correct marker.
func renderClassEdges(b *svgBuilder, l *layout.Layout, th *theme.Theme) {
	for i, edge := range l.Edges {
		if len(edge.Points) < 2 {
			continue
		}

		d := pointsToPath(edge.Points)
		edgeID := fmt.Sprintf("edge-%d", i)

		attrs := []string{
			"id", edgeID,
			"class", "edgePath",
			"d", d,
			"fill", "none",
			"stroke", th.LineColor,
			"stroke-width", "1.5",
			"stroke-linecap", "round",
			"stroke-linejoin", "round",
		}

		// Edge style: dotted or thick.
		switch edge.Style {
		case ir.Dotted:
			attrs = append(attrs, "stroke-dasharray", "5,5")
		case ir.Thick:
			attrs = append(attrs, "stroke-width", "3")
		}

		// Arrow marker references using arrowhead kind.
		if edge.ArrowEnd {
			markerID := markerRefForArrowKind(edge.ArrowEndKind, false)
			attrs = append(attrs, "marker-end", "url(#"+markerID+")")
		}
		if edge.ArrowStart {
			markerID := markerRefForArrowKind(edge.ArrowStartKind, true)
			attrs = append(attrs, "marker-start", "url(#"+markerID+")")
		}

		b.selfClose("path", attrs...)

		// Render edge label if present.
		if edge.Label != nil && len(edge.Label.Lines) > 0 {
			renderEdgeLabel(b, edge, th)
		}
	}
}

// renderClassNodes renders class nodes as UML compartment boxes.
func renderClassNodes(b *svgBuilder, l *layout.Layout, cd *layout.ClassData, th *theme.Theme, cfg *config.Layout) {
	// Sort node IDs for deterministic rendering order.
	ids := make([]string, 0, len(l.Nodes))
	for id := range l.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		n := l.Nodes[id]
		comp, hasComp := cd.Compartments[id]
		members := cd.Members[id]

		if !hasComp || members == nil {
			// No compartment data; render as a simple node.
			fill := th.PrimaryColor
			if n.Style.Fill != nil {
				fill = *n.Style.Fill
			}
			stroke := th.PrimaryBorderColor
			if n.Style.Stroke != nil {
				stroke = *n.Style.Stroke
			}
			textColor := th.PrimaryTextColor
			if n.Style.TextColor != nil {
				textColor = *n.Style.TextColor
			}
			renderNodeShape(b, n, fill, stroke, textColor)
			continue
		}

		annotation := cd.Annotations[id]
		renderClassNode(b, n, members, comp, annotation, th, cfg)
	}
}

// renderClassNode renders a single UML class box with header, attributes, and methods compartments.
func renderClassNode(b *svgBuilder, n *layout.NodeLayout, members *ir.ClassMembers, comp layout.ClassCompartment, annotation string, th *theme.Theme, cfg *config.Layout) {
	// Compute top-left corner from center coordinates.
	x := n.X - n.Width/2
	y := n.Y - n.Height/2
	w := n.Width
	h := n.Height

	fontSize := n.Label.FontSize
	if fontSize <= 0 {
		fontSize = th.FontSize
	}
	memberFontSize := cfg.Class.MemberFontSize
	if memberFontSize <= 0 {
		memberFontSize = fontSize * 0.85
	}
	lineH := fontSize * 1.2
	memberLineH := memberFontSize * 1.2

	padX := cfg.Class.CompartmentPadX
	if padX <= 0 {
		padX = 12
	}

	// Outer rounded rect (full node dimensions).
	b.rect(x, y, w, h, 3,
		"fill", th.ClassBodyBg,
		"stroke", th.ClassBorder,
		"stroke-width", "1",
	)

	// --- Header section ---
	// Header background rect.
	b.rect(x, y, w, comp.HeaderHeight, 3,
		"fill", th.ClassHeaderBg,
		"stroke", th.ClassBorder,
		"stroke-width", "1",
	)

	// Position text within header.
	headerTextY := y + comp.HeaderHeight/2
	textLines := 1
	if annotation != "" {
		textLines = 2
	}

	totalTextH := lineH * float32(textLines)
	startY := headerTextY - totalTextH/2 + lineH*0.75

	lineIdx := 0
	if annotation != "" {
		annotText := "\u00AB" + annotation + "\u00BB" // guillemets
		b.text(x+w/2, startY+float32(lineIdx)*lineH, annotText,
			"text-anchor", "middle",
			"fill", th.PrimaryTextColor,
			"font-size", fmtFloat(fontSize*0.85),
			"font-style", "italic",
		)
		lineIdx++
	}

	// Class name (centered, bold).
	b.text(x+w/2, startY+float32(lineIdx)*lineH, n.Label.Lines[0],
		"text-anchor", "middle",
		"fill", th.PrimaryTextColor,
		"font-size", fmtFloat(fontSize),
		"font-weight", "bold",
	)

	// --- Divider line after header ---
	dividerY := y + comp.HeaderHeight
	b.line(x, dividerY, x+w, dividerY,
		"stroke", th.ClassBorder,
		"stroke-width", "1",
	)

	// --- Attributes section ---
	attrY := dividerY
	for i, attr := range members.Attributes {
		text := attr.Visibility.Symbol() + attr.Type + " " + attr.Name
		ty := attrY + memberLineH*float32(i+1)
		b.text(x+padX, ty, text,
			"text-anchor", "start",
			"fill", th.TextColor,
			"font-size", fmtFloat(memberFontSize),
		)
	}

	// --- Divider line after attributes ---
	dividerY2 := dividerY + comp.AttributeHeight
	b.line(x, dividerY2, x+w, dividerY2,
		"stroke", th.ClassBorder,
		"stroke-width", "1",
	)

	// --- Methods section ---
	methY := dividerY2
	for i, meth := range members.Methods {
		text := meth.Visibility.Symbol() + meth.Name + "(" + meth.Params + ")"
		if meth.Type != "" {
			text += " : " + meth.Type
		}
		ty := methY + memberLineH*float32(i+1)
		b.text(x+padX, ty, text,
			"text-anchor", "start",
			"fill", th.TextColor,
			"font-size", fmtFloat(memberFontSize),
		)
	}
}
