package mermaid

import (
	"os"
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

func TestRenderClassDiagram(t *testing.T) {
	input := `classDiagram
    class Animal {
        +String name
        +isMammal() bool
    }
    class Dog {
        +String breed
    }
    Animal <|-- Dog`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
	for _, label := range []string{"Animal", "Dog", "name", "breed"} {
		if !strings.Contains(svg, label) {
			t.Errorf("missing label %q", label)
		}
	}
}

func TestRenderStateDiagram(t *testing.T) {
	input := `stateDiagram-v2
    [*] --> Still
    Still --> Moving
    Moving --> Crash
    Crash --> [*]`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
	for _, label := range []string{"Still", "Moving", "Crash"} {
		if !strings.Contains(svg, label) {
			t.Errorf("missing label %q", label)
		}
	}
}

func TestRenderERDiagram(t *testing.T) {
	input := `erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE-ITEM : contains`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing <svg")
	}
	for _, label := range []string{"CUSTOMER", "ORDER", "LINE-ITEM"} {
		if !strings.Contains(svg, label) {
			t.Errorf("missing label %q", label)
		}
	}
}

func TestGoldenClassRelationships(t *testing.T) {
	input := `classDiagram
    Animal <|-- Dog : extends
    Car *-- Engine
    Library o-- Book
    Student --> Course
    Class1 ..> Class2
    Interface1 ..|> Impl1`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Count(svg, "edgePath") < 6 {
		t.Errorf("expected at least 6 edge paths, got %d", strings.Count(svg, "edgePath"))
	}
}

func TestGoldenStateComposite(t *testing.T) {
	input := `stateDiagram-v2
    [*] --> First
    state First {
        [*] --> Second
        Second --> [*]
    }
    First --> Third
    Third --> [*]`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "First") {
		t.Error("missing composite state 'First'")
	}
	if !strings.Contains(svg, "Second") {
		t.Error("missing inner state 'Second'")
	}
}

func TestGoldenERAttributes(t *testing.T) {
	input := `erDiagram
    CUSTOMER {
        string name
        int custNumber PK
    }
    ORDER {
        int orderNumber PK
    }
    CUSTOMER ||--o{ ORDER : places`
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(svg, "PK") {
		t.Error("missing PK annotation")
	}
	if !strings.Contains(svg, "places") {
		t.Error("missing relationship label 'places'")
	}
}

func readFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile("testdata/fixtures/" + name)
	if err != nil {
		t.Fatalf("readFixture(%q): %v", name, err)
	}
	return string(data)
}

func TestRenderSequenceDiagram(t *testing.T) {
	input := readFixture(t, "sequence-simple.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(svg, "<svg") {
		t.Error("output should be SVG")
	}
	if !strings.Contains(svg, "Alice") {
		t.Error("SVG should contain Alice")
	}
	if !strings.Contains(svg, "Bob") {
		t.Error("SVG should contain Bob")
	}
}

func TestGoldenSequenceActivations(t *testing.T) {
	input := readFixture(t, "sequence-activations.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "Hello John") {
		t.Error("SVG should contain message text")
	}
}

func TestGoldenSequenceFrames(t *testing.T) {
	input := readFixture(t, "sequence-frames.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "loop") {
		t.Error("SVG should contain loop frame label")
	}
	if !strings.Contains(svg, "alt") {
		t.Error("SVG should contain alt frame label")
	}
}

func TestGoldenSequenceFull(t *testing.T) {
	input := readFixture(t, "sequence-full.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "API Gateway") {
		t.Error("SVG should contain participant alias")
	}
	if !strings.Contains(svg, "Auth flow") {
		t.Error("SVG should contain note text")
	}
}

func TestRenderKanbanDiagram(t *testing.T) {
	input := readFixture(t, "kanban-basic.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Todo") {
		t.Error("expected column label 'Todo'")
	}
	if !strings.Contains(svg, "Write tests") {
		t.Error("expected card label 'Write tests'")
	}
}

func TestRenderKanbanMetadata(t *testing.T) {
	input := readFixture(t, "kanban-metadata.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "Fix login bug") {
		t.Error("expected card label")
	}
}

func TestRenderPacketDiagram(t *testing.T) {
	input := readFixture(t, "packet-tcp.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Source Port") {
		t.Error("expected field label 'Source Port'")
	}
}

