package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	journeyTitleRe   = regexp.MustCompile(`(?i)^\s*title\s+(.+)$`)
	journeySectionRe = regexp.MustCompile(`(?i)^\s*section\s+(.+)$`)
	journeyTaskRe    = regexp.MustCompile(`^\s*(.+?):\s*(\d+)\s*(?::\s*(.*))?$`)
)

func parseJourney(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Journey

	lines := preprocessInput(input)
	var currentSection string

	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))
		if lower == "journey" {
			continue
		}

		if m := journeyTitleRe.FindStringSubmatch(line); m != nil {
			g.JourneyTitle = strings.TrimSpace(m[1])
			continue
		}

		if m := journeySectionRe.FindStringSubmatch(line); m != nil {
			currentSection = strings.TrimSpace(m[1])
			g.JourneySections = append(g.JourneySections, &ir.JourneySection{
				Name: currentSection,
			})
			continue
		}

		if m := journeyTaskRe.FindStringSubmatch(line); m != nil {
			name := strings.TrimSpace(m[1])
			score, _ := strconv.Atoi(m[2]) // regex guarantees \d+
			if score < 1 {
				score = 1
			}
			if score > 5 {
				score = 5
			}
			var actors []string
			if m[3] != "" {
				for _, a := range strings.Split(m[3], ",") {
					a = strings.TrimSpace(a)
					if a != "" {
						actors = append(actors, a)
					}
				}
			}
			taskIdx := len(g.JourneyTasks)
			g.JourneyTasks = append(g.JourneyTasks, &ir.JourneyTask{
				Name:    name,
				Score:   score,
				Actors:  actors,
				Section: currentSection,
			})
			// Add task index to current section.
			if len(g.JourneySections) > 0 {
				sec := g.JourneySections[len(g.JourneySections)-1]
				sec.Tasks = append(sec.Tasks, taskIdx)
			}
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}
