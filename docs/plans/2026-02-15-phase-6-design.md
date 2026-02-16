# Phase 6: Timeline, Gantt & GitGraph — Design

**Date:** 2026-02-15
**Status:** Approved

## Goal

Add Timeline, Gantt chart, and GitGraph diagram support to gomd2svg. Timeline uses horizontal period-based layout, Gantt uses date-mapped horizontal bars with full date arithmetic, and GitGraph uses swim-lane branch layout.

## IR Types

### Timeline (`ir/timeline.go`)

- `TimelineEvent` struct: `Text string`
- `TimelinePeriod` struct: `Title string`, `Events []*TimelineEvent`
- `TimelineSection` struct: `Title string`, `Periods []*TimelinePeriod`
- Graph fields: `TimelineSections []*TimelineSection`, `TimelineTitle string`

### Gantt (`ir/gantt.go`)

- `GanttTask` struct: `ID string`, `Label string`, `StartStr string`, `EndStr string`, `Tags []string` (done/active/crit/milestone), `AfterIDs []string`, `UntilID string`
- `GanttSection` struct: `Title string`, `Tasks []*GanttTask`
- Graph fields: `GanttSections []*GanttSection`, `GanttTitle string`, `GanttDateFormat string`, `GanttAxisFormat string`, `GanttExcludes []string`, `GanttTickInterval string`, `GanttTodayMarker string`, `GanttWeekday string`

### GitGraph (`ir/gitgraph.go`)

- `GitAction` interface with `gitAction()` marker (sealed)
- Concrete types: `GitCommit` (ID, Tag, Type), `GitBranch` (Name, Order), `GitCheckout` (Branch), `GitMerge` (Branch, ID, Tag, Type), `GitCherryPick` (ID, Parent)
- `GitCommitType` enum: `GitCommitNormal`, `GitCommitReverse`, `GitCommitHighlight`
- Graph fields: `GitActions []GitAction`, `GitMainBranch string`

## Parsers

### Timeline

Line-by-line parsing:
- First content line: `timeline` keyword
- `title <text>` sets TimelineTitle
- `section <text>` starts a new TimelineSection
- Period lines: `<period> : <event>` or `<period> : <event> : <event>` for multiple events
- Continuation events: `: <event>` (indented, no period) appends to the current period
- Periods without a section go into an implicit default section

### Gantt

Line-by-line parsing:
- First content line: `gantt` keyword
- Directives: `dateFormat <fmt>`, `axisFormat <fmt>`, `excludes <list>`, `tickInterval <interval>`, `todayMarker <style|off>`, `weekend <day>`, `title <text>`
- `section <text>` starts a new GanttSection
- Task lines: `<taskName> : <metadata>` — metadata is comma-separated
- Metadata parsing order: identify tags first (done/active/crit/milestone), then remaining items are ID, start, and end/duration by pattern matching
- Start patterns: ISO date, `after <id>` (space-separated list), or implicit (follows previous task)
- End patterns: duration (`5d`, `2w`, `8h`, `30m`), ISO date
- Default dateFormat: `YYYY-MM-DD`

### GitGraph

Line-by-line parsing:
- First content line: `gitGraph` with optional direction suffix (`LR:`, `TB:`, `BT:`)
- `commit` with optional `id: "str"`, `tag: "str"`, `type: NORMAL|REVERSE|HIGHLIGHT`
- `branch <name>` with optional `order: <int>` — quoted names for keyword conflicts
- `checkout <name>` or `switch <name>`
- `merge <name>` with optional `id:`, `tag:`, `type:`
- `cherry-pick id: "str"` with optional `parent: "str"`

## Config

### TimelineConfig

| Field | Type | Default | Purpose |
|-------|------|---------|---------|
| PeriodWidth | float32 | 150 | Width of each time period column |
| EventHeight | float32 | 30 | Height per event box |
| SectionPadding | float32 | 10 | Padding around section bands |
| PaddingX | float32 | 20 | Horizontal canvas padding |
| PaddingY | float32 | 20 | Vertical canvas padding |

### GanttConfig

| Field | Type | Default | Purpose |
|-------|------|---------|---------|
| BarHeight | float32 | 20 | Height of task bars |
| BarGap | float32 | 4 | Gap between task bars |
| TopPadding | float32 | 50 | Top margin |
| SidePadding | float32 | 75 | Side margin for section labels |
| GridLineStartPadding | float32 | 35 | Grid vertical offset |
| FontSize | float32 | 11 | Task label font size |
| SectionFontSize | float32 | 11 | Section label font size |
| NumberSectionStyles | int | 4 | Alternating section color count |

