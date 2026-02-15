package theme

// Theme defines the visual appearance for rendering Mermaid diagrams.
// All color fields are CSS color strings (hex, hsl, etc.).
type Theme struct {
	FontFamily           string
	FontSize             float32
	Background           string
	PrimaryColor         string
	PrimaryBorderColor   string
	PrimaryTextColor     string
	SecondaryColor       string
	SecondaryBorderColor string
	SecondaryTextColor   string
	TertiaryColor        string
	TertiaryBorderColor  string
	LineColor            string
	TextColor            string

	ClusterBackground string
	ClusterBorder     string
	NodeBorderColor   string

	NoteBackground  string
	NoteBorderColor string
	NoteTextColor   string

	ActorBorder     string
	ActorBackground string
	ActorTextColor  string
	ActorLineColor  string

	SignalColor     string
	SignalTextColor string

	ActivationBorderColor string
	ActivationBackground  string
	SequenceNumberColor   string

	EdgeLabelBackground string
	LabelTextColor      string
	LoopTextColor       string

	PieTitleTextSize    float32
	PieTitleTextColor   string
	PieSectionTextSize  float32
	PieSectionTextColor string
	PieStrokeColor      string
	PieStrokeWidth      float32
	PieOuterStrokeWidth float32
	PieOuterStrokeColor string
	PieOpacity          float32
	PieColors           []string

	// Quadrant chart colors
	QuadrantFill1     string
	QuadrantFill2     string
	QuadrantFill3     string
	QuadrantFill4     string
	QuadrantPointFill string

	// Class diagram colors
	ClassHeaderBg string
	ClassBodyBg   string
	ClassBorder   string

	// State diagram colors
	StateFill         string
	StateBorder       string
	StateStartEnd     string
	CompositeHeaderBg string

	// ER diagram colors
	EntityHeaderBg string
	EntityBodyBg   string
	EntityBorder   string

	// Timeline diagram colors
	TimelineSectionColors []string
	TimelineEventFill     string
	TimelineEventBorder   string

	// Gantt diagram colors
	GanttTaskFill         string
	GanttTaskBorder       string
	GanttCritFill         string
	GanttCritBorder       string
	GanttDoneFill         string
	GanttActiveFill       string
	GanttMilestoneFill    string
	GanttGridColor        string
	GanttTodayMarkerColor string
	GanttSectionColors    []string

	// GitGraph diagram colors
	GitBranchColors  []string
	GitCommitFill    string
	GitCommitStroke  string
	GitTagFill       string
	GitTagBorder     string
	GitHighlightFill string

	// XYChart colors
	XYChartColors    []string
	XYChartAxisColor string
	XYChartGridColor string

	// Radar colors
	RadarCurveColors    []string
	RadarAxisColor      string
	RadarGraticuleColor string
	RadarCurveOpacity   float32

	// Mindmap colors
	MindmapBranchColors []string
	MindmapNodeFill     string
	MindmapNodeBorder   string

	// Sankey colors
	SankeyNodeColors  []string
	SankeyLinkColor   string
	SankeyLinkOpacity float32

	// Treemap colors
	TreemapColors    []string
	TreemapBorder    string
	TreemapTextColor string

	// Requirement diagram colors
	RequirementFill   string
	RequirementBorder string

	// Block diagram colors
	BlockColors     []string
	BlockNodeBorder string

	// Journey diagram colors
	JourneySectionColors []string
	JourneyTaskFill      string
	JourneyTaskBorder    string
	JourneyTaskText      string
	JourneyScoreColors   [5]string // score 1-5 indicator colors (red to green)

	// Architecture diagram colors
	ArchServiceFill   string
	ArchServiceBorder string
	ArchServiceText   string
	ArchGroupFill     string
	ArchGroupBorder   string
	ArchGroupText     string
	ArchEdgeColor     string
	ArchJunctionFill  string

	// C4 diagram colors
	C4PersonColor    string
	C4SystemColor    string
	C4ContainerColor string
	C4ComponentColor string
	C4ExternalColor  string
	C4BoundaryColor  string
	C4TextColor      string
}

