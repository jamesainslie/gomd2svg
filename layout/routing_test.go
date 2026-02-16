package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestBuildObstacleGrid(t *testing.T) {
	nodes := map[string]*NodeLayout{
		"A": {ID: "A", X: 50, Y: 50, Width: 40, Height: 30},
		"B": {ID: "B", X: 150, Y: 50, Width: 40, Height: 30},
	}

	g := buildGrid(nodes, 10, 4)

	// Grid should cover all nodes plus margin.
	if g.cols <= 0 || g.rows <= 0 {
		t.Fatalf("grid has no cells: cols=%d rows=%d", g.cols, g.rows)
	}

	// Center of node A (50,50) should be blocked.
	r, c := g.worldToCell(50, 50)
	if !g.isBlocked(r, c) {
		t.Error("center of node A should be blocked")
	}

	// A point well outside all nodes should be free.
	r, c = g.worldToCell(100, 100)
	if g.isBlocked(r, c) {
		t.Error("point (100,100) should be free, no node there")
	}
}

func TestAStarStraightPath(t *testing.T) {
	// Two nodes side by side in LR, no obstacle between them.
	// A* should find a short path.
	nodes := map[string]*NodeLayout{
		"A": {ID: "A", X: 50, Y: 50, Width: 40, Height: 30},
		"B": {ID: "B", X: 200, Y: 50, Width: 40, Height: 30},
	}

	edges := []*ir.Edge{edge("A", "B")}
	result := routeEdges(edges, nodes, ir.LeftRight)

	if len(result) != 1 {
		t.Fatalf("got %d edges, want 1", len(result))
	}
	pts := result[0].Points
	if len(pts) < 2 {
		t.Fatalf("got %d points, want >= 2", len(pts))
	}
	// First point should be near right side of A.
	startX := pts[0][0]
	if startX < 50 {
		t.Errorf("start X = %f, want >= right side of A (~70)", startX)
	}
	// Last point should be near left side of B.
	endX := pts[len(pts)-1][0]
	if endX > 200 {
		t.Errorf("end X = %f, want <= left side of B (~180)", endX)
	}
}

func TestAStarAvoidsObstacle(t *testing.T) {
	// Three nodes: A on left, C in middle (obstacle), B on right.
	// Edge from A -> B should route around C.
	nodes := map[string]*NodeLayout{
		"A": {ID: "A", X: 50, Y: 50, Width: 40, Height: 30},
		"C": {ID: "C", X: 125, Y: 50, Width: 40, Height: 30},
		"B": {ID: "B", X: 200, Y: 50, Width: 40, Height: 30},
	}

	edges := []*ir.Edge{edge("A", "B")}
	result := routeEdges(edges, nodes, ir.LeftRight)

	if len(result) != 1 {
		t.Fatalf("got %d edges, want 1", len(result))
	}
	pts := result[0].Points
	if len(pts) < 3 {
		t.Fatalf("path has %d points, need >= 3 to route around obstacle", len(pts))
	}

	// No intermediate point should be inside C's padded bounds.
	cLeft := float32(125 - 20 - 4)  // C left edge minus padding
	cRight := float32(125 + 20 + 4) // C right edge plus padding
	cTop := float32(50 - 15 - 4)    // C top minus padding
	cBottom := float32(50 + 15 + 4) // C bottom plus padding
	for _, pt := range pts[1 : len(pts)-1] {
		if pt[0] > cLeft && pt[0] < cRight && pt[1] > cTop && pt[1] < cBottom {
			t.Errorf("path point (%f, %f) inside obstacle C's padded bounds", pt[0], pt[1])
		}
	}

	// The path must deviate in Y: at least one point should differ from Y=50.
	deviated := false
	for _, pt := range pts[1 : len(pts)-1] {
		if pt[1] < cTop || pt[1] > cBottom {
			deviated = true
			break
		}
	}
	if !deviated {
		t.Error("path does not deviate in Y to avoid obstacle C")
	}
}

func TestAStarTopDown(t *testing.T) {
	// TD layout: A on top, B on bottom.
	nodes := map[string]*NodeLayout{
		"A": {ID: "A", X: 50, Y: 30, Width: 40, Height: 30},
		"B": {ID: "B", X: 50, Y: 130, Width: 40, Height: 30},
	}

	edges := []*ir.Edge{edge("A", "B")}
	result := routeEdges(edges, nodes, ir.TopDown)

	if len(result) != 1 {
		t.Fatalf("got %d edges, want 1", len(result))
	}
	pts := result[0].Points
	if len(pts) < 2 {
		t.Fatalf("got %d points, want >= 2", len(pts))
	}
	// Start should be near bottom of A.
	startY := pts[0][1]
	if startY < 30 {
		t.Errorf("start Y = %f, want >= bottom of A (~45)", startY)
	}
	// End should be near top of B.
	endY := pts[len(pts)-1][1]
	if endY > 130 {
		t.Errorf("end Y = %f, want <= top of B (~115)", endY)
	}
}

