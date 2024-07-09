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

type DepartmentName string

const mathematics DepartmentName = "Mathematics"
const physics DepartmentName = "Physics"
const biotech DepartmentName = "Biotech"
const chemistry DepartmentName = "Chemistry"
const engineering DepartmentName = "Engineering"

type Department struct {
	name        DepartmentName
	applyScores []DepartmentName // list of departments which scores should be applied
	applicants  []*Applicant
	waves       [][]*Applicant
}

func (d *Department) String() string {
	return strings.ToLower(string(d.name))
}

type Applicant struct {
	firstName, lastName string
	specialScore        float64
	scores              map[DepartmentName]float64
	department          *Department
}

func (applicant *Applicant) Score(department *Department) float64 {
	// Choose the best score for a student in the ranking: either the mean score for the final exam(s) or the special exam's score.
	total := 0.0
	for _, applyDepartmentName := range department.applyScores {
		total += applicant.scores[applyDepartmentName]
	}

	return max(total/float64(len(department.applyScores)), applicant.specialScore)
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
	// file, err := os.Open("/Users/divine/code/jetbrains/University Admission Procedure (Go)/University Admission Procedure (Go)/task/applicant_list_7.txt")
	if err != nil {
		log.Fatal(err)
	}

	departments := map[DepartmentName]Department{
		mathematics: {
			mathematics,
			// math for the Mathematics department
			[]DepartmentName{mathematics},
			make([]*Applicant, maxStudentsInDepartment),
			make([][]*Applicant, 3),
		},
		physics: {
			physics,
			// physics and math for the Physics department
			[]DepartmentName{physics, mathematics},
			make([]*Applicant, maxStudentsInDepartment),
			make([][]*Applicant, 3),
		},
		biotech: {
			biotech,
			// chemistry and physics for the Biotech department.
			[]DepartmentName{biotech, physics},
			make([]*Applicant, maxStudentsInDepartment),
			make([][]*Applicant, 3),
		},
		chemistry: {
			chemistry,
			// chemistry for the Chemistry department
			[]DepartmentName{chemistry},
			make([]*Applicant, maxStudentsInDepartment),
			make([][]*Applicant, 3),
		},
		engineering: {
			engineering,
			// computer science and math for the Engineering Department
			[]DepartmentName{engineering, mathematics},
			make([]*Applicant, maxStudentsInDepartment),
			make([][]*Applicant, 3),
		},
	}
	// Each line equals one applicant, their first name, last name, GPA, first priority department, second priority department, and third priority department.
	// Columns with values are separated by whitespace characters. For example, Laura Spungen 3.71 Physics Engineering Mathematics.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		info := strings.Split(scanner.Text(), " ")
		// physics, chemistry, math, computer science
		firstName, lastName := info[0], info[1]

		scoreIndex := map[int][]DepartmentName{4: {mathematics}, 2: {physics}, 3: {biotech, chemistry}, 5: {engineering}}
		scores := make(map[DepartmentName]float64)
		for i := 2; i < 6; i++ {
			score, err := strconv.ParseFloat(info[i], 64)
			if err != nil {
				log.Fatal(err)
			}

			for _, departmentName := range scoreIndex[i] {
				scores[departmentName] = score
			}
		}

		// Read the file named applicants.txt once again.
		// Mind one additional column, right after the last exam's result.
		// This column represents the special exam's score.
		// For example, Willie McBride 76 45 79 80 100 Physics Engineering Mathematics(where 100 is the admission exam's score).
		specialScore, err := strconv.ParseFloat(info[6], 64)
		if err != nil {
			log.Fatal(err)
		}

		applicant := Applicant{
			firstName,
			lastName,
			specialScore,
			scores,
			nil,
		}
		d1, d2, d3 := departments[DepartmentName(info[7])], departments[DepartmentName(info[8])], departments[DepartmentName(info[9])]
		d1.waves[0] = append(d1.waves[0], &applicant)
		d2.waves[1] = append(d2.waves[1], &applicant)
		d3.waves[2] = append(d3.waves[2], &applicant)
	}

	// Sort applicants according to their GPA and priorities (and names, if their GPA scores are the same).
	// As in the previous stage, if two applicants to the same department have the same GPA, sort them by their full names in alphabetical order.

	sortApplicants := func(applicants []*Applicant, department *Department) {
		sort.Slice(applicants, func(i, j int) bool {
			if applicants[i] == nil {
				return false
			} else if applicants[j] == nil {
				return true
			}

			iScore, jScore := applicants[i].Score(department), applicants[j].Score(department)

			if iScore == jScore {
				return applicants[i].fullName() < applicants[j].fullName()
			}

			return iScore > jScore
		})
	}

	// first of all, let's try to sort out first wave, then second, then third
	distribute := func(department *Department, waveNumber int) {
		// sort applicants by GPA + fullName
		sortApplicants(department.waves[waveNumber], department)

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
	var departmentNames []DepartmentName
	for departmentName := range departments {
		departmentNames = append(departmentNames, departmentName)
	}
	sort.Slice(departmentNames, func(i, j int) bool {
		return departmentNames[i] < departmentNames[j]
	})

	// Some departments are less popular than others, so there may be fewer available candidates for a department.
	// However, their number shouldn't be more than N.
	// Separate the name and the GPA with a whitespace character. Here's an example (you may add empty lines between the departments' lists to increase the text readability):
	for _, departmentName := range departmentNames {
		department := departments[departmentName]
		// department_name
		fmt.Println(department.name)
		// sort applicants by GPA & names
		sortApplicants(department.applicants, &department)
		// Instead of printing the results (you may do it if you want), output the admission lists to files.
		file, err := os.Create(department.String() + ".txt")
		if err != nil {
			log.Fatal(err)
		}

		for _, applicant := range department.applicants {
			if applicant == nil {
				break
			}

			// applicant1 GPA1
			// applicant2 GPA2
			// applicant3 GPA3
			// <...>
			display := fmt.Sprintf("%s %.2f\n", applicant.fullName(), applicant.Score(&department))
			fmt.Print(display)

			// Write the names of the students accepted to the department and their mean finals score to the corresponding file (one student per line).
			file.WriteString(display)
		}

		fmt.Println()
		// Close the file
		file.Close()
	}
}
