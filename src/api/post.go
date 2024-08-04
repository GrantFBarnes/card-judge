package api

import (
	"errors"
	"net/http"

	"github.com/grantfbarnes/card-judge/auth"
)

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
