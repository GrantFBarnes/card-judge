package apiDeck

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/database"
)

func Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
		return
	}

	var name string
	var password string
	for key, val := range r.Form {
		if key == "name" {
			name = val[0]
		} else if key == "password" {
			password = val[0]
		}
	}

	if name == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No name found.")
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get player id.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	id, err := database.CreateDeck(dbcs, playerId, name, password)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	err = database.AddPlayerDeckAccess(dbcs, playerId, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to add access.")
		return
	}

	w.Header().Add("HX-Redirect", "/deck/"+id.String())
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func Update(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	err = r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
		return
	}

	var name string
	var password string
	for key, val := range r.Form {
		if key == "name" {
			name = val[0]
		} else if key == "password" {
			password = val[0]
		}
	}

	if name == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No name found.")
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get player id.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	if !database.HasDeckAccess(dbcs, playerId, id) {
		api.WriteBadHeader(w, http.StatusUnauthorized, "Player does not have access.")
		return
	}

	err = database.UpdateDeck(dbcs, playerId, id, name, password)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get player id.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	if !database.HasDeckAccess(dbcs, playerId, id) {
		api.WriteBadHeader(w, http.StatusUnauthorized, "Player does not have access.")
		return
	}

	err = database.DeleteDeck(dbcs, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	w.Header().Add("HX-Redirect", "/decks")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}
