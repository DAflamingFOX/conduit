package nodes_test

import (
	"path/filepath"
	"testing"

	"github.com/daflamingfox/conduit/internal/nodes"
)

func TestRegistryLoadFromDir(t *testing.T) {
	registry := nodes.NewRegistry()
	nodesDir, err := filepath.Abs("../../nodes")
	if err != nil {
		t.Fatalf("failed to resolve nodes directory path: %v", err)
	}

	count, err := registry.LoadFromDir(nodesDir)
	if err != nil {
		t.Fatalf("expected successful loading of node directory, got error: %v", err)
	}

	if count < 4 {
		t.Errorf("expected at least 4 nodes loaded from ./nodes directory, got %d", count)
	}

	// Verify specific loaded node manifest
	moveNode, ok := registry.Get("conduit.file.move")
	if !ok {
		t.Fatalf("expected node 'conduit.file.move' to be registered")
	}

	if moveNode.Name != "Move File" {
		t.Errorf("expected node name 'Move File', got '%s'", moveNode.Name)
	}
}
