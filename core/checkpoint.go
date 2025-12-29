package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/passgen/config"
)

// Checkpoint represents a generation checkpoint
type Checkpoint struct {
	Config        *config.Config `json:"config"`
	LinesWritten  int64          `json:"lines_written"`
	BytesWritten  int64          `json:"bytes_written"`
	LastPassword  string         `json:"last_password"`
	StateHash     string         `json:"state_hash"` // Hash of deduplication state
	mu            sync.Mutex
}

// CheckpointManager manages generation checkpoints
type CheckpointManager struct {
	checkpointDir string
}

// NewCheckpointManager creates a new checkpoint manager
func NewCheckpointManager() (*CheckpointManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	checkpointDir := filepath.Join(homeDir, ".passgen", "checkpoints")
	if err := os.MkdirAll(checkpointDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create checkpoint directory: %w", err)
	}

	return &CheckpointManager{checkpointDir: checkpointDir}, nil
}

// SaveCheckpoint saves a generation checkpoint
func (cm *CheckpointManager) SaveCheckpoint(jobID string, checkpoint *Checkpoint) error {
	checkpoint.mu.Lock()
	defer checkpoint.mu.Unlock()

	filePath := filepath.Join(cm.checkpointDir, jobID+".json")
	
	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write checkpoint: %w", err)
	}

	return nil
}

// LoadCheckpoint loads a generation checkpoint
func (cm *CheckpointManager) LoadCheckpoint(jobID string) (*Checkpoint, error) {
	filePath := filepath.Join(cm.checkpointDir, jobID+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint: %w", err)
	}

	checkpoint := &Checkpoint{}
	if err := json.Unmarshal(data, checkpoint); err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkpoint: %w", err)
	}

	return checkpoint, nil
}

// DeleteCheckpoint deletes a checkpoint
func (cm *CheckpointManager) DeleteCheckpoint(jobID string) error {
	filePath := filepath.Join(cm.checkpointDir, jobID+".json")
	return os.Remove(filePath)
}

// ListCheckpoints returns list of available checkpoints
func (cm *CheckpointManager) ListCheckpoints() ([]string, error) {
	files, err := os.ReadDir(cm.checkpointDir)
	if err != nil {
		return []string{}, nil
	}

	checkpoints := []string{}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			name := file.Name()[:len(file.Name())-5] // Remove .json
			checkpoints = append(checkpoints, name)
		}
	}

	return checkpoints, nil
}

