# Phase 6: Timeline, Gantt & GitGraph Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add Timeline, Gantt chart, and GitGraph diagram support with full date arithmetic for Gantt.

**Architecture:** Three new diagram types, each with IR types, parser, config, theme, layout, and renderer. Timeline uses horizontal period columns. Gantt uses date-mapped horizontal bars with `time.Parse`, duration math, dependency resolution, and excludes. GitGraph uses swim-lane branch layout simulating git operations. All follow the established per-diagram-type pattern.

**Tech Stack:** Go stdlib (`time`, `math`, `fmt`, `regexp`, `strconv`, `strings`, `sort`), existing `textmetrics`, `theme`, `config` packages.

---

### Task 1: IR types — Timeline

**Files:**
- Create: `ir/timeline.go`
- Create: `ir/timeline_test.go`
- Modify: `ir/graph.go:91` (add Timeline fields after Quadrant)

**Step 1: Write the test**

```go
// ir/timeline_test.go
package ir

import "testing"

func TestTimelineEventDefaults(t *testing.T) {
	e := &TimelineEvent{Text: "Launch"}
	if e.Text != "Launch" {
		t.Errorf("Text = %q, want %q", e.Text, "Launch")
	}
}

func TestTimelinePeriodDefaults(t *testing.T) {
	p := &TimelinePeriod{
		Title:  "2024 Q1",
		Events: []*TimelineEvent{{Text: "Start"}, {Text: "Hire"}},
	}
	if p.Title != "2024 Q1" {
		t.Errorf("Title = %q, want %q", p.Title, "2024 Q1")
	}
	if len(p.Events) != 2 {
		t.Errorf("Events = %d, want 2", len(p.Events))
	}
}

func TestGraphTimelineFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Timeline
	g.TimelineTitle = "Project"
	g.TimelineSections = append(g.TimelineSections, &TimelineSection{
		Title: "Phase 1",
		Periods: []*TimelinePeriod{
			{Title: "Jan", Events: []*TimelineEvent{{Text: "Kickoff"}}},
		},
	})
	if g.TimelineTitle != "Project" {
		t.Errorf("TimelineTitle = %q, want %q", g.TimelineTitle, "Project")
	}
	if len(g.TimelineSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(g.TimelineSections))
	}
	if len(g.TimelineSections[0].Periods) != 1 {
		t.Fatalf("Periods = %d, want 1", len(g.TimelineSections[0].Periods))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestTimeline -v`
Expected: FAIL — `TimelineEvent` undefined

**Step 3: Write minimal implementation**

```go
// ir/timeline.go
package ir

// TimelineEvent represents a single event within a time period.
type TimelineEvent struct {
	Text string
}

// TimelinePeriod represents a time period with one or more events.
type TimelinePeriod struct {
	Title  string
	Events []*TimelineEvent
}

// TimelineSection groups time periods under a named section.
type TimelineSection struct {
	Title   string
	Periods []*TimelinePeriod
}
```

Add to `ir/graph.go` after the Quadrant diagram fields block (after line 90):

```go
	// Timeline diagram fields
	TimelineSections []*TimelineSection
	TimelineTitle    string
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestTimeline -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/timeline.go ir/timeline_test.go ir/graph.go
git commit -m "feat(ir): add Timeline diagram types"
```

---

### Task 2: IR types — Gantt

**Files:**
- Create: `ir/gantt.go`
- Create: `ir/gantt_test.go`
- Modify: `ir/graph.go` (add Gantt fields after Timeline)

**Step 1: Write the test**

```go
// ir/gantt_test.go
package ir

import "testing"

func TestGanttTaskDefaults(t *testing.T) {
	task := &GanttTask{
		ID:       "t1",
		Label:    "Design",
		StartStr: "2024-01-01",
		EndStr:   "10d",
		Tags:     []string{"crit"},
	}
	if task.ID != "t1" {
		t.Errorf("ID = %q, want %q", task.ID, "t1")
	}
	if len(task.Tags) != 1 || task.Tags[0] != "crit" {
		t.Errorf("Tags = %v, want [crit]", task.Tags)
	}
}

func TestGanttSectionDefaults(t *testing.T) {
	s := &GanttSection{
		Title: "Development",
		Tasks: []*GanttTask{
			{ID: "d1", Label: "Code"},
		},
	}
	if s.Title != "Development" {
		t.Errorf("Title = %q", s.Title)
	}
	if len(s.Tasks) != 1 {
		t.Errorf("Tasks = %d, want 1", len(s.Tasks))
	}
}

func TestGraphGanttFields(t *testing.T) {
	g := NewGraph()
	g.Kind = Gantt
	g.GanttTitle = "Project"
	g.GanttDateFormat = "YYYY-MM-DD"
	g.GanttAxisFormat = "%Y-%m-%d"
	g.GanttExcludes = []string{"weekends"}
	g.GanttSections = append(g.GanttSections, &GanttSection{
		Title: "Dev",
		Tasks: []*GanttTask{{ID: "t1", Label: "Code"}},
	})

	if g.GanttTitle != "Project" {
		t.Errorf("GanttTitle = %q", g.GanttTitle)
	}
	if g.GanttDateFormat != "YYYY-MM-DD" {
		t.Errorf("GanttDateFormat = %q", g.GanttDateFormat)
	}
	if len(g.GanttExcludes) != 1 {
		t.Errorf("GanttExcludes = %d, want 1", len(g.GanttExcludes))
	}
	if len(g.GanttSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(g.GanttSections))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestGantt -v`
Expected: FAIL — `GanttTask` undefined

**Step 3: Write minimal implementation**

```go
// ir/gantt.go
package ir

// GanttTask represents a single task in a Gantt chart.
type GanttTask struct {
	ID       string   // Optional task ID for dependencies
	Label    string   // Display name
	StartStr string   // Start: date string, "after t1", or empty (follows previous)
	EndStr   string   // End: duration string ("5d", "2w") or date string
	Tags     []string // Status tags: done, active, crit, milestone
	AfterIDs []string // Task IDs this depends on (parsed from "after t1 t2")
	UntilID  string   // Task ID this runs until
}

// GanttSection groups tasks under a named section.
type GanttSection struct {
	Title string
	Tasks []*GanttTask
}
```

Add to `ir/graph.go` after Timeline fields:

```go
	// Gantt diagram fields
	GanttSections    []*GanttSection
	GanttTitle       string
	GanttDateFormat  string
	GanttAxisFormat  string
	GanttExcludes    []string
	GanttTickInterval string
	GanttTodayMarker string
	GanttWeekday     string
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestGantt -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/gantt.go ir/gantt_test.go ir/graph.go
git commit -m "feat(ir): add Gantt diagram types"
```

---

### Task 3: IR types — GitGraph

**Files:**
- Create: `ir/gitgraph.go`
- Create: `ir/gitgraph_test.go`
- Modify: `ir/graph.go` (add GitGraph fields after Gantt)

**Step 1: Write the test**

```go
// ir/gitgraph_test.go
package ir

import "testing"

func TestGitCommitTypeString(t *testing.T) {
	tests := []struct {
		ct   GitCommitType
		want string
	}{
		{GitCommitNormal, "NORMAL"},
		{GitCommitReverse, "REVERSE"},
		{GitCommitHighlight, "HIGHLIGHT"},
	}
	for _, tt := range tests {
		if got := tt.ct.String(); got != tt.want {
			t.Errorf("%d.String() = %q, want %q", tt.ct, got, tt.want)
		}
	}
}

func TestGitActionInterface(t *testing.T) {
	// Verify all action types implement GitAction.
	actions := []GitAction{
		&GitCommit{ID: "c1"},
		&GitBranch{Name: "dev"},
		&GitCheckout{Branch: "dev"},
		&GitMerge{Branch: "dev"},
		&GitCherryPick{ID: "c1"},
	}
	if len(actions) != 5 {
		t.Errorf("actions = %d, want 5", len(actions))
	}
}

func TestGraphGitGraphFields(t *testing.T) {
	g := NewGraph()
	g.Kind = GitGraph
	g.GitMainBranch = "main"
	g.GitActions = append(g.GitActions,
		&GitCommit{ID: "init", Tag: "v0.1"},
		&GitBranch{Name: "develop", Order: 1},
		&GitCheckout{Branch: "develop"},
		&GitCommit{ID: "feat1", Type: GitCommitHighlight},
		&GitCheckout{Branch: "main"},
		&GitMerge{Branch: "develop", ID: "merge1", Tag: "v1.0"},
	)

	if g.GitMainBranch != "main" {
		t.Errorf("GitMainBranch = %q", g.GitMainBranch)
	}
	if len(g.GitActions) != 6 {
		t.Errorf("GitActions = %d, want 6", len(g.GitActions))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./ir/ -run TestGit -v`
