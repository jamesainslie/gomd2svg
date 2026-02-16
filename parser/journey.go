package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

// maxJourneyScore is the maximum valid score for a journey task.
const maxJourneyScore = 5

var (
	journeyTitleRe   = regexp.MustCompile(`(?i)^\s*title\s+(.+)$`)
	journeySectionRe = regexp.MustCompile(`(?i)^\s*section\s+(.+)$`)
	journeyTaskRe    = regexp.MustCompile(`^\s*(.+?):\s*(\d+)\s*(?::\s*(.*))?$`)
)

func parseJourney(input string) (*ParseOutput, error) { //nolint:unparam // error return is part of the parser interface contract used by Parse().
	graph := ir.NewGraph()
	graph.Kind = ir.Journey

	lines := preprocessInput(input)
	var currentSection string

	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))
		if lower == "journey" {
			continue
		}

		if match := journeyTitleRe.FindStringSubmatch(line); match != nil {
			graph.JourneyTitle = strings.TrimSpace(match[1])
			continue
		}

		if match := journeySectionRe.FindStringSubmatch(line); match != nil {
			currentSection = strings.TrimSpace(match[1])
			graph.JourneySections = append(graph.JourneySections, &ir.JourneySection{
				Name: currentSection,
			})
			continue
		}

		if match := journeyTaskRe.FindStringSubmatch(line); match != nil {
			name := strings.TrimSpace(match[1])
			score, _ := strconv.Atoi(match[2]) //nolint:errcheck // regex guarantees \d+.
			if score < 1 {
				score = 1
			}
			if score > maxJourneyScore {
				score = maxJourneyScore
			}
			var actors []string
			if match[3] != "" {
				for _, actor := range strings.Split(match[3], ",") {
					actor = strings.TrimSpace(actor)
					if actor != "" {
						actors = append(actors, actor)
					}
				}
			}
			taskIdx := len(graph.JourneyTasks)
			graph.JourneyTasks = append(graph.JourneyTasks, &ir.JourneyTask{
				Name:    name,
				Score:   score,
				Actors:  actors,
				Section: currentSection,
			})
			// Add task index to current section.
			if len(graph.JourneySections) > 0 {
				sec := graph.JourneySections[len(graph.JourneySections)-1]
				sec.Tasks = append(sec.Tasks, taskIdx)
			}
			continue
		}
	}

	return &ParseOutput{Graph: graph}, nil
}
