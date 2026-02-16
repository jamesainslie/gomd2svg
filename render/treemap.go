package render

import (
	"fmt"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Treemap rendering constants.
const (
	treemapTitleOffsetY     float32 = 20
	treemapSectionLabelPadX float32 = 4
	treemapSectionLabelPadY float32 = 6
)

func renderTreemap(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	td, ok := lay.Diagram.(layout.TreemapData)
	if !ok {
		return
	}

	colors := th.TreemapColors
	if len(colors) == 0 {
		colors = []string{"#4C78A8", "#72B7B2", "#EECA3B", "#F58518"}
	}

	// Draw title if present.
	if td.Title != "" {
		builder.text(lay.Width/2, treemapTitleOffsetY, td.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", "16",
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	// Draw rectangles.
	for _, rect := range td.Rects {
		color := colors[rect.ColorIndex%len(colors)]

		if rect.IsSection {
			// Section: draw a container rect with header.
			builder.rect(rect.X, rect.Y, rect.Width, rect.Height, 2,
				"fill", "none",
				"stroke", th.TreemapBorder,
				"stroke-width", "1",
			)
			// Header background.
			headerH := cfg.Treemap.HeaderHeight
			if headerH > rect.Height {
				headerH = rect.Height
			}
			builder.rect(rect.X, rect.Y, rect.Width, headerH, 0,
				"fill", color,
				"opacity", "0.3",
			)
			// Section label.
			builder.text(rect.X+treemapSectionLabelPadX, rect.Y+headerH-treemapSectionLabelPadY, rect.Label,
				"font-family", th.FontFamily,
				"font-size", fmtFloat(cfg.Treemap.LabelFontSize),
				"fill", th.TextColor,
			)
		} else {
			// Leaf: fill with color, add label and value.
			builder.rect(rect.X, rect.Y, rect.Width, rect.Height, 2,
				"fill", color,
				"stroke", th.TreemapBorder,
				"stroke-width", "1",
			)

			// Only draw label if rect is large enough.
			if rect.Width > 20 && rect.Height > 14 {
				cx := rect.X + rect.Width/2
				cy := rect.Y + rect.Height/2

				builder.text(cx, cy, rect.Label,
					"text-anchor", "middle",
					"font-family", th.FontFamily,
					"font-size", fmtFloat(cfg.Treemap.LabelFontSize),
					"fill", th.TreemapTextColor,
				)

				// Show value below label if there's room.
				if rect.Height > 30 && rect.Value > 0 {
					builder.text(cx, cy+cfg.Treemap.ValueFontSize+2, fmt.Sprintf("%.0f", rect.Value),
						"text-anchor", "middle",
						"font-family", th.FontFamily,
						"font-size", fmtFloat(cfg.Treemap.ValueFontSize),
						"fill", th.TreemapTextColor,
						"opacity", "0.7",
					)
				}
			}
		}
	}
}
