package api_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/daflamingfox/conduit/internal/api"
	"github.com/daflamingfox/conduit/internal/db"
	"github.com/daflamingfox/conduit/internal/nodes"
	"github.com/daflamingfox/conduit/internal/storage"
)

func TestAPINodesAndHealthEndpoints(t *testing.T) {
	tempDB := filepath.Join(t.TempDir(), "test_api.db")
	database, err := db.InitDB(tempDB)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer database.Close()

	registry := nodes.NewRegistry()
	_, _ = registry.LoadFromDir("../../nodes")

	storageMgr := storage.NewManager()
	_, _ = storageMgr.RegisterLocation("storage:default", t.TempDir())

	srv := api.NewServer(database, registry, storageMgr)

	// Test GET /api/v1/health
	reqHealth := httptest.NewRequest("GET", "/api/v1/health", nil)
	recHealth := httptest.NewRecorder()
	srv.ServeHTTP(recHealth, reqHealth)

	if recHealth.Code != http.StatusOK {
		t.Errorf("expected status 200 for health endpoint, got %d", recHealth.Code)
	}

	// Test GET /api/v1/nodes
	reqNodes := httptest.NewRequest("GET", "/api/v1/nodes", nil)
	recNodes := httptest.NewRecorder()
	srv.ServeHTTP(recNodes, reqNodes)

	if recNodes.Code != http.StatusOK {
		t.Errorf("expected status 200 for nodes endpoint, got %d", recNodes.Code)
	}
}
