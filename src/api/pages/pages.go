package apiPages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

type HomeData struct {
	PageTitle string
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/home.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	tmpl.ExecuteTemplate(w, "base", HomeData{
		PageTitle: "Card Judge - Home",
	})
}

type LobbiesData struct {
	PageTitle string
	Lobbies   []database.Lobby
}

func Lobbies(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/lobbies.html",
		"templates/components/dialogs/lobby-create-dialog.html",
		"templates/components/forms/lobby-create-form.html",
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

	tmpl.ExecuteTemplate(w, "base", LobbiesData{
		PageTitle: "Card Judge - Lobbies",
		Lobbies:   lobbies,
	})
}

type LobbyData struct {
	PageTitle string
	LoggedIn  bool
	Player    database.Player
	HasAccess bool
	Lobby     database.Lobby
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/lobby.html",
		"templates/components/forms/player-create-form.html",
		"templates/components/forms/player-update-form.html",
		"templates/components/dialogs/player-update-dialog.html",
		"templates/components/forms/lobby-access-form.html",
		"templates/components/forms/lobby-update-form.html",
		"templates/components/dialogs/lobby-update-dialog.html",
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

	playerId, err := auth.GetPlayerId(r)
	loggedIn := err == nil
	player, err := database.GetPlayer(dbcs, playerId)
	if loggedIn && err != nil {
		fmt.Fprintf(w, "failed to get player\n")
		return
	}

	hasAccess := true
	if lobby.PasswordHash.Valid {
		hasAccess = auth.HasAccess(r, lobby.Id)
	}

	tmpl.ExecuteTemplate(w, "base", LobbyData{
		PageTitle: "Card Judge - Lobby",
		LoggedIn:  loggedIn,
		Player:    player,
		HasAccess: hasAccess,
		Lobby:     lobby,
	})
}

type DecksData struct {
	PageTitle string
	Decks     []database.Deck
}

func Decks(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/decks.html",
		"templates/components/forms/deck-create-form.html",
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

	tmpl.ExecuteTemplate(w, "base", DecksData{
		PageTitle: "Card Judge - Decks",
		Decks:     decks,
	})
}

type DeckData struct {
	PageTitle string
	HasAccess bool
	Deck      database.Deck
	Cards     []database.Card
}

func Deck(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/deck.html",
		"templates/components/forms/deck-access-form.html",
		"templates/components/forms/deck-update-form.html",
		"templates/components/dialogs/deck-update-dialog.html",
		"templates/components/forms/card-create-form.html",
		"templates/components/forms/card-update-form.html",
		"templates/components/dialogs/card-update-dialog.html",
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

	hasAccess := true
	if deck.PasswordHash.Valid {
		hasAccess = auth.HasAccess(r, deck.Id)
	}

	tmpl.ExecuteTemplate(w, "base", DeckData{
		PageTitle: "Card Judge - Deck",
		HasAccess: hasAccess,
		Deck:      deck,
		Cards:     cards,
	})
}
