# System Architecture & Infrastructure Plan — Conduit

Conduit is an open-source, scalable, distributed node-based file processing platform capable of transforming, analyzing, routing, and processing files of any type across distributed processing nodes.

---

## 1. System Overview & Core Architecture

Conduit operates on a **Conduit Manager (Server)** and **Conduit Worker** model driven by a visual **Directed Acyclic Graph (DAG)** execution engine.

```mermaid
flowchart TB
    subgraph Clients["Web Dashboard"]
        UI["Svelte 5 Dashboard (Svelte Flow)"]
    end

    subgraph Manager["Conduit Manager"]
        API["API Gateway & WebSockets"]
        DB[("Embedded SQLite - WAL Mode")]
        SCHED["Job Scheduler & Flow Compiler"]
        MQ["Task Queue Broker (Valkey)"]
        REGISTRY["Node Registry (JSON/YAML Schema)"]
        SM["Storage Location Manager & Health Probes"]
    end

    subgraph NodeGroup["Distributed Worker Nodes (Linux Containers)"]
        subgraph Node1["Worker Server 1"]
            AGENT1["Worker Daemon"]
            EXEC1["Flow Execution Sandbox"]
        end
        subgraph Node2["Worker Server 2"]
            AGENT2["Worker Daemon"]
            EXEC2["Flow Execution Sandbox"]
        end
    end

    subgraph Storage["Managed Storage Pools"]
        NAS["NFS / SMB (CIFS) Shared Storage"]
    end

    UI <-->|WebSocket / REST| API

    API <--> SCHED
    SCHED <--> DB
    SCHED <--> MQ
    SCHED <--> REGISTRY
    SCHED <--> SM

    MQ <-->|gRPC / WebSockets| AGENT1
    MQ <-->|gRPC / WebSockets| AGENT2

    EXEC1 <-->|Pre-flight Verified Read/Write| NAS
    EXEC2 <-->|Pre-flight Verified Read/Write| NAS
```

---

## 2. Component Names & Stack

* **Conduit Manager**: Central server managing flow canvas graphs, SQLite DB, task scheduling, storage health probes, and worker dispatching.
* **Conduit Worker**: Processing node agent that executes flow node steps in container sandboxes and reports hardware capabilities.
* **Conduit UI**: Svelte 5 visual flow canvas and management dashboard.
* **Language & Stack**: Go (Golang) backend (`cmd/`, `internal/`), Svelte 5 frontend (`ui/`), SQLite database, Valkey task queue, gRPC/WebSocket streaming.
