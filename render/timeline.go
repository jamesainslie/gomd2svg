package render

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Timeline rendering constants.
const (
	timelineTitlePadding    float32 = 5
	timelineSectionRadius   float32 = 4
	timelineSectionLabelOff float32 = 10
	timelinePeriodLabelGap  float32 = 4
	timelineEventPadding    float32 = 4
	timelineEventRadius     float32 = 12
)

func renderTimeline(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, _ *config.Layout) {
	td, ok := lay.Diagram.(layout.TimelineData)
	if !ok {
		return
	}

	// Title.
	if td.Title != "" {
		builder.text(lay.Width/2, th.FontSize+timelineTitlePadding, td.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize+2),
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	for _, sec := range td.Sections {
		// Section background.
		builder.rect(sec.X, sec.Y, sec.Width, sec.Height, timelineSectionRadius,
			"fill", sec.Color,
			"stroke", "none",
		)

		// Section label.
		if sec.Title != "" {
			builder.text(sec.X+timelineSectionLabelOff, sec.Y+sec.Height/2, sec.Title,
				"dominant-baseline", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize),
				"font-weight", "bold",
				"fill", th.TextColor,
			)
		}

		for _, period := range sec.Periods {
			// Period column separator.
			builder.rect(period.X, period.Y, period.Width, period.Height, 0,
				"fill", "none",
				"stroke", th.TimelineEventBorder,
				"stroke-width", "0.5",
				"stroke-opacity", "0.3",
			)

			// Period title.
			builder.text(period.X+period.Width/2, period.Y-timelinePeriodLabelGap, period.Title,
				"text-anchor", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize-1),
				"font-weight", "bold",
				"fill", th.TextColor,
			)

			// Events.
			for _, event := range period.Events {
				builder.rect(event.X, event.Y+2, event.Width, event.Height-timelineEventPadding, timelineEventRadius,
					"fill", th.TimelineEventFill,
					"stroke", th.TimelineEventBorder,
					"stroke-width", "1",
				)
				builder.text(event.X+event.Width/2, event.Y+event.Height/2, event.Text,
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
