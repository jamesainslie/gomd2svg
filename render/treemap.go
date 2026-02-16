package render

import (
	"fmt"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderTreemap(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	td, ok := l.Diagram.(layout.TreemapData)
	if !ok {
		return
	}

	colors := th.TreemapColors
	if len(colors) == 0 {
		colors = []string{"#4C78A8", "#72B7B2", "#EECA3B", "#F58518"}
	}

	// Draw title if present.
	if td.Title != "" {
		b.text(l.Width/2, 20, td.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", "16",
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	// Draw rectangles.
	for _, r := range td.Rects {
		color := colors[r.ColorIndex%len(colors)]

		if r.IsSection {
			// Section: draw a container rect with header.
			b.rect(r.X, r.Y, r.Width, r.Height, 2,
				"fill", "none",
				"stroke", th.TreemapBorder,
				"stroke-width", "1",
			)
			// Header background.
			headerH := cfg.Treemap.HeaderHeight
			if headerH > r.Height {
				headerH = r.Height
			}
			b.rect(r.X, r.Y, r.Width, headerH, 0,
				"fill", color,
				"opacity", "0.3",
			)
			// Section label.
			b.text(r.X+4, r.Y+headerH-6, r.Label,
				"font-family", th.FontFamily,
				"font-size", fmtFloat(cfg.Treemap.LabelFontSize),
				"fill", th.TextColor,
			)
		} else {
			// Leaf: fill with color, add label and value.
			b.rect(r.X, r.Y, r.Width, r.Height, 2,
				"fill", color,
				"stroke", th.TreemapBorder,
				"stroke-width", "1",
			)

			// Only draw label if rect is large enough.
			if r.Width > 20 && r.Height > 14 {
				cx := r.X + r.Width/2
				cy := r.Y + r.Height/2

				b.text(cx, cy, r.Label,
					"text-anchor", "middle",
					"font-family", th.FontFamily,
					"font-size", fmtFloat(cfg.Treemap.LabelFontSize),
					"fill", th.TreemapTextColor,
				)

				// Show value below label if there's room.
				if r.Height > 30 && r.Value > 0 {
					b.text(cx, cy+cfg.Treemap.ValueFontSize+2, fmt.Sprintf("%.0f", r.Value),
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
