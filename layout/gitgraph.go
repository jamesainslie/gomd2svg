package layout

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

type gitBranchInfo struct {
	name  string
	order int
	head  string // latest commit ID on this branch
}

type gitCommitInfo struct {
	id     string
	tag    string
	ctype  ir.GitCommitType
	branch string
	seq    int // sequential order
}

type gitPendingConnection struct {
	fromIdx      int
	toIdx        int
	isCherryPick bool
}

func computeGitGraphLayout(graph *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	padX := cfg.GitGraph.PaddingX
	padY := cfg.GitGraph.PaddingY
	commitSpacing := cfg.GitGraph.CommitSpacing
	branchSpacing := cfg.GitGraph.BranchSpacing

	mainBranch := graph.GitMainBranch
	if mainBranch == "" {
		mainBranch = "main"
	}

	// Simulate git operations to build commit graph.
	commits, connections, branches := gitGraphSimulate(graph, mainBranch)

	// Sort branches by order for lane assignment.
	sortedBranches := gitGraphSortBranches(branches)

	branchY := make(map[string]float32, len(sortedBranches))
	for idx, bl := range sortedBranches {
		branchY[bl.name] = padY + float32(idx)*branchSpacing
	}

	// Position commits.
	commitLayouts := make([]GitGraphCommitLayout, len(commits))
	for idx, ci := range commits {
		commitLayouts[idx] = GitGraphCommitLayout{
			ID:     ci.id,
			Tag:    ci.tag,
			Type:   ci.ctype,
			Branch: ci.branch,
			X:      padX + float32(ci.seq)*commitSpacing,
			Y:      branchY[ci.branch],
		}
	}

	// Resolve connection pixel positions.
	connLayouts := make([]GitGraphConnection, 0, len(connections))
	for _, conn := range connections {
		if conn.fromIdx < len(commitLayouts) && conn.toIdx < len(commitLayouts) {
			connLayouts = append(connLayouts, GitGraphConnection{
				FromX:        commitLayouts[conn.fromIdx].X,
				FromY:        commitLayouts[conn.fromIdx].Y,
				ToX:          commitLayouts[conn.toIdx].X,
				ToY:          commitLayouts[conn.toIdx].Y,
				IsCherryPick: conn.isCherryPick,
			})
		}
	}

	// Build branch layouts.
	branchLayouts := gitGraphBuildBranchLayouts(sortedBranches, commitLayouts, branchY, th)

	totalW := padX*2 + float32(len(commits))*commitSpacing
	totalH := padY*2 + float32(len(sortedBranches))*branchSpacing

	return &Layout{
		Kind:   graph.Kind,
		Nodes:  map[string]*NodeLayout{},
		Width:  totalW,
		Height: totalH,
		Diagram: GitGraphData{
			Commits:     commitLayouts,
			Branches:    branchLayouts,
			Connections: connLayouts,
		},
	}
}

// gitGraphSimulate processes git actions and builds the commit graph, connections,
// and branch map.
func gitGraphSimulate(graph *ir.Graph, mainBranch string) ([]gitCommitInfo, []gitPendingConnection, map[string]*gitBranchInfo) {
	branches := map[string]*gitBranchInfo{
		mainBranch: {name: mainBranch, order: 0},
	}
	currentBranch := mainBranch

	var commits []gitCommitInfo
	commitMap := make(map[string]int) // commit ID -> index in commits
	var connections []gitPendingConnection
	autoID := 0

	for _, action := range graph.GitActions {
		switch gitAction := action.(type) {
		case *ir.GitCommit:
			id := gitAction.ID
			if id == "" {
				id = fmt.Sprintf("auto_%d", autoID)
				autoID++
			}
			ci := gitCommitInfo{
				id:     id,
				tag:    gitAction.Tag,
				ctype:  gitAction.Type,
				branch: currentBranch,
				seq:    len(commits),
			}
			commitMap[id] = len(commits)
			commits = append(commits, ci)
			branches[currentBranch].head = id

		case *ir.GitBranch:
			order := gitAction.Order
			if order < 0 {
				order = len(branches)
			}
			branches[gitAction.Name] = &gitBranchInfo{
				name:  gitAction.Name,
				order: order,
				head:  branches[currentBranch].head,
			}
			currentBranch = gitAction.Name

		case *ir.GitCheckout:
			currentBranch = gitAction.Branch

		case *ir.GitMerge:
			id := gitAction.ID
			if id == "" {
				id = fmt.Sprintf("merge_%d", autoID)
				autoID++
			}
			ci := gitCommitInfo{
				id:     id,
				tag:    gitAction.Tag,
				ctype:  gitAction.Type,
				branch: currentBranch,
				seq:    len(commits),
			}
			commitMap[id] = len(commits)
			commits = append(commits, ci)

			// Connection from merged branch head to this merge commit.
			if srcBranch, ok := branches[gitAction.Branch]; ok && srcBranch.head != "" {
				if srcIdx, ok2 := commitMap[srcBranch.head]; ok2 {
					connections = append(connections, gitPendingConnection{
						fromIdx: srcIdx,
						toIdx:   len(commits) - 1,
					})
				}
			}
			branches[currentBranch].head = id

		case *ir.GitCherryPick:
			id := fmt.Sprintf("cp_%d", autoID)
			autoID++
			ci := gitCommitInfo{
				id:     id,
				tag:    gitAction.ID, // show source as tag
				ctype:  ir.GitCommitNormal,
				branch: currentBranch,
				seq:    len(commits),
			}
			commitMap[id] = len(commits)
			commits = append(commits, ci)

			if srcIdx, ok := commitMap[gitAction.ID]; ok {
				connections = append(connections, gitPendingConnection{
					fromIdx:      srcIdx,
					toIdx:        len(commits) - 1,
					isCherryPick: true,
				})
			}
			branches[currentBranch].head = id
		}
	}

	return commits, connections, branches
}

type gitBranchLane struct {
	name  string
	order int
}

// gitGraphSortBranches collects and sorts branches by order.
func gitGraphSortBranches(branches map[string]*gitBranchInfo) []gitBranchLane {
	sorted := make([]gitBranchLane, 0, len(branches))
	for name, bi := range branches {
		sorted = append(sorted, gitBranchLane{name, bi.order})
	}
	sort.Slice(sorted, func(idxA, idxB int) bool {
		return sorted[idxA].order < sorted[idxB].order
	})
	return sorted
}

// gitGraphBuildBranchLayouts builds the branch layout objects with colors and
// start/end X positions.
func gitGraphBuildBranchLayouts(sortedBranches []gitBranchLane, commitLayouts []GitGraphCommitLayout, branchY map[string]float32, th *theme.Theme) []GitGraphBranchLayout {
	branchLayouts := make([]GitGraphBranchLayout, 0, len(sortedBranches))
	for idx, bl := range sortedBranches {
		color := "#4A90D9" // fallback
		if len(th.GitBranchColors) > 0 {
			color = th.GitBranchColors[idx%len(th.GitBranchColors)]
		}
		// Find start and end X for this branch.
		var startX, endX float32
		first := true
		for _, cl := range commitLayouts {
			if cl.Branch == bl.name {
				if first || cl.X < startX {
					startX = cl.X
				}
				if first || cl.X > endX {
					endX = cl.X
				}
				first = false
			}
		}
		branchLayouts = append(branchLayouts, GitGraphBranchLayout{
			Name:   bl.name,
			Y:      branchY[bl.name],
			Color:  color,
			StartX: startX,
			EndX:   endX,
		})
	}
	return branchLayouts
}
