package apiCard

import (
	"net/http"
	"regexp"
	"strings"
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
	var cardTypeNameSearch string
	var textSearch string
	for key, val := range r.Form {
		if key == "deckId" {
			deckId, err = uuid.Parse(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse deck id."))
				return
			}
		} else if key == "cardTypeName" {
			cardTypeNameSearch = val[0]
		} else if key == "text" {
			textSearch = val[0]
		}
	}

	textSearch = "%" + textSearch + "%"

	cards, err := database.GetCardsInDeck(deckId, cardTypeNameSearch, textSearch)
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
	var cardTypeName string
	var text string
	for key, val := range r.Form {
		if key == "deckId" {
			deckId, err = uuid.Parse(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse deck id."))
				return
			}
		} else if key == "cardTypeName" {
			cardTypeName = val[0]
		} else if key == "text" {
			text = val[0]
		}
	}

	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No text found."))
		return
	}

	userId := api.GetUserId(r)
	if userId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get user id."))
		return
	}

	existingCardId, err := database.GetCardId(deckId, text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	existingCardTypeName, err := database.GetCardType(existingCardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if existingCardId != uuid.Nil && existingCardTypeName == cardTypeName {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Card text already exists."))
		return
	}

	if !database.UserHasDeckAccess(userId, deckId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User does not have access."))
		return
	}

	var blankCount int
	if strings.ToLower(cardTypeName) == "judge" {
		text, blankCount, err = processJudgeCardText(text)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	_, err = database.CreateCard(deckId, cardTypeName, text, blankCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusCreated)
}

func SetCardType(w http.ResponseWriter, r *http.Request) {
	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form."))
		return
	}

	var deckId uuid.UUID
	var cardTypeName string
	for key, val := range r.Form {
		if key == "deckId" {
			deckId, err = uuid.Parse(val[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Failed to parse deck id."))
				return
			}
		} else if key == "cardTypeName" {
			cardTypeName = val[0]
		}
	}

	userId := api.GetUserId(r)
	if userId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get user id."))
		return
	}

	if !database.UserHasDeckAccess(userId, deckId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User does not have access."))
		return
	}

	err = database.SetCardType(cardId, cardTypeName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func SetText(w http.ResponseWriter, r *http.Request) {
	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
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

	userId := api.GetUserId(r)
	if userId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get user id."))
		return
	}

	if !database.UserHasDeckAccess(userId, deckId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User does not have access."))
		return
	}

	existingCardId, err := database.GetCardId(deckId, text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	existingCardTypeName, err := database.GetCardType(existingCardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	cardTypeName, err := database.GetCardType(cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if existingCardId != uuid.Nil && existingCardTypeName == cardTypeName {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Card text already exists."))
		return
	}

	var blankCount int
	if strings.ToLower(cardTypeName) == "judge" {
		text, blankCount, err = processJudgeCardText(text)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	err = database.SetCardText(cardId, text, blankCount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	cardIdString := r.PathValue("cardId")
	cardId, err := uuid.Parse(cardIdString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card id from path."))
		return
	}

	userId := api.GetUserId(r)
	if userId == uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get user id."))
		return
	}

	card, err := database.GetCard(cardId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get card."))
		return
	}

	if !database.UserHasDeckAccess(userId, card.DeckId) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User does not have access."))
		return
	}

	err = database.DeleteCard(cardId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func processJudgeCardText(text string) (string, int, error) {
	var normalizedText string = text
	var blankCount int = 0

	blankRegExp, err := regexp.Compile(`__+`)
	if err != nil {
		return normalizedText, blankCount, err
	}

	normalizedText = blankRegExp.ReplaceAllString(text, "_____")
	blankCount = len(blankRegExp.FindAllString(text, -1))

	return normalizedText, blankCount, err
}
