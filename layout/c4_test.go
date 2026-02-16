package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestC4Layout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.C4
	graph.C4SubKind = ir.C4Context

	userLabel := "User"
	sysLabel := "System"
	graph.EnsureNode("user", &userLabel, nil)
	graph.EnsureNode("sys", &sysLabel, nil)

	graph.C4Elements = append(graph.C4Elements,
		&ir.C4Element{ID: "user", Label: "User", Type: ir.C4Person},
		&ir.C4Element{ID: "sys", Label: "System", Type: ir.C4System},
	)
	relLabel := "Uses"
	graph.Edges = append(graph.Edges, &ir.Edge{From: "user", To: "sys", Label: &relLabel, Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := computeC4Layout(graph, th, cfg)

	if lay.Kind != ir.C4 {
		t.Fatalf("Kind = %v", lay.Kind)
	}
	cd, ok := lay.Diagram.(C4Data)
	if !ok {
		t.Fatal("Diagram is not C4Data")
	}
	if len(cd.Elements) != 2 {
		t.Errorf("Elements = %d", len(cd.Elements))
	}
	if cd.SubKind != ir.C4Context {
		t.Errorf("SubKind = %v", cd.SubKind)
	}
}

func TestC4LayoutWithBoundary(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.C4
	graph.C4SubKind = ir.C4Container

	apiLabel := "API"
	dbLabel := "DB"
	graph.EnsureNode("api", &apiLabel, nil)
	graph.EnsureNode("db", &dbLabel, nil)

	graph.C4Elements = append(graph.C4Elements,
		&ir.C4Element{ID: "api", Label: "API", Type: ir.C4ContainerPlain, BoundaryID: "sys"},
		&ir.C4Element{ID: "db", Label: "DB", Type: ir.C4ContainerDb, BoundaryID: "sys"},
	)
	graph.C4Boundaries = append(graph.C4Boundaries, &ir.C4Boundary{
		ID: "sys", Label: "My System", Type: "Software System", Children: []string{"api", "db"},
	})
	relLabel := "reads"
	graph.Edges = append(graph.Edges, &ir.Edge{From: "api", To: "db", Label: &relLabel, Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := computeC4Layout(graph, th, cfg)

	cd, ok := lay.Diagram.(C4Data)
	if !ok {
		t.Fatal("Diagram is not C4Data")
	}
	if len(cd.Boundaries) != 1 {
		t.Fatalf("Boundaries = %d, want 1", len(cd.Boundaries))
	}
	b := cd.Boundaries[0]
	if b.Label != "My System" {
		t.Errorf("boundary label = %q", b.Label)
	}
	if b.Width <= 0 || b.Height <= 0 {
		t.Errorf("boundary size = %vx%v", b.Width, b.Height)
	}
}

func TestC4LayoutEmpty(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.C4
	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := computeC4Layout(graph, th, cfg)
	if len(lay.Nodes) != 0 {
		t.Errorf("Nodes = %d", len(lay.Nodes))
	}
}
