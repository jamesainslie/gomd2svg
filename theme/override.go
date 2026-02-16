package theme

// Overrides specifies selective theme field overrides.
// Non-nil pointer fields replace the corresponding field in the base theme.
type Overrides struct {
	FontFamily *string
	FontSize   *float32
	Background *string

	PrimaryColor       *string
	PrimaryBorderColor *string
	PrimaryTextColor   *string

	SecondaryColor       *string
	SecondaryBorderColor *string
	SecondaryTextColor   *string

	TertiaryColor       *string
	TertiaryBorderColor *string

	LineColor *string
	TextColor *string

	ClusterBackground *string
	ClusterBorder     *string
	NodeBorderColor   *string

	NoteBackground  *string
	NoteBorderColor *string
	NoteTextColor   *string
}

// WithOverrides creates a new Theme by copying base and applying non-nil overrides.
// If base is nil, Modern() is used.
func WithOverrides(base *Theme, overrides Overrides) *Theme {
	if base == nil {
		base = Modern()
	}
	// Shallow copy the base theme.
	result := *base
	// Deep-copy slice fields to prevent mutation.
	result.PieColors = copyStrings(base.PieColors)
	result.TimelineSectionColors = copyStrings(base.TimelineSectionColors)
	result.GanttSectionColors = copyStrings(base.GanttSectionColors)
	result.GitBranchColors = copyStrings(base.GitBranchColors)
	result.XYChartColors = copyStrings(base.XYChartColors)
	result.RadarCurveColors = copyStrings(base.RadarCurveColors)
	result.MindmapBranchColors = copyStrings(base.MindmapBranchColors)
	result.SankeyNodeColors = copyStrings(base.SankeyNodeColors)
	result.TreemapColors = copyStrings(base.TreemapColors)
	result.BlockColors = copyStrings(base.BlockColors)
	result.JourneySectionColors = copyStrings(base.JourneySectionColors)

	// Apply non-nil overrides.
	if overrides.FontFamily != nil {
		result.FontFamily = *overrides.FontFamily
	}
	if overrides.FontSize != nil {
		result.FontSize = *overrides.FontSize
	}
	if overrides.Background != nil {
		result.Background = *overrides.Background
	}
	if overrides.PrimaryColor != nil {
		result.PrimaryColor = *overrides.PrimaryColor
	}
	if overrides.PrimaryBorderColor != nil {
		result.PrimaryBorderColor = *overrides.PrimaryBorderColor
	}
	if overrides.PrimaryTextColor != nil {
		result.PrimaryTextColor = *overrides.PrimaryTextColor
	}
	if overrides.SecondaryColor != nil {
		result.SecondaryColor = *overrides.SecondaryColor
	}
	if overrides.SecondaryBorderColor != nil {
		result.SecondaryBorderColor = *overrides.SecondaryBorderColor
	}
	if overrides.SecondaryTextColor != nil {
		result.SecondaryTextColor = *overrides.SecondaryTextColor
	}
	if overrides.TertiaryColor != nil {
		result.TertiaryColor = *overrides.TertiaryColor
	}
	if overrides.TertiaryBorderColor != nil {
		result.TertiaryBorderColor = *overrides.TertiaryBorderColor
	}
	if overrides.LineColor != nil {
		result.LineColor = *overrides.LineColor
	}
	if overrides.TextColor != nil {
		result.TextColor = *overrides.TextColor
	}
	if overrides.ClusterBackground != nil {
		result.ClusterBackground = *overrides.ClusterBackground
	}
	if overrides.ClusterBorder != nil {
		result.ClusterBorder = *overrides.ClusterBorder
	}
	if overrides.NodeBorderColor != nil {
		result.NodeBorderColor = *overrides.NodeBorderColor
	}
	if overrides.NoteBackground != nil {
		result.NoteBackground = *overrides.NoteBackground
	}
	if overrides.NoteBorderColor != nil {
		result.NoteBorderColor = *overrides.NoteBorderColor
	}
	if overrides.NoteTextColor != nil {
		result.NoteTextColor = *overrides.NoteTextColor
	}
	return &result
}

func copyStrings(s []string) []string {
	if s == nil {
		return nil
	}
	cp := make([]string, len(s))
	copy(cp, s)
	return cp
}
