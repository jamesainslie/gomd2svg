package ir

// BlockDef represents a block in a block diagram.
type BlockDef struct {
	ID       string
	Label    string
	Shape    NodeShape
	Width    int         // column span (1 = default)
	Children []*BlockDef // nested blocks
}
