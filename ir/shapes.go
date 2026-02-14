package ir

type NodeShape int

const (
	Rectangle NodeShape = iota
	ForkJoin
	RoundRect
	Stadium
	Subroutine
	Cylinder
	ActorBox
	Circle
	DoubleCircle
	Diamond
	Hexagon
	Parallelogram
	ParallelogramAlt
	Trapezoid
	TrapezoidAlt
	Asymmetric
	MindmapDefault
	Text
)

type EdgeStyle int

const (
	Solid EdgeStyle = iota
	Dotted
	Thick
)

type EdgeDecoration int

const (
	DecCircle EdgeDecoration = iota
	DecCross
	DecDiamond
	DecDiamondFilled
	DecCrowsFootOne
	DecCrowsFootZeroOne
	DecCrowsFootMany
	DecCrowsFootZeroMany
)

type EdgeArrowhead int

const (
	OpenTriangle EdgeArrowhead = iota
	ClassDependency
)
