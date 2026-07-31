package worker_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/daflamingfox/conduit/internal/engine"
	"github.com/daflamingfox/conduit/internal/nodes"
	"github.com/daflamingfox/conduit/internal/worker"
)

func TestDetectCapabilities(t *testing.T) {
	caps, err := worker.DetectCapabilities()
	if err != nil {
		t.Fatalf("failed to detect worker capabilities: %v", err)
	}

	if caps.WorkerID == "" {
		t.Errorf("worker ID should not be empty")
	}
	if caps.CPUCores <= 0 {
		t.Errorf("CPU cores should be greater than 0")
	}

	t.Logf("Detected Worker: ID=%s, Host=%s, OS=%s, CPU=%d, Tools=%v",
		caps.WorkerID, caps.Hostname, caps.OS, caps.CPUCores, caps.Toolchains)
}

func TestRunnerExecuteStep(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "conduit_runner_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcFile := filepath.Join(tempDir, "source.txt")
	dstFile := filepath.Join(tempDir, "destination.txt")

	_ = os.WriteFile(srcFile, []byte("conduit step execution test"), 0644)

	ctx, _ := engine.NewContext("flow-runner-test", engine.NewFileInfo(srcFile, 27), tempDir)

	step := &engine.ExecutionStep{
		StepID:         "step-copy-1",
		NodeInstanceID: "node-inst-1",
		Manifest: &nodes.Manifest{
			ID:   "conduit.file.copy",
			Name: "Copy File",
			Execution: nodes.ExecutionSpec{
				Type:   nodes.ExecTypeCommand,
				Binary: "cp",
			},
		},
		RenderedArgs: []string{srcFile, dstFile},
		Context:      ctx,
	}

	runner := worker.NewRunner(nil)
	if err := runner.ExecuteStep(context.Background(), step); err != nil {
		t.Fatalf("failed to execute step subprocess: %v", err)
	}

	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Errorf("destination file was not created by step runner")
	}
}
