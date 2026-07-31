package nodes

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

// TestInput specifies input conditions for testing a single node manifest in isolation.
type TestInput struct {
	Manifest      *Manifest              `json:"manifest"`
	InputFilePath string                 `json:"input_file_path"`
	Parameters    map[string]interface{} `json:"parameters"`
	WorkingDir    string                 `json:"working_dir"`
}

// TestResult details the outcome of a node dry-run or subprocess test execution.
type TestResult struct {
	NodeID         string   `json:"node_id"`
	RenderedBinary string   `json:"rendered_binary"`
	RenderedArgs   []string `json:"rendered_args"`
	OutputBranch   string   `json:"output_branch"` // "success" or "error"
	Stdout         string   `json:"stdout"`
	Stderr         string   `json:"stderr"`
	DurationMs     int64    `json:"duration_ms"`
}

// DryRunNode validates manifest parameters and renders argument templates without executing commands.
func DryRunNode(input TestInput) (*TestResult, error) {
	if input.Manifest == nil {
		return nil, fmt.Errorf("test input manifest cannot be nil")
	}

	// Prepare parameter map with defaults
	mergedParams := make(map[string]interface{})
	for _, p := range input.Manifest.Parameters {
		if p.Default != nil {
			mergedParams[p.ID] = p.Default
		}
	}
	for k, v := range input.Parameters {
		mergedParams[k] = v
	}

	// Check required input pins
	for _, pin := range input.Manifest.Inputs {
		if pin.Required && pin.Type == IOTypeFile && input.InputFilePath == "" {
			return nil, fmt.Errorf("required input pin '%s' missing file path", pin.ID)
		}
	}

	fileName := ""
	if input.InputFilePath != "" {
		fileName = filepath.Base(input.InputFilePath)
	}

	tv := TemplateValues{
		FilePath:   input.InputFilePath,
		FileName:   fileName,
		WorkingDir: input.WorkingDir,
		Params:     mergedParams,
	}

	renderedArgs := make([]string, len(input.Manifest.Execution.Args))
	for i, arg := range input.Manifest.Execution.Args {
		renderedArgs[i] = RenderTemplate(arg, tv)
	}

	binary := input.Manifest.Execution.Binary

	return &TestResult{
		NodeID:         input.Manifest.ID,
		RenderedBinary: binary,
		RenderedArgs:   renderedArgs,
		OutputBranch:   "success",
	}, nil
}

// ExecuteTestNode runs DryRunNode and executes the rendered binary in a subprocess, capturing output.
func ExecuteTestNode(ctx context.Context, input TestInput) (*TestResult, error) {
	start := time.Now()

	res, err := DryRunNode(input)
	if err != nil {
		return nil, fmt.Errorf("dry-run failed before execution: %w", err)
	}

	if res.RenderedBinary == "" {
		res.DurationMs = time.Since(start).Milliseconds()
		return res, nil // Internal / passthrough node
	}

	cmd := exec.CommandContext(ctx, res.RenderedBinary, res.RenderedArgs...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	execErr := cmd.Run()
	res.DurationMs = time.Since(start).Milliseconds()
	res.Stdout = stdoutBuf.String()
	res.Stderr = stderrBuf.String()

	if execErr != nil {
		res.OutputBranch = "error"
		return res, fmt.Errorf("subprocess execution returned error: %w (stderr: %s)", execErr, res.Stderr)
	}

	res.OutputBranch = "success"
	return res, nil
}
