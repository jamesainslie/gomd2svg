package ir

// StateAnnotation marks a state as a pseudo-state (choice, fork, or join).
type StateAnnotation int

const (
	StateChoice StateAnnotation = iota
	StateFork
	StateJoin
)

// String returns the lowercase name of the annotation.
func (a StateAnnotation) String() string {
	switch a {
	case StateChoice:
		return "choice"
	case StateFork:
		return "fork"
	case StateJoin:
		return "join"
	default:
		return ""
	}
}

// CompositeState represents a state that contains a nested state machine.
type CompositeState struct {
	ID        string
	Label     string
	Inner     *Graph   // primary nested state machine
	Regions   []*Graph // concurrent regions (separated by --)
	Direction *Direction
}
