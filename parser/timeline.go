package parser

import (
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

func parseTimeline(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Timeline

	lines := preprocessInput(input)

	var currentSection *ir.TimelineSection
	var currentPeriod *ir.TimelinePeriod

	for _, line := range lines {
		lower := strings.ToLower(line)

		if strings.HasPrefix(lower, "timeline") {
			continue
		}

		if strings.HasPrefix(lower, "title ") {
			g.TimelineTitle = strings.TrimSpace(line[len("title "):])
			continue
		}

		if strings.HasPrefix(lower, "section ") {
			currentSection = &ir.TimelineSection{
				Title: strings.TrimSpace(line[len("section "):]),
			}
			g.TimelineSections = append(g.TimelineSections, currentSection)
			currentPeriod = nil
			continue
		}

		// Ensure we have a section.
		if currentSection == nil {
			currentSection = &ir.TimelineSection{}
			g.TimelineSections = append(g.TimelineSections, currentSection)
		}

		// Continuation event line: starts with ":"
		if strings.HasPrefix(strings.TrimSpace(line), ":") {
			if currentPeriod != nil {
				eventText := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), ":"))
				if eventText != "" {
					currentPeriod.Events = append(currentPeriod.Events, &ir.TimelineEvent{Text: eventText})
				}
			}
			continue
		}

		// Period line: "period : event : event ..."
		if idx := strings.Index(line, ":"); idx >= 0 {
			period := strings.TrimSpace(line[:idx])
			rest := line[idx+1:]

			currentPeriod = &ir.TimelinePeriod{Title: period}
			currentSection.Periods = append(currentSection.Periods, currentPeriod)

			// Split remaining by ":" for multiple events.
			parts := strings.Split(rest, ":")
			for _, p := range parts {
				text := strings.TrimSpace(p)
				if text != "" {
					currentPeriod.Events = append(currentPeriod.Events, &ir.TimelineEvent{Text: text})
				}
			}
		}
	}

	return &ParseOutput{Graph: g}, nil
}
