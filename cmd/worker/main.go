package main

import (
	"context"
	"log/slog"
	"os"

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
		} else {
			defer taskQueue.Close()
			slog.Info("Worker registered with Valkey queue", "worker_id", caps.WorkerID)
		}
	}

	_ = context.Background()
	slog.Info("Conduit Worker daemon ready and polling for tasks", "worker_id", caps.WorkerID)
}
