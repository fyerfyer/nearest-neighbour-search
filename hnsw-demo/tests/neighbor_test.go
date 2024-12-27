package tests

import (
    "testing"

    "github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/config"
    "github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/distance"
    "github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/src/algorithm"
)

func TestSelectNeighborsSimple(t *testing.T) {
    // Setup
    cfg := config.NewDefaultConfig()
    hnsw, err := algorithm.New(cfg, distance.Euclidean)
    if err != nil {
        t.Fatalf("Failed to create HNSW: %v", err)
    }

    // Insert test vectors
    vectors := map[int][]float64{
        1: {1.0, 0.0, 0.0},
        2: {0.0, 1.0, 0.0},
        3: {0.0, 0.0, 1.0},
        4: {1.0, 1.0, 0.0},
        5: {0.0, 1.0, 1.0},
    }

    for id, vec := range vectors {
        if err := hnsw.Insert(id, vec); err != nil {
            t.Fatalf("Failed to insert vector %d: %v", id, err)
        }
    }

    // Test cases
    tests := []struct {
        name       string
        query      []float64
        candidates []int
        M          int
        wantLen    int
    }{
        {
            name:       "Normal case",
            query:      []float64{1.0, 0.0, 0.0},
            candidates: []int{1, 2, 3, 4, 5},
            M:          3,
            wantLen:    3,
        },
        {
            name:       "M larger than candidates",
            query:      []float64{1.0, 0.0, 0.0},
            candidates: []int{1, 2},
            M:          3,
            wantLen:    2,
        },
        {
            name:       "Empty candidates",
            query:      []float64{1.0, 0.0, 0.0},
            candidates: []int{},
            M:          3,
            wantLen:    0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := hnsw.SelectNeighborsSimple(tt.query, tt.candidates, tt.M)
            if len(result) != tt.wantLen {
                t.Errorf("got len %d, want %d", len(result), tt.wantLen)
            }
        })
    }
}

func TestSelectNeighborsHeuristic(t *testing.T) {
    cfg := config.NewDefaultConfig()
    hnsw, err := algorithm.New(cfg, distance.Euclidean)
    if err != nil {
        t.Fatalf("Failed to create HNSW: %v", err)
    }

    // Insert test vectors
    vectors := map[int][]float64{
        1: {1.0, 0.0, 0.0},
        2: {0.0, 1.0, 0.0},
        3: {0.0, 0.0, 1.0},
        4: {1.0, 1.0, 0.0},
        5: {0.0, 1.0, 1.0},
    }

    for id, vec := range vectors {
        if err := hnsw.Insert(id, vec); err != nil {
            t.Fatalf("Failed to insert vector %d: %v", id, err)
        }
    }

    tests := []struct {
        name              string
        query            []float64
        candidates       []int
        M                int
        level            int
        extendCandidates bool
        keepPruned       bool
        wantLen          int
    }{
        {
            name:              "Basic heuristic selection",
            query:            []float64{1.0, 0.0, 0.0},
            candidates:       []int{1, 2, 3, 4, 5},
            M:                3,
            level:            0,
            extendCandidates: true,
            keepPruned:       true,
            wantLen:          3,
        },
        {
            name:              "Without extension",
            query:            []float64{1.0, 0.0, 0.0},
            candidates:       []int{1, 2, 3},
            M:                2,
            level:            0,
            extendCandidates: false,
            keepPruned:       false,
            wantLen:          2,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := hnsw.SelectNeighborsHeuristic(
                tt.query,
                tt.candidates,
                tt.M,
                tt.level,
                tt.extendCandidates,
                tt.keepPruned,
            )
            if len(result) != tt.wantLen {
                t.Errorf("got len %d, want %d", len(result), tt.wantLen)
            }
        })
    }
}