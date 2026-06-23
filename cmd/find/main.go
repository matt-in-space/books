package main

import (
	"books"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: find <book_id>")
		return
	}

	catalog, err := books.OpenCatalog("testdata/catalog.json")

	if err != nil {
		fmt.Printf("Opening catalog: %v\n", err)
		return
	}

	ID := os.Args[1]

	book, found := catalog.GetBook(ID)

	if !found {
		fmt.Println("Book not found")
		return
	}

	fmt.Println(book)
}
