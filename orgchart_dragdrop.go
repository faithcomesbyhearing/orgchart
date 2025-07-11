package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	// Check if we're running on Windows
	if os.PathSeparator != '\\' {
		fmt.Println("This program is designed for Windows drag-and-drop functionality.")
		fmt.Println("Please use the regular orgchart.exe for command-line usage.")
		os.Exit(1)
	}

	// Get command line arguments
	args := os.Args[1:]

	// If no arguments, show usage
	if len(args) == 0 {
		showUsage()
		return
	}

	// Check if the first argument is a CSV file
	csvFile := args[0]
	if !strings.HasSuffix(strings.ToLower(csvFile), ".csv") {
		fmt.Printf("Error: '%s' is not a CSV file.\n", csvFile)
		fmt.Println("Please drag and drop a CSV file onto this executable.")
		showUsage()
		return
	}

	// Check if file exists
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist.\n", csvFile)
		return
	}

	// Generate output filename
	outputFile := generateOutputFilename(csvFile)
	keepDot := false

	// Parse additional arguments
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "--keep-dot" || arg == "-k" {
			keepDot = true
		} else if !strings.HasPrefix(arg, "-") {
			// If it's not a flag, treat as output filename
			outputFile = arg
		}
	}

	// Process the CSV file
	if err := processCSV(csvFile, outputFile, keepDot); err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\nPress any key to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	// Success message
	fmt.Printf("\n✅ Org chart generated successfully: %s\n", outputFile)
	if keepDot {
		dotFile := strings.TrimSuffix(outputFile, ".pdf") + ".dot"
		fmt.Printf("📄 DOT file saved: %s\n", dotFile)
	}

	fmt.Println("\nPress any key to exit...")
	fmt.Scanln()
}

func showUsage() {
	fmt.Println("=== OrgChart Generator (Windows Drag-and-Drop) ===")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("1. Drag and drop a CSV file onto this executable")
	fmt.Println("2. Or run from command line: orgchart_dragdrop.exe fcbh.csv")
	fmt.Println("")
	fmt.Println("Optional arguments:")
	fmt.Println("  --keep-dot, -k    Keep the intermediate DOT file")
	fmt.Println("  [output-pdf]      Specify custom output filename")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  orgchart_dragdrop.exe fcbh.csv")
	fmt.Println("  orgchart_dragdrop.exe fcbh.csv my_chart.pdf --keep-dot")
	fmt.Println("")
	fmt.Println("Requirements:")
	fmt.Println("  - Graphviz must be installed (choco install graphviz)")
	fmt.Println("  - CSV file must have the correct format (see README.md)")
}

func generateOutputFilename(csvFile string) string {
	// Get the directory and base name
	dir := filepath.Dir(csvFile)
	base := filepath.Base(csvFile)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// Create output filename in the same directory
	return filepath.Join(dir, name+".pdf")
}

