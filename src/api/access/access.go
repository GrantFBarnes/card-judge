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
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	lobby, err := database.GetLobby(dbcs, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	err = r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
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

	if !auth.PasswordMatchesHash(password, lobby.PasswordHash.String) {
		api.WriteBadHeader(w, http.StatusBadRequest, "Provided password is not valid.")
		return
	}

	playerId, err := auth.GetCookiePlayerId(r)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get player id.")
		return
	}

	err = database.AddPlayerLobbyAccess(dbcs, playerId, lobby.Id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to add access.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func Deck(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	deck, err := database.GetDeck(dbcs, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	err = r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
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

	if !auth.PasswordMatchesHash(password, deck.PasswordHash.String) {
		api.WriteBadHeader(w, http.StatusBadRequest, "Provided password is not valid.")
		return
	}

	playerId, err := auth.GetCookiePlayerId(r)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get player id.")
		return
	}

	err = database.AddPlayerDeckAccess(dbcs, playerId, deck.Id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to add access.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}
