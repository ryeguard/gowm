#!/bin/bash
# Installation script for gowm (unified CLI and MCP server)
# Usage: curl -sSL https://raw.githubusercontent.com/ryeguard/gowm/main/install.sh | bash

set -e

# Configuration
GITHUB_REPO="ryeguard/gowm"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --install-dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --help)
            echo "gowm Installation Script"
            echo ""
            echo "Installs the unified gowm binary with CLI and MCP server functionality."
            echo ""
            echo "Usage: install.sh [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --install-dir DIR Install to specified directory (default: /usr/local/bin)"
            echo "  --help            Show this help message"
            echo ""
            echo "Environment variables:"
            echo "  INSTALL_DIR       Installation directory (default: /usr/local/bin)"
            echo ""
            echo "Commands:"
            echo "  gowm get-weather  Get weather data via CLI"
            echo "  gowm mcp          Start MCP server for LLM integration"
            echo "  gowm version      Show version information"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Detect OS and architecture
detect_platform() {
    local os=""
    local arch=""

    # Detect OS
    case "$(uname -s)" in
        Linux*)
            os="Linux"
            ;;
        Darwin*)
            os="Darwin"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            print_error "Windows is not supported by this script. Please download binaries manually from:"
            print_error "https://github.com/${GITHUB_REPO}/releases/latest"
            exit 1
            ;;
        *)
            print_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac

    # Detect architecture
    case "$(uname -m)" in
        x86_64)
            arch="x86_64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac

    echo "${os}_${arch}"
}

# Get latest release version
get_latest_version() {
    local version=$(curl -sSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | \
                    grep '"tag_name":' | \
                    sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$version" ]; then
        print_error "Failed to fetch latest version"
        exit 1
    fi

    echo "$version"
}

# Download and install binary
install_binary() {
    local binary_name=$1
    local version=$2
    local platform=$3

    local archive_name="gowm_${version#v}_${platform}.tar.gz"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${archive_name}"
    local tmp_dir=$(mktemp -d)

    print_info "Downloading ${binary_name} ${version} for ${platform}..."

    if ! curl -sSL "$download_url" -o "${tmp_dir}/${archive_name}"; then
        print_error "Failed to download ${binary_name}"
        rm -rf "$tmp_dir"
        exit 1
    fi

    print_info "Extracting ${binary_name}..."
    tar -xzf "${tmp_dir}/${archive_name}" -C "$tmp_dir"

    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "${tmp_dir}/${binary_name}" "${INSTALL_DIR}/${binary_name}"
        chmod +x "${INSTALL_DIR}/${binary_name}"
    else
        print_info "Installing to ${INSTALL_DIR} requires sudo..."
        sudo mv "${tmp_dir}/${binary_name}" "${INSTALL_DIR}/${binary_name}"
        sudo chmod +x "${INSTALL_DIR}/${binary_name}"
    fi

    rm -rf "$tmp_dir"
    print_info "${binary_name} installed successfully to ${INSTALL_DIR}/${binary_name}"
}

# Verify installation
verify_installation() {
    local binary_name=$1

    if ! command -v "$binary_name" &> /dev/null; then
        print_warn "${binary_name} is not in PATH. You may need to add ${INSTALL_DIR} to your PATH:"
        echo ""
        echo "    export PATH=\"\$PATH:${INSTALL_DIR}\""
        echo ""
        echo "Add this line to your ~/.bashrc, ~/.zshrc, or equivalent shell config file."
        return 1
    fi

    return 0
}

# Main installation
main() {
    print_info "Starting gowm installation..."
    echo ""

    # Detect platform
    platform=$(detect_platform)
    print_info "Detected platform: $platform"

    # Get latest version
    version=$(get_latest_version)
    print_info "Latest version: $version"
    echo ""

    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        print_info "Creating installation directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR" 2>/dev/null || sudo mkdir -p "$INSTALL_DIR"
    fi

    # Install gowm binary
    install_binary "gowm" "$version" "$platform"

    echo ""
    print_info "Installation complete!"
    echo ""

    # Verify and show usage
    if verify_installation "gowm"; then
        print_info "gowm installed successfully!"
    else
        print_info "gowm installed at: ${INSTALL_DIR}/gowm"
    fi

    echo ""
    print_info "Usage:"
    echo ""
    echo "  CLI - Get weather data:"
    echo "    gowm get-weather 'stockholm,sweden' --api-key=YOUR_API_KEY"
    echo ""
    echo "  MCP Server - For LLM integration:"
    echo "    gowm mcp"
    echo ""
    print_info "For Claude Desktop, add to your claude_desktop_config.json:"
    echo ""
    echo '  {
    "mcpServers": {
      "weather": {
        "command": "'${INSTALL_DIR}'/gowm",
        "args": ["mcp"],
        "env": {
          "OWM_API_KEY": "YOUR_API_KEY"
        }
      }
    }
  }'
    echo ""
    print_info "Get your OpenWeatherMap API key from: https://openweathermap.org/"
}

# Run main function
main
