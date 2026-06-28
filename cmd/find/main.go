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

	ID := os.Args[1]
	client := books.NewClient("localhost:3000")
	book, err := client.GetBook(ID)

	if err != nil {
		fmt.Printf("Request: %v\n", err)
		return
	}

	fmt.Println(book)
}
