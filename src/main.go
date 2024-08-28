package main

import (
	"fmt"
	"net/http"

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
	http.HandleFunc("GET /", apiPages.Home)
	http.HandleFunc("GET /lobbies", apiPages.Lobbies)
	http.HandleFunc("GET /lobby/{id}", apiPages.Lobby)
	http.HandleFunc("GET /decks", apiPages.Decks)
	http.HandleFunc("GET /deck/{id}", apiPages.Deck)

	// player
	http.HandleFunc("POST /api/player/login", apiPlayer.Login)
	http.HandleFunc("POST /api/player/update", apiPlayer.Update)
	http.HandleFunc("POST /api/player/logout", apiPlayer.Logout)

	// lobby
	http.HandleFunc("POST /api/lobby/{id}/access", apiLobby.Access)
	http.HandleFunc("POST /api/lobby/create", apiLobby.Create)
	http.HandleFunc("DELETE /api/lobby/{id}", apiLobby.Delete)

	// deck
	http.HandleFunc("POST /api/deck/{id}/access", apiDeck.Access)
	http.HandleFunc("POST /api/deck/create", apiDeck.Create)
	http.HandleFunc("DELETE /api/deck/{id}", apiDeck.Delete)

	// card
	http.HandleFunc("POST /api/card/create", apiCard.Create)
	http.HandleFunc("PUT /api/card/{id}", apiCard.Update)
	http.HandleFunc("DELETE /api/card/{id}", apiCard.Delete)

	port := ":8080"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
