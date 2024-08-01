# Library Management System Documentation

## Overview

This document provides an overview of the Library Management System implemented in Go. The system allows users to manage a library's collection of books and member activities such as borrowing and returning books.

## Folder Structure

library_management/
├── main.go
├── controllers/
│ └── library_controller.go
├── models/
│ └── book.go
│ └── member.go
├── services/
│ └── library_service.go
├── docs/
│ └── documentation.md
└── go.mod

- main.go: Entry point of the application.
- controllers/library_controller.go: Handles console input and invokes the appropriate service methods.
- models/book.go: Defines the Book struct.
- models/member.go: Defines the Member struct.
- services/library_service.go: Contains business logic and data manipulation functions.
- docs/documentation.md: Contains system documentation and other related information.
- go.mod: Defines the module and its dependencies.

## Structs

### Book

type Book struct {
ID int
Title string
Author string
Status string // can be "Available" or "Borrowed"
}

### Member

type Member struct {
ID int
Name string
BorrowedBooks []Book // a slice to hold borrowed books
}

## Interfaces

### LibraryManager

type LibraryManager interface {
AddBook(book Book)
RemoveBook(bookID int)
BorrowBook(bookID int, memberID int) error
ReturnBook(bookID int, memberID int) error
ListAvailableBooks() []Book
ListBorrowedBooks(memberID int) []Book
}

## Implementation

### Library

The Library struct implements the LibraryManager interface.

type Library struct {
Books map[int]Book
Members map[int]Member
}

### Methods

- AddBook(book Book)
  Adds a new book to the library.

- RemoveBook(bookID int)
  Removes a book from the library by its ID.

- BorrowBook(bookID int, memberID int) error
  Allows a member to borrow a book if it is available.

- ReturnBook(bookID int, memberID int) error
  Allows a member to return a borrowed book.

- ListAvailableBooks() []Book
  Lists all available books in the library.

- ListBorrowedBooks(memberID int) []Book
  Lists all books borrowed by a specific member.

## Console Interaction

### Add a New Book

Prompts the user to enter book details and adds the book to the library.

### Remove an Existing Book

Prompts the user to enter the book ID and removes the book from the library.

### Borrow a Book

Prompts the user to enter the book ID and member ID, and allows the member to borrow the book if it is available.

### Return a Book

Prompts the user to enter the book ID and member ID, and allows the member to return the borrowed book.

### List All Available Books

Displays all available books in the library.

### List All Borrowed Books by a Member

Prompts the user to enter the member ID and displays all books borrowed by the member.

## Error Handling

The system includes error handling for scenarios where books or members are not found, or books are already borrowed. Appropriate error messages are displayed to the user.

## Usage

1. Start the Application: Run the main.go file to start the console application.
2. Follow Prompts: Use the console prompts to add, remove, borrow, and return books, or to list available and borrowed books.

## Dependencies

The go.mod file defines the module and its dependencies. Ensure that Go is installed and properly set up on your machine.

## Conclusion

This Library Management System demonstrates the use of structs, interfaces, methods, slices, and maps in Go. The system is designed to be simple and easy to use, with a focus on showcasing key Go functionalities.
