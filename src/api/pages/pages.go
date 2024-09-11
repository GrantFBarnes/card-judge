package apiPages

import (
	"fmt"
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
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		http.Redirect(w, r, "/lobbies", http.StatusSeeOther)
		return
	}

	lobby, err := database.GetLobby(lobbyId)
	if err != nil {
		http.Redirect(w, r, "/lobbies", http.StatusSeeOther)
		return
	}

	if lobby.Id == uuid.Nil {
		http.Redirect(w, r, "/lobbies", http.StatusSeeOther)
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Lobby"

	if !database.HasLobbyAccess(basePageData.User.Id, lobbyId) {
		http.Redirect(w, r, fmt.Sprintf("/lobby/%s/access", lobbyId), http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/lobby.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	playerId, err := database.AddUserToLobby(lobbyId, basePageData.User.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to join lobby"))
		return
	}

	type data struct {
		api.BasePageData
		Lobby    database.Lobby
		PlayerId uuid.UUID
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
		Lobby:        lobby,
		PlayerId:     playerId,
	})
}

func LobbyAccess(w http.ResponseWriter, r *http.Request) {
	lobbyIdString := r.PathValue("lobbyId")
	lobbyId, err := uuid.Parse(lobbyIdString)
	if err != nil {
		http.Redirect(w, r, "/lobbies", http.StatusSeeOther)
		return
	}

	lobby, err := database.GetLobby(lobbyId)
	if err != nil {
		http.Redirect(w, r, "/lobbies", http.StatusSeeOther)
		return
	}

	if lobby.Id == uuid.Nil {
		http.Redirect(w, r, "/lobbies", http.StatusSeeOther)
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Lobby Access"

	if database.HasLobbyAccess(basePageData.User.Id, lobbyId) {
		http.Redirect(w, r, fmt.Sprintf("/lobby/%s", lobbyId), http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/lobby-access.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
		Lobby database.Lobby
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
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
	deckIdString := r.PathValue("deckId")
	deckId, err := uuid.Parse(deckIdString)
	if err != nil {
		http.Redirect(w, r, "/decks", http.StatusSeeOther)
		return
	}

	deck, err := database.GetDeck(deckId)
	if err != nil {
		http.Redirect(w, r, "/decks", http.StatusSeeOther)
		return
	}

	if deck.Id == uuid.Nil {
		http.Redirect(w, r, "/decks", http.StatusSeeOther)
		return
	}

	cards, err := database.GetCardsInDeck(deckId, "%")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get cards in deck"))
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Deck"

	if !database.HasDeckAccess(basePageData.User.Id, deckId) {
		http.Redirect(w, r, fmt.Sprintf("/deck/%s/access", deckId), http.StatusSeeOther)
		return
	}

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
		Deck  database.Deck
		Cards []database.CardDetails
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
		Deck:         deck,
		Cards:        cards,
	})
}

func DeckAccess(w http.ResponseWriter, r *http.Request) {
	deckIdString := r.PathValue("deckId")
	deckId, err := uuid.Parse(deckIdString)
	if err != nil {
		http.Redirect(w, r, "/decks", http.StatusSeeOther)
		return
	}

	deck, err := database.GetDeck(deckId)
	if err != nil {
		http.Redirect(w, r, "/decks", http.StatusSeeOther)
		return
	}

	if deck.Id == uuid.Nil {
		http.Redirect(w, r, "/decks", http.StatusSeeOther)
		return
	}

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Deck"

	if database.HasDeckAccess(basePageData.User.Id, deckId) {
		http.Redirect(w, r, fmt.Sprintf("/deck/%s", deckId), http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/deck-access.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
		Deck database.Deck
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
		Deck:         deck,
	})
}
