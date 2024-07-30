package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your name:")
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	name = strings.TrimSpace(name)

	fmt.Println("Enter the number of subjects:")
	nStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	nStr = strings.TrimSpace(nStr)
	n, err := strconv.Atoi(nStr)
	if err != nil {
		fmt.Println("Invalid input for number of subjects. Please enter a valid number.")
		return
	}

	subjects := make([]Grade, n)

	for i := 0; i < n; i++ {
		fmt.Printf("Enter the name of subject %d:\n", i+1)
		subjectName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			i--
			continue
		}
		subjects[i].subject = strings.TrimSpace(subjectName)

		fmt.Printf("Enter the grade for subject %d:\n", i+1)
		gradeStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			i--
			continue
		}
		gradeStr = strings.TrimSpace(gradeStr)

		grade, err := strconv.ParseFloat(gradeStr, 64)
		if err != nil {
			fmt.Println("Invalid input for grade. Please enter a valid number.")
			i--
			continue
		}

		if grade < 0 || grade > 100 {
			fmt.Println("Invalid value for grade. Please enter a value between 0 and 100.")
			i--
			continue
		}

		subjects[i].grade = grade
	}

	avg := average(subjects, float64(n))

	fmt.Printf("\nStudent Name: %s\n", name)
	fmt.Println("Subject Report:")
	fmt.Println("-----------------------------------------------------------")
	fmt.Printf("| %-20s | %-6s \n", "Subject", "Grade")
	fmt.Println("-----------------------------------------------------------")
	for i := 0; i < len(subjects); i++ {
		fmt.Printf("| %-20s | %-6.2f \n", subjects[i].subject, subjects[i].grade)
	}
	fmt.Println("-----------------------------------------------------------")
	fmt.Printf("| %-20s | %-6.2f \n", "Average Grade: ", avg)
}
