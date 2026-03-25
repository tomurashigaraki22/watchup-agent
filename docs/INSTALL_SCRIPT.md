# Installation Script Documentation

## Overview

The `install/install.sh` script provides a fully automated installation of the Watchup Server Agent on Linux systems.

## Features

✅ **Automatic OS Detection** - Supports Ubuntu, Debian, CentOS, RHEL, Fedora, Arch, Alpine  
✅ **Automatic Go Installation** - Installs or upgrades Go if needed  
✅ **Multi-Architecture Support** - AMD64, ARM64, ARMv6/v7  
✅ **Dependency Management** - Installs git, curl, wget, tar automatically  
✅ **Source Build** - Clones and builds from GitHub  
✅ **Systemd Integration** - Creates and enables service  
✅ **Idempotent** - Safe to run multiple times  

---

## Usage

### One-Line Installation

```bash
curl -s https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent/main/install/install.sh | sudo bash
```

Or download and run:

```bash
wget https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent/main/install/install.sh
chmod +x install.sh
sudo ./install.sh
```

---

## What It Does

### 1. System Detection

Detects:
- Operating system (Ubuntu, Debian, CentOS, etc.)
- Architecture (AMD64, ARM64, ARMv6/v7)
- Package manager (apt, yum, dnf, pacman, apk)

### 2. Dependency Installation

Installs required packages:
- `git` - For cloning repository
- `curl` or `wget` - For downloading
- `tar` - For extracting archives

**Ubuntu/Debian:**
```bash
apt-get update
apt-get install -y git curl wget tar
```

**CentOS/RHEL/Fedora:**
```bash
dnf install -y git curl wget tar
# or
yum install -y git curl wget tar
```

**Arch/Manjaro:**
```bash
pacman -Sy --noconfirm git curl wget tar
```

**Alpine:**
```bash
apk add --no-cache git curl wget tar bash
```

### 3. Go Installation

Checks if Go is installed:
- If not installed → Installs Go 1.21.0
- If version < 1.20 → Upgrades to Go 1.21.0
- If version >= 1.20 → Skips installation

**Installation process:**
1. Detects system architecture
2. Downloads appropriate Go tarball from https://go.dev/dl/
3. Extracts to `/usr/local/go`
4. Adds to PATH in `/etc/profile` and `~/.bashrc`
5. Verifies installation

**Supported architectures:**
- `amd64` (x86_64)
- `arm64` (aarch64)
- `armv6l` (Raspberry Pi)

### 4. Agent Build

1. Clones repository to `/tmp/watchup-agent-build`
2. Runs `go mod tidy` to install dependencies
3. Builds binary: `go build -o watchup-agent cmd/agent/main.go`
4. Moves binary to `/usr/local/bin/watchup-agent`
5. Sets executable permissions (755)
6. Cleans up build directory

### 5. Configuration

Creates `/etc/watchup/config.yaml` with defaults:

```yaml
server_key: ""
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
```

Sets file permissions to `600` (owner read/write only).

### 6. Systemd Service

Creates `/etc/systemd/system/watchup-agent.service`:

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

Then:
1. Reloads systemd daemon
2. Enables service for auto-start on boot

---

## Installation Output

```
=== Watchup Server Agent Installation ===

Starting installation...

Detected OS: ubuntu

Installing dependencies...
[package installation output]

Go is not installed.
Installing Go 1.21.0...
Downloading Go from https://go.dev/dl/go1.21.0.linux-amd64.tar.gz...
Extracting Go...
Go installed successfully: go version go1.21.0 linux/amd64

Building Watchup Agent from source...
Cloning repository...
Compiling agent...
Installing binary to /usr/local/bin...
Agent built and installed successfully.

Creating configuration directory...
Creating default configuration...
Configuration file created at /etc/watchup/config.yaml

Creating systemd service...
Reloading systemd...
Enabling service...
Systemd service created and enabled.

=== Installation Complete ===

The Watchup Server Agent has been installed successfully!

Next steps:
1. Start the agent: sudo systemctl start watchup-agent
2. View logs: sudo journalctl -u watchup-agent -f
3. The agent will prompt for registration on first run
4. Check status: sudo systemctl status watchup-agent

For more information, visit:
https://github.com/tomurashigaraki22/watchup-agent
```

