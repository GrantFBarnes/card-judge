package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/grantfbarnes/card-judge/database"
)

type PageDataHome struct {
	PageTitle string
}

func PageHome(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/home.html",
		"templates/base.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	tmpl.ExecuteTemplate(w, "base", PageDataHome{
		PageTitle: "Card Judge - Home",
	})
}

func PageLogin(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/login.html",
		"templates/base.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	tmpl.ExecuteTemplate(w, "base", nil)
}

type PageDataLobbyJoin struct {
	PageTitle string
	Lobbies   []database.Lobby
}

func PageLobbyJoin(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/lobby/join.html",
		"templates/base.html",
	)
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

	tmpl.ExecuteTemplate(w, "base", PageDataLobbyJoin{
		PageTitle: "Card Judge - Home",
		Lobbies:   lobbies,
	})
}

type PageDataCardList struct {
	PageTitle string
	Cards     []database.Card
}

func PageCardList(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/cards.html",
		"templates/base.html",
	)
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

	tmpl.ExecuteTemplate(w, "base", PageDataCardList{
		PageTitle: "Card Judge - Home",
		Cards:     cards,
	})
}
