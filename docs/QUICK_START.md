# Watchup Server Agent - Quick Start Guide

## Installation (Production)

```bash
curl -s https://watchup.site/install.sh | sudo bash
```

## First Run

```bash
sudo systemctl start watchup-agent
```

You'll be prompted:
```
Enter your Project ID: proj_12345
Enter your Master API Key: ma_abcdef123
Enter Server Identifier: prod-api-1
```

## Verify Installation

```bash
# Check status
sudo systemctl status watchup-agent

# View logs
sudo journalctl -u watchup-agent -f
```

## Development Setup

### 1. Clone and Build

```bash
git clone https://github.com/tomurashigaraki22/watchup-agent.git
cd watchup-agent
go mod tidy
go build -o watchup-agent cmd/agent/main.go
```

### 2. Create Config

```bash
cp config.yaml my-config.yaml
# Edit my-config.yaml with your settings
```

### 3. Run Locally

```bash
go run cmd/agent/main.go my-config.yaml
```

## Configuration

Edit `/etc/watchup/config.yaml`:

```yaml
server_key: "srv_89sd0a"
project_id: "proj_12345"
server_identifier: "prod-api-1"
sampling_interval: 5
api_endpoint: "https://api.watchup.site"
registered: true

alerts:
  cpu:
    threshold: 80    # Percentage
    duration: 300    # Seconds
  ram:
    threshold: 75
    duration: 600
  process_cpu:
    threshold: 60
    duration: 120
```

## Service Management

```bash
# Start
sudo systemctl start watchup-agent

# Stop
sudo systemctl stop watchup-agent

# Restart
sudo systemctl restart watchup-agent

# Status
sudo systemctl status watchup-agent

# Logs
sudo journalctl -u watchup-agent -f

# Enable auto-start
sudo systemctl enable watchup-agent

# Disable auto-start
sudo systemctl disable watchup-agent
```

## Testing Alerts

Lower thresholds for testing:

```yaml
alerts:
  cpu:
    threshold: 10   # Will trigger easily
    duration: 30    # 30 seconds
```

Restart agent:
```bash
sudo systemctl restart watchup-agent
```

Watch for alerts:
```bash
sudo journalctl -u watchup-agent -f
```

## Troubleshooting

### Agent Won't Start

```bash
# Check logs
sudo journalctl -u watchup-agent -n 50

# Run manually
sudo /usr/local/bin/watchup-agent /etc/watchup/config.yaml
```

### Registration Failed

```bash
# Verify credentials
# Check network connectivity
curl -I https://api.watchup.site

# Try manual registration
sudo systemctl stop watchup-agent
sudo /usr/local/bin/watchup-agent /etc/watchup/config.yaml
```

### Network Issues

```bash
# Test API
curl https://api.watchup.site/health

# Check firewall
sudo iptables -L OUTPUT
```

## Uninstall

```bash
sudo systemctl stop watchup-agent
sudo systemctl disable watchup-agent
sudo rm /etc/systemd/system/watchup-agent.service
sudo systemctl daemon-reload
sudo rm /usr/local/bin/watchup-agent
sudo rm -rf /etc/watchup
```

## Build for Different Platforms

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o watchup-agent-linux-amd64 cmd/agent/main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o watchup-agent-linux-arm64 cmd/agent/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o watchup-agent-darwin-amd64 cmd/agent/main.go
```

## API Endpoints

The agent communicates with:

- `POST /agent/register` - Registration
- `POST /server/metrics` - Metrics (every 5s)
- `POST /server/alerts` - Alerts (on threshold)
- `GET /agent/config` - Config updates (every 60s)

## Support

- Documentation: `docs/`
- Issues: GitHub Issues
- Email: support@watchup.site
