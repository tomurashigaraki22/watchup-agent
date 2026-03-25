# Your Personal Deployment Guide

Follow these exact steps to deploy your Watchup Server Agent.

---

## ✅ Step 1: Push to GitHub (5 minutes)

### 1.1 Open Terminal in Project Directory

```bash
cd C:\Users\HP\Desktop\watchup-agent
```

### 1.2 Initialize Git (if not done)

```bash
git init
```

### 1.3 Add All Files

```bash
git add .
```

### 1.4 Commit

```bash
git commit -m "Initial commit: Watchup Server Agent v1.0.0"
```

### 1.5 Create GitHub Repository

1. Go to https://github.com/new
2. Repository name: `watchup-agent`
3. Description: `Lightweight server monitoring agent for Watchup platform`
4. Choose: **Public**
5. **Do NOT** check "Initialize with README"
6. Click **"Create repository"**

### 1.6 Connect and Push

```bash
# Add remote (use YOUR GitHub username)
git remote add origin https://github.com/tomurashigaraki22/watchup-agent.git

# Push to GitHub
git branch -M main
git push -u origin main
```

### 1.7 Verify

Visit: https://github.com/tomurashigaraki22/watchup-agent

You should see all your files!

---

## ✅ Step 2: Install on Your VPS (10 minutes)

### 2.1 Connect to VPS

```bash
ssh your-username@your-vps-ip
```

Example:
```bash
ssh root@192.168.1.100
# or
ssh ubuntu@your-domain.com
```

### 2.2 Install Dependencies

**For Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y git curl wget
```

**For CentOS/RHEL:**
```bash
sudo yum install -y git curl wget
```

### 2.3 Install Go

```bash
# Download Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# Extract
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

You should see: `go version go1.21.0 linux/amd64`

### 2.4 Clone Your Repository

```bash
git clone https://github.com/tomurashigaraki22/watchup-agent.git
cd watchup-agent
```

### 2.5 Build the Agent

```bash
# Install dependencies
go mod tidy

# Build
go build -o watchup-agent cmd/agent/main.go

# Verify build
ls -lh watchup-agent
```

You should see a file around 10-20MB.

### 2.6 Install Binary

```bash
# Move to system directory
sudo mv watchup-agent /usr/local/bin/

# Make executable
sudo chmod +x /usr/local/bin/watchup-agent

# Verify installation
which watchup-agent
```

Should output: `/usr/local/bin/watchup-agent`

### 2.7 Create Configuration

```bash
# Create config directory
sudo mkdir -p /etc/watchup

# Copy sample config
sudo cp config.yaml /etc/watchup/config.yaml

# Set permissions
sudo chmod 600 /etc/watchup/config.yaml
```

### 2.8 Create Systemd Service

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

### 2.9 Enable and Start Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable auto-start on boot
sudo systemctl enable watchup-agent

# Start the service
sudo systemctl start watchup-agent
```

### 2.10 View Logs and Register

```bash
# View real-time logs
sudo journalctl -u watchup-agent -f
```

You'll see prompts like:
```
Enter your Project ID: 
Enter your Master API Key: 
Enter Server Identifier (press Enter for auto-generated): 
```

**Enter your credentials:**
- **Project ID**: Get from Watchup dashboard (e.g., `proj_12345`)
- **Master API Key**: Get from Watchup project settings (e.g., `ma_abcdef123`)
- **Server Identifier**: Give your server a name (e.g., `production-api-1`)

After registration, you should see:
```
✓ Registration successful!
Server Key: srv_89sd0a
Configuration saved. Agent is ready to start monitoring.

✓ Agent registered successfully
Project ID: proj_12345
Server: production-api-1
Starting monitoring loop...
[14:23:15] CPU: 12.3%, RAM: 45.6%
```

---

## ✅ Step 3: Verify Everything Works (2 minutes)

### 3.1 Check Service Status

```bash
sudo systemctl status watchup-agent
```

Should show: `Active: active (running)`

### 3.2 View Recent Logs

```bash
sudo journalctl -u watchup-agent -n 20
```

Should show monitoring data.

### 3.3 Check Process

```bash
ps aux | grep watchup-agent
```

Should show the running process.

### 3.4 Test API Connectivity

```bash
curl -I https://api.watchup.site
```

Should return HTTP 200 or similar.

---

## ✅ Step 4: Create GitHub Release (Optional, 3 minutes)

### 4.1 Tag Version

```bash
# In your local project directory
cd C:\Users\HP\Desktop\watchup-agent

