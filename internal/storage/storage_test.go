package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/daflamingfox/conduit/internal/storage"
)

func TestStorageHealthProbe(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "conduit_storage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp test directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test Healthy Directory
	res := storage.PerformHealthProbe(tempDir)
	if res.Status != storage.StatusHealthy {
		t.Errorf("expected status 'healthy', got '%s' (error: %s)", res.Status, res.ErrorMessage)
	}

	// Test Non-existent Directory
	badDir := filepath.Join(tempDir, "non_existent_folder")
	resBad := storage.PerformHealthProbe(badDir)
	if resBad.Status != storage.StatusOffline {
		t.Errorf("expected status 'offline' for missing directory, got '%s'", resBad.Status)
	}
}

func TestStorageManagerPathResolution(t *testing.T) {
	mgr := storage.NewManager()
	tempDir, _ := os.MkdirTemp("", "conduit_pool_*")
	defer os.RemoveAll(tempDir)

	_, err := mgr.RegisterLocation("storage:media", tempDir)
	if err != nil {
		t.Fatalf("failed to register storage location: %v", err)
	}

	// Resolve Alias Path
	resolved, err := mgr.ResolvePath("storage:media/videos/input.mp4")
	if err != nil {
		t.Fatalf("failed to resolve storage alias path: %v", err)
	}

	expected := filepath.Join(tempDir, "videos/input.mp4")
	if resolved != expected {
		t.Errorf("expected resolved path '%s', got '%s'", expected, resolved)
	}
}
