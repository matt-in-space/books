package main

import (
	"books"
	"fmt"
)

func main() {
	client := books.NewClient("localhost:3000")
	list, err := client.GetAllBooks()

	if err != nil {
		fmt.Printf("Request: %v\n", err)
		return
	}

	for _, book := range list {
		fmt.Println(book)
	}

}
