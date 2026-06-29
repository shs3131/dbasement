param(
    [string]$Destination = "."
)

$ErrorActionPreference = "Stop"
$Repo = "shs3131/dbasement"

Write-Host "Installing Dbasement to: $Destination"

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "amd64" }

# Get latest release tag
Write-Host "Detecting latest release..."
$Tag = (Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest").tag_name
Write-Host "Latest release: $Tag"

# Download and extract
$Asset = "dbasement-windows-$Arch.zip"
$Url = "https://github.com/$Repo/releases/download/$Tag/$Asset"
Write-Host "Downloading: $Url"
$Zip = Join-Path $env:TEMP "dbasement-install.zip"
Invoke-WebRequest -Uri $Url -OutFile $Zip
Expand-Archive -Path $Zip -DestinationPath $Destination -Force
Remove-Item $Zip

Write-Host ""
Write-Host "Dbasement $Tag installed at: $(Join-Path $Destination 'dbasement.exe')"
Write-Host ""
Write-Host "Next steps:"
Write-Host "  1. Configure your AI client (see MCP_CONFIGS.md)"
Write-Host "  2. Start a new AI session"
Write-Host ""
