package algorithm

import "github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/heap"

// selectNeighborsHeuristic implements neighbor selection with heuristic algorithm
func (h *HNSW) selectNeighborsHeuristic(q []float64, candidates []int, M int,
	level int, extendCandidates bool, keepPrunedConnections bool) []int {

	// Create working queue W
	workingQueue := heap.NewPriorityQueue()
	visited := make(map[int]bool)

	// Add all candidates to working queue
	for _, candidateID := range candidates {
		if !visited[candidateID] {
			dist := h.distFunc(q, h.nodes[candidateID].GetVector())
			workingQueue.PushItem(candidateID, dist)
			visited[candidateID] = true
		}
	}

	// Extend candidates if needed
	if extendCandidates {
		// Store new candidates temporarily
		tempCandidates := heap.NewPriorityQueue()

		// Check neighbors of current candidates
		for _, candidateID := range candidates {
			neighbors, _ := h.nodes[candidateID].GetNeighbors(level)
			for _, neighborID := range neighbors {
				if !visited[neighborID] {
					dist := h.distFunc(q, h.nodes[neighborID].GetVector())
					tempCandidates.PushItem(neighborID, dist)
					visited[neighborID] = true
				}
			}
		}

		// Add new candidates to working queue
		for tempCandidates.Len() > 0 {
			nodeID, dist := tempCandidates.PopItem()
			workingQueue.PushItem(nodeID, dist)
		}
	}

	// Create result set R and discarded queue Wd
	results := make([]int, 0, M)
	discardedQueue := heap.NewPriorityQueue()

	// Main loop: process working queue
	for workingQueue.Len() > 0 && len(results) < M {
		nodeID, dist := workingQueue.PopItem()

		// Check if should add to results
		shouldAdd := true
		if len(results) > 0 {
			// Check relationship with existing results
			for _, resultID := range results {
				resultDist := h.distFunc(h.nodes[resultID].GetVector(), h.nodes[nodeID].GetVector())
				if resultDist < dist {
					shouldAdd = false
					break
				}
			}
		}

		if shouldAdd {
			results = append(results, nodeID)
		} else {
			discardedQueue.PushItem(nodeID, dist)
		}
	}

	// Add pruned connections if needed
	if keepPrunedConnections {
		for discardedQueue.Len() > 0 && len(results) < M {
			nodeID, _ := discardedQueue.PopItem()
			results = append(results, nodeID)
		}
	}

	return results
}

// SelectNeighborsHeuristic is the public interface for neighbor selection
func (h *HNSW) SelectNeighborsHeuristic(q []float64, candidates []int, M int,
	level int, extendCandidates bool, keepPrunedConnections bool) []int {
	return h.selectNeighborsHeuristic(q, candidates, M, level,
		extendCandidates, keepPrunedConnections)
}
