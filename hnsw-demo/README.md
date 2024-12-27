# hnsw-project README.md

# HNSW Nearest Neighbor Search

This project implements a Hierarchical Navigable Small World (HNSW) algorithm for efficient nearest neighbor search in high-dimensional spaces.

## Project Structure

```
hnsw-project
├── src
│   ├── core
│   │   ├── hnsw.go         # Main HNSW implementation
│   │   ├── node.go         # Node structure definition
│   │   ├── heap.go         # Priority queue implementation
│   │   └── distance.go     # Distance calculation functions
│   ├── search
│   │   ├── knn_search.go   # K-nearest neighbor search
│   │   └── layer_search.go  # Layer-wise search implementation
│   ├── neighbor
│   │   ├── select_simple.go    # Simple neighbor selection
│   │   └── select_heuristic.go  # Heuristic neighbor selection
│   └── utils
│       ├── config.go       # Configuration handling
│       ├── stats.go        # Statistics collection
│       └── persistence.go   # Save/Load functionality
├── tests
│   ├── core_test.go        # Tests for core functionalities
│   ├── search_test.go      # Tests for search functionalities
│   └── neighbor_test.go    # Tests for neighbor selection functionalities
├── examples
│   └── main.go             # Example usage of the HNSW implementation
├── go.mod                  # Go module file
├── go.sum                  # Go module dependencies
└── README.md               # Project documentation
```

## Features

- Efficient insertion and deletion of nodes.
- K-nearest neighbor search with both simple and heuristic selection methods.
- Layer-wise search capabilities.
- Configuration handling and statistics collection.

## Getting Started

1. Clone the repository:
   ```
   git clone <repository-url>
   cd hnsw-project
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run examples:
   ```
   go run examples/main.go
   ```

## Testing

To run the tests, use:
```
go test ./tests
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.