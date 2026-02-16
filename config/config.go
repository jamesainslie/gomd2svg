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

// DefaultLayout returns a Layout with default values for diagram rendering.
func DefaultLayout() *Layout {
	return &Layout{
		NodeSpacing:     50,
		RankSpacing:     70,
		LabelLineHeight: 1.2,
		Flowchart: FlowchartConfig{
			OrderPasses:  24,
			PortSideBias: 0.0,
		},
		Padding: PaddingConfig{
			NodeHorizontal: 15,
			NodeVertical:   10,
		},
		Class: ClassConfig{
			CompartmentPadX: 12,
			CompartmentPadY: 6,
			MemberFontSize:  12,
		},
		State: StateConfig{
			CompositePadding:   20,
			RegionSeparatorPad: 10,
			StartEndRadius:     8,
			ForkBarWidth:       80,
			ForkBarHeight:      6,
		},
		ER: ERConfig{
			AttributeRowHeight: 22,
			ColumnPadding:      10,
			HeaderPadding:      8,
		},
		Sequence: SequenceConfig{
			ParticipantSpacing: 80,
			MessageSpacing:     40,
			ActivationWidth:    16,
			NoteMaxWidth:       200,
			BoxPadding:         12,
			FramePadding:       10,
			HeaderHeight:       40,
			SelfMessageWidth:   40,
		},
		Kanban: KanbanConfig{
			Padding:      8,
			SectionWidth: 200,
			CardSpacing:  8,
			HeaderHeight: 36,
		},
		Packet: PacketConfig{
			RowHeight:  32,
			BitWidth:   32,
			BitsPerRow: 32,
			ShowBits:   true,
			PaddingX:   5,
			PaddingY:   5,
		},
		Pie: PieConfig{
			Radius:       150,
			InnerRadius:  0,
			TextPosition: 0.75,
			PaddingX:     20,
			PaddingY:     20,
		},
		Quadrant: QuadrantConfig{
			ChartWidth:            400,
			ChartHeight:           400,
			PointRadius:           5,
			PaddingX:              40,
			PaddingY:              40,
			QuadrantLabelFontSize: 14,
			AxisLabelFontSize:     12,
		},
		Timeline: TimelineConfig{
			PeriodWidth:    150,
			EventHeight:    30,
			SectionPadding: 10,
			PaddingX:       20,
			PaddingY:       20,
		},
		Gantt: GanttConfig{
			BarHeight:            20,
			BarGap:               4,
			TopPadding:           50,
			SidePadding:          75,
			GridLineStartPadding: 35,
			FontSize:             11,
			SectionFontSize:      11,
			NumberSectionStyles:  4,
		},
		GitGraph: GitGraphConfig{
			CommitRadius:  8,
			CommitSpacing: 60,
			BranchSpacing: 40,
			PaddingX:      30,
			PaddingY:      30,
			TagFontSize:   11,
		},
		XYChart: XYChartConfig{
			ChartWidth:    700,
			ChartHeight:   500,
			PaddingX:      60,
			PaddingY:      40,
			BarWidth:      0.6,
			TickLength:    5,
			AxisFontSize:  12,
			TitleFontSize: 16,
		},
		Radar: RadarConfig{
			Radius:       200,
			PaddingX:     40,
			PaddingY:     40,
			DefaultTicks: 5,
			LabelOffset:  20,
			CurveOpacity: 0.3,
		},
		Mindmap: MindmapConfig{
			BranchSpacing: 80,
			LevelSpacing:  60,
			PaddingX:      40,
			PaddingY:      40,
			NodePadding:   12,
		},
		Sankey: SankeyConfig{
			ChartWidth:  800,
			ChartHeight: 400,
			NodeWidth:   20,
			NodePadding: 10,
			PaddingX:    40,
			PaddingY:    20,
		},
		Treemap: TreemapConfig{
			ChartWidth:    600,
			ChartHeight:   400,
			Padding:       4,
			HeaderHeight:  24,
			PaddingX:      10,
			PaddingY:      10,
			LabelFontSize: 12,
			ValueFontSize: 10,
		},
		Requirement: RequirementConfig{
			NodeMinWidth:     180,
			NodePadding:      12,
			MetadataFontSize: 11,
			PaddingX:         10,
			PaddingY:         10,
		},
		Block: BlockConfig{
			ColumnGap:   20,
			RowGap:      20,
			NodePadding: 12,
			PaddingX:    20,
			PaddingY:    20,
		},
		C4: C4Config{
			PersonWidth:     160,
			PersonHeight:    180,
			SystemWidth:     200,
			SystemHeight:    120,
			BoundaryPadding: 20,
			PaddingX:        20,
			PaddingY:        20,
		},
		Journey: JourneyConfig{
			TaskWidth:   120,
			TaskHeight:  50,
			TaskSpacing: 20,
			TrackHeight: 200,
			SectionGap:  10,
			PaddingX:    30,
			PaddingY:    40,
		},
		Architecture: ArchitectureConfig{
			ServiceWidth:  120,
			ServiceHeight: 80,
			GroupPadding:  30,
			JunctionSize:  10,
			ColumnGap:     60,
			RowGap:        60,
			PaddingX:      30,
			PaddingY:      30,
		},
	}
}