---

## Requirements

### System Requirements

- **OS**: Linux with systemd
- **Root Access**: Required (use sudo)
- **Internet**: Required for downloading Go and cloning repository
- **Disk Space**: ~500MB (for Go and build artifacts)

### Network Requirements

- Access to https://go.dev (for Go download)
- Access to https://github.com (for repository clone)
- Access to https://api.watchup.site (for agent operation)

---

## Troubleshooting

### "Permission denied"

**Cause**: Not running as root

**Solution**:
```bash
sudo bash install.sh
```

### "curl: command not found"

**Cause**: curl not installed (rare, script should install it)

**Solution**:
```bash
# Ubuntu/Debian
sudo apt-get install curl

# CentOS/RHEL
sudo yum install curl
```

### "Unsupported architecture"

**Cause**: Running on unsupported CPU architecture

**Supported**: AMD64, ARM64, ARMv6/v7

**Solution**: Build manually on your architecture or request support

### "Go installation failed"

**Cause**: Network issue or unsupported architecture

**Solution**:
1. Check internet connection
2. Verify architecture: `uname -m`
3. Install Go manually: https://go.dev/doc/install

### "git clone failed"

**Cause**: Network issue or repository not accessible

**Solution**:
1. Check internet connection
2. Verify repository exists: https://github.com/tomurashigaraki22/watchup-agent
3. Check firewall settings

### "Build failed"

**Cause**: Missing dependencies or Go version issue

**Solution**:
```bash
# Check Go version
go version  # Should be 1.20+

# Try manual build
git clone https://github.com/tomurashigaraki22/watchup-agent.git
cd watchup-agent
go mod tidy
go build -o watchup-agent cmd/agent/main.go
```

---

## Customization

### Change Go Version

Edit `install.sh`:

```bash
GO_VERSION="1.22.0"  # Change this line
```

### Change GitHub Repository

Edit `install.sh`:

```bash
GITHUB_REPO="your-username/watchup-agent"  # Change this line
```

### Change Installation Directory

Edit `install.sh`:

```bash
INSTALL_DIR="/opt/watchup"  # Change this line
```

---

## Uninstallation

The script doesn't provide uninstallation. To uninstall manually:

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

# Remove Go (optional, if installed by script)
sudo rm -rf /usr/local/go
```

---

## Security Considerations

### Running as Root

The script requires root access to:
- Install system packages
- Write to `/usr/local/bin`
- Create systemd service
- Write to `/etc/watchup`

### Downloaded Content

The script downloads:
- Go from official source: https://go.dev/dl/
- Agent source from GitHub: https://github.com/tomurashigaraki22/watchup-agent

Always verify the script content before running with sudo.

### Verification

Before running, inspect the script:

```bash
curl -s https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent/main/install/install.sh | less
```

---

## Advanced Usage

### Offline Installation

1. Download Go tarball manually
2. Clone repository manually
3. Modify script to use local files
4. Run script

### Custom Build Flags

Edit the build command in `install.sh`:

```bash
go build -ldflags="-s -w" -o "${BINARY_NAME}" cmd/agent/main.go
```

### Skip Go Installation

If Go is already installed:

```bash
# Comment out the Go installation check
# if ! check_go; then
#     install_go
# fi
```

---

## Supported Distributions

| Distribution | Package Manager | Status |
|--------------|----------------|--------|
| Ubuntu 18.04+ | apt | ✅ Tested |
| Debian 9+ | apt | ✅ Tested |
| CentOS 7+ | yum | ✅ Tested |
| CentOS 8+ | dnf | ✅ Tested |
| RHEL 7+ | yum | ✅ Tested |
| RHEL 8+ | dnf | ✅ Tested |
| Fedora 28+ | dnf | ✅ Tested |
| Arch Linux | pacman | ✅ Tested |
| Manjaro | pacman | ✅ Tested |
| Alpine Linux | apk | ✅ Tested |

---

## Contributing

To improve the installation script:

1. Test on your distribution
2. Report issues on GitHub
3. Submit pull requests with fixes
4. Add support for new distributions

---

## License

MIT License - Same as the main project
