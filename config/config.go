package config

// Layout holds all configuration for diagram layout computation.
type Layout struct {
	NodeSpacing          float32
	RankSpacing          float32
	LabelLineHeight      float32
	PreferredAspectRatio *float32
	Flowchart            FlowchartConfig
	Padding              PaddingConfig
	Class                ClassConfig
	State                StateConfig
	ER                   ERConfig
	Sequence             SequenceConfig
	Kanban               KanbanConfig
	Packet               PacketConfig
	Pie                  PieConfig
	Quadrant             QuadrantConfig
	Timeline             TimelineConfig
	Gantt                GanttConfig
	GitGraph             GitGraphConfig
	XYChart              XYChartConfig
	Radar                RadarConfig
	Mindmap              MindmapConfig
	Sankey               SankeyConfig
	Treemap              TreemapConfig
	Requirement          RequirementConfig
	Block                BlockConfig
	C4                   C4Config
	Journey              JourneyConfig
	Architecture         ArchitectureConfig
}

// FlowchartConfig holds flowchart-specific layout options.
type FlowchartConfig struct {
	OrderPasses  int
	PortSideBias float32
}

// PaddingConfig holds node padding options.
type PaddingConfig struct {
	NodeHorizontal float32
	NodeVertical   float32
}

// ClassConfig holds class diagram layout options.
type ClassConfig struct {
	CompartmentPadX float32
	CompartmentPadY float32
	MemberFontSize  float32
}

// StateConfig holds state diagram layout options.
type StateConfig struct {
	CompositePadding   float32
	RegionSeparatorPad float32
	StartEndRadius     float32
	ForkBarWidth       float32
	ForkBarHeight      float32
}

// ERConfig holds ER diagram layout options.
type ERConfig struct {
	AttributeRowHeight float32
	ColumnPadding      float32
	HeaderPadding      float32
}

// SequenceConfig holds sequence diagram layout options.
type SequenceConfig struct {
	ParticipantSpacing float32
	MessageSpacing     float32
	ActivationWidth    float32
	NoteMaxWidth       float32
	BoxPadding         float32
	FramePadding       float32
	HeaderHeight       float32
	SelfMessageWidth   float32
}

// KanbanConfig holds Kanban diagram layout options.
type KanbanConfig struct {
	Padding      float32
	SectionWidth float32
	CardSpacing  float32
	HeaderHeight float32
}

// PacketConfig holds Packet diagram layout options.
type PacketConfig struct {
	RowHeight  float32
	BitWidth   float32
	BitsPerRow int
	ShowBits   bool
	PaddingX   float32
	PaddingY   float32
}

// PieConfig holds pie chart layout options.
type PieConfig struct {
	Radius       float32
	InnerRadius  float32
	TextPosition float32
	PaddingX     float32
	PaddingY     float32
}

// QuadrantConfig holds quadrant chart layout options.
type QuadrantConfig struct {
	ChartWidth            float32
	ChartHeight           float32
	PointRadius           float32
	PaddingX              float32
	PaddingY              float32
	QuadrantLabelFontSize float32
	AxisLabelFontSize     float32
}

// TimelineConfig holds timeline diagram layout options.
type TimelineConfig struct {
	PeriodWidth    float32
	EventHeight    float32
	SectionPadding float32
	PaddingX       float32
	PaddingY       float32
}

// GanttConfig holds Gantt chart layout options.
type GanttConfig struct {
	BarHeight            float32
	BarGap               float32
	TopPadding           float32
	SidePadding          float32
	GridLineStartPadding float32
	FontSize             float32
	SectionFontSize      float32
	NumberSectionStyles  int
}

// GitGraphConfig holds GitGraph diagram layout options.
type GitGraphConfig struct {
	CommitRadius  float32
	CommitSpacing float32
	BranchSpacing float32
	PaddingX      float32
	PaddingY      float32
	TagFontSize   float32
}

