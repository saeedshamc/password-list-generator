package core

import (
	"context"
	"runtime"
	"sync"
)

// ParallelProcessor processes passwords in parallel using multiple workers
type ParallelProcessor struct {
	workers int
}

// NewParallelProcessor creates a new parallel processor
func NewParallelProcessor() *ParallelProcessor {
	workers := runtime.NumCPU()
	if workers > 8 {
		workers = 8 // Cap at 8 workers to avoid excessive overhead
	}
	return &ParallelProcessor{workers: workers}
}

// ProcessParallel processes passwords in parallel batches
func (pp *ParallelProcessor) ProcessParallel(ctx context.Context, input <-chan string, processor func(string) string) <-chan string {
	output := make(chan string, 100*pp.workers)
	
	var wg sync.WaitGroup
	
	// Start worker goroutines
	for i := 0; i < pp.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for word := range input {
				select {
				case <-ctx.Done():
					return
				default:
					processed := processor(word)
					if processed != "" {
						select {
						case <-ctx.Done():
							return
						case output <- processed:
						}
					}
				}
			}
		}()
	}
	
	// Close output when all workers are done
	go func() {
		wg.Wait()
		close(output)
	}()
	
	return output
}

// GetWorkerCount returns the number of workers
func (pp *ParallelProcessor) GetWorkerCount() int {
	return pp.workers
}

// SetWorkerCount sets the number of workers
func (pp *ParallelProcessor) SetWorkerCount(count int) {
	if count > 0 && count <= 16 {
		pp.workers = count
	}
}

