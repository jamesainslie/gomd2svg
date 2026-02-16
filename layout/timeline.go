package layout

import (
	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeTimelineLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	padX := cfg.Timeline.PaddingX
	padY := cfg.Timeline.PaddingY
	periodW := cfg.Timeline.PeriodWidth
	eventH := cfg.Timeline.EventHeight
	secPad := cfg.Timeline.SectionPadding

	// Title height.
	var titleHeight float32
	if g.TimelineTitle != "" {
		titleHeight = th.FontSize + padY
	}

	// Count max periods across all sections for width.
	var maxPeriods int
	for _, sec := range g.TimelineSections {
		if len(sec.Periods) > maxPeriods {
			maxPeriods = len(sec.Periods)
		}
	}

	// Section label width.
	var sectionLabelWidth float32
	for _, sec := range g.TimelineSections {
		if sec.Title != "" {
			sectionLabelWidth = padX * 3 // fixed width for labels
			break
		}
	}

	// Compute layout per section.
	var sections []TimelineSectionLayout
	curY := titleHeight + padY

	for i, sec := range g.TimelineSections {
		// Find max events in any period of this section.
		var maxEvents int
		for _, p := range sec.Periods {
			if len(p.Events) > maxEvents {
				maxEvents = len(p.Events)
			}
		}
		if maxEvents == 0 {
			maxEvents = 1
		}

		sectionH := float32(maxEvents)*eventH + secPad*2

		// Color cycling.
		color := "#F0F4F8" // fallback
		if len(th.TimelineSectionColors) > 0 {
			color = th.TimelineSectionColors[i%len(th.TimelineSectionColors)]
		}

		var periods []TimelinePeriodLayout
		for j, p := range sec.Periods {
			px := padX + sectionLabelWidth + float32(j)*periodW

			var events []TimelineEventLayout
			for k, e := range p.Events {
				events = append(events, TimelineEventLayout{
					Text:   e.Text,
					X:      px + secPad,
					Y:      curY + secPad + float32(k)*eventH,
					Width:  periodW - secPad*2,
					Height: eventH,
				})
			}

			periods = append(periods, TimelinePeriodLayout{
				Title:  p.Title,
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
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: TimelineData{
			Sections: sections,
			Title:    g.TimelineTitle,
		},
	}
}
