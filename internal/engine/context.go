package engine

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// FileInfo represents metadata of the file currently being processed.
type FileInfo struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	SizeBytes  int64  `json:"size_bytes"`
	Extension  string `json:"extension"`
	MimeType   string `json:"mime_type"`
	HashSHA256 string `json:"hash_sha256,omitempty"`
}

// NewFileInfo constructs a FileInfo struct from a target filepath and file size.
func NewFileInfo(filePath string, sizeBytes int64) FileInfo {
	ext := strings.TrimPrefix(filepath.Ext(filePath), ".")
	name := filepath.Base(filePath)
	return FileInfo{
		Path:      filePath,
		Name:      name,
		SizeBytes: sizeBytes,
		Extension: strings.ToLower(ext),
	}
}

// Context is the thread-safe state envelope passed between DAG nodes during a job execution.
type Context struct {
	mu           sync.RWMutex           `json:"-"`
	JobID        string                 `json:"job_id"`
	FlowID       string                 `json:"flow_id"`
	CurrentFile  FileInfo               `json:"current_file"`
	WorkingDir   string                 `json:"working_dir"`
	Variables    map[string]interface{} `json:"variables"`
	OutputBranch string                 `json:"output_branch"`
	CreatedAt    time.Time              `json:"created_at"`
}

// NewContext initializes a new job execution context with UUID v7 IDs.
func NewContext(flowID string, initialFile FileInfo, workingDir string) (*Context, error) {
	jobUUID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID v7 for job: %w", err)
	}

	if flowID == "" {
		flowUUID, err := uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID v7 for flow: %w", err)
		}
		flowID = flowUUID.String()
	}

	return &Context{
		JobID:        jobUUID.String(),
		FlowID:       flowID,
		CurrentFile:  initialFile,
		WorkingDir:   workingDir,
		Variables:    make(map[string]interface{}),
		OutputBranch: "success",
		CreatedAt:    time.Now().UTC(),
	}, nil
}

// GetVariable thread-safely retrieves a custom variable value.
func (c *Context) GetVariable(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.Variables[key]
	return val, ok
}

// SetVariable thread-safely stores a custom variable value.
func (c *Context) SetVariable(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Variables[key] = value
}

// SetOutputBranch thread-safely sets the output branch result (e.g. "success", "error").
func (c *Context) SetOutputBranch(branch string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.OutputBranch = branch
}

// GetOutputBranch thread-safely returns the current output branch.
func (c *Context) GetOutputBranch() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.OutputBranch
}
