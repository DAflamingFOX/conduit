package storage

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// HealthStatus represents the status of a storage location or mount.
type HealthStatus string

const (
	StatusHealthy  HealthStatus = "healthy"
	StatusDegraded HealthStatus = "degraded"
	StatusOffline  HealthStatus = "offline"
)

// ProbeResult details the outcome of a storage mount write/read probe.
type ProbeResult struct {
	Status       HealthStatus  `json:"status"`
	Path         string        `json:"path"`
	ErrorMessage string        `json:"error_message,omitempty"`
	Latency      time.Duration `json:"latency"`
	TestedAt     time.Time     `json:"tested_at"`
}

// PerformHealthProbe executes a pre-flight write/read probe on a target mount directory.
func PerformHealthProbe(dirPath string) ProbeResult {
	start := time.Now()

	result := ProbeResult{
		Status:   StatusHealthy,
		Path:     dirPath,
		TestedAt: start.UTC(),
	}

	// 1. Check directory existence and read permission
	info, err := os.Stat(dirPath)
	if err != nil {
		result.Status = StatusOffline
		if os.IsNotExist(err) {
			result.ErrorMessage = fmt.Sprintf("storage directory does not exist: %s", dirPath)
		} else {
			result.ErrorMessage = fmt.Sprintf("storage directory inaccessible: %v", err)
		}
		result.Latency = time.Since(start)
		return result
	}

	if !info.IsDir() {
		result.Status = StatusOffline
		result.ErrorMessage = fmt.Sprintf("storage path is not a directory: %s", dirPath)
		result.Latency = time.Since(start)
		return result
	}

	// 2. Perform write/read lock probe using a temporary probe file
	probeUUID, _ := uuid.NewV7()
	probeFilename := fmt.Sprintf(".conduit_probe_%s.tmp", probeUUID.String()[:8])
	probePath := filepath.Join(dirPath, probeFilename)

	probeContent := []byte("conduit_mount_health_check_" + time.Now().Format(time.RFC3339))
	if err := os.WriteFile(probePath, probeContent, 0644); err != nil {
		result.Status = StatusDegraded
		result.ErrorMessage = fmt.Sprintf("failed to write to storage mount (mount may be read-only or disconnected): %v", err)
		result.Latency = time.Since(start)
		return result
	}

	// Cleanup probe file
	defer func() {
		_ = os.Remove(probePath)
	}()

	// 3. Verify probe read back
	readBytes, err := os.ReadFile(probePath)
	if err != nil || string(readBytes) != string(probeContent) {
		result.Status = StatusDegraded
		result.ErrorMessage = "storage mount read back verification failed"
		result.Latency = time.Since(start)
		return result
	}

	result.Latency = time.Since(start)
	slog.Debug("storage mount pre-flight health probe passed", "path", dirPath, "latency_ms", result.Latency.Milliseconds())
	return result
}
