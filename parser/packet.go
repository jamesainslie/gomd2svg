package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/gomd2svg/ir"
)

var (
	packetRangeRe    = regexp.MustCompile(`^(\d+)-(\d+)\s*:\s*"([^"]*)"$`)
	packetBitCountRe = regexp.MustCompile(`^\+(\d+)\s*:\s*"([^"]*)"$`)
)

func parsePacket(input string) (*ParseOutput, error) { //nolint:unparam // error return is part of the parser interface contract used by Parse().
	graph := ir.NewGraph()
	graph.Kind = ir.Packet

	lines := preprocessInput(input)
	nextBit := 0

	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "packet") {
			continue
		}

		// Try range notation: 0-15: "Source Port".
		if match := packetRangeRe.FindStringSubmatch(line); match != nil {
			start, _ := strconv.Atoi(match[1]) //nolint:errcheck // regex guarantees digits.
			end, _ := strconv.Atoi(match[2])   //nolint:errcheck // regex guarantees digits.
			desc := match[3]
			graph.Fields = append(graph.Fields, &ir.PacketField{
				Start: start, End: end, Description: desc,
			})
			nextBit = end + 1
			continue
		}

		// Try bit count notation: +16: "Source Port".
		if match := packetBitCountRe.FindStringSubmatch(line); match != nil {
			count, _ := strconv.Atoi(match[1]) //nolint:errcheck // regex guarantees digits.
			desc := match[2]
			start := nextBit
			end := start + count - 1
			graph.Fields = append(graph.Fields, &ir.PacketField{
				Start: start, End: end, Description: desc,
			})
			nextBit = end + 1
			continue
		}
	}

	return &ParseOutput{Graph: graph}, nil
}
