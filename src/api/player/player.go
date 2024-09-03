package apiPlayer

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/auth"
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
		api.WriteBadHeader(w, http.StatusBadRequest, "No name found.")
		return
	}

	if password == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No password found.")
		return
	}

	if password != passwordConfirm {
		api.WriteBadHeader(w, http.StatusBadRequest, "Passwords do not match.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	id, err := database.CreatePlayer(dbcs, name, password)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	err = auth.SetCookiePlayerId(w, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to create player cookie in browser.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func Login(w http.ResponseWriter, r *http.Request) {
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

	if password == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No password found.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	id, err := database.GetPlayerId(dbcs, name, password)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to login.")
		return
	}

	err = auth.SetCookiePlayerId(w, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to create player cookie in browser.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	auth.RemoveCookiePlayerId(w)
	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func SetName(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	if !isPlayerOrAdmin(r, id) {
		api.WriteBadHeader(w, http.StatusUnauthorized, "Player does not have access.")
		return
	}

	err = r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
		return
	}

	var name string
	for key, val := range r.Form {
		if key == "name" {
			name = val[0]
		}
	}

	if name == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No name found.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	err = database.SetPlayerName(dbcs, id, name)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func SetPassword(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	if !isPlayerOrAdmin(r, id) {
		api.WriteBadHeader(w, http.StatusUnauthorized, "Player does not have access.")
		return
	}

	err = r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
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

	if password == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No password found.")
		return
	}

	if password != passwordConfirm {
		api.WriteBadHeader(w, http.StatusBadRequest, "Passwords do not match.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	err = database.SetPlayerPassword(dbcs, id, password)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func SetColorTheme(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	if !isPlayerOrAdmin(r, id) {
		api.WriteBadHeader(w, http.StatusUnauthorized, "Player does not have access.")
		return
	}

	err = r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
		return
	}

	var colorTheme string
	for key, val := range r.Form {
		if key == "colorTheme" {
			colorTheme = val[0]
		}
	}

	if colorTheme == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No color theme found.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	err = database.SetPlayerColorTheme(dbcs, id, colorTheme)
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

	if !isPlayerOrAdmin(r, id) {
		api.WriteBadHeader(w, http.StatusUnauthorized, "Player does not have access.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	err = database.DeletePlayer(dbcs, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	auth.RemoveCookiePlayerId(w)

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func isPlayerOrAdmin(r *http.Request, checkId uuid.UUID) bool {
	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		return false
	}

	if playerId == checkId {
		return true
	}

	dbcs := database.GetDatabaseConnectionString()
	player, err := database.GetPlayer(dbcs, playerId)
	if err != nil {
		return false
	}

	if player.IsAdmin {
		return true
	}

	return false
}
