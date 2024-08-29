package apiLobby

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

func Access(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	lobby, err := database.GetLobby(dbcs, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get lobby from database.")
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

	err = auth.AddCookieAccessId(w, r, lobby.Id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to set cookie in browser.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
		return
	}

	var name string
	var password string
	for key, val := range r.Form {
		if key == "newLobbyName" {
			name = val[0]
		} else if key == "newLobbyPassword" {
			password = val[0]
		}
	}

	if name == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No name found.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	id, err := database.CreateLobby(dbcs, name, password)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to create lobby in database.")
		return
	}

	auth.AddCookieAccessId(w, r, id)

	w.Header().Add("HX-Redirect", "/lobby/"+id.String())
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
		if key == "newLobbyName" {
			name = val[0]
		} else if key == "newLobbyPassword" {
			password = val[0]
		}
	}

	if name == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No name found.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	err = database.UpdateLobby(dbcs, id, name, password)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update lobby in database.")
		return
	}

	auth.AddCookieAccessId(w, r, id)

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

	dbcs := database.GetDatabaseConnectionString()
	err = database.DeleteLobby(dbcs, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to delete lobby in database.")
		return
	}

	w.Header().Add("HX-Redirect", "/lobbies")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}
