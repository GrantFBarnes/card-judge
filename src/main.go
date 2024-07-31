package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/headers", headers)

	port := ":8090"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "home\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
