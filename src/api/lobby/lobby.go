package apiLobby

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/database"
)

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
	for key, val := range r.Form {
		if key == "name" {
			name = val[0]
		} else if key == "password" {
			password = val[0]
		} else if key == "passwordConfirm" {
			passwordConfirm = val[0]
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

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
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

	id, err := database.CreateLobby(playerId, name, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = database.AddPlayerLobbyAccess(playerId, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Redirect", "/lobby/"+id.String())
	w.WriteHeader(http.StatusCreated)
}

func SetName(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
		return
	}

	if !database.HasLobbyAccess(playerId, id) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Player does not have access."))
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

	err = database.SetLobbyName(playerId, id, name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func SetPassword(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
		return
	}

	if !database.HasLobbyAccess(playerId, id) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Player does not have access."))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var password string
	var passwordConfirm string
	for key, val := range r.Form {
		if key == "password" {
			password = val[0]
		} else if key == "passwordConfirm" {
			passwordConfirm = val[0]
		}
	}

	if password != "" {
		if password != passwordConfirm {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Passwords do not match."))
			return
		}
	}

	err = database.SetLobbyPassword(playerId, id, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
		return
	}

	if !database.HasLobbyAccess(playerId, id) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Player does not have access."))
		return
	}

	err = database.DeleteLobby(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to update the database."))
		return
	}

	w.Header().Add("HX-Redirect", "/lobbies")
	w.WriteHeader(http.StatusOK)
}
