package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/textmetrics"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computePacketLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	measurer := textmetrics.New()
	pc := cfg.Packet
	bitsPerRow := pc.BitsPerRow
	bitW := pc.BitWidth
	rowH := pc.RowHeight
	padX := pc.PaddingX
	padY := pc.PaddingY

	totalRowW := float32(bitsPerRow) * bitW

	// Group fields into rows based on bit positions.
	type rowBuilder struct {
		fields []PacketFieldLayout
		y      float32
	}
	var rows []rowBuilder

	cursorY := padY
	// If ShowBits, reserve space for bit number labels at top.
	if pc.ShowBits {
		cursorY += th.FontSize*cfg.LabelLineHeight + padY
	}

	for _, field := range g.Fields {
		startRow := field.Start / bitsPerRow
		endRow := field.End / bitsPerRow

		for row := startRow; row <= endRow; row++ {
			// Ensure we have enough rows
			for len(rows) <= row {
				rows = append(rows, rowBuilder{
					y: cursorY + float32(len(rows))*(rowH+padY),
				})
			}

			// Compute the bit range within this row
			rowStartBit := row * bitsPerRow
			fieldStartInRow := field.Start
			if fieldStartInRow < rowStartBit {
				fieldStartInRow = rowStartBit
			}
			fieldEndInRow := field.End
			if fieldEndInRow >= rowStartBit+bitsPerRow {
				fieldEndInRow = rowStartBit + bitsPerRow - 1
			}

			bitsInRow := fieldEndInRow - fieldStartInRow + 1
			offsetInRow := fieldStartInRow - rowStartBit

			x := padX + float32(offsetInRow)*bitW
			w := float32(bitsInRow) * bitW

			tw := measurer.Width(field.Description, th.FontSize, th.FontFamily)
			tb := TextBlock{
				Lines:    []string{field.Description},
				Width:    tw,
				Height:   th.FontSize * cfg.LabelLineHeight,
				FontSize: th.FontSize,
			}

			rows[row].fields = append(rows[row].fields, PacketFieldLayout{
				Label:    tb,
				X:        x,
				Y:        rows[row].y,
				Width:    w,
				Height:   rowH,
				StartBit: fieldStartInRow,
				EndBit:   fieldEndInRow,
			})
		}
	}

	// Build final row layouts
	resultRows := make([]PacketRowLayout, len(rows))
	for i, rb := range rows {
		resultRows[i] = PacketRowLayout{
			Y:      rb.y,
			Height: rowH,
			Fields: rb.fields,
		}
	}

	totalH := cursorY + float32(len(rows))*(rowH+padY) + padY
	totalW := totalRowW + 2*padX

	return &Layout{
		Kind:   g.Kind,
		Width:  totalW,
		Height: totalH,
		Diagram: PacketData{
			Rows:       resultRows,
			BitsPerRow: bitsPerRow,
			ShowBits:   pc.ShowBits,
		},
	}
}
