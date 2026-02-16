package ir

// XYSeriesType distinguishes bar and line series.
type XYSeriesType int

const (
	XYSeriesNone XYSeriesType = iota
	XYSeriesBar
	XYSeriesLine
)

func (t XYSeriesType) String() string {
	switch t {
	case XYSeriesNone:
		return "none"
	case XYSeriesBar:
		return "bar"
	case XYSeriesLine:
		return "line"
	default:
		return "unknown"
	}
}

// XYAxisMode distinguishes categorical and numeric axes.
type XYAxisMode int

const (
	XYAxisNone XYAxisMode = iota
	XYAxisBand
	XYAxisNumeric
)

func (m XYAxisMode) String() string {
	switch m {
	case XYAxisNone:
		return "none"
	case XYAxisBand:
		return "band"
	case XYAxisNumeric:
		return "numeric"
	default:
		return "unknown"
	}
}

// XYAxis holds configuration for one axis of an XY chart.
type XYAxis struct {
	Mode       XYAxisMode
	Title      string
	Categories []string // for band axis
	Min        float64  // for numeric axis
	Max        float64  // for numeric axis
}

// XYSeries holds one data series (bar or line).
type XYSeries struct {
	Type   XYSeriesType
	Values []float64
}
