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

// PieData holds pie-chart-specific layout data.
type PieData struct {
	Slices   []PieSliceLayout
	CenterX  float32
	CenterY  float32
	Radius   float32
	Title    string
	ShowData bool
}

func (PieData) diagramData() {}

// PieSliceLayout holds computed angles and label position for one slice.
type PieSliceLayout struct {
	Label      string
	Value      float64
	Percentage float32
	StartAngle float32
	EndAngle   float32
	LabelX     float32
	LabelY     float32
	ColorIndex int
}

// QuadrantData holds quadrant-chart-specific layout data.
type QuadrantData struct {
	Points      []QuadrantPointLayout
	ChartX      float32
	ChartY      float32
	ChartWidth  float32
	ChartHeight float32
	Title       string
	Labels      [4]string
	XAxisLeft   string
	XAxisRight  string
	YAxisBottom string
	YAxisTop    string
}

func (QuadrantData) diagramData() {}

// QuadrantPointLayout holds the pixel position of a data point.
type QuadrantPointLayout struct {
	Label string
	X     float32
	Y     float32
}

// TimelineData holds timeline-diagram-specific layout data.
type TimelineData struct {
	Sections []TimelineSectionLayout
	Title    string
}

func (TimelineData) diagramData() {}

// TimelineSectionLayout holds positioned section data.
type TimelineSectionLayout struct {
	Title   string
	X, Y    float32
	Width   float32
	Height  float32
	Color   string
	Periods []TimelinePeriodLayout
}

// TimelinePeriodLayout holds positioned period data.
type TimelinePeriodLayout struct {
	Title  string
	X, Y   float32
	Width  float32
	Height float32
	Events []TimelineEventLayout
}

// TimelineEventLayout holds positioned event data.
type TimelineEventLayout struct {
	Text   string
	X, Y   float32
	Width  float32
	Height float32
}

// GanttData holds Gantt-diagram-specific layout data.
type GanttData struct {
	Sections        []GanttSectionLayout
	Title           string
	AxisTicks       []GanttAxisTick
	TodayMarkerX    float32
	ShowTodayMarker bool
	ChartX          float32
	ChartY          float32
	ChartWidth      float32
	ChartHeight     float32
}

func (GanttData) diagramData() {}

// GanttSectionLayout holds positioned section data.
type GanttSectionLayout struct {
	Title  string
	Y      float32
	Height float32
	Color  string
	Tasks  []GanttTaskLayout
}

// GanttTaskLayout holds positioned task bar data.
type GanttTaskLayout struct {
	ID          string
	Label       string
	X, Y        float32
	Width       float32
	Height      float32
	IsCrit      bool
	IsDone      bool
	IsActive    bool
	IsMilestone bool
}

// GanttAxisTick holds a tick mark on the date axis.
type GanttAxisTick struct {
	Label string
	X     float32
}

// GitGraphData holds GitGraph-diagram-specific layout data.
type GitGraphData struct {
	Commits     []GitGraphCommitLayout
	Branches    []GitGraphBranchLayout
	Connections []GitGraphConnection
}

func (GitGraphData) diagramData() {}

// GitGraphCommitLayout holds positioned commit data.
type GitGraphCommitLayout struct {
	ID     string
	Tag    string
	Type   ir.GitCommitType
	Branch string
	X, Y   float32
}

// GitGraphBranchLayout holds branch lane data.
type GitGraphBranchLayout struct {
	Name   string
	Y      float32
	Color  string
	StartX float32
	EndX   float32
}

// GitGraphConnection holds a line connecting two commits (merge/cherry-pick).
type GitGraphConnection struct {
	FromX, FromY float32
	ToX, ToY     float32
	IsCherryPick bool
}

// XYChartData holds XY chart layout data.
type XYChartData struct {
	Series      []XYSeriesLayout
	XLabels     []XYAxisLabel
	YTicks      []XYAxisTick
	Title       string
	ChartX      float32
	ChartY      float32
	ChartWidth  float32
	ChartHeight float32
	YMin        float64
	YMax        float64
	Horizontal  bool
}

func (XYChartData) diagramData() {}

// XYSeriesLayout holds one positioned data series.
type XYSeriesLayout struct {
	Type       ir.XYSeriesType
	Points     []XYPointLayout
	ColorIndex int
}

// XYPointLayout holds the pixel position and value of one data point.
type XYPointLayout struct {
	X      float32
	Y      float32
	Width  float32 // bar width (0 for line points)
	Height float32 // bar height (0 for line points)
	Value  float64
}

// XYAxisLabel holds a label on the x-axis.
type XYAxisLabel struct {
	Text string
	X    float32
}

// XYAxisTick holds a tick mark on the y-axis.
type XYAxisTick struct {
	Label string
	Y     float32
}

// RadarData holds radar chart layout data.
type RadarData struct {
	Axes           []RadarAxisLayout
	Curves         []RadarCurveLayout
	GraticuleRadii []float32
	GraticuleType  ir.RadarGraticule
	CenterX        float32
	CenterY        float32
	Radius         float32
	Title          string
	ShowLegend     bool
	MaxValue       float64
	MinValue       float64
}

func (RadarData) diagramData() {}

// RadarAxisLayout holds the endpoint and label position of one axis.
type RadarAxisLayout struct {
	Label  string
	EndX   float32
	EndY   float32
	LabelX float32
	LabelY float32
}

// RadarCurveLayout holds one data series polygon.
type RadarCurveLayout struct {
	Label      string
	Points     [][2]float32
	ColorIndex int
}
