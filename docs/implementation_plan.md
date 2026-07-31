# Conduit — Implementation Plan & Roadmap

This document outlines the step-by-step phase plan for implementing **Conduit**.

---

## 🚀 Phase 1: Core Scaffolding & Container Infrastructure
1. **Go Module Setup**: Initialize `go.mod` for `github.com/daflamingfox/conduit`.
2. **Skeleton Entrypoints**:
   - `cmd/manager/main.go` (Conduit Manager server shell)
   - `cmd/worker/main.go` (Conduit Worker daemon shell)
3. **Database Layer (`internal/db`)**:
   - SQLite initialization with WAL mode.
   - Initial database tables (`flows`, `jobs`, `workers`, `storage_locations`).
4. **Local Infrastructure (`docker/docker-compose.yml`)**:
   - Multi-container setup containing Manager, Valkey, and Worker node.

---

## ⚡ Phase 2: Flow Engine & Node Schema Specs
1. **Runtime Context Envelope (`internal/engine/context.go`)**:
   - JSON data envelope passing state between connected nodes (file path, metadata, variables, execution status).
2. **DAG Compiler & Validator (`internal/engine/dag.go`)**:
   - Topological sorting and flow execution graph traversal.
   - Branch evaluator for conditional logic (`success`, `error`, boolean expressions).
3. **YAML Node Parser (`internal/nodes/schema.go`)**:
   - Parse YAML node definitions (inputs, outputs, parameters, execution commands).
   - Standard node library seeds (File Move, Copy, Delete, FFmpeg Command Runner).

---

## 🛰️ Phase 3: Manager-Worker IPC & Storage Manager
1. **gRPC Protocol (`internal/protocol`)**:
   - Protobuf schemas for Worker registration, hardware capability reporting (CPU/GPU/toolchains), and heartbeats.
2. **Storage Location Manager (`internal/storage`)**:
   - Logical storage path mapping per worker host.
   - Pre-flight write/read mount health check probes.
3. **Task Queue & Event Bus (`internal/queue`)**:
   - Valkey task dispatcher for distributing node steps to available workers.
   - Pub/Sub channel for live log streaming and worker metrics.

---

## 🎨 Phase 4: Svelte 5 Web Dashboard (Conduit UI)
1. **SvelteKit + Svelte 5 Scaffolding (`ui/`)**:
   - Setup Svelte 5 project with Svelte Flow for interactive node canvas.
2. **Canvas Flow Builder**:
   - Drag-and-drop node palette, connector lines, parameter sidebars.
3. **Real-time Monitoring**:
   - WebSocket client streaming live execution state, node glow effects, and console log streams.

---

## 🧪 Phase 5: Integration & End-to-End Validation
1. **Sample Flow Creation**:
   - Build a sample media intake flow (File Intake → Probe → Condition Check → Transcode → Move to Target Storage).
2. **End-to-End Test Execution**:
   - Validate distributed task execution across Manager, Valkey, and Worker containers.
