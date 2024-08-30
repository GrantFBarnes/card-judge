package main

import (
	"fmt"
	"net/http"

	"github.com/grantfbarnes/card-judge/api"
	apiCard "github.com/grantfbarnes/card-judge/api/card"
	apiDeck "github.com/grantfbarnes/card-judge/api/deck"
	apiLobby "github.com/grantfbarnes/card-judge/api/lobby"
	apiPages "github.com/grantfbarnes/card-judge/api/pages"
	apiPlayer "github.com/grantfbarnes/card-judge/api/player"
)

func main() {
	// static files
	http.HandleFunc("GET /static/{fileType}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		fileType := r.PathValue("fileType")
		fileName := r.PathValue("fileName")
		http.ServeFile(w, r, "static/"+fileType+"/"+fileName)
	})

	// pages
	http.Handle("GET /", api.Middleware(http.HandlerFunc(apiPages.Home)))
	http.Handle("GET /login", api.Middleware(http.HandlerFunc(apiPages.Login)))
	http.Handle("GET /manage", api.Middleware(http.HandlerFunc(apiPages.Manage)))
	http.Handle("GET /lobbies", api.Middleware(http.HandlerFunc(apiPages.Lobbies)))
	http.Handle("GET /lobby/{id}", api.Middleware(http.HandlerFunc(apiPages.Lobby)))
	http.Handle("GET /decks", api.Middleware(http.HandlerFunc(apiPages.Decks)))
	http.Handle("GET /deck/{id}", api.Middleware(http.HandlerFunc(apiPages.Deck)))

	// player
	http.HandleFunc("POST /api/player/create", apiPlayer.Create)
	http.HandleFunc("PUT /api/player/{id}", apiPlayer.Update)
	http.HandleFunc("DELETE /api/player/{id}", apiPlayer.Delete)

	// lobby
	http.HandleFunc("POST /api/lobby/{id}/access", apiLobby.Access)
	http.HandleFunc("POST /api/lobby/create", apiLobby.Create)
	http.HandleFunc("PUT /api/lobby/{id}", apiLobby.Update)
	http.HandleFunc("DELETE /api/lobby/{id}", apiLobby.Delete)

	// deck
	http.HandleFunc("POST /api/deck/{id}/access", apiDeck.Access)
	http.HandleFunc("POST /api/deck/create", apiDeck.Create)
	http.HandleFunc("PUT /api/deck/{id}", apiDeck.Update)
	http.HandleFunc("DELETE /api/deck/{id}", apiDeck.Delete)

	// card
	http.HandleFunc("POST /api/card/create", apiCard.Create)
	http.HandleFunc("PUT /api/card/{id}", apiCard.Update)
	http.HandleFunc("DELETE /api/card/{id}", apiCard.Delete)

	port := ":8080"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
