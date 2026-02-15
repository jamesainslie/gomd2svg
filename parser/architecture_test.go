package parser

import (
	"testing"

	"github.com/jamesainslie/mermaid-go/ir"
)

func TestParseArchitecture(t *testing.T) {
	input := `architecture-beta
  group api(cloud)[API]
  service db(database)[Database] in api
  service server(server)[Server] in api
  junction junc1 in api
  db:R -- L:server
  server:R --> L:junc1
`
	out, err := parseArchitecture(input)
	if err != nil {
		t.Fatalf("parseArchitecture() error: %v", err)
	}
	g := out.Graph

	if g.Kind != ir.Architecture {
		t.Errorf("Kind = %v, want Architecture", g.Kind)
	}

	// Groups
	if len(g.ArchGroups) != 1 {
		t.Fatalf("len(ArchGroups) = %d, want 1", len(g.ArchGroups))
	}
	grp := g.ArchGroups[0]
	if grp.ID != "api" {
		t.Errorf("group ID = %q, want %q", grp.ID, "api")
	}
	if grp.Label != "API" {
		t.Errorf("group Label = %q, want %q", grp.Label, "API")
	}
	if grp.Icon != "cloud" {
		t.Errorf("group Icon = %q, want %q", grp.Icon, "cloud")
	}
	if grp.ParentID != "" {
		t.Errorf("group ParentID = %q, want empty", grp.ParentID)
	}

	// Services
	if len(g.ArchServices) != 2 {
		t.Fatalf("len(ArchServices) = %d, want 2", len(g.ArchServices))
	}
	svc0 := g.ArchServices[0]
	if svc0.ID != "db" {
		t.Errorf("service[0] ID = %q, want %q", svc0.ID, "db")
	}
	if svc0.Label != "Database" {
		t.Errorf("service[0] Label = %q, want %q", svc0.Label, "Database")
	}
	if svc0.Icon != "database" {
		t.Errorf("service[0] Icon = %q, want %q", svc0.Icon, "database")
	}
	if svc0.GroupID != "api" {
		t.Errorf("service[0] GroupID = %q, want %q", svc0.GroupID, "api")
	}
	svc1 := g.ArchServices[1]
	if svc1.ID != "server" {
		t.Errorf("service[1] ID = %q, want %q", svc1.ID, "server")
	}
	if svc1.Label != "Server" {
		t.Errorf("service[1] Label = %q, want %q", svc1.Label, "Server")
	}
	if svc1.Icon != "server" {
		t.Errorf("service[1] Icon = %q, want %q", svc1.Icon, "server")
	}
	if svc1.GroupID != "api" {
		t.Errorf("service[1] GroupID = %q, want %q", svc1.GroupID, "api")
	}

	// Junctions
	if len(g.ArchJunctions) != 1 {
		t.Fatalf("len(ArchJunctions) = %d, want 1", len(g.ArchJunctions))
	}
	junc := g.ArchJunctions[0]
	if junc.ID != "junc1" {
		t.Errorf("junction ID = %q, want %q", junc.ID, "junc1")
	}
	if junc.GroupID != "api" {
		t.Errorf("junction GroupID = %q, want %q", junc.GroupID, "api")
	}

	// Edges
	if len(g.ArchEdges) != 2 {
		t.Fatalf("len(ArchEdges) = %d, want 2", len(g.ArchEdges))
	}
	e0 := g.ArchEdges[0]
	if e0.FromID != "db" {
		t.Errorf("edge[0] FromID = %q, want %q", e0.FromID, "db")
	}
	if e0.FromSide != ir.ArchRight {
		t.Errorf("edge[0] FromSide = %v, want ArchRight", e0.FromSide)
	}
	if e0.ToID != "server" {
		t.Errorf("edge[0] ToID = %q, want %q", e0.ToID, "server")
	}
	if e0.ToSide != ir.ArchLeft {
		t.Errorf("edge[0] ToSide = %v, want ArchLeft", e0.ToSide)
	}
	if e0.ArrowLeft {
		t.Error("edge[0] ArrowLeft = true, want false")
	}
	if e0.ArrowRight {
		t.Error("edge[0] ArrowRight = true, want false")
	}

	e1 := g.ArchEdges[1]
	if e1.FromID != "server" {
		t.Errorf("edge[1] FromID = %q, want %q", e1.FromID, "server")
	}
	if !e1.ArrowRight {
		t.Error("edge[1] ArrowRight = false, want true")
	}

	// Group children
	wantChildren := map[string]bool{"db": true, "server": true, "junc1": true}
	for _, child := range grp.Children {
		if !wantChildren[child] {
			t.Errorf("unexpected child %q in group", child)
		}
		delete(wantChildren, child)
	}
	for missing := range wantChildren {
		t.Errorf("missing child %q in group children", missing)
	}

	// Nodes created for services and junctions
	for _, id := range []string{"db", "server", "junc1"} {
		if _, ok := g.Nodes[id]; !ok {
			t.Errorf("Nodes[%q] not found", id)
		}
	}

	// Graph edges for layout
	if len(g.Edges) != 2 {
		t.Fatalf("len(Edges) = %d, want 2", len(g.Edges))
	}
}

