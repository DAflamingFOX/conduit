package queue

import (
	"context"
	"time"

	"github.com/daflamingfox/conduit/internal/engine"
)

// JobLogEntry represents a single log line emitted during a node execution step.
type JobLogEntry struct {
	JobID     string    `json:"job_id"`
	NodeID    string    `json:"node_id"`
	Level     string    `json:"level"` // "info", "warn", "error"
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// TaskQueue defines the abstraction interface for job task queueing and log streaming.
type TaskQueue interface {
	// PushStep enqueues an execution step onto the pending task queue.
	PushStep(ctx context.Context, step *engine.ExecutionStep) error

	// PopStep claims and pops the next pending execution step for a worker (blocking wait).
	PopStep(ctx context.Context, workerID string, timeout time.Duration) (*engine.ExecutionStep, error)

	// PublishLog broadcasts a log line entry over the pub/sub event bus.
	PublishLog(ctx context.Context, entry JobLogEntry) error

	// Close cleanly closes the queue connection.
	Close()
}