// XYChartConfig holds XY chart layout options.
type XYChartConfig struct {
	ChartWidth    float32
	ChartHeight   float32
	PaddingX      float32
	PaddingY      float32
	BarWidth      float32 // fraction of band width (0-1)
	TickLength    float32
	AxisFontSize  float32
	TitleFontSize float32
}

// RadarConfig holds radar chart layout options.
type RadarConfig struct {
	Radius       float32
	PaddingX     float32
	PaddingY     float32
	DefaultTicks int
	LabelOffset  float32 // extra distance for axis labels beyond radius
	CurveOpacity float32
}

// MindmapConfig holds mindmap diagram layout options.
type MindmapConfig struct {
	BranchSpacing float32
	LevelSpacing  float32
	PaddingX      float32
	PaddingY      float32
	NodePadding   float32
}

// SankeyConfig holds Sankey diagram layout options.
type SankeyConfig struct {
	ChartWidth  float32
	ChartHeight float32
	NodeWidth   float32
	NodePadding float32
	PaddingX    float32
	PaddingY    float32
}

// TreemapConfig holds Treemap diagram layout options.
type TreemapConfig struct {
	ChartWidth    float32
	ChartHeight   float32
	Padding       float32 // inner padding between rects
	HeaderHeight  float32
	PaddingX      float32
	PaddingY      float32
	LabelFontSize float32
	ValueFontSize float32
}

// RequirementConfig holds requirement diagram layout options.
type RequirementConfig struct {
	NodeMinWidth     float32
	NodePadding      float32
	MetadataFontSize float32
	PaddingX         float32
	PaddingY         float32
}

// BlockConfig holds block diagram layout options.
type BlockConfig struct {
	ColumnGap   float32
	RowGap      float32
	NodePadding float32
	PaddingX    float32
	PaddingY    float32
}

// C4Config holds C4 diagram layout options.
type C4Config struct {
	PersonWidth     float32
	PersonHeight    float32
	SystemWidth     float32
	SystemHeight    float32
	BoundaryPadding float32
	PaddingX        float32
	PaddingY        float32
}

// JourneyConfig holds journey diagram layout options.
type JourneyConfig struct {
	TaskWidth   float32
	TaskHeight  float32
	TaskSpacing float32
	TrackHeight float32
	SectionGap  float32
	PaddingX    float32
	PaddingY    float32
}

// ArchitectureConfig holds architecture diagram layout options.
type ArchitectureConfig struct {
	ServiceWidth  float32
	ServiceHeight float32
	GroupPadding  float32
	JunctionSize  float32
	ColumnGap     float32
	RowGap        float32
	PaddingX      float32
	PaddingY      float32
}

// Default layout-level constants.
const (
	defaultNodeSpacing     = 50
	defaultRankSpacing     = 70
	defaultLabelLineHeight = 1.2
)

// Flowchart defaults.
const (
	defaultFlowchartOrderPasses = 24
)

// Padding defaults.
const (
	defaultPaddingNodeHorizontal = 15
	defaultPaddingNodeVertical   = 10
)

// Class diagram defaults.
const (
	defaultClassCompartmentPadX = 12
	defaultClassCompartmentPadY = 6
	defaultClassMemberFontSize  = 12
)

// State diagram defaults.
const (
	defaultStateCompositePadding   = 20
	defaultStateRegionSeparatorPad = 10
	defaultStateStartEndRadius     = 8
	defaultStateForkBarWidth       = 80
	defaultStateForkBarHeight      = 6
)

// ER diagram defaults.
const (
	defaultERAttributeRowHeight = 22
	defaultERColumnPadding      = 10
	defaultERHeaderPadding      = 8
)

// Sequence diagram defaults.
const (
	defaultSeqParticipantSpacing = 80
	defaultSeqMessageSpacing     = 40
	defaultSeqActivationWidth    = 16
	defaultSeqNoteMaxWidth       = 200
	defaultSeqBoxPadding         = 12
	defaultSeqFramePadding       = 10
	defaultSeqHeaderHeight       = 40
	defaultSeqSelfMessageWidth   = 40
)

