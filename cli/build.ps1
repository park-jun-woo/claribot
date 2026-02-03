# Claritask CLI Build Script for Windows (PowerShell)
# Usage: .\build.ps1 [command]
# Commands: build, install, clean, all (default: build)

param(
    [Parameter(Position=0)]
    [ValidateSet("build", "install", "clean", "all")]
    [string]$Command = "build"
)

$ErrorActionPreference = "Stop"
$BinName = "clari.exe"
$InstallDir = "$env:USERPROFILE\bin"

function Test-Prerequisites {
    # Check Go
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        Write-Error "Go is not installed. Install from https://go.dev/dl/"
        exit 1
    }

    Write-Host "Prerequisites OK" -ForegroundColor Green
}

function Invoke-Build {
    Write-Host "Building $BinName..." -ForegroundColor Cyan

    go build -o $BinName ./cmd/claritask/

    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build successful: $BinName" -ForegroundColor Green
    } else {
        Write-Error "Build failed"
        exit 1
    }
}

function Invoke-Install {
    if (-not (Test-Path $BinName)) {
        Write-Host "Binary not found. Building first..." -ForegroundColor Yellow
        Invoke-Build
    }

    Write-Host "Installing to $InstallDir..." -ForegroundColor Cyan

    # Create install directory if not exists
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    Copy-Item $BinName "$InstallDir\$BinName" -Force

    # Check if InstallDir is in PATH
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$InstallDir*") {
        Write-Host "Adding $InstallDir to PATH..." -ForegroundColor Yellow
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
        Write-Host "PATH updated. Restart PowerShell to use 'clari' command." -ForegroundColor Yellow
    }

    Write-Host "Installed: $InstallDir\$BinName" -ForegroundColor Green
}

function Invoke-Clean {
    Write-Host "Cleaning..." -ForegroundColor Cyan

    if (Test-Path $BinName) {
        Remove-Item $BinName -Force
        Write-Host "Removed $BinName" -ForegroundColor Green
    } else {
        Write-Host "Nothing to clean" -ForegroundColor Yellow
    }
}

# Main
Push-Location $PSScriptRoot
try {
    switch ($Command) {
        "build" {
            Test-Prerequisites
            Invoke-Build
        }
        "install" {
            Test-Prerequisites
            Invoke-Install
        }
        "clean" {
            Invoke-Clean
        }
        "all" {
            Test-Prerequisites
            Invoke-Build
            Invoke-Install
        }
    }
} finally {
    Pop-Location
}
