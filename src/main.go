package main

import (
	"log"
	"net/http"

	"github.com/grantfbarnes/card-judge/api"
	apiAccess "github.com/grantfbarnes/card-judge/api/access"
	apiCard "github.com/grantfbarnes/card-judge/api/card"
	apiDeck "github.com/grantfbarnes/card-judge/api/deck"
	apiLobby "github.com/grantfbarnes/card-judge/api/lobby"
	apiPages "github.com/grantfbarnes/card-judge/api/pages"
	apiPlayer "github.com/grantfbarnes/card-judge/api/player"
	"github.com/grantfbarnes/card-judge/database"
	"github.com/grantfbarnes/card-judge/websocket"
)

func main() {
	database.SetDatabaseConnectionString()
	err := database.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	// static files
	http.HandleFunc("GET /static/{fileType}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		fileType := r.PathValue("fileType")
		fileName := r.PathValue("fileName")
		http.ServeFile(w, r, "static/"+fileType+"/"+fileName)
	})

	// pages
	http.Handle("GET /", api.PageMiddleware(http.HandlerFunc(apiPages.Home)))
	http.Handle("GET /login", api.PageMiddleware(http.HandlerFunc(apiPages.Login)))
	http.Handle("GET /manage", api.PageMiddleware(http.HandlerFunc(apiPages.Manage)))
	http.Handle("GET /admin", api.PageMiddleware(http.HandlerFunc(apiPages.Admin)))
	http.Handle("GET /lobbies", api.PageMiddleware(http.HandlerFunc(apiPages.Lobbies)))
	http.Handle("GET /lobby/{id}", api.PageMiddleware(http.HandlerFunc(apiPages.Lobby)))
	http.Handle("GET /decks", api.PageMiddleware(http.HandlerFunc(apiPages.Decks)))
	http.Handle("GET /deck/{id}", api.PageMiddleware(http.HandlerFunc(apiPages.Deck)))

	// player
	http.Handle("POST /api/player/search", api.ApiMiddleware(http.HandlerFunc(apiPlayer.Search)))
	http.Handle("POST /api/player/create", api.ApiMiddleware(http.HandlerFunc(apiPlayer.Create)))
	http.Handle("POST /api/player/create/default", api.ApiMiddleware(http.HandlerFunc(apiPlayer.CreateDefault)))
	http.Handle("POST /api/player/login", api.ApiMiddleware(http.HandlerFunc(apiPlayer.Login)))
	http.Handle("POST /api/player/logout", api.ApiMiddleware(http.HandlerFunc(apiPlayer.Logout)))
	http.Handle("PUT /api/player/{id}/name", api.ApiMiddleware(http.HandlerFunc(apiPlayer.SetName)))
	http.Handle("PUT /api/player/{id}/password", api.ApiMiddleware(http.HandlerFunc(apiPlayer.SetPassword)))
	http.Handle("PUT /api/player/{id}/password/reset", api.ApiMiddleware(http.HandlerFunc(apiPlayer.ResetPassword)))
	http.Handle("PUT /api/player/{id}/color-theme", api.ApiMiddleware(http.HandlerFunc(apiPlayer.SetColorTheme)))
	http.Handle("PUT /api/player/{id}/is-admin", api.ApiMiddleware(http.HandlerFunc(apiPlayer.SetIsAdmin)))
	http.Handle("DELETE /api/player/{id}", api.ApiMiddleware(http.HandlerFunc(apiPlayer.Delete)))

	// lobby
	http.Handle("GET /api/lobby/{id}/players", api.ApiMiddleware(http.HandlerFunc(apiLobby.GetPlayers)))
	http.Handle("GET /api/lobby/{id}/card-count", api.ApiMiddleware(http.HandlerFunc(apiLobby.GetCardCount)))
	http.Handle("POST /api/lobby/{lobbyId}/player/{playerId}/draw", api.ApiMiddleware(http.HandlerFunc(apiLobby.DrawPlayerHand)))
	http.Handle("POST /api/lobby/search", api.ApiMiddleware(http.HandlerFunc(apiLobby.Search)))
	http.Handle("POST /api/lobby/create", api.ApiMiddleware(http.HandlerFunc(apiLobby.Create)))
	http.Handle("PUT /api/lobby/{id}/name", api.ApiMiddleware(http.HandlerFunc(apiLobby.SetName)))
	http.Handle("PUT /api/lobby/{id}/password", api.ApiMiddleware(http.HandlerFunc(apiLobby.SetPassword)))

	// deck
	http.Handle("POST /api/deck/search", api.ApiMiddleware(http.HandlerFunc(apiDeck.Search)))
	http.Handle("POST /api/deck/create", api.ApiMiddleware(http.HandlerFunc(apiDeck.Create)))
	http.Handle("PUT /api/deck/{id}/name", api.ApiMiddleware(http.HandlerFunc(apiDeck.SetName)))
	http.Handle("PUT /api/deck/{id}/password", api.ApiMiddleware(http.HandlerFunc(apiDeck.SetPassword)))
	http.Handle("DELETE /api/deck/{id}", api.ApiMiddleware(http.HandlerFunc(apiDeck.Delete)))

	// card
	http.Handle("POST /api/card/search", api.ApiMiddleware(http.HandlerFunc(apiCard.Search)))
	http.Handle("POST /api/card/create", api.ApiMiddleware(http.HandlerFunc(apiCard.Create)))
	http.Handle("PUT /api/card/{id}/type", api.ApiMiddleware(http.HandlerFunc(apiCard.SetType)))
	http.Handle("PUT /api/card/{id}/text", api.ApiMiddleware(http.HandlerFunc(apiCard.SetText)))
	http.Handle("DELETE /api/card/{id}", api.ApiMiddleware(http.HandlerFunc(apiCard.Delete)))

	// access
	http.Handle("POST /api/access/lobby/{id}", api.ApiMiddleware(http.HandlerFunc(apiAccess.Lobby)))
	http.Handle("POST /api/access/deck/{id}", api.ApiMiddleware(http.HandlerFunc(apiAccess.Deck)))

	// websocket
	http.HandleFunc("GET /ws/lobby/{id}", websocket.ServeWs)

	port := ":8080"
	log.Printf("running at http://localhost%s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
