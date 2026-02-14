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
	DecNone EdgeDecoration = iota // zero value = no decoration
	DecCircle
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
	OpenTriangle    EdgeArrowhead = iota
	ClassDependency               // dependency
	ClosedTriangle                // inheritance, realization
	FilledDiamond                 // composition
	OpenDiamond                   // aggregation
	Lollipop                      // provided interface
)
