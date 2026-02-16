package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	ganttTaskRe = regexp.MustCompile(`^(.+?)\s*:\s*(.+)$`)
	ganttTagSet = map[string]bool{"done": true, "active": true, "crit": true, "milestone": true}
	ganttDateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	ganttDurRe  = regexp.MustCompile(`^\d+[dwmhDWMH]$`)
)

func parseGantt(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Gantt
	g.GanttDateFormat = "YYYY-MM-DD" // default

	lines := preprocessInput(input)

	var currentSection *ir.GanttSection

	for _, line := range lines {
		lower := strings.ToLower(line)

		if strings.HasPrefix(lower, "gantt") {
			continue
		}

		// Directives.
		if strings.HasPrefix(lower, "title ") {
			g.GanttTitle = strings.TrimSpace(line[len("title "):])
			continue
		}
		if strings.HasPrefix(lower, "dateformat ") {
			g.GanttDateFormat = strings.TrimSpace(line[len("dateformat "):])
			continue
		}
		if strings.HasPrefix(lower, "axisformat ") {
			g.GanttAxisFormat = strings.TrimSpace(line[len("axisformat "):])
			continue
		}
		if strings.HasPrefix(lower, "excludes ") {
			val := strings.TrimSpace(line[len("excludes "):])
			for _, ex := range strings.Split(val, ",") {
				g.GanttExcludes = append(g.GanttExcludes, strings.TrimSpace(ex))
			}
			continue
		}
		if strings.HasPrefix(lower, "tickinterval ") {
			g.GanttTickInterval = strings.TrimSpace(line[len("tickinterval "):])
			continue
		}
		if strings.HasPrefix(lower, "todaymarker ") {
			g.GanttTodayMarker = strings.TrimSpace(line[len("todaymarker "):])
			continue
		}
		if strings.HasPrefix(lower, "weekend ") {
			g.GanttWeekday = strings.TrimSpace(line[len("weekend "):])
			continue
		}

		// Section.
		if strings.HasPrefix(lower, "section ") {
			currentSection = &ir.GanttSection{
				Title: strings.TrimSpace(line[len("section "):]),
			}
			g.GanttSections = append(g.GanttSections, currentSection)
			continue
		}

		// Task line: "Task Name : metadata"
		if m := ganttTaskRe.FindStringSubmatch(line); m != nil {
			if currentSection == nil {
				currentSection = &ir.GanttSection{}
				g.GanttSections = append(g.GanttSections, currentSection)
			}

			label := strings.TrimSpace(m[1])
			metadata := strings.TrimSpace(m[2])
			task := parseGanttTask(label, metadata)
			currentSection.Tasks = append(currentSection.Tasks, task)
		}
	}

	return &ParseOutput{Graph: g}, nil
}

func parseGanttTask(label, metadata string) *ir.GanttTask {
	task := &ir.GanttTask{Label: label}

	parts := strings.Split(metadata, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	// Extract tags first.
	var remaining []string
	for _, p := range parts {
		if ganttTagSet[strings.ToLower(p)] {
			task.Tags = append(task.Tags, strings.ToLower(p))
		} else {
			remaining = append(remaining, p)
		}
	}

	// Classify remaining items.
	for i, p := range remaining {
		lp := strings.ToLower(p)

		if strings.HasPrefix(lp, "after ") {
			ids := strings.Fields(p[len("after "):])
			task.AfterIDs = ids
			task.StartStr = p
		} else if strings.HasPrefix(lp, "until ") {
			task.UntilID = strings.TrimSpace(p[len("until "):])
		} else if ganttDateRe.MatchString(p) {
			if task.StartStr == "" {
				task.StartStr = p
			} else {
				task.EndStr = p
			}
		} else if ganttDurRe.MatchString(p) {
			task.EndStr = p
		} else {
			// Must be a task ID — only if it's the first non-tag item
			// and we haven't set start yet.
			if i == 0 && task.ID == "" {
				task.ID = p
			} else if task.StartStr == "" {
				task.StartStr = p
			} else {
				task.EndStr = p
			}
		}
	}

	return task
}
