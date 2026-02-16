package ir

// C4Kind identifies the specific C4 diagram subtype.
type C4Kind int

const (
	C4Context C4Kind = iota
	C4Container
	C4Component
	C4Dynamic
	C4Deployment
)

func (k C4Kind) String() string {
	switch k {
	case C4Container:
		return "C4Container"
	case C4Component:
		return "C4Component"
	case C4Dynamic:
		return "C4Dynamic"
	case C4Deployment:
		return "C4Deployment"
	default:
		return "C4Context"
	}
}

// C4ElementType identifies the kind of C4 element.
type C4ElementType int

const (
	C4Person C4ElementType = iota
	C4System
	C4SystemDb
	C4SystemQueue
	C4ContainerPlain
	C4ContainerDb
	C4ContainerQueue
	C4ComponentPlain
	C4ExternalPerson
	C4ExternalSystem
	C4ExternalSystemDb
	C4ExternalSystemQueue
	C4ExternalContainer
	C4ExternalContainerDb
	C4ExternalContainerQueue
	C4ExternalComponent
)

func (e C4ElementType) String() string {
	switch e {
	case C4Person:
		return "Person"
	case C4System:
		return "System"
	case C4SystemDb:
		return "SystemDb"
	case C4SystemQueue:
		return "SystemQueue"
	case C4ContainerPlain:
		return "Container"
	case C4ContainerDb:
		return "ContainerDb"
	case C4ContainerQueue:
		return "ContainerQueue"
	case C4ComponentPlain:
		return "Component"
	case C4ExternalPerson:
		return "Person_Ext"
	case C4ExternalSystem:
		return "System_Ext"
	case C4ExternalSystemDb:
		return "SystemDb_Ext"
	case C4ExternalSystemQueue:
		return "SystemQueue_Ext"
	case C4ExternalContainer:
		return "Container_Ext"
	case C4ExternalContainerDb:
		return "ContainerDb_Ext"
	case C4ExternalContainerQueue:
		return "ContainerQueue_Ext"
	case C4ExternalComponent:
		return "Component_Ext"
	default:
		return "System"
	}
}

// IsExternal returns true if the element type represents an external entity.
func (e C4ElementType) IsExternal() bool {
	switch e {
	case C4ExternalPerson, C4ExternalSystem, C4ExternalSystemDb, C4ExternalSystemQueue,
		C4ExternalContainer, C4ExternalContainerDb, C4ExternalContainerQueue, C4ExternalComponent:
		return true
	default:
		return false
	}
}

// IsPerson returns true if the element type is a person (internal or external).
func (e C4ElementType) IsPerson() bool {
	return e == C4Person || e == C4ExternalPerson
}

// IsDatabase returns true if the element type is a database variant.
func (e C4ElementType) IsDatabase() bool {
	switch e {
	case C4SystemDb, C4ContainerDb, C4ExternalSystemDb, C4ExternalContainerDb:
		return true
	default:
		return false
	}
}

// IsQueue returns true if the element type is a queue variant.
func (e C4ElementType) IsQueue() bool {
	switch e {
	case C4SystemQueue, C4ContainerQueue, C4ExternalSystemQueue, C4ExternalContainerQueue:
		return true
	default:
		return false
	}
}

// C4Element represents a single element in a C4 diagram.
type C4Element struct {
	ID          string
	Label       string
	Technology  string
	Description string
	Type        C4ElementType
	BoundaryID  string // empty if top-level
}

// C4Boundary represents a boundary/grouping in a C4 diagram.
type C4Boundary struct {
	ID       string
	Label    string
	Type     string   // e.g. "Enterprise", "Software System"
	Children []string // element IDs within this boundary
}

// C4Rel represents a relationship between C4 elements.
type C4Rel struct {
	From        string
	To          string
	Label       string
	Technology  string
	Description string
}
