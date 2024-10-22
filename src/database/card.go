package database

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	Id            uuid.UUID
	CreatedOnDate time.Time
	ChangedOnDate time.Time

	DeckId   uuid.UUID
	Category string
	Text     string
}

func SearchCardsInDeck(deckId uuid.UUID, categorySearch string, textSearch string, pageNumber int) ([]Card, error) {
	if categorySearch == "" {
		categorySearch = "%"
	}

	if textSearch == "" {
		textSearch = "%"
	}

	pageSize := 10

	if pageNumber < 1 {
		pageNumber = 1
	} else if pageNumber > 100 {
		pageNumber = 100
	}

	sqlString := `
		SELECT
			C.ID,
			C.CREATED_ON_DATE,
			C.CHANGED_ON_DATE,
			C.DECK_ID,
			C.CATEGORY,
			C.TEXT
		FROM CARD AS C
		WHERE C.DECK_ID = ?
			AND C.CATEGORY LIKE ?
			AND C.TEXT LIKE ?
		ORDER BY
			C.CATEGORY ASC,
			TO_DAYS(C.CHANGED_ON_DATE) DESC,
			C.TEXT ASC
		LIMIT ? OFFSET ?
	`
	rows, err := query(sqlString, deckId, categorySearch, textSearch, pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	result := make([]Card, 0)
	for rows.Next() {
		var card Card
		if err := rows.Scan(
			&card.Id,
			&card.CreatedOnDate,
			&card.ChangedOnDate,
			&card.DeckId,
			&card.Category,
			&card.Text); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, card)
	}
	return result, nil
}

func GetCardsInDeckExport(deckId uuid.UUID) ([]Card, error) {
	sqlString := `
		SELECT
			C.CATEGORY,
			C.TEXT
		FROM CARD AS C
		WHERE C.DECK_ID = ?
		ORDER BY
			C.CATEGORY ASC,
			C.TEXT ASC
	`
	rows, err := query(sqlString, deckId)
	if err != nil {
		return nil, err
	}

	result := make([]Card, 0)
	for rows.Next() {
		var card Card
		if err := rows.Scan(&card.Category, &card.Text); err != nil {
			log.Println(err)
			return result, errors.New("failed to scan row in query results")
		}
		result = append(result, card)
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
			CATEGORY,
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
			&card.Category,
			&card.Text); err != nil {
			log.Println(err)
			return card, errors.New("failed to scan row in query results")
		}
	}

	return card, nil
}

func CreateCard(deckId uuid.UUID, category string, text string) (uuid.UUID, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
		return id, errors.New("failed to generate new id")
	}

	sqlString := `
		INSERT INTO CARD (ID, DECK_ID, CATEGORY, TEXT)
		VALUES (?, ?, ?, ?)
	`
	return id, execute(sqlString, id, deckId, category, text)
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

func GetCardTextStart(id uuid.UUID) (string, error) {
	var text string

	sqlString := `
		SELECT
			TEXT
		FROM CARD
		WHERE ID = ?
	`
	rows, err := query(sqlString, id)
	if err != nil {
		return text, err
	}

	for rows.Next() {
		if err := rows.Scan(&text); err != nil {
			log.Println(err)
			return text, errors.New("failed to scan row in query results")
		}
	}

	if len(text) > 100 {
		text = text[:100] + "..."
	}

	return text, nil
}

func SetCardCategory(id uuid.UUID, category string) error {
	sqlString := `
		UPDATE CARD
		SET CATEGORY = ?
		WHERE ID = ?
	`
	return execute(sqlString, category, id)
}

func SetCardText(id uuid.UUID, text string) error {
	sqlString := `
		UPDATE CARD
		SET
			TEXT = ?
		WHERE ID = ?
	`
	return execute(sqlString, text, id)
}

func DeleteCard(id uuid.UUID) error {
	sqlString := `
		DELETE FROM CARD
		WHERE ID = ?
	`
	return execute(sqlString, id)
}
