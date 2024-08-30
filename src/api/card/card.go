package apiCard

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/api"
	"github.com/grantfbarnes/card-judge/database"
)

func Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
		return
	}

	var deckId uuid.UUID
	var cardType database.CardType
	var text string
	for key, val := range r.Form {
		if key == "deckId" {
			deckId, err = uuid.Parse(val[0])
			if err != nil {
				api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse deck id.")
				return
			}
		} else if key == "cardType" {
			if val[0] == "Judge" {
				cardType = database.JudgeCard
			} else if val[0] == "Player" {
				cardType = database.PlayerCard
			} else {
				api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse card type.")
				return
			}
		} else if key == "text" {
			text = val[0]
		}
	}

	if text == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No text found.")
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get player id.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	if !database.HasDeckAccess(dbcs, playerId, deckId) {
		api.WriteBadHeader(w, http.StatusUnauthorized, "Player does not have access.")
		return
	}

	_, err = database.CreateCard(dbcs, playerId, deckId, cardType, text)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func Update(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	err = r.ParseForm()
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse form.")
		return
	}

	var deckId uuid.UUID
	var cardType database.CardType
	var text string
	for key, val := range r.Form {
		if key == "deckId" {
			deckId, err = uuid.Parse(val[0])
			if err != nil {
				api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse deck id.")
				return
			}
		} else if key == "cardType" {
			if val[0] == "Judge" {
				cardType = database.JudgeCard
			} else if val[0] == "Player" {
				cardType = database.PlayerCard
			} else {
				api.WriteBadHeader(w, http.StatusBadRequest, "Failed to parse card type.")
				return
			}
		} else if key == "text" {
			text = val[0]
		}
	}

	if text == "" {
		api.WriteBadHeader(w, http.StatusBadRequest, "No text found.")
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get player id.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	if !database.HasDeckAccess(dbcs, playerId, deckId) {
		api.WriteBadHeader(w, http.StatusUnauthorized, "Player does not have access.")
		return
	}

	err = database.UpdateCard(dbcs, playerId, id, cardType, text)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get id from path.")
		return
	}

	playerId := api.GetPlayerId(r)
	if playerId == uuid.Nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get player id.")
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	card, err := database.GetCard(dbcs, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to get card.")
		return
	}

	if !database.HasDeckAccess(dbcs, playerId, card.DeckId) {
		api.WriteBadHeader(w, http.StatusUnauthorized, "Player does not have access.")
		return
	}

	err = database.DeleteCard(dbcs, id)
	if err != nil {
		api.WriteBadHeader(w, http.StatusBadRequest, "Failed to update the database.")
		return
	}

	w.Header().Add("HX-Refresh", "true")
	api.WriteGoodHeader(w, http.StatusCreated, "Success")
}
