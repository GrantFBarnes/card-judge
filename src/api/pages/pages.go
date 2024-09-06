package apiPages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/database"
)

type homeData struct {
	Data api.BasePageData
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

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Home"

	tmpl.ExecuteTemplate(w, "base", homeData{
		Data: basePageData,
	})
}

type loginData struct {
	Data api.BasePageData
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/login.html",
		"templates/components/forms/player-login-form.html",
		"templates/components/forms/player-create-form.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Login"

	tmpl.ExecuteTemplate(w, "base", loginData{
		Data: basePageData,
	})
}

type manageData struct {
	Data api.BasePageData
}

func Manage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/manage.html",
		"templates/components/forms/player-name-form.html",
		"templates/components/forms/player-password-form.html",
		"templates/components/forms/player-color-theme-form.html",
		"templates/components/forms/player-logout-form.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Manage"

	tmpl.ExecuteTemplate(w, "base", manageData{
		Data: basePageData,
	})
}

type adminData struct {
	Data    api.BasePageData
	Players []database.Player
}

func Admin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/admin.html",
		"templates/components/table-rows/player-table-rows.html",
		"templates/components/forms/player-create-default-form.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	players, err := database.GetPlayers()
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Admin"

	tmpl.ExecuteTemplate(w, "base", adminData{
		Data:    basePageData,
		Players: players,
	})
}

type lobbiesData struct {
	Data    api.BasePageData
	Lobbies []database.Lobby
}

func Lobbies(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/lobbies.html",
		"templates/components/table-rows/lobby-table-rows.html",
		"templates/components/dialogs/lobby-create-dialog.html",
		"templates/components/forms/lobby-create-form.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	lobbies, err := database.GetLobbies()
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Lobbies"

	tmpl.ExecuteTemplate(w, "base", lobbiesData{
		Data:    basePageData,
		Lobbies: lobbies,
	})
}

type lobbyData struct {
	Data      api.BasePageData
	HasAccess bool
	Lobby     database.Lobby
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/lobby.html",
		"templates/components/forms/lobby-access-form.html",
		"templates/components/forms/lobby-name-form.html",
		"templates/components/forms/lobby-password-form.html",
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

	lobby, err := database.GetLobby(id)
	if err != nil {
		fmt.Fprintf(w, "failed to get lobby\n")
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Lobby"

	tmpl.ExecuteTemplate(w, "base", lobbyData{
		Data:      basePageData,
		HasAccess: !lobby.PasswordHash.Valid || basePageData.Player.HasLobbyAccess(lobby.Id),
		Lobby:     lobby,
	})
}

type decksData struct {
	Data  api.BasePageData
	Decks []database.Deck
}

func Decks(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/decks.html",
		"templates/components/table-rows/deck-table-rows.html",
		"templates/components/dialogs/deck-create-dialog.html",
		"templates/components/forms/deck-create-form.html",
	)
	if err != nil {
		fmt.Fprintf(w, "failed to parse HTML\n")
		return
	}

	decks, err := database.GetDecks()
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Decks"

	tmpl.ExecuteTemplate(w, "base", decksData{
		Data:  basePageData,
		Decks: decks,
	})
}

type deckData struct {
	Data      api.BasePageData
	HasAccess bool
	Deck      database.Deck
	Cards     []database.Card
}

func Deck(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/deck.html",
		"templates/components/forms/deck-access-form.html",
		"templates/components/forms/deck-name-form.html",
		"templates/components/forms/deck-password-form.html",
		"templates/components/dialogs/deck-update-dialog.html",
		"templates/components/table-rows/card-table-rows.html",
		"templates/components/forms/card-create-form.html",
		"templates/components/dialogs/card-create-dialog.html",
		"templates/components/forms/card-type-form.html",
		"templates/components/forms/card-text-form.html",
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

	deck, err := database.GetDeck(id)
	if err != nil {
		fmt.Fprintf(w, "failed to get deck\n")
		return
	}

	cards, err := database.GetCardsInDeck(id)
	if err != nil {
		fmt.Fprintf(w, "failed to connect to database\n")
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Deck"

	tmpl.ExecuteTemplate(w, "base", deckData{
		Data:      basePageData,
		HasAccess: !deck.PasswordHash.Valid || basePageData.Player.HasDeckAccess(deck.Id),
		Deck:      deck,
		Cards:     cards,
	})
}
