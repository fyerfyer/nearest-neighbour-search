# HNSW Test Suite

Unit tests for the Hierarchical Navigable Small World (HNSW) algorithm implementation.

## Test Files Structure

```
tests
├── README.md
├── core_test.go
├── neighbor_test.go
└── search_test.go
```


## Test Coverage

### Core Tests (`core_test.go`)
- Basic insertion and search
- Empty index handling
- Duplicate insertion prevention
- Configuration validation

### Neighbor Selection Tests (`neighbor_test.go`)
- Simple neighbor selection
  - Distance-based selection
  - M-nearest neighbors
- Heuristic neighbor selection
  - Candidate extension
  - Pruned connections
  - Level-wise selection

### Search Tests (`search_test.go`)
- K-nearest neighbor search
- Dimension mismatch handling
- Empty index search
- Distance ordering validation

## Running Tests

```bash
# Run all tests
go test ./tests

# Run specific test file
go test ./tests -run TestKNNSearch

# Run with verbose output
go test -v ./tests

# Run with coverage
go test -cover ./tests
```

## Performance Testing 

```go
go test -bench=. ./tests
```

## Test Data

The test suite uses synthetic data points in 2D and 3D space for validation:
- 2D points for search tests
- 3D points for neighbor selection
- Various edge cases (empty sets, dimension mismatches)

### Test Configuration

The test suite uses a default configuration for the HNSW index:
```go
config.Config{
    M:              16,   // Max connections
    MaxM:           32,   // Max construction connections
    EfConstruction: 100,  // Search list size
    ML:             1.0 / 1.5,
    DelayRebuild:   false,
}
```

### Example Test Case 
```go
func TestBasicSearch(t *testing.T) {
    cfg := config.NewDefaultConfig()
    hnsw, _ := algorithm.New(cfg, distance.Euclidean)
    
    // Insert test point
    hnsw.Insert(1, []float64{1.0, 1.0})
    
    // Search
    results, _ := hnsw.Search([]float64{1.1, 1.1}, 1)
    
    if len(results) != 1 {
        t.Error("Expected 1 result")
    }
}
```