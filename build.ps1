# ERP System Build Script for PowerShell
# Supports building frontend and backend with cross-compilation options

param(
    [switch]$FrontendOnly,
    [switch]$BackendOnly,
    [string[]]$Targets = @("linux/amd64"),
    [switch]$All,
    [switch]$Clean,
    [switch]$Help
)

# Show help
if ($Help) {
    Write-Host "Usage: .\build.ps1 [OPTIONS]" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -FrontendOnly           Build only frontend"
    Write-Host "  -BackendOnly            Build only backend"
    Write-Host "  -Targets TARGETS        Array of build targets (default: linux/amd64)"
    Write-Host "                          Available targets:"
    Write-Host "                            linux/amd64, linux/arm64, linux/armv7, windows/amd64"
    Write-Host "  -All                    Build for all supported platforms"
    Write-Host "  -Clean                  Clean build artifacts before building"
    Write-Host "  -Help                   Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\build.ps1                      # Build frontend + linux/amd64 backend"
    Write-Host "  .\build.ps1 -FrontendOnly       # Build only frontend"
    Write-Host "  .\build.ps1 -BackendOnly -All   # Build backend for all platforms"
    Write-Host "  .\build.ps1 -Targets linux/arm64,windows/amd64"
    exit 0
}

# Set build flags
$BUILD_FRONTEND = -not $BackendOnly
$BUILD_BACKEND = -not $FrontendOnly

if ($All) {
    $Targets = @("linux/amd64", "linux/arm64", "linux/armv7", "windows/amd64")
}

Write-Host "========================================" -ForegroundColor Green
Write-Host "  ERP System Build Script" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""

# Get script directory
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $SCRIPT_DIR

# Clean build artifacts if requested
if ($Clean) {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Yellow
    if (Test-Path "frontend\dist") {
        Remove-Item -Recurse -Force "frontend\dist"
    }
    if (Test-Path "erp") {
        Remove-Item -Force "erp"
    }
    if (Test-Path "erp.exe") {
        Remove-Item -Force "erp.exe"
    }
    Get-ChildItem -Filter "erp-*" | Remove-Item -Force
    Write-Host "Clean completed." -ForegroundColor Green
    Write-Host ""
}

# Build frontend
if ($BUILD_FRONTEND) {
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "Building Frontend..." -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green

    if (-not (Test-Path "frontend\node_modules")) {
        Write-Host "Installing frontend dependencies..." -ForegroundColor Yellow
        Set-Location frontend
        npm install
        Set-Location ..
    }

    Write-Host "Running frontend build..." -ForegroundColor Yellow
    Set-Location frontend
    npm run build
    Set-Location ..

    Write-Host "Frontend build completed: frontend\dist\" -ForegroundColor Green
    Write-Host ""
}

# Build backend
if ($BUILD_BACKEND) {
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "Building Backend..." -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green

    # Check if frontend is built
    if (-not (Test-Path "frontend\dist") -and -not $FrontendOnly) {
        Write-Host "[ERROR] Frontend not built. Please build frontend first or use -FrontendOnly." -ForegroundColor Red
        exit 1
    }

    # Set environment variables
    $env:CGO_ENABLED = "0"

    # Build for each target
    foreach ($TARGET in $Targets) {
        $PARTS = $TARGET -split '/'
        $GOOS = $PARTS[0]
        $GOARCH = $PARTS[1]
        $GOARM = ""

        # Set ARM version for armv7
        if ($GOARCH -eq "armv7") {
            $GOARCH = "arm"
            $GOARM = "7"
        }

        # Determine output filename
        if ($GOOS -eq "windows") {
            $OUTPUT_NAME = "erp-${GOOS}-${GOARCH}.exe"
        } else {
            $OUTPUT_NAME = "erp-${GOOS}-${GOARCH}"
        }

        Write-Host "Building for ${GOOS}/${GOARCH}..." -ForegroundColor Yellow

        # Set build environment
        $env:GOOS = $GOOS
        $env:GOARCH = $GOARCH
        if ($GOARM -ne "") {
            $env:GOARM = $GOARM
        } else {
            Remove-Item Env:GOARM -ErrorAction SilentlyContinue
        }

        # Build
        go build -o "$OUTPUT_NAME" -ldflags="-s -w" main.go

        if ($LASTEXITCODE -eq 0) {
            $FILE = Get-Item "$OUTPUT_NAME"
            $SIZE = [math]::Round($FILE.Length / 1MB, 2)
            Write-Host "[OK] Built: $OUTPUT_NAME ($SIZE MB)" -ForegroundColor Green
        } else {
            Write-Host "[ERROR] Failed to build for ${GOOS}/${GOARCH}" -ForegroundColor Red
            exit 1
        }
    }

    Write-Host ""
    Write-Host "Backend build completed!" -ForegroundColor Green
    Write-Host ""
}

# Summary
Write-Host "========================================" -ForegroundColor Green
Write-Host "Build Summary" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

if ($BUILD_FRONTEND) {
    Write-Host "[OK] Frontend: frontend\dist\" -ForegroundColor Green
}

if ($BUILD_BACKEND) {
    Write-Host "[OK] Backend binaries:" -ForegroundColor Green
    Get-ChildItem -Filter "erp-*" -ErrorAction SilentlyContinue | ForEach-Object {
        $SIZE = [math]::Round($_.Length / 1MB, 2)
        Write-Host "  $($_.Name) ($SIZE MB)"
    }
    if (Test-Path "erp.exe") {
        $FILE = Get-Item "erp.exe"
        $SIZE = [math]::Round($FILE.Length / 1MB, 2)
        Write-Host "  erp.exe ($SIZE MB)"
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "All builds completed successfully!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
