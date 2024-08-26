package apiCard

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/grantfbarnes/card-judge/database"
)

func Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse form"))
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
				w.Write([]byte("failed to parse deck id"))
				return
			}
		} else if key == "cardType" {
			if val[0] == "Judge" {
				cardType = database.Judge
			} else {
				cardType = database.Player
			}
		} else if key == "text" {
			text = val[0]
		}
	}

	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no text found"))
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	_, err = database.CreateCard(dbcs, deckId, cardType, text)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to create card"))
		return
	}

	w.Header().Add("HX-Refresh", "true")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to get card id"))
		return
	}

	dbcs := database.GetDatabaseConnectionString()
	err = database.DeleteCard(dbcs, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to delete card"))
		return
	}

	w.Header().Add("HX-Refresh", "true")
}
