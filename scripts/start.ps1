param(
    [string]$ProjectPath = "."
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $PSCommandPath
$RootDir = Split-Path -Parent $ScriptDir
$ExePath = Join-Path $RootDir "dbasement.exe"

if (-not (Test-Path $ExePath)) {
    Write-Host "Downloading Dbasement..."
    $repo = "shs3131/dbasement"
    $tag = (Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest").tag_name
    $url = "https://github.com/$repo/releases/download/$tag/dbasement-windows-amd64.zip"
    $zip = Join-Path $env:TEMP "dbasement.zip"
    Invoke-WebRequest -Uri $url -OutFile $zip
    Expand-Archive -Path $zip -DestinationPath $RootDir -Force
    Remove-Item $zip
    Write-Host "Downloaded: $ExePath"
}

& $ExePath --project $ProjectPath
