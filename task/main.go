package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Department struct {
	name       string
	applicants []*Applicant
	waves      [][]*Applicant
}

type Applicant struct {
	firstName, lastName string
	GPA                 float64
	department          *Department
}

func (applicant *Applicant) fullName() string {
	return fmt.Sprintf("%s %s", applicant.firstName, applicant.lastName)
}

func main() {
	// Read an N integer from the input. This integer represents the maximum number of students for each department.
	var maxStudentsInDepartment int
	fmt.Scan(&maxStudentsInDepartment)
	// Read the file named applicants.txt (this file is already included in the project's files, even though it is not visible; so you only need to download it if you want to take a closer look at it).
	file, err := os.Open("applicants.txt")
	if err != nil {
		log.Fatal(err)
	}

	departments := map[string]Department{
		"Mathematics": {"Mathematics", make([]*Applicant, maxStudentsInDepartment), make([][]*Applicant, 3)},
		"Physics":     {"Physics", make([]*Applicant, maxStudentsInDepartment), make([][]*Applicant, 3)},
		"Biotech":     {"Biotech", make([]*Applicant, maxStudentsInDepartment), make([][]*Applicant, 3)},
		"Chemistry":   {"Chemistry", make([]*Applicant, maxStudentsInDepartment), make([][]*Applicant, 3)},
		"Engineering": {"Engineering", make([]*Applicant, maxStudentsInDepartment), make([][]*Applicant, 3)},
	}
	// Each line equals one applicant, their first name, last name, GPA, first priority department, second priority department, and third priority department.
	// Columns with values are separated by whitespace characters. For example, Laura Spungen 3.71 Physics Engineering Mathematics.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		info := strings.Split(scanner.Text(), " ")
		GPA, err1 := strconv.ParseFloat(info[2], 64)
		if err1 != nil {
			fmt.Println("Error parsing GPA", info)
			continue
		}
		applicant := Applicant{
			info[0],
			info[1],
			GPA,
			nil,
		}
		d1, d2, d3 := departments[info[3]], departments[info[4]], departments[info[5]]
		d1.waves[0] = append(d1.waves[0], &applicant)
		d2.waves[1] = append(d2.waves[1], &applicant)
		d3.waves[2] = append(d3.waves[2], &applicant)
	}

	// Sort applicants according to their GPA and priorities (and names, if their GPA scores are the same).
	// As in the previous stage, if two applicants to the same department have the same GPA, sort them by their full names in alphabetical order.

	sortApplicants := func(applicants []*Applicant) {
		sort.Slice(applicants, func(i, j int) bool {
			if applicants[i] == nil {
				return false
			} else if applicants[j] == nil {
				return true
			}

			if applicants[i].GPA == applicants[j].GPA {
				return applicants[i].fullName() < applicants[j].fullName()
			}

			return applicants[i].GPA > applicants[j].GPA
		})
	}

	// first of all, let's try to sort out first wave, then second, then third
	distribute := func(department *Department, waveNumber int) {
		// sort applicants by GPA + fullName
		sortApplicants(department.waves[waveNumber])

		index, capacity := 0, len(department.applicants)
		for _, applicant := range department.applicants {
			if applicant == nil {
				break
			}
			index++
		}

		for _, applicant := range department.waves[waveNumber] {
			if index >= capacity {
				return
			}

			if applicant.department == nil {
				department.applicants[index] = applicant
				applicant.department = department
				index++
			}
		}
	}

	// distribute first wave
	for _, department := range departments {
		distribute(&department, 0)
	}
	// distribute second wave
	for _, department := range departments {
		distribute(&department, 1)
	}
	// distribute third wave
	for _, department := range departments {
		distribute(&department, 2)
	}

	// Print the departments in the alphabetic order (Biotech, Chemistry, Engineering, Mathematics, Physics), output the names and the GPA of enrolled applicants for each department.
	var departmentNames []string
	for departmentName := range departments {
		departmentNames = append(departmentNames, departmentName)
	}
	sort.Strings(departmentNames)
	// Some departments are less popular than others, so there may be fewer available candidates for a department.
	// However, their number shouldn't be more than N.
	// Separate the name and the GPA with a whitespace character. Here's an example (you may add empty lines between the departments' lists to increase the text readability):
	for _, departmentName := range departmentNames {
		department := departments[departmentName]
		// department_name
		fmt.Println(department.name)
		// sort applicants by GPA & names
		sortApplicants(department.applicants)
		for _, applicant := range department.applicants {
			if applicant == nil {
				break
			}

			// applicant1 GPA1
			// applicant2 GPA2
			// applicant3 GPA3
			// <...>
			fmt.Printf("%s %.2f\n", applicant.fullName(), applicant.GPA)
		}
		fmt.Println()
	}
}
