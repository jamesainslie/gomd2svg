package layout

import (
	"sort"

	"github.com/jamesainslie/gomd2svg/ir"
)

// computeRanks assigns an integer rank to each node using a modified Kahn's
// algorithm. Each node's rank is max(predecessor_rank) + 1, with root nodes
// at rank 0. When a cycle is detected (queue empties but unranked nodes
// remain), the unranked node with the lowest nodeOrder is forced into the
// queue to break the cycle.
func computeRanks(nodes []string, edges []*ir.Edge, nodeOrder map[string]int) map[string]int { //nolint:gocognit // ranking algorithm requires sequential state tracking
	// Build adjacency list and in-degree map.
	adj := make(map[string][]string)
	inDegree := make(map[string]int)
	predecessors := make(map[string][]string)

	nodeSet := make(map[string]bool, len(nodes))
	for _, nodeID := range nodes {
		nodeSet[nodeID] = true
		inDegree[nodeID] = 0
	}

	for _, edge := range edges {
		if !nodeSet[edge.From] || !nodeSet[edge.To] {
			continue
		}
		adj[edge.From] = append(adj[edge.From], edge.To)
		inDegree[edge.To]++
		predecessors[edge.To] = append(predecessors[edge.To], edge.From)
	}

	// Initialize the queue with all nodes that have in-degree 0.
	var queue []string
	for _, nodeID := range nodes {
		if inDegree[nodeID] == 0 {
			queue = append(queue, nodeID)
		}
	}

	ranks := make(map[string]int, len(nodes))
	processed := 0

	for processed < len(nodes) {
		// If queue is empty but we have unprocessed nodes, we have a cycle.
		// Break it by picking the unranked node with the lowest nodeOrder.
		if len(queue) == 0 {
			best, bestOrder := computeRanksFindBestCycleBreak(nodes, ranks, nodeOrder)
			if best == "" {
				break // safety: should not happen
			}
			_ = bestOrder // used only for selection
			queue = append(queue, best)
			// The forced node gets rank 0 (or max of its already-ranked predecessors + 1).
			rank := 0
			for _, pred := range predecessors[best] {
				if predRank, ok := ranks[pred]; ok && predRank+1 > rank {
					rank = predRank + 1
				}
			}
			ranks[best] = rank
			processed++

			// Process successors: decrement their in-degree.
			for _, succ := range adj[best] {
				inDegree[succ]--
				if inDegree[succ] == 0 {
					if _, ranked := ranks[succ]; !ranked {
						queue = append(queue, succ)
					}
				}
			}
			// Remove the processed node from queue head.
			queue = queue[1:]
			continue
		}

		// Normal BFS processing.
		curr := queue[0]
		queue = queue[1:]

		if _, ranked := ranks[curr]; ranked {
			continue // already processed (e.g. by cycle-breaking)
		}

		// Rank = max(predecessor_rank) + 1, or 0 if no predecessors.
		rank := 0
		for _, pred := range predecessors[curr] {
			if predRank, ok := ranks[pred]; ok && predRank+1 > rank {
				rank = predRank + 1
			}
		}
		ranks[curr] = rank
		processed++

		// Add successors whose in-degree reaches 0.
		for _, succ := range adj[curr] {
			inDegree[succ]--
			if inDegree[succ] == 0 {
				if _, ranked := ranks[succ]; !ranked {
					queue = append(queue, succ)
				}
			}
		}
	}

	return ranks
}

// computeRanksFindBestCycleBreak finds the unranked node with the lowest
// nodeOrder to break a cycle in the ranking algorithm.
func computeRanksFindBestCycleBreak(nodes []string, ranks map[string]int, nodeOrder map[string]int) (string, int) {
	var best string
	bestOrder := int(^uint(0) >> 1) // max int
	for _, nodeID := range nodes {
		if _, ranked := ranks[nodeID]; ranked {
			continue
		}
		order, ok := nodeOrder[nodeID]
		if !ok {
			order = 0
		}
		if best == "" || order < bestOrder {
			best = nodeID
			bestOrder = order
		}
	}
	return best, bestOrder
}

// sortedNodeIDs returns the node IDs from a map, sorted by their nodeOrder.
func sortedNodeIDs(nodes map[string]*ir.Node, nodeOrder map[string]int) []string {
	ids := make([]string, 0, len(nodes))
	for id := range nodes {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(idxA, idxB int) bool {
		orderA := nodeOrder[ids[idxA]]
		orderB := nodeOrder[ids[idxB]]
		if orderA != orderB {
			return orderA < orderB
		}
		return ids[idxA] < ids[idxB]
	})
	return ids
}
