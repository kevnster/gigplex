package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kevnster/gigplex"
)

type Backend struct {
	mu       sync.Mutex
	jobs     map[string]*gigplex.Job
	queue    []string // ordered list of pending job IDs
	workers  map[string]*gigplex.WorkerInfo
}

func New() *Backend {
	return &Backend{
		jobs:    make(map[string]*gigplex.Job),
		queue:   []string{},
		workers: make(map[string]*gigplex.WorkerInfo),
	}
}

func (b *Backend) Enqueue(_ context.Context, job gigplex.Job) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	job.Status = gigplex.StatusPending
	job.CreatedAt = time.Now()
	b.jobs[job.ID] = &job
	b.queue = append(b.queue, job.ID)
	return nil
}

func (b *Backend) Dequeue(_ context.Context) (*gigplex.Job, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// find the first pending job in the queue
	for i, id := range b.queue {
		job := b.jobs[id]
		if job.Status == gigplex.StatusPending {
			job.Status = gigplex.StatusProcessing
			job.UpdatedAt = time.Now()
			b.queue = append(b.queue[:i], b.queue[i+1:]...)
			return job, nil
		}
	}
	return nil, nil // nothing to process
}

func (b *Backend) Acknowledge(_ context.Context, jobID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	job, ok := b.jobs[jobID]
	if !ok {
		return fmt.Errorf("job %s not found", jobID)
	}
	job.Status = gigplex.StatusDone
	job.UpdatedAt = time.Now()
	return nil
}

func (b *Backend) Fail(_ context.Context, jobID string, reason string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	job, ok := b.jobs[jobID]
	if !ok {
		return fmt.Errorf("job %s not found", jobID)
	}
	job.Status = gigplex.StatusFailed
	job.Error = reason
	job.UpdatedAt = time.Now()
	return nil
}

func (b *Backend) Heartbeat(_ context.Context, workerID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if w, ok := b.workers[workerID]; ok {
		w.LastBeat = time.Now()
		return nil
	}
	// first heartbeat — register the worker
	b.workers[workerID] = &gigplex.WorkerInfo{
		ID:       workerID,
		LastBeat: time.Now(),
	}
	return nil
}

func (b *Backend) Workers(_ context.Context) ([]gigplex.WorkerInfo, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	result := make([]gigplex.WorkerInfo, 0, len(b.workers))
	for _, w := range b.workers {
		result = append(result, *w)
	}
	return result, nil
}

func (b *Backend) Stats(_ context.Context) (gigplex.Stats, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	var stats gigplex.Stats
	for _, job := range b.jobs {
		switch job.Status {
		case gigplex.StatusPending:
			stats.Pending++
		case gigplex.StatusProcessing:
			stats.Processing++
		case gigplex.StatusDone:
			stats.Done++
		case gigplex.StatusFailed:
			stats.Failed++
		}
	}
	return stats, nil
}