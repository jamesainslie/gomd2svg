package theme

// Dark returns a dark theme with a dark background and light text.
//
//nolint:dupl,funlen // each theme preset is intentionally distinct; flat declaration
func Dark() *Theme {
	return &Theme{
		FontFamily: "Inter, sans-serif",
		FontSize:   FontSizeDefault,
		Background: "#1A1A2E",

		PrimaryColor:       "#4C78A8",
		PrimaryBorderColor: "#6B9BD2",
		PrimaryTextColor:   "#E0E0E0",

		SecondaryColor:       "#72B7B2",
		SecondaryBorderColor: "#8FD0CC",
		SecondaryTextColor:   "#E0E0E0",

		TertiaryColor:       "#EECA3B",
		TertiaryBorderColor: "#F5DC6E",

		LineColor: "#A0AEC0",
		TextColor: "#E0E0E0",

		ClusterBackground: "#2D2D44",
		ClusterBorder:     "#4A4A6A",
		NodeBorderColor:   "#6B9BD2",

		NoteBackground:  "#3D3D1E",
		NoteBorderColor: "#6B6B3D",
		NoteTextColor:   "#E8E8C8",

		ActorBorder:     "#6B9BD2",
		ActorBackground: "#4C78A8",
		ActorTextColor:  "#FFFFFF",
		ActorLineColor:  "#A0AEC0",

		SignalColor:     "#A0AEC0",
		SignalTextColor: "#E0E0E0",

		ActivationBorderColor: "#6B9BD2",
		ActivationBackground:  "#2D3748",
		SequenceNumberColor:   "#FFFFFF",

		EdgeLabelBackground: "#1A1A2E",
		LabelTextColor:      "#E0E0E0",
		LoopTextColor:       "#E0E0E0",

		PieTitleTextSize:    PieTitleSizeDefault,
		PieTitleTextColor:   "#E0E0E0",
		PieSectionTextSize:  PieSectionSizeDefault,
		PieSectionTextColor: "#FFFFFF",
		PieStrokeColor:      "#1A1A2E",
		PieStrokeWidth:      PieStrokeWidthDefault,
		PieOuterStrokeWidth: PieStrokeWidthDefault,
		PieOuterStrokeColor: "#6B9BD2",
		PieOpacity:          PieOpacityDefault,
		PieColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},

		QuadrantFill1:     "#2D2D44",
		QuadrantFill2:     "#2D3D2D",
		QuadrantFill3:     "#3D2D2D",
		QuadrantFill4:     "#3D3D1E",
		QuadrantPointFill: "#72B7B2",

		ClassHeaderBg: "#4C78A8",
		ClassBodyBg:   "#2D2D44",
		ClassBorder:   "#6B9BD2",

		StateFill:         "#2D2D44",
		StateBorder:       "#6B9BD2",
		StateStartEnd:     "#E0E0E0",
		CompositeHeaderBg: "#2D3748",

		EntityHeaderBg: "#4C78A8",
		EntityBodyBg:   "#2D2D44",
		EntityBorder:   "#6B9BD2",

		TimelineSectionColors: []string{"#2D2D44", "#2D3D2D", "#3D2D2D", "#3D3D1E"},
		TimelineEventFill:     "#4C78A8",
		TimelineEventBorder:   "#6B9BD2",

		GanttTaskFill:         "#4C78A8",
		GanttTaskBorder:       "#6B9BD2",
		GanttCritFill:         "#E45756",
		GanttCritBorder:       "#FF7070",
		GanttDoneFill:         "#4A4A6A",
		GanttActiveFill:       "#72B7B2",
		GanttMilestoneFill:    "#F58518",
		GanttGridColor:        "#3D3D5C",
		GanttTodayMarkerColor: "#E45756",
		GanttSectionColors:    []string{"#2D2D44", "#3D3D1E", "#2D3D2D", "#3D2D2D"},

		GitBranchColors: []string{
			"#4C78A8", "#E45756", "#54A24B", "#F58518",
			"#72B7B2", "#B279A2", "#EECA3B", "#FF9DA6",
		},
		GitCommitFill:    "#E0E0E0",
		GitCommitStroke:  "#E0E0E0",
		GitTagFill:       "#EECA3B",
		GitTagBorder:     "#C9A820",
		GitHighlightFill: "#F58518",

		XYChartColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		XYChartAxisColor: "#A0AEC0",
		XYChartGridColor: "#3D3D5C",

		RadarCurveColors: []string{
			"#4C78A8", "#E45756", "#54A24B", "#F58518",
			"#72B7B2", "#B279A2", "#EECA3B", "#FF9DA6",
		},
		RadarAxisColor:      "#A0AEC0",
		RadarGraticuleColor: "#3D3D5C",
		RadarCurveOpacity:   RadarCurveOpacityDefault,

		MindmapBranchColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		MindmapNodeFill:   "#2D2D44",
		MindmapNodeBorder: "#6B9BD2",

		SankeyNodeColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		SankeyLinkColor:   "#A0AEC0",
		SankeyLinkOpacity: SankeyLinkOpacityDark,

		TreemapColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		TreemapBorder:    "#6B9BD2",
		TreemapTextColor: "#E0E0E0",

		RequirementFill:   "#2D2D44",
		RequirementBorder: "#6B9BD2",

		BlockColors: []string{
			"#4C78A8", "#72B7B2", "#EECA3B", "#F58518",
			"#E45756", "#54A24B", "#B279A2", "#FF9DA6",
		},
		BlockNodeBorder: "#6B9BD2",

		JourneySectionColors: []string{"#2D3748", "#1E3A2D", "#3A2D1E", "#2D1E3A", "#1E3A3A"},
		JourneyTaskFill:      "#2D2D44",
		JourneyTaskBorder:    "#4A4A6A",
		JourneyTaskText:      "#E0E0E0",
		JourneyScoreColors:   [5]string{"#E45756", "#F58518", "#EECA3B", "#72B7B2", "#54A24B"},

		ArchServiceFill:   "#2D3748",
		ArchServiceBorder: "#6B9BD2",
		ArchServiceText:   "#E0E0E0",
		ArchGroupFill:     "#2D2D44",
		ArchGroupBorder:   "#4A4A6A",
		ArchGroupText:     "#A0AEC0",
		ArchEdgeColor:     "#A0AEC0",
		ArchJunctionFill:  "#A0AEC0",

		C4PersonColor:    "#6B9BD2",
		C4SystemColor:    "#4C78A8",
		C4ContainerColor: "#72B7B2",
		C4ComponentColor: "#8FD0CC",
		C4ExternalColor:  "#4A4A6A",
		C4BoundaryColor:  "#A0AEC0",
		C4TextColor:      "#FFFFFF",
	}
}

