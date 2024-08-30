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
		"templates/components/forms/player-delete-form.html",
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

type lobbiesData struct {
	Data    api.BasePageData
	Lobbies []database.Lobby
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

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Lobby"

	hasAccess := !lobby.PasswordHash.Valid
	if !hasAccess {
		for _, id := range basePageData.Player.LobbyIds {
			if id == lobby.Id {
				hasAccess = true
				break
			}
		}
	}

	tmpl.ExecuteTemplate(w, "base", lobbyData{
		Data:      basePageData,
		HasAccess: hasAccess,
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

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Deck"

	hasAccess := !deck.PasswordHash.Valid
	if !hasAccess {
		for _, id := range basePageData.Player.DeckIds {
			if id == deck.Id {
				hasAccess = true
				break
			}
		}
	}

	tmpl.ExecuteTemplate(w, "base", deckData{
		Data:      basePageData,
		HasAccess: hasAccess,
		Deck:      deck,
		Cards:     cards,
	})
}
