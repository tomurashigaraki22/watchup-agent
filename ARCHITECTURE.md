# Watchup Server Agent - System Architecture

## Overview

The Watchup Server Agent is a lightweight, terminal-based monitoring tool designed to run directly on user servers. It continuously monitors system resources (CPU, RAM, processes) and reports anomalies to the Watchup platform without modifying or terminating any processes.

## Design Principles

- **Lightweight**: Minimal resource footprint (< 1% CPU, 5-20MB RAM)
- **Secure**: Read-only permissions, HTTPS communication, authenticated via server_key
- **Easy to Install**: Single static binary, one-line installation
- **Highly Configurable**: User-customizable thresholds and alert rules
- **Cross-Platform**: Compatible with most Linux servers

## Technology Stack

- **Language**: Go
- **System Metrics Library**: github.com/shirou/gopsutil
- **Deployment**: Systemd service
- **Communication**: HTTPS REST API

## System Architecture

```
watchup-agent
│
├── cmd/
│   └── main.go                 # Entry point, initialization
│
├── collectors/
│   ├── cpu.go                  # CPU metrics collection
│   ├── memory.go               # RAM metrics collection
│   └── process.go              # Process information collection
│
├── detectors/
│   └── spike_detector.go       # Threshold violation detection
│
├── alerts/
│   └── alert_manager.go        # Alert generation and formatting
│
├── transport/
│   └── api_client.go           # Watchup API communication
│
├── config/
│   └── config.go               # Configuration management
│
├── internal/
│   └── scheduler.go            # Monitoring loop scheduler
│
└── install/
    └── install.sh              # Installation script
```

## Core Components

### 1. Main Entry Point (`cmd/main.go`)

Responsibilities:
- Load configuration from `/etc/watchup/config.yaml`
- Initialize collectors, detectors, and alert manager
- Start the monitoring scheduler
- Coordinate the monitoring loop

### 2. Collectors (`collectors/`)

#### CPU Collector (`cpu.go`)
- Collects total CPU usage percentage
- Optional per-core usage tracking
- Sampling interval: 5 seconds

#### Memory Collector (`memory.go`)
- Tracks total, used, and available memory
- Calculates memory usage percentage
- Monitors memory pressure

#### Process Collector (`process.go`)
- Identifies top resource-consuming processes
- Collects: PID, process name, CPU %, memory %
- Correlates processes with resource spikes

### 3. Spike Detector (`detectors/spike_detector.go`)

#### Detection Algorithm
```
if usage > threshold:
    increase violation counter
else:
    reset counter

if violation counter >= required samples:
    trigger alert
```

#### Sustained Threshold Logic
- Prevents false positives from momentary spikes
- Requires threshold violation for specified duration
- Example: CPU > 80% for 300 seconds (60 consecutive samples)

### 4. Alert Manager (`alerts/alert_manager.go`)

Generates structured alerts containing:
- Server identifier
- Affected metric (CPU/RAM/process)
- Usage percentage
- Duration exceeded
- Top responsible processes
- Timestamp

Example alert payload:
```json
{
  "server": "prod-api-1",
  "metric": "cpu",
  "usage": 87,
  "duration": 300,
  "top_processes": [
    {
      "pid": 4213,
      "name": "node",
      "cpu": 52
    }
  ],
  "timestamp": "2026-03-25T14:11:00Z"
}
```

### 5. API Client (`transport/api_client.go`)

#### Endpoints

**Send Metrics**
```
POST https://api.watchup.site/server/metrics
```

**Send Alerts**
```
POST https://api.watchup.site/server/alerts
```

**Fetch Configuration** (dynamic updates)
```
GET https://api.watchup.site/agent/config?server_key=xxx
```

Authentication: Server key in request headers/payload

### 6. Configuration System (`config/config.go`)

Configuration file: `/etc/watchup/config.yaml`

Example:
```yaml
server_key: "srv_89sd0a"
sampling_interval: 5

alerts:
  cpu:
    threshold: 80
    duration: 300
  ram:
    threshold: 75
    duration: 600
  process_cpu:
    threshold: 60
    duration: 120
```

Features:
- Loaded at startup
- Dynamic reload via API polling (every 60 seconds)
- No restart required for configuration changes

### 7. Scheduler (`internal/scheduler.go`)

Monitoring loop (every 5 seconds):
1. Collect CPU metrics
2. Collect RAM metrics
3. Collect process metrics
4. Send to spike detector
5. Transmit metrics to API
6. Check for configuration updates (every 60 seconds)

## Data Flow

```
┌─────────────────┐
│  System (OS)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Collectors    │ ◄── Every 5 seconds
│  (CPU/RAM/Proc) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Spike Detector  │ ◄── Threshold logic
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Alert Manager   │ ◄── Generate alerts
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   API Client    │ ◄── HTTPS to Watchup
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Watchup Backend │
└─────────────────┘
```

## Installation Architecture

### One-Line Installation
```bash
curl -s https://watchup.site/install.sh | bash
```

### Installation Process
1. Download latest binary from release server
2. Place binary at `/usr/local/bin/watchup-agent`
3. Create configuration directory `/etc/watchup`
4. Generate systemd service `/etc/systemd/system/watchup-agent.service`
5. Enable and start service

### Systemd Service
```ini
[Unit]
Description=Watchup Server Monitoring Agent

[Service]
ExecStart=/usr/local/bin/watchup-agent
Restart=always

[Install]
WantedBy=multi-user.target
```

## Security Model

- **Read-Only Operations**: No system modifications or process termination
- **Minimal Permissions**: Runs with standard user privileges
- **Encrypted Communication**: All API calls via HTTPS
- **Authentication**: Server key-based authentication
- **No Sensitive Data**: Only collects resource metrics and process names

## Performance Targets

| Metric | Target |
|--------|--------|
| CPU Usage | < 1% |
| Memory Usage | 5-20 MB |
| Network Usage | Minimal (small payloads only) |
| Sampling Interval | 5 seconds |
| Config Refresh | 60 seconds |

## Scalability Considerations

- Single static binary deployment
- No external dependencies
- Minimal network overhead
- Efficient metric aggregation
- Configurable sampling rates

## Future Enhancements

- Remote configuration management via dashboard
- Dynamic threshold adjustments
- Metric enable/disable controls
- Extended process analytics
- Custom alert rules
- Multi-server orchestration

## Integration with Watchup Platform

The agent extends Watchup from application monitoring to full server observability, enabling users to monitor both applications and infrastructure from a single platform.

### Dashboard Capabilities
- View real-time server metrics
- Configure alert thresholds
- Manage multiple servers
- Analyze historical trends
- Correlate application and infrastructure events
