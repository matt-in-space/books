package books_test

import (
	"books"
	"cmp"
	"slices"
	"testing"
)

func mockCatalog() books.Catalog {
	catalog := books.Catalog{
		"1": {ID: "1", Title: "In Cold Blood", Author: "Truman Capote", Copies: 10},
		"2": {ID: "2", Title: "The Stand", Author: "Stephen King", Copies: 3},
	}

	return catalog
}

func TestBookToString_FormatsBookInfoAsString(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	book := books.Book{Copies: 1}
	book.SetCopies(5)

	if book.Copies != 5 {
		t.Fatalf("want book with 5 copies, got %d", book.Copies)
	}
}

func TestBookSetCopies_FailsOnNegativeCopies(t *testing.T) {
	t.Parallel()

	book := books.Book{}
	err := book.SetCopies(-5)

	if err == nil {
		t.Fatal("want error, got nil")
	}

}

func TestGetAllBooks_ReturnsAllBooks(t *testing.T) {
	t.Parallel()

	catalog := mockCatalog()
	got := catalog.GetAllBooks()
	assertBooksEqual(t, got)
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

	_, ok := catalog.GetBook("blah")

	if ok {
		t.Fatal("want not found, got book")
	}
}

func TestAddBook_AddsABook(t *testing.T) {
	catalog := books.Catalog{}
	catalog.AddBook(books.Book{})

	if len(catalog) != 1 {
		t.Fatal("want 1 book, got", len(catalog))
	}
}

func TestSetCopies_UpdatesCopies(t *testing.T) {
	catalog := mockCatalog()
	err := catalog.SetCopies("1", 5)

	if err != nil {
		t.Fatal(err)
	}

	if catalog["1"].Copies != 5 {
		t.Fatal("want 5 copies, got", catalog["1"].Copies)
	}
}

func TestSetCopies_RequiresValidCount(t *testing.T) {
	catalog := books.Catalog{}
	catalog.AddBook(books.Book{ID: "1"})

	err := catalog.SetCopies("1", -1)

	if err == nil {
		t.Fatal("want error, got nil")
	}
}

func TestCatalogSync_ReadsAndWritesCatalog(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/catalog.json"
	catalog := mockCatalog()

	err := catalog.Sync(path)

	if err != nil {
		t.Fatal(err)
	}

	written, err := books.OpenCatalog(path)

	if err != nil {
		t.Fatal(err)
	}

	got := written.GetAllBooks()
	assertBooksEqual(t, got)
}

func assertBooksEqual(t *testing.T, got []books.Book) {
	t.Helper()

	want := []books.Book{
		{ID: "1", Title: "In Cold Blood", Author: "Truman Capote", Copies: 10},
		{ID: "2", Title: "The Stand", Author: "Stephen King", Copies: 3},
	}

	slices.SortFunc(got, func(a, b books.Book) int {
		return cmp.Compare(a.ID, b.ID)
	})

	if !slices.Equal(want, got) {
		t.Fatalf("want %v, got %v", want, got)
	}
}
