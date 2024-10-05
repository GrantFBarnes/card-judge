package database

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type CardType struct {
	Id            uuid.UUID
	CreatedOnDate time.Time
	ChangedOnDate time.Time

	Name string
}

type Card struct {
	Id            uuid.UUID
	CreatedOnDate time.Time
	ChangedOnDate time.Time

	DeckId     uuid.UUID
	CardTypeId uuid.UUID
	Text       string
}

type CardDetails struct {
	Card
	CardTypeName string
}

func GetCardsInDeck(deckId uuid.UUID, cardTypeNameSearch string, textSearch string) ([]CardDetails, error) {
	if cardTypeNameSearch == "" {
		cardTypeNameSearch = "%"
	}

	if textSearch == "" {
		textSearch = "%"
	}

	sqlString := `
		SELECT
			C.ID,
			C.CREATED_ON_DATE,
			C.CHANGED_ON_DATE,
			C.DECK_ID,
			C.CARD_TYPE_ID,
			C.TEXT,
			CT.NAME AS CARD_TYPE_NAME
		FROM CARD AS C
			INNER JOIN CARD_TYPE AS CT ON CT.ID = C.CARD_TYPE_ID
		WHERE C.DECK_ID = ?
			AND CT.NAME LIKE ?
			AND C.TEXT LIKE ?
		ORDER BY
			CT.NAME ASC,
			TO_DAYS(C.CHANGED_ON_DATE) DESC,
			C.TEXT ASC
	`
	rows, err := query(sqlString, deckId, cardTypeNameSearch, textSearch)
	if err != nil {
		return nil, err
	}

	result := make([]CardDetails, 0)
	for rows.Next() {
		var cardDetails CardDetails
		if err := rows.Scan(
			&cardDetails.Id,
			&cardDetails.CreatedOnDate,
			&cardDetails.ChangedOnDate,
			&cardDetails.DeckId,
			&cardDetails.CardTypeId,
			&cardDetails.Text,
			&cardDetails.CardTypeName); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, cardDetails)
	}
	return result, nil
}

func GetCard(id uuid.UUID) (Card, error) {
	var card Card

	sqlString := `
		SELECT
			ID,
			CREATED_ON_DATE,
			CHANGED_ON_DATE,
			DECK_ID,
			CARD_TYPE_ID,
			TEXT
		FROM CARD
		WHERE ID = ?
	`
	rows, err := query(sqlString, id)
	if err != nil {
		return card, err
	}

	for rows.Next() {
		if err := rows.Scan(
			&card.Id,
			&card.CreatedOnDate,
			&card.ChangedOnDate,
			&card.DeckId,
			&card.CardTypeId,
			&card.Text); err != nil {
			log.Println(err)
			return card, errors.New("failed to scan row in query results")
		}
	}

	return card, nil
}

func CreateCard(deckId uuid.UUID, cardTypeName string, text string, blankCount int) (uuid.UUID, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		return id, errors.New("failed to generate new id")
	}

	cardTypeId, err := getCardTypeId(cardTypeName)
	if err != nil {
		log.Println(err)
		return id, errors.New("failed to get card type id")
	}

	if cardTypeId == uuid.Nil {
		log.Println(err)
		return id, errors.New("card type name not found")
	}

	sqlString := `
		INSERT INTO CARD (ID, DECK_ID, CARD_TYPE_ID, TEXT, BLANK_COUNT)
		VALUES (?, ?, ?, ?, ?)
	`
	return id, execute(sqlString, id, deckId, cardTypeId, text, blankCount)
}

func GetCardId(deckId uuid.UUID, text string) (uuid.UUID, error) {
	var id uuid.UUID

	sqlString := `
		SELECT
			ID
		FROM CARD
		WHERE DECK_ID = ?
			AND TEXT = ?
	`
	rows, err := query(sqlString, deckId, text)
	if err != nil {
		return id, err
	}

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			log.Println(err)
			return id, errors.New("failed to scan row in query results")
		}
	}

	return id, nil
}

func GetCardType(id uuid.UUID) (string, error) {
	var cardType string

	sqlString := `
		SELECT
			CT.NAME
		FROM CARD AS C
			INNER JOIN CARD_TYPE AS CT ON CT.ID = C.CARD_TYPE_ID
		WHERE C.ID = ?
	`
	rows, err := query(sqlString, id)
	if err != nil {
		return cardType, err
	}

	for rows.Next() {
		if err := rows.Scan(&cardType); err != nil {
			log.Println(err)
			return cardType, errors.New("failed to scan row in query results")
		}

	}

	return cardType, nil
}

func getCardTypeId(cardTypeName string) (uuid.UUID, error) {
	var id uuid.UUID

	sqlString := `
		SELECT
			ID
		FROM CARD_TYPE
		WHERE NAME = ?
	`
	rows, err := query(sqlString, cardTypeName)
	if err != nil {
		return id, err
	}

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			log.Println(err)
			return id, errors.New("failed to scan row in query results")
		}
	}

	return id, nil
}

func SetCardType(id uuid.UUID, cardTypeName string) error {
	cardTypeId, err := getCardTypeId(cardTypeName)
	if err != nil {
		log.Println(err)
		return errors.New("failed to get card type id")
	}

	if cardTypeId == uuid.Nil {
		log.Println(err)
		return errors.New("card type name not found")
	}

	sqlString := `
		UPDATE CARD
		SET
			CARD_TYPE_ID = ?
		WHERE ID = ?
	`
	return execute(sqlString, cardTypeId, id)
}

func SetCardText(id uuid.UUID, text string, blankCount int) error {
	sqlString := `
		UPDATE CARD
		SET
			TEXT = ?,
			BLANK_COUNT = ?
		WHERE ID = ?
	`
	return execute(sqlString, text, blankCount, id)
}

func DeleteCard(id uuid.UUID) error {
	sqlString := `
		DELETE FROM CARD
		WHERE ID = ?
	`
	return execute(sqlString, id)
}
