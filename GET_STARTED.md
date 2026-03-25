# Get Started with Watchup Server Agent

This guide will help you deploy the agent to GitHub and install it on your VPS in under 10 minutes.

---

## 🚀 Quick Start (3 Steps)

### 1️⃣ Deploy to GitHub (2 minutes)

```bash
# In your project directory
cd watchup-agent

# Add to GitHub (replace YOUR_USERNAME)
git init
git add .
git commit -m "Initial commit: Watchup Server Agent"
git remote add origin https://github.com/YOUR_USERNAME/watchup-agent.git
git branch -M main
git push -u origin main
```

✅ **Done!** Your code is now on GitHub.

### 2️⃣ Install on VPS (2 minutes - Fully Automatic!)

```bash
# SSH to your VPS
ssh user@your-vps-ip

# Run one command - it does everything!
curl -s https://raw.githubusercontent.com/YOUR_USERNAME/watchup-agent/main/install/install.sh | sudo bash
```

**What it does automatically:**
- ✅ Detects your OS (Ubuntu, Debian, CentOS, RHEL, Fedora, Arch, Alpine)
- ✅ Installs dependencies (git, curl, wget, tar)
- ✅ Installs Go if not present (or upgrades if < 1.20)
- ✅ Clones your repository
- ✅ Builds the agent from source
- ✅ Installs the binary
- ✅ Creates systemd service
- ✅ Enables auto-start on boot

**Supported architectures:** AMD64, ARM64, ARMv6/v7 (Raspberry Pi)

Or build from source manually:

```bash
# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Clone and build
git clone https://github.com/YOUR_USERNAME/watchup-agent.git
cd watchup-agent
go mod tidy
go build -o watchup-agent cmd/agent/main.go
sudo mv watchup-agent /usr/local/bin/
```

### 3️⃣ Start and Register (2 minutes)

```bash
# Create systemd service
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

# Start service
sudo systemctl daemon-reload
sudo systemctl enable watchup-agent
sudo systemctl start watchup-agent

# View logs and register
sudo journalctl -u watchup-agent -f
```

When prompted, enter:
- **Project ID**: From Watchup dashboard
- **Master API Key**: From Watchup project settings
- **Server Identifier**: A name for your server (e.g., "production-api")

✅ **Done!** Your agent is monitoring your server.

---

## 📋 What You Need

### For GitHub Deployment
- GitHub account
- Git installed locally
- Your project files

### For VPS Installation
- Linux VPS (Ubuntu, Debian, CentOS, etc.)
- SSH access with sudo
- Internet connectivity

---

## 🔍 Verify Everything Works

```bash
# Check service status
sudo systemctl status watchup-agent

# View logs
sudo journalctl -u watchup-agent -n 20

# Check if agent is running
ps aux | grep watchup-agent
```

You should see:
```
✓ Agent registered successfully
Project ID: proj_12345
Server: your-server-name
Starting monitoring loop...
[14:23:15] CPU: 12.3%, RAM: 45.6%
```

---

## 📚 Next Steps

1. **View Metrics**: Check your Watchup dashboard
2. **Customize Alerts**: Edit `/etc/watchup/config.yaml`
3. **Add More Servers**: Repeat VPS installation on other servers
4. **Monitor Logs**: `sudo journalctl -u watchup-agent -f`

---

## 🆘 Need Help?

### Common Issues

**"Agent won't start"**
```bash
sudo journalctl -u watchup-agent -n 50
sudo /usr/local/bin/watchup-agent /etc/watchup/config.yaml
```

**"Registration failed"**
- Check your Project ID and Master API Key
- Verify internet connectivity: `curl -I https://api.watchup.site`
- Free projects are limited to one agent

**"Can't connect to GitHub"**
```bash
# Use HTTPS instead of SSH
git remote set-url origin https://github.com/YOUR_USERNAME/watchup-agent.git
```

### Get Support

- 📖 **Full Documentation**: [README.md](README.md)
- 🚀 **Deployment Guide**: [DEPLOYMENT.md](DEPLOYMENT.md)
- 💻 **Command Reference**: [COMMANDS.md](COMMANDS.md)
- 🐛 **Report Issues**: GitHub Issues
- 📧 **Email**: support@watchup.site

---

## 🎯 What's Next?

Your agent is now:
- ✅ Monitoring CPU, RAM, and processes every 5 seconds
- ✅ Detecting sustained spikes (prevents false positives)
- ✅ Sending alerts to Watchup dashboard
- ✅ Auto-starting on server reboot

**Congratulations!** 🎉 Your server is now monitored by Watchup.
