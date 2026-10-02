[CmdletBinding()]
param(
    [string]$Version = ""
)

$ErrorActionPreference = "Stop"
$repository = "duck-ahiru-Z/DevDuck"
$installDir = Join-Path $env:LOCALAPPDATA "DevDuck\bin"
$installPath = Join-Path $installDir "duck.exe"

$architecture = [Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITEW6432")
if ([string]::IsNullOrWhiteSpace($architecture)) {
    $architecture = [Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")
}
switch ($architecture.ToUpperInvariant()) {
    "AMD64" { $asset = "devduck-windows-amd64.exe" }
    "ARM64" { $asset = "devduck-windows-arm64.exe" }
    default { throw "Unsupported Windows architecture: $architecture" }
}

$releaseBase = if ([string]::IsNullOrWhiteSpace($Version)) {
    "https://github.com/$repository/releases/latest/download"
} else {
    "https://github.com/$repository/releases/download/$Version"
}
$binaryUrl = "$releaseBase/$asset"
$checksumUrl = "$releaseBase/SHA256SUMS"
$tempDir = Join-Path ([IO.Path]::GetTempPath()) ("devduck-install-" + [Guid]::NewGuid())
$tempBinary = Join-Path $tempDir $asset
$tempChecksums = Join-Path $tempDir "SHA256SUMS"

try {
    New-Item -ItemType Directory -Path $tempDir | Out-Null
    Invoke-WebRequest -Uri $binaryUrl -OutFile $tempBinary
    Invoke-WebRequest -Uri $checksumUrl -OutFile $tempChecksums
    $expected = (Get-Content $tempChecksums | Where-Object { $_ -match [regex]::Escape($asset) } | Select-Object -First 1) -split '\s+' | Select-Object -First 1
    if ([string]::IsNullOrWhiteSpace($expected)) { throw "No checksum found for $asset" }
    $actual = (Get-FileHash -Algorithm SHA256 -Path $tempBinary).Hash.ToLowerInvariant()
    if ($actual -ne $expected.ToLowerInvariant()) { throw "SHA256 mismatch for $asset" }
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
    Move-Item -Force $tempBinary $installPath
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $pathEntries = @($userPath -split ';' | Where-Object { $_ })
    if ($pathEntries -notcontains $installDir) {
        $newPath = (($pathEntries + $installDir) -join ';')
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        Write-Host "Added $installDir to your user PATH. Open a new PowerShell window to use duck."
    }
    Write-Host "Installed DevDuck at $installPath"
} finally {
    Remove-Item -Recurse -Force $tempDir -ErrorAction SilentlyContinue
}