func TestRenderPacketBitCount(t *testing.T) {
	input := readFixture(t, "packet-bitcount.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "Checksum") {
		t.Error("expected field label 'Checksum'")
	}
}

func TestRenderPieFixture(t *testing.T) {
	input := readFixture(t, "pie-basic.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Dogs") {
		t.Error("missing Dogs label")
	}
	if !strings.Contains(svg, "<path") {
		t.Error("missing arc path")
	}
}

func TestRenderPieShowDataFixture(t *testing.T) {
	input := readFixture(t, "pie-showdata.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "Calcium") {
		t.Error("missing Calcium label")
	}
	if !strings.Contains(svg, "43") {
		t.Error("missing value display with showData")
	}
}

func TestRenderQuadrantFixture(t *testing.T) {
	input := readFixture(t, "quadrant-campaigns.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Campaign A") {
		t.Error("missing Campaign A")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing data point circles")
	}
	if !strings.Contains(svg, "We should expand") {
		t.Error("missing quadrant label")
	}
}

func TestRenderQuadrantMinimalFixture(t *testing.T) {
	input := readFixture(t, "quadrant-minimal.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "Point A") {
		t.Error("missing Point A")
	}
	if !strings.Contains(svg, "Point B") {
		t.Error("missing Point B")
	}
}

func TestRenderTimelineBasicFixture(t *testing.T) {
	input := readFixture(t, "timeline-basic.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "History of Social Media") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "LinkedIn") {
		t.Error("missing event LinkedIn")
	}
	if !strings.Contains(svg, "Facebook") {
		t.Error("missing event Facebook")
	}
}

func TestRenderTimelineSectionsFixture(t *testing.T) {
	input := readFixture(t, "timeline-sections.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "Research") {
		t.Error("missing event Research")
	}
	if !strings.Contains(svg, "Development") {
		t.Error("missing event Development")
	}
}

func TestRenderGanttBasicFixture(t *testing.T) {
	input := readFixture(t, "gantt-basic.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "A Gantt Diagram") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Research") {
		t.Error("missing task Research")
	}
	if !strings.Contains(svg, "Backend") {
		t.Error("missing task Backend")
	}
}

func TestRenderGanttDependenciesFixture(t *testing.T) {
	input := readFixture(t, "gantt-dependencies.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "Sprint Plan") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "Implement") {
		t.Error("missing task Implement")
	}
}

func TestRenderGitGraphBasicFixture(t *testing.T) {
	input := readFixture(t, "gitgraph-basic.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "v1.0") {
		t.Error("missing tag v1.0")
	}
	if !strings.Contains(svg, "<circle") {
		t.Error("missing commit circles")
	}
	if !strings.Contains(svg, "main") {
		t.Error("missing branch label main")
	}
}

func TestRenderGitGraphBranchesFixture(t *testing.T) {
	input := readFixture(t, "gitgraph-branches.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "v2.0") {
		t.Error("missing tag v2.0")
	}
	if !strings.Contains(svg, "feature") {
		t.Error("missing branch label feature")
	}
	if !strings.Contains(svg, "bugfix") {
		t.Error("missing branch label bugfix")
	}
}

func TestRenderXYChartFixture(t *testing.T) {
	input := readFixture(t, "xychart-basic.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Sales Revenue") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "<rect") {
		t.Error("missing bar rects")
	}
	if !strings.Contains(svg, "<polyline") {
		t.Error("missing line polyline")
	}
}

func TestRenderXYChartHorizontalFixture(t *testing.T) {
	input := readFixture(t, "xychart-horizontal.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Performance") {
		t.Error("missing title")
	}
}

func TestRenderRadarFixture(t *testing.T) {
	input := readFixture(t, "radar-basic.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Language Skills") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "<polygon") {
		t.Error("missing curve polygons")
	}
	if !strings.Contains(svg, "<line") {
		t.Error("missing axis lines")
	}
}

func TestRenderRadarPolygonFixture(t *testing.T) {
	input := readFixture(t, "radar-polygon.mmd")
	svg, err := Render(input)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("expected SVG output")
	}
	if !strings.Contains(svg, "Team Comparison") {
		t.Error("missing title")
	}
	if !strings.Contains(svg, "<polygon") {
		t.Error("missing polygon graticule/curves")
	}
}
