# Windows Installation Guide for OrgChart

This guide will help you set up the OrgChart generator on Windows with drag-and-drop functionality.

## Quick Start

1. **Download the Windows executable** from the build output
2. **Install Graphviz** (required for PDF generation)
3. **Set up file associations** (optional, for drag-and-drop)
4. **Test the installation**

## Step-by-Step Installation

### 1. Download the Executable

Choose the appropriate version for your Windows system:

- **`orgchart_win64.exe`** - For 64-bit Windows (most common)
- **`orgchart_win32.exe`** - For 32-bit Windows (older systems)
- **`orgchart_win_arm64.exe`** - For Windows on ARM

### 2. Install Graphviz

Graphviz is required for PDF generation. Choose one method:

#### Option A: Using Chocolatey (Recommended)
```cmd
# Install Chocolatey first (if not already installed)
# Run PowerShell as Administrator and execute:
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install Graphviz
choco install graphviz
```

#### Option B: Manual Installation
1. Download Graphviz from https://graphviz.org/download/
2. Run the installer
3. Add Graphviz to your PATH (usually done automatically)

#### Option C: Using Winget
```cmd
winget install graphviz
```

### 3. Test the Installation

Open Command Prompt and test:
```cmd
dot --version
```

You should see Graphviz version information.

### 4. Set Up Drag-and-Drop (Optional)

#### Option A: Using the Association Scripts

**Using PowerShell (as Administrator):**
```powershell
.\associate_csv.ps1
```

**Using Command Prompt (as Administrator):**
```cmd
associate_csv.bat
```

#### Option B: Manual Association

1. Right-click on any `.csv` file
2. Select "Open with" → "Choose another app"
3. Browse to your `orgchart_win64.exe`
4. Check "Always use this app to open .csv files"

## Usage

### Command Line Usage

```cmd
# Basic usage
orgchart_win64.exe fcbh.csv

# Custom output name
orgchart_win64.exe fcbh.csv my_organization.pdf

# Keep intermediate DOT file
orgchart_win64.exe fcbh.csv --keep-dot
```

### Drag-and-Drop Usage

1. **Double-click any CSV file** (if file association is set up)
2. **Drag and drop CSV files** onto the executable
3. **Right-click CSV file** → "Open with" → Select the executable

## CSV File Format

Your CSV file must have these columns:

| Column | Description | Example |
|--------|-------------|---------|
| 1 | Employee Status | "Active" |
| 2 | Date Hired | "09/01/1972" |
| 3 | Job Title | "President" |
| 4 | Employee Name | "GERALD A. JACKSON" |
| 5 | Supervisor Name | "GERALD A. JACKSON" |
| 6 | Manager/Director Name | "" |
| 7 | Cost Centers | "69_Corporate" |

### Example CSV Header
```csv
"Employee Status","Date Hired","Default Jobs (HR) Full Path","Employee Name","Supervisor Name","Manager / Director Name","Cost Centers(Department)"
```

## Troubleshooting

### "dot: command not found"
- **Solution**: Install Graphviz using one of the methods above
- **Verify**: Run `dot --version` in Command Prompt

### "Error converting to PDF"
- **Solution**: Make sure Graphviz is installed and in PATH
- **Alternative**: Use `--keep-dot` flag to save the DOT file, then manually convert

### "Access Denied" when running association scripts
- **Solution**: Right-click the script and "Run as administrator"

### CSV file not opening with the program
- **Solution**: Manually associate CSV files with the executable
- **Alternative**: Use drag-and-drop or command line

### Large organizational charts
- **Solution**: Use `--keep-dot` to inspect the DOT file
- **Alternative**: Edit the DOT file manually for custom layouts

## Advanced Usage

### Building from Source

If you want to build the executable yourself:

1. **Install Go** from https://golang.org/dl/
2. **Clone the repository**
3. **Run the build script**:
   ```cmd
   build_windows.bat
   ```

### Custom File Associations

To manually set up file associations:

```cmd
# Associate .csv files (run as Administrator)
assoc .csv=OrgChart.CSV
ftype OrgChart.CSV="C:\path\to\orgchart_win64.exe" "%1"
```

### Batch Processing

Create a batch file to process multiple CSV files:

```cmd
@echo off
for %%f in (*.csv) do (
    echo Processing %%f...
    orgchart_win64.exe "%%f"
)
pause
```

## File Locations

- **Executable**: Place anywhere convenient
- **Input CSV**: Any location
- **Output PDF**: Same directory as input CSV (by default)
- **DOT files**: Same directory (if using `--keep-dot`)

## Performance Tips

- **Large organizations**: Use `--keep-dot` to inspect the graph structure
- **Frequent use**: Set up file associations for convenience
- **Batch processing**: Use command line scripts for multiple files

## Support

If you encounter issues:

1. **Check Graphviz installation**: `dot --version`
2. **Verify CSV format**: Ensure correct column structure
3. **Use `--keep-dot`**: Inspect the generated DOT file
4. **Check file permissions**: Ensure write access to output directory

## Requirements Summary

- **Windows**: 7, 8, 10, or 11
- **Architecture**: 32-bit, 64-bit, or ARM64
- **Graphviz**: Required for PDF generation
- **Permissions**: Administrator rights for file associations 