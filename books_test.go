package books_test

import (
	"books"
	"testing"
)

func mockBook() books.Book {
	return books.Book{
		ID:     "1",
		Title:  "In Cold Blood",
		Author: "Truman Capote",
		Copies: 10,
	}
}

func mockCatalog() books.Catalog {
	catalog := books.Catalog{}

	catalog.AddBook(mockBook())

	return catalog
}

func TestBookToString_FormatsBookInfoAsString(t *testing.T) {
	input := books.Book{
		Title:  "In Cold Blood",
		Author: "Truman Capote",
		Copies: 10,
	}

	want := "In Cold Blood by Truman Capote (copies: 10)"
	got := input.String()

	if got != want {
		t.Fatalf("BookToString: wrong result: got %q, want %q", got, want)
	}
}

func TestBookSetCopies_UpdatesCopies(t *testing.T) {
	book := books.Book{Copies: 1}
	book.SetCopies(5)

	if book.Copies != 5 {
		t.Fatalf("want book with 5 copies, got %d", book.Copies)
	}
}

func TestBookSetCopies_FailsOnNegativeCopies(t *testing.T) {
	book := books.Book{}
	err := book.SetCopies(-5)

	if err == nil {
		t.Fatal("want error, got nil")
	}

}

func TestGetAllBooks_ReturnsAllBooks(t *testing.T) {
	catalog := mockCatalog()

	want := catalog

	got := catalog.GetAllBooks()

	if len(got) != len(want) {
		t.Fatalf("want %#v, got %#v", want, got)
	}
}

func TestGetBook_ReturnsFoundBook(t *testing.T) {
	catalog := mockCatalog()

	_, ok := catalog.GetBook("1")

	if !ok {
		t.Fatal("want found book, got book not found")
	}

}

func TestGetBook_ReturnsNotFound(t *testing.T) {
	catalog := mockCatalog()

	_, ok := catalog.GetBook("2")

	if ok {
		t.Fatal("want not found, got book")
	}
}

func TestAddBook_AddsABook(t *testing.T) {
	catalog := books.Catalog{}

	catalog.AddBook(mockBook())

	if len(catalog) != 1 {
		t.Fatal("want 1 book, got", len(catalog))
	}
}
