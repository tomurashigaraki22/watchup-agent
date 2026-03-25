# Deployment Guide

This guide covers deploying the Watchup Server Agent to GitHub and installing it on your VPS.

---

## Part 1: Deploying to GitHub

### Step 1: Prepare Your Repository

```bash
# Navigate to project directory
cd watchup-agent

# Check current status
git status
```

### Step 2: Initialize Git (if not already done)

```bash
git init
git add .
git commit -m "Initial commit: Watchup Server Agent v1.0.0"
```

### Step 3: Create GitHub Repository

1. Go to https://github.com/new
2. Repository name: `watchup-agent`
3. Description: `Lightweight server monitoring agent for Watchup platform`
4. Choose: **Public** (recommended) or Private
5. **Do NOT** check "Initialize with README" (you already have one)
6. Click **"Create repository"**

### Step 4: Connect Local Repository to GitHub

```bash
# Add GitHub as remote (replace YOUR_USERNAME with your GitHub username)
git remote add origin https://github.com/YOUR_USERNAME/watchup-agent.git

# Verify remote
git remote -v

# Push to GitHub
git branch -M main
git push -u origin main
```

### Step 5: Verify Upload

Visit `https://github.com/YOUR_USERNAME/watchup-agent` to see your repository.

### Step 6: Create First Release

```bash
# Tag the current version
git tag -a v1.0.0 -m "Release v1.0.0: Initial production release"

# Push the tag
git push origin v1.0.0
```

Then on GitHub:
1. Go to your repository
2. Click **"Releases"** (right sidebar)
3. Click **"Create a new release"**
4. Choose tag: `v1.0.0`
5. Release title: `v1.0.0 - Initial Release`
6. Description:
   ```
   ## Features
   - CPU, RAM, and process monitoring
   - Sustained spike detection
   - One-agent-per-free-project enforcement
   - Dynamic configuration updates
   - Systemd service support
   
   ## Installation
   See README.md for installation instructions.
   ```
7. Click **"Publish release"**

---

## Part 2: Installing on VPS

### Prerequisites

- Linux VPS (Ubuntu, Debian, CentOS, etc.)
- SSH access with sudo privileges
- Internet connectivity

### Method 1: Install from GitHub (Recommended)

#### Step 1: Connect to VPS

```bash
ssh user@your-vps-ip
```

#### Step 2: Install Dependencies

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y git curl wget
```

**CentOS/RHEL:**
```bash
sudo yum install -y git curl wget
```

#### Step 3: Install Go

```bash
# Download Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# Extract
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

#### Step 4: Clone and Build

```bash
# Clone your repository
git clone https://github.com/YOUR_USERNAME/watchup-agent.git
cd watchup-agent

# Install dependencies
go mod tidy

# Build the agent
go build -o watchup-agent cmd/agent/main.go

# Verify build
./watchup-agent --help
```

#### Step 5: Install Binary

```bash
# Move to system directory
sudo mv watchup-agent /usr/local/bin/
sudo chmod +x /usr/local/bin/watchup-agent

# Verify installation
which watchup-agent
```

#### Step 6: Create Configuration

```bash
# Create config directory
sudo mkdir -p /etc/watchup

# Copy sample config
sudo cp config.yaml /etc/watchup/config.yaml

# Set proper permissions
sudo chmod 600 /etc/watchup/config.yaml
```

#### Step 7: Create Systemd Service

```bash
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
```

#### Step 8: Enable and Start Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable auto-start
sudo systemctl enable watchup-agent

# Start the service
sudo systemctl start watchup-agent

# Check status
sudo systemctl status watchup-agent
```

#### Step 9: Register Agent

The agent will prompt for registration on first run. View logs:

```bash
sudo journalctl -u watchup-agent -f
```

You'll see:
```
Enter your Project ID: 
Enter your Master API Key: 
Enter Server Identifier: 
```

Enter your credentials from the Watchup dashboard.

#### Step 10: Verify Operation

```bash
# Check service status
sudo systemctl status watchup-agent

# View logs
sudo journalctl -u watchup-agent -n 50

# Check if agent is running
ps aux | grep watchup-agent
```

Expected log output:
```
Watchup Server Agent started
✓ Agent registered successfully
Project ID: proj_12345
Server: my-vps-server
Starting monitoring loop...
[14:23:15] CPU: 12.3%, RAM: 45.6%
```

---

### Method 2: Install from Pre-built Binary

If you've created a GitHub release with binaries:

```bash
# Download binary
wget https://github.com/YOUR_USERNAME/watchup-agent/releases/download/v1.0.0/watchup-agent-linux-amd64

# Make executable
chmod +x watchup-agent-linux-amd64

# Move to system directory
sudo mv watchup-agent-linux-amd64 /usr/local/bin/watchup-agent

