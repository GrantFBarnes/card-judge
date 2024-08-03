package main

import (
	"fmt"
	"net/http"

	"github.com/grantfbarnes/card-judge/api"
)

func main() {
	http.Handle("GET /", api.Middleware(http.HandlerFunc(api.PageHome)))
	http.Handle("GET /login", api.Middleware(http.HandlerFunc(api.PageLogin)))
	http.Handle("GET /lobby/join", api.Middleware(http.HandlerFunc(api.PageLobbyJoin)))
	http.Handle("GET /cards", api.Middleware(http.HandlerFunc(api.PageCardList)))

	http.HandleFunc("POST /api/login", api.PostLogin)

	port := ":8090"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
