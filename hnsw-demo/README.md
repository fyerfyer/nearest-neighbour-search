# HNSW Nearest Neighbor Search

This project implements a Hierarchical Navigable Small World (HNSW) algorithm for efficient nearest neighbor search in high-dimensional spaces. The algorithm are from the paper "Efficient and Robust Approximate Nearest Neighbor Search using Hierarchical Navigable Small World Graphs" by Yu. A. Malkov and D. A. Yashunin.

## Project Structure

```
hnsw-demo
├── main
│   └── main.go             # Entry point for the demo application
├── pkg 
│   ├── config              # Configuration handling
│   ├── distance            # Distance metric implementations  
│   ├── heap                # Priority queue implementation
│   ├── node                # Node data structure
│   └── storage             # Persistence layer
├── src
│   └── algorithm           # HNSW algorithm implementation
├── tests
│   ├── core_test.go        # Core algorithm tests
│   ├── neighbor_test.go    # Neighbor search tests
│   └── search_test.go      # Search tests
├── README.md               # Project README
└── go.mod                  # Go module definition
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
   cd hnsw-demo
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run examples:
   ```
   go run main/main.go
   ```

## Testing

To run the tests, use:
```
go test ./tests
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.
