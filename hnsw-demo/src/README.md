# HNSW Implementation

A Go implementation of the Hierarchical Navigable Small World algorithm for approximate nearest neighbor search, based on the paper ["Efficient and robust approximate nearest neighbor search using Hierarchical Navigable Small World graphs"](https://arxiv.org/pdf/1603.09320.pdf).

## Features

### Core Algorithms
- **HNSW Graph Construction**: Implements multi-layer graph construction with randomized level generation
- **K-Nearest Neighbor Search**: Two search strategies:
  - Basic KNN search with distance calculations
  - Layer-wise search optimization

### Neighbor Selection Methods
1. **Simple Selection** ([`SelectNeighborsSimple`](src/algorithm/neighbour_simple.go))
   - Basic distance-based neighbor selection
   - Returns M closest neighbors based on distance

2. **Heuristic Selection** ([`SelectNeighborsHeuristic`](src/algorithm/neighbour_heuristic.go))
   - Advanced selection with pruning strategies
   - Supports candidate extension and pruned connection preservation
   - Better connection quality for search performance

### Distance Metrics
- Euclidean distance
- Manhattan distance
- Cosine similarity
- Dot product

## Usage

### Creating an Index

```go
// Create configuration
cfg := config.Config{
    M:              16,  // Max number of connections per node
    MaxM:           32,  // Max connections for upper layers
    EfConstruction: 100, // Size of dynamic candidate list
    ML:             1.0 / float64(16),
    DelayRebuild:   false,
}

// Initialize HNSW index
index, err := algorithm.New(cfg, distance.Euclidean)
if err != nil {
    panic(err)
}
```

### Inserting Elements

```go
// Insert a vector with ID
vector := []float64{1.0, 2.0, 3.0}
err := index.Insert(1, vector)
if err != nil {
    log.Fatal(err)
}
```

### Searching for Nearest Neighbors

```go
// Perform K-NN search
query := []float64{1.0, 2.0, 3.0}
k := 10      // Number of nearest neighbors
ef := 50     // Search expansion factor

// Get results with distances
results, distances := index.KNNSearchWithDistances(query, k, ef)
```

## Performance Considerations

1. Layer Generation
    * Uses probabilistic level assignment
    * Higher levels have fewer nodes for efficient navigation

2. Search Optimization
    * Layer-wise search from top to bottom
    * Beam search with dynamic candidate lists
    * Early termination based on distance bounds

3. Neighbor Selection
    * Heuristic method provides better graph quality
    * Supports pruning to maintain connection diversity
    * Optional extension of candidate set

## Testing

Run the test suite:

```bash
go test ./tests
```

## Appendix

### Pseudo-code In the Paper

* **Algorithm 1**: HNSW Construction
```
W ← ∅ // Current list of nearest neighbor elements
ep ← get enter point for hnsw
L ← level of ep // Top level of hnsw
l ← ⌊-ln(unif(0..1))∙mL⌋ // Level of the new element
for lc ← L … l+1
    W ← SEARCH-LAYER(q, ep, ef=1, lc)
    ep ← get the nearest element from W to q
for lc ← min(L, l) … 0
    W ← SEARCH-LAYER(q, ep, efConstruction, lc)
    neighbors ← SELECT-NEIGHBORS(q, W, M, lc) // Algorithm 3 or Algorithm 4
    add bidirectional connections from neighbors to q at layer lc
    for each e ∈ neighbors // Shrink connections if necessary
        eConn ← neighbourhood(e) at layer lc 
        if │eConn│ > Mmax // Shrink connections of e
            eNewConn ← SELECT-NEIGHBORS(e, eConn, Mmax, lc) // Algorithm 3 or Algorithm 4
            set neighbourhood(e) at layer lc to eNewConn
    ep ← W
if l > L
    set enter point for hnsw to q
```

* **Algorithm 2**: Search Layer
```
SEARCH-LAYER(q, ep, ef, lc)
Input: query element q, enter points ep, number of nearest to q elements to return ef, layer number lc
Output: ef closest neighbors to q

v ← ep // Set of visited elements
C ← ep // Set of candidates
W ← ep // Dynamic list of found nearest neighbors

while │C│ > 0
    c ← extract nearest element from C to q
    f ← get furthest element from W to q
    if distance(c, q) > distance(f, q)
        break // All elements in W are evaluated
    for each e ∈ neighbourhood(c) at layer lc // Update C and W
        if e ∉ v
            v ← v ⋃ e
        f ← get furthest element from W to q
        if distance(e, q) < distance(f, q) or │W│ < ef
            C ← C ⋃ e
            W ← W ⋃ e
        if │W│ > ef
            remove furthest element from W to q
return W
```

* **Algorithm 3**: Select Neighbors Simple
```
SELECT-NEIGHBORS-SIMPLE(q, C, M)
Input: base element q, candidate elements C, number of neighbors to return M
Output: M nearest elements to q

return M nearest elements from C to q
```

* **Algorithm 4**: Select Neighbors Heuristic
```
SELECT-NEIGHBORS-HEURISTIC(q, C, M, lc, extendCandidates, keepPrunedConnections)
Input: base element q, candidate elements C, number of neighbors to return M, layer number lc, flag indicating whether or not to extend candidate list extendCandidates, flag indicating whether or not to add discarded elements keepPrunedConnections
Output: M elements selected by the heuristic

R ← ∅
W ← C // Working queue for the candidates

if extendCandidates // Extend candidates by their neighbors
    for each e ∈ C
        for each eadj ∈ neighbourhood(e) at layer lc
            if eadj ∉ W
                W ← W ⋃ eadj

Wd ← ∅ // Queue for the discarded candidates

while │W│ > 0 and │R│ < M
    e ← extract nearest element from W to q
    if e is closer to q compared to any element from R
        R ← R ⋃ e
    else
        Wd ← Wd ⋃ e

if keepPrunedConnections // Add some of the discarded connections from Wd
    while │Wd│ > 0 and │R│ < M
        R ← R ⋃ extract nearest element from Wd to q

return R
```

* **Algorithm 5**: K-Nearest Neighbor Search
```
K-NN-SEARCH(hnsw, q, K, ef)
Input: multilayer graph hnsw, query element q, number of nearest neighbors to return K, size of the dynamic candidate list ef
Output: K nearest elements to q

W ← ∅ // Set for the current nearest elements
ep ← get enter point for hnsw
L ← level of ep // Top layer for hnsw

for lc ← L … 1
    W ← SEARCH-LAYER(q, ep, ef=1, lc)
    ep ← get nearest element from W to q

W ← SEARCH-LAYER(q, ep, ef, lc=0)

return K nearest elements from W to q
```