Expected: FAIL — `GitCommitType` undefined

**Step 3: Write minimal implementation**

```go
// ir/gitgraph.go
package ir

// GitCommitType represents the visual type of a git commit.
type GitCommitType int

const (
	GitCommitNormal GitCommitType = iota
	GitCommitReverse
	GitCommitHighlight
)

// String returns the mermaid keyword for the commit type.
func (t GitCommitType) String() string {
	switch t {
	case GitCommitReverse:
		return "REVERSE"
	case GitCommitHighlight:
		return "HIGHLIGHT"
	default:
		return "NORMAL"
	}
}

// GitAction is a sealed interface for git operations in a gitGraph diagram.
type GitAction interface {
	gitAction()
}

// GitCommit represents a commit operation.
type GitCommit struct {
	ID   string
	Tag  string
	Type GitCommitType
}

func (*GitCommit) gitAction() {}

// GitBranch represents a branch creation.
type GitBranch struct {
	Name  string
	Order int // Display order (-1 = unset)
}

func (*GitBranch) gitAction() {}

// GitCheckout represents switching to a branch.
type GitCheckout struct {
	Branch string
}

func (*GitCheckout) gitAction() {}

// GitMerge represents merging a branch into the current branch.
type GitMerge struct {
	Branch string
	ID     string
	Tag    string
	Type   GitCommitType
}

func (*GitMerge) gitAction() {}

// GitCherryPick represents cherry-picking a commit.
type GitCherryPick struct {
	ID     string // Source commit ID
	Parent string // Parent ID for merge commits
}

func (*GitCherryPick) gitAction() {}
```

Add to `ir/graph.go` after Gantt fields:

```go
	// GitGraph diagram fields
	GitActions    []GitAction
	GitMainBranch string
```

**Step 4: Run test to verify it passes**

Run: `go test ./ir/ -run TestGit -v`
Expected: PASS

**Step 5: Commit**

```bash
git add ir/gitgraph.go ir/gitgraph_test.go ir/graph.go
git commit -m "feat(ir): add GitGraph diagram types"
```

---

### Task 4: Config and Theme — Timeline, Gantt, GitGraph

**Files:**
- Modify: `config/config.go` (add 3 config structs + Layout fields + defaults)
- Modify: `config/config_test.go` (add 3 tests)
- Modify: `theme/theme.go` (add theme fields + values in Modern/MermaidDefault)
- Modify: `theme/theme_test.go` (add tests)

**Step 1: Write the tests**

Add to `config/config_test.go`:

```go
func TestDefaultLayoutTimelineConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Timeline.PeriodWidth != 150 {
		t.Errorf("Timeline.PeriodWidth = %f, want 150", cfg.Timeline.PeriodWidth)
	}
	if cfg.Timeline.EventHeight != 30 {
		t.Errorf("Timeline.EventHeight = %f, want 30", cfg.Timeline.EventHeight)
	}
}

func TestDefaultLayoutGanttConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.Gantt.BarHeight != 20 {
		t.Errorf("Gantt.BarHeight = %f, want 20", cfg.Gantt.BarHeight)
	}
	if cfg.Gantt.SidePadding != 75 {
		t.Errorf("Gantt.SidePadding = %f, want 75", cfg.Gantt.SidePadding)
	}
}

func TestDefaultLayoutGitGraphConfig(t *testing.T) {
	cfg := DefaultLayout()
	if cfg.GitGraph.CommitRadius != 8 {
		t.Errorf("GitGraph.CommitRadius = %f, want 8", cfg.GitGraph.CommitRadius)
	}
	if cfg.GitGraph.CommitSpacing != 60 {
		t.Errorf("GitGraph.CommitSpacing = %f, want 60", cfg.GitGraph.CommitSpacing)
	}
	if cfg.GitGraph.BranchSpacing != 40 {
		t.Errorf("GitGraph.BranchSpacing = %f, want 40", cfg.GitGraph.BranchSpacing)
	}
}
```

Add to `theme/theme_test.go`:

```go
func TestModernTimelineColors(t *testing.T) {
	th := Modern()
	if len(th.TimelineSectionColors) < 4 {
		t.Errorf("TimelineSectionColors = %d, want >= 4", len(th.TimelineSectionColors))
	}
	if th.TimelineEventFill == "" {
		t.Error("TimelineEventFill is empty")
	}
}

func TestModernGanttColors(t *testing.T) {
	th := Modern()
	if th.GanttTaskFill == "" {
		t.Error("GanttTaskFill is empty")
	}
	if th.GanttCritFill == "" {
		t.Error("GanttCritFill is empty")
	}
	if len(th.GanttSectionColors) < 4 {
		t.Errorf("GanttSectionColors = %d, want >= 4", len(th.GanttSectionColors))
	}
}

func TestModernGitGraphColors(t *testing.T) {
	th := Modern()
	if len(th.GitBranchColors) < 8 {
		t.Errorf("GitBranchColors = %d, want >= 8", len(th.GitBranchColors))
	}
	if th.GitCommitFill == "" {
		t.Error("GitCommitFill is empty")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./config/ -v && go test ./theme/ -v`
Expected: FAIL — `cfg.Timeline` undefined

**Step 3: Write minimal implementation**

Add config structs to `config/config.go` after QuadrantConfig:

```go
// TimelineConfig holds timeline diagram layout options.
type TimelineConfig struct {
	PeriodWidth    float32
	EventHeight    float32
	SectionPadding float32
	PaddingX       float32
	PaddingY       float32
}

// GanttConfig holds Gantt chart layout options.
type GanttConfig struct {
	BarHeight            float32
	BarGap               float32
	TopPadding           float32
	SidePadding          float32
	GridLineStartPadding float32
	FontSize             float32
	SectionFontSize      float32
	NumberSectionStyles  int
}

// GitGraphConfig holds GitGraph diagram layout options.
type GitGraphConfig struct {
	CommitRadius  float32
	CommitSpacing float32
	BranchSpacing float32
	PaddingX      float32
	PaddingY      float32
	TagFontSize   float32
}
```

Add fields to `Layout` struct (after `Quadrant QuadrantConfig`):

```go
	Timeline TimelineConfig
	Gantt    GanttConfig
	GitGraph GitGraphConfig
```

Add defaults in `DefaultLayout()` (after Quadrant defaults):

```go
		Timeline: TimelineConfig{
			PeriodWidth:    150,
			EventHeight:    30,
			SectionPadding: 10,
			PaddingX:       20,
			PaddingY:       20,
		},
		Gantt: GanttConfig{
			BarHeight:            20,
			BarGap:               4,
			TopPadding:           50,
			SidePadding:          75,
			GridLineStartPadding: 35,
			FontSize:             11,
			SectionFontSize:      11,
			NumberSectionStyles:  4,
		},
		GitGraph: GitGraphConfig{
			CommitRadius:  8,
			CommitSpacing: 60,
			BranchSpacing: 40,
			PaddingX:      30,
			PaddingY:      30,
			TagFontSize:   11,
		},
```

Add theme fields to `Theme` struct in `theme/theme.go` (after Quadrant fields):

```go
	// Timeline diagram colors
	TimelineSectionColors []string
	TimelineEventFill     string
	TimelineEventBorder   string

	// Gantt diagram colors
	GanttTaskFill        string
	GanttTaskBorder      string
	GanttCritFill        string
	GanttCritBorder      string
	GanttDoneFill        string
	GanttActiveFill      string
	GanttMilestoneFill   string
	GanttGridColor       string
	GanttTodayMarkerColor string
	GanttSectionColors   []string

	// GitGraph diagram colors
	GitBranchColors  []string
	GitCommitFill    string
	GitCommitStroke  string
	GitTagFill       string
	GitTagBorder     string
	GitHighlightFill string
```

Add values in `Modern()`:

