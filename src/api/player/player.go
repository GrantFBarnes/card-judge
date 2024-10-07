package apiPlayer

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/database"
	"github.com/grantfbarnes/card-judge/websocket"
)

func GetGameInterfaceHtml(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

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

func DrawPlayerHand(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	err = database.DrawPlayerHand(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	GetGameInterfaceHtml(w, r)
}

func PlayPlayerCard(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	err = database.PlayPlayerCard(playerId, cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	GetGameInterfaceHtml(w, r)
}

func DiscardPlayerHand(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	err = database.DiscardPlayerHand(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	GetGameInterfaceHtml(w, r)
}

func LockPlayerCard(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	err = database.LockPlayerCard(playerId, cardId, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	GetGameInterfaceHtml(w, r)
}

func UnlockPlayerCard(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	err = database.LockPlayerCard(playerId, cardId, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	GetGameInterfaceHtml(w, r)
}

func DiscardPlayerCard(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	err = database.DiscardPlayerCard(playerId, cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	GetGameInterfaceHtml(w, r)
}

func FlipTable(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	player, err := database.GetPlayer(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	websocket.LobbyBroadcast(player.LobbyId, player.Name+": FLIP THE TABLE!")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
