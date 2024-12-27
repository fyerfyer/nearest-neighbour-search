package node

import (
	"fmt"
	"sync"
)

// Node represents a node in the HNSW graph
type Node struct {
	// Basic properties
	ID     int
	Vector []float64
	Level  int

	// Neighbors at each level
	// map[level][]neighborID
	Neighbors map[int][]int

	// Concurrency control
	mutex sync.RWMutex

	// Soft deletion flag
	deleted bool
}

// NewNode creates a new node instance
func NewNode(id int, vector []float64, level int) *Node {
	return &Node{
		ID:        id,
		Vector:    vector,
		Level:     level,
		Neighbors: make(map[int][]int),
		deleted:   false,
	}
}

// AddNeighbor adds a neighbor at specified level
func (n *Node) AddNeighbor(level int, neighborID int) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.deleted {
		return fmt.Errorf("cannot add neighbor to deleted node %d", n.ID)
	}

	if level < 0 || level > n.Level {
		return fmt.Errorf("invalid level %d for node %d with max level %d", level, n.ID, n.Level)
	}

	// Initialize slice if not exists
	if _, exists := n.Neighbors[level]; !exists {
		n.Neighbors[level] = make([]int, 0)
	}

	// Check if neighbor already exists
	for _, id := range n.Neighbors[level] {
		if id == neighborID {
			return nil // Already exists
		}
	}

	n.Neighbors[level] = append(n.Neighbors[level], neighborID)
	return nil
}

// RemoveNeighbor removes a neighbor at specified level
func (n *Node) RemoveNeighbor(level int, neighborID int) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.deleted {
		return fmt.Errorf("cannot remove neighbor from deleted node %d", n.ID)
	}

	neighbors, exists := n.Neighbors[level]
	if !exists {
		return fmt.Errorf("level %d does not exist in node %d", level, n.ID)
	}

	// Find and remove neighbor
	for i, id := range neighbors {
		if id == neighborID {
			n.Neighbors[level] = append(neighbors[:i], neighbors[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("neighbor %d not found in level %d of node %d", neighborID, level, n.ID)
}

// GetNeighbors returns neighbors at specified level
func (n *Node) GetNeighbors(level int) ([]int, error) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	if n.deleted {
		return nil, fmt.Errorf("cannot get neighbors from deleted node %d", n.ID)
	}

	neighbors, exists := n.Neighbors[level]
	if !exists {
		return nil, fmt.Errorf("level %d does not exist in node %d", level, n.ID)
	}

	// Return copy to prevent external modifications
	result := make([]int, len(neighbors))
	copy(result, neighbors)
	return result, nil
}

// SetNeighbors sets all neighbors at specified level
func (n *Node) SetNeighbors(level int, neighbors []int) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.deleted {
		return fmt.Errorf("cannot set neighbors for deleted node %d", n.ID)
	}

	if level < 0 || level > n.Level {
		return fmt.Errorf("invalid level %d for node %d with max level %d", level, n.ID, n.Level)
	}

	// Create copy to prevent external modifications
	newNeighbors := make([]int, len(neighbors))
	copy(newNeighbors, neighbors)
	n.Neighbors[level] = newNeighbors
	return nil
}

// MarkDeleted marks the node as deleted
func (n *Node) MarkDeleted() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.deleted = true
}

// IsDeleted checks if the node is marked as deleted
func (n *Node) IsDeleted() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.deleted
}

// GetVector returns a copy of the node's vector
func (n *Node) GetVector() []float64 {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	result := make([]float64, len(n.Vector))
	copy(result, n.Vector)
	return result
}

// GetLevel returns the node's level
func (n *Node) GetLevel() int {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.Level
}

// ClearNeighbors removes all neighbors at specified level
func (n *Node) ClearNeighbors(level int) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.deleted {
		return fmt.Errorf("cannot clear neighbors of deleted node %d", n.ID)
	}

	if level < 0 || level > n.Level {
		return fmt.Errorf("invalid level %d for node %d with max level %d", level, n.ID, n.Level)
	}

	n.Neighbors[level] = make([]int, 0)
	return nil
}
