package algorithm

import (
	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/heap"
)

func (h *HNSW) selectNeighborsSimple(q []float64, candidates []int, M int) []int {
	if len(candidates) <= M {
		return candidates
	}

	pq := heap.NewPriorityQueue()

	for _, candidateID := range candidates {
		dist := h.distFunc(q, h.nodes[candidateID].GetVector())
		pq.PushItem(candidateID, dist)
	}

	result := make([]int, 0, M)
	i := 0
	for pq.Len() > 0 && i < M {
		nodeID, _ := pq.PopItem()
		result = append(result, nodeID)
		i++
	}

	return result
}

func (h *HNSW) SelectNeighborsSimple(q []float64, candidates []int, M int) []int {
	return h.selectNeighborsSimple(q, candidates, M)
}
