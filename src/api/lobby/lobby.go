package apiLobby

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/database"
	"github.com/grantfbarnes/card-judge/websocket"
)

func GetGameInterface(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	player, err := getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	writeGameInterfaceHtml(w, player.Id)
}

func Search(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var search string
	for key, val := range r.Form {
		if key == "search" {
			search = val[0]
		}
	}

	search = "%" + search + "%"

	lobbies, err := database.SearchLobbies(search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/table-rows/lobby-table-rows.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "lobby-table-rows", lobbies)
}

func Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var name string
	var password string
	var passwordConfirm string
	var handSize int
	var deckIds []uuid.UUID = make([]uuid.UUID, 0)
	for key, val := range r.Form {
		if key == "name" {
			name = val[0]
		} else if key == "password" {
			password = val[0]
		} else if key == "passwordConfirm" {
			passwordConfirm = val[0]
		} else if key == "handSize" {
			handSize, err = strconv.Atoi(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse hand size."))
				return
			}
		} else if strings.HasPrefix(key, "deckId") {
			deckId, err := uuid.Parse(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse deck id."))
				return
			}
			deckIds = append(deckIds, deckId)
		}
	}

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No name found."))
		return
	}

	if password != "" {
		if password != passwordConfirm {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Passwords do not match."))
			return
		}
	}

	if handSize <= 0 {
		handSize = 1
	}

	if handSize > 16 {
		handSize = 16
	}

	if len(deckIds) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("At least one deck is required."))
		return
	}

	userId := api.GetUserId(r)
	if userId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get user id."))
		return
	}

	existingLobbyId, err := database.GetLobbyId(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if existingLobbyId != uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Lobby name already exists."))
		return
	}

	lobbyId, err := database.CreateLobby(name, password, handSize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.AddCardsToLobby(lobbyId, deckIds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.AddUserLobbyAccess(userId, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Redirect", "/lobby/"+lobbyId.String())
	w.WriteHeader(http.StatusCreated)
}

func DrawPlayerHand(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	player, err := getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.DrawPlayerHand(player.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	writeGameInterfaceHtml(w, player.Id)
}

func PlayPlayerCard(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	player, err := getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.PlayPlayerCard(player.Id, cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	websocket.LobbyBroadcast(lobbyId, "refresh")
	w.WriteHeader(http.StatusOK)
}

func WithdrawalPlayerCard(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	player, err := getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.WithdrawalPlayerCard(player.Id, cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	websocket.LobbyBroadcast(lobbyId, "refresh")
	w.WriteHeader(http.StatusOK)
}

func DiscardPlayerCard(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	player, err := getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.DiscardPlayerCard(player.Id, cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	writeGameInterfaceHtml(w, player.Id)
}

func LockPlayerCard(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	player, err := getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.LockPlayerCard(player.Id, cardId, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	writeGameInterfaceHtml(w, player.Id)
}

func UnlockPlayerCard(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	player, err := getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.LockPlayerCard(player.Id, cardId, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	writeGameInterfaceHtml(w, player.Id)
}

func PickWinner(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	_, err = getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	winnerName, err := database.PickLobbyWinner(lobbyId, cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	cardTextStart, err := database.GetCardTextStart(cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	websocket.LobbyBroadcast(lobbyId, "Winner: "+winnerName+"\nCard Played: "+cardTextStart)
	websocket.LobbyBroadcast(lobbyId, "refresh")
	w.WriteHeader(http.StatusOK)
}

func DiscardPlayerHand(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	player, err := getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.DiscardPlayerHand(player.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	writeGameInterfaceHtml(w, player.Id)
}

func FlipTable(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	player, err := getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	websocket.LobbyBroadcast(lobbyId, player.Name+": FLIP THE TABLE!")

	w.Header().Add("HX-Redirect", "/lobbies")
	w.WriteHeader(http.StatusOK)
}

func SkipJudgeCard(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	_, err = getLobbyRequestPlayer(r, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.SkipJudgeCard(lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	websocket.LobbyBroadcast(lobbyId, "refresh")
	w.WriteHeader(http.StatusOK)
}

func SetName(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	userId := api.GetUserId(r)
	if userId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get user id."))
		return
	}

	if !database.UserHasLobbyAccess(userId, lobbyId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User does not have access."))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var name string
	for key, val := range r.Form {
		if key == "name" {
			name = val[0]
		}
	}

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No name found."))
		return
	}

	existingLobbyId, err := database.GetLobbyId(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if existingLobbyId != uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Lobby name already exists."))
		return
	}

	err = database.SetLobbyName(lobbyId, name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	websocket.LobbyBroadcast(lobbyId, fmt.Sprintf("Lobby name set to %s...", name))
	websocket.LobbyBroadcast(lobbyId, "refresh")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func SetHandSize(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get lobby id from path."))
		return
	}

	userId := api.GetUserId(r)
	if userId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get user id."))
		return
	}

	if !database.UserHasLobbyAccess(userId, lobbyId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User does not have access."))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var handSize int
	for key, val := range r.Form {
		if key == "handSize" {
			handSize, err = strconv.Atoi(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse hand size."))
				return
			}
		}
	}

	if handSize <= 0 {
		handSize = 1
	}

	if handSize > 16 {
		handSize = 16
	}

	err = database.SetLobbyHandSize(lobbyId, handSize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	websocket.LobbyBroadcast(lobbyId, fmt.Sprintf("Lobby hand size set to %d...", handSize))
	websocket.LobbyBroadcast(lobbyId, "refresh")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func getLobbyRequestPlayer(r *http.Request, lobbyId uuid.UUID) (database.Player, error) {
	var player database.Player

	userId := api.GetUserId(r)
	if userId == uuid.Nil {
		return player, errors.New("failed to get user id")
	}

	player, err := database.GetPlayer(lobbyId, userId)
	if err != nil {
		return player, err
	}

	if player.Id == uuid.Nil {
		return player, errors.New("user not found in lobby")
	}

	return player, nil
}

func writeGameInterfaceHtml(w http.ResponseWriter, playerId uuid.UUID) {
	gameData, err := database.GetPlayerGameData(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/game/game-interface.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "game-interface", gameData)
}
