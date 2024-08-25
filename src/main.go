package main

import (
	"fmt"
	"net/http"

	api "github.com/grantfbarnes/card-judge/api"
	apiPages "github.com/grantfbarnes/card-judge/api/pages"
)

func main() {
	http.Handle("GET /", api.Middleware(http.HandlerFunc(apiPages.Home)))
	http.Handle("GET /login", api.Middleware(http.HandlerFunc(apiPages.Login)))
	http.Handle("GET /lobbies", api.Middleware(http.HandlerFunc(apiPages.Lobbies)))
	http.Handle("GET /lobby/{lobbyid}", api.Middleware(http.HandlerFunc(apiPages.Lobby)))
	http.Handle("GET /decks", api.Middleware(http.HandlerFunc(apiPages.Decks)))
	http.Handle("GET /deck/{deckid}", api.Middleware(http.HandlerFunc(apiPages.Deck)))

	http.HandleFunc("POST /api/access/deck/{deckid}", api.PostAccessDeck)
	http.HandleFunc("POST /api/access/lobby/{lobbyid}", api.PostAccessLobby)
	http.HandleFunc("POST /api/player/login", api.PostPlayerLogin)
	http.HandleFunc("POST /api/player/update", api.PostPlayerUpdate)
	http.HandleFunc("POST /api/player/logout", api.PostPlayerLogout)
	http.HandleFunc("POST /api/lobby/create", api.PostLobbyCreate)
	http.HandleFunc("POST /api/deck/create", api.PostDeckCreate)

	port := ":8080"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
