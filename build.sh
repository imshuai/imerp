#!/bin/bash

# ERP System Build Script
# Supports building frontend and backend with cross-compilation options

set -e

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
BUILD_FRONTEND=true
BUILD_BACKEND=true
BACKEND_TARGETS=("linux/amd64")
FRONTEND_ONLY=false
BACKEND_ONLY=false
CLEAN=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--frontend-only)
            BUILD_FRONTEND=true
            BUILD_BACKEND=false
            FRONTEND_ONLY=true
            shift
            ;;
        -b|--backend-only)
            BUILD_FRONTEND=false
            BUILD_BACKEND=true
            BACKEND_ONLY=true
            shift
            ;;
        --targets)
            IFS=',' read -ra BACKEND_TARGETS <<< "$2"
            shift 2
            ;;
        --all)
            BACKEND_TARGETS=(
                "linux/amd64"
                "linux/arm64"
                "linux/armv7"
                "windows/amd64"
            )
            shift
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -f, --frontend-only    Build only frontend"
            echo "  -b, --backend-only     Build only backend"
            echo "  --targets TARGETS      Comma-separated build targets (default: linux/amd64)"
            echo "                          Available targets:"
            echo "                            linux/amd64, linux/arm64, linux/armv7, windows/amd64"
            echo "  --all                   Build for all supported platforms"
            echo "  -c, --clean             Clean build artifacts before building"
            echo "  -h, --help              Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                      # Build frontend + linux/amd64 backend"
            echo "  $0 -f                   # Build only frontend"
            echo "  $0 -b --all            # Build backend for all platforms"
            echo "  $0 --targets linux/arm64,windows/amd64"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  ERP System Build Script${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Clean build artifacts if requested
if [ "$CLEAN" = true ]; then
    echo -e "${YELLOW}Cleaning build artifacts...${NC}"
    rm -rf frontend/dist embedded/dist
    rm -f erp erp.exe erp-*
    echo -e "${GREEN}Clean completed.${NC}"
    echo ""
fi

# Build frontend
if [ "$BUILD_FRONTEND" = true ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Building Frontend...${NC}"
    echo -e "${GREEN}========================================${NC}"

    if [ ! -d "frontend/node_modules" ]; then
        echo -e "${YELLOW}Installing frontend dependencies...${NC}"
        cd frontend
        npm install
        cd ..
    fi

    echo -e "${YELLOW}Running frontend build...${NC}"
    cd frontend
    npm run build
    cd ..

    echo -e "${GREEN}Frontend build completed: embedded/dist/${NC}"
    echo ""
fi

# Build backend
if [ "$BUILD_BACKEND" = true ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Building Backend...${NC}"
    echo -e "${GREEN}========================================${NC}"

    # Check if frontend is built
    if [ ! -d "embedded/dist" ] && [ "$FRONTEND_ONLY" = false ]; then
        echo -e "${RED}Error: Frontend not built. Please build frontend first or use -f flag.${NC}"
        exit 1
    fi

    # Build for each target
    for TARGET in "${BACKEND_TARGETS[@]}"; do
        GOOS="${TARGET%/*}"
        GOARCH="${TARGET#*/}"

        # Set ARM version for armv7
        if [ "$GOARCH" = "armv7" ]; then
            GOARCH="arm"
            GOARM="7"
        else
            GOARM=""
        fi

        # Determine output filename
        if [ "$GOOS" = "windows" ]; then
            OUTPUT_NAME="erp-${GOOS}-${GOARCH}.exe"
        else
            OUTPUT_NAME="erp-${GOOS}-${GOARCH}"
        fi

        echo -e "${YELLOW}Building for ${GOOS}/${GOARCH}...${NC}"

        # Set build environment
        export GOOS="$GOOS"
        export GOARCH="$GOARCH"
        [ -n "$GOARM" ] && export GOARM="$GOARM"
        export CGO_ENABLED=0

        # Build
        go build -o "$OUTPUT_NAME" -ldflags="-s -w" main.go

        if [ $? -eq 0 ]; then
            SIZE=$(du -h "$OUTPUT_NAME" | cut -f1)
            echo -e "${GREEN}✓ Built: $OUTPUT_NAME ($SIZE)${NC}"
        else
            echo -e "${RED}✗ Failed to build for ${GOOS}/${GOARCH}${NC}"
            exit 1
        fi
    done

    echo ""
    echo -e "${GREEN}Backend build completed!${NC}"
    echo ""
fi

# Summary
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Build Summary${NC}"
echo -e "${GREEN}========================================${NC}"

if [ "$BUILD_FRONTEND" = true ]; then
    echo -e "${GREEN}✓ Frontend: embedded/dist/${NC}"
fi

if [ "$BUILD_BACKEND" = true ]; then
    echo -e "${GREEN}✓ Backend binaries:${NC}"
    ls -lh erp-* 2>/dev/null || ls -lh erp 2>/dev/null || ls -lh erp.exe 2>/dev/null
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}All builds completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
