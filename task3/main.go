package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"task3/controllers"
	"task3/services"
)

func main() {
	lib := services.NewLibrary()
	controller := &controllers.LibraryController{Library: lib}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n\t\t\t Library Management System")
		fmt.Println(strings.Repeat("-", 75))
		fmt.Println("\n1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Add Member")
		fmt.Println("4. Borrow Book")
		fmt.Println("5. Return Book")
		fmt.Println("6. List Available Books")
		fmt.Println("7. List Borrowed Books")
		fmt.Println("8. Exit")
		fmt.Println(strings.Repeat("-", 75))
		fmt.Print("Choose an option: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			controller.AddBook(scanner)
		case "2":
			controller.RemoveBook(scanner)
		case "3":
			controller.AddMember(scanner)
		case "4":
			controller.BorrowBook(scanner)
		case "5":
			controller.ReturnBook(scanner)
		case "6":
			controller.ListAvailableBooks()
		case "7":
			controller.ListBorrowedBooks(scanner)
		case "8":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 8.")
		}
	}
}
