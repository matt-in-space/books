package main

import (
	"books"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	res, err := http.Get("http://localhost:3000/v1/list")

	if err != nil {
		fmt.Printf("Request: %v\n", err)
		return
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Request: %v\n", res.Status)
		return
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("Reading response: %v\n", err)
		return
	}

	list := []books.Book{}
	err = json.Unmarshal(data, &list)

	if err != nil {
		fmt.Printf("Parsing response: %v\n", err)
		return
	}

	for _, book := range list {
		fmt.Println(book)
	}
}
