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
	From        string
	To          string
	Label       *TextBlock
	Points      [][2]float32
	LabelAnchor [2]float32
	Style       ir.EdgeStyle
	ArrowStart  bool
	ArrowEnd    bool
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
