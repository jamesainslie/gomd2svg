package render

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderTimeline(b *svgBuilder, l *layout.Layout, th *theme.Theme, _ *config.Layout) {
	td, ok := l.Diagram.(layout.TimelineData)
	if !ok {
		return
	}

	// Title.
	if td.Title != "" {
		b.text(l.Width/2, th.FontSize+5, td.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize+2),
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	for _, sec := range td.Sections {
		// Section background.
		b.rect(sec.X, sec.Y, sec.Width, sec.Height, 4,
			"fill", sec.Color,
			"stroke", "none",
		)

		// Section label.
		if sec.Title != "" {
			b.text(sec.X+10, sec.Y+sec.Height/2, sec.Title,
				"dominant-baseline", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize),
				"font-weight", "bold",
				"fill", th.TextColor,
			)
		}

		for _, p := range sec.Periods {
			// Period column separator.
			b.rect(p.X, p.Y, p.Width, p.Height, 0,
				"fill", "none",
				"stroke", th.TimelineEventBorder,
				"stroke-width", "0.5",
				"stroke-opacity", "0.3",
			)

			// Period title.
			b.text(p.X+p.Width/2, p.Y-4, p.Title,
				"text-anchor", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize-1),
				"font-weight", "bold",
				"fill", th.TextColor,
			)

			// Events.
			for _, e := range p.Events {
				b.rect(e.X, e.Y+2, e.Width, e.Height-4, 12,
					"fill", th.TimelineEventFill,
					"stroke", th.TimelineEventBorder,
					"stroke-width", "1",
				)
				b.text(e.X+e.Width/2, e.Y+e.Height/2, e.Text,
					"text-anchor", "middle",
					"dominant-baseline", "middle",
					"font-family", th.FontFamily,
					"font-size", fmtFloat(th.FontSize-2),
					"fill", "#FFFFFF",
				)
			}
		}
	}
}
