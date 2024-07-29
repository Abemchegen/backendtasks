package main

import (
	"fmt"
)

type Grade struct {
	subject string
	grade   float64
}

func average(s []Grade, num float64) float64 {

	var total float64

	for i := 0; i < len(s); i++ {
		total += s[i].grade
	}

	return total / num

}
func main() {

	fmt.Println("Enter your name:")
	var name string
	fmt.Scan(&name)

	fmt.Println("Enter the number of subjects:")
	var n int
	fmt.Scan(&n)

	subjects := make([]Grade, n)

	for i := 0; i < n; i++ {
		fmt.Printf("Enter the name of subject %d:\n", i+1)
		fmt.Scan(&subjects[i].subject)

		fmt.Printf("Enter the grade for subject %d:\n", i+1)
		fmt.Scan(&subjects[i].grade)
		for subjects[i].grade < 0 || subjects[i].grade > 100 {
			fmt.Println("Invalid value for grade. Please enter a value between 0 and 100.")
			fmt.Scan(&subjects[i].grade)
		}
	}

	average := average(subjects, float64(n))

	fmt.Printf("\nStudent Name: %s\n", name)
	fmt.Println("Subjects and Grades:")
	fmt.Printf("\n")
	for i := 0; i < len(subjects); i++ {
		fmt.Printf("Subject %d: %s Grade: %.2f\n", i+1, subjects[i].subject, subjects[i].grade)
	}
	fmt.Printf("Average Grade: %.2f\n", average)
}
