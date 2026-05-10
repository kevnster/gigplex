package gigplex

import (
	"context"
	"time"
)

// Backend is the only interface you need to implement to use gigplex
// with any storage system — Redis, AWS, Postgres, or your own.
type Backend interface {
	Enqueue(ctx context.Context, job Job) error
	Dequeue(ctx context.Context) (*Job, error)
	Acknowledge(ctx context.Context, jobID string) error
	Fail(ctx context.Context, jobID string, reason string) error
	Heartbeat(ctx context.Context, workerID string) error
	Workers(ctx context.Context) ([]WorkerInfo, error)
	Stats(ctx context.Context) (Stats, error)
	RecentJobs(ctx context.Context, limit int) ([]Job, error)
	KillWorker(ctx context.Context, workerID string) error
	RetryFailed(ctx context.Context) error
}

type WorkerInfo struct {
	ID       string
	LastBeat time.Time
	JobsDone int
	InFlight int
	IsLeader bool
}

type Stats struct {
	Pending    int
	Processing int
	Done       int
	Failed     int
}
