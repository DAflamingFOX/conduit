package storage

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// StorageLocation represents a managed storage pool.
type StorageLocation struct {
	ID          string       `json:"id"`    // UUID v7
	Alias       string       `json:"alias"` // e.g. "storage:media"
	Path        string       `json:"path"`  // Host mount path
	Status      HealthStatus `json:"status"`
	LastChecked time.Time    `json:"last_checked"`
}

// Manager manages storage pools and resolves logical alias paths.
type Manager struct {
	mu        sync.RWMutex
	locations map[string]*StorageLocation // alias -> location
}

// NewManager initializes a Storage Manager.
func NewManager() *Manager {
	return &Manager{
		locations: make(map[string]*StorageLocation),
	}
}

// RegisterLocation registers or updates a storage pool alias mapping.
func (m *Manager) RegisterLocation(alias string, hostPath string) (*StorageLocation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	locID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID v7 for storage location: %w", err)
	}

	cleanPath := filepath.Clean(hostPath)
	loc := &StorageLocation{
		ID:          locID.String(),
		Alias:       alias,
		Path:        cleanPath,
		Status:      StatusHealthy,
		LastChecked: time.Now().UTC(),
	}

	m.locations[alias] = loc
	slog.Info("registered storage pool location", "alias", alias, "path", cleanPath)
	return loc, nil
}

// GetLocation retrieves a storage pool location by alias.
func (m *Manager) GetLocation(alias string) (*StorageLocation, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	loc, ok := m.locations[alias]
	return loc, ok
}

// ResolvePath converts a logical alias path (e.g. "storage:media/output.mp4") to an absolute host path.
func (m *Manager) ResolvePath(aliasPath string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !strings.HasPrefix(aliasPath, "storage:") {
		return filepath.Clean(aliasPath), nil // Regular path
	}

	parts := strings.SplitN(aliasPath, "/", 2)
	alias := parts[0]

	loc, ok := m.locations[alias]
	if !ok {
		return "", fmt.Errorf("unknown storage location alias: %s", alias)
	}

	if loc.Status == StatusOffline {
		return "", fmt.Errorf("storage location '%s' is offline", alias)
	}

	if len(parts) == 1 {
		return loc.Path, nil
	}
	return filepath.Join(loc.Path, parts[1]), nil
}

// ProbeAll executes health checks across all registered storage pool locations.
func (m *Manager) ProbeAll() map[string]ProbeResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	results := make(map[string]ProbeResult)
	for alias, loc := range m.locations {
		probeRes := PerformHealthProbe(loc.Path)
		loc.Status = probeRes.Status
		loc.LastChecked = probeRes.TestedAt
		results[alias] = probeRes

		if probeRes.Status != StatusHealthy {
			slog.Warn("storage pool health check warning", "alias", alias, "status", probeRes.Status, "error", probeRes.ErrorMessage)
		}
	}

	return results
}
