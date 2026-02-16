package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestArchitectureLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Architecture

	// Create 3 services
	dbLabel := "Database"
	srvLabel := "Server"
	webLabel := "WebApp"
	g.EnsureNode("db", &dbLabel, nil)
	g.EnsureNode("server", &srvLabel, nil)
	g.EnsureNode("web", &webLabel, nil)

	g.ArchServices = []*ir.ArchService{
		{ID: "db", Label: "Database", Icon: "database"},
		{ID: "server", Label: "Server", Icon: "server"},
		{ID: "web", Label: "WebApp", Icon: "cloud"},
	}
	g.ArchEdges = []*ir.ArchEdge{
		{FromID: "db", FromSide: ir.ArchRight, ToID: "server", ToSide: ir.ArchLeft},
		{FromID: "server", FromSide: ir.ArchRight, ToID: "web", ToSide: ir.ArchLeft, ArrowRight: true},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeArchitectureLayout(g, th, cfg)

	if l.Kind != ir.Architecture {
		t.Fatalf("Kind = %v, want Architecture", l.Kind)
	}

	// Verify nodes are positioned left-to-right (db -> server -> web)
	dbN := l.Nodes["db"]
	srvN := l.Nodes["server"]
	webN := l.Nodes["web"]
	if dbN == nil || srvN == nil || webN == nil {
		t.Fatal("missing node layouts")
	}
	if dbN.X >= srvN.X {
		t.Errorf("db.X (%f) should be < server.X (%f)", dbN.X, srvN.X)
	}
	if srvN.X >= webN.X {
		t.Errorf("server.X (%f) should be < web.X (%f)", srvN.X, webN.X)
	}

	// All nodes should be on the same row
	if dbN.Y != srvN.Y || srvN.Y != webN.Y {
		t.Errorf("nodes should be on the same row: db.Y=%f, server.Y=%f, web.Y=%f", dbN.Y, srvN.Y, webN.Y)
	}

	// Verify edges
	if len(l.Edges) != 2 {
		t.Fatalf("len(Edges) = %d, want 2", len(l.Edges))
	}

	// Second edge should have ArrowEnd set
	if !l.Edges[1].ArrowEnd {
		t.Error("second edge should have ArrowEnd=true")
	}

	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("invalid dimensions: %f x %f", l.Width, l.Height)
	}
}

