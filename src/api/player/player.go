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
	for key, val := range r.Form {
		if key == "playerName" {
			name = val[0]
		}
	}

	if name == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No name found.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	id, err := database.CreatePlayer(dbcs, name)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to create player in database.")
		return
	}

	err = auth.SetPlayerId(w, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to create player cookie in browser.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
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
	for key, val := range r.Form {
		if key == "playerName" {
			name = val[0]
		}
	}

	if name == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No name found.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	err = database.UpdatePlayer(dbcs, id, name)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update player in database.")
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

	dbcs := database.GetDatabaseConnectionString()
	err = database.DeletePlayer(dbcs, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to delete player in database.")
		return
	}

	auth.RemovePlayerId(w)

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}
