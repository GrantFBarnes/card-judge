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
	Player    database.Player
	LoggedIn  bool
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

	player, err := getRequestPlayer(r)
	loggedIn := err == nil

	tmpl.ExecuteTemplate(w, "base", HomeData{
		PageTitle: "Card Judge - Home",
		Player:    player,
		LoggedIn:  loggedIn,
	})
}

type LoginData struct {
	PageTitle string
	Player    database.Player
	LoggedIn  bool
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/login.html",
		"templates/components/forms/player-create-form.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	player, err := getRequestPlayer(r)
	loggedIn := err == nil

	tmpl.ExecuteTemplate(w, "base", LoginData{
		PageTitle: "Card Judge - Login",
		Player:    player,
		LoggedIn:  loggedIn,
	})
}

type ManageData struct {
	PageTitle string
	Player    database.Player
	LoggedIn  bool
}

func Manage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/manage.html",
		"templates/components/forms/player-update-form.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	player, err := getRequestPlayer(r)
	loggedIn := err == nil

	tmpl.ExecuteTemplate(w, "base", ManageData{
		PageTitle: "Card Judge - Manage",
		Player:    player,
		LoggedIn:  loggedIn,
	})
}

type LobbiesData struct {
	PageTitle string
	Lobbies   []database.Lobby
	Player    database.Player
	LoggedIn  bool
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

	player, err := getRequestPlayer(r)
	loggedIn := err == nil

	tmpl.ExecuteTemplate(w, "base", LobbiesData{
		PageTitle: "Card Judge - Lobbies",
		Lobbies:   lobbies,
		Player:    player,
		LoggedIn:  loggedIn,
	})
}

type LobbyData struct {
	PageTitle string
	HasAccess bool
	Lobby     database.Lobby
	Player    database.Player
	LoggedIn  bool
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/lobby.html",
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

	hasAccess := true
	if lobby.PasswordHash.Valid {
		hasAccess = auth.HasCookieAccess(r, lobby.Id)
	}

	player, err := getRequestPlayer(r)
	loggedIn := err == nil

	tmpl.ExecuteTemplate(w, "base", LobbyData{
		PageTitle: "Card Judge - Lobby",
		HasAccess: hasAccess,
		Lobby:     lobby,
		Player:    player,
		LoggedIn:  loggedIn,
	})
}

type DecksData struct {
	PageTitle string
	Decks     []database.Deck
	Player    database.Player
	LoggedIn  bool
}

func Decks(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/decks.html",
		"templates/components/dialogs/deck-create-dialog.html",
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

	player, err := getRequestPlayer(r)
	loggedIn := err == nil

	tmpl.ExecuteTemplate(w, "base", DecksData{
		PageTitle: "Card Judge - Decks",
		Decks:     decks,
		Player:    player,
		LoggedIn:  loggedIn,
	})
}

type DeckData struct {
	PageTitle string
	HasAccess bool
	Deck      database.Deck
	Cards     []database.Card
	Player    database.Player
	LoggedIn  bool
}

func Deck(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/deck.html",
		"templates/components/forms/deck-access-form.html",
		"templates/components/forms/deck-update-form.html",
		"templates/components/dialogs/deck-update-dialog.html",
		"templates/components/forms/card-create-form.html",
		"templates/components/dialogs/card-create-dialog.html",
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
		hasAccess = auth.HasCookieAccess(r, deck.Id)
	}

	player, err := getRequestPlayer(r)
	loggedIn := err == nil

	tmpl.ExecuteTemplate(w, "base", DeckData{
		PageTitle: "Card Judge - Deck",
		HasAccess: hasAccess,
		Deck:      deck,
		Cards:     cards,
		Player:    player,
		LoggedIn:  loggedIn,
	})
}

func getRequestPlayer(r *http.Request) (player database.Player, err error) {
	playerId, err := auth.GetCookiePlayerId(r)
	if err != nil {
		return player, err
	}

	dbcs := database.GetDatabaseConnectionString()
	player, err = database.GetPlayer(dbcs, playerId)
	if err != nil {
		return player, err
	}

	return player, nil
}
