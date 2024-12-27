package algorithm

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/config"
	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/distance"
	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/heap"
	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/node"
)

// HNSW represents the hierarchical navigable small world graph
type HNSW struct {
	nodes      map[int]*node.Node
	entryPoint int
	maxLevel   int
	config     config.Config
	distFunc   distance.DistanceFunction
	mutex      sync.RWMutex
	nodesMutex sync.RWMutex
	// deletedCount int
}

// New creates a new HNSW index
func New(cfg config.Config, metric string) (*HNSW, error) {
	distFunc, err := distance.GetDistanceFunction(metric)
	if err != nil {
		return nil, fmt.Errorf("failed to get distance function: %v", err)
	}

	return &HNSW{
		nodes:    make(map[int]*node.Node),
		config:   cfg,
		distFunc: distFunc,
	}, nil
}

// generateLevel generates random level for new nodes
func (h *HNSW) generateLevel() int {
	level := int(math.Floor(-math.Log(rand.Float64()) * h.config.ML))
	if level > h.maxLevel {
		h.mutex.Lock()
		if level > h.maxLevel {
			h.maxLevel = level
		}
		h.mutex.Unlock()
	}
	return level
}

// Insert adds a new element to the index
func (h *HNSW) Insert(id int, vector []float64) error {
	// Check if node already exists
	h.nodesMutex.Lock()
	if _, exists := h.nodes[id]; exists {
		h.nodesMutex.Unlock()
		return fmt.Errorf("node %d already exists", id)
	}

	// Create new node
	level := h.generateLevel()
	newNode := node.NewNode(id, vector, level)
	h.nodes[id] = newNode
	h.nodesMutex.Unlock()

	// Handle first node
	h.mutex.Lock()
	if h.entryPoint == 0 {
		h.entryPoint = id
		h.mutex.Unlock()
		return nil
	}
	h.mutex.Unlock()

	// Search for insert
	currObj := h.entryPoint
	for lc := h.maxLevel; lc > level; lc-- {
		changed := false
		currNode := h.nodes[currObj]

		// Find better point to start from
		for _, neighbor := range currNode.Neighbors[lc] {
			if dist := h.distFunc(vector, h.nodes[neighbor].GetVector()); dist < h.distFunc(vector, currNode.GetVector()) {
				currObj = neighbor
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	// Connect on each level
	for lc := min(level, h.maxLevel); lc >= 0; lc-- {
		// Find candidates
		candidates := h.searchLayer(vector, currObj, h.config.EfConstruction, lc)

		// Select neighbors
		neighbors := h.selectNeighborsHeuristic(vector, candidates, h.config.M, lc, true, true)

		// Add connections
		for _, neighborID := range neighbors {
			newNode.AddNeighbor(lc, neighborID)
			h.nodes[neighborID].AddNeighbor(lc, id)
		}

		currObj = id
	}

	return nil
}

// searchLayer implements layer-wise search
func (h *HNSW) searchLayer(q []float64, entryPointID int, ef int, level int) []int {
	visited := make(map[int]bool)
	candidates := heap.NewPriorityQueue()
	results := heap.NewPriorityQueue()

	dist := h.distFunc(q, h.nodes[entryPointID].GetVector())
	candidates.PushItem(entryPointID, dist)
	results.PushItem(entryPointID, dist)
	visited[entryPointID] = true

	for candidates.Len() > 0 {
		// Fix PopItem usage
		nodeID, nodeDist := candidates.PopItem()

		// Fix Top usage
		furthest, exists := results.Top()
		if exists && nodeDist > furthest.Distance {
			break
		}

		neighbors, _ := h.nodes[nodeID].GetNeighbors(level)
		for _, neighborID := range neighbors {
			if !visited[neighborID] {
				visited[neighborID] = true
				dist := h.distFunc(q, h.nodes[neighborID].GetVector())

				furthest, exists := results.Top()
				if !exists || results.Len() < ef || dist < furthest.Distance {
					candidates.PushItem(neighborID, dist)
					results.PushItem(neighborID, dist)
					if results.Len() > ef {
						results.Pop()
					}
				}
			}
		}
	}

	resultIds := make([]int, 0, results.Len())
	for results.Len() > 0 {
		id, _ := results.PopItem()
		resultIds = append(resultIds, id)
	}
	return resultIds
}

// Search performs K-NN search
func (h *HNSW) Search(q []float64, k int) ([]int, []float64) {
	if len(h.nodes) == 0 {
		return nil, nil
	}

	h.mutex.RLock()
	ep := h.entryPoint
	currentLevel := h.maxLevel
	h.mutex.RUnlock()

	currObj := ep
	for level := currentLevel; level > 0; level-- {
		changed := false
		currNode := h.nodes[currObj]

		// Find better point to start from
		for _, neighbor := range currNode.Neighbors[level] {
			if dist := h.distFunc(q, h.nodes[neighbor].GetVector()); dist < h.distFunc(q, currNode.GetVector()) {
				currObj = neighbor
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	// Search at layer 0
	candidates := h.searchLayer(q, currObj, k*2, 0)

	// Get k nearest
	results := make([]int, 0, k)
	distances := make([]float64, 0, k)

	for _, id := range candidates {
		if len(results) >= k {
			break
		}
		results = append(results, id)
		distances = append(distances, h.distFunc(q, h.nodes[id].GetVector()))
	}

	return results, distances
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
