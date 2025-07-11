@echo off
REM CSV File Association Script for OrgChart
REM This script associates .csv files with the orgchart_dragdrop.exe

echo === CSV File Association for OrgChart ===
echo.

REM Check if running as administrator
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo ❌ Error: This script must be run as Administrator
    echo Right-click on this file and select "Run as administrator"
    pause
    exit /b 1
)

REM Check if orgchart_dragdrop.exe exists
if not exist "orgchart_dragdrop.exe" (
    echo ❌ Error: orgchart_dragdrop.exe not found
    echo Please build the Windows executable first using build_windows.bat
    pause
    exit /b 1
)

echo 🔧 Associating .csv files with orgchart_dragdrop.exe...
echo.

REM Get the full path to the executable
for %%i in ("orgchart_dragdrop.exe") do set "EXE_PATH=%%~fi"

REM Create file association
assoc .csv=OrgChart.CSV
if %errorLevel% neq 0 (
    echo ❌ Failed to associate .csv file type
    pause
    exit /b 1
)

REM Create file type description
ftype OrgChart.CSV="%EXE_PATH%" "%%1"
if %errorLevel% neq 0 (
    echo ❌ Failed to set file type handler
    pause
    exit /b 1
)

echo ✅ Successfully associated .csv files with orgchart_dragdrop.exe
echo.
echo === Usage Instructions ===
echo 1. Double-click any .csv file to generate an org chart
echo 2. The PDF will be created in the same folder as the CSV
echo 3. To keep the DOT file, run from command line with --keep-dot
echo.
echo === Command Line Usage ===
echo orgchart_dragdrop.exe fcbh.csv
echo orgchart_dragdrop.exe fcbh.csv my_chart.pdf --keep-dot
echo.
echo === Requirements ===
echo Make sure Graphviz is installed: choco install graphviz
echo.
pause 