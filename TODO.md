# Conduit — Feature Roadmap & TODO

This document tracks features, enhancements, and integrations deferred to future development phases to keep the initial MVP scope clean and focused.

---

## 🎯 Deferred / Future Phase Items

### Storage Protocols & OS Support
- [ ] **SMB / CIFS Protocol Integration**: Dedicated SMB mount health probes for Windows file shares alongside NFS.
- [ ] **Native Windows Worker Daemon**: Standalone non-containerized Windows agent executable for Windows-specific software workflows.

### Triggers & Automation
- [ ] **Automated File Watchers**: Inotify / fsnotify library folder watchers for automatic flow execution on file creation/modification.
- [ ] **Incoming Webhook Triggers**: HTTP endpoints to trigger flow executions from external services (Sonarr, Radarr, Plex, custom scripts).
- [ ] **Cron / Scheduled Triggers**: Time-based flow execution scheduling.

### CLI & External Tooling
- [ ] **Conduit CLI**: Standalone command-line interface for managing jobs, inspecting worker nodes, and triggering flows manually.

### Advanced Storage & Node Features
- [ ] **Cloud Storage Drivers**: Direct S3 / MinIO stream adapters for workers without mounted local shares.
- [ ] **External Node Registries**: Dynamic fetching and hot-reloading of custom node packs from remote Git repositories.