```go
		TimelineSectionColors: []string{"#E8EFF5", "#F0E8F0", "#E8F5E8", "#FFF8E1"},
		TimelineEventFill:     "#4C78A8",
		TimelineEventBorder:   "#3B6492",

		GanttTaskFill:         "#4C78A8",
		GanttTaskBorder:       "#3B6492",
		GanttCritFill:         "#E45756",
		GanttCritBorder:       "#CC3333",
		GanttDoneFill:         "#B0C4DE",
		GanttActiveFill:       "#72B7B2",
		GanttMilestoneFill:    "#F58518",
		GanttGridColor:        "#E0E0E0",
		GanttTodayMarkerColor: "#E45756",
		GanttSectionColors:    []string{"#F0F4F8", "#FFF8E1", "#F0E8F0", "#E8F5E8"},

		GitBranchColors: []string{
			"#4C78A8", "#E45756", "#54A24B", "#F58518",
			"#72B7B2", "#B279A2", "#EECA3B", "#FF9DA6",
		},
		GitCommitFill:    "#333344",
		GitCommitStroke:  "#333344",
		GitTagFill:       "#EECA3B",
		GitTagBorder:     "#C9A820",
		GitHighlightFill: "#F58518",
```

Add values in `MermaidDefault()`:

```go
		TimelineSectionColors: []string{"#ffffde", "#f0ece8", "#e8f0e8", "#ece8f0"},
		TimelineEventFill:     "#ECECFF",
		TimelineEventBorder:   "#9370DB",

		GanttTaskFill:         "#8a90dd",
		GanttTaskBorder:       "#534fbc",
		GanttCritFill:         "#ff8888",
		GanttCritBorder:       "#ff0000",
		GanttDoneFill:         "#d3d3d3",
		GanttActiveFill:       "#8a90dd",
		GanttMilestoneFill:    "#E76F51",
		GanttGridColor:        "#ddd",
		GanttTodayMarkerColor: "#d42",
		GanttSectionColors:    []string{"#ffffde", "#ffffff", "#ffffde", "#ffffff"},

		GitBranchColors: []string{
			"#9370DB", "#ff0000", "#00cc00", "#F58518",
			"#48A9A6", "#E76F51", "#D08AC0", "#F7B7A3",
		},
		GitCommitFill:    "#333",
		GitCommitStroke:  "#333",
		GitTagFill:       "#ffffde",
		GitTagBorder:     "#aaaa33",
		GitHighlightFill: "#ff0000",
```

**Step 4: Run tests to verify they pass**

Run: `go test ./config/ -v && go test ./theme/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add config/config.go config/config_test.go theme/theme.go theme/theme_test.go
git commit -m "feat(config,theme): add Timeline, Gantt, and GitGraph config and theme"
```

---

### Task 5: Parser — Timeline

**Files:**
- Create: `parser/timeline.go`
- Create: `parser/timeline_test.go`
- Modify: `parser/parser.go:37` (add `case ir.Timeline`)

**Step 1: Write the test**

```go
// parser/timeline_test.go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseTimelineBasic(t *testing.T) {
	input := `timeline
    title History of Social Media
    2002 : LinkedIn
    2004 : Facebook : Google
    2005 : YouTube`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Timeline {
		t.Errorf("Kind = %v, want Timeline", g.Kind)
	}
	if g.TimelineTitle != "History of Social Media" {
		t.Errorf("Title = %q", g.TimelineTitle)
	}
	// No explicit sections, so 1 implicit section.
	if len(g.TimelineSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(g.TimelineSections))
	}
	sec := g.TimelineSections[0]
	if len(sec.Periods) != 3 {
		t.Fatalf("Periods = %d, want 3", len(sec.Periods))
	}
	if sec.Periods[0].Title != "2002" {
		t.Errorf("Period[0].Title = %q", sec.Periods[0].Title)
	}
	if len(sec.Periods[0].Events) != 1 {
		t.Fatalf("Period[0].Events = %d, want 1", len(sec.Periods[0].Events))
	}
	if sec.Periods[0].Events[0].Text != "LinkedIn" {
		t.Errorf("Period[0].Events[0] = %q", sec.Periods[0].Events[0].Text)
	}
	// 2004 has 2 events.
	if len(sec.Periods[1].Events) != 2 {
		t.Fatalf("Period[1].Events = %d, want 2", len(sec.Periods[1].Events))
	}
}

func TestParseTimelineSections(t *testing.T) {
	input := `timeline
    title Product Timeline
    section Phase 1
        January : Research
        February : Design
    section Phase 2
        March : Development`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if len(g.TimelineSections) != 2 {
		t.Fatalf("Sections = %d, want 2", len(g.TimelineSections))
	}
	if g.TimelineSections[0].Title != "Phase 1" {
		t.Errorf("Section[0].Title = %q", g.TimelineSections[0].Title)
	}
	if len(g.TimelineSections[0].Periods) != 2 {
		t.Errorf("Section[0].Periods = %d, want 2", len(g.TimelineSections[0].Periods))
	}
	if g.TimelineSections[1].Title != "Phase 2" {
		t.Errorf("Section[1].Title = %q", g.TimelineSections[1].Title)
	}
}

func TestParseTimelineContinuationEvents(t *testing.T) {
	input := `timeline
    2023 : Event A
         : Event B
         : Event C`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	sec := g.TimelineSections[0]
	if len(sec.Periods) != 1 {
		t.Fatalf("Periods = %d, want 1", len(sec.Periods))
	}
	if len(sec.Periods[0].Events) != 3 {
		t.Fatalf("Events = %d, want 3", len(sec.Periods[0].Events))
	}
	if sec.Periods[0].Events[2].Text != "Event C" {
		t.Errorf("Events[2] = %q", sec.Periods[0].Events[2].Text)
	}
}

func TestParseTimelineMinimal(t *testing.T) {
	input := `timeline
    2024 : Launch`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.TimelineSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(out.Graph.TimelineSections))
	}
	if len(out.Graph.TimelineSections[0].Periods) != 1 {
		t.Fatalf("Periods = %d, want 1", len(out.Graph.TimelineSections[0].Periods))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParseTimeline -v`
Expected: FAIL — falls through to parseFlowchart

**Step 3: Write minimal implementation**

```go
// parser/timeline.go
package parser

import (
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
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
```

Add to `parser/parser.go` Parse switch (before `default`):

```go
	case ir.Timeline:
		return parseTimeline(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParseTimeline -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/timeline.go parser/timeline_test.go parser/parser.go
git commit -m "feat(parser): add Timeline diagram parser"
```

---

### Task 6: Parser — Gantt

**Files:**
- Create: `parser/gantt.go`
- Create: `parser/gantt_test.go`
- Modify: `parser/parser.go` (add `case ir.Gantt`)

**Step 1: Write the test**

```go
// parser/gantt_test.go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseGanttBasic(t *testing.T) {
	input := `gantt
    title A Gantt Diagram
    dateFormat YYYY-MM-DD
    section Development
        Design :d1, 2024-01-01, 10d
        Coding :d2, after d1, 20d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.Gantt {
		t.Errorf("Kind = %v, want Gantt", g.Kind)
	}
	if g.GanttTitle != "A Gantt Diagram" {
		t.Errorf("Title = %q", g.GanttTitle)
	}
	if g.GanttDateFormat != "YYYY-MM-DD" {
		t.Errorf("DateFormat = %q", g.GanttDateFormat)
	}
	if len(g.GanttSections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(g.GanttSections))
	}
	sec := g.GanttSections[0]
	if sec.Title != "Development" {
		t.Errorf("Section.Title = %q", sec.Title)
	}
	if len(sec.Tasks) != 2 {
		t.Fatalf("Tasks = %d, want 2", len(sec.Tasks))
	}
	t1 := sec.Tasks[0]
	if t1.ID != "d1" || t1.Label != "Design" {
		t.Errorf("Task[0] = %+v", t1)
	}
	if t1.StartStr != "2024-01-01" || t1.EndStr != "10d" {
		t.Errorf("Task[0] start=%q end=%q", t1.StartStr, t1.EndStr)
	}
	t2 := sec.Tasks[1]
	if t2.ID != "d2" || len(t2.AfterIDs) != 1 || t2.AfterIDs[0] != "d1" {
		t.Errorf("Task[1] = %+v", t2)
	}
}

