@echo off
REM Windows Build Script for OrgChart (Batch)
REM Builds the orgchart program for Windows

echo === Windows Build Script for OrgChart ===
echo Building for Windows...
echo.

REM Check if Go is available
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Error: Go is not installed or not in PATH
    echo Please install Go 1.24+ from https://golang.org/dl/
    pause
    exit /b 1
)

REM Check if orgchart_dragdrop.go exists
if not exist "orgchart_dragdrop.go" (
    echo ❌ Error: orgchart_dragdrop.go not found in current directory
    echo Please run this script from the directory containing orgchart_dragdrop.go
    pause
    exit /b 1
)

echo ✅ Go version: 
go version
echo.

echo === Building Windows Executables ===

REM Create build directory
if not exist "build\windows" mkdir "build\windows"

REM Build for Windows 64-bit (most common)
echo 🔨 Building for Windows 64-bit...
set GOOS=windows
set GOARCH=amd64
go build -o "build\windows\orgchart_win64.exe" orgchart_dragdrop.go
if %errorlevel% neq 0 (
    echo ❌ Failed to build Windows 64-bit
    pause
    exit /b 1
)
echo ✅ Windows 64-bit: build\windows\orgchart_win64.exe

REM Build for Windows 32-bit (for older systems)
echo 🔨 Building for Windows 32-bit...
set GOOS=windows
set GOARCH=386
go build -o "build\windows\orgchart_win32.exe" orgchart_dragdrop.go
if %errorlevel% neq 0 (
    echo ❌ Failed to build Windows 32-bit
    pause
    exit /b 1
)
echo ✅ Windows 32-bit: build\windows\orgchart_win32.exe

REM Build for Windows ARM64 (for newer Windows on ARM)
echo 🔨 Building for Windows ARM64...
set GOOS=windows
set GOARCH=arm64
go build -o "build\windows\orgchart_win_arm64.exe" orgchart_dragdrop.go
if %errorlevel% neq 0 (
    echo ❌ Failed to build Windows ARM64
    pause
    exit /b 1
)
echo ✅ Windows ARM64: build\windows\orgchart_win_arm64.exe

echo.
echo === Build Summary ===
echo 📁 Build directory: build\windows
echo 📦 Executables created:
dir "build\windows\*.exe"

echo.
echo === Installation Instructions ===
echo 1. Install Graphviz on Windows:
echo    - Using Chocolatey: choco install graphviz
echo    - Or download from: https://graphviz.org/download/
echo 2. Run the program: orgchart_win64.exe fcbh.csv
echo.
echo === Usage Examples ===
echo orgchart_win64.exe fcbh.csv
echo orgchart_win64.exe fcbh.csv my_chart.pdf
echo orgchart_win64.exe fcbh.csv --keep-dot
echo.
echo ✅ Windows build completed successfully!
pause 