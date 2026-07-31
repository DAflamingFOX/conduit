package worker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"time"

	"github.com/daflamingfox/conduit/internal/engine"
	"github.com/daflamingfox/conduit/internal/queue"
)

// Runner executes binary steps in a subprocess and streams output logs.
type Runner struct {
	queue queue.TaskQueue
}

// NewRunner initializes a step Runner.
func NewRunner(q queue.TaskQueue) *Runner {
	return &Runner{queue: q}
}

// ExecuteStep runs a prepared ExecutionStep binary command and streams line-by-line output logs.
func (r *Runner) ExecuteStep(ctx context.Context, step *engine.ExecutionStep) error {
	binary := step.Manifest.Execution.Binary
	args := step.RenderedArgs

	slog.Info("executing step subprocess", "node_id", step.NodeInstanceID, "binary", binary, "args", args)

	cmd := exec.CommandContext(ctx, binary, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to open stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start subprocess '%s': %w", binary, err)
	}

	// Stream stdout & stderr concurrently
	go r.streamPipe(ctx, step.Context.JobID, step.NodeInstanceID, "info", stdout)
	go r.streamPipe(ctx, step.Context.JobID, step.NodeInstanceID, "error", stderr)

	if err := cmd.Wait(); err != nil {
		step.Context.SetOutputBranch("error")
		r.publishLog(ctx, step.Context.JobID, step.NodeInstanceID, "error", fmt.Sprintf("step failed: %v", err))
		return fmt.Errorf("step execution failed: %w", err)
	}

	step.Context.SetOutputBranch("success")
	r.publishLog(ctx, step.Context.JobID, step.NodeInstanceID, "info", "step execution completed successfully")
	return nil
}

func (r *Runner) streamPipe(ctx context.Context, jobID, nodeID, level string, pipe io.Reader) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		r.publishLog(ctx, jobID, nodeID, level, line)
	}
}

func (r *Runner) publishLog(ctx context.Context, jobID, nodeID, level, message string) {
	entry := queue.JobLogEntry{
		JobID:     jobID,
		NodeID:    nodeID,
		Level:     level,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}

	if level == "error" {
		slog.Error("step log", "job_id", jobID, "node_id", nodeID, "msg", message)
	} else {
		slog.Info("step log", "job_id", jobID, "node_id", nodeID, "msg", message)
	}

	if r.queue != nil {
		_ = r.queue.PublishLog(ctx, entry)
	}
}
