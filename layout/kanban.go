package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeKanbanLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	kc := cfg.Kanban
	pad := kc.Padding

	columns := make([]KanbanColumnLayout, len(g.Columns))
	cursorX := pad

	maxColHeight := float32(0)

	for i, col := range g.Columns {
		// Measure column header
		headerTW := measurer.Width(col.Label, th.FontSize, th.FontFamily)
		headerTB := TextBlock{
			Lines:    []string{col.Label},
			Width:    headerTW,
			Height:   th.FontSize * cfg.LabelLineHeight,
			FontSize: th.FontSize,
		}

		// Layout cards vertically within column
		cardY := kc.HeaderHeight + pad
		cards := make([]KanbanCardLayout, len(col.Cards))

		for j, card := range col.Cards {
			cardTW := measurer.Width(card.Label, th.FontSize, th.FontFamily)
			cardH := th.FontSize*cfg.LabelLineHeight + 2*pad
			cardTB := TextBlock{
				Lines:    []string{card.Label},
				Width:    cardTW,
				Height:   th.FontSize * cfg.LabelLineHeight,
				FontSize: th.FontSize,
			}

			// Build metadata map for renderer
			meta := make(map[string]string)
			if card.Assigned != "" {
				meta["assigned"] = card.Assigned
			}
			if card.Ticket != "" {
				meta["ticket"] = card.Ticket
			}
			if card.Icon != "" {
				meta["icon"] = card.Icon
			}

			cards[j] = KanbanCardLayout{
				ID:       card.ID,
				Label:    cardTB,
				Priority: card.Priority,
				X:        cursorX + pad,
				Y:        cardY,
				Width:    kc.SectionWidth - 2*pad,
				Height:   cardH,
				Metadata: meta,
			}

			cardY += cardH + kc.CardSpacing
		}

		colHeight := cardY + pad
		if colHeight < kc.HeaderHeight+2*pad {
			colHeight = kc.HeaderHeight + 2*pad
		}

		columns[i] = KanbanColumnLayout{
			ID:     col.ID,
			Label:  headerTB,
			X:      cursorX,
			Y:      0,
			Width:  kc.SectionWidth,
			Height: colHeight,
			Cards:  cards,
		}

		if colHeight > maxColHeight {
			maxColHeight = colHeight
		}

		cursorX += kc.SectionWidth + pad
	}

	// Normalize all columns to same height
	for i := range columns {
		columns[i].Height = maxColHeight
	}

	totalW := cursorX
	totalH := maxColHeight + pad

	return &Layout{
		Kind:   g.Kind,
		Width:  totalW,
		Height: totalH,
		Diagram: KanbanData{
			Columns: columns,
		},
	}
}
