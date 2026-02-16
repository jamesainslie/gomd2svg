package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeTimelineLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	padX := cfg.Timeline.PaddingX
	padY := cfg.Timeline.PaddingY
	periodW := cfg.Timeline.PeriodWidth
	eventH := cfg.Timeline.EventHeight
	secPad := cfg.Timeline.SectionPadding

	// Title height.
	var titleHeight float32
	if graph.TimelineTitle != "" {
		titleHeight = th.FontSize + padY
	}

	// Count max periods across all sections for width.
	var maxPeriods int
	for _, sec := range graph.TimelineSections {
		if len(sec.Periods) > maxPeriods {
			maxPeriods = len(sec.Periods)
		}
	}

	// Section label width.
	var sectionLabelWidth float32
	for _, sec := range graph.TimelineSections {
		if sec.Title != "" {
			sectionLabelWidth = padX * 3 // fixed width for labels
			break
		}
	}

	// Compute layout per section.
	sections := make([]TimelineSectionLayout, 0, len(graph.TimelineSections))
	curY := titleHeight + padY

	for secIdx, sec := range graph.TimelineSections {
		// Find max events in any period of this section.
		var maxEvents int
		for _, period := range sec.Periods {
			if len(period.Events) > maxEvents {
				maxEvents = len(period.Events)
			}
		}
		if maxEvents == 0 {
			maxEvents = 1
		}

		sectionH := float32(maxEvents)*eventH + secPad*2

		// Color cycling.
		color := "#F0F4F8" // fallback
		if len(th.TimelineSectionColors) > 0 {
			color = th.TimelineSectionColors[secIdx%len(th.TimelineSectionColors)]
		}

		periods := make([]TimelinePeriodLayout, 0, len(sec.Periods))
		for periodIdx, period := range sec.Periods {
			px := padX + sectionLabelWidth + float32(periodIdx)*periodW

			events := make([]TimelineEventLayout, 0, len(period.Events))
			for evIdx, ev := range period.Events {
				events = append(events, TimelineEventLayout{
					Text:   ev.Text,
					X:      px + secPad,
					Y:      curY + secPad + float32(evIdx)*eventH,
					Width:  periodW - secPad*2,
					Height: eventH,
				})
			}

			periods = append(periods, TimelinePeriodLayout{
				Title:  period.Title,
				X:      px,
				Y:      curY,
				Width:  periodW,
				Height: sectionH,
				Events: events,
			})
		}

		sections = append(sections, TimelineSectionLayout{
			Title:   sec.Title,
			X:       padX,
			Y:       curY,
			Width:   sectionLabelWidth + float32(len(sec.Periods))*periodW,
			Height:  sectionH,
			Color:   color,
			Periods: periods,
		})

		curY += sectionH
	}

	totalW := padX*2 + sectionLabelWidth + float32(maxPeriods)*periodW
	totalH := curY + padY

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: TimelineData{
			Sections: sections,
			Title:    graph.TimelineTitle,
		},
	}
}
