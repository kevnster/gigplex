package gigplex

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Config holds the configuration for a Gigplex instance.
type Config struct {
	Backend     Backend
	Workers     int
	RetryPolicy RetryPolicy
}

// Gigplex is the main entry point for the library.
type Gigplex struct {
	config   Config
	handlers map[string]Handler
	mu       sync.RWMutex
}

// Handler is a function that processes a job payload.
type Handler func(ctx context.Context, payload []byte) error

// New creates a new Gigplex instance with the given config.
func New(cfg Config) *Gigplex {
	if cfg.Workers == 0 {
		cfg.Workers = 3
	}
	if cfg.RetryPolicy == nil {
		cfg.RetryPolicy = DefaultRetry
	}
	return &Gigplex{
		config:   cfg,
		handlers: make(map[string]Handler),
	}
}

// Register associates a handler function with a job type.
func (g *Gigplex) Register(jobType string, h Handler) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.handlers[jobType] = h
}

// Enqueue adds a new job to the queue.
func (g *Gigplex) Enqueue(ctx context.Context, jobType string, payload []byte) error {
	job := Job{
		ID:        fmt.Sprintf("%s-%d", jobType, time.Now().UnixNano()),
		Type:      jobType,
		Payload:   payload,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
	return g.config.Backend.Enqueue(ctx, job)
}

// Start launches the worker pool. Blocks until the context is cancelled.
func (g *Gigplex) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < g.config.Workers; i++ {
		wg.Add(1)
		workerID := fmt.Sprintf("worker-%d", i+1)
		go func(id string) {
			defer wg.Done()
			g.runWorker(ctx, id)
		}(workerID)
	}

	wg.Wait()
}

// runWorker is the main loop for a single worker goroutine.
func (g *Gigplex) runWorker(ctx context.Context, workerID string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			job, err := g.config.Backend.Dequeue(ctx)
			if err != nil {
				log.Printf("[%s] dequeue error: %v", workerID, err)
				time.Sleep(time.Second)
				continue
			}

			// no jobs available — wait a bit before polling again
			if job == nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			g.processJob(ctx, workerID, job)
		}
	}
}

// processJob runs the handler for a job and handles retries.
func (g *Gigplex) processJob(ctx context.Context, workerID string, job *Job) {
	g.mu.RLock()
	handler, ok := g.handlers[job.Type]
	g.mu.RUnlock()

	if !ok {
		log.Printf("[%s] no handler registered for job type: %s", workerID, job.Type)
		g.config.Backend.Fail(ctx, job.ID, "no handler registered")
		return
	}

	log.Printf("[%s] processing job %s (type: %s, attempt: %d)", workerID, job.ID, job.Type, job.Attempt+1)

	err := handler(ctx, job.Payload)
	if err != nil {
		job.Attempt++
		if g.config.RetryPolicy.ShouldRetry(job.Attempt, err) {
			backoff := g.config.RetryPolicy.Backoff(job.Attempt)
			log.Printf("[%s] job %s failed, retrying in %s: %v", workerID, job.ID, backoff, err)
			time.Sleep(backoff)
			// re-enqueue for retry
			job.Status = StatusPending
			g.config.Backend.Enqueue(ctx, *job)
			return
		}
		log.Printf("[%s] job %s failed permanently: %v", workerID, job.ID, err)
		g.config.Backend.Fail(ctx, job.ID, err.Error())
		return
	}

	log.Printf("[%s] job %s done", workerID, job.ID)
	g.config.Backend.Acknowledge(ctx, job.ID)
}