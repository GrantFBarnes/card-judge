package api

import (
	"net/http"

	"github.com/grantfbarnes/card-judge/auth"
)

func PostLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
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
		return
	}

	err := auth.SetPlayerName(w, playerName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to set cookie"))
		return
	}

	w.Header().Add("HX-Redirect", auth.GetRedirectURL(r))
}
