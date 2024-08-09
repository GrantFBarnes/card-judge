package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
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

type PageDataLobbies struct {
	PageTitle  string
	PlayerName string
	Lobbies    []database.Lobby
}

func PageLobbies(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/base.html",
		"templates/pages/body/lobbies.html",
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

	tmpl.ExecuteTemplate(w, "base", PageDataLobbies{
		PageTitle:  "Card Judge - Lobbies",
		PlayerName: playerName,
		Lobbies:    lobbies,
	})
}

type PageDataDecks struct {
	PageTitle  string
	PlayerName string
	Decks      []database.Deck
}

func PageDecks(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/base.html",
		"templates/pages/body/decks.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	decks, err := database.GetDecks(dbcs)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	// playerName will be defined because of middleware check
	playerName, _ := auth.GetPlayerName(r)

	tmpl.ExecuteTemplate(w, "base", PageDataDecks{
		PageTitle:  "Card Judge - Decks",
		PlayerName: playerName,
		Decks:      decks,
	})
}

type PageDataCards struct {
	PageTitle  string
	PlayerName string
	Deck       database.Deck
	Cards      []database.Card
}

func PageCards(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/base.html",
		"templates/pages/body/cards.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	deckIdString := r.PathValue("deckid")
	deckId, err := uuid.Parse(deckIdString)
	if err != nil {
		fmt.Fprintf(w, "deck id invalid\n")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	deck, err := database.GetDeck(dbcs, deckId)
	if err != nil {
		fmt.Fprintf(w, "failed to get deck\n")
		return
	}

	cards, err := database.GetCards(dbcs, deckId)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	// playerName will be defined because of middleware check
	playerName, _ := auth.GetPlayerName(r)

	tmpl.ExecuteTemplate(w, "base", PageDataCards{
		PageTitle:  "Card Judge - Cards",
		PlayerName: playerName,
		Deck:       deck,
		Cards:      cards,
	})
}
