package controllers

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"task3/models"
	"task3/services"
)

type LibraryController struct {
	Library *services.Library
}

func (lc *LibraryController) AddBook(scanner *bufio.Scanner) {
	fmt.Print("Enter book title: ")
	scanner.Scan()
	title := scanner.Text()

	fmt.Print("Enter book author: ")
	scanner.Scan()
	author := scanner.Text()

	book := models.Book{Title: title, Author: author, Status: "Available"}
	lc.Library.AddBook(book)
	fmt.Println("Book added successfully.")
	fmt.Println(strings.Repeat("-", 75))

}

func (lc *LibraryController) RemoveBook(scanner *bufio.Scanner) {
	fmt.Print("Enter book ID to remove: ")
	scanner.Scan()
	id, err := strconv.Atoi(scanner.Text())

	for err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		fmt.Print("Enter book ID to remove: ")
		scanner.Scan()
		id, err = strconv.Atoi(scanner.Text())
	}
	e := lc.Library.RemoveBook(id)

	if e != nil {
		fmt.Println("error removing book: ", e)
	} else {
		fmt.Printf("removed book with id: %d\n", id)
	}
	fmt.Println(strings.Repeat("-", 75))
}
func (lc *LibraryController) AddMember(scanner *bufio.Scanner) {
	fmt.Println("enter name of the member")
	scanner.Scan()
	name := scanner.Text()
	id := lc.Library.AddMember(name)
	fmt.Printf("Your ID is %v\n", id)
	fmt.Println(strings.Repeat("-", 75))

}
func (lc *LibraryController) BorrowBook(scanner *bufio.Scanner) {
	fmt.Print("Enter book ID to borrow: ")
	scanner.Scan()
	id, err := strconv.Atoi(scanner.Text())

	for err != nil {
		fmt.Println("Invalid book ID. Please enter a number.")
		fmt.Print("Enter book ID to borrow: ")
		scanner.Scan()
		id, err = strconv.Atoi(scanner.Text())
	}

	fmt.Print("Enter member ID: ")
	scanner.Scan()
	memberID, err := strconv.Atoi(scanner.Text())
	for err != nil {
		fmt.Println("Invalid member ID. Please enter a number.")
		fmt.Print("Enter member ID: ")
		scanner.Scan()
		memberID, err = strconv.Atoi(scanner.Text())
	}

	e := lc.Library.BorrowBook(id, memberID)
	if e != nil {
		fmt.Println("error borrowing book: ", e)
	} else {
		fmt.Printf("book with ID %d was borrowed by member %d", id, memberID)
	}
	fmt.Println(strings.Repeat("-", 75))

}

func (lc *LibraryController) ReturnBook(scanner *bufio.Scanner) {
	fmt.Print("Enter book ID to return: ")
	scanner.Scan()
	id, err := strconv.Atoi(scanner.Text())
	for err != nil {
		fmt.Println("Invalid book ID. Please enter a number.")
		fmt.Print("Enter book ID to return: ")
		scanner.Scan()
		id, err = strconv.Atoi(scanner.Text())
	}

	fmt.Print("Enter member ID: ")
	scanner.Scan()
	memberID, err := strconv.Atoi(scanner.Text())
	for err != nil {
		fmt.Println("Invalid member ID. Please enter a number.")
		fmt.Print("Enter member ID: ")
		scanner.Scan()
		memberID, err = strconv.Atoi(scanner.Text())
	}

	e := lc.Library.ReturnBook(id, memberID)

	if e != nil {
		fmt.Println("error return book: ", e)
	} else {
		fmt.Printf("book with IF %d returned by member with ID %d\n", id, memberID)
	}
	fmt.Println(strings.Repeat("-", 75))

}
func (lc *LibraryController) ListAvailableBooks() {
	books := lc.Library.ListAvailableBooks()
	fmt.Println("List of Available Books:")
	fmt.Println(strings.Repeat("-", 75))
	if len(books) != 0 {
		fmt.Printf("| %-5s | %-30s | %-30s |\n", "ID", "Title", "Author")
		fmt.Println(strings.Repeat("-", 75))
		for _, book := range books {
			fmt.Printf("| %-5d | %-30s | %-30s |\n", book.ID, book.Title, book.Author)
		}
	} else {
		fmt.Println("None")
	}
	fmt.Println(strings.Repeat("-", 75))
}

func (lc *LibraryController) ListBorrowedBooks(scanner *bufio.Scanner) {
	fmt.Print("Enter member ID to list borrowed books: ")
	scanner.Scan()
	memberIDStr := strings.TrimSpace(scanner.Text())
	memberID, err := strconv.Atoi(memberIDStr)
	for err != nil {
		fmt.Println("Invalid member ID. Please enter a number.")
		fmt.Print("Enter member ID to list borrowed books: ")
		scanner.Scan()
		memberIDStr = strings.TrimSpace(scanner.Text())
		memberID, err = strconv.Atoi(memberIDStr)
	}

	books := lc.Library.ListBorrowedBooks(memberID)
	fmt.Printf("Books borrowed by member with ID %d:\n", memberID)
	fmt.Println(strings.Repeat("-", 75))
	if len(books) != 0 {
		fmt.Printf("| %-5s | %-30s | %-30s |\n", "ID", "Title", "Author")
		fmt.Println(strings.Repeat("-", 75))
		for _, book := range books {
			fmt.Printf("| %-5d | %-30s | %-30s |\n", book.ID, book.Title, book.Author)
		}
	} else {
		fmt.Println("None")
	}
	fmt.Println(strings.Repeat("-", 75))
}
