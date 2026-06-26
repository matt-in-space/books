package books

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"sync"
)

type Book struct {
	ID     string
	Title  string
	Author string
	Copies int
}

type Catalog struct {
	mu   *sync.RWMutex
	data map[string]Book
	Path string
}

func NewCatalog() *Catalog {
	return &Catalog{
		mu:   &sync.RWMutex{},
		data: map[string]Book{},
	}
}

func OpenCatalog(path string) (*Catalog, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	// `defer` ensures that `file.Close()` is called when the function returns, regardless
	// of where that return statement is reached
	// This can be used for cleanup in numerous scenarios
	defer file.Close()
	catalog := NewCatalog()
	catalog.Path = path
	err = json.NewDecoder(file).Decode(&catalog.data)

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

func (catalog *Catalog) GetAllBooks() []Book {
	catalog.mu.RLock()
	defer catalog.mu.RUnlock()
	return slices.Collect(maps.Values(catalog.data))
}

func (catalog *Catalog) AddBook(book Book) error {
	catalog.mu.Lock()
	defer catalog.mu.Unlock()
	if _, ok := catalog.data[book.ID]; ok {
		return fmt.Errorf("book already exists: %s", book.ID)
	}

	catalog.data[book.ID] = book
	return nil
}

func (catalog *Catalog) GetBook(ID string) (Book, bool) {
	catalog.mu.RLock()
	defer catalog.mu.RUnlock()
	book, ok := catalog.data[ID]
	return book, ok
}

func (catalog *Catalog) RemoveBook(ID string) {
	catalog.mu.Lock()
	defer catalog.mu.Unlock()
	delete(catalog.data, ID)
}

func (catalog *Catalog) SetCopies(ID string, copies int) error {
	if copies < 0 {
		return fmt.Errorf("negative number of copies: %d", copies)
	}

	catalog.mu.Lock()
	defer catalog.mu.Unlock()
	book, ok := catalog.data[ID]

	if !ok {
		return fmt.Errorf("book not found: %s", ID)
	}

	book.Copies = copies
	catalog.data[ID] = book

	return nil
}

func (catalog *Catalog) GetCopies(ID string) (int, error) {
	catalog.mu.RLock()
	defer catalog.mu.RUnlock()
	book, ok := catalog.data[ID]

	if !ok {
		return 0, fmt.Errorf("book not found: %s", ID)
	}

	return book.Copies, nil
}

func (catalog *Catalog) Sync() error {
	catalog.mu.RLock()
	defer catalog.mu.RUnlock()
	file, err := os.Create(catalog.Path)

	if err != nil {
		return err
	}

	defer file.Close()
	err = json.NewEncoder(file).Encode(catalog.data)

	if err != nil {
		return err
	}

	return nil
}
