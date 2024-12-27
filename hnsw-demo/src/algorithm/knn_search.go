package algorithm

// KNNSearch implements k-nearest neighbor search
func (h *HNSW) KNNSearch(q []float64, K int, ef int) []int {
	// Check if graph is empty
	if len(h.nodes) == 0 {
		return []int{}
	}

	// Get entry point
	h.mutex.RLock()
	ep := h.entryPoint
	currentLevel := h.maxLevel
	h.mutex.RUnlock()

	if ep == 0 {
		return []int{}
	}

	// Search from top layer
	currObj := ep
	for level := currentLevel; level >= 1; level-- {
		// Search layer with ef=1 to find better entry point
		candidates := h.searchLayer(q, currObj, 1, level)
		if len(candidates) > 0 {
			currObj = candidates[0]
		}
	}

	// Search bottom layer with specified ef
	finalResults := h.searchLayer(q, currObj, ef, 0)

	// Return K nearest elements
	if K > len(finalResults) {
		K = len(finalResults)
	}
	return finalResults[:K]
}

// KNNSearchWithDistances returns K nearest neighbors with distances
func (h *HNSW) KNNSearchWithDistances(q []float64, K int, ef int) ([]int, []float64) {
	// Get K nearest neighbors
	ids := h.KNNSearch(q, K, ef)

	// Calculate distances
	distances := make([]float64, len(ids))
	for i, id := range ids {
		distances[i] = h.distFunc(q, h.nodes[id].GetVector())
	}

	return ids, distances
}
