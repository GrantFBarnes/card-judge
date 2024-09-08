package apiAccess

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

func Lobby(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	lobbyPasswordHash, err := database.GetLobbyPasswordHash(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var password string
	for key, val := range r.Form {
		if key != "password" {
			continue
		}
		password = val[0]
		break
	}

	if !auth.PasswordMatchesHash(password, lobbyPasswordHash.String) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Provided password is not valid."))
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
		return
	}

	err = database.AddPlayerLobbyAccess(playerId, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to add access."))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func Deck(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	deckPasswordHash, err := database.GetDeckPasswordHash(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var password string
	for key, val := range r.Form {
		if key != "password" {
			continue
		}
		password = val[0]
		break
	}

	if !auth.PasswordMatchesHash(password, deckPasswordHash.String) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Provided password is not valid."))
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
		return
	}

	err = database.AddPlayerDeckAccess(playerId, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to add access."))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
