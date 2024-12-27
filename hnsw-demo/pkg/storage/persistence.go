package storage

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/config"
	"github.com/fyerfyer/nearest-neighbour-search/hnsw-demo/pkg/node"
)

// IndexMetadata stores metadata about the saved index
type IndexMetadata struct {
	Version     string        // Version of the index format
	CreatedAt   time.Time     // When the index was created
	NodesCount  int           // Number of nodes in the index
	MaxLevel    int           // Maximum level in the index
	Config      config.Config // Index configuration
	Description string        // Optional description
}

// SaveData represents the complete state of the index
type SaveData struct {
	Metadata   IndexMetadata
	Nodes      map[int]*node.Node
	EntryPoint int
}

// SaveIndex saves the index state to a file
func SaveIndex(filename string, nodes map[int]*node.Node, entryPoint int,
	maxLevel int, cfg config.Config, description string) error {

	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Prepare metadata
	metadata := IndexMetadata{
		Version:     "1.0",
		CreatedAt:   time.Now(),
		NodesCount:  len(nodes),
		MaxLevel:    maxLevel,
		Config:      cfg,
		Description: description,
	}

	// Prepare save data
	data := SaveData{
		Metadata:   metadata,
		Nodes:      nodes,
		EntryPoint: entryPoint,
	}

	// Create encoder and encode data
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode data: %v", err)
	}

	return nil
}

// LoadIndex loads the index state from a file
func LoadIndex(filename string) (map[int]*node.Node, int, config.Config, error) {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, config.Config{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create decoder
	decoder := gob.NewDecoder(file)

	// Decode data
	var data SaveData
	if err := decoder.Decode(&data); err != nil {
		return nil, 0, config.Config{}, fmt.Errorf("failed to decode data: %v", err)
	}

	// Validate loaded data
	if err := validateLoadedData(&data); err != nil {
		return nil, 0, config.Config{}, fmt.Errorf("invalid data: %v", err)
	}

	return data.Nodes, data.EntryPoint, data.Metadata.Config, nil
}

// CreateBackup creates a backup of the index file
func CreateBackup(sourceFile string) error {
	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102150405")
	backupFile := fmt.Sprintf("%s.%s.backup", sourceFile, timestamp)

	// Copy source file to backup
	source, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer source.Close()

	destination, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer destination.Close()

	if _, err := destination.ReadFrom(source); err != nil {
		return fmt.Errorf("failed to copy data: %v", err)
	}

	return nil
}

// validateLoadedData performs validation checks on loaded data
func validateLoadedData(data *SaveData) error {
	// Check version compatibility
	if data.Metadata.Version != "1.0" {
		return fmt.Errorf("unsupported index version: %s", data.Metadata.Version)
	}

	// Validate node count
	if len(data.Nodes) != data.Metadata.NodesCount {
		return fmt.Errorf("node count mismatch: metadata=%d, actual=%d",
			data.Metadata.NodesCount, len(data.Nodes))
	}

	// Validate config
	if err := data.Metadata.Config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	// Validate entry point
	if _, exists := data.Nodes[data.EntryPoint]; !exists && len(data.Nodes) > 0 {
		return fmt.Errorf("invalid entry point: %d", data.EntryPoint)
	}

	return nil
}

// GetIndexInfo returns information about a saved index
func GetIndexInfo(filename string) (*IndexMetadata, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var data SaveData
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode data: %v", err)
	}

	return &data.Metadata, nil
}
