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

// Gantt chart layout constants.
const (
	ganttHoursPerDay      = 24
	ganttDaysPerWeek      = 7
	ganttTitlePadding     = 10
	ganttPixelsPerDay     = 20
	ganttMinChartWidth    = 200
	ganttMaxChartWidth    = 2000
	ganttDailyTickLimit   = 14
	ganttMonthlyTickLimit = 90
	ganttMonthlyTickDays  = 30
)

var ganttDurationRe = regexp.MustCompile(`^(\d+)([dwmhDWMH])$`)

// mermaidDateToGoLayout converts mermaid dateFormat tokens to Go time layout.
func mermaidDateToGoLayout(format string) string {
	replacer := strings.NewReplacer(
		"YYYY", "2006", "YY", "06",
		"MM", "01", "DD", "02",
		"HH", "15", "mm", "04", "ss", "05",
	)
	return replacer.Replace(format)
}

// parseMermaidDuration converts a mermaid duration string to time.Duration.
func parseMermaidDuration(durStr string) time.Duration {
	match := ganttDurationRe.FindStringSubmatch(durStr)
	if match == nil {
		return 0
	}
	count, err := strconv.Atoi(match[1]) // regex guarantees digits
	if err != nil {
		return 0
	}
	switch strings.ToLower(match[2]) {
	case "d":
		return time.Duration(count) * ganttHoursPerDay * time.Hour
	case "w":
		return time.Duration(count) * ganttDaysPerWeek * ganttHoursPerDay * time.Hour
	case "h":
		return time.Duration(count) * time.Hour
	case "m":
		return time.Duration(count) * time.Minute
	default:
		return 0
	}
}

// isExcluded checks if a date should be excluded based on the excludes list.
func isExcluded(checkTime time.Time, excludes []string, goLayout string) bool {
	dayName := strings.ToLower(checkTime.Weekday().String())
	for _, ex := range excludes {
		ex = strings.ToLower(strings.TrimSpace(ex))
		if ex == "weekends" && (checkTime.Weekday() == time.Saturday || checkTime.Weekday() == time.Sunday) {
			return true
		}
		if ex == dayName {
			return true
		}
		// Try parsing as a date.
		if exDate, err := time.Parse(goLayout, ex); err == nil {
			if checkTime.Year() == exDate.Year() && checkTime.YearDay() == exDate.YearDay() {
				return true
			}
		}
	}
	return false
}

// addWorkingDays adds n working days to start, skipping excluded days.
func addWorkingDays(start time.Time, days int, excludes []string, goLayout string) time.Time {
	if len(excludes) == 0 {
		return start.Add(time.Duration(days) * ganttHoursPerDay * time.Hour)
	}
	current := start
	added := 0
	for added < days {
		current = current.Add(ganttHoursPerDay * time.Hour)
		if !isExcluded(current, excludes, goLayout) {
			added++
		}
	}
	return current
}

type resolvedTask struct {
	Start time.Time
	End   time.Time
}

func computeGanttLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	goLayout := mermaidDateToGoLayout(graph.GanttDateFormat)
	sidePad := cfg.Gantt.SidePadding
	topPad := cfg.Gantt.TopPadding
	barH := cfg.Gantt.BarHeight
	barGap := cfg.Gantt.BarGap

	// Title height.
	var titleHeight float32
	if graph.GanttTitle != "" {
		titleHeight = th.FontSize + ganttTitlePadding
	}

	// Collect all tasks in order for resolution.
	totalTaskCount := 0
	for _, sec := range graph.GanttSections {
		totalTaskCount += len(sec.Tasks)
	}
	allTasks := make([]*ir.GanttTask, 0, totalTaskCount)
	for _, sec := range graph.GanttSections {
		allTasks = append(allTasks, sec.Tasks...)
	}

	// Resolve all task dates.
	resolved := make(map[string]resolvedTask)
	var prevEnd time.Time

	for _, task := range allTasks {
		start, end := resolveGanttTaskDates(task, resolved, &prevEnd, goLayout, graph.GanttExcludes)
		if task.ID != "" {
			resolved[task.ID] = resolvedTask{Start: start, End: end}
		}
		prevEnd = end
	}

	// Find global date range.
	minDate, maxDate := ganttDateRange(allTasks, resolved, goLayout, graph.GanttExcludes)

	totalDays := maxDate.Sub(minDate).Hours() / ganttHoursPerDay
	if totalDays < 1 {
		totalDays = 1
	}

	chartW := float32(totalDays) * ganttPixelsPerDay
	if chartW < ganttMinChartWidth {
		chartW = ganttMinChartWidth
	}
	if chartW > ganttMaxChartWidth {
		chartW = ganttMaxChartWidth
	}

	chartX := sidePad
	chartY := titleHeight + topPad

	// dateToX converts a date to an X pixel position.
	dateToX := func(dateVal time.Time) float32 {
		days := dateVal.Sub(minDate).Hours() / ganttHoursPerDay
		return chartX + float32(days/totalDays)*chartW
	}

	// Build sections and tasks.
	sections := make([]GanttSectionLayout, 0, len(graph.GanttSections))
	curY := chartY
	prevEnd = time.Time{}

	for secIdx, sec := range graph.GanttSections {
		tasks := make([]GanttTaskLayout, 0, len(sec.Tasks))
		secStartY := curY

		for _, task := range sec.Tasks {
			start, end := resolveGanttTaskDates(task, resolved, &prevEnd, goLayout, graph.GanttExcludes)

			taskX := dateToX(start)
			taskW := dateToX(end) - taskX
			if taskW < 1 {
				taskW = 1
			}

			tasks = append(tasks, GanttTaskLayout{
				ID:          task.ID,
				Label:       task.Label,
				X:           taskX,
				Y:           curY,
				Width:       taskW,
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
			color = th.GanttSectionColors[secIdx%len(th.GanttSectionColors)]
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
	tickDays := ganttDaysPerWeek
	if totalDays < ganttDailyTickLimit {
		tickDays = 1
	} else if totalDays > ganttMonthlyTickLimit {
		tickDays = ganttMonthlyTickDays
	}
	for dateVal := minDate; !dateVal.After(maxDate); dateVal = dateVal.AddDate(0, 0, tickDays) {
		axisTicks = append(axisTicks, GanttAxisTick{
			Label: dateVal.Format("2006-01-02"),
			X:     dateToX(dateVal),
		})
	}

	// Today marker.
	today := time.Now()
	showToday := graph.GanttTodayMarker != "off" && !today.Before(minDate) && !today.After(maxDate)
	var todayX float32
	if showToday {
		todayX = dateToX(today)
	}

	totalW := sidePad*2 + chartW
	totalH := curY + topPad

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: GanttData{
			Sections:        sections,
			Title:           graph.GanttTitle,
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
		if parsed, err := time.Parse(goLayout, task.StartStr); err == nil {
			start = parsed
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
		days := int(dur.Hours() / ganttHoursPerDay)
		if days > 0 {
			end = addWorkingDays(start, days, excludes, goLayout)
		} else {
			end = start.Add(dur)
		}
	} else if parsed, err := time.Parse(goLayout, task.EndStr); err == nil {
		end = parsed
	} else {
		end = start.Add(ganttHoursPerDay * time.Hour)
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
		maxDate = minDate.Add(ganttHoursPerDay * time.Hour)
	}
	return minDate, maxDate
}

func hasTag(tags []string, tag string) bool {
	for _, tagVal := range tags {
		if tagVal == tag {
			return true
		}
	}
	return false
}
