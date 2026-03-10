// app2/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from App2!")
	})

	log.Println("App2 running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}