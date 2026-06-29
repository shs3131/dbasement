# Installation Guide

**Prerequisite:** Git must be installed (for change detection). No other dependencies.

---

## Pre-built Binary (Recommended)

Download from [GitHub Releases](https://github.com/shs3131/dbasement/releases).

### macOS (Apple Silicon)

```bash
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-arm64.tar.gz | tar xz
sudo mv dbasement /usr/local/bin/
```

### macOS (Intel)

```bash
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-amd64.tar.gz | tar xz
sudo mv dbasement /usr/local/bin/
```

### Linux (x86_64)

```bash
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-amd64.tar.gz | tar xz
sudo mv dbasement /usr/local/bin/
```

### Linux (ARM64)

```bash
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-arm64.tar.gz | tar xz
sudo mv dbasement /usr/local/bin/
```

### Windows (PowerShell)

```powershell
curl.exe -LO https://github.com/shs3131/dbasement/releases/latest/download/dbasement-windows-amd64.zip
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath . -Force
Move-Item .\dbasement.exe C:\Windows\System32\
```

### Verify

```bash
dbasement --help
```

---

## Alternative (requires Go 1.26+)

```bash
go install github.com/shs3131/dbasement/cmd/dbasement@latest
```

---

## AI Configuration

After installing, see [MCP_CONFIGS.md](MCP_CONFIGS.md) to configure your AI client.

If you are an AI agent installing Dbasement automatically, read [AI_QUICKSTART.md](AI_QUICKSTART.md) first.
