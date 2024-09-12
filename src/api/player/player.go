package apiPlayer

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/database"
)

func GetPlayerData(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	playerData, err := database.GetPlayerData(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/game/player-data.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "player-data", playerData)
}

func GetGameBoard(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	data, err := database.GetPlayerGameBoard(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/game/game-board.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "game-board", data)
}

func BecomeJudge(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	playerData, err := database.PlayerBecomeJudge(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/game/player-data.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "player-data", playerData)
}

func DrawPlayerHand(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	playerData, err := database.DrawPlayerHand(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/game/player-data.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "player-data", playerData)
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

	playerData, err := database.PlayPlayerCard(playerId, cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/game/player-data.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "player-data", playerData)
}

func DiscardPlayerHand(w http.ResponseWriter, r *http.Request) {
	playerIdString := r.PathValue("playerId")
	playerId, err := uuid.Parse(playerIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id from path."))
		return
	}

	playerData, err := database.DiscardPlayerHand(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/game/player-data.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "player-data", playerData)
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

	playerData, err := database.DiscardPlayerCard(playerId, cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/game/player-data.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "player-data", playerData)
}
