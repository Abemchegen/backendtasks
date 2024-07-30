package services

import (
	"errors"
	"task3/models"
)

var bookcount, membercount = 0, 0

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member
}

func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

func (l *Library) AddBook(book models.Book) {
	book.ID = bookcount
	l.Books[bookcount] = book
	bookcount += 1
}

func (l *Library) RemoveBook(bookID int) error {

	_, exists := l.Books[bookID]

	if !exists {
		return errors.New("book not found")
	}
	delete(l.Books, bookID)
	return nil
}
func (l *Library) AddMember(name string) int {

	mem := models.Member{ID: membercount, Name: name}
	l.Members[membercount] = mem
	membercount += 1
	return mem.ID
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, exists := l.Books[bookID]

	if !exists {
		return errors.New("book not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("book is borrowed")
	}

	member, exists := l.Members[memberID]

	if !exists {
		return errors.New("member not found")
	}
	book.Status = "Borrowed"
	l.Books[bookID] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member
	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, exists := l.Books[bookID]

	if !exists {
		return errors.New("book not found")
	}
	if book.Status == "Available" {
		return errors.New("book was not borrowed ")
	}

	member, exists := l.Members[memberID]

	if !exists {
		return errors.New("member not found")
	}

	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			break
		}
	}

	book.Status = "Available"
	l.Books[bookID] = book
	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {

	var availablebooks []models.Book

	for _, book := range l.Books {

		if book.Status == "Available" {
			availablebooks = append(availablebooks, book)
		}
	}
	return availablebooks
}
func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.Members[memberID]

	if !exists {
		return nil
	}

	return member.BorrowedBooks
}
