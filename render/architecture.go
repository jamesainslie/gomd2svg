package render

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

// renderArchitecture renders all architecture diagram elements: groups,
// edges, service nodes, and junctions.
func renderArchitecture(b *svgBuilder, l *layout.Layout, th *theme.Theme, _ *config.Layout) {
	data, ok := l.Diagram.(layout.ArchitectureData)
	if !ok {
		return
	}

	// 1. Render groups first (behind everything).
	for _, grp := range data.Groups {
		if grp.Width <= 0 || grp.Height <= 0 {
			continue
		}
		// Dashed border group rectangle.
		b.rect(grp.X, grp.Y, grp.Width, grp.Height, 8,
			"fill", th.ArchGroupFill,
			"stroke", th.ArchGroupBorder,
			"stroke-width", "1",
			"stroke-dasharray", "6,3",
		)
		// Group label at top-left.
		if grp.Label != "" {
			b.text(grp.X+10, grp.Y+18, grp.Label,
				"text-anchor", "start",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize*0.9),
				"font-weight", "bold",
				"fill", th.ArchGroupText,
			)
		}
		// Simple icon next to label at top-right corner.
		if grp.Icon != "" {
			renderArchIcon(b, grp.Icon, grp.X+grp.Width-20, grp.Y+14, 10)
		}
	}

	// 2. Render edges using architecture-specific edge color.
	for i, edge := range l.Edges {
		if len(edge.Points) < 2 {
			continue
		}

		d := pointsToPath(edge.Points)
		edgeID := fmt.Sprintf("edge-%d", i)

		strokeColor := th.ArchEdgeColor
		strokeWidth := "1.5"

		attrs := []string{
			"id", edgeID,
			"class", "edgePath",
			"d", d,
			"fill", "none",
			"stroke", strokeColor,
			"stroke-width", strokeWidth,
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

		// Arrow marker references.
		if edge.ArrowEnd {
			attrs = append(attrs, "marker-end", "url(#arrowhead)")
		}
		if edge.ArrowStart {
			attrs = append(attrs, "marker-start", "url(#arrowhead-start)")
		}

		b.selfClose("path", attrs...)

		// Render edge label if present.
		if edge.Label != nil && len(edge.Label.Lines) > 0 {
			renderEdgeLabel(b, edge, th)
		}
	}

	// 3. Render service nodes sorted by ID for deterministic output.
	junctionSet := make(map[string]bool, len(data.Junctions))
	for _, junc := range data.Junctions {
		junctionSet[junc.ID] = true
	}

	ids := make([]string, 0, len(l.Nodes))
	for id := range l.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		if junctionSet[id] {
			continue
		}

		n := l.Nodes[id]

		// Service rectangle.
		nx := n.X - n.Width/2
		ny := n.Y - n.Height/2
		b.rect(nx, ny, n.Width, n.Height, 6,
			"fill", th.ArchServiceFill,
			"stroke", th.ArchServiceBorder,
			"stroke-width", "1.5",
		)

		// Render icon above label if available.
		svcInfo, hasSvcInfo := data.Services[id]
		if hasSvcInfo && svcInfo.Icon != "" {
			renderArchIcon(b, svcInfo.Icon, n.X, n.Y-10, 12)
		}

		// Render label centered in the service box.
		labelY := n.Y + 8
		if hasSvcInfo && svcInfo.Icon != "" {
			labelY = n.Y + 14 // shift down when icon is present
		}
		if len(n.Label.Lines) > 0 {
			b.text(n.X, labelY, n.Label.Lines[0],
				"text-anchor", "middle",
				"dominant-baseline", "auto",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize),
				"fill", th.ArchServiceText,
			)
		}
	}

	// 4. Render junctions as small filled circles.
	for _, junc := range data.Junctions {
		b.circle(junc.X, junc.Y, junc.Size/2,
			"fill", th.ArchJunctionFill,
		)
	}
}

// renderArchIcon renders a simple SVG icon shape at the given center position.
func renderArchIcon(b *svgBuilder, icon string, cx, cy, size float32) {
	half := size / 2
	switch icon {
	case "database":
		// Simple cylinder representation: ellipse.
		b.ellipse(cx, cy, half, half*0.6,
			"fill", "#78909C",
			"stroke", "none",
		)
	case "server":
		// Simple box.
		b.rect(cx-half, cy-half, size, size, 2,
			"fill", "#78909C",
			"stroke", "none",
		)
		// Two horizontal lines inside the box.
		third := size / 3
		b.line(cx-half+2, cy-half+third, cx+half-2, cy-half+third,
			"stroke", "#fff",
			"stroke-width", "1",
		)
		b.line(cx-half+2, cy-half+2*third, cx+half-2, cy-half+2*third,
			"stroke", "#fff",
			"stroke-width", "1",
		)
	case "cloud":
		// Cloud-like ellipse.
		b.ellipse(cx, cy, half*1.2, half*0.7,
			"fill", "#78909C",
			"stroke", "none",
		)
	case "internet":
		// Globe: circle with a vertical and horizontal cross line.
		b.circle(cx, cy, half,
			"fill", "none",
			"stroke", "#78909C",
			"stroke-width", "1.5",
		)
		b.line(cx-half, cy, cx+half, cy,
			"stroke", "#78909C",
			"stroke-width", "1",
		)
		b.line(cx, cy-half, cx, cy+half,
			"stroke", "#78909C",
			"stroke-width", "1",
		)
		// Curved arc approximation for globe effect.
		d := fmt.Sprintf("M %s,%s Q %s,%s %s,%s",
			fmtFloat(cx), fmtFloat(cy-half),
			fmtFloat(cx+half*0.5), fmtFloat(cy),
			fmtFloat(cx), fmtFloat(cy+half),
		)
		b.path(d,
			"fill", "none",
			"stroke", "#78909C",
			"stroke-width", "1",
		)
	case "disk":
		// Filled circle.
		b.circle(cx, cy, half,
			"fill", "#78909C",
			"stroke", "none",
		)
	}
}
