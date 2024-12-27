package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/config"
	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/distance"
	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/src/algorithm"
)

func main() {
	// Create configuration
	cfg := config.Config{
		M:              16,
		MaxM:           32,
		EfConstruction: 100,
		ML:             1.0 / float64(16),
		DelayRebuild:   false,
	}

	// Create HNSW index
	index, err := algorithm.New(cfg, distance.Euclidean)
	if err != nil {
		panic(fmt.Sprintf("Failed to create index: %v", err))
	}

	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Generate random vectors
	dim := 3
	numVectors := 1000
	vectors := make(map[int][]float64)

	for i := 0; i < numVectors; i++ {
		vec := make([]float64, dim)
		for j := 0; j < dim; j++ {
			vec[j] = rand.Float64()
		}
		vectors[i] = vec
	}

	// Insert vectors
	fmt.Println("Inserting vectors...")
	start := time.Now()
	for id, vec := range vectors {
		if err := index.Insert(id, vec); err != nil {
			fmt.Printf("Failed to insert vector %d: %v\n", id, err)
		}
	}
	fmt.Printf("Insertion took: %v\n", time.Since(start))

	// Perform search
	fmt.Println("\nPerforming search...")
	query := make([]float64, dim)
	for i := 0; i < dim; i++ {
		query[i] = rand.Float64()
	}

	k := 10
	ef := 50
	start = time.Now()
	results, distances := index.KNNSearchWithDistances(query, k, ef)
	searchTime := time.Since(start)

	// Print results
	fmt.Printf("\nSearch took: %v\n", searchTime)
	fmt.Printf("Query vector: %v\n", query)
	fmt.Printf("\nNearest %d neighbors:\n", k)
	for i, id := range results {
		fmt.Printf("%d. ID: %d, Distance: %.4f, Vector: %v\n",
			i+1, id, distances[i], vectors[id])
	}
}
