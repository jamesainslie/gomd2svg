package layout

import (
	"sort"

	"github.com/jamesainslie/gomd2svg/ir"
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
	for _, rank := range ranks {
		if rank > maxRank {
			maxRank = rank
		}
	}

	// Group nodes by rank.
	layers := make([][]string, maxRank+1)
	for id, rank := range ranks {
		layers[rank] = append(layers[rank], id)
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
	for _, edge := range edges {
		fromRank, fromOK := ranks[edge.From]
		toRank, toOK := ranks[edge.To]
		if !fromOK || !toOK {
			continue
		}
		if toRank == fromRank+1 {
			successors[edge.From] = append(successors[edge.From], edge.To)
			prevs[edge.To] = append(prevs[edge.To], edge.From)
		}
	}

	// Crossing minimization: median heuristic.
	for pass := range passes {
		if pass%2 == 0 {
			// Forward sweep: use predecessors to order each rank.
			for rank := 1; rank <= maxRank; rank++ {
				posInPrev := positionMap(layers[rank-1])
				sortByMedian(layers[rank], prevs, posInPrev)
			}
		} else {
			// Backward sweep: use successors to order each rank.
			for rank := maxRank - 1; rank >= 0; rank-- {
				posInNext := positionMap(layers[rank+1])
				sortByMedian(layers[rank], successors, posInNext)
			}
		}
	}

	return layers
}

// positionMap builds a map from node ID to its index within a layer.
func positionMap(layer []string) map[string]int {
	positions := make(map[string]int, len(layer))
	for idx, id := range layer {
		positions[id] = idx
	}
	return positions
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
		for _, nbr := range nbrs {
			if pos, ok := refPos[nbr]; ok {
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
			medians[id] = float32(positions[mid-1]+positions[mid]) / 2
		} else {
			medians[id] = float32(positions[mid])
		}
	}

	// Stable sort: nodes without neighbors keep their relative order.
	sort.SliceStable(layer, func(idxA, idxB int) bool {
		medianA := medians[layer[idxA]]
		medianB := medians[layer[idxB]]
		if medianA < 0 && medianB < 0 {
			return false // both have no neighbors, keep original order
		}
		if medianA < 0 {
			return false // push unconnected nodes after connected
		}
		if medianB < 0 {
			return true
		}
		return medianA < medianB
	})
}
