package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.GetPlayerName(r)
		loggedIn := err == nil

		if r.URL.Path == "/login" {
			if loggedIn {
				http.Redirect(w, r, auth.GetRedirectURL(r), http.StatusSeeOther)
				return
			}
		} else {
			if !loggedIn {
				auth.SetRedirectURL(w, r.URL.Path)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			auth.RemoveRedirectURL(w)
		}

		next.ServeHTTP(w, r)
	})
}

type PageDataLogin struct {
	PageTitle string
}

func PageLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/login.html",
		"templates/pages/body/login.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	tmpl.ExecuteTemplate(w, "base", PageDataLogin{
		PageTitle: "Card Judge - Login",
	})
}

type PageDataHome struct {
	PageTitle  string
	PlayerName string
}

func PageHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/base.html",
		"templates/pages/body/home.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	// playerName will be defined because of middleware check
	playerName, _ := auth.GetPlayerName(r)

	tmpl.ExecuteTemplate(w, "base", PageDataHome{
		PageTitle:  "Card Judge - Home",
		PlayerName: playerName,
	})
}

type PageDataLobbyJoin struct {
	PageTitle  string
	PlayerName string
	Lobbies    []database.Lobby
}

func PageLobbyJoin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/base.html",
		"templates/pages/body/lobby-join.html",
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

	// playerName will be defined because of middleware check
	playerName, _ := auth.GetPlayerName(r)

	tmpl.ExecuteTemplate(w, "base", PageDataLobbyJoin{
		PageTitle:  "Card Judge - Join Lobby",
		PlayerName: playerName,
		Lobbies:    lobbies,
	})
}

type PageDataCardList struct {
	PageTitle  string
	PlayerName string
	Cards      []database.Card
}

func PageCardList(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/base.html",
		"templates/pages/body/cards.html",
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

	// playerName will be defined because of middleware check
	playerName, _ := auth.GetPlayerName(r)

	tmpl.ExecuteTemplate(w, "base", PageDataCardList{
		PageTitle:  "Card Judge - Cards",
		PlayerName: playerName,
		Cards:      cards,
	})
}
