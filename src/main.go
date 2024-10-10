package main

import (
	"log"
	"net/http"
	"os"

	"github.com/grantfbarnes/card-judge/api"
	apiAccess "github.com/grantfbarnes/card-judge/api/access"
	apiCard "github.com/grantfbarnes/card-judge/api/card"
	apiDeck "github.com/grantfbarnes/card-judge/api/deck"
	apiLobby "github.com/grantfbarnes/card-judge/api/lobby"
	apiPages "github.com/grantfbarnes/card-judge/api/pages"
	apiPlayer "github.com/grantfbarnes/card-judge/api/player"
	apiUser "github.com/grantfbarnes/card-judge/api/user"
	"github.com/grantfbarnes/card-judge/database"
	"github.com/grantfbarnes/card-judge/websocket"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	db, err := database.CreateDatabaseConnection()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()

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
	http.Handle("GET /lobby/{lobbyId}", api.PageMiddleware(http.HandlerFunc(apiPages.Lobby)))
	http.Handle("GET /lobby/{lobbyId}/access", api.PageMiddleware(http.HandlerFunc(apiPages.LobbyAccess)))
	http.Handle("GET /decks", api.PageMiddleware(http.HandlerFunc(apiPages.Decks)))
	http.Handle("GET /deck/{deckId}", api.PageMiddleware(http.HandlerFunc(apiPages.Deck)))
	http.Handle("GET /deck/{deckId}/access", api.PageMiddleware(http.HandlerFunc(apiPages.DeckAccess)))

	// user
	http.Handle("POST /api/user/search", api.ApiMiddleware(http.HandlerFunc(apiUser.Search)))
	http.Handle("POST /api/user/create", api.ApiMiddleware(http.HandlerFunc(apiUser.Create)))
	http.Handle("POST /api/user/create/default", api.ApiMiddleware(http.HandlerFunc(apiUser.CreateDefault)))
	http.Handle("POST /api/user/login", api.ApiMiddleware(http.HandlerFunc(apiUser.Login)))
	http.Handle("POST /api/user/logout", api.ApiMiddleware(http.HandlerFunc(apiUser.Logout)))
	http.Handle("PUT /api/user/{userId}/name", api.ApiMiddleware(http.HandlerFunc(apiUser.SetName)))
	http.Handle("PUT /api/user/{userId}/password", api.ApiMiddleware(http.HandlerFunc(apiUser.SetPassword)))
	http.Handle("PUT /api/user/{userId}/password/reset", api.ApiMiddleware(http.HandlerFunc(apiUser.ResetPassword)))
	http.Handle("PUT /api/user/{userId}/color-theme", api.ApiMiddleware(http.HandlerFunc(apiUser.SetColorTheme)))
	http.Handle("PUT /api/user/{userId}/is-admin", api.ApiMiddleware(http.HandlerFunc(apiUser.SetIsAdmin)))
	http.Handle("DELETE /api/user/{userId}", api.ApiMiddleware(http.HandlerFunc(apiUser.Delete)))

	// deck
	http.Handle("POST /api/deck/search", api.ApiMiddleware(http.HandlerFunc(apiDeck.Search)))
	http.Handle("POST /api/deck/create", api.ApiMiddleware(http.HandlerFunc(apiDeck.Create)))
	http.Handle("PUT /api/deck/{deckId}/name", api.ApiMiddleware(http.HandlerFunc(apiDeck.SetName)))
	http.Handle("PUT /api/deck/{deckId}/password", api.ApiMiddleware(http.HandlerFunc(apiDeck.SetPassword)))
	http.Handle("DELETE /api/deck/{deckId}", api.ApiMiddleware(http.HandlerFunc(apiDeck.Delete)))

	// card
	http.Handle("POST /api/card/search", api.ApiMiddleware(http.HandlerFunc(apiCard.Search)))
	http.Handle("POST /api/card/create", api.ApiMiddleware(http.HandlerFunc(apiCard.Create)))
	http.Handle("PUT /api/card/{cardId}/category", api.ApiMiddleware(http.HandlerFunc(apiCard.SetCategory)))
	http.Handle("PUT /api/card/{cardId}/text", api.ApiMiddleware(http.HandlerFunc(apiCard.SetText)))
	http.Handle("DELETE /api/card/{cardId}", api.ApiMiddleware(http.HandlerFunc(apiCard.Delete)))

	// lobby
	http.Handle("POST /api/lobby/{lobbyId}/skip-judge-card", api.ApiMiddleware(http.HandlerFunc(apiLobby.SkipJudgeCard)))
	http.Handle("POST /api/lobby/{lobbyId}/card/{cardId}/winner", api.ApiMiddleware(http.HandlerFunc(apiLobby.PickLobbyWinner)))
	http.Handle("POST /api/lobby/search", api.ApiMiddleware(http.HandlerFunc(apiLobby.Search)))
	http.Handle("POST /api/lobby/create", api.ApiMiddleware(http.HandlerFunc(apiLobby.Create)))
	http.Handle("PUT /api/lobby/{lobbyId}/name", api.ApiMiddleware(http.HandlerFunc(apiLobby.SetName)))
	http.Handle("PUT /api/lobby/{lobbyId}/hand-size", api.ApiMiddleware(http.HandlerFunc(apiLobby.SetHandSize)))

	// access
	http.Handle("POST /api/access/lobby/{lobbyId}", api.ApiMiddleware(http.HandlerFunc(apiAccess.Lobby)))
	http.Handle("POST /api/access/deck/{deckId}", api.ApiMiddleware(http.HandlerFunc(apiAccess.Deck)))

	// player
	http.Handle("GET /api/player/{playerId}/game-interface", api.ApiMiddleware(http.HandlerFunc(apiPlayer.GetGameInterfaceHtml)))
	http.Handle("POST /api/player/{playerId}/draw", api.ApiMiddleware(http.HandlerFunc(apiPlayer.DrawPlayerHand)))
	http.Handle("POST /api/player/{playerId}/card/{cardId}/play", api.ApiMiddleware(http.HandlerFunc(apiPlayer.PlayPlayerCard)))
	http.Handle("POST /api/player/{playerId}/card/{cardId}/withdrawal", api.ApiMiddleware(http.HandlerFunc(apiPlayer.WithdrawalPlayerCard)))
	http.Handle("POST /api/player/{playerId}/discard", api.ApiMiddleware(http.HandlerFunc(apiPlayer.DiscardPlayerHand)))
	http.Handle("POST /api/player/{playerId}/flip", api.ApiMiddleware(http.HandlerFunc(apiPlayer.FlipTable)))
	http.Handle("PUT /api/player/{playerId}/card/{cardId}/lock", api.ApiMiddleware(http.HandlerFunc(apiPlayer.LockPlayerCard)))
	http.Handle("PUT /api/player/{playerId}/card/{cardId}/unlock", api.ApiMiddleware(http.HandlerFunc(apiPlayer.UnlockPlayerCard)))
	http.Handle("DELETE /api/player/{playerId}/card/{cardId}", api.ApiMiddleware(http.HandlerFunc(apiPlayer.DiscardPlayerCard)))

	// websocket
	http.HandleFunc("GET /ws/lobby/{lobbyId}", websocket.ServeWs)

	if os.Getenv("CARD_JUDGE_ENV") == "PROD" {
		logFileName := os.Getenv("CARD_JUDGE_LOG_FILE")
		certFileName := os.Getenv("CARD_JUDGE_CERT_FILE")
		keyFileName := os.Getenv("CARD_JUDGE_KEY_FILE")

		logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)

		port := ":443"
		log.Println("server is running...")
		err = http.ListenAndServeTLS(port, certFileName, keyFileName, nil)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		port := ":8080"
		log.Println("server is running...")
		err = http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
