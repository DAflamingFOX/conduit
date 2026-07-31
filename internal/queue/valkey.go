package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/daflamingfox/conduit/internal/engine"
	"github.com/valkey-io/valkey-go"
)

const (
	PendingTasksKey = "conduit:tasks:pending"
	LogChannelKey   = "conduit:logs:stream"
)

// ValkeyQueue implements TaskQueue using the official Valkey Go client.
type ValkeyQueue struct {
	client valkey.Client
}

// NewValkeyQueue initializes a Valkey task queue connection.
func NewValkeyQueue(addr string) (*ValkeyQueue, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Valkey at %s: %w", addr, err)
	}

	// Verify connection with Ping
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping Valkey at %s: %w", addr, err)
	}

	slog.Info("connected to Valkey task queue", "address", addr)
	return &ValkeyQueue{client: client}, nil
}

// PushStep serializes and pushes an execution step to the pending list queue.
func (vq *ValkeyQueue) PushStep(ctx context.Context, step *engine.ExecutionStep) error {
	data, err := json.Marshal(step)
	if err != nil {
		return fmt.Errorf("failed to marshal execution step: %w", err)
	}

	cmd := vq.client.B().Rpush().Key(PendingTasksKey).Element(string(data)).Build()
	if err := vq.client.Do(ctx, cmd).Error(); err != nil {
		return fmt.Errorf("failed to push task to Valkey queue: %w", err)
	}

	slog.Debug("enqueued execution step", "step_id", step.StepID, "node_instance_id", step.NodeInstanceID)
	return nil
}

// PopStep pops and claims the next available pending execution step (blocking with timeout).
func (vq *ValkeyQueue) PopStep(ctx context.Context, workerID string, timeout time.Duration) (*engine.ExecutionStep, error) {
	secs := float64(timeout.Seconds())
	if secs <= 0 {
		secs = 1.0
	}

	cmd := vq.client.B().Blpop().Key(PendingTasksKey).Timeout(secs).Build()
	resp := vq.client.Do(ctx, cmd)
	if err := resp.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return nil, nil // Timeout, no pending task
		}
		return nil, fmt.Errorf("error popping step from Valkey: %w", err)
	}

	values, err := resp.AsStrSlice()
	if err != nil || len(values) < 2 {
		return nil, nil
	}

	var step engine.ExecutionStep
	if err := json.Unmarshal([]byte(values[1]), &step); err != nil {
		return nil, fmt.Errorf("failed to unmarshal popped execution step: %w", err)
	}

	slog.Debug("worker claimed execution step", "worker_id", workerID, "step_id", step.StepID)
	return &step, nil
}

// PublishLog broadcasts a log entry over Valkey Pub/Sub channel.
func (vq *ValkeyQueue) PublishLog(ctx context.Context, entry JobLogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	cmd := vq.client.B().Publish().Channel(LogChannelKey).Message(string(data)).Build()
	return vq.client.Do(ctx, cmd).Error()
}

// Close closes the Valkey client connection.
func (vq *ValkeyQueue) Close() {
	if vq.client != nil {
		vq.client.Close()
	}
}
