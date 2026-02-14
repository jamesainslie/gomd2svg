package ir

// AttributeKey represents the key type of an ER entity attribute.
type AttributeKey int

const (
	KeyNone    AttributeKey = iota
	KeyPrimary              // PK
	KeyForeign              // FK
	KeyUnique               // UK
)

// String returns the short label for the attribute key.
func (k AttributeKey) String() string {
	switch k {
	case KeyPrimary:
		return "PK"
	case KeyForeign:
		return "FK"
	case KeyUnique:
		return "UK"
	default:
		return ""
	}
}

// EntityAttribute describes a single attribute within an ER entity.
type EntityAttribute struct {
	Type    string
	Name    string
	Keys    []AttributeKey
	Comment string
}

// Entity represents an entity (table) in an ER diagram.
type Entity struct {
	ID         string
	Label      string // display name (alias); empty means use ID
	Attributes []EntityAttribute
}

// DisplayName returns the Label if non-empty, otherwise the ID.
func (e *Entity) DisplayName() string {
	if e.Label != "" {
		return e.Label
	}
	return e.ID
}