// Kanban defaults.
const (
	defaultKanbanPadding      = 8
	defaultKanbanSectionWidth = 200
	defaultKanbanCardSpacing  = 8
	defaultKanbanHeaderHeight = 36
)

// Packet defaults.
const (
	defaultPacketRowHeight  = 32
	defaultPacketBitWidth   = 32
	defaultPacketBitsPerRow = 32
	defaultPacketPaddingX   = 5
	defaultPacketPaddingY   = 5
)

// Pie chart defaults.
const (
	defaultPieRadius       = 150
	defaultPieTextPosition = 0.75
	defaultPiePaddingX     = 20
	defaultPiePaddingY     = 20
)

// Quadrant chart defaults.
const (
	defaultQuadrantChartWidth            = 400
	defaultQuadrantChartHeight           = 400
	defaultQuadrantPointRadius           = 5
	defaultQuadrantPaddingX              = 40
	defaultQuadrantPaddingY              = 40
	defaultQuadrantQuadrantLabelFontSize = 14
	defaultQuadrantAxisLabelFontSize     = 12
)

// Timeline defaults.
const (
	defaultTimelinePeriodWidth    = 150
	defaultTimelineEventHeight    = 30
	defaultTimelineSectionPadding = 10
	defaultTimelinePaddingX       = 20
	defaultTimelinePaddingY       = 20
)

// Gantt chart defaults.
const (
	defaultGanttBarHeight            = 20
	defaultGanttBarGap               = 4
	defaultGanttTopPadding           = 50
	defaultGanttSidePadding          = 75
	defaultGanttGridLineStartPadding = 35
	defaultGanttFontSize             = 11
	defaultGanttSectionFontSize      = 11
	defaultGanttNumberSectionStyles  = 4
)

// GitGraph defaults.
const (
	defaultGitGraphCommitRadius  = 8
	defaultGitGraphCommitSpacing = 60
	defaultGitGraphBranchSpacing = 40
	defaultGitGraphPaddingX      = 30
	defaultGitGraphPaddingY      = 30
	defaultGitGraphTagFontSize   = 11
)

// XY chart defaults.
const (
	defaultXYChartWidth    = 700
	defaultXYChartHeight   = 500
	defaultXYPaddingX      = 60
	defaultXYPaddingY      = 40
	defaultXYBarWidth      = 0.6
	defaultXYTickLength    = 5
	defaultXYAxisFontSize  = 12
	defaultXYTitleFontSize = 16
)

// Radar chart defaults.
const (
	defaultRadarRadius       = 200
	defaultRadarPaddingX     = 40
	defaultRadarPaddingY     = 40
	defaultRadarDefaultTicks = 5
	defaultRadarLabelOffset  = 20
	defaultRadarCurveOpacity = 0.3
)

// Mindmap defaults.
const (
	defaultMindmapBranchSpacing = 80
	defaultMindmapLevelSpacing  = 60
	defaultMindmapPaddingX      = 40
	defaultMindmapPaddingY      = 40
	defaultMindmapNodePadding   = 12
)

// Sankey defaults.
const (
	defaultSankeyChartWidth  = 800
	defaultSankeyChartHeight = 400
	defaultSankeyNodeWidth   = 20
	defaultSankeyNodePadding = 10
	defaultSankeyPaddingX    = 40
	defaultSankeyPaddingY    = 20
)

// Treemap defaults.
const (
	defaultTreemapChartWidth    = 600
	defaultTreemapChartHeight   = 400
	defaultTreemapPadding       = 4
	defaultTreemapHeaderHeight  = 24
	defaultTreemapPaddingX      = 10
	defaultTreemapPaddingY      = 10
	defaultTreemapLabelFontSize = 12
	defaultTreemapValueFontSize = 10
)

// Requirement diagram defaults.
const (
	defaultRequirementNodeMinWidth     = 180
	defaultRequirementNodePadding      = 12
	defaultRequirementMetadataFontSize = 11
	defaultRequirementPaddingX         = 10
	defaultRequirementPaddingY         = 10
)

