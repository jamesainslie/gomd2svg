package ir

// MindmapShape represents the visual shape of a mindmap node.
type MindmapShape int

const (
	MindmapShapeDefault MindmapShape = iota
	MindmapSquare
	MindmapRounded
	MindmapCircle
	MindmapBang
	MindmapCloud
	MindmapHexagon
)

func (s MindmapShape) String() string {
	switch s {
	case MindmapShapeDefault:
		return "default"
	case MindmapSquare:
		return "square"
	case MindmapRounded:
		return "rounded"
	case MindmapCircle:
		return "circle"
	case MindmapBang:
		return "bang"
	case MindmapCloud:
		return "cloud"
	case MindmapHexagon:
		return "hexagon"
	default:
		return "unknown"
	}
}

// MindmapNode represents a node in the mindmap tree.
type MindmapNode struct {
	ID       string
	Label    string
	Shape    MindmapShape
	Icon     string // CSS class from ::icon()
	Class    string // CSS class from :::
	Children []*MindmapNode
}
