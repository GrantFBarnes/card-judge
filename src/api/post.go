package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

func PostAccessLobby(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyid")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get lobby id"))
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	lobby, err := database.GetLobby(dbcs, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get lobby"))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse form"))
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

	if lobby.Password.String != password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("password not valid"))
		return
	}

	err = auth.AddAccessId(w, r, lobby.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to set cookie"))
		return
	}

	w.Header().Add("HX-Refresh", "true")
}

func PostAccessDeck(w http.ResponseWriter, r *http.Request) {
	deckIdString := r.PathValue("deckid")
	deckId, err := uuid.Parse(deckIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get deck id"))
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	deck, err := database.GetDeck(dbcs, deckId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get deck"))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse form"))
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

	if deck.Password.String != password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("password not valid"))
		return
	}

	err = auth.AddAccessId(w, r, deck.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to set cookie"))
		return
	}

	w.Header().Add("HX-Refresh", "true")
}

func PostPlayerLogin(w http.ResponseWriter, r *http.Request) {
	err := setPlayerName(w, r)
	if err != nil {
		return
	}
	w.Header().Add("HX-Redirect", auth.GetRedirectURL(r))
}

func PostPlayerUpdate(w http.ResponseWriter, r *http.Request) {
	err := setPlayerName(w, r)
	if err != nil {
		return
	}
	w.Header().Add("HX-Refresh", "true")
}

func PostPlayerLogout(w http.ResponseWriter, r *http.Request) {
	auth.RemovePlayerName(w)
	w.Header().Add("HX-Refresh", "true")
}

func setPlayerName(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse form"))
		return err
	}

	var playerName string
	for key, val := range r.Form {
		if key != "playerName" {
			continue
		}
		playerName = val[0]
		break
	}

	if playerName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no name found"))
		return errors.New("no name found")
	}

	err = auth.SetPlayerName(w, playerName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to set cookie"))
		return errors.New("failed to set cookie")
	}

	return nil
}
