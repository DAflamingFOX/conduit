# Conduit Repository & File Tree Layout

This document outlines the complete directory layout and module organization for **Conduit**.

```text
conduit/
├── .github/                    # GitHub actions CI/CD workflows
│   └── workflows/
│       ├── build.yml           # Go build & test workflow
│       └── docker.yml          # Container image publish workflow
│
├── cmd/                        # Application Entry Points
│   ├── manager/                # Conduit Manager daemon binary
│   │   └── main.go             # Entry point for Manager
│   └── worker/                 # Conduit Worker daemon binary
│       └── main.go             # Entry point for Worker
│
├── internal/                   # Private Application Code (Enforced by Go compiler)
│   ├── api/                    # REST HTTP handlers & WebSocket gateways
│   │   ├── handlers/           # Flow, Job, Worker, Storage HTTP routes
│   │   └── websocket/          # Real-time log & status streaming hub
│   │
│   ├── db/                     # SQLite Database layer
│   │   ├── migrations/         # SQL migration scripts
│   │   ├── models/             # Database structs (Flows, Jobs, Workers, Storage)
│   │   └── db.go               # SQLite connection & WAL mode config
│   │
│   ├── engine/                 # Flow Engine Core
│   │   ├── dag.go              # Directed Acyclic Graph validator & compiler
│   │   ├── context.go          # Runtime context envelope (JSON data state)
│   │   └── evaluator.go        # Flow branch condition & logic evaluator
│   │
│   ├── nodes/                  # Node Schema & Execution Registry
│   │   ├── schema.go           # YAML/JSON Node Manifest parser
│   │   ├── registry.go         # Git-backed & local node registry loader
│   │   └── executor/           # Execution handlers (Command, Script, Docker)
│   │       ├── command.go      # Binary command runner (e.g. FFmpeg)
│   │       ├── script.go       # Python / JS / Bash script runner
│   │       └── docker.go       # Isolated container task runner
│   │
│   ├── protocol/               # Manager-Worker gRPC Protocol
│   │   ├── proto/              # Protobuf definitions (.proto)
│   │   │   ├── heartbeat.proto
│   │   │   └── task.proto
│   │   └── pb/                 # Generated Go Protobuf code
│   │
│   ├── queue/                  # Task Queue Layer
│   │   ├── queue.go            # Queue interface abstraction
│   │   └── valkey.go           # Valkey (Redis) client implementation
│   │
│   └── storage/                # Storage Pool & Health Manager
│       ├── manager.go          # Storage pool path resolver & mapping
│       └── probe.go            # Pre-flight mount read/write health checker
│
├── ui/                         # Svelte 5 Web Dashboard
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/     # Flow canvas, Node palette, Log viewers, Worker status
│   │   │   ├── stores/         # Svelte 5 runes ($state) for canvas & job states
│   │   │   └── api.ts          # REST & WebSocket client
│   │   ├── routes/             # SvelteKit page routes (Flows, Jobs, Workers, Settings)
│   │   ├── app.html            # Main HTML wrapper
│   │   └── app.css             # Tailwind / Modern CSS design system
│   ├── package.json            # Node dependencies
│   ├── svelte.config.js        # Svelte 5 configuration
│   └── vite.config.ts          # Vite build config
│
├── nodes/                      # Built-in Standard Node Definitions (YAML)
│   ├── file/                   # Move, Copy, Delete, Hash, Rename
│   ├── media/                  # FFmpeg Transcode, FFprobe, Audio Extract
│   ├── logic/                  # If/Else Condition, Switch, Passthrough
│   └── notification/           # Discord, Webhook, Email
│
├── docker/                     # Container Configurations
│   ├── Dockerfile.manager      # Multi-stage Go build for Manager
│   ├── Dockerfile.worker       # Multi-stage Go build for Worker (with FFmpeg/tools)
│   └── docker-compose.yml      # Complete stack (Manager + Valkey + Worker)
│
├── docs/                       # Project Documentation
│   ├── architecture.md         # Architecture specification
│   └── project_structure.md    # Repository layout reference
│
├── go.mod                      # Go module definition
└── go.sum                      # Go dependency checksums
```
