package ir

type NodeStyle struct {
	Fill            *string
	Stroke          *string
	TextColor       *string
	StrokeWidth     *float32
	StrokeDasharray *string
	LineColor       *string
}

type EdgeStyleOverride struct {
	Stroke      *string
	StrokeWidth *float32
	Dasharray   *string
	LabelColor  *string
}

type NodeLink struct {
	URL    string
	Title  *string
	Target *string
}
