package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	ganttTaskRe = regexp.MustCompile(`^(.+?)\s*:\s*(.+)$`)
	//nolint:gochecknoglobals // package-level lookup table is idiomatic for constant sets.
	ganttTagSet = map[string]bool{"done": true, "active": true, "crit": true, "milestone": true}
	ganttDateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	ganttDurRe  = regexp.MustCompile(`^\d+[dwmhDWMH]$`)
)

//nolint:unparam // error return is part of the parser interface contract used by Parse().
func parseGantt(input string) (*ParseOutput, error) {
	graph := ir.NewGraph()
	graph.Kind = ir.Gantt
	graph.GanttDateFormat = "YYYY-MM-DD" // default

	lines := preprocessInput(input)

	var currentSection *ir.GanttSection

	for _, line := range lines {
		lower := strings.ToLower(line)

		if strings.HasPrefix(lower, "gantt") {
			continue
		}

		// Directives.
		if strings.HasPrefix(lower, "title ") {
			graph.GanttTitle = strings.TrimSpace(line[len("title "):])
			continue
		}
		if strings.HasPrefix(lower, "dateformat ") {
			graph.GanttDateFormat = strings.TrimSpace(line[len("dateformat "):])
			continue
		}
		if strings.HasPrefix(lower, "axisformat ") {
			graph.GanttAxisFormat = strings.TrimSpace(line[len("axisformat "):])
			continue
		}
		if strings.HasPrefix(lower, "excludes ") {
			val := strings.TrimSpace(line[len("excludes "):])
			for _, ex := range strings.Split(val, ",") {
				graph.GanttExcludes = append(graph.GanttExcludes, strings.TrimSpace(ex))
			}
			continue
		}
		if strings.HasPrefix(lower, "tickinterval ") {
			graph.GanttTickInterval = strings.TrimSpace(line[len("tickinterval "):])
			continue
		}
		if strings.HasPrefix(lower, "todaymarker ") {
			graph.GanttTodayMarker = strings.TrimSpace(line[len("todaymarker "):])
			continue
		}
		if strings.HasPrefix(lower, "weekend ") {
			graph.GanttWeekday = strings.TrimSpace(line[len("weekend "):])
			continue
		}

		// Section.
		if strings.HasPrefix(lower, "section ") {
			currentSection = &ir.GanttSection{
				Title: strings.TrimSpace(line[len("section "):]),
			}
			graph.GanttSections = append(graph.GanttSections, currentSection)
			continue
		}

		// Task line: "Task Name : metadata".
		if match := ganttTaskRe.FindStringSubmatch(line); match != nil {
			if currentSection == nil {
				currentSection = &ir.GanttSection{}
				graph.GanttSections = append(graph.GanttSections, currentSection)
			}

			label := strings.TrimSpace(match[1])
			metadata := strings.TrimSpace(match[2])
			task := parseGanttTask(label, metadata)
			currentSection.Tasks = append(currentSection.Tasks, task)
		}
	}

	return &ParseOutput{Graph: graph}, nil
}

func parseGanttTask(label, metadata string) *ir.GanttTask {
	task := &ir.GanttTask{Label: label}

	parts := strings.Split(metadata, ",")
	for idx := range parts {
		parts[idx] = strings.TrimSpace(parts[idx])
	}

	// Extract tags first.
	var remaining []string
	for _, part := range parts {
		if ganttTagSet[strings.ToLower(part)] {
			task.Tags = append(task.Tags, strings.ToLower(part))
		} else {
			remaining = append(remaining, part)
		}
	}

	// Classify remaining items.
	for idx, part := range remaining {
		lp := strings.ToLower(part)

		switch {
		case strings.HasPrefix(lp, "after "):
			ids := strings.Fields(part[len("after "):])
			task.AfterIDs = ids
			task.StartStr = part
		case strings.HasPrefix(lp, "until "):
			task.UntilID = strings.TrimSpace(part[len("until "):])
		case ganttDateRe.MatchString(part):
			if task.StartStr == "" {
				task.StartStr = part
			} else {
				task.EndStr = part
			}
		case ganttDurRe.MatchString(part):
			task.EndStr = part
		default:
			// Must be a task ID -- only if it's the first non-tag item
			// and we haven't set start yet.
			switch {
			case idx == 0 && task.ID == "":
				task.ID = part
			case task.StartStr == "":
				task.StartStr = part
			default:
				task.EndStr = part
			}
		}
	}

	return task
}