func TestParseGanttTags(t *testing.T) {
	input := `gantt
    dateFormat YYYY-MM-DD
    section Tasks
        Done task :done, 2024-01-01, 5d
        Critical :crit, active, t3, 2024-01-06, 10d
        Milestone :milestone, m1, 2024-02-01, 0d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	tasks := out.Graph.GanttSections[0].Tasks

	if len(tasks[0].Tags) != 1 || tasks[0].Tags[0] != "done" {
		t.Errorf("Task[0].Tags = %v", tasks[0].Tags)
	}
	if len(tasks[1].Tags) != 2 {
		t.Errorf("Task[1].Tags = %v, want 2 tags", tasks[1].Tags)
	}
	if tasks[1].ID != "t3" {
		t.Errorf("Task[1].ID = %q, want t3", tasks[1].ID)
	}
	if len(tasks[2].Tags) != 1 || tasks[2].Tags[0] != "milestone" {
		t.Errorf("Task[2].Tags = %v", tasks[2].Tags)
	}
	if tasks[2].ID != "m1" {
		t.Errorf("Task[2].ID = %q, want m1", tasks[2].ID)
	}
}

func TestParseGanttDirectives(t *testing.T) {
	input := `gantt
    dateFormat YYYY-MM-DD
    axisFormat %m/%d
    excludes weekends
    tickInterval 1week
    todayMarker off
    weekend friday
    section A
        Task :2024-01-01, 5d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.GanttAxisFormat != "%m/%d" {
		t.Errorf("AxisFormat = %q", g.GanttAxisFormat)
	}
	if len(g.GanttExcludes) != 1 || g.GanttExcludes[0] != "weekends" {
		t.Errorf("Excludes = %v", g.GanttExcludes)
	}
	if g.GanttTickInterval != "1week" {
		t.Errorf("TickInterval = %q", g.GanttTickInterval)
	}
	if g.GanttTodayMarker != "off" {
		t.Errorf("TodayMarker = %q", g.GanttTodayMarker)
	}
	if g.GanttWeekday != "friday" {
		t.Errorf("Weekday = %q", g.GanttWeekday)
	}
}

func TestParseGanttMultipleAfter(t *testing.T) {
	input := `gantt
    dateFormat YYYY-MM-DD
    section A
        Task A :a1, 2024-01-01, 5d
        Task B :b1, 2024-01-01, 3d
        Task C :c1, after a1 b1, 10d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	tasks := out.Graph.GanttSections[0].Tasks
	if len(tasks[2].AfterIDs) != 2 {
		t.Errorf("Task[2].AfterIDs = %v, want 2", tasks[2].AfterIDs)
	}
}

func TestParseGanttNoSection(t *testing.T) {
	input := `gantt
    dateFormat YYYY-MM-DD
    Task A :2024-01-01, 5d
    Task B :2024-01-06, 3d`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.GanttSections) != 1 {
		t.Fatalf("Sections = %d, want 1 (implicit)", len(out.Graph.GanttSections))
	}
	if len(out.Graph.GanttSections[0].Tasks) != 2 {
		t.Errorf("Tasks = %d, want 2", len(out.Graph.GanttSections[0].Tasks))
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParseGantt -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// parser/gantt.go
package parser

import (
	"regexp"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
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
			g.GanttExcludes = append(g.GanttExcludes, strings.Split(val, ",")...)
			for i := range g.GanttExcludes {
				g.GanttExcludes[i] = strings.TrimSpace(g.GanttExcludes[i])
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
```

Add to `parser/parser.go` Parse switch:

```go
	case ir.Gantt:
		return parseGantt(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParseGantt -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/gantt.go parser/gantt_test.go parser/parser.go
git commit -m "feat(parser): add Gantt chart parser"
```

---

### Task 7: Parser — GitGraph

**Files:**
- Create: `parser/gitgraph.go`
- Create: `parser/gitgraph_test.go`
- Modify: `parser/parser.go` (add `case ir.GitGraph`)

**Step 1: Write the test**

```go
// parser/gitgraph_test.go
package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseGitGraphBasic(t *testing.T) {
	input := `gitGraph
    commit
    commit id: "feat1"
    branch develop
    checkout develop
    commit
    checkout main
    merge develop`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g := out.Graph
	if g.Kind != ir.GitGraph {
		t.Errorf("Kind = %v, want GitGraph", g.Kind)
	}
	if len(g.GitActions) != 7 {
		t.Fatalf("Actions = %d, want 7", len(g.GitActions))
	}

	// First action should be a commit.
	if _, ok := g.GitActions[0].(*ir.GitCommit); !ok {
		t.Errorf("Action[0] type = %T, want *GitCommit", g.GitActions[0])
	}
	// Second commit has ID.
	c1 := g.GitActions[1].(*ir.GitCommit)
	if c1.ID != "feat1" {
		t.Errorf("Action[1].ID = %q, want feat1", c1.ID)
	}
	// Branch.
	br := g.GitActions[2].(*ir.GitBranch)
	if br.Name != "develop" {
		t.Errorf("Branch.Name = %q", br.Name)
	}
	// Merge.
	mg := g.GitActions[6].(*ir.GitMerge)
	if mg.Branch != "develop" {
		t.Errorf("Merge.Branch = %q", mg.Branch)
	}
}

func TestParseGitGraphCommitOptions(t *testing.T) {
	input := `gitGraph
    commit id: "c1" tag: "v1.0" type: HIGHLIGHT
    commit type: REVERSE`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	c0 := out.Graph.GitActions[0].(*ir.GitCommit)
	if c0.ID != "c1" || c0.Tag != "v1.0" || c0.Type != ir.GitCommitHighlight {
		t.Errorf("commit[0] = %+v", c0)
	}
	c1 := out.Graph.GitActions[1].(*ir.GitCommit)
	if c1.Type != ir.GitCommitReverse {
		t.Errorf("commit[1].Type = %v, want REVERSE", c1.Type)
	}
}

func TestParseGitGraphBranchOrder(t *testing.T) {
	input := `gitGraph
    commit
    branch develop order: 2
    branch feature order: 1`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	b0 := out.Graph.GitActions[1].(*ir.GitBranch)
	if b0.Name != "develop" || b0.Order != 2 {
		t.Errorf("branch[0] = %+v", b0)
	}
	b1 := out.Graph.GitActions[2].(*ir.GitBranch)
	if b1.Name != "feature" || b1.Order != 1 {
		t.Errorf("branch[1] = %+v", b1)
	}
}

func TestParseGitGraphMergeOptions(t *testing.T) {
	input := `gitGraph
    commit
    branch dev
    checkout dev
    commit
    checkout main
    merge dev id: "m1" tag: "v2.0" type: HIGHLIGHT`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	mg := out.Graph.GitActions[5].(*ir.GitMerge)
	if mg.Branch != "dev" || mg.ID != "m1" || mg.Tag != "v2.0" || mg.Type != ir.GitCommitHighlight {
		t.Errorf("merge = %+v", mg)
	}
}

func TestParseGitGraphCherryPick(t *testing.T) {
	input := `gitGraph
    commit id: "src"
    branch dev
    checkout dev
    cherry-pick id: "src"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	cp := out.Graph.GitActions[3].(*ir.GitCherryPick)
	if cp.ID != "src" {
		t.Errorf("cherry-pick.ID = %q", cp.ID)
	}
}

func TestParseGitGraphSwitch(t *testing.T) {
	input := `gitGraph
    commit
    branch dev
    switch dev
    commit`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	co := out.Graph.GitActions[2].(*ir.GitCheckout)
	if co.Branch != "dev" {
		t.Errorf("switch.Branch = %q", co.Branch)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./parser/ -run TestParseGitGraph -v`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
// parser/gitgraph.go
package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	gitKeyValRe = regexp.MustCompile(`(\w+)\s*:\s*(?:"([^"]+)"|(\S+))`)
)

func parseGitGraph(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"

	lines := preprocessInput(input)

	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))

		if strings.HasPrefix(lower, "gitgraph") {
			continue
		}

		if strings.HasPrefix(lower, "commit") {
			g.GitActions = append(g.GitActions, parseGitCommit(line))
			continue
		}

		if strings.HasPrefix(lower, "branch ") {
			g.GitActions = append(g.GitActions, parseGitBranch(line))
			continue
		}

		if strings.HasPrefix(lower, "checkout ") || strings.HasPrefix(lower, "switch ") {
			var rest string
			if strings.HasPrefix(lower, "checkout ") {
				rest = strings.TrimSpace(line[len("checkout "):])
			} else {
				rest = strings.TrimSpace(line[len("switch "):])
			}
			g.GitActions = append(g.GitActions, &ir.GitCheckout{Branch: rest})
			continue
		}

		if strings.HasPrefix(lower, "merge ") {
			g.GitActions = append(g.GitActions, parseGitMerge(line))
			continue
		}

		if strings.HasPrefix(lower, "cherry-pick") {
			g.GitActions = append(g.GitActions, parseGitCherryPick(line))
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}