// Modern returns a theme with a clean, modern color palette using the Inter font.
func Modern() *Theme {
	return &Theme{
		FontFamily: "Inter, sans-serif",
		FontSize:   14,
		Background: "#FFFFFF",

		PrimaryColor:       "#4C78A8",
		PrimaryBorderColor: "#3B6492",
		PrimaryTextColor:   "#1A1A2E",

		SecondaryColor:       "#72B7B2",
		SecondaryBorderColor: "#5A9994",
		SecondaryTextColor:   "#1A1A2E",

		TertiaryColor:       "#EECA3B",
		TertiaryBorderColor: "#C9A820",

		LineColor: "#6E7B8B",
		TextColor: "#333344",

		ClusterBackground: "#F0F4F8",
		ClusterBorder:     "#B0C4DE",
		NodeBorderColor:   "#3B6492",

		NoteBackground:  "#FFF3CD",
		NoteBorderColor: "#FFECB5",
		NoteTextColor:   "#664D03",

		ActorBorder:     "#3B6492",
		ActorBackground: "#4C78A8",
		ActorTextColor:  "#FFFFFF",
		ActorLineColor:  "#6E7B8B",

		SignalColor:     "#6E7B8B",
		SignalTextColor: "#333344",

		ActivationBorderColor: "#3B6492",
		ActivationBackground:  "#E8EFF5",
		SequenceNumberColor:   "#FFFFFF",

		EdgeLabelBackground: "#FFFFFF",
		LabelTextColor:      "#333344",
		LoopTextColor:       "#333344",

		PieTitleTextSize:    18,
		PieTitleTextColor:   "#333344",
		PieSectionTextSize:  14,
		PieSectionTextColor: "#FFFFFF",
		PieStrokeColor:      "#FFFFFF",
		PieStrokeWidth:      2,
		PieOuterStrokeWidth: 2,
		PieOuterStrokeColor: "#3B6492",
		PieOpacity:          0.85,
		PieColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},

		QuadrantFill1:     "#E8EFF5",
		QuadrantFill2:     "#F0F4F8",
		QuadrantFill3:     "#F5F5F5",
		QuadrantFill4:     "#FFF8E1",
		QuadrantPointFill: "#4C78A8",

		ClassHeaderBg: "#4C78A8",
		ClassBodyBg:   "#F0F4F8",
		ClassBorder:   "#3B6492",

		StateFill:         "#F0F4F8",
		StateBorder:       "#3B6492",
		StateStartEnd:     "#333344",
		CompositeHeaderBg: "#E8EFF5",

		EntityHeaderBg: "#4C78A8",
		EntityBodyBg:   "#F0F4F8",
		EntityBorder:   "#3B6492",

		TimelineSectionColors: []string{"#E8EFF5", "#F0E8F0", "#E8F5E8", "#FFF8E1"},
		TimelineEventFill:     "#4C78A8",
		TimelineEventBorder:   "#3B6492",

		GanttTaskFill:         "#4C78A8",
		GanttTaskBorder:       "#3B6492",
		GanttCritFill:         "#E45756",
		GanttCritBorder:       "#CC3333",
		GanttDoneFill:         "#B0C4DE",
		GanttActiveFill:       "#72B7B2",
		GanttMilestoneFill:    "#F58518",
		GanttGridColor:        "#E0E0E0",
		GanttTodayMarkerColor: "#E45756",
		GanttSectionColors:    []string{"#F0F4F8", "#FFF8E1", "#F0E8F0", "#E8F5E8"},

		GitBranchColors: []string{
			"#4C78A8", "#E45756", "#54A24B", "#F58518",
			"#72B7B2", "#B279A2", "#EECA3B", "#FF9DA6",
		},
		GitCommitFill:    "#333344",
		GitCommitStroke:  "#333344",
		GitTagFill:       "#EECA3B",
		GitTagBorder:     "#C9A820",
		GitHighlightFill: "#F58518",

		XYChartColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		XYChartAxisColor: "#6E7B8B",
		XYChartGridColor: "#E0E0E0",

		RadarCurveColors: []string{
			"#4C78A8", "#E45756", "#54A24B", "#F58518",
			"#72B7B2", "#B279A2", "#EECA3B", "#FF9DA6",
		},
		RadarAxisColor:      "#6E7B8B",
		RadarGraticuleColor: "#E0E0E0",
		RadarCurveOpacity:   0.3,

		MindmapBranchColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		MindmapNodeFill:   "#F0F4F8",
		MindmapNodeBorder: "#3B6492",

		SankeyNodeColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		SankeyLinkColor:   "#6E7B8B",
		SankeyLinkOpacity: 0.4,

		TreemapColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		TreemapBorder:    "#3B6492",
		TreemapTextColor: "#FFFFFF",

		RequirementFill:   "#F0F4F8",
		RequirementBorder: "#3B6492",

		BlockColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		BlockNodeBorder: "#3B6492",

		JourneySectionColors: []string{"#E8F5E9", "#E3F2FD", "#FFF3E0", "#F3E5F5", "#E0F7FA"},
		JourneyTaskFill:      "#FFFFFF",
		JourneyTaskBorder:    "#90A4AE",
		JourneyTaskText:      "#263238",
		JourneyScoreColors:   [5]string{"#F44336", "#FF9800", "#FFC107", "#8BC34A", "#4CAF50"},

		ArchServiceFill:   "#E3F2FD",
		ArchServiceBorder: "#1565C0",
		ArchServiceText:   "#0D47A1",
		ArchGroupFill:     "#F5F5F5",
		ArchGroupBorder:   "#BDBDBD",
		ArchGroupText:     "#424242",
		ArchEdgeColor:     "#546E7A",
		ArchJunctionFill:  "#546E7A",

		C4PersonColor:    "#08427B",
		C4SystemColor:    "#1168BD",
		C4ContainerColor: "#438DD5",
		C4ComponentColor: "#85BBF0",
		C4ExternalColor:  "#999999",
		C4BoundaryColor:  "#444444",
		C4TextColor:      "#FFFFFF",
	}
}

