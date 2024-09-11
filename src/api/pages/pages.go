package apiPages

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/database"
)

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

	type data struct {
		api.BasePageData
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Login"

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/login.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
	})
}

func Manage(w http.ResponseWriter, r *http.Request) {
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Manage"

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/manage.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
	})
}

func Admin(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetUsers("%")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get users"))
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Admin"

	tmpl, err := template.ParseFiles(
		"templates/components/table-rows/user-table-rows.html",
		"templates/pages/base.html",
		"templates/pages/body/admin.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
		Users []database.User
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
		Users:        users,
	})
}

func Lobbies(w http.ResponseWriter, r *http.Request) {
	lobbies, err := database.GetLobbies("%")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get lobbies"))
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Lobbies"

	decks, err := database.GetUserDecks(basePageData.User.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get user decks"))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/table-rows/lobby-table-rows.html",
		"templates/pages/base.html",
		"templates/pages/body/lobbies.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
		Lobbies []database.LobbyDetails
		Decks   []database.Deck
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
		Lobbies:      lobbies,
		Decks:        decks,
	})
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
		"templates/pages/base.html",
		"templates/pages/body/lobby.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
		HasAccess bool
		Lobby     database.Lobby
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
		HasAccess:    database.HasLobbyAccess(basePageData.User.Id, lobby.Id),
		Lobby:        lobby,
	})
}

func Decks(w http.ResponseWriter, r *http.Request) {
	decks, err := database.GetDecks("%")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get decks"))
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Decks"

	tmpl, err := template.ParseFiles(
		"templates/components/table-rows/deck-table-rows.html",
		"templates/pages/base.html",
		"templates/pages/body/decks.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
		Decks []database.DeckDetails
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
		Decks:        decks,
	})
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

	cards, err := database.GetCardsInDeck(id, "%")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get cards in deck"))
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Deck"

	tmpl, err := template.ParseFiles(
		"templates/components/table-rows/card-table-rows.html",
		"templates/pages/base.html",
		"templates/pages/body/deck.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
		HasAccess bool
		Deck      database.Deck
		Cards     []database.CardDetails
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
		HasAccess:    database.HasDeckAccess(basePageData.User.Id, deck.Id),
		Deck:         deck,
		Cards:        cards,
	})
}
