package ir

// RadarGraticule distinguishes graticule shapes.
type RadarGraticule int

const (
	RadarGraticuleNone RadarGraticule = iota
	RadarGraticuleCircle
	RadarGraticulePolygon
)

func (g RadarGraticule) String() string {
	switch g {
	case RadarGraticuleNone:
		return "none"
	case RadarGraticuleCircle:
		return "circle"
	case RadarGraticulePolygon:
		return "polygon"
	default:
		return "unknown"
	}
}

// RadarAxis defines one radial axis.
type RadarAxis struct {
	ID    string
	Label string
}

// RadarCurve defines one data series on the radar chart.
type RadarCurve struct {
	ID     string
	Label  string
	Values []float64
}
