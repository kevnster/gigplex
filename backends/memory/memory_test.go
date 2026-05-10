package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/kevnster/gigplex"
	"github.com/kevnster/gigplex/backends/memory"
)

func TestEnqueueDequeue(t *testing.T) {
	b := memory.New()
	ctx := context.Background()

	job := gigplex.Job{
		ID:      "job-1",
		Type:    "send-email",
		Payload: []byte(`{"to":"test@example.com"}`),
	}

	if err := b.Enqueue(ctx, job); err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	got, err := b.Dequeue(ctx)
	if err != nil {
		t.Fatalf("dequeue failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected a job, got nil")
	}
	if got.ID != job.ID {
		t.Fatalf("expected job ID %s, got %s", job.ID, got.ID)
	}
	if got.Status != gigplex.StatusProcessing {
		t.Fatalf("expected status processing, got %s", got.Status)
	}
}

func TestAcknowledge(t *testing.T) {
	b := memory.New()
	ctx := context.Background()

	b.Enqueue(ctx, gigplex.Job{ID: "job-2", Type: "resize-image"})
	b.Dequeue(ctx)

	if err := b.Acknowledge(ctx, "job-2"); err != nil {
		t.Fatalf("acknowledge failed: %v", err)
	}

	stats, _ := b.Stats(ctx)
	if stats.Done != 1 {
		t.Fatalf("expected 1 done, got %d", stats.Done)
	}
}

func TestFail(t *testing.T) {
	b := memory.New()
	ctx := context.Background()

	b.Enqueue(ctx, gigplex.Job{ID: "job-3", Type: "send-email"})
	b.Dequeue(ctx)
	b.Fail(ctx, "job-3", "smtp connection refused")

	stats, _ := b.Stats(ctx)
	if stats.Failed != 1 {
		t.Fatalf("expected 1 failed, got %d", stats.Failed)
	}
}

func TestHeartbeat(t *testing.T) {
	b := memory.New()
	ctx := context.Background()

	b.Heartbeat(ctx, "worker-1")
	workers, _ := b.Workers(ctx)

	if len(workers) != 1 {
		t.Fatalf("expected 1 worker, got %d", len(workers))
	}
	if time.Since(workers[0].LastBeat) > time.Second {
		t.Fatal("heartbeat timestamp is too old")
	}
}