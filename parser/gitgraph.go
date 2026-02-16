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

func parseGitGraph(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.GitGraph
	g.GitMainBranch = "main"

	lines := preprocessInput(input)

	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))

		if strings.HasPrefix(lower, "gitgraph") {
			continue
		}

		if strings.HasPrefix(lower, "commit") {
			g.GitActions = append(g.GitActions, parseGitCommit(line))
			continue
		}

		if strings.HasPrefix(lower, "branch ") {
			g.GitActions = append(g.GitActions, parseGitBranch(line))
			continue
		}

		if strings.HasPrefix(lower, "checkout ") || strings.HasPrefix(lower, "switch ") {
			var rest string
			if strings.HasPrefix(lower, "checkout ") {
				rest = strings.TrimSpace(line[len("checkout "):])
			} else {
				rest = strings.TrimSpace(line[len("switch "):])
			}
			g.GitActions = append(g.GitActions, &ir.GitCheckout{Branch: rest})
			continue
		}

		if strings.HasPrefix(lower, "merge ") {
			g.GitActions = append(g.GitActions, parseGitMerge(line))
			continue
		}

		if strings.HasPrefix(lower, "cherry-pick") {
			g.GitActions = append(g.GitActions, parseGitCherryPick(line))
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}

func parseGitCommit(line string) *ir.GitCommit {
	c := &ir.GitCommit{}
	opts := gitKeyValRe.FindAllStringSubmatch(line, -1)
	for _, m := range opts {
		key := strings.ToLower(m[1])
		val := m[2]
		if val == "" {
			val = m[3]
		}
		switch key {
		case "id":
			c.ID = val
		case "tag":
			c.Tag = val
		case "type":
			c.Type = parseGitCommitType(val)
		}
	}
	return c
}

func parseGitBranch(line string) *ir.GitBranch {
	// "branch <name> order: <n>"
	rest := strings.TrimSpace(line[len("branch "):])
	b := &ir.GitBranch{Order: -1}

	// Check for "order:" option.
	if idx := strings.Index(strings.ToLower(rest), "order:"); idx >= 0 {
		b.Name = strings.TrimSpace(rest[:idx])
		orderStr := strings.TrimSpace(rest[idx+len("order:"):])
		if n, err := strconv.Atoi(orderStr); err == nil {
			b.Order = n
		}
	} else {
		b.Name = rest
	}

	// Strip quotes from name if present.
	b.Name = strings.Trim(b.Name, `"`)

	return b
}

func parseGitMerge(line string) *ir.GitMerge {
	rest := strings.TrimSpace(line[len("merge "):])
	mg := &ir.GitMerge{}

	// Extract the branch name (first word before any key: value pairs).
	parts := strings.Fields(rest)
	if len(parts) > 0 {
		mg.Branch = strings.Trim(parts[0], `"`)
	}

	// Parse key-value options.
	opts := gitKeyValRe.FindAllStringSubmatch(rest, -1)
	for _, m := range opts {
		key := strings.ToLower(m[1])
		val := m[2]
		if val == "" {
			val = m[3]
		}
		switch key {
		case "id":
			mg.ID = val
		case "tag":
			mg.Tag = val
		case "type":
			mg.Type = parseGitCommitType(val)
		}
	}

	return mg
}

func parseGitCherryPick(line string) *ir.GitCherryPick {
	cp := &ir.GitCherryPick{}
	opts := gitKeyValRe.FindAllStringSubmatch(line, -1)
	for _, m := range opts {
		key := strings.ToLower(m[1])
		val := m[2]
		if val == "" {
			val = m[3]
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

func parseGitCommitType(s string) ir.GitCommitType {
	switch strings.ToUpper(s) {
	case "REVERSE":
		return ir.GitCommitReverse
	case "HIGHLIGHT":
		return ir.GitCommitHighlight
	default:
		return ir.GitCommitNormal
	}
}
