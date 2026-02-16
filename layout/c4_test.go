package layout

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/config"
	"github.com/jamesainslie/mermaid-go/ir"
	"github.com/jamesainslie/mermaid-go/theme"
)

func TestC4Layout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.C4
	g.C4SubKind = ir.C4Context

	userLabel := "User"
	sysLabel := "System"
	g.EnsureNode("user", &userLabel, nil)
	g.EnsureNode("sys", &sysLabel, nil)

	g.C4Elements = append(g.C4Elements,
		&ir.C4Element{ID: "user", Label: "User", Type: ir.C4Person},
		&ir.C4Element{ID: "sys", Label: "System", Type: ir.C4System},
	)
	relLabel := "Uses"
	g.Edges = append(g.Edges, &ir.Edge{From: "user", To: "sys", Label: &relLabel, Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeC4Layout(g, th, cfg)

	if l.Kind != ir.C4 {
		t.Fatalf("Kind = %v", l.Kind)
	}
	cd, ok := l.Diagram.(C4Data)
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
	g := ir.NewGraph()
	g.Kind = ir.C4
	g.C4SubKind = ir.C4Container

	apiLabel := "API"
	dbLabel := "DB"
	g.EnsureNode("api", &apiLabel, nil)
	g.EnsureNode("db", &dbLabel, nil)

	g.C4Elements = append(g.C4Elements,
		&ir.C4Element{ID: "api", Label: "API", Type: ir.C4ContainerPlain, BoundaryID: "sys"},
		&ir.C4Element{ID: "db", Label: "DB", Type: ir.C4ContainerDb, BoundaryID: "sys"},
	)
	g.C4Boundaries = append(g.C4Boundaries, &ir.C4Boundary{
		ID: "sys", Label: "My System", Type: "Software System", Children: []string{"api", "db"},
	})
	relLabel := "reads"
	g.Edges = append(g.Edges, &ir.Edge{From: "api", To: "db", Label: &relLabel, Directed: true, ArrowEnd: true})

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeC4Layout(g, th, cfg)

	cd, ok := l.Diagram.(C4Data)
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
	g := ir.NewGraph()
	g.Kind = ir.C4
	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeC4Layout(g, th, cfg)
	if len(l.Nodes) != 0 {
		t.Errorf("Nodes = %d", len(l.Nodes))
	}
}