git tag -a v1.0.0 -m "Release v1.0.0: Initial production release"
git push origin v1.0.0
```

### 4.2 Create Release on GitHub

1. Go to https://github.com/tomurashigaraki22/watchup-agent
2. Click **"Releases"** (right sidebar)
3. Click **"Create a new release"**
4. Choose tag: `v1.0.0`
5. Release title: `v1.0.0 - Initial Release`
6. Description:
   ```
   ## Watchup Server Agent v1.0.0
   
   First production release of the Watchup Server Agent.
   
   ### Features
   - CPU, RAM, and process monitoring
   - Sustained spike detection
   - One-agent-per-free-project enforcement
   - Dynamic configuration updates
   - Systemd service support
   
   ### Installation
   See [README.md](README.md) for installation instructions.
   
   ### Quick Install
   ```bash
   git clone https://github.com/tomurashigaraki22/watchup-agent.git
   cd watchup-agent
   go build -o watchup-agent cmd/agent/main.go
   sudo mv watchup-agent /usr/local/bin/
   ```
   ```
7. Click **"Publish release"**

---

## 🎉 You're Done!

Your Watchup Server Agent is now:

✅ Deployed to GitHub  
✅ Running on your VPS  
✅ Monitoring your server  
✅ Sending data to Watchup  

---

## 📋 Quick Reference

### Service Commands

```bash
# Start
sudo systemctl start watchup-agent

# Stop
sudo systemctl stop watchup-agent

# Restart
sudo systemctl restart watchup-agent

# Status
sudo systemctl status watchup-agent

# Logs (real-time)
sudo journalctl -u watchup-agent -f

# Logs (last 50 lines)
sudo journalctl -u watchup-agent -n 50
```

### Configuration

```bash
# Edit config
sudo nano /etc/watchup/config.yaml

# After editing, restart
sudo systemctl restart watchup-agent
```

### Troubleshooting

```bash
# Check if binary exists
ls -l /usr/local/bin/watchup-agent

# Check if config exists
ls -l /etc/watchup/config.yaml

# Run manually (for debugging)
sudo /usr/local/bin/watchup-agent /etc/watchup/config.yaml

# Check resource usage
top -p $(pgrep watchup-agent)
```

---

## 🆘 Common Issues

### "Permission denied" when pushing to GitHub

**Solution**: Use HTTPS and enter your GitHub username/password (or personal access token)

```bash
git remote set-url origin https://github.com/tomurashigaraki22/watchup-agent.git
git push -u origin main
```

### "Agent won't start" on VPS

**Solution**: Check logs for errors

```bash
sudo journalctl -u watchup-agent -n 100
```

### "Registration failed"

**Solution**: Verify credentials and network

```bash
# Test API connectivity
curl -I https://api.watchup.site

# Check credentials in Watchup dashboard
```

### "This project already has an agent"

**Solution**: Free projects are limited to one agent. Either:
- Use the existing agent
- Upgrade to paid plan
- Use a different project

---

## 📞 Need Help?

- **Documentation**: See all files in `docs/` folder
- **Quick Start**: [GET_STARTED.md](GET_STARTED.md)
- **Full Guide**: [DEPLOYMENT.md](DEPLOYMENT.md)
- **Commands**: [COMMANDS.md](COMMANDS.md)
- **GitHub Issues**: https://github.com/tomurashigaraki22/watchup-agent/issues

---

## 🎯 What's Next?

1. **Monitor**: Check your Watchup dashboard for metrics
2. **Customize**: Adjust alert thresholds in config.yaml
3. **Scale**: Install on more servers
4. **Enhance**: Add more features (see PATTERN.md)

**Congratulations! Your server is now monitored by Watchup! 🚀**
