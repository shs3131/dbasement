param(
  [string]$Version = "dev"
)

$OutDir = "dist"
Write-Output "Building Dbasement $Version"
New-Item -ItemType Directory -Path $OutDir -Force | Out-Null

$Targets = @(
  @{ GOOS="linux";  GOARCH="amd64"; Output="dbasement-linux-amd64" }
  @{ GOOS="linux";  GOARCH="arm64"; Output="dbasement-linux-arm64" }
  @{ GOOS="darwin"; GOARCH="amd64"; Output="dbasement-darwin-amd64" }
  @{ GOOS="darwin"; GOARCH="arm64"; Output="dbasement-darwin-arm64" }
  @{ GOOS="windows"; GOARCH="amd64"; Output="dbasement-windows-amd64.exe" }
)

foreach ($t in $Targets) {
  $env:GOOS = $t.GOOS
  $env:GOARCH = $t.GOARCH
  $output = Join-Path $OutDir $t.Output

  Write-Output "Building $($t.GOOS)/$($t.GOARCH) -> $output"

  go build -ldflags="-s -w -X main.version=$Version" -o $output ./cmd/dbasement/
}

Write-Output ""
Write-Output "Build complete. Artifacts in $OutDir/:"
Get-ChildItem -Path $OutDir | Select-Object Name, Length | Format-Table -AutoSize
