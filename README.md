# OrgChart Generator

A command-line tool that converts CSV employee data into organizational charts as PDF files.

## Features

- **CSV to PDF**: Converts employee CSV data directly to PDF organizational charts
- **Cross-platform**: Works on macOS, Linux, and Windows
- **Flexible output**: Optional intermediate DOT file generation for debugging
- **Automatic layout**: Uses Graphviz for professional organizational chart layout
- **Windows drag-and-drop**: On Windows, you can drag and drop a CSV file onto the executable to generate a PDF instantly (no command line needed)

## Requirements

- **Go 1.24+**: Required to build the program
- **Graphviz**: Required for PDF generation (the `dot` command must be available)

### Installing Graphviz

#### macOS
```bash
brew install graphviz
```

#### Ubuntu/Debian
```bash
sudo apt-get install graphviz
```

#### Windows
```bash
# Using Chocolatey
choco install graphviz

# Or download from https://graphviz.org/download/
```

## Building

```bash
go build orgchart.go
```

This creates an executable named `orgchart` (or `orgchart.exe` on Windows).

## Usage

### Basic Usage

Generate a PDF from a CSV file:

```bash
./orgchart fcbh.csv
```

This creates `fcbh.pdf` in the same directory.

### Custom Output Name

Specify a custom output filename:

```bash
./orgchart fcbh.csv my_organization.pdf
```

### Keep Intermediate DOT File

To debug or inspect the generated graph structure, keep the intermediate DOT file:

```bash
./orgchart fcbh.csv --keep-dot
# or
./orgchart fcbh.csv -k
```

This creates both `fcbh.pdf` and `fcbh.dot` files.

### Complete Example

```bash
./orgchart fcbh.csv my_org_chart.pdf --keep-dot
```

Creates:
- `my_org_chart.pdf` - The organizational chart
- `my_org_chart.dot` - The Graphviz DOT file (for inspection)

### Windows Drag-and-Drop

- On Windows, you can drag and drop a CSV file onto the executable (`orgchart_win64.exe`) to generate a PDF in the same folder.
- You can also double-click a CSV file (after file association) or right-click → "Open with" → select the executable.
- All command line options are also available.

## CSV Format

**The input CSV should be exported from the UKG "Title name supervisor" report (export without Display Header/Footer).**

### CSV Example

```csv
"Employee Status","Date Hired","Default Jobs (HR) Full Path","Employee Name","Supervisor Name","Manager / Director Name","Cost Centers(Department)"
"Active","09/01/1972","President","GERALD A. JACKSON","GERALD A. JACKSON","","69_Corporate"
"Active","06/08/2015","Vice President of Technology","JONATHAN R. STEARLEY","JOSHUA A. MEE","CLAY JACKSON","771_Technology Development"
```

## Output

The program generates an organizational chart showing:

- **Employee boxes**: Each supervisor's direct reports in a table format
- **Hierarchy**: Lines connecting supervisors to their managers
- **Edge labels**: Number of direct reports and total reports (direct + indirect)
- **Statistics**: Summary information including total people, average span of control, and job title count

## Command Line Options

| Option | Description |
|--------|-------------|
| `--keep-dot` or `-k` | Keep the intermediate DOT file for inspection |
| `[output-pdf]` | Specify custom output filename (optional) |

## Error Handling

The program provides helpful error messages for common issues:

- **Missing CSV file**: "Error opening CSV file: ..."
- **Invalid CSV format**: "Error reading CSV: ..."
- **Graphviz not found**: "Error converting to PDF: ..." (with manual conversion instructions)

## Cross-Platform Support

### Building for Windows

#### Using Build Scripts (Recommended)

**From macOS/Linux:**
```bash
./build_windows.sh
```

**From Windows (PowerShell):**
```powershell
.\build_windows.ps1
```

**From Windows (Command Prompt):**
```cmd
build_windows.bat
```

These scripts build for multiple Windows architectures:
- Windows 64-bit (`orgchart_win64.exe`) - Most common
- Windows 32-bit (`orgchart_win32.exe`) - For older systems  
- Windows ARM64 (`orgchart_win_arm64.exe`) - For Windows on ARM

#### Windows Installation and Setup

For detailed Windows installation instructions, including drag-and-drop functionality, see [WINDOWS_INSTALL.md](WINDOWS_INSTALL.md).

#### Manual Build

From macOS/Linux, build a Windows executable:

```bash
GOOS=windows GOARCH=amd64 go build orgchart_dragdrop.go
```

This creates `orgchart_dragdrop.exe` for Windows.

### Building for Different Architectures

```bash
# Windows 64-bit
GOOS=windows GOARCH=amd64 go build orgchart_dragdrop.go

# Windows 32-bit
GOOS=windows GOARCH=386 go build orgchart_dragdrop.go

# Linux 64-bit
GOOS=linux GOARCH=amd64 go build orgchart.go

# macOS ARM (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build orgchart.go
```

## Troubleshooting

### "dot: command not found"
Install Graphviz on your system (see Requirements section).

### "Error converting to PDF"
The program will save the DOT file and provide manual conversion instructions:
```bash
dot -Tpdf fcbh.dot -o fcbh.pdf
```

### Large organizational charts
For very large organizations, consider:
- Using `--keep-dot` to inspect the DOT file
- Manually editing the DOT file for custom layouts
- Using Graphviz's layout engines: `dot`, `neato`, `fdp`, `sfdp`, `twopi`, `circo`

## License

This program is provided as-is for organizational chart generation. 