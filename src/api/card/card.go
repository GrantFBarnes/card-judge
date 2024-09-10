package apiCard

import (
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/database"
)

func Search(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var deckId uuid.UUID
	var cardTypeSearch string
	var textSearch string
	for key, val := range r.Form {
		if key == "deckId" {
			deckId, err = uuid.Parse(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse deck id."))
				return
			}
		} else if key == "cardType" {
			cardTypeSearch = val[0]
		} else if key == "text" {
			textSearch = val[0]
		}
	}

	if cardTypeSearch == "" {
		cardTypeSearch = "%"
	}

	textSearch = "%" + textSearch + "%"

	cards, err := database.GetCardsInDeck(deckId, textSearch, cardTypeSearch)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/components/table-rows/card-table-rows.html",
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse HTML."))
		return
	}

	tmpl.ExecuteTemplate(w, "card-table-rows", cards)
}

func Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var deckId uuid.UUID
	var cardType database.CardType
	var text string
	for key, val := range r.Form {
		if key == "deckId" {
			deckId, err = uuid.Parse(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse deck id."))
				return
			}
		} else if key == "cardType" {
			if val[0] == "Judge" {
				cardType = database.JudgeCard
			} else if val[0] == "Player" {
				cardType = database.PlayerCard
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse card type."))
				return
			}
		} else if key == "text" {
			text = val[0]
		}
	}

	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No text found."))
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
		return
	}

	existingCardId, err := database.GetCardId(deckId, text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if existingCardId != uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Card text already exists."))
		return
	}

	if !database.HasDeckAccess(playerId, deckId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Player does not have access."))
		return
	}

	_, err = database.CreateCard(deckId, cardType, text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusCreated)
}

func SetType(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var deckId uuid.UUID
	var cardType database.CardType
	for key, val := range r.Form {
		if key == "deckId" {
			deckId, err = uuid.Parse(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse deck id."))
				return
			}
		} else if key == "cardType" {
			if val[0] == "Judge" {
				cardType = database.JudgeCard
			} else if val[0] == "Player" {
				cardType = database.PlayerCard
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse card type."))
				return
			}
		}
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
		return
	}

	if !database.HasDeckAccess(playerId, deckId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Player does not have access."))
		return
	}

	err = database.SetCardType(id, cardType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func SetText(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var deckId uuid.UUID
	var text string
	for key, val := range r.Form {
		if key == "deckId" {
			deckId, err = uuid.Parse(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse deck id."))
				return
			}
		} else if key == "text" {
			text = val[0]
		}
	}

	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No text found."))
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
		return
	}

	if !database.HasDeckAccess(playerId, deckId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Player does not have access."))
		return
	}

	existingCardId, err := database.GetCardId(deckId, text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if existingCardId != uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Card text already exists."))
		return
	}

	err = database.SetCardText(id, text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get id from path."))
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get player id."))
		return
	}

	card, err := database.GetCard(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card."))
		return
	}

	if !database.HasDeckAccess(playerId, card.DeckId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Player does not have access."))
		return
	}

	err = database.DeleteCard(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