// Forest returns a nature-inspired green theme.
//
//nolint:dupl,funlen // each theme preset is intentionally distinct; flat declaration
func Forest() *Theme {
	return &Theme{
		FontFamily: "Inter, sans-serif",
		FontSize:   FontSizeDefault,
		Background: "#FFFFFF",

		PrimaryColor:       "#2D6A4F",
		PrimaryBorderColor: "#1B4332",
		PrimaryTextColor:   "#1A1A2E",

		SecondaryColor:       "#95D5B2",
		SecondaryBorderColor: "#74C69D",
		SecondaryTextColor:   "#1A1A2E",

		TertiaryColor:       "#D8F3DC",
		TertiaryBorderColor: "#B7E4C7",

		LineColor: "#40916C",
		TextColor: "#1B4332",

		ClusterBackground: "#D8F3DC",
		ClusterBorder:     "#74C69D",
		NodeBorderColor:   "#1B4332",

		NoteBackground:  "#FEFAE0",
		NoteBorderColor: "#DDA15E",
		NoteTextColor:   "#606C38",

		ActorBorder:     "#1B4332",
		ActorBackground: "#2D6A4F",
		ActorTextColor:  "#FFFFFF",
		ActorLineColor:  "#40916C",

		SignalColor:     "#40916C",
		SignalTextColor: "#1B4332",

		ActivationBorderColor: "#1B4332",
		ActivationBackground:  "#D8F3DC",
		SequenceNumberColor:   "#FFFFFF",

		EdgeLabelBackground: "#FFFFFF",
		LabelTextColor:      "#1B4332",
		LoopTextColor:       "#1B4332",

		PieTitleTextSize:    PieTitleSizeDefault,
		PieTitleTextColor:   "#1B4332",
		PieSectionTextSize:  PieSectionSizeDefault,
		PieSectionTextColor: "#FFFFFF",
		PieStrokeColor:      "#FFFFFF",
		PieStrokeWidth:      PieStrokeWidthDefault,
		PieOuterStrokeWidth: PieStrokeWidthDefault,
		PieOuterStrokeColor: "#1B4332",
		PieOpacity:          PieOpacityDefault,
		PieColors: []string{
			"#2D6A4F", "#52B788", "#DDA15E", "#BC6C25",
			"#E76F51", "#606C38", "#283618", "#FEFAE0",
		},

		QuadrantFill1:     "#D8F3DC",
		QuadrantFill2:     "#B7E4C7",
		QuadrantFill3:     "#FEFAE0",
		QuadrantFill4:     "#95D5B2",
		QuadrantPointFill: "#2D6A4F",

		ClassHeaderBg: "#2D6A4F",
		ClassBodyBg:   "#D8F3DC",
		ClassBorder:   "#1B4332",

		StateFill:         "#D8F3DC",
		StateBorder:       "#1B4332",
		StateStartEnd:     "#1B4332",
		CompositeHeaderBg: "#B7E4C7",

		EntityHeaderBg: "#2D6A4F",
		EntityBodyBg:   "#D8F3DC",
		EntityBorder:   "#1B4332",

		TimelineSectionColors: []string{"#D8F3DC", "#FEFAE0", "#B7E4C7", "#95D5B2"},
		TimelineEventFill:     "#2D6A4F",
		TimelineEventBorder:   "#1B4332",

		GanttTaskFill:         "#2D6A4F",
		GanttTaskBorder:       "#1B4332",
		GanttCritFill:         "#E76F51",
		GanttCritBorder:       "#BC6C25",
		GanttDoneFill:         "#B7E4C7",
		GanttActiveFill:       "#52B788",
		GanttMilestoneFill:    "#DDA15E",
		GanttGridColor:        "#D8F3DC",
		GanttTodayMarkerColor: "#E76F51",
		GanttSectionColors:    []string{"#D8F3DC", "#FEFAE0", "#B7E4C7", "#95D5B2"},

		GitBranchColors: []string{
			"#2D6A4F", "#E76F51", "#52B788", "#DDA15E",
			"#40916C", "#BC6C25", "#606C38", "#283618",
		},
		GitCommitFill:    "#1B4332",
		GitCommitStroke:  "#1B4332",
		GitTagFill:       "#FEFAE0",
		GitTagBorder:     "#DDA15E",
		GitHighlightFill: "#E76F51",

		XYChartColors: []string{
			"#2D6A4F", "#52B788", "#DDA15E", "#BC6C25",
			"#E76F51", "#606C38", "#283618", "#95D5B2",
		},
		XYChartAxisColor: "#40916C",
		XYChartGridColor: "#D8F3DC",

		RadarCurveColors: []string{
			"#2D6A4F", "#E76F51", "#52B788", "#DDA15E",
			"#40916C", "#BC6C25", "#606C38", "#283618",
		},
		RadarAxisColor:      "#40916C",
		RadarGraticuleColor: "#D8F3DC",
		RadarCurveOpacity:   RadarCurveOpacityDefault,

		MindmapBranchColors: []string{
			"#2D6A4F", "#52B788", "#DDA15E", "#BC6C25",
			"#E76F51", "#606C38", "#40916C", "#95D5B2",
		},
		MindmapNodeFill:   "#D8F3DC",
		MindmapNodeBorder: "#1B4332",

		SankeyNodeColors: []string{
			"#2D6A4F", "#52B788", "#DDA15E", "#BC6C25",
			"#E76F51", "#606C38", "#283618", "#95D5B2",
		},
		SankeyLinkColor:   "#40916C",
		SankeyLinkOpacity: SankeyLinkOpacityLight,

		TreemapColors: []string{
			"#2D6A4F", "#52B788", "#DDA15E", "#BC6C25",
			"#E76F51", "#606C38", "#283618", "#95D5B2",
		},
		TreemapBorder:    "#1B4332",
		TreemapTextColor: "#FFFFFF",

		RequirementFill:   "#D8F3DC",
		RequirementBorder: "#1B4332",

		BlockColors: []string{
			"#2D6A4F", "#52B788", "#DDA15E", "#BC6C25",
			"#E76F51", "#606C38", "#40916C", "#95D5B2",
		},
		BlockNodeBorder: "#1B4332",

		JourneySectionColors: []string{"#D8F3DC", "#FEFAE0", "#B7E4C7", "#95D5B2", "#D4E09B"},
		JourneyTaskFill:      "#FFFFFF",
		JourneyTaskBorder:    "#74C69D",
		JourneyTaskText:      "#1B4332",
		JourneyScoreColors:   [5]string{"#E76F51", "#DDA15E", "#FEFAE0", "#95D5B2", "#2D6A4F"},

		ArchServiceFill:   "#D8F3DC",
		ArchServiceBorder: "#1B4332",
		ArchServiceText:   "#1B4332",
		ArchGroupFill:     "#F0F9F0",
		ArchGroupBorder:   "#74C69D",
		ArchGroupText:     "#2D6A4F",
		ArchEdgeColor:     "#40916C",
		ArchJunctionFill:  "#40916C",

		C4PersonColor:    "#1B4332",
		C4SystemColor:    "#2D6A4F",
		C4ContainerColor: "#52B788",
		C4ComponentColor: "#95D5B2",
		C4ExternalColor:  "#74C69D",
		C4BoundaryColor:  "#40916C",
		C4TextColor:      "#FFFFFF",
	}
}

