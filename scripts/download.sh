#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-latest}"
OUT_DIR="${2:-.}"
REPO="shs3131/dbasement"

if [ "$VERSION" = "latest" ]; then
    TAG=$(curl -sL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
else
    TAG="$VERSION"
fi

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
esac

case "$OS" in
    linux) FILE="dbasement-linux-$ARCH.tar.gz" ;;
    darwin) FILE="dbasement-darwin-$ARCH.tar.gz" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

URL="https://github.com/$REPO/releases/download/$TAG/$FILE"

echo "Downloading Dbasement $TAG for $OS/$ARCH..."
curl -sL "$URL" -o "/tmp/$FILE"
tar xzf "/tmp/$FILE" -C "$OUT_DIR"
rm "/tmp/$FILE"
chmod +x "$OUT_DIR/dbasement"
echo "Installed at: $OUT_DIR/dbasement"
