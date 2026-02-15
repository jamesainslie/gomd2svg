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
	g := NewGraph()
	g.Kind = Quadrant
	g.QuadrantTitle = "Campaigns"
	g.QuadrantLabels = [4]string{"Expand", "Promote", "Re-evaluate", "Improve"}
	g.XAxisLeft = "Low Reach"
	g.XAxisRight = "High Reach"
	g.YAxisBottom = "Low Engagement"
	g.YAxisTop = "High Engagement"
	g.QuadrantPoints = append(g.QuadrantPoints, &QuadrantPoint{Label: "A", X: 0.3, Y: 0.6})

	if g.QuadrantTitle != "Campaigns" {
		t.Errorf("QuadrantTitle = %q, want %q", g.QuadrantTitle, "Campaigns")
	}
	if len(g.QuadrantPoints) != 1 {
		t.Fatalf("QuadrantPoints = %d, want 1", len(g.QuadrantPoints))
	}
	if g.QuadrantLabels[0] != "Expand" {
		t.Errorf("QuadrantLabels[0] = %q, want %q", g.QuadrantLabels[0], "Expand")
	}
}
