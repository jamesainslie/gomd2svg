// Package layout computes the geometric layout of Mermaid diagrams.
// It takes an intermediate representation (ir.Graph) and produces
// positioned nodes, routed edges, and a bounding box.
package layout

import (
	"github.com/jamesainslie/mermaid-go/ir"
)

// Layout holds the fully computed positions and dimensions for a diagram.
type Layout struct {
	Kind      ir.DiagramKind
	Nodes     map[string]*NodeLayout
	Edges     []*EdgeLayout
	Subgraphs []*SubgraphLayout
	Width     float32
	Height    float32
	Diagram   DiagramData
}

// NodeLayout holds the position, size, and style of a single node.
type NodeLayout struct {
	ID     string
	Label  TextBlock
	Shape  ir.NodeShape
	X, Y   float32
	Width  float32
	Height float32
	Style  ir.NodeStyle
}

// EdgeLayout holds the route, label, and style of a single edge.
type EdgeLayout struct {
	From           string
	To             string
	Label          *TextBlock
	Points         [][2]float32
	LabelAnchor    [2]float32
	Style          ir.EdgeStyle
	ArrowStart     bool
	ArrowEnd       bool
	ArrowStartKind *ir.EdgeArrowhead
	ArrowEndKind   *ir.EdgeArrowhead
}

// SubgraphLayout holds the position and size of a subgraph container.
type SubgraphLayout struct {
	ID     string
	Label  string
	X, Y   float32
	Width  float32
	Height float32
}

// TextBlock holds measured text for rendering inside nodes or on edges.
type TextBlock struct {
	Lines    []string
	Width    float32
	Height   float32
	FontSize float32
}

// DiagramData is a sealed interface for diagram-specific layout data.
type DiagramData interface {
	diagramData()
}

// GraphData holds flowchart/graph-specific layout state.
type GraphData struct{}

func (GraphData) diagramData() {}

// ClassData holds class-diagram-specific layout data.
type ClassData struct {
	Compartments map[string]ClassCompartment
	Members      map[string]*ir.ClassMembers
	Annotations  map[string]string
}

func (ClassData) diagramData() {}

// ClassCompartment stores section dimensions for UML class boxes.
type ClassCompartment struct {
	HeaderHeight    float32
	AttributeHeight float32
	MethodHeight    float32
}

// ERData holds ER-diagram-specific layout data.
type ERData struct {
	EntityDims map[string]EntityDimensions
	Entities   map[string]*ir.Entity
}

func (ERData) diagramData() {}

// EntityDimensions stores column widths for entity rendering.
type EntityDimensions struct {
	TypeColWidth float32
	NameColWidth float32
	KeyColWidth  float32
	HeaderHeight float32
	RowCount     int
}

// StateData holds state-diagram-specific layout data.
type StateData struct {
	InnerLayouts    map[string]*Layout
	Descriptions    map[string]string
	Annotations     map[string]ir.StateAnnotation
	CompositeStates map[string]*ir.CompositeState
}

func (StateData) diagramData() {}

// SequenceData holds sequence-diagram-specific layout data.
type SequenceData struct {
	Participants  []SeqParticipantLayout
	Lifelines     []SeqLifeline
	Messages      []SeqMessageLayout
	Activations   []SeqActivationLayout
	Notes         []SeqNoteLayout
	Frames        []SeqFrameLayout
	Boxes         []SeqBoxLayout
	Autonumber    bool
	DiagramHeight float32
}

func (SequenceData) diagramData() {}

// SeqParticipantLayout holds the position of a participant header.
type SeqParticipantLayout struct {
	ID     string
	Label  TextBlock
	Kind   ir.SeqParticipantKind
	X      float32 // center X
	Y      float32 // top Y (0 for normal, mid-diagram for created)
	Width  float32
	Height float32
}

// SeqLifeline is a vertical dashed line from participant to diagram bottom.
type SeqLifeline struct {
	ParticipantID string
	X             float32
	TopY          float32
	BottomY       float32
}

// SeqMessageLayout holds the position of a message arrow.
type SeqMessageLayout struct {
	From   string
	To     string
	Text   TextBlock
	Kind   ir.SeqMessageKind
	Y      float32 // vertical position
	FromX  float32
	ToX    float32
	Number int // autonumber (0 if disabled)
}

// SeqActivationLayout holds the bounds of an activation bar.
type SeqActivationLayout struct {
	ParticipantID string
	X             float32
	TopY          float32
	BottomY       float32
	Width         float32
}

// SeqNoteLayout holds the position and content of a note.
type SeqNoteLayout struct {
	Text   TextBlock
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// SeqFrameLayout holds the bounds and label of a frame (combined fragment).
type SeqFrameLayout struct {
	Kind     ir.SeqFrameKind
	Label    string
	Color    string
	X        float32
	Y        float32
	Width    float32
	Height   float32
	Dividers []float32 // Y positions of else/and/option divider lines
}

// SeqBoxLayout holds the bounds and label of a participant box group.
type SeqBoxLayout struct {
	Label  string
	Color  string
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// KanbanData holds Kanban-diagram-specific layout data.
type KanbanData struct {
	Columns []KanbanColumnLayout
}

func (KanbanData) diagramData() {}

// KanbanColumnLayout holds the position of a Kanban column.
type KanbanColumnLayout struct {
	ID     string
	Label  TextBlock
	X, Y   float32
	Width  float32
	Height float32
	Cards  []KanbanCardLayout
}

// KanbanCardLayout holds the position of a single Kanban card.
type KanbanCardLayout struct {
	ID       string
	Label    TextBlock
	Priority ir.KanbanPriority
	X, Y     float32
	Width    float32
	Height   float32
	Metadata map[string]string
}

// PacketData holds Packet-diagram-specific layout data.
type PacketData struct {
	Rows       []PacketRowLayout
	BitsPerRow int
	ShowBits   bool
}

func (PacketData) diagramData() {}

// PacketRowLayout holds the position of a row of packet fields.
type PacketRowLayout struct {
	Y      float32
	Height float32
	Fields []PacketFieldLayout
}

// PacketFieldLayout holds the position of a single packet field cell.
type PacketFieldLayout struct {
	Label    TextBlock
	X, Y     float32
	Width    float32
	Height   float32
	StartBit int
	EndBit   int
}
