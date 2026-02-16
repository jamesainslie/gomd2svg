package render

import (
	"strconv"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Journey diagram rendering constants.
const (
	journeyScoreMaxIdx     int     = 4
	journeyScoreRange      float32 = 4.0
	journeyScoreLabelGap   float32 = 15
	journeySectionRadius   float32 = 4
	journeySectionLabelGap float32 = 5
	journeyTaskRadius      float32 = 6
	journeyIndicatorOffset float32 = 10
	journeyIndicatorRadius float32 = 5
	journeyTaskTextOffset  float32 = 5
	journeyLegendGap       float32 = 20
	journeyLegendDotOffset float32 = 5
	journeyLegendDotR      float32 = 4
	journeyLegendTextOff   float32 = 15
	journeyTextBaselineOff float32 = 4
)

func renderJourney(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	data, ok := lay.Diagram.(layout.JourneyData)
	if !ok {
		return
	}

	// Title
	if data.Title != "" {
		builder.text(lay.Width/2, cfg.Journey.PaddingY, data.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize+2),
			"font-weight", "bold",
			"fill", th.JourneyTaskText,
		)
	}

	// Dashed horizontal score guidelines (1-5)
	for score := 1; score <= 5; score++ {
		scoreRatio := float32(score-1) / journeyScoreRange
		posY := data.TrackY + data.TrackH*(1-scoreRatio)
		builder.line(cfg.Journey.PaddingX, posY, lay.Width-cfg.Journey.PaddingX, posY,
			"stroke", "#ddd",
			"stroke-dasharray", "4,4",
		)
		// Score label on left
		builder.text(cfg.Journey.PaddingX-journeyScoreLabelGap, posY+journeyTextBaselineOff, strconv.Itoa(score),
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", "10",
			"fill", "#999",
		)
	}

	// Section backgrounds
	for _, sec := range data.Sections {
		fill := sec.Color
		if fill == "" {
			fill = "#f5f5f5"
		}
		builder.rect(sec.X, sec.Y, sec.Width, sec.Height, journeySectionRadius,
			"fill", fill,
			"stroke", "none",
		)
		// Section label at top
		if sec.Label != "" {
			builder.text(sec.X+sec.Width/2, sec.Y-journeySectionLabelGap, sec.Label,
				"text-anchor", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize-2),
				"fill", th.JourneyTaskText,
			)
		}

		// Task rectangles
		for _, task := range sec.Tasks {
			// Score-colored indicator
			scoreIdx := task.Score - 1
			if scoreIdx < 0 {
				scoreIdx = 0
			}
			if scoreIdx > journeyScoreMaxIdx {
				scoreIdx = journeyScoreMaxIdx
			}
			scoreColor := th.JourneyScoreColors[scoreIdx]

			// Task rectangle
			tx := task.X - task.Width/2
			ty := task.Y - task.Height/2
			builder.rect(tx, ty, task.Width, task.Height, journeyTaskRadius,
				"fill", th.JourneyTaskFill,
				"stroke", th.JourneyTaskBorder,
				"stroke-width", "1",
			)

			// Score indicator circle
			builder.circle(tx+journeyIndicatorOffset, task.Y, journeyIndicatorRadius,
				"fill", scoreColor,
				"stroke", scoreColor,
			)

			// Task label
			builder.text(task.X+journeyTaskTextOffset, task.Y+journeyTextBaselineOff, task.Label,
				"text-anchor", "middle",
				"dominant-baseline", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize-2),
				"fill", th.JourneyTaskText,
			)
		}
	}

	// Actor legend at bottom
	if len(data.Actors) > 0 {
		legendY := data.TrackY + data.TrackH + journeyLegendGap
		legendX := cfg.Journey.PaddingX
		for _, actor := range data.Actors {
			colorIdx := actor.ColorIndex
			color := "#666"
			if len(th.JourneySectionColors) > 0 {
				color = th.JourneySectionColors[colorIdx%len(th.JourneySectionColors)]
			}
			builder.circle(legendX+journeyLegendDotOffset, legendY, journeyLegendDotR,
				"fill", color,
				"stroke", color,
			)
			builder.text(legendX+journeyLegendTextOff, legendY+journeyTextBaselineOff, actor.Name,
				"font-family", th.FontFamily,
				"font-size", "11",
				"fill", th.JourneyTaskText,
			)
			legendX += 80
		}
	}
}
