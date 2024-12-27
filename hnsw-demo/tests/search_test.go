package tests

import (
	"testing"

	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/config"
	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/distance"
	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/src/algorithm"
)

func TestKNNSearch(t *testing.T) {
	// Setup
	cfg := config.NewDefaultConfig()
	hnsw, err := algorithm.New(cfg, distance.Euclidean)
	if err != nil {
		t.Fatalf("Failed to create HNSW: %v", err)
	}

	// Insert test data
	vectors := map[int][]float64{
		1: {1.0, 2.0},
		2: {2.0, 1.0},
		3: {3.0, 4.0},
		4: {4.0, 3.0},
		5: {5.0, 5.0},
	}

	for id, vec := range vectors {
		if err := hnsw.Insert(id, vec); err != nil {
			t.Fatalf("Failed to insert vector %d: %v", id, err)
		}
	}

	tests := []struct {
		name      string
		query     []float64
		k         int
		ef        int
		wantLen   int
		wantFirst int // ID of expected first result
	}{
		{
			name:      "Basic KNN search",
			query:     []float64{1.0, 1.0},
			k:         2,
			ef:        10,
			wantLen:   2,
			wantFirst: 2, // Closest to (1,1)
		},
		{
			name:      "K larger than dataset",
			query:     []float64{1.0, 1.0},
			k:         10,
			ef:        20,
			wantLen:   5, // Should return all vectors
			wantFirst: 2,
		},
		{
			name:      "Single nearest neighbor",
			query:     []float64{5.0, 5.0},
			k:         1,
			ef:        10,
			wantLen:   1,
			wantFirst: 5, // Exactly matches (5,5)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, distances := hnsw.KNNSearchWithDistances(tt.query, tt.k, tt.ef)

			// Check result length
			if len(results) != tt.wantLen {
				t.Errorf("got len %d, want %d", len(results), tt.wantLen)
			}

			// Check distances are sorted
			for i := 1; i < len(distances); i++ {
				if distances[i] < distances[i-1] {
					t.Error("distances not sorted")
				}
			}

			// Check first result if specified
			if len(results) > 0 && results[0] != tt.wantFirst {
				t.Errorf("got first result %d, want %d", results[0], tt.wantFirst)
			}
		})
	}
}

func TestSearchEmpty(t *testing.T) {
	cfg := config.NewDefaultConfig()
	hnsw, err := algorithm.New(cfg, distance.Euclidean)
	if err != nil {
		t.Fatalf("Failed to create HNSW: %v", err)
	}

	query := []float64{1.0, 1.0}
	results, distances := hnsw.KNNSearchWithDistances(query, 1, 10)

	if len(results) != 0 || len(distances) != 0 {
		t.Error("Expected empty results for empty index")
	}
}

func TestSearchDimensionMismatch(t *testing.T) {
	cfg := config.NewDefaultConfig()
	hnsw, err := algorithm.New(cfg, distance.Euclidean)
	if err != nil {
		t.Fatalf("Failed to create HNSW: %v", err)
	}

	// Insert 2D vector
	err = hnsw.Insert(1, []float64{1.0, 1.0})
	if err != nil {
		t.Fatalf("Failed to insert vector: %v", err)
	}

	// Search with 3D query
	query := []float64{1.0, 1.0, 1.0}
	results, distances := hnsw.KNNSearchWithDistances(query, 1, 10)

	if len(results) != 0 || len(distances) != 0 {
		t.Error("Expected empty results for dimension mismatch")
	}
}
