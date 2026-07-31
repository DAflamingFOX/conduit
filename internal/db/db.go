package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log/slog"

	_ "modernc.org/sqlite"
)

//go:embed migrations/001_initial_schema.sql
var initialSchemaSQL string

// Database wraps standard sql.DB connection for Conduit
type Database struct {
	*sql.DB
}

// InitDB initializes SQLite database at dbPath and enables WAL mode
func InitDB(dbPath string) (*Database, error) {
	connStr := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)", dbPath)
	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database at %s: %w", dbPath, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("connected to SQLite database", "path", dbPath, "journal_mode", "WAL")

	// Run initial schema migration
	if _, err := db.Exec(initialSchemaSQL); err != nil {
		return nil, fmt.Errorf("failed to execute initial schema migration: %w", err)
	}

	slog.Info("initial schema migrations applied successfully")
	return &Database{DB: db}, nil
}
