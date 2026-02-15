package render

import (
	"fmt"
	"math"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func renderPie(b *svgBuilder, l *layout.Layout, th *theme.Theme, _ *config.Layout) {
	pd, ok := l.Diagram.(layout.PieData)
	if !ok {
		return
	}

	// Title.
	if pd.Title != "" {
		b.text(l.Width/2, th.PieTitleTextSize, pd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.PieTitleTextSize),
			"font-weight", "bold",
			"fill", th.PieTitleTextColor,
		)
	}

	cx := pd.CenterX
	cy := pd.CenterY
	r := pd.Radius

	for _, s := range pd.Slices {
		color := "#888888" // fallback
		if len(th.PieColors) > 0 {
			color = th.PieColors[s.ColorIndex%len(th.PieColors)]
		}
		opacity := fmt.Sprintf("%.2f", th.PieOpacity)

		// Full circle special case.
		span := s.EndAngle - s.StartAngle
		if span >= 2*math.Pi-0.01 {
			b.circle(cx, cy, r,
				"fill", color,
				"fill-opacity", opacity,
				"stroke", th.PieStrokeColor,
				"stroke-width", fmtFloat(th.PieStrokeWidth),
			)
		} else {
			// Arc path.
			x1 := cx + r*float32(math.Cos(float64(s.StartAngle)))
			y1 := cy + r*float32(math.Sin(float64(s.StartAngle)))
			x2 := cx + r*float32(math.Cos(float64(s.EndAngle)))
			y2 := cy + r*float32(math.Sin(float64(s.EndAngle)))

			largeArc := "0"
			if span > math.Pi {
				largeArc = "1"
			}

			d := fmt.Sprintf("M %s,%s L %s,%s A %s,%s 0 %s,1 %s,%s Z",
				fmtFloat(cx), fmtFloat(cy),
				fmtFloat(x1), fmtFloat(y1),
				fmtFloat(r), fmtFloat(r),
				largeArc,
				fmtFloat(x2), fmtFloat(y2),
			)

			b.path(d,
				"fill", color,
				"fill-opacity", opacity,
				"stroke", th.PieStrokeColor,
				"stroke-width", fmtFloat(th.PieStrokeWidth),
			)
		}

		// Slice label.
		labelText := s.Label
		if pd.ShowData {
			labelText = fmt.Sprintf("%s (%.0f)", s.Label, s.Value)
		}
		b.text(s.LabelX, s.LabelY, labelText,
			"text-anchor", "middle",
			"dominant-baseline", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.PieSectionTextSize),
			"fill", th.PieSectionTextColor,
		)
	}

	// Outer stroke.
	if th.PieOuterStrokeWidth > 0 {
		b.circle(cx, cy, r,
			"fill", "none",
			"stroke", th.PieOuterStrokeColor,
			"stroke-width", fmtFloat(th.PieOuterStrokeWidth),
		)
	}
}
