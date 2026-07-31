package nodes_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/daflamingfox/conduit/internal/nodes"
)

func TestBuiltinNodeFileCopy(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "node_test_copy_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, "source_document.txt")
	dstDir := filepath.Join(tempDir, "backup")
	_ = os.MkdirAll(dstDir, 0755)
	_ = os.WriteFile(srcPath, []byte("conduit copy test content"), 0644)

	registry := nodes.NewRegistry()
	nodesDir, _ := filepath.Abs("../../nodes")
	_, err = registry.LoadFromDir(nodesDir)
	if err != nil {
		t.Fatalf("failed to load nodes dir: %v", err)
	}

	manifest, ok := registry.Get("conduit.file.copy")
	if !ok {
		t.Fatalf("node manifest 'conduit.file.copy' not found in registry")
	}

	input := nodes.TestInput{
		Manifest:      manifest,
		InputFilePath: srcPath,
		Parameters:    map[string]interface{}{"destination": dstDir},
		WorkingDir:    tempDir,
	}

	// 1. Test Dry Run
	dryRes, err := nodes.DryRunNode(input)
	if err != nil {
		t.Fatalf("dry-run failed for conduit.file.copy: %v", err)
	}
	if dryRes.RenderedBinary != "cp" {
		t.Errorf("expected binary 'cp', got '%s'", dryRes.RenderedBinary)
	}

	// 2. Test Execution
	execRes, err := nodes.ExecuteTestNode(context.Background(), input)
	if err != nil {
		t.Fatalf("execution failed for conduit.file.copy: %v", err)
	}
	if execRes.OutputBranch != "success" {
		t.Errorf("expected output branch 'success', got '%s'", execRes.OutputBranch)
	}

	// Verify copied file exists
	expectedDstFile := filepath.Join(dstDir, "source_document.txt")
	if _, err := os.Stat(expectedDstFile); os.IsNotExist(err) {
		t.Errorf("expected copied file to exist at '%s'", expectedDstFile)
	}
}

func TestBuiltinNodeFileMove(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "node_test_move_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, "file_to_move.txt")
	dstDir := filepath.Join(tempDir, "archive")
	_ = os.MkdirAll(dstDir, 0755)
	_ = os.WriteFile(srcPath, []byte("conduit move test content"), 0644)

	registry := nodes.NewRegistry()
	nodesDir, _ := filepath.Abs("../../nodes")
	_, _ = registry.LoadFromDir(nodesDir)

	manifest, _ := registry.Get("conduit.file.move")

	input := nodes.TestInput{
		Manifest:      manifest,
		InputFilePath: srcPath,
		Parameters:    map[string]interface{}{"destination": dstDir},
		WorkingDir:    tempDir,
	}

	execRes, err := nodes.ExecuteTestNode(context.Background(), input)
	if err != nil {
		t.Fatalf("execution failed for conduit.file.move: %v", err)
	}
	if execRes.OutputBranch != "success" {
		t.Errorf("expected output branch 'success', got '%s'", execRes.OutputBranch)
	}

	expectedDstFile := filepath.Join(dstDir, "file_to_move.txt")
	if _, err := os.Stat(expectedDstFile); os.IsNotExist(err) {
		t.Errorf("expected moved file to exist at '%s'", expectedDstFile)
	}
}

func TestBuiltinNodeFileRename(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "node_test_rename_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, "old_name.txt")
	_ = os.WriteFile(srcPath, []byte("conduit rename test content"), 0644)

	registry := nodes.NewRegistry()
	nodesDir, _ := filepath.Abs("../../nodes")
	_, _ = registry.LoadFromDir(nodesDir)

	manifest, _ := registry.Get("conduit.file.rename")

	input := nodes.TestInput{
		Manifest:      manifest,
		InputFilePath: srcPath,
		Parameters:    map[string]interface{}{"new_name": "new_name.txt"},
		WorkingDir:    tempDir,
	}

	execRes, err := nodes.ExecuteTestNode(context.Background(), input)
	if err != nil {
		t.Fatalf("execution failed for conduit.file.rename: %v", err)
	}
	if execRes.OutputBranch != "success" {
		t.Errorf("expected output branch 'success', got '%s'", execRes.OutputBranch)
	}

	expectedDstFile := filepath.Join(tempDir, "new_name.txt")
	if _, err := os.Stat(expectedDstFile); os.IsNotExist(err) {
		t.Errorf("expected renamed file to exist at '%s'", expectedDstFile)
	}
}

func TestBuiltinNodeFileDelete(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "node_test_delete_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, "file_to_delete.txt")
	_ = os.WriteFile(srcPath, []byte("conduit delete test content"), 0644)

	registry := nodes.NewRegistry()
	nodesDir, _ := filepath.Abs("../../nodes")
	_, _ = registry.LoadFromDir(nodesDir)

	manifest, _ := registry.Get("conduit.file.delete")

	input := nodes.TestInput{
		Manifest:      manifest,
		InputFilePath: srcPath,
		WorkingDir:    tempDir,
	}

	execRes, err := nodes.ExecuteTestNode(context.Background(), input)
	if err != nil {
		t.Fatalf("execution failed for conduit.file.delete: %v", err)
	}
	if execRes.OutputBranch != "success" {
		t.Errorf("expected output branch 'success', got '%s'", execRes.OutputBranch)
	}

	if _, err := os.Stat(srcPath); !os.IsNotExist(err) {
		t.Errorf("expected target file to be deleted")
	}
}
