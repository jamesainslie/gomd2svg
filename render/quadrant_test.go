package render

import (
	"strings"
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/layout"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestRenderQuadrant(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Quadrant
	g.QuadrantTitle = "Campaigns"
	g.XAxisLeft = "Low Reach"
	g.XAxisRight = "High Reach"
	g.YAxisBottom = "Low Engagement"
	g.YAxisTop = "High Engagement"
	g.QuadrantLabels = [4]string{"Expand", "Promote", "Re-evaluate", "Improve"}
	g.QuadrantPoints = []*ir.QuadrantPoint{
		{Label: "Campaign A", X: 0.3, Y: 0.6},
		{Label: "Campaign B", X: 0.7, Y: 0.4},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "Campaigns") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Expand") {
		t.Error("missing quadrant label Q1")
	}
	if !strings.Contains(svg, "Low Reach") {
		t.Error("missing x-axis left label")
	}
	if !strings.Contains(svg, "High Engagement") {
		t.Error("missing y-axis top label")
	}
	if !strings.Contains(svg, "Campaign A") {
		t.Error("missing point label")
	}
	// Should contain circles for data points.
	if !strings.Contains(svg, "<circle") {
		t.Error("missing <circle> for data points")
	}
	// Should contain rects for quadrant backgrounds.
	count := strings.Count(svg, "<rect")
	// 1 background + 4 quadrant rects + 1 border rect = at least 5
	if count < 5 {
		t.Errorf("rect count = %d, want >= 5", count)
	}
}

func TestRenderQuadrantMinimal(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Quadrant
	g.QuadrantPoints = []*ir.QuadrantPoint{
		{Label: "P", X: 0.5, Y: 0.5},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := layout.ComputeLayout(g, th, cfg)
	svg := RenderSVG(l, th, cfg)

	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg> tag")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing data point circle")
	}
}