### GitGraphConfig

| Field | Type | Default | Purpose |
|-------|------|---------|---------|
| CommitRadius | float32 | 8 | Commit circle radius |
| CommitSpacing | float32 | 60 | Horizontal space between commits |
| BranchSpacing | float32 | 40 | Vertical space between branch lanes |
| PaddingX | float32 | 30 | Horizontal canvas padding |
| PaddingY | float32 | 30 | Vertical canvas padding |
| TagFontSize | float32 | 11 | Font size for tag labels |

## Layout

### Timeline (`layout/timeline.go`)

Horizontal layout. Sections are color-banded rows. Periods are evenly-spaced columns. Events stack vertically within each period cell. No Sugiyama — pure geometric positioning.

Layout types: `TimelineData` (implements DiagramData), `TimelineSectionLayout`, `TimelinePeriodLayout`, `TimelineEventLayout`.

### Gantt (`layout/gantt.go`)

Date-based horizontal layout. X-axis maps date range to pixel positions. Y-axis stacks sections and tasks vertically.

Key logic:
- `time.Parse` with configurable `dateFormat` (map mermaid tokens YYYY/MM/DD to Go layout)
- Duration parsing: `5d`, `2w`, `8h`, `30m` -> `time.Duration`
- `after` dependency resolution: topological walk to find latest predecessor end date
- `excludes` filtering: skip weekends, named days, specific dates when computing durations
- Milestone: rendered as diamond at start_date + duration/2
- Today marker: vertical line at current date position

Layout types: `GanttData` (implements DiagramData), `GanttSectionLayout`, `GanttTaskLayout`, `GanttAxisTick`.

### GitGraph (`layout/gitgraph.go`)

Swim-lane layout. Each branch gets a horizontal lane. Commits are evenly spaced along X-axis. Merge and cherry-pick lines connect across lanes.

Key logic:
- Simulate git operations: track current branch, HEAD per branch, commit graph
- Assign each branch a lane (Y position) respecting `order` attribute
- Position commits sequentially on X-axis
- Compute merge/cherry-pick connection points between lanes

Layout types: `GitGraphData` (implements DiagramData), `GitGraphCommitLayout`, `GitGraphBranchLayout`, `GitGraphConnection`.

## Rendering

### Timeline (`render/timeline.go`)

- Section background bands with alternating colors
- Period labels on X-axis (bottom of each column)
- Event boxes stacked vertically per period with rounded rects
- Section labels on the left side
- Title centered above

### Gantt (`render/gantt.go`)

- X-axis date labels using `axisFormat` strftime conversion
- Section labels on Y-axis with alternating background bands
- Task bars colored by status: crit=red stroke, done=gray fill, active=blue fill, default=theme primary
- Milestone diamonds
- Optional today marker (vertical dashed line)
- Grid lines at tick intervals
- Task labels inside or beside bars

### GitGraph (`render/gitgraph.go`)

- Branch lanes as colored horizontal lines
- Commit circles: normal=filled, reverse=crossed, highlight=filled rect
- Branch labels at lane start
- Tag labels above/below commits
- Merge arrows: solid lines connecting lanes
- Cherry-pick arrows: dashed lines connecting lanes

## Theme Fields

### Timeline
- `TimelineSectionColors []string` — alternating section band colors (4+ colors)
- `TimelineEventFill string`, `TimelineEventBorder string`

### Gantt
- `GanttTaskFill string`, `GanttTaskBorder string`
- `GanttCritFill string`, `GanttCritBorder string`
- `GanttDoneFill string`, `GanttActiveFill string`
- `GanttMilestoneFill string`
- `GanttGridColor string`
- `GanttTodayMarkerColor string`
- `GanttSectionColors []string` — alternating section backgrounds

### GitGraph
- `GitBranchColors []string` — per-branch colors (8+ colors)
- `GitCommitFill string`, `GitCommitStroke string`
- `GitTagFill string`, `GitTagBorder string`
- `GitHighlightFill string`

## Deferred

- Gantt `displayMode: compact` (stacks non-overlapping tasks on same row)
- Gantt `topAxis` (date labels above chart)
- Gantt `vert` vertical reference lines
- GitGraph `BT:` direction (bottom-to-top)
- GitGraph `TB:` direction (top-to-bottom) — LR only for now
- GitGraph `parallelCommits` configuration
- GitGraph `rotateCommitLabel` configuration
- Per-task/per-commit classDef styling — Phase 12 Theme System
