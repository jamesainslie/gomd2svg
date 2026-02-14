package mermaid

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	svg, err := Render("flowchart LR; A-->B-->C")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("missing </svg>")
	}
}

func TestRenderWithOptions(t *testing.T) {
	opts := Options{}
	svg, err := RenderWithOptions("flowchart TD; X-->Y", opts)
	if err != nil {
		t.Fatalf("RenderWithOptions() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
}

func TestRenderWithTiming(t *testing.T) {
	result, err := RenderWithTiming("flowchart LR; A-->B", Options{})
	if err != nil {
		t.Fatalf("RenderWithTiming() error: %v", err)
	}
	if !strings.Contains(result.SVG, "<svg") {
		t.Error("missing <svg")
	}
	if result.TotalUs() <= 0 {
		t.Error("TotalUs should be > 0")
	}
}

func TestRenderInvalidInput(t *testing.T) {
	_, err := Render("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestRenderContainsNodeLabels(t *testing.T) {
	svg, err := Render("flowchart LR\n  A[Start] --> B[End]")
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "Start") {
		t.Error("missing label 'Start'")
	}
	if !strings.Contains(svg, "End") {
		t.Error("missing label 'End'")
	}
}

func TestGoldenFlowchartSimple(t *testing.T) {
	input := "flowchart LR; A-->B-->C"
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "viewBox") {
		t.Error("missing viewBox")
	}
	if strings.Count(svg, "<rect") < 3 {
		t.Errorf("expected at least 3 rects (nodes), got %d", strings.Count(svg, "<rect"))
	}
	if strings.Count(svg, "edgePath") < 2 {
		t.Errorf("expected at least 2 edge paths, got %d", strings.Count(svg, "edgePath"))
	}
	for _, label := range []string{"A", "B", "C"} {
		if !strings.Contains(svg, ">"+label+"<") {
			t.Errorf("missing node label %q in SVG", label)
		}
	}
}

func TestGoldenFlowchartLabels(t *testing.T) {
	input := "flowchart TD\n    A[Start] --> B{Decision}\n    B -->|Yes| C[OK]\n    B -->|No| D[Cancel]"
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	for _, label := range []string{"Start", "Decision", "OK", "Cancel"} {
		if !strings.Contains(svg, label) {
			t.Errorf("missing label %q in SVG", label)
		}
	}
	if strings.Count(svg, "edgePath") < 3 {
		t.Errorf("expected at least 3 edge paths, got %d", strings.Count(svg, "edgePath"))
	}
}

func TestGoldenFlowchartShapes(t *testing.T) {
	input := "flowchart LR\n    A[Rectangle] --> B(Rounded)\n    B --> C([Stadium])\n    C --> D{Diamond}\n    D --> E((Circle))"
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	// Should have nodes for all 5
	for _, label := range []string{"Rectangle", "Rounded", "Stadium", "Diamond", "Circle"} {
		if !strings.Contains(svg, label) {
			t.Errorf("missing label %q in SVG", label)
		}
	}
	if strings.Count(svg, "edgePath") < 4 {
		t.Errorf("expected at least 4 edge paths, got %d", strings.Count(svg, "edgePath"))
	}
}
