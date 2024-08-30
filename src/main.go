package main

import (
	"fmt"
	"net/http"

	apiAccess "github.com/grantfbarnes/card-judge/api/access"
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
	http.Handle("GET /", apiPages.Middleware(http.HandlerFunc(apiPages.Home)))
	http.Handle("GET /login", apiPages.Middleware(http.HandlerFunc(apiPages.Login)))
	http.Handle("GET /manage", apiPages.Middleware(http.HandlerFunc(apiPages.Manage)))
	http.Handle("GET /lobbies", apiPages.Middleware(http.HandlerFunc(apiPages.Lobbies)))
	http.Handle("GET /lobby/{id}", apiPages.Middleware(http.HandlerFunc(apiPages.Lobby)))
	http.Handle("GET /decks", apiPages.Middleware(http.HandlerFunc(apiPages.Decks)))
	http.Handle("GET /deck/{id}", apiPages.Middleware(http.HandlerFunc(apiPages.Deck)))

	// player
	http.HandleFunc("POST /api/player/create", apiPlayer.Create)
	http.HandleFunc("POST /api/player/login", apiPlayer.Login)
	http.HandleFunc("POST /api/player/logout", apiPlayer.Logout)
	http.HandleFunc("PUT /api/player/{id}/name", apiPlayer.SetName)
	http.HandleFunc("PUT /api/player/{id}/password", apiPlayer.SetPassword)
	http.HandleFunc("PUT /api/player/{id}/color-theme", apiPlayer.SetColorTheme)
	http.HandleFunc("DELETE /api/player/{id}", apiPlayer.Delete)

	// lobby
	http.HandleFunc("POST /api/lobby/create", apiLobby.Create)
	http.HandleFunc("PUT /api/lobby/{id}", apiLobby.Update)
	http.HandleFunc("DELETE /api/lobby/{id}", apiLobby.Delete)

	// deck
	http.HandleFunc("POST /api/deck/create", apiDeck.Create)
	http.HandleFunc("PUT /api/deck/{id}", apiDeck.Update)
	http.HandleFunc("DELETE /api/deck/{id}", apiDeck.Delete)

	// card
	http.HandleFunc("POST /api/card/create", apiCard.Create)
	http.HandleFunc("PUT /api/card/{id}", apiCard.Update)
	http.HandleFunc("DELETE /api/card/{id}", apiCard.Delete)

	// access
	http.HandleFunc("POST /api/access/lobby/{id}", apiAccess.Lobby)
	http.HandleFunc("POST /api/access/deck/{id}", apiAccess.Deck)

	port := ":8080"
	fmt.Printf("running at http://localhost%s\n", port)
	http.ListenAndServe(port, nil)
}
