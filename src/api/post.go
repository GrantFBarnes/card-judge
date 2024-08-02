package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/grantfbarnes/card-judge/auth"
)

func PostLogin(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var playerName string
	for key, val := range req.Form {
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

	tmpl, err := template.ParseFiles(
		"templates/components/login.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	tmpl.ExecuteTemplate(w, "login", PageDataHome{
		LoggedIn:   true,
		PlayerName: playerName,
	})
}
