package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestTreemapLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Treemap
	graph.TreemapTitle = "Budget"
	graph.TreemapRoot = &ir.TreemapNode{
		Label: "Root",
		Children: []*ir.TreemapNode{
			{Label: "A", Value: 60},
			{Label: "B", Value: 30},
			{Label: "C", Value: 10},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", lay.Width, lay.Height)
	}

	td, ok := lay.Diagram.(TreemapData)
	if !ok {
		t.Fatal("Diagram is not TreemapData")
	}
	if td.Title != "Budget" {
		t.Errorf("title = %q, want Budget", td.Title)
	}
	if len(td.Rects) != 3 {
		t.Fatalf("rects = %d, want 3", len(td.Rects))
	}
	for _, rect := range td.Rects {
		if rect.Width <= 0 || rect.Height <= 0 {
			t.Errorf("rect %q has zero dimension: %v x %v", rect.Label, rect.Width, rect.Height)
		}
	}
}

func TestTreemapLayoutAreaProportional(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Treemap
	graph.TreemapRoot = &ir.TreemapNode{
		Label: "Root",
		Children: []*ir.TreemapNode{
			{Label: "Big", Value: 75},
			{Label: "Small", Value: 25},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	td, ok := lay.Diagram.(TreemapData)
	if !ok {
		t.Fatal("expected TreemapData")
	}
	if len(td.Rects) != 2 {
		t.Fatalf("rects = %d, want 2", len(td.Rects))
	}

	var bigArea, smallArea float64
	for _, rect := range td.Rects {
		area := float64(rect.Width) * float64(rect.Height)
		if rect.Label == "Big" {
			bigArea = area
		} else {
			smallArea = area
		}
	}
	if bigArea <= smallArea {
		t.Errorf("Big area (%f) should be larger than Small area (%f)", bigArea, smallArea)
	}

	// The ratio of areas should approximately reflect the value ratio (3:1),
	// allowing generous tolerance for padding effects.
	ratio := bigArea / smallArea
	if ratio < 1.5 || ratio > 5.0 {
		t.Errorf("area ratio = %f, want roughly 3.0 (within 1.5-5.0)", ratio)
	}
}

func TestTreemapLayoutNested(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Treemap
	graph.TreemapRoot = &ir.TreemapNode{
		Label: "Root",
		Children: []*ir.TreemapNode{
			{Label: "Section", Children: []*ir.TreemapNode{
				{Label: "X", Value: 20},
				{Label: "Y", Value: 30},
			}},
			{Label: "Z", Value: 50},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	td, ok := lay.Diagram.(TreemapData)
	if !ok {
		t.Fatal("Diagram is not TreemapData")
	}
	// Expect: Section (IsSection), X, Y, Z = 4 rects.
	if len(td.Rects) != 4 {
		t.Errorf("rects = %d, want 4", len(td.Rects))
	}

	// Verify the section rect exists and is marked.
	foundSection := false
	for _, rect := range td.Rects {
		if rect.Label == "Section" {
			foundSection = true
			if !rect.IsSection {
				t.Error("Section rect should have IsSection=true")
			}
			if rect.Depth != 0 {
				t.Errorf("Section depth = %d, want 0", rect.Depth)
			}
		}
		if rect.Label == "X" || rect.Label == "Y" {
			if rect.Depth != 1 {
				t.Errorf("leaf %q depth = %d, want 1", rect.Label, rect.Depth)
			}
		}
	}
	if !foundSection {
		t.Error("did not find Section rect")
	}
}

func TestTreemapLayoutEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Treemap

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", lay.Width, lay.Height)
	}

	td, ok := lay.Diagram.(TreemapData)
	if !ok {
		t.Fatal("Diagram is not TreemapData")
	}
	if len(td.Rects) != 0 {
		t.Errorf("rects = %d, want 0", len(td.Rects))
	}
}

func TestTreemapLayoutSingleLeafRoot(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Treemap
	graph.TreemapRoot = &ir.TreemapNode{
		Label: "Only",
		Value: 100,
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	td, ok := lay.Diagram.(TreemapData)
	if !ok {
		t.Fatal("expected TreemapData")
	}
	if len(td.Rects) != 1 {
		t.Fatalf("rects = %d, want 1", len(td.Rects))
	}
	if td.Rects[0].Label != "Only" {
		t.Errorf("label = %q, want Only", td.Rects[0].Label)
	}
	if td.Rects[0].Width <= 0 || td.Rects[0].Height <= 0 {
		t.Errorf("single rect has zero dimension")
	}
}

func TestTreemapLayoutRectsWithinBounds(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.Treemap
	graph.TreemapRoot = &ir.TreemapNode{
		Label: "Root",
		Children: []*ir.TreemapNode{
			{Label: "A", Value: 40},
			{Label: "B", Value: 30},
			{Label: "C", Value: 20},
			{Label: "D", Value: 10},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	td, ok := lay.Diagram.(TreemapData)
	if !ok {
		t.Fatal("expected TreemapData")
	}
	for _, rect := range td.Rects {
		if rect.X < 0 || rect.Y < 0 {
			t.Errorf("rect %q has negative position: (%v, %v)", rect.Label, rect.X, rect.Y)
		}
		if rect.X+rect.Width > lay.Width {
			t.Errorf("rect %q exceeds right bound: %v + %v > %v", rect.Label, rect.X, rect.Width, lay.Width)
		}
		if rect.Y+rect.Height > lay.Height {
			t.Errorf("rect %q exceeds bottom bound: %v + %v > %v", rect.Label, rect.Y, rect.Height, lay.Height)
		}
	}
}

func TestTreemapSquarifyAspectRatios(t *testing.T) {
	// Verify the squarify algorithm produces reasonable aspect ratios.
	items := []treemapItem{
		{value: 6, idx: 0},
		{value: 6, idx: 1},
		{value: 4, idx: 2},
		{value: 3, idx: 3},
		{value: 2, idx: 4},
		{value: 2, idx: 5},
		{value: 1, idx: 6},
	}
	rects := treemapSquarify(items, 0, 0, 600, 400, 24)
	if len(rects) != len(items) {
		t.Fatalf("got %d rects, want %d", len(rects), len(items))
	}
	for _, rect := range rects {
		if rect.w <= 0 || rect.h <= 0 {
			t.Errorf("rect idx %d has zero dimension: %v x %v", rect.item.idx, rect.w, rect.h)
		}
		aspect := float64(rect.w) / float64(rect.h)
		if aspect < 1 {
			aspect = 1 / aspect
		}
		// Squarified algorithm should keep aspect ratios reasonable (< 10:1).
		if aspect > 10 {
			t.Errorf("rect idx %d has bad aspect ratio: %f (w=%f, h=%f)", rect.item.idx, aspect, rect.w, rect.h)
		}
	}
}

func TestTreemapWorstAspect(t *testing.T) {
	// A single square item should have aspect ratio 1.0.
	items := []treemapItem{{value: 100, idx: 0}}
	aspect := treemapWorstAspect(items, 100, 100, 100, 100)
	if aspect != 1.0 {
		t.Errorf("single square aspect = %f, want 1.0", aspect)
	}

	// Empty row should return MaxFloat64.
	aspect = treemapWorstAspect(nil, 0, 100, 100, 100)
	if aspect < 1e300 {
		t.Errorf("empty row aspect = %f, want MaxFloat64", aspect)
	}
}
