package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/daflamingfox/conduit/internal/db"
	"github.com/daflamingfox/conduit/internal/engine"
	"github.com/daflamingfox/conduit/internal/nodes"
	"github.com/daflamingfox/conduit/internal/storage"
)

// Server represents the Conduit Manager HTTP REST API and UI static server.
type Server struct {
	db       *db.Database
	registry *nodes.Registry
	storage  *storage.Manager
	mux      *http.ServeMux
}

// NewServer initializes HTTP REST API and static UI routes.
func NewServer(database *db.Database, reg *nodes.Registry, storageMgr *storage.Manager) *Server {
	s := &Server{
		db:       database,
		registry: reg,
		storage:  storageMgr,
		mux:      http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) routes() {
	// API v1 Endpoints
	s.mux.HandleFunc("GET /api/v1/nodes", s.handleListNodes)
	s.mux.HandleFunc("GET /api/v1/storage", s.handleListStorage)
	s.mux.HandleFunc("GET /api/v1/health", s.handleHealthCheck)
	s.mux.HandleFunc("POST /api/v1/jobs/trigger", s.handleTriggerJob)

	// Static UI File Server (from ui/dist if present)
	if _, err := os.Stat("ui/dist"); err == nil {
		fileServer := http.FileServer(http.Dir("ui/dist"))
		s.mux.Handle("/", fileServer)
		slog.Info("serving Svelte 5 dashboard UI from ui/dist")
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS header for UI development mode
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)})
}

func (s *Server) handleListNodes(w http.ResponseWriter, r *http.Request) {
	nodeList := s.registry.List()
	writeJSON(w, http.StatusOK, nodeList)
}

func (s *Server) handleListStorage(w http.ResponseWriter, r *http.Request) {
	probes := s.storage.ProbeAll()
	writeJSON(w, http.StatusOK, probes)
}

type TriggerJobRequest struct {
	FlowID   string `json:"flow_id"`
	FilePath string `json:"file_path"`
}

func (s *Server) handleTriggerJob(w http.ResponseWriter, r *http.Request) {
	var req TriggerJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	if req.FilePath == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "file_path is required"})
		return
	}

	fileInfo := engine.NewFileInfo(req.FilePath, 0)
	ctx, err := engine.NewContext(req.FlowID, fileInfo, "/tmp/conduit/scratch")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	slog.Info("triggered job execution via API", "job_id", ctx.JobID, "file_path", req.FilePath)
	writeJSON(w, http.StatusAccepted, map[string]string{
		"job_id":  ctx.JobID,
		"status":  "queued",
		"message": fmt.Sprintf("Job %s triggered for %s", ctx.JobID, req.FilePath),
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
