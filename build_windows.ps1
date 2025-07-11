# Windows Build Script for OrgChart (PowerShell)
# Builds the orgchart program for Windows

param(
    [switch]$Help
)

if ($Help) {
    Write-Host "Windows Build Script for OrgChart"
    Write-Host "Usage: .\build_windows.ps1"
    Write-Host "This script builds the orgchart program for Windows"
    exit 0
}

Write-Host "=== Windows Build Script for OrgChart ===" -ForegroundColor Green
Write-Host "Building for Windows..." -ForegroundColor Yellow
Write-Host ""

# Check if Go is available
try {
    $goVersion = go version
    Write-Host "✅ Go version: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Error: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go 1.24+ from https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Check if orgchart.go exists
if (-not (Test-Path "orgchart.go")) {
    Write-Host "❌ Error: orgchart.go not found in current directory" -ForegroundColor Red
    Write-Host "Please run this script from the directory containing orgchart.go" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "=== Building Windows Executables ===" -ForegroundColor Green

# Create build directory
$BUILD_DIR = "build\windows"
if (-not (Test-Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR -Force | Out-Null
}

# Build for Windows 64-bit (most common)
Write-Host "🔨 Building for Windows 64-bit..." -ForegroundColor Yellow
try {
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    go build -o "$BUILD_DIR\orgchart_win64.exe" orgchart_dragdrop.go
    Write-Host "✅ Windows 64-bit: $BUILD_DIR\orgchart_win64.exe" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to build Windows 64-bit" -ForegroundColor Red
    exit 1
}

# Build for Windows 32-bit (for older systems)
Write-Host "🔨 Building for Windows 32-bit..." -ForegroundColor Yellow
try {
    $env:GOOS = "windows"
    $env:GOARCH = "386"
    go build -o "$BUILD_DIR\orgchart_win32.exe" orgchart_dragdrop.go
    Write-Host "✅ Windows 32-bit: $BUILD_DIR\orgchart_win32.exe" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to build Windows 32-bit" -ForegroundColor Red
    exit 1
}

# Build for Windows ARM64 (for newer Windows on ARM)
Write-Host "🔨 Building for Windows ARM64..." -ForegroundColor Yellow
try {
    $env:GOOS = "windows"
    $env:GOARCH = "arm64"
    go build -o "$BUILD_DIR\orgchart_win_arm64.exe" orgchart_dragdrop.go
    Write-Host "✅ Windows ARM64: $BUILD_DIR\orgchart_win_arm64.exe" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to build Windows ARM64" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=== Build Summary ===" -ForegroundColor Green
Write-Host "📁 Build directory: $BUILD_DIR" -ForegroundColor Cyan
Write-Host "📦 Executables created:" -ForegroundColor Cyan
Get-ChildItem "$BUILD_DIR\*.exe" | ForEach-Object {
    Write-Host "   $($_.Name) ($([math]::Round($_.Length/1MB, 2)) MB)" -ForegroundColor White
}

Write-Host ""
Write-Host "=== Installation Instructions ===" -ForegroundColor Green
Write-Host "1. Install Graphviz on Windows:" -ForegroundColor Yellow
Write-Host "   - Using Chocolatey: choco install graphviz" -ForegroundColor White
Write-Host "   - Or download from: https://graphviz.org/download/" -ForegroundColor White
Write-Host "2. Run the program: .\orgchart_win64.exe fcbh.csv" -ForegroundColor Yellow
Write-Host ""
Write-Host "=== Usage Examples ===" -ForegroundColor Green
Write-Host ".\orgchart_win64.exe fcbh.csv" -ForegroundColor White
Write-Host ".\orgchart_win64.exe fcbh.csv my_chart.pdf" -ForegroundColor White
Write-Host ".\orgchart_win64.exe fcbh.csv --keep-dot" -ForegroundColor White
Write-Host ""
Write-Host "✅ Windows build completed successfully!" -ForegroundColor Green 