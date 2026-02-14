package ir

type Node struct {
	ID    string
	Label string
	Shape NodeShape
	Value *float32
	Icon  *string
}

type Edge struct {
	From            string
	To              string
	Label           *string
	StartLabel      *string
	EndLabel        *string
	Directed        bool
	ArrowStart      bool
	ArrowEnd        bool
	ArrowStartKind  *EdgeArrowhead
	ArrowEndKind    *EdgeArrowhead
	StartDecoration *EdgeDecoration
	EndDecoration   *EdgeDecoration
	Style           EdgeStyle
}

type Subgraph struct {
	ID        *string
	Label     string
	Nodes     []string
	Direction *Direction
	Icon      *string
}

type Graph struct {
	Kind      DiagramKind
	Direction Direction
	Nodes     map[string]*Node
	NodeOrder map[string]int
	Edges     []*Edge
	Subgraphs []*Subgraph

	ClassDefs        map[string]*NodeStyle
	NodeClasses      map[string][]string
	NodeStyles       map[string]*NodeStyle
	SubgraphStyles   map[string]*NodeStyle
	SubgraphClasses  map[string][]string
	NodeLinks        map[string]*NodeLink
	EdgeStyles       map[int]*EdgeStyleOverride
	EdgeStyleDefault *EdgeStyleOverride

	// Class diagram fields
	Members     map[string]*ClassMembers
	Annotations map[string]string // node ID -> stereotype text
	Namespaces  []*Namespace
	Notes       []*DiagramNote

	// ER diagram fields
	Entities map[string]*Entity

	// State diagram fields
	CompositeStates   map[string]*CompositeState
	StateDescriptions map[string]string
	StateAnnotations  map[string]StateAnnotation
}

func NewGraph() *Graph {
	return &Graph{
		Kind:              Flowchart,
		Direction:         TopDown,
		Nodes:             make(map[string]*Node),
		NodeOrder:         make(map[string]int),
		ClassDefs:         make(map[string]*NodeStyle),
		NodeClasses:       make(map[string][]string),
		NodeStyles:        make(map[string]*NodeStyle),
		SubgraphStyles:    make(map[string]*NodeStyle),
		SubgraphClasses:   make(map[string][]string),
		NodeLinks:         make(map[string]*NodeLink),
		EdgeStyles:        make(map[int]*EdgeStyleOverride),
		Members:           make(map[string]*ClassMembers),
		Annotations:       make(map[string]string),
		Entities:          make(map[string]*Entity),
		CompositeStates:   make(map[string]*CompositeState),
		StateDescriptions: make(map[string]string),
		StateAnnotations:  make(map[string]StateAnnotation),
	}
}

func (g *Graph) EnsureNode(id string, label *string, shape *NodeShape) {
	n, exists := g.Nodes[id]
	if !exists {
		n = &Node{
			ID:    id,
			Label: id,
			Shape: Rectangle,
		}
		g.Nodes[id] = n
		g.NodeOrder[id] = len(g.NodeOrder)
	}
	if label != nil {
		n.Label = *label
	}
	if shape != nil {
		n.Shape = *shape
	}
}
