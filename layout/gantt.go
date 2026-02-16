package layout

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

var ganttDurationRe = regexp.MustCompile(`^(\d+)([dwmhDWMH])$`)

// mermaidDateToGoLayout converts mermaid dateFormat tokens to Go time layout.
func mermaidDateToGoLayout(format string) string {
	r := strings.NewReplacer(
		"YYYY", "2006", "YY", "06",
		"MM", "01", "DD", "02",
		"HH", "15", "mm", "04", "ss", "05",
	)
	return r.Replace(format)
}

// parseMermaidDuration converts a mermaid duration string to time.Duration.
func parseMermaidDuration(s string) time.Duration {
	m := ganttDurationRe.FindStringSubmatch(s)
	if m == nil {
		return 0
	}
	n, _ := strconv.Atoi(m[1]) // regex guarantees digits
	switch strings.ToLower(m[2]) {
	case "d":
		return time.Duration(n) * 24 * time.Hour
	case "w":
		return time.Duration(n) * 7 * 24 * time.Hour
	case "h":
		return time.Duration(n) * time.Hour
	case "m":
		return time.Duration(n) * time.Minute
	default:
		return 0
	}
}

// isExcluded checks if a date should be excluded based on the excludes list.
func isExcluded(t time.Time, excludes []string, goLayout string) bool {
	dayName := strings.ToLower(t.Weekday().String())
	for _, ex := range excludes {
		ex = strings.ToLower(strings.TrimSpace(ex))
		if ex == "weekends" && (t.Weekday() == time.Saturday || t.Weekday() == time.Sunday) {
			return true
		}
		if ex == dayName {
			return true
		}
		// Try parsing as a date.
		if exDate, err := time.Parse(goLayout, ex); err == nil {
			if t.Year() == exDate.Year() && t.YearDay() == exDate.YearDay() {
				return true
			}
		}
	}
	return false
}

// addWorkingDays adds n working days to start, skipping excluded days.
func addWorkingDays(start time.Time, days int, excludes []string, goLayout string) time.Time {
	if len(excludes) == 0 {
		return start.Add(time.Duration(days) * 24 * time.Hour)
	}
	t := start
	added := 0
	for added < days {
		t = t.Add(24 * time.Hour)
		if !isExcluded(t, excludes, goLayout) {
			added++
		}
	}
	return t
}

type resolvedTask struct {
	Start time.Time
	End   time.Time
}

func computeGanttLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	goLayout := mermaidDateToGoLayout(g.GanttDateFormat)
	sidePad := cfg.Gantt.SidePadding
	topPad := cfg.Gantt.TopPadding
	barH := cfg.Gantt.BarHeight
	barGap := cfg.Gantt.BarGap

	// Title height.
	var titleHeight float32
	if g.GanttTitle != "" {
		titleHeight = th.FontSize + 10
	}

	// Collect all tasks in order for resolution.
	var allTasks []*ir.GanttTask
	for _, sec := range g.GanttSections {
		allTasks = append(allTasks, sec.Tasks...)
	}

	// Resolve all task dates.
	resolved := make(map[string]resolvedTask)
	var prevEnd time.Time

	for _, task := range allTasks {
		start, end := resolveGanttTaskDates(task, resolved, &prevEnd, goLayout, g.GanttExcludes)
		if task.ID != "" {
			resolved[task.ID] = resolvedTask{Start: start, End: end}
		}
		prevEnd = end
	}

	// Find global date range.
	minDate, maxDate := ganttDateRange(allTasks, resolved, goLayout, g.GanttExcludes)

	totalDays := maxDate.Sub(minDate).Hours() / 24
	if totalDays < 1 {
		totalDays = 1
	}

	chartW := float32(totalDays) * 20 // 20px per day
	if chartW < 200 {
		chartW = 200
	}
	if chartW > 2000 {
		chartW = 2000
	}

	chartX := sidePad
	chartY := titleHeight + topPad

	// dateToX converts a date to an X pixel position.
	dateToX := func(t time.Time) float32 {
		days := t.Sub(minDate).Hours() / 24
		return chartX + float32(days/totalDays)*chartW
	}

	// Build sections and tasks.
	var sections []GanttSectionLayout
	curY := chartY
	prevEnd = time.Time{}

	for i, sec := range g.GanttSections {
		var tasks []GanttTaskLayout
		secStartY := curY

		for _, task := range sec.Tasks {
			start, end := resolveGanttTaskDates(task, resolved, &prevEnd, goLayout, g.GanttExcludes)

			x := dateToX(start)
			w := dateToX(end) - x
			if w < 1 {
				w = 1
			}

			tasks = append(tasks, GanttTaskLayout{
				ID:          task.ID,
				Label:       task.Label,
				X:           x,
				Y:           curY,
				Width:       w,
				Height:      barH,
				IsCrit:      hasTag(task.Tags, "crit"),
				IsDone:      hasTag(task.Tags, "done"),
				IsActive:    hasTag(task.Tags, "active"),
				IsMilestone: hasTag(task.Tags, "milestone"),
			})

			prevEnd = end
			curY += barH + barGap
		}

		secH := curY - secStartY
		color := "#F0F4F8" // fallback
		if len(th.GanttSectionColors) > 0 {
			color = th.GanttSectionColors[i%len(th.GanttSectionColors)]
		}
		sections = append(sections, GanttSectionLayout{
			Title:  sec.Title,
			Y:      secStartY,
			Height: secH,
			Color:  color,
			Tasks:  tasks,
		})
	}

	// Axis ticks.
	var axisTicks []GanttAxisTick
	tickDays := 7
	if totalDays < 14 {
		tickDays = 1
	} else if totalDays > 90 {
		tickDays = 30
	}
	for d := minDate; !d.After(maxDate); d = d.AddDate(0, 0, tickDays) {
		axisTicks = append(axisTicks, GanttAxisTick{
			Label: d.Format("2006-01-02"),
			X:     dateToX(d),
		})
	}

	// Today marker.
	today := time.Now()
	showToday := g.GanttTodayMarker != "off" && !today.Before(minDate) && !today.After(maxDate)
	var todayX float32
	if showToday {
		todayX = dateToX(today)
	}

	totalW := sidePad*2 + chartW
	totalH := curY + topPad

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: GanttData{
			Sections:        sections,
			Title:           g.GanttTitle,
			AxisTicks:       axisTicks,
			TodayMarkerX:    todayX,
			ShowTodayMarker: showToday,
			ChartX:          chartX,
			ChartY:          chartY,
			ChartWidth:      chartW,
			ChartHeight:     curY - chartY,
		},
	}
}

// resolveGanttTaskDates computes start and end times for a single task.
func resolveGanttTaskDates(task *ir.GanttTask, resolved map[string]resolvedTask, prevEnd *time.Time, goLayout string, excludes []string) (time.Time, time.Time) {
	var start, end time.Time

	// Check if already resolved by ID.
	if task.ID != "" {
		if rt, ok := resolved[task.ID]; ok {
			return rt.Start, rt.End
		}
	}

	// Resolve start.
	if len(task.AfterIDs) > 0 {
		for _, depID := range task.AfterIDs {
			if dep, ok := resolved[depID]; ok {
				if dep.End.After(start) {
					start = dep.End
				}
			}
		}
	} else if task.StartStr != "" && !strings.HasPrefix(strings.ToLower(task.StartStr), "after ") {
		if t, err := time.Parse(goLayout, task.StartStr); err == nil {
			start = t
		}
	}

	if start.IsZero() && prevEnd != nil && !prevEnd.IsZero() {
		start = *prevEnd
	}
	if start.IsZero() {
		start = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	// Resolve end.
	dur := parseMermaidDuration(task.EndStr)
	if dur > 0 {
		days := int(dur.Hours() / 24)
		if days > 0 {
			end = addWorkingDays(start, days, excludes, goLayout)
		} else {
			end = start.Add(dur)
		}
	} else if t, err := time.Parse(goLayout, task.EndStr); err == nil {
		end = t
	} else {
		end = start.Add(24 * time.Hour)
	}

	return start, end
}

// ganttDateRange finds the global min and max dates across all tasks.
func ganttDateRange(allTasks []*ir.GanttTask, resolved map[string]resolvedTask, goLayout string, excludes []string) (time.Time, time.Time) {
	var minDate, maxDate time.Time
	first := true
	var prevEnd time.Time

	for _, task := range allTasks {
		start, end := resolveGanttTaskDates(task, resolved, &prevEnd, goLayout, excludes)
		if first || start.Before(minDate) {
			minDate = start
		}
		if first || end.After(maxDate) {
			maxDate = end
		}
		first = false
		prevEnd = end
	}

	if minDate.IsZero() {
		minDate = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if maxDate.IsZero() || !maxDate.After(minDate) {
		maxDate = minDate.Add(24 * time.Hour)
	}
	return minDate, maxDate
}

func hasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
