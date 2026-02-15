package render

import (
	"fmt"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func renderPacket(b *svgBuilder, l *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	pd, ok := l.Diagram.(layout.PacketData)
	if !ok {
		return
	}

	pc := cfg.Packet
	smallFontSize := th.FontSize * 0.7

	// Render bit numbers at top if enabled
	if pd.ShowBits {
		for bit := 0; bit < pd.BitsPerRow; bit++ {
			x := pc.PaddingX + float32(bit)*pc.BitWidth + pc.BitWidth/2
			y := th.FontSize * cfg.LabelLineHeight
			b.text(x, y, fmt.Sprintf("%d", bit),
				"text-anchor", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(smallFontSize),
				"fill", th.SecondaryTextColor,
			)
		}
	}

	// Render rows and fields
	for _, row := range pd.Rows {
		for _, field := range row.Fields {
			// Field rectangle
			b.rect(field.X, field.Y, field.Width, field.Height, 0,
				"fill", th.PrimaryColor,
				"stroke", th.NodeBorderColor,
				"stroke-width", "1",
			)

			// Field label (centered)
			textX := field.X + field.Width/2
			textY := field.Y + field.Height/2 + field.Label.FontSize/3
			b.text(textX, textY, field.Label.Lines[0],
				"text-anchor", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(field.Label.FontSize),
				"fill", th.PrimaryTextColor,
			)

			// Bit range labels at bottom-left and bottom-right of field
			bitY := field.Y + field.Height - 2
			b.text(field.X+2, bitY, fmt.Sprintf("%d", field.StartBit),
				"font-family", th.FontFamily,
				"font-size", fmtFloat(smallFontSize),
				"fill", th.SecondaryTextColor,
			)
			if field.EndBit != field.StartBit {
				b.text(field.X+field.Width-2, bitY, fmt.Sprintf("%d", field.EndBit),
					"text-anchor", "end",
					"font-family", th.FontFamily,
					"font-size", fmtFloat(smallFontSize),
					"fill", th.SecondaryTextColor,
				)
			}
		}
	}
}
