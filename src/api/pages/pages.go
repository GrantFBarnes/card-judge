package apiPages

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/auth"
	"github.com/grantfbarnes/card-judge/database"
)

type requestContextKey string

const basePageDataKey requestContextKey = "basePageDataKey"

type basePageData struct {
	Player   database.Player
	LoggedIn bool
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		basePageData := basePageData{
			Player:   database.Player{},
			LoggedIn: false,
		}

		playerId, err := auth.GetCookiePlayerId(r)
		if err == nil {
			dbcs := database.GetDatabaseConnectionString()
			player, err := database.GetPlayer(dbcs, playerId)
			if err == nil {
				basePageData.Player = player
				basePageData.LoggedIn = true
			}
		}

		// required to be logged in
		if r.URL.Path == "/manage" ||
			strings.HasPrefix(r.URL.Path, "/lobby/") ||
			strings.HasPrefix(r.URL.Path, "/deck/") {
			if !basePageData.LoggedIn {
				auth.SetCookieRedirectURL(w, r.URL.Path)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}

		// required to not be logged in
		if r.URL.Path == "/login" {
			if basePageData.LoggedIn {
				http.Redirect(w, r, auth.GetCookieRedirectURL(r), http.StatusSeeOther)
				return
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), basePageDataKey, basePageData))

		next.ServeHTTP(w, r)
	})
}

type homeData struct {
	PageTitle string
	Data      basePageData
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

	tmpl.ExecuteTemplate(w, "base", homeData{
		PageTitle: "Card Judge - Home",
		Data:      r.Context().Value(basePageDataKey).(basePageData),
	})
}

type loginData struct {
	PageTitle string
	Data      basePageData
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

	tmpl.ExecuteTemplate(w, "base", loginData{
		PageTitle: "Card Judge - Login",
		Data:      r.Context().Value(basePageDataKey).(basePageData),
	})
}

type manageData struct {
	PageTitle string
	Data      basePageData
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

	tmpl.ExecuteTemplate(w, "base", manageData{
		PageTitle: "Card Judge - Manage",
		Data:      r.Context().Value(basePageDataKey).(basePageData),
	})
}

type lobbiesData struct {
	PageTitle string
	Lobbies   []database.Lobby
	Data      basePageData
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

	tmpl.ExecuteTemplate(w, "base", lobbiesData{
		PageTitle: "Card Judge - Lobbies",
		Lobbies:   lobbies,
		Data:      r.Context().Value(basePageDataKey).(basePageData),
	})
}

type lobbyData struct {
	PageTitle string
	HasAccess bool
	Lobby     database.Lobby
	Data      basePageData
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

	tmpl.ExecuteTemplate(w, "base", lobbyData{
		PageTitle: "Card Judge - Lobby",
		HasAccess: hasAccess,
		Lobby:     lobby,
		Data:      r.Context().Value(basePageDataKey).(basePageData),
	})
}

type decksData struct {
	PageTitle string
	Decks     []database.Deck
	Data      basePageData
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

	tmpl.ExecuteTemplate(w, "base", decksData{
		PageTitle: "Card Judge - Decks",
		Decks:     decks,
		Data:      r.Context().Value(basePageDataKey).(basePageData),
	})
}

type deckData struct {
	PageTitle string
	HasAccess bool
	Deck      database.Deck
	Cards     []database.Card
	Data      basePageData
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

	tmpl.ExecuteTemplate(w, "base", deckData{
		PageTitle: "Card Judge - Deck",
		HasAccess: hasAccess,
		Deck:      deck,
		Cards:     cards,
		Data:      r.Context().Value(basePageDataKey).(basePageData),
	})
}