func TestParseArchitectureMinimal(t *testing.T) {
	input := `architecture-beta
  service a(server)[ServiceA]
  service b(database)[ServiceB]
  a:R -- L:b
`
	out, err := parseArchitecture(input)
	if err != nil {
		t.Fatalf("parseArchitecture() error: %v", err)
	}
	g := out.Graph

	if g.Kind != ir.Architecture {
		t.Errorf("Kind = %v, want Architecture", g.Kind)
	}
	if len(g.ArchServices) != 2 {
		t.Fatalf("len(ArchServices) = %d, want 2", len(g.ArchServices))
	}
	if len(g.ArchEdges) != 1 {
		t.Fatalf("len(ArchEdges) = %d, want 1", len(g.ArchEdges))
	}

	edge := g.ArchEdges[0]
	if edge.FromID != "a" || edge.ToID != "b" {
		t.Errorf("edge From/To = %q/%q, want a/b", edge.FromID, edge.ToID)
	}
	if edge.ArrowLeft || edge.ArrowRight {
		t.Error("expected no arrows on minimal edge")
	}
}

func TestParseArchitectureNested(t *testing.T) {
	input := `architecture-beta
  group outer(cloud)[Outer]
  group inner(server)[Inner] in outer
  service svc(database)[DB] in inner
`
	out, err := parseArchitecture(input)
	if err != nil {
		t.Fatalf("parseArchitecture() error: %v", err)
	}
	g := out.Graph

	if len(g.ArchGroups) != 2 {
		t.Fatalf("len(ArchGroups) = %d, want 2", len(g.ArchGroups))
	}

	var inner *ir.ArchGroup
	for _, grp := range g.ArchGroups {
		if grp.ID == "inner" {
			inner = grp
		}
	}
	if inner == nil {
		t.Fatal("inner group not found")
	}
	if inner.ParentID != "outer" {
		t.Errorf("inner.ParentID = %q, want %q", inner.ParentID, "outer")
	}

	if len(g.ArchServices) != 1 {
		t.Fatalf("len(ArchServices) = %d, want 1", len(g.ArchServices))
	}
	if g.ArchServices[0].GroupID != "inner" {
		t.Errorf("svc.GroupID = %q, want %q", g.ArchServices[0].GroupID, "inner")
	}

	// Outer group should have "inner" as child
	var outer *ir.ArchGroup
	for _, grp := range g.ArchGroups {
		if grp.ID == "outer" {
			outer = grp
		}
	}
	if outer == nil {
		t.Fatal("outer group not found")
	}
	found := false
	for _, c := range outer.Children {
		if c == "inner" {
			found = true
		}
	}
	if !found {
		t.Error("outer.Children should contain 'inner'")
	}

	// Inner group should have "svc" as child
	found = false
	for _, c := range inner.Children {
		if c == "svc" {
			found = true
		}
	}
	if !found {
		t.Error("inner.Children should contain 'svc'")
	}
}

func TestParseArchitectureBidirectionalArrow(t *testing.T) {
	input := `architecture-beta
  service a(server)[A]
  service b(server)[B]
  a:R <--> L:b
`
	out, err := parseArchitecture(input)
	if err != nil {
		t.Fatalf("parseArchitecture() error: %v", err)
	}
	g := out.Graph

	if len(g.ArchEdges) != 1 {
		t.Fatalf("len(ArchEdges) = %d, want 1", len(g.ArchEdges))
	}
	edge := g.ArchEdges[0]
	if !edge.ArrowLeft {
		t.Error("ArrowLeft = false, want true")
	}
	if !edge.ArrowRight {
		t.Error("ArrowRight = false, want true")
	}
}

func TestParseArchitectureLeftArrow(t *testing.T) {
	input := `architecture-beta
  service a(server)[A]
  service b(server)[B]
  a:R <-- L:b
`
	out, err := parseArchitecture(input)
	if err != nil {
		t.Fatalf("parseArchitecture() error: %v", err)
	}
	g := out.Graph

	if len(g.ArchEdges) != 1 {
		t.Fatalf("len(ArchEdges) = %d, want 1", len(g.ArchEdges))
	}
	edge := g.ArchEdges[0]
	if !edge.ArrowLeft {
		t.Error("ArrowLeft = false, want true")
	}
	if edge.ArrowRight {
		t.Error("ArrowRight = true, want false")
	}
}
