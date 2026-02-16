package render

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Architecture diagram rendering constants.
const (
	archGroupBorderRadius   float32 = 8
	archGroupLabelOffsetX   float32 = 10
	archGroupLabelOffsetY   float32 = 18
	archGroupFontScale      float32 = 0.9
	archGroupIconOffsetX    float32 = 20
	archGroupIconOffsetY    float32 = 14
	archGroupIconSize       float32 = 10
	archServiceBorderRadius float32 = 6
	archServiceIconOffsetY  float32 = 10
	archServiceIconSize     float32 = 12
	archLabelOffsetY        float32 = 8
	archLabelIconOffsetY    float32 = 14
	archEllipseYScale       float32 = 0.6
	archCloudXScale         float32 = 1.2
	archCloudYScale         float32 = 0.7
)

// renderArchitecture renders all architecture diagram elements: groups,
// edges, service nodes, and junctions.
func renderArchitecture(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, _ *config.Layout) {
	data, ok := lay.Diagram.(layout.ArchitectureData)
	if !ok {
		return
	}

	// 1. Render groups first (behind everything).
	for _, grp := range data.Groups {
		if grp.Width <= 0 || grp.Height <= 0 {
			continue
		}
		// Dashed border group rectangle.
		builder.rect(grp.X, grp.Y, grp.Width, grp.Height, archGroupBorderRadius,
			"fill", th.ArchGroupFill,
			"stroke", th.ArchGroupBorder,
			"stroke-width", "1",
			"stroke-dasharray", "6,3",
		)
		// Group label at top-left.
		if grp.Label != "" {
			builder.text(grp.X+archGroupLabelOffsetX, grp.Y+archGroupLabelOffsetY, grp.Label,
				"text-anchor", "start",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize*archGroupFontScale),
				"font-weight", "bold",
				"fill", th.ArchGroupText,
			)
		}
		// Simple icon next to label at top-right corner.
		if grp.Icon != "" {
			renderArchIcon(builder, grp.Icon, grp.X+grp.Width-archGroupIconOffsetX, grp.Y+archGroupIconOffsetY, archGroupIconSize)
		}
	}

	// 2. Render edges using architecture-specific edge color.
	for idx, edge := range lay.Edges {
		if len(edge.Points) < 2 {
			continue
		}

		pathData := pointsToPath(edge.Points)
		edgeID := fmt.Sprintf("edge-%d", idx)

		strokeColor := th.ArchEdgeColor
		strokeWidth := "1.5"

		attrs := []string{
			"id", edgeID,
			"class", "edgePath",
			"d", pathData,
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

		builder.selfClose("path", attrs...)

		// Render edge label if present.
		if edge.Label != nil && len(edge.Label.Lines) > 0 {
			renderEdgeLabel(builder, edge, th)
		}
	}

	// 3. Render service nodes sorted by ID for deterministic output.
	junctionSet := make(map[string]bool, len(data.Junctions))
	for _, junc := range data.Junctions {
		junctionSet[junc.ID] = true
	}

	ids := make([]string, 0, len(lay.Nodes))
	for id := range lay.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		if junctionSet[id] {
			continue
		}

		node := lay.Nodes[id]

		// Service rectangle.
		nx := node.X - node.Width/2
		ny := node.Y - node.Height/2
		builder.rect(nx, ny, node.Width, node.Height, archServiceBorderRadius,
			"fill", th.ArchServiceFill,
			"stroke", th.ArchServiceBorder,
			"stroke-width", "1.5",
		)

		// Render icon above label if available.
		svcInfo, hasSvcInfo := data.Services[id]
		if hasSvcInfo && svcInfo.Icon != "" {
			renderArchIcon(builder, svcInfo.Icon, node.X, node.Y-archServiceIconOffsetY, archServiceIconSize)
		}

		// Render label centered in the service box.
		labelY := node.Y + archLabelOffsetY
		if hasSvcInfo && svcInfo.Icon != "" {
			labelY = node.Y + archLabelIconOffsetY // shift down when icon is present
		}
		if len(node.Label.Lines) > 0 {
			builder.text(node.X, labelY, node.Label.Lines[0],
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
		builder.circle(junc.X, junc.Y, junc.Size/2,
			"fill", th.ArchJunctionFill,
		)
	}
}

// renderArchIcon renders a simple SVG icon shape at the given center position.
func renderArchIcon(builder *svgBuilder, icon string, cx, cy, size float32) {
	half := size / 2
	switch icon {
	case "database":
		// Simple cylinder representation: ellipse.
		builder.ellipse(cx, cy, half, half*archEllipseYScale,
			"fill", "#78909C",
			"stroke", "none",
		)
	case "server":
		// Simple box.
		builder.rect(cx-half, cy-half, size, size, 2,
			"fill", "#78909C",
			"stroke", "none",
		)
		// Two horizontal lines inside the box.
		third := size / 3
		builder.line(cx-half+2, cy-half+third, cx+half-2, cy-half+third,
			"stroke", "#fff",
			"stroke-width", "1",
		)
		builder.line(cx-half+2, cy-half+2*third, cx+half-2, cy-half+2*third,
			"stroke", "#fff",
			"stroke-width", "1",
		)
	case "cloud":
		// Cloud-like ellipse.
		builder.ellipse(cx, cy, half*archCloudXScale, half*archCloudYScale,
			"fill", "#78909C",
			"stroke", "none",
		)
	case "internet":
		// Globe: circle with a vertical and horizontal cross line.
		builder.circle(cx, cy, half,
			"fill", "none",
			"stroke", "#78909C",
			"stroke-width", "1.5",
		)
		builder.line(cx-half, cy, cx+half, cy,
			"stroke", "#78909C",
			"stroke-width", "1",
		)
		builder.line(cx, cy-half, cx, cy+half,
			"stroke", "#78909C",
			"stroke-width", "1",
		)
		// Curved arc approximation for globe effect.
		pathData := fmt.Sprintf("M %s,%s Q %s,%s %s,%s",
			fmtFloat(cx), fmtFloat(cy-half),
			fmtFloat(cx+half*0.5), fmtFloat(cy),
			fmtFloat(cx), fmtFloat(cy+half),
		)
		builder.path(pathData,
			"fill", "none",
			"stroke", "#78909C",
			"stroke-width", "1",
		)
	case "disk":
		// Filled circle.
		builder.circle(cx, cy, half,
			"fill", "#78909C",
			"stroke", "none",
		)
	}
}
