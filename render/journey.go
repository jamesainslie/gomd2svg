package render

import (
	"fmt"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func renderJourney(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	data, ok := l.Diagram.(layout.JourneyData)
	if !ok {
		return
	}

	// Title
	if data.Title != "" {
		b.text(l.Width/2, cfg.Journey.PaddingY, data.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize+2),
			"font-weight", "bold",
			"fill", th.JourneyTaskText,
		)
	}

	// Dashed horizontal score guidelines (1-5)
	for score := 1; score <= 5; score++ {
		scoreRatio := float32(score-1) / 4.0
		y := data.TrackY + data.TrackH*(1-scoreRatio)
		b.line(cfg.Journey.PaddingX, y, l.Width-cfg.Journey.PaddingX, y,
			"stroke", "#ddd",
			"stroke-dasharray", "4,4",
		)
		// Score label on left
		b.text(cfg.Journey.PaddingX-15, y+4, fmt.Sprintf("%d", score),
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
		b.rect(sec.X, sec.Y, sec.Width, sec.Height, 4,
			"fill", fill,
			"stroke", "none",
		)
		// Section label at top
		if sec.Label != "" {
			b.text(sec.X+sec.Width/2, sec.Y-5, sec.Label,
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
			if scoreIdx > 4 {
				scoreIdx = 4
			}
			scoreColor := th.JourneyScoreColors[scoreIdx]

			// Task rectangle
			tx := task.X - task.Width/2
			ty := task.Y - task.Height/2
			b.rect(tx, ty, task.Width, task.Height, 6,
				"fill", th.JourneyTaskFill,
				"stroke", th.JourneyTaskBorder,
				"stroke-width", "1",
			)

			// Score indicator circle
			b.circle(tx+10, task.Y, 5,
				"fill", scoreColor,
				"stroke", scoreColor,
			)

			// Task label
			b.text(task.X+5, task.Y+4, task.Label,
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
		legendY := data.TrackY + data.TrackH + 20
		legendX := cfg.Journey.PaddingX
		for _, actor := range data.Actors {
			colorIdx := actor.ColorIndex
			color := "#666"
			if len(th.JourneySectionColors) > 0 {
				color = th.JourneySectionColors[colorIdx%len(th.JourneySectionColors)]
			}
			b.circle(legendX+5, legendY, 4,
				"fill", color,
				"stroke", color,
			)
			b.text(legendX+15, legendY+4, actor.Name,
				"font-family", th.FontFamily,
				"font-size", "11",
				"fill", th.JourneyTaskText,
			)
			legendX += 80
		}
	}
}
