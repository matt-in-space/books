package main

import (
	"books"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: copies <book_id> <copies>")
		return
	}

	catalog, err := books.OpenCatalog("testdata/catalog.json")

	if err != nil {
		fmt.Printf("Opening catalog: %v\n", err)
		return
	}

	id := os.Args[1]
	copies := os.Args[2]

	copiesCount, err := strconv.Atoi(copies)

	if err != nil {
		fmt.Printf("Invalid copies: %v\n", err)
		return
	}

	err = catalog.SetCopies(id, copiesCount)

	if err != nil {
		fmt.Printf("Setting copies: %v\n", err)
		return
	}

	fmt.Printf("Copies updated to %d\n", copiesCount)
}
