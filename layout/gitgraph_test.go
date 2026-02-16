package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestGitGraphLayout(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"
	g.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1"},
		&ir.GitBranch{Name: "develop"},
		&ir.GitCheckout{Branch: "develop"},
		&ir.GitCommit{ID: "c2"},
		&ir.GitCheckout{Branch: "main"},
		&ir.GitMerge{Branch: "develop", ID: "m1"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)

	if l.Kind != ir.GitGraph {
		t.Errorf("Kind = %v, want GitGraph", l.Kind)
	}
	if l.Width <= 0 || l.Height <= 0 {
		t.Errorf("dimensions = %f x %f", l.Width, l.Height)
	}

	ggd, ok := l.Diagram.(GitGraphData)
	if !ok {
		t.Fatalf("Diagram type = %T, want GitGraphData", l.Diagram)
	}
	if len(ggd.Commits) < 3 {
		t.Errorf("Commits = %d, want >= 3", len(ggd.Commits))
	}
	if len(ggd.Branches) < 2 {
		t.Errorf("Branches = %d, want >= 2", len(ggd.Branches))
	}
}

func TestGitGraphLayoutBranchLanes(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"
	g.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1"},
		&ir.GitBranch{Name: "dev"},
		&ir.GitCheckout{Branch: "dev"},
		&ir.GitCommit{ID: "c2"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)
	ggd := l.Diagram.(GitGraphData)

	// main and dev should be on different Y lanes.
	var mainY, devY float32
	for _, br := range ggd.Branches {
		if br.Name == "main" {
			mainY = br.Y
		}
		if br.Name == "dev" {
			devY = br.Y
		}
	}
	if mainY == devY {
		t.Errorf("main.Y=%f == dev.Y=%f, want different lanes", mainY, devY)
	}
}

func TestGitGraphLayoutCommitOrder(t *testing.T) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"
	g.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1"},
		&ir.GitCommit{ID: "c2"},
		&ir.GitCommit{ID: "c3"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	l := ComputeLayout(g, th, cfg)
	ggd := l.Diagram.(GitGraphData)

	// Commits should be left-to-right.
	if len(ggd.Commits) != 3 {
		t.Fatalf("Commits = %d, want 3", len(ggd.Commits))
	}
	if ggd.Commits[0].X >= ggd.Commits[1].X || ggd.Commits[1].X >= ggd.Commits[2].X {
		t.Errorf("commits not left-to-right: %f, %f, %f",
			ggd.Commits[0].X, ggd.Commits[1].X, ggd.Commits[2].X)
	}
}
