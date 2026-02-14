package ir

type DiagramKind int

const (
	Flowchart DiagramKind = iota
	Class
	State
	Sequence
	Er
	Pie
	Mindmap
	Journey
	Timeline
	Gantt
	Requirement
	GitGraph
	C4
	Sankey
	Quadrant
	ZenUML
	Block
	Packet
	Kanban
	Architecture
	Radar
	Treemap
	XYChart
)

type Direction int

const (
	TopDown Direction = iota
	LeftRight
	BottomTop
	RightLeft
)

func DirectionFromToken(token string) (Direction, bool) {
	switch token {
	case "TD", "TB":
		return TopDown, true
	case "LR":
		return LeftRight, true
	case "RL":
		return RightLeft, true
	case "BT":
		return BottomTop, true
	default:
		return TopDown, false
	}
}
