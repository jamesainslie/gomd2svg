package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

const directedArrow = "-->"

var (
	blockColumnsRe = regexp.MustCompile(`^columns\s+(\d+)$`)
	blockEdgeRe    = regexp.MustCompile(`^(\w+)\s*(-->|---)\s*(?:\|"?([^"|]*)"?\|\s*)?(\w+)$`)
	blockDefRe     = regexp.MustCompile(`^(\w+)(?:\["([^"]*)"\]|\("([^"]*)"\)|\(\("([^"]*)"\)\)|\{"([^"]*)"\}|>\["([^"]*)"\])?(?::(\d+))?\s*$`)
)

//nolint:unparam // error return is part of the parser interface contract used by Parse().
func parseBlock(input string) (*ParseOutput, error) {
	lines := preprocessInput(input)
	graph := ir.NewGraph()
	graph.Kind = ir.Block

	if len(lines) > 0 {
		lower := strings.ToLower(lines[0])
		if strings.HasPrefix(lower, "block") {
			lines = lines[1:]
		}
	}

	for _, line := range lines {
		if match := blockColumnsRe.FindStringSubmatch(line); match != nil {
			cols, errConv := strconv.Atoi(match[1])
			if errConv == nil {
				graph.BlockColumns = cols
			}
			continue
		}

		if match := blockEdgeRe.FindStringSubmatch(line); match != nil {
			from, arrow, label, to := match[1], match[2], match[3], match[4]
			edge := &ir.Edge{
				From:     from,
				To:       to,
				Directed: arrow == directedArrow,
				ArrowEnd: arrow == directedArrow,
			}
			if label != "" {
				edge.Label = &label
			}
			graph.Edges = append(graph.Edges, edge)
			graph.EnsureNode(from, nil, nil)
			graph.EnsureNode(to, nil, nil)
			continue
		}

		parseBlockDefs(line, graph)
	}

	return &ParseOutput{Graph: graph}, nil
}

func parseBlockDefs(line string, graph *ir.Graph) {
	tokens := strings.Fields(line)
	for _, token := range tokens {
		match := blockDefRe.FindStringSubmatch(token)
		if match == nil {
			continue
		}
		id := match[1]
		label := id
		shape := ir.Rectangle
		switch {
		case match[2] != "":
			label = match[2]
			shape = ir.Rectangle
		case match[3] != "":
			label = match[3]
			shape = ir.RoundRect
		case match[4] != "":
			label = match[4]
			shape = ir.Circle
		case match[5] != "":
			label = match[5]
			shape = ir.Diamond
		case match[6] != "":
			label = match[6]
			shape = ir.Asymmetric
		}

		width := 1
		if match[7] != "" {
			parsed, errConv := strconv.Atoi(match[7])
			if errConv == nil {
				width = parsed
			}
		}

		block := &ir.BlockDef{
			ID:    id,
			Label: label,
			Shape: shape,
			Width: width,
		}
		graph.Blocks = append(graph.Blocks, block)
		graph.EnsureNode(id, &label, &shape)
	}
}
