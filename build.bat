@echo off
setlocal EnableDelayedExpansion

REM ERP System Build Script for Windows
REM Supports building frontend and backend with cross-compilation options

set "BUILD_FRONTEND=true"
set "BUILD_BACKEND=true"
set "BACKEND_TARGETS=linux/amd64"
set "CLEAN=false"

REM Parse arguments
:parse_args
if "%~1"=="" goto end_parse
if /i "%~1"=="-f" (
    set "BUILD_FRONTEND=true"
    set "BUILD_BACKEND=false"
    shift
    goto parse_args
)
if /i "%~1"=="--frontend-only" (
    set "BUILD_FRONTEND=true"
    set "BUILD_BACKEND=false"
    shift
    goto parse_args
)
if /i "%~1"=="-b" (
    set "BUILD_FRONTEND=false"
    set "BUILD_BACKEND=true"
    shift
    goto parse_args
)
if /i "%~1"=="--backend-only" (
    set "BUILD_FRONTEND=false"
    set "BUILD_BACKEND=true"
    shift
    goto parse_args
)
if /i "%~1"=="--targets" (
    set "BACKEND_TARGETS=%~2"
    shift
    shift
    goto parse_args
)
if /i "%~1"=="--all" (
    set "BACKEND_TARGETS=linux/amd64 linux/arm64 linux/armv7 windows/amd64"
    shift
    goto parse_args
)
if /i "%~1"=="-c" (
    set "CLEAN=true"
    shift
    goto parse_args
)
if /i "%~1"=="--clean" (
    set "CLEAN=true"
    shift
    goto parse_args
)
if /i "%~1"=="-h" goto show_help
if /i "%~1"=="--help" goto show_help
echo [ERROR] Unknown option: %~1
echo Use -h or --help for usage information
exit /b 1

:show_help
echo Usage: %~nx0 [OPTIONS]
echo.
echo Options:
echo   -f, --frontend-only    Build only frontend
echo   -b, --backend-only     Build only backend
echo   --targets TARGETS      Comma-separated build targets (default: linux/amd64)
echo                          Available targets:
echo                            linux/amd64, linux/arm64, linux/armv7, windows/amd64
echo   --all                   Build for all supported platforms
echo   -c, --clean             Clean build artifacts before building
echo   -h, --help              Show this help message
echo.
echo Examples:
echo   %~nx0                   # Build frontend + linux/amd64 backend
echo   %~nx0 -f                # Build only frontend
echo   %~nx0 -b --all         # Build backend for all platforms
echo   %~nx0 --targets linux/arm64,windows/amd64
exit /b 0

:end_parse

echo ========================================
echo   ERP System Build Script
echo ========================================
echo.

REM Clean build artifacts if requested
if "%CLEAN%"=="true" (
    echo Cleaning build artifacts...
    if exist "embedded\dist" rmdir /s /q "embedded\dist"
    if exist "erp.exe" del /q "erp.exe"
    if exist "erp-*" del /q "erp-*"
    echo Clean completed.
    echo.
)

REM Build frontend
if "%BUILD_FRONTEND%"=="true" (
    echo ========================================
    echo Building Frontend...
    echo ========================================

    if not exist "frontend\node_modules" (
        echo Installing frontend dependencies...
        cd frontend
        call npm install
        cd ..
    )

    echo Running frontend build...
    cd frontend
    call npm run build
    cd ..

    echo Frontend build completed: embedded\dist\
    echo.
)

REM Build backend
if "%BUILD_BACKEND%"=="true" (
    echo ========================================
    echo Building Backend...
    echo ========================================

    REM Check if frontend is built
    if not exist "embedded\dist" (
        echo [ERROR] Frontend not built. Please build frontend first or use -f flag.
        exit /b 1
    )

    REM Parse targets
    for %%T in (%BACKEND_TARGETS%) do (
        set "TARGET=%%T"
        for /f "tokens=1,2 delims=/" %%a in ("%%T") do (
            set "GOOS=%%a"
            set "GOARCH=%%b"
        )

        REM Determine output filename
        if "!GOOS!"=="windows" (
            set "OUTPUT_NAME=erp-!GOOS!-!GOARCH!.exe"
        ) else (
            set "OUTPUT_NAME=erp-!GOOS!-!GOARCH!"
        )

        echo Building for !GOOS!/!GOARCH!...

        REM Set ARM version for armv7
        if "!GOARCH!"=="armv7" (
            set "GOARCH=arm"
            set "GOARM=7"
        ) else (
            set "GOARM="
        )

        REM Set build environment and compile
        set "CGO_ENABLED=0"
        set "GOOS=!GOOS!"
        set "GOARCH=!GOARCH!"
        if defined GOARM set "GOARM=!GOARM!"

        go build -o "!OUTPUT_NAME!" -ldflags="-s -w" main.go

        if !errorlevel! equ 0 (
            echo Built: !OUTPUT_NAME!
        ) else (
            echo [ERROR] Failed to build for !GOOS!/!GOARCH!
            exit /b 1
        )
    )

    echo.
    echo Backend build completed!
    echo.
)

REM Summary
echo ========================================
echo Build Summary
echo ========================================

if "%BUILD_FRONTEND%"=="true" (
    echo [OK] Frontend: embedded\dist\
)

if "%BUILD_BACKEND%"=="true" (
    echo [OK] Backend binaries:
    dir /b erp-* 2>nul
    dir /b erp.exe 2>nul
)

echo.
echo ========================================
echo All builds completed successfully!
echo ========================================
endlocal