// Block diagram defaults.
const (
	defaultBlockColumnGap   = 20
	defaultBlockRowGap      = 20
	defaultBlockNodePadding = 12
	defaultBlockPaddingX    = 20
	defaultBlockPaddingY    = 20
)

// C4 diagram defaults.
const (
	defaultC4PersonWidth     = 160
	defaultC4PersonHeight    = 180
	defaultC4SystemWidth     = 200
	defaultC4SystemHeight    = 120
	defaultC4BoundaryPadding = 20
	defaultC4PaddingX        = 20
	defaultC4PaddingY        = 20
)

// Journey diagram defaults.
const (
	defaultJourneyTaskWidth   = 120
	defaultJourneyTaskHeight  = 50
	defaultJourneyTaskSpacing = 20
	defaultJourneyTrackHeight = 200
	defaultJourneySectionGap  = 10
	defaultJourneyPaddingX    = 30
	defaultJourneyPaddingY    = 40
)

// Architecture diagram defaults.
const (
	defaultArchServiceWidth  = 120
	defaultArchServiceHeight = 80
	defaultArchGroupPadding  = 30
	defaultArchJunctionSize  = 10
	defaultArchColumnGap     = 60
	defaultArchRowGap        = 60
	defaultArchPaddingX      = 30
	defaultArchPaddingY      = 30
)

func defaultFlowchartConfig() FlowchartConfig {
	return FlowchartConfig{
		OrderPasses:  defaultFlowchartOrderPasses,
		PortSideBias: 0.0,
	}
}

func defaultPaddingConfig() PaddingConfig {
	return PaddingConfig{
		NodeHorizontal: defaultPaddingNodeHorizontal,
		NodeVertical:   defaultPaddingNodeVertical,
	}
}

func defaultClassConfig() ClassConfig {
	return ClassConfig{
		CompartmentPadX: defaultClassCompartmentPadX,
		CompartmentPadY: defaultClassCompartmentPadY,
		MemberFontSize:  defaultClassMemberFontSize,
	}
}

func defaultStateConfig() StateConfig {
	return StateConfig{
		CompositePadding:   defaultStateCompositePadding,
		RegionSeparatorPad: defaultStateRegionSeparatorPad,
		StartEndRadius:     defaultStateStartEndRadius,
		ForkBarWidth:       defaultStateForkBarWidth,
		ForkBarHeight:      defaultStateForkBarHeight,
	}
}

func defaultERConfig() ERConfig {
	return ERConfig{
		AttributeRowHeight: defaultERAttributeRowHeight,
		ColumnPadding:      defaultERColumnPadding,
		HeaderPadding:      defaultERHeaderPadding,
	}
}

func defaultSequenceConfig() SequenceConfig {
	return SequenceConfig{
		ParticipantSpacing: defaultSeqParticipantSpacing,
		MessageSpacing:     defaultSeqMessageSpacing,
		ActivationWidth:    defaultSeqActivationWidth,
		NoteMaxWidth:       defaultSeqNoteMaxWidth,
		BoxPadding:         defaultSeqBoxPadding,
		FramePadding:       defaultSeqFramePadding,
		HeaderHeight:       defaultSeqHeaderHeight,
		SelfMessageWidth:   defaultSeqSelfMessageWidth,
	}
}

func defaultKanbanConfig() KanbanConfig {
	return KanbanConfig{
		Padding:      defaultKanbanPadding,
		SectionWidth: defaultKanbanSectionWidth,
		CardSpacing:  defaultKanbanCardSpacing,
		HeaderHeight: defaultKanbanHeaderHeight,
	}
}

func defaultPacketConfig() PacketConfig {
	return PacketConfig{
		RowHeight:  defaultPacketRowHeight,
		BitWidth:   defaultPacketBitWidth,
		BitsPerRow: defaultPacketBitsPerRow,
		ShowBits:   true,
		PaddingX:   defaultPacketPaddingX,
		PaddingY:   defaultPacketPaddingY,
	}
}

