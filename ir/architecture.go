package ir

// ArchSide represents a connection side on a service or junction.
type ArchSide int

const (
	ArchLeft ArchSide = iota
	ArchRight
	ArchTop
	ArchBottom
)

func (s ArchSide) String() string {
	switch s {
	case ArchRight:
		return "R"
	case ArchTop:
		return "T"
	case ArchBottom:
		return "B"
	default:
		return "L"
	}
}

// ArchService represents a service node in an architecture diagram.
type ArchService struct {
	ID      string
	Label   string
	Icon    string
	GroupID string // empty if top-level
}

// ArchGroup represents a grouping container.
type ArchGroup struct {
	ID       string
	Label    string
	Icon     string
	ParentID string   // for nested groups
	Children []string // service/junction IDs
}

// ArchJunction is a connection point between edges.
type ArchJunction struct {
	ID      string
	GroupID string
}

// ArchEdge represents a connection between services/junctions.
type ArchEdge struct {
	FromID     string
	FromSide   ArchSide
	ToID       string
	ToSide     ArchSide
	ArrowLeft  bool
	ArrowRight bool
}
