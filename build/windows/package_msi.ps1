<#
    .SYNOPSIS
        This script packages NRJMX
#>
param (
    # Target architecture: amd64 (default) or 386
    [ValidateSet("amd64", "386")]
    [string]$arch="amd64",
    [string]$tag="v0.0.0",
    [string]$scriptPath=$(split-path -parent $MyInvocation.MyCommand.Definition)
)

# Trim v from tag.
$version = $(echo $tag | %{if($_ -match "^v") { $_.Substring(1); }})

# Source build Functions.
. $scriptPath/functions.ps1

Function DownloadNrjmx {
    Write-Output "--- Downloading nrjmx"

    # download
    [string]$file="nrjmx_windows_0.0.0_noarch.zip"
    $url="https://github.com/newrelic/nrjmx/releases/download/${tag}/${file}"
 
    DownloadAndExtractZip -dest:"$downloadPath\nrjmx" -url:"$url"

    Copy-Item -Path "$downloadPath\nrjmx\Program Files\New Relic\nrjmx\bin\*" -Destination "$downloadPath\nrjmx\" -Recurse -Force
    Remove-Item -Path "$downloadPath\nrjmx\Program Files" -Force -Recurse
}
# Call all the steps.
$downloadPath = "$scriptPath\..\..\target\"

Write-Output "--- Cleaning..."

Remove-Item $downloadPath -Recurse -ErrorAction Ignore
New-Item -ItemType Directory -Force -Path "$downloadPath"

DownloadNrjmx

echo "Checking MSBuild.exe..."
$msBuild = (Get-ItemProperty hklm:\software\Microsoft\MSBuild\ToolsVersions\4.0).MSBuildToolsPath
if ($msBuild.Length -eq 0) {
    echo "Can't find MSBuild tool. .NET Framework 4.0.x must be installed"
    exit -1
}
echo $msBuild

echo "--- Building Installer"

Push-Location -Path "$scriptPath\pkg\windows\"
$env:NRJMX_VERSION = $version
. $msBuild/MSBuild.exe nrjmx-installer.wixproj

if (-not $?)
{
    echo "Failed building installer"
    Pop-Location
    exit -1
}

echo "Making versioned installed copy"

cd target\msi\Release

cp "nrjmx.msi" "nrjmx-$arch.$version.msi"

Pop-Location
