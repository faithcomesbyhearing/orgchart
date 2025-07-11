# CSV File Association Script for OrgChart (PowerShell)
# This script associates .csv files with the orgchart_dragdrop.exe

param(
    [switch]$Help
)

if ($Help) {
    Write-Host "CSV File Association for OrgChart"
    Write-Host "Usage: .\associate_csv.ps1"
    Write-Host "This script associates .csv files with orgchart_dragdrop.exe"
    exit 0
}

Write-Host "=== CSV File Association for OrgChart ===" -ForegroundColor Green
Write-Host ""

# Check if running as administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")
if (-not $isAdmin) {
    Write-Host "❌ Error: This script must be run as Administrator" -ForegroundColor Red
    Write-Host "Right-click on this file and select 'Run as administrator'" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

# Check if orgchart_dragdrop.exe exists
if (-not (Test-Path "orgchart_dragdrop.exe")) {
    Write-Host "❌ Error: orgchart_dragdrop.exe not found" -ForegroundColor Red
    Write-Host "Please build the Windows executable first using build_windows.bat" -ForegroundColor Yellow
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "🔧 Associating .csv files with orgchart_dragdrop.exe..." -ForegroundColor Yellow
Write-Host ""

# Get the full path to the executable
$exePath = (Get-Item "orgchart_dragdrop.exe").FullName

try {
    # Create file association
    cmd /c "assoc .csv=OrgChart.CSV"
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to associate .csv file type"
    }

    # Create file type description
    cmd /c "ftype OrgChart.CSV=`"$exePath`" `"%1`""
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to set file type handler"
    }

    Write-Host "✅ Successfully associated .csv files with orgchart_dragdrop.exe" -ForegroundColor Green
    Write-Host ""
    Write-Host "=== Usage Instructions ===" -ForegroundColor Cyan
    Write-Host "1. Double-click any .csv file to generate an org chart" -ForegroundColor White
    Write-Host "2. The PDF will be created in the same folder as the CSV" -ForegroundColor White
    Write-Host "3. To keep the DOT file, run from command line with --keep-dot" -ForegroundColor White
    Write-Host ""
    Write-Host "=== Command Line Usage ===" -ForegroundColor Cyan
    Write-Host ".\orgchart_dragdrop.exe fcbh.csv" -ForegroundColor White
    Write-Host ".\orgchart_dragdrop.exe fcbh.csv my_chart.pdf --keep-dot" -ForegroundColor White
    Write-Host ""
    Write-Host "=== Requirements ===" -ForegroundColor Cyan
    Write-Host "Make sure Graphviz is installed: choco install graphviz" -ForegroundColor White
    Write-Host ""
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

Read-Host "Press Enter to exit" 