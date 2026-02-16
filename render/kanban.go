package render

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Kanban diagram rendering constants.
const (
	kanbanColumnRadius float32 = 4
)

func renderKanban(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	kd, ok := lay.Diagram.(layout.KanbanData)
	if !ok {
		return
	}

	for _, col := range kd.Columns {
		// Column background
		builder.rect(col.X, col.Y, col.Width, col.Height, kanbanColumnRadius,
			"fill", th.ClusterBackground,
			"stroke", th.NodeBorderColor,
			"stroke-width", "1",
		)

		// Column header text (centered, bold)
		headerX := col.X + col.Width/2
		headerY := col.Y + cfg.Kanban.HeaderHeight/2 + col.Label.FontSize/3
		builder.text(headerX, headerY, col.Label.Lines[0],
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(col.Label.FontSize),
			"font-weight", "bold",
			"fill", th.PrimaryTextColor,
		)

		// Header divider line
		divY := col.Y + cfg.Kanban.HeaderHeight
		builder.line(col.X, divY, col.X+col.Width, divY,
			"stroke", th.NodeBorderColor,
			"stroke-width", "1",
		)

		// Cards
		for _, card := range col.Cards {
			builder.rect(card.X, card.Y, card.Width, card.Height, 3,
				"fill", th.Background,
				"stroke", th.NodeBorderColor,
				"stroke-width", "1",
			)

			// Card label
			textX := card.X + cfg.Kanban.Padding
			textY := card.Y + card.Height/2 + card.Label.FontSize/3
			builder.text(textX, textY, card.Label.Lines[0],
				"font-family", th.FontFamily,
				"font-size", fmtFloat(card.Label.FontSize),
				"fill", th.PrimaryTextColor,
			)
		}
	}
}
