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

	tmpl.ExecuteTemplate(w, "base", basePageData)
}

func About(w http.ResponseWriter, r *http.Request) {
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - About"

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/about.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	tmpl.ExecuteTemplate(w, "base", basePageData)
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

	tmpl.ExecuteTemplate(w, "base", basePageData)
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

	tmpl.ExecuteTemplate(w, "base", basePageData)
}

func Admin(w http.ResponseWriter, r *http.Request) {
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Admin"

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/admin.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	tmpl.ExecuteTemplate(w, "base", basePageData)
}

func Stats(w http.ResponseWriter, r *http.Request) {
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Stats"

	mostWinsByPlayer, err := database.GetMostWinsByPlayer()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get most wins by player"))
		return
	}

	mostWinsByCard, err := database.GetMostWinsByCard(basePageData.User.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get most wins by card"))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/stats.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	type data struct {
		api.BasePageData
		MostWinsByPlayer []database.MostWins
		MostWinsByCard   []database.MostWins
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData:     basePageData,
		MostWinsByPlayer: mostWinsByPlayer,
		MostWinsByCard:   mostWinsByCard,
	})
}

func Lobbies(w http.ResponseWriter, r *http.Request) {
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Lobbies"

	decks, err := database.GetReadableDecks(basePageData.User.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get user decks"))
		return
	}

	tmpl, err := template.ParseFiles(
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
		Decks []database.Deck
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
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

	hasLobbyAccess, err := database.UserHasLobbyAccess(basePageData.User.Id, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to check lobby access"))
		return
	}

	if !hasLobbyAccess {
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

	hasLobbyAccess, err := database.UserHasLobbyAccess(basePageData.User.Id, lobbyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to check lobby access"))
		return
	}

	if hasLobbyAccess {
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
	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Decks"

	tmpl, err := template.ParseFiles(
		"templates/pages/base.html",
		"templates/pages/body/decks.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to parse HTML"))
		return
	}

	tmpl.ExecuteTemplate(w, "base", basePageData)
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

	basePageData := api.GetBasePageData(r)
	basePageData.PageTitle = "Card Judge - Deck"

	hasDeckAccess, err := database.UserHasDeckAccess(basePageData.User.Id, deckId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to check deck access"))
		return
	}

	if !hasDeckAccess {
		http.Redirect(w, r, fmt.Sprintf("/deck/%s/access", deckId), http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles(
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
		Deck database.Deck
	}

	tmpl.ExecuteTemplate(w, "base", data{
		BasePageData: basePageData,
		Deck:         deck,
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

	hasDeckAccess, err := database.UserHasDeckAccess(basePageData.User.Id, deckId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to check deck access"))
		return
	}

	if hasDeckAccess {
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
