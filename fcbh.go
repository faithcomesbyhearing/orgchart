package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"time"
)

func main() {
	// Open the CSV file.
	// visualized order of people matches order in fcbh.csv (eg alphabetical first name)
	f, err := os.Open("fcbh.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// Create a new CSV reader.
	r := csv.NewReader(f)

	// Read the records from the file.
	data, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	type staffgroup struct {
		staff string
		count int
		total int
	}

	currentTime := time.Now()
	time := currentTime.Format("2006-01-02")
	fmt.Println("digraph \"orgchart\" {\n graph [ rankdir=\"LR\", splines=true]; overlap=false; node[shape=box color=none]; edge[color=grey];")

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
		fmt.Printf("%s [label=<<table>%s</table>>];\n", supervisor, sgroup.staff)
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
				fmt.Printf("%s:\"%s\":e -> %s:w [label=\"%d\"];\n", personsSupervisor[supervisor], supervisor, supervisor, sgroup.count)
			} else {
				fmt.Printf("%s:\"%s\":e -> %s:w [label=\"%d/%d\"];\n", personsSupervisor[supervisor], supervisor, supervisor, sgroup.count, sgroup.total)
			}
		}
	}

	fmt.Println("title [label=<<table border=\"1\" cellborder=\"0\" color=\"#888888\"><tr><td><b>Faith Comes By Hearing</b></td></tr>")
	fmt.Printf("<tr><td>rows show Name, Title, Hire Date (sorted increasing)</td></tr>\n", titles)
	fmt.Printf("<tr><td>lines show timecard-approval role in UKG</td></tr>\n")
	fmt.Printf("<tr><td>lines labels show the number of direct/indirect timecards</td></tr>\n")
	fmt.Printf("<tr><td>%d people / %d timecard-signers = %.1f avg span</td></tr>\n", people, groups, float32(people)/float32(groups))
	fmt.Printf("<tr><td>%d distinct job titles</td></tr>\n", titles)
	fmt.Printf("<tr><td>Chart generated %s</td></tr></table>>];\n", time)
	fmt.Println("{rank=same; GOD title GERALDAJACKSON;}")

	fmt.Println("}")
}
