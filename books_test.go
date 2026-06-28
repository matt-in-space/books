package books_test

import (
	"books"
	"cmp"
	"encoding/json"
	"net"
	"net/http"
	"slices"
	"testing"
)

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

func TestGetAllBooks_OnClientReturnsAllBooks(t *testing.T) {
	t.Parallel()
	client := mockClient(t)
	books, err := client.GetAllBooks()

	if err != nil {
		t.Fatalf("Get all books failed with error %v", err)
	}

	assertBooksEqual(t, books)
}

func TestGetBook_ReturnsFoundBook(t *testing.T) {
	t.Parallel()
	catalog := mockCatalog()

	_, ok := catalog.GetBook("1")

	if !ok {
		t.Fatal("want found book, got book not found")
	}

}

func TestGetBook_ReturnsNotFound(t *testing.T) {
	t.Parallel()
	catalog := mockCatalog()

	_, ok := catalog.GetBook("blah")

	if ok {
		t.Fatal("want not found, got book")
	}
}

func TestGetBook_OnClientFindsBookById(t *testing.T) {
	t.Parallel()
	client := mockClient(t)
	got, err := client.GetBook("1")

	if err != nil {
		t.Fatal(err)
	}

	if Book1 != got {
		t.Fatalf("want %#v, got %#v", Book1, got)
	}
}

func TestAddBook_AddsABook(t *testing.T) {
	t.Parallel()
	catalog := books.NewCatalog()
	catalog.AddBook(books.Book{})
	total := len(catalog.GetAllBooks())

	if total != 1 {
		t.Fatal("want 1 book, got", total)
	}
}

func TestAddBook_FailsOnDuplicateID(t *testing.T) {
	t.Parallel()
	catalog := books.NewCatalog()
	catalog.AddBook(books.Book{ID: "1"})
	err := catalog.AddBook(books.Book{ID: "1"})

	if err == nil {
		t.Fatal("want error, got nil")
	}

	if len(catalog.GetAllBooks()) != 1 {
		t.Fatal("want 1 book, got", len(catalog.GetAllBooks()))
	}
}

func TestSetCopies_UpdatesCopies(t *testing.T) {
	t.Parallel()
	catalog := books.NewCatalog()
	catalog.AddBook(Book1)
	err := catalog.SetCopies("1", 5)

	if err != nil {
		t.Fatal(err)
	}

	book, _ := catalog.GetBook("1")

	if book.Copies != 5 {
		t.Fatal("want 5 copies, got", book.Copies)
	}
}

func TestSetCopies_RequiresValidCount(t *testing.T) {
	t.Parallel()
	catalog := mockCatalog()

	err := catalog.SetCopies("1", -1)

	if err == nil {
		t.Fatal("want error, got nil")
	}
}

func TestSetCopies_IsRaceFree(t *testing.T) {
	t.Parallel()
	catalog := mockCatalog()

	go func() {
		for range 100 {
			err := catalog.SetCopies("1", 0)
			if err != nil {
				panic(err)
			}
		}
	}()

	for range 100 {
		err := catalog.SetCopies("1", 0)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCatalogSync_ReadsAndWritesCatalog(t *testing.T) {
	t.Parallel()
	path := t.TempDir() + "/catalog.json"
	catalog := mockCatalog()
	catalog.Path = path

	err := catalog.Sync()

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

func TestNewCatalog_CreatesEmptyCatalog(t *testing.T) {
	t.Parallel()
	catalog := books.NewCatalog()
	books := catalog.GetAllBooks()
	if len(books) != 0 {
		t.Fatal("want empty catalog, got non-empty")
	}
}

func TestServer_ListsAllBooks(t *testing.T) {
	t.Parallel()
	catalog := mockCatalog()
	catalog.Path = t.TempDir() + "/catalog.json"
	addr := randomLocalAddr(t)

	go func() {
		err := books.ListenAndServe(addr, catalog)
		if err != nil {
			panic(err)
		}
	}()

	resp, err := http.Get("http://" + addr + "/v1/list")

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want status OK, got %v", resp.Status)
	}

	bookList := []books.Book{}
	err = json.NewDecoder(resp.Body).Decode(&bookList)

	if err != nil {
		t.Fatal(err)
	}

	assertBooksEqual(t, bookList)
}

func TestServer_FindsBook(t *testing.T) {
	t.Parallel()
	catalog := mockCatalog()
	catalog.Path = t.TempDir() + "/catalog.json"
	addr := randomLocalAddr(t)

	go func() {
		err := books.ListenAndServe(addr, catalog)
		if err != nil {
			panic(err)
		}
	}()

	resp, err := http.Get("http://" + addr + "/v1/find/1")

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("want status OK, got %v", resp.Status)
	}

	book := books.Book{}
	err = json.NewDecoder(resp.Body).Decode(&book)

	if err != nil {
		t.Fatal(err)
	}

	if book != Book1 {
		t.Fatalf("want book %v, got %v", Book1, book)
	}
}

func TestServer_FindBookReturnsNotFound(t *testing.T) {
	t.Parallel()
	catalog := mockCatalog()
	catalog.Path = t.TempDir() + "/catalog.json"
	addr := randomLocalAddr(t)

	go func() {
		err := books.ListenAndServe(addr, catalog)
		if err != nil {
			panic(err)
		}
	}()

	resp, err := http.Get("http://" + addr + "/v1/find/oof")

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("want status NotFound, got %v", resp.Status)
	}
}

// *** Helper functions ***
var (
	Book1 = books.Book{ID: "1", Title: "In Cold Blood", Author: "Truman Capote", Copies: 10}
	Book2 = books.Book{ID: "2", Title: "The Stand", Author: "Stephen King", Copies: 3}
)

func mockCatalog() *books.Catalog {
	catalog := books.NewCatalog()
	catalog.AddBook(Book1)
	catalog.AddBook(Book2)
	return catalog
}

func mockClient(t *testing.T) *books.Client {
	addr := randomLocalAddr(t)
	catalog := mockCatalog()
	catalog.Path = t.TempDir() + "/catalog.json"

	go func() {
		err := books.ListenAndServe(addr, catalog)
		if err != nil {
			panic(err)
		}
	}()

	return books.NewClient(addr)
}

func assertBooksEqual(t *testing.T, got []books.Book) {
	t.Helper()

	want := []books.Book{
		Book1,
		Book2,
	}

	slices.SortFunc(got, func(a, b books.Book) int {
		return cmp.Compare(a.ID, b.ID)
	})

	if !slices.Equal(want, got) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func randomLocalAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}
