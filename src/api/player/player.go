package apiPlayer

import (
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

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

	players, err := database.SearchPlayers(search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/table-rows/player-table-rows.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "player-table-rows", players)
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

	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No password found."))
		return
	}

	if password != passwordConfirm {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passwords do not match."))
		return
	}

	existingPlayerId, err := database.GetPlayerId(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if existingPlayerId != uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Player name already exists."))
		return
	}

	id, err := database.CreatePlayer(name, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = auth.SetCookiePlayerId(w, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusCreated)
}

func CreateDefault(w http.ResponseWriter, r *http.Request) {
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

	if password == "" {
		password = "password"
	} else if password != passwordConfirm {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passwords do not match."))
		return
	}

	existingPlayerId, err := database.GetPlayerId(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if existingPlayerId != uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Player name already exists."))
		return
	}

	_, err = database.CreatePlayer(name, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No name found."))
		return
	}

	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No password found."))
		return
	}

	playerId, err := database.GetPlayerId(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Player name does not exist."))
		return
	}

	passwordHash, err := database.GetPlayerPasswordHash(playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if !auth.PasswordMatchesHash(password, passwordHash) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Provided password is not valid."))
		return
	}

	err = auth.SetCookiePlayerId(w, playerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	auth.RemoveCookiePlayerId(w)
	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func SetName(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	if !isCurrentPlayer(r, id) {
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

	existingPlayerId, err := database.GetPlayerId(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if existingPlayerId != uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Player name already exists."))
		return
	}

	err = database.SetPlayerName(id, name)
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

	if !isCurrentPlayer(r, id) {
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

	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No password found."))
		return
	}

	if password != passwordConfirm {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passwords do not match."))
		return
	}

	err = database.SetPlayerPassword(id, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	if !isAdmin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Player does not have access."))
		return
	}

	err = database.SetPlayerPassword(id, "password")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("&#9989;"))
}

func SetColorTheme(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	if !isCurrentPlayer(r, id) {
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

	var colorTheme string
	for key, val := range r.Form {
		if key == "colorTheme" {
			colorTheme = val[0]
		}
	}

	if colorTheme == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No color theme found."))
		return
	}

	err = database.SetPlayerColorTheme(id, colorTheme)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func SetIsAdmin(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	if !isAdmin(r) {
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

	var isAdmin bool
	for key, val := range r.Form {
		if key == "isAdmin" {
			isAdmin = val[0] == "1"
		}
	}

	err = database.SetPlayerIsAdmin(id, isAdmin)
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

	isCurrentPlayer := isCurrentPlayer(r, id)
	if !isCurrentPlayer && !isAdmin(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Player does not have access."))
		return
	}

	err = database.DeletePlayer(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if isCurrentPlayer {
		auth.RemoveCookiePlayerId(w)
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func isCurrentPlayer(r *http.Request, checkId uuid.UUID) bool {
	playerId := api.GetPlayerId(r)
	return playerId == checkId
}

func isAdmin(r *http.Request) bool {
	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		return false
	}

	isAdmin, _ := database.GetPlayerIsAdmin(playerId)
	return isAdmin
}
