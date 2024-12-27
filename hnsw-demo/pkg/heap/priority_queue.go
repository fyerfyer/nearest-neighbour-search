package heap

import (
	"container/heap"
)

// Item represents an item in the priority queue
type Item struct {
	NodeID   int     // The ID of the node
	Distance float64 // Distance value (priority)
	Index    int     // Index in the heap (used by heap.Interface)
}

// PriorityQueue implements heap.Interface and holds Items
type PriorityQueue []*Item

// NewPriorityQueue creates a new priority queue
func NewPriorityQueue() *PriorityQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &pq
}

// Len returns the length of the queue
func (pq PriorityQueue) Len() int {
	return len(pq)
}

// Less determines the priority of items
// For min-heap (nearest neighbors), we want smaller distances to have higher priority
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Distance < pq[j].Distance
}

// Swap swaps two items in the queue
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

// Push adds an item to the queue
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the highest priority item
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// PushItem adds a new item to the queue
func (pq *PriorityQueue) PushItem(nodeID int, distance float64) {
	item := &Item{
		NodeID:   nodeID,
		Distance: distance,
	}
	heap.Push(pq, item)
}

// PopItem removes and returns the highest priority node ID and its distance
func (pq *PriorityQueue) PopItem() (int, float64) {
	if pq.Len() == 0 {
		return -1, -1
	}
	item := heap.Pop(pq).(*Item)
	return item.NodeID, item.Distance
}

// Top returns the highest priority item without removing it
func (pq PriorityQueue) Top() (*Item, bool) {
	if pq.Len() == 0 {
		return nil, false
	}
	return pq[0], true
}

// Clear removes all items from the queue
func (pq *PriorityQueue) Clear() {
	*pq = make(PriorityQueue, 0)
}

// Contains checks if a nodeID exists in the queue
func (pq PriorityQueue) Contains(nodeID int) bool {
	for _, item := range pq {
		if item.NodeID == nodeID {
			return true
		}
	}
	return false
}

// Update modifies the distance of an existing item
func (pq *PriorityQueue) Update(nodeID int, distance float64) bool {
	for _, item := range *pq {
		if item.NodeID == nodeID {
			item.Distance = distance
			heap.Fix(pq, item.Index)
			return true
		}
	}
	return false
}
