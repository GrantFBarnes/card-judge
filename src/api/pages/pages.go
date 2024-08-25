package apiPages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

type LoginData struct {
	PageTitle string
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/login.html",
		"templates/pages/body/login.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	tmpl.ExecuteTemplate(w, "base", LoginData{
		PageTitle: "Card Judge - Login",
	})
}

type HomeData struct {
	PageTitle  string
	PlayerName string
}

func Home(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "base", HomeData{
		PageTitle:  "Card Judge - Home",
		PlayerName: playerName,
	})
}

type LobbiesData struct {
	PageTitle  string
	PlayerName string
	Lobbies    []database.Lobby
}

func Lobbies(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "base", LobbiesData{
		PageTitle:  "Card Judge - Lobbies",
		PlayerName: playerName,
		Lobbies:    lobbies,
	})
}

type LobbyData struct {
	PageTitle  string
	PlayerName string
	HasAccess  bool
	Lobby      database.Lobby
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/base.html",
		"templates/pages/body/lobby.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		fmt.Fprintf(w, "lobby id invalid\n")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	lobby, err := database.GetLobby(dbcs, id)
	if err != nil {
		fmt.Fprintf(w, "failed to get lobby\n")
		return
	}

	// playerName will be defined because of middleware check
	playerName, _ := auth.GetPlayerName(r)

	hasAccess := true
	if lobby.Password.Valid {
		hasAccess = auth.HasAccess(r, lobby.Id)
	}

	tmpl.ExecuteTemplate(w, "base", LobbyData{
		PageTitle:  "Card Judge - Lobby",
		PlayerName: playerName,
		HasAccess:  hasAccess,
		Lobby:      lobby,
	})
}

type DecksData struct {
	PageTitle  string
	PlayerName string
	Decks      []database.Deck
}

func Decks(w http.ResponseWriter, r *http.Request) {
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

	tmpl.ExecuteTemplate(w, "base", DecksData{
		PageTitle:  "Card Judge - Decks",
		PlayerName: playerName,
		Decks:      decks,
	})
}

type DeckData struct {
	PageTitle  string
	PlayerName string
	HasAccess  bool
	Deck       database.Deck
	Cards      []database.Card
}

func Deck(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/topbar/base.html",
		"templates/pages/body/deck.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		fmt.Fprintf(w, "deck id invalid\n")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	deck, err := database.GetDeck(dbcs, id)
	if err != nil {
		fmt.Fprintf(w, "failed to get deck\n")
		return
	}

	cards, err := database.GetCardsInDeck(dbcs, id)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	// playerName will be defined because of middleware check
	playerName, _ := auth.GetPlayerName(r)

	hasAccess := true
	if deck.Password.Valid {
		hasAccess = auth.HasAccess(r, deck.Id)
	}

	tmpl.ExecuteTemplate(w, "base", DeckData{
		PageTitle:  "Card Judge - Deck",
		PlayerName: playerName,
		HasAccess:  hasAccess,
		Deck:       deck,
		Cards:      cards,
	})
}
