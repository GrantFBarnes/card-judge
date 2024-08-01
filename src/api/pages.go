package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

func PageHome(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/pages/home.html", "templates/base.html")
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	type HomePageData struct {
		LoggedIn   bool
		PlayerName string
	}

	playerName, err := auth.GetPlayerName(req)
	tmpl.ExecuteTemplate(w, "base", HomePageData{
		LoggedIn:   err == nil,
		PlayerName: playerName,
	})
}

func PageLobbyJoin(w http.ResponseWriter, req *http.Request) {
	dbcs := database.GetDatabaseConnectionString()
	lobbies, err := database.GetLobbies(dbcs)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	tmpl, err := template.ParseFiles("templates/pages/lobby/join.html", "templates/base.html")
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	tmpl.ExecuteTemplate(w, "base", lobbies)
}

func PageCardList(w http.ResponseWriter, req *http.Request) {
	dbcs := database.GetDatabaseConnectionString()
	cards, err := database.GetJudgeCards(dbcs)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	tmpl, err := template.ParseFiles("templates/cards.html", "templates/base.html")
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	tmpl.ExecuteTemplate(w, "base", cards)
}
