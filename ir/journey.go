package ir

// JourneyTask represents a single task in a journey diagram.
type JourneyTask struct {
	Name    string
	Score   int      // 1-5 satisfaction score
	Actors  []string // participating actors
	Section string   // section name this task belongs to
}

// JourneySection represents a named section grouping tasks.
type JourneySection struct {
	Name  string
	Tasks []int // indices into Graph.JourneyTasks
}
