package distance

import (
	"fmt"
	"math"
)

// DistanceFunction defines interface for distance calculation
type DistanceFunction func(a, b []float64) float64

// Available distance metrics
const (
	Euclidean  = "euclidean"
	Manhattan  = "manhattan"
	Cosine     = "cosine"
	DotProduct = "dot"
)

// GetDistanceFunction returns the corresponding distance function
func GetDistanceFunction(metric string) (DistanceFunction, error) {
	switch metric {
	case Euclidean:
		return EuclideanDistance, nil
	case Manhattan:
		return ManhattanDistance, nil
	case Cosine:
		return CosineDistance, nil
	case DotProduct:
		return DotProductDistance, nil
	default:
		return nil, fmt.Errorf("unsupported distance metric: %s", metric)
	}
}

// EuclideanDistance calculates Euclidean distance between vectors
func EuclideanDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return math.Inf(1)
	}

	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

// ManhattanDistance calculates Manhattan distance between vectors
func ManhattanDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return math.Inf(1)
	}

	sum := 0.0
	for i := range a {
		sum += math.Abs(a[i] - b[i])
	}
	return sum
}

// CosineDistance calculates Cosine distance between vectors
func CosineDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return math.Inf(1)
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return math.Inf(1)
	}

	similarity := dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
	if similarity > 1 {
		similarity = 1
	}

	return 1 - similarity
}

// DotProductDistance calculates negative dot product distance
func DotProductDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return math.Inf(1)
	}

	var dotProduct float64
	for i := range a {
		dotProduct += a[i] * b[i]
	}
	return -dotProduct
}

// ValidateVectors checks if vectors have same dimension
func ValidateVectors(a, b []float64) error {
	if len(a) != len(b) {
		return fmt.Errorf("vector dimensions mismatch: %d != %d", len(a), len(b))
	}
	return nil
}

// NormalizeVector normalizes vector to unit length
func NormalizeVector(v []float64) []float64 {
	norm := 0.0
	for _, val := range v {
		norm += val * val
	}
	norm = math.Sqrt(norm)

	if norm == 0 {
		return v
	}

	normalized := make([]float64, len(v))
	for i, val := range v {
		normalized[i] = val / norm
	}
	return normalized
}
