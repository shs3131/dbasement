#!/usr/bin/env bash
set -euo pipefail

# Dbasement installer
# Usage: curl -sL https://raw.githubusercontent.com/shs3131/dbasement/main/install.sh | bash
# Or:    bash install.sh [destination]

DEST="${1:-.}"
REPO="shs3131/dbasement"

echo "Installing Dbasement to: $DEST"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Get latest release tag
echo "Detecting latest release..."
TAG=$(curl -sL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
if [ -z "$TAG" ]; then
    echo "Failed to detect latest release"
    exit 1
fi
echo "Latest release: $TAG"

# Download and extract
ASSET="dbasement-$OS-$ARCH.tar.gz"
URL="https://github.com/$REPO/releases/download/$TAG/$ASSET"
echo "Downloading: $URL"
curl -sL "$URL" -o "/tmp/$ASSET"
tar xzf "/tmp/$ASSET" -C "$DEST"
rm "/tmp/$ASSET"
chmod +x "$DEST/dbasement"

echo ""
echo "Dbasement $TAG installed at: $DEST/dbasement"
echo ""
echo "Next steps:"
echo "  1. Configure your AI client (see MCP_CONFIGS.md)"
echo "  2. Start a new AI session"
echo ""
