package layout

import (
	"sort"

	"github.com/jamesainslie/mermaid-go/ir"
)

// orderRankNodes organizes nodes into ranked layers and minimizes edge
// crossings using a median heuristic over multiple passes. Returns a
// slice of slices where each inner slice contains the node IDs in a
// single rank, ordered to minimize crossings.
func orderRankNodes(
	ranks map[string]int,
	edges []*ir.Edge,
	passes int,
) [][]string {
	if len(ranks) == 0 {
		return nil
	}

	// Determine the maximum rank.
	maxRank := 0
	for _, r := range ranks {
		if r > maxRank {
			maxRank = r
		}
	}

	// Group nodes by rank.
	layers := make([][]string, maxRank+1)
	for id, r := range ranks {
		layers[r] = append(layers[r], id)
	}

	// Sort each layer alphabetically as a stable starting point.
	for _, layer := range layers {
		sort.Strings(layer)
	}

	// Build adjacency structures for median computation.
	// successors[node] = list of nodes in the next rank
	// predecessors[node] = list of nodes in the previous rank
	successors := make(map[string][]string)
	prevs := make(map[string][]string)
	for _, e := range edges {
		fromRank, fromOK := ranks[e.From]
		toRank, toOK := ranks[e.To]
		if !fromOK || !toOK {
			continue
		}
		if toRank == fromRank+1 {
			successors[e.From] = append(successors[e.From], e.To)
			prevs[e.To] = append(prevs[e.To], e.From)
		}
	}

	// Crossing minimization: median heuristic.
	for pass := 0; pass < passes; pass++ {
		if pass%2 == 0 {
			// Forward sweep: use predecessors to order each rank.
			for r := 1; r <= maxRank; r++ {
				posInPrev := positionMap(layers[r-1])
				sortByMedian(layers[r], prevs, posInPrev)
			}
		} else {
			// Backward sweep: use successors to order each rank.
			for r := maxRank - 1; r >= 0; r-- {
				posInNext := positionMap(layers[r+1])
				sortByMedian(layers[r], successors, posInNext)
			}
		}
	}

	return layers
}

// positionMap builds a map from node ID to its index within a layer.
func positionMap(layer []string) map[string]int {
	m := make(map[string]int, len(layer))
	for i, id := range layer {
		m[id] = i
	}
	return m
}

// sortByMedian sorts a layer's nodes by the median position of their
// connected neighbors in the reference layer.
func sortByMedian(layer []string, neighbors map[string][]string, refPos map[string]int) {
	medians := make(map[string]float32, len(layer))

	for _, id := range layer {
		nbrs := neighbors[id]
		if len(nbrs) == 0 {
			// No neighbors: keep current relative position.
			medians[id] = -1
			continue
		}

		// Collect positions of neighbors in the reference layer.
		positions := make([]int, 0, len(nbrs))
		for _, n := range nbrs {
			if pos, ok := refPos[n]; ok {
				positions = append(positions, pos)
			}
		}
		if len(positions) == 0 {
			medians[id] = -1
			continue
		}

		sort.Ints(positions)
		mid := len(positions) / 2
		if len(positions)%2 == 0 {
			medians[id] = float32(positions[mid-1]+positions[mid]) / 2.0
		} else {
			medians[id] = float32(positions[mid])
		}
	}

	// Stable sort: nodes without neighbors keep their relative order.
	sort.SliceStable(layer, func(i, j int) bool {
		mi := medians[layer[i]]
		mj := medians[layer[j]]
		if mi < 0 && mj < 0 {
			return false // both have no neighbors, keep original order
		}
		if mi < 0 {
			return false // push unconnected nodes after connected
		}
		if mj < 0 {
			return true
		}
		return mi < mj
	})
}
