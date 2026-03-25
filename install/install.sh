#!/bin/bash

set -e

BINARY_NAME="watchup-agent"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/watchup"
SERVICE_FILE="/etc/systemd/system/watchup-agent.service"
GO_VERSION="1.21.0"
GITHUB_REPO="tomurashigaraki22/watchup-agent"

echo "=== Watchup Server Agent Installation ==="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Please run as root (use sudo)"
    exit 1
fi

# Detect OS and distribution
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        VERSION=$VERSION_ID
    elif [ -f /etc/redhat-release ]; then
        OS="rhel"
    elif [ -f /etc/debian_version ]; thenn
        OS="debian"
    else
        OS=$(uname -s)
    fi
    
    echo "Detected OS: $OS"
}

# Install dependencies based on OS
install_dependencies() {
    echo "Installing dependencies..."
    
    case "$OS" in
        ubuntu|debian)
            apt-get update
            apt-get install -y git curl wget tar
            ;;
        centos|rhel|fedora)
            if command -v dnf &> /dev/null; then
                dnf install -y git curl wget tar
            else
                yum install -y git curl wget tar
            fi
            ;;
        arch|manjaro)
            pacman -Sy --noconfirm git curl wget tar
            ;;
        alpine)
            apk add --no-cache git curl wget tar bash
            ;;
        *)
            echo "Warning: Unknown OS. Attempting to continue..."
            ;;
    esac
}

# Check if Go is installed
check_go() {
    if command -v go &> /dev/null; then
        GO_INSTALLED_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        echo "Go is already installed: $GO_INSTALLED_VERSION"
        
        # Check if version is sufficient (1.20+)
        REQUIRED_VERSION="1.20"
        if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_INSTALLED_VERSION" | sort -V | head -n1)" = "$REQUIRED_VERSION" ]; then
            echo "Go version is sufficient."
            return 0
        else
            echo "Go version is too old. Installing newer version..."
            return 1
        fi
    else
        echo "Go is not installed."
        return 1
    fi
}

# Install Go
install_go() {
    echo "Installing Go ${GO_VERSION}..."
    
    # Detect architecture
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)
            GO_ARCH="amd64"
            ;;
        aarch64|arm64)
            GO_ARCH="arm64"
            ;;
        armv7l|armv6l)
            GO_ARCH="armv6l"
            ;;
        *)
            echo "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    GO_TARBALL="go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    GO_URL="https://go.dev/dl/${GO_TARBALL}"
    
    echo "Downloading Go from ${GO_URL}..."
    
    # Download Go
    if command -v curl &> /dev/null; then
        curl -L -o "/tmp/${GO_TARBALL}" "${GO_URL}"
    elif command -v wget &> /dev/null; then
        wget -O "/tmp/${GO_TARBALL}" "${GO_URL}"
    else
        echo "Error: curl or wget is required"
        exit 1
    fi
    
    # Remove old Go installation if exists
    if [ -d "/usr/local/go" ]; then
        echo "Removing old Go installation..."
        rm -rf /usr/local/go
    fi
    
    # Extract Go
    echo "Extracting Go..."
    tar -C /usr/local -xzf "/tmp/${GO_TARBALL}"
    
    # Add Go to PATH
    if ! grep -q "/usr/local/go/bin" /etc/profile; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    fi
    
    # Add to current session
    export PATH=$PATH:/usr/local/go/bin
    
    # Add to root's bashrc
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    fi
    
    # Cleanup
    rm -f "/tmp/${GO_TARBALL}"
    
    # Verify installation
    if command -v go &> /dev/null; then
        echo "Go installed successfully: $(go version)"
    else
        echo "Error: Go installation failed"
        exit 1
    fi
}

# Build agent from source
build_agent() {
    echo "Building Watchup Agent from source..."
    
    BUILD_DIR="/tmp/watchup-agent-build"
    
    # Clean up any previous build
    rm -rf "$BUILD_DIR"
    
    # Clone repository
    echo "Cloning repository..."
    git clone "https://github.com/${GITHUB_REPO}.git" "$BUILD_DIR"
    
    cd "$BUILD_DIR"
    
    # Build
    echo "Compiling agent..."
    go mod tidy
    go build -o "${BINARY_NAME}" cmd/agent/main.go
    
    # Install binary
    echo "Installing binary to ${INSTALL_DIR}..."
    mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    
    # Cleanup
    cd /
    rm -rf "$BUILD_DIR"
    
    echo "Agent built and installed successfully."
}

# Create configuration
create_config() {
    echo "Creating configuration directory..."
    mkdir -p "${CONFIG_DIR}"
    
    if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
        echo "Creating default configuration..."
        cat > "${CONFIG_DIR}/config.yaml" << 'EOF'
server_key: ""
project_id: ""
server_identifier: ""
sampling_interval: 5
api_endpoint: "https://watchup.space"
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
        chmod 600 "${CONFIG_DIR}/config.yaml"
        echo "Configuration file created at ${CONFIG_DIR}/config.yaml"
    else
        echo "Configuration file already exists. Skipping..."
    fi
}

# Create systemd service
create_service() {
    echo "Creating systemd service..."
    cat > "${SERVICE_FILE}" << EOF
[Unit]
Description=Watchup Server Monitoring Agent
After=network.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/${BINARY_NAME}
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

    # Reload systemd
    echo "Reloading systemd..."
    systemctl daemon-reload
    
    # Enable service
    echo "Enabling service..."
    systemctl enable watchup-agent
    
    echo "Systemd service created and enabled."
}

# Main installation flow
main() {
    echo "Starting installation..."
    echo ""
    
    # Detect OS
    detect_os
    echo ""
    
    # Install dependencies
    install_dependencies
    echo ""
    
    # Check and install Go if needed
    if ! check_go; then
        install_go
    fi
    echo ""
    
    # Build agent from source
    build_agent
    echo ""
    
    # Create configuration
    create_config
    echo ""
    
    # Create systemd service
    create_service
    echo ""
    
    echo "=== Installation Complete ==="
    echo ""
    echo "The Watchup Server Agent has been installed successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Start the agent: sudo systemctl start watchup-agent"
    echo "2. View logs: sudo journalctl -u watchup-agent -f"
    echo "3. The agent will prompt for registration on first run"
    echo "4. Check status: sudo systemctl status watchup-agent"
    echo ""
    echo "For more information, visit:"
    echo "https://github.com/${GITHUB_REPO}"
    echo ""
}

# Run main installation
main
