package main

import (
	"books"
	"cmp"
	"fmt"
	"slices"
)

func main() {
	catalog, err := books.OpenCatalog("testdata/catalog.json")

	if err != nil {
		fmt.Printf("Opening catalog: %v\n", err)
		return
	}

	catalogBooks := catalog.GetAllBooks()

	slices.SortFunc(catalogBooks, func(a, b books.Book) int {
		return cmp.Compare(a.Author, b.Author)
	})

	for _, book := range catalogBooks {
		fmt.Println(book)
	}
}
