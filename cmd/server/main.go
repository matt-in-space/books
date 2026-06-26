package main

import (
	"books"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("Usage: server <catalog_path>")
		os.Exit(1)
	}

	path := args[0]
	catalog, err := books.OpenCatalog(path)

	if err != nil {
		panic(err)
	}

	books.ListenAndServe(":3000", catalog)
}
