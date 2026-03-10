// app1/main.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://app2:8080") // call app2 service
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error connecting to app2: %v", err)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "Hello from App1! App2 says: %s", string(body))
	})

	log.Println("App1 running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}