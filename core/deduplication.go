package core

import (
	"context"
	"crypto/sha256"
	"hash"
	"sync"

	"github.com/passgen/config"
)

// Deduplicator removes duplicate passwords using a Bloom filter approach
// For very large datasets, we use a hash set with limited memory
type Deduplicator struct {
	cfg      *config.DeduplicationConfig
	seen     map[string]bool
	seenHash map[[32]byte]bool // SHA256 hash for memory efficiency
	mu       sync.RWMutex
	maxSize  int // Maximum entries in memory
}

func NewDeduplicator(cfg *config.DeduplicationConfig) *Deduplicator {
	return &Deduplicator{
		cfg:      cfg,
		seen:     make(map[string]bool),
		seenHash: make(map[[32]byte]bool),
		maxSize:  1000000, // 1M entries in memory
	}
}

// Deduplicate filters out duplicate passwords
func (d *Deduplicator) Deduplicate(ctx context.Context, input <-chan string) <-chan string {
	output := make(chan string, 100)
	
	go func() {
		defer close(output)
		
		if !d.cfg.Enabled {
			// Pass through unchanged
			for word := range input {
				select {
				case <-ctx.Done():
					return
				case output <- word:
				}
			}
			return
		}
		
		hasher := sha256.New()
		
		for word := range input {
			select {
			case <-ctx.Done():
				return
			default:
				if d.isUnique(word, hasher) {
					output <- word
				}
			}
		}
	}()
	
	return output
}

func (d *Deduplicator) isUnique(word string, hasher hash.Hash) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// If we have space, use exact string matching
	if len(d.seen) < d.maxSize {
		if d.seen[word] {
			return false
		}
		d.seen[word] = true
		return true
	}
	
	// Otherwise, use hash-based deduplication
	hasher.Reset()
	hasher.Write([]byte(word))
	var hash [32]byte
	hasher.Sum(hash[:0])
	
	if d.seenHash[hash] {
		return false
	}
	
	d.seenHash[hash] = true
	return true
}

// Reset clears the deduplication state
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]bool)
	d.seenHash = make(map[[32]byte]bool)
}