// MermaidDefault returns the classic mermaid.js default theme.
func MermaidDefault() *Theme {
	return &Theme{
		FontFamily: "trebuchet ms, verdana, arial, sans-serif",
		FontSize:   16,
		Background: "#FFFFFF",

		PrimaryColor:       "#ECECFF",
		PrimaryBorderColor: "#9370DB",
		PrimaryTextColor:   "#333",

		SecondaryColor:       "#ffffde",
		SecondaryBorderColor: "#aaaa33",
		SecondaryTextColor:   "#333",

		TertiaryColor:       "#fff0f0",
		TertiaryBorderColor: "#BB0000",

		LineColor: "#333",
		TextColor: "#333",

		ClusterBackground: "#ffffde",
		ClusterBorder:     "#aaaa33",
		NodeBorderColor:   "#9370DB",

		NoteBackground:  "#fff5ad",
		NoteBorderColor: "#aaaa33",
		NoteTextColor:   "#333",

		ActorBorder:     "#9370DB",
		ActorBackground: "#ECECFF",
		ActorTextColor:  "#333",
		ActorLineColor:  "#888",

		SignalColor:     "#333",
		SignalTextColor: "#333",

		ActivationBorderColor: "#666",
		ActivationBackground:  "#f4f4f4",
		SequenceNumberColor:   "#fff",

		EdgeLabelBackground: "#e8e8e8",
		LabelTextColor:      "#333",
		LoopTextColor:       "#333",

		PieTitleTextSize:    25,
		PieTitleTextColor:   "#333",
		PieSectionTextSize:  17,
		PieSectionTextColor: "#333",
		PieStrokeColor:      "#ccc",
		PieStrokeWidth:      2,
		PieOuterStrokeWidth: 2,
		PieOuterStrokeColor: "#999",
		PieOpacity:          0.7,
		PieColors: []string{
			"#4C78A8", "#48A9A6", "#E4E36A", "#F4A261",
			"#E76F51", "#7FB069", "#D08AC0", "#F7B7A3",
		},

		QuadrantFill1:     "#f0ece8",
		QuadrantFill2:     "#e8f0e8",
		QuadrantFill3:     "#f0f0f0",
		QuadrantFill4:     "#ece8f0",
		QuadrantPointFill: "#4C78A8",

		ClassHeaderBg: "#ECECFF",
		ClassBodyBg:   "#FFFFFF",
		ClassBorder:   "#9370DB",

		StateFill:         "#ECECFF",
		StateBorder:       "#9370DB",
		StateStartEnd:     "#333",
		CompositeHeaderBg: "#f4f4f4",

		EntityHeaderBg: "#ECECFF",
		EntityBodyBg:   "#FFFFFF",
		EntityBorder:   "#9370DB",

		TimelineSectionColors: []string{"#ffffde", "#f0ece8", "#e8f0e8", "#ece8f0"},
		TimelineEventFill:     "#ECECFF",
		TimelineEventBorder:   "#9370DB",

		GanttTaskFill:         "#8a90dd",
		GanttTaskBorder:       "#534fbc",
		GanttCritFill:         "#ff8888",
		GanttCritBorder:       "#ff0000",
		GanttDoneFill:         "#d3d3d3",
		GanttActiveFill:       "#8a90dd",
		GanttMilestoneFill:    "#E76F51",
		GanttGridColor:        "#ddd",
		GanttTodayMarkerColor: "#d42",
		GanttSectionColors:    []string{"#ffffde", "#ffffff", "#ffffde", "#ffffff"},

		GitBranchColors: []string{
			"#9370DB", "#ff0000", "#00cc00", "#F58518",
			"#48A9A6", "#E76F51", "#D08AC0", "#F7B7A3",
		},
		GitCommitFill:    "#333",
		GitCommitStroke:  "#333",
		GitTagFill:       "#ffffde",
		GitTagBorder:     "#aaaa33",
		GitHighlightFill: "#ff0000",

		XYChartColors: []string{
			"#4C78A8", "#48A9A6", "#E4E36A", "#F4A261",
			"#E76F51", "#7FB069", "#D08AC0", "#F7B7A3",
		},
		XYChartAxisColor: "#333",
		XYChartGridColor: "#ddd",

		RadarCurveColors: []string{
			"#9370DB", "#E76F51", "#7FB069", "#F4A261",
			"#48A9A6", "#D08AC0", "#E4E36A", "#F7B7A3",
		},
		RadarAxisColor:      "#888",
		RadarGraticuleColor: "#ddd",
		RadarCurveOpacity:   0.3,

		MindmapBranchColors: []string{
			"#9370DB", "#E76F51", "#7FB069", "#F4A261",
			"#48A9A6", "#D08AC0", "#E4E36A", "#F7B7A3",
		},
		MindmapNodeFill:   "#ECECFF",
		MindmapNodeBorder: "#9370DB",

		SankeyNodeColors: []string{
			"#4C78A8", "#48A9A6", "#E4E36A", "#F4A261",
			"#E76F51", "#7FB069", "#D08AC0", "#F7B7A3",
		},
		SankeyLinkColor:   "#888",
		SankeyLinkOpacity: 0.4,

		TreemapColors: []string{
			"#4C78A8", "#48A9A6", "#E4E36A", "#F4A261",
			"#E76F51", "#7FB069", "#D08AC0", "#F7B7A3",
		},
		TreemapBorder:    "#9370DB",
		TreemapTextColor: "#333",

		RequirementFill:   "#ECECFF",
		RequirementBorder: "#9370DB",

		BlockColors: []string{
			"#9370DB", "#E76F51", "#7FB069", "#F4A261",
			"#48A9A6", "#D08AC0", "#E4E36A", "#F7B7A3",
		},
		BlockNodeBorder: "#9370DB",

		JourneySectionColors: []string{"#f4f4de", "#dde4ff", "#ffe6cc", "#e8d5f5", "#d5f5e3"},
		JourneyTaskFill:      "#fff",
		JourneyTaskBorder:    "#999",
		JourneyTaskText:      "#333",
		JourneyScoreColors:   [5]string{"#d32f2f", "#e65100", "#fbc02d", "#689f38", "#388e3c"},

		ArchServiceFill:   "#cde4ff",
		ArchServiceBorder: "#326ce5",
		ArchServiceText:   "#1a237e",
		ArchGroupFill:     "#ececff",
		ArchGroupBorder:   "#9370db",
		ArchGroupText:     "#333",
		ArchEdgeColor:     "#666",
		ArchJunctionFill:  "#666",

		C4PersonColor:    "#08427B",
		C4SystemColor:    "#1168BD",
		C4ContainerColor: "#438DD5",
		C4ComponentColor: "#85BBF0",
		C4ExternalColor:  "#999999",
		C4BoundaryColor:  "#444444",
		C4TextColor:      "#FFFFFF",
	}
}
