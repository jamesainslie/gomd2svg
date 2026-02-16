package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	gitKeyValRe = regexp.MustCompile(`(\w+)\s*:\s*(?:"([^"]+)"|(\S+))`)
)

//nolint:unparam // error return is part of the parser interface contract used by Parse().
func parseGitGraph(input string) (*ParseOutput, error) {
	graph := ir.NewGraph()
	graph.Kind = ir.GitGraph
	graph.GitMainBranch = "main"

	lines := preprocessInput(input)

	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))

		if strings.HasPrefix(lower, "gitgraph") {
			continue
		}

		if strings.HasPrefix(lower, "commit") {
			graph.GitActions = append(graph.GitActions, parseGitCommit(line))
			continue
		}

		if strings.HasPrefix(lower, "branch ") {
			graph.GitActions = append(graph.GitActions, parseGitBranch(line))
			continue
		}

		if strings.HasPrefix(lower, "checkout ") || strings.HasPrefix(lower, "switch ") {
			var rest string
			if strings.HasPrefix(lower, "checkout ") {
				rest = strings.TrimSpace(line[len("checkout "):])
			} else {
				rest = strings.TrimSpace(line[len("switch "):])
			}
			graph.GitActions = append(graph.GitActions, &ir.GitCheckout{Branch: rest})
			continue
		}

		if strings.HasPrefix(lower, "merge ") {
			graph.GitActions = append(graph.GitActions, parseGitMerge(line))
			continue
		}

		if strings.HasPrefix(lower, "cherry-pick") {
			graph.GitActions = append(graph.GitActions, parseGitCherryPick(line))
			continue
		}
	}

	return &ParseOutput{Graph: graph}, nil
}

func parseGitCommit(line string) *ir.GitCommit {
	commit := &ir.GitCommit{}
	opts := gitKeyValRe.FindAllStringSubmatch(line, -1)
	for _, match := range opts {
		key := strings.ToLower(match[1])
		val := match[2]
		if val == "" {
			val = match[3]
		}
		switch key {
		case "id":
			commit.ID = val
		case "tag":
			commit.Tag = val
		case "type":
			commit.Type = parseGitCommitType(val)
		}
	}
	return commit
}

func parseGitBranch(line string) *ir.GitBranch {
	// "branch <name> order: <n>"
	rest := strings.TrimSpace(line[len("branch "):])
	branch := &ir.GitBranch{Order: -1}

	// Check for "order:" option.
	if idx := strings.Index(strings.ToLower(rest), "order:"); idx >= 0 {
		branch.Name = strings.TrimSpace(rest[:idx])
		orderStr := strings.TrimSpace(rest[idx+len("order:"):])
		if orderNum, err := strconv.Atoi(orderStr); err == nil {
			branch.Order = orderNum
		}
	} else {
		branch.Name = rest
	}

	// Strip quotes from name if present.
	branch.Name = strings.Trim(branch.Name, `"`)

	return branch
}

func parseGitMerge(line string) *ir.GitMerge {
	rest := strings.TrimSpace(line[len("merge "):])
	merge := &ir.GitMerge{}

	// Extract the branch name (first word before any key: value pairs).
	parts := strings.Fields(rest)
	if len(parts) > 0 {
		merge.Branch = strings.Trim(parts[0], `"`)
	}

	// Parse key-value options.
	opts := gitKeyValRe.FindAllStringSubmatch(rest, -1)
	for _, match := range opts {
		key := strings.ToLower(match[1])
		val := match[2]
		if val == "" {
			val = match[3]
		}
		switch key {
		case "id":
			merge.ID = val
		case "tag":
			merge.Tag = val
		case "type":
			merge.Type = parseGitCommitType(val)
		}
	}

	return merge
}

func parseGitCherryPick(line string) *ir.GitCherryPick {
	cp := &ir.GitCherryPick{}
	opts := gitKeyValRe.FindAllStringSubmatch(line, -1)
	for _, match := range opts {
		key := strings.ToLower(match[1])
		val := match[2]
		if val == "" {
			val = match[3]
		}
		switch key {
		case "id":
			cp.ID = val
		case "parent":
			cp.Parent = val
		}
	}
	return cp
}

func parseGitCommitType(str string) ir.GitCommitType {
	switch strings.ToUpper(str) {
	case "REVERSE":
		return ir.GitCommitReverse
	case "HIGHLIGHT":
		return ir.GitCommitHighlight
	default:
		return ir.GitCommitNormal
	}
}
