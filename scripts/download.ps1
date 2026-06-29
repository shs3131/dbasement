param(
    [string]$Version = "latest",
    [string]$OutDir = "."
)

$repo = "shs3131/dbasement"
if ($Version -eq "latest") {
    $tag = (Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest").tag_name
} else {
    $tag = $Version
}

$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$url = "https://github.com/$repo/releases/download/$tag/dbasement-windows-$arch.zip"

Write-Output "Downloading Dbasement $tag for Windows..."
Invoke-WebRequest -Uri $url -OutFile "$OutDir\dbasement-windows-$arch.zip"
Expand-Archive -Path "$OutDir\dbasement-windows-$arch.zip" -DestinationPath $OutDir -Force
Remove-Item "$OutDir\dbasement-windows-$arch.zip"
Write-Output "Installed at: $OutDir\dbasement.exe"
