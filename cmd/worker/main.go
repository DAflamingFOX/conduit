package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/daflamingfox/conduit/internal/logger"
	"github.com/daflamingfox/conduit/internal/queue"
	"github.com/daflamingfox/conduit/internal/version"
	"github.com/daflamingfox/conduit/internal/worker"
)

func main() {
	logLevel := os.Getenv("CONDUIT_LOG_LEVEL")
	jsonFormat := os.Getenv("CONDUIT_LOG_FORMAT") == "json"

	logger.InitLogger(logLevel, jsonFormat)

	slog.Info("starting Conduit Worker daemon", "version", version.Get())

	// Auto-detect worker system capabilities
	caps, err := worker.DetectCapabilities()
	if err != nil {
		logger.Fatal("failed to detect worker hardware capabilities", "error", err)
	}

	slog.Info("Worker capabilities detected",
		"worker_id", caps.WorkerID,
		"hostname", caps.Hostname,
		"ip", caps.IPAddress,
		"os", caps.OS,
		"cpu_cores", caps.CPUCores,
		"toolchains", caps.Toolchains,
	)

	// Connect to Valkey Task Queue if address provided
	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr != "" {
		taskQueue, err := queue.NewValkeyQueue(valkeyAddr)
		if err != nil {
			slog.Warn("could not connect to Valkey queue (operating in standalone mode)", "address", valkeyAddr, "error", err)
			select {} // Block forever in standalone mode
		}
		defer taskQueue.Close()

		slog.Info("Worker registered with Valkey queue", "worker_id", caps.WorkerID)
		runner := worker.NewRunner(taskQueue)
		slog.Info("Conduit Worker daemon ready and polling for tasks", "worker_id", caps.WorkerID)

		// Active Worker Polling Loop
		for {
			step, err := taskQueue.PopStep(context.Background(), caps.WorkerID, 5*time.Second)
			if err != nil {
				slog.Error("error polling task queue", "error", err)
				time.Sleep(2 * time.Second)
				continue
			}
			if step == nil {
				continue // No task available, continue polling
			}

			slog.Info("claimed task step", "step_id", step.StepID, "node_id", step.NodeInstanceID)
			if err := runner.ExecuteStep(context.Background(), step); err != nil {
				slog.Error("step execution failed", "step_id", step.StepID, "error", err)
			}
		}
	} else {
		slog.Info("Conduit Worker running in standalone mode (no Valkey address specified)", "worker_id", caps.WorkerID)
		select {} // Block forever
	}
}
