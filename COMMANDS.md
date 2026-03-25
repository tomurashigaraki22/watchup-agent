# Quick Command Reference

## GitHub Deployment

```bash
# Initialize repository
git init
git add .
git commit -m "Initial commit"

# Add remote (replace YOUR_USERNAME)
git remote add origin https://github.com/YOUR_USERNAME/watchup-agent.git

# Push to GitHub
git branch -M main
git push -u origin main

# Create release tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## VPS Installation

```bash
# Connect to VPS
ssh user@your-vps-ip

# Install dependencies (Ubuntu/Debian)
sudo apt update && sudo apt install -y git curl wget

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Clone and build
git clone https://github.com/YOUR_USERNAME/watchup-agent.git
cd watchup-agent
go mod tidy
go build -o watchup-agent cmd/agent/main.go

# Install
sudo mv watchup-agent /usr/local/bin/
sudo chmod +x /usr/local/bin/watchup-agent
sudo mkdir -p /etc/watchup
sudo cp config.yaml /etc/watchup/config.yaml
```

## Systemd Service

```bash
# Create service file
sudo tee /etc/systemd/system/watchup-agent.service > /dev/null <<'EOF'
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

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable watchup-agent
sudo systemctl start watchup-agent
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

# Enable auto-start
sudo systemctl enable watchup-agent

# Disable auto-start
sudo systemctl disable watchup-agent
```

## Logs

```bash
# Real-time logs
sudo journalctl -u watchup-agent -f

# Last 50 lines
sudo journalctl -u watchup-agent -n 50

# Last 100 lines
sudo journalctl -u watchup-agent -n 100

# Today's logs
sudo journalctl -u watchup-agent --since today

# Logs with timestamps
sudo journalctl -u watchup-agent -o short-precise
```

## Troubleshooting

```bash
# Check if binary exists
ls -l /usr/local/bin/watchup-agent

# Check if config exists
ls -l /etc/watchup/config.yaml

# Run manually
sudo /usr/local/bin/watchup-agent /etc/watchup/config.yaml

# Check process
ps aux | grep watchup-agent

# Check resource usage
top -p $(pgrep watchup-agent)

# Test API connectivity
curl -I https://api.watchup.site
```

## Configuration

```bash
# Edit config
sudo nano /etc/watchup/config.yaml

# View config
sudo cat /etc/watchup/config.yaml

# Validate config
sudo /usr/local/bin/watchup-agent /etc/watchup/config.yaml --validate
```

## Updates

```bash
# Pull latest code
cd watchup-agent
git pull origin main

# Rebuild
go build -o watchup-agent cmd/agent/main.go

# Stop service
sudo systemctl stop watchup-agent

# Update binary
sudo mv watchup-agent /usr/local/bin/

# Start service
sudo systemctl start watchup-agent

# Verify
sudo systemctl status watchup-agent
```

## Uninstall

```bash
# Stop and disable
sudo systemctl stop watchup-agent
sudo systemctl disable watchup-agent

# Remove service
sudo rm /etc/systemd/system/watchup-agent.service
sudo systemctl daemon-reload

# Remove binary
sudo rm /usr/local/bin/watchup-agent

# Remove config (optional)
sudo rm -rf /etc/watchup
```

## Build Commands

```bash
# Build for current platform
go build -o watchup-agent cmd/agent/main.go

# Build for Linux AMD64
GOOS=linux GOARCH=amd64 go build -o watchup-agent-linux-amd64 cmd/agent/main.go

# Build for Linux ARM64
GOOS=linux GOARCH=arm64 go build -o watchup-agent-linux-arm64 cmd/agent/main.go

# Build with optimizations
go build -ldflags="-s -w" -o watchup-agent cmd/agent/main.go
```

## Git Commands

```bash
# Check status
git status

# Add all changes
git add .

# Commit
git commit -m "Your message"

# Push
git push origin main

# Create tag
git tag -a v1.0.0 -m "Version 1.0.0"
git push origin v1.0.0

# View tags
git tag -l

# Delete tag
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
```
