package tests

import (
    "testing"

    "github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/config"
    "github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/distance"
    "github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/src/algorithm"
)

func TestInsertAndSearch(t *testing.T) {
    // Create config
    cfg := config.Config{
        M:              16,
        MaxM:           32,
        EfConstruction: 100,
        ML:             1.0 / 1.5,
        DelayRebuild:   false,
    }

    // Create HNSW index with Euclidean distance
    hnsw, err := algorithm.New(cfg, distance.Euclidean)
    if err != nil {
        t.Fatalf("Failed to create HNSW: %v", err)
    }

    // Test vectors
    vectors := map[int][]float64{
        1: {1.0, 1.0, 1.0},
        2: {2.0, 2.0, 2.0},
        3: {3.0, 3.0, 3.0},
    }

    // Insert vectors
    for id, vector := range vectors {
        err := hnsw.Insert(id, vector)
        if err != nil {
            t.Errorf("Failed to insert vector %d: %v", id, err)
        }
    }

    // Test search
    query := []float64{1.1, 1.1, 1.1}
    k := 2
    results, distances := hnsw.Search(query, k)

    // Verify results
    if len(results) != k {
        t.Errorf("Expected %d results, got %d", k, len(results))
    }

    if results[0] != 1 { // Closest should be vector 1
        t.Errorf("Expected closest vector to be 1, got %d", results[0])
    }

    // Verify distances are sorted
    for i := 1; i < len(distances); i++ {
        if distances[i] < distances[i-1] {
            t.Error("Distances are not sorted")
        }
    }
}

func TestEmptySearch(t *testing.T) {
    cfg := config.NewDefaultConfig()
    hnsw, err := algorithm.New(cfg, distance.Euclidean)
    if err != nil {
        t.Fatalf("Failed to create HNSW: %v", err)
    }

    query := []float64{1.0, 1.0, 1.0}
    results, distances := hnsw.Search(query, 1)

    if len(results) != 0 || len(distances) != 0 {
        t.Error("Empty index should return empty results")
    }
}

func TestDuplicateInsert(t *testing.T) {
    cfg := config.NewDefaultConfig()
    hnsw, err := algorithm.New(cfg, distance.Euclidean)
    if err != nil {
        t.Fatalf("Failed to create HNSW: %v", err)
    }

    vector := []float64{1.0, 1.0, 1.0}
    
    // First insert should succeed
    err = hnsw.Insert(1, vector)
    if err != nil {
        t.Errorf("First insert failed: %v", err)
    }

    // Second insert should fail
    err = hnsw.Insert(1, vector)
    if err == nil {
        t.Error("Expected error on duplicate insert")
    }
}