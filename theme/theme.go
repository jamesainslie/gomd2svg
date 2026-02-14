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
	}
}
