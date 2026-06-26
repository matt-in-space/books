package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.ListenAndServe(":3000", http.HandlerFunc(ping))
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}
