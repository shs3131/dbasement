#!/bin/bash
set -euo pipefail

VERSION="${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}"
OUTDIR="dist"

echo "Building Dbasement $VERSION"

mkdir -p "$OUTDIR"

declare -A TARGETS=(
  ["linux/amd64"]="dbasement-linux-amd64"
  ["linux/arm64"]="dbasement-linux-arm64"
  ["darwin/amd64"]="dbasement-darwin-amd64"
  ["darwin/arm64"]="dbasement-darwin-arm64"
  ["windows/amd64"]="dbasement-windows-amd64.exe"
)

for target in "${!TARGETS[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "$target"
  output="${TARGETS[$target]}"
  
  echo "Building $GOOS/$GOARCH -> $OUTDIR/$output"
  
  GOOS="$GOOS" GOARCH="$GOARCH" go build \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o "$OUTDIR/$output" \
    ./cmd/dbasement/
done

echo ""
echo "Build complete. Artifacts in $OUTDIR/:"
ls -lh "$OUTDIR/"
