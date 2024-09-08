package apiPages

import (
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
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Home"

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/home.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	tmpl.ExecuteTemplate(w, "base", homeData{
		Data: basePageData,
	})
}

type loginData struct {
	Data api.BasePageData
}

func Login(w http.ResponseWriter, r *http.Request) {
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Login"

	tmpl, err := template.ParseFiles(
		"templates/components/forms/player-create-form.html",
		"templates/components/forms/player-login-form.html",
		"templates/pages/base.html",
		"templates/pages/body/login.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	tmpl.ExecuteTemplate(w, "base", loginData{
		Data: basePageData,
	})
}

type manageData struct {
	Data api.BasePageData
}

func Manage(w http.ResponseWriter, r *http.Request) {
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Manage"

	tmpl, err := template.ParseFiles(
		"templates/components/forms/player-color-theme-form.html",
		"templates/components/forms/player-logout-form.html",
		"templates/components/forms/player-name-form.html",
		"templates/components/forms/player-password-form.html",
		"templates/pages/base.html",
		"templates/pages/body/manage.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	tmpl.ExecuteTemplate(w, "base", manageData{
		Data: basePageData,
	})
}

type adminData struct {
	Data    api.BasePageData
	Players []database.Player
}

func Admin(w http.ResponseWriter, r *http.Request) {
	players, err := database.GetPlayers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get players"))
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Admin"

	tmpl, err := template.ParseFiles(
		"templates/components/forms/player-create-default-form.html",
		"templates/components/table-rows/player-table-rows.html",
		"templates/pages/base.html",
		"templates/pages/body/admin.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	tmpl.ExecuteTemplate(w, "base", adminData{
		Data:    basePageData,
		Players: players,
	})
}

type lobbiesData struct {
	Data    api.BasePageData
	Lobbies []database.Lobby
	Decks   []database.Deck
}

func Lobbies(w http.ResponseWriter, r *http.Request) {
	lobbies, err := database.GetLobbies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get lobbies"))
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Lobbies"

	decks, err := database.GetPlayerDecks(basePageData.Player.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get player decks"))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/dialogs/lobby-create-dialog.html",
		"templates/components/forms/lobby-create-form.html",
		"templates/components/table-rows/lobby-table-rows.html",
		"templates/pages/base.html",
		"templates/pages/body/lobbies.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	tmpl.ExecuteTemplate(w, "base", lobbiesData{
		Data:    basePageData,
		Lobbies: lobbies,
		Decks:   decks,
	})
}

type lobbyData struct {
	Data      api.BasePageData
	HasAccess bool
	Lobby     database.Lobby
}

func Lobby(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		Lobbies(w, r)
		return
	}

	lobby, err := database.GetLobby(id)
	if err != nil {
		Lobbies(w, r)
		return
	}

	if lobby.Id == uuid.Nil {
		Lobbies(w, r)
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Lobby"

	tmpl, err := template.ParseFiles(
		"templates/components/dialogs/lobby-update-dialog.html",
		"templates/components/forms/lobby-access-form.html",
		"templates/components/forms/lobby-name-form.html",
		"templates/components/forms/lobby-password-form.html",
		"templates/pages/base.html",
		"templates/pages/body/lobby.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

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
	decks, err := database.GetDecks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get decks"))
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Decks"

	tmpl, err := template.ParseFiles(
		"templates/components/dialogs/deck-create-dialog.html",
		"templates/components/forms/deck-create-form.html",
		"templates/components/table-rows/deck-table-rows.html",
		"templates/pages/base.html",
		"templates/pages/body/decks.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

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
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		Decks(w, r)
		return
	}

	deck, err := database.GetDeck(id)
	if err != nil {
		Decks(w, r)
		return
	}

	if deck.Id == uuid.Nil {
		Decks(w, r)
		return
	}

	cards, err := database.GetCardsInDeck(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get cards in deck"))
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Deck"

	tmpl, err := template.ParseFiles(
		"templates/components/dialogs/card-create-dialog.html",
		"templates/components/dialogs/card-update-dialog.html",
		"templates/components/dialogs/deck-update-dialog.html",
		"templates/components/forms/card-create-form.html",
		"templates/components/forms/card-text-form.html",
		"templates/components/forms/card-type-form.html",
		"templates/components/forms/deck-access-form.html",
		"templates/components/forms/deck-name-form.html",
		"templates/components/forms/deck-password-form.html",
		"templates/components/table-rows/card-table-rows.html",
		"templates/pages/base.html",
		"templates/pages/body/deck.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	tmpl.ExecuteTemplate(w, "base", deckData{
		Data:      basePageData,
		HasAccess: !deck.PasswordHash.Valid || basePageData.Player.HasDeckAccess(deck.Id),
		Deck:      deck,
		Cards:     cards,
	})
}
