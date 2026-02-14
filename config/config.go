package config

// Layout holds all configuration for diagram layout computation.
type Layout struct {
	NodeSpacing          float32
	RankSpacing          float32
	LabelLineHeight      float32
	PreferredAspectRatio *float32
	Flowchart            FlowchartConfig
	Padding              PaddingConfig
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
	}
}