# Follow steps 6-10 from Method 1
```

---

### Method 3: One-Line Installation (Future)

Once you host the install script:

```bash
curl -s https://raw.githubusercontent.com/YOUR_USERNAME/watchup-agent/main/install/install.sh | sudo bash
```

---

## Part 3: Multiple VPS Deployment

### Deploy to Multiple Servers

Create a deployment script:

```bash
#!/bin/bash
# deploy-to-servers.sh

SERVERS=(
    "user@server1.example.com"
    "user@server2.example.com"
    "user@server3.example.com"
)

for server in "${SERVERS[@]}"; do
    echo "Deploying to $server..."
    
    ssh "$server" << 'ENDSSH'
        # Install dependencies
        sudo apt update && sudo apt install -y git curl
        
        # Clone and build
        cd /tmp
        git clone https://github.com/YOUR_USERNAME/watchup-agent.git
        cd watchup-agent
        go build -o watchup-agent cmd/agent/main.go
        
        # Install
        sudo mv watchup-agent /usr/local/bin/
        sudo mkdir -p /etc/watchup
        
        # Create service
        sudo tee /etc/systemd/system/watchup-agent.service > /dev/null <<'EOF'
[Unit]
Description=Watchup Server Monitoring Agent
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/watchup-agent
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
        
        # Enable and start
        sudo systemctl daemon-reload
        sudo systemctl enable watchup-agent
        sudo systemctl start watchup-agent
        
        echo "Deployment complete on $(hostname)"
ENDSSH
    
    echo "✓ Deployed to $server"
done

echo "All deployments complete!"
```

Run:
```bash
chmod +x deploy-to-servers.sh
./deploy-to-servers.sh
```

---

## Part 4: Updating the Agent

### Update on Single VPS

```bash
# SSH to VPS
ssh user@your-vps-ip

# Stop service
sudo systemctl stop watchup-agent

# Pull latest code
cd watchup-agent
git pull origin main

# Rebuild
go build -o watchup-agent cmd/agent/main.go

# Reinstall
sudo mv watchup-agent /usr/local/bin/

# Restart service
sudo systemctl start watchup-agent

# Verify
sudo systemctl status watchup-agent
```

### Update via GitHub Release

```bash
# Download new version
wget https://github.com/YOUR_USERNAME/watchup-agent/releases/download/v1.1.0/watchup-agent-linux-amd64

# Stop service
sudo systemctl stop watchup-agent

# Replace binary
sudo mv watchup-agent-linux-amd64 /usr/local/bin/watchup-agent
sudo chmod +x /usr/local/bin/watchup-agent

# Start service
sudo systemctl start watchup-agent
```

---

## Part 5: Monitoring Deployment

### Check Agent Status Across Servers

```bash
#!/bin/bash
# check-agents.sh

SERVERS=(
    "user@server1.example.com"
    "user@server2.example.com"
)

for server in "${SERVERS[@]}"; do
    echo "Checking $server..."
    ssh "$server" "sudo systemctl status watchup-agent --no-pager | head -n 5"
    echo "---"
done
```

### View Logs from All Servers

```bash
#!/bin/bash
# view-logs.sh

SERVER=$1

if [ -z "$SERVER" ]; then
    echo "Usage: ./view-logs.sh user@server.com"
    exit 1
fi

ssh "$SERVER" "sudo journalctl -u watchup-agent -f"
```

---

## Part 6: Troubleshooting Deployment

### Build Fails

```bash
# Check Go version
go version  # Should be 1.20+

# Clean and rebuild
go clean
go mod tidy
go build -v -o watchup-agent cmd/agent/main.go
```

### Service Won't Start

```bash
# Check logs
sudo journalctl -u watchup-agent -n 100

# Run manually
sudo /usr/local/bin/watchup-agent /etc/watchup/config.yaml

# Check permissions
ls -l /usr/local/bin/watchup-agent
ls -l /etc/watchup/config.yaml
```

### Registration Issues

```bash
# Stop service
sudo systemctl stop watchup-agent

# Run interactively
sudo /usr/local/bin/watchup-agent /etc/watchup/config.yaml

# After registration, restart service
sudo systemctl start watchup-agent
```

---

## Part 7: Production Checklist

Before deploying to production:

- [ ] Code is tested locally
- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Version is tagged in Git
- [ ] GitHub release is created
- [ ] Binaries are built for target platforms
- [ ] Installation script is tested
- [ ] Systemd service is configured
- [ ] Firewall allows HTTPS to api.watchup.site
- [ ] Monitoring is configured in Watchup dashboard
- [ ] Alert thresholds are appropriate
- [ ] Backup plan exists for rollback

---

## Support

- **Documentation**: See `docs/` folder
- **Issues**: https://github.com/YOUR_USERNAME/watchup-agent/issues
- **Email**: support@watchup.site
