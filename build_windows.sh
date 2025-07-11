#!/bin/bash

# Windows Build Script for OrgChart
# Cross-compiles the orgchart program for Windows from macOS/Linux

set -e  # Exit on any error

echo "=== Windows Build Script for OrgChart ==="
echo "Building from $(uname -s) to Windows..."
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    echo "Please install Go 1.24+ from https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Check if orgchart_dragdrop.go exists
if [ ! -f "orgchart_dragdrop.go" ]; then
	echo "❌ Error: orgchart_dragdrop.go not found in current directory"
	echo "Please run this script from the directory containing orgchart_dragdrop.go"
	exit 1
fi

echo ""
echo "=== Building Windows Executables ==="

# Create build directory
BUILD_DIR="build/windows"
mkdir -p "$BUILD_DIR"

# Build for Windows 64-bit (most common)
echo "🔨 Building for Windows 64-bit..."
GOOS=windows GOARCH=amd64 go build -o "$BUILD_DIR/orgchart_win64.exe" orgchart_dragdrop.go
if [ $? -eq 0 ]; then
    echo "✅ Windows 64-bit: $BUILD_DIR/orgchart_win64.exe"
else
    echo "❌ Failed to build Windows 64-bit"
    exit 1
fi

# Build for Windows 32-bit (for older systems)
echo "🔨 Building for Windows 32-bit..."
GOOS=windows GOARCH=386 go build -o "$BUILD_DIR/orgchart_win32.exe" orgchart_dragdrop.go
if [ $? -eq 0 ]; then
    echo "✅ Windows 32-bit: $BUILD_DIR/orgchart_win32.exe"
else
    echo "❌ Failed to build Windows 32-bit"
    exit 1
fi

# Build for Windows ARM64 (for newer Windows on ARM)
echo "🔨 Building for Windows ARM64..."
GOOS=windows GOARCH=arm64 go build -o "$BUILD_DIR/orgchart_win_arm64.exe" orgchart_dragdrop.go
if [ $? -eq 0 ]; then
    echo "✅ Windows ARM64: $BUILD_DIR/orgchart_win_arm64.exe"
else
    echo "❌ Failed to build Windows ARM64"
    exit 1
fi

echo ""
echo "=== Build Summary ==="
echo "📁 Build directory: $BUILD_DIR"
echo "📦 Executables created:"
ls -la "$BUILD_DIR"/*.exe

echo ""
echo "=== Windows Installation Instructions ==="
echo "1. Copy the appropriate .exe file to your Windows machine"
echo "2. Install Graphviz on Windows:"
echo "   - Using Chocolatey: choco install graphviz"
echo "   - Or download from: https://graphviz.org/download/"
echo "3. Run the program: orgchart_win64.exe fcbh.csv"
echo ""
echo "=== Usage Examples ==="
echo "orgchart_win64.exe fcbh.csv"
echo "orgchart_win64.exe fcbh.csv my_chart.pdf"
echo "orgchart_win64.exe fcbh.csv --keep-dot"
echo ""
echo "✅ Windows build completed successfully!" 