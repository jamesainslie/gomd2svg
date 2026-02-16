package render

import (
	"fmt"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func renderGantt(b *svgBuilder, l *layout.Layout, th *theme.Theme, _ *config.Layout) {
	gd, ok := l.Diagram.(layout.GanttData)
	if !ok {
		return
	}

	// Title.
	if gd.Title != "" {
		b.text(l.Width/2, th.FontSize+5, gd.Title,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", fmtFloat(th.FontSize+2),
			"font-weight", "bold",
			"fill", th.TextColor,
		)
	}

	// Section backgrounds.
	for _, sec := range gd.Sections {
		b.rect(gd.ChartX, sec.Y, gd.ChartWidth, sec.Height, 0,
			"fill", sec.Color,
			"stroke", "none",
		)
		// Section label.
		if sec.Title != "" {
			b.text(gd.ChartX-5, sec.Y+sec.Height/2, sec.Title,
				"text-anchor", "end",
				"dominant-baseline", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize-2),
				"fill", th.TextColor,
			)
		}
	}

	// Grid lines.
	for _, tick := range gd.AxisTicks {
		b.line(tick.X, gd.ChartY, tick.X, gd.ChartY+gd.ChartHeight,
			"stroke", th.GanttGridColor,
			"stroke-width", "0.5",
		)
		// Axis label.
		b.text(tick.X, gd.ChartY+gd.ChartHeight+12, tick.Label,
			"text-anchor", "middle",
			"font-family", th.FontFamily,
			"font-size", "9",
			"fill", th.TextColor,
		)
	}

	// Task bars.
	for _, sec := range gd.Sections {
		for _, task := range sec.Tasks {
			fill := th.GanttTaskFill
			border := th.GanttTaskBorder
			if task.IsCrit {
				fill = th.GanttCritFill
				border = th.GanttCritBorder
			}
			if task.IsDone {
				fill = th.GanttDoneFill
			}
			if task.IsActive {
				fill = th.GanttActiveFill
			}

			if task.IsMilestone {
				// Render as diamond.
				cx := task.X
				cy := task.Y + task.Height/2
				size := task.Height / 2
				d := fmt.Sprintf("M %s,%s L %s,%s L %s,%s L %s,%s Z",
					fmtFloat(cx), fmtFloat(cy-size),
					fmtFloat(cx+size), fmtFloat(cy),
					fmtFloat(cx), fmtFloat(cy+size),
					fmtFloat(cx-size), fmtFloat(cy),
				)
				b.path(d,
					"fill", th.GanttMilestoneFill,
					"stroke", border,
					"stroke-width", "1",
				)
			} else {
				b.rect(task.X, task.Y, task.Width, task.Height, 2,
					"fill", fill,
					"stroke", border,
					"stroke-width", "1",
				)
			}

			// Task label.
			b.text(task.X+task.Width+4, task.Y+task.Height/2+1, task.Label,
				"dominant-baseline", "middle",
				"font-family", th.FontFamily,
				"font-size", fmtFloat(th.FontSize-3),
				"fill", th.TextColor,
			)
		}
	}

	// Today marker.
	if gd.ShowTodayMarker {
		b.line(gd.TodayMarkerX, gd.ChartY, gd.TodayMarkerX, gd.ChartY+gd.ChartHeight,
			"stroke", th.GanttTodayMarkerColor,
			"stroke-width", "2",
			"stroke-dasharray", "4,4",
		)
	}
}
