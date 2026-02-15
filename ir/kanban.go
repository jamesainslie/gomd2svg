package ir

// KanbanPriority represents the priority level of a Kanban card.
type KanbanPriority int

const (
	PriorityNone KanbanPriority = iota
	PriorityVeryLow
	PriorityLow
	PriorityHigh
	PriorityVeryHigh
)

func (p KanbanPriority) String() string {
	switch p {
	case PriorityVeryLow:
		return "Very Low"
	case PriorityLow:
		return "Low"
	case PriorityHigh:
		return "High"
	case PriorityVeryHigh:
		return "Very High"
	default:
		return ""
	}
}

// KanbanCard represents a single card/task on a Kanban board.
type KanbanCard struct {
	ID          string
	Label       string
	Assigned    string
	Ticket      string
	Priority    KanbanPriority
	Icon        string
	Description string
}

// KanbanColumn represents a column on a Kanban board.
type KanbanColumn struct {
	ID    string
	Label string
	Cards []*KanbanCard
}
