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

	sqlFiles := []string{
		// database
		"../database/settings.sql",

		// tables
		"../database/tables/USER.sql",
		"../database/tables/DECK.sql",
		"../database/tables/CARD.sql",
		"../database/tables/LOBBY.sql",
		"../database/tables/DRAW_PILE.sql",
		"../database/tables/PLAYER.sql",
		"../database/tables/HAND.sql",
		"../database/tables/JUDGE.sql",
		"../database/tables/BOARD.sql",
		"../database/tables/WIN.sql",
		"../database/tables/USER_ACCESS_DECK.sql",
		"../database/tables/USER_ACCESS_LOBBY.sql",
		"../database/tables/LOGIN_ATTEMPT.sql",
		"../database/tables/LOG_DRAW.sql",
		"../database/tables/LOG_DISCARD.sql",
		"../database/tables/LOG_PLAY.sql",
		"../database/tables/LOG_SKIP.sql",
		"../database/tables/LOG_WIN.sql",

		// functions
		"../database/functions/FN_USER_HAS_DECK_ACCESS.sql",
		"../database/functions/FN_USER_HAS_LOBBY_ACCESS.sql",
		"../database/functions/FN_GET_RANDOM_LOBBY_PLAYER.sql",
		"../database/functions/FN_GET_RANDOM_PROMPT_CARD.sql",
		"../database/functions/FN_GET_LOGIN_ATTEMPT_IS_ALLOWED.sql",

		// procedures
		"../database/procedures/SP_GET_READABLE_DECKS.sql",
		"../database/procedures/SP_SET_PLAYER_ACTIVE.sql",
		"../database/procedures/SP_SET_PLAYER_INACTIVE.sql",
		"../database/procedures/SP_SET_JUDGE.sql",
		"../database/procedures/SP_SET_MISSING_JUDGE.sql",
		"../database/procedures/SP_DRAW_HAND.sql",
		"../database/procedures/SP_SKIP_PROMPT.sql",
		"../database/procedures/SP_PLAY_CARD.sql",
		"../database/procedures/SP_PLAY_STEAL_CARD.sql",
		"../database/procedures/SP_PLAY_SURPRISE_CARD.sql",
		"../database/procedures/SP_PLAY_WILD_CARD.sql",
		"../database/procedures/SP_WITHDRAWAL_CARD.sql",
		"../database/procedures/SP_DISCARD_HAND.sql",
		"../database/procedures/SP_DISCARD_CARD.sql",
		"../database/procedures/SP_PICK_WINNER.sql",

		// events
		"../database/events/EVT_CLEAN_LOGIN_ATTEMPTS.sql",

		// triggers
		"../database/triggers/TR_BF_UP_USER_SET_CHANGED_ON_DATE.sql",
		"../database/triggers/TR_BF_UP_DECK_SET_CHANGED_ON_DATE.sql",
		"../database/triggers/TR_BF_UP_CARD_SET_CHANGED_ON_DATE.sql",
		"../database/triggers/TR_BF_UP_LOBBY_SET_CHANGED_ON_DATE.sql",
		"../database/triggers/TR_BF_UP_PLAYER_SET_CHANGED_ON_DATE.sql",
		"../database/triggers/TR_BF_UP_HAND_SET_CHANGED_ON_DATE.sql",
		"../database/triggers/TR_BF_UP_JUDGE_SET_CHANGED_ON_DATE.sql",
		"../database/triggers/TR_AF_UP_DECK_REVOKE_ACCESS.sql",
		"../database/triggers/TR_AF_UP_LOBBY_REVOKE_ACCESS.sql",
		"../database/triggers/TR_AF_IN_PLAYER_DRAW_HAND.sql",
		"../database/triggers/TR_AF_IN_PLAYER_SET_MISSING_JUDGE.sql",
		"../database/triggers/TR_AF_UP_PLAYER_SET_MISSING_JUDGE.sql",
		"../database/triggers/TR_AF_DL_PLAYER_SET_MISSING_JUDGE.sql",
		"../database/triggers/TR_AF_IN_WIN_SET_JUDGE.sql",

		// populate
		"../database/populate/admin.sql",
		"../database/populate/deck.sql",
	}
	for _, sqlFile := range sqlFiles {
		err = database.RunFile(sqlFile)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

	// static files
	http.HandleFunc("GET /static/{fileType}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		fileType := r.PathValue("fileType")
		fileName := r.PathValue("fileName")
		http.ServeFile(w, r, "static/"+fileType+"/"+fileName)
	})

	// pages
	http.Handle("GET /", api.PageMiddleware(http.HandlerFunc(apiPages.Home)))
	http.Handle("GET /about", api.PageMiddleware(http.HandlerFunc(apiPages.About)))
	http.Handle("GET /login", api.PageMiddleware(http.HandlerFunc(apiPages.Login)))
	http.Handle("GET /manage", api.PageMiddleware(http.HandlerFunc(apiPages.Manage)))
	http.Handle("GET /admin", api.PageMiddleware(http.HandlerFunc(apiPages.Admin)))
	http.Handle("GET /stats", api.PageMiddleware(http.HandlerFunc(apiPages.Stats)))
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
	http.Handle("GET /api/deck/{deckId}/card-export", api.ApiMiddleware(http.HandlerFunc(apiDeck.GetCardExport)))
	http.Handle("POST /api/deck/search", api.ApiMiddleware(http.HandlerFunc(apiDeck.Search)))
	http.Handle("POST /api/deck/create", api.ApiMiddleware(http.HandlerFunc(apiDeck.Create)))
	http.Handle("PUT /api/deck/{deckId}/name", api.ApiMiddleware(http.HandlerFunc(apiDeck.SetName)))
	http.Handle("PUT /api/deck/{deckId}/password", api.ApiMiddleware(http.HandlerFunc(apiDeck.SetPassword)))
	http.Handle("PUT /api/deck/{deckId}/is-public-read-only", api.ApiMiddleware(http.HandlerFunc(apiDeck.SetIsPublicReadOnly)))
	http.Handle("DELETE /api/deck/{deckId}", api.ApiMiddleware(http.HandlerFunc(apiDeck.Delete)))

	// card
	http.Handle("POST /api/card/search", api.ApiMiddleware(http.HandlerFunc(apiCard.Search)))
	http.Handle("POST /api/card/create", api.ApiMiddleware(http.HandlerFunc(apiCard.Create)))
	http.Handle("PUT /api/card/{cardId}/category", api.ApiMiddleware(http.HandlerFunc(apiCard.SetCategory)))
	http.Handle("PUT /api/card/{cardId}/text", api.ApiMiddleware(http.HandlerFunc(apiCard.SetText)))
	http.Handle("DELETE /api/card/{cardId}", api.ApiMiddleware(http.HandlerFunc(apiCard.Delete)))

	// lobby
	http.Handle("GET /api/lobby/{lobbyId}/game-interface", api.ApiMiddleware(http.HandlerFunc(apiLobby.GetGameInterface)))
	http.Handle("POST /api/lobby/search", api.ApiMiddleware(http.HandlerFunc(apiLobby.Search)))
	http.Handle("POST /api/lobby/create", api.ApiMiddleware(http.HandlerFunc(apiLobby.Create)))
	http.Handle("POST /api/lobby/{lobbyId}/draw-hand", api.ApiMiddleware(http.HandlerFunc(apiLobby.DrawHand)))
	http.Handle("POST /api/lobby/{lobbyId}/card/{cardId}/play", api.ApiMiddleware(http.HandlerFunc(apiLobby.PlayCard)))
	http.Handle("POST /api/lobby/{lobbyId}/card/steal/play", api.ApiMiddleware(http.HandlerFunc(apiLobby.PlayStealCard)))
	http.Handle("POST /api/lobby/{lobbyId}/card/surprise/play", api.ApiMiddleware(http.HandlerFunc(apiLobby.PlaySurpriseCard)))
	http.Handle("POST /api/lobby/{lobbyId}/card/wild/play", api.ApiMiddleware(http.HandlerFunc(apiLobby.PlayWildCard)))
	http.Handle("POST /api/lobby/{lobbyId}/card/{cardId}/withdrawal", api.ApiMiddleware(http.HandlerFunc(apiLobby.WithdrawalCard)))
	http.Handle("POST /api/lobby/{lobbyId}/card/{cardId}/discard", api.ApiMiddleware(http.HandlerFunc(apiLobby.DiscardCard)))
	http.Handle("POST /api/lobby/{lobbyId}/card/{cardId}/lock", api.ApiMiddleware(http.HandlerFunc(apiLobby.LockCard)))
	http.Handle("POST /api/lobby/{lobbyId}/card/{cardId}/unlock", api.ApiMiddleware(http.HandlerFunc(apiLobby.UnlockCard)))
	http.Handle("POST /api/lobby/{lobbyId}/card/{cardId}/pick-winner", api.ApiMiddleware(http.HandlerFunc(apiLobby.PickWinner)))
	http.Handle("POST /api/lobby/{lobbyId}/discard-hand", api.ApiMiddleware(http.HandlerFunc(apiLobby.DiscardHand)))
	http.Handle("POST /api/lobby/{lobbyId}/flip", api.ApiMiddleware(http.HandlerFunc(apiLobby.FlipTable)))
	http.Handle("POST /api/lobby/{lobbyId}/skip-prompt", api.ApiMiddleware(http.HandlerFunc(apiLobby.SkipPrompt)))
	http.Handle("PUT /api/lobby/{lobbyId}/name", api.ApiMiddleware(http.HandlerFunc(apiLobby.SetName)))
	http.Handle("PUT /api/lobby/{lobbyId}/hand-size", api.ApiMiddleware(http.HandlerFunc(apiLobby.SetHandSize)))
	http.Handle("PUT /api/lobby/{lobbyId}/credit-limit", api.ApiMiddleware(http.HandlerFunc(apiLobby.SetCreditLimit)))

	// access
	http.Handle("POST /api/access/lobby/{lobbyId}", api.ApiMiddleware(http.HandlerFunc(apiAccess.Lobby)))
	http.Handle("POST /api/access/deck/{deckId}", api.ApiMiddleware(http.HandlerFunc(apiAccess.Deck)))

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