func parseGitCommit(line string) *ir.GitCommit {
	c := &ir.GitCommit{}
	opts := gitKeyValRe.FindAllStringSubmatch(line, -1)
	for _, m := range opts {
		key := strings.ToLower(m[1])
		val := m[2]
		if val == "" {
			val = m[3]
		}
		switch key {
		case "id":
			c.ID = val
		case "tag":
			c.Tag = val
		case "type":
			c.Type = parseGitCommitType(val)
		}
	}
	return c
}

func parseGitBranch(line string) *ir.GitBranch {
	// "branch <name> order: <n>"
	rest := strings.TrimSpace(line[len("branch "):])
	b := &ir.GitBranch{Order: -1}

	// Check for "order:" option.
	if idx := strings.Index(strings.ToLower(rest), "order:"); idx >= 0 {
		b.Name = strings.TrimSpace(rest[:idx])
		orderStr := strings.TrimSpace(rest[idx+len("order:"):])
		if n, err := strconv.Atoi(orderStr); err == nil {
			b.Order = n
		}
	} else {
		b.Name = rest
	}

	// Strip quotes from name if present.
	b.Name = strings.Trim(b.Name, `"`)

	return b
}

func parseGitMerge(line string) *ir.GitMerge {
	rest := strings.TrimSpace(line[len("merge "):])
	mg := &ir.GitMerge{}

	// Extract the branch name (first word before any key: value pairs).
	parts := strings.Fields(rest)
	if len(parts) > 0 {
		mg.Branch = strings.Trim(parts[0], `"`)
	}

	// Parse key-value options.
	opts := gitKeyValRe.FindAllStringSubmatch(rest, -1)
	for _, m := range opts {
		key := strings.ToLower(m[1])
		val := m[2]
		if val == "" {
			val = m[3]
		}
		switch key {
		case "id":
			mg.ID = val
		case "tag":
			mg.Tag = val
		case "type":
			mg.Type = parseGitCommitType(val)
		}
	}

	return mg
}

func parseGitCherryPick(line string) *ir.GitCherryPick {
	cp := &ir.GitCherryPick{}
	opts := gitKeyValRe.FindAllStringSubmatch(line, -1)
	for _, m := range opts {
		key := strings.ToLower(m[1])
		val := m[2]
		if val == "" {
			val = m[3]
		}
		switch key {
		case "id":
			cp.ID = val
		case "parent":
			cp.Parent = val
		}
	}
	return cp
}

func parseGitCommitType(s string) ir.GitCommitType {
	switch strings.ToUpper(s) {
	case "REVERSE":
		return ir.GitCommitReverse
	case "HIGHLIGHT":
		return ir.GitCommitHighlight
	default:
		return ir.GitCommitNormal
	}
}
```

Add to `parser/parser.go` Parse switch:

```go
	case ir.GitGraph:
		return parseGitGraph(input)
