package main

import (
	"books"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: find <book_id>")
		return
	}

	ID := os.Args[1]
	res, err := http.Get("http://localhost:3000/v1/find/" + ID)

	if err != nil {
		fmt.Printf("Request: %v\n", err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Reading response: %v\n", err)
		return
	}

	book := books.Book{}
	err = json.Unmarshal(body, &book)
	if err != nil {
		fmt.Printf("Unmarshalling response: %v\n", err)
		return
	}

	fmt.Println(book)
}
