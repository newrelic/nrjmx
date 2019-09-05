<#
    .SYNOPSIS
        This script packages NRJMX
#>
param (
    # Target architecture: amd64 (default) or 386
    [ValidateSet("amd64", "386")]
    [string]$arch="amd64",
    [string]$version="0.0.0"
    # Creates a signed installer
    #[switch]$installer=$false,
    # Skip tests
    #[switch]$skipTests=$false
)

echo "Checking MSBuild.exe..."
$msBuild = (Get-ItemProperty hklm:\software\Microsoft\MSBuild\ToolsVersions\4.0).MSBuildToolsPath
if ($msBuild.Length -eq 0) {
    echo "Can't find MSBuild tool. .NET Framework 4.0.x must be installed"
    exit -1
}
echo $msBuild

echo "--- Building Installer"

Push-Location -Path "pkg\windows\"
$env:NRJMX_VERSION = $version
. $msBuild/MSBuild.exe nrjmx-installer.wixproj

if (-not $?)
{
    echo "Failed building installer"
    Pop-Location
    exit -1
}

echo "Making versioned installed copy"

cd ..\..\target\msi\Release

cp "nrjmx.msi" "nrjmx-$arch.$version.msi"

Pop-Location
