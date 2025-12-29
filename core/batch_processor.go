package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/passgen/config"
)

// BatchJob represents a single batch generation job
type BatchJob struct {
	ID          string
	Name        string
	Config      *config.Config
	Status      string // "pending", "running", "completed", "failed", "paused"
	Progress    float64
	LinesWritten int64
	BytesWritten int64
	Error       error
	mu          sync.Mutex
}

// BatchProcessor manages multiple generation jobs
type BatchProcessor struct {
	jobs    map[string]*BatchJob
	mu      sync.RWMutex
	running int
	maxJobs int // Maximum concurrent jobs
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(maxConcurrentJobs int) *BatchProcessor {
	if maxConcurrentJobs <= 0 {
		maxConcurrentJobs = 1 // Default to 1 job at a time
	}
	return &BatchProcessor{
		jobs:    make(map[string]*BatchJob),
		maxJobs: maxConcurrentJobs,
	}
}

// AddJob adds a new job to the batch
func (bp *BatchProcessor) AddJob(id, name string, cfg *config.Config) *BatchJob {
	bp.mu.Lock()
	defer bp.mu.Unlock()

	job := &BatchJob{
		ID:     id,
		Name:   name,
		Config: cfg,
		Status: "pending",
	}

	bp.jobs[id] = job
	return job
}

// GetJob returns a job by ID
func (bp *BatchProcessor) GetJob(id string) (*BatchJob, bool) {
	bp.mu.RLock()
	defer bp.mu.RUnlock()
	job, ok := bp.jobs[id]
	return job, ok
}

// ListJobs returns all jobs
func (bp *BatchProcessor) ListJobs() []*BatchJob {
	bp.mu.RLock()
	defer bp.mu.RUnlock()

	jobs := make([]*BatchJob, 0, len(bp.jobs))
	for _, job := range bp.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// ProcessBatch processes all pending jobs
func (bp *BatchProcessor) ProcessBatch(ctx context.Context, progressCallback func(jobID string, lines, bytes int64)) error {
	// Process jobs with concurrency limit
	semaphore := make(chan struct{}, bp.maxJobs)
	var wg sync.WaitGroup

	for _, job := range bp.ListJobs() {
		if job.Status != "pending" {
			continue
		}

		// Wait for available slot
		semaphore <- struct{}{}
		wg.Add(1)

		go func(j *BatchJob) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			bp.processJob(ctx, j, progressCallback)
		}(job)
	}

	wg.Wait()
	return nil
}

// processJob processes a single job
func (bp *BatchProcessor) processJob(ctx context.Context, job *BatchJob, progressCallback func(jobID string, lines, bytes int64)) {
	job.mu.Lock()
	job.Status = "running"
	job.mu.Unlock()

	generator := NewGenerator(job.Config)
	generator.SetJobID(job.ID)

	err := generator.Generate(ctx, func(lines, bytes int64) {
		job.mu.Lock()
		job.LinesWritten = lines
		job.BytesWritten = bytes
		if job.Config.Limits.MaxOutputLines > 0 {
			job.Progress = float64(lines) / float64(job.Config.Limits.MaxOutputLines)
		}
		job.mu.Unlock()

		if progressCallback != nil {
			progressCallback(job.ID, lines, bytes)
		}
	})

	job.mu.Lock()
	if err != nil {
		if err.Error() == "generation paused" {
			job.Status = "paused"
		} else {
			job.Status = "failed"
			job.Error = err
		}
	} else {
		job.Status = "completed"
		job.Progress = 1.0
	}
	job.mu.Unlock()
}

// CancelJob cancels a running job
func (bp *BatchProcessor) CancelJob(jobID string) {
	job, ok := bp.GetJob(jobID)
	if !ok {
		return
	}

	job.mu.Lock()
	if job.Status == "running" {
		job.Status = "cancelled"
	}
	job.mu.Unlock()
}

// PauseJob pauses a running job
func (bp *BatchProcessor) PauseJob(jobID string) error {
	job, ok := bp.GetJob(jobID)
	if !ok {
		return fmt.Errorf("job not found")
	}

	job.mu.Lock()
	defer job.mu.Unlock()

	if job.Status != "running" {
		return fmt.Errorf("job is not running")
	}

	job.Status = "paused"
	return nil
}

// ResumeJob resumes a paused job
func (bp *BatchProcessor) ResumeJob(jobID string) error {
	job, ok := bp.GetJob(jobID)
	if !ok {
		return fmt.Errorf("job not found")
	}

	job.mu.Lock()
	defer job.mu.Unlock()

	if job.Status != "paused" {
		return fmt.Errorf("job is not paused")
	}

	job.Status = "pending"
	return nil
}

