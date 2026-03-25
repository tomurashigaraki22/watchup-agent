# Phase 9 - Installation System

## Status: ✓ Complete

## Overview

Provides one-line installation script for Linux servers with systemd support.

## File: install/install.sh

## Installation Command

```bash
curl -s https://watchup.site/install.sh | sudo bash
```

## Installation Steps

### 1. Root Check
```bash
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root (use sudo)"
    exit 1
fi
```

### 2. Download Binary
```bash
DOWNLOAD_URL="https://watchup.site/releases/latest/watchup-agent"

if command -v curl &> /dev/null; then
    curl -L -o "/tmp/watchup-agent" "${DOWNLOAD_URL}"
elif command -v wget &> /dev/null; then
    wget -O "/tmp/watchup-agent" "${DOWNLOAD_URL}"
fi
```

### 3. Install Binary
```bash
mv "/tmp/watchup-agent" "/usr/local/bin/watchup-agent"
chmod +x "/usr/local/bin/watchup-agent"
```

**Location**: `/usr/local/bin/watchup-agent`
**Permissions**: `755` (executable)

### 4. Create Configuration Directory
```bash
mkdir -p "/etc/watchup"
```

### 5. Create Default Configuration
```bash
cat > "/etc/watchup/config.yaml" << EOF
server_key: "REPLACE_WITH_YOUR_SERVER_KEY"
project_id: ""
server_identifier: ""
sampling_interval: 5
api_endpoint: "https://api.watchup.site"
registered: false

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
EOF
```

**Location**: `/etc/watchup/config.yaml`
**Permissions**: `600` (owner read/write only)

### 6. Create Systemd Service
```bash
cat > "/etc/systemd/system/watchup-agent.service" << EOF
[Unit]
Description=Watchup Server Monitoring Agent
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/watchup-agent
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
```

**Location**: `/etc/systemd/system/watchup-agent.service`

### 7. Enable Service
```bash
systemctl daemon-reload
systemctl enable watchup-agent
```

## Systemd Service Configuration

```ini
[Unit]
Description=Watchup Server Monitoring Agent
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/watchup-agent
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

### Service Features

- **Type**: `simple` (foreground process)
- **Restart**: `always` (auto-restart on failure)
- **RestartSec**: `10` (wait 10s before restart)
- **Logging**: systemd journal
- **Auto-start**: Enabled on boot

## Post-Installation Steps

### 1. Start Agent (First Time)
```bash
sudo systemctl start watchup-agent
```

Agent will prompt for registration:
```
Enter your Project ID: proj_12345
Enter your Master API Key: ma_abcdef123
Enter Server Identifier: prod-api-1
```

### 2. Check Status
```bash
sudo systemctl status watchup-agent
```

Expected output:
```
● watchup-agent.service - Watchup Server Monitoring Agent
   Loaded: loaded (/etc/systemd/system/watchup-agent.service; enabled)
   Active: active (running) since Wed 2026-03-25 14:00:00 UTC
   Main PID: 12345
   CGroup: /system.slice/watchup-agent.service
           └─12345 /usr/local/bin/watchup-agent
```

### 3. View Logs
```bash
sudo journalctl -u watchup-agent -f
```

Output:
```
Mar 25 14:00:00 server watchup-agent[12345]: Watchup Server Agent started
Mar 25 14:00:00 server watchup-agent[12345]: ✓ Agent registered successfully
Mar 25 14:00:00 server watchup-agent[12345]: Starting monitoring loop...
Mar 25 14:01:00 server watchup-agent[12345]: [14:01:00] CPU: 45.2%, RAM: 67.8%
```

## Service Management

### Start
```bash
sudo systemctl start watchup-agent
```

### Stop
```bash
sudo systemctl stop watchup-agent
```

### Restart
```bash
sudo systemctl restart watchup-agent
```

### Enable (auto-start on boot)
```bash
sudo systemctl enable watchup-agent
```

### Disable
```bash
sudo systemctl disable watchup-agent
```

### Status
```bash
sudo systemctl status watchup-agent
```

## Uninstallation

```bash
# Stop and disable service
sudo systemctl stop watchup-agent
sudo systemctl disable watchup-agent