func processCSV(csvFile, outputFile string, keepDot bool) error {
	// Open the CSV file
	f, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("opening CSV file: %v", err)
	}
	defer f.Close()

	// Create a new CSV reader
	r := csv.NewReader(f)

	// Read the records from the file
	data, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("reading CSV: %v", err)
	}

	type staffgroup struct {
		staff string
		count int
		total int
	}

	currentTime := time.Now()
	timeStr := currentTime.Format("2006-01-02")

	// Generate DOT content
	dotContent := "digraph \"orgchart\" {\n graph [ rankdir=\"LR\", splines=true]; overlap=false; node[shape=box color=none]; edge[color=grey];\n"

	// make lists
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z]+`)
	var yearRegex = regexp.MustCompile(`[0-9]{4}`)
	var ampRegex = regexp.MustCompile(`&`)
	supervisorsStaff := make(map[string]*staffgroup)
	personsSupervisor := make(map[string]string)
	titleCount := make(map[string]int)
	var person, supervisor, title string
	var groups, people, titles int

	fmt.Println("📊 Processing CSV data...")

	for _, row := range data {
		person = nonAlphanumericRegex.ReplaceAllString(row[3], "")
		if person == "EmployeeName" || person == "ALANDHOOKER" {
			continue
		}
		supervisor = nonAlphanumericRegex.ReplaceAllString(row[4], "")
		title = row[2]

		if person == supervisor {
			supervisor = "GOD"
		}
		if _, ok := supervisorsStaff[supervisor]; !ok { // new supervisor
			supervisorsStaff[supervisor] = &staffgroup{"", 0, 0}
			groups++
		}
		if _, ok := titleCount[title]; !ok {
			titleCount[title]++
			titles++
		}

		people++
		supervisorsStaff[supervisor].staff += "\n<tr><td align=\"left\">" + row[3] + "</td>\n <td align=\"left\">" + ampRegex.ReplaceAllString(row[2], "&amp;") + "</td>\n <td align=\"left\" port=\"" + person + "\">" + yearRegex.FindString(row[1]) + "</td></tr>"
		supervisorsStaff[supervisor].count++
		personsSupervisor[person] = supervisor
	}

	var first bool
	var pcount int
	for supervisor, sgroup := range supervisorsStaff {
		first = true
		dotContent += fmt.Sprintf("%s [label=<<table>%s</table>>];\n", supervisor, sgroup.staff)
	NEXTSUPERVISOR:
		if len(personsSupervisor[supervisor]) > 0 {
			if first {
				pcount = sgroup.count
				first = false
			} else {
				sgroup = supervisorsStaff[supervisor]
				sgroup.total += pcount
				supervisor = personsSupervisor[supervisor]
			}
			goto NEXTSUPERVISOR
		}
	}
	for supervisor, sgroup := range supervisorsStaff {
		if len(personsSupervisor[supervisor]) > 0 {
			if sgroup.total == sgroup.count {
				dotContent += fmt.Sprintf("%s:\"%s\":e -> %s:w [label=\"%d\"];\n", personsSupervisor[supervisor], supervisor, supervisor, sgroup.count)
			} else {
				dotContent += fmt.Sprintf("%s:\"%s\":e -> %s:w [label=\"%d/%d\"];\n", personsSupervisor[supervisor], supervisor, supervisor, sgroup.count, sgroup.total)
			}
		}
	}

	dotContent += "title [label=<<table border=\"1\" cellborder=\"0\" color=\"#888888\"><tr><td><b>Faith Comes By Hearing</b></td></tr>"
	dotContent += fmt.Sprintf("<tr><td>rows show Name, Title, Hire Date (sorted increasing)</td></tr>\n", titles)
	dotContent += fmt.Sprintf("<tr><td>lines show timecard-approval role in UKG</td></tr>\n")
	dotContent += fmt.Sprintf("<tr><td>lines labels show the number of direct/indirect timecards</td></tr>\n")
	dotContent += fmt.Sprintf("<tr><td>%d people / %d timecard-signers = %.1f avg span</td></tr>\n", people, groups, float32(people)/float32(groups))
	dotContent += fmt.Sprintf("<tr><td>%d distinct job titles</td></tr>\n", titles)
	dotContent += fmt.Sprintf("<tr><td>Chart generated %s</td></tr></table>>];\n", timeStr)
	dotContent += "{rank=same; GOD title GERALDAJACKSON;}\n"
	dotContent += "}"

	// Write DOT file
	dotFile := strings.TrimSuffix(outputFile, ".pdf") + ".dot"
	if err := os.WriteFile(dotFile, []byte(dotContent), 0644); err != nil {
		return fmt.Errorf("writing DOT file: %v", err)
	}

	fmt.Println("🔨 Converting to PDF...")

	// Convert to PDF using CLI
	cmd := exec.Command("dot", "-Tpdf", dotFile, "-o", outputFile)
	if err := cmd.Run(); err != nil {
		// If dot command fails, provide helpful error message
		return fmt.Errorf("converting to PDF: %v\n\nMake sure Graphviz is installed:\nchoco install graphviz\n\nOr download from: https://graphviz.org/download/", err)
	}

	// Clean up DOT file unless --keep-dot flag is used
	if !keepDot {
		os.Remove(dotFile)
	}

	return nil
}
