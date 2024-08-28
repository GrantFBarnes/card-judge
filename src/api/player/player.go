package apiPlayer

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

func Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse form"))
		return
	}

	var name string
	for key, val := range r.Form {
		if key == "playerName" {
			name = val[0]
		}
	}

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no name found"))
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	id, err := database.CreatePlayer(dbcs, name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to create player"))
		return
	}

	err = auth.SetPlayerId(w, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to set cookie"))
		return
	}

	w.Header().Add("HX-Refresh", "true")
}

func Update(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get player id"))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse form"))
		return
	}

	var name string
	for key, val := range r.Form {
		if key == "playerName" {
			name = val[0]
		}
	}

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no name found"))
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	err = database.UpdatePlayer(dbcs, id, name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to update player"))
		return
	}

	w.Header().Add("HX-Refresh", "true")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get player id"))
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	err = database.DeletePlayer(dbcs, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to delete player"))
		return
	}

	auth.RemovePlayerId(w)

	w.Header().Add("HX-Refresh", "true")
}
