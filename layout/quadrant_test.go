package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestQuadrantLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Quadrant
	graph.QuadrantTitle = "Campaigns"
	graph.XAxisLeft = "Low"
	graph.XAxisRight = "High"
	graph.YAxisBottom = "Low"
	graph.YAxisTop = "High"
	graph.QuadrantLabels = [4]string{"Q1", "Q2", "Q3", "Q4"}
	graph.QuadrantPoints = []*ir.QuadrantPoint{
		{Label: "A", X: 0.0, Y: 0.0},
		{Label: "B", X: 1.0, Y: 1.0},
		{Label: "C", X: 0.5, Y: 0.5},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.Quadrant {
		t.Errorf("Kind = %v, want Quadrant", lay.Kind)
	}
	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("dimensions = %f x %f", lay.Width, lay.Height)
	}

	qd, ok := lay.Diagram.(QuadrantData)
	if !ok {
		t.Fatalf("Diagram type = %T, want QuadrantData", lay.Diagram)
	}
	if len(qd.Points) != 3 {
		t.Fatalf("Points = %d, want 3", len(qd.Points))
	}

	// Point A (0,0) should be at bottom-left; Point B (1,1) at top-right.
	// In pixel space: A.X < B.X and A.Y > B.Y (SVG Y is inverted).
	if qd.Points[0].X >= qd.Points[1].X {
		t.Errorf("A.X=%f >= B.X=%f, want A left of B", qd.Points[0].X, qd.Points[1].X)
	}
	if qd.Points[0].Y <= qd.Points[1].Y {
		t.Errorf("A.Y=%f <= B.Y=%f, want A below B (higher Y)", qd.Points[0].Y, qd.Points[1].Y)
	}

	// Point C (0.5,0.5) should be at center.
	midX := (qd.Points[0].X + qd.Points[1].X) / 2
	midY := (qd.Points[0].Y + qd.Points[1].Y) / 2
	if abs32(qd.Points[2].X-midX) > 1 {
		t.Errorf("C.X=%f not near midX=%f", qd.Points[2].X, midX)
	}
	if abs32(qd.Points[2].Y-midY) > 1 {
		t.Errorf("C.Y=%f not near midY=%f", qd.Points[2].Y, midY)
	}
}

func TestQuadrantLayoutNoPoints(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Quadrant
	graph.QuadrantLabels = [4]string{"Q1", "Q2", "Q3", "Q4"}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	qd, ok := lay.Diagram.(QuadrantData)
	if !ok {
		t.Fatal("expected QuadrantData")
	}
	if len(qd.Points) != 0 {
		t.Errorf("Points = %d, want 0", len(qd.Points))
	}
	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("dimensions = %f x %f", lay.Width, lay.Height)
	}
}

func abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
