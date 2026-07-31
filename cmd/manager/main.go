package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/daflamingfox/conduit/internal/api"
	"github.com/daflamingfox/conduit/internal/db"
	"github.com/daflamingfox/conduit/internal/logger"
	"github.com/daflamingfox/conduit/internal/nodes"
	"github.com/daflamingfox/conduit/internal/queue"
	"github.com/daflamingfox/conduit/internal/storage"
	"github.com/daflamingfox/conduit/internal/version"
)

func main() {
	logLevel := os.Getenv("CONDUIT_LOG_LEVEL")
	jsonFormat := os.Getenv("CONDUIT_LOG_FORMAT") == "json"

	logger.InitLogger(logLevel, jsonFormat)

	slog.Info("starting Conduit Manager", "version", version.Get())

	// Initialize SQLite Database
	dbPath := os.Getenv("CONDUIT_DB_PATH")
	if dbPath == "" {
		dbPath = "./conduit.db"
	}

	database, err := db.InitDB(dbPath)
	if err != nil {
		logger.Fatal("failed to initialize database", "path", dbPath, "error", err)
	}
	defer database.Close()

	// Initialize Node Registry & Load Built-in Manifests
	nodesDir := os.Getenv("CONDUIT_NODES_DIR")
	if nodesDir == "" {
		nodesDir = "./nodes"
	}

	registry := nodes.NewRegistry()
	absNodesDir, _ := filepath.Abs(nodesDir)
	count, err := registry.LoadFromDir(absNodesDir)
	if err != nil {
		slog.Warn("could not load node manifests directory", "path", absNodesDir, "error", err)
	} else {
		slog.Info("Node Registry initialized", "loaded_nodes", count, "path", absNodesDir)
	}

	// Initialize Storage Location Manager
	storageMgr := storage.NewManager()
	defaultStorage := os.Getenv("CONDUIT_STORAGE_PATH")
	if defaultStorage == "" {
		defaultStorage = "/tmp/conduit/storage"
	}
	_ = os.MkdirAll(defaultStorage, 0755)

	_, _ = storageMgr.RegisterLocation("storage:default", defaultStorage)
	probes := storageMgr.ProbeAll()
	slog.Info("Storage Location Manager initialized", "pools", len(probes))

	// Connect to Valkey Task Queue if configured
	valkeyAddr := os.Getenv("VALKEY_ADDR")
	if valkeyAddr != "" {
		taskQueue, err := queue.NewValkeyQueue(valkeyAddr)
		if err != nil {
			slog.Warn("could not connect to Valkey queue (operating in standalone mode)", "address", valkeyAddr, "error", err)
		} else {
			defer taskQueue.Close()
			slog.Info("Conduit Manager connected to Valkey queue broker", "address", valkeyAddr)
		}
	}

	// Initialize & Start HTTP API Server
	apiServer := api.NewServer(database, registry, storageMgr)
	port := os.Getenv("PORT")
	if port == "" {
		port = "5252"
	}

	slog.Info("Conduit Manager initialized successfully", "db_path", dbPath, "http_port", port)

	if err := http.ListenAndServe(":"+port, apiServer); err != nil {
		logger.Fatal("HTTP API server stopped unexpectedly", "error", err)
	}
}
