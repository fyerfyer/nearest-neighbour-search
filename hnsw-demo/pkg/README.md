# HNSW Packages

Core packages for the Hierarchical Navigable Small World (HNSW) implementation.

## Package Structure

```
pkg
├── config
│   └── config.go
├── distance
│   └── metric.go
├── heap
│   └── priority_queue.go
├── node
│   └── node.go
├── storage
│   └── persistence.go
└── README.md
```


## Core Components

### config
- `config` - HNSW parameters configuration
- Default settings:
  - M: 16 (max connections per node)
  - MaxM: 32 (max connections during construction)
  - EfConstruction: 100 (dynamic candidate list size)

```go
cfg := config.NewDefaultConfig()
// or
cfg, err := config.NewConfig(16, 32, 100, false)
```

### distance
Supported metrics:

- Euclidean
- Manhattan
- Cosine
- Dot Product

```go
distFunc, err := distance.GetDistanceFunction("euclidean")
dist := distFunc(vector1, vector2)
```

### node 

Graph node implementation with thread-safe operations:

- Neighbor management
- Soft deletion support
- Vector storage

```go
node := node.NewNode(id, vector, level)
node.AddNeighbor(level, neighborID)
```

### heap
Priority queue for efficient nearest neighbor search:
```go
pq := heap.NewPriorityQueue()
pq.PushItem(nodeID, distance)
```

### storage
Index persistence with metadata:

```go
// Save index
err := storage.SaveIndex("index.hnsw", nodes, entryPoint, maxLevel, cfg, "")
// Load index
nodes, entryPoint, cfg, err := storage.LoadIndex("index.hnsw")
```

### Configuration Options

```go
type Config struct {
    M              int     // Max connections per element
    MaxM           int     // Max connections at construction
    EfConstruction int     // Dynamic candidate list size
    ML             float64 // Level generation parameter
    DelayRebuild   bool    // Delayed index rebuilding flag
}
```

### Thread Safety
- All node operations are protected by RWMutex
- Safe for concurrent search operations
- Write operations (insert/delete) should be synchronized externally

### Dependencies
- Standard library only
- No external dependencies required

### Error Handling
- All operations return error types when applicable
- Validation for vector dimensions, configuration parameters
- Thread-safe operations with proper locking