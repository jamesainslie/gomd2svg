package render

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Class diagram rendering constants.
const (
	classMemberFontScale     float32 = 0.85
	classLineHeightScale     float32 = 1.2
	classHeaderBaselineShift float32 = 0.75
)

// renderClass renders all class diagram elements: edges and UML class nodes.
func renderClass(builder *svgBuilder, computed *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	cd, ok := computed.Diagram.(layout.ClassData)
	if !ok {
		renderGraph(builder, computed, th, cfg)
		return
	}

	// Render edges with class-specific markers.
	renderClassEdges(builder, computed, th)

	// Render class nodes as UML compartment boxes.
	renderClassNodes(builder, computed, &cd, th, cfg)
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
func renderClassEdges(builder *svgBuilder, computed *layout.Layout, th *theme.Theme) {
	for edgeIdx, edge := range computed.Edges {
		if len(edge.Points) < 2 {
			continue
		}

		pathData := pointsToPath(edge.Points)
		edgeID := fmt.Sprintf("edge-%d", edgeIdx)

		attrs := []string{
			"id", edgeID,
			"class", "edgePath",
			"d", pathData,
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

		builder.selfClose("path", attrs...)

		// Render edge label if present.
		if edge.Label != nil && len(edge.Label.Lines) > 0 {
			renderEdgeLabel(builder, edge, th)
		}
	}
}

// renderClassNodes renders class nodes as UML compartment boxes.
func renderClassNodes(builder *svgBuilder, computed *layout.Layout, cd *layout.ClassData, th *theme.Theme, cfg *config.Layout) {
	// Sort node IDs for deterministic rendering order.
	ids := make([]string, 0, len(computed.Nodes))
	for id := range computed.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		node := computed.Nodes[id]
		comp, hasComp := cd.Compartments[id]
		members := cd.Members[id]

		if !hasComp || members == nil {
			// No compartment data; render as a simple node.
			fill := th.PrimaryColor
			if node.Style.Fill != nil {
				fill = *node.Style.Fill
			}
			stroke := th.PrimaryBorderColor
			if node.Style.Stroke != nil {
				stroke = *node.Style.Stroke
			}
			textColor := th.PrimaryTextColor
			if node.Style.TextColor != nil {
				textColor = *node.Style.TextColor
			}
			renderNodeShape(builder, node, fill, stroke, textColor)
			continue
		}

		annotation := cd.Annotations[id]
		renderClassNode(builder, node, members, comp, annotation, th, cfg)
	}
}

// renderClassNode renders a single UML class box with header, attributes, and methods compartments.
func renderClassNode(builder *svgBuilder, node *layout.NodeLayout, members *ir.ClassMembers, comp layout.ClassCompartment, annotation string, th *theme.Theme, cfg *config.Layout) {
	// Compute top-left corner from center coordinates.
	posX := node.X - node.Width/2
	posY := node.Y - node.Height/2
	width := node.Width
	height := node.Height

	fontSize := node.Label.FontSize
	if fontSize <= 0 {
		fontSize = th.FontSize
	}
	memberFontSize := cfg.Class.MemberFontSize
	if memberFontSize <= 0 {
		memberFontSize = fontSize * classMemberFontScale
	}
	lineH := fontSize * classLineHeightScale
	memberLineH := memberFontSize * classLineHeightScale

	padX := cfg.Class.CompartmentPadX
	if padX <= 0 {
		padX = 12
	}

	// Outer rounded rect (full node dimensions).
	builder.rect(posX, posY, width, height, defaultRectRadius,
		"fill", th.ClassBodyBg,
		"stroke", th.ClassBorder,
		"stroke-width", "1",
	)

	// --- Header section ---
	// Header background rect.
	builder.rect(posX, posY, width, comp.HeaderHeight, defaultRectRadius,
		"fill", th.ClassHeaderBg,
		"stroke", th.ClassBorder,
		"stroke-width", "1",
	)

	// Position text within header.
	headerTextY := posY + comp.HeaderHeight/2
	textLines := 1
	if annotation != "" {
		textLines = 2
	}

	totalTextH := lineH * float32(textLines)
	startY := headerTextY - totalTextH/2 + lineH*classHeaderBaselineShift

	lineIdx := 0
	if annotation != "" {
		annotText := "\u00AB" + annotation + "\u00BB" // guillemets
		builder.text(posX+width/2, startY+float32(lineIdx)*lineH, annotText,
			"text-anchor", "middle",
			"fill", th.PrimaryTextColor,
			"font-size", fmtFloat(fontSize*classMemberFontScale),
			"font-style", "italic",
		)
		lineIdx++
	}

	// Class name (centered, bold).
	builder.text(posX+width/2, startY+float32(lineIdx)*lineH, node.Label.Lines[0],
		"text-anchor", "middle",
		"fill", th.PrimaryTextColor,
		"font-size", fmtFloat(fontSize),
		"font-weight", "bold",
	)

	// --- Divider line after header ---
	dividerY := posY + comp.HeaderHeight
	builder.line(posX, dividerY, posX+width, dividerY,
		"stroke", th.ClassBorder,
		"stroke-width", "1",
	)

	// --- Attributes section ---
	attrY := dividerY
	for idx, attr := range members.Attributes {
		text := attr.Visibility.Symbol() + attr.Type + " " + attr.Name
		ty := attrY + memberLineH*float32(idx+1)
		builder.text(posX+padX, ty, text,
			"text-anchor", "start",
			"fill", th.TextColor,
			"font-size", fmtFloat(memberFontSize),
		)
	}

	// --- Divider line after attributes ---
	dividerY2 := dividerY + comp.AttributeHeight
	builder.line(posX, dividerY2, posX+width, dividerY2,
		"stroke", th.ClassBorder,
		"stroke-width", "1",
	)

	// --- Methods section ---
	methY := dividerY2
	for idx, meth := range members.Methods {
		text := meth.Visibility.Symbol() + meth.Name + "(" + meth.Params + ")"
		if meth.Type != "" {
			text += " : " + meth.Type
		}
		ty := methY + memberLineH*float32(idx+1)
		builder.text(posX+padX, ty, text,
			"text-anchor", "start",
			"fill", th.TextColor,
			"font-size", fmtFloat(memberFontSize),
		)
	}
}