// Neutral returns a desaturated, accessibility-focused theme with high contrast.
//
//nolint:dupl,funlen // each theme preset is intentionally distinct; flat declaration
func Neutral() *Theme {
	return &Theme{
		FontFamily: "Inter, sans-serif",
		FontSize:   FontSizeDefault,
		Background: "#FFFFFF",

		PrimaryColor:       "#5D6D7E",
		PrimaryBorderColor: "#4A5568",
		PrimaryTextColor:   "#2D3748",

		SecondaryColor:       "#A0AEC0",
		SecondaryBorderColor: "#718096",
		SecondaryTextColor:   "#2D3748",

		TertiaryColor:       "#E2E8F0",
		TertiaryBorderColor: "#CBD5E0",

		LineColor: "#4A5568",
		TextColor: "#2D3748",

		ClusterBackground: "#EDF2F7",
		ClusterBorder:     "#A0AEC0",
		NodeBorderColor:   "#4A5568",

		NoteBackground:  "#FEFCBF",
		NoteBorderColor: "#ECC94B",
		NoteTextColor:   "#744210",

		ActorBorder:     "#4A5568",
		ActorBackground: "#5D6D7E",
		ActorTextColor:  "#FFFFFF",
		ActorLineColor:  "#718096",

		SignalColor:     "#718096",
		SignalTextColor: "#2D3748",

		ActivationBorderColor: "#4A5568",
		ActivationBackground:  "#EDF2F7",
		SequenceNumberColor:   "#FFFFFF",

		EdgeLabelBackground: "#FFFFFF",
		LabelTextColor:      "#2D3748",
		LoopTextColor:       "#2D3748",

		PieTitleTextSize:    PieTitleSizeDefault,
		PieTitleTextColor:   "#2D3748",
		PieSectionTextSize:  PieSectionSizeDefault,
		PieSectionTextColor: "#FFFFFF",
		PieStrokeColor:      "#FFFFFF",
		PieStrokeWidth:      PieStrokeWidthDefault,
		PieOuterStrokeWidth: PieStrokeWidthDefault,
		PieOuterStrokeColor: "#4A5568",
		PieOpacity:          PieOpacityDefault,
		PieColors: []string{
			"#5D6D7E", "#A0AEC0", "#718096", "#4A5568",
			"#2D3748", "#CBD5E0", "#E2E8F0", "#1A202C",
		},

		QuadrantFill1:     "#EDF2F7",
		QuadrantFill2:     "#E2E8F0",
		QuadrantFill3:     "#F7FAFC",
		QuadrantFill4:     "#FEFCBF",
		QuadrantPointFill: "#5D6D7E",

		ClassHeaderBg: "#5D6D7E",
		ClassBodyBg:   "#EDF2F7",
		ClassBorder:   "#4A5568",

		StateFill:         "#EDF2F7",
		StateBorder:       "#4A5568",
		StateStartEnd:     "#2D3748",
		CompositeHeaderBg: "#E2E8F0",

		EntityHeaderBg: "#5D6D7E",
		EntityBodyBg:   "#EDF2F7",
		EntityBorder:   "#4A5568",

		TimelineSectionColors: []string{"#EDF2F7", "#E2E8F0", "#F7FAFC", "#FEFCBF"},
		TimelineEventFill:     "#5D6D7E",
		TimelineEventBorder:   "#4A5568",

		GanttTaskFill:         "#5D6D7E",
		GanttTaskBorder:       "#4A5568",
		GanttCritFill:         "#E53E3E",
		GanttCritBorder:       "#C53030",
		GanttDoneFill:         "#CBD5E0",
		GanttActiveFill:       "#A0AEC0",
		GanttMilestoneFill:    "#718096",
		GanttGridColor:        "#E2E8F0",
		GanttTodayMarkerColor: "#E53E3E",
		GanttSectionColors:    []string{"#EDF2F7", "#F7FAFC", "#E2E8F0", "#FEFCBF"},

		GitBranchColors: []string{
			"#5D6D7E", "#E53E3E", "#718096", "#A0AEC0",
			"#4A5568", "#2D3748", "#CBD5E0", "#1A202C",
		},
		GitCommitFill:    "#2D3748",
		GitCommitStroke:  "#2D3748",
		GitTagFill:       "#FEFCBF",
		GitTagBorder:     "#ECC94B",
		GitHighlightFill: "#E53E3E",

		XYChartColors: []string{
			"#5D6D7E", "#A0AEC0", "#718096", "#4A5568",
			"#2D3748", "#CBD5E0", "#E2E8F0", "#1A202C",
		},
		XYChartAxisColor: "#718096",
		XYChartGridColor: "#E2E8F0",

		RadarCurveColors: []string{
			"#5D6D7E", "#E53E3E", "#718096", "#A0AEC0",
			"#4A5568", "#2D3748", "#CBD5E0", "#1A202C",
		},
		RadarAxisColor:      "#718096",
		RadarGraticuleColor: "#E2E8F0",
		RadarCurveOpacity:   RadarCurveOpacityDefault,

		MindmapBranchColors: []string{
			"#5D6D7E", "#A0AEC0", "#718096", "#4A5568",
			"#2D3748", "#CBD5E0", "#E2E8F0", "#1A202C",
		},
		MindmapNodeFill:   "#EDF2F7",
		MindmapNodeBorder: "#4A5568",

		SankeyNodeColors: []string{
			"#5D6D7E", "#A0AEC0", "#718096", "#4A5568",
			"#2D3748", "#CBD5E0", "#E2E8F0", "#1A202C",
		},
		SankeyLinkColor:   "#718096",
		SankeyLinkOpacity: SankeyLinkOpacityLight,

		TreemapColors: []string{
			"#5D6D7E", "#A0AEC0", "#718096", "#4A5568",
			"#2D3748", "#CBD5E0", "#E2E8F0", "#1A202C",
		},
		TreemapBorder:    "#4A5568",
		TreemapTextColor: "#FFFFFF",

		RequirementFill:   "#EDF2F7",
		RequirementBorder: "#4A5568",

		BlockColors: []string{
			"#5D6D7E", "#A0AEC0", "#718096", "#4A5568",
			"#2D3748", "#CBD5E0", "#E2E8F0", "#1A202C",
		},
		BlockNodeBorder: "#4A5568",

		JourneySectionColors: []string{"#EDF2F7", "#E2E8F0", "#F7FAFC", "#FEFCBF", "#F0FFF4"},
		JourneyTaskFill:      "#FFFFFF",
		JourneyTaskBorder:    "#A0AEC0",
		JourneyTaskText:      "#2D3748",
		JourneyScoreColors:   [5]string{"#E53E3E", "#DD6B20", "#ECC94B", "#48BB78", "#38A169"},

		ArchServiceFill:   "#EDF2F7",
		ArchServiceBorder: "#4A5568",
		ArchServiceText:   "#2D3748",
		ArchGroupFill:     "#F7FAFC",
		ArchGroupBorder:   "#A0AEC0",
		ArchGroupText:     "#4A5568",
		ArchEdgeColor:     "#718096",
		ArchJunctionFill:  "#718096",

		C4PersonColor:    "#2D3748",
		C4SystemColor:    "#5D6D7E",
		C4ContainerColor: "#A0AEC0",
		C4ComponentColor: "#CBD5E0",
		C4ExternalColor:  "#718096",
		C4BoundaryColor:  "#4A5568",
		C4TextColor:      "#FFFFFF",
	}
}
