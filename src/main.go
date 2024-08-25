package main

import (
	"fmt"
	"net/http"

	api "github.com/grantfbarnes/card-judge/api"
	apiDeck "github.com/grantfbarnes/card-judge/api/deck"
	apiLobby "github.com/grantfbarnes/card-judge/api/lobby"
	apiPages "github.com/grantfbarnes/card-judge/api/pages"
	apiPlayer "github.com/grantfbarnes/card-judge/api/player"
)

func main() {
	// pages
	http.Handle("GET /", api.Middleware(http.HandlerFunc(apiPages.Home)))
	http.Handle("GET /login", api.Middleware(http.HandlerFunc(apiPages.Login)))
	http.Handle("GET /lobbies", api.Middleware(http.HandlerFunc(apiPages.Lobbies)))
	http.Handle("GET /lobby/{lobbyid}", api.Middleware(http.HandlerFunc(apiPages.Lobby)))
	http.Handle("GET /decks", api.Middleware(http.HandlerFunc(apiPages.Decks)))
	http.Handle("GET /deck/{deckid}", api.Middleware(http.HandlerFunc(apiPages.Deck)))

	// player
	http.HandleFunc("POST /api/player/login", apiPlayer.Login)
	http.HandleFunc("POST /api/player/update", apiPlayer.Update)
	http.HandleFunc("POST /api/player/logout", apiPlayer.Logout)

	// lobby
	http.HandleFunc("POST /api/lobby/{lobbyid}/access", apiLobby.Access)
	http.HandleFunc("POST /api/lobby/create", apiLobby.Create)

	// deck
	http.HandleFunc("POST /api/deck/{deckid}/access", apiDeck.Access)
	http.HandleFunc("POST /api/deck/create", apiDeck.Create)

	port := ":8080"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
