package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	blockColumnsRe = regexp.MustCompile(`^columns\s+(\d+)$`)
	blockEdgeRe    = regexp.MustCompile(`^(\w+)\s*(-->|---)\s*(?:\|"?([^"|]*)"?\|\s*)?(\w+)$`)
	blockDefRe     = regexp.MustCompile(`^(\w+)(?:\["([^"]*)"\]|\("([^"]*)"\)|\(\("([^"]*)"\)\)|\{"([^"]*)"\}|>\["([^"]*)"\])?(?::(\d+))?\s*$`)
)

func parseBlock(input string) (*ParseOutput, error) {
	lines := preprocessInput(input)
	g := ir.NewGraph()
	g.Kind = ir.Block

	if len(lines) > 0 {
		lower := strings.ToLower(lines[0])
		if strings.HasPrefix(lower, "block") {
			lines = lines[1:]
		}
	}

	for _, line := range lines {
		if m := blockColumnsRe.FindStringSubmatch(line); m != nil {
			cols, _ := strconv.Atoi(m[1])
			g.BlockColumns = cols
			continue
		}

		if m := blockEdgeRe.FindStringSubmatch(line); m != nil {
			from, arrow, label, to := m[1], m[2], m[3], m[4]
			edge := &ir.Edge{
				From:     from,
				To:       to,
				Directed: arrow == "-->",
				ArrowEnd: arrow == "-->",
			}
			if label != "" {
				edge.Label = &label
			}
			g.Edges = append(g.Edges, edge)
			g.EnsureNode(from, nil, nil)
			g.EnsureNode(to, nil, nil)
			continue
		}

		parseBlockDefs(line, g)
	}

	return &ParseOutput{Graph: g}, nil
}

func parseBlockDefs(line string, g *ir.Graph) {
	tokens := strings.Fields(line)
	for _, token := range tokens {
		m := blockDefRe.FindStringSubmatch(token)
		if m == nil {
			continue
		}
		id := m[1]
		label := id
		shape := ir.Rectangle
		switch {
		case m[2] != "":
			label = m[2]
			shape = ir.Rectangle
		case m[3] != "":
			label = m[3]
			shape = ir.RoundRect
		case m[4] != "":
			label = m[4]
			shape = ir.Circle
		case m[5] != "":
			label = m[5]
			shape = ir.Diamond
		case m[6] != "":
			label = m[6]
			shape = ir.Asymmetric
		}

		width := 1
		if m[7] != "" {
			width, _ = strconv.Atoi(m[7])
		}

		block := &ir.BlockDef{
			ID:    id,
			Label: label,
			Shape: shape,
			Width: width,
		}
		g.Blocks = append(g.Blocks, block)
		g.EnsureNode(id, &label, &shape)
	}
}
