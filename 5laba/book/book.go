package book

import "fmt"

type Book struct {
	Title  string
	Author string
	Year   int
}

type Stringer interface {
	String() string
}

func NewBook(title string, author string, year int) *Book {
	return &Book{Title: title, Author: author, Year: year}
}

func (b Book) String() string {
	return fmt.Sprintf("Название: %s, Автор: %s, год выхода: книги %d", b.Title, b.Author, b.Year)

}
