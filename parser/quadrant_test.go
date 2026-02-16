package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseQuadrantBasic(t *testing.T) {
	input := `quadrantChart
    title Reach and engagement of campaigns
    x-axis Low Reach --> High Reach
    y-axis Low Engagement --> High Engagement
    quadrant-1 We should expand
    quadrant-2 Need to promote
    quadrant-3 Re-evaluate
    quadrant-4 May be improved
    Campaign A: [0.3, 0.6]
    Campaign B: [0.45, 0.23]
    Campaign C: [0.57, 0.69]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if graph.Kind != ir.Quadrant {
		t.Errorf("Kind = %v, want Quadrant", graph.Kind)
	}
	if graph.QuadrantTitle != "Reach and engagement of campaigns" {
		t.Errorf("Title = %q", graph.QuadrantTitle)
	}
	if graph.XAxisLeft != "Low Reach" || graph.XAxisRight != "High Reach" {
		t.Errorf("XAxis = %q / %q", graph.XAxisLeft, graph.XAxisRight)
	}
	if graph.YAxisBottom != "Low Engagement" || graph.YAxisTop != "High Engagement" {
		t.Errorf("YAxis = %q / %q", graph.YAxisBottom, graph.YAxisTop)
	}
	if graph.QuadrantLabels[0] != "We should expand" {
		t.Errorf("Q1 = %q", graph.QuadrantLabels[0])
	}
	if graph.QuadrantLabels[2] != "Re-evaluate" {
		t.Errorf("Q3 = %q", graph.QuadrantLabels[2])
	}
	if len(graph.QuadrantPoints) != 3 {
		t.Fatalf("Points = %d, want 3", len(graph.QuadrantPoints))
	}
	p := graph.QuadrantPoints[0]
	if p.Label != "Campaign A" || p.X != 0.3 || p.Y != 0.6 {
		t.Errorf("point[0] = %+v", p)
	}
}

func TestParseQuadrantSingleAxisLabel(t *testing.T) {
	input := `quadrantChart
    x-axis Effort
    y-axis Impact
    Task A: [0.5, 0.5]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if graph.XAxisLeft != "Effort" || graph.XAxisRight != "" {
		t.Errorf("XAxis = %q / %q", graph.XAxisLeft, graph.XAxisRight)
	}
	if graph.YAxisBottom != "Impact" || graph.YAxisTop != "" {
		t.Errorf("YAxis = %q / %q", graph.YAxisBottom, graph.YAxisTop)
	}
}

func TestParseQuadrantMinimal(t *testing.T) {
	input := `quadrantChart
    Point: [0.1, 0.9]`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(out.Graph.QuadrantPoints) != 1 {
		t.Fatalf("Points = %d, want 1", len(out.Graph.QuadrantPoints))
	}
}
