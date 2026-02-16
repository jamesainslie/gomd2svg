package render

import (
	"strconv"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

// Packet diagram rendering constants.
const (
	packetSmallFontScale float32 = 0.7
)

func renderPacket(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, cfg *config.Layout) {
	pd, ok := lay.Diagram.(layout.PacketData)
	if !ok {
		return
	}

	pc := cfg.Packet
	smallFontSize := th.FontSize * packetSmallFontScale

	// Render bit numbers at top if enabled
	if pd.ShowBits {
		for bit := range pd.BitsPerRow {
			posX := pc.PaddingX + float32(bit)*pc.BitWidth + pc.BitWidth/2
			posY := th.FontSize * cfg.LabelLineHeight
			builder.text(posX, posY, strconv.Itoa(bit),
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
			builder.rect(field.X, field.Y, field.Width, field.Height, 0,
				"fill", th.PrimaryColor,
				"stroke", th.NodeBorderColor,
				"stroke-width", "1",
			)

			// Field label (centered)
			textX := field.X + field.Width/2
			textY := field.Y + field.Height/2 + field.Label.FontSize/3
			builder.text(textX, textY, field.Label.Lines[0],
				"text-anchor", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(field.Label.FontSize),
				"fill", th.PrimaryTextColor,
			)

			// Bit range labels at bottom-left and bottom-right of field
			bitY := field.Y + field.Height - 2
			builder.text(field.X+2, bitY, strconv.Itoa(field.StartBit),
				"font-family", th.FontFamily,
				"font-size", fmtFloat(smallFontSize),
				"fill", th.SecondaryTextColor,
			)
			if field.EndBit != field.StartBit {
				builder.text(field.X+field.Width-2, bitY, strconv.Itoa(field.EndBit),
					"text-anchor", "end",
					"font-family", th.FontFamily,
					"font-size", fmtFloat(smallFontSize),
					"fill", th.SecondaryTextColor,
				)
			}
		}
	}
}
