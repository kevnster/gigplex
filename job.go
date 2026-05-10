package gigplex

import "time"

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusDone       Status = "done"
	StatusFailed     Status = "failed"
)

type Job struct {
	ID        string
	Type      string
	Payload   []byte
	Status    Status
	Attempt   int
	CreatedAt time.Time
	Error     string
	UpdatedAt time.Time
}