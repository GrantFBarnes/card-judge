package main

import (
	"fmt"
	"net/http"

	"github.com/grantfbarnes/card-judge/api"
)

func main() {
	http.Handle("GET /", api.Middleware(http.HandlerFunc(api.PageHome)))
	http.Handle("GET /login", api.Middleware(http.HandlerFunc(api.PageLogin)))
	http.Handle("GET /lobbies", api.Middleware(http.HandlerFunc(api.PageLobbies)))
	http.Handle("GET /game/{lobbyid}", api.Middleware(http.HandlerFunc(api.PageGame)))
	http.Handle("GET /decks", api.Middleware(http.HandlerFunc(api.PageDecks)))
	http.Handle("GET /deck/{deckid}", api.Middleware(http.HandlerFunc(api.PageDeck)))

	http.HandleFunc("POST /api/access/deck/{deckid}", api.PostAccessDeck)
	http.HandleFunc("POST /api/access/lobby/{lobbyid}", api.PostAccessLobby)
	http.HandleFunc("POST /api/player/login", api.PostPlayerLogin)
	http.HandleFunc("POST /api/player/update", api.PostPlayerUpdate)
	http.HandleFunc("POST /api/player/logout", api.PostPlayerLogout)

	port := ":8080"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
