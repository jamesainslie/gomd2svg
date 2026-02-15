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
	}
}
