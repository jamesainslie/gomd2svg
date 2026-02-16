package layout

import (
	"fmt"
	"sort"

	"github.com/jamesainslie/gomd2svg/config"
	"github.com/jamesainslie/gomd2svg/ir"
	"github.com/jamesainslie/gomd2svg/theme"
)

func computeGitGraphLayout(g *ir.Graph, th *theme.Theme, cfg *config.Layout) *Layout {
	padX := cfg.GitGraph.PaddingX
	padY := cfg.GitGraph.PaddingY
	commitSpacing := cfg.GitGraph.CommitSpacing
	branchSpacing := cfg.GitGraph.BranchSpacing

	mainBranch := g.GitMainBranch
	if mainBranch == "" {
		mainBranch = "main"
	}

	// Simulate git operations to build commit graph.
	type branchInfo struct {
		name  string
		order int
		head  string // latest commit ID on this branch
	}

	branches := map[string]*branchInfo{
		mainBranch: {name: mainBranch, order: 0},
	}
	currentBranch := mainBranch

	type commitInfo struct {
		id     string
		tag    string
		ctype  ir.GitCommitType
		branch string
		seq    int // sequential order
	}

	var commits []commitInfo
	commitMap := make(map[string]int) // commit ID -> index in commits
	type pendingConnection struct {
		fromIdx      int
		toIdx        int
		isCherryPick bool
	}
	var connections []pendingConnection
	autoID := 0

	for _, action := range g.GitActions {
		switch a := action.(type) {
		case *ir.GitCommit:
			id := a.ID
			if id == "" {
				id = fmt.Sprintf("auto_%d", autoID)
				autoID++
			}
			ci := commitInfo{
				id:     id,
				tag:    a.Tag,
				ctype:  a.Type,
				branch: currentBranch,
				seq:    len(commits),
			}
			commitMap[id] = len(commits)
			commits = append(commits, ci)
			branches[currentBranch].head = id

		case *ir.GitBranch:
			order := a.Order
			if order < 0 {
				order = len(branches)
			}
			branches[a.Name] = &branchInfo{
				name:  a.Name,
				order: order,
				head:  branches[currentBranch].head,
			}
			currentBranch = a.Name

		case *ir.GitCheckout:
			currentBranch = a.Branch

		case *ir.GitMerge:
			id := a.ID
			if id == "" {
				id = fmt.Sprintf("merge_%d", autoID)
				autoID++
			}
			ci := commitInfo{
				id:     id,
				tag:    a.Tag,
				ctype:  a.Type,
				branch: currentBranch,
				seq:    len(commits),
			}
			commitMap[id] = len(commits)
			commits = append(commits, ci)

			// Connection from merged branch head to this merge commit.
			if srcBranch, ok := branches[a.Branch]; ok && srcBranch.head != "" {
				if srcIdx, ok2 := commitMap[srcBranch.head]; ok2 {
					connections = append(connections, pendingConnection{
						fromIdx: srcIdx,
						toIdx:   len(commits) - 1,
					})
				}
			}
			branches[currentBranch].head = id

		case *ir.GitCherryPick:
			id := fmt.Sprintf("cp_%d", autoID)
			autoID++
			ci := commitInfo{
				id:     id,
				tag:    a.ID, // show source as tag
				ctype:  ir.GitCommitNormal,
				branch: currentBranch,
				seq:    len(commits),
			}
			commitMap[id] = len(commits)
			commits = append(commits, ci)

			if srcIdx, ok := commitMap[a.ID]; ok {
				connections = append(connections, pendingConnection{
					fromIdx:      srcIdx,
					toIdx:        len(commits) - 1,
					isCherryPick: true,
				})
			}
			branches[currentBranch].head = id
		}
	}

	// Sort branches by order for lane assignment.
	type branchLane struct {
		name  string
		order int
	}
	var sortedBranches []branchLane
	for name, bi := range branches {
		sortedBranches = append(sortedBranches, branchLane{name, bi.order})
	}
	sort.Slice(sortedBranches, func(i, j int) bool {
		return sortedBranches[i].order < sortedBranches[j].order
	})

	branchY := make(map[string]float32)
	for i, bl := range sortedBranches {
		branchY[bl.name] = padY + float32(i)*branchSpacing
	}

	// Position commits.
	var commitLayouts []GitGraphCommitLayout
	for _, ci := range commits {
		x := padX + float32(ci.seq)*commitSpacing
		y := branchY[ci.branch]
		commitLayouts = append(commitLayouts, GitGraphCommitLayout{
			ID:     ci.id,
			Tag:    ci.tag,
			Type:   ci.ctype,
			Branch: ci.branch,
			X:      x,
			Y:      y,
		})
	}

	// Resolve connection pixel positions.
	var connLayouts []GitGraphConnection
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
	var branchLayouts []GitGraphBranchLayout
	for i, bl := range sortedBranches {
		color := "#4A90D9" // fallback
		if len(th.GitBranchColors) > 0 {
			color = th.GitBranchColors[i%len(th.GitBranchColors)]
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

	totalW := padX*2 + float32(len(commits))*commitSpacing
	totalH := padY*2 + float32(len(sortedBranches))*branchSpacing

	return &Layout{
		Kind:   g.Kind,
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