```

**Step 4: Run test to verify it passes**

Run: `go test ./parser/ -run TestParseGitGraph -v`
Expected: PASS

**Step 5: Commit**

```bash
git add parser/gitgraph.go parser/gitgraph_test.go parser/parser.go
git commit -m "feat(parser): add GitGraph diagram parser"
```

---

### Task 8: Layout — Timeline

**Files:**
- Create: `layout/timeline.go`
- Create: `layout/timeline_test.go`
- Modify: `layout/types.go` (add TimelineData types)
- Modify: `layout/layout.go` (add `case ir.Timeline`)

**Step 1: Write the test**

```go
// layout/timeline_test.go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestTimelineLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Timeline
	g.TimelineTitle = "History"
	g.TimelineSections = []*ir.TimelineSection{
		{
			Title: "Early",
			Periods: []*ir.TimelinePeriod{
				{Title: "2002", Events: []*ir.TimelineEvent{{Text: "LinkedIn"}}},
				{Title: "2004", Events: []*ir.TimelineEvent{{Text: "Facebook"}, {Text: "Google"}}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Timeline {
		t.Errorf("Kind = %v, want Timeline", l.Kind)
	}
	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("dimensions = %f x %f", l.Width, l.Height)
	}

	td, ok := l.Diagram.(TimelineData)
	if !ok {
		t.Fatalf("Diagram type = %T, want TimelineData", l.Diagram)
	}
	if len(td.Sections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(td.Sections))
	}
	if len(td.Sections[0].Periods) != 2 {
		t.Fatalf("Periods = %d, want 2", len(td.Sections[0].Periods))
	}
	// Second period has 2 events, so should be taller.
	p0 := td.Sections[0].Periods[0]
	p1 := td.Sections[0].Periods[1]
	if p0.X >= p1.X {
		t.Errorf("Period[0].X=%f >= Period[1].X=%f, want left-to-right", p0.X, p1.X)
	}
}

func TestTimelineLayoutMultipleSections(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Timeline
	g.TimelineSections = []*ir.TimelineSection{
		{Title: "S1", Periods: []*ir.TimelinePeriod{
			{Title: "P1", Events: []*ir.TimelineEvent{{Text: "E1"}}},
		}},
		{Title: "S2", Periods: []*ir.TimelinePeriod{
			{Title: "P2", Events: []*ir.TimelineEvent{{Text: "E2"}}},
		}},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	td := l.Diagram.(TimelineData)
	if len(td.Sections) != 2 {
		t.Fatalf("Sections = %d, want 2", len(td.Sections))
	}
	// S2 should be below S1.
	if td.Sections[0].Y >= td.Sections[1].Y {
		t.Errorf("S1.Y=%f >= S2.Y=%f, want S1 above S2", td.Sections[0].Y, td.Sections[1].Y)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestTimelineLayout -v`
Expected: FAIL — `TimelineData` undefined

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// TimelineData holds timeline-diagram-specific layout data.
type TimelineData struct {
	Sections []TimelineSectionLayout
	Title    string
}

func (TimelineData) diagramData() {}

// TimelineSectionLayout holds positioned section data.
type TimelineSectionLayout struct {
	Title   string
	X, Y    float32
	Width   float32
	Height  float32
	Color   string
	Periods []TimelinePeriodLayout
}

// TimelinePeriodLayout holds positioned period data.
type TimelinePeriodLayout struct {
	Title  string
	X, Y   float32
	Width  float32
	Height float32
	Events []TimelineEventLayout
}

// TimelineEventLayout holds positioned event data.
type TimelineEventLayout struct {
	Text   string
	X, Y   float32
	Width  float32
	Height float32
}
```

Create `layout/timeline.go`:

```go
package layout

import (
	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
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

	// Count total periods across all sections for width.
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
		colorIdx := i % len(th.TimelineSectionColors)
		color := th.TimelineSectionColors[colorIdx]

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
```

Add to `layout/layout.go` ComputeLayout switch (before `default`):

```go
	case ir.Timeline:
		return computeTimelineLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./layout/ -run TestTimelineLayout -v`
Expected: PASS

**Step 5: Commit**

```bash
git add layout/timeline.go layout/timeline_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Timeline diagram layout"
```

---

### Task 9: Layout — Gantt (date arithmetic + dependency resolution)

**Files:**
- Create: `layout/gantt.go`
- Create: `layout/gantt_test.go`
- Modify: `layout/types.go` (add GanttData types)
- Modify: `layout/layout.go` (add `case ir.Gantt`)

This is the most complex layout task. It includes date parsing, duration math, dependency resolution, and excludes.

**Step 1: Write the test**

```go
// layout/gantt_test.go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestGanttLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Gantt
	g.GanttDateFormat = "YYYY-MM-DD"
	g.GanttTitle = "Project"
	g.GanttSections = []*ir.GanttSection{
		{
			Title: "Dev",
			Tasks: []*ir.GanttTask{
				{ID: "t1", Label: "Design", StartStr: "2024-01-01", EndStr: "10d"},
				{ID: "t2", Label: "Code", StartStr: "2024-01-11", EndStr: "20d"},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.Gantt {
		t.Errorf("Kind = %v, want Gantt", l.Kind)
	}
	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("dimensions = %f x %f", l.Width, l.Height)
	}

	gd, ok := l.Diagram.(GanttData)
	if !ok {
		t.Fatalf("Diagram type = %T, want GanttData", l.Diagram)
	}
	if len(gd.Sections) != 1 {
		t.Fatalf("Sections = %d, want 1", len(gd.Sections))
	}
	if len(gd.Sections[0].Tasks) != 2 {
		t.Fatalf("Tasks = %d, want 2", len(gd.Sections[0].Tasks))
	}
	// Task 1 should start before Task 2.
	if gd.Sections[0].Tasks[0].X >= gd.Sections[0].Tasks[1].X {
		t.Errorf("Task1.X=%f >= Task2.X=%f", gd.Sections[0].Tasks[0].X, gd.Sections[0].Tasks[1].X)
	}
	// Task 2 should be wider (20d vs 10d).
	if gd.Sections[0].Tasks[1].Width <= gd.Sections[0].Tasks[0].Width {
		t.Errorf("Task2.Width=%f <= Task1.Width=%f", gd.Sections[0].Tasks[1].Width, gd.Sections[0].Tasks[0].Width)
	}
}

func TestGanttLayoutAfterDependency(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Gantt
	g.GanttDateFormat = "YYYY-MM-DD"
	g.GanttSections = []*ir.GanttSection{
		{
			Tasks: []*ir.GanttTask{
				{ID: "a", Label: "Task A", StartStr: "2024-01-01", EndStr: "5d"},
				{ID: "b", Label: "Task B", StartStr: "after a", EndStr: "3d", AfterIDs: []string{"a"}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	gd := l.Diagram.(GanttData)
	tasks := gd.Sections[0].Tasks
	// Task B should start where Task A ends.
	if tasks[1].X <= tasks[0].X+tasks[0].Width-1 {
		t.Errorf("TaskB.X=%f should start after TaskA ends at %f", tasks[1].X, tasks[0].X+tasks[0].Width)
	}
}

func TestGanttLayoutMilestone(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Gantt
	g.GanttDateFormat = "YYYY-MM-DD"
	g.GanttSections = []*ir.GanttSection{
		{
			Tasks: []*ir.GanttTask{
				{Label: "Release", StartStr: "2024-02-01", EndStr: "0d", Tags: []string{"milestone"}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	gd := l.Diagram.(GanttData)
	task := gd.Sections[0].Tasks[0]
	if !task.IsMilestone {
		t.Error("expected milestone flag")
	}
}

func TestGanttLayoutExcludesWeekends(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Gantt
	g.GanttDateFormat = "YYYY-MM-DD"
	g.GanttExcludes = []string{"weekends"}
	g.GanttSections = []*ir.GanttSection{
		{
			Tasks: []*ir.GanttTask{
				{ID: "t1", Label: "Work", StartStr: "2024-01-01", EndStr: "5d"},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	// Should succeed without panic.
	if l.Width <= 0 {
		t.Errorf("Width = %f", l.Width)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestGanttLayout -v`
Expected: FAIL — `GanttData` undefined

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// GanttData holds Gantt-diagram-specific layout data.
type GanttData struct {
	Sections       []GanttSectionLayout
	Title          string
	AxisTicks      []GanttAxisTick
	TodayMarkerX   float32
	ShowTodayMarker bool
	ChartX         float32
	ChartY         float32
	ChartWidth     float32
	ChartHeight    float32
}

func (GanttData) diagramData() {}

// GanttSectionLayout holds positioned section data.
type GanttSectionLayout struct {
	Title  string
	Y      float32
	Height float32
	Color  string
	Tasks  []GanttTaskLayout
}

// GanttTaskLayout holds positioned task bar data.
type GanttTaskLayout struct {
	ID          string
	Label       string
	X, Y        float32
	Width       float32
	Height      float32
	IsCrit      bool
	IsDone      bool
	IsActive    bool
	IsMilestone bool
}

// GanttAxisTick holds a tick mark on the date axis.
type GanttAxisTick struct {
	Label string
	X     float32
}
```

Create `layout/gantt.go`:

```go
package layout

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

var ganttDurationRe = regexp.MustCompile(`^(\d+)([dwmhDWMH])$`)

// mermaidDateToGoLayout converts mermaid dateFormat tokens to Go time layout.
func mermaidDateToGoLayout(fmt string) string {
	r := strings.NewReplacer(
		"YYYY", "2006", "YY", "06",
		"MM", "01", "DD", "02",
		"HH", "15", "mm", "04", "ss", "05",
	)
	return r.Replace(fmt)
}

// parseDuration converts a mermaid duration string to time.Duration.
func parseMermaidDuration(s string) time.Duration {
	m := ganttDurationRe.FindStringSubmatch(s)
	if m == nil {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
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

	// Resolve all task dates.
	resolved := make(map[string]resolvedTask)
	var allTasks []*ir.GanttTask
	for _, sec := range g.GanttSections {
		allTasks = append(allTasks, sec.Tasks...)
	}

	var prevEnd time.Time
	for _, task := range allTasks {
		var start, end time.Time

		// Resolve start.
		if len(task.AfterIDs) > 0 {
			// Start after latest dependency.
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

		if start.IsZero() && !prevEnd.IsZero() {
			start = prevEnd
		}
		if start.IsZero() {
			start = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		}

		// Resolve end.
		dur := parseMermaidDuration(task.EndStr)
		if dur > 0 {
			days := int(dur.Hours() / 24)
			if days > 0 {
				end = addWorkingDays(start, days, g.GanttExcludes, goLayout)
			} else {
				end = start.Add(dur)
			}
		} else if t, err := time.Parse(goLayout, task.EndStr); err == nil {
			end = t
		} else {
			end = start.Add(24 * time.Hour)
		}

		if task.ID != "" {
			resolved[task.ID] = resolvedTask{Start: start, End: end}
		}
		prevEnd = end
	}

	// Find global date range.
	var minDate, maxDate time.Time
	first := true
	for _, rt := range resolved {
		if first || rt.Start.Before(minDate) {
			minDate = rt.Start
		}
		if first || rt.End.After(maxDate) {
			maxDate = rt.End
		}
		first = false
	}
	// Also check tasks without IDs.
	taskIdx := 0
	prevEnd = time.Time{}
	for _, sec := range g.GanttSections {
		for _, task := range sec.Tasks {
			var start, end time.Time
			if task.ID != "" {
				rt := resolved[task.ID]
				start, end = rt.Start, rt.End
			} else {
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
				if start.IsZero() && !prevEnd.IsZero() {
					start = prevEnd
				}
				if start.IsZero() {
					start = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				}
				dur := parseMermaidDuration(task.EndStr)
				if dur > 0 {
					days := int(dur.Hours() / 24)
					if days > 0 {
						end = addWorkingDays(start, days, g.GanttExcludes, goLayout)
					} else {
						end = start.Add(dur)
					}
				} else if t, err := time.Parse(goLayout, task.EndStr); err == nil {
					end = t
				} else {
					end = start.Add(24 * time.Hour)
				}
			}

			if first || start.Before(minDate) {
				minDate = start
			}
			if first || end.After(maxDate) {
				maxDate = end
			}
			first = false
			prevEnd = end
			taskIdx++
		}
	}

	if minDate.IsZero() {
		minDate = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if maxDate.IsZero() || !maxDate.After(minDate) {
		maxDate = minDate.Add(24 * time.Hour)
	}

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
	taskIdx = 0
	prevEnd = time.Time{}

	for i, sec := range g.GanttSections {
		var tasks []GanttTaskLayout
		secStartY := curY

		for _, task := range sec.Tasks {
			var start, end time.Time
			if task.ID != "" {
				if rt, ok := resolved[task.ID]; ok {
					start, end = rt.Start, rt.End
				}
			}
			if start.IsZero() {
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
				if start.IsZero() && !prevEnd.IsZero() {
					start = prevEnd
				}
				if start.IsZero() {
					start = minDate
				}
				dur := parseMermaidDuration(task.EndStr)
				if dur > 0 {
					days := int(dur.Hours() / 24)
					if days > 0 {
						end = addWorkingDays(start, days, g.GanttExcludes, goLayout)
					} else {
						end = start.Add(dur)
					}
				} else if t, err := time.Parse(goLayout, task.EndStr); err == nil {
					end = t
				} else {
					end = start.Add(24 * time.Hour)
				}
			}

			x := dateToX(start)
			w := dateToX(end) - x
			if w < 1 {
				w = 1
			}

			hasTags := func(tag string) bool {
				for _, t := range task.Tags {
					if t == tag {
						return true
					}
				}
				return false
			}

			tasks = append(tasks, GanttTaskLayout{
				ID:          task.ID,
				Label:       task.Label,
				X:           x,
				Y:           curY,
				Width:       w,
				Height:      barH,
				IsCrit:      hasTags("crit"),
				IsDone:      hasTags("done"),
				IsActive:    hasTags("active"),
				IsMilestone: hasTags("milestone"),
			})

			prevEnd = end
			curY += barH + barGap
			taskIdx++
		}

		secH := curY - secStartY
		colorIdx := i % len(th.GanttSectionColors)
		sections = append(sections, GanttSectionLayout{
			Title:  sec.Title,
			Y:      secStartY,
			Height: secH,
			Color:  th.GanttSectionColors[colorIdx],
			Tasks:  tasks,
		})
	}

	// Axis ticks (one per week or per day depending on range).
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
			TodayMarkerX:   todayX,
			ShowTodayMarker: showToday,
			ChartX:          chartX,
			ChartY:          chartY,
			ChartWidth:      chartW,
			ChartHeight:     curY - chartY,
		},
	}
}
```

Add to `layout/layout.go` ComputeLayout switch:

```go
	case ir.Gantt:
		return computeGanttLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./layout/ -run TestGanttLayout -v`
Expected: PASS

**Step 5: Commit**

```bash
git add layout/gantt.go layout/gantt_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add Gantt chart layout with date arithmetic"
```

---

### Task 10: Layout — GitGraph

**Files:**
- Create: `layout/gitgraph.go`
- Create: `layout/gitgraph_test.go`
- Modify: `layout/types.go` (add GitGraphData types)
- Modify: `layout/layout.go` (add `case ir.GitGraph`)

**Step 1: Write the test**

```go
// layout/gitgraph_test.go
package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestGitGraphLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"
	g.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1"},
		&ir.GitBranch{Name: "develop"},
		&ir.GitCheckout{Branch: "develop"},
		&ir.GitCommit{ID: "c2"},
		&ir.GitCheckout{Branch: "main"},
		&ir.GitMerge{Branch: "develop", ID: "m1"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.GitGraph {
		t.Errorf("Kind = %v, want GitGraph", l.Kind)
	}
	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("dimensions = %f x %f", l.Width, l.Height)
	}

	ggd, ok := l.Diagram.(GitGraphData)
	if !ok {
		t.Fatalf("Diagram type = %T, want GitGraphData", l.Diagram)
	}
	if len(ggd.Commits) < 3 {
		t.Errorf("Commits = %d, want >= 3", len(ggd.Commits))
	}
	if len(ggd.Branches) < 2 {
		t.Errorf("Branches = %d, want >= 2", len(ggd.Branches))
	}
}

func TestGitGraphLayoutBranchLanes(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"
	g.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1"},
		&ir.GitBranch{Name: "dev"},
		&ir.GitCheckout{Branch: "dev"},
		&ir.GitCommit{ID: "c2"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)
	ggd := l.Diagram.(GitGraphData)

	// main and dev should be on different Y lanes.
	var mainY, devY float32
	for _, br := range ggd.Branches {
		if br.Name == "main" {
			mainY = br.Y
		}
		if br.Name == "dev" {
			devY = br.Y
		}
	}
	if mainY == devY {
		t.Errorf("main.Y=%f == dev.Y=%f, want different lanes", mainY, devY)
	}
}

func TestGitGraphLayoutCommitOrder(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"
	g.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1"},
		&ir.GitCommit{ID: "c2"},
		&ir.GitCommit{ID: "c3"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)
	ggd := l.Diagram.(GitGraphData)

	// Commits should be left-to-right.
	if len(ggd.Commits) != 3 {
		t.Fatalf("Commits = %d, want 3", len(ggd.Commits))
	}
	if ggd.Commits[0].X >= ggd.Commits[1].X || ggd.Commits[1].X >= ggd.Commits[2].X {
		t.Errorf("commits not left-to-right: %f, %f, %f",
			ggd.Commits[0].X, ggd.Commits[1].X, ggd.Commits[2].X)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./layout/ -run TestGitGraphLayout -v`
Expected: FAIL — `GitGraphData` undefined

**Step 3: Write minimal implementation**

Add to `layout/types.go`:

```go
// GitGraphData holds GitGraph-diagram-specific layout data.
type GitGraphData struct {
	Commits     []GitGraphCommitLayout
	Branches    []GitGraphBranchLayout
	Connections []GitGraphConnection
}

func (GitGraphData) diagramData() {}

// GitGraphCommitLayout holds positioned commit data.
type GitGraphCommitLayout struct {
	ID     string
	Tag    string
	Type   ir.GitCommitType
	Branch string
	X, Y   float32
}

// GitGraphBranchLayout holds branch lane data.
type GitGraphBranchLayout struct {
	Name   string
	Y      float32
	Color  string
	StartX float32
	EndX   float32
}

// GitGraphConnection holds a line connecting two commits (merge/cherry-pick).
type GitGraphConnection struct {
	FromX, FromY float32
	ToX, ToY     float32
	IsCherryPick bool
}
```

Create `layout/gitgraph.go`:

```go
package layout

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func computeGitGraphLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	padX := cfg.GitGraph.PaddingX
	padY := cfg.GitGraph.PaddingY
	commitSpacing := cfg.GitGraph.CommitSpacing
	branchSpacing := cfg.GitGraph.BranchSpacing

	mainBranch := g.GitMainBranch
	if mainBranch == "" {
		mainBranch = "main"
	}

	// Simulate git operations to build commit graph.
	type branchInfo struct {
		name  string
		order int
		head  string // latest commit ID on this branch
	}

	branches := map[string]*branchInfo{
		mainBranch: {name: mainBranch, order: 0},
	}
	currentBranch := mainBranch

	type commitInfo struct {
		id     string
		tag    string
		ctype  ir.GitCommitType
		branch string
		seq    int // sequential order
	}

	var commits []commitInfo
	commitMap := make(map[string]int) // commit ID -> index in commits
	var connections []GitGraphConnection
	autoID := 0

	for _, action := range g.GitActions {
		switch a := action.(type) {
		case *ir.GitCommit:
			id := a.ID
			if id == "" {
				id = fmt.Sprintf("auto_%d", autoID)
				autoID++
			}
			ci := commitInfo{
				id:     id,
				tag:    a.Tag,
				ctype:  a.Type,
				branch: currentBranch,
				seq:    len(commits),
			}
			commitMap[id] = len(commits)
			commits = append(commits, ci)
			branches[currentBranch].head = id

		case *ir.GitBranch:
			order := a.Order
			if order < 0 {
				order = len(branches)
			}
			branches[a.Name] = &branchInfo{
				name:  a.Name,
				order: order,
				head:  branches[currentBranch].head,
			}
			currentBranch = a.Name

		case *ir.GitCheckout:
			currentBranch = a.Branch

		case *ir.GitMerge:
			id := a.ID
			if id == "" {
				id = fmt.Sprintf("merge_%d", autoID)
				autoID++
			}
			ci := commitInfo{
				id:     id,
				tag:    a.Tag,
				ctype:  a.Type,
				branch: currentBranch,
				seq:    len(commits),
			}
			commitMap[id] = len(commits)
			commits = append(commits, ci)

			// Connection from merged branch head to this merge commit.
			if srcBranch, ok := branches[a.Branch]; ok && srcBranch.head != "" {
				if srcIdx, ok2 := commitMap[srcBranch.head]; ok2 {
					connections = append(connections, GitGraphConnection{
						FromX: float32(srcIdx), // placeholder, resolved below
						FromY: 0,
						ToX:   float32(len(commits) - 1),
						ToY:   0,
					})
				}
			}
			branches[currentBranch].head = id

		case *ir.GitCherryPick:
			id := fmt.Sprintf("cp_%d", autoID)
			autoID++
			ci := commitInfo{
				id:     id,
				tag:    a.ID, // show source as tag
				ctype:  ir.GitCommitNormal,
				branch: currentBranch,
				seq:    len(commits),
			}
			commitMap[id] = len(commits)
			commits = append(commits, ci)

			if srcIdx, ok := commitMap[a.ID]; ok {
				connections = append(connections, GitGraphConnection{
					FromX:        float32(srcIdx),
					ToX:          float32(len(commits) - 1),
					IsCherryPick: true,
				})
			}
			branches[currentBranch].head = id
		}
	}

	// Sort branches by order for lane assignment.
	type branchLane struct {
		name  string
		order int
	}
	var sortedBranches []branchLane
	for name, bi := range branches {
		sortedBranches = append(sortedBranches, branchLane{name, bi.order})
	}
	sort.Slice(sortedBranches, func(i, j int) bool {
		return sortedBranches[i].order < sortedBranches[j].order
	})

	branchY := make(map[string]float32)
	for i, bl := range sortedBranches {
		branchY[bl.name] = padY + float32(i)*branchSpacing
	}

	// Position commits.
	var commitLayouts []GitGraphCommitLayout
	for _, ci := range commits {
		x := padX + float32(ci.seq)*commitSpacing
		y := branchY[ci.branch]
		commitLayouts = append(commitLayouts, GitGraphCommitLayout{
			ID:     ci.id,
			Tag:    ci.tag,
			Type:   ci.ctype,
			Branch: ci.branch,
			X:      x,
			Y:      y,
		})
	}

	// Resolve connection pixel positions.
	var connLayouts []GitGraphConnection
	for _, conn := range connections {
		fromIdx := int(conn.FromX)
		toIdx := int(conn.ToX)
		if fromIdx < len(commitLayouts) && toIdx < len(commitLayouts) {
			connLayouts = append(connLayouts, GitGraphConnection{
				FromX:        commitLayouts[fromIdx].X,
				FromY:        commitLayouts[fromIdx].Y,
				ToX:          commitLayouts[toIdx].X,
				ToY:          commitLayouts[toIdx].Y,
				IsCherryPick: conn.IsCherryPick,
			})
		}
	}

	// Build branch layouts.
	var branchLayouts []GitGraphBranchLayout
	for i, bl := range sortedBranches {
		colorIdx := i % len(th.GitBranchColors)
		// Find start and end X for this branch.
		var startX, endX float32
		first := true
		for _, cl := range commitLayouts {
			if cl.Branch == bl.name {
				if first || cl.X < startX {
					startX = cl.X
				}
				if first || cl.X > endX {
					endX = cl.X
				}
				first = false
			}
		}
		branchLayouts = append(branchLayouts, GitGraphBranchLayout{
			Name:   bl.name,
			Y:      branchY[bl.name],
			Color:  th.GitBranchColors[colorIdx],
			StartX: startX,
			EndX:   endX,
		})
	}

	totalW := padX*2 + float32(len(commits))*commitSpacing
	totalH := padY*2 + float32(len(sortedBranches))*branchSpacing

	return &Layout{
		Kind:   g.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: GitGraphData{
			Commits:     commitLayouts,
			Branches:    branchLayouts,
			Connections: connLayouts,
		},
	}
}
```

Add to `layout/layout.go` ComputeLayout switch:

```go
	case ir.GitGraph:
		return computeGitGraphLayout(g, th, cfg)
```

**Step 4: Run test to verify it passes**

Run: `go test ./layout/ -run TestGitGraphLayout -v`
Expected: PASS

**Step 5: Commit**

```bash
git add layout/gitgraph.go layout/gitgraph_test.go layout/types.go layout/layout.go
git commit -m "feat(layout): add GitGraph swim-lane layout"
```

---

### Task 11: Renderer — Timeline

**Files:**
- Create: `render/timeline.go`
- Create: `render/timeline_test.go`
- Modify: `render/svg.go` (add `case layout.TimelineData`)

**Step 1: Write the test**

```go
// render/timeline_test.go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderTimeline(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Timeline
	g.TimelineTitle = "History"
	g.TimelineSections = []*ir.TimelineSection{
		{
			Title: "Early",
			Periods: []*ir.TimelinePeriod{
				{Title: "2002", Events: []*ir.TimelineEvent{{Text: "LinkedIn"}}},
				{Title: "2004", Events: []*ir.TimelineEvent{{Text: "Facebook"}}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "History") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "LinkedIn") {
		t.Error("missing event text")
	}
	if !strings.Contains(svg, "2002") {
		t.Error("missing period label")
	}
}
```

**Step 2-5: Implement, test, commit**

Create `render/timeline.go` with `renderTimeline()` that draws section bands, period labels, and event boxes. Add `case layout.TimelineData: renderTimeline(&b, l, th, cfg)` to svg.go. Use comma-ok assertion. Commit as `feat(render): add Timeline diagram SVG renderer`.

---

### Task 12: Renderer — Gantt

**Files:**
- Create: `render/gantt.go`
- Create: `render/gantt_test.go`
- Modify: `render/svg.go` (add `case layout.GanttData`)

**Step 1: Write the test**

```go
// render/gantt_test.go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderGantt(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Gantt
	g.GanttTitle = "Project"
	g.GanttDateFormat = "YYYY-MM-DD"
	g.GanttSections = []*ir.GanttSection{
		{
			Title: "Dev",
			Tasks: []*ir.GanttTask{
				{ID: "t1", Label: "Design", StartStr: "2024-01-01", EndStr: "10d"},
				{ID: "t2", Label: "Code", StartStr: "2024-01-11", EndStr: "20d", Tags: []string{"crit"}},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "Project") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Design") {
		t.Error("missing task label")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing task bars")
	}
}
```

**Step 2-5: Implement, test, commit**

Create `render/gantt.go` with `renderGantt()` that draws section backgrounds, task bars (colored by status), axis ticks, grid lines, today marker, milestone diamonds, and task labels. Add `case layout.GanttData`. Use comma-ok assertion. Commit as `feat(render): add Gantt chart SVG renderer`.

---

### Task 13: Renderer — GitGraph

**Files:**
- Create: `render/gitgraph.go`
- Create: `render/gitgraph_test.go`
- Modify: `render/svg.go` (add `case layout.GitGraphData`)

**Step 1: Write the test**

```go
// render/gitgraph_test.go
package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/layout"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestRenderGitGraph(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"
	g.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1", Tag: "v1.0"},
		&ir.GitBranch{Name: "develop"},
		&ir.GitCheckout{Branch: "develop"},
		&ir.GitCommit{ID: "c2"},
		&ir.GitCheckout{Branch: "main"},
		&ir.GitMerge{Branch: "develop"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "v1.0") {
		t.Error("missing tag label")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing commit circles")
	}
	if !strings.Contains(svg, "main") {
		t.Error("missing branch label")
	}
}
```

**Step 2-5: Implement, test, commit**

Create `render/gitgraph.go` with `renderGitGraph()` that draws branch lines, commit circles (normal/reverse/highlight), branch labels, tag labels, and merge/cherry-pick connections. Add `case layout.GitGraphData`. Use comma-ok assertion. Commit as `feat(render): add GitGraph diagram SVG renderer`.

---

### Task 14: Integration tests and fixtures

**Files:**
- Create: `testdata/fixtures/timeline-basic.mmd`
- Create: `testdata/fixtures/timeline-sections.mmd`
- Create: `testdata/fixtures/gantt-basic.mmd`
- Create: `testdata/fixtures/gantt-dependencies.mmd`
- Create: `testdata/fixtures/gitgraph-basic.mmd`
- Create: `testdata/fixtures/gitgraph-branches.mmd`
- Modify: `mermaid_test.go` (add integration tests)

**Step 1: Create fixtures and tests**

Fixture `timeline-basic.mmd`:
```
timeline
    title History of Social Media
    2002 : LinkedIn
    2004 : Facebook : Google
    2005 : YouTube
    2006 : Twitter
```

Fixture `gantt-basic.mmd`:
```
gantt
    title A Gantt Diagram
    dateFormat YYYY-MM-DD
    section Design
        Research :d1, 2024-01-01, 10d
        Mockups  :d2, after d1, 5d
    section Development
        Backend  :dev1, after d2, 20d
        Frontend :dev2, after d2, 15d
```

Fixture `gitgraph-basic.mmd`:
```
gitGraph
    commit id: "initial"
    commit id: "feature-start"
    branch develop
    checkout develop
    commit id: "dev-work"
    commit id: "dev-done" tag: "v0.1"
    checkout main
    merge develop tag: "v1.0"
    commit id: "hotfix" type: HIGHLIGHT
```

Add integration tests to `mermaid_test.go` using `readFixture()` pattern.

**Step 2: Run full test suite**

Run: `go test ./... -v`
Expected: ALL PASS

**Step 3: Commit**

```bash
git add testdata/fixtures/ mermaid_test.go
git commit -m "test: add integration tests and fixtures for Timeline, Gantt, and GitGraph"
```

---

### Task 15: Final validation

**Step 1:** Run: `go test ./... -v` — ALL PASS
**Step 2:** Run: `go vet ./... && gofmt -l .` — Clean
**Step 3:** Run: `go build ./...` — Clean
