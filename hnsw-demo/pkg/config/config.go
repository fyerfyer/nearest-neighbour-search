package config

import (
	"fmt"
	"math"
)

// Config stores HNSW parameters
type Config struct {
	// Maximum number of connections per element in the graph
	M int

	// Maximum number of connections at construction time
	MaxM int

	// Size of the dynamic candidate list for construction
	EfConstruction int

	// Level generation parameter
	ML float64

	// Whether to delay index rebuilding after deletions
	DelayRebuild bool
}

// NewDefaultConfig creates a Config with default values
func NewDefaultConfig() Config {
	return Config{
		M:              16,
		MaxM:           32,
		EfConstruction: 100,
		ML:             1.0 / math.Log(16),
		DelayRebuild:   false,
	}
}

// NewConfig creates a Config with custom values
func NewConfig(m, maxM, efConstruction int, delayRebuild bool) (Config, error) {
	if m <= 0 {
		return Config{}, fmt.Errorf("M must be positive, got %d", m)
	}

	if maxM < m {
		return Config{}, fmt.Errorf("MaxM must be >= M, got M=%d, MaxM=%d", m, maxM)
	}

	if efConstruction <= 0 {
		return Config{}, fmt.Errorf("efConstruction must be positive, got %d", efConstruction)
	}

	return Config{
		M:              m,
		MaxM:           maxM,
		EfConstruction: efConstruction,
		ML:             1.0 / math.Log(float64(m)),
		DelayRebuild:   delayRebuild,
	}, nil
}

// Validate checks if the configuration is valid
func (c Config) Validate() error {
	if c.M <= 0 {
		return fmt.Errorf("M must be positive, got %d", c.M)
	}
	if c.MaxM < c.M {
		return fmt.Errorf("MaxM must be >= M, got M=%d, MaxM=%d", c.M, c.MaxM)
	}
	if c.EfConstruction <= 0 {
		return fmt.Errorf("efConstruction must be positive, got %d", c.EfConstruction)
	}
	if c.ML <= 0 {
		return fmt.Errorf("ML must be positive, got %f", c.ML)
	}
	return nil
}

// String returns a string representation of the config
func (c Config) String() string {
	return fmt.Sprintf("Config{M: %d, MaxM: %d, EfConstruction: %d, ML: %f, DelayRebuild: %v}",
		c.M, c.MaxM, c.EfConstruction, c.ML, c.DelayRebuild)
}
