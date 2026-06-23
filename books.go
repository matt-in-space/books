package books

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
)

type Book struct {
	ID     string
	Title  string
	Author string
	Copies int
}

type Catalog map[string]Book

func GetCatalog() Catalog {
	return map[string]Book{
		"1": {ID: "1", Title: "In Cold Blood", Author: "Truman Capote", Copies: 10},
		"2": {ID: "2", Title: "The Phantom Toolbooth", Author: "Norton Juster", Copies: 4},
		"3": {ID: "3", Title: "The Way of Kings", Author: "Brandon Sanderson", Copies: 1},
	}
}

func OpenCatalog(path string) (Catalog, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	// `defer` ensures that `file.Close()` is called when the function returns, regardless
	// of where that return statement is reached
	// This can be used for cleanup in numerous scenarios
	defer file.Close()
	catalog := Catalog{}
	err = json.NewDecoder(file).Decode(&catalog)

	if err != nil {
		return nil, err
	}

	return catalog, nil
}

func (book Book) String() string {
	return fmt.Sprintf("%s by %s (copies: %d)", book.Title, book.Author, book.Copies)
}

func (book *Book) SetCopies(copies int) error {
	if copies < 0 {
		return fmt.Errorf("negative number of copies: %d", copies)
	}

	book.Copies = copies
	return nil
}

func (catalog Catalog) GetAllBooks() []Book {
	return slices.Collect(maps.Values(catalog))
}

func (catalog Catalog) AddBook(book Book) {
	catalog[book.ID] = book
}

func (catalog Catalog) GetBook(ID string) (Book, bool) {
	book, ok := catalog[ID]
	return book, ok
}

func (catalog Catalog) RemoveBook(ID string) {
	delete(catalog, ID)
}
