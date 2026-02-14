// Package layout computes the geometric layout of Mermaid diagrams.
// It takes an intermediate representation (ir.Graph) and produces
// positioned nodes, routed edges, and a bounding box.
package layout

import "github.com/jamesainslie/mermaid-go/ir"

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
