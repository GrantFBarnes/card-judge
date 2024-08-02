package api

import (
	"fmt"
	"net/http"

	"github.com/grantfbarnes/card-judge/auth"
)

func PostLogin(w http.ResponseWriter, req *http.Request) {
	playerName := "Grant"
	err := auth.SetPlayerName(w, playerName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failure"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf("Welcome %s!", playerName)))
}
