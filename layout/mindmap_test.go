package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestMindmapLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap
	g.MindmapRoot = &ir.MindmapNode{
		ID: "root", Label: "Central", Shape: ir.MindmapCircle,
		Children: []*ir.MindmapNode{
			{ID: "a", Label: "Branch A", Shape: ir.MindmapShapeDefault},
			{ID: "b", Label: "Branch B", Shape: ir.MindmapSquare,
				Children: []*ir.MindmapNode{
					{ID: "c", Label: "Leaf C", Shape: ir.MindmapRounded},
				},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %v x %v", l.Width, l.Height)
	}

	md, ok := l.Diagram.(MindmapData)
	if !ok {
		t.Fatal("Diagram is not MindmapData")
	}
	if md.Root == nil {
		t.Fatal("Root is nil")
	}
	if md.Root.Label != "Central" {
		t.Errorf("root label = %q", md.Root.Label)
	}
	if len(md.Root.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(md.Root.Children))
	}
}

func TestMindmapLayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions for empty mindmap: %v x %v", l.Width, l.Height)
	}
}

func TestMindmapLayoutSingleRoot(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap
	g.MindmapRoot = &ir.MindmapNode{
		ID: "root", Label: "Only Root", Shape: ir.MindmapShapeDefault,
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeMindmapLayout(g, th, cfg)

	md, ok := l.Diagram.(MindmapData)
	if !ok {
		t.Fatal("Diagram is not MindmapData")
	}
	if md.Root == nil {
		t.Fatal("Root is nil")
	}
	if md.Root.Label != "Only Root" {
		t.Errorf("root label = %q, want %q", md.Root.Label, "Only Root")
	}
	if len(md.Root.Children) != 0 {
		t.Errorf("root children = %d, want 0", len(md.Root.Children))
	}
	if md.Root.Width <= 0 || md.Root.Height <= 0 {
		t.Errorf("root has zero dimensions: %v x %v", md.Root.Width, md.Root.Height)
	}
}

func TestMindmapLayoutBranchColors(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap
	g.MindmapRoot = &ir.MindmapNode{
		ID: "root", Label: "Root", Shape: ir.MindmapShapeDefault,
		Children: []*ir.MindmapNode{
			{ID: "a", Label: "A", Shape: ir.MindmapShapeDefault,
				Children: []*ir.MindmapNode{
					{ID: "a1", Label: "A1", Shape: ir.MindmapShapeDefault},
				},
			},
			{ID: "b", Label: "B", Shape: ir.MindmapShapeDefault},
			{ID: "c", Label: "C", Shape: ir.MindmapShapeDefault},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeMindmapLayout(g, th, cfg)

	md := l.Diagram.(MindmapData)

	// Root should have ColorIndex 0 (branchIdx passed as 0).
	if md.Root.ColorIndex != 0 {
		t.Errorf("root ColorIndex = %d, want 0", md.Root.ColorIndex)
	}

	// Each top-level child should have a different branch index.
	if md.Root.Children[0].ColorIndex != 0 {
		t.Errorf("child A ColorIndex = %d, want 0", md.Root.Children[0].ColorIndex)
	}
	if md.Root.Children[1].ColorIndex != 1 {
		t.Errorf("child B ColorIndex = %d, want 1", md.Root.Children[1].ColorIndex)
	}
	if md.Root.Children[2].ColorIndex != 2 {
		t.Errorf("child C ColorIndex = %d, want 2", md.Root.Children[2].ColorIndex)
	}

	// Grandchild A1 should inherit branch A's index.
	a1 := md.Root.Children[0].Children[0]
	if a1.ColorIndex != 0 {
		t.Errorf("grandchild A1 ColorIndex = %d, want 0", a1.ColorIndex)
	}
}

func TestMindmapLayoutNoOverlap(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap
	g.MindmapRoot = &ir.MindmapNode{
		ID: "root", Label: "Center", Shape: ir.MindmapCircle,
		Children: []*ir.MindmapNode{
			{ID: "a", Label: "North", Shape: ir.MindmapShapeDefault},
			{ID: "b", Label: "East", Shape: ir.MindmapShapeDefault},
			{ID: "c", Label: "South", Shape: ir.MindmapShapeDefault},
			{ID: "d", Label: "West", Shape: ir.MindmapShapeDefault},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeMindmapLayout(g, th, cfg)

	md := l.Diagram.(MindmapData)

	// All child nodes should be at different positions from root.
	for i, child := range md.Root.Children {
		if child.X == md.Root.X && child.Y == md.Root.Y {
			t.Errorf("child %d overlaps root at (%v, %v)", i, child.X, child.Y)
		}
	}

	// All child nodes should be at different positions from each other.
	for i := 0; i < len(md.Root.Children); i++ {
		for j := i + 1; j < len(md.Root.Children); j++ {
			ci := md.Root.Children[i]
			cj := md.Root.Children[j]
			if ci.X == cj.X && ci.Y == cj.Y {
				t.Errorf("children %d and %d overlap at (%v, %v)", i, j, ci.X, ci.Y)
			}
		}
	}
}

func TestMindmapLayoutPositiveCoordinates(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Mindmap
	g.MindmapRoot = &ir.MindmapNode{
		ID: "root", Label: "Root", Shape: ir.MindmapShapeDefault,
		Children: []*ir.MindmapNode{
			{ID: "a", Label: "Alpha", Shape: ir.MindmapSquare},
			{ID: "b", Label: "Beta", Shape: ir.MindmapRounded,
				Children: []*ir.MindmapNode{
					{ID: "b1", Label: "Beta-1", Shape: ir.MindmapShapeDefault},
					{ID: "b2", Label: "Beta-2", Shape: ir.MindmapShapeDefault},
				},
			},
		},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeMindmapLayout(g, th, cfg)

	md := l.Diagram.(MindmapData)

	// After normalization, all node top-left corners should be non-negative.
	var checkPositive func(n *MindmapNodeLayout, path string)
	checkPositive = func(n *MindmapNodeLayout, path string) {
		left := n.X - n.Width/2
		top := n.Y - n.Height/2
		if left < 0 {
			t.Errorf("%s left edge = %v (negative)", path, left)
		}
		if top < 0 {
			t.Errorf("%s top edge = %v (negative)", path, top)
		}
		for _, child := range n.Children {
			checkPositive(child, path+"/"+child.Label)
		}
	}
	checkPositive(md.Root, "root")
}