func TestSimplifyPath(t *testing.T) {
	tests := []struct {
		name string
		in   [][2]float32
		want int // expected number of points after simplification
	}{
		{
			name: "already simple",
			in:   [][2]float32{{0, 0}, {100, 0}},
			want: 2,
		},
		{
			name: "collinear horizontal",
			in:   [][2]float32{{0, 0}, {50, 0}, {100, 0}},
			want: 2, // middle point removed
		},
		{
			name: "collinear vertical",
			in:   [][2]float32{{0, 0}, {0, 50}, {0, 100}},
			want: 2,
		},
		{
			name: "L-shape preserved",
			in:   [][2]float32{{0, 0}, {50, 0}, {50, 50}},
			want: 3, // corner point kept
		},
		{
			name: "five collinear to two",
			in:   [][2]float32{{0, 0}, {25, 0}, {50, 0}, {75, 0}, {100, 0}},
			want: 2,
		},
		{
			name: "staircase",
			in:   [][2]float32{{0, 0}, {50, 0}, {50, 50}, {100, 50}, {100, 100}},
			want: 5, // all corners needed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := simplifyPath(tt.in)
			if len(got) != tt.want {
				t.Errorf("simplifyPath(%v) = %d points, want %d", tt.in, len(got), tt.want)
			}
		})
	}
}

func TestRouteLabelAnchor(t *testing.T) {
	// Edge with a label should have a label anchor near the path midpoint.
	nodes := map[string]*NodeLayout{
		"A": {ID: "A", X: 50, Y: 50, Width: 40, Height: 30,
			Label: TextBlock{FontSize: 14}},
		"B": {ID: "B", X: 200, Y: 50, Width: 40, Height: 30,
			Label: TextBlock{FontSize: 14}},
	}

	lbl := "yes"
	edges := []*ir.Edge{{
		From: "A", To: "B", Directed: true, ArrowEnd: true,
		Style: ir.Solid, Label: &lbl,
	}}
	result := routeEdges(edges, nodes, ir.LeftRight)

	if len(result) != 1 {
		t.Fatalf("got %d edges, want 1", len(result))
	}
	e := result[0]
	if e.Label == nil {
		t.Fatal("edge label is nil, want non-nil")
	}
	// Label anchor should be between the two nodes' X positions.
	if e.LabelAnchor[0] <= 70 || e.LabelAnchor[0] >= 180 {
		t.Errorf("LabelAnchor X = %f, want between ~70 and ~180", e.LabelAnchor[0])
	}
}

func TestRouteEdgesPreservesMetadata(t *testing.T) {
	nodes := map[string]*NodeLayout{
		"A": {ID: "A", X: 50, Y: 50, Width: 40, Height: 30},
		"B": {ID: "B", X: 200, Y: 50, Width: 40, Height: 30},
	}

	arrow := ir.FilledDiamond
	edges := []*ir.Edge{{
		From: "A", To: "B", Directed: true,
		ArrowStart: true, ArrowEnd: true,
		ArrowStartKind: &arrow,
		Style:          ir.Dotted,
	}}
	result := routeEdges(edges, nodes, ir.LeftRight)

	if len(result) != 1 {
		t.Fatalf("got %d edges, want 1", len(result))
	}
	e := result[0]
	if !e.ArrowStart {
		t.Error("ArrowStart should be true")
	}
	if !e.ArrowEnd {
		t.Error("ArrowEnd should be true")
	}
	if e.ArrowStartKind == nil || *e.ArrowStartKind != ir.FilledDiamond {
		t.Error("ArrowStartKind should be FilledDiamond")
	}
	if e.Style != ir.Dotted {
		t.Errorf("Style = %v, want Dotted", e.Style)
	}
}

func TestRouteEdgesMissingNode(t *testing.T) {
	// Edge referencing a missing node should be skipped.
	nodes := map[string]*NodeLayout{
		"A": {ID: "A", X: 50, Y: 50, Width: 40, Height: 30},
	}
	edges := []*ir.Edge{edge("A", "MISSING")}
	result := routeEdges(edges, nodes, ir.LeftRight)
	if len(result) != 0 {
		t.Errorf("got %d edges, want 0 (missing node)", len(result))
	}
}

