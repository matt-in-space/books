package main

import (
	"books"
	"fmt"
	"os"
)

func main() {
	ID := os.Args[1]

	catalog := books.GetCatalog()
	book, found := catalog.GetBook(ID)

	if !found {
		fmt.Println("Book not found")
		return
	}

	fmt.Println(book)
}
