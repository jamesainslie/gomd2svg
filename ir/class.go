package ir

// Visibility represents the access modifier of a class member.
type Visibility int

const (
	VisNone      Visibility = iota
	VisPublic               // +
	VisPrivate              // -
	VisProtected            // #
	VisPackage              // ~
)

// Symbol returns the Mermaid syntax character for this visibility.
func (v Visibility) Symbol() string {
	switch v {
	case VisPublic:
		return "+"
	case VisPrivate:
		return "-"
	case VisProtected:
		return "#"
	case VisPackage:
		return "~"
	default:
		return ""
	}
}

// MemberClassifier distinguishes abstract and static members.
type MemberClassifier int

const (
	ClassifierNone     MemberClassifier = iota
	ClassifierAbstract                  // *
	ClassifierStatic                    // $
)

// ClassMember represents an attribute or method of a class.
type ClassMember struct {
	Name       string
	Type       string // attribute type or return type
	Params     string // raw parameter string (empty for attributes)
	IsMethod   bool
	Visibility Visibility
	Classifier MemberClassifier
	Generic    string // generic type parameter e.g. "T"
}

// ClassMembers groups a class's attributes and methods.
type ClassMembers struct {
	Attributes []ClassMember
	Methods    []ClassMember
}

// Namespace groups classes under a named scope.
type Namespace struct {
	Name    string
	Classes []string // node IDs
}

// DiagramNote represents a note attached to a class diagram.
type DiagramNote struct {
	Text     string
	Position string // "right of", "left of", or "" for floating
	Target   string // node ID, empty for floating notes
}
