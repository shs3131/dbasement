#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
EXE_PATH="$ROOT_DIR/dbasement"
PROJECT_PATH="${1:-.}"

if [ ! -f "$EXE_PATH" ]; then
    echo "Downloading Dbasement..."
    REPO="shs3131/dbasement"
    TAG=$(curl -sL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
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
    curl -sL "$URL" -o "/tmp/$FILE"
    tar xzf "/tmp/$FILE" -C "$ROOT_DIR"
    rm "/tmp/$FILE"
    chmod +x "$EXE_PATH"
    echo "Downloaded: $EXE_PATH"
fi

exec "$EXE_PATH" --project "$PROJECT_PATH"
