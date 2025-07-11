package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func main() {
	// Check command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: orgchart <csv-file> [output-pdf]")
		fmt.Println("If output-pdf is not specified, will use csv-file name with .pdf extension")
		return
	}

	csvFile := os.Args[1]
	outputFile := ""
	if len(os.Args) >= 3 {
		outputFile = os.Args[2]
	} else {
		// Default to same name as CSV but with .pdf extension
		outputFile = csvFile[:len(csvFile)-4] + ".pdf"
	}

	// Open the CSV file
	f, err := os.Open(csvFile)
	if err != nil {
		fmt.Printf("Error opening CSV file: %v\n", err)
		return
	}
	defer f.Close()

	// Create a new CSV reader
	r := csv.NewReader(f)

	// Read the records from the file
	data, err := r.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
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
	dotFile := outputFile[:len(outputFile)-4] + ".dot"
	if err := os.WriteFile(dotFile, []byte(dotContent), 0644); err != nil {
		fmt.Printf("Error writing DOT file: %v\n", err)
		return
	}

	// Convert to PDF using CLI
	cmd := exec.Command("dot", "-Tpdf", dotFile, "-o", outputFile)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error converting to PDF: %v\n", err)
		fmt.Printf("DOT file generated: %s\n", dotFile)
		fmt.Printf("You can manually convert it with: dot -Tpdf %s -o %s\n", dotFile, outputFile)
		return
	}

	// Clean up DOT file
	os.Remove(dotFile)

	fmt.Printf("Org chart generated successfully: %s\n", outputFile)
}