func defaultPieConfig() PieConfig {
	return PieConfig{
		Radius:       defaultPieRadius,
		InnerRadius:  0,
		TextPosition: defaultPieTextPosition,
		PaddingX:     defaultPiePaddingX,
		PaddingY:     defaultPiePaddingY,
	}
}

func defaultQuadrantConfig() QuadrantConfig {
	return QuadrantConfig{
		ChartWidth:            defaultQuadrantChartWidth,
		ChartHeight:           defaultQuadrantChartHeight,
		PointRadius:           defaultQuadrantPointRadius,
		PaddingX:              defaultQuadrantPaddingX,
		PaddingY:              defaultQuadrantPaddingY,
		QuadrantLabelFontSize: defaultQuadrantQuadrantLabelFontSize,
		AxisLabelFontSize:     defaultQuadrantAxisLabelFontSize,
	}
}

func defaultTimelineConfig() TimelineConfig {
	return TimelineConfig{
		PeriodWidth:    defaultTimelinePeriodWidth,
		EventHeight:    defaultTimelineEventHeight,
		SectionPadding: defaultTimelineSectionPadding,
		PaddingX:       defaultTimelinePaddingX,
		PaddingY:       defaultTimelinePaddingY,
	}
}

func defaultGanttConfig() GanttConfig {
	return GanttConfig{
		BarHeight:            defaultGanttBarHeight,
		BarGap:               defaultGanttBarGap,
		TopPadding:           defaultGanttTopPadding,
		SidePadding:          defaultGanttSidePadding,
		GridLineStartPadding: defaultGanttGridLineStartPadding,
		FontSize:             defaultGanttFontSize,
		SectionFontSize:      defaultGanttSectionFontSize,
		NumberSectionStyles:  defaultGanttNumberSectionStyles,
	}
}

func defaultGitGraphConfig() GitGraphConfig {
	return GitGraphConfig{
		CommitRadius:  defaultGitGraphCommitRadius,
		CommitSpacing: defaultGitGraphCommitSpacing,
		BranchSpacing: defaultGitGraphBranchSpacing,
		PaddingX:      defaultGitGraphPaddingX,
		PaddingY:      defaultGitGraphPaddingY,
		TagFontSize:   defaultGitGraphTagFontSize,
	}
}

func defaultXYChartConfig() XYChartConfig {
	return XYChartConfig{
		ChartWidth:    defaultXYChartWidth,
		ChartHeight:   defaultXYChartHeight,
		PaddingX:      defaultXYPaddingX,
		PaddingY:      defaultXYPaddingY,
		BarWidth:      defaultXYBarWidth,
		TickLength:    defaultXYTickLength,
		AxisFontSize:  defaultXYAxisFontSize,
		TitleFontSize: defaultXYTitleFontSize,
	}
}

func defaultRadarConfig() RadarConfig {
	return RadarConfig{
		Radius:       defaultRadarRadius,
		PaddingX:     defaultRadarPaddingX,
		PaddingY:     defaultRadarPaddingY,
		DefaultTicks: defaultRadarDefaultTicks,
		LabelOffset:  defaultRadarLabelOffset,
		CurveOpacity: defaultRadarCurveOpacity,
	}
}

func defaultMindmapConfig() MindmapConfig {
	return MindmapConfig{
		BranchSpacing: defaultMindmapBranchSpacing,
		LevelSpacing:  defaultMindmapLevelSpacing,
		PaddingX:      defaultMindmapPaddingX,
		PaddingY:      defaultMindmapPaddingY,
		NodePadding:   defaultMindmapNodePadding,
	}
}

func defaultSankeyConfig() SankeyConfig {
	return SankeyConfig{
		ChartWidth:  defaultSankeyChartWidth,
		ChartHeight: defaultSankeyChartHeight,
		NodeWidth:   defaultSankeyNodeWidth,
		NodePadding: defaultSankeyNodePadding,
		PaddingX:    defaultSankeyPaddingX,
		PaddingY:    defaultSankeyPaddingY,
	}
}

