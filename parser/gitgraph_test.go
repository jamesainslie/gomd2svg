package parser

import (
	"testing"

	"github.com/jamesainslie/gomd2svg/ir"
)

func TestParseGitGraphBasic(t *testing.T) {
	//nolint:dupword // mermaid syntax
	input := `gitGraph
    commit
    commit id: "feat1"
    branch develop
    checkout develop
    commit
    checkout main
    merge develop`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	graph := out.Graph
	if graph.Kind != ir.GitGraph {
		t.Errorf("Kind = %v, want GitGraph", graph.Kind)
	}
	if len(graph.GitActions) != 7 {
		t.Fatalf("Actions = %d, want 7", len(graph.GitActions))
	}

	// First action should be a commit.
	if _, ok := graph.GitActions[0].(*ir.GitCommit); !ok {
		t.Errorf("Action[0] type = %T, want *GitCommit", graph.GitActions[0])
	}
	// Second commit has ID.
	c1, ok := graph.GitActions[1].(*ir.GitCommit)
	if !ok {
		t.Fatal("expected *ir.GitCommit")
	}
	if c1.ID != "feat1" {
		t.Errorf("Action[1].ID = %q, want feat1", c1.ID)
	}
	// Branch.
	br, ok := graph.GitActions[2].(*ir.GitBranch)
	if !ok {
		t.Fatal("expected *ir.GitBranch")
	}
	if br.Name != "develop" {
		t.Errorf("Branch.Name = %q", br.Name)
	}
	// Merge.
	mg, ok := graph.GitActions[6].(*ir.GitMerge)
	if !ok {
		t.Fatal("expected *ir.GitMerge")
	}
	if mg.Branch != "develop" {
		t.Errorf("Merge.Branch = %q", mg.Branch)
	}
}

func TestParseGitGraphCommitOptions(t *testing.T) {
	input := `gitGraph
    commit id: "c1" tag: "v1.0" type: HIGHLIGHT
    commit type: REVERSE`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	c0, ok := out.Graph.GitActions[0].(*ir.GitCommit)
	if !ok {
		t.Fatal("expected *ir.GitCommit")
	}
	if c0.ID != "c1" || c0.Tag != "v1.0" || c0.Type != ir.GitCommitHighlight {
		t.Errorf("commit[0] = %+v", c0)
	}
	c1, ok := out.Graph.GitActions[1].(*ir.GitCommit)
	if !ok {
		t.Fatal("expected *ir.GitCommit")
	}
	if c1.Type != ir.GitCommitReverse {
		t.Errorf("commit[1].Type = %v, want REVERSE", c1.Type)
	}
}

func TestParseGitGraphBranchOrder(t *testing.T) {
	input := `gitGraph
    commit
    branch develop order: 2
    branch feature order: 1`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	b0, ok := out.Graph.GitActions[1].(*ir.GitBranch)
	if !ok {
		t.Fatal("expected *ir.GitBranch")
	}
	if b0.Name != "develop" || b0.Order != 2 {
		t.Errorf("branch[0] = %+v", b0)
	}
	b1, ok := out.Graph.GitActions[2].(*ir.GitBranch)
	if !ok {
		t.Fatal("expected *ir.GitBranch")
	}
	if b1.Name != "feature" || b1.Order != 1 {
		t.Errorf("branch[1] = %+v", b1)
	}
}

func TestParseGitGraphMergeOptions(t *testing.T) {
	input := `gitGraph
    commit
    branch dev
    checkout dev
    commit
    checkout main
    merge dev id: "m1" tag: "v2.0" type: HIGHLIGHT`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	mg, ok := out.Graph.GitActions[5].(*ir.GitMerge)
	if !ok {
		t.Fatal("expected *ir.GitMerge")
	}
	if mg.Branch != "dev" || mg.ID != "m1" || mg.Tag != "v2.0" || mg.Type != ir.GitCommitHighlight {
		t.Errorf("merge = %+v", mg)
	}
}

func TestParseGitGraphCherryPick(t *testing.T) {
	input := `gitGraph
    commit id: "src"
    branch dev
    checkout dev
    cherry-pick id: "src"`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	cp, ok := out.Graph.GitActions[3].(*ir.GitCherryPick)
	if !ok {
		t.Fatal("expected *ir.GitCherryPick")
	}
	if cp.ID != "src" {
		t.Errorf("cherry-pick.ID = %q", cp.ID)
	}
}

func TestParseGitGraphSwitch(t *testing.T) {
	input := `gitGraph
    commit
    branch dev
    switch dev
    commit`

	out, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	co, ok := out.Graph.GitActions[2].(*ir.GitCheckout)
	if !ok {
		t.Fatal("expected *ir.GitCheckout")
	}
	if co.Branch != "dev" {
		t.Errorf("switch.Branch = %q", co.Branch)
	}
}
