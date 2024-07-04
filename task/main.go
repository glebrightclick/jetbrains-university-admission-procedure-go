package main

import (
	"fmt"
	"sort"
)

type Applicant struct {
	firstName, lastName string
	GPA                 float64
}

func (applicant *Applicant) fullName() string {
	return fmt.Sprintf("%s %s", applicant.firstName, applicant.lastName)
}

func main() {
	// Read the first input, an N integer representing the total number of applicants
	var numberOfApplicants int
	fmt.Scan(&numberOfApplicants)

	// Read the second input, an M integer representing the number of applicants that should be accepted to the university.
	var shouldBeAccepted int
	fmt.Scan(&shouldBeAccepted)

	if numberOfApplicants < shouldBeAccepted || shouldBeAccepted == 0 {
		return
	}

	// Read N lines from the input.
	// Each line contains the first name, the last name, and the applicant's GPA.
	// These values are separated by one whitespace character. A GPA is a floating-point number with two decimals.

	// Write your code solution for the project below.
	applicants := make([]Applicant, numberOfApplicants)
	for i := range applicants {
		fmt.Scanf("%s %s %f", &applicants[i].firstName, &applicants[i].lastName, &applicants[i].GPA)
	}

	// Output the Successful applicants: message.
	fmt.Println("Successful applicants:")

	sort.Slice(applicants, func(i, j int) bool {
		if applicants[i].GPA == applicants[j].GPA {
			return applicants[i].fullName() < applicants[j].fullName()
		}

		return applicants[i].GPA > applicants[j].GPA
	})
	// Output M lines for applicants with the top-ranking GPAs.
	// Each line should contain the first and the last name of the applicant separated by a whitespace character.
	// Successful applicants' names should be printed in descending order depending on the GPA — first, the applicant with the best GPA, then the second-best, and so on.
	// If two applicants have the same GPA, rank them in alphabetical order using their full names (we know it's not really fair to choose students by their names, but what choice do we have?)
	for _, applicant := range applicants[:shouldBeAccepted] {
		fmt.Println(applicant.fullName())
	}
}