func TestArchitectureLayoutVertical(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Architecture

	topLabel := "Frontend"
	midLabel := "API"
	botLabel := "Database"
	g.EnsureNode("top", &topLabel, nil)
	g.EnsureNode("mid", &midLabel, nil)
	g.EnsureNode("bot", &botLabel, nil)

	g.ArchServices = []*ir.ArchService{
		{ID: "top", Label: "Frontend"},
		{ID: "mid", Label: "API"},
		{ID: "bot", Label: "Database"},
	}
	g.ArchEdges = []*ir.ArchEdge{
		{FromID: "top", FromSide: ir.ArchBottom, ToID: "mid", ToSide: ir.ArchTop},
		{FromID: "mid", FromSide: ir.ArchBottom, ToID: "bot", ToSide: ir.ArchTop},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeArchitectureLayout(g, th, cfg)

	topN := l.Nodes["top"]
	midN := l.Nodes["mid"]
	botN := l.Nodes["bot"]
	if topN == nil || midN == nil || botN == nil {
		t.Fatal("missing node layouts")
	}

	// Nodes should be stacked vertically
	if topN.Y >= midN.Y {
		t.Errorf("top.Y (%f) should be < mid.Y (%f)", topN.Y, midN.Y)
	}
	if midN.Y >= botN.Y {
		t.Errorf("mid.Y (%f) should be < bot.Y (%f)", midN.Y, botN.Y)
	}

	// All on the same column
	if topN.X != midN.X || midN.X != botN.X {
		t.Errorf("nodes should be on the same column: top.X=%f, mid.X=%f, bot.X=%f", topN.X, midN.X, botN.X)
	}
}

func TestArchitectureLayoutGroups(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Architecture

	dbLabel := "Database"
	srvLabel := "Server"
	g.EnsureNode("db", &dbLabel, nil)
	g.EnsureNode("srv", &srvLabel, nil)

	g.ArchGroups = []*ir.ArchGroup{
		{ID: "api", Label: "API", Icon: "cloud", Children: []string{"db", "srv"}},
	}
	g.ArchServices = []*ir.ArchService{
		{ID: "db", Label: "Database", GroupID: "api"},
		{ID: "srv", Label: "Server", GroupID: "api"},
	}
	g.ArchEdges = []*ir.ArchEdge{
		{FromID: "db", FromSide: ir.ArchRight, ToID: "srv", ToSide: ir.ArchLeft},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeArchitectureLayout(g, th, cfg)

	data, ok := l.Diagram.(ArchitectureData)
	if !ok {
		t.Fatal("Diagram is not ArchitectureData")
	}
	if len(data.Groups) != 1 {
		t.Fatalf("len(Groups) = %d, want 1", len(data.Groups))
	}
	grp := data.Groups[0]
	if grp.Width <= 0 || grp.Height <= 0 {
		t.Errorf("group has invalid dimensions: %f x %f", grp.Width, grp.Height)
	}
	if grp.Label != "API" {
		t.Errorf("group label = %q, want %q", grp.Label, "API")
	}
	if grp.Icon != "cloud" {
		t.Errorf("group icon = %q, want %q", grp.Icon, "cloud")
	}

	// Group should encompass both nodes
	dbN := l.Nodes["db"]
	srvN := l.Nodes["srv"]
	if dbN.X-dbN.Width/2 < grp.X {
		t.Errorf("db left edge (%f) should be >= group X (%f)", dbN.X-dbN.Width/2, grp.X)
	}
	if srvN.X+srvN.Width/2 > grp.X+grp.Width {
		t.Errorf("srv right edge (%f) should be <= group right (%f)", srvN.X+srvN.Width/2, grp.X+grp.Width)
	}
}

func TestArchitectureLayoutJunctions(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Architecture

	aLabel := "ServiceA"
	bLabel := "ServiceB"
	g.EnsureNode("a", &aLabel, nil)
	g.EnsureNode("b", &bLabel, nil)
	g.EnsureNode("j1", nil, nil)

	g.ArchServices = []*ir.ArchService{
		{ID: "a", Label: "ServiceA"},
		{ID: "b", Label: "ServiceB"},
	}
	g.ArchJunctions = []*ir.ArchJunction{
		{ID: "j1"},
	}
	g.ArchEdges = []*ir.ArchEdge{
		{FromID: "a", FromSide: ir.ArchRight, ToID: "j1", ToSide: ir.ArchLeft},
		{FromID: "j1", FromSide: ir.ArchRight, ToID: "b", ToSide: ir.ArchLeft},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeArchitectureLayout(g, th, cfg)

	data, ok := l.Diagram.(ArchitectureData)
	if !ok {
		t.Fatal("Diagram is not ArchitectureData")
	}
	if len(data.Junctions) != 1 {
		t.Fatalf("len(Junctions) = %d, want 1", len(data.Junctions))
	}
	junc := data.Junctions[0]
	if junc.ID != "j1" {
		t.Errorf("junction ID = %q, want %q", junc.ID, "j1")
	}
	if junc.Size <= 0 {
		t.Errorf("junction size should be > 0, got %f", junc.Size)
	}

	// Junction should be between a and b horizontally
	aN := l.Nodes["a"]
	bN := l.Nodes["b"]
	if aN.X >= junc.X || junc.X >= bN.X {
		t.Errorf("junction X (%f) should be between a.X (%f) and b.X (%f)", junc.X, aN.X, bN.X)
	}
}

func TestArchitectureLayoutEmpty(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.Architecture

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := computeArchitectureLayout(g, th, cfg)

	if l.Kind != ir.Architecture {
		t.Fatalf("Kind = %v, want Architecture", l.Kind)
	}
	data, ok := l.Diagram.(ArchitectureData)
	if !ok {
		t.Fatal("Diagram is not ArchitectureData")
	}
	if len(data.Groups) != 0 {
		t.Errorf("len(Groups) = %d, want 0", len(data.Groups))
	}
	if len(data.Junctions) != 0 {
		t.Errorf("len(Junctions) = %d, want 0", len(data.Junctions))
	}
}

func TestArchSideOffset(t *testing.T) {
	tests := []struct {
		side   ir.ArchSide
		wantDC int
		wantDR int
		name   string
	}{
		{ir.ArchRight, 1, 0, "right"},
		{ir.ArchLeft, -1, 0, "left"},
		{ir.ArchBottom, 0, 1, "bottom"},
		{ir.ArchTop, 0, -1, "top"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc, dr := archSideOffset(tt.side)
			if dc != tt.wantDC || dr != tt.wantDR {
				t.Errorf("archSideOffset(%v) = (%d, %d), want (%d, %d)", tt.side, dc, dr, tt.wantDC, tt.wantDR)
			}
		})
	}
}

func TestArchSideOffsetReverse(t *testing.T) {
	tests := []struct {
		side   ir.ArchSide
		wantDC int
		wantDR int
		name   string
	}{
		{ir.ArchLeft, -1, 0, "left"},
		{ir.ArchRight, 1, 0, "right"},
		{ir.ArchTop, 0, -1, "top"},
		{ir.ArchBottom, 0, 1, "bottom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc, dr := archSideOffsetReverse(tt.side)
			if dc != tt.wantDC || dr != tt.wantDR {
				t.Errorf("archSideOffsetReverse(%v) = (%d, %d), want (%d, %d)", tt.side, dc, dr, tt.wantDC, tt.wantDR)
			}
		})
	}
}

func TestArchAnchorPoint(t *testing.T) {
	n := &NodeLayout{
		X:      100,
		Y:      50,
		Width:  60,
		Height: 40,
	}

	tests := []struct {
		side  ir.ArchSide
		wantX float32
		wantY float32
		name  string
	}{
		{ir.ArchLeft, 70, 50, "left"},
		{ir.ArchRight, 130, 50, "right"},
		{ir.ArchTop, 100, 30, "top"},
		{ir.ArchBottom, 100, 70, "bottom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := archAnchorPoint(n, tt.side)
			if x != tt.wantX || y != tt.wantY {
				t.Errorf("archAnchorPoint(%v) = (%f, %f), want (%f, %f)", tt.side, x, y, tt.wantX, tt.wantY)
			}
		})
	}
}
