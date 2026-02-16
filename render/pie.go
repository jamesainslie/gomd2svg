package render

import (
	"fmt"
	"math"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderPie(builder *svgBuilder, lay *layout.Layout, th *theme.Theme, _ *config.Layout) {
	pd, ok := lay.Diagram.(layout.PieData)
	if !ok {
		return
	}

	// Title.
	if pd.Title != "" {
		builder.text(lay.Width/2, th.PieTitleTextSize, pd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.PieTitleTextSize),
			"font-weight", "bold",
			"fill", th.PieTitleTextColor,
		)
	}

	cx := pd.CenterX
	cy := pd.CenterY
	radius := pd.Radius

	for _, slice := range pd.Slices {
		color := "#888888" // fallback
		if len(th.PieColors) > 0 {
			color = th.PieColors[slice.ColorIndex%len(th.PieColors)]
		}
		opacity := fmt.Sprintf("%.2f", th.PieOpacity)

		// Full circle special case.
		span := slice.EndAngle - slice.StartAngle
		if span >= 2*math.Pi-0.01 {
			builder.circle(cx, cy, radius,
				"fill", color,
				"fill-opacity", opacity,
				"stroke", th.PieStrokeColor,
				"stroke-width", fmtFloat(th.PieStrokeWidth),
			)
		} else {
			// Arc path.
			x1 := cx + radius*float32(math.Cos(float64(slice.StartAngle)))
			y1 := cy + radius*float32(math.Sin(float64(slice.StartAngle)))
			x2 := cx + radius*float32(math.Cos(float64(slice.EndAngle)))
			y2 := cy + radius*float32(math.Sin(float64(slice.EndAngle)))

			largeArc := "0"
			if span > math.Pi {
				largeArc = "1"
			}

			pathData := fmt.Sprintf("M %s,%s L %s,%s A %s,%s 0 %s,1 %s,%s Z",
				fmtFloat(cx), fmtFloat(cy),
				fmtFloat(x1), fmtFloat(y1),
				fmtFloat(radius), fmtFloat(radius),
				largeArc,
				fmtFloat(x2), fmtFloat(y2),
			)

			builder.path(pathData,
				"fill", color,
				"fill-opacity", opacity,
				"stroke", th.PieStrokeColor,
				"stroke-width", fmtFloat(th.PieStrokeWidth),
			)
		}

		// Slice label.
		labelText := slice.Label
		if pd.ShowData {
			labelText = fmt.Sprintf("%s (%.0f)", slice.Label, slice.Value)
		}
		builder.text(slice.LabelX, slice.LabelY, labelText,
			"text-anchor", "middle",
			"dominant-baseline", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.PieSectionTextSize),
			"fill", th.PieSectionTextColor,
		)
	}

	// Outer stroke.
	if th.PieOuterStrokeWidth > 0 {
		builder.circle(cx, cy, radius,
			"fill", "none",
			"stroke", th.PieOuterStrokeColor,
			"stroke-width", fmtFloat(th.PieOuterStrokeWidth),
		)
	}
}
