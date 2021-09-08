<#
    .SYNOPSIS
        This script contains common functions for building the Windows New Relic Infrastructure Agent.
#>
Function SignExecutable {
    param (
        # Signing tool
        [string]$signtool='"C:\Program Files (x86)\Windows Kits\10\bin\x64\signtool.exe"',
        [string]$executable=$(throw "-executable path is required")
    )

    Invoke-Expression "& $signtool sign /d 'New Relic Infrastructure Agent' /n 'New Relic, Inc.' $executable"
    if ($lastExitCode -ne 0) {
       throw "Failed to sign $executable"
    }
}

Function DownloadFile {
    param (
        [string]$url=$(throw "-url is required"),
        # $dest is that destination path where the file will be downloaded.
        [string]$dest=".\",
        # Pass outFile if you want to rename the outputFile. By default will use the file name from the url.
        [string]$outFile=""
    )

    if ([string]::IsNullOrWhitespace($outFile)) {
        $outFile = $url.Substring($url.LastIndexOf("/") + 1)
    }

    # Download zip file.
    $ProgressPreference = 'SilentlyContinue'
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

    Write-Output "Downloading $url"

    New-Item -path $dest -type directory -Force
    Invoke-WebRequest $url -OutFile "$dest\$outFile"
}

Function DownloadAndExtractZip {
    param (
        [string]$url=$(throw "-url is required"),
        [string]$dest=$(throw "-dest is required")
    )

    DownloadFile -dest:"$dest" -url:"$url"

    $file = $url.Substring($url.LastIndexOf("/") + 1)

    # extract
    expand-archive -path "$dest\$file" -destinationpath $dest
    Remove-Item "$dest\$file"
}