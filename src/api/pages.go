package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

func PageHome(w http.ResponseWriter, req *http.Request) {
	type PageDataHome struct {
		LoggedIn   bool
		PlayerName string
	}

	tmpl, err := template.ParseFiles("templates/pages/home.html", "templates/base.html")
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	playerName, err := auth.GetPlayerName(req)

	tmpl.ExecuteTemplate(w, "base", PageDataHome{
		LoggedIn:   err == nil,
		PlayerName: playerName,
	})
}

func PageLobbyJoin(w http.ResponseWriter, req *http.Request) {
	type PageDataLobbyJoin struct {
		LoggedIn   bool
		PlayerName string
		Lobbies    []database.Lobby
	}

	tmpl, err := template.ParseFiles("templates/pages/lobby/join.html", "templates/base.html")
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	lobbies, err := database.GetLobbies(dbcs)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	playerName, err := auth.GetPlayerName(req)

	tmpl.ExecuteTemplate(w, "base", PageDataLobbyJoin{
		LoggedIn:   err == nil,
		PlayerName: playerName,
		Lobbies:    lobbies,
	})
}

func PageCardList(w http.ResponseWriter, req *http.Request) {
	type PageDataCardList struct {
		LoggedIn   bool
		PlayerName string
		Cards      []database.Card
	}

	tmpl, err := template.ParseFiles("templates/cards.html", "templates/base.html")
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	cards, err := database.GetJudgeCards(dbcs)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	playerName, err := auth.GetPlayerName(req)

	tmpl.ExecuteTemplate(w, "base", PageDataCardList{
		LoggedIn:   err == nil,
		PlayerName: playerName,
		Cards:      cards,
	})
}
