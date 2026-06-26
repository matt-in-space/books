package books

import (
	"encoding/json"
	"net/http"
)

func ListenAndServe(addr string, catalog *Catalog) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/find/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		book, ok := catalog.GetBook(id)

		if !ok {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}

		err := json.NewEncoder(w).Encode(book)

		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/v1/list", func(w http.ResponseWriter, r *http.Request) {
		books := catalog.GetAllBooks()
		err := json.NewEncoder(w).Encode(books)
		if err != nil {
			panic(err)
		}
	})

	return http.ListenAndServe(addr, mux)
}
