#!/bin/bash

# vssh installation script for Linux and macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/ncecere/vssh/main/install.sh | bash

set -e

# Configuration
REPO="ncecere/vssh"
BINARY_NAME="vssh"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os
    local arch
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        *)          log_error "Unsupported operating system: $(uname -s)"; exit 1 ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        *)              log_error "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release version
get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$version" ]; then
        log_error "Failed to get latest version"
        exit 1
    fi
    
    echo "$version"
}

# Download and install binary
install_binary() {
    local platform="$1"
    local version="$2"
    local binary_name="${BINARY_NAME}-${version}-${platform}"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${binary_name}"
    local temp_file="/tmp/${binary_name}"
    
    log_info "Downloading ${binary_name}..."
    
    # Download binary
    if ! curl -fsSL "$download_url" -o "$temp_file"; then
        log_error "Failed to download binary from $download_url"
        exit 1
    fi
    
    # Make executable
    chmod +x "$temp_file"
    
    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        log_info "Creating install directory: $INSTALL_DIR"
        sudo mkdir -p "$INSTALL_DIR"
    fi
    
    # Install binary
    log_info "Installing to $INSTALL_DIR/$BINARY_NAME..."
    
    if [ -w "$INSTALL_DIR" ]; then
        mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
    else
        sudo mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    log_success "vssh installed successfully!"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version
        version=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        log_success "vssh is installed and available in PATH"
        log_info "Version: $version"
    else
        log_warning "vssh is installed but not in PATH"
        log_info "Add $INSTALL_DIR to your PATH or run: export PATH=\"$INSTALL_DIR:\$PATH\""
    fi
}

# Show post-installation instructions
show_instructions() {
    echo
    log_info "Next steps:"
    echo "  1. Initialize configuration: ${BINARY_NAME} init"
    echo "  2. Edit config file: ~/.config/vssh/config.yaml"
    echo "  3. Connect to a server: ${BINARY_NAME} user@server.com"
    echo
    log_info "Documentation:"
    echo "  - README: https://github.com/${REPO}/blob/main/README.md"
    echo "  - Config: https://github.com/${REPO}/blob/main/CONFIG.md"
    echo
}

# Main installation function
main() {
    log_info "Installing vssh..."
    
    # Check dependencies
    if ! command -v curl >/dev/null 2>&1; then
        log_error "curl is required but not installed"
        exit 1
    fi
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    log_info "Detected platform: $platform"
    
    # Get latest version
    local version
    version=$(get_latest_version)
    log_info "Latest version: $version"
    
    # Install binary
    install_binary "$platform" "$version"
    
    # Verify installation
    verify_installation
    
    # Show instructions
    show_instructions
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "vssh installation script"
        echo
        echo "Usage: $0 [options]"
        echo
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --version, -v  Install specific version"
        echo
        echo "Environment variables:"
        echo "  INSTALL_DIR    Installation directory (default: /usr/local/bin)"
        echo
        echo "Examples:"
        echo "  $0                    # Install latest version"
        echo "  $0 --version v1.0.0   # Install specific version"
        echo "  INSTALL_DIR=~/bin $0  # Install to custom directory"
        exit 0
        ;;
    --version|-v)
        if [ -z "${2:-}" ]; then
            log_error "Version argument required"
            exit 1
        fi
        VERSION="$2"
        ;;
esac

# Run main installation
main
