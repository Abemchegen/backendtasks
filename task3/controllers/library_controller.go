package controllers

import (
	"bufio"
	"fmt"
	"strconv"
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

}
func (lc *LibraryController) AddMember(scanner *bufio.Scanner) {
	fmt.Println("enter name of the member")
	scanner.Scan()
	name := scanner.Text()
	id := lc.Library.AddMember(name)
	fmt.Printf("Your ID is %v\n", id)
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
}
func (lc *LibraryController) ListAvailableBooks() {
	books := lc.Library.ListAvailableBooks()
	fmt.Println("the list of available books are: ")

	if len(books) != 0 {
		for _, book := range books {
			fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
		}
	} else {
		fmt.Println("none")
	}

}

func (lc *LibraryController) ListBorrowedBooks(scanner *bufio.Scanner) {
	fmt.Print("Enter member ID to list borrowed books: ")
	scanner.Scan()
	memberID, err := strconv.Atoi(scanner.Text())
	for err != nil {
		fmt.Println("Invalid member ID. Please enter a number.")
		fmt.Print("Enter member ID to list borrowed books: ")
		scanner.Scan()
		memberID, err = strconv.Atoi(scanner.Text())
	}

	books := lc.Library.ListBorrowedBooks(memberID)
	fmt.Printf("books borrowed by member with if %d: \n", memberID)

	if len(books) != 0 {
		for _, book := range books {
			fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
		}
	} else {
		fmt.Println("none")
	}

}
