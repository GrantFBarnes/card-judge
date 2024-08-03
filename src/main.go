package main

import (
	"fmt"
	"net/http"

	"github.com/grantfbarnes/card-judge/api"
)

func main() {
	http.HandleFunc("GET /", api.PageHome)
	http.HandleFunc("GET /login", api.PageLogin)
	http.HandleFunc("GET /lobby/join", api.PageLobbyJoin)
	http.HandleFunc("GET /cards", api.PageCardList)

	http.HandleFunc("POST /api/login", api.PostLogin)

	port := ":8090"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
