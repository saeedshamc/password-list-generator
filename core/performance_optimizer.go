package core

import (
	"context"
	"runtime"
	"sync"
)

// PerformanceOptimizer optimizes generation performance
type PerformanceOptimizer struct {
	maxWorkers int
	bufferSize int
}

// NewPerformanceOptimizer creates a new performance optimizer
func NewPerformanceOptimizer() *PerformanceOptimizer {
	// Use CPU count for optimal performance
	numCPU := runtime.NumCPU()
	maxWorkers := numCPU
	if maxWorkers > 8 {
		maxWorkers = 8 // Cap at 8 to avoid overhead
	}
	if maxWorkers < 2 {
		maxWorkers = 2 // At least 2 workers
	}

	return &PerformanceOptimizer{
		maxWorkers: maxWorkers,
		bufferSize: 1000, // Larger buffer for better throughput
	}
}

// OptimizeChannel optimizes a channel with better buffering
func (p *PerformanceOptimizer) OptimizeChannel(input <-chan string) <-chan string {
	output := make(chan string, p.bufferSize)

	go func() {
		defer close(output)
		for word := range input {
			output <- word
		}
	}()

	return output
}

// ParallelProcess processes passwords in parallel batches
func (p *PerformanceOptimizer) ParallelProcess(ctx context.Context, input <-chan string, processor func(string) string) <-chan string {
	output := make(chan string, p.bufferSize*p.maxWorkers)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, p.maxWorkers)

	go func() {
		defer close(output)

		for word := range input {
			select {
			case <-ctx.Done():
				return
			default:
				semaphore <- struct{}{}
				wg.Add(1)

				go func(w string) {
					defer func() {
						<-semaphore
						wg.Done()
					}()

					processed := processor(w)
					if processed != "" {
						select {
						case <-ctx.Done():
							return
						case output <- processed:
						}
					}
				}(word)
			}
		}

		wg.Wait()
	}()

	return output
}

// GetOptimalBufferSize returns optimal buffer size based on system
func GetOptimalBufferSize() int {
	numCPU := runtime.NumCPU()
	// Larger buffer for systems with more CPUs
	return 100 * numCPU
}

// GetOptimalWorkerCount returns optimal worker count
func GetOptimalWorkerCount() int {
	numCPU := runtime.NumCPU()
	if numCPU > 8 {
		return 8
	}
	if numCPU < 2 {
		return 2
	}
	return numCPU
}

