# Installation

Dbasement is a single native binary with zero runtime dependencies. No Go, no Node, no Python required.

## Download (recommended)

Download the binary for your platform from the [latest release](https://github.com/shs3131/dbasement/releases/latest):

### Windows

```powershell
curl.exe -LO https://github.com/shs3131/dbasement/releases/latest/download/dbasement-windows-amd64.zip
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath .
```

### Linux

```bash
curl -LO https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-amd64.tar.gz
tar xzf dbasement-linux-amd64.tar.gz
chmod +x dbasement
```

### macOS

```bash
# Intel
curl -LO https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-amd64.tar.gz
tar xzf dbasement-darwin-amd64.tar.gz
chmod +x dbasement

# Apple Silicon
curl -LO https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-arm64.tar.gz
tar xzf dbasement-darwin-arm64.tar.gz
chmod +x dbasement
```

### Verify

```bash
./dbasement --help
```

## Build from Source

Requires [Go 1.26+](https://go.dev/dl).

```bash
git clone https://github.com/shs3131/dbasement.git
cd dbasement
go build -o dbasement ./cmd/dbasement/
```

## Platform Notes

- **macOS**: `xattr -d com.apple.quarantine dbasement` if blocked by Gatekeeper
- **Windows**: Windows Defender false positive is common with Go binaries; add an exclusion or build from source
- **Linux**: Ensure `git` is installed for change detection support
