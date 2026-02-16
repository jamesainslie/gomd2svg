package layout

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func TestGitGraphLayout(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.GitGraph
	graph.GitMainBranch = "main"
	graph.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1"},
		&ir.GitBranch{Name: "develop"},
		&ir.GitCheckout{Branch: "develop"},
		&ir.GitCommit{ID: "c2"},
		&ir.GitCheckout{Branch: "main"},
		&ir.GitMerge{Branch: "develop", ID: "m1"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	if lay.Kind != ir.GitGraph {
		t.Errorf("Kind = %v, want GitGraph", lay.Kind)
	}
	if lay.Width <= 0 || lay.Height <= 0 {
		t.Errorf("dimensions = %f x %f", lay.Width, lay.Height)
	}

	ggd, ok := lay.Diagram.(GitGraphData)
	if !ok {
		t.Fatalf("Diagram type = %T, want GitGraphData", lay.Diagram)
	}
	if len(ggd.Commits) < 3 {
		t.Errorf("Commits = %d, want >= 3", len(ggd.Commits))
	}
	if len(ggd.Branches) < 2 {
		t.Errorf("Branches = %d, want >= 2", len(ggd.Branches))
	}
}

func TestGitGraphLayoutBranchLanes(t *testing.T) {
	graph := ir.NewGraph()
	graph.Kind = ir.GitGraph
	graph.GitMainBranch = "main"
	graph.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1"},
		&ir.GitBranch{Name: "dev"},
		&ir.GitCheckout{Branch: "dev"},
		&ir.GitCommit{ID: "c2"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	ggd, ok := lay.Diagram.(GitGraphData)
	if !ok {
		t.Fatal("expected GitGraphData")
	}

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
	graph := ir.NewGraph()
	graph.Kind = ir.GitGraph
	graph.GitMainBranch = "main"
	graph.GitActions = []ir.GitAction{
		&ir.GitCommit{ID: "c1"},
		&ir.GitCommit{ID: "c2"},
		&ir.GitCommit{ID: "c3"},
	}

	th := theme.Modern()
	cfg := config.DefaultLayout()
	lay := ComputeLayout(graph, th, cfg)

	ggd, ok := lay.Diagram.(GitGraphData)
	if !ok {
		t.Fatal("expected GitGraphData")
	}

	// Commits should be left-to-right.
	if len(ggd.Commits) != 3 {
		t.Fatalf("Commits = %d, want 3", len(ggd.Commits))
	}
	if ggd.Commits[0].X >= ggd.Commits[1].X || ggd.Commits[1].X >= ggd.Commits[2].X {
		t.Errorf("commits not left-to-right: %f, %f, %f",
			ggd.Commits[0].X, ggd.Commits[1].X, ggd.Commits[2].X)
	}
}
