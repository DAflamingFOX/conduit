-- Initial Conduit Database Schema

-- Flows Table (DAG Canvas Data)
CREATE TABLE IF NOT EXISTS flows (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    graph_json TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Jobs Table (Flow Execution Runs)
CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    flow_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, running, completed, failed
    current_file_path TEXT NOT NULL,
    error_message TEXT DEFAULT '',
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (flow_id) REFERENCES flows(id) ON DELETE CASCADE
);

-- Job Step Logs Table (Line-by-line Console Outputs)
CREATE TABLE IF NOT EXISTS job_step_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id TEXT NOT NULL,
    node_id TEXT NOT NULL,
    level TEXT NOT NULL DEFAULT 'info', -- info, warn, error
    message TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
);

-- Workers Table (Registered Nodes)
CREATE TABLE IF NOT EXISTS workers (
    id TEXT PRIMARY KEY,
    hostname TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    os TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'offline', -- online, busy, offline
    capabilities_json TEXT DEFAULT '{}',    -- Reported tools/hardware (ffmpeg, GPU, CPU count)
    last_heartbeat DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Storage Pools Table (Storage Locations & Health Status)
CREATE TABLE IF NOT EXISTS storage_locations (
    id TEXT PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,            -- e.g. storage:media
    path TEXT NOT NULL,                    -- Host directory path
    status TEXT NOT NULL DEFAULT 'healthy',-- healthy, degraded, offline
    last_checked DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
