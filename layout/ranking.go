package layout

import (
	"sort"

	"github.com/jamesainslie/mermaid-go/ir"
)

// computeRanks assigns an integer rank to each node using a modified Kahn's
// algorithm. Each node's rank is max(predecessor_rank) + 1, with root nodes
// at rank 0. When a cycle is detected (queue empties but unranked nodes
// remain), the unranked node with the lowest nodeOrder is forced into the
// queue to break the cycle.
func computeRanks(nodes []string, edges []*ir.Edge, nodeOrder map[string]int) map[string]int {
	// Build adjacency list and in-degree map.
	adj := make(map[string][]string)
	inDegree := make(map[string]int)
	predecessors := make(map[string][]string)

	nodeSet := make(map[string]bool, len(nodes))
	for _, n := range nodes {
		nodeSet[n] = true
		inDegree[n] = 0
	}

	for _, e := range edges {
		if !nodeSet[e.From] || !nodeSet[e.To] {
			continue
		}
		adj[e.From] = append(adj[e.From], e.To)
		inDegree[e.To]++
		predecessors[e.To] = append(predecessors[e.To], e.From)
	}

	// Initialize the queue with all nodes that have in-degree 0.
	var queue []string
	for _, n := range nodes {
		if inDegree[n] == 0 {
			queue = append(queue, n)
		}
	}

	ranks := make(map[string]int, len(nodes))
	processed := 0

	for processed < len(nodes) {
		// If queue is empty but we have unprocessed nodes, we have a cycle.
		// Break it by picking the unranked node with the lowest nodeOrder.
		if len(queue) == 0 {
			var best string
			bestOrder := int(^uint(0) >> 1) // max int
			for _, n := range nodes {
				if _, ranked := ranks[n]; ranked {
					continue
				}
				order, ok := nodeOrder[n]
				if !ok {
					order = 0
				}
				if best == "" || order < bestOrder {
					best = n
					bestOrder = order
				}
			}
			if best == "" {
				break // safety: should not happen
			}
			queue = append(queue, best)
			// The forced node gets rank 0 (or max of its already-ranked predecessors + 1).
			rank := 0
			for _, pred := range predecessors[best] {
				if pr, ok := ranks[pred]; ok && pr+1 > rank {
					rank = pr + 1
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
			if pr, ok := ranks[pred]; ok && pr+1 > rank {
				rank = pr + 1
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

// sortedNodeIDs returns the node IDs from a map, sorted by their nodeOrder.
func sortedNodeIDs(nodes map[string]*ir.Node, nodeOrder map[string]int) []string {
	ids := make([]string, 0, len(nodes))
	for id := range nodes {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		oi := nodeOrder[ids[i]]
		oj := nodeOrder[ids[j]]
		if oi != oj {
			return oi < oj
		}
		return ids[i] < ids[j]
	})
	return ids
}