# Remove service file
sudo rm /etc/systemd/system/watchup-agent.service
sudo systemctl daemon-reload

# Remove binary
sudo rm /usr/local/bin/watchup-agent

# Remove configuration (optional)
sudo rm -rf /etc/watchup
```

## File Locations

| File | Location | Purpose |
|------|----------|---------|
| Binary | `/usr/local/bin/watchup-agent` | Executable |
| Config | `/etc/watchup/config.yaml` | Configuration |
| Service | `/etc/systemd/system/watchup-agent.service` | Systemd unit |
| Logs | `journalctl -u watchup-agent` | System logs |

## Permissions

| File | Owner | Permissions | Reason |
|------|-------|-------------|--------|
| Binary | root | 755 | Executable by all |
| Config | root | 600 | Protect server_key |
| Service | root | 644 | Standard systemd |

## Troubleshooting

### Service Won't Start
```bash
# Check service status
sudo systemctl status watchup-agent

# View detailed logs
sudo journalctl -u watchup-agent -n 50

# Check binary exists
ls -l /usr/local/bin/watchup-agent

# Check config exists
ls -l /etc/watchup/config.yaml
```

### Registration Issues
```bash
# Stop service
sudo systemctl stop watchup-agent

# Run manually to see prompts
sudo /usr/local/bin/watchup-agent /etc/watchup/config.yaml

# After registration, restart service
sudo systemctl start watchup-agent
```

### Network Issues
```bash
# Test API connectivity
curl -I https://api.watchup.site

# Check firewall
sudo iptables -L OUTPUT

# View network errors in logs
sudo journalctl -u watchup-agent | grep "Failed to send"
```

## Platform Support

### Supported
- Ubuntu 18.04+
- Debian 9+
- CentOS 7+
- RHEL 7+
- Fedora 28+
- Any Linux with systemd

### Requirements
- systemd (for service management)
- curl or wget (for installation)
- HTTPS connectivity to api.watchup.site

## Build for Distribution

### Build Binary
```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o watchup-agent cmd/agent/main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o watchup-agent-arm64 cmd/agent/main.go
```

### Create Release
```bash
# Upload to release server
aws s3 cp watchup-agent s3://watchup-releases/latest/watchup-agent
aws s3 cp install.sh s3://watchup-releases/install.sh
```

## Installation Output

```
=== Watchup Server Agent Installation ===

Downloading Watchup Agent...
Installing binary to /usr/local/bin...
Creating configuration directory...
Creating default configuration...

⚠️  IMPORTANT: Edit /etc/watchup/config.yaml and set your server_key

Creating systemd service...
Reloading systemd...
Enabling service...

=== Installation Complete ===

Next steps:
1. Edit your server key: sudo nano /etc/watchup/config.yaml
2. Start the agent: sudo systemctl start watchup-agent
3. Check status: sudo systemctl status watchup-agent
4. View logs: sudo journalctl -u watchup-agent -f
```

## Complete Installation Flow

```
User runs: curl -s https://watchup.site/install.sh | sudo bash
    │
    ├─ Download binary
    ├─ Install to /usr/local/bin
    ├─ Create /etc/watchup directory
    ├─ Create default config
    ├─ Create systemd service
    ├─ Enable service
    │
    ▼
User runs: sudo systemctl start watchup-agent
    │
    ▼
Agent prompts for registration
    │
    ├─ Project ID
    ├─ Master API Key
    └─ Server Identifier
    │
    ▼
Agent registers with backend
    │
    ▼
Agent starts monitoring
    │
    ▼
Metrics sent to Watchup dashboard
```

## Next Steps

After installation, the agent is ready to monitor the server and send data to the Watchup platform.
