package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/public-key", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("static-public-key-stub"))
	})
	fmt.Println("Node listening on :8080...")
	http.ListenAndServe(":8080", nil)
}
