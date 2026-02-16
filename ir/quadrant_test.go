package ir

import "testing"

func TestQuadrantPointDefaults(t *testing.T) {
	p := &QuadrantPoint{Label: "Campaign A", X: 0.3, Y: 0.6}
	if p.Label != "Campaign A" {
		t.Errorf("Label = %q, want %q", p.Label, "Campaign A")
	}
	if p.X != 0.3 || p.Y != 0.6 {
		t.Errorf("X,Y = %f,%f, want 0.3,0.6", p.X, p.Y)
	}
}

func TestGraphQuadrantFields(t *testing.T) {
	graph := NewGraph()
	graph.Kind = Quadrant
	graph.QuadrantTitle = "Campaigns"
	graph.QuadrantLabels = [4]string{"Expand", "Promote", "Re-evaluate", "Improve"}
	graph.XAxisLeft = "Low Reach"
	graph.XAxisRight = "High Reach"
	graph.YAxisBottom = "Low Engagement"
	graph.YAxisTop = "High Engagement"
	graph.QuadrantPoints = append(graph.QuadrantPoints, &QuadrantPoint{Label: "A", X: 0.3, Y: 0.6})

	if graph.QuadrantTitle != "Campaigns" {
		t.Errorf("QuadrantTitle = %q, want %q", graph.QuadrantTitle, "Campaigns")
	}
	if len(graph.QuadrantPoints) != 1 {
		t.Fatalf("QuadrantPoints = %d, want 1", len(graph.QuadrantPoints))
	}
	if graph.QuadrantLabels[0] != "Expand" {
		t.Errorf("QuadrantLabels[0] = %q, want %q", graph.QuadrantLabels[0], "Expand")
	}
}