func TestPathMidpoint(t *testing.T) {
	tests := []struct {
		name         string
		pts          [][2]float32
		wantX, wantY float32
	}{
		{
			name:  "empty",
			pts:   nil,
			wantX: 0, wantY: 0,
		},
		{
			name:  "single point",
			pts:   [][2]float32{{10, 20}},
			wantX: 10, wantY: 20,
		},
		{
			name:  "two points horizontal",
			pts:   [][2]float32{{0, 0}, {100, 0}},
			wantX: 50, wantY: 0,
		},
		{
			name:  "two points vertical",
			pts:   [][2]float32{{0, 0}, {0, 100}},
			wantX: 0, wantY: 50,
		},
		{
			name:  "L-shape midpoint on second segment",
			pts:   [][2]float32{{0, 0}, {100, 0}, {100, 100}},
			wantX: 100, wantY: 0, // midpoint at 100 along a 200-length path
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pathMidpoint(tt.pts)
			const tol = 1.0 // tolerance for float comparison
			dx := got[0] - tt.wantX
			dy := got[1] - tt.wantY
			if dx < -tol || dx > tol || dy < -tol || dy > tol {
				t.Errorf("pathMidpoint = (%f, %f), want (%f, %f)", got[0], got[1], tt.wantX, tt.wantY)
			}
		})
	}
}

func TestBuildGridEmpty(t *testing.T) {
	g := buildGrid(nil, 10, 4)
	if g.rows != 0 || g.cols != 0 {
		t.Errorf("empty grid should have 0 rows/cols, got rows=%d cols=%d", g.rows, g.cols)
	}
	// findPath on empty grid should return nil.
	result := g.findPath(0, 0, 100, 100, "A", "B")
	if result != nil {
		t.Errorf("findPath on empty grid should return nil, got %v", result)
	}
}

func TestGridOverlappingNodePadding(t *testing.T) {
	// Two nodes close enough that their padding regions overlap.
	// A cell in the overlap must allow edges from either node to pass through.
	nodes := map[string]*NodeLayout{
		"A": {ID: "A", X: 50, Y: 50, Width: 40, Height: 30},
		"B": {ID: "B", X: 80, Y: 50, Width: 40, Height: 30}, // 10px gap, but 4px padding each side means overlap
	}

	g := buildGrid(nodes, 8, 4)

	// Find a cell in the overlap zone (between X=66 and X=64: A's right pad edge and B's left pad edge).
	// A right edge + pad = 50+20+4 = 74, B left edge - pad = 80-20-4 = 56
	// So cells between X=56 and X=74 could be in the overlap.
	overlapX := float32(65) // should be in both A and B's padded region
	overlapY := float32(50)
	r, c := g.worldToCell(overlapX, overlapY)

	// Cell should be blocked.
	if !g.isBlocked(r, c) {
		t.Fatal("overlap cell should be blocked")
	}

	// Cell should be passable when excluding A (edge from A goes through).
	if g.isBlockedExcluding(r, c, "A", "X") {
		t.Error("overlap cell should be passable when excluding A")
	}

	// Cell should be passable when excluding B (edge from B goes through).
	if g.isBlockedExcluding(r, c, "X", "B") {
		t.Error("overlap cell should be passable when excluding B")
	}

	// Cell should be blocked when excluding neither A nor B.
	if !g.isBlockedExcluding(r, c, "X", "Y") {
		t.Error("overlap cell should be blocked when excluding unrelated nodes")
	}
}

func TestAStarFallbackToLShape(t *testing.T) {
	// Target node completely surrounded by obstacles. A* should fail and
	// routeEdges should fall back to L-shaped routing.
	nodes := map[string]*NodeLayout{
		"A": {ID: "A", X: 50, Y: 50, Width: 40, Height: 30},
		"B": {ID: "B", X: 150, Y: 50, Width: 40, Height: 30},
		// Surround B with wall nodes on all sides.
		"W1": {ID: "W1", X: 150, Y: 10, Width: 80, Height: 10},
		"W2": {ID: "W2", X: 150, Y: 90, Width: 80, Height: 10},
		"W3": {ID: "W3", X: 110, Y: 50, Width: 10, Height: 80},
		"W4": {ID: "W4", X: 190, Y: 50, Width: 10, Height: 80},
	}

	edges := []*ir.Edge{edge("A", "B")}
	result := routeEdges(edges, nodes, ir.LeftRight)

	if len(result) != 1 {
		t.Fatalf("got %d edges, want 1", len(result))
	}
	// Should still produce a valid path (via fallback).
	pts := result[0].Points
	if len(pts) < 2 {
		t.Errorf("fallback should produce at least 2 points, got %d", len(pts))
	}
}
