#!/bin/bash

GITHUB_USER="FloomAI"
GITHUB_REPO="FloomCLI"
INSTALL_DIR="$HOME/.floom/bin"

mkdir -p "$INSTALL_DIR"

error() {
  echo "Error: $1" >&2
  exit 1
}

info() {
  echo "Info: $1"
}

# Detect OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case $ARCH in
    x86_64) ARCH="amd64";;
    arm64) ARCH="arm64";;
    *) error "Unsupported architecture: $ARCH";;
esac

# Fetch the latest release data from GitHub
RELEASE_DATA=$(curl -s "https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/releases/latest")

# Extract version from the release tag
VERSION=$(echo "$RELEASE_DATA" | grep '"tag_name":' | sed -E 's/.*"v([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    error "Failed to extract version from the latest release"
fi

# Construct the binary name based on OS, Architecture, and dynamically determined version
BINARY="floom-v${VERSION}-${OS}-${ARCH}"

# Use the GitHub API data to find the correct binary URL
DOWNLOAD_URL=$(echo "$RELEASE_DATA" | grep "browser_download_url.*${BINARY}" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    error "Failed to find a download URL for ${OS}/${ARCH}"
fi

info "Downloading Floom CLI from ${DOWNLOAD_URL}"

# Download and install the binary
curl -sSL "$DOWNLOAD_URL" -o "${INSTALL_DIR}/floom"
if [ $? -ne 0 ]; then
    error "Download failed"
fi

chmod +x "${INSTALL_DIR}/floom"

# Update PATH
if ! grep -q "$INSTALL_DIR" <<< "$PATH"; then
    echo 'export PATH="$HOME/.floom/bin:$PATH"' >> "$HOME/.bashrc"
    echo 'export PATH="$HOME/.floom/bin:$PATH"' >> "$HOME/.zshrc"
    info "Added $INSTALL_DIR to your PATH."
    exec "${SHELL}" # Reload shell
else
    info "$INSTALL_DIR is already in your PATH."
fi

info "Installation completed successfully. Type 'floom' to get started!"
