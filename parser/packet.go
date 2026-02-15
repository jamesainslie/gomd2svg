package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jamesainslie/mermaid-go/ir"
)

var (
	packetRangeRe    = regexp.MustCompile(`^(\d+)-(\d+)\s*:\s*"([^"]*)"$`)
	packetBitCountRe = regexp.MustCompile(`^\+(\d+)\s*:\s*"([^"]*)"$`)
)

func parsePacket(input string) (*ParseOutput, error) {
	g := ir.NewGraph()
	g.Kind = ir.Packet

	lines := preprocessInput(input)
	nextBit := 0

	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "packet") {
			continue
		}

		// Try range notation: 0-15: "Source Port"
		if m := packetRangeRe.FindStringSubmatch(line); m != nil {
			start, _ := strconv.Atoi(m[1]) // regex guarantees digits
			end, _ := strconv.Atoi(m[2])   // regex guarantees digits
			desc := m[3]
			g.Fields = append(g.Fields, &ir.PacketField{
				Start: start, End: end, Description: desc,
			})
			nextBit = end + 1
			continue
		}

		// Try bit count notation: +16: "Source Port"
		if m := packetBitCountRe.FindStringSubmatch(line); m != nil {
			count, _ := strconv.Atoi(m[1]) // regex guarantees digits
			desc := m[2]
			start := nextBit
			end := start + count - 1
			g.Fields = append(g.Fields, &ir.PacketField{
				Start: start, End: end, Description: desc,
			})
			nextBit = end + 1
			continue
		}
	}

	return &ParseOutput{Graph: g}, nil
}