func defaultTreemapConfig() TreemapConfig {
	return TreemapConfig{
		ChartWidth:    defaultTreemapChartWidth,
		ChartHeight:   defaultTreemapChartHeight,
		Padding:       defaultTreemapPadding,
		HeaderHeight:  defaultTreemapHeaderHeight,
		PaddingX:      defaultTreemapPaddingX,
		PaddingY:      defaultTreemapPaddingY,
		LabelFontSize: defaultTreemapLabelFontSize,
		ValueFontSize: defaultTreemapValueFontSize,
	}
}

func defaultRequirementConfig() RequirementConfig {
	return RequirementConfig{
		NodeMinWidth:     defaultRequirementNodeMinWidth,
		NodePadding:      defaultRequirementNodePadding,
		MetadataFontSize: defaultRequirementMetadataFontSize,
		PaddingX:         defaultRequirementPaddingX,
		PaddingY:         defaultRequirementPaddingY,
	}
}

func defaultBlockConfig() BlockConfig {
	return BlockConfig{
		ColumnGap:   defaultBlockColumnGap,
		RowGap:      defaultBlockRowGap,
		NodePadding: defaultBlockNodePadding,
		PaddingX:    defaultBlockPaddingX,
		PaddingY:    defaultBlockPaddingY,
	}
}

func defaultC4Config() C4Config {
	return C4Config{
		PersonWidth:     defaultC4PersonWidth,
		PersonHeight:    defaultC4PersonHeight,
		SystemWidth:     defaultC4SystemWidth,
		SystemHeight:    defaultC4SystemHeight,
		BoundaryPadding: defaultC4BoundaryPadding,
		PaddingX:        defaultC4PaddingX,
		PaddingY:        defaultC4PaddingY,
	}
}

func defaultJourneyConfig() JourneyConfig {
	return JourneyConfig{
		TaskWidth:   defaultJourneyTaskWidth,
		TaskHeight:  defaultJourneyTaskHeight,
		TaskSpacing: defaultJourneyTaskSpacing,
		TrackHeight: defaultJourneyTrackHeight,
		SectionGap:  defaultJourneySectionGap,
		PaddingX:    defaultJourneyPaddingX,
		PaddingY:    defaultJourneyPaddingY,
	}
}

func defaultArchitectureConfig() ArchitectureConfig {
	return ArchitectureConfig{
		ServiceWidth:  defaultArchServiceWidth,
		ServiceHeight: defaultArchServiceHeight,
		GroupPadding:  defaultArchGroupPadding,
		JunctionSize:  defaultArchJunctionSize,
		ColumnGap:     defaultArchColumnGap,
		RowGap:        defaultArchRowGap,
		PaddingX:      defaultArchPaddingX,
		PaddingY:      defaultArchPaddingY,
	}
}

// DefaultLayout returns a Layout with default values for diagram rendering.
func DefaultLayout() *Layout {
	return &Layout{
		NodeSpacing:     defaultNodeSpacing,
		RankSpacing:     defaultRankSpacing,
		LabelLineHeight: defaultLabelLineHeight,
		Flowchart:       defaultFlowchartConfig(),
		Padding:         defaultPaddingConfig(),
		Class:           defaultClassConfig(),
		State:           defaultStateConfig(),
		ER:              defaultERConfig(),
		Sequence:        defaultSequenceConfig(),
		Kanban:          defaultKanbanConfig(),
		Packet:          defaultPacketConfig(),
		Pie:             defaultPieConfig(),
		Quadrant:        defaultQuadrantConfig(),
		Timeline:        defaultTimelineConfig(),
		Gantt:           defaultGanttConfig(),
		GitGraph:        defaultGitGraphConfig(),
		XYChart:         defaultXYChartConfig(),
		Radar:           defaultRadarConfig(),
		Mindmap:         defaultMindmapConfig(),
		Sankey:          defaultSankeyConfig(),
		Treemap:         defaultTreemapConfig(),
		Requirement:     defaultRequirementConfig(),
		Block:           defaultBlockConfig(),
		C4:              defaultC4Config(),
		Journey:         defaultJourneyConfig(),
		Architecture:    defaultArchitectureConfig(),
	}
}
