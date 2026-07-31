package nodes

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Registry maintains an in-memory thread-safe map of available node manifests.
type Registry struct {
	mu    sync.RWMutex
	nodes map[string]*Manifest
}

// NewRegistry initializes an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		nodes: make(map[string]*Manifest),
	}
}

// Register adds or updates a node manifest in the registry.
func (r *Registry) Register(m *Manifest) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[m.ID] = m
	slog.Debug("registered node manifest", "id", m.ID, "name", m.Name, "category", m.Category)
}

// Get retrieves a registered node manifest by ID.
func (r *Registry) Get(id string) (*Manifest, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.nodes[id]
	return m, ok
}

// List returns a slice of all registered node manifests.
func (r *Registry) List() []*Manifest {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Manifest, 0, len(r.nodes))
	for _, m := range r.nodes {
		result = append(result, m)
	}
	return result
}

// LoadFromDir recursively scans target directory for .json files, parses with hujson, and registers them.
func (r *Registry) LoadFromDir(dirPath string) (int, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return 0, fmt.Errorf("node directory does not exist: %s", dirPath)
	}

	count := 0
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			slog.Error("failed to read node file", "path", path, "error", err)
			return nil
		}

		m, err := ParseManifest(data)
		if err != nil {
			slog.Warn("skipping invalid node manifest file", "path", path, "error", err)
			return nil
		}

		r.Register(m)
		count++
		return nil
	})

	if err != nil {
		return count, fmt.Errorf("failed to walk node directory %s: %w", dirPath, err)
	}

	slog.Info("loaded node manifests from directory", "path", dirPath, "total_nodes", count)
	return count, nil
}
