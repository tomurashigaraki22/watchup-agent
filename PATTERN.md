# Watchup Server Agent - AI Build Instructions

## Project Goal

Build the **Watchup Server Agent**, a lightweight, cross-platform monitoring agent that:

- Collects server metrics (CPU, RAM, top processes)
- Detects sustained spikes based on configurable thresholds
- Sends alerts and metrics to the Watchup backend
- Enforces **one agent per free project** using project-based registration
- Operates securely with read-only permissions
- Runs as a system service (systemd)

---

## Technology Stack

- Language: **Go**
- Metrics: `github.com/shirou/gopsutil`
- Communication: HTTPS REST API
- Deployment: Static binary
- Scheduler: internal Go scheduler
- Configuration: `/etc/watchup/config.yaml`

---

## Repository Structure


watchup-agent
│
├── cmd/agent/main.go # Entry point
├── collectors/
│ ├── cpu.go
│ ├── memory.go
│ └── process.go
├── detectors/spike_detector.go
├── alerts/alert_manager.go
├── transport/api_client.go
├── config/config.go
├── internal/scheduler.go
└── install/install.sh


---

## Development Phases

### Phase 1 — Agent Foundation

- Implement `cmd/agent/main.go`
- Load configuration from `/etc/watchup/config.yaml`
- Initialize scheduler
- Agent prints: `Watchup Server Agent started`

**Verification:** Run `go run cmd/agent/main.go`

---

### Phase 2 — Project Registration (One-Time Setup)

1. During first run, agent asks user for:

   - `project_id`
   - `master_api_key`

2. Agent sends registration request:

```http
POST https://api.watchup.site/agent/register
{
  "project_id": "proj_12345",
  "master_api_key": "ma_abcdef123",
  "server_identifier": "prod-api-1"
}
Backend checks:
master_api_key validity
If agent already installed for this project (agent_installed flag)
Responses:
Success: backend returns a unique server_key
Failure: backend returns error if project already has an agent
Agent saves server_key locally (/etc/watchup/config.yaml) for future authentication.

Result: Registration happens only once per project.

Phase 3 — System Metrics Collection
Implement collectors using gopsutil:

Files:

collectors/cpu.go: CPU usage, optional per-core tracking
collectors/memory.go: total/used/available memory
collectors/process.go: top 5 processes by CPU usage

Sampling interval: 5 seconds

Phase 4 — Scheduler Loop
Implement internal/scheduler.go
Loop every 5 seconds:
while true:
    collect CPU/RAM/process metrics
    send metrics to spike detector
    sleep 5 seconds
Phase 5 — Spike Detection Engine

File: detectors/spike_detector.go

Detect sustained threshold violations:
if usage > threshold:
    increment violation counter
else:
    reset counter

if counter >= required_samples:
    trigger alert
Example: CPU > 80% for 300 seconds → 60 consecutive samples at 5s interval
Phase 6 — Alert System

File: alerts/alert_manager.go

Generate structured alerts:
{
  "server": "prod-api-1",
  "metric": "cpu",
  "usage": 87,
  "duration": 300,
  "top_processes": [{"pid": 4213, "name": "node", "cpu": 52}],
  "timestamp": "2026-03-25T14:11:00Z"
}
Phase 7 — Watchup API Client

File: transport/api_client.go

Sends metrics and alerts to Watchup backend:
POST /server/metrics
POST /server/alerts
GET /agent/config (optional dynamic updates)
Authentication via server_key
All communication via HTTPS
Phase 8 — Configuration System

File: config/config.go

Stores settings in /etc/watchup/config.yaml:
server_key: "srv_89sd0a"
project_id: "proj_12345"
sampling_interval: 5
alerts:
  cpu:
    threshold: 80
    duration: 300
  ram:
    threshold: 75
    duration: 600
Reload configuration every 60 seconds
No restart required
Phase 9 — Installation System

File: install/install.sh

Installation steps:

Download latest binary
Place binary at /usr/local/bin/watchup-agent
Create config directory /etc/watchup
Create systemd service /etc/systemd/system/watchup-agent.service
Enable and start service automatically

Example systemd service:

[Unit]
Description=Watchup Server Monitoring Agent

[Service]
ExecStart=/usr/local/bin/watchup-agent
Restart=always

[Install]
WantedBy=multi-user.target
Agent Operation Flow
User runs agent → provides project_id + master_api_key
         │
         ▼
Backend checks registration → issues server_key
         │
         ▼
Agent stores server_key locally
         │
         ▼
Scheduler collects metrics every 5s
         │
         ▼
Spike detector evaluates thresholds
         │
         ▼
Alert manager sends JSON alerts via API
         │
         ▼
Watchup backend stores data → dashboard displays alerts
Free vs Paid Users
Free users: Only 1 agent per project. Backend rejects additional registrations.
Paid users: Multiple agents allowed. Each gets a unique server_key.
Security & Performance
Read-only monitoring
Minimal resource usage (<1% CPU, 5–20MB RAM)
HTTPS communication only
Server key-based authentication
Configurable thresholds and sampling
Notes for AI Implementation
Implement the project phase by phase.
Do not skip registration logic.
All data must be sent to backend for dashboard display.
Agent must refuse to run if registration fails (enforces free-user limits).
Final output is a single, cross-platform static binary.